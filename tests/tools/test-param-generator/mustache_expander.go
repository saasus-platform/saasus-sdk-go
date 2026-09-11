package main

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
)

var mustachePattern = regexp.MustCompile(`^\{\{\s*([^{}\s]+)\s*\}\}$`)

// expandMustachePlaceholders は iterationData の値を使用してプレースホルダを解決します。
// フォールバック解析で抽出したパラメータマップが {{variable}} のようなキー／値を含んでいる
// 場合に、実際の値へ展開してから型マッピングを行えるようにするためのプリプロセスです。
func expandMustachePlaceholders(paramGroups map[string]*ParamGroup, moduleName string) {
	if len(paramGroups) == 0 {
		return
	}

	loader := NewIterationDataLoader(moduleName, false)
	iterationData, err := loader.LoadIterationData()
	if err != nil {
		fmt.Fprintf(os.Stderr, "警告: iterationDataの読み込みに失敗しました: %v\n", err)
		return
	}

	for groupName, group := range paramGroups {
		if group == nil {
			continue
		}

		records := findIterationRecordsForGroup(iterationData, groupName)
		if len(records) == 0 {
			// iterationData が存在しない場合でも、{{field}} -> field への正規化だけは実施
			if containsMustacheEntries(group.CreateParams) {
				group.CreateParams = normalizeMustacheEntries(group.CreateParams)
			}
			if containsMustacheEntries(group.UpdateParams) {
				group.UpdateParams = normalizeMustacheEntries(group.UpdateParams)
			}
			continue
		}

		if len(group.CreateParams) > 0 {
			group.CreateParams = expandParamMapWithRecords(group.CreateParams, records)
		}
		if len(group.UpdateParams) > 0 {
			group.UpdateParams = expandParamMapWithRecords(group.UpdateParams, records)
		}
	}
}

// findIterationRecordsForGroup はリソース名にマッチする iterationData エントリを検索します。
func findIterationRecordsForGroup(iterationData map[string][]map[string]interface{}, groupName string) []map[string]interface{} {
	if len(iterationData) == 0 {
		return nil
	}

	normalizedGroup := normalizeKey(groupName)

	if records, ok := iterationData[groupName]; ok {
		return records
	}

	exactKeys := make([]string, 0)
	partialKeys := make([]string, 0)
	for key := range iterationData {
		normalizedKey := normalizeKey(key)
		if normalizedKey == normalizedGroup {
			exactKeys = append(exactKeys, key)
			continue
		}
		if strings.Contains(normalizedKey, normalizedGroup) || strings.Contains(normalizedGroup, normalizedKey) {
			partialKeys = append(partialKeys, key)
		}
	}

	if len(exactKeys) > 0 {
		sort.Strings(exactKeys)
		return collectRecordsInOrder(iterationData, exactKeys)
	}

	if len(partialKeys) > 0 {
		sort.Strings(partialKeys)
		return collectRecordsInOrder(iterationData, partialKeys)
	}

	// フォールバック: 先頭のiterationDataを使用（キー順を固定）
	allKeys := make([]string, 0, len(iterationData))
	for key := range iterationData {
		allKeys = append(allKeys, key)
	}
	sort.Strings(allKeys)
	for _, key := range allKeys {
		records := iterationData[key]
		if len(records) > 0 {
			return records
		}
	}

	return nil
}

func collectRecordsInOrder(iterationData map[string][]map[string]interface{}, keys []string) []map[string]interface{} {
	var result []map[string]interface{}
	for _, key := range keys {
		records := iterationData[key]
		if len(records) == 0 {
			continue
		}
		result = append(result, records...)
	}
	return result
}

func normalizeKey(value string) string {
	lower := strings.ToLower(value)
	replacer := strings.NewReplacer("-", "_", " ", "")
	return replacer.Replace(lower)
}

// expandParamMapWithRecords は Mustache プレースホルダを実際の iterationData に展開します。
func expandParamMapWithRecords(params map[string]any, records []map[string]interface{}) map[string]any {
	if len(params) == 0 {
		return params
	}

	keys := make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	changed := false
	result := make(map[string]any)

	for _, key := range keys {
		value := params[key]
		placeholderName, isPlaceholder := extractMustacheName(key)
		if !isPlaceholder {
			continue
		}
		expanded := expandPlaceholderKey(placeholderName, value, records)
		if len(expanded) == 0 {
			continue
		}
		for expandedKey, expandedValue := range expanded {
			result[expandedKey] = expandedValue
		}
		changed = true
	}

	for _, key := range keys {
		value := params[key]
		if _, isPlaceholder := extractMustacheName(key); isPlaceholder {
			// 既にプレースホルダ展開済み
			continue
		}
		resolvedValue, valueChanged := resolvePlaceholdersInValue(value, firstRecord(records))
		if valueChanged {
			changed = true
		}
		result[key] = resolvedValue
	}

	if !changed {
		return params
	}

	return result
}

