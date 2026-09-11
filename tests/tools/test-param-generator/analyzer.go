package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

var (
	// globalTypeCache はグローバルな型情報キャッシュです
	globalTypeCache *TypeCache
)

// initGlobalCache はグローバルキャッシュを初期化します
func initGlobalCache() {
	if globalTypeCache == nil {
		// デフォルト設定: 最大100エントリ、TTL 1時間、メモリ制限 100MB
		globalTypeCache = NewTypeCache(100, 3600, 100)
	}
}

// GetGlobalCache はグローバルキャッシュを取得します
func GetGlobalCache() *TypeCache {
	initGlobalCache()
	return globalTypeCache
}

// shouldForceMappingDiagnostics は診断用マッピングの強制実行が有効かを判断します
func shouldForceMappingDiagnostics() bool {
	envVars := []string{
		"TPG_FORCE_MAPPING_DIAGNOSTICS",
		"TEST_PARAM_GENERATOR_FORCE_DIAGNOSTICS",
	}

	for _, key := range envVars {
		value := strings.TrimSpace(strings.ToLower(os.Getenv(key)))
		if value == "" {
			continue
		}

		switch value {
		case "1", "true", "yes", "on", "enable", "enabled":
			return true
		}
	}

	return false
}

func appendUniqueModule(modules []string, module string) []string {
	module = strings.TrimSpace(module)
	if module == "" {
		return modules
	}
	for _, m := range modules {
		if m == module {
			return modules
		}
	}
	return append(modules, module)
}

func collectModuleCandidates(report *PostmanReport, defaultModule string, groups map[string]*ParamGroup) []string {
	seen := make(map[string]struct{})
	ordered := make([]string, 0)

	add := func(module string) {
		module = strings.TrimSpace(module)
		if module == "" {
			return
		}
		if _, exists := seen[module]; exists {
			return
		}
		seen[module] = struct{}{}
		ordered = append(ordered, module)
	}

	add(defaultModule)

	if report != nil {
		for _, step := range report.TestFlow {
			add(step.Module)
		}
		for _, method := range report.ClientMethods {
			add(method.Module)
		}
	}

	for _, group := range groups {
		for _, module := range group.ModuleCandidates {
			add(module)
		}
	}

	return ordered
}

func buildModuleLoadOrder(primary string, modules []string) []string {
	order := make([]string, 0, len(modules)+1)
	order = appendUniqueModule(order, primary)
	for _, module := range modules {
		order = appendUniqueModule(order, module)
	}
	return order
}

func mergeGroupModuleOrder(groupModules, globalOrder []string) []string {
	merged := make([]string, 0, len(groupModules)+len(globalOrder))
	for _, module := range groupModules {
		merged = appendUniqueModule(merged, module)
	}
	for _, module := range globalOrder {
		merged = appendUniqueModule(merged, module)
	}
	return merged
}

// extractModuleNameFromPath はファイルパスからモジュール名を抽出します
func extractModuleNameFromPath(filePath string) string {
	// パスを正規化
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		absPath = filePath
	}

	// ディレクトリ名を取得（例: generated/billingapi）
	dir := filepath.Dir(absPath)
	moduleName := filepath.Base(dir)

	// "api"サフィックスを除去（例: billingapi -> billing）
	moduleName = strings.TrimSuffix(moduleName, "api")

	return moduleName
}

// isDebugMode はデバッグモードが有効かどうかを返します
func isDebugMode() bool {
	return os.Getenv("DEBUG") == "1" || os.Getenv("DEBUG") == "true"
}

// debugLog はデバッグモードの場合のみログを出力します
func debugLog(format string, args ...interface{}) {
	if isDebugMode() {
		fmt.Fprintf(os.Stderr, "[DEBUG] "+format, args...)
	}
}

// parseGeneratedTypesEnhanced は生成された型定義ファイルを拡張解析します
func parseGeneratedTypesEnhanced(typesFilePath string) (map[string]EnhancedTypeInfo, error) {
	// ファイルの存在確認
	if _, err := os.Stat(typesFilePath); os.IsNotExist(err) {
		return nil, NewFileNotFoundError(typesFilePath, err)
	}

	// モジュール名をキャッシュキーとして使用
	moduleName := extractModuleNameFromPath(typesFilePath)

	// キャッシュから取得を試みる
	cache := GetGlobalCache()
	if cached, found := cache.Get(moduleName); found {
		debugLog("キャッシュから型情報を取得: %s (%d 型)\n", moduleName, len(cached))
		return cached, nil
	}

	// ファイルセットを作成
	fset := token.NewFileSet()

	// Goファイルを解析
	node, err := parser.ParseFile(fset, typesFilePath, nil, parser.ParseComments)
	if err != nil {
		return nil, NewTypeAnalysisError(typesFilePath, err)
	}

	// 型情報を格納するマップ
	typeInfoMap := make(map[string]EnhancedTypeInfo)
	// 型エイリアスを格納するマップ
	typeAliases := make(map[string]string)

	// 第1パス: 構造体定義と型エイリアスを収集
	ast.Inspect(node, func(n ast.Node) bool {
		typeSpec, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}

		typeName := typeSpec.Name.Name

		// 構造体定義を処理
		if structType, ok := typeSpec.Type.(*ast.StructType); ok {
			// EnhancedTypeInfoを作成
			typeInfo := EnhancedTypeInfo{
				Name:   typeName,
				Fields: make(map[string]EnhancedFieldInfo),
			}

			// 型名からメタデータを推測
			metadata := analyzeTypeMetadata(typeName)
			typeInfo.IsRequestBody = metadata.TypeCategory == "RequestBody"
			typeInfo.IsParams = metadata.TypeCategory == "Params"
			typeInfo.HTTPMethod = metadata.HTTPMethod
			typeInfo.ResourceName = metadata.ResourceName

			// フィールドを解析
			for _, field := range structType.Fields.List {
				fieldInfo := extractEnhancedFieldInfo(field)
				if fieldInfo.Name != "" {
					typeInfo.Fields[fieldInfo.Name] = fieldInfo
				}
			}

			typeInfoMap[typeName] = typeInfo
		} else if ident, ok := typeSpec.Type.(*ast.Ident); ok {
			// 型エイリアス
			typeAliases[typeName] = ident.Name
			debugLog("型エイリアス発見: %s -> %s\n", typeName, ident.Name)
		}

		return true
	})

	// 第2パス: 型エイリアスを解決
	var resolveAlias func(aliasName string, visited map[string]struct{}) (EnhancedTypeInfo, string, bool)
	resolveAlias = func(aliasName string, visited map[string]struct{}) (EnhancedTypeInfo, string, bool) {
		if info, exists := typeInfoMap[aliasName]; exists {
			return info, aliasName, true
		}

		targetName, exists := typeAliases[aliasName]
		if !exists {
			return EnhancedTypeInfo{}, "", false
		}

		if _, seen := visited[aliasName]; seen {
			return EnhancedTypeInfo{}, "", false
		}
		visited[aliasName] = struct{}{}

		info, resolvedTarget, ok := resolveAlias(targetName, visited)
		if !ok {
			return EnhancedTypeInfo{}, "", false
		}

		typeInfoMap[aliasName] = info
		return info, resolvedTarget, true
	}

	for aliasName := range typeAliases {
		if _, exists := typeInfoMap[aliasName]; exists {
			continue
		}

		info, resolvedTarget, ok := resolveAlias(aliasName, make(map[string]struct{}))
		if ok {
			debugLog("型エイリアス解決: %s -> %s (%d フィールド)\n",
				aliasName, resolvedTarget, len(info.Fields))
		} else {
			debugLog("型エイリアス解決失敗: %s -> %s\n", aliasName, typeAliases[aliasName])
		}
	}

	if len(typeInfoMap) == 0 {
		return nil, NewTypeAnalysisError(typesFilePath, fmt.Errorf("構造体定義が見つかりません"))
	}

	// キャッシュに保存
	if err := cache.Set(moduleName, typeInfoMap); err != nil {
		debugLog("キャッシュへの保存に失敗: %v\n", err)
	} else {
		debugLog("キャッシュに型情報を保存: %s (%d 型)\n", moduleName, len(typeInfoMap))
	}

	return typeInfoMap, nil
}

