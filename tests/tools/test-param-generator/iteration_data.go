package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// IterationDataLoader はiterationDataの読み込みを担当します
type IterationDataLoader struct {
	ModuleName string
	BaseDir    string // iterationDataのベースディレクトリ
	Verbose    bool   // 詳細ログ出力
}

// iterationDataModuleAliases はモジュール名とiterationDataディレクトリ名のエイリアス対応を保持します
var iterationDataModuleAliases = map[string][]string{
	"integration": {"eventbridge"},
	"eventbridge": {"integration"},
}

// uniqueStringsForIterationData は重複を排除した文字列スライスを返します
func uniqueStringsForIterationData(values []string) []string {
	set := make(map[string]struct{})
	result := make([]string, 0, len(values))
	for _, v := range values {
		if _, exists := set[v]; exists {
			continue
		}
		set[v] = struct{}{}
		result = append(result, v)
	}
	return result
}

// NewIterationDataLoader は新しいIterationDataLoaderを作成します
func NewIterationDataLoader(moduleName string, verbose bool) *IterationDataLoader {
	// 現在の作業ディレクトリを取得
	cwd, err := os.Getwd()
	baseDir := "iterationData"

	// tests/tools/test-param-generatorディレクトリから実行されていない場合
	if err != nil || !strings.HasSuffix(cwd, "tests/tools/test-param-generator") {
		baseDir = "tests/tools/test-param-generator/iterationData"
	}

	return &IterationDataLoader{
		ModuleName: moduleName,
		BaseDir:    baseDir,
		Verbose:    verbose,
	}
}

// findIterationDataFiles はiterationDataディレクトリ内のJSONファイルを検索します
func (l *IterationDataLoader) findIterationDataFiles() ([]string, error) {
	candidateDirs := []string{l.ModuleName}
	if aliases, ok := iterationDataModuleAliases[l.ModuleName]; ok {
		candidateDirs = append(candidateDirs, aliases...)
	}
	candidateDirs = uniqueStringsForIterationData(candidateDirs)

	var jsonFiles []string

	for _, dirName := range candidateDirs {
		moduleDir := filepath.Join(l.BaseDir, dirName)

		if l.Verbose {
			if dirName == l.ModuleName {
				fmt.Printf("iterationDataディレクトリを検索中: %s\n", moduleDir)
			} else {
				fmt.Printf("iterationDataエイリアス '%s' を検索中: %s\n", dirName, moduleDir)
			}
		}

		if _, err := os.Stat(moduleDir); os.IsNotExist(err) {
			if l.Verbose {
				fmt.Printf("  エイリアス '%s' のディレクトリは存在しません\n", dirName)
			}
			continue
		}

		err := filepath.Walk(moduleDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			if strings.HasSuffix(strings.ToLower(info.Name()), ".json") {
				jsonFiles = append(jsonFiles, path)
				if l.Verbose {
					fmt.Printf("  見つかったJSONファイル: %s\n", path)
				}
			}

			return nil
		})

		if err != nil {
			return nil, fmt.Errorf("iterationDataディレクトリの探索に失敗: %w", err)
		}
	}

	if l.Verbose {
		fmt.Printf("合計 %d 個のJSONファイルが見つかりました\n", len(jsonFiles))
	}

	return jsonFiles, nil
}

// loadJSONFile はJSONファイルを読み込んで配列として返します
func (l *IterationDataLoader) loadJSONFile(filePath string) ([]map[string]interface{}, error) {
	if l.Verbose {
		fmt.Printf("JSONファイルを読み込み中: %s\n", filePath)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("ファイルの読み込みに失敗: %w", err)
	}

	// まず配列形式として解析を試みる
	var arrayData []map[string]interface{}
	if err := json.Unmarshal(data, &arrayData); err == nil {
		if l.Verbose {
			fmt.Printf("  配列形式として解析成功: %d 要素\n", len(arrayData))
		}
		return arrayData, nil
	}

	// 配列形式でない場合、単一オブジェクトとして解析
	var singleData map[string]interface{}
	if err := json.Unmarshal(data, &singleData); err != nil {
		return nil, fmt.Errorf("JSONの解析に失敗: %w", err)
	}

	if l.Verbose {
		fmt.Printf("  単一オブジェクトとして解析成功\n")
	}

	// 単一オブジェクトを配列にラップして返す
	return []map[string]interface{}{singleData}, nil
}

