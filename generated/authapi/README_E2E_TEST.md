# SaaSus Auth API E2Eテスト

このディレクトリには、SaaSus Auth APIのEnd-to-Endテストが含まれています。

## 概要

このE2Eテストは、実際のSaaSus Auth APIサーバー（https://api.dev.saasus.io/v1/auth）に対してテストを実行し、APIの動作を検証します。

## テスト対象

以下のClientWithResponsesInterfaceで定義されている全てのメソッドをテストします：

### 基本情報・認証情報
- `GetAuthInfoWithResponse` - 認証情報取得
- `UpdateAuthInfoWithResponse` - 認証情報更新
- `GetBasicInfoWithResponse` - 基本設定情報取得
- `UpdateBasicInfoWithResponse` - 基本設定情報更新

### SaaSユーザー管理
- `GetSaasUsersWithResponse` - SaaSユーザー一覧取得
- `CreateSaasUserWithResponse` - SaaSユーザー作成
- `GetSaasUserWithResponse` - SaaSユーザー詳細取得
- `DeleteSaasUserWithResponse` - SaaSユーザー削除
- `UpdateSaasUserPasswordWithResponse` - パスワード変更
- `UpdateSaasUserEmailWithResponse` - メールアドレス変更
- `UpdateSaasUserAttributesWithResponse` - ユーザー属性更新

### テナント管理
- `GetTenantsWithResponse` - テナント一覧取得
- `CreateTenantWithResponse` - テナント作成
- `GetTenantWithResponse` - テナント詳細取得
- `UpdateTenantWithResponse` - テナント更新
- `DeleteTenantWithResponse` - テナント削除

### テナントユーザー管理
- `GetAllTenantUsersWithResponse` - 全テナントユーザー取得
- `GetAllTenantUserWithResponse` - 全テナントユーザー詳細取得
- `GetTenantUsersWithResponse` - テナントユーザー一覧取得
- `CreateTenantUserWithResponse` - テナントユーザー作成
- `GetTenantUserWithResponse` - テナントユーザー詳細取得
- `UpdateTenantUserWithResponse` - テナントユーザー更新
- `DeleteTenantUserWithResponse` - テナントユーザー削除

### 役割・権限管理
- `GetRolesWithResponse` - 役割一覧取得
- `CreateRoleWithResponse` - 役割作成
- `DeleteRoleWithResponse` - 役割削除
- `CreateTenantUserRolesWithResponse` - テナントユーザー役割作成
- `DeleteTenantUserRoleWithResponse` - テナントユーザー役割削除

### 属性管理
- `GetUserAttributesWithResponse` - ユーザー属性一覧取得
- `CreateUserAttributeWithResponse` - ユーザー属性作成
- `DeleteUserAttributeWithResponse` - ユーザー属性削除
- `GetTenantAttributesWithResponse` - テナント属性一覧取得
- `CreateTenantAttributeWithResponse` - テナント属性作成
- `DeleteTenantAttributeWithResponse` - テナント属性削除
- `CreateSaasUserAttributeWithResponse` - SaaSユーザー属性作成

### 環境管理
- `GetEnvsWithResponse` - 環境一覧取得
- `CreateEnvWithResponse` - 環境作成
- `GetEnvWithResponse` - 環境詳細取得
- `UpdateEnvWithResponse` - 環境更新
- `DeleteEnvWithResponse` - 環境削除

### 認証・認可設定
- `GetAuthCredentialsWithResponse` - 認証情報取得
- `CreateAuthCredentialsWithResponse` - 認証情報作成
- `GetSignInSettingsWithResponse` - サインイン設定取得
- `UpdateSignInSettingsWithResponse` - サインイン設定更新
- `GetIdentityProvidersWithResponse` - 外部IDプロバイダー取得
- `UpdateIdentityProviderWithResponse` - 外部IDプロバイダー更新

### ページカスタマイズ
- `GetCustomizePagesWithResponse` - カスタマイズページ取得
- `UpdateCustomizePagesWithResponse` - カスタマイズページ更新
- `GetCustomizePageSettingsWithResponse` - ページ設定取得
- `UpdateCustomizePageSettingsWithResponse` - ページ設定更新

