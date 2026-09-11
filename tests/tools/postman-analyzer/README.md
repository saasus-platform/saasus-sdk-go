# Postman Collection Analyzer

> **最終更新日**: 2025年10月29日  
> **バージョン**: 2.0

Postmanコレクションを分析し、テストフローとクライアントメソッド定義を生成するツールです。OpenAPI仕様と連携して正確なメソッド名を生成します。

## Usage

```bash
# プロジェクトルートから実行（推奨）
go run tests/tools/postman-analyzer/*.go -module billing -spec api/billingapi.yml

# postman-analyzerディレクトリから実行する場合
cd tests/tools/postman-analyzer
go run . -module communication -spec ../../../api/communicationapi.yml

# コレクションを明示的に指定
go run . -collection "collections/SaaSus Billing API Schema.postman_collection.json" \
         -spec ../../../api/billingapi.yml \
         -module billing

# カスタム出力ディレクトリ
go run . -module billing -spec ../../../api/billingapi.yml -output reports/custom
```

### Options

- `-module <name>` - モジュール名（billing, auth, communication 等）
- `-spec <path>` - OpenAPI 仕様ファイルのパス（オプション、モジュール名から自動推定）
- `-collection <path>` - Postman コレクション JSON ファイルのパス
- `-output <dir>` - 出力ディレクトリ（デフォルト: reports/{module}）

### Examples

```bash
# Billing APIを分析（OpenAPI仕様と連携）
go run tests/tools/postman-analyzer/*.go -module billing -spec api/billingapi.yml

# Communication APIを分析
go run tests/tools/postman-analyzer/*.go -module communication -spec api/communicationapi.yml

# カスタム出力ディレクトリを指定
go run tests/tools/postman-analyzer/*.go -module billing -spec api/billingapi.yml -output reports/custom

# 任意のコレクションファイルを分析
go run tests/tools/postman-analyzer/*.go \
  -collection /path/to/collection.json \
  -spec /path/to/openapi.yml \
  -output reports/custom
```

### Building (Optional)

If you want to build a binary for repeated use:

```bash
go build -o postman-analyzer
./postman-analyzer -module communication
```

## Output Files

分析結果は 2 つのファイルとして出力されます:

1. `reports/{module}/postman_analysis.json` - 構造化された JSON データ
2. `reports/{module}/postman_analysis_report.txt` - 人間が読めるテキストレポート

### テキストレポート出力例

```
TEST FLOW
----------------------------------------
1. 基本設定情報の取得(Get Basic Configurations)
   Method: GET /basic-info
   Module: auth
   Operation ID: GetBasicInfo
   Client Method: GetBasicInfo
   Expected Status: 200
   Tests: 18 assertions
   Sets Variables: domain_name, from_email_address

2. 基本設定情報の更新(Update Basic Configurations)
   Method: PUT /basic-info
   Module: auth
   Operation ID: UpdateBasicInfo
   Client Method: UpdateBasicInfo
   Parameter Type: UpdateBasicInfoJSONRequestBody
   Type Fields: 3 fields
     - domain_name: string (required)
     - from_email_address: string (required)
     - reply_email_address: *string
   Expected Status: 200
   Tests: 2 assertions
   Required Variables: domain_name, from_email_address
   Body Fields: domain_name, from_email_address

GENERATED METHODS
----------------------------------------
✅ GetBasicInfo (standard)
✅ GetBasicInfoWithResponse (with_response)
✅ UpdateBasicInfo (standard)
   Params: UpdateBasicInfoJSONRequestBody
✅ UpdateBasicInfoWithResponse (with_response)
   Params: UpdateBasicInfoJSONRequestBody
✅ UpdateBasicInfoWithBody (with_body)
   Params: UpdateBasicInfoJSONRequestBody
✅ UpdateBasicInfoWithBodyWithResponse (with_body_with_response)
   Params: UpdateBasicInfoJSONRequestBody
```

**注**: 
- `Operation ID`: OpenAPI仕様から取得したoperationId
- `Client Method`: `client.gen.go`で生成されたメソッド名
- `Parameter Type`: リクエストボディまたはパラメータの型名（`types.gen.go`から抽出）
- `Type Fields`: パラメータ型のフィールド情報（型、必須/任意の区別を含む）

