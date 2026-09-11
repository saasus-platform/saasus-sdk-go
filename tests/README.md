# SaaSus SDK Go - Tests

テストフレームワークとE2Eテスト

## ディレクトリ構造

```
tests/
├── testlib/              # 共通テストライブラリ
├── tools/                # 分析・レポートツール
└── e2e/                  # E2Eテスト
    ├── billingapi/       # Billing API E2Eテスト
    ├── integrationapi/   # Integration API E2Eテスト
    └── snapshot/         # スナップショットテスト
```

## クイックスタート

```bash
# Billing API E2Eテスト
cd tests/e2e/billingapi
go test -v

# Integration API E2Eテスト
cd tests/e2e/integrationapi
go test -v

# スナップショットテスト
cd tests/e2e/billingapi
go test -v -run TestBillingAPIWithSnapshot -snapshot-mode=capture
```

## 詳細ドキュメント

- [E2E Testing Framework](./e2e/README.md)
- [共通テストライブラリ](./testlib/README.md)
- [分析ツール](./tools/README.md)
