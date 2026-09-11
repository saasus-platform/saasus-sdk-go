# Test Tools - 分析・レポートツール

> **最終更新日**: 2025年10月29日  
> **バージョン**: 2.0

## 概要

Test Toolsは、E2Eテスト設計・検証・実装を支援する統合ツール群です。Postmanコレクションの分析からテストパラメータの自動生成まで、テスト開発ワークフロー全体をカバーします。

### 目的と位置づけ

- **テスト設計の効率化**: Postmanコレクションからテストフローを自動抽出
- **実装の正確性向上**: OpenAPI仕様とoapi-codegen生成コードの自動マッピング
- **保守性の向上**: 型安全なテストパラメータとヘルパー関数の自動生成
- **品質保証**: 80%以上のマッピング成功率で高精度なテストデータを生成

### E2Eテスト開発における役割

```
Postman Collection → postman-analyzer → Test Param Generator → E2E Tests
     (手動設計)         (自動分析)         (自動生成)        (実装・実行)
```

このツール群により、手動でのテストパラメータ定義やヘルパー関数の実装が不要になり、OpenAPI仕様の変更にも迅速に対応できます。

## 🔧 ツール一覧

### openapi-resolver-lib/ - OpenAPI解析共有ライブラリ

**役割**: OpenAPI仕様の解析ロジックを提供する共有ライブラリ

**主な機能**:
- OpenAPI仕様ファイルのパース
- エンドポイントからoperationIDの解決
- パス正規化とメソッドバリアント生成
- キャッシング機構による高速処理

**使用箇所**: `openapi-resolver`と`postman-analyzer`の両方で使用

### openapi-resolver/ - OpenAPI解析CLIツール

**役割**: OpenAPI仕様からoperationIDとメソッド情報を取得するCLIツール

**主な機能**:
- エンドポイントとHTTPメソッドからoperationIDを取得
- メソッドバリアント（WithResponse、WithBody等）の自動生成
- client.gen.goとtypes.gen.goのパース
- JSON/テキスト形式での出力

**主な用途**:
- 手動でのoperationID確認
- メソッド名の検証
- デバッグとトラブルシューティング

### postman-analyzer/ - Postman分析ツール

**役割**: Postmanコレクションを分析してテストフローとAPI情報を抽出

**主な機能**:
- Postmanコレクションの構造解析
- テストフローの抽出（リクエスト順序、パラメータ、期待値）
- OpenAPI仕様との連携による正確なメソッド名生成
- モジュール自動識別
- カバレッジ分析（標準メソッド、WithResponseメソッド）
- JSON/テキスト形式でのレポート出力

**主な用途**:
- テストストーリー設計の参考資料生成
- APIカバレッジの可視化
- test-param-generatorへの入力データ生成

### test-param-generator/ - テストパラメータ自動生成ツール

**役割**: postman-analyzerのレポートから、E2Eテスト用のパラメータファイルとGoヘルパー関数を自動生成

**主な機能**:
- postman-analyzerレポートの解析
- types.gen.goとの高精度マッピング（80%以上の成功率）
- カスタムマッピングルールのサポート
- 複数モジュール対応
- 診断レポートによるデバッグ支援
- test_params.jsonとGoヘルパー関数の生成

**主な用途**:
- E2Eテストのパラメータファイル生成
- テストヘルパー関数の自動生成
- 型安全なテストデータの提供

## 🚀 クイックスタート

### 基本的なワークフロー

E2Eテストのパラメータファイルを生成する標準的な手順：

```bash
# ステップ1: Postman Analyzerでレポート生成
cd tests/tools/postman-analyzer
go run . -module communication -spec ../../api/communicationapi.yml

# ステップ2: Test Param Generatorでパラメータ生成
cd ../test-param-generator
./run.sh -report ../postman-analyzer/reports/communication/postman_analysis.json -module communication

# ステップ3: 生成されたファイルを確認
ls -la ../../tests/tools/test-param-generator/tests/e2e/communication/
# → test_params.json
# → billing_test_helpers.go
```

