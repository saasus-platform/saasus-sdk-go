# API モジュール E2E・スナップショットテスト作成ガイド

他の API モジュール（auth、pricing、communication 等）の e2e とスナップショットテストを作成するための手順書

## 📋 目次

1. [概要](#概要)
2. [前提条件](#前提条件)
3. [作成手順](#作成手順)
4. [ディレクトリ構造](#ディレクトリ構造)
5. [各ステップの詳細](#各ステップの詳細)
6. [チェックリスト](#チェックリスト)

## 概要

`tests/tools/INTEGRATION.md` で定義された OpenAPI Resolver 統合フローに沿って、Postman Collection から E2E / スナップショットテストを構築します。

**パイプライン概要:**

1. OpenAPI 仕様と生成済みクライアントを確認 (`openapi-resolver`)
2. Postman Collection を解析してテストフローとメソッドを抽出 (`postman-analyzer`)
3. 解析結果を元にテストパラメータとヘルパーコードを自動生成 (`test-param-generator`)
4. 生成物を活用して stories / メインテスト / スナップショット設定を仕上げる

**参考実装:** `tests/e2e/billingapi/`, `tests/tools/test-param-generator/tests/e2e/billing/`

## 前提条件

### 必要なツール

- Go 1.20 以上
- OpenAPI 仕様ファイル `api/<api-name>api.yml`
- 生成済みクライアント `generated/<api-name>api/` （`./generate_all.sh` 実行済みであること）
- Postman Collection JSON ファイル
- ツールチェーン
  - `tests/tools/openapi-resolver/`
  - `tests/tools/postman-analyzer/`
  - `tests/tools/test-param-generator/`

### 必要な知識

- Go 言語の基本
- テスト対象 API の仕様理解
- testlib と snapshot framework の使い方（`tests/testlib/README.md`, `tests/e2e/snapshot/docs/` 参照）
- `tests/tools/INTEGRATION.md` に記載の統合アーキテクチャ

## 作成手順

> このセクションの各ステップは `tests/tools/INTEGRATION.md` のツールチェーン構成を前提にしています。必要に応じて同ドキュメントの「データフロー」図と突き合わせながら進めてください。

### Step 1: OpenAPI / クライアント / Postman コレクション準備

```bash
# 1-0. OpenAPI仕様と生成クライアントを最新化
./generate_all.sh

# 1-0-1. operationID やメソッドバリアントを確認（任意）
tests/tools/openapi-resolver/openapi-resolver \
  -spec api/<api-name>api.yml \
  -list

# 1-1. Postman CollectionをJSON形式でエクスポート
# Postman UI: Collections > Export > Collection v2.1

# 1-2. エクスポートしたファイルを配置
cp /path/to/downloaded.json "tests/tools/postman-analyzer/collections/SaaSus <Service Name> API Schema.postman_collection.json"
# 例: SaaSus Analysis API の場合は "SaaSus Analysis API Schema.postman_collection.json"
```

> ℹ️ `openapi-resolver` は `tests/tools/INTEGRATION.md` の「ステップバイステップデバッグ」セクションを参照すると、パス正規化や operationID の確認に活用できます。

### Step 2: Postman 分析 & テストパラメータ生成

```bash
# 2-1. アナライザーを実行（例: SaaSus Analysis API）
cd tests/tools/postman-analyzer
go run analyzer.go \
  -collection "collections/SaaSus Analysis API Schema.postman_collection.json" \
  -output reports/analysis
# 他のAPIを分析する場合はファイル名と出力ディレクトリを適切なものに置き換えてください

# 2-2. 生成されたレポートを確認
cat reports/analysis/postman_analysis_report.txt
cat reports/analysis/postman_analysis.json | jq '.client_methods'

# 2-3. テストパラメータとヘルパーを自動生成
cd ../test-param-generator
./run.sh \
  -report ../postman-analyzer/reports/analysis/postman_analysis.json \
  -module <api-name> \
  -output ../../e2e/<api-name>api/testdata \
  -verbose
cd ../../..
```

**確認ポイント:**

- レポート `postman_analysis.json` の `test_flow` / `client_methods`
- 生成ファイル `tests/e2e/<api-name>api/testdata/test_params.json`
- 生成ファイル `tests/e2e/<api-name>api/testdata/<api-name>_test_helpers.go`

> 📁 `test-param-generator` の実行詳細や診断レポートの読み方は `tests/tools/INTEGRATION.md` および `tests/tools/test-param-generator/README.md` を参照してください。

### Step 3: テストディレクトリ作成

```bash
# 3-1. APIモジュール用ディレクトリ作成（未作成の場合）
mkdir -p tests/e2e/<api-name>api

# 3-2. test-param-generator の出力先ディレクトリを確認
ls tests/e2e/<api-name>api/testdata
# -> test_params.json / <api-name>_test_helpers.go が生成されていること

# 3-3. 必要なファイルを作成
cd tests/e2e/<api-name>api
touch <api-name>_test.go          # メインテストファイル
touch stories.go                  # ストーリー定義とメソッド一覧
touch validation.go               # バリデーション・ヘルパー
touch story_snapshot_config.go    # スナップショット設定

# 3-4. testdata を go module から参照できるようにする（必要に応じて go:embed 等を利用）
# 例: testdata ディレクトリを残したまま go test が実行されることを確認
```

### Step 4: ストーリーファイルの基本構造

**ファイル:** `tests/e2e/<api-name>api/stories.go`

このファイルには以下の 2 つを含めます：

1. ストーリー定義関数
2. メソッド一覧関数

> 💡 `test-param-generator` で生成された `testdata/test_params.json` / `<api-name>_test_helpers.go` をそのまま利用する場合は、
>
> - パッケージ名を `<api-name>api` または `testdata` に合わせて調整
> - 必要なヘルパー関数をエクスポート（例: `LoadTestParams`）
>   といった最小修正を加えて stories から呼び出せるようにしてください。

```go
package <api-name>api

import (
	"fmt"
	"github.com/saasus-platform/saasus-sdk-go/tests/testlib"
)

// Get<ApiName>Methods returns all available methods for <ApiName> API
func Get<ApiName>Methods() []string {
	return []string{
		// Standard methods (without WithBody suffix)
		"GetResource",
		"UpdateResource",
		"DeleteResource",

		// Standard methods that require request bodies (WithBody suffix)
		"CreateResourceWithBody",
		"UpdateResourceWithBody",

		// WithResponse methods (without WithBody suffix)
		"GetResourceWithResponse",
		"UpdateResourceWithResponse",
		"DeleteResourceWithResponse",

		// WithResponse methods that require request bodies (WithBody suffix)
		"CreateResourceWithBodyWithResponse",
		"UpdateResourceWithBodyWithResponse",
	}
}

// ... ストーリー定義は次のステップで追加 ...
```

**データソース:**

- Postman 分析レポートの「Method Coverage Analysis」セクション
- または`generated/<api-name>api/client.gen.go`のメソッド定義
- 自動生成された `testdata/test_params.json` と `<api-name>_test_helpers.go`

### Step 5: パラメータ関数の実装

**ファイル:** `tests/e2e/<api-name>api/validation.go`

ストーリーのステップで使用するパラメータ関数を実装します。パラメータ関数は`variables`を受け取り、適切な型のパラメータを返す必要があります。

**重要な実装パターン:**

1. **単一パラメータメソッド（例: GetFeedback）**

   ```go
   // GetFeedback(ctx, feedbackId, reqEditors) のようなメソッド
   func getFeedbackParams(variables map[string]interface{}) interface{} {
       feedbackID, _ := variables["feedback_id"].(string)
       return feedbackID  // 直接文字列を返す
   }
   ```

2. **複数パラメータメソッド（例: GetFeedbackComment）**

   ```go
   // GetFeedbackComment(ctx, feedbackId, commentId, reqEditors) のようなメソッド
   func getFeedbackCommentParams(variables map[string]interface{}) interface{} {
       feedbackID, _ := variables["feedback_id"].(string)
       commentID, _ := variables["comment_id"].(string)
       // 構造体を返す - StandardArgumentBuilderが自動的にフィールドを展開
       return struct {
           FeedbackId string
           CommentId  string
       }{
           FeedbackId: feedbackID,
           CommentId:  commentID,
       }
   }
   ```

3. **リクエストボディを持つメソッド（例: CreateFeedback）**

   ```go
   // CreateFeedback(ctx, body, reqEditors) のようなメソッド
   func createTestFeedbackParams(variables map[string]interface{}) interface{} {
       userID := getTestUserID()
       if uid, ok := variables["user_id"].(string); ok && uid != "" {
           userID = uid
       }

       // ✅ 推奨: 構造体を直接作成（型安全、パフォーマンス向上）
       return communicationapi.CreateFeedbackJSONRequestBody{
           UserId:              userID,
           FeedbackTitle:       "Test Feedback Title",
           FeedbackDescription: "Test Feedback Description",
       }
   }

   // ❌ 非推奨: JSON経由の変換（不要なオーバーヘッド）
   // bodyData := map[string]any{...}
   // bodyBytes, _ := json.Marshal(bodyData)
   // var body communicationapi.CreateFeedbackJSONRequestBody
   // _ = json.Unmarshal(bodyBytes, &body)
   // return body
   ```

4. **WithBody メソッド（例: UpdateFeedbackWithBody）**

   ```go
   // UpdateFeedbackWithBody(ctx, feedbackId, contentType, body, reqEditors) のようなメソッド
   func updateFeedbackWithBodyParams(variables map[string]interface{}) interface{} {
       feedbackID, _ := variables["feedback_id"].(string)
       param := communicationapi.UpdateFeedbackParam{
           FeedbackTitle:       "Updated Title",
           FeedbackDescription: "Updated Description",
       }
       body, _ := json.Marshal(param)

       // WithBodyArgumentBuilderが処理する構造体
       // FeedbackIdフィールドが最初に展開され、その後ContentTypeとBodyが追加される
       return struct {
           FeedbackId  string
           ContentType string
           Body        io.Reader
       }{
           FeedbackId:  feedbackID,
           ContentType: "application/json",
           Body:        bytes.NewReader(body),
       }
   }
   ```

**⚠️ 重要な注意点:**

- **関数型を返さない**: パラメータ関数は関数自体を返してはいけません。必ず具体的な値（文字列、構造体など）を返してください。
- **構造体フィールドの順序**: 複数パラメータメソッドの場合、構造体のフィールド順序はメソッドシグネチャのパラメータ順序と一致させてください。
- **WithBody メソッドの特殊ルール**: WithBody メソッドの場合、追加パラメータ（feedbackId など）を最初に配置し、その後に`ContentType`と`Body`フィールドを配置してください。

#### test-param-generator 生成物の活用

`test-param-generator`は以下のファイルを自動生成します：

1. **test_params.json** - テストパラメータのJSON
2. **<api-name>_test_helpers.go** - ヘルパー関数とLoadTestParams関数

### 推奨アプローチ: test_params.json の活用

**メリット:**
- ✅ **データとコードの分離**: テストデータをJSONで管理し、コードをシンプルに保つ
- ✅ **複数のテストで共有**: 同じパラメータを複数のストーリーやテストで再利用
- ✅ **保守性の向上**: パラメータの変更がJSONファイルの編集だけで完結
- ✅ **複雑なパラメータ構造**: ネストした構造や大量のフィールドを持つパラメータに最適
- ✅ **環境別の設定**: 開発、ステージング、本番環境で異なるパラメータを使い分け可能

**実装例（推奨パターン）:**

```go
// tests/e2e/<api-name>api/testdata/<api-name>_test_helpers.go
package testdata

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// TestParams はテストパラメータファイルの構造体
type TestParams struct {
	BasicInfo struct {
		UpdateParams map[string]any `json:"updateParams"`
	} `json:"basic-info"`
	
	Users struct {
		CreateParams map[string]any `json:"createParams"`
		UpdateParams map[string]any `json:"updateParams"`
	} `json:"users"`
	
	Tenants struct {
		CreateParams map[string]any `json:"createParams"`
		UpdateParams map[string]any `json:"updateParams"`
	} `json:"tenants"`
	
	// ... 他のリソース
}

// LoadTestParams はテストパラメータファイルを読み込みます
func LoadTestParams(t *testing.T) *TestParams {
	paramFile := filepath.Join("testdata", "test_params.json")
	data, err := os.ReadFile(paramFile)
	if err != nil {
		t.Fatalf("テストパラメータファイルの読み込みに失敗: %v", err)
	}

	var params TestParams
	if err := json.Unmarshal(data, &params); err != nil {
		t.Fatalf("テストパラメータファイルの解析に失敗: %v", err)
	}

	return &params
}

// ポインタヘルパー関数
func StringPtr(s string) *string { return &s }
func IntPtr(i int) *int { return &i }
func BoolPtr(b bool) *bool { return &b }
```

**stories.go での使用例:**

```go
// アプローチB: test_params.json から読み込み（推奨）
func GetPostmanStoryStandardMethods(t *testing.T) testlib.Story {
	// test_params.json からパラメータを読み込み
	params := testdata.LoadTestParams(t)
	
	vars := map[string]any{
		// JSONから値を取得
		"domain_name":        params.BasicInfo.UpdateParams["domain_name"],
		"from_email_address": params.BasicInfo.UpdateParams["from_email_address"],
		"reply_email_address": params.BasicInfo.UpdateParams["reply_email_address"],
		"email":              params.Users.CreateParams["email"],
		"password":           params.Users.CreateParams["password"],
		"tenant_name":        params.Tenants.CreateParams["name"],
		// ... 他のパラメータ
	}

	return testlib.Story{
		Name:        "Postman Collection Story - Standard Methods",
		Description: "test_params.json から読み込んだパラメータを使用",
		Variables:   vars,
		Steps:       getPostmanStandardSteps(),
	}
}
```

**validation.go での使用例:**

```go
// test_params.json のデータを直接使用
func updateBasicInfoParams(variables map[string]interface{}) interface{} {
	// variables から取得（test_params.json 由来）
	domainName, _ := variables["domain_name"].(string)
	fromEmail, _ := variables["from_email_address"].(string)
	replyEmail, _ := variables["reply_email_address"].(string)

	return authapi.UpdateBasicInfoParam{
		DomainName:        domainName,
		FromEmailAddress:  fromEmail,
		ReplyEmailAddress: testdata.StringPtr(replyEmail),
	}
}
```

**test_params.json の構造例:**

```json
{
  "basic-info": {
    "updateParams": {
      "domain_name": "example.dev.saasus.io",
      "from_email_address": "noreply@example.com",
      "reply_email_address": "support@example.com"
    }
  },
  "users": {
    "createParams": {
      "email": "test@example.com",
      "password": "Passw0rd!"
    },
    "updateParams": {
      "email": "updated@example.com"
    }
  },
  "tenants": {
    "createParams": {
      "name": "Test Tenant",
      "back_office_staff_email": "admin@example.com",
      "attributes": {}
    }
  }
}
```

### 代替アプローチ: Variables 直接定義

シンプルなパラメータの場合は、ストーリーの`Variables`フィールドで直接定義することも可能です：

```go
// アプローチA: Variables 直接定義（シンプルなケース向け）
func GetPostmanStoryStandardMethods() testlib.Story {
	vars := map[string]any{
		"domain_name": "example.dev.saasus.io",
		"email": "test@example.com",
		"password": "Passw0rd!",
		// ... 他のパラメータ
	}

	return testlib.Story{
		Name:      "Postman Collection Story - Standard Methods",
		Variables: vars,
		Steps:     getPostmanStandardSteps(),
	}
}
```

### アプローチの選択基準

| 基準 | アプローチA（Variables直接定義） | アプローチB（test_params.json活用）★推奨 |
|------|--------------------------------|----------------------------------------|
| **パラメータの複雑さ** | シンプル（10個以下） | 複雑（10個以上、ネスト構造） |
| **再利用性** | 単一ストーリー | 複数のストーリー・テストで共有 |
| **保守性** | コード変更が必要 | JSONファイルの編集のみ |
| **可読性** | コード内で完結 | データとコードが分離 |
| **環境対応** | ハードコード | 環境別JSONファイルで切り替え可能 |
| **推奨ケース** | プロトタイプ、小規模テスト | 本番運用、大規模テスト |

### 実装のポイント

1. **ポインタヘルパー関数は必須**
   ```go
   import "github.com/saasus-platform/saasus-sdk-go/tests/e2e/<api-name>api/testdata"
   
   // オプショナルフィールドに値を設定
   ReplyEmailAddress: testdata.StringPtr(replyEmail),
   DisplayName: testdata.StringPtr(displayName),
   Enabled: testdata.BoolPtr(false),
   ```

2. **test_params.json の構造を API に合わせて調整**
   - test-param-generator が生成した構造をベースに、必要に応じてカスタマイズ
   - リソースごとに `createParams`, `updateParams` を分けて管理

3. **環境別の設定ファイル（オプション）**
   ```bash
   testdata/
   ├── test_params.json          # デフォルト（開発環境）
   ├── test_params.staging.json  # ステージング環境
   └── test_params.prod.json     # 本番環境
   ```

4. **LoadTestParams の拡張**
   ```go
   func LoadTestParams(t *testing.T) *TestParams {
       env := os.Getenv("TEST_ENV")
       if env == "" {
           env = "default"
       }
       
       var paramFile string
       if env == "default" {
           paramFile = filepath.Join("testdata", "test_params.json")
       } else {
           paramFile = filepath.Join("testdata", fmt.Sprintf("test_params.%s.json", env))
       }
       
       // ... ファイル読み込み処理
   }
   ```

> 🔁 生成された JSON のキー構成は API ごとに異なるため、`test_params.json` を確認しながらフィールドアクセスを更新してください。

> 🛠️ 自動生成ファイルは土台として提供されるため、パッケージ名の変更や必要に応じた修正を行った上でテストコードに統合してください。

> ⭐ **推奨**: 本番運用を見据えた実装では、アプローチB（test_params.json活用）を採用することで、長期的な保守性と拡張性が向上します。

### Step 6: ストーリー実装

**ファイル:** `tests/e2e/<api-name>api/stories.go` (続き)

Step 4 で作成したファイルにストーリー定義を追加します。Standard, WithResponse, Standard + WithBody, WithResponse + WithBody の 4 パターンをそれぞれ別ストーリーとして用意します（リクエストボディが必要なメソッドは `WithBody` / `WithBodyWithResponse` サフィックスが付いている想定）。

```go
// Get<ApiName>Stories returns test stories for <ApiName> API
func Get<ApiName>Stories() []testlib.Story {
	return []testlib.Story{
		{
			Name:        "Postman Collection Story - Standard Methods",
			Description: "Test story based on Postman collection using standard methods",
			Steps:       getPostmanStandardSteps(),
		},
		{
			Name:        "Postman Collection Story - Standard Methods With Body",
			Description: "Standard methods that send request bodies",
			Steps:       getPostmanStandardWithBodySteps(),
		},
		{
			Name:        "Postman Collection Story - WithResponse Methods",
			Description: "Test story based on Postman collection using WithResponse methods",
			Steps:       getPostmanWithResponseSteps(),
		},
		{
			Name:        "Postman Collection Story - WithResponse Methods With Body",
			Description: "WithResponse methods that send request bodies",
			Steps:       getPostmanWithResponseWithBodySteps(),
		},
	}
}

func getPostmanStandardSteps() []testlib.Step {
	return []testlib.Step{
		{
			Name:           "Step1: Get Initial Resource",
			ClientMethod:   "GetResource",
			Parameters:     nil,  // パラメータ不要の場合
			ExpectedStatus: 200,
			Validation: func(response interface{}) error {
				// カスタムバリデーション（必要に応じて）
				return nil
			},
		},
		{
			Name:           "Step2: Update Resource",
			ClientMethod:   "UpdateResource",
			Parameters:     updateResourceParams,  // パラメータ関数を指定
			ExpectedStatus: 200,
		},
		{
			Name:           "Step3: Delete Resource",
			ClientMethod:   "DeleteResource",
			Parameters:     deleteResourceParams,  // パラメータ関数を指定
			ExpectedStatus: 200,
		},
	}
}

func getPostmanStandardWithBodySteps() []testlib.Step {
	return []testlib.Step{
		{
			Name:           "Step1: Create Resource",
			ClientMethod:   "CreateResourceWithBody",
			ExpectedStatus: 201,
		},
		{
			Name:           "Step2: Update Resource",
			ClientMethod:   "UpdateResourceWithBody",
			ExpectedStatus: 200,
		},
		{
			Name:           "Step3: Delete Resource",
			ClientMethod:   "DeleteResource",
			ExpectedStatus: 200,
		},
	}
}

func getPostmanWithResponseSteps() []testlib.Step {
	return []testlib.Step{
		{
			Name:           "Step1: Get Initial Resource",
			ClientMethod:   "GetResourceWithResponse",
			ExpectedStatus: 200,
		},
		{
			Name:           "Step2: Update Resource",
			ClientMethod:   "UpdateResourceWithResponse",
			ExpectedStatus: 200,
		},
		{
			Name:           "Step3: Delete Resource",
			ClientMethod:   "DeleteResourceWithResponse",
			ExpectedStatus: 200,
		},
	}
}

func getPostmanWithResponseWithBodySteps() []testlib.Step {
	return []testlib.Step{
		{
			Name:           "Step1: Create Resource",
			ClientMethod:   "CreateResourceWithBodyWithResponse",
			ExpectedStatus: 201,
		},
		{
			Name:           "Step2: Update Resource",
			ClientMethod:   "UpdateResourceWithBodyWithResponse",
			ExpectedStatus: 200,
		},
		{
			Name:           "Step3: Delete Resource",
			ClientMethod:   "DeleteResourceWithResponse",
			ExpectedStatus: 200,
		},
	}
}
```

**設計ポイント:**

- Postman 分析レポートの「Test Flow」を参考にステップを定義
- Step 順序は Postman Collection の実行順序に従う
- ExpectedStatus は Postman のレスポンステストから取得
- `WithBody` / `WithBodyWithResponse` サフィックスが付くメソッドは WithBody 系ストーリーに、サフィックス無しのメソッドは Standard / WithResponse ストーリーに割り当てる（例: `UpdateStripeInfo` は Standard、`UpdateStripeInfoWithBody` は Standard With Body）
  - Standard Methods - 基本メソッドのみ
  - Standard WithBody Methods - WithBody バリアントのみ
  - WithResponse Methods - WithResponse メソッドのみ
  - WithBodyWithResponse Methods - WithBodyWithResponse バリアントのみ

> ⚠️ **重複操作に注意**
>
> Standard 系と WithResponse 系、さらに `WithBody` バリアントを同一ストーリーに混在させると、Postman 由来のシナリオをそのまま再現した際に状態が衝突する恐れがあります。例えば **Step 14** で `CreateVoteUserWithBody` を実行した直後に、同じストーリー内で **Step 15** として `CreateVoteUserWithBodyWithResponse` を再度呼び出すと、同一ユーザーによる重複投票と判定され **400 Bad Request**（重複投票エラー）が発生します。4 つのストーリーに分割して調整することで、こうした副作用を避けながら各メソッドを網羅的にテストできます。

### Step 5.5: クライアントメソッド使用確認

作成したストーリーで生成済みクライアントのメソッドを漏れなく使用しているか必ずチェックします（`ReturnInternalServerError` および `ReturnInternalServerErrorWithResponse` だけは除外可）。

#### メソッドカバレッジの確認方法

1. **Postman分析レポートの確認**
   - `tests/tools/postman-analyzer/reports/<api-name>/postman_analysis.json` を開き、`client_methods` の一覧を取得します
   - 例: Auth の場合は `tests/tools/postman-analyzer/reports/auth/postman_analysis.json`
   - `jq '.client_methods | length'` でメソッド数を確認できます

2. **生成コードのメソッド確認**
   - `generated/<api-name>api/client.gen.go` に定義された `ClientInterface` と `ClientWithResponsesInterface` のメソッド名を全て洗い出します
   - 以下のコマンドでメソッド数を確認できます：
     ```bash
     # ClientInterface のメソッド数
     sed -n '/^type ClientInterface interface/,/^}/p' generated/<api-name>api/client.gen.go | grep -c "("
     
     # ClientWithResponsesInterface のメソッド数
     sed -n '/^type ClientWithResponsesInterface interface/,/^}/p' generated/<api-name>api/client.gen.go | grep -c "("
     
     # 合計（ReturnInternalServerError系を除く）
     echo "$(( $(sed -n '/^type ClientInterface interface/,/^}/p' generated/<api-name>api/client.gen.go | grep -c "(") + $(sed -n '/^type ClientWithResponsesInterface interface/,/^}/p' generated/<api-name>api/client.gen.go | grep -c "(") - 2 ))"
     ```

3. **メソッド一覧関数の実装**
   - `stories.go` に `Get<ApiName>Methods()` 関数を実装し、全メソッドを配列で返します
   - Standard, WithBody, WithResponse, WithBodyWithResponse の順序で整理します
   - 例（Auth API）：
     ```go
     func GetAuthMethods() []string {
         return []string{
             // Standard メソッド
             "GetBasicInfo", "UpdateBasicInfo", "GetAuthInfo", "UpdateAuthInfo",
             // ... 他のStandardメソッド ...
             
             // WithBody メソッド
             "UpdateBasicInfoWithBody", "UpdateAuthInfoWithBody",
             // ... 他のWithBodyメソッド ...
             
             // WithResponse メソッド
             "GetBasicInfoWithResponse", "UpdateBasicInfoWithResponse",
             // ... 他のWithResponseメソッド ...
             
             // WithBodyWithResponse メソッド
             "UpdateBasicInfoWithBodyWithResponse", "UpdateAuthInfoWithBodyWithResponse",
             // ... 他のWithBodyWithResponseメソッド ...
         }
     }
     ```

4. **メソッドカバレッジ検証関数の実装**
   - `stories.go` に `VerifyMethodCoverage()` 関数を実装し、未使用メソッドを検出します
   - この関数は `Get<ApiName>Methods()` で定義された全メソッドが、4つのストーリーのいずれかで少なくとも1回呼び出されることを確認します
   - 例（Auth API）：
     ```go
     func VerifyMethodCoverage() error {
         allMethods := GetAuthMethods()
         usedMethods := make(map[string]bool, len(allMethods))
         
         stories := GetAuthStories()
         for _, story := range stories {
             for _, step := range story.Steps {
                 usedMethods[step.ClientMethod] = true
             }
         }
         
         var missingMethods []string
         for _, method := range allMethods {
             if !usedMethods[method] {
                 missingMethods = append(missingMethods, method)
             }
         }
         
         if len(missingMethods) > 0 {
             sort.Strings(missingMethods)
             return fmt.Errorf("以下のメソッドがストーリーでカバーされていません (%d/%d 個): %v",
                 len(missingMethods), len(allMethods), missingMethods)
         }
         
         return nil
     }
     ```

5. **メインテストでの検証**
   - メインテストファイル（`<api-name>_test.go`）で `VerifyMethodCoverage()` を呼び出します
   - 例：
     ```go
     func TestAuthAPIE2E(t *testing.T) {
         // ... 既存のコード ...
         
         // メソッドカバレッジの検証
         if err := VerifyMethodCoverage(); err != nil {
             t.Errorf("❌ Method coverage verification failed: %v", err)
         }
         
         // ... 既存のコード ...
     }
     ```

6. **ストーリーへのメソッド割り当て**
   - Step 5 で定義した 4 種類のストーリー（Standard / Standard With Body / WithResponse / WithResponse With Body）の `ClientMethod` が、(2) で洗い出したメソッドのうち `ReturnInternalServerError` および `ReturnInternalServerErrorWithResponse` を除く全てを少なくとも一度ずつ呼び出していることを確認します
   - `WithBody` / `WithBodyWithResponse` サフィックス付きメソッドは WithBody 系ストーリーへ、サフィックス無しメソッドは Standard / WithResponse ストーリーへ割り当ててください
   - 例: `UpdateStripeInfo` → Standard、`UpdateStripeInfoWithBody` → Standard With Body
   - 呼び出していないメソッドがあればステップを追加するか、既存ステップを修正してください

7. **テストパラメータの確認**
   - `testdata/test_params.json` に自動生成された値と stories / validation で参照しているフィールドが噛み合っているか確認します

8. **レビュー時の確認**
   - レビュー時には、上記確認結果（使用状況のメモやチェックリストなど）を PR コメントに残すことを推奨します
   - メソッド数の一致を確認した結果を記載します（例: "generated code: 250メソッド、GetAuthMethods(): 250メソッド、Postman分析: 248メソッド"）

#### メソッドカバレッジの重要性

- **完全性の保証**: 全メソッドがテストされることで、APIの完全性が保証されます
- **リグレッション防止**: 新しいメソッドが追加された際に、テストの追加漏れを防ぎます
- **ドキュメント化**: メソッド一覧がコードとして明示的に定義されます
- **自動検証**: テスト実行時に自動的にカバレッジが検証されます

### Step 6: メインテスト実装

**ファイル:** `tests/e2e/<api-name>api/<api-name>_test.go`

```go
package <api-name>api

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/saasus-platform/saasus-sdk-go/modules/<api-name>"
	"github.com/saasus-platform/saasus-sdk-go/tests/testlib"
)

func Test<ApiName>APIE2E(t *testing.T) {
	fmt.Println("🚀 Starting SaaSus <ApiName> API E2E Tests")
	fmt.Println("========================================")

	// Load environment variables from .env file
	if err := godotenv.Load("../../../.env"); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	// Create client
	client, err := <api-name>.<ApiName>WithResponse()
	if err != nil {
		t.Fatalf("❌ Failed to create <api-name> client: %v", err)
	}

	engine := testlib.NewE2EEngine(client, Get<ApiName>Methods())

	// Parse command line arguments
	if err := engine.Config.ParseArgs(os.Args[1:]); err != nil {
		t.Fatalf("❌ Error parsing arguments: %v", err)
	}

	// Validate configuration
	if err := engine.Config.Validate(); err != nil {
		t.Fatalf("❌ Configuration error: %v", err)
	}

	// Get stories
	stories := Get<ApiName>Stories()

	// Execute stories
	results := engine.ExecuteStories(stories)

	// Print results
	engine.PrintResults(results)

	// Check if all stories passed
	for _, result := range results {
		if result.Status != testlib.StatusPassed {
			t.Errorf("❌ Story '%s' failed: %v", result.StoryName, result.Error)
		}
	}

	fmt.Println("\n✅ All <ApiName> API E2E tests completed!")
}
```

### Step 7: （オプション）将来の拡張に備えた設定

現在、スナップショット機能は使用していませんが、将来の拡張に備えて設定ファイルを残すことができます。

### Step 8: テスト実行

```bash
# 8-1. 基本的なE2Eテスト実行
cd tests/e2e/<api-name>api
go test -v

# 8-2. 詳細ログ付きテスト実行
LOG_LEVEL=debug go test -v

# 8-3. 特定のストーリーのみ実行
go test -v -stories="Postman Collection Story - Standard Methods"
```

### Step 9: テスト結果の確認

```bash
# 9-1. テスト結果を確認
# 全てのストーリーがPASSEDであることを確認

# 9-2. カバレッジを確認
# テスト実行時に表示されるカバレッジ情報を確認
```

### Step 10: ドキュメントの更新

テスト実装が完了したら、READMEファイルを作成または更新します。

**ファイル:** `tests/e2e/<api-name>api/README.md`

READMEには以下の情報を含めます：

1. **概要**
   - API の説明
   - テストスイートの目的

2. **ファイル構成**
   - 各ファイルの役割を説明

3. **テストストーリー**
   - 4つのストーリーの説明
   - カバーするメソッドカテゴリ

4. **前提条件**
   - 環境変数の設定方法
   - テストパラメータの準備

5. **テストの実行**
   - 基本的な実行方法
   - 詳細ログ付き実行
   - 特定のストーリーのみ実行
   - スナップショットテスト（将来の拡張）

6. **メソッドカバレッジ**
   - カバーするメソッド数
   - ClientInterface と ClientWithResponsesInterface の内訳
   - 除外されるメソッド（ReturnInternalServerError系）
   - 自動検証の説明

7. **トラブルシューティング**
   - よくある問題と解決策

8. **開発ガイド**
   - 新しいパラメータ関数の追加方法
   - 新しいストーリーの追加方法
   - バリデーションルールの追加方法

9. **参考資料**
   - 関連ドキュメントへのリンク

**メソッドカバレッジセクションの例（Auth API）:**

```markdown
## メソッドカバレッジ

このテストスイートは、Auth API の全メソッド（250メソッド）をカバーします。`ReturnInternalServerError` および `ReturnInternalServerErrorWithResponse` メソッドのみ除外されます。

### カバレッジ詳細

- **ClientInterface**: 126メソッド（標準メソッド + WithBody メソッド）
- **ClientWithResponsesInterface**: 126メソッド（WithResponse メソッド + WithBodyWithResponse メソッド）
- **合計**: 252メソッド（ReturnInternalServerError系2つを除いて250メソッド）

メソッドカバレッジは `VerifyMethodCoverage()` 関数によって自動的に検証され、未使用のメソッドがある場合はエラーが報告されます。
```

**参考実装:**
- `tests/e2e/authapi/README.md` - Auth API の完全な例
- `tests/e2e/billingapi/README.md` - Billing API の例
- `tests/e2e/pricingapi/README.md` - Pricing API の例

## ディレクトリ構造

作成後の構造:

```
tests/
├── e2e/
│   ├── <api-name>api/                    # 新規作成
│   │   ├── <api-name>_test.go           # メインテストファイル
│   │   ├── stories.go                    # ストーリー定義とメソッド一覧
│   │   ├── validation.go                # パラメータ/バリデーション
│   │   ├── story_snapshot_config.go     # スナップショット設定（オプション）
│   │   └── testdata/                    # 自動生成されたテストパラメータ
│   │       ├── test_params.json
│   │       └── <api-name>_test_helpers.go
│   ├── billingapi/                       # 参考実装
│   │   ├── billing_test.go
│   │   ├── stories.go                    # ストーリー定義とメソッド一覧
│   │   └── story_snapshot_config.go
│   └── snapshot/                         # 共通スナップショット基盤
│       ├── main.go                       # CLI
│       ├── <api-name>/                   # モジュールごとの出力
│       │   ├── story_snapshots/          # キャプチャデータ
│       │   ├── story_comparisons/        # 比較結果
│       │   ├── story_validations/        # 検証結果
│       │   └── story_reports/            # HTMLレポート
│       └── docs/                         # スナップショット関連ドキュメント
├── testlib/                              # 共通テストライブラリ
│   ├── types.go
│   ├── config.go
│   ├── e2e.go
│   └── snapshot/                         # スナップショット機能
└── tools/
    └── postman-analyzer/                 # Postman分析ツール
        ├── analyzer.go
        ├── collections/
        │   └── SaaSus <Service Name> API Schema.postman_collection.json  # 新規追加
        └── reports/
            └── <api-name>/               # 分析結果
                ├── postman_analysis.json
                └── postman_analysis_report.txt
```

## 各ステップの詳細

### Postman Collection の要件

✅ **必須項目:**

- リクエスト URL
- HTTP メソッド
- 期待ステータスコード

✅ **推奨項目:**

- テスト名
- テスト説明
- レスポンステスト
- 前処理/後処理スクリプト

### ストーリー設計のベストプラクティス

1. **ストーリーの分割:**

   - Standard Methods 用ストーリー
   - WithResponse Methods 用ストーリー
   - エラーケース用ストーリー（オプション）

2. **ステップの命名:**

   ```go
   "Step1: Get Initial Resource"
   "Step2: Create New Resource"
   "Step3: Update Resource"
   "Step4: Delete Resource"
   ```

3. **状態管理:**

   ```go
   {
       Name:         "Step2: Create Resource",
       ClientMethod: "CreateResource",
       StateUpdate: func(response interface{}, variables map[string]interface{}) error {
           // レスポンスからIDを抽出して保存
           variables["resource_id"] = extractID(response)
           return nil
       },
   }
   ```

4. **バリデーション:**

   ```go
   Validation: func(response interface{}) error {
       if response == nil {
           return fmt.Errorf("response is nil")
       }
       // カスタム検証ロジック
       return nil
   }
   ```

5. **パラメータ関数の実装:**

   - パラメータ関数は`variables`を受け取り、具体的な値を返す
   - 関数型を返さない（重要！）
   - 単一パラメータの場合は直接値を返す
   - 複数パラメータの場合は構造体を返す
   - WithBody メソッドの場合は、追加パラメータ + ContentType + Body の順序で構造体を定義
   - **構造体を直接作成し、不要な JSON 経由の変換を避ける**（パフォーマンスと型安全性の向上）

   ```go
   // ✅ 単一パラメータ
   func getResourceParams(variables map[string]interface{}) interface{} {
       return variables["resource_id"].(string)
   }

   // ✅ 複数パラメータ
   func getResourceDetailParams(variables map[string]interface{}) interface{} {
       return struct {
           ResourceId string
           DetailId   string
       }{
           ResourceId: variables["resource_id"].(string),
           DetailId:   variables["detail_id"].(string),
       }
   }

   // ✅ 推奨: 構造体を直接作成（型安全、パフォーマンス向上）
   func createResourceParams(variables map[string]interface{}) interface{} {
       return api.CreateResourceJSONRequestBody{
           Name:        "Resource Name",
           Description: "Resource Description",
           Tags:        []string{"tag1", "tag2"},
       }
   }

   // ✅ Union型の場合: メソッドを使用して型変換
   func createUnionTypeParams(variables map[string]interface{}) interface{} {
       fixedUnit := api.FixedUnitForSave{
           Name:        "Fixed Unit",
           UnitAmount:  uint64(1000),
           Currency:    api.JPY,
       }
       var body api.UnitForSave
       _ = body.FromFixedUnitForSave(fixedUnit)
       return body
   }

   // ❌ 非推奨: JSON経由の変換（不要なオーバーヘッド）
   func createResourceParamsOld(variables map[string]interface{}) interface{} {
       bodyData := map[string]any{
           "name":        "Resource Name",
           "description": "Resource Description",
           "tags":        []string{"tag1", "tag2"},
       }
       bodyBytes, _ := json.Marshal(bodyData)
       var body api.CreateResourceJSONRequestBody
       _ = json.Unmarshal(bodyBytes, &body)
       return body
   }

   // ✅ WithBodyメソッド
   func updateResourceWithBodyParams(variables map[string]interface{}) interface{} {
       resourceID := variables["resource_id"].(string)
       param := api.UpdateResourceParam{
           Name: "Updated Name",
       }
       body, _ := json.Marshal(param)

       return struct {
           ResourceId  string    // 追加パラメータ
           ContentType string    // ContentType
           Body        io.Reader // Body
       }{
           ResourceId:  resourceID,
           ContentType: "application/json",
           Body:        bytes.NewReader(body),
       }
   }
   ```

### メソッド一覧の作成方法

**方法 1: Postman 分析レポートから取得**

```bash
cd tests/tools/postman-analyzer
go run analyzer.go -collection "collections/SaaSus <Service Name> API Schema.postman_collection.json" -output reports/<api-name>
cat reports/<api-name>/postman_analysis_report.txt | grep -A 20 "=== Covered Methods ==="
```

**方法 2: 生成されたコードから取得**

```bash
grep "^func.*Client)" generated/<api-name>api/client.gen.go | \
  sed 's/func (c \*Client) \(.*\)(.*/\1/'
```

**作成先:** `tests/e2e/<api-name>api/stories.go` 内の `Get<ApiName>Methods()` 関数

### test-param-generator の活用ポイント

- `-verbose` で実行すると、マッピング経路や iterationData との突合情報が表示される
- 生成された `test_params.json` のキーは `client_methods` と揃うため、Stories のステップ名と突き合わせて抜け漏れを検出しやすい
- 変換に失敗したパラメータは診断レポート（`diagnostic_report_*.txt`）に出力されるので、`tests/tools/INTEGRATION.md` の「診断レポートの活用」を参照して手動で補完する
- `-overwrite` フラグを付けない限り既存ファイルは保護される。再生成時は差分を確認した上でフラグを付与する
- カスタムマッピングが必要な場合は `mapping_rules.json` を用意し、`-rules` オプションで指定する

### 環境変数設定

```bash
# .envファイルを作成（プロジェクトルート）
cat > .env << 'EOF'
# 必須
SAASUS_SAAS_ID=your-saas-id
SAASUS_API_KEY=your-api-key
SAASUS_SECRET_KEY=your-secret-key

# オプション
LOG_LEVEL=info
E2E_LOG_LEVEL=debug
SAASUS_BASE_URL=https://api.dev.saasus.io/v1
EOF
```

## チェックリスト

### Step 1-2: ツールチェーン準備 ✓

- [ ] `./generate_all.sh` を実行し、OpenAPI / client.gen.go を最新化
- [ ] `openapi-resolver -spec api/<api-name>api.yml -list` で operationID を確認（必要に応じて）
- [ ] Postman Collection を JSON 形式でエクスポートし、`tests/tools/postman-analyzer/collections/` に配置
- [ ] `postman-analyzer` を実行し、`postman_analysis.json` / `postman_analysis_report.txt` を確認
- [ ] `test-param-generator` を実行し、`tests/e2e/<api-name>api/testdata/` にファイルが生成されたことを確認

### Step 3-4: テストファイル作成 ✓

- [ ] テストディレクトリ作成
- [ ] `<api-name>_test.go`作成
- [ ] `stories.go`作成
- [ ] メソッド一覧を`stories.go`に実装

### Step 5: ストーリー実装 ✓

- [ ] Standard Methods ストーリー作成
- [ ] WithResponse Methods ストーリー作成
- [ ] ステップの順序を確認
- [ ] ExpectedStatus を設定
- [ ] バリデーション追加（必要に応じて）

### Step 6-7: テスト完成 ✓

- [ ] メインテスト関数実装
- [ ] クライアント初期化コード追加
- [ ] エラーハンドリング実装

### Step 8-9: 実行確認 ✓

- [ ] 基本的な E2E テスト実行
- [ ] 全ストーリーが成功
- [ ] カバレッジ確認

### Step 10: ドキュメント更新 ✓

- [ ] README.md を作成または更新
- [ ] メソッドカバレッジ情報を記載
- [ ] テスト実行方法を記載
- [ ] トラブルシューティングを記載
- [ ] 開発ガイドを記載

### 最終確認 ✓

- [ ] 全メソッドがカバーされている（`VerifyMethodCoverage()` で自動検証）
- [ ] メソッド数が一致している（generated code, GetMethods(), Postman分析）
- [ ] テストログが適切に出力される
- [ ] エラーケースも適切に処理される
- [ ] ドキュメントが更新されている
- [ ] コードが整形されている（`gofmt -l` でエラーなし）
- [ ] 静的解析でエラーがない（`go vet` でエラーなし）
- [ ] ビルドが成功する（`go build` でエラーなし）

## 参考リソース

### ドキュメント

- `tests/testlib/README.md` - testlib 使用方法
- `tests/tools/README.md` - ツール概要
- `tests/tools/postman-analyzer/README.md` - アナライザー使用方法
- `tests/e2e/snapshot/docs/guides/STORY_SNAPSHOT_CLI_USAGE.md` - スナップショット CLI 使用方法
- `tests/e2e/snapshot/docs/guides/STORY_SNAPSHOT_CONFIG_GUIDE.md` - スナップショット設定ガイド

### 参考実装

- `tests/e2e/billingapi/` - Billing API 実装例
- `tests/e2e/billingapi/billing_test.go` - メインテスト
- `tests/e2e/billingapi/stories.go` - ストーリー定義とメソッド一覧
- `tests/e2e/pricingapi/` - Pricing API 実装例（最新）

### テスト実行例

```bash
# 詳細ログでテスト実行
LOG_LEVEL=debug go test ./tests/e2e/<api-name>api/

# 特定のストーリーのみ
go test -v -stories="Postman Collection Story - Standard Methods"
```

## トラブルシューティング

### よくある問題

**問題 0: test-param-generator が失敗する**

```
❌ 設定の検証に失敗: 出力ディレクトリが指定されていません
```

**解決策:**

- `-module` `-report` `-output` をすべて指定しているか確認
- `tests/tools/INTEGRATION.md` の「デバッグ方法」を参照し、`-verbose` / `-debug` オプションで詳細ログと診断レポートを確認
- 既存ファイルを上書きしたい場合は `-overwrite` を付与

**問題 1: クライアント作成エラー**

```
❌ Failed to create client: configuration error
```

**解決策:** 環境変数が正しく設定されているか確認

**問題 2: メソッドが見つからない**

```
❌ Error: method 'MethodName' not found on client
```

**解決策:** `<api-name>_methods.go`のメソッド名が正しいか確認

**問題 3: スナップショット比較エラー**

```
Error: snapshot file not found
```

**解決策:** `--mode=capture`でスナップショットを先にキャプチャ

**問題 4: 期待ステータスコード不一致**

```
❌ Status code mismatch: expected 200, got 404
```

**解決策:** API の仕様を確認し、ExpectedStatus を修正

**問題 5: パラメータ型変換エラー**

```
❌ reflect: CallSlice using func(map[string]interface {}) interface {} as type string
```

**解決策:** パラメータ関数が関数型を返している可能性があります。パラメータ関数は具体的な値（文字列、構造体など）を返す必要があります。

```go
// ❌ 間違った実装
func getResourceParams(variables map[string]interface{}) interface{} {
    return func(vars map[string]interface{}) interface{} {
        return vars["resource_id"].(string)
    }
}

// ✅ 正しい実装
func getResourceParams(variables map[string]interface{}) interface{} {
    return variables["resource_id"].(string)
}
```

**問題 6: WithBody メソッドのパラメータエラー**

```
❌ Strategy Enhanced failed: parameters must be a struct with ContentType and Body fields
```

**解決策:** WithBody メソッドのパラメータ構造体に、追加パラメータ（feedbackId など）、`ContentType`、`Body`フィールドが正しい順序で含まれているか確認してください。

```go
// ✅ 正しい実装（UpdateFeedbackWithBodyの例）
return struct {
    FeedbackId  string    // 追加パラメータを最初に
    ContentType string    // ContentTypeフィールド
    Body        io.Reader // Bodyフィールド
}{
    FeedbackId:  feedbackID,
    ContentType: "application/json",
    Body:        bytes.NewReader(body),
}
```

**問題 7: JSON シリアライズエラー**

```
❌ Failed to save story snapshot: json: unsupported type: func(map[string]interface {}) interface {}
```

**解決策:** この問題は通常、内部的に自動修正されます。もし発生した場合は、パラメータ関数が正しく実装されているか確認してください（問題 5 を参照）。

**問題 8: スナップショットモードでのレスポンスボディ読み取りエラー**

```
❌ failed to read response body: context canceled
❌ failed to read response body: http: read on closed response body
```

**原因:** スナップショットモードでは、エンジンがレスポンスボディを先に読み取る場合があります。その後、`extractResponseMap`などのヘルパー関数が再度ボディを読み取ろうとするとエラーが発生します。

**解決策（2025年1月更新）:** この問題は`tests/testlib/snapshot/engine.go`と`tests/testlib/snapshot/capture.go`の改善により自動的に解決されています。エンジンがStateUpdate関数にキャプチャ済みのJSONデータを直接渡すため、レスポンスボディの二重読み取りは発生しません。

**重要な改善（2025年1月）:** 

1. **エンジン側の改善**: `tests/testlib/snapshot/engine.go`の`executeStepWithResponseCapture`メソッドが更新され、StateUpdate処理でキャプチャされたJSONデータを自動的に優先使用するようになりました。これにより、StateUpdate関数内で`extractResponseMap`を呼び出す必要がなくなり、レスポンスボディの二重読み取りエラーが自動的に回避されます。

2. **キャプチャ処理の簡素化**: `tests/testlib/snapshot/capture.go`の`captureHTTPResponse`メソッドが簡素化され、以下の改善が行われました：
   - HTTPレスポンスボディの読み取りと再設定処理を削除
   - 不要なボディ操作を排除し、シンプルな実装に変更
   - パフォーマンスの向上とコードの保守性改善
   - レスポンスボディの二重読み取りリスクを根本的に解消

```go
// StateUpdate関数の実装例（最新版）
StateUpdate: func(response any, variables map[string]any) error {
    // エンジンが自動的にキャプチャされたJSONDataを渡すため、
    // extractResponseMapを呼び出す必要はありません
    respMap, ok := response.(map[string]interface{})
    if !ok {
        return fmt.Errorf("unexpected response type: %T", response)
    }
    
    if id, ok := respMap["id"].(string); ok {
        variables["resource_id"] = id
    }
    return nil
}
```

**注意（2025年1月更新）:** 最新の`tests/testlib/snapshot/engine.go`では、StateUpdate関数に渡されるレスポンスは自動的にキャプチャされたJSONDataに変換されます。そのため、StateUpdate内で`extractResponseMap`を呼び出す必要はなく、直接`map[string]interface{}`として型アサーションできます。

また、`tests/testlib/snapshot/capture.go`の`captureHTTPResponse`メソッドも簡素化され、以下の改善が行われています：

- **HTTPレスポンスボディの処理を簡素化**: ボディの読み取りと再設定処理を削除
- **不要な操作を排除**: レスポンスボディを直接操作せず、リフレクションベースの抽出のみに依存
- **パフォーマンスの向上**: 不要なI/O操作を削減
- **コードの保守性改善**: シンプルな実装により、バグのリスクを低減
- **二重読み取りリスクの解消**: ボディ操作を行わないため、StateUpdateでの読み取りエラーが発生しない

**実装のポイント:**

- リフレクションを使用してレスポンス構造を動的に検査
- スナップショットエンジンが提供する`JSONData`フィールドを検出
- 通常のレスポンス型とスナップショットデータの両方に対応
- 型安全性を保ちながらフォールバック処理を実装

**キャプチャ処理の簡素化（2025年1月）:**

`tests/testlib/snapshot/capture.go`の`captureHTTPResponse`メソッドから以下の処理が削除されました：

- HTTPレスポンスボディの読み取り（`io.ReadAll`）
- ボディのクローズと再設定（`io.NopCloser`）
- ボディバイト数のデバッグログ

この簡素化により：
- レスポンスボディを直接操作しないため、StateUpdateでの読み取りエラーが発生しない
- コードがシンプルになり、保守性が向上
- 不要なI/O操作が削減され、パフォーマンスが改善

この改善により、通常の E2E テストとスナップショットテストの両方で`extractResponseMap`が正しく動作します。

**最新の推奨パターン（2025年1月）:**

`tests/testlib/snapshot/engine.go`の改善により、StateUpdate関数では`extractResponseMap`を使用せず、直接レスポンスを`map[string]interface{}`として扱うことが推奨されます：

```go
// ✅ 推奨: 直接型アサーション（エンジンが自動変換）
StateUpdate: func(response any, variables map[string]any) error {
    respMap, ok := response.(map[string]interface{})
    if !ok {
        return fmt.Errorf("unexpected response type: %T", response)
    }
    
    if id, ok := respMap["id"].(string); ok {
        variables["resource_id"] = id
    }
    return nil
}

// ❌ 非推奨: extractResponseMapの使用（不要なオーバーヘッド）
StateUpdate: func(response any, variables map[string]any) error {
    respMap, err := extractResponseMap(response)
    if err != nil {
        return err
    }
    // ...
}
```

**デバッグログの活用:**

デバッグログは以下の場合に役立ちます：

- スナップショットモードでのレスポンス処理の追跡
- 新しいレスポンス型の追加時の動作確認
- エラー発生時の原因特定

本番環境では、これらのログを削除するか、ログレベルで制御することを推奨します。

**注意:** この実装には`reflect`パッケージのインポートが必要です：

```go
import (
    "reflect"
    // ... 他のインポート
)
```

## 補足: 複数 API モジュール対応例

### Auth API (実装済み)

```bash
# Postman Collection配置
cp /path/to/exported.json "tests/tools/postman-analyzer/collections/SaaSus Auth API Schema.postman_collection.json"

# 分析実行
cd tests/tools/postman-analyzer
go run analyzer.go -collection "collections/SaaSus Auth API Schema.postman_collection.json" -output reports/auth

# テストパラメータ生成
cd ../test-param-generator
./run.sh -report ../postman-analyzer/reports/auth/postman_analysis.json -module auth -output ../../e2e/authapi/testdata -verbose
cd ../../..
# （ここでプロジェクトルートに戻る）

# テスト作成
mkdir -p tests/e2e/authapi
# ... ファイル作成 ...

# テスト実行
cd tests/e2e/authapi
go test -v

# 詳細ログ付き実行
LOG_LEVEL=debug go test -v

# 特定のストーリーのみ実行
go test -v -stories="Postman Collection Story - Standard Methods"
go test -v -stories="Postman Collection Story - Standard Methods With Body"
go test -v -stories="Postman Collection Story - WithResponse Methods"
go test -v -stories="Postman Collection Story - WithResponse Methods With Body"
```

**実装済みファイル:**

- `auth_test.go` - メインテストファイル
- `stories.go` - ストーリー定義とメソッド一覧
- `validation.go` - パラメータ関数とバリデーション
- `story_snapshot_config.go` - スナップショット設定（将来の拡張用）
- `testdata/test_params.json` - テストパラメータ（自動生成）
- `testdata/auth_test_helpers.go` - ヘルパー関数
- `README.md` - テストドキュメント

**特徴:**

- 全 Auth API メソッドをカバー（250メソッド）
- 4 つのストーリー（Standard, Standard With Body, WithResponse, WithResponse With Body）
- メソッドカバレッジの自動検証（`VerifyMethodCoverage()` 関数）
- testdata/test_params.json からテストデータを読み込み
- 全メソッドバリアントを網羅的にカバー：
  - Standard メソッド（基本的なHTTPレスポンス）
  - WithBody メソッド（io.Readerボディパラメータ）
  - WithResponse バリアント（詳細レスポンス）
  - WithBodyWithResponse バリアント（io.Reader + 詳細レスポンス）

**メソッドカバレッジ:**

- **ClientInterface**: 126メソッド（標準メソッド + WithBody メソッド）
- **ClientWithResponsesInterface**: 126メソッド（WithResponse メソッド + WithBodyWithResponse メソッド）
- **合計**: 252メソッド（ReturnInternalServerError系2つを除いて250メソッド）
- **除外**: `ReturnInternalServerError`, `ReturnInternalServerErrorWithResponse`

**メソッドカバレッジ検証:**

```go
// stories.go に実装
func GetAuthMethods() []string {
    return []string{
        // Standard メソッド（63メソッド）
        "GetBasicInfo", "UpdateBasicInfo", "GetAuthInfo", "UpdateAuthInfo",
        // ... 他のStandardメソッド ...
        
        // WithBody メソッド（31メソッド）
        "UpdateBasicInfoWithBody", "UpdateAuthInfoWithBody",
        // ... 他のWithBodyメソッド ...
        
        // WithResponse メソッド（125メソッド）
        "GetBasicInfoWithResponse", "UpdateBasicInfoWithResponse",
        // ... 他のWithResponseメソッド ...
        
        // WithBodyWithResponse メソッド（31メソッド）
        "UpdateBasicInfoWithBodyWithResponse", "UpdateAuthInfoWithBodyWithResponse",
        // ... 他のWithBodyWithResponseメソッド ...
    }
}

func VerifyMethodCoverage() error {
    allMethods := GetAuthMethods()
    usedMethods := make(map[string]bool, len(allMethods))
    
    stories := GetAuthStories()
    for _, story := range stories {
        for _, step := range story.Steps {
            usedMethods[step.ClientMethod] = true
        }
    }
    
    var missingMethods []string
    for _, method := range allMethods {
        if !usedMethods[method] {
            missingMethods = append(missingMethods, method)
        }
    }
    
    if len(missingMethods) > 0 {
        sort.Strings(missingMethods)
        return fmt.Errorf("以下のメソッドがストーリーでカバーされていません (%d/%d 個): %v",
            len(missingMethods), len(allMethods), missingMethods)
    }
    
    return nil
}
```

**テスト実行例:**

```bash
# Auth APIテストディレクトリから実行
cd tests/e2e/authapi

# 基本的なE2Eテスト実行
go test -v

# 詳細ログ付き実行
LOG_LEVEL=debug go test -v

# 特定のストーリーのみ実行
go test -v -stories="Postman Collection Story - Standard Methods"

# メソッドカバレッジの確認
# テスト実行時に自動的に検証され、未使用メソッドがあればエラーが報告されます
```

**メソッド数の確認コマンド:**

```bash
# generated/authapi/client.gen.go のメソッド数
echo "ClientInterface のメソッド数:"
sed -n '/^type ClientInterface interface/,/^}/p' generated/authapi/client.gen.go | grep -c "("

echo "ClientWithResponsesInterface のメソッド数:"
sed -n '/^type ClientWithResponsesInterface interface/,/^}/p' generated/authapi/client.gen.go | grep -c "("

echo "合計（ReturnInternalServerError系を除く）:"
echo "$(( $(sed -n '/^type ClientInterface interface/,/^}/p' generated/authapi/client.gen.go | grep -c "(") + $(sed -n '/^type ClientWithResponsesInterface interface/,/^}/p' generated/authapi/client.gen.go | grep -c "(") - 2 ))"

# Postman 分析レポートのメソッド数
jq '.client_methods | length' tests/tools/postman-analyzer/reports/auth/postman_analysis.json

# GetAuthMethods() で定義されたメソッド数
sed -n '/^func GetAuthMethods/,/^}/p' tests/e2e/authapi/stories.go | grep -o '"[^"]*"' | wc -l
```

### Pricing API

```bash
# Postman Collection配置
cp /path/to/exported.json "tests/tools/postman-analyzer/collections/SaaSus Pricing API Schema.postman_collection.json"

# 分析実行
cd tests/tools/postman-analyzer
go run analyzer.go -collection "collections/SaaSus Pricing API Schema.postman_collection.json" -output reports/pricing

# テストパラメータ生成
cd ../test-param-generator
./run.sh -report ../postman-analyzer/reports/pricing/postman_analysis.json -module pricing -output ../../e2e/pricingapi/testdata -verbose
cd ../../..
# （ここでプロジェクトルートに戻る）

# テスト作成
mkdir -p tests/e2e/pricingapi
# ... ファイル作成 ...

# テスト実行
go test ./tests/e2e/pricingapi/
```

### Communication API

```bash
# Postman Collection配置
cp /path/to/exported.json "tests/tools/postman-analyzer/collections/SaaSus Communication API Schema.postman_collection.json"

# 分析実行
cd tests/tools/postman-analyzer
go run analyzer.go -collection "collections/SaaSus Communication API Schema.postman_collection.json" -output reports/communication

# テストパラメータ生成
cd ../test-param-generator
./run.sh -report ../postman-analyzer/reports/communication/postman_analysis.json -module communication -output ../../e2e/communicationapi/testdata -verbose
cd ../../..
# （ここでプロジェクトルートに戻る）

# テスト作成
mkdir -p tests/e2e/communicationapi
# ... ファイル作成 ...

# テスト実行
go test ./tests/e2e/communicationapi/
```

### Integration API (実装済み)

```bash
# テストディレクトリ作成
mkdir -p tests/e2e/integrationapi

# E2Eテスト実行
cd tests/e2e/integrationapi
go test -v

# 詳細ログ付き実行
LOG_LEVEL=debug go test -v

# 特定のストーリーのみ実行
go test -v -stories="Postman Collection Story - Standard Methods"
go test -v -stories="Postman Collection Story - Standard Methods With Body"
go test -v -stories="Postman Collection Story - WithResponse Methods"
go test -v -stories="Postman Collection Story - WithResponse Methods With Body"

# スナップショットテスト実行（integration_test.goに統合済み）

# 通常のE2Eテスト（スナップショットモード未指定）
go test -v

# 1. スナップショットキャプチャ
go test -v -snapshot-mode=capture

# 2. スナップショット比較
go test -v -snapshot-mode=comparison

# 3. スナップショットレポート生成
go test -v -snapshot-mode=reporting

# 4. 統合実行（全フェーズ）
go test -v -snapshot-mode=integrated

# 生成されたファイルを確認
ls -la ../snapshot/integration/story_snapshots/     # キャプチャデータ
ls -la ../snapshot/integration/story_comparisons/   # 比較結果
ls -la ../snapshot/integration/story_validations/   # 検証結果
ls -la ../snapshot/integration/story_reports/       # HTMLレポート
```

**実装済みファイル:**

- `integration_test.go` - E2Eテスト + スナップショットテスト（統合済み）
- `stories.go` - ストーリー定義とメソッド一覧
- `validation.go` - パラメータ関数とバリデーション
- `story_snapshot_config.go` - スナップショット設定（完全実装済み）
- `testdata/test_params.json` - テストパラメータ（自動生成）

**実装済み機能:**

`integration_test.go`には以下が実装されています：
1. `-snapshot-mode`フラグの定義
2. `getSnapshotConfigForMode()`関数の実装
3. スナップショットモードとE2Eモードの分岐処理
4. `snapshot.NewE2EEngineWithSnapshot`の使用

**特徴:**

- 全Integration APIメソッドをカバー（EventBridge設定、イベント送信）
- 4つのストーリー（Standard, Standard With Body, WithResponse, WithResponse With Body）
- 各ストーリー10ステップ、合計14メソッドカバレッジ
- testdata/test_params.jsonからテストデータを読み込み
- 全メソッドバリアントを網羅的にカバー：
  - Standard メソッド（基本的なHTTPレスポンス）
  - WithBody メソッド（io.Readerボディパラメータ）
  - WithResponse バリアント（詳細レスポンス）
  - WithBodyWithResponse バリアント（io.Reader + 詳細レスポンス）

**スナップショット設定:**

- 4つのモード別設定関数を提供：
  - `GetIntegrationSnapshotConfigForCapture()` - キャプチャ専用
  - `GetIntegrationSnapshotConfigForComparison()` - 比較専用
  - `GetIntegrationSnapshotConfigForReporting()` - レポート専用
  - `GetIntegrationSnapshotConfigForIntegrated()` - 統合実行（全フェーズ）
- 4つのバリデーションルール：
  - Completion（完了検証）- Error
  - Sequence（シーケンス検証）- Error
  - StateTransition（状態遷移検証）- Warning
  - Timing（タイミング検証）- Info（デフォルト無効）

**スナップショットテストの実行例:**

```bash
# Integration APIテストディレクトリから実行
cd tests/e2e/integrationapi

# キャプチャモード: 現在のAPIレスポンスをキャプチャ
go test -v -snapshot-mode=capture
# 出力: ../snapshot/integration/story_snapshots/*.json

# 比較モード: 以前のスナップショットと比較
go test -v -snapshot-mode=comparison
# 出力: ../snapshot/integration/story_comparisons/*.json

# レポートモード: HTMLレポートを生成
go test -v -snapshot-mode=reporting
# 出力: ../snapshot/integration/story_reports/*.html

# 統合モード: 全フェーズを一度に実行
go test -v -snapshot-mode=integrated
# 出力: 上記すべてのファイル

# 特定のストーリーのみ実行
go test -v -snapshot-mode=capture -stories="Postman Collection Story - Standard Methods"
```

**スナップショットファイルの構造:**

```
tests/e2e/snapshot/integration/
├── story_snapshots/          # キャプチャされたAPIレスポンス
│   ├── story_snapshot_v1.0.0_postman_collection_story_standard_methods.json
│   ├── story_snapshot_v1.0.0_postman_collection_story_standard_methods_with_body.json
│   ├── story_snapshot_v1.0.0_postman_collection_story_withresponse_methods.json
│   └── story_snapshot_v1.0.0_postman_collection_story_withresponse_methods_with_body.json
├── story_comparisons/        # 比較結果
│   ├── story_comparison_v1.0.0_postman_collection_story_standard_methods.json
│   └── ...
├── story_validations/        # 検証結果
│   ├── story_validation_v1.0.0_postman_collection_story_standard_methods.json
│   └── ...
└── story_reports/            # HTMLレポート
    ├── story_report_v1.0.0_postman_collection_story_standard_methods.html
    └── ...
```

**ファイル名フォーマット:**

スナップショットファイルは`story_snapshot_{tag}_{story_name}.json`形式で生成されます：
- `{tag}`: バージョンタグ（例: `v1.0.0`、`v1.10.19-125-g03232a4`）
- `{story_name}`: ストーリー名（スペースとハイフンはアンダースコアに変換）

**スナップショットテストのユースケース:**

1. **後方互換性の検証**: APIの変更が既存の動作を壊していないか確認
2. **リグレッションテスト**: 過去のバージョンとの比較で予期しない変更を検出
3. **ドキュメント生成**: HTMLレポートでAPIの動作を可視化
4. **CI/CD統合**: 自動テストパイプラインでスナップショット比較を実行

## まとめ

このガイドに従うことで、任意の API モジュールに対して：

1. Postman Collection ベースのストーリー作成
2. E2E テスト実装
3. スナップショットテスト実装

を体系的に行うことができます。

**重要ポイント:**

- Postman 分析結果を活用してストーリー設計
- testlib を活用して再利用可能なテスト基盤を構築
- スナップショットテストで互換性を保証

### Pricing API (実装済み)

```bash
# テストディレクトリ作成
mkdir -p tests/e2e/pricingapi

# テスト実行
cd tests/e2e/pricingapi
go test -v

# 詳細ログ付き実行
LOG_LEVEL=debug go test -v

# 特定のストーリーのみ実行
go test -v -stories="Postman Collection Story - Standard Methods"
go test -v -stories="Postman Collection Story - Standard Methods With Body"
go test -v -stories="Postman Collection Story - WithResponse Methods"
go test -v -stories="Postman Collection Story - WithResponse Methods With Body"
```

**実装済みファイル:**

- `pricing_test.go` - メインテストファイル
- `stories.go` - ストーリー定義とメソッド一覧
- `validation.go` - パラメータ関数とバリデーション
- `story_snapshot_config.go` - スナップショット設定（完全実装済み）
- `testdata/test_params.json` - テストパラメータ（自動生成）

**特徴:**

- 全 Pricing API メソッドをカバー（Pricing Units, Menus, Plans, Metering, Tax Rate, Stripe 連携）
- 4 つのストーリー（Standard, Standard With Body, WithResponse, WithResponse With Body）
- 各ストーリー 26 ステップ、合計 122 メソッドカバレッジ
- testdata/test_params.json からテストデータを読み込み
- 全メソッドバリアントを網羅的にカバー：
  - Struct wrapper メソッド（型付き構造体パラメータ）
  - WithBody メソッド（io.Reader ボディパラメータ）
  - WithResponse バリアント（詳細レスポンス）
  - WithBodyWithResponse バリアント（io.Reader + 詳細レスポンス）

**スナップショット設定:**

- 4 つのモード別設定関数を提供：
  - `GetPricingSnapshotConfigForCapture()` - キャプチャ専用
  - `GetPricingSnapshotConfigForComparison()` - 比較専用
  - `GetPricingSnapshotConfigForReporting()` - レポート専用
  - `GetPricingSnapshotConfigForIntegrated()` - 統合実行（全フェーズ）
- 4 つのバリデーションルール：
  - Completion（完了検証）- Error
  - Sequence（シーケンス検証）- Error
  - StateTransition（状態遷移検証）- Warning
  - Timing（タイミング検証）- Info（デフォルト無効）

#
# LoadTestParams 使用例と実装パターン

このセクションでは、`LoadTestParams()` 関数とポインタヘルパー関数の具体的な使用例を示します。

### 基本的な使用例

#### 1. テストパラメータの読み込み

```go
package authapi

import (
    "testing"
    "github.com/saasus-platform/saasus-sdk-go/tests/e2e/authapi/testdata"
    "github.com/saasus-platform/saasus-sdk-go/tests/testlib"
)

func GetPostmanStoryStandardMethods(t *testing.T) testlib.Story {
    // ステップ1: テストパラメータを読み込む
    params := testdata.LoadTestParams(t)
    
    // ステップ2: variables マップを構築
    variables := map[string]interface{}{
        // 基本設定
        "domain_name":           params.BasicInfo.UpdateParams["domain_name"],
        "from_email_address":    params.BasicInfo.UpdateParams["from_email_address"],
        "reply_email_address":   params.BasicInfo.UpdateParams["reply_email_address"],
        
        // 認証情報
        "callback_url":          params.AuthInfo.UpdateParams["callback_url"],
        
        // ユーザー
        "email":                 params.Users.CreateParams["email"],
        "password":              params.Users.CreateParams["password"],
        
        // ロール
        "role_name":             params.Roles.CreateParams["role_name"],
        "display_name":          params.Roles.CreateParams["display_name"],
        
        // 属性
        "attribute_name":        params.UserAttributes.CreateParams["attribute_name"],
        "attribute_type":        params.UserAttributes.CreateParams["attribute_type"],
        
        // 環境
        "env_id":                params.Envs.CreateParams["id"],
        "name":                  params.Envs.CreateParams["name"],
        
        // テナント
        "attributes":            params.Tenants.CreateParams["attributes"],
        "back_office_staff_email": params.Tenants.CreateParams["back_office_staff_email"],
    }
    
    // ステップ3: ストーリーを返す
    return testlib.Story{
        Name:        "Postman Collection Story - Standard Methods",
        Description: "Postman Collection のテストフローに基づいた Standard メソッドのテスト",
        Variables:   variables,
        Steps: []testlib.Step{
            {
                Name:           "基本設定の更新",
                ClientMethod:   "UpdateBasicInfo",
                Parameters:     updateBasicInfoParams,
                ExpectedStatus: 200,
            },
            {
                Name:           "ユーザーの作成",
                ClientMethod:   "CreateSaasUser",
                Parameters:     createSaasUserParams,
                ExpectedStatus: 201,
                StateUpdate: func(variables map[string]interface{}, response interface{}) error {
                    // レスポンスから user_id を抽出
                    if resp, ok := response.(*authapi.CreateSaasUserResponse); ok {
                        if resp.JSON201 != nil && resp.JSON201.Id != nil {
                            variables["user_id"] = *resp.JSON201.Id
                        }
                    }
                    return nil
                },
            },
            // その他のステップ...
        },
    }
}
```

### パラメータ関数での使用例

#### 2. 単純なパラメータ関数

```go
// getSaasUserParams は GetSaasUser メソッドのパラメータを返します
// メソッドシグネチャ: GetSaasUser(ctx context.Context, userId string, reqEditors ...RequestEditorFn) (*http.Response, error)
func getSaasUserParams(variables map[string]interface{}) interface{} {
    // variables マップから値を取得（LoadTestParams から読み込まれた値）
    userID, _ := variables["user_id"].(string)
    return userID
}
```

#### 3. ポインタヘルパー関数を使用するパラメータ関数

```go
// updateBasicInfoParams は UpdateBasicInfo メソッドのパラメータを返します
func updateBasicInfoParams(variables map[string]interface{}) interface{} {
    domainName, _ := variables["domain_name"].(string)
    fromEmail, _ := variables["from_email_address"].(string)
    replyEmail, _ := variables["reply_email_address"].(string)

    return authapi.UpdateBasicInfoParam{
        DomainName:        domainName,
        FromEmailAddress:  fromEmail,
        ReplyEmailAddress: testdata.StringPtr(replyEmail), // ポインタヘルパー関数を使用
    }
}
```

#### 4. 複数パラメータを持つ関数

```go
// updateSaasUserPasswordParams は UpdateSaasUserPassword メソッドのパラメータを返します
func updateSaasUserPasswordParams(variables map[string]interface{}) interface{} {
    userID, _ := variables["user_id"].(string)
    password, _ := variables["password"].(string)

    // 構造体を返す（フィールド順序はメソッドシグネチャと一致）
    return struct {
        UserId string
        Body   authapi.UpdateSaasUserPasswordParam
    }{
        UserId: userID,
        Body: authapi.UpdateSaasUserPasswordParam{
            Password: password,
        },
    }
}
```

#### 5. WithBody メソッドのパラメータ関数

```go
// createSaasUserWithBodyParams は CreateSaasUserWithBody メソッドのパラメータを返します
func createSaasUserWithBodyParams(variables map[string]interface{}) interface{} {
    email, _ := variables["email"].(string)
    password, _ := variables["password"].(string)

    // リクエストボディを作成
    param := authapi.CreateSaasUserParam{
        Email:    email,
        Password: testdata.StringPtr(password),
    }
    bodyBytes, _ := json.Marshal(param)

    // WithBody メソッドのパラメータ構造
    return struct {
        ContentType string
        Body        io.Reader
    }{
        ContentType: "application/json",
        Body:        bytes.NewReader(bodyBytes),
    }
}
```

#### 6. 型変換を含むパラメータ関数

```go
// createEnvParams は CreateEnv メソッドのパラメータを返します
func createEnvParams(variables map[string]interface{}) interface{} {
    // int 型への変換（JSON からの読み込みでは float64 になる可能性がある）
    envID, ok := variables["env_id"].(int)
    if !ok {
        if envIDFloat, ok := variables["env_id"].(float64); ok {
            envID = int(envIDFloat)
        } else {
            envID = 10 // デフォルト値
        }
    }
    name, _ := variables["name"].(string)
    displayName, _ := variables["display_name"].(string)

    return authapi.Env{
        Id:          authapi.Id(uint64(envID)),
        Name:        name,
        DisplayName: testdata.StringPtr(displayName),
    }
}
```

### 高度な使用例

#### 7. デフォルト値の提供

```go
func createRoleParams(variables map[string]interface{}) interface{} {
    // 型アサーションが失敗した場合のデフォルト値
    roleName, ok := variables["role_name"].(string)
    if !ok || roleName == "" {
        roleName = "default_role"
    }
    
    displayName, ok := variables["display_name"].(string)
    if !ok || displayName == "" {
        displayName = "Default Role"
    }

    return authapi.CreateRoleParam{
        RoleName:    roleName,
        DisplayName: displayName,
    }
}
```

#### 8. 条件付きポインタの設定

```go
func updateBasicInfoParams(variables map[string]interface{}) interface{} {
    domainName, _ := variables["domain_name"].(string)
    fromEmail, _ := variables["from_email_address"].(string)
    
    // オプショナルなフィールド
    var replyEmail *string
    if reply, ok := variables["reply_email_address"].(string); ok && reply != "" {
        replyEmail = testdata.StringPtr(reply)
    }

    return authapi.UpdateBasicInfoParam{
        DomainName:        domainName,
        FromEmailAddress:  fromEmail,
        ReplyEmailAddress: replyEmail, // nil または値
    }
}
```

#### 9. 複雑なネストされた構造

```go
func createTenantParams(variables map[string]interface{}) interface{} {
    name, _ := variables["name"].(string)
    backOfficeEmail, _ := variables["back_office_staff_email"].(string)
    
    // map[string]interface{} 型の attributes
    attributes, ok := variables["attributes"].(map[string]interface{})
    if !ok {
        attributes = make(map[string]interface{})
    }

    return authapi.CreateTenantParam{
        Name:                 name,
        BackOfficeStaffEmail: backOfficeEmail,
        Attributes:           attributes,
    }
}
```

#### 10. StateUpdate 関数でのパラメータ更新

```go
{
    Name:           "ユーザーの作成",
    ClientMethod:   "CreateSaasUser",
    Parameters:     createSaasUserParams,
    ExpectedStatus: 201,
    StateUpdate: func(variables map[string]interface{}, response interface{}) error {
        // レスポンスから user_id を抽出して variables に追加
        if resp, ok := response.(*authapi.CreateSaasUserResponse); ok {
            if resp.JSON201 != nil && resp.JSON201.Id != nil {
                variables["user_id"] = *resp.JSON201.Id
                // 次のステップで getSaasUserParams が user_id を使用できる
            }
        }
        return nil
    },
},
{
    Name:           "ユーザーの取得",
    ClientMethod:   "GetSaasUser",
    Parameters:     getSaasUserParams, // 前のステップで設定された user_id を使用
    ExpectedStatus: 200,
},
```

### ベストプラクティス

#### 1. 型アサーションの安全な使用

```go
// ❌ 悪い例: パニックの可能性
userID := variables["user_id"].(string)

// ✅ 良い例: 安全な型アサーション
userID, ok := variables["user_id"].(string)
if !ok {
    userID = "default_user_id"
}
```

#### 2. ポインタヘルパー関数の活用

```go
// ❌ 悪い例: 冗長
replyEmail := "support@example.com"
param := authapi.UpdateBasicInfoParam{
    ReplyEmailAddress: &replyEmail,
}

// ✅ 良い例: ヘルパー関数を使用
param := authapi.UpdateBasicInfoParam{
    ReplyEmailAddress: testdata.StringPtr("support@example.com"),
}
```

#### 3. デフォルト値の提供

```go
// ✅ 良い例: デフォルト値を提供
envID, ok := variables["env_id"].(int)
if !ok {
    if envIDFloat, ok := variables["env_id"].(float64); ok {
        envID = int(envIDFloat)
    } else {
        envID = 10 // デフォルト値
    }
}
```

#### 4. 構造体フィールドの順序

```go
// ✅ 良い例: メソッドシグネチャと同じ順序
// UpdateSaasUserPassword(ctx, userId, body, reqEditors)
return struct {
    UserId string  // 1番目のパラメータ
    Body   authapi.UpdateSaasUserPasswordParam  // 2番目のパラメータ
}{
    UserId: userID,
    Body:   body,
}
```

#### 5. WithBody メソッドのパラメータ順序

```go
// ✅ 良い例: 追加パラメータ → ContentType → Body の順序
// UpdateSaasUserPasswordWithBody(ctx, userId, contentType, body, reqEditors)
return struct {
    UserId      string    // 追加パラメータ
    ContentType string    // ContentType
    Body        io.Reader // Body
}{
    UserId:      userID,
    ContentType: "application/json",
    Body:        bytes.NewReader(bodyBytes),
}
```

### トラブルシューティング

#### 問題: 型アサーションが失敗する

```go
// エラー: panic: interface conversion: interface {} is nil, not string
userID := variables["user_id"].(string)
```

**解決方法:**

```go
userID, ok := variables["user_id"].(string)
if !ok {
    // エラーハンドリングまたはデフォルト値
    return fmt.Errorf("user_id not found in variables")
}
```

#### 問題: JSON の数値が float64 になる

```go
// エラー: interface conversion: interface {} is float64, not int
envID := variables["env_id"].(int)
```

**解決方法:**

```go
envID, ok := variables["env_id"].(int)
if !ok {
    if envIDFloat, ok := variables["env_id"].(float64); ok {
        envID = int(envIDFloat)
    } else {
        envID = 10 // デフォルト値
    }
}
```

#### 問題: ポインタが nil になる

```go
// エラー: runtime error: invalid memory address or nil pointer dereference
fmt.Println(*replyEmail)
```

**解決方法:**

```go
var replyEmail *string
if reply, ok := variables["reply_email_address"].(string); ok && reply != "" {
    replyEmail = testdata.StringPtr(reply)
}

// 使用前に nil チェック
if replyEmail != nil {
    fmt.Println(*replyEmail)
}
```

### まとめ

- `LoadTestParams(t)` を使用してテストパラメータを読み込む
- 読み込んだパラメータを `variables` マップに格納する
- パラメータ関数は `variables` マップから値を取得する
- ポインタヘルパー関数（`StringPtr`, `IntPtr`, `BoolPtr`）を活用する
- 型アサーションは安全に行う（`value, ok := ...` パターン）
- デフォルト値を提供してロバストなコードを書く
