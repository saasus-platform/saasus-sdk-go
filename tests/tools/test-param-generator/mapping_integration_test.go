package main

import (
	"fmt"
	"os"
	"testing"
)

// TestMappingSuccessRateWithRealTypes は実際のtypes.gen.goファイルを使用したマッピング成功率テスト
func TestMappingSuccessRateWithRealTypes(t *testing.T) {
	// デバッグモードを有効化
	os.Setenv("DEBUG", "1")
	defer os.Unsetenv("DEBUG")

	tests := []struct {
		name             string
		moduleName       string
		minSuccessRate   float64 // 最低成功率（80%）
		skipReason       string
		createTestParams func() map[string]*ParamGroup
	}{
		{
			name:           "auth API",
			moduleName:     "auth",
			minSuccessRate: 0.80,
			createTestParams: func() map[string]*ParamGroup {
				return map[string]*ParamGroup{
					"saasusers": {
						Name: "saasusers",
						CreateParams: map[string]any{
							"email":    "test@example.com",
							"password": "password123",
						},
						UpdateParams: map[string]any{
							"email": "updated@example.com",
						},
					},
					"tenants": {
						Name: "tenants",
						CreateParams: map[string]any{
							"name":                    "Test Tenant",
							"attributes":              map[string]any{},
							"back_office_staff_email": "admin@example.com",
						},
						UpdateParams: map[string]any{
							"name":                    "Updated Tenant",
							"attributes":              map[string]any{},
							"back_office_staff_email": "admin2@example.com",
						},
					},
				}
			},
		},
		{
			name:           "billing API",
			moduleName:     "billing",
			minSuccessRate: 0.80,
			createTestParams: func() map[string]*ParamGroup {
				return map[string]*ParamGroup{
					"stripe": {
						Name: "stripe",
						CreateParams: map[string]any{
							"secret_key": "sk_test_xxx",
						},
						UpdateParams: map[string]any{
							"secret_key": "sk_test_yyy",
						},
					},
				}
			},
		},
		{
			name:           "pricing API",
			moduleName:     "pricing",
			minSuccessRate: 0.80,
			createTestParams: func() map[string]*ParamGroup {
				return map[string]*ParamGroup{
					"plans": {
						Name: "plans",
						CreateParams: map[string]any{
							"name":         "test plan",
							"display_name": "Test Plan",
							"description":  "test description",
							"menu_ids":     []string{"menu-1", "menu-2"},
						},
						UpdateParams: map[string]any{
							"name":         "updated plan",
							"display_name": "Updated Plan",
							"description":  "updated description",
							"menu_ids":     []string{"menu-1", "menu-2"},
						},
					},
					"menus": {
						Name: "menus",
						CreateParams: map[string]any{
							"name":         "test menu",
							"display_name": "Test Menu",
							"description":  "test description",
							"unit_ids":     []string{"unit-1", "unit-2"},
						},
						UpdateParams: map[string]any{
							"name":         "updated menu",
							"display_name": "Updated Menu",
							"description":  "updated description",
							"unit_ids":     []string{"unit-1", "unit-2"},
						},
					},
				}
			},
		},
		{
			name:           "communication API",
			moduleName:     "communication",
			minSuccessRate: 0.80,
			createTestParams: func() map[string]*ParamGroup {
				return map[string]*ParamGroup{
					"feedbacks": {
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
					},
					"comments": {
						Name: "comments",
						CreateParams: map[string]any{
							"body": "test comment",
						},
						UpdateParams: map[string]any{
							"body": "updated comment",
						},
					},
				}
			},
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
			paramGroups := tt.createTestParams()

			// マッピング前のフィールド数をカウント
			totalFieldsBefore := countTotalFields(paramGroups)
			t.Logf("マッピング前の総フィールド数: %d", totalFieldsBefore)

			// マッピングを実行
			err = mapToGeneratedTypes(tt.moduleName, paramGroups, []string{tt.moduleName})
			if err != nil {
				t.Errorf("mapToGeneratedTypes() エラー: %v", err)
			}

			// マッピング成功率を計算
			successRate := calculateMappingSuccessRate(paramGroups, totalFieldsBefore)
			t.Logf("マッピング成功率: %.1f%% (%d/%d フィールド)",
				successRate*100,
				int(float64(totalFieldsBefore)*successRate),
				totalFieldsBefore)

			// 最低成功率を検証
			if successRate < tt.minSuccessRate {
				t.Errorf("マッピング成功率が目標値を下回っています: %.1f%% < %.1f%%",
					successRate*100, tt.minSuccessRate*100)
			}

			// 詳細な結果を出力
			printDetailedMappingResults(t, paramGroups)
		})
	}
}