// LoadIterationData はモジュールのiterationDataを読み込みます
func (l *IterationDataLoader) LoadIterationData() (map[string][]map[string]interface{}, error) {
	files, err := l.findIterationDataFiles()
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		if l.Verbose {
			fmt.Printf("iterationDataファイルが見つかりませんでした\n")
		}
		return make(map[string][]map[string]interface{}), nil
	}

	result := make(map[string][]map[string]interface{})

	for _, filePath := range files {
		data, err := l.loadJSONFile(filePath)
		if err != nil {
			if l.Verbose {
				fmt.Printf("警告: %s の読み込みに失敗しました: %v\n", filePath, err)
			}
			continue
		}

		// ファイル名からキーを生成（拡張子を除く）
		fileName := filepath.Base(filePath)
		key := strings.TrimSuffix(fileName, filepath.Ext(fileName))

		result[key] = data
	}

	if l.Verbose {
		fmt.Printf("合計 %d 個のiterationDataファイルを読み込みました\n", len(result))
	}

	return result, nil
}

// extractFieldsBySuffix はサフィックスに基づいてフィールドを抽出します
// 例: _update, _create, _delete
func extractFieldsBySuffix(data map[string]interface{}, suffix string) map[string]interface{} {
	result := make(map[string]interface{})

	for key, value := range data {
		// サフィックスが一致する場合
		if strings.HasSuffix(key, suffix) {
			// サフィックスを除いた正規化されたフィールド名を使用
			normalizedKey := removeFieldSuffix(key, suffix)
			result[normalizedKey] = value
		}
	}

	return result
}

// removeFieldSuffix はフィールド名からサフィックスを除去します
// 例: "field_name_update" -> "field_name"
func removeFieldSuffix(fieldName, suffix string) string {
	if strings.HasSuffix(fieldName, suffix) {
		return strings.TrimSuffix(fieldName, suffix)
	}
	return fieldName
}

// matchIterationDataToParamGroup はiterationDataをパラメータグループにマッチングします
// 要件: 11.3 - JSONタグとフィールド名の両方でマッチングを試行し、大文字小文字を無視
func matchIterationDataToParamGroup(
	iterationData map[string][]map[string]interface{},
	paramGroup *ParamGroup,
	resourceName string,
) error {
	if len(iterationData) == 0 {
		return nil
	}

	normalizedResource := normalizeIterationDataKey(resourceName)
	debugEnabled := shouldDebugIterationResource(resourceName)

	type candidate struct {
		key      string
		priority int
		distance int
		record   map[string]interface{}
	}

	candidates := make([]candidate, 0, len(iterationData))

	for key, dataArray := range iterationData {
		if len(dataArray) == 0 {
			continue
		}

		normalizedKey := normalizeIterationDataKey(key)
		priority, distance, ok := classifyIterationDataMatch(normalizedResource, normalizedKey)
		if !ok {
			continue
		}

		candidates = append(candidates, candidate{
			key:      key,
			priority: priority,
			distance: distance,
			record:   dataArray[0],
		})
	}

	if len(candidates) == 0 {
		if len(iterationData) == 1 {
			for _, dataArray := range iterationData {
				if len(dataArray) > 0 {
					mergeIterationRecordIntoGroup(paramGroup, dataArray[0])
				}
			}
		}
		return nil
	}

	sort.SliceStable(candidates, func(i, j int) bool {
		if candidates[i].priority != candidates[j].priority {
			return candidates[i].priority < candidates[j].priority
		}
		if candidates[i].distance != candidates[j].distance {
			return candidates[i].distance < candidates[j].distance
		}
		return candidates[i].key < candidates[j].key
	})

	if debugEnabled {
		fmt.Printf("[TPG DEBUG] iterationData candidates for '%s' (normalized: %s)\n", resourceName, normalizedResource)
		for _, c := range candidates {
			preview := summarizeIterationRecord(c.record)
			fmt.Printf("  - key=%s priority=%d distance=%d preview=%s\n", c.key, c.priority, c.distance, preview)
		}
	}

	for _, candidate := range candidates {
		mergeIterationRecordIntoGroup(paramGroup, candidate.record)
	}

	return nil
}

func shouldDebugIterationResource(resourceName string) bool {
	debugEnv := os.Getenv("TPG_DEBUG_ITERATION")
	if debugEnv == "" {
		return false
	}
	segments := strings.Split(debugEnv, ",")
	targets := make(map[string]struct{}, len(segments))
	for _, seg := range segments {
		seg = strings.TrimSpace(seg)
		if seg == "" {
			continue
		}
		targets[seg] = struct{}{}
	}
	_, enabled := targets[resourceName]
	return enabled
}