### 実行例

#### 1. Billing APIの完全なワークフロー

```bash
# プロジェクトルートから実行
# 1. Postman分析
go run tests/tools/postman-analyzer/*.go \
  -module billing \
  -spec api/billingapi.yml

# 2. テストパラメータ生成
tests/tools/test-param-generator/run.sh \
  -report tests/tools/postman-analyzer/reports/billing/postman_analysis.json \
  -module billing

# 3. E2Eテストで使用
cd tests/e2e/billingapi
go test -v
```

#### 2. 複数モジュールの一括処理

```bash
# 全モジュールのテストパラメータを生成
for module in auth billing pricing communication; do
  echo "Processing ${module}..."
  
  # Postman分析
  go run tests/tools/postman-analyzer/*.go \
    -module ${module} \
    -spec api/${module}api.yml
  
  # パラメータ生成
  tests/tools/test-param-generator/run.sh \
    -report tests/tools/postman-analyzer/reports/${module}/postman_analysis.json \
    -module ${module} \
    -overwrite
done
```

#### 3. OpenAPI Resolverでの個別確認

```bash
# 特定のエンドポイントのoperationIDを確認
cd tests/tools/openapi-resolver
go build -o openapi-resolver .

./openapi-resolver \
  -spec ../../../api/billingapi.yml \
  -method PUT \
  -path /stripe/info

# 出力例:
# HTTPメソッド: PUT
# パス: /stripe/info
# OperationID: UpdateStripeInfo
# メソッドバリアント:
#   - UpdateStripeInfo
#   - UpdateStripeInfoWithResponse
#   - UpdateStripeInfoWithBody
#   - UpdateStripeInfoWithBodyWithResponse
```

### 推奨される実行順序

1. **OpenAPI仕様の更新**: `api/{module}api.yml`を編集
2. **コード生成**: `./generate_all.sh`を実行
3. **Postman分析**: `postman-analyzer`でレポート生成
4. **パラメータ生成**: `test-param-generator`でテストファイル生成
5. **E2Eテスト実装**: 生成されたヘルパー関数を使用してテストを記述

## 📁 ディレクトリ構造

```
tests/tools/
├── openapi-resolver-lib/              # 共有ライブラリ
│   ├── lib.go                         # OpenAPI解析の核心ロジック
│   ├── lib_test.go                    # ライブラリテスト
│   └── CHANGELOG.md                   # 変更履歴
│
├── openapi-resolver/                  # CLIツール
│   ├── main.go                        # エントリーポイント
│   ├── cli.go                         # CLI処理
│   ├── client_parser.go               # client.gen.goパーサー
│   ├── type_parser.go                 # types.gen.goパーサー
│   ├── *_test.go                      # 各種テスト
│   ├── testdata/                      # テストデータ
│   ├── README.md                      # 詳細な使用方法
│   └── VERIFICATION_RESULTS.md        # 検証結果
│
├── postman-analyzer/                  # Postman分析ツール
│   ├── analyzer.go                    # メイン処理
│   ├── module_resolver.go             # モジュール解決
│   ├── report_generator.go            # レポート生成
│   ├── collections/                   # Postmanコレクション格納
│   │   ├── SaaSus Auth API Schema.postman_collection.json
│   │   ├── SaaSus Billing API Schema.postman_collection.json
│   │   ├── SaaSus Communication API Schema.postman_collection.json
│   │   └── SaaSus Pricing API Schema.postman_collection.json
│   ├── reports/                       # 分析レポート出力先
│   │   ├── auth/
│   │   ├── billing/
│   │   ├── communication/
│   │   └── pricing/
│   └── README.md                      # 使用方法
│
├── test-param-generator/              # テストパラメータ自動生成ツール
│   ├── main.go                        # エントリーポイント
│   ├── cli.go                         # CLI処理
│   ├── analyzer.go                    # レポート解析
│   ├── parser.go                      # types.gen.goパーサー
│   ├── generator.go                   # ファイル生成
│   ├── openapi_mapper.go              # OpenAPIマッピング
│   ├── mapping_rules.go               # カスタムマッピングルール
│   ├── diagnostic_reporter.go         # 診断レポート生成
│   ├── cache.go                       # キャッシング機構
│   ├── *_test.go                      # 各種テスト
│   ├── testdata/                      # テストデータ
│   ├── tests/                         # 統合テスト
│   ├── mapping_rules.json             # カスタムマッピングルール
│   ├── mapping_rules.json.example     # サンプルルール
│   ├── run.sh                         # 実行スクリプト
│   ├── Makefile                       # ビルド設定
│   ├── README.md                      # 詳細な使用方法
│   └── ERROR_HANDLING.md              # エラーハンドリング
│
└── README.md                          # このファイル（概要・クイックスタート）
```