// parseGeneratedTypes は生成された型定義ファイルを解析します
//
// 将来の機能拡張用:
// - 生成された型定義から正確なフィールド型を取得
// - より精密なサンプル値生成
// - 型安全性の向上
//
// 現在は mapToGeneratedTypes から呼び出されていますが、
// エラーは警告として扱われ、処理は継続されます。
func parseGeneratedTypes(typesFilePath string) (map[string]TypeInfo, error) {
	// ファイルの存在確認
	if _, err := os.Stat(typesFilePath); os.IsNotExist(err) {
		return nil, NewFileNotFoundError(typesFilePath, err)
	}

	// ファイルセットを作成
	fset := token.NewFileSet()

	// Goファイルを解析
	node, err := parser.ParseFile(fset, typesFilePath, nil, parser.ParseComments)
	if err != nil {
		return nil, NewTypeAnalysisError(typesFilePath, err)
	}

	// 型情報を格納するマップ
	typeInfoMap := make(map[string]TypeInfo)
	// 型エイリアスを格納するマップ (エイリアス名 -> 実際の型名)
	typeAliases := make(map[string]string)

	// 第1パス: 構造体定義と型エイリアスを収集
	ast.Inspect(node, func(n ast.Node) bool {
		// 型定義を探す
		typeSpec, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}

		typeName := typeSpec.Name.Name

		// 構造体定義を処理
		if structType, ok := typeSpec.Type.(*ast.StructType); ok {
			// TypeInfoを作成
			typeInfo := TypeInfo{
				Name:   typeName,
				Fields: make(map[string]FieldInfo),
			}

			// フィールドを解析
			for _, field := range structType.Fields.List {
				fieldInfo := extractFieldInfo(field)
				if fieldInfo.Name != "" {
					typeInfo.Fields[fieldInfo.Name] = fieldInfo
				}
			}

			// マップに追加
			typeInfoMap[typeName] = typeInfo
		} else if ident, ok := typeSpec.Type.(*ast.Ident); ok {
			// 型エイリアス (type A = B の形式)
			typeAliases[typeName] = ident.Name
			debugLog("型エイリアス発見: %s -> %s\n", typeName, ident.Name)
		}

		return true
	})

	// 第2パス: 型エイリアスを解決
	for aliasName, targetName := range typeAliases {
		if targetInfo, exists := typeInfoMap[targetName]; exists {
			// エイリアスに対して実際の型情報をコピー
			typeInfoMap[aliasName] = targetInfo
			debugLog("型エイリアス解決: %s -> %s (%d フィールド)\n",
				aliasName, targetName, len(targetInfo.Fields))
		}
	}

	if len(typeInfoMap) == 0 {
		return nil, NewTypeAnalysisError(typesFilePath, fmt.Errorf("構造体定義が見つかりません"))
	}

	return typeInfoMap, nil
}

// extractEnhancedFieldInfo は拡張されたフィールド情報を抽出します
func extractEnhancedFieldInfo(field *ast.Field) EnhancedFieldInfo {
	fieldInfo := EnhancedFieldInfo{}

	// フィールド名を取得
	if len(field.Names) > 0 {
		fieldInfo.Name = field.Names[0].Name
	}

	// フィールド型を取得
	fieldInfo.Type = extractTypeName(field.Type)

	// ポインタ型かどうかを判定
	if _, ok := field.Type.(*ast.StarExpr); ok {
		fieldInfo.IsPointer = true
	}

	// 配列型かどうかを判定
	if _, ok := field.Type.(*ast.ArrayType); ok {
		fieldInfo.IsArray = true
	}

	// ネストした構造体かどうかを判定
	if ident, ok := field.Type.(*ast.Ident); ok {
		// 大文字で始まる型名はカスタム型（ネストの可能性）
		if len(ident.Name) > 0 && ident.Name[0] >= 'A' && ident.Name[0] <= 'Z' {
			fieldInfo.IsNested = true
			fieldInfo.NestedType = ident.Name
		}
	}

	// JSONタグを取得
	if field.Tag != nil {
		tag := field.Tag.Value
		jsonTag := extractJSONTag(tag)
		fieldInfo.JSONTag = jsonTag

		// requiredフラグを判定
		fieldInfo.Required = !strings.Contains(jsonTag, "omitempty")
	}

	// フィールド名の別名を生成
	if fieldInfo.JSONTag != "" {
		normalizer := &FieldNameNormalizer{}
		fieldInfo.Aliases = normalizer.generateAliases(fieldInfo.JSONTag)
	}

	return fieldInfo
}

// extractFieldInfo はフィールド情報を抽出します
func extractFieldInfo(field *ast.Field) FieldInfo {
	fieldInfo := FieldInfo{}

	// フィールド名を取得
	if len(field.Names) > 0 {
		fieldInfo.Name = field.Names[0].Name
	}

	// フィールド型を取得
	fieldInfo.Type = extractTypeName(field.Type)

	// JSONタグを取得
	if field.Tag != nil {
		tag := field.Tag.Value
		// タグから`json:"..."`部分を抽出
		jsonTag := extractJSONTag(tag)
		fieldInfo.JSONTag = jsonTag

		// requiredフラグを判定（omitemptyがない場合はrequired）
		fieldInfo.Required = !strings.Contains(jsonTag, "omitempty")
	}

	return fieldInfo
}

// extractTypeName は型名を文字列として抽出します
func extractTypeName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		// 基本型（string, int, bool等）
		return t.Name
	case *ast.SelectorExpr:
		// パッケージ修飾型（time.Time等）
		if x, ok := t.X.(*ast.Ident); ok {
			return x.Name + "." + t.Sel.Name
		}
		return t.Sel.Name
	case *ast.StarExpr:
		// ポインタ型
		return "*" + extractTypeName(t.X)
	case *ast.ArrayType:
		// 配列/スライス型
		return "[]" + extractTypeName(t.Elt)
	case *ast.MapType:
		// マップ型
		keyType := extractTypeName(t.Key)
		valueType := extractTypeName(t.Value)
		return fmt.Sprintf("map[%s]%s", keyType, valueType)
	case *ast.InterfaceType:
		// インターフェース型
		return "interface{}"
	default:
		return "unknown"
	}
}

// extractJSONTag はタグ文字列からJSON部分を抽出します
func extractJSONTag(tag string) string {
	// タグ文字列は `json:"field_name"` または `json:"field_name,omitempty"` の形式
	tag = strings.Trim(tag, "`")

	// json:"..." 部分を探す
	parts := strings.Split(tag, " ")
	for _, part := range parts {
		// CutPrefixを使用してより効率的に処理
		if jsonPart, found := strings.CutPrefix(part, "json:"); found {
			// json:"..." から引用符を除去
			jsonPart = strings.Trim(jsonPart, `"`)
			return jsonPart
		}
	}

	return ""
}

// getTypesFilePath はModule情報からtypes.gen.goファイルのパスを構築します
// Module情報が存在する場合は優先的に使用し、存在しない場合はフォールバックします
// 要件: 2.1, 2.2, 6.1, 6.2
func getTypesFilePath(moduleName string) (string, error) {
	if moduleName == "" {
		return "", fmt.Errorf("モジュール名が空です")
	}

	// Module情報を使用してパスを構築
	// 形式: generated/{module}api/types.gen.go
	primaryPath := filepath.Join("generated", moduleName+"api", "types.gen.go")

	debugLog("Module情報からパスを構築: %s\n", primaryPath)

	// パスの存在確認
	if _, err := os.Stat(primaryPath); err == nil {
		absPath, _ := filepath.Abs(primaryPath)
		debugLog("✓ Module情報からファイルを発見: %s\n", absPath)
		return absPath, nil
	}

	// フォールバック: 従来のパス探索ロジック
	debugLog("Module情報からのパスが見つからない、フォールバック探索を開始\n")
	return findGeneratedTypesFile(moduleName)
}

// findGeneratedTypesFile はモジュール名から生成された型定義ファイルのパスを探します
// これは getTypesFilePath のフォールバックとして使用されます
func findGeneratedTypesFile(moduleName string) (string, error) {
	// 生成ファイルのパスパターン
	patterns := []string{
		filepath.Join("generated", moduleName+"api", "types.gen.go"),
		filepath.Join("generated", moduleName, "types.gen.go"),
		filepath.Join("..", "..", "..", "generated", moduleName+"api", "types.gen.go"),
		filepath.Join("..", "..", "..", "generated", moduleName, "types.gen.go"),
	}

	debugLog("フォールバック: 以下のパターンを試行:\n")
	for _, pattern := range patterns {
		debugLog("  - %s\n", pattern)
		if _, err := os.Stat(pattern); err == nil {
			absPath, _ := filepath.Abs(pattern)
			debugLog("✓ フォールバックでファイルを発見: %s\n", absPath)
			return absPath, nil
		}
	}

	return "", NewFileNotFoundError(
		fmt.Sprintf("generated/%sapi/types.gen.go", moduleName),
		fmt.Errorf("モジュール '%s' の型定義ファイルが見つかりません", moduleName),
	)
}

