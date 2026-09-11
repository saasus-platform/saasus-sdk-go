package main

import (
	"os"
	"testing"
)

// TestMapToGeneratedTypesIntegration は実際のgenerated typesファイルを使用した統合テスト
func TestMapToGeneratedTypesIntegration(t *testing.T) {
	// デバッグモードを有効化
	os.Setenv("DEBUG", "1")
	defer os.Unsetenv("DEBUG")

	tests := []struct {
		name       string
		moduleName string
		skipReason string
	}{
		{
			name:       "communication API",
			moduleName: "communication",
		},
		{
			name:       "billing API",
			moduleName: "billing",
		},
		{
			name:       "auth API",
			moduleName: "auth",
		},
		{
			name:       "pricing API",
			moduleName: "pricing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skipReason != "" {
				t.Skip(tt.skipReason)
			}

			// 型定義ファイルの存在確認
			typesFilePath, err := findGeneratedTypesFile(tt.moduleName)
			if err != nil {
				t.Skipf("型定義ファイルが見つかりません: %v", err)
			}

			t.Logf("型定義ファイル: %s", typesFilePath)

			// テスト用のパラメータグループを作成
			paramGroups := createTestParamGroups(tt.moduleName)

			// マッピングを実行
			err = mapToGeneratedTypes(tt.moduleName, paramGroups, []string{tt.moduleName})
			if err != nil {
				t.Errorf("mapToGeneratedTypes() エラー: %v", err)
			}

			// マッピング結果を検証
			verifyMappingResults(t, paramGroups)
		})
	}
}

// createTestParamGroups はテスト用のパラメータグループを作成します
func createTestParamGroups(moduleName string) map[string]*ParamGroup {
	groups := make(map[string]*ParamGroup)

	switch moduleName {
	case "communication":
		groups["feedbacks"] = &ParamGroup{
			Name: "feedbacks",
			CreateParams: map[string]any{
				"user_id":              "test-user-id",
				"feedback_title":       "test title",
				"feedback_description": "test description",
			},
			UpdateParams: map[string]any{
				"feedback_title":       "updated title",
				"feedback_description": "updated description",
			},
		}
		groups["comments"] = &ParamGroup{
			Name: "comments",
			CreateParams: map[string]any{
				"body": "test comment",
			},
			UpdateParams: map[string]any{
				"body": "updated comment",
			},
		}

	case "billing":
		groups["stripe"] = &ParamGroup{
			Name: "stripe",
			CreateParams: map[string]any{
				"secret_key": "sk_test_xxx",
			},
			UpdateParams: map[string]any{
				"secret_key": "sk_test_yyy",
			},
		}

	case "auth":
		groups["users"] = &ParamGroup{
			Name: "users",
			CreateParams: map[string]any{
				"email":    "test@example.com",
				"password": "password123",
			},
			UpdateParams: map[string]any{
				"email": "updated@example.com",
			},
		}

	case "pricing":
		groups["plans"] = &ParamGroup{
			Name: "plans",
			CreateParams: map[string]any{
				"name":         "test plan",
				"display_name": "Test Plan",
				"description":  "test description",
			},
			UpdateParams: map[string]any{
				"name":         "updated plan",
				"display_name": "Updated Plan",
				"description":  "updated description",
			},
		}
	}

	return groups
}

// verifyMappingResults はマッピング結果を検証します
func verifyMappingResults(t *testing.T, paramGroups map[string]*ParamGroup) {
	for groupName, group := range paramGroups {
		t.Logf("グループ '%s' の検証", groupName)

		// CreateParamsの検証
		if len(group.CreateParams) > 0 {
			t.Logf("  CreateParams: %d個のフィールド", len(group.CreateParams))
			for fieldName, value := range group.CreateParams {
				if value == nil {
					t.Errorf("    フィールド '%s' の値がnilです", fieldName)
				} else {
					t.Logf("    フィールド '%s': %v (%T)", fieldName, value, value)
				}
			}
		}

		// UpdateParamsの検証
		if len(group.UpdateParams) > 0 {
			t.Logf("  UpdateParams: %d個のフィールド", len(group.UpdateParams))
			for fieldName, value := range group.UpdateParams {
				if value == nil {
					t.Errorf("    フィールド '%s' の値がnilです", fieldName)
				} else {
					t.Logf("    フィールド '%s': %v (%T)", fieldName, value, value)
				}
			}
		}
	}
}