### 各ディレクトリの役割

| ディレクトリ | 役割 | 編集可否 |
|------------|------|---------|
| `openapi-resolver-lib/` | OpenAPI解析の共有ロジック | 可 |
| `openapi-resolver/` | CLIツール実装 | 可 |
| `postman-analyzer/` | Postman分析ツール実装 | 可 |
| `postman-analyzer/collections/` | Postmanコレクション格納 | 可（追加のみ） |
| `postman-analyzer/reports/` | 分析レポート出力先 | 自動生成 |
| `test-param-generator/` | パラメータ生成ツール実装 | 可 |
| `test-param-generator/testdata/` | テストデータ | 可 |
| `test-param-generator/tests/` | 統合テスト | 可 |

### 主要ファイルの説明

#### openapi-resolver-lib/
- `lib.go`: OpenAPI解析の核心ロジック（パース、解決、正規化、バリアント生成）
- `lib_test.go`: ライブラリの単体テスト

#### openapi-resolver/
- `main.go`: CLIのエントリーポイント
- `cli.go`: コマンドライン引数処理とメイン処理
- `client_parser.go`: client.gen.goのパース（メソッド名抽出）
- `type_parser.go`: types.gen.goのパース（型定義抽出）

#### postman-analyzer/
- `analyzer.go`: Postmanコレクションの解析メインロジック
- `module_resolver.go`: モジュール名の自動識別
- `report_generator.go`: JSON/テキストレポートの生成

#### test-param-generator/
- `analyzer.go`: postman-analyzerレポートの解析
- `parser.go`: types.gen.goのパースと型情報抽出
- `generator.go`: test_params.jsonとGoヘルパー関数の生成
- `openapi_mapper.go`: 型名とフィールド名の自動マッピング
- `mapping_rules.go`: カスタムマッピングルールの処理
- `diagnostic_reporter.go`: マッピング失敗時の診断レポート生成
- `cache.go`: 型情報のキャッシング機構

## 📊 主要機能の詳細

### postman-analyzer の分析機能

#### 1. テストフロー抽出
Postmanコレクションからテストの実行順序とパラメータを抽出：

```
=== Test Flow ===
1. Pre_GetStripeConnectionInformation
   Method: GET /stripe/info
   Expected Status: 200
   Description: Initial check of Stripe connection status

2. UpdateStripeConnectionInfo
   Method: PUT /stripe/info
   Expected Status: 200
   Description: Update Stripe connection information
   Body Parameters:
     - secret_key: string
```

#### 2. メソッドカバレッジ分析
生成されたクライアントメソッドのカバレッジを可視化：

```
=== Method Coverage Analysis ===
Total Client Methods: 8
Standard Methods: 4
WithResponse Methods: 4
Coverage Percentage: 100.0%

=== Covered Methods ===
✅ DeleteStripeInfo
✅ GetStripeInfo
✅ UpdateStripeInfo
✅ UpdateStripeInfoWithBody
```

#### 3. モジュール自動識別
コレクション名やエンドポイントパスからモジュール名を自動識別：

```
Detected Module: billing
Spec Path: api/billingapi.yml
Types Path: generated/billingapi/types.gen.go
```

#### 4. レポート出力
- **JSON形式** (`reports/{module}/postman_analysis.json`): プログラムでの処理用
- **テキスト形式** (`reports/{module}/postman_analysis_report.txt`): 人間が読む用

