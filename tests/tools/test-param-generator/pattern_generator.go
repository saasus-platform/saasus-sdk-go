package main

import "strings"

// PatternGenerator は型名パターンを生成します
type PatternGenerator struct {
	resourceName string
	operation    string
	typeInfoMap  map[string]EnhancedTypeInfo
	cache        *TypeCache // パターンキャッシュ
}

// NewPatternGenerator は新しいPatternGeneratorを作成します
func NewPatternGenerator(resourceName, operation string, typeInfoMap map[string]EnhancedTypeInfo) *PatternGenerator {
	return &PatternGenerator{
		resourceName: resourceName,
		operation:    operation,
		typeInfoMap:  typeInfoMap,
		cache:        nil, // キャッシュは外部から設定
	}
}

// NewPatternGeneratorWithCache はキャッシュ付きのPatternGeneratorを作成します
func NewPatternGeneratorWithCache(resourceName, operation string, typeInfoMap map[string]EnhancedTypeInfo, cache *TypeCache) *PatternGenerator {
	return &PatternGenerator{
		resourceName: resourceName,
		operation:    operation,
		typeInfoMap:  typeInfoMap,
		cache:        cache,
	}
}

// SetCache はキャッシュを設定します
func (pg *PatternGenerator) SetCache(cache *TypeCache) {
	pg.cache = cache
}

// GeneratePatterns は優先順位付きの型名パターンを生成します
func (pg *PatternGenerator) GeneratePatterns() []TypeNamePattern {
	// キャッシュキーを生成
	cacheKey := pg.resourceName + ":" + pg.operation

	// キャッシュから取得を試みる
	if pg.cache != nil {
		if cached, found := pg.cache.GetPattern(cacheKey); found {
			return cached
		}
	}

	patterns := make([]TypeNamePattern, 0, 50)

	// 各カテゴリのパターンを生成
	patterns = append(patterns, pg.generateExactPatterns()...)
	patterns = append(patterns, pg.generateStandardPatterns()...)
	patterns = append(patterns, pg.generateVariantPatterns()...)
	patterns = append(patterns, pg.generateFuzzyPatterns()...)

	// 重複を除去（最適化版）
	unique := pg.removeDuplicates(patterns)

	// キャッシュに保存
	if pg.cache != nil {
		pg.cache.SetPattern(cacheKey, unique)
	}

	return unique
}

// removeDuplicates は重複パターンを効率的に除去します
func (pg *PatternGenerator) removeDuplicates(patterns []TypeNamePattern) []TypeNamePattern {
	// パターン文字列をキーとしたマップで重複チェック
	seen := make(map[string]int, len(patterns)) // パターン -> 最初に出現したインデックス
	unique := make([]TypeNamePattern, 0, len(patterns))

	for _, p := range patterns {
		if existingIdx, exists := seen[p.Pattern]; exists {
			// 既に存在する場合、優先度が低い方を採用
			if p.Priority < patterns[existingIdx].Priority {
				// 新しいパターンの方が優先度が高い場合、置き換え
				unique[existingIdx] = p
			}
			// 優先度が同じか低い場合はスキップ
		} else {
			// 新しいパターンを追加
			seen[p.Pattern] = len(unique)
			unique = append(unique, p)
		}
	}

	return unique
}

// generateExactPatterns は完全一致パターンを生成します（Priority: 0-9）
func (pg *PatternGenerator) generateExactPatterns() []TypeNamePattern {
	patterns := make([]TypeNamePattern, 0, 10)
	normalized := normalizeResourceName(pg.resourceName)
	priority := 0

	// 最も標準的なoapi-codegenパターン
	patterns = append(patterns, TypeNamePattern{
		Pattern:  pg.operation + normalized + "JSONRequestBody",
		Priority: priority,
		Category: "exact",
	})
	priority++

	patterns = append(patterns, TypeNamePattern{
		Pattern:  pg.operation + normalized + "JSONBody",
		Priority: priority,
		Category: "exact",
	})
	priority++

	patterns = append(patterns, TypeNamePattern{
		Pattern:  pg.operation + normalized + "Body",
		Priority: priority,
		Category: "exact",
	})
	priority++

	// 複数形バージョン
	patterns = append(patterns, TypeNamePattern{
		Pattern:  pg.operation + normalized + "sJSONRequestBody",
		Priority: priority,
		Category: "exact",
	})
	priority++

	patterns = append(patterns, TypeNamePattern{
		Pattern:  pg.operation + normalized + "sJSONBody",
		Priority: priority,
		Category: "exact",
	})

	return patterns
}

