# Apilog API E2E Tests

**SaaSus Apilog API専用E2Eテスト** - testlibを使用したシンプルなストーリーベーステスト

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
tests/e2e/apilogapi/
├── apilog_test.go              # メインE2Eテスト (TestApilogAPIE2E)
├── stories.go                  # ストーリー定義とメソッド一覧
├── validation.go               # バリデーション・ヘルパー関数
├── story_snapshot_config.go   # スナップショット設定
├── testdata/                   # テストデータ
│   ├── test_params.json        # テストパラメータ
│   └── apilog_test_helpers.go  # ヘルパー関数
└── README.md                   # このドキュメント
```

**シンプルな構成:**
- テスト関数は1つ（`TestApilogAPIE2E`）のみ
- すべてtestlib経由で実行
- 統合レイヤーなし

## 🎯 テストストーリー

### 1. Postman Collection Story - Standard Methods
Postmanコレクションのフローを標準メソッドで再現

**3ステップ:**
1. Pre_GetApiLogs - 初期ログ取得と変数抽出
2. GetApiLogs - 基本的なログ取得
3. GetApiLogs With QueryParameters - クエリパラメータ付きログ取得

### 2. Postman Collection Story - WithResponse Methods
同じフローをWithResponseメソッドで実行

**3ステップ:**
- 上記と同じフロー、すべてWithResponseメソッドを使用

## 📊 メソッドカバレッジ

**全4メソッドの100%カバレッジ達成**

| メソッド | 用途 | テストストーリー |
|---------|------|-----------------|
| `GetLogs` | APIログ一覧取得 | Standard |
| `GetLog` | 個別APIログ取得 | Standard |
| `GetLogsWithResponse` | レスポンス解析付きログ一覧取得 | WithResponse |
| `GetLogWithResponse` | レスポンス解析付き個別ログ取得 | WithResponse |

## ⚙️ 設定

### 環境変数（必須）

プロジェクトルートに`.env`ファイルを作成：

```bash
# SaaSus認証情報（必須）
SAASUS_SAAS_ID=your-saas-id
SAASUS_API_KEY=your-api-key
SAASUS_SECRET_KEY=your-secret-key

# オプション設定
SAASUS_BASE_URL=https://api.saasus.io/v1/apilog  # デフォルトは本番環境
LOG_LEVEL=info                                    # error/warn/info/debug
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
🚀 Starting SaaSus Apilog API E2E Tests
========================================

🎬 STORY START: Postman Collection Story - Standard Methods
═══════════════════════════════════════════════════════════════

📋 STEP 1: Pre_GetApiLogs
   Method: GetLogs
   ✅ Step Result: passed | Status Code: 200 | Duration: 467ms

📋 STEP 2: GetApiLogs
   Method: GetLogs
   ✅ Step Result: passed | Status Code: 200 | Duration: 465ms

📋 STEP 3: GetApiLogs With QueryParameters
   Method: GetLogs
   ✅ Step Result: passed | Status Code: 200 | Duration: 423ms

✅ STORY COMPLETED: Postman Collection Story - Standard Methods (Duration: 1.2s)

🎬 STORY START: Postman Collection Story - WithResponse Methods
═══════════════════════════════════════════════════════════════

[... 同様のステップ ...]

📊 Test Execution Summary:
Stories: 2/2 passed (100.0%)
Steps: 6/6 passed (100.0%)
Method Coverage: 4/4 methods (100.0%)
Execution Time: 1.8s

📊 Method Coverage Summary:
Total Coverage: 4/4 (100.0%)
🎉 All methods have been tested!

🎉 ALL TESTS PASSED - FULL METHOD COVERAGE ACHIEVED!
--- PASS: TestApilogAPIE2E (1.73s)
PASS
```

### 失敗時の出力

```
📋 STEP 1: Pre_GetApiLogs
   Method: GetLogs
   ❌ Step Result: failed | Status Code: 401 | Duration: 234ms
   Error: authentication failed: received 401 Unauthorized - check SAASUS_SECRET_KEY, SAASUS_API_KEY, and SAASUS_SAAS_ID

❌ STORY FAILED: Postman Collection Story - Standard Methods
Error: step 'Pre_GetApiLogs' failed

--- FAIL: TestApilogAPIE2E (0.87s)
    apilog_test.go:62: ❌ Some tests failed or coverage is incomplete
FAIL
```

## 🔧 技術詳細

### testlib使用パターン

```go
import (
    "github.com/saasus-platform/saasus-sdk-go/modules/apilog"
    "github.com/saasus-platform/saasus-sdk-go/tests/testlib"
)

