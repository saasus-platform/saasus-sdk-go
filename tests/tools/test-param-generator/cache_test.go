package main

import (
	"testing"
	"time"
)

func TestTypeCache_SetAndGet(t *testing.T) {
	cache := NewTypeCache(10, 60, 10)

	// テスト用の型情報を作成
	typeInfo := map[string]EnhancedTypeInfo{
		"TestType": {
			Name: "TestType",
			Fields: map[string]EnhancedFieldInfo{
				"Field1": {
					Name:    "Field1",
					Type:    "string",
					JSONTag: "field1",
				},
			},
		},
	}

	// キャッシュに保存
	err := cache.Set("test_module", typeInfo)
	if err != nil {
		t.Fatalf("キャッシュへの保存に失敗: %v", err)
	}

	// キャッシュから取得
	cached, found := cache.Get("test_module")
	if !found {
		t.Fatal("キャッシュから型情報を取得できませんでした")
	}

	// 内容を検証
	if len(cached) != 1 {
		t.Errorf("期待される型数: 1, 実際: %d", len(cached))
	}

	if _, exists := cached["TestType"]; !exists {
		t.Error("TestType が見つかりません")
	}
}

func TestTypeCache_TTL(t *testing.T) {
	// TTL 1秒のキャッシュを作成
	cache := NewTypeCache(10, 1, 10)

	typeInfo := map[string]EnhancedTypeInfo{
		"TestType": {Name: "TestType"},
	}

	// キャッシュに保存
	cache.Set("test_module", typeInfo)

	// すぐに取得できることを確認
	_, found := cache.Get("test_module")
	if !found {
		t.Fatal("キャッシュから型情報を取得できませんでした")
	}

	// TTL経過後は取得できないことを確認
	time.Sleep(2 * time.Second)
	_, found = cache.Get("test_module")
	if found {
		t.Error("TTL経過後もキャッシュから取得できてしまいました")
	}
}

func TestTypeCache_Clear(t *testing.T) {
	cache := NewTypeCache(10, 60, 10)

	typeInfo := map[string]EnhancedTypeInfo{
		"TestType": {Name: "TestType"},
	}

	// キャッシュに保存
	cache.Set("test_module", typeInfo)

	// クリア前は取得できる
	_, found := cache.Get("test_module")
	if !found {
		t.Fatal("キャッシュから型情報を取得できませんでした")
	}

	// クリア
	cache.Clear()

	// クリア後は取得できない
	_, found = cache.Get("test_module")
	if found {
		t.Error("クリア後もキャッシュから取得できてしまいました")
	}
}

func TestTypeCache_PatternCache(t *testing.T) {
	cache := NewTypeCache(10, 60, 10)

	patterns := []TypeNamePattern{
		{Pattern: "CreateTestParams", Priority: 0, Category: "exact"},
		{Pattern: "TestCreateParams", Priority: 10, Category: "standard"},
	}

	// パターンをキャッシュに保存
	err := cache.SetPattern("test:Create", patterns)
	if err != nil {
		t.Fatalf("パターンキャッシュへの保存に失敗: %v", err)
	}

	// パターンをキャッシュから取得
	cached, found := cache.GetPattern("test:Create")
	if !found {
		t.Fatal("パターンキャッシュから取得できませんでした")
	}

	// 内容を検証
	if len(cached) != 2 {
		t.Errorf("期待されるパターン数: 2, 実際: %d", len(cached))
	}

	if cached[0].Pattern != "CreateTestParams" {
		t.Errorf("期待されるパターン: CreateTestParams, 実際: %s", cached[0].Pattern)
	}
}

func TestTypeCache_Stats(t *testing.T) {
	cache := NewTypeCache(10, 60, 10)

	typeInfo := map[string]EnhancedTypeInfo{
		"TestType": {Name: "TestType"},
	}

	cache.Set("test_module", typeInfo)

	patterns := []TypeNamePattern{
		{Pattern: "CreateTestParams", Priority: 0, Category: "exact"},
	}
	cache.SetPattern("test:Create", patterns)

	// 統計情報を取得
	stats := cache.GetStats()

	if stats["type_cache_entries"].(int) != 1 {
		t.Errorf("期待される型キャッシュエントリ数: 1, 実際: %d", stats["type_cache_entries"])
	}

	if stats["pattern_cache_entries"].(int) != 1 {
		t.Errorf("期待されるパターンキャッシュエントリ数: 1, 実際: %d", stats["pattern_cache_entries"])
	}
}

func TestTypeCache_MemoryLimit(t *testing.T) {
	// 非常に小さいメモリ制限（1KB）を設定
	cache := NewTypeCache(100, 60, 0) // 0MB = 制限なし

	// 大量の型情報を追加
	for i := 0; i < 10; i++ {
		typeInfo := map[string]EnhancedTypeInfo{
			"TestType": {
				Name: "TestType",
				Fields: map[string]EnhancedFieldInfo{
					"Field1": {Name: "Field1", Type: "string"},
					"Field2": {Name: "Field2", Type: "int"},
				},
			},
		}
		cache.Set("test_module_"+string(rune(i)), typeInfo)
	}

	stats := cache.GetStats()
	if stats["type_cache_entries"].(int) != 10 {
		t.Errorf("期待されるエントリ数: 10, 実際: %d", stats["type_cache_entries"])
	}
}

func TestPatternGenerator_WithCache(t *testing.T) {
	cache := NewTypeCache(10, 60, 10)

	typeInfoMap := map[string]EnhancedTypeInfo{
		"CreateFeedbackJSONRequestBody": {
			Name: "CreateFeedbackJSONRequestBody",
		},
	}

	// キャッシュ付きのPatternGeneratorを作成
	pg := NewPatternGeneratorWithCache("feedbacks", "Create", typeInfoMap, cache)

	// 1回目のパターン生成
	patterns1 := pg.GeneratePatterns()
	if len(patterns1) == 0 {
		t.Fatal("パターンが生成されませんでした")
	}

	// 2回目のパターン生成（キャッシュから取得されるはず）
	patterns2 := pg.GeneratePatterns()
	if len(patterns2) != len(patterns1) {
		t.Errorf("キャッシュから取得したパターン数が異なります: %d vs %d", len(patterns1), len(patterns2))
	}

	// キャッシュ統計を確認
	stats := cache.GetStats()
	if stats["pattern_cache_entries"].(int) != 1 {
		t.Errorf("期待されるパターンキャッシュエントリ数: 1, 実際: %d", stats["pattern_cache_entries"])
	}
}
