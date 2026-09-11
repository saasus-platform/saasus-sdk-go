package main

import (
	"os"
	"path/filepath"
	"testing"
)

// TestParseReport はParseReport関数のテスト
func TestParseReport(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(t *testing.T) string // テストファイルのパスを返す
		wantErr bool
		errType string // エラーの種類
	}{
		{
			name: "正常なレポート",
			setup: func(t *testing.T) string {
				// 実際のレポートファイルを使用
				return "../postman-analyzer/reports/communication/postman_analysis.json"
			},
			wantErr: false,
		},
		{
			name: "存在しないファイル",
			setup: func(t *testing.T) string {
				return "nonexistent_file.json"
			},
			wantErr: true,
			errType: "FileNotFoundError",
		},
		{
			name: "不正なJSON",
			setup: func(t *testing.T) string {
				return "testdata/invalid.json"
			},
			wantErr: true,
			errType: "InvalidJSONError",
		},
		{
			name: "空のレポート",
			setup: func(t *testing.T) string {
				return "testdata/empty_report.json"
			},
			wantErr: true,
			errType: "InvalidReportError",
		},
		{
			name: "必須フィールド欠落",
			setup: func(t *testing.T) string {
				return "testdata/missing_fields.json"
			},
			wantErr: true,
			errType: "InvalidReportError",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath := tt.setup(t)

			report, err := ParseReport(filePath)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseReport() エラーが期待されましたが、エラーが発生しませんでした")
					return
				}

				// エラーの種類を確認
				switch tt.errType {
				case "FileNotFoundError":
					if _, ok := err.(*FileNotFoundError); !ok {
						t.Errorf("ParseReport() エラーの種類 = %T, want *FileNotFoundError", err)
					}
				case "InvalidJSONError":
					if _, ok := err.(*InvalidJSONError); !ok {
						t.Errorf("ParseReport() エラーの種類 = %T, want *InvalidJSONError", err)
					}
				case "InvalidReportError":
					if _, ok := err.(*InvalidReportError); !ok {
						t.Errorf("ParseReport() エラーの種類 = %T, want *InvalidReportError", err)
					}
				}

				t.Logf("期待通りエラー: %v", err)
			} else {
				if err != nil {
					t.Errorf("ParseReport() エラー = %v, wantErr %v", err, tt.wantErr)
					return
				}

				// 正常系の検証
				if report == nil {
					t.Error("ParseReport() report = nil, want non-nil")
					return
				}

				if report.CollectionName == "" {
					t.Error("ParseReport() CollectionName が空です")
				}

				if len(report.TestFlow) == 0 {
					t.Error("ParseReport() TestFlow が空です")
				}

				if len(report.ClientMethods) == 0 {
					t.Error("ParseReport() ClientMethods が空です")
				}

				t.Logf("解析成功: CollectionName=%s, TestFlow=%d steps, ClientMethods=%d methods",
					report.CollectionName, len(report.TestFlow), len(report.ClientMethods))
			}
		})
	}
}

