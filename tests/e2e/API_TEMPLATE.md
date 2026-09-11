# 新しいAPI E2Eテスト実装テンプレート

**新しいAPIのE2Eテスト実装ガイド** - 段階的な実装手順とテンプレート

## 🎯 このテンプレートについて

このドキュメントは、新しいAPIのE2Eテストを実装する際の完全なテンプレートとガイドです。Auth API、Pricing API、Communication API等の実装時に使用してください。

## 📋 実装前チェックリスト

### 1. 事前準備

- [ ] APIクライアントが `modules/{api}/client.go` に存在することを確認
- [ ] 生成されたクライアントが `generated/{api}api/` に存在することを確認
- [ ] OpenAPI仕様が `api/{api}api.yml` に存在することを確認
- [ ] 環境変数に必要な認証情報が設定されていることを確認

### 2. 実装計画

- [ ] 対象APIの全メソッドをリストアップ
- [ ] 基本的なCRUDフローを設計
- [ ] 認証・権限テストを設計
- [ ] エラーケーステストを設計

## 🚀 Step-by-Step実装ガイド

### Step 1: ディレクトリ構造の作成

```bash
# 新しいAPIディレクトリを作成（例：authapi）
mkdir tests/e2e/authapi
cd tests/e2e/authapi
```

### Step 2: テストエントリーポイントの作成

**ファイル名**: `{api}_test.go` (例: `auth_test.go`)

```go
package main

import (
    "testing"
    "github.com/saasus-platform/saasus-sdk-go/tests/testlib"
    "github.com/saasus-platform/saasus-sdk-go/modules/auth" // APIに応じて変更
)

func TestAuthAPI(t *testing.T) {
    // クライアント初期化
    client := auth.NewClient() // APIに応じて変更
    
    // E2Eエンジン作成
    engine := testlib.NewE2EEngine(client, GetAuthMethods()) // APIに応じて変更
    
    // ストーリー実行
    results := engine.ExecuteStories(GetAuthStories()) // APIに応じて変更
    
    // 結果出力
    engine.PrintResults(results)
    
    // テスト結果の検証
    for _, result := range results {
        if result.Status != testlib.TestStatusPassed {
            t.Errorf("Story '%s' failed: %v", result.StoryName, result.Error)
        }
    }
}
```

### Step 3: メソッド定義の作成

**ファイル名**: `methods.go`

```go
package main

// GetAuthMethods returns all Auth API methods for coverage tracking
// 注意: generated/{api}api/client.gen.go から実際のメソッド名を取得してください
func GetAuthMethods() []string {
    return []string{
        // 基本的なCRUDメソッド（例）
        "GetResources",
        "GetResourcesWithResponse",
        "CreateResource",
        "CreateResourceWithResponse",
        "GetResource",
        "GetResourceWithResponse",
        "UpdateResource",
        "UpdateResourceWithResponse",
        "DeleteResource",
        "DeleteResourceWithResponse",
        
        // API固有のメソッド（例：Auth APIの場合）
        "SignIn",
        "SignInWithResponse",
        "SignOut",
        "SignOutWithResponse",
        "RefreshToken",
        "RefreshTokenWithResponse",
        
        // 他のメソッドを追加...
        // 注意: 実際のAPIクライアントから正確なメソッド名を取得してください
    }
}
```

### Step 4: ストーリー定義の作成

**ファイル名**: `stories.go`