// countTotalFields はパラメータグループ内の総フィールド数をカウントします
func countTotalFields(paramGroups map[string]*ParamGroup) int {
	total := 0
	for _, group := range paramGroups {
		total += len(group.CreateParams)
		total += len(group.UpdateParams)
	}
	return total
}

// calculateMappingSuccessRate はマッピング成功率を計算します
func calculateMappingSuccessRate(paramGroups map[string]*ParamGroup, totalFields int) float64 {
	if totalFields == 0 {
		return 0.0
	}

	successCount := 0
	for _, group := range paramGroups {
		// CreateParamsの成功数をカウント
		for _, value := range group.CreateParams {
			if value != nil {
				successCount++
			}
		}

		// UpdateParamsの成功数をカウント
		for _, value := range group.UpdateParams {
			if value != nil {
				successCount++
			}
		}
	}

	return float64(successCount) / float64(totalFields)
}

// printDetailedMappingResults は詳細なマッピング結果を出力します
func printDetailedMappingResults(t *testing.T, paramGroups map[string]*ParamGroup) {
	t.Log("=== 詳細なマッピング結果 ===")

	for groupName, group := range paramGroups {
		t.Logf("\nグループ: %s", groupName)

		// CreateParamsの結果
		if len(group.CreateParams) > 0 {
			t.Log("  CreateParams:")
			successCount := 0
			for fieldName, value := range group.CreateParams {
				if value != nil {
					successCount++
					t.Logf("    ✓ %s: %v (%T)", fieldName, value, value)
				} else {
					t.Logf("    ✗ %s: マッピング失敗", fieldName)
				}
			}
			successRate := float64(successCount) / float64(len(group.CreateParams)) * 100
			t.Logf("  CreateParams成功率: %.1f%% (%d/%d)",
				successRate, successCount, len(group.CreateParams))
		}

		// UpdateParamsの結果
		if len(group.UpdateParams) > 0 {
			t.Log("  UpdateParams:")
			successCount := 0
			for fieldName, value := range group.UpdateParams {
				if value != nil {
					successCount++
					t.Logf("    ✓ %s: %v (%T)", fieldName, value, value)
				} else {
					t.Logf("    ✗ %s: マッピング失敗", fieldName)
				}
			}
			successRate := float64(successCount) / float64(len(group.UpdateParams)) * 100
			t.Logf("  UpdateParams成功率: %.1f%% (%d/%d)",
				successRate, successCount, len(group.UpdateParams))
		}
	}
}

