package testdata

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// TestParams はテストパラメータファイルの構造体
// test_params.jsonの構造に対応しています
type TestParams struct {
	// 基本設定関連
	AuthInfo struct {
		UpdateParams map[string]any `json:"updateParams"` // callback_url
	} `json:"auth-info"`

	BasicInfo struct {
		UpdateParams map[string]any `json:"updateParams"` // domain_name, from_email_address, reply_email_address
	} `json:"basic-info"`

	// ユーザー管理関連
	Users struct {
		CreateParams map[string]any `json:"createParams"` // email, password, 属性情報
	} `json:"users"`

	UsersEmail struct {
		UpdateParams map[string]any `json:"updateParams"` // email
	} `json:"users-email"`

	UsersPassword struct {
		UpdateParams map[string]any `json:"updateParams"` // password
	} `json:"users-password"`

	// MFA設定関連
	UsersMfaPreference struct {
		UpdateParams map[string]any `json:"updateParams"` // enabled, method
	} `json:"users-mfa-preference"`

	UsersMfaSoftwareToken struct {
		UpdateParams map[string]any `json:"updateParams"` // access_token, verification_code
	} `json:"users-mfa-software-token"`

	UsersMfaSoftwareTokenSecretCode struct {
		CreateParams map[string]any `json:"createParams"` // access_token
	} `json:"users-mfa-software-token-secret-code"`

	// ロール管理関連
	Roles struct {
		CreateParams map[string]any `json:"createParams"` // role_name, display_name
	} `json:"roles"`

	// 属性管理関連
	UserAttributes struct {
		CreateParams map[string]any `json:"createParams"` // attribute_name, display_name, attribute_type
	} `json:"user-attributes"`

	TenantAttributes struct {
		CreateParams map[string]any `json:"createParams"` // attribute_name, display_name, attribute_type
	} `json:"tenant-attributes"`

	// 環境管理関連
	Envs struct {
		CreateParams map[string]any `json:"createParams"` // id, name, display_name
		UpdateParams map[string]any `json:"updateParams"` // name, display_name
	} `json:"envs"`

	// テナント管理関連
	Tenants struct {
		CreateParams map[string]any `json:"createParams"` // name, attributes, back_office_staff_email
		UpdateParams map[string]any `json:"updateParams"` // name, attributes, back_office_staff_email
	} `json:"tenants"`

	TenantsBillingInfo struct {
		CreateParams map[string]any `json:"createParams"` // address, invoice_language, name
		UpdateParams map[string]any `json:"updateParams"` // address, invoice_language, name
	} `json:"tenants-billing-info"`

	TenantsIdentityProviders struct {
		UpdateParams map[string]any `json:"updateParams"` // identity_provider_props, provider_type
	} `json:"tenants-identity-providers"`

	TenantsPlans struct {
		CreateParams map[string]any `json:"createParams"` // プラン関連パラメータ
		UpdateParams map[string]any `json:"updateParams"` // delete_usage, next_plan_id, etc.
	} `json:"tenants-plans"`

	TenantsUsers struct {
		CreateParams map[string]any `json:"createParams"` // email, attributes
		UpdateParams map[string]any `json:"updateParams"` // attributes
	} `json:"tenants-users"`

	TenantsUsersEnvsRoles struct {
		CreateParams map[string]any `json:"createParams"` // role_names
	} `json:"tenants-users-envs-roles"`

	// カスタマイズ関連
	CustomizePages struct {
		CreateParams map[string]any `json:"createParams"` // カスタマイズページ設定
		UpdateParams map[string]any `json:"updateParams"` // password_reset_page, sign_in_page, sign_up_page
	} `json:"customize-pages"`

	CustomizePageSettings struct {
		CreateParams map[string]any `json:"createParams"` // ページ設定
		UpdateParams map[string]any `json:"updateParams"` // favicon, icon, title, etc.
	} `json:"customize-page-settings"`

	// 通知関連
	NotificationMessages struct {
		CreateParams map[string]any `json:"createParams"` // 通知メッセージ
		UpdateParams map[string]any `json:"updateParams"` // authentication_mfa, create_user, etc.
	} `json:"notification-messages"`

	// サインイン設定関連
	SignInSettings struct {
		UpdateParams map[string]any `json:"updateParams"` // account_verification, device_configuration, etc.
	} `json:"sign-in-settings"`

	// ID プロバイダー関連
	IdentityProviders struct {
		CreateParams map[string]any `json:"createParams"` // provider, tenant
		UpdateParams map[string]any `json:"updateParams"` // identity_provider_props, provider
	} `json:"identity-providers"`

	// 認証情報関連
	Credentials struct {
		CreateParams map[string]any `json:"createParams"` // access_token, id_token, refresh_token
	} `json:"credentials"`

	// AWS Marketplace関連
	AwsMarketplaceSignUp struct {
		CreateParams map[string]any `json:"createParams"` // email, registration_token, etc.
	} `json:"aws-marketplace-sign-up"`

	ListingStatus struct {
		UpdateParams map[string]any `json:"updateParams"` // listing_status
	} `json:"listing-status"`

	Settings struct {
		UpdateParams map[string]any `json:"updateParams"` // AWS Marketplace設定
	} `json:"settings"`

	// Single Tenant関連
	SingleTenantSettings struct {
		CreateParams map[string]any `json:"createParams"` // cloudformation_template, ddl_template, etc.
		UpdateParams map[string]any `json:"updateParams"` // cloudformation_template, ddl_template, etc.
	} `json:"single-tenant-settings"`

	// Stripe関連
	StripeInfo struct {
		UpdateParams map[string]any `json:"updateParams"` // secret_key
	} `json:"stripe-info"`

	// Pricing関連（他のAPIで使用）
	Menus struct {
		CreateParams map[string]any `json:"createParams"`
	} `json:"menus"`

	MeteringTenantsUnitsTimestamp struct {
		UpdateParams map[string]any `json:"updateParams"`
	} `json:"metering-tenants-units-timestamp"`

	MeteringUnits struct {
		CreateParams map[string]any `json:"createParams"`
	} `json:"metering-units"`

	Plans struct {
		CreateParams map[string]any `json:"createParams"`
	} `json:"plans"`

	TaxRates struct {
		CreateParams map[string]any `json:"createParams"`
	} `json:"tax-rates"`

	Units struct {
		CreateParams map[string]any `json:"createParams"`
	} `json:"units"`

	// 追加ストーリー用パラメータ
	SaasUserAttributes struct {
		UpdateParams map[string]any `json:"updateParams"` // attributes
	} `json:"saas-user-attributes"`

	TenantInvitations struct {
		CreateParams map[string]any `json:"createParams"` // email, envs, access_token
	} `json:"tenant-invitations"`

	ExternalUserLink struct {
		RequestParams map[string]any `json:"requestParams"` // access_token
		ConfirmParams map[string]any `json:"confirmParams"` // access_token, code
	} `json:"external-user-link"`

	EmailUpdate struct {
		RequestParams map[string]any `json:"requestParams"` // email, access_token
		ConfirmParams map[string]any `json:"confirmParams"` // code, access_token
	} `json:"email-update"`

	SignUp struct {
		CreateParams map[string]any `json:"createParams"` // email
	} `json:"sign-up"`

	ProviderManagement struct {
		UnlinkParams map[string]any `json:"unlinkParams"` // provider_name
	} `json:"provider-management"`
}

// StringPtr は文字列のポインタを返すヘルパー関数
func StringPtr(s string) *string {
	return &s
}

// IntPtr は整数のポインタを返すヘルパー関数
func IntPtr(i int) *int {
	return &i
}

// BoolPtr はブール値のポインタを返すヘルパー関数
func BoolPtr(b bool) *bool {
	return &b
}

// LoadTestParams はテストパラメータファイルを読み込みます
// testdataディレクトリ内のtest_params.jsonファイルを読み込み、
// TestParams構造体にデシリアライズして返します。
func LoadTestParams(t *testing.T) *TestParams {
	// 現在のディレクトリからの相対パスを試す
	paramFile := "test_params.json"
	data, err := os.ReadFile(paramFile)
	if err != nil {
		// testdataディレクトリからの相対パスを試す
		paramFile = filepath.Join("testdata", "test_params.json")
		data, err = os.ReadFile(paramFile)
		if err != nil {
			t.Fatalf("テストパラメータファイルの読み込みに失敗: %v", err)
		}
	}

	var params TestParams
	if err := json.Unmarshal(data, &params); err != nil {
		t.Fatalf("テストパラメータファイルの解析に失敗: %v", err)
	}

	return &params
}
