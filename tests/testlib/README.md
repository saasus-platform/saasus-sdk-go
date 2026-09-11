# testlib - 共通テストライブラリ

**SaaSus SDK Go用共通テストライブラリ** - 全APIモジュールで再利用可能なE2Eテスト基盤

## 📦 パッケージ構成

```
testlib/
├── types.go      # 共通データ型定義
├── config.go     # 設定管理（.env自動読み込み）
├── logger.go     # 拡張ログ機能
├── coverage.go   # カバレッジ追跡
├── e2e.go        # E2Eエンジン
└── reporter.go   # レポート機能
```

## 🚀 基本使用方法

```bash
# 詳細ログでテスト実行
LOG_LEVEL=debug go test ./tests/e2e/billingapi/

# 標準ログでテスト実行
LOG_LEVEL=info go test ./tests/e2e/billingapi/
```

```go
import "github.com/saasus-platform/saasus-sdk-go/tests/testlib"

// E2Eエンジン作成
client := yourapi.NewClient()
methods := []string{"Method1", "Method2", "Method3"}
engine := testlib.NewE2EEngine(client, methods)

// ストーリー定義
stories := []testlib.Story{
    {
        Name: "Test Story",
        Steps: []testlib.Step{
            {
                Name:           "Test Step",
                ClientMethod:   "Method1",
                ExpectedStatus: 200,
                Validation:     func(response interface{}) error { return nil },
            },
        },
    },
}

// 実行
results := engine.ExecuteStories(stories)
engine.PrintResults(results)
```

## 🎯 主要コンポーネント

### E2EEngine

テストストーリーの実行エンジン

```go
type E2EEngine struct {
    Client   interface{}
    Config   *Config
    Logger   *Logger
    Coverage *CoverageTracker
    Reporter *Reporter
}
```

**主要メソッド:**
- `NewE2EEngine(client, methods)` - エンジン作成
- `ExecuteStories(stories)` - ストーリー実行
- `PrintResults(results)` - 結果出力

### Config

設定管理（環境変数・コマンドライン引数）

```go
type Config struct {
    SaaSID     string
    APIKey     string
    SecretKey  string
    LogLevel   LogLevel  // ログレベル (DEBUG, INFO, WARN, ERROR)
    DryRun     bool
    // ...
}
```

**機能:**
- `.env`ファイル自動読み込み
- 環境変数処理
- コマンドライン引数解析

### Logger

拡張ログ機能

```go
logger := NewLogger(logLevel)
logger.LogStoryStart("Story Name")
logger.LogStepResult("Step", "Method", status, 200, duration)
```

**機能:**
- ストーリー・ステップレベルログ
- 機密情報自動マスキング
- 設定可能なログレベル

### CoverageTracker

メソッドカバレッジ追跡

```go
tracker := NewCoverageTracker(methods)
tracker.RecordExecution("Method1", "Story", "Step", 200, duration, true, "")
covered, total, percentage := tracker.GetCoverage()
```

**機能:**
- メソッド実行状況追跡
- 成功率・実行時間統計
- 未テストメソッド検出

### Reporter

テスト結果レポート

```go
reporter := NewReporter(coverage)
reporter.PrintSummary(results)
reporter.PrintDetailed(results)
jsonData, _ := reporter.ExportJSON(results)
```

**機能:**
- サマリー・詳細・JSONレポート
- 実行メトリクス
- ファイル出力

## 📋 データ型

### Story

```go
type Story struct {
    Name        string
    Description string
    Steps       []Step
    Variables   map[string]interface{}
    Setup       func() error
    Cleanup     func() error
}
```

### Step

```go
type Step struct {
    Name           string
    ClientMethod   string
    Parameters     interface{}
    ExpectedStatus int
    Validation     func(response interface{}) error
    StateUpdate    func(response interface{}, variables map[string]interface{}) error
}
```

### Result

```go
type StoryResult struct {
    StoryName string
    Status    TestStatus
    Duration  time.Duration
    Steps     []StepResult
    Error     error
    Variables map[string]interface{}
}
```

## ⚙️ 設定オプション

### 環境変数