// TestEdgeCases はエッジケースのテスト
func TestEdgeCases(t *testing.T) {
	// デバッグモードを有効化
	os.Setenv("DEBUG", "1")
	defer os.Unsetenv("DEBUG")

	t.Run("空のrequest_body", func(t *testing.T) {
		moduleName := "communication"

		// 型定義ファイルの存在確認
		_, err := findGeneratedTypesFile(moduleName)
		if err != nil {
			t.Skipf("型定義ファイルが見つかりません: %v", err)
		}

		// 空のパラメータグループを作成
		paramGroups := map[string]*ParamGroup{
			"empty": {
				Name:         "empty",
				CreateParams: map[string]any{},
				UpdateParams: map[string]any{},
			},
		}

		// マッピングを実行（エラーが発生しないことを確認）
		err = mapToGeneratedTypes(moduleName, paramGroups, []string{moduleName})
		if err != nil {
			t.Errorf("空のrequest_bodyでエラーが発生: %v", err)
		}

		t.Log("空のrequest_bodyを正常に処理しました")
	})

	t.Run("ネストした構造体", func(t *testing.T) {
		moduleName := "auth"

		// 型定義ファイルの存在確認
		_, err := findGeneratedTypesFile(moduleName)
		if err != nil {
			t.Skipf("型定義ファイルが見つかりません: %v", err)
		}

		// ネストした構造体を含むパラメータグループを作成
		paramGroups := map[string]*ParamGroup{
			"tenants": {
				Name: "tenants",
				CreateParams: map[string]any{
					"name": "Test Tenant",
					"attributes": map[string]any{
						"key1": "value1",
						"key2": "value2",
					},
					"back_office_staff_email": "admin@example.com",
				},
			},
		}

		// マッピングを実行
		err = mapToGeneratedTypes(moduleName, paramGroups, []string{moduleName})
		if err != nil {
			t.Errorf("ネストした構造体でエラーが発生: %v", err)
		}

		// attributesフィールドが正しくマッピングされたか確認
		if attrs, ok := paramGroups["tenants"].CreateParams["attributes"]; ok {
			if attrs == nil {
				t.Error("attributesフィールドがnilです")
			} else {
				t.Logf("attributesフィールドが正しくマッピングされました: %v", attrs)
			}
		}
	})

	t.Run("配列型フィールド", func(t *testing.T) {
		moduleName := "pricing"

		// 型定義ファイルの存在確認
		_, err := findGeneratedTypesFile(moduleName)
		if err != nil {
			t.Skipf("型定義ファイルが見つかりません: %v", err)
		}

		// 配列型フィールドを含むパラメータグループを作成
		paramGroups := map[string]*ParamGroup{
			"plans": {
				Name: "plans",
				CreateParams: map[string]any{
					"name":         "test plan",
					"display_name": "Test Plan",
					"description":  "test description",
					"menu_ids":     []string{"menu-1", "menu-2", "menu-3"},
				},
			},
		}

		// マッピングを実行
		err = mapToGeneratedTypes(moduleName, paramGroups, []string{moduleName})
		if err != nil {
			t.Errorf("配列型フィールドでエラーが発生: %v", err)
		}

		// menu_idsフィールドが正しくマッピングされたか確認
		if menuIds, ok := paramGroups["plans"].CreateParams["menu_ids"]; ok {
			if menuIds == nil {
				t.Error("menu_idsフィールドがnilです")
			} else {
				t.Logf("menu_idsフィールドが正しくマッピングされました: %v", menuIds)
			}
		}
	})

	t.Run("オプショナルフィールド", func(t *testing.T) {
		moduleName := "auth"

		// 型定義ファイルの存在確認
		_, err := findGeneratedTypesFile(moduleName)
		if err != nil {
			t.Skipf("型定義ファイルが見つかりません: %v", err)
		}

		// オプショナルフィールドを含むパラメータグループを作成
		// passwordはオプショナル（omitempty）の可能性がある
		paramGroups := map[string]*ParamGroup{
			"saasusers": {
				Name: "saasusers",
				CreateParams: map[string]any{
					"email": "test@example.com",
					// passwordは省略
				},
			},
		}

		// マッピングを実行
		err = mapToGeneratedTypes(moduleName, paramGroups, []string{moduleName})
		if err != nil {
			t.Errorf("オプショナルフィールドでエラーが発生: %v", err)
		}

		// emailフィールドが正しくマッピングされたか確認
		if email, ok := paramGroups["saasusers"].CreateParams["email"]; ok {
			if email == nil {
				t.Error("emailフィールドがnilです")
			} else {
				t.Logf("emailフィールドが正しくマッピングされました: %v", email)
			}
		}

		t.Log("オプショナルフィールドを正常に処理しました")
	})
}

