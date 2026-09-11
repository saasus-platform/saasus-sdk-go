package integration

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// TestParams はテストパラメータファイルの構造体
type TestParams struct {

	EventbridgeEvent struct {

		CreateParams map[string]any `json:"createParams"`


	} `json:"eventbridge-event"`

	EventbridgeInfo struct {

		CreateParams map[string]any `json:"createParams"`


		UpdateParams map[string]any `json:"updateParams"`

	} `json:"eventbridge-info"`

}


// stringPtr は文字列のポインタを返すヘルパー関数
// 要件: 5.1, 5.2
func stringPtr(s string) *string {
	return &s
}

// loadTestParams はテストパラメータファイルを読み込みます
// 要件: 5.3, 5.4, 5.5
func loadTestParams(t *testing.T) *TestParams {
	paramFile := filepath.Join("test_params.json")
	data, err := os.ReadFile(paramFile)
	if err != nil {
		t.Fatalf("テストパラメータファイルの読み込みに失敗: %v", err)
	}

	var params TestParams
	if err := json.Unmarshal(data, &params); err != nil {
		t.Fatalf("テストパラメータファイルの解析に失敗: %v", err)
	}

	return &params
}
