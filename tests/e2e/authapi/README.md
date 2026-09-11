# Auth API E2E テスト

このディレクトリには、SaaSus SDK for Go の Auth API モジュールに対する E2E テストとスナップショットテストが含まれています。

## 概要

Auth API は認証、認可、マルチテナンシー管理を提供する SaaSus Platform の API サービスです。このテストスイートは、実際の API エンドポイントを使用して、Auth API の全機能を検証します。

## ファイル構成

```
tests/e2e/authapi/
├── auth_test.go                    # メインテストファイル
├── stories.go                      # ストーリー定義とメソッド一覧
├── validation.go                   # パラメータ関数とバリデーション
├── story_snapshot_config.go        # スナップショット設定
├── testdata/                       # テストデータ
│   ├── test_params.json            # 自動生成されたテストパラメータ
│   ├── auth_test_helpers.go        # ヘルパー関数
│   └── README.md                   # testdata の説明
└── README.md                       # このファイル
```

## テストストーリー

このテストスイートは、9つのストーリーで Auth API の全メソッドをカバーします：

### Postman Collection ベースストーリー（既存）

1. **Postman Collection Story - Standard Methods**
   - WithBody サフィックスを持たない標準メソッドのテスト
   - 基本的な HTTP レスポンスを返すメソッド
   - カバレッジ: 約48メソッド

2. **Postman Collection Story - Standard Methods With Body**
   - WithBody サフィックスを持つメソッドのテスト
   - `io.Reader` でリクエストボディを受け取るメソッド
   - カバレッジ: 約48メソッド

3. **Postman Collection Story - WithResponse Methods**
   - WithResponse サフィックスを持ち、WithBody を持たないメソッドのテスト
   - 詳細なレスポンス情報を返すメソッド
   - カバレッジ: 約48メソッド

4. **Postman Collection Story - WithResponse Methods With Body**
   - WithBodyWithResponse サフィックスを持つメソッドのテスト
   - `io.Reader` + 詳細レスポンスの組み合わせ
   - カバレッジ: 約49メソッド

### 追加ストーリー（未カバーメソッド対応）

5. **SaaS User Attributes Management Story**
   - SaaS ユーザーの追加属性を管理するメソッドのテスト
   - UpdateSaasUserAttributes 系の全バリアント（Standard, WithBody, WithResponse, WithBodyWithResponse）
   - カバレッジ: 4メソッド

6. **Tenant Invitations Story**
   - テナント招待機能に関連するメソッドのテスト
   - CreateTenantInvitation, GetTenantInvitations, ValidateInvitation 系の全バリアント
   - カバレッジ: 約24メソッド

7. **External User Link and Email Update Story**
   - 外部ユーザーリンクとメール更新確認機能のテスト
   - RequestExternalUserLink, ConfirmExternalUserLink, RequestEmailUpdate, ConfirmEmailUpdate 系の全バリアント
   - カバレッジ: 約16メソッド

8. **Provider Management Story**
   - 外部プロバイダー連携機能のテスト
   - UnlinkProvider, GetIdentityProviders, UpdateIdentityProvider, GetTenantIdentityProviders 系の全バリアント
   - カバレッジ: 約20メソッド

9. **Other Management Features Story**
   - その他の管理機能メソッドのテスト
   - 残りの未カバーメソッドの全バリアント
   - カバレッジ: 残りのメソッド

## 前提条件

### 環境変数の設定

プロジェクトルートに `.env` ファイルを作成し、以下の環境変数を設定してください：

```bash
SAASUS_SECRET_KEY=your_secret_key
SAASUS_API_KEY=your_api_key
SAASUS_SAAS_ID=your_saas_id

# Cognito credentials for obtaining ID/access tokens used in the stories
E2E_COGNITO_USER_POOL_ID=your_user_pool_id
E2E_COGNITO_CLIENT_ID=your_app_client_id
E2E_COGNITO_USERNAME=your_test_account
E2E_COGNITO_PASSWORD=your_test_password
E2E_COGNITO_REGION=ap-northeast-1          # optional (defaults to ap-northeast-1)
E2E_COGNITO_ENDPOINT=http://localhost:8081 # optional (set when testing against localstack)
```

