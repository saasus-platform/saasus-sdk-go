# Auth API テストデータ

このディレクトリには、Auth API の E2E テストで使用するテストパラメータとヘルパー関数が含まれています。

## ファイル構成

- `test_params.json` - 自動生成されたテストパラメータ
- `auth_test_helpers.go` - テストパラメータの読み込みとヘルパー関数
- `load_test_params_test.go` - LoadTestParams 関数のテスト

## test_params.json

`test-param-generator` ツールによって自動生成されたテストパラメータファイルです。
Auth API の各エンドポイントに必要なパラメータが含まれています。

### 構造

```json
{
  "users": {
    "createParams": {
      "email": "テスト用文字列",
      "password": "Passw0rd!"
    }
  },
  "roles": {
    "createParams": {
      "role_name": "テスト用文字列",
      "display_name": "テスト用文字列"
    }
  },
  ...
}
```

## auth_test_helpers.go

### TestParams 構造体

`test_params.json` の構造に対応した Go 構造体です。各 API カテゴリごとに `CreateParams` と `UpdateParams` を持ちます。

```go
type TestParams struct {
    Users struct {
        CreateParams map[string]any `json:"createParams"`
    } `json:"users"`
    
    Roles struct {
        CreateParams map[string]any `json:"createParams"`
    } `json:"roles"`
    
    // その他のカテゴリ...
}
```

### LoadTestParams 関数

テストパラメータファイルを読み込み、`TestParams` 構造体にデシリアライズします。

```go
func LoadTestParams(t *testing.T) *TestParams
```

**使用例:**

```go
func TestExample(t *testing.T) {
    // テストパラメータを読み込む
    params := testdata.LoadTestParams(t)
    
    // パラメータを使用
    email := params.Users.CreateParams["email"].(string)
    password := params.Users.CreateParams["password"].(string)
}
```

### ポインタヘルパー関数

Go の API では、オプショナルなフィールドにポインタ型を使用することがあります。
これらのヘルパー関数を使用すると、簡単にポインタを作成できます。

```go
func StringPtr(s string) *string
func IntPtr(i int) *int
func BoolPtr(b bool) *bool
```

**使用例:**

```go
// オプショナルなフィールドにポインタを設定
param := authapi.UpdateBasicInfoParam{
    DomainName:        "example.com",
    FromEmailAddress:  "noreply@example.com",
    ReplyEmailAddress: testdata.StringPtr("support@example.com"), // ポインタ
}
```

## ストーリーでの使用方法

ストーリー関数では、`LoadTestParams()` を使用してテストパラメータを読み込み、
`variables` マップに格納します。

```go
func GetPostmanStoryStandardMethods(t *testing.T) testlib.Story {
    // テストパラメータを読み込む
    params := testdata.LoadTestParams(t)
    
    // variables マップを構築
    variables := map[string]interface{}{
        "email":                 params.Users.CreateParams["email"],
        "password":              params.Users.CreateParams["password"],
        "domain_name":           params.BasicInfo.UpdateParams["domain_name"],
        "from_email_address":    params.BasicInfo.UpdateParams["from_email_address"],
        "reply_email_address":   params.BasicInfo.UpdateParams["reply_email_address"],
        "callback_url":          params.AuthInfo.UpdateParams["callback_url"],
        "role_name":             params.Roles.CreateParams["role_name"],
        "display_name":          params.Roles.CreateParams["display_name"],
        "attribute_name":        params.UserAttributes.CreateParams["attribute_name"],
        "attribute_type":        params.UserAttributes.CreateParams["attribute_type"],
        // その他の変数...
    }
    
    return testlib.Story{
        Name:        "Postman Collection Story - Standard Methods",
        Description: "Postman Collection のテストフローに基づいた Standard メソッドのテスト",
        Variables:   variables,
        Steps: []testlib.Step{
            // ステップ定義...
        },
    }
}
```

## パラメータ関数での使用方法

パラメータ関数は `variables` マップから値を取得します。
`variables` マップには、ストーリー関数で `LoadTestParams()` から読み込まれた値が含まれています。

```go
// createSaasUserParams は CreateSaasUser メソッドのパラメータを返します
func createSaasUserParams(variables map[string]interface{}) interface{} {
    // variables マップから値を取得（LoadTestParams から読み込まれた値）
    email, _ := variables["email"].(string)
    password, _ := variables["password"].(string)

    // ポインタヘルパー関数を使用
    return authapi.CreateSaasUserParam{
        Email:    email,
        Password: testdata.StringPtr(password),
    }
}
```

## テストパラメータのカスタマイズ

`test_params.json` の値を変更することで、テストパラメータをカスタマイズできます。

### 例: ユーザーのメールアドレスを変更

```json
{
  "users": {
    "createParams": {
      "email": "custom-user@example.com",
      "password": "CustomP@ssw0rd!"
    }
  }
}
```

### 環境別設定（オプション）

開発、ステージング、本番環境で異なるパラメータが必要な場合は、
環境別の JSON ファイルを作成できます：

- `test_params.dev.json` - 開発環境用
- `test_params.staging.json` - ステージング環境用
- `test_params.prod.json` - 本番環境用

環境変数 `TEST_ENV` を使用して、読み込むファイルを切り替えることができます：

```go
func LoadTestParams(t *testing.T) *TestParams {
    env := os.Getenv("TEST_ENV")
    if env == "" {
        env = "default"
    }
    
    filename := fmt.Sprintf("test_params.%s.json", env)
    paramFile := filepath.Join("testdata", filename)
    
    // ファイルが存在しない場合はデフォルトを使用
    if _, err := os.Stat(paramFile); os.IsNotExist(err) {
        paramFile = filepath.Join("testdata", "test_params.json")
    }
    
    // ファイルを読み込む...
}
```

## ベストプラクティス

1. **ポインタヘルパー関数を活用する**: オプショナルなフィールドには `StringPtr`, `IntPtr`, `BoolPtr` を使用
2. **型アサーションを安全に行う**: `value, ok := variables["key"].(type)` パターンを使用
3. **デフォルト値を提供する**: 値が存在しない場合のデフォルト値を設定
4. **テストパラメータを再利用する**: 同じパラメータを複数のテストで使用
5. **環境別設定を検討する**: 異なる環境で異なるパラメータが必要な場合

## トラブルシューティング

### テストパラメータが読み込めない

```
failed to read test params: open testdata/test_params.json: no such file or directory
```

**解決方法**: テストを実行するディレクトリが正しいことを確認してください。
テストは `tests/e2e/authapi/` ディレクトリから実行する必要があります。

### 型アサーションが失敗する

```
panic: interface conversion: interface {} is nil, not string
```

**解決方法**: 型アサーションの前に値が存在することを確認してください：

```go
email, ok := variables["email"].(string)
if !ok {
    email = "default@example.com" // デフォルト値
}
```

### ポインタが nil になる

```
panic: runtime error: invalid memory address or nil pointer dereference
```

**解決方法**: ポインタヘルパー関数を使用する前に、値が存在することを確認してください：

```go
var replyEmail *string
if reply, ok := variables["reply_email_address"].(string); ok && reply != "" {
    replyEmail = testdata.StringPtr(reply)
}
```

## 参考資料

- [Auth API 仕様](../../../../api/authapi.yml)
- [E2E テスト実装ガイド](../../SNAPSHOT_TEST_IMPLEMENTATION_GUIDE.md)
- [test-param-generator ツール](../../../tools/test-param-generator/)
- [Postman 分析レポート](../../../tools/postman-analyzer/reports/auth/)
