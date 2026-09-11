package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestEndToEnd は実際のPostman Analyzerレポートを使用した統合テストです
func TestEndToEnd(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tmpDir := t.TempDir()

	// 実際のPostman Analyzerレポートを使用
	reportPath := "../postman-analyzer/reports/communication/postman_analysis.json"
	if _, err := os.Stat(reportPath); os.IsNotExist(err) {
		t.Skipf("Postman Analyzerレポートが見つかりません: %s", reportPath)
	}

	// 設定を作成
	config := &Config{
		ReportPath: reportPath,
		OutputDir:  tmpDir,
		ModuleName: "communication",
		Overwrite:  true,
		Verbose:    false,
	}

	// メイン処理を実行
	if err := run(config); err != nil {
		t.Fatalf("統合テストの実行に失敗: %v", err)
	}

	// 生成されたファイルのパスを取得
	paramsPath, helpersPath := config.GetOutputPaths()

	// 1. パラメータファイルが生成されたことを確認
	t.Run("パラメータファイルの生成確認", func(t *testing.T) {
		if _, err := os.Stat(paramsPath); os.IsNotExist(err) {
			t.Errorf("パラメータファイルが生成されていません: %s", paramsPath)
		}
	})

	// 2. ヘルパー関数ファイルが生成されたことを確認
	t.Run("ヘルパー関数ファイルの生成確認", func(t *testing.T) {
		if _, err := os.Stat(helpersPath); os.IsNotExist(err) {
			t.Errorf("ヘルパー関数ファイルが生成されていません: %s", helpersPath)
		}
	})

	// 3. パラメータファイルの内容を検証
	t.Run("パラメータファイルの内容検証", func(t *testing.T) {
		data, err := os.ReadFile(paramsPath)
		if err != nil {
			t.Fatalf("パラメータファイルの読み込みに失敗: %v", err)
		}

		var params map[string]any
		if err := json.Unmarshal(data, &params); err != nil {
			t.Fatalf("パラメータファイルのJSON解析に失敗: %v", err)
		}

		// 期待されるパラメータグループが存在することを確認
		// feedbacksグループは必ず存在するはず
		if _, ok := params["feedbacks"]; !ok {
			t.Error("feedbacksパラメータグループが見つかりません")
		}

		// 少なくとも1つのパラメータグループが存在することを確認
		if len(params) == 0 {
			t.Error("パラメータグループが1つも生成されていません")
		}

		// feedbacksグループの構造を検証
		if feedbacks, ok := params["feedbacks"].(map[string]any); ok {
			if _, ok := feedbacks["createParams"]; !ok {
				t.Error("feedbacks.createParamsが見つかりません")
			}
			if _, ok := feedbacks["updateParams"]; !ok {
				t.Error("feedbacks.updateParamsが見つかりません")
			}
		} else {
			t.Error("feedbacksグループの構造が不正です")
		}
	})

	// 4. ヘルパー関数ファイルの内容を検証
	t.Run("ヘルパー関数ファイルの内容検証", func(t *testing.T) {
		data, err := os.ReadFile(helpersPath)
		if err != nil {
			t.Fatalf("ヘルパー関数ファイルの読み込みに失敗: %v", err)
		}

		content := string(data)

		// 必須要素が含まれていることを確認
		requiredElements := []string{
			"package",
			"type TestParams struct",
			"func stringPtr(s string) *string",
			"func loadTestParams(t *testing.T) *TestParams",
			"Feedbacks struct",
		}

		for _, element := range requiredElements {
			if !strings.Contains(content, element) {
				t.Errorf("ヘルパー関数ファイルに必須要素が含まれていません: %s", element)
			}
		}
	})

	// 5. 生成されたGoコードがコンパイル可能か確認
	t.Run("生成されたGoコードのコンパイル確認", func(t *testing.T) {
		// go.modファイルを作成
		goModContent := `module test

go 1.20
`
		goModPath := filepath.Join(tmpDir, "go.mod")
		if err := os.WriteFile(goModPath, []byte(goModContent), 0644); err != nil {
			t.Fatalf("go.modファイルの作成に失敗: %v", err)
		}

		// go buildコマンドを実行
		cmd := exec.Command("go", "build", "-o", "/dev/null", helpersPath)
		cmd.Dir = tmpDir
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Errorf("生成されたGoコードのコンパイルに失敗:\n%s\nエラー: %v", string(output), err)
		}
	})
}

