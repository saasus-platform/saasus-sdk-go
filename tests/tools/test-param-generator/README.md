# Test Param Generator

Test Param Generator は、postman-analyzer が生成したレポートファイルを解析し、E2E テスト用のパラメータファイル（JSON）と Go のヘルパー関数を自動生成するコマンドラインツールです。

## 概要

このツールは以下の機能を提供します：

- Postman Analyzer レポート（JSON）の解析
- 生成された API 型定義（`generated/{module}api/types.gen.go`）との自動マッピング
- テストパラメータの JSON 形式での出力
- テストパラメータを読み込む Go ヘルパー関数の生成
- エンドポイントごとのパラメータグループ化

## クイックスタート（5分）

1. Postman Analyzer を実行し、`tests/tools/postman-analyzer/reports/<module>/postman_analysis.json` を用意する。
2. `cd tests/tools/test-param-generator`。
3. `go run . -report ../postman-analyzer/reports/<module>/postman_analysis.json -module <module> -overwrite` を実行する。
4. `tests/tools/test-param-generator/tests/e2e/<module>/` に生成された `test_params.json` と `<module>_test_helpers.go` を確認してテストに取り込む。

繰り返し実行する場合は `go build -o test-param-generator .` でバイナリを作成しておくと便利です。

### 新機能（v2.0）

**拡張型マッピングシステム**

- **高度なパターンマッチング**: 50 以上の型名パターンを自動生成し、優先順位付きでマッチング
- **スマートフィールドマッチング**: 6 段階のマッチング戦略（完全一致、正規化、大文字小文字無視、部分一致、意味的一致、あいまい一致）
- **詳細な診断レポート**: マッピング失敗時の詳細情報と改善提案を自動生成
- **カスタムマッピングルール**: JSON 形式でマッピングルールを定義可能
- **型メタデータ分析**: 型名から HTTP メソッド、リソース名、型カテゴリを自動推測
- **パフォーマンス最適化**: 型情報のキャッシング機構により高速処理を実現

**マッピング成功率の大幅改善**

従来のマッピング成功率を大幅に改善：

| モジュール    | 従来  | v2.0 | 改善率 |
| ------------- | ----- | ---- | ------ |
| auth          | 13.6% | 80%+ | +66.4% |
| billing       | 0%    | 80%+ | +80%   |
| pricing       | 0%    | 80%+ | +80%   |
| communication | 71.4% | 90%+ | +18.6% |

**診断とデバッグ機能**

- マッピング統計の詳細レポート
- 失敗したマッピングの試行履歴
- パターン使用統計とマッチング戦略統計
- 改善提案の自動生成
- 診断レポートのファイル出力

**Module情報の活用（v2.1）**

- **正確なファイルパス解決**: postman-analyzerが提供するmodule情報を優先的に使用して、types.gen.goファイルを正確に特定
  - 形式: `generated/{module}api/types.gen.go`
  - Module情報からのパス構築が失敗した場合は、従来のフォールバック探索を実行
  - デバッグログで試行したパスと発見したファイルを詳細に出力
- **モジュール固有マッピングルール**: 各モジュールに特化したカスタムマッピングルールを定義可能
- **モジュール別統計**: モジュールごとのマッピング成功率と失敗パターンを集計
- **改善されたエラーメッセージ**: Module情報と試行したファイルパスを含む詳細なエラーメッセージ

## インストール

### 前提条件

- Go 1.20 以上
- postman-analyzer ツール（レポート生成用）

### ビルド方法

```bash
# プロジェクトルートからツールのディレクトリへ移動
cd tests/tools/test-param-generator

# go run で直接実行（最も簡単）
go run . -report <report_path> -module <module>

# 再利用する場合はバイナリをビルド
go build -o test-param-generator .
# 例: ./test-param-generator -report ...
```

## 使用方法

### 基本的な使い方

```bash
# 実行スクリプトを使用（自動ビルド）
./run.sh -report <レポートファイルパス> -module <モジュール名>

# または、ビルド済みバイナリを直接実行
./test-param-generator -report <レポートファイルパス> -module <モジュール名>

# または、プロジェクトルートから
tests/tools/test-param-generator/run.sh -report <レポートファイルパス> -module <モジュール名>
```

### コマンドラインオプション

