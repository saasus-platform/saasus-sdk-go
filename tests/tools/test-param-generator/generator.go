package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Generator はコード生成を担当します
type Generator struct {
	Config          *Config
	Analysis        *AnalysisResult
	IterationLoader *IterationDataLoader
	Report          *PostmanReport // レポート情報を保持（要件11.4）
}

// NewGenerator は新しいGeneratorを作成します
func NewGenerator(config *Config, analysis *AnalysisResult) *Generator {
	return &Generator{
		Config:          config,
		Analysis:        analysis,
		IterationLoader: NewIterationDataLoader(analysis.ModuleName, config.Verbose),
	}
}

// NewGeneratorWithReport はレポート情報を含む新しいGeneratorを作成します
// 要件: 11.4
func NewGeneratorWithReport(config *Config, analysis *AnalysisResult, report *PostmanReport) *Generator {
	return &Generator{
		Config:          config,
		Analysis:        analysis,
		IterationLoader: NewIterationDataLoader(analysis.ModuleName, config.Verbose),
		Report:          report,
	}
}

// GenerateParamsFile はパラメータJSONファイルを生成します
// 要件: 2.1, 2.2, 7.3, 8.3
func (g *Generator) GenerateParamsFile() error {
	// 出力ファイルパスを取得
	paramsPath, _ := g.Config.GetOutputPaths()

	// 出力ディレクトリが存在することを確認
	if err := g.Config.EnsureOutputDir(); err != nil {
		return WrapError(err, "出力ディレクトリの作成に失敗")
	}

	// 既存ファイルの上書き確認
	if !g.Config.ShouldOverwriteFile(paramsPath) {
		return NewInvalidConfigError("overwrite", fmt.Sprintf("ファイルが既に存在します。-overwrite フラグを使用してください: %s", paramsPath))
	}

	// パラメータグループからJSON構造を構築
	paramsData := g.buildParamsStructure()

	// パラメータグループが空の場合は警告
	if len(paramsData) == 0 {
		fmt.Fprintf(os.Stderr, "警告: パラメータグループが空です。レポートにテストフローが含まれていない可能性があります。\n")
	}

	// JSONを整形（インデント付き）
	jsonData, err := formatJSON(paramsData)
	if err != nil {
		return WrapError(err, "JSONの整形に失敗")
	}

	// ファイルに書き込み（パーミッション: 0644）
	if err := os.WriteFile(paramsPath, jsonData, 0644); err != nil {
		return NewFileWriteError(paramsPath, err)
	}

	// 成功メッセージを表示
	fmt.Printf("✓ パラメータファイルを生成しました: %s\n", paramsPath)

	// 詳細ログ出力
	if g.Config.Verbose {
		fmt.Printf("  - パラメータグループ数: %d\n", len(g.Analysis.ParamGroups))
		for groupName, group := range g.Analysis.ParamGroups {
			fmt.Printf("  - グループ '%s':\n", groupName)
			if len(group.CreateParams) > 0 {
				fmt.Printf("    - CreateParams: %d フィールド\n", len(group.CreateParams))
			}
			if len(group.UpdateParams) > 0 {
				fmt.Printf("    - UpdateParams: %d フィールド\n", len(group.UpdateParams))
			}
		}
	}

	return nil
}

// buildParamsStructure はパラメータグループからJSON構造を構築します
// 要件: 2.1, 2.2, 11.4
func (g *Generator) buildParamsStructure() map[string]any {
	result := make(map[string]any)

	// ClientMethodNameマッピングを構築（要件11.4）
	methodNameToGroup := g.buildMethodNameMapping()

	// ClientMethodNameが利用可能な場合はそれをキーとして使用（要件11.4）
	if len(methodNameToGroup) > 0 {
		for methodName, groupData := range methodNameToGroup {
			if len(groupData) > 0 {
				result[methodName] = groupData
			}
		}

		// 従来のリソース名キーも保持して後方互換性を確保
		for groupName, group := range g.Analysis.ParamGroups {
			if _, exists := result[groupName]; exists {
				continue
			}
			fallback := make(map[string]any)
			if len(group.CreateParams) > 0 {
				fallback["createParams"] = cloneParamMap(group.CreateParams)
			}
			if len(group.UpdateParams) > 0 {
				fallback["updateParams"] = cloneParamMap(group.UpdateParams)
			}
			if len(fallback) > 0 {
				result[groupName] = fallback
			}
		}

		return result
	}

	// フォールバック: 従来のグループ名をキーとして使用
	for groupName, group := range g.Analysis.ParamGroups {
		groupData := make(map[string]any)

		// CreateParamsを追加（存在する場合）
		if len(group.CreateParams) > 0 {
			groupData["createParams"] = group.CreateParams
		}

		// UpdateParamsを追加（存在する場合）
		if len(group.UpdateParams) > 0 {
			groupData["updateParams"] = group.UpdateParams
		}

		// グループデータが空でない場合のみ追加
		if len(groupData) > 0 {
			result[groupName] = groupData
		}
	}

	return result
}

