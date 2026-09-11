package main

import (
	"fmt"
	"sync"
	"time"
)

// TypeCache は型情報のキャッシュを管理します
type TypeCache struct {
	cache         map[string]*TypeCacheEntry
	patternCache  map[string]*PatternCacheEntry
	mutex         sync.RWMutex
	maxSize       int   // 最大キャッシュサイズ
	ttl           int64 // Time To Live（秒）
	memoryLimit   int64 // メモリ使用量の上限（バイト）
	currentMemory int64 // 現在のメモリ使用量推定値
}

// NewTypeCache は新しいTypeCacheを作成します
func NewTypeCache(maxSize int, ttlSeconds int64, memoryLimitMB int64) *TypeCache {
	return &TypeCache{
		cache:        make(map[string]*TypeCacheEntry),
		patternCache: make(map[string]*PatternCacheEntry),
		maxSize:      maxSize,
		ttl:          ttlSeconds,
		memoryLimit:  memoryLimitMB * 1024 * 1024, // MBをバイトに変換
	}
}

// Get はキャッシュから型情報を取得します
func (tc *TypeCache) Get(moduleName string) (map[string]EnhancedTypeInfo, bool) {
	tc.mutex.RLock()
	entry, exists := tc.cache[moduleName]
	tc.mutex.RUnlock()

	if !exists {
		return nil, false
	}

	// TTLチェック
	if tc.ttl > 0 && time.Now().Unix()-entry.Timestamp > tc.ttl {
		// TTL切れの場合は削除
		tc.mutex.Lock()
		if entry, stillExists := tc.cache[moduleName]; stillExists {
			estimatedSize := tc.estimateTypeInfoSize(entry.TypeInfo)
			delete(tc.cache, moduleName)
			tc.currentMemory -= estimatedSize
			if tc.currentMemory < 0 {
				tc.currentMemory = 0
			}
		}
		tc.mutex.Unlock()
		return nil, false
	}

	return entry.TypeInfo, true
}

// Set はキャッシュに型情報を保存します
func (tc *TypeCache) Set(moduleName string, typeInfo map[string]EnhancedTypeInfo) error {
	tc.mutex.Lock()
	defer tc.mutex.Unlock()

	// 既存エントリがある場合は古いサイズを減算
	if existingEntry, exists := tc.cache[moduleName]; exists {
		oldSize := tc.estimateTypeInfoSize(existingEntry.TypeInfo)
		tc.currentMemory -= oldSize
		if tc.currentMemory < 0 {
			tc.currentMemory = 0
		}
	}

	// メモリ使用量の推定（簡易版）
	estimatedSize := tc.estimateTypeInfoSize(typeInfo)

	// メモリ制限内に収まるまで削除を繰り返す
	if tc.memoryLimit > 0 {
		for tc.currentMemory+estimatedSize > tc.memoryLimit && len(tc.cache) > 0 {
			tc.evictOldestEntries()
		}
	}

	// サイズ制限チェック
	if tc.maxSize > 0 {
		for len(tc.cache) >= tc.maxSize && len(tc.cache) > 0 {
			tc.evictOldestEntries()
		}
	}

	tc.cache[moduleName] = &TypeCacheEntry{
		TypeInfo:  typeInfo,
		Timestamp: time.Now().Unix(),
	}
	tc.currentMemory += estimatedSize

	return nil
}

// GetPattern はキャッシュからパターン生成結果を取得します
func (tc *TypeCache) GetPattern(key string) ([]TypeNamePattern, bool) {
	tc.mutex.RLock()
	entry, exists := tc.patternCache[key]
	tc.mutex.RUnlock()

	if !exists {
		return nil, false
	}

	// TTLチェック
	if tc.ttl > 0 && time.Now().Unix()-entry.Timestamp > tc.ttl {
		// TTL切れの場合は削除
		tc.mutex.Lock()
		if entry, stillExists := tc.patternCache[key]; stillExists {
			estimatedSize := tc.estimatePatternSize(entry.Patterns)
			delete(tc.patternCache, key)
			tc.currentMemory -= estimatedSize
			if tc.currentMemory < 0 {
				tc.currentMemory = 0
			}
		}
		tc.mutex.Unlock()
		return nil, false
	}

	return entry.Patterns, true
}

// SetPattern はキャッシュにパターン生成結果を保存します
func (tc *TypeCache) SetPattern(key string, patterns []TypeNamePattern) error {
	tc.mutex.Lock()
	defer tc.mutex.Unlock()

	// 既存エントリがある場合は古いサイズを減算
	if existingEntry, exists := tc.patternCache[key]; exists {
		oldSize := tc.estimatePatternSize(existingEntry.Patterns)
		tc.currentMemory -= oldSize
		if tc.currentMemory < 0 {
			tc.currentMemory = 0
		}
	}

	// メモリ使用量の推定
	estimatedSize := tc.estimatePatternSize(patterns)

	// メモリ制限内に収まるまで削除を繰り返す
	if tc.memoryLimit > 0 {
		for tc.currentMemory+estimatedSize > tc.memoryLimit && len(tc.patternCache) > 0 {
			tc.evictOldestPatternEntries()
		}
	}

	tc.patternCache[key] = &PatternCacheEntry{
		Patterns:  patterns,
		Timestamp: time.Now().Unix(),
	}
	tc.currentMemory += estimatedSize

	return nil
}

