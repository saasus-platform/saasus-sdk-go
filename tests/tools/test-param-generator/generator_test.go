package main

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestGenerateParamsFile(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tmpDir := t.TempDir()

	// テスト用の設定を作成
	config := &Config{
		ReportPath: "testdata/empty_report.json",
		OutputDir:  tmpDir,
		ModuleName: "communication",
		Overwrite:  true,
		Verbose:    false,
	}

	// テスト用の分析結果を作成
	analysis := &AnalysisResult{
		ModuleName: "communication",
		ParamGroups: map[string]*ParamGroup{
			"feedbacks": {
				Name: "feedbacks",
				CreateParams: map[string]interface{}{
					"user_id":              "69e732d6-8ecc-45c4-c2eb-8438f7ffe775",
					"feedback_title":       "テスト用フィードバック",
					"feedback_description": "これはE2Eテスト用のフィードバックです。",
				},
				UpdateParams: map[string]interface{}{
					"feedback_title":       "更新されたテスト用フィードバック",
					"feedback_description": "これは更新されたE2Eテスト用のフィードバックです。",
				},
			},
		},
		UtilityFuncs: []string{},
	}

	// Generatorを作成
	generator := NewGenerator(config, analysis)

	// パラメータファイルを生成
	err := generator.GenerateParamsFile()
	if err != nil {
		t.Fatalf("GenerateParamsFile() エラー: %v", err)
	}

	// 生成されたファイルを確認
	paramsPath, _ := config.GetOutputPaths()
	if _, err := os.Stat(paramsPath); os.IsNotExist(err) {
		t.Fatalf("パラメータファイルが生成されていません: %s", paramsPath)
	}

	// ファイルの内容を読み込み
	data, err := os.ReadFile(paramsPath)
	if err != nil {
		t.Fatalf("パラメータファイルの読み込みに失敗: %v", err)
	}

	// JSONとして解析可能か確認
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("生成されたJSONの解析に失敗: %v", err)
	}

	// feedbacksグループが存在するか確認
	feedbacks, ok := result["feedbacks"].(map[string]interface{})
	if !ok {
		t.Fatal("feedbacksグループが見つかりません")
	}

	// createParamsが存在するか確認
	createParams, ok := feedbacks["createParams"].(map[string]interface{})
	if !ok {
		t.Fatal("createParamsが見つかりません")
	}

	// user_idが存在するか確認
	if _, ok := createParams["user_id"]; !ok {
		t.Error("user_idが見つかりません")
	}

	// updateParamsが存在するか確認
	if _, ok := feedbacks["updateParams"]; !ok {
		t.Error("updateParamsが見つかりません")
	}
}

func TestBuildParamsStructure(t *testing.T) {
	// テスト用の設定を作成
	config := &Config{
		ModuleName: "test",
		OutputDir:  "/tmp",
		Overwrite:  true,
		Verbose:    false,
	}

	// テスト用の分析結果を作成
	analysis := &AnalysisResult{
		ModuleName: "test",
		ParamGroups: map[string]*ParamGroup{
			"items": {
				Name: "items",
				CreateParams: map[string]interface{}{
					"name": "テスト項目",
				},
				UpdateParams: map[string]interface{}{
					"name": "更新されたテスト項目",
				},
			},
		},
	}

	// Generatorを作成
	generator := NewGenerator(config, analysis)

	// パラメータ構造を構築
	result := generator.buildParamsStructure()

	// itemsグループが存在するか確認
	items, ok := result["items"].(map[string]interface{})
	if !ok {
		t.Fatal("itemsグループが見つかりません")
	}

	// createParamsが存在するか確認
	if _, ok := items["createParams"]; !ok {
		t.Error("createParamsが見つかりません")
	}

	// updateParamsが存在するか確認
	if _, ok := items["updateParams"]; !ok {
		t.Error("updateParamsが見つかりません")
	}
}

func TestFormatJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		wantErr bool
	}{
		{
			name: "simple map",
			input: map[string]interface{}{
				"key": "value",
			},
			wantErr: false,
		},
		{
			name: "nested structure",
			input: map[string]interface{}{
				"group": map[string]interface{}{
					"createParams": map[string]interface{}{
						"name": "test",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "with array",
			input: map[string]interface{}{
				"ids": []string{"id1", "id2"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := formatJSON(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("formatJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(got) == 0 {
				t.Error("formatJSON() returned empty result")
			}

			// JSONとして解析可能か確認
			if !tt.wantErr {
				var result interface{}
				if err := json.Unmarshal(got, &result); err != nil {
					t.Errorf("生成されたJSONの解析に失敗: %v", err)
				}
			}
		})
	}
}

func TestGenerateHelpersFile(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tmpDir := t.TempDir()

	// テスト用の設定を作成
	config := &Config{
		ReportPath: "testdata/empty_report.json",
		OutputDir:  tmpDir,
		ModuleName: "communication",
		Overwrite:  true,
		Verbose:    false,
	}

	// テスト用の分析結果を作成
	analysis := &AnalysisResult{
		ModuleName: "communication",
		ParamGroups: map[string]*ParamGroup{
			"feedbacks": {
				Name: "feedbacks",
				CreateParams: map[string]interface{}{
					"user_id": "test-id",
				},
				UpdateParams: map[string]interface{}{
					"title": "test",
				},
			},
		},
	}

	// Generatorを作成
	generator := NewGenerator(config, analysis)

	// ヘルパーファイルを生成
	err := generator.GenerateHelpersFile()
	if err != nil {
		t.Fatalf("GenerateHelpersFile() エラー: %v", err)
	}

	// 生成されたファイルを確認
	_, helpersPath := config.GetOutputPaths()
	if _, err := os.Stat(helpersPath); os.IsNotExist(err) {
		t.Fatalf("ヘルパーファイルが生成されていません: %s", helpersPath)
	}

	// ファイルの内容を読み込み
	data, err := os.ReadFile(helpersPath)
	if err != nil {
		t.Fatalf("ヘルパーファイルの読み込みに失敗: %v", err)
	}

	content := string(data)

	// 必須要素が含まれているか確認
	requiredElements := []string{
		"package ",
		"type TestParams struct",
		"func stringPtr",
		"func loadTestParams",
		"Feedbacks struct",
	}

	for _, elem := range requiredElements {
		if !strings.Contains(content, elem) {
			t.Errorf("生成されたコードに '%s' が含まれていません", elem)
		}
	}
}
