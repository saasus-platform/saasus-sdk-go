package main

import "strings"

// FieldMatcher はフィールド名のマッチングを行います
type FieldMatcher struct {
	typeFields map[string]EnhancedFieldInfo
	verbose    bool
	normalizer *FieldNameNormalizer
}

// NewFieldMatcher は新しいFieldMatcherを作成します
func NewFieldMatcher(typeFields map[string]EnhancedFieldInfo, verbose bool) *FieldMatcher {
	return &FieldMatcher{
		typeFields: typeFields,
		verbose:    verbose,
		normalizer: &FieldNameNormalizer{},
	}
}

// FindBestMatch は最適なフィールドマッチを見つけます
func (fm *FieldMatcher) FindBestMatch(requestFieldName string) (*FieldMatchResult, error) {
	var triedStrategies []string
	var bestMatch *FieldMatchResult

	// 1. 完全一致（JSONタグ）
	if result := fm.tryExactJSONTagMatch(requestFieldName); result != nil {
		return result, nil
	}
	triedStrategies = append(triedStrategies, "Exact JSON Tag Match")

	// 2. 正規化後の一致
	if result := fm.tryNormalizedMatch(requestFieldName); result != nil {
		if bestMatch == nil || result.Confidence > bestMatch.Confidence {
			bestMatch = result
		}
		if result.Confidence >= 0.9 {
			return result, nil
		}
	}
	triedStrategies = append(triedStrategies, "Normalized Match")

	// 3. 大文字小文字を無視した一致
	if result := fm.tryCaseInsensitiveMatch(requestFieldName); result != nil {
		if bestMatch == nil || result.Confidence > bestMatch.Confidence {
			bestMatch = result
		}
		if result.Confidence >= 0.8 {
			return result, nil
		}
	}
	triedStrategies = append(triedStrategies, "Case Insensitive Match")

	// 4. 部分一致
	if result := fm.tryPartialMatch(requestFieldName); result != nil {
		if bestMatch == nil || result.Confidence > bestMatch.Confidence {
			bestMatch = result
		}
		if result.Confidence >= 0.6 {
			return result, nil
		}
	}
	triedStrategies = append(triedStrategies, "Partial Match")

	// 5. 意味的一致
	if result := fm.trySemanticMatch(requestFieldName); result != nil {
		if bestMatch == nil || result.Confidence > bestMatch.Confidence {
			bestMatch = result
		}
	}
	triedStrategies = append(triedStrategies, "Semantic Match")

	// 最良のマッチが見つかった場合
	if bestMatch != nil {
		bestMatch.TriedStrategies = triedStrategies
		return bestMatch, nil
	}

	// マッチが見つからなかった場合
	return nil, &FieldMatchError{
		FieldName:       requestFieldName,
		TriedStrategies: triedStrategies,
		Message:         "no matching field found",
	}
}

// FieldMatchError はフィールドマッチングエラーを表します
type FieldMatchError struct {
	FieldName       string
	TriedStrategies []string
	Message         string
}

func (e *FieldMatchError) Error() string {
	return e.Message
}

// tryExactJSONTagMatch はJSONタグとの完全一致を試みます
func (fm *FieldMatcher) tryExactJSONTagMatch(fieldName string) *FieldMatchResult {
	for _, field := range fm.typeFields {
		if field.JSONTag == fieldName {
			if fm.verbose {
				debugLog("Exact JSON tag match found: %s -> %s", fieldName, field.Name)
			}
			return &FieldMatchResult{
				Field:         field,
				MatchStrategy: "Exact JSON Tag Match",
				Confidence:    1.0,
			}
		}
	}
	return nil
}