| オプション   | 説明                                       | 必須 | デフォルト値         |
| ------------ | ------------------------------------------ | ---- | -------------------- |
| `-report`    | Postman Analyzer レポートファイルのパス    | ✓    | -                    |
| `-module`    | モジュール名（例: communication, billing） | ✓    | -                    |
| `-output`    | 出力ディレクトリのパス                     |      | `tests/tools/test-param-generator/tests/e2e/{module}` |
| `-overwrite` | 既存ファイルを確認なしで上書き             |      | false                |
| `-verbose`   | 詳細ログを出力（マッピング試行履歴を含む） |      | false                |
| `-help`      | ヘルプメッセージを表示                     |      | -                    |

> ※ 詳細診断を強制したい場合は `TPG_FORCE_MAPPING_DIAGNOSTICS=1` 環境変数を指定してください。カスタムマッピングルールは `mapping_rules.json` を編集するか `-output` 先に配置することで読み込まれます。

### 使用例

#### 1. Communication API のテストパラメータ生成

```bash
# 実行スクリプトを使用
./run.sh \
  -report tests/tools/postman-analyzer/reports/communication/postman_analysis.json \
  -module communication

# または、プロジェクトルートから
tests/tools/test-param-generator/run.sh \
  -report tests/tools/postman-analyzer/reports/communication/postman_analysis.json \
  -module communication
```

**出力:**

- `tests/tools/test-param-generator/tests/e2e/communication/test_params.json` - テストパラメータファイル
- `tests/tools/test-param-generator/tests/e2e/communication/communication_test_helpers.go` - ヘルパー関数

#### 2. Billing API のテストパラメータ生成（カスタム出力先）

```bash
./test-param-generator \
  -report tests/tools/postman-analyzer/reports/billing/postman_analysis.json \
  -module billing \
  -output tests/tools/test-param-generator/tests/e2e/custom-billing
```

#### 3. 既存ファイルを上書き（確認なし）

```bash
./test-param-generator \
  -report tests/tools/postman-analyzer/reports/auth/postman_analysis.json \
  -module auth \
  -overwrite
```

#### 4. 詳細ログ付きで実行

```bash
./test-param-generator \
  -report tests/tools/postman-analyzer/reports/pricing/postman_analysis.json \
  -module pricing \
  -verbose
```

#### 5. 一括実行してマッピング率を可視化

```bash
# プロジェクトルートで実行
./run_test_param_generator_all_modules_with_metrics.sh
```

**出力:**

- `tests/tools/test-param-generator/tests/e2e/{module}/test_params.json` および `{module}_test_helpers.go`
- `tests/tools/test-param-generator/mapping_metrics_*.txt`（モジュールごとのマッピング成功率レポート）

## 出力ファイル

### 1. test_params.json

テストで使用するパラメータを定義した JSON ファイル。エンドポイントごとにグループ化され、以下の構造を持ちます：

```json
{
  "feedbacks": {
    "createParams": {
      "user_id": "69e732d6-8ecc-45c4-c2eb-8438f7ffe775",
      "feedback_title": "テスト用フィードバック",
      "feedback_description": "これはE2Eテスト用のフィードバックです。"
    },
    "updateParams": {
      "feedback_title": "更新されたテスト用フィードバック",
      "feedback_description": "これは更新されたE2Eテスト用のフィードバックです。"
    },
    "testFeedbackIds": ["12345678-1234-1234-1234-123456789012"],
    "invalidFeedbackIds": ["invalid-id", "not-a-uuid"]
  }
}
```

### 2. {module}\_test_helpers.go

テストパラメータを読み込む Go ヘルパー関数を含むファイル：

```go
package communicationapi

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// TestParams はテストパラメータファイルの構造体
type TestParams struct {
	Feedbacks struct {
		CreateParams struct {
			UserId              string `json:"user_id"`
			FeedbackTitle       string `json:"feedback_title"`
			FeedbackDescription string `json:"feedback_description"`
		} `json:"createParams"`
		UpdateParams struct {
			FeedbackTitle       string `json:"feedback_title"`
			FeedbackDescription string `json:"feedback_description"`
		} `json:"updateParams"`
		TestFeedbackIds    []string `json:"testFeedbackIds"`
		InvalidFeedbackIds []string `json:"invalidFeedbackIds"`
	} `json:"feedbacks"`
}

// loadTestParams はテストパラメータファイルを読み込みます
func loadTestParams(t *testing.T) *TestParams {
	paramFile := filepath.Join("test_params.json")
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

// stringPtr は文字列のポインタを返すヘルパー関数
func stringPtr(s string) *string {
	return &s
}
```