// mapToGeneratedTypes は生成された型定義ファイルを読み込んでパラメータ構造をマッピングします
// 要件: 1.4, 1.5, 2.1, 5.1, 5.2, 5.3, 5.4, 5.5
//
// エラーハンドリング戦略:
// - 型定義ファイルが見つからない場合: 警告を出力して継続（デフォルト値を使用）
// - 型解析に失敗した場合: 警告を出力して継続（デフォルト値を使用）
// - 個別のマッピングに失敗した場合: 警告を出力して継続（既存の値を保持）
//
// カスタムマッピングルール:
// - mapping_rules.jsonファイルが存在する場合、カスタムルールを最優先で適用
// - ルールファイルのパス: ./mapping_rules.json または ../mapping_rules.json
func mapToGeneratedTypes(moduleName string, paramGroups map[string]*ParamGroup, moduleCandidates []string) error {
	// DiagnosticReporterの初期化
	reporter := NewDiagnosticReporter(moduleName, isDebugMode())

	// デバッグ: マッピング開始
	debugLog("mapToGeneratedTypes 開始: moduleName=%s, paramGroups数=%d\n", moduleName, len(paramGroups))

	moduleOrder := buildModuleLoadOrder(moduleName, moduleCandidates)
	if isDebugMode() {
		debugLog("参照モジュール候補: %v\n", moduleOrder)
	}

	// カスタムマッピングルールを読み込む
	mappingRules, err := loadMappingRulesIfExists()
	if err != nil {
		fmt.Fprintf(os.Stderr, "警告: マッピングルールの読み込みに失敗: %v\n", err)
		mappingRules = NewMappingRules()
	} else if mappingRules != nil {
		debugLog("カスタムマッピングルールを読み込みました\n")
	}

	// モジュールごとに型定義を読み込み結合
	enhancedTypeInfoMap := make(map[string]EnhancedTypeInfo)
	loadedModules := make([]string, 0, len(moduleOrder))

	for _, candidateModule := range moduleOrder {
		typesFilePath, err := getTypesFilePath(candidateModule)
		if err != nil {
			if candidateModule == moduleName {
				if fileNotFoundErr, ok := err.(*FileNotFoundError); ok {
					fmt.Fprintf(os.Stderr, "警告: %s\n", fileNotFoundErr.Error())
					fmt.Fprintf(os.Stderr, "  → フォールバック: デフォルト値を使用してテストパラメータを生成します\n")
				} else {
					fmt.Fprintf(os.Stderr, "警告: %v\n", err)
					fmt.Fprintf(os.Stderr, "  → フォールバック: デフォルト値を使用してテストパラメータを生成します\n")
				}
			} else if isDebugMode() {
				debugLog("追加モジュールの型定義探索に失敗 (%s): %v\n", candidateModule, err)
			}
			continue
		}

		debugLog("型定義ファイル発見: %s (module=%s)\n", typesFilePath, candidateModule)

		moduleTypeInfo, err := parseGeneratedTypesEnhanced(typesFilePath)
		if err != nil {
			if candidateModule == moduleName {
				if typeAnalysisErr, ok := err.(*TypeAnalysisError); ok {
					fmt.Fprintf(os.Stderr, "警告: %s\n", typeAnalysisErr.Error())
					fmt.Fprintf(os.Stderr, "  → フォールバック: デフォルト値を使用してテストパラメータを生成します\n")
					if isDebugMode() {
						fmt.Fprintf(os.Stderr, "  詳細: %v\n", typeAnalysisErr.Unwrap())
					}
				} else {
					fmt.Fprintf(os.Stderr, "警告: 型定義の解析に失敗: %v\n", err)
					fmt.Fprintf(os.Stderr, "  → フォールバック: デフォルト値を使用してテストパラメータを生成します\n")
				}
			} else if isDebugMode() {
				debugLog("追加モジュールの型定義解析に失敗 (%s): %v\n", candidateModule, err)
			}
			continue
		}

		loadedModules = append(loadedModules, candidateModule)

		for typeName, typeInfo := range moduleTypeInfo {
			if _, exists := enhancedTypeInfoMap[typeName]; exists {
				continue
			}
			enhancedTypeInfoMap[typeName] = typeInfo
		}
	}

	if len(enhancedTypeInfoMap) == 0 {
		debugLog("有効な型定義を読み込めませんでした\n")
		reporter.ReportMappingResults()
		return nil
	}

	debugLog("型定義解析成功: %d個の型を集約 (モジュール: %v)\n", len(enhancedTypeInfoMap), loadedModules)

	// OpenAPIマッパーをモジュールごとに初期化
	openAPIMappers := make(map[string]*OpenAPIMapper)
	for _, candidateModule := range moduleOrder {
		mapper, err := NewOpenAPIMapper(candidateModule)
		if err != nil {
			debugLog("OpenAPIマッパーの初期化失敗 (%s): %v\n", candidateModule, err)
			continue
		}
		openAPIMappers[candidateModule] = mapper
	}

	// デバッグ: 発見された型名をリスト表示
	if isDebugMode() {
		debugLog("発見された型名:\n")
		typesByCategory := make(map[string][]string)
		for typeName, typeInfo := range enhancedTypeInfoMap {
			category := "Other"
			if typeInfo.IsRequestBody {
				category = "RequestBody"
			} else if typeInfo.IsParams {
				category = "Params"
			}
			typesByCategory[category] = append(typesByCategory[category], typeName)
		}

		for category, types := range typesByCategory {
			debugLog("  [%s] (%d個)\n", category, len(types))
			for _, typeName := range types {
				typeInfo := enhancedTypeInfoMap[typeName]
				debugLog("    - %s (フィールド数: %d, HTTPメソッド: %s, リソース: %s)\n",
					typeName, len(typeInfo.Fields), typeInfo.HTTPMethod, typeInfo.ResourceName)
			}
		}
	}

	// 各パラメータグループに対して型情報をマッピング
	for groupName, group := range paramGroups {
		debugLog("\n=== グループ '%s' のマッピング開始 ===\n", groupName)
		reporter.metrics.TotalGroups++
		groupSuccess := false
		moduleOrderForGroup := mergeGroupModuleOrder(group.ModuleCandidates, moduleOrder)

		// CreateParams用の型を探す
		if len(group.CreateParams) > 0 {
			debugLog("  [CreateParams] フィールド数: %d\n", len(group.CreateParams))
			if isDebugMode() {
				debugLog("  フィールド一覧: ")
				fieldNames := make([]string, 0, len(group.CreateParams))
				for fieldName := range group.CreateParams {
					fieldNames = append(fieldNames, fieldName)
				}
				debugLog("%v\n", fieldNames)
			}

			reporter.metrics.TotalFields += len(group.CreateParams)

			// OpenAPIマッパーを優先的に使用
			mappedCount, failed := mapParamsToEnhancedTypeWithOpenAPI(
				group.CreateParams, enhancedTypeInfoMap, mappingRules, groupName, "Create",
				group.CreateEndpoint, group.CreateMethod, moduleOrderForGroup, openAPIMappers, reporter)
			reporter.metrics.MappedFields += mappedCount

			if len(failed) == 0 {
				debugLog("  ✓ CreateParams マッピング成功 (%d/%d フィールド)\n", mappedCount, len(group.CreateParams))
				groupSuccess = true
			} else {
				debugLog("  ✗ CreateParams マッピング部分的に失敗 (%d/%d フィールド成功, %d失敗)\n",
					mappedCount, len(group.CreateParams), len(failed))
			}
		}

		// UpdateParams用の型を探す
		if len(group.UpdateParams) > 0 {
			debugLog("  [UpdateParams] フィールド数: %d\n", len(group.UpdateParams))
			if isDebugMode() {
				debugLog("  フィールド一覧: ")
				fieldNames := make([]string, 0, len(group.UpdateParams))
				for fieldName := range group.UpdateParams {
					fieldNames = append(fieldNames, fieldName)
				}
				debugLog("%v\n", fieldNames)
			}

			reporter.metrics.TotalFields += len(group.UpdateParams)

			// OpenAPIマッパーを優先的に使用
			mappedCount, failed := mapParamsToEnhancedTypeWithOpenAPI(
				group.UpdateParams, enhancedTypeInfoMap, mappingRules, groupName, "Update",
				group.UpdateEndpoint, group.UpdateMethod, moduleOrderForGroup, openAPIMappers, reporter)
			reporter.metrics.MappedFields += mappedCount

			if len(failed) == 0 {
				debugLog("  ✓ UpdateParams マッピング成功 (%d/%d フィールド)\n", mappedCount, len(group.UpdateParams))
				groupSuccess = true
			} else {
				debugLog("  ✗ UpdateParams マッピング部分的に失敗 (%d/%d フィールド成功, %d失敗)\n",
					mappedCount, len(group.UpdateParams), len(failed))
			}
		}

		if groupSuccess {
			reporter.metrics.SuccessfulGroups++
			debugLog("=== グループ '%s' マッピング完了: 成功 ===\n", groupName)
		} else {
			debugLog("=== グループ '%s' マッピング完了: 失敗 ===\n", groupName)
		}
	}

	debugLog("\n=== mapToGeneratedTypes 完了 ===\n")

	// デバッグモードの場合は中間結果を出力
	if isDebugMode() {
		debugLog("\n--- 中間結果サマリー ---\n")
		debugLog("処理したグループ数: %d\n", reporter.metrics.TotalGroups)
		debugLog("成功したグループ数: %d\n", reporter.metrics.SuccessfulGroups)
		debugLog("総フィールド数: %d\n", reporter.metrics.TotalFields)
		debugLog("マッピング成功フィールド数: %d\n", reporter.metrics.MappedFields)

		if len(reporter.metrics.TypePatternStats) > 0 {
			debugLog("\n型パターン使用統計:\n")
			for pattern, count := range reporter.metrics.TypePatternStats {
				debugLog("  %s: %d回\n", pattern, count)
			}
		}

		if len(reporter.metrics.MatchStrategyStats) > 0 {
			debugLog("\nマッチング戦略使用統計:\n")
			for strategy, count := range reporter.metrics.MatchStrategyStats {
				debugLog("  %s: %d回\n", strategy, count)
			}
		}

		debugLog("------------------------\n\n")
	}

	reporter.ReportMappingResults()
	return nil
}

