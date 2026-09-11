package main

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

// TemplateEngine はテンプレートの管理と実行を担当します
// 要件: 3.1, 3.2
type TemplateEngine struct {
	templates map[string]*template.Template
}

// NewTemplateEngine は新しいTemplateEngineを作成します
// 要件: 3.1, 3.2
func NewTemplateEngine() *TemplateEngine {
	engine := &TemplateEngine{
		templates: make(map[string]*template.Template),
	}

	// テンプレート関数を定義
	funcMap := template.FuncMap{
		"capitalize":        capitalize,
		"singularize":       singularize,
		"hasCreateParams":   hasCreateParams,
		"hasUpdateParams":   hasUpdateParams,
		"structFieldName":   structFieldName,
		"sanitizeJSONField": sanitizeJSONField,
	}

	// ヘルパー関数テンプレートをコンパイル
	helpersTemplate := template.Must(template.New("helpers").Funcs(funcMap).Parse(helpersTemplateStr))
	engine.templates["helpers"] = helpersTemplate

	// 構造体定義テンプレートをコンパイル
	structsTemplate := template.Must(template.New("structs").Funcs(funcMap).Parse(structsTemplateStr))
	engine.templates["structs"] = structsTemplate

	return engine
}

// RenderHelpers はヘルパー関数のコードを生成します
// 要件: 3.4, 5.1, 5.2, 5.3, 5.4, 5.5
func (t *TemplateEngine) RenderHelpers(data *TemplateData) (string, error) {
	tmpl, exists := t.templates["helpers"]
	if !exists {
		return "", fmt.Errorf("ヘルパーテンプレートが見つかりません")
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("テンプレートの実行に失敗: %w", err)
	}

	return buf.String(), nil
}

// RenderStructs は構造体定義のコードを生成します
// 要件: 3.1, 3.2, 3.3, 3.4, 3.5
func (t *TemplateEngine) RenderStructs(data *TemplateData) (string, error) {
	tmpl, exists := t.templates["structs"]
	if !exists {
		return "", fmt.Errorf("構造体テンプレートが見つかりません")
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("テンプレートの実行に失敗: %w", err)
	}

	return buf.String(), nil
}

// capitalize は文字列の先頭を大文字にします
func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// singularize は複数形を単数形に変換します（簡易版）
func singularize(s string) string {
	if len(s) > 0 && s[len(s)-1] == 's' {
		return s[:len(s)-1]
	}
	return s
}

// structFieldName はグループ名から有効なGo識別子を生成します
func structFieldName(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return "Group"
	}

	parts := strings.FieldsFunc(name, func(r rune) bool {
		return !(r >= 'a' && r <= 'z') && !(r >= 'A' && r <= 'Z') && !(r >= '0' && r <= '9')
	})

	var builder strings.Builder
	for _, part := range parts {
		if part == "" {
			continue
		}
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		builder.WriteString(strings.ToUpper(part[:1]))
		if len(part) > 1 {
			builder.WriteString(strings.ToLower(part[1:]))
		}
	}

	result := builder.String()
	if result == "" {
		result = "Group"
	}
	if result[0] >= '0' && result[0] <= '9' {
		result = "Group" + result
	}

	return result
}

// sanitizeJSONField はJSONタグに使用できない文字をエスケープします
func sanitizeJSONField(name string) string {
	return strings.ReplaceAll(name, "`", "")
}

// hasCreateParams はCreateParamsが存在するかチェック
func hasCreateParams(group *ParamGroup) bool {
	return group != nil && len(group.CreateParams) > 0
}

// hasUpdateParams はUpdateParamsが存在するかチェック
func hasUpdateParams(group *ParamGroup) bool {
	return group != nil && len(group.UpdateParams) > 0
}

// structsTemplateStr は構造体定義のテンプレート文字列
// 要件: 3.1, 3.2, 3.3, 3.4, 3.5
const structsTemplateStr = `// TestParams はテストパラメータファイルの構造体
type TestParams struct {
{{range $key, $group := .ParamGroups}}
	{{structFieldName $key}} struct {
{{if hasCreateParams $group}}
		CreateParams map[string]any ` + "`json:\"createParams\"`" + `
{{end}}
{{if hasUpdateParams $group}}
		UpdateParams map[string]any ` + "`json:\"updateParams\"`" + `
{{end}}
	} ` + "`json:\"{{sanitizeJSONField $key}}\"`" + `
{{end}}
}
`

// helpersTemplateStr はヘルパー関数のテンプレート文字列
// 要件: 3.4, 5.1, 5.2, 5.3, 5.4, 5.5
const helpersTemplateStr = `package {{.PackageName}}

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

{{.StructsCode}}

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
`