`E2E_COGNITO_*` はストーリー開始前に実ユーザーの ID トークン／アクセストークンを取得するために使用されます。

### テストパラメータの準備

`testdata/test_params.json` にテストパラメータが含まれています。必要に応じて値を調整してください。

#### test_params.json の構造

テストパラメータは以下のカテゴリに分類されています：

```json
{
  "basic_info": {
    "domain_name": "example.dev.saasus.io",
    "from_email_address": "noreply@example.com",
    "reply_email_address": "support@example.com"
  },
  "auth_info": {
    "callback_url": "https://example.com/callback"
  },
  "user": {
    "email": "test@example.com",
    "password": "Passw0rd!"
  },
  "role": {
    "role_name": "test_role",
    "display_name": "Test Role"
  },
  "user_attribute": {
    "attribute_name": "test_attribute",
    "display_name": "Test Attribute",
    "attribute_type": "string"
  },
  "tenant_attribute": {
    "attribute_name": "test_tenant_attribute",
    "display_name": "Test Tenant Attribute",
    "attribute_type": "string"
  },
  "env": {
    "id": 10,
    "name": "test_env",
    "display_name": "Test Environment"
  },
  "tenant": {
    "name": "Test Tenant",
    "attributes": {},
    "back_office_staff_email": "admin@example.com"
  }
}
```

#### LoadTestParams() の使用方法

各ストーリー関数は `LoadTestParams(t)` を使用してテストパラメータを読み込みます：

```go
func GetPostmanStoryStandardMethods(t *testing.T) testlib.Story {
    // テストパラメータを読み込む
    params := testdata.LoadTestParams(t)
    
    // Variables マップを構築
    variables := map[string]interface{}{
        "domain_name":        params.BasicInfo.DomainName,
        "from_email_address": params.BasicInfo.FromEmailAddress,
        "email":              params.User.Email,
        "password":           params.User.Password,
        // ... その他の変数
    }
    
    return testlib.Story{
        Name:        "Postman Collection Story - Standard Methods",
        Description: "Standard メソッドのテスト",
        Variables:   variables,
        Steps:       []testlib.Step{
            // ステップ定義
        },
    }
}
```

#### テストパラメータのカスタマイズ

実際の API 仕様に合わせて `testdata/test_params.json` の値を調整できます：

1. `testdata/test_params.json` を編集
2. 必要なフィールドを追加または変更
3. `testdata/auth_test_helpers.go` の `TestParams` 構造体を更新（新しいフィールドを追加した場合）
4. テストを再実行

#### 環境別設定ファイル（オプション）

開発、ステージング、本番環境で異なるパラメータが必要な場合は、環境別の JSON ファイルを作成できます：

```bash
testdata/
├── test_params.json              # デフォルト（開発環境）
├── test_params.staging.json      # ステージング環境
└── test_params.production.json   # 本番環境
```

環境変数で使用するファイルを切り替え：

```bash
# ステージング環境でテスト
TEST_PARAMS_FILE=test_params.staging.json go test -v

# 本番環境でテスト
TEST_PARAMS_FILE=test_params.production.json go test -v
```

`LoadTestParams()` 関数を拡張して環境変数をサポート：

```go
func LoadTestParams(t *testing.T) *TestParams {
    filename := os.Getenv("TEST_PARAMS_FILE")
    if filename == "" {
        filename = "test_params.json"
    }
    
    data, err := os.ReadFile(filepath.Join("testdata", filename))
    if err != nil {
        t.Fatalf("failed to read test params: %v", err)
    }
    
    var params TestParams
    if err := json.Unmarshal(data, &params); err != nil {
        t.Fatalf("failed to unmarshal test params: %v", err)
    }
    
    return &params
}

## テストの実行

### 基本的な実行

```bash
# authapi ディレクトリに移動
cd tests/e2e/authapi

# テストを実行
go test -v
```

### 詳細ログ付き実行

```bash
# デバッグログを有効にして実行
LOG_LEVEL=debug go test -v
```

### 特定のストーリーのみ実行

```bash
# Standard Methods のみ実行
go test -v -stories="Postman Collection Story - Standard Methods"

# 複数のストーリーを実行（カンマ区切り）
go test -v -stories="Postman Collection Story - Standard Methods,Postman Collection Story - WithResponse Methods"

