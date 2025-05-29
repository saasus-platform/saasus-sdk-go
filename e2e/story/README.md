# SaaSus Auth API E2Eテストストーリー

## 概要
このディレクトリには、SaaSus Auth APIの包括的なE2Eテストストーリーが含まれています。各ストーリーは特定のAPI群をテストし、チェックリスト形式でテストの進捗を管理できます。

## ストーリー一覧

### 1. [基本セットアップストーリー](./01_basic_setup.md)
**対象API**: basicInfo.yml, authInfo.yml, env.yml
- 基本設定情報の管理
- 認証情報の設定
- 環境情報の管理
- カスタマイズページ設定
- 通知メールテンプレート設定

### 2. [ユーザー管理ストーリー](./02_user_management.md)
**対象API**: saasUser.yml, userAttribute.yml, userInfo.yml
- ユーザー属性定義の管理
- SaaSユーザーの作成・管理
- 外部アカウント連携
- AWS Marketplace連携
- MFA設定
- パスワード管理

### 3. [役割(ロール)管理ストーリー](./03_role_management.md)
**対象API**: role.yml
- 役割(ロール)の作成・削除
- 役割ベースの認可準備

### 4. [テナント管理ストーリー](./04_tenant_management.md)
**対象API**: tenant.yml, tenantAttribute.yml
- テナント属性定義の管理
- テナントの基本管理
- 請求情報の管理
- プラン管理
- 外部IDプロバイダ設定
- Stripe連携管理

### 5. [テナントユーザー管理ストーリー](./05_tenant_user_management.md)
**対象API**: tenantUser.yml
- テナントユーザーの基本管理
- 全テナントユーザーの管理
- 役割の付与・削除
- 環境毎の権限設定

### 6. [招待管理ストーリー](./06_invitation_management.md)
**対象API**: invitation.yml
- テナント招待の作成・管理
- 招待の有効性確認
- 新規・既存ユーザーの招待フロー
- 招待の期限管理

### 7. [認証フローストーリー](./07_authentication_flow.md)
**対象API**: credential.yml
- 認証・認可情報の保存・取得
- 一時コード認証フロー
- リフレッシュトークン認証フロー
- 認証情報の有効期限管理

### 8. [シングルテナント管理ストーリー](./08_single_tenant_management.md)
**対象API**: singleTenant.yml
- シングルテナント設定の管理
- CloudFormationテンプレート管理
- AWS連携設定
- DDLテンプレート管理

### 9. [エラーハンドリングストーリー](./09_error_handling.md)
**対象API**: error.yml
- テスト用エラーエンドポイント
- エラーレスポンス形式の確認
- 各種エラーケースのテスト

### 10. [統合テストストーリー](./10_integration_test.md)
**対象API**: 全API
- 完全なSaaSセットアップフロー
- エンドツーエンドシナリオ
- パフォーマンステスト
- 障害復旧テスト

## テスト実行順序

### 推奨実行順序
1. **01_basic_setup.md** - 基本設定（必須）
2. **03_role_management.md** - 役割定義（ユーザー管理の前に必要）
3. **02_user_management.md** - ユーザー管理
4. **04_tenant_management.md** - テナント管理
5. **05_tenant_user_management.md** - テナントユーザー管理
6. **06_invitation_management.md** - 招待機能
7. **07_authentication_flow.md** - 認証フロー
8. **08_single_tenant_management.md** - シングルテナント機能（オプション）
9. **09_error_handling.md** - エラーハンドリング
10. **10_integration_test.md** - 統合テスト（最終確認）

### 並行実行可能
- 基本設定完了後、以下は並行実行可能：
  - ユーザー管理とテナント管理
  - 認証フローとエラーハンドリング
  - シングルテナント管理（独立機能）

## チェックリストの使用方法

各ストーリーファイル内のチェックボックス `- [ ]` を使用してテストの進捗を管理します：

```markdown
- [ ] **GET /basic-info** - 基本設定情報を取得
- [x] **PUT /basic-info** - 基本設定情報を更新 ✅
```

### チェック状態
- `- [ ]` : 未実行
- `- [x]` : 実行完了
- `- [x] ✅` : 実行完了（成功）
- `- [x] ❌` : 実行完了（失敗）

## API網羅状況

### 対象APIファイル（14ファイル）
- ✅ authInfo.yml (4 APIs)
- ✅ basicInfo.yml (8 APIs) 
- ✅ credential.yml (2 APIs)
- ✅ env.yml (5 APIs)
- ✅ error.yml (1 API)
- ✅ invitation.yml (6 APIs)
- ✅ role.yml (3 APIs)
- ✅ saasUser.yml (22 APIs)
- ✅ singleTenant.yml (3 APIs)
- ✅ tenant.yml (12 APIs)
- ✅ tenantAttribute.yml (3 APIs)
- ✅ tenantUser.yml (8 APIs)
- ✅ userAttribute.yml (4 APIs)
- ✅ userInfo.yml (2 APIs)

**総API数**: 81 APIs

## 前提条件

### 環境要件
- SaaSus Platformアカウント
- APIキーの発行
- テスト用ドメインの準備
- AWS環境（シングルテナント機能使用時）
- Stripe環境（課金機能使用時）

### テストデータ
- テスト用メールアドレス
- テスト用ユーザー情報
- CloudFormationテンプレート（シングルテナント用）
- DDLテンプレート（シングルテナント用）

## 注意事項

### セキュリティ
- 本番環境のデータは使用しない
- テスト用の認証情報を使用する
- テスト後は不要なデータを削除する

### 課金
- AWS Marketplace連携テスト時は課金に注意
- Stripe連携テスト時は課金に注意
- テスト環境での実行を推奨

### データ管理
- テストデータの作成と削除を適切に行う
- 他のテストに影響しないよう注意する
- テスト実行順序を守る

## トラブルシューティング

### よくある問題
1. **認証エラー**: APIキーの確認、有効期限の確認
2. **権限エラー**: 適切な権限でのAPIアクセス確認
3. **データ不整合**: テスト実行順序の確認
4. **外部サービス連携エラー**: 外部サービスの設定確認

### ログ確認
- APIレスポンスの詳細確認
- エラーメッセージの分析
- ネットワーク接続の確認

## 貢献

### ストーリーの追加・修正
1. 新しいAPIが追加された場合は対応するストーリーを作成
2. 既存ストーリーの改善提案
3. テストケースの追加

### レポート
- バグ発見時の詳細レポート
- パフォーマンス問題の報告
- 改善提案