```go
package main

import "github.com/saasus-platform/saasus-sdk-go/tests/testlib"

// GetAuthStories returns test stories for Auth API
func GetAuthStories() []testlib.Story {
    return []testlib.Story{
        // ストーリー1: 基本的なCRUDフロー
        {
            Name: "Resource Management Flow",
            Variables: map[string]interface{}{
                "resourceId":   "",
                "resourceName": "Test Resource",
            },
            Steps: []testlib.Step{
                {
                    Name:           "List Resources",
                    ClientMethod:   "GetResources",
                    ExpectedStatus: 200,
                    Validation:     validateResourcesList,
                },
                {
                    Name:         "Create Resource",
                    ClientMethod: "CreateResource",
                    Parameters: map[string]interface{}{
                        "name":        "{{resourceName}}",
                        "description": "Test resource for E2E testing",
                    },
                    ExpectedStatus: 201,
                    Validation:     validateResourceCreation,
                    StateUpdate:    updateResourceIdFromResponse,
                },
                {
                    Name:           "Get Created Resource",
                    ClientMethod:   "GetResource",
                    Parameters:     "{{resourceId}}", // 変数参照
                    ExpectedStatus: 200,
                    Validation:     validateResourceDetails,
                },
                {
                    Name:         "Update Resource",
                    ClientMethod: "UpdateResource",
                    Parameters: map[string]interface{}{
                        "id":          "{{resourceId}}",
                        "name":        "Updated Test Resource",
                        "description": "Updated description",
                    },
                    ExpectedStatus: 200,
                    Validation:     validateResourceUpdate,
                },
                {
                    Name:           "Delete Resource",
                    ClientMethod:   "DeleteResource",
                    Parameters:     "{{resourceId}}",
                    ExpectedStatus: 204,
                },
            },
            Cleanup: func() error {
                // テスト後のクリーンアップ処理
                // 必要に応じて実装
                return nil
            },
        },
        
        // ストーリー2: API固有のフロー（例：認証フロー）
        {
            Name: "Authentication Flow",
            Variables: map[string]interface{}{
                "accessToken":  "",
                "refreshToken": "",
            },
            Steps: []testlib.Step{
                {
                    Name:         "Sign In",
                    ClientMethod: "SignIn",
                    Parameters: map[string]interface{}{
                        "email":    "test@example.com",
                        "password": "testpassword",
                    },
                    ExpectedStatus: 200,
                    Validation:     validateSignInResponse,
                    StateUpdate:    updateTokensFromResponse,
                },
                {
                    Name:           "Refresh Token",
                    ClientMethod:   "RefreshToken",
                    Parameters:     "{{refreshToken}}",
                    ExpectedStatus: 200,
                    Validation:     validateTokenRefresh,
                },
                {
                    Name:           "Sign Out",
                    ClientMethod:   "SignOut",
                    ExpectedStatus: 200,
                },
            },
        },
        
        // ストーリー3: エラーケーステスト
        {
            Name: "Error Handling Tests",
            Steps: []testlib.Step{
                {
                    Name:         "Invalid Parameter Test",
                    ClientMethod: "CreateResource",
                    Parameters: map[string]interface{}{
                        "name": "", // 無効なパラメータ
                    },
                    ExpectedStatus: 400,
                    Validation:     validateErrorResponse,
                },
                {
                    Name:           "Not Found Test",
                    ClientMethod:   "GetResource",
                    Parameters:     "non-existent-id",
                    ExpectedStatus: 404,
                    Validation:     validateNotFoundResponse,
                },
            },
        },
    }
}
```

### Step 5: バリデーション関数の作成

**ファイル名**: `validation.go`

