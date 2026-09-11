# Billing API E2E Tests

**SaaSus Billing API専用E2Eテスト** - testlibを使用したシンプルなストーリーベーステスト

## 🚀 クイックスタート

```bash
# 基本実行（プロジェクトルートまたはこのディレクトリから）
go test -v

# 詳細ログ付き実行
LOG_LEVEL=debug go test -v

# タイムアウト延長
go test -v -timeout 30m
```

## 📁 ファイル構成

```
tests/e2e/billingapi/
├── billing_test.go              # メインE2Eテスト (TestBillingAPIE2E)
├── stories.go                   # ストーリー定義とメソッド一覧
├── validation.go                # バリデーション・ヘルパー関数
├── story_snapshot_config.go    # スナップショット設定（将来使用）
└── README.md                    # このドキュメント
```

**シンプルな構成:**
- テスト関数は1つ（`TestBillingAPIE2E`）のみ
- すべてtestlib経由で実行
- 統合レイヤーなし

## 🎯 テストストーリー

### 1. Postman Collection Story - Standard Methods
Postmanコレクションのフローを標準メソッドで再現

**6ステップ:**
1. Pre_GetStripeConnectionInformation - 初期状態確認
2. UpdateStripeConnectionInfo - Stripe設定更新（型付き）
3. UpdateStripeConnectionInfo_WithBody - Stripe設定更新（任意ボディ）
4. DeleteStripeConnectionInfo - Stripe統合削除
5. RecreateStripeConnectionInfo - Stripe統合再作成
6. Final_GetStripeConnectionInformation - 最終状態確認

### 2. Postman Collection Story - WithResponse Methods
同じフローをWithResponseメソッドで実行

**6ステップ:**
- 上記と同じフロー、すべてWithResponseメソッドを使用

## 📊 メソッドカバレッジ

**全8メソッドの100%カバレッジ達成**

| メソッド | 用途 | テストストーリー |
|---------|------|-----------------|
| `GetStripeInfo` | Stripe統合状況取得 | Standard |
| `UpdateStripeInfo` | Stripe設定更新（型付き） | Standard |
| `UpdateStripeInfoWithBody` | Stripe設定更新（任意ボディ） | Standard |
| `DeleteStripeInfo` | Stripe統合削除 | Standard |
| `GetStripeInfoWithResponse` | レスポンス解析付き取得 | WithResponse |
| `UpdateStripeInfoWithResponse` | レスポンス解析付き更新（型付き） | WithResponse |
| `UpdateStripeInfoWithBodyWithResponse` | レスポンス解析付き更新（任意） | WithResponse |
| `DeleteStripeInfoWithResponse` | レスポンス解析付き削除 | WithResponse |

## ⚙️ 設定

### 環境変数（必須）

プロジェクトルートに`.env`ファイルを作成：

```bash
# SaaSus認証情報（必須）
SAASUS_SAAS_ID=your-saas-id
SAASUS_API_KEY=your-api-key
SAASUS_SECRET_KEY=your-secret-key

# Stripe認証情報（Billing APIテスト用、必須）
STRIPE_SECRET_KEY=sk_test_your_stripe_test_key

# オプション設定
SAASUS_BASE_URL=https://api.saasus.io/v1/billing  # デフォルトは本番環境
LOG_LEVEL=info                                     # error/warn/info/debug
```

### ログレベル

```bash
# デバッグログ（最も詳細）
LOG_LEVEL=debug go test -v

# 情報ログ（デフォルト）
LOG_LEVEL=info go test -v

# 警告以上のみ
LOG_LEVEL=warn go test -v

# エラーのみ
LOG_LEVEL=error go test -v
```

## 📈 実行結果の見方

### 成功時の出力

```
🚀 Starting SaaSus Billing API E2E Tests
========================================

🎬 STORY START: Postman Collection Story - Standard Methods
═══════════════════════════════════════════════════════════════

📋 STEP 1: Pre_GetStripeConnectionInformation
   Method: GetStripeInfo
   ✅ Step Result: passed | Status Code: 200 | Duration: 467ms

📋 STEP 2: UpdateStripeConnectionInfo
   Method: UpdateStripeInfo
   ✅ Step Result: passed | Status Code: 200 | Duration: 465ms

[... 続く ...]

✅ STORY COMPLETED: Postman Collection Story - Standard Methods (Duration: 1.2s)

🎬 STORY START: Postman Collection Story - WithResponse Methods
═══════════════════════════════════════════════════════════════

[... 同様のステップ ...]

📊 Test Execution Summary:
Stories: 2/2 passed (100.0%)
Steps: 12/12 passed (100.0%)
Method Coverage: 8/8 methods (100.0%)
Execution Time: 1.8s

📊 Method Coverage Summary:
Total Coverage: 8/8 (100.0%)
🎉 All methods have been tested!

🎉 ALL TESTS PASSED - FULL METHOD COVERAGE ACHIEVED!
--- PASS: TestBillingAPIE2E (1.73s)
PASS
```

### 失敗時の出力