// TestMappingWithDifferentFieldNameFormats は異なるフィールド名形式のテスト
func TestMappingWithDifferentFieldNameFormats(t *testing.T) {
	// デバッグモードを有効化
	os.Setenv("DEBUG", "1")
	defer os.Unsetenv("DEBUG")

	moduleName := "communication"

	// 型定義ファイルの存在確認
	_, err := findGeneratedTypesFile(moduleName)
	if err != nil {
		t.Skipf("型定義ファイルが見つかりません: %v", err)
	}

	tests := []struct {
		name        string
		fieldName   string
		fieldValue  string
		expectMatch bool
	}{
		{
			name:        "snake_case",
			fieldName:   "user_id",
			fieldValue:  "test-user-id",
			expectMatch: true,
		},
		{
			name:        "camelCase",
			fieldName:   "userId",
			fieldValue:  "test-user-id",
			expectMatch: true,
		},
		{
			name:        "PascalCase",
			fieldName:   "UserId",
			fieldValue:  "test-user-id",
			expectMatch: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// パラメータグループを作成
			paramGroups := map[string]*ParamGroup{
				"feedbacks": {
					Name: "feedbacks",
					CreateParams: map[string]any{
						tt.fieldName: tt.fieldValue,
					},
				},
			}

			// マッピングを実行
			err = mapToGeneratedTypes(moduleName, paramGroups, []string{moduleName})
			if err != nil {
				t.Errorf("マッピングでエラーが発生: %v", err)
			}

			// マッピング結果を確認
			value := paramGroups["feedbacks"].CreateParams[tt.fieldName]
			matched := value != nil

			if tt.expectMatch && !matched {
				t.Errorf("フィールド '%s' がマッピングされませんでした", tt.fieldName)
			} else if !tt.expectMatch && matched {
				t.Errorf("フィールド '%s' がマッピングされるべきではありませんでした", tt.fieldName)
			} else if matched {
				t.Logf("フィールド '%s' が正しくマッピングされました: %v", tt.fieldName, value)
			}
		})
	}
}

// TestMappingPerformance はマッピングのパフォーマンステスト
func TestMappingPerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("パフォーマンステストをスキップします")
	}

	// デバッグモードを無効化（パフォーマンス測定のため）
	os.Unsetenv("DEBUG")

	moduleName := "communication"

	// 型定義ファイルの存在確認
	_, err := findGeneratedTypesFile(moduleName)
	if err != nil {
		t.Skipf("型定義ファイルが見つかりません: %v", err)
	}

	// 大量のパラメータグループを作成
	paramGroups := make(map[string]*ParamGroup)
	for i := 0; i < 100; i++ {
		groupName := fmt.Sprintf("group%d", i)
		paramGroups[groupName] = &ParamGroup{
			Name: groupName,
			CreateParams: map[string]any{
				"user_id":              fmt.Sprintf("user-%d", i),
				"feedback_title":       fmt.Sprintf("title-%d", i),
				"feedback_description": fmt.Sprintf("description-%d", i),
			},
			UpdateParams: map[string]any{
				"feedback_title":       fmt.Sprintf("updated-title-%d", i),
				"feedback_description": fmt.Sprintf("updated-description-%d", i),
			},
		}
	}

	// マッピングを実行
	err = mapToGeneratedTypes(moduleName, paramGroups, []string{moduleName})
	if err != nil {
		t.Errorf("マッピングでエラーが発生: %v", err)
	}

	t.Logf("100個のパラメータグループのマッピングが完了しました")
}