### JSON 出力例

```json
{
  "collection_name": "SaaSus Auth API Schema",
  "test_flow": [
    {
      "name": "基本設定情報の取得(Get Basic Configurations)",
      "method": "GET",
      "endpoint": "/basic-info",
      "module": "auth",
      "description": "",
      "tests": ["pm.test(\"リクエストはステータス コードが 200 で成功\", ...)"],
      "expected_status": 200,
      "variables": {
        "domain_name": "extracted_from_response",
        "from_email_address": "extracted_from_response"
      },
      "required_variables": [],
      "request_body": {},
      "body_fields": [],
      "operation_id": "GetBasicInfo",
      "client_method_name": "GetBasicInfo"
    },
    {
      "name": "基本設定情報の更新(Update Basic Configurations)",
      "method": "PUT",
      "endpoint": "/basic-info",
      "module": "auth",
      "description": "",
      "tests": ["pm.test(\"リクエストはステータス コードが 200 で成功\", ...)"],
      "expected_status": 200,
      "variables": {},
      "required_variables": ["domain_name", "from_email_address"],
      "request_body": {
        "domain_name": "{{domain_name}}",
        "from_email_address": "{{from_email_address}}"
      },
      "body_fields": ["domain_name", "from_email_address"],
      "operation_id": "UpdateBasicInfo",
      "client_method_name": "UpdateBasicInfo",
      "parameter_type_name": "UpdateBasicInfoJSONRequestBody",
      "parameter_type_fields": {
        "domain_name": {
          "name": "DomainName",
          "type": "string",
          "json_tag": "domain_name",
          "required": true
        },
        "from_email_address": {
          "name": "FromEmailAddress",
          "type": "string",
          "json_tag": "from_email_address",
          "required": true
        }
      }
    }
  ],
  "client_methods": [
    {
      "name": "GetBasicInfo",
      "type": "standard",
      "http_method": "GET",
      "endpoint": "/basic-info",
      "module": "auth",
      "has_body": false,
      "operation_id": "GetBasicInfo",
      "param_types": []
    },
    {
      "name": "UpdateBasicInfo",
      "type": "standard",
      "http_method": "PUT",
      "endpoint": "/basic-info",
      "module": "auth",
      "has_body": true,
      "operation_id": "UpdateBasicInfo",
      "param_types": ["UpdateBasicInfoJSONRequestBody"]
    }
  ]
}
```

**JSON出力の主要フィールド:**

- `operation_id`: OpenAPI仕様から取得したoperationId
- `client_method_name`: `client.gen.go`で生成されたメソッド名
- `parameter_type_name`: リクエストボディまたはパラメータの型名
- `parameter_type_fields`: パラメータ型の詳細なフィールド情報（型、JSONタグ、必須/任意）
- `param_types`: クライアントメソッドのパラメータ型リスト

## Directory Structure

```
postman-analyzer/
├── analyzer.go                    # Main analyzer tool
├── collections/                   # Collection files directory
│   ├── SaaSus Analysis API Schema.postman_collection.json
│   └── SaaSus Billing API Schema.postman_collection.json
├── output/                        # Default output directory
│   ├── postman_analysis.json     # Generated analysis (JSON)
│   └── postman_analysis_report.txt # Generated report (text)
├── reports/                       # Custom output directories
│   └── analysis/                  # API-specific reports
└── README.md                      # This file
```

## Features

- Postman コレクションからテストフローを抽出
- OpenAPI 仕様と連携して正確なメソッド名を生成
- operationID からメソッドバリアントを自動生成
  - 標準メソッド: `GetStripeInfo`
  - WithResponse: `GetStripeInfoWithResponse`
  - WithBody: `UpdateStripeInfoWithBody`
  - WithBodyWithResponse: `UpdateStripeInfoWithBodyWithResponse`
- モジュール名を自動的に各テストステップとクライアントメソッドに付与
- パラメータ型情報を`types.gen.go`から抽出（`param_types`フィールド）
- ネストされたフォルダ構造をサポート
- JSON とテキスト形式でエクスポート
- 任意の Postman コレクション形式に対応