func TestApilogAPIE2E(t *testing.T) {
    // 1. クライアント作成
    client, err := apilog.ApiLogWithResponse()

    // 2. エンジン作成（クライアントとメソッド一覧を渡す）
    engine := testlib.NewE2EEngine(client, GetApilogMethods())

    // 3. 設定の検証
    if err := engine.Config.Validate(); err != nil {
        t.Fatalf("Configuration error: %v", err)
    }

    // 4. ストーリー実行
    results := engine.ExecuteStories(GetApilogStories())

    // 5. 結果の確認
    engine.PrintResults(results)
}
```

### ストーリー定義例

```go
// stories.go内
func GetApilogStories() []testlib.Story {
    return []testlib.Story{
        {
            Name:        "Postman Collection Story - Standard Methods",
            Description: "Test story based on Postman collection using standard methods",
            Steps: []testlib.Step{
                {
                    Name:           "Pre_GetApiLogs",
                    ClientMethod:   "GetLogs",
                    Parameters:     getLogsParams,
                    ExpectedStatus: 200,
                    Validation:     validateGetLogsResponse,
                    StateUpdate:    extractLogsVariables,
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
func GetApilogMethods() []string {
    return []string{
        // Standard methods
        "GetLogs",
        "GetLog",

        // WithResponse methods
        "GetLogsWithResponse",
        "GetLogWithResponse",
    }
}
```

### バリデーション関数

```go
// validation.go内
func validateGetLogsResponse(response interface{}) error {
    // HTTPレスポンスのステータスコード検証
    if httpResp, ok := response.(*http.Response); ok {
        if httpResp.StatusCode == 401 {
            return fmt.Errorf("authentication failed: received 401 Unauthorized")
        }
        if httpResp.StatusCode < 200 || httpResp.StatusCode >= 300 {
            return fmt.Errorf("unexpected HTTP status: %d", httpResp.StatusCode)
        }
    }
    return nil
}
```

### 変数抽出関数

```go
// validation.go内
func extractLogsVariables(response interface{}, variables map[string]interface{}) error {
    // レスポンスから変数を抽出
    if apiLogs, ok := response.(*apilogapi.ApiLogs); ok {
        if apiLogs.Cursor != nil {
            variables["cursor"] = *apiLogs.Cursor
        }
        if len(apiLogs.ApiLogs) > 0 {
            firstLog := apiLogs.ApiLogs[0]
            variables["api_log_id"] = firstLog.ApiLogId
            variables["created_date"] = firstLog.CreatedDate
            variables["timestamp"] = firstLog.CreatedAt
            variables["created_at"] = convertUnixToISO8601(firstLog.CreatedAt)
        }
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

#### 認証エラー

```
❌ authentication failed: received 401 Unauthorized
```

**解決策:**
```bash
# 環境変数を確認
echo "SAASUS_SAAS_ID=$SAASUS_SAAS_ID"
echo "SAASUS_API_KEY=$SAASUS_API_KEY"
echo "SAASUS_SECRET_KEY=$SAASUS_SECRET_KEY"

# .envファイルを確認
cat ../../../.env
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
❌ Error: method 'GetLogs' not found on client
```

**解決策:**
```bash
# プロジェクトルートで再生成
cd ../../../
./generate_all.sh
go mod tidy
```

#### ログが存在しない

```
❌ Error: no logs found in response
```

**解決策:**
- APIを使用してログを生成してからテストを実行
- または、空のレスポンスを許容するようにバリデーションを修正

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

## 📝 API仕様

### GetLogs

APIログの一覧を取得します。

**クエリパラメータ:**
- `created_date` (optional): ログ取得日（YYYY-MM-DD形式）
- `created_at` (optional): ログ取得日時（ISO 8601形式）
- `limit` (optional): 取得件数の上限
- `cursor` (optional): カーソルページネーション用

**レスポンス:**
```json
{
  "api_logs": [
    {
      "api_log_id": "uuid",
      "created_at": 1234567890,
      "created_date": "2024-01-01",
      "api_key": "key",
      "request_method": "GET",
      "request_uri": "/api/endpoint",
      "response_status": "200",
      ...
    }
  ],
  "cursor": "next_page_cursor"
}
```

### GetLog

個別のAPIログを取得します。

**パスパラメータ:**
- `api_log_id` (required): APIログID

**レスポンス:**
```json
{
  "api_log_id": "uuid",
  "created_at": 1234567890,
  "created_date": "2024-01-01",
  ...
}
```