```go
package main

import (
    "fmt"
    "reflect"
)

// validateResourcesList validates the resources list response
func validateResourcesList(response interface{}) error {
    if response == nil {
        return fmt.Errorf("response is nil")
    }
    
    respValue := reflect.ValueOf(response)
    if respValue.Kind() == reflect.Ptr {
        respValue = respValue.Elem()
    }
    
    // 必要なフィールドの存在確認
    resourcesField := respValue.FieldByName("Resources")
    if !resourcesField.IsValid() {
        return fmt.Errorf("Resources field not found in response")
    }
    
    return nil
}

// validateResourceCreation validates resource creation response
func validateResourceCreation(response interface{}) error {
    if response == nil {
        return fmt.Errorf("resource creation response is nil")
    }
    
    respValue := reflect.ValueOf(response)
    if respValue.Kind() == reflect.Ptr {
        respValue = respValue.Elem()
    }
    
    // リソースIDの存在確認
    idField := respValue.FieldByName("ResourceId")
    if !idField.IsValid() || idField.String() == "" {
        return fmt.Errorf("ResourceId not found or empty in response")
    }
    
    return nil
}

// updateResourceIdFromResponse updates the resource ID in test variables
func updateResourceIdFromResponse(response interface{}, variables map[string]interface{}) error {
    respValue := reflect.ValueOf(response)
    if respValue.Kind() == reflect.Ptr {
        respValue = respValue.Elem()
    }
    
    idField := respValue.FieldByName("ResourceId")
    if idField.IsValid() {
        variables["resourceId"] = idField.String()
    }
    
    return nil
}

// validateResourceDetails validates resource details response
func validateResourceDetails(response interface{}) error {
    if response == nil {
        return fmt.Errorf("resource details response is nil")
    }
    
    respValue := reflect.ValueOf(response)
    if respValue.Kind() == reflect.Ptr {
        respValue = respValue.Elem()
    }
    
    // 必要なリソース情報の確認
    fields := []string{"ResourceId", "Name", "Description"}
    for _, field := range fields {
        fieldValue := respValue.FieldByName(field)
        if !fieldValue.IsValid() {
            return fmt.Errorf("%s field not found in resource details", field)
        }
    }
    
    return nil
}

// validateResourceUpdate validates resource update response
func validateResourceUpdate(response interface{}) error {
    if response == nil {
        return fmt.Errorf("resource update response is nil")
    }
    
    respValue := reflect.ValueOf(response)
    if respValue.Kind() == reflect.Ptr {
        respValue = respValue.Elem()
    }
    
    // 更新されたリソース情報の確認
    nameField := respValue.FieldByName("Name")
    if nameField.IsValid() && nameField.String() != "Updated Test Resource" {
        return fmt.Errorf("resource name was not updated correctly")
    }
    
    return nil
}

// validateSignInResponse validates sign in response (API固有の例)
func validateSignInResponse(response interface{}) error {
    if response == nil {
        return fmt.Errorf("sign in response is nil")
    }
    
    respValue := reflect.ValueOf(response)
    if respValue.Kind() == reflect.Ptr {
        respValue = respValue.Elem()
    }
    
    // アクセストークンの存在確認
    tokenField := respValue.FieldByName("AccessToken")
    if !tokenField.IsValid() || tokenField.String() == "" {
        return fmt.Errorf("AccessToken not found or empty")
    }
    
    // リフレッシュトークンの存在確認
    refreshTokenField := respValue.FieldByName("RefreshToken")
    if !refreshTokenField.IsValid() || refreshTokenField.String() == "" {
        return fmt.Errorf("RefreshToken not found or empty")
    }
    
    return nil
}

// updateTokensFromResponse updates tokens in test variables (API固有の例)
func updateTokensFromResponse(response interface{}, variables map[string]interface{}) error {
    respValue := reflect.ValueOf(response)
    if respValue.Kind() == reflect.Ptr {
        respValue = respValue.Elem()
    }
    
    tokenField := respValue.FieldByName("AccessToken")
    if tokenField.IsValid() {
        variables["accessToken"] = tokenField.String()
    }
    
    refreshTokenField := respValue.FieldByName("RefreshToken")
    if refreshTokenField.IsValid() {
        variables["refreshToken"] = refreshTokenField.String()
    }
    
    return nil
}

// validateTokenRefresh validates token refresh response (API固有の例)
func validateTokenRefresh(response interface{}) error {
    if response == nil {
        return fmt.Errorf("token refresh response is nil")
    }
    
    respValue := reflect.ValueOf(response)
    if respValue.Kind() == reflect.Ptr {
        respValue = respValue.Elem()
    }
    
    // 新しいアクセストークンの確認
    tokenField := respValue.FieldByName("AccessToken")
    if !tokenField.IsValid() || tokenField.String() == "" {
        return fmt.Errorf("new AccessToken not found or empty")
    }
    
    return nil
}

// validateErrorResponse validates error response
func validateErrorResponse(response interface{}) error {
    if response == nil {
        return fmt.Errorf("error response is nil")
    }
    
    respValue := reflect.ValueOf(response)
    if respValue.Kind() == reflect.Ptr {
        respValue = respValue.Elem()
    }
    
    // エラーメッセージの存在確認
    errorField := respValue.FieldByName("Error")
    if !errorField.IsValid() {
        return fmt.Errorf("Error field not found in error response")
    }
    
    return nil
}

// validateNotFoundResponse validates 404 not found response
func validateNotFoundResponse(response interface{}) error {
    // 404レスポンスの場合、responseがnilの場合もある
    // ステータスコードで判定されるため、特別な検証は不要
    return nil
}
```

### Step 6: README作成

**ファイル名**: `README.md`

```markdown
# {API Name} API E2E Tests

**SaaSus {API Name} API専用E2Eテスト** - {API機能の説明}

> 📖 **上位ドキュメント**: [E2E Testing Framework](../README.md)

## 🚀 クイックスタート

```bash
# テスト実行
go test -v

# 詳細ログ付き実行
VERBOSE=true go test -v