### test-param-generator の生成機能

#### 1. 高精度マッピング（80%以上の成功率）
50以上の型名パターンと6段階のフィールドマッチング戦略：

```
マッピング統計:
グループ: 20/26 成功 (76.9%)
フィールド: 70/88 マッピング成功 (79.5%)

パターン使用統計:
  完全一致パターン: 15回使用
  標準パターン: 8回使用
  バリエーションパターン: 3回使用
```

#### 2. カスタムマッピングルール
JSON形式でマッピングルールを定義可能：

```json
{
  "typeNameMappings": {
    "users.Create": "CreateSaasUserJSONRequestBody"
  },
  "fieldNameMappings": {
    "CreateSaasUserJSONRequestBody": {
      "email": "Email",
      "password": "Password"
    }
  },
  "resourceAliases": {
    "users": ["saasusers", "saas_users"]
  }
}
```

#### 3. 診断レポート
マッピング失敗時の詳細情報と改善提案：

```
--- 失敗したマッピング ---

1. グループ: users, 操作: Create, フィールド: email
   試行したパターン:
     - CreateUserJSONRequestBody (存在しない)
     - CreateUserParams (存在しない)
   最も近かったマッチ: CreateSaasUserJSONRequestBody.Email (信頼度: 0.6)
   提案: リソース名を "saasusers" に変更するか、カスタムマッピングルールを追加
```

#### 4. 複数モジュール対応
レポート内で参照された全モジュールの型定義を自動的に読み込み：

```
エントリーモジュール: pricing
追加モジュール検出: auth, billing
読み込んだ型定義:
  - generated/pricingapi/types.gen.go (125 types)
  - generated/authapi/types.gen.go (89 types)
  - generated/billingapi/types.gen.go (45 types)
```

### openapi-resolver の解析機能

#### 1. OperationID解決
エンドポイントとHTTPメソッドからoperationIDを取得：

```bash
$ ./openapi-resolver -spec api/billingapi.yml -method PUT -path /stripe/info

HTTPメソッド: PUT
パス: /stripe/info
OperationID: UpdateStripeInfo
リクエストボディ: true
概要: Update Stripe Connection Info
```

#### 2. メソッドバリアント生成
oapi-codegenが生成する全メソッドバリアントを列挙：

```
メソッドバリアント:
  - UpdateStripeInfo
  - UpdateStripeInfoWithResponse
  - UpdateStripeInfoWithBody
  - UpdateStripeInfoWithBodyWithResponse
```

#### 3. 全エンドポイント列挙
API仕様内の全エンドポイントをリスト表示：

```bash
$ ./openapi-resolver -spec api/billingapi.yml -list

=== All Endpoints ===
GET    /stripe/info          GetStripeInfo
PUT    /stripe/info          UpdateStripeInfo
DELETE /stripe/info          DeleteStripeInfo
...
```

## 🎯 使用場面

### 1. テスト設計段階

**目的**: Postmanコレクションからテストストーリーを設計

```bash
# ステップ1: Postmanコレクションを分析
cd tests/tools/postman-analyzer
go run . -module billing -spec ../../api/billingapi.yml

# ステップ2: 生成されたレポートを確認
cat reports/billing/postman_analysis_report.txt

# ステップ3: テストストーリーを設計
# - テストフローセクションを参考にストーリーを作成
# - カバレッジ分析から不足しているテストケースを特定
# - メソッドバリアントの使い分けを計画
```

**活用例**:
- E2Eテストのストーリー構成を決定
- テストの実行順序を最適化
- 必要なテストデータを特定

### 2. 実装検証

**目的**: OpenAPI仕様変更後の影響確認とテストパラメータの更新