## ワークフロー

Test Param Generator は以下のワークフローで動作します：

```
1. Postman Analyzerレポートの読み込み
   ↓
2. test_flowセクションからパラメータ情報を抽出
   ↓
3. Module情報を使用して型定義ファイル（types.gen.go）のパスを構築
   - 優先: generated/{module}api/types.gen.go
   - フォールバック: 複数のパターンで探索
   ↓
4. 生成された型定義ファイルの解析
   ↓
5. エンドポイントごとにパラメータをグループ化
   ↓
6. 型定義とのマッピング（カスタムルール → 自動パターン）
   ↓
7. サンプル値の生成
   ↓
8. test_params.jsonの出力
   ↓
9. Goヘルパー関数の生成
```

## パラメータグループ化のルール

ツールは以下のルールでパラメータをグループ化します：

1. **エンドポイントパスからリソース名を抽出**

   - `/feedbacks` → `feedbacks`
   - `/feedbacks/{feedback_id}/comments` → `comments`

2. **HTTP メソッドに基づく分類**

   - `POST` → `createParams`
   - `PATCH`/`PUT` → `updateParams`
   - `GET`/`DELETE` → `testIds`

3. **同じリソースに対する複数の操作をグループ化**

## カスタムマッピングルール

標準のパターンマッチングで型定義が見つからない場合、カスタムマッピングルールを使用できます。カスタムルールは最優先で適用されるため、特殊な命名規則や複雑なマッピングに対応できます。

### ルールファイルの作成

`mapping_rules.json`ファイルを作成して、カスタムマッピングルールを定義します：

```json
{
  "typeNameMappings": {
    "users.Create": "CreateSaasUserJSONRequestBody",
    "users.Update": "UpdateSaasUserJSONRequestBody",
    "tenants.Create": "CreateTenantJSONRequestBody"
  },
  "fieldNameMappings": {
    "CreateSaasUserJSONRequestBody": {
      "email": "Email",
      "password": "Password"
    }
  },
  "resourceAliases": {
    "users": ["saasusers", "saas_users"],
    "tenants": ["tenant"]
  },
  "fieldAliases": {
    "desc": ["description", "Description"],
    "id": ["ID", "Id", "identifier"]
  }
}
```

### ルールファイルの配置

ルールファイルは以下のいずれかの場所に配置してください（優先順位順）：

1. `-rules`オプションで指定したパス（最優先）
2. `./mapping_rules.json` (カレントディレクトリ)
3. `../mapping_rules.json` (親ディレクトリ)
4. `../../mapping_rules.json` (2 階層上)
5. `./tests/tools/test-param-generator/mapping_rules.json`

ツールは上記の順序でファイルを検索し、最初に見つかったファイルを使用します。

### ルールの種類

#### 1. typeNameMappings（型名マッピング）

リソース名と操作から型名への直接マッピングを定義します。これは最も強力なルールで、標準パターンマッチングを完全にバイパスします。

**形式:**

```json
{
  "typeNameMappings": {
    "リソース名.操作": "型名"
  }
}
```

**使用例:**

```json
{
  "typeNameMappings": {
    "users.Create": "CreateSaasUserJSONRequestBody",
    "users.Update": "UpdateSaasUserJSONRequestBody",
    "feedbacks.Update": "UpdateFeedbackJSONRequestBody",
    "tenants.Create": "CreateTenantJSONRequestBody"
  }
}
```

**適用場面:**

- 標準パターンと異なる命名規則を使用している場合
- リソース名と型名の対応が複雑な場合
- 特定のエンドポイントだけ特殊な型を使用している場合

#### 2. fieldNameMappings（フィールド名マッピング）

型名とリクエストフィールド名から実際の Go フィールド名への直接マッピングを定義します。

**形式:**

```json
{
  "fieldNameMappings": {
    "型名": {
      "リクエストフィールド名": "実際のGoフィールド名"
    }
  }
}
```

**使用例:**

```json
{
  "fieldNameMappings": {
    "CreateSaasUserJSONRequestBody": {
      "email": "Email",
      "password": "Password",
      "user_attributes": "Attributes"
    },
    "UpdateTenantJSONRequestBody": {
      "tenant_name": "Name",
      "back_office_staff_email": "BackOfficeStaffEmail"
    }
  }
}
```

**適用場面:**

- JSON タグと異なるフィールド名を使用している場合
- 省略形や特殊な命名規則を使用している場合
- フィールド名の自動変換が失敗する場合