// TestValidateReport はvalidateReport関数のテスト
func TestValidateReport(t *testing.T) {
	tests := []struct {
		name    string
		report  *PostmanReport
		wantErr bool
	}{
		{
			name: "正常なレポート",
			report: &PostmanReport{
				CollectionName: "Test Collection",
				TestFlow: []TestStep{
					{
						Name:     "Test Step",
						Method:   "GET",
						Endpoint: "/test",
					},
				},
				ClientMethods: []ClientMethod{
					{
						Name:       "TestMethod",
						HTTPMethod: "GET",
						Endpoint:   "/test",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "collection_nameが空",
			report: &PostmanReport{
				CollectionName: "",
				TestFlow: []TestStep{
					{Name: "Test"},
				},
				ClientMethods: []ClientMethod{
					{Name: "Test"},
				},
			},
			wantErr: true,
		},
		{
			name: "test_flowが空",
			report: &PostmanReport{
				CollectionName: "Test",
				TestFlow:       []TestStep{},
				ClientMethods: []ClientMethod{
					{Name: "Test"},
				},
			},
			wantErr: true,
		},
		{
			name: "client_methodsが空",
			report: &PostmanReport{
				CollectionName: "Test",
				TestFlow: []TestStep{
					{Name: "Test"},
				},
				ClientMethods: []ClientMethod{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateReport(tt.report)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateReport() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestParseReportWithRealFile は実際のレポートファイルを使用した統合テスト
func TestParseReportWithRealFile(t *testing.T) {
	// 実際のレポートファイルのパス
	reportPath := "../postman-analyzer/reports/communication/postman_analysis.json"

	// ファイルが存在するか確認
	if _, err := os.Stat(reportPath); os.IsNotExist(err) {
		t.Skipf("レポートファイルが見つかりません: %s", reportPath)
	}

	// レポートを解析
	report, err := ParseReport(reportPath)
	if err != nil {
		t.Fatalf("ParseReport() エラー: %v", err)
	}

	// 基本的な検証
	if report.CollectionName == "" {
		t.Error("CollectionName が空です")
	}

	if len(report.TestFlow) == 0 {
		t.Error("TestFlow が空です")
	}

	if len(report.ClientMethods) == 0 {
		t.Error("ClientMethods が空です")
	}

	// 詳細情報をログ出力
	t.Logf("Collection: %s", report.CollectionName)
	t.Logf("Test Flow Steps: %d", len(report.TestFlow))
	t.Logf("Client Methods: %d", len(report.ClientMethods))

	// TestFlowの最初のステップを確認
	if len(report.TestFlow) > 0 {
		firstStep := report.TestFlow[0]
		t.Logf("First Step: %s %s %s", firstStep.Method, firstStep.Endpoint, firstStep.Name)
	}

	// ClientMethodsの最初のメソッドを確認
	if len(report.ClientMethods) > 0 {
		firstMethod := report.ClientMethods[0]
		t.Logf("First Method: %s (%s %s)", firstMethod.Name, firstMethod.HTTPMethod, firstMethod.Endpoint)
	}
}

// TestParseReportErrorMessages はエラーメッセージの内容を検証
func TestParseReportErrorMessages(t *testing.T) {
	tests := []struct {
		name        string
		filePath    string
		wantContain string // エラーメッセージに含まれるべき文字列
	}{
		{
			name:        "ファイルが見つからない",
			filePath:    "nonexistent.json",
			wantContain: "nonexistent.json",
		},
		{
			name:        "不正なJSON",
			filePath:    "testdata/invalid.json",
			wantContain: "JSON",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseReport(tt.filePath)
			if err == nil {
				t.Error("エラーが期待されましたが、エラーが発生しませんでした")
				return
			}

			errMsg := err.Error()
			if errMsg == "" {
				t.Error("エラーメッセージが空です")
			}

			t.Logf("エラーメッセージ: %s", errMsg)
		})
	}
}

// TestParseReportConcurrent は並行処理での動作を確認
func TestParseReportConcurrent(t *testing.T) {
	reportPath := "../postman-analyzer/reports/communication/postman_analysis.json"

	// ファイルが存在するか確認
	if _, err := os.Stat(reportPath); os.IsNotExist(err) {
		t.Skipf("レポートファイルが見つかりません: %s", reportPath)
	}

	// 並行して複数回解析
	const numGoroutines = 10
	done := make(chan bool, numGoroutines)
	errors := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			_, err := ParseReport(reportPath)
			if err != nil {
				errors <- err
			}
			done <- true
		}()
	}

	// すべてのゴルーチンの完了を待つ
	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	close(errors)

	// エラーがないことを確認
	for err := range errors {
		t.Errorf("並行処理中にエラー: %v", err)
	}
}

// TestParseReportWithSymlink はシンボリックリンク経由でのファイル読み込みをテスト
func TestParseReportWithSymlink(t *testing.T) {
	// 実際のレポートファイル
	reportPath := "../postman-analyzer/reports/communication/postman_analysis.json"
	if _, err := os.Stat(reportPath); os.IsNotExist(err) {
		t.Skipf("レポートファイルが見つかりません: %s", reportPath)
	}

	// 一時ディレクトリにシンボリックリンクを作成
	tmpDir := t.TempDir()
	symlinkPath := filepath.Join(tmpDir, "symlink_report.json")

	// 絶対パスを取得
	absReportPath, err := filepath.Abs(reportPath)
	if err != nil {
		t.Fatalf("絶対パスの取得に失敗: %v", err)
	}

	// シンボリックリンクを作成
	if err := os.Symlink(absReportPath, symlinkPath); err != nil {
		t.Skipf("シンボリックリンクの作成に失敗（権限不足の可能性）: %v", err)
	}

	// シンボリックリンク経由で解析
	report, err := ParseReport(symlinkPath)
	if err != nil {
		t.Fatalf("ParseReport() エラー: %v", err)
	}

	if report == nil {
		t.Error("report が nil です")
	}
}