# 追加ストーリーを実行
go test -v -stories="SaaS User Attributes Management Story"
go test -v -stories="Tenant Invitations Story"
go test -v -stories="External User Link and Email Update Story"
go test -v -stories="Provider Management Story"
go test -v -stories="Other Management Features Story"
```

### スナップショットテスト

スナップショットテストは、API レスポンスをキャプチャして後方互換性を検証する機能です。

```bash
# スナップショットをキャプチャ（デフォルトモード）
go test -v -run TestAuthAPIWithSnapshot

# キャプチャモード（明示的）
go test -v -run TestAuthAPIWithSnapshot -snapshot-mode=capture

# 比較モード（既存のスナップショットと比較）
go test -v -run TestAuthAPIWithSnapshot -snapshot-mode=compare

# レポートモード（レポートを生成）
go test -v -run TestAuthAPIWithSnapshot -snapshot-mode=report

# 統合モード（キャプチャ + 比較 + レポート）
go test -v -run TestAuthAPIWithSnapshot -snapshot-mode=full

# 詳細ログ付き実行
go test -v -run TestAuthAPIWithSnapshot -snapshot-mode=capture -snapshot-verbose
```

#### スナップショット出力先

スナップショットファイルは以下のディレクトリに出力されます：

- **キャプチャファイル:** `tests/e2e/snapshot/auth/story_snapshots/`
- **比較結果:** `tests/e2e/snapshot/auth/story_comparisons/`
- **バリデーション結果:** `tests/e2e/snapshot/auth/story_validations/`
- **レポート:** `tests/e2e/snapshot/auth/story_reports/`

## テスト結果の確認

テスト実行後、以下の情報が表示されます：

- 各ストーリーの実行結果（成功/失敗）
- 各ステップの実行時間
- メソッドカバレッジ情報
- エラーメッセージ（失敗した場合）

### 成功例

```
🚀 Starting SaaSus Auth API E2E Tests
========================================
✅ Story: Postman Collection Story - Standard Methods (Passed)
✅ Story: Postman Collection Story - Standard Methods With Body (Passed)
✅ Story: Postman Collection Story - WithResponse Methods (Passed)
✅ Story: Postman Collection Story - WithResponse Methods With Body (Passed)
✅ Story: SaaS User Attributes Management Story (Passed)
✅ Story: Tenant Invitations Story (Passed)
✅ Story: External User Link and Email Update Story (Passed)
✅ Story: Provider Management Story (Passed)
✅ Story: Other Management Features Story (Passed)

✅ All Auth API E2E tests completed!
✅ Method Coverage: 250/250 methods (100%)
```

## メソッドカバレッジ

このテストスイートは、Auth API の全メソッド（250メソッド）をカバーします。`ReturnInternalServerError` および `ReturnInternalServerErrorWithResponse` メソッドのみ除外されます。

### カバレッジ詳細

- **ClientInterface**: 126メソッド（標準メソッド + WithBody メソッド）
- **ClientWithResponsesInterface**: 126メソッド（WithResponse メソッド + WithBodyWithResponse メソッド）
- **合計**: 252メソッド（ReturnInternalServerError系2つを除いて250メソッド）

メソッドカバレッジは `VerifyMethodCoverage()` 関数によって自動的に検証され、未使用のメソッドがある場合はエラーが報告されます。

## トラブルシューティング

### 環境変数が見つからない

```
Warning: .env file not found
```

→ プロジェクトルートに `.env` ファイルを作成してください。

### 認証エラー

```
❌ Failed to create auth client: authentication failed
```

→ `.env` ファイルの認証情報が正しいか確認してください。

### テストパラメータエラー

```
failed to read test params: no such file or directory
```

→ `testdata/test_params.json` が存在することを確認してください。

### SignUp エンドポイントが 400 エラーを返す

```
Error: user can not register yourself.
```

**原因**: セルフサインアップ（self_regist）が無効になっています。

**解決方法**: 
1. `UpdateSignInSettings` でセルフサインアップを有効化する必要があります
2. "Sign Up and Provider Management Story" では自動的に有効化されます
3. 手動で有効化する場合は、以下のパラメータで `UpdateSignInSettings` を呼び出してください：

```json
{
  "self_regist": {
    "enable": true
  }
}
```

### 外部ユーザーリンク/メール更新の Confirm 系メソッドが 401 エラーを返す

```
Status: 401 Unauthorized
```

**原因**: これらのエンドポイントは実際の確認コードが必要です。

**説明**:
- `RequestEmailUpdate` / `RequestExternalUserLink`: 確認コードをメールで送信するため、200 を返します
- `ConfirmEmailUpdate` / `ConfirmExternalUserLink`: 実際の確認コードが必要なため、E2E 環境では 401 エラーになります

**対応方法**:
1. **テスト環境**: Confirm 系メソッドの ExpectedStatus を 401 に設定（現在の実装）
2. **本番環境**: メールサーバーと統合して実際の確認コードを取得する仕組みが必要
3. **モック環境**: テスト用の確認コード生成 API を用意する

現在の実装では、Request 系メソッドは正常に動作し、Confirm 系メソッドは期待通り 401 を返すことを確認しています。

## 開発ガイド

### 新しいパラメータ関数の追加

`validation.go` にパラメータ関数を追加します：

```go
// 単一パラメータメソッドの例
func getXxxParams(variables map[string]interface{}) interface{} {
    xxxID, _ := variables["xxx_id"].(string)
    return xxxID  // 直接値を返す
}

