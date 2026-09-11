# OpenAPI Operation Resolver

OpenAPI 仕様ファイルから operationId を取得し、HTTP メソッドとパスをキーとしてマッピングする CLI ツールです。

## 概要

このツールは、OpenAPI 3.0 形式の YAML 仕様ファイルを解析し、以下の機能を提供します：

- HTTP メソッドとパスから対応する`operationId`を取得
- `operationId`から派生するメソッドバリアント（標準、WithResponse、WithBody 等）の生成
- すべてのエンドポイントと operationId のマッピングをリスト表示
- テキスト形式または JSON 形式での出力

### 主な用途

- `postman-analyzer`が生成するメソッド名を、`client.gen.go`の実際のメソッド名と一致させる
- OpenAPI 仕様から正確なクライアントメソッド名を特定
- API 仕様の確認とドキュメント生成

## インストール

### ビルド

```bash
cd tests/tools/openapi-resolver
go build -o openapi-resolver
```

### 実行可能ファイルの配置（オプション）

```bash
# PATHが通っているディレクトリに配置
sudo mv openapi-resolver /usr/local/bin/
```

## 使用方法

### 基本構文

```bash
openapi-resolver [オプション]
```

### オプション

| オプション         | 説明                                           | 必須              |
| ------------------ | ---------------------------------------------- | ----------------- |
| `-spec <path>`     | OpenAPI 仕様ファイルのパス                     | `-module`と排他   |
| `-module <name>`   | モジュール名（billing, auth 等）               | `-spec`と排他     |
| `-method <method>` | HTTP メソッド（GET, POST, PUT, DELETE, PATCH） | `-list`なしの場合 |
| `-path <path>`     | エンドポイントパス（例: /stripe/info）         | `-list`なしの場合 |
| `-list`            | すべてのエンドポイントをリスト表示             | いいえ            |
| `-json`            | JSON 形式で出力                                | いいえ            |

## 使用例

### 1. 特定のエンドポイントの operationID を取得

```bash
# 仕様ファイルのパスを直接指定
openapi-resolver -spec ../../../api/billingapi.yml -method PUT -path /stripe/info
```

**出力例:**

```
HTTPメソッド: PUT
パス: /stripe/info
OperationID: UpdateStripeInfo
リクエストボディ: true
概要: Update Stripe Connection Info
説明: Updates information on connection with external billing SaaS. Currently possible to connect to Stripe.

メソッドバリアント:
  - UpdateStripeInfo
  - UpdateStripeInfoWithResponse
  - UpdateStripeInfoWithBody
  - UpdateStripeInfoWithBodyWithResponse
```

### 2. モジュール名から自動的に仕様ファイルを読み込む

```bash
# モジュール名を指定（自動的に ../../../api/{module}api.yml を読み込む）
openapi-resolver -module billing -method PUT -path /stripe/info
```

### 3. すべてのエンドポイントをリスト表示

```bash
openapi-resolver -spec ../../../api/billingapi.yml -list
```

**出力例:**

```
合計 15 個のエンドポイント:

[1] GET /stripe/info
    OperationID: GetStripeInfo
    概要: Get Stripe Connection information
    バリアント: GetStripeInfo, GetStripeInfoWithResponse

[2] PUT /stripe/info
    OperationID: UpdateStripeInfo
    概要: Update Stripe Connection Info
    バリアント: UpdateStripeInfo, UpdateStripeInfoWithResponse, UpdateStripeInfoWithBody, UpdateStripeInfoWithBodyWithResponse

[3] DELETE /stripe/info
    OperationID: DeleteStripeInfo
    概要: Delete Stripe Connection
    バリアント: DeleteStripeInfo, DeleteStripeInfoWithResponse

...
```

### 4. JSON 形式で出力

```bash
openapi-resolver -module billing -method PUT -path /stripe/info -json
```

**出力例:**

