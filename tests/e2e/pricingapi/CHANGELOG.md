# Pricing API E2E Tests - 変更履歴

## 2024年11月 - スナップショットキャプチャ処理の簡素化

### 概要

`tests/testlib/snapshot/capture.go`の`CaptureStepResponse`メソッドを強化し、HTTPレスポンスボディの読み取りと再設定処理を追加しました。これにより、StateUpdateでレスポンスボディを再読み取り可能になります。

### 変更内容

**変更箇所**: `tests/testlib/snapshot/capture.go`の`captureHTTPResponse`メソッド

**追加されたコード**:
```go
// Parse JSON body if available
if httpResp.Body != nil {
    bodyBytes, err := io.ReadAll(httpResp.Body)
    if err == nil {
        _ = httpResp.Body.Close()

        // Replace body with a new reader so it can be read again by StateUpdate
        httpResp.Body = io.NopCloser(bytes.NewReader(bodyBytes))

        // Try to parse as JSON
        var jsonData map[string]interface{}
        if err := json.Unmarshal(bodyBytes, &jsonData); err == nil {
            // ...
            returnValue.JSONData = maskedMap
        }
    }
}
```

### 改善効果

1. **レスポンスボディの再読み取りをサポート**: ボディを読み取った後、新しいReaderで置き換えることで、StateUpdateでの再読み取りを可能に
2. **JSONパースの自動化**: レスポンスボディがJSONの場合、自動的にパースしてスナップショットに含める
3. **機密情報のマスキング**: パースされたJSONデータに対してもマスキングを適用
4. **エンジンとの統合**: `tests/testlib/snapshot/engine.go`がStateUpdateにキャプチャ済みJSONDataを渡すため、ボディ操作が不要（推奨パターン）

### 関連ドキュメント更新

- `tests/e2e/SNAPSHOT_TEST_IMPLEMENTATION_GUIDE.md`: トラブルシューティングセクションを更新
- `tests/testlib/README.md`: StateUpdate実装パターンの説明を更新

### 影響範囲

- すべてのスナップショットテスト（pricing、billing等）
- StateUpdate関数の実装パターン

---

## 2024年11月 - Pricing Plan更新時の名前フィールド修正

### 概要

`updatePricingPlanParams`関数で、Pricing Plan更新時に元の名前を保持するように修正しました。これにより、Pricing Unit、Pricing Menu、Pricing Planすべてで名前フィールドの扱いが統一されました。

### 問題の背景

Pricing APIでは、Pricing Planの`name`フィールドは作成後に変更できません。しかし、以前の実装では更新時に`test_params.json`の`plan_name`を使用していたため、API仕様に違反していました。

### 変更内容

**変更箇所**: `tests/e2e/pricingapi/validation.go`の`updatePricingPlanParams`関数

**変更前**:
```go
func updatePricingPlanParams(variables map[string]any) any {
    pricingPlanID := mustGetStringVariable(variables, "pricing_plan_id")
    params := loadTestParams()
    updateParams := params.UpdatePricingPlan.UpdateParams
    
    body := map[string]any{
        "name":         updateParams["plan_name"],  // ❌ 新しい名前を設定しようとする
        "display_name": updateParams["plan_display_name"],
        "description":  updateParams["plan_description"],
        "menu_ids":     []string{menuID},
    }
}
```

**変更後**:
```go
func updatePricingPlanParams(variables map[string]any) any {
    pricingPlanID := mustGetStringVariable(variables, "pricing_plan_id")
    planName := mustGetStringVariable(variables, "pricing_plan_name")  // ✅ 元の名前を取得
    params := loadTestParams()
    updateParams := params.UpdatePricingPlan.UpdateParams
    
    body := map[string]any{
        "name":         planName,  // ✅ 元の名前を使用（変更不可）
        "display_name": updateParams["plan_display_name"],
        "description":  updateParams["plan_description"],
        "menu_ids":     []string{menuID},
    }
}
```

### 影響

1. **API仕様への準拠**: Pricing Planの`name`フィールドが変更不可という制約に正しく対応
2. **一貫性の向上**: Pricing Unit、Pricing Menu、Pricing Planで名前フィールドの扱いが統一
3. **テストの安定性**: 更新操作が正しく動作することを検証
4. **エラーの防止**: 名前変更によるAPIエラーを回避