#### 3. resourceAliases（リソース名の別名）

リソース名の別名を定義します。これにより、複数の表記を同じリソースとして扱えます。

**形式:**

```json
{
  "resourceAliases": {
    "リソース名": ["別名1", "別名2", "別名3"]
  }
}
```

**使用例:**

```json
{
  "resourceAliases": {
    "users": ["saasusers", "saas_users", "user"],
    "tenants": ["tenant"],
    "feedbacks": ["feedback"]
  }
}
```

**適用場面:**

- エンドポイントパスと types.gen.go の型名でリソース名の表記が異なる場合
- 単数形・複数形の両方に対応したい場合
- プレフィックス付きのリソース名を使用している場合

#### 4. fieldAliases（フィールド名の別名）

フィールド名の別名を定義します。これにより、複数の表記を同じフィールドとして扱えます。

**形式:**

```json
{
  "fieldAliases": {
    "フィールド名": ["別名1", "別名2", "別名3"]
  }
}
```

**使用例:**

```json
{
  "fieldAliases": {
    "desc": ["description", "Description", "desc"],
    "id": ["ID", "Id", "identifier"],
    "name": ["Name", "name", "display_name"],
    "email": ["Email", "email", "email_address"]
  }
}
```

**適用場面:**

- 省略形と完全形の両方に対応したい場合
- 複数の命名規則が混在している場合
- 意味的に同じフィールドに異なる名前が使われている場合

#### 5. moduleSpecificRules（モジュール固有ルール）

各モジュールに特化したマッピングルールを定義します。グローバルルールよりも優先されます。

**形式:**

```json
{
  "moduleSpecificRules": {
    "モジュール名": {
      "typeNameMappings": { ... },
      "fieldNameMappings": { ... },
      "resourceAliases": { ... },
      "fieldAliases": { ... }
    }
  }
}
```

**使用例:**

```json
{
  "moduleSpecificRules": {
    "auth": {
      "typeNameMappings": {
        "saasusers.Create": "CreateSaasUserJSONRequestBody",
        "tenants.Create": "CreateTenantJSONRequestBody"
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
    },
    "communication": {
      "typeNameMappings": {
        "feedbacks.Create": "CreateFeedbackJSONRequestBody",
        "comments.Create": "CreateFeedbackCommentJSONRequestBody"
      },
      "fieldNameMappings": {
        "CreateFeedbackJSONRequestBody": {
          "user_id": "UserId",
          "feedback_title": "FeedbackTitle"
        }
      }
    }
  }
}
```

**適用場面:**

- モジュールごとに異なる命名規則を使用している場合
- 特定のモジュールでのみ必要なマッピングルールがある場合
- グローバルルールでは対応できないモジュール固有の問題がある場合

**サンプルファイル:**

完全なモジュール固有ルールの例は `mapping_rules.json.module_example` を参照してください。

## 型定義ファイルの複数モジュール対応

Postman Analyzer のレポートには、各ステップに `module` フィールドが含まれています。Test Param Generator はレポート全体を解析し、該当のステップで参照されたモジュールを自動的に収集します。収集したモジュールごとに `generated/{module}api/types.gen.go` を読み込み、型定義を統合したうえでマッピングを実行します。

- エントリーモジュール（CLI で指定した `-module`）は常に最優先で読み込み
- レポート内で参照された追加モジュールは順に読み込み、既存の型を上書きしない形でマージ
- OpenAPI マッピングやカスタムマッピングルールは、参照モジュールの一覧を順番に試行

例えば `pricing` モジュールのテストフローに `auth` モジュールの API が含まれている場合、両モジュールの `types.gen.go` を読み込み、`CreateTenantJSONRequestBody` など `auth` 定義の型もマッピング対象として利用できるようになります。追加設定は不要です。

### マッピングの優先順位

ツールは以下の優先順位でマッピングを試行します：

1. **モジュール固有カスタムマッピングルール** (最優先)

   - moduleSpecificRules[moduleName].typeNameMappings
   - moduleSpecificRules[moduleName].fieldNameMappings
   - moduleSpecificRules[moduleName].resourceAliases

2. **グローバルカスタムマッピングルール**

   - typeNameMappings（型名の直接マッピング）
   - fieldNameMappings（フィールド名の直接マッピング）
   - resourceAliases（リソース名の別名展開）
   - fieldAliases（フィールド名の別名展開）