```json
{
  "http_method": "PUT",
  "path": "/stripe/info",
  "operation_id": "UpdateStripeInfo",
  "has_body": true,
  "summary": "Update Stripe Connection Info",
  "description": "Updates information on connection with external billing SaaS. Currently possible to connect to Stripe.",
  "variants": [
    "UpdateStripeInfo",
    "UpdateStripeInfoWithResponse",
    "UpdateStripeInfoWithBody",
    "UpdateStripeInfoWithBodyWithResponse"
  ]
}
```

### 5. すべてのエンドポイントを JSON 形式でリスト表示

```bash
openapi-resolver -module billing -list -json
```

## メソッドバリアント

このツールは、`oapi-codegen`が生成するメソッドバリアントを自動的に識別します。

### バリアント生成ルール

| HTTP メソッド  | リクエストボディ | 生成されるバリアント                                                                                       |
| -------------- | ---------------- | ---------------------------------------------------------------------------------------------------------- |
| GET/DELETE     | なし             | `{OperationID}`, `{OperationID}WithResponse`                                                               |
| POST/PUT/PATCH | あり             | `{OperationID}`, `{OperationID}WithResponse`, `{OperationID}WithBody`, `{OperationID}WithBodyWithResponse` |
| POST/PUT/PATCH | なし             | `{OperationID}`, `{OperationID}WithResponse`                                                               |

### バリアントの説明

- **標準メソッド** (`{OperationID}`): 基本的なメソッド
- **WithResponse** (`{OperationID}WithResponse`): HTTP レスポンス全体を返すメソッド
- **WithBody** (`{OperationID}WithBody`): リクエストボディを`io.Reader`として受け取るメソッド
- **WithBodyWithResponse** (`{OperationID}WithBodyWithResponse`): 両方の機能を持つメソッド

## postman-analyzer との統合

このツールは`postman-analyzer`から利用され、Postman コレクションの解析時に正確なメソッド名を生成します。

### 統合方法

`postman-analyzer`は以下のように`openapi-resolver`を呼び出します：

```go
// OpenAPI仕様からoperationIdを取得
resolver, err := NewOperationResolver(specPath)
if err != nil {
    // フォールバック: 従来の方法でメソッド名を生成
    methodName = generateMethodName(method, endpoint)
} else {
    mapping, err := resolver.Resolve(method, endpoint)
    if err != nil || mapping == nil {
        // フォールバック: 従来の方法でメソッド名を生成
        methodName = generateMethodName(method, endpoint)
    } else {
        // OpenAPI仕様から取得したoperationIDを使用
        methodName = mapping.OperationID
        variants = mapping.Variants
    }
}
```

### postman-analyzer での使用例

```bash
# プロジェクトルートから実行する場合
go run tests/tools/postman-analyzer/*.go \
  -collection "tests/tools/postman-analyzer/collections/SaaSus Billing API Schema.postman_collection.json" \
  -spec api/billingapi.yml \
  -module billing

# postman-analyzerディレクトリから実行する場合
cd tests/tools/postman-analyzer
go run . -collection "collections/SaaSus Billing API Schema.postman_collection.json" \
         -spec ../../../api/billingapi.yml \
         -module billing
```

**出力例:**

```
🔍 SaaSus Billing API Schema - Postman Collection Analysis
============================================================
Collection: SaaSus Billing API Schema
Total Test Steps: 5
Total Client Methods: 8

=== Generated Methods ===
✅ GetStripeInfo (standard)
✅ GetStripeInfoWithResponse (with_response)
✅ UpdateStripeInfo (standard)
✅ UpdateStripeInfoWithResponse (with_response)
✅ UpdateStripeInfoWithBody (with_body)
✅ UpdateStripeInfoWithBodyWithResponse (with_body_with_response)
✅ DeleteStripeInfo (standard)
✅ DeleteStripeInfoWithResponse (with_response)

=== File Export ===
📄 JSON analysis exported to tests/tools/postman-analyzer/reports/billing/postman_analysis.json
📄 Text report exported to tests/tools/postman-analyzer/reports/billing/postman_analysis_report.txt
```