### 統一された実装パターン

**Pricing Unit**:
```go
// 作成時
variables["pricing_unit_name"] = unitName
// 更新時
unitName := mustGetStringVariable(variables, "pricing_unit_name_2")
```

**Pricing Plan**:
```go
// 作成時
variables["pricing_plan_name"] = planName
// 更新時
planName := mustGetStringVariable(variables, "pricing_plan_name")
```

**Pricing Menu**:
```go
// 作成時
variables["pricing_menu_name"] = menuName
// 更新時
menuName := mustGetStringVariable(variables, "pricing_menu_name")
```

### ベストプラクティス

Pricing Plan、Unit、Menuを更新する際は、以下の点に注意してください：

```go
// ✅ 正しい実装
func updateResourceParams(variables map[string]any) any {
    resourceName := mustGetStringVariable(variables, "resource_name")  // 元の名前
    
    body := map[string]any{
        "name":         resourceName,                    // 変更不可
        "display_name": updateParams["display_name"],    // 更新可能
        "description":  updateParams["description"],     // 更新可能
    }
}

// ❌ 誤った実装
body := map[string]any{
    "name": "new-name",  // これはエラーになる
}
```

### 関連する変更

- Pricing Unit更新時の名前フィールド修正（2024年12月）
- Pricing Unit名前管理の改善（2024年12月）

### 関連ドキュメント