2. **自動パターンマッチング**

   - **完全一致パターン** (優先度: 0-9)
     - `CreateFeedbackJSONRequestBody`
     - `CreateFeedbacksJSONRequestBody`
   - **標準パターン** (優先度: 10-19)
     - `CreateFeedbackParams`
     - `CreateFeedbackRequest`
     - `FeedbackCreateParams`
   - **バリエーションパターン** (優先度: 20-29)
     - `CreateFeedbacksParams`
     - `NewFeedbackParams`
   - **あいまい一致パターン** (優先度: 30-39)
     - `*Feedback*Create*`
     - `*Create*Feedback*`

3. **フィールドマッチング戦略**

   - **完全一致** (信頼度: 1.0) - JSON タグとの完全一致
   - **正規化一致** (信頼度: 0.9) - snake_case ↔ camelCase 変換
   - **大文字小文字無視** (信頼度: 0.8) - 大文字小文字を無視
   - **部分一致** (信頼度: 0.6) - 前方一致、後方一致
   - **意味的一致** (信頼度: 0.5) - シノニム対応
   - **あいまい一致** (信頼度: 0.3) - レーベンシュタイン距離

4. **デフォルト値** (すべて失敗した場合)
   - 型推論に基づくサンプル値を生成

### カスタムルールの作成手順

#### ステップ 1: 診断レポートを確認

まず、`-debug`オプションで診断レポートを生成します：

```bash
./test-param-generator \
  -report tests/tools/postman-analyzer/reports/auth/postman_analysis.json \
  -module auth \
  -debug
```

診断レポート（`diagnostic_report_{timestamp}.txt`）を確認し、失敗したマッピングを特定します。

#### ステップ 2: サンプルファイルをコピー

```bash
cp tests/tools/test-param-generator/mapping_rules.json.example \
   tests/tools/test-param-generator/mapping_rules.json
```

#### ステップ 3: ルールを追加

診断レポートの情報を基に、必要なルールを追加します：

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
  }
}
```

#### ステップ 4: 再実行して確認

```bash
./test-param-generator \
  -report tests/tools/postman-analyzer/reports/auth/postman_analysis.json \
  -module auth \
  -rules ./mapping_rules.json \
  -verbose
```

マッピング成功率が改善されたことを確認します。

### サンプルファイル

`mapping_rules.json.example`ファイルがサンプルとして提供されています：

```bash
# サンプルファイルをコピーして使用
cp tests/tools/test-param-generator/mapping_rules.json.example \
   tests/tools/test-param-generator/mapping_rules.json

