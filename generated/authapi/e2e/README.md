# SaaSus Auth API E2Eテストスイート

このディレクトリには、SaaSus Auth APIの包括的なEnd-to-End（E2E）テストが含まれています。

## 概要

E2Eテストスイートは、SaaSus Auth APIの全機能を実際のAPIエンドポイントに対してテストし、統合的な動作を検証します。全10のストーリーテストで構成され、基本機能から高度な統合シナリオまでを網羅的にテストします。

## テスト構成

### ストーリーテスト（全10ストーリー）

1. **01_basic_setup_test.go** - 基本セットアップテスト
   - 基本情報管理（ドメイン名、メールアドレス設定）
   - 認証情報管理（コールバックURL設定）
   - 環境管理（開発・本番環境の作成・管理）

2. **02_user_management_test.go** - ユーザー管理テスト
   - ユーザー属性管理（カスタム属性の作成・削除）
   - SaaSユーザー管理（ユーザーの作成・更新・削除）
   - ユーザー情報管理（メール・パスワード更新）

3. **03_role_management_test.go** - ロール管理テスト
   - ロール作成・削除
   - ロール一覧取得
   - ロール権限管理

4. **04_tenant_management_test.go** - テナント管理テスト
   - テナント属性管理（カスタム属性の作成・削除）
   - テナント管理（テナントの作成・更新・削除）
   - 請求情報管理（請求先情報の設定・更新）

5. **05_tenant_user_management_test.go** - テナントユーザー管理テスト
   - テナントユーザー管理（ユーザーのテナント追加・削除）
   - 全テナントユーザー管理（テナント横断ユーザー管理）
   - ロール管理（ユーザーへのロール割り当て・削除）

6. **06_invitation_management_test.go** - 招待管理テスト
   - 招待管理（招待の作成・削除・一覧取得）
   - 招待検証（招待の有効性確認・コード検証）

7. **07_authentication_flow_test.go** - 認証フローテスト
   - 認証・認可情報の保存・取得
   - 一時コード認証フロー（10秒有効期限）
   - リフレッシュトークン認証フロー
   - エラーハンドリング（無効トークン、期限切れ等）

8. **08_single_tenant_management_test.go** - シングルテナント管理テスト
   - シングルテナント設定管理（カスタマイズ設定）
   - CloudFormationテンプレート管理
   - DDLテンプレート管理
   - AWS連携設定（IAMロール、ExternalID）

9. **09_error_handling_test.go** - エラーハンドリングテスト
   - HTTPエラーステータステスト（400, 401, 403, 404, 500等）
   - エラーレスポンス形式の統一性確認
   - 認証・認可エラーテスト
   - バリデーションエラーテスト
   - リソース不存在エラーテスト

10. **10_integration_test.go** - 統合テスト
    - 完全なSaaSセットアップフロー（基本設定からユーザー招待まで）
    - エンドツーエンドシナリオ（ユーザージャーニー、管理者ワークフロー）
    - パフォーマンステスト（負荷テスト、ストレステスト）
    - 障害復旧テスト（データ整合性、トランザクション整合性）

### 共通ライブラリ

- **common/client.go** - 認証付きクライアント管理
- **common/assertions.go** - テストアサーション（レスポンス時間、ステータスコード、データ検証）
- **common/cleanup.go** - リソースクリーンアップ（自動追跡・削除）
- **common/types.go** - テストデータ型定義（全ストーリー対応）
- **common/testdata.go** - テストデータローダー（YAML読み込み）

### テストデータ

各ストーリーに対応するYAMLファイルでテストデータを管理：

- `testdata/basic_setup/test_data.yml` - 基本セットアップ用データ
- `testdata/user_management/test_data.yml` - ユーザー管理用データ
- `testdata/role_management/test_data.yml` - ロール管理用データ
- `testdata/tenant_management/test_data.yml` - テナント管理用データ
- `testdata/tenant_user_management/test_data.yml` - テナントユーザー管理用データ
- `testdata/invitation_management/test_data.yml` - 招待管理用データ
- `testdata/authentication_flow/test_data.yml` - 認証フロー用データ
- `testdata/single_tenant_management/test_data.yml` - シングルテナント管理用データ
- `testdata/error_handling/test_data.yml` - エラーハンドリング用データ
- `testdata/integration_test/test_data.yml` - 統合テスト用データ

## 実行方法

### 前提条件

以下の環境変数を設定してください：

```bash
export SAASUS_SAAS_ID="your-saas-id"
export SAASUS_API_KEY="your-api-key"
export SAASUS_SECRET_KEY="your-secret-key"
```