```bash
# ステップ1: OpenAPI仕様を更新
vim api/billingapi.yml

# ステップ2: コード生成
./generate_all.sh

# ステップ3: Postman分析を再実行
go run tests/tools/postman-analyzer/*.go -module billing -spec api/billingapi.yml

# ステップ4: テストパラメータを再生成
tests/tools/test-param-generator/run.sh \
  -report tests/tools/postman-analyzer/reports/billing/postman_analysis.json \
  -module billing \
  -overwrite

# ステップ5: E2Eテストを実行して検証
cd tests/e2e/billingapi
go test -v
```

**活用例**:
- 新規エンドポイント追加時のテストパラメータ自動生成
- 型定義変更時のマッピング検証
- メソッド名変更の影響範囲確認

### 3. CI/CD統合

**目的**: 継続的なテストデータの検証と自動更新

#### パターンA: テストパラメータの自動検証

```bash
#!/bin/bash
# .github/workflows/validate-test-params.yml

# Postman分析
for module in auth billing pricing communication; do
  go run tests/tools/postman-analyzer/*.go \
    -module ${module} \
    -spec api/${module}api.yml
done

# テストパラメータ生成（dry-run）
for module in auth billing pricing communication; do
  tests/tools/test-param-generator/test-param-generator \
    -report tests/tools/postman-analyzer/reports/${module}/postman_analysis.json \
    -module ${module} \
    -output /tmp/test-params-${module}
done

# マッピング成功率を確認
# 80%未満の場合はビルド失敗
```

#### パターンB: OpenAPI仕様変更時の自動更新

```bash
#!/bin/bash
# scripts/update-test-params.sh

# OpenAPI仕様が変更された場合に実行
if git diff --name-only HEAD~1 | grep -q "^api/.*\.yml$"; then
  echo "OpenAPI仕様が変更されました。テストパラメータを更新します..."
  
  # コード生成
  ./generate_all.sh
  
  # 全モジュールのテストパラメータを再生成
  for module in auth billing pricing communication; do
    go run tests/tools/postman-analyzer/*.go \
      -module ${module} \
      -spec api/${module}api.yml
    
    tests/tools/test-param-generator/run.sh \
      -report tests/tools/postman-analyzer/reports/${module}/postman_analysis.json \
      -module ${module} \
      -overwrite
  done
  
  # 変更をコミット
  git add tests/tools/test-param-generator/tests/e2e/*/test_params.json
  git add tests/tools/test-param-generator/tests/e2e/*/*_test_helpers.go
  git commit -m "chore: update test parameters after OpenAPI spec changes"
fi
```

#### パターンC: カバレッジレポートの生成

```bash
#!/bin/bash
# CI/CDでカバレッジレポートを生成

# 全モジュールのカバレッジを集計
for module in auth billing pricing communication; do
  go run tests/tools/postman-analyzer/*.go \
    -module ${module} \
    -spec api/${module}api.yml
  
  # カバレッジ情報を抽出
  coverage=$(cat tests/tools/postman-analyzer/reports/${module}/postman_analysis.json | \
    jq '.coverage_percentage')
  
  echo "${module}: ${coverage}%"
done
```

**活用例**:
- プルリクエスト時のテストパラメータ検証
- OpenAPI仕様変更の自動検出と対応
- テストカバレッジの継続的な監視

### 4. デバッグとトラブルシューティング

**目的**: マッピング失敗の原因特定と解決

```bash
# ステップ1: 診断レポートを生成
tests/tools/test-param-generator/test-param-generator \
  -report tests/tools/postman-analyzer/reports/auth/postman_analysis.json \
  -module auth \
  -debug \
  -verbose

# ステップ2: 診断レポートを確認
cat diagnostic_report_*.txt

# ステップ3: カスタムマッピングルールを作成
cp tests/tools/test-param-generator/mapping_rules.json.example \
   tests/tools/test-param-generator/mapping_rules.json

# ステップ4: ルールを編集
vim tests/tools/test-param-generator/mapping_rules.json

# ステップ5: 再実行して検証
tests/tools/test-param-generator/test-param-generator \
  -report tests/tools/postman-analyzer/reports/auth/postman_analysis.json \
  -module auth \
  -rules tests/tools/test-param-generator/mapping_rules.json \
  -verbose
```