# ドライラン実行
E2E_DRY_RUN=true go test -v
```

## 🎯 テストストーリー

### 1. Resource Management Flow
基本的なCRUD操作をテスト

### 2. {API固有のフロー名}
{API固有の機能}をテスト

### 3. Error Handling Tests
エラーケースの処理をテスト

## 📊 メソッドカバレッジ

**全XXメソッドの完全カバレッジ**

- 基本CRUD操作メソッド
- {API固有の機能}メソッド
- エラーハンドリングメソッド

## ⚙️ 設定

### 環境変数

```bash
# 基本認証情報
SAASUS_SAAS_ID=your-saas-id
SAASUS_API_KEY=your-api-key
SAASUS_SECRET_KEY=your-secret-key

# API固有の設定
SAASUS_{API_NAME}_BASE_URL=https://api.dev.saasus.io/v1/{api}

# 追加の認証情報（必要に応じて）
{ADDITIONAL_ENV_VARS}=your-values
```

### テスト実行オプション

```bash
# 基本実行
go test -v

# 詳細ログ
VERBOSE=true go test -v

# デバッグモード
VERBOSE=true LOG_LEVEL=debug go test -v

# ドライラン
E2E_DRY_RUN=true go test -v

# タイムアウト設定
go test -v -timeout 10m
```

## 🔧 カスタマイズ

### 新しいストーリーの追加

1. `stories.go`に新しいストーリーを追加
2. 必要なバリデーション関数を`validation.go`に実装
3. テストを実行して動作確認

### バリデーションの強化

1. `validation.go`の既存関数を拡張
2. より詳細なレスポンス検証を実装
3. ビジネスロジックの検証を追加

## 🐛 トラブルシューティング

### よくある問題

1. **認証エラー (401)**
   - 環境変数の設定を確認
   - API キーの有効性を確認

2. **メソッドが見つからない (Method not found)**
   - `methods.go`のメソッド名を確認
   - 生成されたクライアントとの整合性を確認

3. **バリデーションエラー**
   - レスポンス構造を確認
   - フィールド名の正確性を確認

### デバッグ手順

```bash
# 1. ドライランで基本構造を確認
E2E_DRY_RUN=true go test -v

# 2. 詳細ログで実行内容を確認
VERBOSE=true LOG_LEVEL=debug go test -v

# 3. 特定のテストのみ実行
go test -v -run TestSpecificStory
```

## 📚 関連ドキュメント

- [E2E Testing Framework](../README.md) - フレームワーク全体の概要
- [Getting Started Guide](../GETTING_STARTED.md) - 初心者向けガイド
- [Contributing Guide](../CONTRIBUTING.md) - 開発貢献ガイド
- [Architecture Document](../ARCHITECTURE.md) - 技術アーキテクチャ

---

**Happy Testing! 🎉**
```

## 🔧 実装後の確認事項

### 1. テスト実行確認

```bash
# 基本テスト実行
cd tests/e2e/{api}api
go test -v

# 期待される出力例
🚀 Starting SaaSus {API Name} API E2E Tests
========================================

📊 Test Execution Summary:
Stories: X/X passed (100.0%)
Steps: XX/XX passed (100.0%)
Method Coverage: XX/XX methods (100.0%)
Execution Time: X.Xs

🎉 ALL TESTS PASSED - FULL METHOD COVERAGE ACHIEVED!
```

### 2. ドライラン確認

```bash
E2E_DRY_RUN=true go test -v
```

### 3. 詳細ログ確認

```bash
VERBOSE=true LOG_LEVEL=debug go test -v
```

## 📋 実装完了チェックリスト

- [ ] ディレクトリ構造の作成完了
- [ ] テストエントリーポイント実装完了
- [ ] メソッド定義実装完了
- [ ] ストーリー定義実装完了
- [ ] バリデーション関数実装完了
- [ ] README作成完了
- [ ] 基本テスト実行成功
- [ ] ドライラン実行成功
- [ ] 詳細ログ出力確認
- [ ] 100%メソッドカバレッジ達成
- [ ] エラーケーステスト実装完了
- [ ] ドキュメント更新完了

## 🎯 次のステップ

1. **実装したAPIテストの実行**
2. **他のAPIテストの実装**
3. **スナップショットテストへの対応**
4. **CI/CD統合の確認**

---

このテンプレートを使用して、効率的に新しいAPIのE2Eテストを実装してください！