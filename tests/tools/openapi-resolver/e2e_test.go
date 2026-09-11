package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	openapiresolver "github.com/saasus-platform/saasus-sdk-go/tests/tools/openapi-resolver-lib"
)

// TestEndToEndMappingAccuracy tests the complete flow and measures mapping success rate
func TestEndToEndMappingAccuracy(t *testing.T) {
	modules := []struct {
		name           string
		minSuccessRate float64 // 最低期待成功率
	}{
		{"communication", 0.80}, // 80%以上
		{"billing", 0.80},
		{"auth", 0.80},
		{"pricing", 0.70}, // 型定義不足により70%
	}

	overallStats := struct {
		totalTests    int
		successTests  int
		failedModules []string
	}{}

	for _, mod := range modules {
		t.Run(mod.name, func(t *testing.T) {
			// OpenAPI仕様ファイルのパスを確認
			specPath := filepath.Join("../../../api", mod.name+"api.yml")
			if _, err := os.Stat(specPath); os.IsNotExist(err) {
				t.Skipf("Spec file not found: %s", specPath)
				return
			}

			// 統計情報
			stats := struct {
				totalOperations     int
				resolvedOperations  int
				foundMethods        int
				foundTypes          int
				completeIntegration int
			}{}

			// OpenAPI resolverを初期化
			resolver, err := openapiresolver.NewOperationResolver(specPath)
			if err != nil {
				t.Fatalf("Failed to create resolver: %v", err)
			}

			// すべてのoperationを取得
			mappings, err := resolver.ResolveAll()
			if err != nil {
				t.Fatalf("Failed to resolve all operations: %v", err)
			}

			stats.totalOperations = len(mappings)

			// Client parserを初期化
			clientParser, err := NewClientInterfaceParser(mod.name)
			if err != nil {
				t.Fatalf("Failed to create client parser: %v", err)
			}

			// Type parserを初期化
			typeParser, err := NewTypeDefinitionParser(mod.name)
			if err != nil {
				t.Fatalf("Failed to create type parser: %v", err)
			}

			// 各operationをテスト
			for _, mapping := range mappings {
				stats.resolvedOperations++

				// Client methodを検索
				methodSig, err := clientParser.FindMethodByOperationID(mapping.OperationID)
				if err != nil {
					t.Logf("Warning: Could not find method for %s: %v", mapping.OperationID, err)
					continue
				}
				stats.foundMethods++

				// パラメータ型を検索
				var paramTypeName string
				for _, param := range methodSig.ParameterTypes {
					if strings.Contains(param.Type, "Param") ||
						strings.Contains(param.Type, "JSONRequestBody") {
						paramTypeName = param.Type
						break
					}
				}

				// GET/DELETEメソッドはbodyパラメータがないので成功とカウント
				if paramTypeName == "" && (mapping.HTTPMethod == "GET" || mapping.HTTPMethod == "DELETE") {
					stats.completeIntegration++
					continue
				}

				if paramTypeName != "" {
					// JSONRequestBodyをParamに変換
					actualTypeName := strings.Replace(paramTypeName, "JSONRequestBody", "Param", 1)

					// 型定義を検索
					typeDef, err := typeParser.GetTypeDefinition(actualTypeName)
					if err != nil {
						t.Logf("Warning: Could not find type %s: %v", actualTypeName, err)
						continue
					}

					if len(typeDef.Fields) > 0 {
						stats.foundTypes++
						stats.completeIntegration++
					}
				}
			}

			// 成功率を計算
			successRate := 0.0
			if stats.totalOperations > 0 {
				successRate = float64(stats.completeIntegration) / float64(stats.totalOperations)
			}

			// 統計情報をログ出力
			t.Logf("Module: %s", mod.name)
			t.Logf("  Total Operations: %d", stats.totalOperations)
			t.Logf("  Resolved Operations: %d", stats.resolvedOperations)
			t.Logf("  Found Methods: %d", stats.foundMethods)
			t.Logf("  Found Types: %d", stats.foundTypes)
			t.Logf("  Complete Integration: %d", stats.completeIntegration)
			t.Logf("  Success Rate: %.2f%%", successRate*100)

			// 全体統計を更新
			overallStats.totalTests += stats.totalOperations
			overallStats.successTests += stats.completeIntegration

			// 最低成功率をチェック
			if successRate < mod.minSuccessRate {
				overallStats.failedModules = append(overallStats.failedModules, mod.name)
				t.Errorf("Success rate %.2f%% is below minimum %.2f%%", successRate*100, mod.minSuccessRate*100)
			}
		})
	}

	// 全体の統計を出力
	t.Logf("\n=== Overall Statistics ===")
	t.Logf("Total Tests: %d", overallStats.totalTests)
	t.Logf("Successful Tests: %d", overallStats.successTests)
	if overallStats.totalTests > 0 {
		overallSuccessRate := float64(overallStats.successTests) / float64(overallStats.totalTests)
		t.Logf("Overall Success Rate: %.2f%%", overallSuccessRate*100)

		if overallSuccessRate < 0.80 {
			t.Errorf("Overall success rate %.2f%% is below minimum 80%%", overallSuccessRate*100)
		}
	}

	if len(overallStats.failedModules) > 0 {
		t.Logf("Failed Modules: %v", overallStats.failedModules)
	}
}