## モジュール識別

Postman コレクションの URL 変数から自動的にモジュール名を識別します。

### URL 変数とモジュール名のマッピング

| URL 変数                | モジュール名   | 説明                   |
| ----------------------- | -------------- | ---------------------- |
| `{{UrlAuth}}`           | auth           | 認証・認可 API         |
| `{{UrlPricing}}`        | pricing        | 価格設定 API           |
| `{{UrlBilling}}`        | billing        | 請求 API               |
| `{{UrlCommunication}}`  | communication  | コミュニケーション API |
| `{{UrlIntegration}}`    | integration    | 統合 API               |
| `{{UrlAwsMarketplace}}` | awsmarketplace | AWS MarketplaceAPI     |
| `{{UrlApiLog}}`         | apilog         | API ログ API           |
| `{{UrlApiGateway}}`     | api-gateway    | API GatewayAPI         |

URL 変数が認識できない場合、モジュール名は"unknown"として表示されます。

### OpenAPI 仕様ファイルの自動解決

モジュール名から OpenAPI 仕様ファイルのパスを自動的に推定します。以下のマッピングルールが適用されます：

| モジュール名     | OpenAPI 仕様ファイル    | 備考                   |
| ---------------- | ----------------------- | ---------------------- |
| `apigateway`     | `apigatewayapi.yml`     | API GatewayAPI         |
| `apilogs`        | `apilogsapi.yml`        | API ログ API           |
| `auth`           | `authapi.yml`           | 認証・認可 API         |
| `awsmarketplace` | `awsmarketplaceapi.yml` | AWS MarketplaceAPI     |
| `billing`        | `billingapi.yml`        | 請求 API               |
| `communication`  | `communicationapi.yml`  | コミュニケーション API |
| `integration`    | `integration.yml`       | 統合 API（特殊ケース） |
| `pricing`        | `pricingapi.yml`        | 価格設定 API           |

**標準命名規則**: `{module}api.yml`（例: `communicationapi.yml`）

**特殊ケース**: `integration`モジュールのみ標準命名規則と異なるファイル名（`integration.yml`）を使用します（`moduleSpecOverrides`で定義）。

## OpenAPI 連携

`openapi-resolver-lib`を使用して OpenAPI 仕様から operationID を取得します。

### パラメータ型情報の抽出

各クライアントメソッドに対して、`generated/{module}api/types.gen.go`からパラメータ型情報を抽出し、`param_types`フィールドに格納します。これにより、test-param-generator などの下流ツールが型情報を活用できます。

**出力形式:**

- JSON レポート: `param_types`配列に型名を格納
- テキストレポート: メソッド名の後に`[Params: 型名1, 型名2, ...]`形式で表示
- パラメータがない場合は何も表示されません

### フォールバック動作

OpenAPI 仕様が指定されていない、または operationID が見つからない場合、従来の方法（エンドポイントからメソッド名を生成）にフォールバックします。

## Adding New Collections

1. Postman コレクション JSON ファイルを`collections/`ディレクトリに配置
2. 対応する OpenAPI 仕様ファイルを`api/`ディレクトリに配置
3. アナライザーを実行

```bash
go run tests/tools/postman-analyzer/*.go \
  -module newapi \
  -spec api/newapiapi.yml
```

## 実践的な使用例

### 複数モジュールの分析

プロジェクト全体のAPIを一括で分析する場合：

```bash
#!/bin/bash
# 全モジュールを分析するスクリプト例

modules=("auth" "billing" "pricing" "communication" "integration" "awsmarketplace" "apilog")

for module in "${modules[@]}"; do
  echo "Analyzing $module module..."
  go run tests/tools/postman-analyzer/*.go \
    -module "$module" \
    -spec "api/${module}api.yml" \
    -output "tests/tools/postman-analyzer/reports/$module"
done

echo "All modules analyzed successfully!"
```

### CI/CDパイプラインでの使用

GitHub Actionsでの使用例：

```yaml
- name: Analyze Postman Collections
  run: |
    go run tests/tools/postman-analyzer/*.go \
      -module billing \
      -spec api/billingapi.yml \
      -output reports/billing
    
    # レポートをアーティファクトとして保存
    - uses: actions/upload-artifact@v3
      with:
        name: postman-analysis
        path: reports/
```

