package main

import (
	"encoding/json"
	"os"
)

// ParseReport はPostman Analyzerレポートファイルを読み込んで解析します
// 要件: 1.1, 1.2, 8.1, 8.2
func ParseReport(filePath string) (*PostmanReport, error) {
	// ファイルの読み込み
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, NewFileNotFoundError(filePath, err)
		}
		return nil, WrapError(err, "レポートファイルの読み込みに失敗: %s", filePath)
	}

	// JSON解析
	var report PostmanReport
	if err := json.Unmarshal(data, &report); err != nil {
		return nil, NewInvalidJSONError(filePath, err)
	}

	// レポート構造の検証
	if err := validateReport(&report); err != nil {
		return nil, err
	}

	return &report, nil
}

// validateReport はレポートの必須フィールドを検証します
// 要件: 1.3
func validateReport(report *PostmanReport) error {
	// collection_nameの検証
	if report.CollectionName == "" {
		return NewInvalidReportError("collection_name", "必須フィールドが空です")
	}

	// test_flowの検証
	if len(report.TestFlow) == 0 {
		return NewInvalidReportError("test_flow", "テストフローが定義されていません")
	}

	// client_methodsの検証
	if len(report.ClientMethods) == 0 {
		return NewInvalidReportError("client_methods", "クライアントメソッドが定義されていません")
	}

	return nil
}