// tryNormalizedMatch は正規化後のマッチングを試みます
func (fm *FieldMatcher) tryNormalizedMatch(fieldName string) *FieldMatchResult {
	normalized := fm.normalizer.Normalize(fieldName)

	for _, field := range fm.typeFields {
		// フィールド名の正規化
		fieldNormalized := fm.normalizer.Normalize(field.Name)

		// 正規化後の比較
		for _, norm1 := range normalized {
			for _, norm2 := range fieldNormalized {
				if norm1 == norm2 {
					if fm.verbose {
						debugLog("Normalized match found: %s -> %s (normalized: %s)", fieldName, field.Name, norm1)
					}
					return &FieldMatchResult{
						Field:         field,
						MatchStrategy: "Normalized Match",
						Confidence:    0.9,
					}
				}
			}
		}

		// JSONタグの正規化との比較
		if field.JSONTag != "" {
			jsonTagNormalized := fm.normalizer.Normalize(field.JSONTag)
			for _, norm1 := range normalized {
				for _, norm2 := range jsonTagNormalized {
					if norm1 == norm2 {
						if fm.verbose {
							debugLog("Normalized match with JSON tag: %s -> %s (tag: %s)", fieldName, field.Name, field.JSONTag)
						}
						return &FieldMatchResult{
							Field:         field,
							MatchStrategy: "Normalized Match",
							Confidence:    0.9,
						}
					}
				}
			}
		}
	}

	return nil
}

// tryCaseInsensitiveMatch は大文字小文字を無視したマッチングを試みます
func (fm *FieldMatcher) tryCaseInsensitiveMatch(fieldName string) *FieldMatchResult {
	lowerFieldName := strings.ToLower(fieldName)

	for _, field := range fm.typeFields {
		// フィールド名との比較
		if strings.ToLower(field.Name) == lowerFieldName {
			if fm.verbose {
				debugLog("Case insensitive match found: %s -> %s", fieldName, field.Name)
			}
			return &FieldMatchResult{
				Field:         field,
				MatchStrategy: "Case Insensitive Match",
				Confidence:    0.8,
			}
		}

		// JSONタグとの比較
		if field.JSONTag != "" && strings.ToLower(field.JSONTag) == lowerFieldName {
			if fm.verbose {
				debugLog("Case insensitive match with JSON tag: %s -> %s (tag: %s)", fieldName, field.Name, field.JSONTag)
			}
			return &FieldMatchResult{
				Field:         field,
				MatchStrategy: "Case Insensitive Match",
				Confidence:    0.8,
			}
		}
	}

	return nil
}

// tryPartialMatch は部分一致を試みます
func (fm *FieldMatcher) tryPartialMatch(fieldName string) *FieldMatchResult {
	lowerFieldName := strings.ToLower(fieldName)
	var bestMatch *FieldMatchResult

	for _, field := range fm.typeFields {
		lowerName := strings.ToLower(field.Name)
		lowerJSONTag := strings.ToLower(field.JSONTag)

		// 前方一致
		if strings.HasPrefix(lowerName, lowerFieldName) || strings.HasPrefix(lowerFieldName, lowerName) {
			confidence := 0.7
			if len(lowerFieldName) > len(lowerName)/2 {
				confidence = 0.65
			}

			result := &FieldMatchResult{
				Field:         field,
				MatchStrategy: "Partial Match (Prefix)",
				Confidence:    confidence,
			}

			if bestMatch == nil || result.Confidence > bestMatch.Confidence {
				bestMatch = result
			}
		}

		// 後方一致
		if strings.HasSuffix(lowerName, lowerFieldName) || strings.HasSuffix(lowerFieldName, lowerName) {
			confidence := 0.65

			result := &FieldMatchResult{
				Field:         field,
				MatchStrategy: "Partial Match (Suffix)",
				Confidence:    confidence,
			}

			if bestMatch == nil || result.Confidence > bestMatch.Confidence {
				bestMatch = result
			}
		}

		// 含む
		if strings.Contains(lowerName, lowerFieldName) || strings.Contains(lowerFieldName, lowerName) {
			confidence := 0.6

			result := &FieldMatchResult{
				Field:         field,
				MatchStrategy: "Partial Match (Contains)",
				Confidence:    confidence,
			}

			if bestMatch == nil || result.Confidence > bestMatch.Confidence {
				bestMatch = result
			}
		}

		// JSONタグとの部分一致
		if field.JSONTag != "" {
			if strings.Contains(lowerJSONTag, lowerFieldName) || strings.Contains(lowerFieldName, lowerJSONTag) {
				confidence := 0.6

				result := &FieldMatchResult{
					Field:         field,
					MatchStrategy: "Partial Match (JSON Tag)",
					Confidence:    confidence,
				}

				if bestMatch == nil || result.Confidence > bestMatch.Confidence {
					bestMatch = result
				}
			}
		}
	}

	if bestMatch != nil && fm.verbose {
		debugLog("Partial match found: %s -> %s (strategy: %s, confidence: %.2f)",
			fieldName, bestMatch.Field.Name, bestMatch.MatchStrategy, bestMatch.Confidence)
	}

	return bestMatch
}

