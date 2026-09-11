# SaaSus SDK Go - E2E Testing Framework

E2E テストフレームワーク - 全 SaaSus API モジュール対応の統合テストシステム

## 概要

実際の API 環境で SDK の動作を検証し、リリース間の互換性を保証するテストシステムです。

## クイックスタート

### 環境設定

`.env`ファイルを作成：

```bash
SAASUS_SAAS_ID=your-saas-id
SAASUS_API_KEY=your-api-key
SAASUS_SECRET_KEY=your-secret-key
STRIPE_SECRET_KEY=your-stripe-secret-key
```

### テスト実行

#### 通常のE2Eテスト

```bash
# Billing API E2Eテスト
cd tests/e2e/billingapi
go test -v -run TestBillingAPIE2E

# Pricing API E2Eテスト
cd tests/e2e/pricingapi
go test -v -run TestPricingAPIE2E
```

#### スナップショットテスト

```bash
# スナップショット取得
cd tests/e2e/pricingapi
go test -v -run TestPricingAPIWithSnapshot -snapshot-mode=capture

# スナップショット比較
go test -v -run TestPricingAPIWithSnapshot -snapshot-mode=compare

# 統合実行（取得 + 比較 + レポート）
go test -v -run TestPricingAPIWithSnapshot -snapshot-mode=full

# Billing APIスナップショットテスト
cd tests/e2e/billingapi
go test -v -run TestBillingAPIWithSnapshot -snapshot-mode=capture
```

## テスト機能

### 通常のE2Eテスト

- **ストーリーベーステスト**: 実際のユーザーシナリオに基づくテスト
- **カバレッジ追跡**: SDK メソッドの実行状況監視
- **自動パラメータ生成**: test-param-generator による testdata 自動生成

### スナップショットテスト

- **スナップショット取得**: ストーリー実行の詳細な記録（リクエスト・レスポンス・状態変化）
- **互換性チェック**: 前回リリースとの自動差分検出
- **レポート生成**: 差分の可視化と影響範囲分析
- **検証ルール**: 完了ステータス、シーケンス、状態遷移の検証

## ディレクトリ構造

```
tests/e2e/
├── billingapi/                  # Billing API E2Eテスト
│   ├── billing_test.go          # 通常のE2Eテスト
│   ├── billing_snapshot_test.go # スナップショットテスト
│   ├── stories.go
│   ├── validation.go
│   ├── story_snapshot_config.go # スナップショット設定
│   └── testdata/
├── pricingapi/                  # Pricing API E2Eテスト
│   ├── pricing_test.go          # 通常のE2Eテスト
│   ├── pricing_snapshot_test.go # スナップショットテスト
│   ├── stories.go               # 6ストーリー（動的メソッド名変換）
│   ├── validation.go            # 122メソッド対応
│   ├── story_snapshot_config.go # スナップショット設定
│   └── testdata/
├── snapshot/                    # スナップショットテストシステム
│   ├── billingapi/              # Billing APIスナップショット出力
│   │   ├── story_snapshots/
│   │   ├── story_comparisons/
│   │   ├── story_validations/
│   │   └── story_reports/
│   ├── pricingapi/              # Pricing APIスナップショット出力
│   │   ├── story_snapshots/
│   │   ├── story_comparisons/
│   │   ├── story_validations/
│   │   └── story_reports/
│   └── docs/                    # スナップショットテストドキュメント
├── SNAPSHOT_TEST_IMPLEMENTATION_GUIDE.md  # スナップショットテスト実装ガイド
└── API_TEMPLATE.md              # 新API実装テンプレート
```

## 実装済みモジュール

### Billing API

- **ストーリー数**: 4
- **メソッド数**: 24
- **カバレッジ**: 100%
- **詳細**: [billingapi/README.md](./billingapi/README.md)

### Pricing API

- **ストーリー数**: 6（Postman Collection系4 + Individual Operations系2）
- **ステップ数**: 30/ストーリー（Postman Collection系）、14/ストーリー（Individual Operations系）
- **メソッド数**: 122（ReturnInternalServerError 除く）
- **カバレッジ**: 100%
- **機能**: Pricing Units, Menus, Plans, Metering, Tax Rate, Stripe 連携
- **特徴**: 
  - WithResponse ストーリーは Standard ストーリーから動的にメソッド名を変換
  - 全メソッドバリアント対応（Struct wrapper, WithBody, WithResponse, WithBodyWithResponse）
  - スナップショットテスト完全対応
- **詳細**: [pricingapi/README.md](./pricingapi/README.md)

## 関連ドキュメント

- [Billing API E2E Tests](./billingapi/README.md)
- [Pricing API E2E Tests](./pricingapi/README.md)
- [Snapshot Testing System](./snapshot/README.md)
- [API 実装テンプレート](./API_TEMPLATE.md)
- [共通テストライブラリ](../testlib/README.md)

## 実装状況

- **完了**: Billing API, Pricing API (2/7 APIs)
- **次の優先度**: Auth API → Communication API
- **テンプレート**: [API_TEMPLATE.md](./API_TEMPLATE.md)