// Clear はキャッシュをクリアします
func (tc *TypeCache) Clear() {
	tc.mutex.Lock()
	defer tc.mutex.Unlock()

	tc.cache = make(map[string]*TypeCacheEntry)
	tc.patternCache = make(map[string]*PatternCacheEntry)
	tc.currentMemory = 0
}

// ClearModule は特定のモジュールのキャッシュをクリアします
func (tc *TypeCache) ClearModule(moduleName string) {
	tc.mutex.Lock()
	defer tc.mutex.Unlock()

	if entry, exists := tc.cache[moduleName]; exists {
		estimatedSize := tc.estimateTypeInfoSize(entry.TypeInfo)
		delete(tc.cache, moduleName)
		tc.currentMemory -= estimatedSize
		if tc.currentMemory < 0 {
			tc.currentMemory = 0
		}
	}
}

// GetStats はキャッシュの統計情報を取得します
func (tc *TypeCache) GetStats() map[string]interface{} {
	tc.mutex.RLock()
	defer tc.mutex.RUnlock()

	return map[string]interface{}{
		"type_cache_entries":    len(tc.cache),
		"pattern_cache_entries": len(tc.patternCache),
		"memory_usage_bytes":    tc.currentMemory,
		"memory_usage_mb":       float64(tc.currentMemory) / (1024 * 1024),
		"memory_limit_mb":       float64(tc.memoryLimit) / (1024 * 1024),
		"max_size":              tc.maxSize,
		"ttl_seconds":           tc.ttl,
	}
}

// evictOldestEntries は最も古いエントリを削除します
func (tc *TypeCache) evictOldestEntries() {
	if len(tc.cache) == 0 {
		return
	}

	// 最も古いエントリを見つける
	var oldestKey string
	var oldestTime int64 = time.Now().Unix()

	for key, entry := range tc.cache {
		if entry.Timestamp < oldestTime {
			oldestTime = entry.Timestamp
			oldestKey = key
		}
	}

	// 削除
	if oldestKey != "" {
		if entry, exists := tc.cache[oldestKey]; exists {
			estimatedSize := tc.estimateTypeInfoSize(entry.TypeInfo)
			delete(tc.cache, oldestKey)
			tc.currentMemory -= estimatedSize
			if tc.currentMemory < 0 {
				tc.currentMemory = 0
			}
		}
	}
}

// evictOldestPatternEntries は最も古いパターンエントリを削除します
func (tc *TypeCache) evictOldestPatternEntries() {
	if len(tc.patternCache) == 0 {
		return
	}

	// 最も古いエントリを見つける
	var oldestKey string
	var oldestTime int64 = time.Now().Unix()

	for key, entry := range tc.patternCache {
		if entry.Timestamp < oldestTime {
			oldestTime = entry.Timestamp
			oldestKey = key
		}
	}

	// 削除
	if oldestKey != "" {
		if entry, exists := tc.patternCache[oldestKey]; exists {
			estimatedSize := tc.estimatePatternSize(entry.Patterns)
			delete(tc.patternCache, oldestKey)
			tc.currentMemory -= estimatedSize
			if tc.currentMemory < 0 {
				tc.currentMemory = 0
			}
		}
	}
}

// estimateTypeInfoSize は型情報のメモリ使用量を推定します（簡易版）
func (tc *TypeCache) estimateTypeInfoSize(typeInfo map[string]EnhancedTypeInfo) int64 {
	var size int64

	for typeName, info := range typeInfo {
		// 型名のサイズ
		size += int64(len(typeName))

		// フィールド情報のサイズ
		for fieldName, field := range info.Fields {
			size += int64(len(fieldName))
			size += int64(len(field.Name))
			size += int64(len(field.Type))
			size += int64(len(field.JSONTag))
			size += int64(len(field.NestedType))

			// エイリアスのサイズ
			for _, alias := range field.Aliases {
				size += int64(len(alias))
			}
		}

		// その他のフィールド
		size += int64(len(info.HTTPMethod))
		size += int64(len(info.ResourceName))
	}

	// 構造体のオーバーヘッドを追加（概算）
	size += int64(len(typeInfo)) * 100

	return size
}

// estimatePatternSize はパターンのメモリ使用量を推定します
func (tc *TypeCache) estimatePatternSize(patterns []TypeNamePattern) int64 {
	var size int64

	for _, pattern := range patterns {
		size += int64(len(pattern.Pattern))
		size += int64(len(pattern.Category))
		size += 16 // Priority用
	}

	// 構造体のオーバーヘッドを追加
	size += int64(len(patterns)) * 50

	return size
}

// PrintStats はキャッシュの統計情報を出力します
func (tc *TypeCache) PrintStats() {
	stats := tc.GetStats()
	fmt.Println("\n=== キャッシュ統計情報 ===")
	fmt.Printf("型キャッシュエントリ数: %d\n", stats["type_cache_entries"])
	fmt.Printf("パターンキャッシュエントリ数: %d\n", stats["pattern_cache_entries"])
	fmt.Printf("メモリ使用量: %.2f MB / %.2f MB\n",
		stats["memory_usage_mb"],
		stats["memory_limit_mb"])
	fmt.Printf("最大サイズ: %d\n", stats["max_size"])
	fmt.Printf("TTL: %d 秒\n", stats["ttl_seconds"])
	fmt.Println("========================")
}