// loadMappingRulesIfExists はマッピングルールファイルが存在する場合に読み込みます
func loadMappingRulesIfExists() (*MappingRules, error) {
	// 複数のパスを試行
	paths := []string{
		"mapping_rules.json",
		"../mapping_rules.json",
		"../../mapping_rules.json",
		"./tests/tools/test-param-generator/mapping_rules.json",
	}

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			debugLog("マッピングルールファイルを発見: %s\n", path)
			return LoadMappingRules(path)
		}
	}

	// ファイルが見つからない場合は空のルールを返す
	return NewMappingRules(), nil
}

// generateTypeNamePatterns は型名パターンを生成します（優先順位順）
func generateTypeNamePatterns(groupName, operation string) []string {
	normalized := normalizeResourceName(groupName)
	lowercaseOp := strings.ToLower(operation)

	// 優先順位の高い順にパターンを定義
	patterns := []string{
		// 1. 標準的なoapi-codegenパターン
		operation + normalized + "JSONRequestBody",
		operation + normalized + "JSONBody",
		operation + normalized + "Body",

		// 2. パラメータ型パターン
		operation + normalized + "Params",
		operation + normalized + "Param",
		operation + normalized + "Request",

		// 3. 逆順パターン
		normalized + operation + "Params",
		normalized + operation + "Param",
		normalized + operation + "Request",
		normalized + operation + "Body",

		// 4. 小文字operationパターン
		normalized + strings.Title(lowercaseOp) + "Params",
		normalized + strings.Title(lowercaseOp) + "Param",

		// 5. 複数形対応
		operation + normalized + "sParams",
		operation + normalized + "sParam",

		// 6. シンプルパターン
		normalized + "Params",
		normalized + "Param",
		normalized + "Body",
	}

	// 重複を除去
	seen := make(map[string]bool)
	unique := make([]string, 0, len(patterns))
	for _, p := range patterns {
		if !seen[p] {
			seen[p] = true
			unique = append(unique, p)
		}
	}

	return unique
}

// mapParamsToEnhancedTypeWithOpenAPI はOpenAPIマッパーを優先して使用します
func mapParamsToEnhancedTypeWithOpenAPI(params map[string]any, typeInfoMap map[string]EnhancedTypeInfo, mappingRules *MappingRules, groupName, operation, endpoint, method string, moduleOrder []string, openAPIMappers map[string]*OpenAPIMapper, reporter *DiagnosticReporter) (mappedCount int, failed []DetailedFailedMapping) {
	var typeInfo EnhancedTypeInfo
	var found bool
	var matchedTypeName string
	var matchSource string

	// 優先度1: OpenAPIマッパーを使用
	if endpoint != "" && method != "" {
		for _, module := range moduleOrder {
			if mapper := openAPIMappers[module]; mapper != nil {
				if typeName, ok := mapper.GetTypeNameForEndpoint(method, endpoint); ok {
					if ti, exists := typeInfoMap[typeName]; exists {
						typeInfo = ti
						found = true
						matchedTypeName = typeName
						matchSource = "OpenAPI"
						reporter.metrics.TypePatternStats["OpenAPI"]++
						debugLog("    OpenAPIマッピング成功: %s %s -> %s (module=%s)\n", method, endpoint, typeName, module)
						break
					}
				}
			}
			if found {
				break
			}
		}
	}

	// 優先度2: カスタムマッピングルール
	if !found && mappingRules != nil {
		for _, module := range moduleOrder {
			if customTypeName, ok := mappingRules.ApplyTypeNameMappingWithModule(module, groupName, operation); ok {
				if ti, exists := typeInfoMap[customTypeName]; exists {
					typeInfo = ti
					found = true
					matchedTypeName = customTypeName
					matchSource = "CustomRule"
					reporter.metrics.TypePatternStats["CustomRule"]++
					break
				}
			}
		}
	}

	// 優先度3: パターンマッチング（フォールバック）
	if !found {
		return mapParamsToEnhancedTypeWithReporter(params, typeInfoMap, mappingRules, groupName, operation, moduleOrder, reporter)
	}

	// マッピング成功 - フィールドをマッピング
	debugLog("    型マッピング成功: %s (ソース: %s)\n", matchedTypeName, matchSource)
	return mapFieldsToType(params, typeInfo, mappingRules, groupName, operation, moduleOrder, reporter)
}

