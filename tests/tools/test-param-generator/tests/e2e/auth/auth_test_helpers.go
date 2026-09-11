package auth

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// TestParams はテストパラメータファイルの構造体
type TestParams struct {

	AuthInfo struct {


		UpdateParams map[string]any `json:"updateParams"`

	} `json:"auth-info"`

	AwsMarketplaceSignUp struct {

		CreateParams map[string]any `json:"createParams"`


	} `json:"aws-marketplace-sign-up"`

	BasicInfo struct {


		UpdateParams map[string]any `json:"updateParams"`

	} `json:"basic-info"`

	Credentials struct {

		CreateParams map[string]any `json:"createParams"`


	} `json:"credentials"`

	CustomizePageSettings struct {

		CreateParams map[string]any `json:"createParams"`


		UpdateParams map[string]any `json:"updateParams"`

	} `json:"customize-page-settings"`

	CustomizePages struct {

		CreateParams map[string]any `json:"createParams"`


		UpdateParams map[string]any `json:"updateParams"`

	} `json:"customize-pages"`

	Envs struct {

		CreateParams map[string]any `json:"createParams"`


		UpdateParams map[string]any `json:"updateParams"`

	} `json:"envs"`

	IdentityProviders struct {

		CreateParams map[string]any `json:"createParams"`


		UpdateParams map[string]any `json:"updateParams"`

	} `json:"identity-providers"`

	ListingStatus struct {


		UpdateParams map[string]any `json:"updateParams"`

	} `json:"listing-status"`

	Menus struct {

		CreateParams map[string]any `json:"createParams"`


	} `json:"menus"`

	MeteringTenantsUnitsTimestamp struct {


		UpdateParams map[string]any `json:"updateParams"`

	} `json:"metering-tenants-units-timestamp"`

	MeteringUnits struct {

		CreateParams map[string]any `json:"createParams"`


	} `json:"metering-units"`

	NotificationMessages struct {

		CreateParams map[string]any `json:"createParams"`


		UpdateParams map[string]any `json:"updateParams"`

	} `json:"notification-messages"`

	Plans struct {

		CreateParams map[string]any `json:"createParams"`


	} `json:"plans"`

	Roles struct {

		CreateParams map[string]any `json:"createParams"`


	} `json:"roles"`

	Settings struct {


		UpdateParams map[string]any `json:"updateParams"`

	} `json:"settings"`

	SignInSettings struct {


		UpdateParams map[string]any `json:"updateParams"`

	} `json:"sign-in-settings"`

	SingleTenantSettings struct {

		CreateParams map[string]any `json:"createParams"`


		UpdateParams map[string]any `json:"updateParams"`

	} `json:"single-tenant-settings"`

	StripeInfo struct {


		UpdateParams map[string]any `json:"updateParams"`

	} `json:"stripe-info"`

	TaxRates struct {

		CreateParams map[string]any `json:"createParams"`


	} `json:"tax-rates"`

	TenantAttributes struct {

		CreateParams map[string]any `json:"createParams"`


	} `json:"tenant-attributes"`

	Tenants struct {

		CreateParams map[string]any `json:"createParams"`


		UpdateParams map[string]any `json:"updateParams"`

	} `json:"tenants"`

	TenantsBillingInfo struct {

		CreateParams map[string]any `json:"createParams"`


		UpdateParams map[string]any `json:"updateParams"`

	} `json:"tenants-billing-info"`

	TenantsIdentityProviders struct {


		UpdateParams map[string]any `json:"updateParams"`

	} `json:"tenants-identity-providers"`

	TenantsPlans struct {

		CreateParams map[string]any `json:"createParams"`


		UpdateParams map[string]any `json:"updateParams"`

	} `json:"tenants-plans"`

	TenantsUsers struct {

		CreateParams map[string]any `json:"createParams"`


		UpdateParams map[string]any `json:"updateParams"`

	} `json:"tenants-users"`

	TenantsUsersEnvsRoles struct {

		CreateParams map[string]any `json:"createParams"`


	} `json:"tenants-users-envs-roles"`

	Units struct {

		CreateParams map[string]any `json:"createParams"`


	} `json:"units"`

	UserAttributes struct {

		CreateParams map[string]any `json:"createParams"`


	} `json:"user-attributes"`

	Users struct {

		CreateParams map[string]any `json:"createParams"`


	} `json:"users"`

	UsersEmail struct {


		UpdateParams map[string]any `json:"updateParams"`

	} `json:"users-email"`

	UsersMfaPreference struct {


		UpdateParams map[string]any `json:"updateParams"`

	} `json:"users-mfa-preference"`

	UsersMfaSoftwareToken struct {


		UpdateParams map[string]any `json:"updateParams"`

	} `json:"users-mfa-software-token"`

	UsersMfaSoftwareTokenSecretCode struct {

		CreateParams map[string]any `json:"createParams"`


	} `json:"users-mfa-software-token-secret-code"`

	UsersPassword struct {


		UpdateParams map[string]any `json:"updateParams"`

	} `json:"users-password"`

}


// stringPtr は文字列のポインタを返すヘルパー関数
// 要件: 5.1, 5.2
func stringPtr(s string) *string {
	return &s
}

// loadTestParams はテストパラメータファイルを読み込みます
// 要件: 5.3, 5.4, 5.5
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
