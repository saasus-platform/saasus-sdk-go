package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindIterationDataFiles(t *testing.T) {
	tests := []struct {
		name       string
		moduleName string
		wantFiles  int
		wantErr    bool
	}{
		{
			name:       "communication module",
			moduleName: "communication",
			wantFiles:  1, // communication.json
			wantErr:    false,
		},
		{
			name:       "pricing module",
			moduleName: "pricing",
			wantFiles:  15, // 複数のサブディレクトリにJSONファイルがある
			wantErr:    false,
		},
		{
			name:       "non-existent module",
			moduleName: "nonexistent",
			wantFiles:  0,
			wantErr:    false, // エラーではなく空の結果を返す
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loader := NewIterationDataLoader(tt.moduleName, false)
			files, err := loader.findIterationDataFiles()

			if (err != nil) != tt.wantErr {
				t.Errorf("findIterationDataFiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(files) < tt.wantFiles {
				t.Errorf("findIterationDataFiles() got %d files, want at least %d", len(files), tt.wantFiles)
			}

			// 見つかったファイルがすべて.json拡張子を持つことを確認
			for _, file := range files {
				if filepath.Ext(file) != ".json" {
					t.Errorf("findIterationDataFiles() returned non-JSON file: %s", file)
				}
			}
		})
	}
}

func TestLoadJSONFile(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tmpDir := t.TempDir()

	tests := []struct {
		name        string
		content     string
		wantLen     int
		wantErr     bool
		description string
	}{
		{
			name:        "array format",
			content:     `[{"field1": "value1"}, {"field2": "value2"}]`,
			wantLen:     2,
			wantErr:     false,
			description: "配列形式のJSONを正しく読み込む",
		},
		{
			name:        "single object",
			content:     `{"field1": "value1", "field2": "value2"}`,
			wantLen:     1,
			wantErr:     false,
			description: "単一オブジェクトを配列にラップして返す",
		},
		{
			name:        "invalid json",
			content:     `{invalid json}`,
			wantLen:     0,
			wantErr:     true,
			description: "不正なJSONはエラーを返す",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// テストファイルを作成
			filePath := filepath.Join(tmpDir, tt.name+".json")
			if err := os.WriteFile(filePath, []byte(tt.content), 0644); err != nil {
				t.Fatalf("テストファイルの作成に失敗: %v", err)
			}

			loader := NewIterationDataLoader("test", false)
			data, err := loader.loadJSONFile(filePath)

			if (err != nil) != tt.wantErr {
				t.Errorf("loadJSONFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && len(data) != tt.wantLen {
				t.Errorf("loadJSONFile() got %d elements, want %d", len(data), tt.wantLen)
			}
		})
	}
}

func TestLoadIterationData(t *testing.T) {
	tests := []struct {
		name       string
		moduleName string
		wantEmpty  bool
		wantErr    bool
	}{
		{
			name:       "communication module",
			moduleName: "communication",
			wantEmpty:  false,
			wantErr:    false,
		},
		{
			name:       "non-existent module",
			moduleName: "nonexistent",
			wantEmpty:  true,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loader := NewIterationDataLoader(tt.moduleName, false)
			data, err := loader.LoadIterationData()

			if (err != nil) != tt.wantErr {
				t.Errorf("LoadIterationData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantEmpty && len(data) > 0 {
				t.Errorf("LoadIterationData() expected empty result, got %d items", len(data))
			}

			if !tt.wantEmpty && len(data) == 0 {
				t.Errorf("LoadIterationData() expected non-empty result, got empty")
			}
		})
	}
}

func TestNewIterationDataLoader(t *testing.T) {
	moduleName := "test-module"
	verbose := true

	loader := NewIterationDataLoader(moduleName, verbose)

	if loader.ModuleName != moduleName {
		t.Errorf("NewIterationDataLoader() ModuleName = %v, want %v", loader.ModuleName, moduleName)
	}

	if loader.Verbose != verbose {
		t.Errorf("NewIterationDataLoader() Verbose = %v, want %v", loader.Verbose, verbose)
	}

	// BaseDir should be either "iterationData" or "tests/tools/test-param-generator/iterationData"
	// depending on where the test is run from
	if loader.BaseDir != "iterationData" && loader.BaseDir != "tests/tools/test-param-generator/iterationData" {
		t.Errorf("NewIterationDataLoader() BaseDir = %v, want either 'iterationData' or 'tests/tools/test-param-generator/iterationData'", loader.BaseDir)
	}
}