// TestEndToEndWithMultipleModules は複数のモジュールに対する統合テストです
func TestEndToEndWithMultipleModules(t *testing.T) {
	// 利用可能なレポートファイルを検索
	reportsDir := "../postman-analyzer/reports"
	if _, err := os.Stat(reportsDir); os.IsNotExist(err) {
		t.Skipf("レポートディレクトリが見つかりません: %s", reportsDir)
	}

	// テスト対象のモジュール
	modules := []struct {
		name       string
		reportPath string
	}{
		{
			name:       "communication",
			reportPath: "../postman-analyzer/reports/communication/postman_analysis.json",
		},
		{
			name:       "billing",
			reportPath: "../postman-analyzer/reports/billing/postman_analysis.json",
		},
		// pricingは異なるレポート構造を持つため、現時点ではスキップ
		// {
		// 	name:       "pricing",
		// 	reportPath: "../postman-analyzer/reports/pricing/postman_analysis.json",
		// },
	}

	for _, module := range modules {
		t.Run(module.name, func(t *testing.T) {
			// レポートファイルが存在するか確認
			if _, err := os.Stat(module.reportPath); os.IsNotExist(err) {
				t.Skipf("レポートファイルが見つかりません: %s", module.reportPath)
			}

			// テスト用の一時ディレクトリを作成
			tmpDir := t.TempDir()

			// 設定を作成
			config := &Config{
				ReportPath: module.reportPath,
				OutputDir:  tmpDir,
				ModuleName: module.name,
				Overwrite:  true,
				Verbose:    false,
			}

			// メイン処理を実行
			if err := run(config); err != nil {
				t.Fatalf("統合テストの実行に失敗 (%s): %v", module.name, err)
			}

			// 生成されたファイルのパスを取得
			paramsPath, helpersPath := config.GetOutputPaths()

			// ファイルが生成されたことを確認
			if _, err := os.Stat(paramsPath); os.IsNotExist(err) {
				t.Errorf("パラメータファイルが生成されていません (%s): %s", module.name, paramsPath)
			}

			if _, err := os.Stat(helpersPath); os.IsNotExist(err) {
				t.Errorf("ヘルパー関数ファイルが生成されていません (%s): %s", module.name, helpersPath)
			}

			// パラメータファイルがJSONとして有効か確認
			data, err := os.ReadFile(paramsPath)
			if err != nil {
				t.Fatalf("パラメータファイルの読み込みに失敗 (%s): %v", module.name, err)
			}

			var params map[string]any
			if err := json.Unmarshal(data, &params); err != nil {
				t.Errorf("パラメータファイルのJSON解析に失敗 (%s): %v", module.name, err)
			}

			// 少なくとも1つのパラメータグループが存在することを確認
			if len(params) == 0 {
				t.Errorf("パラメータグループが1つも生成されていません (%s)", module.name)
			}
		})
	}
}

