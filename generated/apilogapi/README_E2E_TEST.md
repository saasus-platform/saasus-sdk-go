# SaaSus ApiLog API End-to-End テスト

このディレクトリには、SaaSus ApiLog APIのEnd-to-Endテストが含まれています。

## ファイル構成

- `client_e2e_test.go` - 実際のSaaSus APIエンドポイントに対するE2Eテスト
- `test_params.json` - テストパラメータ設定ファイル
- `client_test.go` - 統合テスト（既存）
- `client.gen.go` - 自動生成されたAPIクライアント
- `types.gen.go` - 自動生成された型定義

## 環境変数の設定

E2Eテストを実行するには、以下の環境変数を設定する必要があります：

```bash
export SAASUS_SAAS_ID="your-saas-id"
export SAASUS_API_KEY="your-api-key"
export SAASUS_SECRET_KEY="your-secret-key"
```

これらの環境変数が設定されていない場合、テストは自動的にスキップされます。

## テストパラメータの設定

テストで使用するパラメータは `test_params.json` ファイルで設定できます：

```json
{
  "getLogs": {
    "testCases": [
      {
        "name": "パラメータなし",
        "params": {}
      },
      {
        "name": "制限数指定",
        "params": {
          "limit": 5
        }
      },
      {
        "name": "日付指定",
        "params": {
          "created_date": "2023-12-14"
        }
      }
    ],
    "pagination": {
      "limit": 2
    }
  },
  "getLog": {
    "invalidIds": [
      "invalid-uuid-format",
      "",
      "00000000-0000-0000-0000-000000000000"
    ]
  },
  "performance": {
    "maxResponseTime": "5s",
    "testLimit": 1
  }
}
```

### パラメータファイルの構造

- **getLogs.testCases**: GetLogsAPIのテストケース一覧
  - `name`: テストケース名
  - `params`: APIパラメータ（limit, created_date, created_at, cursor）
- **getLogs.pagination**: ページネーションテスト用の設定
  - `limit`: 1ページあたりの取得件数
- **getLog.invalidIds**: GetLogAPIで使用する無効なID一覧
- **performance**: パフォーマンステスト用の設定
  - `maxResponseTime`: 許容される最大レスポンス時間
  - `testLimit`: パフォーマンステストでの取得件数

## テストの実行

### 全てのE2Eテストを実行

```bash
go test ./generated/apilogapi -v -run "TestE2E"
```

### 特定のテストを実行

```bash
# GetLogs のテスト
go test ./generated/apilogapi -v -run "TestE2E_GetLogsWithResponse"

# GetLog のテスト
go test ./generated/apilogapi -v -run "TestE2E_GetLogWithResponse"

# ページネーションのテスト
go test ./generated/apilogapi -v -run "TestE2E_GetLogsWithPagination"

# エラーハンドリングのテスト
go test ./generated/apilogapi -v -run "TestE2E_ErrorHandling"

# パフォーマンステスト
go test ./generated/apilogapi -v -run "TestE2E_PerformanceBasic"
```

### 統合テストも含めて全て実行

```bash
go test ./generated/apilogapi -v
```

## テスト内容

### 1. GetLogsWithResponse テスト
- **基本機能テスト**: APIログ一覧取得の動作確認
- **パラメータバリエーション**: 22種類のパラメータ組み合わせテスト
  - 制限数: 最小値(1)、標準値(10)、大きな値(100)、最大値(1000)
  - 日付指定: 今日、昨日、1週間前、1ヶ月前
  - 日時指定: 午前、正午、午後、深夜、UTC午前0時
  - カーソル指定: 複数のサンプルカーソル
  - 組み合わせ: 制限数+日付、制限数+日時、日付+日時、全パラメータ
- **レスポンス検証**: データ構造、必須フィールド、制限数の遵守
- **日付フィルタ検証**: 指定日付でのフィルタリング動作確認

### 2. GetLogsWithResponse エッジケーステスト
- **境界値テスト**: 制限数0、負の値での400エラー確認
- **無効データテスト**: 不正な日付・日時形式での400エラー確認
- **カーソルテスト**: 空・無効なカーソルでのエラー確認
- **時間範囲テスト**: 未来の日付、過去の古い日付での動作確認

