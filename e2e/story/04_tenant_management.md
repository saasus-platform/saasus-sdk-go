# 04. テナント管理ストーリー

## 概要
テナントの作成、管理、設定を行うE2Eテストストーリー

## テスト対象API
- tenant.yml
- tenantAttribute.yml

## ストーリー

### 1. テナント属性定義の管理
- [ ] **GET /tenant-attributes** - テナント属性一覧を取得
  - 既存のテナント属性定義の確認
- [ ] **POST /tenant-attributes** - テナント属性を作成
  - 会社規模、業界、設立年などの属性定義を作成
- [ ] **DELETE /tenant-attributes/{attribute_name}** - テナント属性を削除
  - 不要な属性定義の削除

### 2. テナントの基本管理
- [ ] **GET /tenants** - テナント一覧を取得
  - 全テナントの一覧取得
  - 請求情報、プラン情報の確認
- [ ] **POST /tenants** - テナントを作成
  - 新規テナントの作成
  - 属性情報、事務管理部門メールアドレスの設定
- [ ] **GET /tenants/{tenant_id}** - テナント詳細を取得
  - 特定テナントの詳細情報取得
  - プラン履歴、現在のプラン期間の確認
- [ ] **PATCH /tenants/{tenant_id}** - テナント情報を更新
  - テナント名、属性情報の更新
- [ ] **DELETE /tenants/{tenant_id}** - テナントを削除
  - テナントの完全削除

### 3. テナント請求情報の管理
- [ ] **PUT /tenants/{tenant_id}/billing-info** - 請求情報を更新
  - 請求先住所、請求書言語の設定
  - 請求用テナント名の設定

### 4. テナントプラン管理
- [ ] **PUT /tenants/{tenant_id}/plans** - プラン情報を更新
  - 次回プランID、税率IDの設定
  - プラン変更時の比例配分設定
  - 従量課金アイテム削除設定

### 5. テナント外部IDプロバイダ設定
- [ ] **GET /tenants/{tenant_id}/identity-providers** - 外部IDプロバイダ取得
  - テナント毎のSAML設定確認
- [ ] **PUT /tenants/{tenant_id}/identity-providers** - 外部IDプロバイダ更新
  - SAML設定の更新
  - メタデータURL、メール属性の設定

### 6. Stripe連携管理
- [ ] **GET /tenants/{tenant_id}/stripe-customer** - Stripe顧客情報取得
  - 顧客ID、サブスクリプションスケジュールIDの確認
- [ ] **PATCH /stripe/init** - Stripe初期設定
  - billing経由でのStripe初期情報設定
- [ ] **DELETE /stripe** - Stripe情報削除
  - 顧客情報・商品情報の削除

### 7. プラン関連の管理
- [ ] **PUT /plans/reset** - プラン情報リセット
  - 全プラン関連情報の削除
  - Stripe連携の解除

## 前提条件
- 基本セットアップが完了している
- 役割(ロール)が定義されている
- ユーザー管理機能が動作している

## 期待結果
- テナントが正常に作成・管理できる
- テナント属性が適切に設定される
- 請求情報が正しく管理される
- プラン管理が正常に動作する
- 外部IDプロバイダ設定が機能する
- Stripe連携が正常に動作する

## 注意事項
- テナント削除時は関連するユーザー情報も削除される
- Stripe連携時は実際の課金が発生する可能性がある
- プラン変更時の比例配分設定は慎重に行う