### 全テスト実行

```bash
# E2Eテストディレクトリに移動
cd generated/authapi/e2e

# 全ストーリーテスト実行（推奨）
go test -v ./...

# 詳細ログ付きで実行
go test -v -verbose ./...

# タイムアウト延長（大規模テスト用）
go test -v -timeout 30m ./...
```

### 個別ストーリーテスト実行

```bash
# 基本セットアップテストのみ実行
go test -v ./stories -run TestStory01_BasicSetup

# 認証フローテストのみ実行
go test -v ./stories -run TestStory07_AuthenticationFlow

# 統合テストのみ実行
go test -v ./stories -run TestStory10_IntegrationTest

# エラーハンドリングテストのみ実行
go test -v ./stories -run TestStory09_ErrorHandling
```

### 特定カテゴリのテスト実行

```bash
# 基本機能テスト（ストーリー1-6）
go test -v ./stories -run "TestStory0[1-6]"

# 高度機能テスト（ストーリー7-10）
go test -v ./stories -run "TestStory(07|08|09|10)"

# パフォーマンステストのみ
go test -v ./stories -run "Performance"
```

### ヘルスチェック

```bash
# API接続確認
go test -v -run TestHealthCheck

# 環境確認
go test -v -run TestE2EStoryExecution -short
```

## テスト設計

### テストパターン

1. **正常系テスト**: 期待される動作の確認
2. **異常系テスト**: エラーハンドリングの確認
3. **境界値テスト**: 制限値での動作確認
4. **統合テスト**: 複数API間の連携確認
5. **パフォーマンステスト**: 負荷・レスポンス時間確認
6. **セキュリティテスト**: 認証・認可の確認

### リソース管理

- テスト実行時に作成されたリソースは自動的に追跡
- ストーリー別・リソース種別でのクリーンアップ
- テスト終了時に自動クリーンアップ
- 失敗時のリソース残存を防止
- 並行実行時の競合回避

### アサーション

- **レスポンス時間の検証**: API応答性能の確認
- **ステータスコードの検証**: HTTP応答の正確性確認
- **レスポンスボディの検証**: データ形式・内容の確認
- **データ整合性の検証**: 作成・更新・削除の整合性確認
- **エラーメッセージの検証**: エラー応答の適切性確認

## 設定

### タイムアウト設定

- **デフォルトタイムアウト**: 30秒
- **長時間処理**: 2-3分（認証フロー、統合テスト）
- **クリーンアップ**: 5分
- **統合テスト**: 10分
- **パフォーマンステスト**: 5分

### リトライ設定

- **デフォルトリトライ回数**: 3回
- **リトライ間隔**: 1秒
- **指数バックオフ**: 対応

### パフォーマンス設定

- **同時実行数**: 10-20（負荷テスト）
- **テスト継続時間**: 60秒（負荷テスト）
- **成功率閾値**: 95%以上

## 高度な機能

### 認証フローテスト

- **一時コード認証**: 10秒有効期限の確認
- **リフレッシュトークン認証**: トークン更新フローの確認
- **認証エラーハンドリング**: 無効トークン・期限切れの処理確認

### シングルテナント管理

- **カスタマイズ設定**: ページタイトル、CSS、ロゴ等の設定
- **AWS連携**: CloudFormation、DDLテンプレート、IAMロール設定
- **テンプレート管理**: インフラ構築用テンプレートの取得・検証

### エラーハンドリング

- **包括的エラーテスト**: 全HTTPステータスコードの確認
- **エラーレスポンス統一性**: 一貫したエラー形式の確認
- **多言語エラーメッセージ**: 日本語エラーメッセージの確認

### 統合テスト

- **完全セットアップフロー**: SaaS環境の完全構築テスト
- **マルチテナント分離**: テナント間データ分離の確認
- **パフォーマンス統合**: 負荷状況下での全機能動作確認

## トラブルシューティング

### よくある問題

1. **認証エラー**
   ```bash
   # 環境変数の確認
   echo $SAASUS_SAAS_ID
   echo $SAASUS_API_KEY
   echo $SAASUS_SECRET_KEY
   
   # APIキーの有効性確認
   go test -v -run TestHealthCheck
   ```

2. **タイムアウトエラー**
   ```bash
   # ネットワーク接続確認
   curl -I https://api.saasus.io/v1/auth/health
   
   # タイムアウト延長
   go test -v -timeout 30m ./...
   ```

3. **リソース競合**
   ```bash
   # 順次実行（並行実行を避ける）
   go test -v -p 1 ./...
   
   # クリーンアップ強制実行
   go test -v -run TestHealthCheck  # クリーンアップのみ
   ```