# 必要に応じて編集
vim tests/tools/test-param-generator/mapping_rules.json
```

### ベストプラクティス

1. **段階的にルールを追加**: 一度にすべてのルールを追加せず、失敗したマッピングから順番に対応
2. **診断レポートを活用**: `-debug`オプションで詳細情報を確認してから対応
3. **汎用的なルールを優先**: resourceAliases や fieldAliases で汎用的に対応できる場合は、それを優先
4. **ルールファイルをバージョン管理**: プロジェクトに含めて、チーム全体で共有
5. **定期的に見直し**: types.gen.go の更新に合わせてルールも更新

## 型推論ルール

パラメータ名から型を自動推論します：

| パラメータ名パターン            | 推論される型  | サンプル値                             |
| ------------------------------- | ------------- | -------------------------------------- |
| `*_id`, `*Id`                   | UUID (string) | `69e732d6-8ecc-45c4-c2eb-8438f7ffe775` |
| `*_title`, `*Title`             | string        | `テスト用タイトル`                     |
| `*_description`, `*Description` | string        | `これはテスト用の説明です。`           |
| `*_body`, `*Body`               | string        | `これはテスト用の本文です。`           |
| `*_name`, `*Name`               | string        | `テスト用名前`                         |
| `status`                        | int/uint64    | `0` または `1`                         |
| `count`, `limit`                | int           | `10`                                   |
| `enabled`, `is_*`               | bool          | `true`                                 |

## エラーハンドリング

ツールは以下のエラーを適切に処理します：

| エラー種別             | 発生条件                     | 対応方法                                     |
| ---------------------- | ---------------------------- | -------------------------------------------- |
| ファイル読み込みエラー | レポートファイルが存在しない | ファイルパスを確認してエラーメッセージを表示 |
| JSON 解析エラー        | JSON 形式が不正              | 解析エラーの詳細を表示                       |
| レポート検証エラー     | レポート構造が不正           | 必須フィールドの欠落を指摘                   |
| ディレクトリ作成エラー | 出力ディレクトリ作成失敗     | パーミッションエラーを表示                   |
| ファイル書き込みエラー | ファイル書き込み失敗         | 書き込み先パスとエラー詳細を表示             |

## トラブルシューティング

このセクションでは、よくある問題とその解決方法を説明します。

### 1. レポートファイルが見つからない

**症状:**

```bash
Error: レポートファイルの読み込みに失敗: reports/communication/postman_analysis.json: no such file or directory
```

**原因:**

- レポートファイルのパスが間違っている
- postman-analyzer がまだ実行されていない

**解決方法:**

1. レポートファイルのパスを確認：

   ```bash
   ls -la tests/tools/postman-analyzer/reports/communication/
   ```

2. postman-analyzer を実行してレポートを生成：

   ```bash
   cd tests/tools/postman-analyzer
   go run . -collection collections/SaaSus\ Communication\ API\ Schema.postman_collection.json \
            -module communication
   ```

### 2. JSON 解析エラー

**症状:**

```bash
Error: JSONの解析に失敗: invalid character '}' looking for beginning of object key string
```

**原因:**

- レポートファイルの JSON 形式が不正
- ファイルが破損している

**解決方法:**

1. JSON ファイルの妥当性を確認：

   ```bash
   cat tests/tools/postman-analyzer/reports/communication/postman_analysis.json | jq .
   ```

2. postman-analyzer を再実行してレポートを再生成：

   ```bash
   cd tests/tools/postman-analyzer
   go run . -collection collections/SaaSus\ Communication\ API\ Schema.postman_collection.json \
            -module communication \
            -overwrite
   ```

### 3. 出力ディレクトリの作成失敗

**症状:**

```bash
Error: 出力ディレクトリの作成に失敗: tests/tools/test-param-generator/tests/e2e/communication: permission denied
```

**原因:**

- 出力先ディレクトリへの書き込み権限がない
- 親ディレクトリが存在しない

**解決方法:**

1. ディレクトリの権限を確認：

   ```bash
   ls -ld tests/tools/test-param-generator/tests/e2e/
   ```

2. 必要に応じて権限を変更：

   ```bash
   chmod 755 tests/tools/test-param-generator/tests/e2e/
   ```

3. 親ディレクトリを作成：

   ```bash
   mkdir -p tests/tools/test-param-generator/tests/e2e/communication
   ```

### 4. 生成された型定義ファイルが見つからない

**症状:**

```bash
Warning: 生成された型定義ファイルが見つかりません: generated/communicationapi/types.gen.go
```

**原因:**

- oapi-codegen でコード生成がまだ実行されていない
- モジュール名が間違っている
- ファイルパスの構造が想定と異なる

**解決方法:**

1. プロジェクトルートで型定義を生成：

   ```bash
   ./generate_all.sh
   ```

2. 特定のモジュールだけ生成：

   ```bash
   ./generate.sh communicationapi api/communicationapi.yml
   ```

3. 生成されたファイルを確認：

   ```bash
   ls -la generated/communicationapi/
   ```

4. デバッグモードで詳細なパス探索ログを確認：

   ```bash
   ./test-param-generator \
     -report tests/tools/postman-analyzer/reports/communication/postman_analysis.json \
     -module communication \
     -verbose
   ```

   出力例：
   ```
   Module情報からパスを構築: generated/communicationapi/types.gen.go
   ✓ Module情報からファイルを発見: /path/to/project/generated/communicationapi/types.gen.go
   ```

   または、フォールバック時：
   ```
   Module情報からのパスが見つからない、フォールバック探索を開始
   フォールバック: 以下のパターンを試行:
     - generated/communicationapi/types.gen.go
     - generated/communication/types.gen.go
     - ../generated/communicationapi/types.gen.go
   ✓ フォールバックでファイルを発見: /path/to/project/generated/communicationapi/types.gen.go
   ```

### 5. マッピング成功率が低い

**症状:**

```bash
=== マッピング統計 ===
グループ: 3/26 成功 (11.5%)
フィールド: 12/88 マッピング成功 (13.6%)
```

**原因:**

- types.gen.go の型名が標準パターンと異なる
- フィールド名の命名規則が特殊
- リソース名の表記が異なる

**解決方法:**

#### ステップ 1: 診断レポートを生成

```bash
./test-param-generator \
  -report tests/tools/postman-analyzer/reports/auth/postman_analysis.json \
  -module auth \
  -debug