// generateStandardPatterns は標準的なパターンを生成します（Priority: 10-19）
func (pg *PatternGenerator) generateStandardPatterns() []TypeNamePattern {
	patterns := make([]TypeNamePattern, 0, 15)
	normalized := normalizeResourceName(pg.resourceName)
	priority := 10

	// Paramsパターン
	patterns = append(patterns, TypeNamePattern{
		Pattern:  pg.operation + normalized + "Params",
		Priority: priority,
		Category: "standard",
	})
	priority++

	patterns = append(patterns, TypeNamePattern{
		Pattern:  pg.operation + normalized + "Param",
		Priority: priority,
		Category: "standard",
	})
	priority++

	// Requestパターン
	patterns = append(patterns, TypeNamePattern{
		Pattern:  pg.operation + normalized + "Request",
		Priority: priority,
		Category: "standard",
	})
	priority++

	// 逆順パターン（リソース名が先）
	patterns = append(patterns, TypeNamePattern{
		Pattern:  normalized + pg.operation + "Params",
		Priority: priority,
		Category: "standard",
	})
	priority++

	patterns = append(patterns, TypeNamePattern{
		Pattern:  normalized + pg.operation + "Param",
		Priority: priority,
		Category: "standard",
	})
	priority++

	patterns = append(patterns, TypeNamePattern{
		Pattern:  normalized + pg.operation + "Request",
		Priority: priority,
		Category: "standard",
	})
	priority++

	patterns = append(patterns, TypeNamePattern{
		Pattern:  normalized + pg.operation + "Body",
		Priority: priority,
		Category: "standard",
	})

	return patterns
}

// generateVariantPatterns はバリエーションパターンを生成します（Priority: 20-29）
func (pg *PatternGenerator) generateVariantPatterns() []TypeNamePattern {
	patterns := make([]TypeNamePattern, 0, 20)
	normalized := normalizeResourceName(pg.resourceName)
	lowercaseOp := strings.ToLower(pg.operation)
	priority := 20

	// 小文字operationパターン
	capitalizedOp := strings.ToUpper(lowercaseOp[:1]) + lowercaseOp[1:]
	patterns = append(patterns, TypeNamePattern{
		Pattern:  normalized + capitalizedOp + "Params",
		Priority: priority,
		Category: "variant",
	})
	priority++

	patterns = append(patterns, TypeNamePattern{
		Pattern:  normalized + capitalizedOp + "Param",
		Priority: priority,
		Category: "variant",
	})
	priority++

	// 複数形対応
	patterns = append(patterns, TypeNamePattern{
		Pattern:  pg.operation + normalized + "sParams",
		Priority: priority,
		Category: "variant",
	})
	priority++

	patterns = append(patterns, TypeNamePattern{
		Pattern:  pg.operation + normalized + "sParam",
		Priority: priority,
		Category: "variant",
	})
	priority++

	patterns = append(patterns, TypeNamePattern{
		Pattern:  normalized + "s" + pg.operation + "Params",
		Priority: priority,
		Category: "variant",
	})
	priority++

	// シンプルパターン（operationなし）
	patterns = append(patterns, TypeNamePattern{
		Pattern:  normalized + "Params",
		Priority: priority,
		Category: "variant",
	})
	priority++

	patterns = append(patterns, TypeNamePattern{
		Pattern:  normalized + "Param",
		Priority: priority,
		Category: "variant",
	})
	priority++

	patterns = append(patterns, TypeNamePattern{
		Pattern:  normalized + "Body",
		Priority: priority,
		Category: "variant",
	})
	priority++

	// 動詞のバリエーション（Create -> New, Add等）
	if pg.operation == "Create" {
		patterns = append(patterns, TypeNamePattern{
			Pattern:  "New" + normalized + "Params",
			Priority: priority,
			Category: "variant",
		})
		priority++

		patterns = append(patterns, TypeNamePattern{
			Pattern:  "Add" + normalized + "Params",
			Priority: priority,
			Category: "variant",
		})
		priority++

		patterns = append(patterns, TypeNamePattern{
			Pattern:  normalized + "NewParams",
			Priority: priority,
			Category: "variant",
		})
		priority++

		patterns = append(patterns, TypeNamePattern{
			Pattern:  normalized + "AddParams",
			Priority: priority,
			Category: "variant",
		})
	}

	// Update -> Modify, Edit等
	if pg.operation == "Update" {
		patterns = append(patterns, TypeNamePattern{
			Pattern:  "Modify" + normalized + "Params",
			Priority: priority,
			Category: "variant",
		})
		priority++

		patterns = append(patterns, TypeNamePattern{
			Pattern:  "Edit" + normalized + "Params",
			Priority: priority,
			Category: "variant",
		})
		priority++

		patterns = append(patterns, TypeNamePattern{
			Pattern:  normalized + "ModifyParams",
			Priority: priority,
			Category: "variant",
		})
		priority++

		patterns = append(patterns, TypeNamePattern{
			Pattern:  normalized + "EditParams",
			Priority: priority,
			Category: "variant",
		})
	}

	return patterns
}

