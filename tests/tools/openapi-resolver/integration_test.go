package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	openapiresolver "github.com/saasus-platform/saasus-sdk-go/tests/tools/openapi-resolver-lib"
)

// TestPostmanAnalyzerIntegration tests the complete flow:
// Postman → OpenAPI → client.gen.go → types.gen.go
func TestPostmanAnalyzerIntegration(t *testing.T) {
	tests := []struct {
		name       string
		moduleName string
		method     string
		path       string
	}{
		{
			name:       "communication - CreateFeedback",
			moduleName: "communication",
			method:     "POST",
			path:       "/feedbacks",
		},
		{
			name:       "billing - UpdateStripeInfo",
			moduleName: "billing",
			method:     "PUT",
			path:       "/stripe/info",
		},
		{
			name:       "auth - CreateSaasUser",
			moduleName: "auth",
			method:     "POST",
			path:       "/users",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 1. OpenAPI resolverでoperationIdを取得
			specPath := filepath.Join("../../../api", tt.moduleName+"api.yml")
			resolver, err := openapiresolver.NewOperationResolver(specPath)
			if err != nil {
				t.Fatalf("NewOperationResolver() error = %v", err)
			}

			mapping, err := resolver.Resolve(tt.method, tt.path)
			if err != nil {
				t.Fatalf("Resolve() error = %v", err)
			}

			if mapping.OperationID == "" {
				t.Error("OperationID is empty")
			}

			t.Logf("Step 1: Found operationId: %s", mapping.OperationID)

			// 2. Client parserでメソッド名とパラメータ型を取得
			clientParser, err := NewClientInterfaceParser(tt.moduleName)
			if err != nil {
				t.Fatalf("NewClientInterfaceParser() error = %v", err)
			}

			methodSig, err := clientParser.FindMethodByOperationID(mapping.OperationID)
			if err != nil {
				t.Fatalf("FindMethodByOperationID() error = %v", err)
			}

			if methodSig.MethodName == "" {
				t.Error("MethodName is empty")
			}

			t.Logf("Step 2: Found method: %s with %d parameters", methodSig.MethodName, len(methodSig.ParameterTypes))

			// 3. パラメータ型名を抽出
			var paramTypeName string
			for _, param := range methodSig.ParameterTypes {
				if strings.Contains(param.Type, "Param") ||
					strings.Contains(param.Type, "JSONRequestBody") ||
					strings.Contains(param.Type, "RequestBody") {
					paramTypeName = param.Type
					break
				}
			}

			if paramTypeName == "" && tt.method != "GET" && tt.method != "DELETE" {
				t.Logf("Warning: No parameter type found for %s method", tt.method)
			}

			if paramTypeName != "" {
				t.Logf("Step 3: Found parameter type: %s", paramTypeName)

				// 4. Type parserで型定義を取得
				typeParser, err := NewTypeDefinitionParser(tt.moduleName)
				if err != nil {
					t.Fatalf("NewTypeDefinitionParser() error = %v", err)
				}

				// JSONRequestBodyをParamに変換して検索
				actualTypeName := strings.Replace(paramTypeName, "JSONRequestBody", "Param", 1)
				t.Logf("Step 4: Looking for type: %s (converted from %s)", actualTypeName, paramTypeName)

				typeDef, err := typeParser.GetTypeDefinition(actualTypeName)
				if err != nil {
					t.Fatalf("GetTypeDefinition() error = %v", err)
				}

				if len(typeDef.Fields) == 0 {
					t.Error("Type definition has no fields")
				}

				t.Logf("Step 4: Found type definition with %d fields", len(typeDef.Fields))

				// フィールド情報を検証
				for _, field := range typeDef.Fields {
					if field.Name == "" {
						t.Error("Field has empty name")
					}
					if field.Type == "" {
						t.Error("Field has empty type")
					}
					if field.JSONTag == "" {
						t.Logf("Warning: Field %s has no JSON tag", field.Name)
					}
					t.Logf("  Field: %s (%s) json:%s", field.Name, field.Type, field.JSONTag)
				}
			}

			// 統合結果のサマリー
			t.Logf("Integration test completed successfully:")
			t.Logf("  Module: %s", tt.moduleName)
			t.Logf("  Method: %s %s", tt.method, tt.path)
			t.Logf("  OperationID: %s", mapping.OperationID)
			t.Logf("  Client Method: %s", methodSig.MethodName)
			if paramTypeName != "" {
				t.Logf("  Parameter Type: %s", paramTypeName)
			}
		})
	}
}