4. **認証フローテストの失敗**
   ```bash
   # 一時コード有効期限の確認（10秒制限）
   go test -v ./stories -run "TempCode"
   
   # リフレッシュトークンの確認
   go test -v ./stories -run "RefreshToken"
   ```

### ログ確認

```bash
# 詳細ログ出力
go test -v -verbose ./...

# 特定ストーリーの詳細ログ
go test -v -run TestStory07_AuthenticationFlow -verbose

# エラーのみ表示
go test ./... 2>&1 | grep -E "(FAIL|ERROR)"

# パフォーマンス情報表示
go test -v ./stories -run "Performance" -verbose
```

### デバッグモード

```bash
# デバッグ情報付きで実行
SAASUS_DEBUG=true go test -v ./...

# 特定のAPIコールをトレース
SAASUS_TRACE=true go test -v -run TestStory01_BasicSetup
```

## 開発者向け情報

### 新しいテストの追加

1. **ストーリーテストファイルの作成**
   ```bash
   # 新しいストーリーテストファイル
   touch stories/11_new_feature_test.go
   ```

2. **テストデータの作成**
   ```bash
   # テストデータディレクトリ作成
   mkdir testdata/new_feature
   touch testdata/new_feature/test_data.yml
   ```

3. **型定義の追加**
   ```go
   // common/types.go に新しい型を追加
   type NewFeatureTestData struct {
       // フィールド定義
   }
   ```

4. **データローダーの追加**
   ```go
   // common/testdata.go にローダーメソッドを追加
   func (loader *TestDataLoader) LoadNewFeatureData() (*NewFeatureTestData, error) {
       // 実装
   }
   ```

### テストデータの管理

- **構造化データ**: YAMLファイルで階層的にデータを管理
- **環境固有データ**: `test_data_dev.yml`, `test_data_prod.yml` で環境別管理
- **テストケース定義**: 期待値を明確に定義
- **パラメータ化**: 複数パターンのテストデータを効率的に管理

### CI/CD統合

```yaml
# GitHub Actions例
name: E2E Tests
on: [push, pull_request]
jobs:
  e2e:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
        with:
          go-version: 1.21
      - name: Run E2E Tests
        env:
          SAASUS_SAAS_ID: ${{ secrets.SAASUS_SAAS_ID }}
          SAASUS_API_KEY: ${{ secrets.SAASUS_API_KEY }}
          SAASUS_SECRET_KEY: ${{ secrets.SAASUS_SECRET_KEY }}
        run: |
          cd generated/authapi/e2e
          go test -v -timeout 30m ./...
```

### パフォーマンス最適化

```bash
# 並行実行数の調整
go test -v -p 4 ./...

# メモリ使用量の監視
go test -v -memprofile=mem.prof ./...

# CPU使用量の監視
go test -v -cpuprofile=cpu.prof ./...
```

## ベストプラクティス

### テスト設計

1. **独立性**: 各テストは他のテストに依存しない
2. **冪等性**: 同じテストを何度実行しても同じ結果
3. **クリーンアップ**: リソースの適切な削除
4. **エラーハンドリング**: 予期しないエラーへの対応
5. **ドキュメント**: テストの目的と期待値を明確に記述

### データ管理

1. **外部化**: テストデータをコードから分離
2. **バージョン管理**: テストデータの変更履歴を管理
3. **環境対応**: 開発・テスト・本番環境への対応
4. **セキュリティ**: 機密情報の適切な管理

### 実行戦略

1. **段階的実行**: 基本機能から高度機能へ
2. **並行実行**: 独立したテストの並行実行
3. **失敗時対応**: 失敗時の詳細情報収集
4. **継続的実行**: CI/CDパイプラインでの自動実行

## 貢献

テストの改善や新機能のテスト追加は以下の手順で行ってください：

1. **既存パターンの踏襲**: 既存のテストパターンに従う
2. **適切なクリーンアップ**: リソースの確実な削除を実装
3. **テストデータ外部化**: YAMLファイルでのデータ管理
4. **ドキュメント更新**: README.mdの更新
5. **レビュー**: コードレビューでの品質確保

### コントリビューションガイドライン

- テストの目的を明確に記述
- エラーケースを含む包括的なテスト
- パフォーマンスへの配慮
- セキュリティベストプラクティスの遵守

## ライセンス

このテストスイートは、SaaSus SDKと同じライセンスの下で提供されます。

---

**注意**: このE2Eテストスイートは実際のSaaSus Auth APIに対してテストを実行します。テスト実行時は適切な環境設定と認証情報の管理を行ってください。