// generateFuzzyPatterns はあいまい一致パターンを生成します（Priority: 30-39）
func (pg *PatternGenerator) generateFuzzyPatterns() []TypeNamePattern {
	patterns := make([]TypeNamePattern, 0, 10)
	normalized := normalizeResourceName(pg.resourceName)
	priority := 30

	// ワイルドカードパターン
	patterns = append(patterns, TypeNamePattern{
		Pattern:  "*" + normalized + "*" + pg.operation + "*",
		Priority: priority,
		Category: "fuzzy",
	})
	priority++

	patterns = append(patterns, TypeNamePattern{
		Pattern:  "*" + pg.operation + "*" + normalized + "*",
		Priority: priority,
		Category: "fuzzy",
	})
	priority++

	patterns = append(patterns, TypeNamePattern{
		Pattern:  normalized + "*Params",
		Priority: priority,
		Category: "fuzzy",
	})
	priority++

	patterns = append(patterns, TypeNamePattern{
		Pattern:  normalized + "*Body",
		Priority: priority,
		Category: "fuzzy",
	})
	priority++

	patterns = append(patterns, TypeNamePattern{
		Pattern:  normalized + "*Request",
		Priority: priority,
		Category: "fuzzy",
	})
	priority++

	patterns = append(patterns, TypeNamePattern{
		Pattern:  "*" + normalized + "*",
		Priority: priority,
		Category: "fuzzy",
	})

	return patterns
}

// matchPattern は型名がパターンにマッチするか判定し、一致度を返します
// 戻り値: (マッチしたか, 一致度 0.0-1.0)
func matchPattern(typeName string, pattern TypeNamePattern) (bool, float64) {
	// 完全一致
	if typeName == pattern.Pattern {
		return true, 1.0
	}

	// ワイルドカードパターンの処理
	if strings.Contains(pattern.Pattern, "*") {
		return matchWildcard(typeName, pattern.Pattern)
	}

	// 大文字小文字を無視した一致
	if strings.EqualFold(typeName, pattern.Pattern) {
		return true, 0.9
	}

	// 部分一致（前方一致）
	if strings.HasPrefix(typeName, pattern.Pattern) {
		matchRatio := float64(len(pattern.Pattern)) / float64(len(typeName))
		return true, 0.7 * matchRatio
	}

	// 部分一致（後方一致）
	if strings.HasSuffix(typeName, pattern.Pattern) {
		matchRatio := float64(len(pattern.Pattern)) / float64(len(typeName))
		return true, 0.6 * matchRatio
	}

	// 部分一致（含む）
	if strings.Contains(typeName, pattern.Pattern) {
		matchRatio := float64(len(pattern.Pattern)) / float64(len(typeName))
		return true, 0.5 * matchRatio
	}

	return false, 0.0
}

// matchWildcard はワイルドカードパターンとのマッチングを行います
func matchWildcard(typeName, pattern string) (bool, float64) {
	// パターンを*で分割
	parts := strings.Split(pattern, "*")

	// 空のパーツを除去
	nonEmptyParts := make([]string, 0, len(parts))
	for _, part := range parts {
		if part != "" {
			nonEmptyParts = append(nonEmptyParts, part)
		}
	}

	if len(nonEmptyParts) == 0 {
		// パターンが"*"のみの場合は全てにマッチ
		return true, 0.3
	}

	// 各パーツが順番に含まれているかチェック
	currentPos := 0
	matchedLength := 0

	for _, part := range nonEmptyParts {
		index := strings.Index(typeName[currentPos:], part)
		if index == -1 {
			return false, 0.0
		}
		currentPos += index + len(part)
		matchedLength += len(part)
	}

	// 一致度を計算（マッチした部分の長さ / 型名の長さ）
	confidence := float64(matchedLength) / float64(len(typeName))

	// ワイルドカードマッチは最大でも0.8の信頼度
	if confidence > 0.8 {
		confidence = 0.8
	}

	return true, confidence
}