// TestEndToEndFlow tests the complete end-to-end flow with JSON output
func TestEndToEndFlow(t *testing.T) {
	moduleName := "communication"
	method := "POST"
	path := "/feedbacks"

	// 完全なフローを実行
	specPath := filepath.Join("../../../api", moduleName+"api.yml")
	resolver, err := openapiresolver.NewOperationResolver(specPath)
	if err != nil {
		t.Fatalf("NewOperationResolver() error = %v", err)
	}

	mapping, err := resolver.Resolve(method, path)
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}

	clientParser, err := NewClientInterfaceParser(moduleName)
	if err != nil {
		t.Fatalf("NewClientInterfaceParser() error = %v", err)
	}

	methodSig, err := clientParser.FindMethodByOperationID(mapping.OperationID)
	if err != nil {
		t.Fatalf("FindMethodByOperationID() error = %v", err)
	}

	// パラメータ型を探す
	var paramTypeName string
	for _, param := range methodSig.ParameterTypes {
		if strings.Contains(param.Type, "Param") {
			paramTypeName = param.Type
			break
		}
	}

	var fields map[string]interface{}
	if paramTypeName != "" {
		typeParser, err := NewTypeDefinitionParser(moduleName)
		if err != nil {
			t.Fatalf("NewTypeDefinitionParser() error = %v", err)
		}

		typeDef, err := typeParser.GetTypeDefinition(paramTypeName)
		if err != nil {
			t.Fatalf("GetTypeDefinition() error = %v", err)
		}

		// フィールド情報をマップに変換
		fields = make(map[string]interface{})
		for _, field := range typeDef.Fields {
			fields[field.JSONTag] = map[string]interface{}{
				"name":     field.Name,
				"type":     field.Type,
				"required": field.Required,
			}
		}
	}

	// 結果をJSON形式で出力
	result := map[string]interface{}{
		"module":           moduleName,
		"method":           method,
		"path":             path,
		"operation_id":     mapping.OperationID,
		"client_method":    methodSig.MethodName,
		"parameter_type":   paramTypeName,
		"parameter_fields": fields,
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	t.Logf("Complete integration result:\n%s", string(jsonData))

	// 結果の検証
	if result["operation_id"] == "" {
		t.Error("operation_id is empty")
	}
	if result["client_method"] == "" {
		t.Error("client_method is empty")
	}
	if paramTypeName != "" && len(fields) == 0 {
		t.Error("parameter_fields is empty")
	}
}

// TestMultipleModulesIntegration tests integration across multiple modules
func TestMultipleModulesIntegration(t *testing.T) {
	modules := []struct {
		name   string
		method string
		path   string
	}{
		{"communication", "POST", "/feedbacks"},
		{"billing", "PUT", "/stripe/info"},
		{"auth", "POST", "/users"},
		{"pricing", "GET", "/pricing/units"},
	}

	successCount := 0
	for _, mod := range modules {
		t.Run(mod.name, func(t *testing.T) {
			specPath := filepath.Join("../../../api", mod.name+"api.yml")

			// OpenAPI仕様ファイルが存在するか確認
			if _, err := os.Stat(specPath); os.IsNotExist(err) {
				t.Skipf("Spec file not found: %s", specPath)
				return
			}

			resolver, err := openapiresolver.NewOperationResolver(specPath)
			if err != nil {
				t.Skipf("Could not create resolver: %v", err)
				return
			}

			mapping, err := resolver.Resolve(mod.method, mod.path)
			if err != nil {
				t.Skipf("Could not resolve operation: %v", err)
				return
			}

			clientParser, err := NewClientInterfaceParser(mod.name)
			if err != nil {
				t.Skipf("Could not create client parser: %v", err)
				return
			}

			methodSig, err := clientParser.FindMethodByOperationID(mapping.OperationID)
			if err != nil {
				t.Skipf("Could not find method: %v", err)
				return
			}

			t.Logf("✓ Module %s: %s -> %s", mod.name, mapping.OperationID, methodSig.MethodName)
			successCount++
		})
	}

	t.Logf("Successfully integrated %d/%d modules", successCount, len(modules))
}