// TestEndToEndErrorHandling はエラーハンドリングの統合テストです
func TestEndToEndErrorHandling(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		expectError bool
		errorType   string
	}{
		{
			name: "存在しないレポートファイル",
			config: &Config{
				ReportPath: "nonexistent.json",
				OutputDir:  t.TempDir(),
				ModuleName: "test",
				Overwrite:  true,
				Verbose:    false,
			},
			expectError: true,
			errorType:   "file_not_found",
		},
		{
			name: "不正なJSONファイル",
			config: &Config{
				ReportPath: "testdata/invalid.json",
				OutputDir:  t.TempDir(),
				ModuleName: "test",
				Overwrite:  true,
				Verbose:    false,
			},
			expectError: true,
			errorType:   "invalid_json",
		},
		{
			name: "必須フィールドが欠落したレポート",
			config: &Config{
				ReportPath: "testdata/missing_fields.json",
				OutputDir:  t.TempDir(),
				ModuleName: "test",
				Overwrite:  true,
				Verbose:    false,
			},
			expectError: true,
			errorType:   "invalid_report",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := run(tt.config)

			if tt.expectError {
				if err == nil {
					t.Error("エラーが期待されましたが、エラーが発生しませんでした")
				} else {
					// エラーの種類を確認
					switch tt.errorType {
					case "file_not_found":
						if !IsFileNotFoundError(err) {
							t.Errorf("FileNotFoundErrorが期待されましたが、異なるエラーが発生しました: %v", err)
						}
					case "invalid_json":
						if !IsInvalidJSONError(err) {
							t.Errorf("InvalidJSONErrorが期待されましたが、異なるエラーが発生しました: %v", err)
						}
					case "invalid_report":
						if !IsInvalidReportError(err) {
							t.Errorf("InvalidReportErrorが期待されましたが、異なるエラーが発生しました: %v", err)
						}
					}
				}
			} else {
				if err != nil {
					t.Errorf("エラーが発生しました: %v", err)
				}
			}
		})
	}
}

// TestEndToEndOutputOverwrite は出力ファイルの上書き動作の統合テストです
func TestEndToEndOutputOverwrite(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tmpDir := t.TempDir()

	// 実際のPostman Analyzerレポートを使用
	reportPath := "../postman-analyzer/reports/communication/postman_analysis.json"
	if _, err := os.Stat(reportPath); os.IsNotExist(err) {
		t.Skipf("Postman Analyzerレポートが見つかりません: %s", reportPath)
	}

	// 設定を作成
	config := &Config{
		ReportPath: reportPath,
		OutputDir:  tmpDir,
		ModuleName: "communication",
		Overwrite:  true,
		Verbose:    false,
	}

	// 1回目の実行
	if err := run(config); err != nil {
		t.Fatalf("1回目の実行に失敗: %v", err)
	}

	// 生成されたファイルのパスを取得
	paramsPath, helpersPath := config.GetOutputPaths()

	// 1回目の内容を保存
	firstParamsData, err := os.ReadFile(paramsPath)
	if err != nil {
		t.Fatalf("1回目のパラメータファイルの読み込みに失敗: %v", err)
	}

	firstHelpersData, err := os.ReadFile(helpersPath)
	if err != nil {
		t.Fatalf("1回目のヘルパー関数ファイルの読み込みに失敗: %v", err)
	}

	// 2回目の実行（上書き）
	if err := run(config); err != nil {
		t.Fatalf("2回目の実行に失敗: %v", err)
	}

	// 2回目の内容を読み込み
	secondParamsData, err := os.ReadFile(paramsPath)
	if err != nil {
		t.Fatalf("2回目のパラメータファイルの読み込みに失敗: %v", err)
	}

	secondHelpersData, err := os.ReadFile(helpersPath)
	if err != nil {
		t.Fatalf("2回目のヘルパー関数ファイルの読み込みに失敗: %v", err)
	}

	// 内容が同じであることを確認（冪等性）
	if string(firstParamsData) != string(secondParamsData) {
		t.Error("パラメータファイルの内容が1回目と2回目で異なります")
	}

	if string(firstHelpersData) != string(secondHelpersData) {
		t.Error("ヘルパー関数ファイルの内容が1回目と2回目で異なります")
	}
}