// TestParseGeneratedTypesWithRealFiles は実際の型定義ファイルを解析するテスト
func TestParseGeneratedTypesWithRealFiles(t *testing.T) {
	modules := []string{"communication", "billing", "auth", "pricing"}

	for _, moduleName := range modules {
		t.Run(moduleName, func(t *testing.T) {
			// 型定義ファイルを探す
			typesFilePath, err := findGeneratedTypesFile(moduleName)
			if err != nil {
				t.Skipf("型定義ファイルが見つかりません: %v", err)
			}

			// 型定義を解析
			typeInfoMap, err := parseGeneratedTypes(typesFilePath)
			if err != nil {
				t.Fatalf("型定義の解析に失敗: %v", err)
			}

			if len(typeInfoMap) == 0 {
				t.Error("型定義が1つも見つかりませんでした")
			}

			t.Logf("発見された型: %d個", len(typeInfoMap))

			// 各型の詳細をログ出力
			for typeName, typeInfo := range typeInfoMap {
				t.Logf("  型名: %s (フィールド数: %d)", typeName, len(typeInfo.Fields))

				// フィールドの詳細（最初の3つのみ）
				count := 0
				for fieldName, fieldInfo := range typeInfo.Fields {
					if count >= 3 {
						break
					}
					t.Logf("    - %s: %s (json:%s)", fieldName, fieldInfo.Type, fieldInfo.JSONTag)
					count++
				}
				if len(typeInfo.Fields) > 3 {
					t.Logf("    ... 他 %d個のフィールド", len(typeInfo.Fields)-3)
				}
			}
		})
	}
}

// TestTypeNamePatternGeneration は型名パターン生成のテスト
func TestTypeNamePatternGeneration(t *testing.T) {
	tests := []struct {
		name         string
		groupName    string
		operation    string
		wantContains []string
		wantNotEmpty bool
	}{
		{
			name:      "feedbacks + Create",
			groupName: "feedbacks",
			operation: "Create",
			wantContains: []string{
				"CreateFeedbackJSONRequestBody",
				"CreateFeedbackJSONBody",
				"CreateFeedbackParams",
			},
			wantNotEmpty: true,
		},
		{
			name:      "stripe + Update",
			groupName: "stripe",
			operation: "Update",
			wantContains: []string{
				"UpdateStripeJSONRequestBody",
				"UpdateStripeParams",
			},
			wantNotEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			patterns := generateTypeNamePatterns(tt.groupName, tt.operation)

			if tt.wantNotEmpty && len(patterns) == 0 {
				t.Error("パターンが生成されませんでした")
			}

			t.Logf("生成されたパターン数: %d", len(patterns))
			for i, pattern := range patterns {
				t.Logf("  %d: %s", i+1, pattern)
			}

			// 期待されるパターンが含まれているか確認
			for _, want := range tt.wantContains {
				found := false
				for _, pattern := range patterns {
					if pattern == want {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("期待されるパターン '%s' が見つかりません", want)
				}
			}
		})
	}
}

// TestFieldNameNormalization はフィールド名正規化のテスト
func TestFieldNameNormalization(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantSnake  string
		wantCamel  string
		wantPascal string
	}{
		{
			name:       "snake_case",
			input:      "user_id",
			wantSnake:  "user_id",
			wantCamel:  "userId",
			wantPascal: "UserId",
		},
		{
			name:       "camelCase",
			input:      "userId",
			wantSnake:  "user_id",
			wantCamel:  "userId",
			wantPascal: "UserId",
		},
		{
			name:       "PascalCase",
			input:      "UserId",
			wantSnake:  "user_id",
			wantCamel:  "userId",
			wantPascal: "UserId",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSnake := normalizeFieldName(tt.input)
			if gotSnake != tt.wantSnake {
				t.Errorf("normalizeFieldName(%q) = %q, want %q", tt.input, gotSnake, tt.wantSnake)
			}

			gotCamel := snakeToCamel(tt.wantSnake)
			if gotCamel != tt.wantCamel {
				t.Errorf("snakeToCamel(%q) = %q, want %q", tt.wantSnake, gotCamel, tt.wantCamel)
			}

			gotPascal := snakeToPascal(tt.wantSnake)
			if gotPascal != tt.wantPascal {
				t.Errorf("snakeToPascal(%q) = %q, want %q", tt.wantSnake, gotPascal, tt.wantPascal)
			}
		})
	}
}