// TestMappingWithCustomRules はカスタムマッピングルールのテスト
func TestMappingWithCustomRules(t *testing.T) {
	// デバッグモードを有効化
	os.Setenv("DEBUG", "1")
	defer os.Unsetenv("DEBUG")

	moduleName := "auth"

	// 型定義ファイルの存在確認
	_, err := findGeneratedTypesFile(moduleName)
	if err != nil {
		t.Skipf("型定義ファイルが見つかりません: %v", err)
	}

	// カスタムマッピングルールを作成
	rules := &MappingRules{
		TypeNameMappings: map[string]string{
			"saasusers.Create": "CreateSaasUserJSONRequestBody",
			"saasusers.Update": "UpdateSaasUserEmailJSONRequestBody",
		},
		FieldNameMappings: map[string]map[string]string{
			"CreateSaasUserJSONRequestBody": {
				"email":    "Email",
				"password": "Password",
			},
		},
		ResourceAliases: map[string][]string{
			"users": {"saasusers", "saas_users"},
		},
		FieldAliases: map[string][]string{
			"email": {"Email", "email_address"},
		},
	}

	// カスタムルールを一時ファイルに保存
	tmpFile := t.TempDir() + "/mapping_rules.json"
	if err := rules.SaveMappingRules(tmpFile); err != nil {
		t.Fatalf("カスタムルールの保存に失敗: %v", err)
	}

	// カスタムルールを読み込み
	loadedRules, err := LoadMappingRules(tmpFile)
	if err != nil {
		t.Fatalf("カスタムルールの読み込みに失敗: %v", err)
	}

	// ルールが正しく読み込まれたか確認
	if len(loadedRules.TypeNameMappings) != len(rules.TypeNameMappings) {
		t.Errorf("TypeNameMappingsの数が一致しません: %d != %d",
			len(loadedRules.TypeNameMappings), len(rules.TypeNameMappings))
	}

	t.Log("カスタムマッピングルールが正しく読み込まれました")
}

// TestModuleInfoUsage はModule情報を使用したtypes.gen.goファイル参照のテスト
// 要件: 8.1, 8.2, 8.3
func TestModuleInfoUsage(t *testing.T) {
	// デバッグモードを有効化
	os.Setenv("DEBUG", "1")
	defer os.Unsetenv("DEBUG")

	tests := []struct {
		name       string
		moduleName string
		expectPath string
	}{
		{
			name:       "auth module",
			moduleName: "auth",
			expectPath: "generated/authapi/types.gen.go",
		},
		{
			name:       "billing module",
			moduleName: "billing",
			expectPath: "generated/billingapi/types.gen.go",
		},
		{
			name:       "pricing module",
			moduleName: "pricing",
			expectPath: "generated/pricingapi/types.gen.go",
		},
		{
			name:       "communication module",
			moduleName: "communication",
			expectPath: "generated/communicationapi/types.gen.go",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// getTypesFilePathを使用してパスを取得
			path, err := getTypesFilePath(tt.moduleName)
			if err != nil {
				t.Skipf("型定義ファイルが見つかりません: %v", err)
			}

			t.Logf("Module情報から取得したパス: %s", path)

			// パスが期待値を含むか確認
			if !contains(path, tt.expectPath) {
				t.Errorf("パスが期待値を含みません: %s (期待: %s を含む)", path, tt.expectPath)
			}

			// ファイルが実際に存在するか確認
			if _, err := os.Stat(path); os.IsNotExist(err) {
				t.Errorf("ファイルが存在しません: %s", path)
			}
		})
	}
}

