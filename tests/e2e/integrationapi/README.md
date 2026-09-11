# Integration API E2Eテスト

Integration APIモジュールのE2Eテストとスナップショットテストの実装。

## 概要

Integration APIは、Amazon EventBridgeとの連携機能を提供し、リアルタイムイベント配信を実現します。このE2Eテストは、testlibフレームワークを活用し、Postman Collectionのテストフローを再現する4つのストーリーを実装しています。

## ファイル構成

```
tests/e2e/integrationapi/
├── integration_test.go          # メインテストファイル
├── stories.go                   # ストーリー定義とメソッド一覧
├── validation.go                # パラメータ関数とバリデーション
├── story_snapshot_config.go     # スナップショット設定
└── README.md                    # このファイル
```

## テストストーリー

### 1. Postman Collection Story - Standard Methods
標準メソッドを使用したEventBridge設定のライフサイクルテスト。

**ステップ:**
1. 初期状態確認（設定なし）
2. 設定保存
3. 保存後確認（設定あり）
4. 設定削除
5. 削除後確認（設定なし）
6. テスト準備
7. テスト用設定保存
8. テストイベント送信
9. 実際のイベント送信
10. クリーンアップ

### 2. Postman Collection Story - Standard Methods With Body
WithBodyメソッドを使用したテスト（SaveEventBridgeSettingsWithBody、CreateEventBridgeEventWithBody）。

### 3. Postman Collection Story - WithResponse Methods
WithResponseメソッドを使用したテスト（詳細なレスポンス情報を取得）。

### 4. Postman Collection Story - WithResponse Methods With Body
WithBodyWithResponseメソッドを使用したテスト。

## テスト実行方法

### 基本的なE2Eテスト実行

```bash
cd tests/e2e/integrationapi
go test -v
```

### 詳細ログ付きテスト実行

```bash
LOG_LEVEL=debug go test -v
```

### 特定のストーリーのみ実行

```bash
go test -v -stories="Postman Collection Story - Standard Methods"
go test -v -stories="Postman Collection Story - Standard Methods With Body"
go test -v -stories="Postman Collection Story - WithResponse Methods"
go test -v -stories="Postman Collection Story - WithResponse Methods With Body"
```

## 環境変数

テスト実行には以下の環境変数が必要です：

```bash
SAASUS_SAAS_ID=your-saas-id
SAASUS_API_KEY=your-api-key
SAASUS_SECRET_KEY=your-secret-key
```

オプションの環境変数：

```bash
TEST_AWS_ACCOUNT_ID=267185063265  # テスト用AWS Account ID
TEST_AWS_REGION=ap-northeast-1    # テスト用AWS Region
LOG_LEVEL=debug                   # ログレベル（debug, info, warn, error）
```

## メソッドカバレッジ

Integration APIの全14メソッドをカバー：

**Standard Methods (5メソッド):**
- GetEventBridgeSettings
- SaveEventBridgeSettings
- DeleteEventBridgeSettings
- CreateEventBridgeEvent
- CreateEventBridgeTestEvent

**WithBody Methods (2メソッド):**
- SaveEventBridgeSettingsWithBody
- CreateEventBridgeEventWithBody

**WithResponse Methods (5メソッド):**
- GetEventBridgeSettingsWithResponse
- SaveEventBridgeSettingsWithResponse
- DeleteEventBridgeSettingsWithResponse
- CreateEventBridgeEventWithResponse
- CreateEventBridgeTestEventWithResponse

**WithBodyWithResponse Methods (2メソッド):**
- SaveEventBridgeSettingsWithBodyWithResponse
- CreateEventBridgeEventWithBodyWithResponse

## 実装の特徴

### パラメータ生成の最適化

構造体を直接作成することで、JSON経由の変換を避け、型安全性とパフォーマンスを向上：

```go
// ✅ 推奨パターン
func createSaveEventBridgeSettingsParams() integrationapi.SaveEventBridgeSettingsJSONRequestBody {
    return integrationapi.SaveEventBridgeSettingsParam{
        AwsAccountId: getTestAWSAccountID(),
        AwsRegion:    integrationapi.ApNortheast1,
    }
}
```

### WithBodyメソッドのパラメータ構造