```
📋 STEP 2: UpdateStripeConnectionInfo
   Method: UpdateStripeInfo
   ❌ Step Result: failed | Status Code: 400 | Duration: 234ms
   Error: Bad Request - Invalid Stripe key format

❌ STORY FAILED: Postman Collection Story - Standard Methods
Error: step 'UpdateStripeConnectionInfo' failed

--- FAIL: TestBillingAPIE2E (0.87s)
    billing_test.go:62: ❌ Some tests failed or coverage is incomplete
FAIL
```

## 🔧 技術詳細

### testlib使用パターン

```go
import (
    "github.com/saasus-platform/saasus-sdk-go/modules/billing"
    "github.com/saasus-platform/saasus-sdk-go/tests/testlib"
)

func TestBillingAPIE2E(t *testing.T) {
    // 1. クライアント作成
    client, err := billing.BillingWithResponse()

    // 2. エンジン作成（クライアントとメソッド一覧を渡す）
    engine := testlib.NewE2EEngine(client, GetBillingMethods())

    // 3. 設定の検証
    if err := engine.Config.Validate(); err != nil {
        t.Fatalf("Configuration error: %v", err)
    }

    // 4. ストーリー実行
    results := engine.ExecuteStories(GetBillingStories())

    // 5. 結果の確認
    engine.PrintResults(results)
}
```

### ストーリー定義例

```go
// stories.go内
func GetBillingStories() []testlib.Story {
    return []testlib.Story{
        {
            Name:        "Postman Collection Story - Standard Methods",
            Description: "Test story based on Postman collection using standard methods",
            Steps: []testlib.Step{
                {
                    Name:           "Pre_GetStripeConnectionInformation",
                    ClientMethod:   "GetStripeInfo",
                    ExpectedStatus: 200,
                    Validation:     validateStripeInfoResponse,
                    StateUpdate:    updateStripeInfoState,
                },
                // ... 他のステップ
            },
        },
    }
}
```

### メソッド一覧定義

```go
// stories.go内
func GetBillingMethods() []string {
    return []string{
        // Standard methods
        "DeleteStripeInfo",
        "GetStripeInfo",
        "UpdateStripeInfoWithBody",
        "UpdateStripeInfo",

        // WithResponse methods
        "DeleteStripeInfoWithResponse",
        "GetStripeInfoWithResponse",
        "UpdateStripeInfoWithBodyWithResponse",
        "UpdateStripeInfoWithResponse",
    }
}
```

### バリデーション関数

```go
// validation.go内
func validateStripeInfoResponse(response interface{}) error {
    // レスポンス検証ロジック
    if response == nil {
        return fmt.Errorf("response is nil")
    }
    return nil
}
```

## 🔄 開発ワークフロー

### 1. ローカルでテスト開発

```bash
# 詳細ログで実行
LOG_LEVEL=debug go test -v

# 特定のエラーをデバッグ
LOG_LEVEL=debug go test -v 2>&1 | grep "ERROR\|FAILED"
```

### 2. テスト修正

```go
// stories.goを編集してストーリーやステップを追加・修正
// validation.goを編集してバリデーションロジックを追加・修正
```

### 3. コミット前確認

```bash
# 標準ログレベルでテスト実行
go test -v

# すべてのE2Eテストを実行
cd ../..  # tests/ディレクトリへ
go test ./e2e/...
```

### 4. CI/CDでの実行

CI/CDパイプラインでは自動的に実行されます（詳細は[TEST_EXECUTION_GUIDE.md](../TEST_EXECUTION_GUIDE.md)参照）

## 🐛 トラブルシューティング

### よくある問題

#### Stripe認証エラー

```
❌ Status Code: 400 - Invalid Stripe key
```

**解決策:**
```bash
# テスト用Stripeキーを設定
echo "STRIPE_SECRET_KEY=sk_test_your_key" >> ../../../.env
```

#### タイムアウト

```
panic: test timed out after 10m0s
```

**解決策:**
```bash
go test -v -timeout 30m
```

#### メソッドが見つからない

```
❌ Error: method 'GetStripeInfo' not found on client
```

**解決策:**
```bash
# プロジェクトルートで再生成
cd ../../../
./generate_all.sh
go mod tidy
```

## 📚 関連ドキュメント

- **[TEST_EXECUTION_GUIDE.md](../TEST_EXECUTION_GUIDE.md)** - テスト実行方法の完全ガイド
- **[testlib README](../../testlib/README.md)** - 共通テストライブラリの詳細
- **[API_MODULE_E2E_SETUP_GUIDE.md](../API_MODULE_E2E_SETUP_GUIDE.md)** - 他のAPIモジュールへの展開方法

## 🎯 次のステップ

このテストを参考に他のAPIモジュールのE2Eテストを作成：

1. **Auth API** - ユーザー認証・管理API
2. **Pricing API** - 料金プラン・使用量API
3. **Communication API** - フィードバック・通知API

作成手順は[API_MODULE_E2E_SETUP_GUIDE.md](../API_MODULE_E2E_SETUP_GUIDE.md)を参照。