// buildMethodNameMapping はClientMethodNameをキーとしたマッピングを構築します
// 要件: 11.4
func (g *Generator) buildMethodNameMapping() map[string]map[string]any {
	if g.Report == nil || len(g.Report.TestFlow) == 0 {
		return nil
	}

	result := make(map[string]map[string]any)

	for _, step := range g.Report.TestFlow {
		// ClientMethodNameが利用可能かチェック
		if step.ClientMethodName == "" {
			continue
		}

		// メソッド名に対応するデータを取得または作成
		methodData, exists := result[step.ClientMethodName]
		if !exists {
			methodData = make(map[string]any)
			result[step.ClientMethodName] = methodData
		}

		paramGroup := g.findParamGroupForStep(step)

		// HTTPメソッドに基づいてパラメータを分類
		switch step.Method {
		case "POST":
			// iterationDataなどでマージ済みのパラメータを優先
			if paramGroup != nil && len(paramGroup.CreateParams) > 0 {
				methodData["createParams"] = cloneParamMap(paramGroup.CreateParams)
				continue
			}

			// 型情報ベースのフォールバック
			if len(step.ParameterTypeFields) > 0 {
				createParams := make(map[string]any)
				for jsonTag, fieldInfo := range step.ParameterTypeFields {
					sampleValue := generateSampleValueFromFieldInfo(fieldInfo)
					createParams[jsonTag] = sampleValue
				}
				methodData["createParams"] = createParams
			} else if len(step.RequestBody) > 0 {
				// フォールバック: RequestBodyを使用
				methodData["createParams"] = step.RequestBody
			}

		case "PATCH", "PUT":
			if paramGroup != nil && len(paramGroup.UpdateParams) > 0 {
				methodData["updateParams"] = cloneParamMap(paramGroup.UpdateParams)
				continue
			}

			if len(step.ParameterTypeFields) > 0 {
				updateParams := make(map[string]any)
				for jsonTag, fieldInfo := range step.ParameterTypeFields {
					sampleValue := generateSampleValueFromFieldInfo(fieldInfo)
					updateParams[jsonTag] = sampleValue
				}
				methodData["updateParams"] = updateParams
			} else if len(step.RequestBody) > 0 {
				methodData["updateParams"] = step.RequestBody
			}
		}
	}

	return result
}

// findParamGroupForStep はテストステップに対応するパラメータグループを探索します
func (g *Generator) findParamGroupForStep(step TestStep) *ParamGroup {
	if g.Analysis == nil || len(g.Analysis.ParamGroups) == 0 {
		return nil
	}

	candidates := []string{}
	if step.Endpoint != "" {
		if resource := extractResourceName(step.Endpoint); resource != "" {
			candidates = append(candidates, resource)
		}
	}

	if step.ClientMethodName != "" {
		candidates = append(candidates, step.ClientMethodName)
		if inferred := inferResourceNameFromMethodName(step.ClientMethodName); inferred != "" {
			candidates = append(candidates, inferred)
		}
	}

	if step.ParameterTypeName != "" {
		candidates = append(candidates, step.ParameterTypeName)
	}

	for _, candidate := range candidates {
		if group := g.lookupParamGroupByName(candidate); group != nil {
			return group
		}
	}

	return nil
}

// lookupParamGroupByName は候補名を正規化してパラメータグループを検索します
func (g *Generator) lookupParamGroupByName(candidate string) *ParamGroup {
	if candidate == "" || g.Analysis == nil {
		return nil
	}

	normalizedCandidate := normalizeGroupKey(candidate, g.Analysis.ModuleName)
	for groupName, group := range g.Analysis.ParamGroups {
		if normalizeGroupKey(groupName, g.Analysis.ModuleName) == normalizedCandidate {
			return group
		}
	}

	return nil
}

// normalizeGroupKey は比較用にグループ名を正規化します
func normalizeGroupKey(name, module string) string {
	if name == "" {
		return ""
	}

	key := strings.ToLower(name)
	replacements := []string{"-", "_", "/", " ", "{", "}", ":"}
	for _, r := range replacements {
		key = strings.ReplaceAll(key, r, "")
	}

	if module != "" {
		moduleKey := strings.ToLower(module)
		if strings.HasPrefix(key, moduleKey) {
			key = strings.TrimPrefix(key, moduleKey)
		}
	}

	switch {
	case strings.HasSuffix(key, "ies"):
		key = strings.TrimSuffix(key, "ies") + "y"
	case strings.HasSuffix(key, "ses"), strings.HasSuffix(key, "xes"), strings.HasSuffix(key, "zes"),
		strings.HasSuffix(key, "ches"), strings.HasSuffix(key, "shes"):
		key = strings.TrimSuffix(key, "es")
	case strings.HasSuffix(key, "s") && !strings.HasSuffix(key, "ss"):
		key = strings.TrimSuffix(key, "s")
	}

	return key
}

// cloneParamMap はパラメータマップの浅いコピーを返します
func cloneParamMap(src map[string]any) map[string]any {
	if len(src) == 0 {
		return nil
	}
	clone := make(map[string]any, len(src))
	for k, v := range src {
		clone[k] = v
	}
	return clone
}