// TestModuleSpecificMappingRules はモジュール固有マッピングルールのテスト
// 要件: 8.1, 8.2, 8.3
func TestModuleSpecificMappingRules(t *testing.T) {
	// デバッグモードを有効化
	os.Setenv("DEBUG", "1")
	defer os.Unsetenv("DEBUG")

	// モジュール固有ルールを含むカスタムマッピングルールを作成
	rules := NewMappingRules()

	// authモジュール固有のルール
	rules.AddModuleTypeNameMapping("auth", "saasusers", "Create", "CreateSaasUserJSONRequestBody")
	rules.AddModuleFieldNameMapping("auth", "CreateSaasUserJSONRequestBody", "email", "Email")

	// communicationモジュール固有のルール
	rules.AddModuleTypeNameMapping("communication", "feedbacks", "Create", "CreateFeedbackJSONRequestBody")
	rules.AddModuleFieldNameMapping("communication", "CreateFeedbackJSONRequestBody", "user_id", "UserId")

	// ルールを一時ファイルに保存
	tmpFile := t.TempDir() + "/mapping_rules.json"
	if err := rules.SaveMappingRules(tmpFile); err != nil {
		t.Fatalf("カスタムルールの保存に失敗: %v", err)
	}

	// ルールを読み込み
	loadedRules, err := LoadMappingRules(tmpFile)
	if err != nil {
		t.Fatalf("カスタムルールの読み込みに失敗: %v", err)
	}

	// モジュール固有ルールが正しく読み込まれたか確認
	if len(loadedRules.ModuleSpecificRules) != 2 {
		t.Errorf("ModuleSpecificRulesの数が一致しません: %d != 2", len(loadedRules.ModuleSpecificRules))
	}

	// authモジュールのルールを確認
	if authRules, ok := loadedRules.ModuleSpecificRules["auth"]; ok {
		if typeName, ok := authRules.TypeNameMappings["saasusers.Create"]; !ok || typeName != "CreateSaasUserJSONRequestBody" {
			t.Errorf("authモジュールの型名マッピングが正しくありません")
		}
	} else {
		t.Error("authモジュールのルールが見つかりません")
	}

	// communicationモジュールのルールを確認
	if commRules, ok := loadedRules.ModuleSpecificRules["communication"]; ok {
		if typeName, ok := commRules.TypeNameMappings["feedbacks.Create"]; !ok || typeName != "CreateFeedbackJSONRequestBody" {
			t.Errorf("communicationモジュールの型名マッピングが正しくありません")
		}
	} else {
		t.Error("communicationモジュールのルールが見つかりません")
	}

	t.Log("モジュール固有マッピングルールが正しく読み込まれました")
}

// TestModuleStatistics はモジュール別統計の集計テスト
// 要件: 8.1, 8.2, 8.3
func TestModuleStatistics(t *testing.T) {
	// デバッグモードを有効化
	os.Setenv("DEBUG", "1")
	defer os.Unsetenv("DEBUG")

	tests := []struct {
		name       string
		moduleName string
	}{
		{name: "auth module", moduleName: "auth"},
		{name: "billing module", moduleName: "billing"},
		{name: "pricing module", moduleName: "pricing"},
		{name: "communication module", moduleName: "communication"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 型定義ファイルの存在確認
			_, err := getTypesFilePath(tt.moduleName)
			if err != nil {
				t.Skipf("型定義ファイルが見つかりません: %v", err)
			}

			// DiagnosticReporterを作成
			reporter := NewDiagnosticReporter(tt.moduleName, true)

			// テスト用のメトリクスを設定
			reporter.metrics.TotalGroups = 10
			reporter.metrics.SuccessfulGroups = 8
			reporter.metrics.TotalFields = 50
			reporter.metrics.MappedFields = 40

			// モジュール別統計を取得
			stats := reporter.GetModuleStats()

			// 統計が正しく計算されているか確認
			if stats.ModuleName != tt.moduleName {
				t.Errorf("モジュール名が一致しません: %s != %s", stats.ModuleName, tt.moduleName)
			}

			if stats.TotalGroups != 10 {
				t.Errorf("総グループ数が一致しません: %d != 10", stats.TotalGroups)
			}

			if stats.SuccessRate != 80.0 {
				t.Errorf("成功率が一致しません: %.1f != 80.0", stats.SuccessRate)
			}

			t.Logf("モジュール '%s' の統計: 成功率 %.1f%% (%d/%d フィールド)",
				stats.ModuleName, stats.SuccessRate, stats.MappedFields, stats.TotalFields)
		})
	}
}

// contains はstringがsubstrを含むか確認します
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