```bash
# 必須
SAASUS_SAAS_ID=your-saas-id
SAASUS_API_KEY=your-api-key
SAASUS_SECRET_KEY=your-secret-key

# オプション
LOG_LEVEL=debug              # ログレベル (debug, info, warn, error)
E2E_LOG_LEVEL=debug          # E2E固有のログレベル（LOG_LEVELをオーバーライド）
E2E_DRY_RUN=true
E2E_TIMEOUT=300
E2E_MAX_RETRIES=3
E2E_FAIL_FAST=true
SAASUS_BASE_URL=https://api.dev.saasus.io/v1  # ベースURL
```

### コマンドライン

```bash
--verbose          # 詳細ログ（LOG_LEVEL=debugと同等）
--dry-run          # ドライラン
--timeout 600      # タイムアウト（秒）
--fail-fast        # 失敗時即停止
--real-stripe      # 実際のStripeキー使用
```

### ログレベル設定

ログの詳細度は`LOG_LEVEL`環境変数で制御できます：

```bash
# 最も詳細なログ（デバッグ情報を含む）
LOG_LEVEL=debug go test ./...

# 標準ログ（デフォルト）
LOG_LEVEL=info go test ./...

# 警告以上のみ
LOG_LEVEL=warn go test ./...

# エラーのみ
LOG_LEVEL=error go test ./...
```

E2E固有の設定が必要な場合は`E2E_LOG_LEVEL`を使用：

```bash
# 一般的なログはinfoレベル、E2Eテストのみdebugレベル
LOG_LEVEL=info E2E_LOG_LEVEL=debug go test ./tests/e2e/...
```

## 🔧 拡張方法

### 新しいAPIでの使用

```go
// 1. クライアント準備
client := newapi.NewClient()

// 2. メソッドリスト定義
methods := []string{"GetData", "UpdateData", "DeleteData"}

// 3. エンジン作成
engine := testlib.NewE2EEngine(client, methods)

// 4. ストーリー定義
stories := []testlib.Story{
    {
        Name: "CRUD Operations",
        Steps: []testlib.Step{
            {Name: "Get", ClientMethod: "GetData", ExpectedStatus: 200},
            {Name: "Update", ClientMethod: "UpdateData", ExpectedStatus: 200},
            {Name: "Delete", ClientMethod: "DeleteData", ExpectedStatus: 200},
        },
    },
}

// 5. 実行
results := engine.ExecuteStories(stories)
engine.PrintResults(results)
```

### カスタムバリデーション

```go
func validateResponse(response interface{}) error {
    // カスタム検証ロジック
    if response == nil {
        return fmt.Errorf("response is nil")
    }
    return nil
}

step := testlib.Step{
    Name:       "Custom Validation",
    Validation: validateResponse,
}
```

### 状態管理

```go
func updateState(response interface{}, variables map[string]interface{}) error {
    // レスポンスから値を抽出して変数に保存
    variables["extracted_id"] = extractID(response)
    return nil
}

step := testlib.Step{
    Name:        "Extract ID",
    StateUpdate: updateState,
}
```

**重要な改善（2024年11月）:** スナップショットモードでは、StateUpdate関数に渡されるレスポンスは自動的にキャプチャされたJSONデータ（`map[string]interface{}`）に変換されます。これにより、レスポンスボディの二重読み取りエラーが回避され、StateUpdate関数の実装が簡素化されます。

また、`tests/testlib/snapshot/capture.go`の`CaptureStepResponse`メソッドが強化され、HTTPレスポンスボディの読み取りと再設定処理が追加されました。これにより、StateUpdateでレスポンスボディを再読み取り可能になっていますが、推奨される方法は自動変換された`map[string]interface{}`を使用することです。

```go
// ✅ 推奨: スナップショットモード対応のStateUpdate
func updateState(response interface{}, variables map[string]interface{}) error {
    // スナップショットモードでは、responseは自動的にmap[string]interface{}に変換される
    respMap, ok := response.(map[string]interface{})
    if !ok {
        // 通常モードの場合は、従来の型アサーションを試行
        return fmt.Errorf("unexpected response type: %T", response)
    }
    
    if id, ok := respMap["id"].(string); ok {
        variables["resource_id"] = id
    }
    return nil
}
```

## 🎯 利点

1. **再利用性**: 全APIで共通利用
2. **一貫性**: 統一されたテスト体験
3. **拡張性**: 新機能の簡単追加
4. **保守性**: 中央集約された機能管理
5. **柔軟性**: カスタマイズ可能な各種機能