```

診断レポート（`diagnostic_report_{timestamp}.txt`）が生成されます。

#### ステップ 2: 診断レポートを確認

レポートには以下の情報が含まれます：

- 失敗したマッピングの詳細
- 試行した型名パターン
- 試行したフィールドマッチング戦略
- 改善提案

例：

```
--- 失敗したマッピング ---

1. グループ: users, 操作: Create, フィールド: email
   試行したパターン:
     - CreateUserJSONRequestBody (存在しない)
     - CreateUserParams (存在しない)
   最も近かったマッチ: CreateSaasUserJSONRequestBody.Email (信頼度: 0.6)
   提案: リソース名を "saasusers" に変更するか、カスタムマッピングルールを追加してください
```

#### ステップ 3: カスタムマッピングルールを作成

1. サンプルファイルをコピー：

   ```bash
   cp tests/tools/test-param-generator/mapping_rules.json.example \
      tests/tools/test-param-generator/mapping_rules.json
   ```

2. 診断レポートの提案に基づいてルールを追加：

   ```json
   {
     "typeNameMappings": {
       "users.Create": "CreateSaasUserJSONRequestBody",
       "users.Update": "UpdateSaasUserJSONRequestBody"
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

3. ルールを適用して再実行：

   ```bash
   ./test-param-generator \
     -report tests/tools/postman-analyzer/reports/auth/postman_analysis.json \
     -module auth \
     -rules ./mapping_rules.json \
     -verbose
   ```

4. マッピング成功率を確認：

   ```bash
   === マッピング統計 ===
   グループ: 20/26 成功 (76.9%)
   フィールド: 70/88 マッピング成功 (79.5%)
   ```

### 6. 型名パターンが見つからない

**症状:**

```bash
Warning: 型名パターンが見つかりません: users.Create
試行したパターン: CreateUserJSONRequestBody, CreateUserParams, ...
```

**原因:**

- types.gen.go に該当する型が存在しない
- 型名の命名規則が特殊

**解決方法:**

1. types.gen.go で実際の型名を確認：

   ```bash
   grep "type.*struct" generated/authapi/types.gen.go | grep -i user | grep -i create
   ```

2. 見つかった型名をカスタムマッピングルールに追加：

   ```json
   {
     "typeNameMappings": {
       "users.Create": "実際の型名"
     }
   }
   ```

3. 型が存在しない場合は、OpenAPI 仕様を確認：

   ```bash
   cat api/authapi.yml | grep -A 10 "operationId.*createUser"
   ```

### 7. フィールド名がマッチしない

**症状:**

```bash
Warning: フィールド名がマッチしません: email
試行した戦略: Exact JSON Tag Match, Normalized Match, ...
```

**原因:**

- JSON タグとフィールド名が異なる
- フィールド名の命名規則が特殊

**解決方法:**

1. types.gen.go で実際のフィールド名を確認：

   ```bash
   grep -A 20 "type CreateSaasUserJSONRequestBody struct" generated/authapi/types.gen.go
   ```

2. 見つかったフィールド名をカスタムマッピングルールに追加：

   ```json
   {
     "fieldNameMappings": {
       "CreateSaasUserJSONRequestBody": {
         "email": "実際のフィールド名"
       }
     }
   }
   ```

### 8. ビルドエラー

**症状:**

```bash
Error: go build failed
```

**原因:**

- Go 環境が正しく設定されていない
- 依存パッケージが不足している

**解決方法:**

1. Go 環境を確認：

   ```bash
   go version
   go env
   ```

2. 依存パッケージをインストール：

   ```bash
   cd tests/tools/test-param-generator
   go mod tidy
   go mod download
   ```

3. 再ビルド：

   ```bash
   make clean
   make build
   ```

### 9. パフォーマンスが遅い

**症状:**

- 処理に 5 秒以上かかる
- メモリ使用量が多い

**原因:**

- types.gen.go ファイルが非常に大きい
- キャッシュが無効になっている

**解決方法:**

1. キャッシュが有効か確認（デフォルトで有効）

2. 不要な詳細ログを無効化：

   ```bash
   # -verboseオプションを外す
   ./test-param-generator -report <レポート> -module <モジュール>
   ```

3. 環境変数でキャッシュサイズを調整：

   ```bash
   export TYPE_CACHE_SIZE=1000
   ./test-param-generator -report <レポート> -module <モジュール>
   ```

### 10. 生成された Go コードがコンパイルエラー

**症状:**

```bash
# 生成されたヘルパー関数でコンパイルエラー
./communication_test_helpers.go:15:2: undefined: stringPtr
```

**原因:**

- 生成されたコードに不足している関数がある
- パッケージ名が間違っている

**解決方法:**

1. 生成されたファイルを確認：

   ```bash
   cat tests/tools/test-param-generator/tests/e2e/communication/communication_test_helpers.go
   ```

2. 不足している関数を追加：

   ```go
   // stringPtr は文字列のポインタを返すヘルパー関数
   func stringPtr(s string) *string {
       return &s
   }
   ```

3. パッケージ名を確認して修正

### よくある質問（FAQ）

#### Q1: カスタムマッピングルールはどこに配置すべきですか？

A: 以下の優先順位で検索されます：

1. `-rules`オプションで指定したパス
2. カレントディレクトリの`mapping_rules.json`
3. 親ディレクトリの`mapping_rules.json`
4. `tests/tools/test-param-generator/mapping_rules.json`

プロジェクト全体で共有する場合は、4 番目の場所に配置してバージョン管理に含めることを推奨します。

#### Q2: マッピング成功率の目標値は？

A: 80%以上を目標としています。カスタムマッピングルールを適切に設定することで、ほとんどのモジュールで 80-90%の成功率を達成できます。

#### Q3: 診断レポートはどこに出力されますか？

A: `-debug`オプションを使用すると、カレントディレクトリに`diagnostic_report_{timestamp}.txt`が生成されます。

#### Q4: 既存の test_params.json を上書きしたくない場合は？

A: `-overwrite`オプションを指定しない場合、確認プロンプトが表示されます。または、`-output`オプションで別のディレクトリを指定してください。

#### Q5: 複数のモジュールを一括処理できますか？

A: 現在は 1 モジュールずつ処理する必要があります。シェルスクリプトでループ処理することで一括実行できます：

```bash
for module in auth billing pricing communication; do
  ./test-param-generator \
    -report tests/tools/postman-analyzer/reports/${module}/postman_analysis.json \
    -module ${module} \
    -overwrite
done
```

### サポート

上記の方法で解決しない場合は、以下の情報を含めて Issue を作成してください：

1. 実行したコマンド
2. エラーメッセージの全文
3. 診断レポート（`-debug`オプションで生成）
4. 使用している Go のバージョン（`go version`）
5. 対象のモジュール名

## テストでの使用例

生成されたヘルパー関数をテストで使用する例：

```go
func TestCreateFeedback(t *testing.T) {
	// テストパラメータを読み込み
	params := loadTestParams(t)

	// クライアントの初期化
	client := setupClient(t)

	// フィードバックを作成
	ctx := context.Background()
	result, err := client.CreateFeedback(ctx, communicationapi.CreateFeedbackJSONRequestBody{
		UserId:              stringPtr(params.Feedbacks.CreateParams.UserId),
		FeedbackTitle:       stringPtr(params.Feedbacks.CreateParams.FeedbackTitle),
		FeedbackDescription: stringPtr(params.Feedbacks.CreateParams.FeedbackDescription),
	})

	if err != nil {
		t.Fatalf("フィードバックの作成に失敗: %v", err)
	}

	// 検証
	if result.StatusCode != 201 {
		t.Errorf("期待されるステータスコード: 201, 実際: %d", result.StatusCode)
	}
}
```

## 関連ドキュメント

- **ドキュメント入口**: [docs/test-tools/README.md](../../docs/test-tools/README.md)
- **エラーハンドリング**: [ERROR_HANDLING.md](ERROR_HANDLING.md)
- **改善履歴メモ**: [IMPROVEMENTS.md](IMPROVEMENTS.md)

## 関連ツール

- **postman-analyzer**: Postman コレクションを解析してレポートを生成
- **oapi-codegen**: OpenAPI 仕様から Go クライアントコードを生成

## ライセンス

このツールはプロジェクトのライセンスに従います。

## サポート

問題が発生した場合は、以下を確認してください：

1. postman-analyzer が正常に実行され、レポートが生成されているか
2. `./generate_all.sh`を実行して、API 型定義が最新か
3. `-verbose`オプションで詳細ログを確認

それでも解決しない場合は、プロジェクトの Issue トラッカーに報告してください。