**活用例**:
- マッピング失敗の原因特定
- カスタムルールの作成と検証
- 型定義の不整合の発見

## 📋 出力ファイル一覧

### postman-analyzer の出力

| ファイル | 形式 | 用途 | 配置場所 |
|---------|------|------|---------|
| `postman_analysis.json` | JSON | test-param-generatorへの入力、プログラム処理 | `reports/{module}/` |
| `postman_analysis_report.txt` | テキスト | 人間が読む、ドキュメント、レビュー資料 | `reports/{module}/` |

**JSON構造例**:
```json
{
  "collection_name": "SaaSus Billing API Schema",
  "module": "billing",
  "test_flow": [
    {
      "name": "GetStripeInfo",
      "method": "GET",
      "endpoint": "/stripe/info",
      "expected_status": 200,
      "body_params": null
    }
  ],
  "client_methods": ["GetStripeInfo", "GetStripeInfoWithResponse", ...],
  "coverage_percentage": 100.0
}
```

### test-param-generator の出力

| ファイル | 形式 | 用途 | 配置場所 |
|---------|------|------|---------|
| `test_params.json` | JSON | E2Eテストのパラメータデータ | `tests/tools/test-param-generator/tests/e2e/{module}/` |
| `{module}_test_helpers.go` | Go | テストヘルパー関数 | 同上 |
| `diagnostic_report_{timestamp}.txt` | テキスト | マッピング診断レポート（-debugオプション時） | カレントディレクトリ |

**test_params.json構造例**:
```json
{
  "stripe": {
    "updateParams": {
      "secret_key": "sk_test_example123"
    }
  }
}
```

**ヘルパー関数例**:
```go
package billingapi

type TestParams struct {
    Stripe struct {
        UpdateParams struct {
            SecretKey string `json:"secret_key"`
        } `json:"updateParams"`
    } `json:"stripe"`
}

func loadTestParams(t *testing.T) *TestParams { ... }
func stringPtr(s string) *string { ... }
```

## 🎯 利点

### 開発効率の向上

1. **手動作業の削減**: テストパラメータとヘルパー関数を自動生成
2. **迅速な対応**: OpenAPI仕様変更後、数分でテストパラメータを更新
3. **学習コストの低減**: 標準的なワークフローで誰でも使用可能

### 品質の向上

4. **型安全性**: types.gen.goとの自動マッピングによる型の整合性保証
5. **高精度**: 80%以上のマッピング成功率
6. **一貫性**: 共有ライブラリにより全ツールで同じロジックを使用

### 保守性の向上

7. **変更の局所化**: OpenAPI仕様の変更が自動的にテストに反映
8. **トレーサビリティ**: Postmanコレクション → レポート → テストパラメータの追跡可能
9. **診断機能**: マッピング失敗時の詳細な診断レポート

### CI/CD統合

10. **自動化**: パイプラインに組み込み可能
11. **検証**: テストパラメータの妥当性を自動チェック
12. **可視化**: カバレッジレポートの自動生成

## ⚠️ 制限事項と注意点

### 1. Postmanコレクションの構造

- **前提**: Postmanコレクションが適切に構造化されている必要があります
- **推奨**: フォルダ構造でエンドポイントをグループ化
- **注意**: テスト名とエンドポイントパスの一貫性が重要

### 2. マッピング精度

- **成功率**: 標準的な命名規則で80%以上
- **カスタムルール**: 特殊な命名規則にはカスタムマッピングルールが必要
- **診断**: `-debug`オプションで失敗原因を特定可能

### 3. 複数モジュール対応

- **自動検出**: レポート内のmodule情報から自動的に複数モジュールを読み込み
- **優先順位**: エントリーモジュール → 追加モジュールの順で型定義をマージ
- **注意**: 同名の型が複数モジュールに存在する場合、最初に見つかった型を使用

### 4. パフォーマンス

- **キャッシング**: 型情報はキャッシュされ、2回目以降は高速
- **大規模API**: 数百のエンドポイントでも数秒で処理完了
- **メモリ**: types.gen.goのサイズに比例してメモリ使用量が増加