// generateSampleValueFromFieldInfo はFieldInfoからサンプル値を生成します
// 要件: 11.4
func generateSampleValueFromFieldInfo(fieldInfo FieldInfo) any {
	// 既存のgenerateSampleValue関数を使用
	return generateSampleValue(fieldInfo.Type, fieldInfo.JSONTag)
}

// formatJSON はJSONを整形します
func formatJSON(data any) ([]byte, error) {
	// インデント付きでJSON化
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return nil, WrapError(err, "JSONのマーシャルに失敗")
	}

	return jsonData, nil
}

// sanitizePackageName はパッケージ名を有効なGo識別子に変換します
// ハイフンをアンダースコアに変換し、数字で始まる場合は"pkg_"を前置します
func sanitizePackageName(name string) string {
	if name == "" {
		return "main"
	}

	// strings.Builderを使用してより効率的に文字列を構築
	var result strings.Builder
	result.Grow(len(name)) // 事前にメモリを確保

	for _, ch := range name {
		switch {
		case ch == '-':
			result.WriteRune('_')
		case (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') || ch == '_':
			result.WriteRune(ch)
			// 無効な文字はスキップ
		}
	}

	sanitized := result.String()

	// 空の場合はデフォルト値を返す
	if sanitized == "" {
		return "main"
	}

	// 数字で始まる場合は"pkg_"を前置
	if sanitized[0] >= '0' && sanitized[0] <= '9' {
		return "pkg_" + sanitized
	}

	return sanitized
}

// GenerateHelpersFile はGoヘルパー関数ファイルを生成します
// 要件: 3.1, 3.2, 3.5, 7.4, 8.3, 8.4
func (g *Generator) GenerateHelpersFile() error {
	// 出力ファイルパスを取得
	_, helpersPath := g.Config.GetOutputPaths()

	// 既存ファイルの上書き確認
	if !g.Config.ShouldOverwriteFile(helpersPath) {
		return NewInvalidConfigError("overwrite", fmt.Sprintf("ファイルが既に存在します。-overwrite フラグを使用してください: %s", helpersPath))
	}

	// テンプレートエンジンを作成
	templateEngine := NewTemplateEngine()

	// テンプレートデータを準備
	// パッケージ名を生成（ハイフンをアンダースコアに変換）
	packageName := filepath.Base(g.Config.OutputDir)
	packageName = sanitizePackageName(packageName)

	templateData := &TemplateData{
		PackageName: packageName,
		ModuleName:  g.Analysis.ModuleName,
		ParamGroups: g.Analysis.ParamGroups,
		ImportPaths: []string{},
	}

	// 構造体定義のコードを生成
	structsCode, err := templateEngine.RenderStructs(templateData)
	if err != nil {
		return NewTemplateRenderError("structs", err)
	}
	templateData.StructsCode = structsCode

	// ヘルパー関数のコードを生成
	helpersCode, err := templateEngine.RenderHelpers(templateData)
	if err != nil {
		return NewTemplateRenderError("helpers", err)
	}

	// ファイルに書き込み
	if err := os.WriteFile(helpersPath, []byte(helpersCode), 0644); err != nil {
		return NewFileWriteError(helpersPath, err)
	}

	// 詳細ログ出力
	if g.Config.Verbose {
		fmt.Printf("ヘルパー関数ファイルを生成しました: %s\n", helpersPath)
	}

	return nil
}

// mergeIterationData はiterationDataをパラメータグループにマージします
// 要件: 9.5, 9.6
func (g *Generator) mergeIterationData() error {
	// iterationDataを読み込む
	iterationData, err := g.IterationLoader.LoadIterationData()
	if err != nil {
		// エラーが発生した場合は警告を表示してフォールバック
		if g.Config.Verbose {
			fmt.Fprintf(os.Stderr, "警告: iterationDataの読み込みに失敗しました: %v\n", err)
			fmt.Fprintf(os.Stderr, "自動生成値を使用します\n")
		}
		return nil // エラーを返さずに処理を継続
	}

	// iterationDataが空の場合は何もしない
	if len(iterationData) == 0 {
		if g.Config.Verbose {
			fmt.Printf("iterationDataが見つかりませんでした。自動生成値を使用します\n")
		}
		return nil
	}

	if g.Config.Verbose {
		fmt.Printf("iterationDataをパラメータグループにマージ中...\n")
	}

	// 各パラメータグループにiterationDataをマッチング
	for groupName, paramGroup := range g.Analysis.ParamGroups {
		if err := matchIterationDataToParamGroup(iterationData, paramGroup, groupName); err != nil {
			if g.Config.Verbose {
				fmt.Fprintf(os.Stderr, "警告: グループ '%s' へのマッチングに失敗: %v\n", groupName, err)
			}
			continue
		}

		if g.Config.Verbose {
			fmt.Printf("  - グループ '%s' にiterationDataをマージしました\n", groupName)
		}
	}

	return nil
}