### カスタムコレクションの分析

独自のPostmanコレクションを分析する場合：

```bash
# カスタムコレクションを分析
go run tests/tools/postman-analyzer/*.go \
  -collection "/path/to/custom-collection.json" \
  -spec "/path/to/openapi-spec.yml" \
  -output "reports/custom"

# OpenAPI仕様なしで分析（フォールバックモード）
go run tests/tools/postman-analyzer/*.go \
  -collection "/path/to/custom-collection.json" \
  -output "reports/custom"
```

### test-param-generatorとの連携

分析結果を使用してE2Eテストパラメータを生成：

```bash
# 1. Postman Collectionを分析
go run tests/tools/postman-analyzer/*.go \
  -module billing \
  -spec api/billingapi.yml

# 2. 分析結果からテストパラメータを生成
go run tests/tools/test-param-generator/*.go \
  -report tests/tools/postman-analyzer/reports/billing/postman_analysis.json \
  -module billing \
  -output tests/e2e/billing/test_params.json
```

## トラブルシューティング

### よくあるエラーと解決方法

#### エラー: "module not found"

**症状:**
```
Error resolving module 'billing': module "billing" not found; available modules: auth, pricing, communication
```

**原因:**
- 指定したモジュール名に対応するPostmanコレクションファイルが見つからない
- コレクションファイルの命名規則が正しくない

**解決方法:**
1. `collections/`ディレクトリにコレクションファイルが存在するか確認
2. ファイル名が`SaaSus {ModuleName} API Schema.postman_collection.json`形式であることを確認
3. 利用可能なモジュール一覧を確認して正しいモジュール名を使用

```bash
# 利用可能なコレクションを確認
ls tests/tools/postman-analyzer/collections/

# 正しいモジュール名で再実行
go run tests/tools/postman-analyzer/*.go -module auth
```

#### エラー: "failed to initialize resolver"

**症状:**
```
Warning: Failed to initialize OpenAPI resolver for module 'billing': failed to read OpenAPI spec file
Falling back to best-effort analysis; operationId resolution may be partial
```

**原因:**
- OpenAPI仕様ファイルが見つからない
- ファイルパスが間違っている
- 実行ディレクトリが正しくない

**解決方法:**
1. OpenAPI仕様ファイルが`api/`ディレクトリに存在するか確認
2. プロジェクトルートから実行する（推奨）
3. `-spec`オプションで明示的にパスを指定

```bash
# プロジェクトルートから実行（推奨）
cd /path/to/project-root
go run tests/tools/postman-analyzer/*.go -module billing -spec api/billingapi.yml

# または絶対パスを指定
go run tests/tools/postman-analyzer/*.go \
  -module billing \
  -spec /absolute/path/to/api/billingapi.yml
```

#### エラー: "failed to parse collection JSON"

**症状:**
```
Error analyzing collection: failed to parse collection JSON: invalid character '}' looking for beginning of object key string
```

**原因:**
- Postmanコレクションファイルが破損している
- JSON形式が不正

**解決方法:**
1. Postmanから最新のコレクションをエクスポート
2. JSONファイルの構文を検証

```bash
# JSONファイルの構文を検証
jq . tests/tools/postman-analyzer/collections/collection.json > /dev/null

# または
python -m json.tool collection.json > /dev/null
```

#### 警告: "Method not found for operationId"

**症状:**
```
Warning: Method not found for operationId 'GetBasicInfo' in module 'auth': method not found
```

**原因:**
- `client.gen.go`が最新でない
- OpenAPI仕様とコード生成が同期していない

**解決方法:**
1. コード生成を再実行

```bash
# 全APIクライアントを再生成
./generate_all.sh

# または特定のモジュールのみ
./generate.sh authapi api/authapi.yml
```

#### 警告: "Type definition not found"

**症状:**
```
Warning: Type definition not found for type 'UpdateBasicInfoJSONRequestBody' in module 'auth'
```

**原因:**
- `types.gen.go`が最新でない
- 型名が変更されている

**解決方法:**
1. コード生成を再実行
2. OpenAPI仕様を確認