**注意:** 出力ファイルは`tests/tools/postman-analyzer/reports/{module}/`ディレクトリに保存されます。

## パス正規化

このツールは、異なるパスパラメータ形式を自動的に正規化してマッチングします。

### サポートされる形式

| 入力形式                     | 正規化後                   |
| ---------------------------- | -------------------------- |
| `/feedbacks/{feedback_id}`   | `/feedbacks/{feedback_id}` |
| `/feedbacks/:feedback_id`    | `/feedbacks/{feedback_id}` |
| `/feedbacks/{{feedback_id}}` | `/feedbacks/{feedback_id}` |

これにより、Postman コレクションなどで異なる形式のパスが使用されていても、正確にマッチングできます。

## エラーハンドリング

### 一般的なエラーと対処法

#### 1. ファイルが見つからない

```
エラー: failed to initialize resolver: failed to read OpenAPI spec file: open ../../../api/billingapi.yml: no such file or directory
```

**対処法:**

- ファイルパスが正しいか確認してください
- 相対パスの場合、実行ディレクトリからの相対パスを指定してください
- `-module`オプションを使用する場合、`tests/tools/openapi-resolver/`ディレクトリから実行してください

#### 2. YAML 解析エラー

```
エラー: failed to initialize resolver: failed to parse YAML: yaml: line 10: mapping values are not allowed in this context
```

**対処法:**

- OpenAPI 仕様ファイルの YAML 構文が正しいか確認してください
- YAML リンターを使用して構文エラーをチェックしてください

#### 3. operation が見つからない

```
エラー: failed to resolve operation: Resolve: /invalid/path: operation not found for given method and path
```

**対処法:**

- HTTP メソッドとパスが正しいか確認してください
- `-list`オプションを使用して、利用可能なエンドポイントを確認してください
- パスパラメータの形式が正しいか確認してください（`{param}`形式を推奨）

#### 4. HTTP メソッドが不正

```
エラー: invalid HTTP method: get (must be GET, POST, PUT, DELETE, or PATCH)
```

**対処法:**

- HTTP メソッドは大文字で指定してください（GET, POST, PUT, DELETE, PATCH）
- ツールは自動的に大文字に変換しますが、有効なメソッド名である必要があります

#### 5. 必須パラメータが不足

```
エラー: -method is required when -list is not specified
```

**対処法:**

- `-list`オプションを使用しない場合、`-method`と`-path`の両方を指定してください
- または、`-list`オプションを使用してすべてのエンドポイントを表示してください

## トラブルシューティング

### デバッグ方法

#### 1. 利用可能なエンドポイントを確認

```bash
# すべてのエンドポイントをリスト表示
openapi-resolver -module billing -list
```

#### 2. JSON 形式で詳細情報を取得

```bash
# JSON形式で出力して、詳細情報を確認
openapi-resolver -module billing -method PUT -path /stripe/info -json | jq .
```

#### 3. OpenAPI 仕様ファイルを直接確認

```bash
# 仕様ファイルの内容を確認
cat ../../../api/billingapi.yml | grep -A 10 "/stripe/info"
```

### よくある問題

#### Q1: operationId が空で返される

**原因:** OpenAPI 仕様ファイルで`operationId`が定義されていない

**対処法:**

- OpenAPI 仕様ファイルを確認し、該当するエンドポイントに`operationId`を追加してください
- 例:

```yaml
paths:
  /stripe/info:
    put:
      operationId: UpdateStripeInfo # これを追加
      summary: Update Stripe Connection Info
```

#### Q2: パスパラメータを含むパスがマッチしない

**原因:** パスパラメータの形式が異なる

**対処法:**

- パスパラメータは`{param}`形式を使用してください
- ツールは自動的に`:param`や`{{param}}`を`{param}`に正規化します
- 例: `/feedbacks/:feedback_id` → `/feedbacks/{feedback_id}`

#### Q3: モジュール名から仕様ファイルが見つからない

**原因:** 実行ディレクトリが正しくない、またはファイル構造が異なる