## 🔧 トラブルシューティング

よくある問題と解決方法：

| 問題 | 原因 | 解決方法 |
|-----|------|---------|
| レポートファイルが見つからない | postman-analyzerが未実行 | postman-analyzerを実行してレポート生成 |
| 型定義ファイルが見つからない | コード生成が未実行 | `./generate_all.sh`を実行 |
| マッピング成功率が低い | 命名規則が特殊 | カスタムマッピングルールを作成 |
| 生成されたGoコードがエラー | パッケージ名の不一致 | 生成されたファイルのパッケージ名を確認 |

詳細なトラブルシューティングは各ツールのREADMEを参照してください。

## 📚 詳細ドキュメント

### 統合・アーキテクチャ
- [docs/test-tools/README.md](../docs/test-tools/README.md) - クイックスタート、診断、改善ポイントのサマリー

### 各ツールの詳細
- [openapi-resolver/README.md](openapi-resolver/README.md) - CLIツールの詳細な使用方法、オプション、トラブルシューティング
- [openapi-resolver/VERIFICATION_RESULTS.md](openapi-resolver/VERIFICATION_RESULTS.md) - 検証結果と動作確認
- [postman-analyzer/README.md](postman-analyzer/README.md) - 分析ツールの使用方法、出力形式、モジュール識別
- [test-param-generator/README.md](test-param-generator/README.md) - パラメータ生成ツールの包括的なドキュメント
- [test-param-generator/ERROR_HANDLING.md](test-param-generator/ERROR_HANDLING.md) - エラーハンドリングの詳細

### 変更履歴
- [openapi-resolver-lib/CHANGELOG.md](openapi-resolver-lib/CHANGELOG.md) - 共有ライブラリの変更履歴

## 🚀 次のステップ

### 初めて使用する場合

1. **クイックスタートを実行**: 上記の「クイックスタート」セクションに従って基本的なワークフローを体験
2. **各ツールのREADMEを確認**: 詳細な機能とオプションを理解
3. **サンプルを確認**: `tests/e2e/`配下の既存のテストパラメータを参照

### 既存プロジェクトへの導入

1. **Postmanコレクションを準備**: 既存のPostmanコレクションを`postman-analyzer/collections/`に配置
2. **分析を実行**: postman-analyzerでレポートを生成
3. **パラメータを生成**: test-param-generatorでテストパラメータを生成
4. **E2Eテストを更新**: 生成されたヘルパー関数を使用してテストを書き換え

### カスタマイズが必要な場合

1. **診断レポートを確認**: `-debug`オプションでマッピング失敗の原因を特定
2. **カスタムルールを作成**: `mapping_rules.json`でマッピングルールを定義
3. **検証**: `-verbose`オプションで詳細ログを確認しながら調整

## 💡 ベストプラクティス

### Postmanコレクションの管理

- **命名規則**: テスト名にエンドポイントパスを含める（例: `GET /users/{id}`）
- **フォルダ構造**: リソースごとにフォルダを作成
- **説明の記載**: 各テストに目的と期待値を記載

### テストパラメータの管理

- **バージョン管理**: `test_params.json`と`*_test_helpers.go`をGitで管理
- **定期的な更新**: OpenAPI仕様変更時は必ず再生成
- **レビュー**: 生成されたパラメータが適切か確認

### CI/CDでの活用

- **自動検証**: プルリクエスト時にマッピング成功率をチェック
- **自動更新**: OpenAPI仕様変更時にテストパラメータを自動更新
- **カバレッジ監視**: 定期的にカバレッジレポートを生成

## 📞 サポート

問題が発生した場合：

1. **ドキュメントを確認**: 各ツールのREADMEとトラブルシューティングセクション
2. **診断レポートを生成**: `-debug`オプションで詳細情報を取得
3. **Issueを作成**: 上記で解決しない場合は、診断レポートを添付してIssueを作成

## 📄 ライセンス

このツール群はプロジェクトのライセンスに従います。