```go
func createSaveEventBridgeSettingsBodyParams() struct {
    ContentType string
    Body        io.Reader
} {
    param := integrationapi.SaveEventBridgeSettingsParam{
        AwsAccountId: getTestAWSAccountID(),
        AwsRegion:    integrationapi.ApNortheast1,
    }
    body, _ := json.Marshal(param)
    
    return struct {
        ContentType string
        Body        io.Reader
    }{
        ContentType: "application/json",
        Body:        bytes.NewReader(body),
    }
}
```

### クリーンアップ処理

各ストーリーの前後でEventBridge設定をクリーンアップし、テストの独立性を保証：

```go
func cleanupEventBridgeSettings() error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    client, err := integration.IntegrationWithResponse()
    if err != nil {
        return nil // クリーンアップなのでエラーは無視
    }
    
    _, _ = client.DeleteEventBridgeSettings(ctx)
    return nil
}
```

## トラブルシューティング

### 認証エラー（401 Unauthorized）

**症状:**
```
authentication failed: received 401 Unauthorized
```

**解決策:**
1. 環境変数が正しく設定されているか確認
2. `.env`ファイルが存在し、正しいパスにあるか確認
3. API KeyとSecret Keyが有効か確認

### 501 Not Implemented エラー

**症状:**
```
unexpected status code: want 201, got 501
```

**原因:**
CreateEventBridgeEventエンドポイントがサーバー側で実装されていない可能性があります。

**解決策:**
1. APIサーバーのバージョンを確認
2. エンドポイントが利用可能か確認
3. 一時的にこのステップをスキップするか、期待ステータスコードを調整

### テストパラメータの読み込みエラー

**症状:**
```
failed to load test parameters
```

**解決策:**
1. `tests/tools/test-param-generator`を実行してパラメータを生成
2. `test_params.json`が存在するか確認
3. 環境変数`TEST_AWS_ACCOUNT_ID`と`TEST_AWS_REGION`を設定

## スナップショットテスト

スナップショット機能は`integration_test.go`に統合されており、`-snapshot-mode`フラグで実行できます。

### スナップショットテストの実行

```bash
# 通常のE2Eテスト（スナップショットモード未指定）
cd tests/e2e/integrationapi
go test -v

# 1. スナップショットキャプチャ（ベースライン作成）
go test -v -snapshot-mode=capture

# 2. スナップショット比較（変更検出）
go test -v -snapshot-mode=comparison

# 3. スナップショットレポート生成（可視化）
go test -v -snapshot-mode=reporting

# 4. 統合実行（全フェーズ）
go test -v -snapshot-mode=integrated

# 特定のストーリーのみ実行
go test -v -snapshot-mode=capture -stories="Postman Collection Story - Standard Methods"
```

### スナップショット設定

`story_snapshot_config.go`で4つのモード別設定を提供：

```go
// Captureモード - 現在のAPIレスポンスをキャプチャ
config := GetIntegrationSnapshotConfigForCapture()

// Comparisonモード - 以前のスナップショットと比較
config := GetIntegrationSnapshotConfigForComparison()

// Reportingモード - HTMLレポートを生成
config := GetIntegrationSnapshotConfigForReporting()

// Integratedモード - 全フェーズを一度に実行
config := GetIntegrationSnapshotConfigForIntegrated()
```

### スナップショットファイルの出力先

スナップショットファイルは`tests/e2e/snapshot/integration/`に出力されます：

```
tests/e2e/snapshot/integration/
├── story_snapshots/          # キャプチャされたAPIレスポンス
│   ├── story_snapshot_v1.0.0_postman_collection_story_standard_methods.json
│   ├── story_snapshot_v1.0.0_postman_collection_story_standard_methods_with_body.json
│   ├── story_snapshot_v1.0.0_postman_collection_story_withresponse_methods.json
│   └── story_snapshot_v1.0.0_postman_collection_story_withresponse_methods_with_body.json
├── story_comparisons/        # 比較結果
├── story_validations/        # 検証結果
└── story_reports/            # HTMLレポート
```

### バリデーションルール

スナップショット設定には4つのバリデーションルールが含まれます：

- **Completion（完了検証）** - Error: ストーリーが完全に実行されたことを確認
- **Sequence（シーケンス検証）** - Error: ステップが正しい順序で実行されたことを確認
- **StateTransition（状態遷移検証）** - Warning: 状態遷移が期待通りであることを確認
- **Timing（タイミング検証）** - Info（デフォルト無効）: 実行時間を記録

## 参考資料

- [testlibフレームワーク](../../testlib/README.md)
- [E2Eテストガイド](../README.md)
- [Integration API仕様](../../../api/integration.yml)