func normalizeIterationDataKey(key string) string {
	return normalizeGroupKey(key, "")
}

func classifyIterationDataMatch(resourceKey, candidateKey string) (int, int, bool) {
	if resourceKey == "" || candidateKey == "" {
		return 0, 0, false
	}

	const (
		priorityPartial = 0
		priorityAffix   = 1
		priorityExact   = 2
	)

	if candidateKey == resourceKey {
		return priorityExact, 0, true
	}

	if strings.HasPrefix(candidateKey, resourceKey) || strings.HasSuffix(candidateKey, resourceKey) ||
		strings.HasPrefix(resourceKey, candidateKey) || strings.HasSuffix(resourceKey, candidateKey) {
		distance := absInt(len(candidateKey) - len(resourceKey))
		return priorityAffix, distance, true
	}

	if strings.Contains(candidateKey, resourceKey) || strings.Contains(resourceKey, candidateKey) {
		distance := absInt(len(candidateKey) - len(resourceKey))
		return priorityPartial, distance, true
	}

	return 0, 0, false
}

func mergeIterationRecordIntoGroup(paramGroup *ParamGroup, record map[string]interface{}) {
	if len(record) == 0 {
		return
	}

	createFields := extractFieldsBySuffix(record, "_create")
	if len(createFields) > 0 {
		if paramGroup.CreateParams == nil {
			paramGroup.CreateParams = make(map[string]interface{})
		}
		for _, key := range sortedIterationKeys(createFields) {
			matchAndAssignValue(paramGroup.CreateParams, key, createFields[key], record)
		}
	} else {
		keys := sortedIterationKeys(record)
		for _, key := range keys {
			if strings.Contains(key, "_update") || strings.Contains(key, "_delete") || strings.Contains(key, "_vote") {
				continue
			}
			if paramGroup.CreateParams == nil {
				paramGroup.CreateParams = make(map[string]interface{})
			}
			matchAndAssignValue(paramGroup.CreateParams, key, record[key], record)
		}
	}

	updateFields := extractFieldsBySuffix(record, "_update")
	if len(updateFields) > 0 {
		if paramGroup.UpdateParams == nil {
			paramGroup.UpdateParams = make(map[string]interface{})
		}
		for _, key := range sortedIterationKeys(updateFields) {
			matchAndAssignValue(paramGroup.UpdateParams, key, updateFields[key], record)
		}
	}
}

func sortedIterationKeys(data map[string]interface{}) []string {
	keys := make([]string, 0, len(data))
	for key := range data {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func absInt(value int) int {
	if value < 0 {
		return -value
	}
	return value
}

func summarizeIterationRecord(record map[string]interface{}) string {
	if len(record) == 0 {
		return "{}"
	}
	keys := make([]string, 0, len(record))
	for key := range record {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	parts := make([]string, 0, minInt(3, len(keys)))
	for idx, key := range keys {
		if idx >= 3 {
			break
		}
		parts = append(parts, fmt.Sprintf("%s=%v", key, record[key]))
	}
	return strings.Join(parts, ", ")
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// matchAndAssignValue は改善されたマッチングロジックを使用して値を割り当てます
// 要件: 11.3 - JSONタグとフィールド名の両方でマッチングを試行し、大文字小文字を無視
func matchAndAssignValue(params map[string]interface{}, key string, value interface{}, allData map[string]interface{}) {
	// 1. 完全一致を試行
	if _, exists := params[key]; exists {
		params[key] = value
		return
	}

	// 2. 大文字小文字を無視した一致を試行（要件11.3）
	keyLower := strings.ToLower(key)
	for paramKey := range params {
		if strings.ToLower(paramKey) == keyLower {
			params[paramKey] = value
			return
		}
	}

	// 3. スネークケース <-> キャメルケース変換を試行（要件11.3）
	// スネークケースをキャメルケースに変換
	camelKey := snakeToCamel(key)
	for paramKey := range params {
		if strings.EqualFold(paramKey, camelKey) {
			params[paramKey] = value
			return
		}
	}

	// キャメルケースをスネークケースに変換
	snakeKey := camelToSnake(key)
	for paramKey := range params {
		if strings.EqualFold(paramKey, snakeKey) {
			params[paramKey] = value
			return
		}
	}

	// 4. マッチしない場合は新しいキーとして追加
	params[key] = value
}