- [tests/e2e/pricingapi/README.md](README.md#重要な注意点---pricing-unitplanmenu更新時の名前フィールド) - 名前フィールドの制約に関する詳細
- [tests/e2e/pricingapi/validation.go](validation.go) - 修正されたコード

## 2024年11月 - Pricing Unit名前管理の改善

### 概要

`createPricingUnitStructParams`関数で、生成されたユニット名を`variables`に保存するように改善しました。これにより、更新操作時に元の名前を正確に参照できるようになりました。

### 変更内容

**変更箇所**: `tests/e2e/pricingapi/validation.go`の`createPricingUnitStructParams`関数

**変更前**:
```go
func createPricingUnitStructParams(variables map[string]any) any {
    // ...
    body := map[string]any{
        "type":               createParams["type"],
        "name":               generateUniqueName(createParams["name"]),  // 生成するだけ
        "display_name":       createParams["display_name"],
        // ...
    }
}
```

**変更後**:
```go
func createPricingUnitStructParams(variables map[string]any) any {
    // ...
    // Generate unique name and store it for later use
    unitName := generateUniqueName(createParams["name"])
    variables["pricing_unit_name"] = unitName  // ✅ 名前を保存

    body := map[string]any{
        "type":               createParams["type"],
        "name":               unitName,  // 保存した名前を使用
        "display_name":       createParams["display_name"],
        // ...
    }
}
```

### 影響

1. **更新操作の正確性**: `updatePricingUnitStructParams`で元の名前を正確に参照可能
2. **エラーの防止**: 名前フィールドの変更不可制約に正しく対応
3. **テストの安定性**: 名前の一貫性が保証される

### 関連する変更

この変更は、以前の「Pricing Unit更新時の名前フィールド修正」と組み合わせて動作します：

1. **作成時** (`createPricingUnitStructParams`): 生成された名前を`variables`に保存
2. **更新時** (`updatePricingUnitStructParams`): 保存された名前を`variables`から取得

### ベストプラクティス

リソース作成時に後で参照が必要な値は、必ず`variables`に保存してください：

```go
// ✅ 推奨パターン
func createResourceParams(variables map[string]any) any {
    uniqueName := generateUniqueName(baseName)
    variables["resource_name"] = uniqueName  // 後で使用するために保存
    
    return ResourceBody{
        Name: uniqueName,
        // ...
    }
}
```

## 2024年11月 - Pricing Unit更新時の名前フィールド修正

### 概要

`updatePricingUnitStructParams`関数で、Pricing Unit更新時に元の名前を保持するように修正しました。

### 問題の背景

Pricing APIでは、Pricing Unitの`name`フィールドは作成後に変更できません。しかし、以前の実装では更新時に`createParams`から新しい名前を取得しようとしていたため、API仕様に違反していました。

### 変更内容

**変更前**:
```go
func updatePricingUnitStructParams(variables map[string]any) any {
    pricingUnitID := mustGetStringVariable(variables, "pricing_unit_id_2")
    params := loadTestParams()
    createParams := params.CreatePricingUnit.CreateParams
    
    fixedUnit := pricingapi.PricingFixedUnitForSave{
        Name: createParams["name1"].(string), // ❌ 新しい名前を設定しようとする
        // ...
    }
}
```

**変更後**:
```go
func updatePricingUnitStructParams(variables map[string]any) any {
    pricingUnitID := mustGetStringVariable(variables, "pricing_unit_id_2")
    unitName := mustGetStringVariable(variables, "pricing_unit_name_2")  // ✅ 元の名前を取得
    params := loadTestParams()
    
    fixedUnit := pricingapi.PricingFixedUnitForSave{
        Name: unitName, // ✅ 元の名前を使用（変更不可）
        // ...
    }
}
```

### 影響

1. **API仕様への準拠**: Pricing Unitの`name`フィールドが変更不可という制約に正しく対応
2. **テストの正確性**: 更新操作が正しく動作することを検証
3. **エラーの防止**: 名前変更によるAPIエラーを回避

### 関連する変更

- `stories.go`の`CreatePricingUnit`ステップで、作成時の名前を`pricing_unit_name_2`として`variables`に保存
- `updatePricingUnitStructParams`で、保存された名前を使用して更新

### ベストプラクティス

Pricing Unitを更新する際は、以下の点に注意してください：

```go
// ✅ 正しい実装
func updatePricingUnitStructParams(variables map[string]any) any {
    unitName := mustGetStringVariable(variables, "pricing_unit_name_2")  // 元の名前
    
    fixedUnit := pricingapi.PricingFixedUnitForSave{
        Name:        unitName,                                    // 変更不可
        DisplayName: updateParams["display_name"].(string),       // 更新可能
        Description: updateParams["description"].(string),        // 更新可能
        // ...
    }
}

// ❌ 誤った実装
fixedUnit := pricingapi.PricingFixedUnitForSave{
    Name: "new-name",  // これはエラーになる
    // ...
}
```

### 関連ドキュメント

- [tests/e2e/pricingapi/README.md](README.md#問題-3-パラメータ型変換エラー) - トラブルシューティングセクションに詳細を追加
- [tests/e2e/pricingapi/validation.go](validation.go) - 修正されたコード

## 2024年11月 - Individual Operations Storyの拡張

### 概要

Individual Operations Story（Standard / WithResponse）に8つの新しい更新メソッドのテストステップを追加し、ステップ数を9から14に拡張しました。

### 主な変更内容

#### 1. 新規追加されたテストステップ（Step 6-13）

**更新メソッドの網羅的なカバレッジ**:
- Step 6: UpdatePricingUnit（Struct wrapper）
- Step 7: UpdatePricingMenu（Struct wrapper）
- Step 8: UpdatePricingPlan（Struct wrapper）
- Step 9: UpdatePricingPlansUsed（Struct wrapper）
- Step 10: UpdateMeteringUnitByID（Struct wrapper）
- Step 11: UpdateMeteringUnitTimestampCount（Struct wrapper）
- Step 12: UpdateMeteringUnitTimestampCountNow（Struct wrapper）
- Step 13: UpdateTaxRate（Struct wrapper）

#### 2. テストフローの再構成

**3フェーズ構造**:
1. **作成フェーズ** (Step 1-5): 全リソースの作成
2. **更新フェーズ** (Step 6-13): 全更新メソッドのテスト
3. **クリーンアップフェーズ** (Step 14): 一括削除

#### 3. 影響を受けたストーリー

- Individual Operations Story - Standard Methods: 9 → 14ステップ
- Individual Operations Story - WithResponse Methods: 9 → 14ステップ

### 利点

1. **カバレッジの向上**: 全ての主要な更新メソッド（Struct wrapper版）を網羅的にテスト
2. **型安全性の確認**: Struct wrapperバリアントの動作を詳細に検証
3. **問題の早期発見**: 更新操作の問題を個別に特定可能
4. **テストの粒度向上**: より細かい単位でのテストが可能に

### 技術的詳細

**実装パターン**:
```go
// 各更新メソッドは対応するパラメータ関数を使用
{
    Name:           "Additional6: Update Pricing Unit (UpdatePricingUnit)",
    ClientMethod:   "UpdatePricingUnit",
    Parameters:     updatePricingUnitStructParams,
    ExpectedStatus: 200,
}
```

**パラメータ関数の命名規則**:
- `update{Resource}StructParams`: Struct wrapperバリアント用
- 例: `updatePricingUnitStructParams`, `updateTaxRateStructParams`

## 2024年11月 - validation.goのコード品質改善

### 概要

パラメータ生成関数の実装を改善し、型安全性とパフォーマンスを向上させました。

### 主な変更内容

#### 1. 構造体の直接作成

**変更前**: JSON経由のマーシャル/アンマーシャル
```go
bodyData := map[string]any{
    "name":             generateUniqueName(createParams["menu_name"]),
    "display_name":     createParams["menu_display_name"],
    "description":      createParams["menu_description"],
    "pricing_unit_ids": []string{pricingUnitID},
}
bodyBytes, _ := json.Marshal(bodyData)
var body pricingapi.CreatePricingMenuJSONRequestBody
_ = json.Unmarshal(bodyBytes, &body)
return body
```

**変更後**: 構造体を直接作成
```go
body := pricingapi.CreatePricingMenuJSONRequestBody{
    Name:        generateUniqueName(createParams["menu_name"]),
    DisplayName: createParams["menu_display_name"].(string),
    Description: createParams["menu_description"].(string),
    UnitIds:     []pricingapi.Uuid{pricingapi.Uuid(pricingUnitID)},
}
return body
```

#### 2. 影響を受けた関数

- `createPricingMenuStructParams`: JSON変換を削除、構造体を直接作成
- `createPricingPlanStructParams`: JSON変換を削除、構造体を直接作成
- `updatePricingPlansUsedStructParams`: 不要な構造体ラッパーを削除

#### 3. コードフォーマットの改善

- インデントの統一
- コメントの整理

### 利点

1. **パフォーマンス向上**
   - JSON経由の変換処理が不要になり、実行速度が向上
   - メモリ使用量の削減

2. **型安全性の向上**
   - コンパイル時に型チェックが行われる
   - 実行時エラーのリスクが減少

3. **可読性の向上**
   - コードが簡潔になり、意図が明確に
   - デバッグが容易に

4. **メンテナンス性の向上**
   - 型定義の変更が即座に検出される
   - リファクタリングが安全に

### 影響範囲

- **テスト動作**: 変更なし（すべてのテストは引き続き正常に動作）
- **API互換性**: 影響なし（内部実装の改善のみ）
- **パフォーマンス**: 向上（JSON変換のオーバーヘッドを削減）

### ベストプラクティス

今後のパラメータ関数実装では、以下のパターンを推奨します：

```go
// ✅ 推奨: 構造体を直接作成
func createResourceParams(variables map[string]any) any {
    return api.CreateResourceJSONRequestBody{
        Name:        "Resource Name",
        Description: "Resource Description",
        Tags:        []string{"tag1", "tag2"},
    }
}

// ❌ 非推奨: JSON経由の変換
func createResourceParamsOld(variables map[string]any) any {
    bodyData := map[string]any{...}
    bodyBytes, _ := json.Marshal(bodyData)
    var body api.CreateResourceJSONRequestBody
    _ = json.Unmarshal(bodyBytes, &body)
    return body
}
```

### 関連ドキュメント

- [tests/e2e/pricingapi/README.md](README.md) - Pricing API E2Eテストの概要
- [tests/e2e/SNAPSHOT_TEST_IMPLEMENTATION_GUIDE.md](../SNAPSHOT_TEST_IMPLEMENTATION_GUIDE.md) - E2Eテスト実装ガイド
- [tests/e2e/pricingapi/validation.go](validation.go) - パラメータ関数の実装

### 今後の改善予定

- [ ] 他のパラメータ関数も同様のパターンに統一
- [ ] 型アサーションのエラーハンドリングを強化
- [ ] テストパラメータの型定義を強化