### 3. GetLogWithResponse テスト
- **正常系テスト**: 実際のAPIログIDを使用した取得テスト
- **異常系テスト**: 11種類の無効なID形式でのテスト
  - 短すぎるID、長すぎるID、無効な文字を含むID
  - ハイフンなしID、アンダースコア区切りID
  - NULL文字、特殊文字を含むID
- **エッジケーステスト**: 3種類の特殊なID形式
- **レスポンス検証**: 詳細なフィールド検証、IDの一致確認

### 4. ページネーション テスト
- **カーソルベースページネーション**: 設定可能なページサイズ
- **複数ページ取得**: 異なるページで異なるデータの確認
- **カーソル検証**: 次ページカーソルの存在確認

### 5. エラーハンドリング テスト
- **無効ID形式**: 様々な不正なUUID形式での400/404エラー確認
- **空ID**: 空文字列でのエラー確認
- **特殊文字**: NULL文字、制御文字を含むIDでのエラー確認

### 6. パフォーマンス テスト
- **基本パフォーマンス**: 設定可能な最大レスポンス時間での検証
- **負荷別テスト**: 3段階の負荷レベル
  - 軽負荷(1件): 1秒以内
  - 中負荷(50件): 3秒以内
  - 高負荷(100件): 5秒以内
- **スループット測定**: 件数/秒の計算
- **同時リクエストテスト**: 3つの並行リクエストでの性能確認
- **成功率検証**: 80%以上の成功率確認

## 認証について

このテストでは、SaaSus APIの認証方式（SigV1）を使用しています：

- `SAASUS_SECRET_KEY`: 署名生成用の秘密鍵
- `SAASUS_API_KEY`: APIキー
- `SAASUS_SAAS_ID`: SaaSID

認証は `client.SetSigV1()` 関数を使用して自動的に処理されます。

## テスト結果の例

```
=== RUN   TestE2E_GetLogsWithResponse
=== RUN   TestE2E_GetLogsWithResponse/パラメータなし
    client_e2e_test.go:102: ステータスコード: 200
    client_e2e_test.go:103: ステータス: 200 OK
    client_e2e_test.go:112: APIログ数: 10
    client_e2e_test.go:119: ログ 1:
    client_e2e_test.go:120:   - ID: 12345678-1234-1234-1234-123456789abc
    client_e2e_test.go:121:   - トレースID: trace-12345
    client_e2e_test.go:122:   - メソッド: GET
    client_e2e_test.go:123:   - URI: /v1/auth/userinfo
    client_e2e_test.go:124:   - ステータス: 200
--- PASS: TestE2E_GetLogsWithResponse/パラメータなし (0.43s)
```

## 注意事項

1. **実際のAPIエンドポイント**: このテストは実際のSaaSus APIエンドポイント（`https://api.saasus.io/v1/apilog`）に対して実行されます。

2. **認証情報**: 有効なSaaSus認証情報が必要です。無効な認証情報では401/403エラーが発生します。

3. **データの依存性**: 一部のテスト（特にGetLogWithResponse）は、実際に存在するAPIログデータに依存します。

4. **レート制限**: 大量のテストを実行する場合は、APIのレート制限に注意してください。

5. **ネットワーク**: インターネット接続が必要です。

## トラブルシューティング

### 環境変数が設定されていない場合
```
--- SKIP: TestE2E_GetLogsWithResponse (0.00s)
    client_e2e_test.go:25: E2Eテストをスキップ: SAASUS_SAAS_ID, SAASUS_API_KEY, SAASUS_SECRET_KEY 環境変数が設定されていません
```

### 認証エラーの場合
```
client_e2e_test.go:40: リクエストエラー: 401 Unauthorized
```

### ネットワークエラーの場合
```
client_e2e_test.go:40: リクエストエラー: dial tcp: lookup api.saasus.io: no such host
```

これらのエラーが発生した場合は、環境変数の設定、認証情報の確認、ネットワーク接続を確認してください。