// trySemanticMatch は意味的な一致を試みます（シノニム対応）
func (fm *FieldMatcher) trySemanticMatch(fieldName string) *FieldMatchResult {
	// 一般的な省略形と完全形のマッピング
	synonyms := map[string][]string{
		"desc":        {"description", "Description"},
		"description": {"desc", "Desc"},
		"id":          {"ID", "Id", "identifier", "Identifier"},
		"num":         {"number", "Number"},
		"addr":        {"address", "Address"},
		"tel":         {"telephone", "Telephone", "phone", "Phone"},
		"email":       {"mail", "Mail", "e_mail", "EMail"},
		"pwd":         {"password", "Password"},
		"msg":         {"message", "Message"},
		"info":        {"information", "Information"},
		"config":      {"configuration", "Configuration"},
		"param":       {"parameter", "Parameter"},
		"attr":        {"attribute", "Attribute"},
		"val":         {"value", "Value"},
		"cnt":         {"count", "Count"},
		"qty":         {"quantity", "Quantity"},
		"amt":         {"amount", "Amount"},
		"max":         {"maximum", "Maximum"},
		"min":         {"minimum", "Minimum"},
		"avg":         {"average", "Average"},
		"tmp":         {"temporary", "Temporary", "temp", "Temp"},
		"src":         {"source", "Source"},
		"dst":         {"destination", "Destination"},
		"req":         {"request", "Request"},
		"res":         {"response", "Response"},
		"auth":        {"authentication", "Authentication", "authorization", "Authorization"},
	}

	lowerFieldName := strings.ToLower(fieldName)

	// フィールド名のシノニムを取得
	var fieldSynonyms []string
	if syns, ok := synonyms[lowerFieldName]; ok {
		fieldSynonyms = append(fieldSynonyms, syns...)
	}

	// 逆引き（完全形から省略形）
	for key, values := range synonyms {
		for _, val := range values {
			if strings.ToLower(val) == lowerFieldName {
				fieldSynonyms = append(fieldSynonyms, key)
				fieldSynonyms = append(fieldSynonyms, values...)
				break
			}
		}
	}

	// シノニムでマッチング
	for _, field := range fm.typeFields {
		lowerName := strings.ToLower(field.Name)
		lowerJSONTag := strings.ToLower(field.JSONTag)

		for _, synonym := range fieldSynonyms {
			lowerSynonym := strings.ToLower(synonym)

			if lowerName == lowerSynonym || lowerJSONTag == lowerSynonym {
				if fm.verbose {
					debugLog("Semantic match found: %s -> %s (synonym: %s)", fieldName, field.Name, synonym)
				}
				return &FieldMatchResult{
					Field:         field,
					MatchStrategy: "Semantic Match",
					Confidence:    0.5,
				}
			}

			// 部分一致も試す
			if strings.Contains(lowerName, lowerSynonym) || strings.Contains(lowerSynonym, lowerName) {
				if fm.verbose {
					debugLog("Semantic partial match found: %s -> %s (synonym: %s)", fieldName, field.Name, synonym)
				}
				return &FieldMatchResult{
					Field:         field,
					MatchStrategy: "Semantic Match (Partial)",
					Confidence:    0.4,
				}
			}
		}
	}

	return nil
}

// FieldNameNormalizer はフィールド名の正規化を行います
type FieldNameNormalizer struct{}

