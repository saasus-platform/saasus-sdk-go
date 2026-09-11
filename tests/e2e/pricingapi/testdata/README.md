# Pricing API テストパラメータ

このディレクトリには、Pricing API E2Eテストで使用するテストパラメータが含まれています。

## ファイル

### test_params.json

test-param-generatorツールによって自動生成されたテストパラメータファイルです。

#### 構造

```json
{
  "CreatePricingUnit": {
    "createParams": {
      "type": "tiered_usage",
      "name": "tiered_usage_unit-name",
      "display_name": "段階使用量ユニット",
      "description": "段階使用量ユニットdescription",
      "currency": "JPY",
      "upper_count": 5000,
      "metering_unit_name": "tiered_unit_metter",
      "aggregate_usage": "max",
      "up_to": 0,
      "unit_amount": 500,
      "flat_amount": 400,
      "inf": true
    }
  },
  "UpdatePricingUnit": {
    "updateParams": {
      // 更新用パラメータ
    }
  },
  // ... 他のメソッドのパラメータ
}
```

#### 主要なパラメータグループ

1. **CreatePricingUnit / UpdatePricingUnit**
   - Pricing Unit の作成・更新パラメータ
   - type: ユニットタイプ（tiered_usage, tiered, usage, fixed）
   - currency: 通貨（JPY, USD）
   - metering_unit_name: 関連するメータリングユニット名

2. **CreatePricingMenu / UpdatePricingMenu**
   - Pricing Menu の作成・更新パラメータ
   - menu_name, menu_display_name, menu_description
   - unit_ids: 含まれる Pricing Unit の ID リスト

3. **CreatePricingPlan / UpdatePricingPlan**
   - Pricing Plan の作成・更新パラメータ
   - plan_name, plan_display_name, plan_description
   - menu_ids: 含まれる Pricing Menu の ID リスト

4. **CreateMeteringUnit / UpdateMeteringUnitByID**
   - Metering Unit の作成・更新パラメータ
   - unit_name, display_name, description
   - aggregate_usage: 集計方法（sum, max）

5. **UpdateMeteringUnitTimestampCount / UpdateMeteringUnitTimestampCountNow**
   - メータリングカウントの更新パラメータ
   - method: 更新方法（add, sub, direct）
   - count: カウント値

6. **CreateTaxRate / UpdateTaxRate**
   - Tax Rate の作成・更新パラメータ
   - name, display_name, description
   - percentage: 税率（パーセント）
   - inclusive: 内税かどうか
   - country: 国コード（ISO 3166-1 alpha-2）

7. **UpdatePricingPlansUsed**
   - Pricing Plan の使用済みフラグ更新パラメータ
   - plan_ids: 更新対象の Plan ID リスト

## 使用方法

### validation.go での読み込み

```go
params := loadTestParams()

// Pricing Unit 作成パラメータの取得
createParams := params.CreatePricingUnit.CreateParams
unitName := createParams["name"].(string)
displayName := createParams["display_name"].(string)

// Pricing Menu 作成パラメータの取得
menuParams := params.CreatePricingMenu.CreateParams
menuName := menuParams["menu_name"].(string)
```

### パラメータのカスタマイズ

テストパラメータをカスタマイズする場合は、test_params.json を直接編集してください。

**注意**: このファイルは test-param-generator によって自動生成されるため、再生成時に上書きされる可能性があります。カスタマイズした内容を保持したい場合は、バックアップを取るか、別のファイルに保存してください。

## 関連ファイル

- `../validation.go` - パラメータ関数の実装
- `../stories.go` - テストストーリーの定義
- `../../tools/test-param-generator/` - パラメータ生成ツール

## トラブルシューティング

### パラメータの型エラー

```
panic: interface conversion: interface {} is nil, not string
```

**原因**: test_params.json に期待するキーが存在しない、または型が異なる

**解決策**:
1. test_params.json の構造を確認
2. validation.go のパラメータアクセスコードを修正
3. 型アサーション時にデフォルト値を設定

### パラメータの欠落

```
panic: failed to read test params: no such file or directory
```

**原因**: test_params.json ファイルが存在しない

**解決策**:
1. test-param-generator を実行してファイルを生成
2. testdata ディレクトリが正しい場所にあるか確認
3. ファイルパスが正しいか確認（`testdata/test_params.json`）
