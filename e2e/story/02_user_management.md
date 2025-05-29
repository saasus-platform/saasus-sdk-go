# 02. ユーザー管理ストーリー

## 概要
SaaSユーザーの作成、管理、属性設定を行うE2Eテストストーリー

## テスト対象API
- saasUser.yml
- userAttribute.yml
- userInfo.yml

## ストーリー

### 1. ユーザー属性定義の管理
- [ ] **GET /user-attributes** - ユーザー属性一覧を取得
  - 既存のユーザー属性定義の確認
- [ ] **POST /user-attributes** - ユーザー属性を作成
  - 誕生日、住所などの属性定義を作成
- [ ] **DELETE /user-attributes/{attribute_name}** - ユーザー属性を削除
  - 不要な属性定義の削除

### 2. SaaSユーザー属性の管理
- [ ] **POST /saas-user-attributes** - SaaSユーザー属性を作成
  - 全テナント共通の属性定義を作成

### 3. ユーザーの新規登録
- [ ] **POST /sign-up** - ユーザー新規登録
  - メールアドレスでの新規登録
- [ ] **POST /sign-up/resend** - 確認メール再送信
  - 仮パスワードの再送信

### 4. 外部アカウント連携
- [ ] **POST /external-users/request** - 外部アカウント連携要求
  - 外部アカウントとの連携要求
- [ ] **POST /external-users/confirm** - 外部アカウント連携確認
  - 検証コードによる連携確認

### 5. AWS Marketplace連携
- [ ] **POST /aws-marketplace/sign-up** - AWS Marketplace新規登録
  - AWS Marketplaceからの新規登録
- [ ] **POST /aws-marketplace/sign-up-confirm** - AWS Marketplace登録確認
  - 登録の確定とテナント作成
- [ ] **PATCH /aws-marketplace/link** - 既存テナントとの連携
  - 既存テナントとAWS Marketplaceの連携

### 6. SaaSユーザーの管理
- [ ] **GET /users** - ユーザー一覧を取得
  - 全SaaSユーザーの一覧取得
- [ ] **POST /users** - SaaSユーザーを作成
  - 管理者によるユーザー作成
- [ ] **GET /users/{user_id}** - ユーザー情報を取得
  - 特定ユーザーの詳細情報取得
- [ ] **DELETE /users/{user_id}** - ユーザーを削除
  - ユーザーの完全削除

### 7. ユーザー情報の更新
- [ ] **PATCH /users/{user_id}/attributes** - ユーザー属性を更新
  - 追加属性情報の更新
- [ ] **PATCH /users/{user_id}/email** - メールアドレスを変更
  - ユーザーのメールアドレス変更
- [ ] **POST /users/{user_id}/email/request** - メールアドレス変更要求
  - メールアドレス変更の要求
- [ ] **POST /users/{user_id}/email/confirm** - メールアドレス変更確認
  - 検証コードによる変更確認

### 8. パスワード管理
- [ ] **PATCH /users/{user_id}/password** - パスワード変更
  - ユーザーのログインパスワード変更

### 9. MFA設定
- [ ] **GET /users/{user_id}/mfa/preference** - MFA設定を取得
  - ユーザーのMFA設定状況確認
- [ ] **PATCH /users/{user_id}/mfa/preference** - MFA設定を更新
  - MFAの有効化・無効化
- [ ] **POST /users/{user_id}/mfa/software-token/secret-code** - シークレットコード作成
  - 認証アプリ用のシークレットコード生成
- [ ] **PUT /users/{user_id}/mfa/software-token** - 認証アプリ登録
  - 認証アプリケーションの登録

### 10. 外部IDプロバイダ連携解除
- [ ] **DELETE /users/{user_id}/providers/{provider_name}** - プロバイダ連携解除
  - Google等の外部IDプロバイダとの連携解除

### 11. ユーザー情報取得
- [ ] **GET /userinfo** - IDトークンからユーザー情報取得
  - IDトークンを使用したユーザー情報取得
- [ ] **GET /userinfo/search/email** - メールアドレスからユーザー情報取得
  - メールアドレスを条件としたユーザー検索

## 前提条件
- 基本セットアップが完了している
- ユーザー属性が定義されている
- 認証設定が完了している

## 期待結果
- ユーザーが正常に作成・管理できる
- ユーザー属性が適切に設定される
- MFA設定が正常に動作する
- 外部アカウント連携が機能する
- AWS Marketplace連携が機能する