**対処法:**

- `tests/tools/openapi-resolver/`ディレクトリから実行してください
- または、`-spec`オプションで絶対パスを指定してください
- ファイル構造が標準と異なる場合は、`-spec`オプションを使用してください

#### Q4: バリアントが期待通りに生成されない

**原因:** リクエストボディの定義が不正確

**対処法:**

- OpenAPI 仕様ファイルで`requestBody`が正しく定義されているか確認してください
- `required: true`が設定されているか確認してください
- 例:

```yaml
requestBody:
  required: true
  content:
    application/json:
      schema:
        type: object
```

## 開発者向け情報

### プロジェクト構造

```
tests/tools/openapi-resolver/
├── main.go                    # CLIエントリーポイント
├── cli.go                     # CLIハンドラー
├── parser.go                  # OpenAPI仕様パーサー
├── resolver.go                # Operation resolver
├── client_parser.go           # client.gen.goパーサー
├── type_parser.go             # types.gen.goパーサー（NEW）
├── variant_generator.go       # メソッドバリアント生成
├── path_normalizer.go         # パス正規化
├── types.go                   # データ型定義
├── errors.go                  # エラー定義
├── *_test.go                  # テストファイル
└── testdata/                  # テストデータ
    ├── test_spec.yml          # 基本テスト用
    ├── edge_cases_spec.yml    # エッジケーステスト用
    ├── invalid_spec.yml       # エラーハンドリングテスト用
    └── malformed.yml          # YAML解析エラーテスト用
```

### テストの実行

```bash
# すべてのテストを実行
go test ./...

# 詳細な出力でテストを実行
go test -v ./...

# カバレッジレポートを生成
go test -cover ./...

# 特定のテストを実行
go test -run TestPathNormalizer
```

### コードの拡張

#### 新しいバリアントルールの追加

`variant_generator.go`の`GenerateVariants`メソッドを修正してください：

```go
func (g *MethodVariantGenerator) GenerateVariants(operationID, method string, hasBody bool) []string {
    // 新しいバリアントロジックを追加
}
```

#### 新しいパス正規化ルールの追加

`path_normalizer.go`の`Normalize`メソッドを修正してください：

```go
func (n *PathNormalizer) Normalize(path string) string {
    // 新しい正規化ルールを追加
}
```

#### client.gen.goからメソッド情報を抽出

`client_parser.go`は、生成されたクライアントコードからメソッドシグネチャを抽出します：

```go
// ClientInterfaceParserを作成
parser, err := NewClientInterfaceParser("communication")
if err != nil {
    log.Fatal(err)
}

// すべてのメソッドを解析
methods, err := parser.ParseAllMethods()
if err != nil {
    log.Fatal(err)
}

// operationIDからメソッドを検索
method, err := parser.FindMethodByOperationID("CreateFeedback")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("メソッド名: %s\n", method.MethodName)
fmt.Printf("パラメータ: %v\n", method.ParameterTypes)
```

**主な機能:**

- `ClientWithResponses`インターフェースからメソッド定義を抽出
- メソッドシグネチャ（パラメータ型、戻り値型）を取得
- operationIDからメソッド名を逆引き
- WithResponse、WithBodyなどのサフィックスを考慮した検索

#### types.gen.goからパラメータ型情報を抽出

`type_parser.go`は、生成された型定義ファイルから構造体の詳細情報を抽出します：

```go
// TypeDefinitionParserを作成
parser, err := NewTypeDefinitionParser("communication")
if err != nil {
    log.Fatal(err)
}

// 型定義を取得
typeDef, err := parser.GetTypeDefinition("CreateFeedbackJSONRequestBody")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("型名: %s\n", typeDef.Name)
for _, field := range typeDef.Fields {
    fmt.Printf("  フィールド: %s (型: %s, JSONタグ: %s)\n", 
        field.Name, field.Type, field.JSONTag)
}
```

**主な機能:**