func TestExtractFieldsBySuffix(t *testing.T) {
	tests := []struct {
		name     string
		data     map[string]interface{}
		suffix   string
		wantKeys []string
	}{
		{
			name: "extract _update fields",
			data: map[string]interface{}{
				"field1":        "value1",
				"field2_update": "value2",
				"field3_update": "value3",
				"field4_create": "value4",
			},
			suffix:   "_update",
			wantKeys: []string{"field2", "field3"},
		},
		{
			name: "extract _create fields",
			data: map[string]interface{}{
				"field1":        "value1",
				"field2_create": "value2",
				"field3_update": "value3",
			},
			suffix:   "_create",
			wantKeys: []string{"field2"},
		},
		{
			name: "no matching fields",
			data: map[string]interface{}{
				"field1": "value1",
				"field2": "value2",
			},
			suffix:   "_update",
			wantKeys: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractFieldsBySuffix(tt.data, tt.suffix)

			if len(result) != len(tt.wantKeys) {
				t.Errorf("extractFieldsBySuffix() got %d fields, want %d", len(result), len(tt.wantKeys))
			}

			for _, key := range tt.wantKeys {
				if _, exists := result[key]; !exists {
					t.Errorf("extractFieldsBySuffix() missing expected key: %s", key)
				}
			}
		})
	}
}

func TestRemoveFieldSuffix(t *testing.T) {
	tests := []struct {
		name      string
		fieldName string
		suffix    string
		want      string
	}{
		{
			name:      "remove _update suffix",
			fieldName: "field_name_update",
			suffix:    "_update",
			want:      "field_name",
		},
		{
			name:      "remove _create suffix",
			fieldName: "field_name_create",
			suffix:    "_create",
			want:      "field_name",
		},
		{
			name:      "no suffix to remove",
			fieldName: "field_name",
			suffix:    "_update",
			want:      "field_name",
		},
		{
			name:      "different suffix",
			fieldName: "field_name_delete",
			suffix:    "_update",
			want:      "field_name_delete",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := removeFieldSuffix(tt.fieldName, tt.suffix)
			if got != tt.want {
				t.Errorf("removeFieldSuffix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMatchIterationDataToParamGroup(t *testing.T) {
	tests := []struct {
		name          string
		iterationData map[string][]map[string]interface{}
		resourceName  string
		wantCreate    bool
		wantUpdate    bool
	}{
		{
			name: "match with _create and _update suffixes",
			iterationData: map[string][]map[string]interface{}{
				"feedbacks": {
					{
						"title":        "test title",
						"title_create": "create title",
						"title_update": "update title",
						"body":         "test body",
					},
				},
			},
			resourceName: "feedbacks",
			wantCreate:   true,
			wantUpdate:   true,
		},
		{
			name: "match without suffixes",
			iterationData: map[string][]map[string]interface{}{
				"comments": {
					{
						"body":   "test comment",
						"author": "test author",
					},
				},
			},
			resourceName: "comments",
			wantCreate:   true,
			wantUpdate:   false,
		},
		{
			name: "no match with multiple entries",
			iterationData: map[string][]map[string]interface{}{
				"other1": {
					{
						"field": "value",
					},
				},
				"other2": {
					{
						"field": "value",
					},
				},
			},
			resourceName: "feedbacks",
			wantCreate:   false,
			wantUpdate:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			paramGroup := &ParamGroup{
				Name: tt.resourceName,
			}

			err := matchIterationDataToParamGroup(tt.iterationData, paramGroup, tt.resourceName)
			if err != nil {
				t.Errorf("matchIterationDataToParamGroup() error = %v", err)
				return
			}

			hasCreate := len(paramGroup.CreateParams) > 0
			hasUpdate := len(paramGroup.UpdateParams) > 0

			if hasCreate != tt.wantCreate {
				t.Errorf("matchIterationDataToParamGroup() CreateParams exists = %v, want %v", hasCreate, tt.wantCreate)
			}

			if hasUpdate != tt.wantUpdate {
				t.Errorf("matchIterationDataToParamGroup() UpdateParams exists = %v, want %v", hasUpdate, tt.wantUpdate)
			}
		})
	}
}