func expandPlaceholderKey(name string, rawValue any, records []map[string]interface{}) map[string]any {
	output := make(map[string]any)
	for _, record := range records {
		actual := lookupString(record, name)
		if actual == "" {
			continue
		}
		resolvedValue, _ := resolvePlaceholdersInValue(rawValue, record)
		output[actual] = resolvedValue
	}

	if len(output) == 0 {
		resolvedValue, _ := resolvePlaceholdersInValue(rawValue, nil)
		output[name] = resolvedValue
	}

	return output
}

func resolvePlaceholdersInValue(value any, record map[string]interface{}) (any, bool) {
	switch typed := value.(type) {
	case string:
		if name, ok := extractMustacheName(typed); ok {
			if record != nil {
				if actual := lookupString(record, name); actual != "" {
					return actual, true
				}
			}
			return name, true
		}
		return value, false
	case map[string]any:
		if len(typed) == 0 {
			return value, false
		}
		changed := false
		expanded := make(map[string]any, len(typed))
		for k, v := range typed {
			newKey := k
			if name, ok := extractMustacheName(k); ok {
				if record != nil {
					if actual := lookupString(record, name); actual != "" {
						newKey = actual
					} else {
						newKey = name
					}
				} else {
					newKey = name
				}
				changed = true
			}
			resolvedValue, valueChanged := resolvePlaceholdersInValue(v, record)
			if valueChanged {
				changed = true
			}
			expanded[newKey] = resolvedValue
		}
		if changed {
			return expanded, true
		}
		return value, false
	case []any:
		if len(typed) == 0 {
			return value, false
		}
		changed := false
		resolved := make([]any, len(typed))
		for i, elem := range typed {
			resolvedElem, elemChanged := resolvePlaceholdersInValue(elem, record)
			if elemChanged {
				changed = true
			}
			resolved[i] = resolvedElem
		}
		if changed {
			return resolved, true
		}
		return value, false
	default:
		return value, false
	}
}

func lookupString(record map[string]interface{}, key string) string {
	if record == nil {
		return ""
	}
	if value, exists := record[key]; exists {
		return strings.TrimSpace(fmt.Sprint(value))
	}
	return ""
}

func firstRecord(records []map[string]interface{}) map[string]interface{} {
	if len(records) == 0 {
		return nil
	}
	return records[0]
}

func extractMustacheName(value string) (string, bool) {
	if !mustachePattern.MatchString(value) {
		return "", false
	}
	matches := mustachePattern.FindStringSubmatch(value)
	if len(matches) != 2 {
		return "", false
	}
	return strings.TrimSpace(matches[1]), true
}

func containsMustacheEntries(params map[string]any) bool {
	for key, value := range params {
		if _, ok := extractMustacheName(key); ok {
			return true
		}
		if str, ok := value.(string); ok {
			if _, ok := extractMustacheName(str); ok {
				return true
			}
		}
		if nestedMap, ok := value.(map[string]any); ok {
			if containsMustacheEntries(nestedMap) {
				return true
			}
		}
		if array, ok := value.([]any); ok {
			for _, elem := range array {
				if elemMap, ok := elem.(map[string]any); ok {
					if containsMustacheEntries(elemMap) {
						return true
					}
				}
				if elemStr, ok := elem.(string); ok {
					if _, ok := extractMustacheName(elemStr); ok {
						return true
					}
				}
			}
		}
	}
	return false
}

func normalizeMustacheEntries(params map[string]any) map[string]any {
	if len(params) == 0 {
		return params
	}
	changed := false
	result := make(map[string]any, len(params))
	for key, value := range params {
		newKey := key
		if name, ok := extractMustacheName(key); ok {
			newKey = name
			changed = true
		}
		resolvedValue, valueChanged := resolvePlaceholdersInValue(value, nil)
		if valueChanged {
			changed = true
		}
		result[newKey] = resolvedValue
	}
	if changed {
		return result
	}
	return params
}