- Go ASTパーサーを使用した構造体定義の解析
- フィールド名、型、JSONタグ、コメントの抽出
- ネストした型の再帰的解析
- 型定義のキャッシング機構

## 完全なデータフロー

OpenAPI Operation Resolverは、以下のツールチェーンの中核を担います：

```
┌─────────────────────────────────────────────────────────────────┐
│                        Complete Data Flow                        │
└─────────────────────────────────────────────────────────────────┘

1. Postman Collection
   ↓ (URL + Method)
   
2. OpenAPI Spec (api/{module}api.yml)
   ↓ (operationId) ← OpenAPI Operation Resolver
   
3. client.gen.go (generated/{module}api/client.gen.go)
   ↓ (Method Name + Parameter Types) ← ClientInterfaceParser
   
4. types.gen.go (generated/{module}api/types.gen.go)
   ↓ (Type Definition + Fields) ← TypeDefinitionParser
   
5. postman-analyzer Report (JSON)
   ↓ (Complete Type Information)
   
6. test-param-generator
   ↓ (Accurate Type Mapping)
   
7. test_params.json
   ✓ (Correct Types and Values)
```

### ツール間の連携

#### 1. postman-analyzer との統合

postman-analyzerは、OpenAPI Operation Resolverを使用して正確なoperationIdを取得します：

```go
// OpenAPI仕様からoperationIdを取得
resolver, err := NewOperationResolver(specPath)
if err == nil {
    mapping, err := resolver.Resolve(method, endpoint)
    if err == nil && mapping != nil {
        step.OperationID = mapping.OperationID
        
        // client.gen.goからメソッド名を取得
        clientParser, _ := NewClientInterfaceParser(moduleName)
        methodSig, _ := clientParser.FindMethodByOperationID(step.OperationID)
        if methodSig != nil {
            step.ClientMethodName = methodSig.MethodName
            
            // パラメータ型名を抽出
            for _, param := range methodSig.ParameterTypes {
                if strings.HasSuffix(param.Type, "JSONRequestBody") {
                    step.ParameterTypeName = param.Type
                    
                    // types.gen.goから型定義を取得
                    typeParser, _ := NewTypeDefinitionParser(moduleName)
                    typeDef, _ := typeParser.GetTypeDefinition(param.Type)
                    if typeDef != nil {
                        // フィールド情報をレポートに追加
                        step.ParameterTypeFields = convertToFieldInfo(typeDef.Fields)
                    }
                }
            }
        }
    }
}
```

#### 2. test-param-generator での活用

test-param-generatorは、postman-analyzerが生成した型情報を使用して正確なパラメータマッピングを実現します：

```go
// postman-analyzerレポートから型情報を読み込む
for _, step := range report.TestFlow {
    if step.ParameterTypeName != "" && len(step.ParameterTypeFields) > 0 {
        // 型情報を直接使用してマッピング
        params := make(map[string]interface{})
        
        for jsonTag, fieldInfo := range step.ParameterTypeFields {
            // iterationDataから対応する値を検索
            if value, ok := iterationData[jsonTag]; ok {
                params[fieldInfo.Name] = value
            } else {
                params[fieldInfo.Name] = generateDefaultValue(fieldInfo.Type)
            }
        }
        
        testParams[step.ClientMethodName] = params
    }
}
```

### 利点

この統合により、以下の利点が得られます：

1. **型安全性の向上**: 実際の型定義に基づいたパラメータ生成
2. **マッピング精度の向上**: 80%以上のマッピング成功率を達成
3. **メンテナンス性の向上**: OpenAPI仕様の変更が自動的に反映される
4. **エラーの早期発見**: 型の不一致をコンパイル時に検出

## ライセンス

このプロジェクトのライセンスについては、プロジェクトルートの`LICENSE.txt`を参照してください。

## 高度なトラブルシューティング

### client.gen.goパース失敗時の対処法

#### 症状

```bash
Error: failed to parse client.gen.go: file not found or invalid syntax
```

#### 原因と対処法

**原因1: ファイルが存在しない**