// TestTypeNameConversion tests the conversion between JSONRequestBody and Param suffixes
func TestTypeNameConversion(t *testing.T) {
	tests := []struct {
		name       string
		moduleName string
		operations []struct {
			operationID string
			expectType  string
		}
	}{
		{
			name:       "communication",
			moduleName: "communication",
			operations: []struct {
				operationID string
				expectType  string
			}{
				{"CreateFeedback", "CreateFeedbackParam"},
				{"CreateFeedbackComment", "CreateFeedbackCommentParam"},
			},
		},
		{
			name:       "billing",
			moduleName: "billing",
			operations: []struct {
				operationID string
				expectType  string
			}{
				{"UpdateStripeInfo", "UpdateStripeInfoParam"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientParser, err := NewClientInterfaceParser(tt.moduleName)
			if err != nil {
				t.Skipf("Could not create client parser: %v", err)
				return
			}

			typeParser, err := NewTypeDefinitionParser(tt.moduleName)
			if err != nil {
				t.Skipf("Could not create type parser: %v", err)
				return
			}

			for _, op := range tt.operations {
				t.Run(op.operationID, func(t *testing.T) {
					// Client methodを検索
					methodSig, err := clientParser.FindMethodByOperationID(op.operationID)
					if err != nil {
						t.Skipf("Could not find method: %v", err)
						return
					}

					// パラメータ型を検索
					var paramTypeName string
					for _, param := range methodSig.ParameterTypes {
						if strings.Contains(param.Type, "JSONRequestBody") {
							paramTypeName = param.Type
							break
						}
					}

					if paramTypeName == "" {
						t.Skip("No JSONRequestBody parameter found")
						return
					}

					// JSONRequestBodyをParamに変換
					actualTypeName := strings.Replace(paramTypeName, "JSONRequestBody", "Param", 1)

					// 期待される型名と一致するか確認
					if actualTypeName != op.expectType {
						t.Errorf("Converted type name = %s, want %s", actualTypeName, op.expectType)
					}

					// 型定義が存在するか確認
					typeDef, err := typeParser.GetTypeDefinition(actualTypeName)
					if err != nil {
						t.Errorf("Type definition not found: %v", err)
						return
					}

					if len(typeDef.Fields) == 0 {
						t.Error("Type definition has no fields")
					}

					t.Logf("✓ %s -> %s (%d fields)", paramTypeName, actualTypeName, len(typeDef.Fields))
				})
			}
		})
	}
}

// TestJSONOutputFormat tests that the integration result can be serialized to JSON
func TestJSONOutputFormat(t *testing.T) {
	moduleName := "communication"
	operationID := "CreateFeedback"

	// Client parserを初期化
	clientParser, err := NewClientInterfaceParser(moduleName)
	if err != nil {
		t.Fatalf("Failed to create client parser: %v", err)
	}

	// Type parserを初期化
	typeParser, err := NewTypeDefinitionParser(moduleName)
	if err != nil {
		t.Fatalf("Failed to create type parser: %v", err)
	}

	// Client methodを検索
	methodSig, err := clientParser.FindMethodByOperationID(operationID)
	if err != nil {
		t.Fatalf("Failed to find method: %v", err)
	}

	// パラメータ型を検索
	var paramTypeName string
	for _, param := range methodSig.ParameterTypes {
		if strings.Contains(param.Type, "JSONRequestBody") {
			paramTypeName = param.Type
			break
		}
	}

	if paramTypeName == "" {
		t.Fatal("No JSONRequestBody parameter found")
	}

	// 型定義を取得
	actualTypeName := strings.Replace(paramTypeName, "JSONRequestBody", "Param", 1)
	typeDef, err := typeParser.GetTypeDefinition(actualTypeName)
	if err != nil {
		t.Fatalf("Failed to get type definition: %v", err)
	}

	// フィールド情報をマップに変換
	fields := make(map[string]interface{})
	for _, field := range typeDef.Fields {
		fields[field.JSONTag] = map[string]interface{}{
			"name":     field.Name,
			"type":     field.Type,
			"required": field.Required,
		}
	}

	// 結果をJSON形式で出力
	result := map[string]interface{}{
		"module":           moduleName,
		"operation_id":     operationID,
		"client_method":    methodSig.MethodName,
		"parameter_type":   paramTypeName,
		"actual_type_name": actualTypeName,
		"parameter_fields": fields,
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	t.Logf("JSON output:\n%s", string(jsonData))

	// JSONが有効であることを確認
	var decoded map[string]interface{}
	if err := json.Unmarshal(jsonData, &decoded); err != nil {
		t.Errorf("Failed to unmarshal JSON: %v", err)
	}

	// 必須フィールドが存在することを確認
	requiredFields := []string{"module", "operation_id", "client_method", "parameter_type", "parameter_fields"}
	for _, field := range requiredFields {
		if _, ok := decoded[field]; !ok {
			t.Errorf("Required field %s not found in JSON output", field)
		}
	}
}