// 複数パラメータメソッドの例
func getXxxYyyParams(variables map[string]interface{}) interface{} {
    xxxID, _ := variables["xxx_id"].(string)
    yyyID, _ := variables["yyy_id"].(string)
    return struct {
        XxxId string
        YyyId string
    }{
        XxxId: xxxID,
        YyyId: yyyID,
    }
}

// リクエストボディを持つメソッドの例
func createXxxParams(variables map[string]interface{}) interface{} {
    name, _ := variables["name"].(string)
    
    // 構造体を直接作成（JSON 経由の変換を避ける）
    return authapi.CreateXxxJSONRequestBody{
        Name: name,
    }
}

// WithBody メソッドの例
func createXxxWithBodyParams(variables map[string]interface{}) interface{} {
    name, _ := variables["name"].(string)
    
    param := authapi.CreateXxxParam{
        Name: name,
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

// WithBody メソッド（追加パラメータあり）の例
func updateXxxWithBodyParams(variables map[string]interface{}) interface{} {
    xxxID, _ := variables["xxx_id"].(string)
    name, _ := variables["name"].(string)
    
    param := authapi.UpdateXxxParam{
        Name: name,
    }
    body, _ := json.Marshal(param)
    
    return struct {
        XxxId       string    // 追加パラメータを最初に
        ContentType string
        Body        io.Reader
    }{
        XxxId:       xxxID,
        ContentType: "application/json",
        Body:        bytes.NewReader(body),
    }
}
```

#### パラメータ関数の実装ルール

1. **関数型を返さない**: 具体的な値（文字列、構造体など）を返す
2. **構造体フィールドの順序**: メソッドシグネチャのパラメータ順序と一致させる
3. **WithBody メソッドの特殊ルール**: 追加パラメータ → ContentType → Body の順序
4. **構造体を直接作成**: JSON 経由の変換を避ける（パフォーマンス向上）
5. **ポインタヘルパー関数の活用**: `testdata.StringPtr()`, `testdata.IntPtr()`, `testdata.BoolPtr()` を使用

### 新しいストーリーの追加

`stories.go` にストーリーを追加します：

```go
func GetNewStory(t *testing.T) testlib.Story {
    // テストパラメータを読み込む
    params := testdata.LoadTestParams(t)
    
    // Variables マップを構築
    variables := map[string]interface{}{
        "xxx_id": "",  // 動的に設定される値
        "name":   params.Xxx.Name,
        // ... その他の変数
    }
    
    return testlib.Story{
        Name:        "New Story",
        Description: "新しいストーリーの説明",
        Variables:   variables,
        Steps: []testlib.Step{
            {
                Name:           "Create Xxx",
                ClientMethod:   "CreateXxx",
                Parameters:     createXxxParams,
                ExpectedStatus: 201,
                StateUpdate: func(variables map[string]interface{}, resp interface{}) error {
                    // レスポンスから ID を抽出して変数に保存
                    if r, ok := resp.(*authapi.CreateXxxResponse); ok && r.JSON201 != nil {
                        variables["xxx_id"] = r.JSON201.Id
                    }
                    return nil
                },
            },
            {
                Name:           "Get Xxx",
                ClientMethod:   "GetXxx",
                Parameters:     getXxxParams,
                ExpectedStatus: 200,
            },
            // ... その他のステップ
        },
        Setup: func() error {
            // セットアップ処理（必要に応じて）
            return nil
        },
        Cleanup: func() error {
            // クリーンアップ処理（必要に応じて）
            return nil
        },
    }
}
```

#### ストーリー実装のポイント

1. **testing.T パラメータ**: `GetNewStory(t *testing.T)` のように testing.T を受け取る
2. **LoadTestParams の使用**: `testdata.LoadTestParams(t)` でパラメータを読み込む
3. **Variables マップの構築**: test_params.json から読み込んだ値を使用
4. **StateUpdate 関数**: レスポンスから値を抽出して変数に保存
5. **Setup/Cleanup 関数**: テストデータの準備と削除を実装
6. **独立性の確保**: 各ストーリーは独立して実行できるように設計

### バリデーションルールの追加

`story_snapshot_config.go` にバリデーションルールを追加します：

```go
{
    Type:     snapshot.ValidationRuleCustom,
    Enabled:  true,
    Severity: snapshot.ValidationSeverityError,
    Parameters: map[string]interface{}{
        "description": "カスタムバリデーション",
    },
}
```

## よくある質問（FAQ）

### Q1: test_params.json の値を変更したらテストが失敗する

A: test_params.json の値は実際の API 仕様に合わせて調整する必要があります。特に以下の点に注意してください：

- `domain_name`: 有効なドメイン形式（例: `example.dev.saasus.io`）
- `email`: 有効なメールアドレス形式
- `password`: パスワードポリシーに準拠（大文字、小文字、数字を含む8文字以上など）
- `env.id`: 既存の環境 ID と重複しない値（10以上を推奨）

### Q2: 特定のストーリーだけ実行したい

A: `-stories` フラグを使用してストーリー名を指定できます：

```bash
go test -v -stories="Postman Collection Story - Standard Methods"
```

複数のストーリーを実行する場合はカンマ区切りで指定します。

### Q3: メソッドカバレッジが 100% にならない

A: `VerifyMethodCoverage()` 関数が未使用のメソッドを検出します。以下を確認してください：

1. 全ストーリーが `GetAuthStories(t)` に含まれているか
2. 各ストーリーのステップで正しいメソッド名が指定されているか
3. `ReturnInternalServerError` 系のメソッドは除外されているか

### Q4: 環境別にテストパラメータを切り替えたい

A: 環境別の JSON ファイルを作成し、環境変数で切り替えます：

```bash
# ステージング環境
TEST_PARAMS_FILE=test_params.staging.json go test -v

# 本番環境
TEST_PARAMS_FILE=test_params.production.json go test -v
```

`LoadTestParams()` 関数を拡張して環境変数をサポートする必要があります（上記の「環境別設定ファイル」セクションを参照）。

### Q5: 新しいメソッドが追加されたらどうすればいい？

A: 以下の手順で対応します：

1. `stories.go` の `GetAuthMethods()` に新しいメソッドを追加
2. 適切なストーリーに新しいステップを追加
3. `validation.go` に必要なパラメータ関数を追加
4. `testdata/test_params.json` に必要なパラメータを追加
5. テストを実行して `VerifyMethodCoverage()` がエラーを返さないことを確認

## 参考資料

- [実装ガイド](../SNAPSHOT_TEST_IMPLEMENTATION_GUIDE.md)
- [要件定義書](../../../.kiro/specs/auth-e2e-snapshot-tests/requirements.md)
- [設計書](../../../.kiro/specs/auth-e2e-snapshot-tests/design.md)
- [タスクリスト](../../../.kiro/specs/auth-e2e-snapshot-tests/tasks.md)
- [未カバーメソッド分析](./UNCOVERED_METHODS_ANALYSIS.md)
- [testdata README](./testdata/README.md)

## ライセンス

このテストスイートは SaaSus SDK for Go の一部であり、同じライセンスが適用されます。