```bash
# ファイルの存在を確認
ls -la generated/{module}api/client.gen.go
```

対処法：
```bash
# コード生成を実行
./generate_all.sh

# または特定のモジュールのみ
./generate.sh {module}api api/{module}api.yml
```

**原因2: Go構文エラー**

```bash
# 構文エラーをチェック
go build generated/{module}api/client.gen.go
```

対処法：
- OpenAPI仕様ファイルを確認し、不正な定義を修正
- oapi-codegenのバージョンを確認（v1.12.4以上を推奨）
- 再生成を実行

**原因3: ClientWithResponsesインターフェースが見つからない**

```bash
# インターフェースの存在を確認
grep "type ClientWithResponses interface" generated/{module}api/client.gen.go
```

対処法：
- oapi-codegenの設定を確認
- `generate: client: true`が設定されているか確認
- OpenAPI仕様に`operationId`が定義されているか確認

#### デバッグ方法

1. **詳細ログを有効化**

```go
parser, err := NewClientInterfaceParser("communication")
if err != nil {
    log.Printf("Parser initialization failed: %v", err)
    log.Printf("Tried path: generated/communicationapi/client.gen.go")
}
```

2. **手動でファイルを確認**

```bash
# ClientWithResponsesインターフェースを検索
grep -A 50 "type ClientWithResponses interface" generated/communicationapi/client.gen.go

# メソッド定義を確認
grep "CreateFeedback" generated/communicationapi/client.gen.go
```

3. **AST解析のテスト**

```bash
# パーサーのテストを実行
cd tests/tools/openapi-resolver
go test -v -run TestClientInterfaceParser
```

### types.gen.goパース失敗時の対処法

#### 症状

```bash
Error: failed to parse types.gen.go: type definition not found
```

#### 原因と対処法

**原因1: ファイルが存在しない**

```bash
# ファイルの存在を確認
ls -la generated/{module}api/types.gen.go
```

対処法：
```bash
# コード生成を実行
./generate_all.sh
```

**原因2: 型定義が見つからない**

```bash
# 型定義を検索
grep "type CreateFeedbackJSONRequestBody struct" generated/communicationapi/types.gen.go
```

対処法：
- OpenAPI仕様で`requestBody`が正しく定義されているか確認
- `operationId`が正しく設定されているか確認
- 型名の命名規則を確認（通常は`{OperationID}JSONRequestBody`）

**原因3: ネスト型の解析失敗**

```bash
# ネストした型を確認
grep -A 20 "type CreateFeedbackJSONRequestBody struct" generated/communicationapi/types.gen.go
```

対処法：
- 循環参照がないか確認
- ネストの深さが適切か確認（推奨: 5階層以下）
- 必要に応じてOpenAPI仕様を簡素化

#### デバッグ方法

1. **型定義の一覧を取得**

```go
parser, err := NewTypeDefinitionParser("communication")
if err != nil {
    log.Fatal(err)
}

allTypes, err := parser.ParseAllTypes()
if err != nil {
    log.Fatal(err)
}

for typeName := range allTypes {
    fmt.Println(typeName)
}
```

2. **特定の型を詳細に確認**

```bash
# 型定義を抽出
awk '/type CreateFeedbackJSONRequestBody struct/,/^}/' \
    generated/communicationapi/types.gen.go
```

3. **JSONタグの確認**

```bash
# JSONタグを抽出
grep -A 10 "type CreateFeedbackJSONRequestBody struct" \
    generated/communicationapi/types.gen.go | \
    grep "json:"
```

### 型情報が不足している場合の対処法

#### 症状

postman-analyzerのレポートに型情報が含まれていない：

```json
{
  "operation_id": "CreateFeedback",
  "client_method_name": "",
  "parameter_type_name": "",
  "parameter_type_fields": {}
}
```

#### 原因と対処法

**原因1: OpenAPI仕様にoperationIdが定義されていない**

```bash
# operationIdを確認
grep -A 5 "post:" api/communicationapi.yml | grep "operationId"
```