// mapParamsToEnhancedTypeWithReporter は拡張型情報とDiagnosticReporterを使用してパラメータをマッピングします
func mapParamsToEnhancedTypeWithReporter(params map[string]any, typeInfoMap map[string]EnhancedTypeInfo, mappingRules *MappingRules, groupName, operation string, moduleOrder []string, reporter *DiagnosticReporter) (mappedCount int, failed []DetailedFailedMapping) {
	patternGen := NewPatternGenerator(groupName, operation, typeInfoMap)
	patterns := patternGen.GeneratePatterns()

	debugLog("    生成されたパターン数: %d\n", len(patterns))

	var typeInfo EnhancedTypeInfo
	var found bool

	if mappingRules != nil {
		for _, module := range moduleOrder {
			if customTypeName, ok := mappingRules.ApplyTypeNameMappingWithModule(module, groupName, operation); ok {
				debugLog("    カスタムルールで型名マッピング発見: %s.%s.%s -> %s\n", module, groupName, operation, customTypeName)
				if ti, exists := typeInfoMap[customTypeName]; exists {
					typeInfo = ti
					found = true
					reporter.metrics.TypePatternStats["CustomRule"]++
					break
				} else {
					debugLog("    警告: カスタムルールで指定された型名 '%s' が見つかりません (module=%s)\n", customTypeName, module)
				}
			}
		}

		if !found {
			aliases := mappingRules.GetResourceAliases(groupName)
			for _, alias := range aliases {
				for _, module := range moduleOrder {
					if customTypeName, ok := mappingRules.ApplyTypeNameMappingWithModule(module, alias, operation); ok {
						debugLog("    カスタムルール（別名）で型名マッピング発見: %s.%s.%s -> %s\n", module, alias, operation, customTypeName)
						if ti, exists := typeInfoMap[customTypeName]; exists {
							typeInfo = ti
							found = true
							reporter.metrics.TypePatternStats["CustomRule(Alias)"]++
							break
						}
					}
				}
				if found {
					break
				}
			}
		}
	}

	if !found {
		for _, pattern := range patterns {
			if ti, exists := typeInfoMap[pattern.Pattern]; exists {
				typeInfo = ti
				found = true
				reporter.metrics.TypePatternStats[pattern.Category]++
				debugLog("    パターンマッチ成功: %s (カテゴリ: %s, 優先度: %d)\n",
					pattern.Pattern, pattern.Category, pattern.Priority)
				break
			}
		}
	}

	if !found {
		triedPatterns := make([]TypeNamePattern, 0, len(patterns))
		for _, p := range patterns {
			triedPatterns = append(triedPatterns, p)
		}

		warningLevel := "警告"
		if len(patterns) > 10 {
			warningLevel = "エラー"
		}

		fmt.Fprintf(os.Stderr, "%s: グループ '%s' の%s操作用の型定義が見つかりません\n", warningLevel, groupName, operation)
		if isDebugMode() {
			fmt.Fprintf(os.Stderr, "  試行したパターン数: %d\n", len(patterns))
			fmt.Fprintf(os.Stderr, "  最初の3パターン: ")
			for i := 0; i < 3 && i < len(patterns); i++ {
				if i > 0 {
					fmt.Fprintf(os.Stderr, ", ")
				}
				fmt.Fprintf(os.Stderr, "%s", patterns[i].Pattern)
			}
			fmt.Fprintf(os.Stderr, "\n")
		}
		fmt.Fprintf(os.Stderr, "  → フォールバック: デフォルト値を使用します\n")

		baseModule := reporter.moduleName
		if len(moduleOrder) > 0 {
			baseModule = moduleOrder[0]
		}

		triedFilePaths := make([]string, 0, len(moduleOrder))
		for _, module := range moduleOrder {
			triedFilePaths = append(triedFilePaths, fmt.Sprintf("generated/%sapi/types.gen.go", module))
		}
		if len(triedFilePaths) == 0 {
			triedFilePaths = append(triedFilePaths, fmt.Sprintf("generated/%sapi/types.gen.go", reporter.moduleName))
		}

		suggestPattern := "<unknown>"
		if len(patterns) > 0 {
			suggestPattern = patterns[0].Pattern
		}

		for fieldName := range params {
			failedMapping := DetailedFailedMapping{
				ModuleName:     baseModule,
				GroupName:      groupName,
				Operation:      operation,
				FieldName:      fieldName,
				TriedPatterns:  triedPatterns,
				TriedFilePaths: triedFilePaths,
				Reason:         "型定義が見つかりません",
				Suggestion:     fmt.Sprintf("モジュール '%s' の types.gen.go に以下のいずれかの型定義を追加してください: %s", baseModule, suggestPattern),
			}
			reporter.ReportFailedMapping(failedMapping)
			failed = append(failed, failedMapping)
		}
		return 0, failed
	}

	return mapFieldsToType(params, typeInfo, mappingRules, groupName, operation, moduleOrder, reporter)
}

// mapFieldsToType はフィールドを型定義にマッピングします
func mapFieldsToType(params map[string]any, typeInfo EnhancedTypeInfo, mappingRules *MappingRules, groupName, operation string, moduleOrder []string, reporter *DiagnosticReporter) (mappedCount int, failed []DetailedFailedMapping) {
	matchedTypeName := typeInfo.Name
	debugLog("    マッチした型名: %s (フィールド数: %d)\n", matchedTypeName, len(typeInfo.Fields))

	// FieldMatcherを使用してフィールドをマッピング
	fieldMatcher := NewFieldMatcher(typeInfo.Fields, isDebugMode())

	// パラメータの各フィールドに対して型情報を適用
	for fieldName := range params {
		// カスタムルールでフィールド名マッピングを確認（モジュール固有ルールを優先）
		var matchResult *FieldMatchResult
		if mappingRules != nil {
			for _, module := range moduleOrder {
				if mappedFieldName, ok := mappingRules.ApplyFieldNameMappingWithModule(module, matchedTypeName, fieldName); ok {
					debugLog("      カスタムルールでフィールド名マッピング発見: %s.%s.%s -> %s\n", module, matchedTypeName, fieldName, mappedFieldName)
					for _, field := range typeInfo.Fields {
						if field.Name == mappedFieldName || field.JSONTag == mappedFieldName {
							matchResult = &FieldMatchResult{
								Field:         field,
								MatchStrategy: "CustomRule",
								Confidence:    1.0,
							}
							reporter.metrics.MatchStrategyStats["CustomRule"]++
							break
						}
					}
				}
				if matchResult != nil {
					break
				}
			}
		}

		// カスタムルールで見つからない場合はFieldMatcherを使用
		if matchResult == nil {
			var err error
			matchResult, err = fieldMatcher.FindBestMatch(fieldName)
			if err != nil {
				// 警告レベルの判定: 信頼度が低い場合は警告、見つからない場合はエラー
				warningLevel := "警告"
				if fieldMatchErr, ok := err.(*FieldMatchError); ok {
					if len(fieldMatchErr.TriedStrategies) > 3 {
						warningLevel = "エラー"
					}
				}

				if !isDebugMode() {
					// 非デバッグモードでは簡潔なメッセージのみ
					fmt.Fprintf(os.Stderr, "%s: フィールド '%s.%s' のマッピング失敗\n", warningLevel, groupName, fieldName)
				} else {
					debugLog("      フィールド '%s' のマッピング失敗: %v\n", fieldName, err)
				}

				// 失敗情報を記録
				triedFilePaths := make([]string, 0, len(moduleOrder))
				for _, module := range moduleOrder {
					triedFilePaths = append(triedFilePaths, fmt.Sprintf("generated/%sapi/types.gen.go", module))
				}
				if len(triedFilePaths) == 0 {
					triedFilePaths = append(triedFilePaths, fmt.Sprintf("generated/%sapi/types.gen.go", reporter.moduleName))
				}

				baseModule := reporter.moduleName
				if len(moduleOrder) > 0 {
					baseModule = moduleOrder[0]
				}

				failedMapping := DetailedFailedMapping{
					ModuleName:      baseModule,
					GroupName:       groupName,
					Operation:       operation,
					FieldName:       fieldName,
					TriedPatterns:   []TypeNamePattern{},
					TriedStrategies: []string{},
					TriedFilePaths:  triedFilePaths,
					Reason:          "フィールドが型定義に存在しません",
					Suggestion:      fmt.Sprintf("モジュール '%s' の types.gen.go の %s 型に '%s' フィールドを追加するか、mapping_rules.json でモジュール固有のカスタムマッピングを定義してください", baseModule, matchedTypeName, fieldName),
				}

				if fieldMatchErr, ok := err.(*FieldMatchError); ok {
					failedMapping.TriedStrategies = fieldMatchErr.TriedStrategies
				}

				reporter.ReportFailedMapping(failedMapping)
				failed = append(failed, failedMapping)
				continue
			}

			// 信頼度が低い場合は警告を表示
			if matchResult.Confidence < 0.7 {
				fmt.Fprintf(os.Stderr, "警告: フィールド '%s.%s' のマッチング信頼度が低い (%.2f): %s\n",
					groupName, fieldName, matchResult.Confidence, matchResult.MatchStrategy)
			}

			// マッチング戦略の統計を記録
			reporter.metrics.MatchStrategyStats[matchResult.MatchStrategy]++
		}

		// 型に基づいてサンプル値を生成
		sampleValue := generateEnhancedSampleValue(&matchResult.Field)
		params[fieldName] = sampleValue
		mappedCount++
		debugLog("      フィールド '%s' -> 型 '%s' にマッピング (戦略: %s, 信頼度: %.2f, 値: %v)\n",
			fieldName, matchResult.Field.Type, matchResult.MatchStrategy, matchResult.Confidence, sampleValue)
	}

	debugLog("    マッピング完了: %d/%d フィールド成功\n", mappedCount, len(params))
	return mappedCount, failed
}



// normalizeResourceName はリソース名を正規化します（PascalCaseに変換）
func normalizeResourceName(name string) string {
	if name == "" {
		return ""
	}

	// ハイフンやアンダースコアで分割してPascalCaseに変換
	parts := strings.FieldsFunc(name, func(r rune) bool {
		return r == '-' || r == '_'
	})

	var result strings.Builder
	for _, part := range parts {
		if part == "" {
			continue
		}
		// 複数形を単数形に変換（簡易版）
		singular := strings.TrimSuffix(part, "s")
		// 先頭を大文字に
		result.WriteString(strings.ToUpper(singular[:1]) + singular[1:])
	}

	return result.String()
}