```bash
# コード生成を再実行
./generate_all.sh

# types.gen.goで型定義を確認
grep "UpdateBasicInfoJSONRequestBody" generated/authapi/types.gen.go
```

### デバッグ方法

#### 詳細なログ出力

分析プロセスの詳細を確認する場合：

```bash
# 標準エラー出力をファイルに保存
go run tests/tools/postman-analyzer/*.go -module billing 2> debug.log

# ログを確認
cat debug.log
```

#### 中間ファイルの確認

生成されたJSONファイルを確認して問題を特定：

```bash
# JSON出力を整形して確認
jq . tests/tools/postman-analyzer/reports/billing/postman_analysis.json

# 特定のフィールドを抽出
jq '.test_flow[] | {name, operation_id, client_method_name}' \
  tests/tools/postman-analyzer/reports/billing/postman_analysis.json
```

#### OpenAPI Resolverの動作確認

OpenAPI Resolverが正しく動作しているか確認：

```bash
# OpenAPI Resolverを直接実行
go run tests/tools/openapi-resolver/*.go \
  -spec api/billingapi.yml \
  -method GET \
  -endpoint "/stripe/info"
```

## エラーハンドリングとフォールバック動作

### OpenAPI連携のフォールバック

Postman Analyzerは、OpenAPI仕様が利用できない場合でも動作するように設計されています。

#### フォールバックの動作フロー

```
1. OpenAPI仕様ファイルの読み込み試行
   ↓ 失敗
2. 警告メッセージを表示
   "Warning: Failed to initialize OpenAPI resolver..."
   "Falling back to best-effort analysis..."
   ↓
3. エンドポイントからメソッド名を生成
   例: GET /stripe/info → GetStripeInfo
   ↓
4. 標準的なメソッドバリアントを生成
   - GetStripeInfo (standard)
   - GetStripeInfoWithResponse (with_response)
   - UpdateStripeInfo (standard)
   - UpdateStripeInfoWithBody (with_body)
   等
```

#### フォールバックモードの制限事項

OpenAPI仕様なしで動作する場合、以下の情報が取得できません：

- **operationId**: OpenAPI仕様で定義されたoperationId
- **正確なメソッド名**: `client.gen.go`で生成された実際のメソッド名
- **パラメータ型情報**: `types.gen.go`から抽出される型定義
- **型フィールド詳細**: パラメータ型のフィールド情報

**フォールバックモードでの出力例:**

```json
{
  "name": "GetStripeInfo",
  "method": "GET",
  "endpoint": "/stripe/info",
  "module": "billing",
  "operation_id": "",
  "client_method_name": "",
  "parameter_type_name": "",
  "parameter_type_fields": {}
}
```

### モジュール識別のフォールバック

URL変数からモジュール名を識別できない場合：

```
1. URL変数の抽出試行
   例: {{UrlBilling}} → billing
   ↓ 失敗
2. "unknown"モジュールとして処理
   ↓
3. 分析は継続（モジュール名以外の情報は取得可能）
```

### 型情報抽出のエラーハンドリング

`client.gen.go`や`types.gen.go`の解析に失敗した場合：

```
1. ファイルの存在確認
   ↓ 失敗
2. 警告メッセージを表示
   "Warning: Failed to initialize client parser..."
   ↓
3. 該当テストステップの型情報をスキップ
   ↓
4. 他のテストステップの処理は継続
```

### ベストプラクティス

1. **常にOpenAPI仕様を指定**: 最も正確な分析結果を得るため
2. **プロジェクトルートから実行**: パス解決の問題を回避
3. **コード生成を最新に保つ**: `./generate_all.sh`を定期的に実行
4. **警告メッセージを確認**: フォールバック動作が発生していないか確認
5. **CI/CDで自動化**: 定期的な分析で変更を追跡

## 関連ドキュメント

- [../docs/test-tools/README.md](../docs/test-tools/README.md) - ツールチェーン全体の概要と統合手順
- [../openapi-resolver/README.md](../openapi-resolver/README.md) - OpenAPI Resolver CLI ツール
- [../test-param-generator/README.md](../test-param-generator/README.md) - Test Param Generator ツール