### 通知・メッセージ
- `FindNotificationMessagesWithResponse` - 通知メッセージ取得
- `UpdateNotificationMessagesWithResponse` - 通知メッセージ更新

### 招待機能
- `GetTenantInvitationsWithResponse` - テナント招待一覧取得
- `CreateTenantInvitationWithResponse` - テナント招待作成
- `GetTenantInvitationWithResponse` - テナント招待詳細取得
- `DeleteTenantInvitationWithResponse` - テナント招待削除
- `GetInvitationValidityWithResponse` - 招待有効性確認
- `ValidateInvitationWithResponse` - 招待検証

### サインアップ
- `SignUpWithResponse` - サインアップ
- `ResendSignUpConfirmationEmailWithResponse` - 確認メール再送信

### シングルテナント機能
- `GetSingleTenantSettingsWithResponse` - シングルテナント設定取得
- `UpdateSingleTenantSettingsWithResponse` - シングルテナント設定更新
- `GetCloudFormationLaunchStackLinkForSingleTenantWithResponse` - CloudFormationリンク取得

### Stripe連携
- `CreateTenantAndPricingWithResponse` - Stripe初期設定
- `DeleteStripeTenantAndPricingWithResponse` - Stripe情報削除
- `ResetPlanWithResponse` - プランリセット
- `GetStripeCustomerWithResponse` - Stripe顧客情報取得

### その他
- `LinkAwsMarketplaceWithResponse` - AWS Marketplace連携
- `SignUpWithAwsMarketplaceWithResponse` - AWS Marketplaceサインアップ
- `ConfirmSignUpWithAwsMarketplaceWithResponse` - AWS Marketplaceサインアップ確認
- `RequestExternalUserLinkWithResponse` - 外部ユーザー連携要求
- `ConfirmExternalUserLinkWithResponse` - 外部ユーザー連携確認

**注意**: `ReturnInternalServerError`はテスト用途のエラーエンドポイントのため、テストコードは生成していません。

## 前提条件

### 環境変数の設定

テストを実行する前に、以下の環境変数を設定する必要があります：

```bash
export SAASUS_SAAS_ID="your-saas-id"
export SAASUS_API_KEY="your-api-key"
export SAASUS_SECRET_KEY="your-secret-key"
```

これらの環境変数が設定されていない場合、テストは自動的にスキップされます。

### 依存関係

- Go 1.19以上
- SaaSus SDK for Go
- インターネット接続（実際のAPIサーバーにアクセスするため）

## テストファイル

- `client_e2e_test.go` - E2Eテストの実装
- `test_params.json` - テストパラメータの設定ファイル
- `README_E2E_TEST.md` - このファイル

## テストパラメータ設定

`test_params.json`ファイルには、各テストで使用するパラメータが定義されています：

```json
{
  "getUserInfo": {
    "testCases": [...],
    "edgeCases": [...]
  },
  "basicInfo": {
    "updateParams": {
      "domainName": "test.example.com",
      "fromEmailAddress": "noreply@test.example.com"
    }
  },
  "saasUsers": {
    "createParams": {
      "email": "newuser@example.com",
      "password": "TempPassword123!"
    }
  },
  ...
}
```

必要に応じて、このファイルを編集してテストパラメータをカスタマイズできます。

## テスト実行方法

### 全てのE2Eテストを実行

```bash
cd generated/authapi
go test -v -run "TestE2E_"
```

### 特定のテストのみ実行

```bash
# 基本情報関連のテストのみ実行
go test -v -run "TestE2E_.*BasicInfo.*"

# SaaSユーザー関連のテストのみ実行
go test -v -run "TestE2E_.*SaasUser.*"

# テナント関連のテストのみ実行
go test -v -run "TestE2E_.*Tenant.*"

# パフォーマンステストのみ実行
go test -v -run "TestE2E_Performance.*"
```

### 並列実行

```bash
# 並列実行（注意：APIレート制限に注意）
go test -v -parallel 2 -run "TestE2E_"
```