// normalizeFieldName はフィールド名を正規化します（snake_case対応）
func normalizeFieldName(name string) string {
	// 既にsnake_caseの場合はそのまま返す
	if strings.Contains(name, "_") {
		return name
	}

	// camelCaseをsnake_caseに変換
	var result strings.Builder
	result.Grow(len(name) + 5) // アンダースコア分を考慮

	for i, r := range name {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}

	return strings.ToLower(result.String())
}

// camelToSnake はキャメルケースをスネークケースに変換します
func camelToSnake(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

// snakeToCamel はsnake_caseをcamelCaseに変換します
func snakeToCamel(name string) string {
	if !strings.Contains(name, "_") {
		return name
	}

	parts := strings.Split(name, "_")
	var result strings.Builder
	result.Grow(len(name))

	for i, part := range parts {
		if part == "" {
			continue
		}
		if i == 0 {
			result.WriteString(strings.ToLower(part))
		} else {
			result.WriteString(strings.Title(part))
		}
	}

	return result.String()
}

// snakeToPascal はsnake_caseをPascalCaseに変換します
func snakeToPascal(name string) string {
	if !strings.Contains(name, "_") {
		return strings.Title(name)
	}

	parts := strings.Split(name, "_")
	var result strings.Builder
	result.Grow(len(name))

	for _, part := range parts {
		if part == "" {
			continue
		}
		result.WriteString(strings.Title(part))
	}

	return result.String()
}

// findMatchingField はフィールド情報から一致するフィールドを探します
func findMatchingField(fields map[string]FieldInfo, originalName, normalizedName string) *FieldInfo {
	// 1. JSONタグで完全一致を探す
	for _, field := range fields {
		jsonName := strings.Split(field.JSONTag, ",")[0]
		if jsonName == originalName {
			return &field
		}
	}

	// 2. 正規化されたJSONタグで一致を探す
	for _, field := range fields {
		jsonName := strings.Split(field.JSONTag, ",")[0]
		if jsonName == normalizedName {
			return &field
		}
	}

	// 3. フィールド名で一致を探す（大文字小文字を無視）
	lowerOriginal := strings.ToLower(originalName)
	lowerNormalized := strings.ToLower(normalizedName)

	for _, field := range fields {
		lowerFieldName := strings.ToLower(field.Name)
		if lowerFieldName == lowerOriginal || lowerFieldName == lowerNormalized {
			return &field
		}
	}

	// 4. camelCase/PascalCase変換を試す
	camelName := snakeToCamel(originalName)
	pascalName := snakeToPascal(originalName)

	for _, field := range fields {
		if field.Name == camelName || field.Name == pascalName {
			return &field
		}

		jsonName := strings.Split(field.JSONTag, ",")[0]
		if jsonName == camelName || jsonName == pascalName {
			return &field
		}
	}

	// 5. 部分一致を試す（最後の手段）
	for _, field := range fields {
		jsonName := strings.Split(field.JSONTag, ",")[0]
		if strings.Contains(strings.ToLower(jsonName), lowerOriginal) ||
			strings.Contains(lowerOriginal, strings.ToLower(jsonName)) {
			debugLog("        部分一致でマッチ: %s ≈ %s\n", originalName, jsonName)
			return &field
		}
	}

	return nil
}

// analyzeTypeMetadata は型名から推測される情報を分析します
func analyzeTypeMetadata(typeName string) TypeMetadata {
	metadata := TypeMetadata{}

	// HTTPメソッドを推測
	if strings.HasPrefix(typeName, "Create") {
		metadata.HTTPMethod = "Create"
	} else if strings.HasPrefix(typeName, "Update") {
		metadata.HTTPMethod = "Update"
	} else if strings.HasPrefix(typeName, "Get") {
		metadata.HTTPMethod = "Get"
	} else if strings.HasPrefix(typeName, "Delete") {
		metadata.HTTPMethod = "Delete"
	} else if strings.HasPrefix(typeName, "List") {
		metadata.HTTPMethod = "List"
	}

	// 型カテゴリを推測
	if strings.Contains(typeName, "JSONRequestBody") || strings.Contains(typeName, "JSONBody") {
		metadata.TypeCategory = "RequestBody"
	} else if strings.Contains(typeName, "Params") || strings.Contains(typeName, "Param") {
		metadata.TypeCategory = "Params"
	} else if strings.Contains(typeName, "Response") || strings.Contains(typeName, "Resp") {
		metadata.TypeCategory = "Response"
	}

	// リソース名を推測（HTTPメソッドと型カテゴリを除去）
	resourceName := typeName
	// HTTPメソッドを除去
	for _, method := range []string{"Create", "Update", "Get", "Delete", "List"} {
		resourceName = strings.TrimPrefix(resourceName, method)
	}
	// 型カテゴリを除去
	for _, category := range []string{"JSONRequestBody", "JSONBody", "Params", "Param", "Request", "Response", "Resp"} {
		resourceName = strings.TrimSuffix(resourceName, category)
	}
	metadata.ResourceName = resourceName

	return metadata
}

// FieldNameNormalizer と generateAliases は field_matcher.go に移動しました

// AnalyzeReport はレポートを分析してパラメータ構造を抽出します
//
// 処理フロー:
// 1. 型情報が利用可能な場合は型情報ベースでパラメータグループを構築（要件11.2）
// 2. 型情報が不足している場合はテストフローからパラメータグループを抽出（フォールバック）
// 3. 生成された型定義にマッピング（オプショナル、失敗しても継続）
// 4. 分析結果を構築して返却
func AnalyzeReport(report *PostmanReport, moduleName string) (*AnalysisResult, error) {
	var paramGroups map[string]*ParamGroup
	var usedFallback bool
	forceDiagnostics := shouldForceMappingDiagnostics()

	// 型情報ベースのパラメータグループ構築を試行（要件11.2）
	typeBasedGroups, typeInfoAvailable := extractParamGroupsFromTypeInfo(report.TestFlow)

	if typeInfoAvailable && len(typeBasedGroups) > 0 {
		// 型情報が利用可能な場合は型情報ベースのグループを使用
		paramGroups = typeBasedGroups
		debugLog("✓ 型情報ベースでパラメータグループを構築しました (%d グループ)\n", len(paramGroups))

		// 成功メッセージを表示
		fmt.Printf("✓ 型情報を使用してパラメータを生成しました (%d グループ)\n", len(paramGroups))

		if forceDiagnostics {
			fallbackGroups := extractParamGroups(report.TestFlow)
			if len(fallbackGroups) > 0 {
				enrichParamGroupsWithRequestBodies(fallbackGroups, report.TestFlow, moduleName)
				expandMustachePlaceholders(fallbackGroups, moduleName)
				cleanParamGroups(fallbackGroups)
				debugLog("診断用にフォールバックマッピングを実行します (%d グループ)\n", len(fallbackGroups))
				fallbackModules := collectModuleCandidates(report, moduleName, fallbackGroups)
				_ = mapToGeneratedTypes(moduleName, fallbackGroups, fallbackModules)
			}
		}
	} else {
		// フォールバック: 従来の方法でパラメータグループを抽出（要件11.5）
		usedFallback = true
		paramGroups = extractParamGroups(report.TestFlow)

		// フォールバック理由を診断（要件11.5）
		if !typeInfoAvailable {
			fmt.Fprintf(os.Stderr, "警告: 型情報が利用できないため、フォールバック処理を使用します\n")
			fmt.Fprintf(os.Stderr, "  理由: postman-analyzerレポートに型情報が含まれていません\n")
			fmt.Fprintf(os.Stderr, "  推奨: postman-analyzerを最新版で実行し、-spec フラグでOpenAPI仕様を指定してください\n")
		} else if len(typeBasedGroups) == 0 {
			fmt.Fprintf(os.Stderr, "警告: 型情報からパラメータグループを構築できなかったため、フォールバック処理を使用します\n")
			fmt.Fprintf(os.Stderr, "  理由: テストフローに有効な型情報が含まれていません\n")
		}

		debugLog("フォールバック: テストフローからパラメータグループを抽出しました (%d グループ)\n", len(paramGroups))

	}

	if report != nil {
		enrichParamGroupsWithRequestBodies(paramGroups, report.TestFlow, moduleName)
	}
	expandMustachePlaceholders(paramGroups, moduleName)
	cleanParamGroups(paramGroups)

	moduleCandidates := collectModuleCandidates(report, moduleName, paramGroups)

	if usedFallback {
		// 生成された型定義にマッピング（エラーは内部で警告として処理される）
		_ = mapToGeneratedTypes(moduleName, paramGroups, moduleCandidates)
	}

	// 診断情報を記録（要件11.5）
	if isDebugMode() {
		debugLog("\n=== 分析結果サマリー ===\n")
		debugLog("モジュール名: %s\n", moduleName)
		debugLog("パラメータグループ数: %d\n", len(paramGroups))
		debugLog("型情報ベース: %v\n", !usedFallback)
		debugLog("フォールバック使用: %v\n", usedFallback)

		for groupName, group := range paramGroups {
			debugLog("\nグループ '%s':\n", groupName)
			if len(group.CreateParams) > 0 {
				debugLog("  CreateParams: %d フィールド\n", len(group.CreateParams))
			}
			if len(group.UpdateParams) > 0 {
				debugLog("  UpdateParams: %d フィールド\n", len(group.UpdateParams))
			}
		}
		debugLog("========================\n\n")
	}

	// 分析結果を作成
	result := &AnalysisResult{
		ModuleName:         moduleName,
		ParamGroups:        paramGroups,
		UtilityFuncs:       []string{"stringPtr", "loadTestParams"},
		ModuleDependencies: moduleCandidates,
	}

	return result, nil
}

func cleanParamGroups(paramGroups map[string]*ParamGroup) {
	if len(paramGroups) == 0 {
		return
	}

	if group, ok := paramGroups["notification-messages"]; ok && group != nil && len(group.UpdateParams) > 0 {
		delete(group.UpdateParams, "sign_up_test")
	}

	if group, ok := paramGroups["customize-pages"]; ok && group != nil && len(group.UpdateParams) > 0 {
		delete(group.UpdateParams, "sign_up_page1")
	}

	if group, ok := paramGroups["tenants"]; ok && group != nil && len(group.UpdateParams) > 0 {
		delete(group.UpdateParams, "next_plan_id")
		delete(group.UpdateParams, "using_next_plan_from")
	}
}

func enrichParamGroupsWithRequestBodies(paramGroups map[string]*ParamGroup, testFlow []TestStep, moduleName string) {
	if len(paramGroups) == 0 || len(testFlow) == 0 {
		return
	}

	for _, step := range testFlow {
		if len(step.RequestBody) == 0 {
			continue
		}

		candidates := make([]string, 0, 4)
		if resource := extractResourceName(step.Endpoint); resource != "" {
			candidates = append(candidates, resource)
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

		group := findParamGroupByCandidates(paramGroups, moduleName, candidates)
		if group == nil {
			continue
		}

		var target map[string]any
		switch step.Method {
		case "POST":
			if group.CreateParams == nil {
				group.CreateParams = make(map[string]any)
			}
			target = group.CreateParams
		case "PUT", "PATCH":
			if group.UpdateParams == nil {
				group.UpdateParams = make(map[string]any)
			}
			target = group.UpdateParams
		default:
			continue
		}

		for key, value := range step.RequestBody {
			normalizedKey := strings.TrimSpace(key)
			if normalizedKey == "" {
				continue
			}
			if existing, exists := target[normalizedKey]; exists && !shouldOverrideRequestBodyValue(existing) {
				continue
			}
			target[normalizedKey] = value
		}
	}
}

func findParamGroupByCandidates(paramGroups map[string]*ParamGroup, moduleName string, candidates []string) *ParamGroup {
	for _, candidate := range candidates {
		candidate = strings.TrimSpace(candidate)
		if candidate == "" {
			continue
		}

		if group, exists := paramGroups[candidate]; exists {
			return group
		}

		normalizedCandidate := normalizeGroupKey(candidate, moduleName)
		for groupName, group := range paramGroups {
			if normalizeGroupKey(groupName, moduleName) == normalizedCandidate {
				return group
			}
		}
	}

	return nil
}

// extractParamGroupsFromTypeInfo は型情報を使用してパラメータグループを構築します
// 要件: 11.2, 11.3, 11.5
//
// 型情報が利用可能な場合、extractParamGroupsをスキップして型情報を直接使用します。
// 戻り値: パラメータグループのマップと型情報が利用可能かどうかのフラグ
func extractParamGroupsFromTypeInfo(testFlow []TestStep) (map[string]*ParamGroup, bool) {
	paramGroups := make(map[string]*ParamGroup)
	hasTypeInfo := false
	missingTypeInfoCount := 0
	totalSteps := len(testFlow)

	for _, step := range testFlow {
		// 型情報が利用可能かチェック（要件11.1）
		if step.ClientMethodName == "" || step.ParameterTypeName == "" || len(step.ParameterTypeFields) == 0 {
			missingTypeInfoCount++
			// 診断情報を記録（要件11.5）
			if isDebugMode() {
				debugLog("型情報が不足しているステップ: %s %s (ClientMethodName: %s, ParameterTypeName: %s, Fields: %d)\n",
					step.Method, step.Endpoint, step.ClientMethodName, step.ParameterTypeName, len(step.ParameterTypeFields))
			}
			continue
		}

		hasTypeInfo = true

		// エンドポイントからリソース名を抽出
		resourceName := extractResourceName(step.Endpoint)
		if resourceName == "" {
			// ClientMethodNameからリソース名を推測
			resourceName = inferResourceNameFromMethodName(step.ClientMethodName)
		}

		if resourceName == "" {
			debugLog("リソース名を推測できませんでした: %s %s (ClientMethodName: %s)\n",
				step.Method, step.Endpoint, step.ClientMethodName)
			continue
		}

		// パラメータグループを取得または作成
		group, exists := paramGroups[resourceName]
		if !exists {
			group = &ParamGroup{
				Name:         resourceName,
				CreateParams: make(map[string]any),
				UpdateParams: make(map[string]any),
			}
			paramGroups[resourceName] = group
		}

		group.ModuleCandidates = appendUniqueModule(group.ModuleCandidates, step.Module)

		// HTTPメソッドに基づいてパラメータを分類
		switch step.Method {
		case "POST":
			// 作成用パラメータ
			group.CreateEndpoint = step.Endpoint
			group.CreateMethod = step.Method
			// 型情報からパラメータを構築（要件11.2）
			buildParamsFromTypeInfo(step.ParameterTypeFields, group.CreateParams)
			debugLog("型情報からCreateParamsを構築: %s (%d フィールド)\n", resourceName, len(step.ParameterTypeFields))

		case "PATCH", "PUT":
			// 更新用パラメータ
			group.UpdateEndpoint = step.Endpoint
			group.UpdateMethod = step.Method
			// 型情報からパラメータを構築（要件11.2）
			buildParamsFromTypeInfo(step.ParameterTypeFields, group.UpdateParams)
			debugLog("型情報からUpdateParamsを構築: %s (%d フィールド)\n", resourceName, len(step.ParameterTypeFields))
		}
	}

	// 診断情報を出力（要件11.5）
	if missingTypeInfoCount > 0 {
		typeInfoCoverage := float64(totalSteps-missingTypeInfoCount) / float64(totalSteps) * 100
		if isDebugMode() {
			debugLog("\n型情報カバレッジ: %.1f%% (%d/%d ステップ)\n", typeInfoCoverage, totalSteps-missingTypeInfoCount, totalSteps)
			debugLog("型情報が不足しているステップ数: %d\n", missingTypeInfoCount)
		}

		// カバレッジが低い場合は警告を表示
		if typeInfoCoverage < 50.0 {
			fmt.Fprintf(os.Stderr, "警告: 型情報のカバレッジが低いです (%.1f%%)。postman-analyzerで型情報を含むレポートを生成してください。\n", typeInfoCoverage)
		}
	}

	return paramGroups, hasTypeInfo
}

// buildParamsFromTypeInfo は型情報からパラメータマップを構築します
// 要件: 11.2, 11.3
func buildParamsFromTypeInfo(typeFields map[string]FieldInfo, params map[string]any) {
	for jsonTag, fieldInfo := range typeFields {
		// JSONタグをキーとして使用
		key := jsonTag
		if key == "" {
			// JSONタグがない場合はフィールド名を使用
			key = fieldInfo.Name
		}

		// 型に基づいてサンプル値を生成
		sampleValue := generateSampleValueFromFieldInfo(fieldInfo)
		params[key] = sampleValue
	}
}

// inferResourceNameFromMethodName はメソッド名からリソース名を推測します
func inferResourceNameFromMethodName(methodName string) string {
	// "Create", "Update", "Get", "Delete" などのプレフィックスを除去
	prefixes := []string{"Create", "Update", "Get", "Delete", "List", "Find"}

	for _, prefix := range prefixes {
		if strings.HasPrefix(methodName, prefix) {
			resourceName := strings.TrimPrefix(methodName, prefix)
			// 最初の文字を小文字に変換
			if len(resourceName) > 0 {
				resourceName = strings.ToLower(resourceName[:1]) + resourceName[1:]
			}
			return resourceName
		}
	}

	// プレフィックスが見つからない場合はメソッド名をそのまま使用
	return strings.ToLower(methodName)
}

// extractParamGroups はエンドポイントごとにパラメータをグループ化します
// 要件: 2.1, 2.3
//
// 処理対象:
// - POST: CreateParamsとして抽出
// - PATCH/PUT: UpdateParamsとして抽出
// - GET/DELETE: 現在は処理対象外（将来の拡張用）
//
// 戻り値: リソース名をキーとしたパラメータグループのマップ
func extractParamGroups(testFlow []TestStep) map[string]*ParamGroup {
	paramGroups := make(map[string]*ParamGroup)

	for _, step := range testFlow {
		// エンドポイントからリソース名を抽出
		resourceName := extractResourceName(step.Endpoint)
		if resourceName == "" {
			continue
		}

		// パラメータグループを取得または作成
		group, exists := paramGroups[resourceName]
		if !exists {
			group = &ParamGroup{
				Name:         resourceName,
				CreateParams: make(map[string]any),
				UpdateParams: make(map[string]any),
			}
			paramGroups[resourceName] = group
		}

		// HTTPメソッドに基づいてパラメータを分類
		switch step.Method {
		case "POST":
			// 作成用パラメータ
			if len(step.RequestBody) > 0 {
				group.CreateEndpoint = step.Endpoint
				group.CreateMethod = step.Method
				for key, value := range step.RequestBody {
					cleanKey := strings.TrimSpace(key)
					if shouldIncludeRequestValue(value) {
						group.CreateParams[cleanKey] = value
					}
				}
			}

		case "PATCH", "PUT":
			// 更新用パラメータ
			if len(step.RequestBody) > 0 {
				group.UpdateEndpoint = step.Endpoint
				group.UpdateMethod = step.Method
				for key, value := range step.RequestBody {
					cleanKey := strings.TrimSpace(key)
					if shouldIncludeRequestValue(value) {
						group.UpdateParams[cleanKey] = value
					}
				}
			}

		}
	}

	return paramGroups
}

func shouldIncludeRequestValue(value any) bool {
	if value == nil {
		return false
	}
	if str, ok := value.(string); ok {
		return strings.TrimSpace(str) != ""
	}
	return true
}

func shouldOverrideRequestBodyValue(existing any) bool {
	if existing == nil {
		return true
	}
	switch v := existing.(type) {
	case string:
		return strings.TrimSpace(v) == ""
	case []interface{}:
		return len(v) == 0
	case map[string]interface{}:
		return len(v) == 0
	}
	return false
}

// extractResourceName はエンドポイントパスからリソース名を抽出します
func extractResourceName(endpoint string) string {
	// エンドポイントをスラッシュで分割
	parts := strings.Split(strings.Trim(endpoint, "/"), "/")
	if len(parts) == 0 {
		return ""
	}

	// パラメータ部分（{xxx}や:xxxで始まる）を除外して結合
	var resourceParts []string
	for _, part := range parts {
		// パラメータ部分をスキップ
		if strings.HasPrefix(part, ":") || strings.HasPrefix(part, "{") || part == "" {
			continue
		}
		resourceParts = append(resourceParts, part)
	}

	// 複数のパーツがある場合はハイフンで結合
	// 例: /stripe/info -> stripe-info
	// 例: /tenants/{tenant_id}/users -> tenants-users
	if len(resourceParts) > 0 {
		return strings.Join(resourceParts, "-")
	}

	return ""
}

// generateEnhancedSampleValue は拡張フィールド情報からサンプル値を生成します
func generateEnhancedSampleValue(field *EnhancedFieldInfo) any {
	// デフォルト値が設定されている場合はそれを使用
	if field.DefaultValue != nil {
		return field.DefaultValue
	}

	// 既存のgenerateSampleValue関数を使用
	value := generateSampleValue(field.Type, field.JSONTag)

	// ポインタ型の場合は値をそのまま返す（ポインタ化は呼び出し側で行う）
	// 配列型の場合は空配列を返す
	if field.IsArray && value != nil {
		return []any{value}
	}

	return value
}

// generateSampleValue はフィールド型からサンプル値を生成します
// 要件: 2.4, 2.5
//
// サンプル値生成ルール:
// 1. フィールド名のパターンマッチング（_id, title, description等）
// 2. 型に基づくデフォルト値（string, int, bool等）
// 3. 日本語のサンプルテキストを優先
func generateSampleValue(fieldType, jsonTag string) any {
	// JSONタグからフィールド名を取得
	jsonName := strings.Split(jsonTag, ",")[0]

	// ポインタ型の場合は基底型を取得
	baseType := strings.TrimPrefix(fieldType, "*")

	// UUID型の判定
	if baseType == "Uuid" || baseType == "string" {
		// ID系のフィールド
		if strings.HasSuffix(jsonName, "_id") || strings.HasSuffix(jsonName, "Id") {
			return "69e732d6-8ecc-45c4-c2eb-8438f7ffe775"
		}

		// タイトル系のフィールド
		if strings.Contains(jsonName, "title") || strings.Contains(jsonName, "Title") {
			return "テスト用タイトル"
		}

		// 説明系のフィールド
		if strings.Contains(jsonName, "description") || strings.Contains(jsonName, "Description") {
			return "これはテスト用の説明です。E2Eテストで使用されるサンプルテキストです。"
		}

		// 本文系のフィールド
		if strings.Contains(jsonName, "body") || strings.Contains(jsonName, "Body") {
			return "これはテスト用の本文です。コメントやフィードバックのテストに使用されます。"
		}

		// 名前系のフィールド
		if strings.Contains(jsonName, "name") || strings.Contains(jsonName, "Name") {
			return "テスト用名前"
		}

		// メールアドレス系のフィールド
		if strings.Contains(jsonName, "email") || strings.Contains(jsonName, "Email") {
			return "test@example.com"
		}

		// URL系のフィールド
		if strings.Contains(jsonName, "url") || strings.Contains(jsonName, "Url") || strings.Contains(jsonName, "URL") {
			return "https://example.com/test"
		}

		// デフォルトの文字列
		return "テスト用文字列"
	}

	// 数値型
	if baseType == "int" || baseType == "int64" || baseType == "uint64" || baseType == "int32" || baseType == "uint32" || baseType == "Status" {
		// ステータス系のフィールド
		if jsonName == "status" || jsonName == "Status" {
			return 0
		}

		// カウント・リミット系のフィールド
		if strings.Contains(jsonName, "count") || strings.Contains(jsonName, "Count") ||
			strings.Contains(jsonName, "limit") || strings.Contains(jsonName, "Limit") {
			return 10
		}

		// ページ番号系のフィールド
		if strings.Contains(jsonName, "page") || strings.Contains(jsonName, "Page") {
			return 1
		}

		// デフォルトの数値
		return 0
	}

	// 浮動小数点型
	if baseType == "float32" || baseType == "float64" {
		if strings.Contains(jsonName, "rate") || strings.Contains(jsonName, "Rate") {
			return 0.0
		}
		if strings.Contains(jsonName, "amount") || strings.Contains(jsonName, "Amount") {
			return 100.0
		}
		return 0.0
	}

	// ブール型
	if baseType == "bool" {
		// enabled系のフィールドはtrueをデフォルトに
		if strings.Contains(jsonName, "enabled") || strings.Contains(jsonName, "Enabled") ||
			strings.Contains(jsonName, "is_") || strings.Contains(jsonName, "Is") {
			return true
		}
		return false
	}

	// 配列型
	if strings.HasPrefix(baseType, "[]") {
		return []any{}
	}

	// マップ型
	if strings.HasPrefix(baseType, "map[") {
		return map[string]any{}
	}

	// time.Time型
	if baseType == "time.Time" || baseType == "Time" {
		return "2024-01-01T00:00:00Z"
	}

	// デフォルト
	return nil
}