// Normalize はフィールド名を正規化します
func (n *FieldNameNormalizer) Normalize(name string) []string {
	if name == "" {
		return []string{}
	}

	results := []string{name} // 元の名前も含める

	// snake_case の正規化
	if strings.Contains(name, "_") {
		results = append(results, n.normalizeSnakeCase(name))
	}

	// camelCase の正規化
	if hasLowerThenUpper(name) {
		results = append(results, n.normalizeCamelCase(name))
	}

	// PascalCase の正規化
	if len(name) > 0 && isUpper(name[0]) {
		results = append(results, n.normalizePascalCase(name))
	}

	// 別名の生成
	aliases := n.generateAliases(name)
	results = append(results, aliases...)

	// 重複を除去
	return uniqueStrings(results)
}

// normalizeSnakeCase はsnake_caseを正規化します
func (n *FieldNameNormalizer) normalizeSnakeCase(name string) string {
	// snake_case -> camelCase
	parts := strings.Split(name, "_")
	if len(parts) <= 1 {
		return name
	}

	result := parts[0]
	for i := 1; i < len(parts); i++ {
		if len(parts[i]) > 0 {
			result += strings.ToUpper(string(parts[i][0])) + parts[i][1:]
		}
	}

	return result
}

// normalizeCamelCase はcamelCaseを正規化します
func (n *FieldNameNormalizer) normalizeCamelCase(name string) string {
	// camelCase -> snake_case
	var result strings.Builder

	for i, r := range name {
		if i > 0 && isUpper(byte(r)) {
			result.WriteRune('_')
		}
		result.WriteRune(toLower(r))
	}

	return result.String()
}

// normalizePascalCase はPascalCaseを正規化します
func (n *FieldNameNormalizer) normalizePascalCase(name string) string {
	// PascalCase -> snake_case
	var result strings.Builder

	for i, r := range name {
		if i > 0 && isUpper(byte(r)) {
			result.WriteRune('_')
		}
		result.WriteRune(toLower(r))
	}

	return result.String()
}

// generateAliases はフィールド名の別名を生成します
func (n *FieldNameNormalizer) generateAliases(name string) []string {
	var aliases []string

	// snake_case -> PascalCase
	if strings.Contains(name, "_") {
		parts := strings.Split(name, "_")
		var pascalCase strings.Builder
		for _, part := range parts {
			if len(part) > 0 {
				pascalCase.WriteString(strings.ToUpper(string(part[0])) + part[1:])
			}
		}
		if pascalCase.Len() > 0 {
			aliases = append(aliases, pascalCase.String())
		}

		// snake_case -> camelCase
		camelCase := n.normalizeSnakeCase(name)
		if camelCase != name {
			aliases = append(aliases, camelCase)
		}
	}

	// camelCase -> PascalCase
	if len(name) > 0 && !isUpper(name[0]) && hasLowerThenUpper(name) {
		pascalCase := strings.ToUpper(string(name[0])) + name[1:]
		aliases = append(aliases, pascalCase)

		// camelCase -> snake_case
		snakeCase := n.normalizeCamelCase(name)
		if snakeCase != name {
			aliases = append(aliases, snakeCase)
		}
	}

	// PascalCase -> camelCase
	if len(name) > 0 && isUpper(name[0]) {
		camelCase := strings.ToLower(string(name[0])) + name[1:]
		aliases = append(aliases, camelCase)

		// PascalCase -> snake_case
		snakeCase := n.normalizePascalCase(name)
		if snakeCase != name {
			aliases = append(aliases, snakeCase)
		}
	}

	return aliases
}

// ヘルパー関数

// hasLowerThenUpper は文字列に小文字の後に大文字が続くパターンがあるか確認します
func hasLowerThenUpper(s string) bool {
	for i := 0; i < len(s)-1; i++ {
		if !isUpper(s[i]) && isUpper(s[i+1]) {
			return true
		}
	}
	return false
}

// isUpper は文字が大文字かどうか確認します
func isUpper(b byte) bool {
	return b >= 'A' && b <= 'Z'
}

// toLower はruneを小文字に変換します
func toLower(r rune) rune {
	if r >= 'A' && r <= 'Z' {
		return r + ('a' - 'A')
	}
	return r
}

// uniqueStrings は文字列スライスから重複を除去します
func uniqueStrings(strs []string) []string {
	seen := make(map[string]bool)
	result := []string{}

	for _, s := range strs {
		if s != "" && !seen[s] {
			seen[s] = true
			result = append(result, s)
		}
	}

	return result
}