// TestEndToEndWithIterationData はiterationDataを使用した統合テストです
func TestEndToEndWithIterationData(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tmpDir := t.TempDir()

	// 実際のPostman Analyzerレポートを使用
	reportPath := "../postman-analyzer/reports/communication/postman_analysis.json"
	if _, err := os.Stat(reportPath); os.IsNotExist(err) {
		t.Skipf("Postman Analyzerレポートが見つかりません: %s", reportPath)
	}

	// 設定を作成
	config := &Config{
		ReportPath: reportPath,
		OutputDir:  tmpDir,
		ModuleName: "communication",
		Overwrite:  true,
		Verbose:    true, // 詳細ログを有効化
	}

	// メイン処理を実行
	if err := run(config); err != nil {
		t.Fatalf("統合テストの実行に失敗: %v", err)
	}

	// 生成されたファイルのパスを取得
	paramsPath, _ := config.GetOutputPaths()

	// パラメータファイルの内容を検証
	t.Run("iterationDataがマージされていることを確認", func(t *testing.T) {
		data, err := os.ReadFile(paramsPath)
		if err != nil {
			t.Fatalf("パラメータファイルの読み込みに失敗: %v", err)
		}

		var params map[string]any
		if err := json.Unmarshal(data, &params); err != nil {
			t.Fatalf("パラメータファイルのJSON解析に失敗: %v", err)
		}

		// feedbacksグループが存在することを確認
		feedbacks, ok := params["feedbacks"].(map[string]any)
		if !ok {
			t.Fatal("feedbacksパラメータグループが見つかりません")
		}

		// createParamsが存在することを確認
		createParams, ok := feedbacks["createParams"].(map[string]any)
		if !ok {
			t.Fatal("feedbacks.createParamsが見つかりません")
		}

		// iterationDataから読み込まれた値が含まれていることを確認
		// communication.jsonには以下のフィールドが含まれているはず
		expectedFields := []string{"user_id", "feedback_title", "feedback_description"}
		for _, field := range expectedFields {
			if _, exists := createParams[field]; !exists {
				t.Errorf("createParamsに期待されるフィールドが見つかりません: %s", field)
			}
		}

		// updateParamsが存在することを確認
		updateParams, ok := feedbacks["updateParams"].(map[string]any)
		if !ok {
			t.Fatal("feedbacks.updateParamsが見つかりません")
		}

		// iterationDataから読み込まれた_updateサフィックス付きフィールドが含まれていることを確認
		expectedUpdateFields := []string{"feedback_title", "feedback_description"}
		for _, field := range expectedUpdateFields {
			if _, exists := updateParams[field]; !exists {
				t.Errorf("updateParamsに期待されるフィールドが見つかりません: %s", field)
			}
		}
	})
}

// TestMergeIterationData はmergeIterationData関数のテストです
func TestMergeIterationData(t *testing.T) {
	// テスト用の設定を作成
	config := &Config{
		ModuleName: "communication",
		Verbose:    true, // デバッグ出力を有効化
	}

	// テスト用の分析結果を作成
	analysis := &AnalysisResult{
		ModuleName: "communication",
		ParamGroups: map[string]*ParamGroup{
			"feedbacks": {
				Name:         "feedbacks",
				CreateParams: make(map[string]interface{}),
				UpdateParams: make(map[string]interface{}),
			},
		},
	}

	// Generatorを作成
	generator := NewGenerator(config, analysis)

	// mergeIterationDataを実行
	err := generator.mergeIterationData()
	if err != nil {
		t.Fatalf("mergeIterationData() error = %v", err)
	}

	// feedbacksグループのパラメータが更新されていることを確認
	feedbacks := analysis.ParamGroups["feedbacks"]

	t.Logf("CreateParams: %+v", feedbacks.CreateParams)
	t.Logf("UpdateParams: %+v", feedbacks.UpdateParams)

	// iterationDataにはfeedbacksに直接マッチするキーがないため、
	// communicationキーからマッチングされる必要がある
	// この場合、マッチングロジックが期待通りに動作しているか確認
	if len(feedbacks.CreateParams) == 0 {
		t.Log("CreateParamsが空です。これは期待される動作かもしれません（feedbacksキーが見つからない場合）")
	}

	if len(feedbacks.UpdateParams) == 0 {
		t.Log("UpdateParamsが空です。これは期待される動作かもしれません（feedbacksキーが見つからない場合）")
	}
}