対処法：
```yaml
# api/communicationapi.ymlに追加
paths:
  /feedbacks:
    post:
      operationId: CreateFeedback  # これを追加
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
```

**原因2: requestBodyが定義されていない**

```bash
# requestBodyを確認
grep -A 10 "operationId: CreateFeedback" api/communicationapi.yml | grep "requestBody"
```

対処法：
```yaml
# requestBodyを追加
operationId: CreateFeedback
requestBody:
  required: true
  content:
    application/json:
      schema:
        type: object
        properties:
          user_id:
            type: string
          feedback_title:
            type: string
```

**原因3: postman-analyzerが古いバージョン**

対処法：
```bash
# 最新版を使用
cd tests/tools/postman-analyzer
go build -o postman-analyzer .

# 型情報を含むレポートを生成
./postman-analyzer \
  -collection collections/SaaSus\ Communication\ API\ Schema.postman_collection.json \
  -spec ../../../api/communicationapi.yml \
  -module communication
```

#### デバッグ方法

1. **OpenAPI Resolverを直接テスト**

```bash
cd tests/tools/openapi-resolver
go run . -module communication -method POST -path /feedbacks
```

期待される出力：
```
HTTPメソッド: POST
パス: /feedbacks
OperationID: CreateFeedback
リクエストボディ: true
```

2. **ClientInterfaceParserを直接テスト**

```go
parser, _ := NewClientInterfaceParser("communication")
method, err := parser.FindMethodByOperationID("CreateFeedback")
if err != nil {
    log.Printf("Method not found: %v", err)
} else {
    log.Printf("Found method: %s", method.MethodName)
    log.Printf("Parameters: %v", method.ParameterTypes)
}
```

3. **TypeDefinitionParserを直接テスト**

```go
parser, _ := NewTypeDefinitionParser("communication")
typeDef, err := parser.GetTypeDefinition("CreateFeedbackJSONRequestBody")
if err != nil {
    log.Printf("Type not found: %v", err)
} else {
    log.Printf("Found type: %s", typeDef.Name)
    for _, field := range typeDef.Fields {
        log.Printf("  Field: %s (%s)", field.Name, field.JSONTag)
    }
}
```

### 統合テストの実行

完全なフローをテストする方法：

```bash
# 1. OpenAPI仕様を更新
vim api/communicationapi.yml

# 2. コード生成
./generate_all.sh

# 3. postman-analyzerを実行
cd tests/tools/postman-analyzer
go run . -module communication -spec ../../../api/communicationapi.yml

# 4. レポートを確認
cat reports/communication/postman_analysis.json | jq '.test_flow[0]'

# 5. test-param-generatorを実行
cd ../test-param-generator
go run . -report ../postman-analyzer/reports/communication/postman_analysis.json \
         -module communication \
         -verbose

# 6. 生成されたファイルを確認
cat tests/e2e/communicationapi/test_params.json | jq .
```

### パフォーマンス最適化

#### 大きなファイルの処理

types.gen.goが非常に大きい場合（10MB以上）：

```go
// キャッシュサイズを増やす
parser, err := NewTypeDefinitionParser("communication")
parser.SetCacheSize(2000) // デフォルト: 1000
```

#### 並列処理

複数のモジュールを並列処理する場合：

```bash
# 並列実行スクリプト
for module in auth billing pricing communication; do
  (
    cd tests/tools/postman-analyzer
    go run . -module ${module} -spec ../../../api/${module}api.yml
  ) &
done
wait
```

## サポート

問題が発生した場合や機能リクエストがある場合は、プロジェクトの Issue トラッカーに報告してください。

報告時には以下の情報を含めてください：

1. 実行したコマンド
2. エラーメッセージの全文
3. 使用しているGoのバージョン（`go version`）
4. 対象のモジュール名
5. OpenAPI仕様ファイルの該当部分（可能であれば）
6. 生成されたclient.gen.goとtypes.gen.goのファイルサイズ