## テスト結果の解釈

### 成功例

```
=== RUN   TestE2E_GetBasicInfoWithResponse
    client_e2e_test.go:XXX: ステータスコード: 200
    client_e2e_test.go:XXX: 基本情報: DomainName=example.com
    client_e2e_test.go:XXX: DNS検証済み: true
--- PASS: TestE2E_GetBasicInfoWithResponse (0.50s)
```

### スキップ例

```
=== RUN   TestE2E_GetSaasUsersWithResponse
    client_e2e_test.go:XXX: E2Eテストをスキップ: SAASUS_SAAS_ID, SAASUS_API_KEY, SAASUS_SECRET_KEY 環境変数が設定されていません
--- SKIP: TestE2E_GetSaasUsersWithResponse (0.00s)
```

### エラー例

```
=== RUN   TestE2E_CreateSaasUserWithResponse
    client_e2e_test.go:XXX: ステータスコード: 400
    client_e2e_test.go:XXX: バリデーションエラー: Type=validation_error, Message=Email already exists
--- PASS: TestE2E_CreateSaasUserWithResponse (0.30s)
```

## パフォーマンステスト

パフォーマンステストでは以下を測定します：

- **レスポンス時間**: 各APIの応答時間
- **スループット**: 単位時間あたりの処理件数
- **同時実行**: 複数リクエストの並列処理性能
- **成功率**: 同時実行時の成功率

### パフォーマンス基準

- 基本的なGETリクエスト: 5秒以内
- 軽負荷テスト: 2秒以内
- 中負荷テスト: 3秒以内
- 高負荷テスト: 5秒以内
- 同時実行成功率: 80%以上

## トラブルシューティング

### 認証エラー

```
401 Unauthorized
```

- 環境変数が正しく設定されているか確認
- APIキーとシークレットキーが有効か確認
- SaaS IDが正しいか確認

### ネットワークエラー

```
dial tcp: lookup api.dev.saasus.io: no such host
```

- インターネット接続を確認
- ファイアウォール設定を確認
- プロキシ設定を確認

### レート制限エラー

```
429 Too Many Requests
```

- テストの実行間隔を空ける
- 並列実行数を減らす
- APIレート制限の緩和を検討

### データ依存エラー

一部のテストは既存のデータに依存します：

- ユーザー一覧が空の場合、個別ユーザー取得テストはスキップされます
- テナント一覧が空の場合、個別テナント取得テストはスキップされます

## 注意事項

1. **実際のAPIサーバーを使用**: このテストは実際のSaaSus APIサーバーに対して実行されるため、テストデータが実際に作成・変更・削除されます。

2. **テストデータの管理**: テストで作成されたデータは、可能な限りテスト内で削除されますが、エラーが発生した場合は手動で削除が必要な場合があります。

3. **APIレート制限**: 大量のテストを短時間で実行すると、APIレート制限に達する可能性があります。

4. **環境分離**: 本番環境ではなく、開発・テスト環境で実行することを強く推奨します。

5. **データの一意性**: テストで作成するデータ（メールアドレス、テナント名など）は、タイムスタンプを使用して一意性を保っています。

## 継続的インテグレーション

CI/CDパイプラインでこれらのテストを実行する場合：

```yaml
# GitHub Actions例
- name: Run E2E Tests
  env:
    SAASUS_SAAS_ID: ${{ secrets.SAASUS_SAAS_ID }}
    SAASUS_API_KEY: ${{ secrets.SAASUS_API_KEY }}
    SAASUS_SECRET_KEY: ${{ secrets.SAASUS_SECRET_KEY }}
  run: |
    cd generated/authapi
    go test -v -run "TestE2E_" -timeout 10m
```

## 貢献

テストの改善や新しいテストケースの追加は歓迎します：

1. `test_params.json`でテストパラメータを調整
2. 新しいエッジケースやバリエーションを追加
3. パフォーマンステストの基準を調整
4. エラーハンドリングの改善

## ライセンス

このテストコードは、SaaSus SDK for Goと同じライセンスの下で提供されます。