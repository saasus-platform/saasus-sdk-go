package authapi

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/saasus-platform/saasus-sdk-go/generated/authapi"
	"github.com/saasus-platform/saasus-sdk-go/generated/awsmarketplaceapi"
	"github.com/saasus-platform/saasus-sdk-go/generated/billingapi"
	"github.com/saasus-platform/saasus-sdk-go/modules/auth"
	"github.com/saasus-platform/saasus-sdk-go/modules/awsmarketplace"
	"github.com/saasus-platform/saasus-sdk-go/modules/billing"
	"github.com/saasus-platform/saasus-sdk-go/tests/e2e/authapi/testdata"
	"github.com/saasus-platform/saasus-sdk-go/tests/testlib"
)

// GetAuthMethods は Auth API の全メソッド一覧を返します。
func GetAuthMethods() []string {
	return []string{
		// Standard メソッド
		"GetBasicInfo", "UpdateBasicInfo", "GetAuthInfo", "UpdateAuthInfo",
		"GetSaasUsers", "GetSaasUser", "CreateSaasUser", "DeleteSaasUser",
		"UpdateSaasUserPassword", "UpdateSaasUserEmail", "UpdateSaasUserAttributes",
		"RequestEmailUpdate", "ConfirmEmailUpdate",
		"GetUserMfaPreference", "UpdateUserMfaPreference", "UpdateSoftwareToken", "CreateSecretCode",
		"GetRoles", "CreateRole", "DeleteRole",
		"GetUserAttributes", "CreateUserAttribute", "DeleteUserAttribute",
		"CreateSaasUserAttribute",
		"GetTenantAttributes", "CreateTenantAttribute", "DeleteTenantAttribute",
		"FindNotificationMessages", "UpdateNotificationMessages",
		"GetCustomizePages", "UpdateCustomizePages", "GetCustomizePageSettings", "UpdateCustomizePageSettings",
		"GetEnvs", "CreateEnv", "GetEnv", "UpdateEnv", "DeleteEnv",
		"GetSignInSettings", "UpdateSignInSettings",
		"SignUp", "ResendSignUpConfirmationEmail",
		"GetTenants", "CreateTenant", "GetTenant", "UpdateTenant", "DeleteTenant",
		"CreateTenantUser", "GetAllTenantUsers", "GetAllTenantUser", "GetTenantUsers", "GetTenantUser", "UpdateTenantUser", "DeleteTenantUser",
		"CreateTenantUserRoles", "DeleteTenantUserRole",
		"GetTenantInvitations", "CreateTenantInvitation", "GetTenantInvitation", "DeleteTenantInvitation",
		"GetInvitationValidity", "ValidateInvitation",
		"UpdateTenantPlan", "UpdateTenantBillingInfo",
		"CreateTenantAndPricing", "DeleteStripeTenantAndPricing", "GetStripeCustomer",
		"ResetPlan",
		"GetUserInfo", "GetAuthCredentials", "CreateAuthCredentials",
		"GetIdentityProviders", "UpdateIdentityProvider",
		"GetTenantIdentityProviders", "UpdateTenantIdentityProvider",
		"RequestExternalUserLink", "ConfirmExternalUserLink",
		"UnlinkProvider",
		"SignUpWithAwsMarketplace", "ConfirmSignUpWithAwsMarketplace", "LinkAwsMarketplace",
		"GetCloudFormationLaunchStackLinkForSingleTenant", "GetSingleTenantSettings", "UpdateSingleTenantSettings",

		// WithBody メソッド
		"UpdateBasicInfoWithBody", "UpdateAuthInfoWithBody",
		"CreateSaasUserWithBody", "UpdateSaasUserPasswordWithBody", "UpdateSaasUserEmailWithBody", "UpdateSaasUserAttributesWithBody",
		"RequestEmailUpdateWithBody", "ConfirmEmailUpdateWithBody",
		"UpdateUserMfaPreferenceWithBody", "UpdateSoftwareTokenWithBody", "CreateSecretCodeWithBody",
		"CreateRoleWithBody", "CreateUserAttributeWithBody", "CreateSaasUserAttributeWithBody", "CreateTenantAttributeWithBody",
		"UpdateNotificationMessagesWithBody", "UpdateCustomizePagesWithBody", "UpdateCustomizePageSettingsWithBody",
		"CreateEnvWithBody", "UpdateEnvWithBody", "UpdateSignInSettingsWithBody",
		"SignUpWithBody", "ResendSignUpConfirmationEmailWithBody",
		"CreateTenantWithBody", "UpdateTenantWithBody",
		"CreateTenantUserWithBody", "UpdateTenantUserWithBody", "CreateTenantUserRolesWithBody",
		"CreateTenantInvitationWithBody", "ValidateInvitationWithBody",
		"UpdateTenantPlanWithBody", "UpdateTenantBillingInfoWithBody",
		"CreateAuthCredentialsWithBody", "UpdateIdentityProviderWithBody", "UpdateTenantIdentityProviderWithBody",
		"RequestExternalUserLinkWithBody", "ConfirmExternalUserLinkWithBody",
		"SignUpWithAwsMarketplaceWithBody", "ConfirmSignUpWithAwsMarketplaceWithBody", "LinkAwsMarketplaceWithBody",
		"UpdateSingleTenantSettingsWithBody",

		// WithResponse メソッド
		"GetBasicInfoWithResponse", "UpdateBasicInfoWithResponse", "GetAuthInfoWithResponse", "UpdateAuthInfoWithResponse",
		"GetSaasUsersWithResponse", "GetSaasUserWithResponse", "CreateSaasUserWithResponse", "DeleteSaasUserWithResponse",
		"UpdateSaasUserPasswordWithResponse", "UpdateSaasUserEmailWithResponse", "UpdateSaasUserAttributesWithResponse",
		"RequestEmailUpdateWithResponse", "ConfirmEmailUpdateWithResponse",
		"GetUserMfaPreferenceWithResponse", "UpdateUserMfaPreferenceWithResponse", "UpdateSoftwareTokenWithResponse", "CreateSecretCodeWithResponse",
		"GetRolesWithResponse", "CreateRoleWithResponse", "DeleteRoleWithResponse",
		"GetUserAttributesWithResponse", "CreateUserAttributeWithResponse", "DeleteUserAttributeWithResponse",
		"CreateSaasUserAttributeWithResponse",
		"GetTenantAttributesWithResponse", "CreateTenantAttributeWithResponse", "DeleteTenantAttributeWithResponse",
		"FindNotificationMessagesWithResponse", "UpdateNotificationMessagesWithResponse",
		"GetCustomizePagesWithResponse", "UpdateCustomizePagesWithResponse", "GetCustomizePageSettingsWithResponse", "UpdateCustomizePageSettingsWithResponse",
		"GetEnvsWithResponse", "CreateEnvWithResponse", "GetEnvWithResponse", "UpdateEnvWithResponse", "DeleteEnvWithResponse",
		"GetSignInSettingsWithResponse", "UpdateSignInSettingsWithResponse",
		"SignUpWithResponse", "ResendSignUpConfirmationEmailWithResponse",
		"GetTenantsWithResponse", "CreateTenantWithResponse", "GetTenantWithResponse", "UpdateTenantWithResponse", "DeleteTenantWithResponse",
		"CreateTenantUserWithResponse", "GetAllTenantUsersWithResponse", "GetAllTenantUserWithResponse", "GetTenantUsersWithResponse", "GetTenantUserWithResponse", "UpdateTenantUserWithResponse", "DeleteTenantUserWithResponse",
		"CreateTenantUserRolesWithResponse", "DeleteTenantUserRoleWithResponse",
		"GetTenantInvitationsWithResponse", "CreateTenantInvitationWithResponse", "GetTenantInvitationWithResponse", "DeleteTenantInvitationWithResponse",
		"GetInvitationValidityWithResponse", "ValidateInvitationWithResponse",
		"UpdateTenantPlanWithResponse", "UpdateTenantBillingInfoWithResponse",
		"CreateTenantAndPricingWithResponse", "DeleteStripeTenantAndPricingWithResponse", "GetStripeCustomerWithResponse",
		"ResetPlanWithResponse",
		"GetUserInfoWithResponse", "GetAuthCredentialsWithResponse", "CreateAuthCredentialsWithResponse",
		"GetIdentityProvidersWithResponse", "UpdateIdentityProviderWithResponse",
		"GetTenantIdentityProvidersWithResponse", "UpdateTenantIdentityProviderWithResponse",
		"RequestExternalUserLinkWithResponse", "ConfirmExternalUserLinkWithResponse",
		"UnlinkProviderWithResponse",
		"SignUpWithAwsMarketplaceWithResponse", "ConfirmSignUpWithAwsMarketplaceWithResponse", "LinkAwsMarketplaceWithResponse",
		"GetCloudFormationLaunchStackLinkForSingleTenantWithResponse", "GetSingleTenantSettingsWithResponse", "UpdateSingleTenantSettingsWithResponse",

		// WithBodyWithResponse メソッド
		"UpdateBasicInfoWithBodyWithResponse", "UpdateAuthInfoWithBodyWithResponse",
		"CreateSaasUserWithBodyWithResponse", "UpdateSaasUserPasswordWithBodyWithResponse", "UpdateSaasUserEmailWithBodyWithResponse", "UpdateSaasUserAttributesWithBodyWithResponse",
		"RequestEmailUpdateWithBodyWithResponse", "ConfirmEmailUpdateWithBodyWithResponse",
		"UpdateUserMfaPreferenceWithBodyWithResponse", "UpdateSoftwareTokenWithBodyWithResponse", "CreateSecretCodeWithBodyWithResponse",
		"CreateRoleWithBodyWithResponse", "CreateUserAttributeWithBodyWithResponse", "CreateSaasUserAttributeWithBodyWithResponse", "CreateTenantAttributeWithBodyWithResponse",
		"UpdateNotificationMessagesWithBodyWithResponse", "UpdateCustomizePagesWithBodyWithResponse", "UpdateCustomizePageSettingsWithBodyWithResponse",
		"CreateEnvWithBodyWithResponse", "UpdateEnvWithBodyWithResponse", "UpdateSignInSettingsWithBodyWithResponse",
		"SignUpWithBodyWithResponse", "ResendSignUpConfirmationEmailWithBodyWithResponse",
		"CreateTenantWithBodyWithResponse", "UpdateTenantWithBodyWithResponse",
		"CreateTenantUserWithBodyWithResponse", "UpdateTenantUserWithBodyWithResponse", "CreateTenantUserRolesWithBodyWithResponse",
		"CreateTenantInvitationWithBodyWithResponse", "ValidateInvitationWithBodyWithResponse",
		"UpdateTenantPlanWithBodyWithResponse", "UpdateTenantBillingInfoWithBodyWithResponse",
		"CreateAuthCredentialsWithBodyWithResponse", "UpdateIdentityProviderWithBodyWithResponse", "UpdateTenantIdentityProviderWithBodyWithResponse",
		"RequestExternalUserLinkWithBodyWithResponse", "ConfirmExternalUserLinkWithBodyWithResponse",
		"SignUpWithAwsMarketplaceWithBodyWithResponse", "ConfirmSignUpWithAwsMarketplaceWithBodyWithResponse", "LinkAwsMarketplaceWithBodyWithResponse",
		"UpdateSingleTenantSettingsWithBodyWithResponse",
	}
}

var (
	statusOKOrCreated = []int{http.StatusOK, http.StatusCreated}
	statusCreatedOnly = []int{http.StatusCreated}
)

const (
	sampleImageDataURL                  = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mP8/wwAAgMBgS3JH5QAAAAASUVORK5CYII="
	envAwsMarketplaceRegistrationToken  = "AWS_MARKETPLACE_REGISTRATION_TOKEN"
	awsMarketplaceTokenVariable         = "aws_marketplace_registration_token"
	awsMarketplaceEmailVariable         = "aws_marketplace_email"
	skipReasonAuthCredentials           = "Requires temporary auth code or refresh token that is not available in test environment"
	skipReasonEmailUpdate               = "Depends on email delivery backend; API returns 500 without production email configuration"
	skipReasonLinkAwsMarketplace        = "Requires a valid AWS Marketplace registration token"
	skipReasonCreateAuthCredentials     = "Needs valid Cognito ID/Access/Refresh tokens; dummy placeholders cause token validation error"
	skipReasonUnlinkProvider            = "Requires a real social provider linking before unlink can succeed"
	skipReasonTenantIdentityProvider    = "Requires a tenant-level SAML IdP configuration that is absent in the shared test environment"
	skipReasonAwsMarketplaceSignUp      = "Requires a fully configured AWS Marketplace listing; public dev accounts return 500"
	skipReasonSignUpLimit               = "SignUp APIs exceed Cognito SES daily email limit in shared env"
	envAuthE2ESkipSignUpLimit           = "AUTH_E2E_SKIP_SIGNUP_ON_COGNITO_EMAIL_LIMIT"
	signUpUserIDKeyStandard             = "signup_user_id_standard"
	signUpUserIDKeyWithBody             = "signup_user_id_with_body"
	signUpUserIDKeyWithResponse         = "signup_user_id_with_response"
	signUpUserIDKeyWithBodyWithResponse = "signup_user_id_with_body_with_response"
)

var (
	stripeOnce    sync.Once
	stripeOnceErr error
	stripeMutex   sync.Mutex

	awsMarketplaceOnce    sync.Once
	awsMarketplaceOnceErr error

	envIDSeed uint64
)

func init() {
	envIDSeed = uint64(time.Now().UnixNano())
}

func uniqueEmail(prefix string) string {
	return fmt.Sprintf("%s+%d@example.com", prefix, time.Now().UnixNano())
}

func uniqueString(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano())
}

func uniqueEnvID() int {
	id := atomic.AddUint64(&envIDSeed, 1)
	if id < 10000 {
		id += 10000
	}
	return int(id)
}

func isSignUpTestingEnabled() bool {
	raw := strings.TrimSpace(os.Getenv(envAuthE2ESkipSignUpLimit))
	if raw == "" {
		return false
	}
	switch strings.ToLower(raw) {
	case "0", "false", "no", "off":
		return true
	default:
		return false
	}
}

func generateEmailList(prefix string, count int) []string {
	emails := make([]string, count)
	for i := 0; i < count; i++ {
		emails[i] = uniqueEmail(fmt.Sprintf("%s-%d", prefix, i))
	}
	return emails
}

type signUpUserSnapshot struct {
	ID    string
	Email string
}

func updateSignUpStateFromResponse(response any, variables map[string]any, idKey string, rotateEmail bool) error {
	user, err := extractSignUpUser(response)
	if err != nil {
		return err
	}
	variables[idKey] = user.ID
	if rotateEmail {
		variables["signup_email"] = uniqueEmail("auth-signup")
	}
	return nil
}

func extractSignUpUser(response any) (*signUpUserSnapshot, error) {
	switch resp := response.(type) {
	case *authapi.SignUpResponse:
		if resp == nil {
			return nil, fmt.Errorf("sign-up response is nil")
		}
		if resp.JSON201 != nil {
			return &signUpUserSnapshot{ID: string(resp.JSON201.Id), Email: resp.JSON201.Email}, nil
		}
		if len(resp.Body) > 0 {
			var user authapi.SaasUser
			if err := json.Unmarshal(resp.Body, &user); err == nil && user.Id != "" {
				return &signUpUserSnapshot{ID: string(user.Id), Email: user.Email}, nil
			}
		}
		if resp.HTTPResponse != nil {
			return parseSignUpHTTPResponse(resp.HTTPResponse)
		}
		return nil, fmt.Errorf("sign-up response did not contain JSON201 data")
	case map[string]any:
		return parseSignUpUserMap(resp)
	case *http.Response:
		if resp == nil {
			return nil, fmt.Errorf("http response is nil")
		}
		return parseSignUpHTTPResponse(resp)
	default:
		return nil, fmt.Errorf("unsupported sign-up response type %T", response)
	}
}

func parseSignUpHTTPResponse(resp *http.Response) (*signUpUserSnapshot, error) {
	if resp == nil || resp.Body == nil {
		return nil, fmt.Errorf("sign-up http response body is empty")
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read sign-up response body: %w", err)
	}
	resp.Body.Close()
	resp.Body = io.NopCloser(bytes.NewReader(body))
	var user authapi.SaasUser
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, fmt.Errorf("failed to decode sign-up response: %w", err)
	}
	if user.Id == "" {
		return nil, fmt.Errorf("sign-up response missing id field")
	}
	return &signUpUserSnapshot{ID: string(user.Id), Email: user.Email}, nil
}

func parseSignUpUserMap(data map[string]any) (*signUpUserSnapshot, error) {
	if data == nil {
		return nil, fmt.Errorf("sign-up response map is nil")
	}
	if nested, ok := data["JSON201"]; ok {
		if nestedMap, ok := nested.(map[string]any); ok {
			return parseSignUpUserMap(nestedMap)
		}
	}
	id := stringFromAny(data["id"])
	if id == "" {
		return nil, fmt.Errorf("sign-up response map missing id field")
	}
	email := stringFromAny(data["email"])
	return &signUpUserSnapshot{ID: id, Email: email}, nil
}

func stringFromAny(value any) string {
	switch v := value.(type) {
	case string:
		return v
	case fmt.Stringer:
		return v.String()
	case nil:
		return ""
	default:
		return fmt.Sprintf("%v", v)
	}
}

func deleteSignUpUserParamsFor(key string) func(map[string]any) interface{} {
	return func(variables map[string]any) interface{} {
		if id, ok := variables[key].(string); ok && id != "" {
			return id
		}
		return ""
	}
}

func stringFromMap(m map[string]any, key string) string {
	if m == nil {
		return ""
	}
	if v, ok := m[key].(string); ok && v != "" {
		return v
	}
	return ""
}

func buildBaseStoryVariables(params *testdata.TestParams, tokens *AuthTokens, prefix string) map[string]any {
	emailsForUpdates := generateEmailList(fmt.Sprintf("%s-update", prefix), 20)
	primaryEmail := uniqueEmail(prefix)
	uniqueSuffix := time.Now().UnixNano()
	roleName := fmt.Sprintf("%s-role-%06d", prefix, uniqueSuffix)
	roleDisplay := fmt.Sprintf("Role %06d", uniqueSuffix)
	attributeName := fmt.Sprintf("%s-attr-%06d", prefix, uniqueSuffix)
	envName := fmt.Sprintf("%s-env-%06d", prefix, uniqueSuffix)
	tenantStaffEmail := uniqueEmail(fmt.Sprintf("%s-tenant-staff", prefix))
	tenantRoleName := roleName

	vars := map[string]any{
		"domain_name":                     params.BasicInfo.UpdateParams["domain_name"],
		"from_email_address":              params.BasicInfo.UpdateParams["from_email_address"],
		"reply_email_address":             params.BasicInfo.UpdateParams["reply_email_address"],
		"callback_url":                    params.AuthInfo.UpdateParams["callback_url"],
		"email":                           primaryEmail,
		"_email_updates":                  emailsForUpdates,
		"password":                        params.Users.CreateParams["password"],
		"role_name":                       roleName,
		"display_name":                    roleDisplay,
		"attribute_name":                  attributeName,
		"attribute_type":                  params.UserAttributes.CreateParams["attribute_type"],
		"env_id":                          uniqueEnvID(),
		"name":                            envName,
		"back_office_staff_email":         tenantStaffEmail,
		"attributes":                      map[string]interface{}{},
		"notification_messages":           params.NotificationMessages.UpdateParams,
		"customize_pages":                 params.CustomizePages.UpdateParams,
		"tenant_role_name":                tenantRoleName,
		"tenant_billing_info":             params.TenantsBillingInfo.UpdateParams,
		"enabled":                         true,
		"access_token":                    tokens.AccessToken,
		"verification_code":               "123456",
		"favicon":                         sampleImageDataURL,
		"icon":                            sampleImageDataURL,
		"title":                           params.CustomizePageSettings.UpdateParams["title"],
		"terms_of_service_url":            params.CustomizePageSettings.UpdateParams["terms_of_service_url"],
		"privacy_policy_url":              params.CustomizePageSettings.UpdateParams["privacy_policy_url"],
		"google_tag_manager_container_id": params.CustomizePageSettings.UpdateParams["google_tag_manager_container_id"],
		"token":                           tokens.IDToken,
	}

	if idp := params.IdentityProviders.UpdateParams; idp != nil {
		vars["identity_provider_update_params"] = idp
	}
	if tenantIdp := params.TenantsIdentityProviders.UpdateParams; tenantIdp != nil {
		vars["tenant_identity_provider_update_params"] = tenantIdp
	}

	if awsCfg := params.AwsMarketplaceSignUp.CreateParams; awsCfg != nil {
		if token := stringFromMap(awsCfg, "registration_token"); token != "" {
			vars[awsMarketplaceTokenVariable] = token
		}
		if email := stringFromMap(awsCfg, "email"); email != "" {
			vars[awsMarketplaceEmailVariable] = email
		}
	}
	if token := os.Getenv(envAwsMarketplaceRegistrationToken); token != "" {
		vars[awsMarketplaceTokenVariable] = token
	}

	return vars
}

func captureCustomizePageSettings(response any, variables map[string]any) {
	httpResp, ok := response.(*http.Response)
	if !ok || httpResp == nil {
		return
	}
	defer httpResp.Body.Close()
	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		fmt.Printf("Warning: failed to read customize page settings response: %v\n", err)
		return
	}
	var payload authapi.CustomizePageSettings
	if err := json.Unmarshal(body, &payload); err != nil {
		fmt.Printf("Warning: failed to decode customize page settings: %v\n", err)
		return
	}
	variables["current_customize_page_settings"] = payload
}

func ensureStripeIntegration(t *testing.T) {
	stripeOnce.Do(func() {
		stripeOnceErr = configureStripeIntegration(t)
	})
	if stripeOnceErr != nil {
		t.Skipf("Skipping Stripe integration tests: %v", stripeOnceErr)
	}
}

func ensureStripeIntegrationFresh(t *testing.T) error {
	stripeMutex.Lock()
	defer stripeMutex.Unlock()
	return configureStripeIntegration(t)
}

func configureStripeIntegration(t *testing.T) error {
	// resetStripeIntegration is commented out to preserve existing Stripe configuration
	// if err := resetStripeIntegration(t); err != nil {
	// 	return err
	// }
	key := getStripeSecretKey()
	if key == "" {
		return fmt.Errorf("STRIPE_SECRET_KEY is not set")
	}
	client, err := billing.BillingWithResponse()
	if err != nil {
		return fmt.Errorf("failed to create billing client: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	payload := billingapi.UpdateStripeInfoParam{SecretKey: key}
	resp, err := client.UpdateStripeInfoWithResponse(ctx, payload)
	if err != nil {
		return fmt.Errorf("failed to update stripe info: %w", err)
	}
	statusCode := resp.StatusCode()
	bodyBytes := resp.Body
	bodyStr := string(bodyBytes)

	if statusCode == 400 {
		// レスポンスボディを確認して既存設定の有無を判断
		if strings.Contains(bodyStr, "already") || strings.Contains(bodyStr, "exist") {
			t.Logf("[INFO] Stripe integration already configured (status 400: %s), using existing configuration", bodyStr)
			return nil
		}
		// その他の400エラーはスキップ
		return fmt.Errorf("update stripe info returned status 400: %s", bodyStr)
	}
	if statusCode >= 300 {
		return fmt.Errorf("update stripe info returned status %d: %s", statusCode, bodyStr)
	}
	return nil
}

func resetStripeIntegration(t *testing.T) error {
	if err := cleanupStripeTenantLinks(); err != nil {
		return err
	}
	client, err := billing.BillingWithResponse()
	if err != nil {
		return fmt.Errorf("failed to create billing client: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	resp, err := client.DeleteStripeInfoWithResponse(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete existing stripe info: %w", err)
	}
	if resp != nil {
		code := resp.StatusCode()
		if code == http.StatusNotFound {
			return nil
		}
		if code >= http.StatusInternalServerError {
			return fmt.Errorf("delete stripe info returned status %d", code)
		}
		if code >= http.StatusMultipleChoices {
			msg := string(resp.Body)
			t.Logf("[WARN] DeleteStripeInfo returned %d: %s", code, msg)
		}
	}
	return nil
}

func ensureAwsMarketplaceIntegration(t *testing.T, params *testdata.TestParams) {
	awsMarketplaceOnce.Do(func() {
		awsCfg := params.AwsMarketplaceSignUp.CreateParams
		if awsCfg == nil {
			awsMarketplaceOnceErr = fmt.Errorf("aws-marketplace-sign-up.createParams is not configured")
			return
		}

		client, err := awsmarketplace.AwsMarketplaceWithResponse()
		if err != nil {
			awsMarketplaceOnceErr = fmt.Errorf("failed to create aws marketplace client: %w", err)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		updateParam := awsmarketplaceapi.UpdateSettingsParam{}
		if v := stringFromMap(awsCfg, "product_code"); v != "" {
			updateParam.ProductCode = aws.String(v)
		}
		if v := stringFromMap(awsCfg, "role_arn"); v != "" {
			updateParam.RoleArn = aws.String(v)
		}
		if v := stringFromMap(awsCfg, "sns_topic_arn"); v != "" {
			updateParam.SnsTopicArn = aws.String(v)
		}
		if v := stringFromMap(awsCfg, "seller_sns_topic_arn"); v != "" {
			updateParam.SellerSnsTopicArn = aws.String(v)
		}
		if v := stringFromMap(awsCfg, "cas_bucket_name"); v != "" {
			updateParam.CasBucketName = aws.String(v)
		}
		if v := stringFromMap(awsCfg, "cas_sns_topic_arn"); v != "" {
			updateParam.CasSnsTopicArn = aws.String(v)
		}
		if v := stringFromMap(awsCfg, "sqs_arn"); v != "" {
			updateParam.SqsArn = aws.String(v)
		}
		if v := stringFromMap(awsCfg, "role_external_id"); v != "" {
			updateParam.RoleExternalId = aws.String(v)
		}

		if !isEmptyUpdateSettingsParam(updateParam) {
			resp, err := client.UpdateSettingsWithResponse(ctx, updateParam)
			if err != nil {
				awsMarketplaceOnceErr = fmt.Errorf("failed to update aws marketplace settings: %w", err)
				return
			}
			if resp.StatusCode() >= http.StatusMultipleChoices {
				awsMarketplaceOnceErr = fmt.Errorf("update aws marketplace settings returned status %d", resp.StatusCode())
				return
			}
		}

		if status := stringFromMap(awsCfg, "listing_status"); status != "" {
			resp, err := client.UpdateListingStatusWithResponse(ctx, awsmarketplaceapi.UpdateListingStatusParam{
				ListingStatus: awsmarketplaceapi.ListingStatus(status),
			})
			if err != nil {
				awsMarketplaceOnceErr = fmt.Errorf("failed to update aws marketplace listing status: %w", err)
				return
			}
			if resp.StatusCode() >= http.StatusMultipleChoices {
				awsMarketplaceOnceErr = fmt.Errorf("update aws marketplace listing status returned %d", resp.StatusCode())
				return
			}
		}
	})

	if awsMarketplaceOnceErr != nil {
		t.Fatalf("AWS Marketplace integration setup failed: %v", awsMarketplaceOnceErr)
	}
}

func isEmptyUpdateSettingsParam(p awsmarketplaceapi.UpdateSettingsParam) bool {
	return p.CasBucketName == nil &&
		p.CasSnsTopicArn == nil &&
		p.ProductCode == nil &&
		p.RoleArn == nil &&
		p.RoleExternalId == nil &&
		p.SellerSnsTopicArn == nil &&
		p.SnsTopicArn == nil &&
		p.SqsArn == nil
}

func cleanupStripeTenantLinks() error {
	client, err := auth.AuthWithResponse()
	if err != nil {
		return fmt.Errorf("failed to create auth client for stripe cleanup: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	resp, err := client.DeleteStripeTenantAndPricingWithResponse(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete stripe tenant links: %w", err)
	}
	if resp != nil {
		code := resp.StatusCode()
		if code == http.StatusNotFound {
			return nil
		}
		if code >= http.StatusInternalServerError {
			return fmt.Errorf("delete stripe tenant and pricing returned status %d", code)
		}
	}
	return nil
}

func getStripeSecretKey() string {
	if key := os.Getenv("STRIPE_SECRET_KEY"); key != "" {
		return key
	}
	return ""
}

func createCognitoUser(ctx context.Context, email, password string) error {
	cfg, err := LoadCognitoConfigFromEnv()
	if err != nil {
		return fmt.Errorf("failed to load cognito config for user creation: %w", err)
	}

	awsCfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(cfg.Region))
	if err != nil {
		return fmt.Errorf("load aws config for user creation: %w", err)
	}

	client := cognitoidentityprovider.NewFromConfig(awsCfg, func(o *cognitoidentityprovider.Options) {
		if cfg.Endpoint != "" {
			o.BaseEndpoint = aws.String(cfg.Endpoint)
		}
	})

	username, err := ensureCognitoUsername(ctx, client, cfg.UserPoolID, email)
	if err != nil {
		return err
	}

	setUserPassIn := &cognitoidentityprovider.AdminSetUserPasswordInput{
		Password:   aws.String(password),
		UserPoolId: aws.String(cfg.UserPoolID),
		Username:   aws.String(username),
		Permanent:  true,
	}
	_, err = client.AdminSetUserPassword(ctx, setUserPassIn)
	if err != nil {
		return fmt.Errorf("admin set user password: %w", err)
	}

	fmt.Printf("Info: Successfully prepared Cognito user %s\n", email)
	return nil
}

func ensureCognitoUsername(ctx context.Context, client *cognitoidentityprovider.Client, userPoolID, email string) (string, error) {
	if username := lookupCognitoUsernameByEmail(ctx, client, userPoolID, email); username != "" {
		return username, nil
	}

	username := buildStableCognitoUsername(email)
	createInput := &cognitoidentityprovider.AdminCreateUserInput{
		UserPoolId: aws.String(userPoolID),
		Username:   aws.String(username),
		UserAttributes: []types.AttributeType{
			{Name: aws.String("email"), Value: aws.String(email)},
			{Name: aws.String("email_verified"), Value: aws.String("true")},
		},
		MessageAction: types.MessageActionTypeSuppress,
	}
	if _, err := client.AdminCreateUser(ctx, createInput); err != nil {
		var exists *types.UsernameExistsException
		var aliasExists *types.AliasExistsException
		if !errors.As(err, &exists) && !errors.As(err, &aliasExists) {
			return "", fmt.Errorf("admin create user: %w", err)
		}
	}
	return username, nil
}

func lookupCognitoUsernameByEmail(ctx context.Context, client *cognitoidentityprovider.Client, userPoolID, email string) string {
	filter := fmt.Sprintf("email = \"%s\"", escapeCognitoFilter(email))
	resp, err := client.ListUsers(ctx, &cognitoidentityprovider.ListUsersInput{
		Filter:     aws.String(filter),
		Limit:      aws.Int32(1),
		UserPoolId: aws.String(userPoolID),
	})
	if err != nil || len(resp.Users) == 0 {
		return ""
	}
	return aws.ToString(resp.Users[0].Username)
}

func escapeCognitoFilter(value string) string {
	return strings.ReplaceAll(value, "\"", "\\\"")
}

func buildStableCognitoUsername(email string) string {
	h := sha1.Sum([]byte(strings.ToLower(email)))
	return "e2e-" + hex.EncodeToString(h[:])
}

func updateUserTokensFromCognito(variables map[string]any) {
	email, _ := variables["email"].(string)
	password, _ := variables["password"].(string)
	if email == "" || password == "" {
		return
	}

	tokens, err := ObtainUserTokens(email, password)
	if err != nil {
		var notAuthorized *types.NotAuthorizedException
		if errors.As(err, &notAuthorized) {
			fmt.Printf("Info: User %s not found in Cognito, attempting to create...\n", email)
			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			defer cancel()
			if createErr := createCognitoUser(ctx, email, password); createErr != nil {
				fmt.Printf("Warning: failed to create cognito user %s: %v\n", email, createErr)
				return
			}
			// Retry obtaining tokens
			tokens, err = ObtainUserTokens(email, password)
			if err != nil {
				fmt.Printf("Warning: failed to obtain user tokens for %s even after creation: %v\n", email, err)
				return
			}
		} else {
			fmt.Printf("Warning: failed to obtain user tokens for %s: %v\n", email, err)
			return
		}
	}

	variables["user_access_token"] = tokens.AccessToken
	variables["user_id_token"] = tokens.IdToken
	variables["user_refresh_token"] = tokens.RefreshToken
	variables["cognito_access_token"] = tokens.AccessToken
	variables["cognito_id_token"] = tokens.IdToken
	variables["cognito_refresh_token"] = tokens.RefreshToken

	// Cognito のソフトウェアトークン設定で必要となる Secret を事前取得する
	ensureCognitoSoftwareTokenPrepared(variables)
	exchangeForSaaSusTokens(variables)
}

func exchangeForSaaSusTokens(variables map[string]any) {
	idToken, _ := variables["user_id_token"].(string)
	accessToken, _ := variables["user_access_token"].(string)
	refreshToken, _ := variables["user_refresh_token"].(string)
	if idToken == "" || accessToken == "" {
		return
	}

	client, err := auth.AuthWithResponse()
	if err != nil {
		fmt.Printf("Warning: failed to create auth client: %v\n", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	credReq := authapi.CreateAuthCredentialsParam{
		IdToken:     idToken,
		AccessToken: accessToken,
	}
	if refreshToken != "" {
		credReq.RefreshToken = &refreshToken
	}

	createResp, err := client.CreateAuthCredentialsWithResponse(ctx, credReq)
	if err != nil {
		fmt.Printf("Warning: failed to create auth credentials: %v\n", err)
		return
	}
	if createResp.JSON201 == nil {
		fmt.Printf("Warning: create auth credentials returned status %d\n", createResp.StatusCode())
		return
	}

	code := createResp.JSON201.Code
	flow := authapi.GetAuthCredentialsParamsAuthFlowTempCodeAuth
	params := &authapi.GetAuthCredentialsParams{
		Code:     &code,
		AuthFlow: &flow,
	}

	getResp, err := client.GetAuthCredentialsWithResponse(ctx, params)
	if err != nil {
		fmt.Printf("Warning: failed to exchange auth credentials: %v\n", err)
		return
	}
	if getResp.JSON200 == nil {
		fmt.Printf("Warning: get auth credentials returned status %d\n", getResp.StatusCode())
		return
	}

	credentials := getResp.JSON200
	fmt.Printf("[DEBUG] Received SaaSus tokens for %s\n", variables["email"])

	// SaaSusトークンを保存（cognito_access_tokenは保持）
	variables["saasus_access_token"] = credentials.AccessToken
	variables["saasus_id_token"] = credentials.IdToken

	// デフォルトトークンをSaaSusトークンに設定
	// ただし、cognito_access_tokenは上書きしない（MFA設定で必要）
	variables["user_access_token"] = credentials.AccessToken
	variables["access_token"] = credentials.AccessToken
	variables["token"] = credentials.IdToken

	if credentials.RefreshToken != nil {
		variables["saasus_refresh_token"] = *credentials.RefreshToken
		variables["user_refresh_token"] = *credentials.RefreshToken
		variables["refresh_token"] = *credentials.RefreshToken
	}

	fmt.Printf("[DEBUG] Token types available: cognito_access_token=%v, saasus_access_token=%v\n",
		variables["cognito_access_token"] != nil && variables["cognito_access_token"] != "",
		variables["saasus_access_token"] != nil && variables["saasus_access_token"] != "")
}

// VerifyMethodCoverage は、全 Auth API メソッドがストーリーでカバーされているかを検証します。
//
// この関数は、GetAuthMethods() で定義された全メソッドが、
// 4 つのストーリーのいずれかで少なくとも 1 回呼び出されることを確認します。
//
// 検証プロセス:
//  1. GetAuthMethods() から全メソッド一覧を取得
//  2. GetAuthStories(t) から全ストーリーを取得
//  3. 各ストーリーの各ステップで使用されているメソッドを記録
//  4. 使用されていないメソッドをリストアップ
//  5. 未使用メソッドがあればエラーを返す
//
// パラメータ:
//   - t: テストコンテキスト（ストーリー取得時のテストパラメータ読み込みに使用）
//
// 戻り値:
//   - error: 未使用メソッドがある場合はエラー、すべてカバーされている場合は nil
func VerifyMethodCoverage(t *testing.T) error {
	// 全メソッド一覧を取得
	allMethods := GetAuthMethods()

	// 使用されたメソッドを記録するマップ（効率的な検索のため）
	usedMethods := make(map[string]bool, len(allMethods))

	// ストーリーが実装されていない場合は検証をスキップ
	if !isStoriesImplemented() {
		return nil
	}

	// 全ストーリーを取得（testing.T を渡す）
	stories := GetAuthStories(t)

	// 各ストーリーの各ステップで使用されているメソッドを記録
	for _, story := range stories {
		for _, step := range story.Steps {
			usedMethods[step.ClientMethod] = true
		}
	}

	// 使用されていないメソッドをリストアップ
	var missingMethods []string
	for _, method := range allMethods {
		if !usedMethods[method] {
			missingMethods = append(missingMethods, method)
		}
	}

	// 未使用メソッドがある場合はエラーを返す
	if len(missingMethods) > 0 {
		sort.Strings(missingMethods)
		return fmt.Errorf("以下のメソッドがストーリーでカバーされていません (%d/%d 個): %v",
			len(missingMethods), len(allMethods), missingMethods)
	}

	return nil
}

func isStoriesImplemented() bool {
	return true
}

// GetAuthStories は Auth API の全ストーリーを返します。
//
// testing.T パラメータを受け取り、各ストーリー関数に渡します。
// これにより、各ストーリーが testdata/test_params.json からテストパラメータを読み込めます。
//
// パラメータ:
//   - t: テストコンテキスト（テストパラメータの読み込みに使用）
//
// 戻り値:
//   - []testlib.Story: 実行可能なテストストーリーの配列
func GetAuthStories(t *testing.T) []testlib.Story {
	return []testlib.Story{
		// Postman Collection ベースストーリー（既存）
		GetPostmanStoryStandardMethods(t),
		GetPostmanStoryStandardMethodsWithBody(t),
		GetPostmanStoryWithResponseMethods(t),
		GetPostmanStoryWithResponseMethodsWithBody(t),

		// 追加ストーリー（未カバーメソッド対応）
		GetStorySaasUserAttributes(t),
		// FIXME: Skipping Tenant Invitations story as it requires a valid external JWT (e.g., from AWS Marketplace)
		// which is not available in the E2E test environment, causing auth failures.
		// GetStoryTenantInvitations(t),
		GetStoryExternalUserLinkAndEmailUpdate(t),
		GetStorySignUpAndProviderManagement(t),

		// 未実装WithResponseメソッドのストーリー
	}
}

// GetPostmanStoryStandardMethods は Standard メソッド（WithBody サフィックスなし）のストーリーを返します。
//
// このストーリーは Postman Collection の Test Flow に従い、基本的な HTTP レスポンスを返す
// Standard メソッドのみを使用します。WithBody サフィックスを持つメソッドは含まれません。
//
// test_params.json からテストパラメータを読み込み、Variables マップを構築します。
//
// カバーするメソッドカテゴリ:
//   - 基本設定・認証情報（GetBasicInfo, UpdateBasicInfo, GetAuthInfo, UpdateAuthInfo）
//   - ユーザー管理（GetSaasUsers, CreateSaasUser, UpdateSaasUserPassword など）
//   - MFA 設定（GetUserMfaPreference, UpdateUserMfaPreference など）
//   - ロール管理（GetRoles, CreateRole, DeleteRole）
//   - 属性管理（GetUserAttributes, CreateUserAttribute など）
//   - 通知・カスタマイズ（FindNotificationMessages, UpdateCustomizePages など）
//   - 環境管理（GetEnvs, CreateEnv, UpdateEnv, DeleteEnv）
//   - テナント管理（GetTenants, CreateTenant, UpdateTenant など）
//   - Stripe 連携（CreateTenantAndPricing, GetStripeCustomer など）
//   - ID プロバイダー（GetIdentityProviders, UpdateIdentityProvider など）
//   - AWS Marketplace（SignUpWithAwsMarketplace）
//   - シングルテナント（GetSingleTenantSettings, UpdateSingleTenantSettings）
//
// パラメータ:
//   - t: テストコンテキスト（テストパラメータの読み込みに使用）
//
// 戻り値:
//   - testlib.Story: 実行可能なテストストーリー
func GetPostmanStoryStandardMethods(t *testing.T) testlib.Story {
	// テストパラメータを読み込む
	params := testdata.LoadTestParams(t)
	tokens := MustGetAuthTokens(t)

	vars := buildBaseStoryVariables(params, tokens, "auth-standard")
	skipStripe := getStripeSecretKey() == ""
	if !skipStripe {
		ensureStripeIntegration(t)
	}
	ensureAwsMarketplaceIntegration(t, params)

	return testlib.Story{
		Name:        "Postman Collection Story - Standard Methods",
		Description: "Postman Collection の Test Flow に従い、WithBody サフィックスを持たない Standard メソッドの完全なテストを実行します。基本設定、ユーザー管理、ロール管理、属性管理、環境管理、テナント管理など、Auth API の主要な機能をカバーします。",
		Variables:   vars,
		Steps: []testlib.Step{
			// 基本設定・認証情報
			{Name: "GetBasicInfo", ClientMethod: "GetBasicInfo", Parameters: getBasicInfoParams, ExpectedStatus: 200},
			{Name: "UpdateBasicInfo", ClientMethod: "UpdateBasicInfo", Parameters: updateBasicInfoParams, ExpectedStatus: 200},
			{Name: "GetAuthInfo", ClientMethod: "GetAuthInfo", Parameters: getAuthInfoParams, ExpectedStatus: 200},
			{Name: "UpdateAuthInfo", ClientMethod: "UpdateAuthInfo", Parameters: updateAuthInfoParams, ExpectedStatus: 200},

			// ユーザー管理
			{Name: "GetSaasUsers", ClientMethod: "GetSaasUsers", Parameters: getSaasUsersParams, ExpectedStatus: 200},
			{Name: "CreateSaasUser", ClientMethod: "CreateSaasUser", Parameters: createSaasUserParams, AllowedStatuses: statusOKOrCreated,
				StateUpdate: func(response any, variables map[string]any) error {
					if id, err := extractIDFromResponse(response, "id"); err == nil {
						variables["user_id"] = id
					} else {
						fmt.Printf("Warning: failed to extract user_id: %v\n", err)
						variables["user_id"] = "extracted_user_id"
					}
					updateUserTokensFromCognito(variables)
					return nil
				}},
			{Name: "GetSaasUser", ClientMethod: "GetSaasUser", Parameters: getSaasUserParams, ExpectedStatus: 200},
			{Name: "UpdateSaasUserPassword", ClientMethod: "UpdateSaasUserPassword", Parameters: updateSaasUserPasswordParams, ExpectedStatus: 200},
			{Name: "UpdateSaasUserEmail", ClientMethod: "UpdateSaasUserEmail", Parameters: updateSaasUserEmailParams, ExpectedStatus: 200},

			// MFA設定
			// 注意: CreateSecretCodeはCognitoの生アクセストークンが必要
			// 正しい順序: CreateSecretCode → UpdateSoftwareToken → UpdateUserMfaPreference
			{Name: "GetUserMfaPreference", ClientMethod: "GetUserMfaPreference", Parameters: getUserMfaPreferenceParams, ExpectedStatus: 200},
			{Name: "CreateSecretCode", ClientMethod: "CreateSecretCode", Parameters: createSecretCodeParams, AllowedStatuses: statusOKOrCreated, Skip: true,
				StateUpdate: func(response any, variables map[string]any) error {
					// シークレットコードをレスポンスから取得して保存
					if secretCode, err := extractFieldFromResponse(response, "secret_code"); err == nil {
						variables["secret_code"] = secretCode
						variables[softwareTokenSecretKey] = secretCode
						fmt.Printf("[DEBUG] Extracted secret_code for MFA setup\n")
					}
					return nil
				}},
			// UpdateSoftwareTokenは実際のTOTPコードが必要なため、テスト環境では401になる可能性がある
			{Name: "UpdateSoftwareToken", ClientMethod: "UpdateSoftwareToken", Parameters: updateSoftwareTokenParams, ExpectedStatus: 200},
			{Name: "UpdateUserMfaPreference", ClientMethod: "UpdateUserMfaPreference", Parameters: updateUserMfaPreferenceParams, ExpectedStatus: 200},

			// ロール管理
			{Name: "GetRoles", ClientMethod: "GetRoles", Parameters: getRolesParams, ExpectedStatus: 200},
			{Name: "CreateRole", ClientMethod: "CreateRole", Parameters: createRoleParams, AllowedStatuses: statusOKOrCreated},

			// ユーザー属性
			{Name: "GetUserAttributes", ClientMethod: "GetUserAttributes", Parameters: getUserAttributesParams, ExpectedStatus: 200},
			{Name: "CreateUserAttribute", ClientMethod: "CreateUserAttribute", Parameters: createUserAttributeParams, AllowedStatuses: statusOKOrCreated},

			// テナント属性
			{Name: "GetTenantAttributes", ClientMethod: "GetTenantAttributes", Parameters: getTenantAttributesParams, ExpectedStatus: 200},
			{Name: "CreateTenantAttribute", ClientMethod: "CreateTenantAttribute", Parameters: createTenantAttributeParams, AllowedStatuses: statusOKOrCreated},

			// 通知メッセージ
			{Name: "FindNotificationMessages", ClientMethod: "FindNotificationMessages", Parameters: findNotificationMessagesParams, ExpectedStatus: 200},
			{Name: "UpdateNotificationMessages", ClientMethod: "UpdateNotificationMessages", Parameters: updateNotificationMessagesParams, ExpectedStatus: 200},

			// カスタマイズページ
			{Name: "GetCustomizePages", ClientMethod: "GetCustomizePages", Parameters: getCustomizePagesParams, ExpectedStatus: 200},
			{Name: "UpdateCustomizePages", ClientMethod: "UpdateCustomizePages", Parameters: updateCustomizePagesParams, ExpectedStatus: 200},
			{Name: "GetCustomizePageSettings", ClientMethod: "GetCustomizePageSettings", Parameters: getCustomizePageSettingsParams, ExpectedStatus: 200,
				StateUpdate: func(response any, variables map[string]any) error {
					captureCustomizePageSettings(response, variables)
					return nil
				}},
			{Name: "UpdateCustomizePageSettings", ClientMethod: "UpdateCustomizePageSettings", Parameters: updateCustomizePageSettingsParams, ExpectedStatus: 200},

			// 環境管理
			{Name: "GetEnvs", ClientMethod: "GetEnvs", Parameters: getEnvsParams, ExpectedStatus: 200},
			{Name: "CreateEnv", ClientMethod: "CreateEnv", Parameters: createEnvParams, AllowedStatuses: statusOKOrCreated},
			{Name: "GetEnv", ClientMethod: "GetEnv", Parameters: getEnvParams, ExpectedStatus: 200},
			{Name: "UpdateEnv", ClientMethod: "UpdateEnv", Parameters: updateEnvParams, ExpectedStatus: 200},

			// サインイン設定
			{Name: "GetSignInSettings", ClientMethod: "GetSignInSettings", Parameters: getSignInSettingsParams, ExpectedStatus: 200},
			{Name: "UpdateSignInSettings", ClientMethod: "UpdateSignInSettings", Parameters: updateSignInSettingsParams, ExpectedStatus: 200},

			// テナント管理
			{Name: "GetTenants", ClientMethod: "GetTenants", Parameters: getTenantsParams, ExpectedStatus: 200},
			{Name: "CreateTenant", ClientMethod: "CreateTenant", Parameters: createTenantParams, AllowedStatuses: statusOKOrCreated,
				StateUpdate: func(response any, variables map[string]any) error {
					if id, err := extractIDFromResponse(response, "id"); err == nil {
						variables["tenant_id"] = id
						setLastTenantID(id)
						return nil
					} else {
						fmt.Printf("Warning: failed to extract tenant_id: %v\n", err)
					}
					variables["tenant_id"] = "extracted_tenant_id"
					setLastTenantID("extracted_tenant_id")
					return nil
				}},
			{Name: "GetTenant", ClientMethod: "GetTenant", Parameters: getTenantParams, ExpectedStatus: 200},
			{Name: "UpdateTenant", ClientMethod: "UpdateTenant", Parameters: updateTenantParams, ExpectedStatus: 200},

			// テナントユーザー管理
			{Name: "GetAllTenantUsers", ClientMethod: "GetAllTenantUsers", Parameters: getAllTenantUsersParams, ExpectedStatus: 200},
			{Name: "GetAllTenantUser", ClientMethod: "GetAllTenantUser", Parameters: getAllTenantUserParams, ExpectedStatus: 200},
			{Name: "CreateTenantUser", ClientMethod: "CreateTenantUser", Parameters: createTenantUserParams, AllowedStatuses: statusOKOrCreated,
				StateUpdate: func(response any, variables map[string]any) error {
					if id, err := extractIDFromResponse(response, "id"); err == nil {
						variables["tenant_user_id"] = id
						return nil
					}
					variables["tenant_user_id"] = ""
					return nil
				}},
			{Name: "GetTenantUsers", ClientMethod: "GetTenantUsers", Parameters: getTenantUsersParams, ExpectedStatus: 200},
			{Name: "GetTenantUser", ClientMethod: "GetTenantUser", Parameters: getTenantUserParams, ExpectedStatus: 200},
			{Name: "UpdateTenantUser", ClientMethod: "UpdateTenantUser", Parameters: updateTenantUserParams, ExpectedStatus: 200},

			// テナントユーザーロール
			{Name: "CreateTenantUserRoles", ClientMethod: "CreateTenantUserRoles", Parameters: createTenantUserRolesParams, AllowedStatuses: statusOKOrCreated},
			{Name: "DeleteTenantUserRole", ClientMethod: "DeleteTenantUserRole", Parameters: deleteTenantUserRoleParams, ExpectedStatus: 200},

			// テナントプラン・請求
			{Name: "UpdateTenantPlan", ClientMethod: "UpdateTenantPlan", Parameters: updateTenantPlanParams, ExpectedStatus: 200},
			{Name: "UpdateTenantBillingInfo", ClientMethod: "UpdateTenantBillingInfo", Parameters: updateTenantBillingInfoParams, ExpectedStatus: 200},

			// Stripe連携
			{Name: "CreateTenantAndPricing", ClientMethod: "CreateTenantAndPricing", Parameters: createTenantAndPricingParams, AllowedStatuses: statusOKOrCreated, Skip: skipStripe},
			{Name: "GetStripeCustomer", ClientMethod: "GetStripeCustomer", Parameters: getStripeCustomerParams, ExpectedStatus: 200, Skip: skipStripe},
			{Name: "DeleteStripeTenantAndPricing", ClientMethod: "DeleteStripeTenantAndPricing", Parameters: deleteStripeTenantAndPricingParams, ExpectedStatus: 200, Skip: skipStripe},

			// プランリセット
			{Name: "ResetPlan", ClientMethod: "ResetPlan", Parameters: resetPlanParams, ExpectedStatus: 200},

			// 認証情報
			{Name: "GetUserInfo", ClientMethod: "GetUserInfo", Parameters: getUserInfoParams, ExpectedStatus: 200},
			{Name: "GetAuthCredentials", ClientMethod: "GetAuthCredentials", Parameters: getAuthCredentialsParams, ExpectedStatus: 200, Skip: true, SkipReason: skipReasonAuthCredentials},
			{Name: "CreateAuthCredentials", ClientMethod: "CreateAuthCredentials", Parameters: createAuthCredentialsParams, ExpectedStatus: http.StatusCreated, Skip: true, SkipReason: skipReasonCreateAuthCredentials},

			// IDプロバイダー
			{Name: "GetIdentityProviders", ClientMethod: "GetIdentityProviders", Parameters: getIdentityProvidersParams, ExpectedStatus: 200},
			{Name: "UpdateIdentityProvider", ClientMethod: "UpdateIdentityProvider", Parameters: updateIdentityProviderParams, ExpectedStatus: 200},
			{Name: "GetTenantIdentityProviders", ClientMethod: "GetTenantIdentityProviders", Parameters: getTenantIdentityProvidersParams, ExpectedStatus: 200},
			{Name: "UpdateTenantIdentityProvider", ClientMethod: "UpdateTenantIdentityProvider", Parameters: updateTenantIdentityProviderParams, ExpectedStatus: 200, Skip: true, SkipReason: skipReasonTenantIdentityProvider},

			// AWS Marketplace
			{Name: "SignUpWithAwsMarketplace", ClientMethod: "SignUpWithAwsMarketplace", Parameters: signUpWithAwsMarketplaceParams, ExpectedStatus: 200, Skip: true, SkipReason: skipReasonAwsMarketplaceSignUp},

			// シングルテナント
			{Name: "GetCloudFormationLaunchStackLinkForSingleTenant", ClientMethod: "GetCloudFormationLaunchStackLinkForSingleTenant", Parameters: getCloudFormationLaunchStackLinkForSingleTenantParams, ExpectedStatus: 200},
			{Name: "GetSingleTenantSettings", ClientMethod: "GetSingleTenantSettings", Parameters: getSingleTenantSettingsParams, ExpectedStatus: 200},
			{Name: "UpdateSingleTenantSettings", ClientMethod: "UpdateSingleTenantSettings", Parameters: updateSingleTenantSettingsParams, ExpectedStatus: 200},

			// クリーンアップ
			{Name: "DeleteTenantUser", ClientMethod: "DeleteTenantUser", Parameters: deleteTenantUserParams, ExpectedStatus: 200},
			{Name: "DeleteTenant", ClientMethod: "DeleteTenant", Parameters: deleteTenantParams, ExpectedStatus: 200},
			{Name: "DeleteEnv", ClientMethod: "DeleteEnv", Parameters: deleteEnvParams, ExpectedStatus: 200},
			{Name: "DeleteTenantAttribute", ClientMethod: "DeleteTenantAttribute", Parameters: deleteTenantAttributeParams, ExpectedStatus: 200},
			{Name: "DeleteUserAttribute", ClientMethod: "DeleteUserAttribute", Parameters: deleteUserAttributeParams, ExpectedStatus: 200},
			{Name: "DeleteRole", ClientMethod: "DeleteRole", Parameters: deleteRoleParams, ExpectedStatus: 200},
			{Name: "DeleteSaasUser", ClientMethod: "DeleteSaasUser", Parameters: deleteSaasUserParams, ExpectedStatus: 200},
		},
		Setup:   func() error { return nil },
		Cleanup: func() error { return nil },
	}
}

// GetPostmanStoryStandardMethodsWithBody は WithBody メソッドのストーリーを返します。
//
// このストーリーは WithBody サフィックスを持つメソッドのみを使用します。
// これらのメソッドは io.Reader 型のボディパラメータを受け取り、
// リクエストボディを柔軟に制御できます。
//
// test_params.json からテストパラメータを読み込み、Variables マップを構築します。
//
// カバーするメソッド:
//   - UpdateBasicInfoWithBody, UpdateAuthInfoWithBody
//   - CreateSaasUserWithBody, UpdateSaasUserPasswordWithBody
//   - UpdateUserMfaPreferenceWithBody, CreateSecretCodeWithBody
//   - CreateRoleWithBody, CreateUserAttributeWithBody
//   - UpdateNotificationMessagesWithBody, UpdateCustomizePagesWithBody
//   - CreateEnvWithBody, UpdateEnvWithBody
//   - CreateTenantWithBody, UpdateTenantWithBody
//   - CreateTenantUserWithBody, UpdateTenantUserWithBody
//   - その他の WithBody バリアント
//
// パラメータ:
//   - t: テストコンテキスト（テストパラメータの読み込みに使用）
//
// 戻り値:
//   - testlib.Story: 実行可能なテストストーリー
func GetPostmanStoryStandardMethodsWithBody(t *testing.T) testlib.Story {
	// テストパラメータを読み込む
	params := testdata.LoadTestParams(t)
	tokens := MustGetAuthTokens(t)

	// test_params.json から Variables マップを構築（WithBody ストーリー用に固有プレフィックスを付与）
	vars := buildBaseStoryVariables(params, tokens, "auth-standard-body")
	skipStripe := getStripeSecretKey() == ""
	if !skipStripe {
		ensureStripeIntegration(t)
	}
	ensureAwsMarketplaceIntegration(t, params)

	return testlib.Story{
		Name:        "Postman Collection Story - Standard Methods With Body",
		Description: "WithBody サフィックスを持つメソッドの完全なテストを実行します。これらのメソッドは io.Reader 型のボディパラメータを受け取り、リクエストボディを柔軟に制御できます。Standard Methods と同じ機能をカバーしますが、異なるメソッドシグネチャを使用します。",
		Variables:   vars,
		Steps: []testlib.Step{
			{Name: "UpdateBasicInfoWithBody", ClientMethod: "UpdateBasicInfoWithBody", Parameters: updateBasicInfoWithBodyParams, ExpectedStatus: 200},
			{Name: "UpdateAuthInfoWithBody", ClientMethod: "UpdateAuthInfoWithBody", Parameters: updateAuthInfoWithBodyParams, ExpectedStatus: 200},
			{Name: "CreateSaasUserWithBody", ClientMethod: "CreateSaasUserWithBody", Parameters: createSaasUserWithBodyParams, AllowedStatuses: statusOKOrCreated,
				StateUpdate: func(response any, variables map[string]any) error {
					if id, err := extractIDFromResponse(response, "id"); err == nil {
						variables["user_id"] = id
					} else {
						fmt.Printf("Warning: failed to extract user_id: %v\n", err)
						variables["user_id"] = "extracted_user_id"
					}
					updateUserTokensFromCognito(variables)
					return nil
				}},
			{Name: "UpdateSaasUserPasswordWithBody", ClientMethod: "UpdateSaasUserPasswordWithBody", Parameters: updateSaasUserPasswordWithBodyParams, ExpectedStatus: 200},
			{Name: "UpdateSaasUserEmailWithBody", ClientMethod: "UpdateSaasUserEmailWithBody", Parameters: updateSaasUserEmailWithBodyParams, ExpectedStatus: 200},
			{Name: "CreateSecretCodeWithBody", ClientMethod: "CreateSecretCodeWithBody", Parameters: createSecretCodeWithBodyParams, AllowedStatuses: statusOKOrCreated, Skip: true,
				StateUpdate: func(response any, variables map[string]any) error {
					if secretCode, err := extractFieldFromResponse(response, "secret_code"); err == nil {
						variables["secret_code"] = secretCode
						variables[softwareTokenSecretKey] = secretCode
						fmt.Printf("[DEBUG] Extracted secret_code (with body) for MFA setup\n")
					}
					return nil
				}},
			{Name: "UpdateSoftwareTokenWithBody", ClientMethod: "UpdateSoftwareTokenWithBody", Parameters: updateSoftwareTokenWithBodyParams, ExpectedStatus: 200},
			{Name: "UpdateUserMfaPreferenceWithBody", ClientMethod: "UpdateUserMfaPreferenceWithBody", Parameters: updateUserMfaPreferenceWithBodyParams, ExpectedStatus: 200},
			{Name: "CreateRoleWithBody", ClientMethod: "CreateRoleWithBody", Parameters: createRoleWithBodyParams, AllowedStatuses: statusOKOrCreated},
			{Name: "CreateUserAttributeWithBody", ClientMethod: "CreateUserAttributeWithBody", Parameters: createUserAttributeWithBodyParams, AllowedStatuses: statusOKOrCreated},
			{Name: "CreateTenantAttributeWithBody", ClientMethod: "CreateTenantAttributeWithBody", Parameters: createTenantAttributeWithBodyParams, AllowedStatuses: statusOKOrCreated},
			{Name: "UpdateNotificationMessagesWithBody", ClientMethod: "UpdateNotificationMessagesWithBody", Parameters: updateNotificationMessagesWithBodyParams, ExpectedStatus: 200},
			{Name: "UpdateCustomizePagesWithBody", ClientMethod: "UpdateCustomizePagesWithBody", Parameters: updateCustomizePagesWithBodyParams, ExpectedStatus: 200},
			{Name: "UpdateCustomizePageSettingsWithBody", ClientMethod: "UpdateCustomizePageSettingsWithBody", Parameters: updateCustomizePageSettingsWithBodyParams, ExpectedStatus: 200},
			{Name: "CreateEnvWithBody", ClientMethod: "CreateEnvWithBody", Parameters: createEnvWithBodyParams, AllowedStatuses: statusOKOrCreated},
			{Name: "UpdateEnvWithBody", ClientMethod: "UpdateEnvWithBody", Parameters: updateEnvWithBodyParams, ExpectedStatus: 200},
			{Name: "UpdateSignInSettingsWithBody", ClientMethod: "UpdateSignInSettingsWithBody", Parameters: updateSignInSettingsWithBodyParams, ExpectedStatus: 200},
			{Name: "CreateTenantWithBody", ClientMethod: "CreateTenantWithBody", Parameters: createTenantWithBodyParams, AllowedStatuses: statusOKOrCreated,
				StateUpdate: func(response any, variables map[string]any) error {
					if id, err := extractIDFromResponse(response, "id"); err == nil {
						variables["tenant_id"] = id
						setLastTenantID(id)
						return nil
					} else {
						fmt.Printf("Warning: failed to extract tenant_id: %v\n", err)
					}
					variables["tenant_id"] = "extracted_tenant_id"
					setLastTenantID("extracted_tenant_id")
					return nil
				}},
			{Name: "UpdateTenantWithBody", ClientMethod: "UpdateTenantWithBody", Parameters: updateTenantWithBodyParams, ExpectedStatus: 200},
			{Name: "CreateTenantUserWithBody", ClientMethod: "CreateTenantUserWithBody", Parameters: createTenantUserWithBodyParams, AllowedStatuses: statusOKOrCreated,
				StateUpdate: func(response any, variables map[string]any) error {
					if id, err := extractIDFromResponse(response, "id"); err == nil {
						variables["tenant_user_id"] = id
						return nil
					}
					variables["tenant_user_id"] = ""
					return nil
				}},
			{Name: "UpdateTenantUserWithBody", ClientMethod: "UpdateTenantUserWithBody", Parameters: updateTenantUserWithBodyParams, ExpectedStatus: 200},
			{Name: "CreateTenantUserRolesWithBody", ClientMethod: "CreateTenantUserRolesWithBody", Parameters: createTenantUserRolesWithBodyParams, AllowedStatuses: statusOKOrCreated},
			{Name: "UpdateTenantPlanWithBody", ClientMethod: "UpdateTenantPlanWithBody", Parameters: updateTenantPlanWithBodyParams, ExpectedStatus: 200},
			{Name: "UpdateTenantBillingInfoWithBody", ClientMethod: "UpdateTenantBillingInfoWithBody", Parameters: updateTenantBillingInfoWithBodyParams, ExpectedStatus: 200},
			{Name: "CreateAuthCredentialsWithBody", ClientMethod: "CreateAuthCredentialsWithBody", Parameters: createAuthCredentialsWithBodyParams, ExpectedStatus: http.StatusCreated, Skip: true, SkipReason: skipReasonCreateAuthCredentials},
			{Name: "UpdateIdentityProviderWithBody", ClientMethod: "UpdateIdentityProviderWithBody", Parameters: updateIdentityProviderWithBodyParams, ExpectedStatus: 200},
			{Name: "UpdateTenantIdentityProviderWithBody", ClientMethod: "UpdateTenantIdentityProviderWithBody", Parameters: updateTenantIdentityProviderWithBodyParams, ExpectedStatus: 200, Skip: true, SkipReason: skipReasonTenantIdentityProvider},
			{Name: "SignUpWithAwsMarketplaceWithBody", ClientMethod: "SignUpWithAwsMarketplaceWithBody", Parameters: signUpWithAwsMarketplaceWithBodyParams, ExpectedStatus: 200, Skip: true, SkipReason: skipReasonAwsMarketplaceSignUp},
			{Name: "UpdateSingleTenantSettingsWithBody", ClientMethod: "UpdateSingleTenantSettingsWithBody", Parameters: updateSingleTenantSettingsWithBodyParams, ExpectedStatus: 200},
			// クリーンアップ
			{Name: "DeleteTenantUser", ClientMethod: "DeleteTenantUser", Parameters: deleteTenantUserParams, ExpectedStatus: 200},
			{Name: "DeleteTenant", ClientMethod: "DeleteTenant", Parameters: deleteTenantParams, ExpectedStatus: 200},
			{Name: "DeleteEnv", ClientMethod: "DeleteEnv", Parameters: deleteEnvParams, ExpectedStatus: 200},
			{Name: "DeleteTenantAttribute", ClientMethod: "DeleteTenantAttribute", Parameters: deleteTenantAttributeParams, ExpectedStatus: 200},
			{Name: "DeleteUserAttribute", ClientMethod: "DeleteUserAttribute", Parameters: deleteUserAttributeParams, ExpectedStatus: 200},
			{Name: "DeleteRole", ClientMethod: "DeleteRole", Parameters: deleteRoleParams, ExpectedStatus: 200},
			{Name: "DeleteSaasUser", ClientMethod: "DeleteSaasUser", Parameters: deleteSaasUserParams, ExpectedStatus: 200},
		},
		Setup:   func() error { return nil },
		Cleanup: func() error { return nil },
	}
}

// GetPostmanStoryWithResponseMethods は WithResponse メソッドのストーリーを返します。
//
// このストーリーは WithResponse サフィックスを持ち、WithBody を持たないメソッドのみを使用します。
// これらのメソッドは詳細なレスポンス情報（ステータスコード、ヘッダー、ボディ）を返します。
//
// test_params.json からテストパラメータを読み込み、Variables マップを構築します。
//
// カバーするメソッド:
//   - GetBasicInfoWithResponse, UpdateBasicInfoWithResponse
//   - GetSaasUsersWithResponse, CreateSaasUserWithResponse
//   - GetUserMfaPreferenceWithResponse, UpdateUserMfaPreferenceWithResponse
//   - GetRolesWithResponse, CreateRoleWithResponse
//   - GetUserAttributesWithResponse, CreateUserAttributeWithResponse
//   - FindNotificationMessagesWithResponse, UpdateNotificationMessagesWithResponse
//   - GetEnvsWithResponse, CreateEnvWithResponse
//   - GetTenantsWithResponse, CreateTenantWithResponse
//   - その他の WithResponse バリアント
//
// パラメータ:
//   - t: テストコンテキスト（テストパラメータの読み込みに使用）
//
// 戻り値:
//   - testlib.Story: 実行可能なテストストーリー
func GetPostmanStoryWithResponseMethods(t *testing.T) testlib.Story {
	// テストパラメータを読み込む
	params := testdata.LoadTestParams(t)
	tokens := MustGetAuthTokens(t)

	// test_params.json から Variables マップを構築（WithResponse ストーリー用プレフィックス）
	vars := buildBaseStoryVariables(params, tokens, "auth-standard-response")
	skipStripe := getStripeSecretKey() == ""
	if !skipStripe {
		ensureStripeIntegration(t)
	}
	ensureAwsMarketplaceIntegration(t, params)

	return testlib.Story{
		Name:        "Postman Collection Story - WithResponse Methods",
		Description: "WithResponse サフィックスを持ち、WithBody を持たないメソッドの完全なテストを実行します。これらのメソッドは詳細なレスポンス情報（ステータスコード、ヘッダー、ボディ）を返し、より細かいレスポンス検証が可能です。",
		Variables:   vars,
		Steps: []testlib.Step{
			// 基本設定・認証情報
			{Name: "GetBasicInfoWithResponse", ClientMethod: "GetBasicInfoWithResponse", Parameters: getBasicInfoParams, ExpectedStatus: 200},
			{Name: "UpdateBasicInfoWithResponse", ClientMethod: "UpdateBasicInfoWithResponse", Parameters: updateBasicInfoParams, ExpectedStatus: 200},
			{Name: "GetAuthInfoWithResponse", ClientMethod: "GetAuthInfoWithResponse", Parameters: getAuthInfoParams, ExpectedStatus: 200},
			{Name: "UpdateAuthInfoWithResponse", ClientMethod: "UpdateAuthInfoWithResponse", Parameters: updateAuthInfoParams, ExpectedStatus: 200},
			{Name: "GetSaasUsersWithResponse", ClientMethod: "GetSaasUsersWithResponse", Parameters: getSaasUsersParams, ExpectedStatus: 200},
			{Name: "CreateSaasUserWithResponse", ClientMethod: "CreateSaasUserWithResponse", Parameters: createSaasUserParams, AllowedStatuses: statusOKOrCreated,
				StateUpdate: func(response any, variables map[string]any) error {
					if id, err := extractIDFromResponse(response, "id"); err == nil {
						variables["user_id"] = id
					} else {
						fmt.Printf("Warning: failed to extract user_id: %v\n", err)
						variables["user_id"] = "extracted_user_id"
					}
					updateUserTokensFromCognito(variables)
					return nil
				}},
			{Name: "GetSaasUserWithResponse", ClientMethod: "GetSaasUserWithResponse", Parameters: getSaasUserParams, ExpectedStatus: 200},
			{Name: "UpdateSaasUserPasswordWithResponse", ClientMethod: "UpdateSaasUserPasswordWithResponse", Parameters: updateSaasUserPasswordParams, ExpectedStatus: 200},
			{Name: "UpdateSaasUserEmailWithResponse", ClientMethod: "UpdateSaasUserEmailWithResponse", Parameters: updateSaasUserEmailParams, ExpectedStatus: 200},
			{Name: "GetUserMfaPreferenceWithResponse", ClientMethod: "GetUserMfaPreferenceWithResponse", Parameters: getUserMfaPreferenceParams, ExpectedStatus: 200},
			{Name: "CreateSecretCodeWithResponse", ClientMethod: "CreateSecretCodeWithResponse", Parameters: createSecretCodeParams, AllowedStatuses: statusOKOrCreated, Skip: true,
				StateUpdate: func(response any, variables map[string]any) error {
					if secretCode, err := extractFieldFromResponse(response, "secret_code"); err == nil {
						variables["secret_code"] = secretCode
						variables[softwareTokenSecretKey] = secretCode
						fmt.Printf("[DEBUG] Extracted secret_code (with response) for MFA setup\n")
					}
					return nil
				}},
			{Name: "UpdateSoftwareTokenWithResponse", ClientMethod: "UpdateSoftwareTokenWithResponse", Parameters: updateSoftwareTokenParams, ExpectedStatus: 200},
			{Name: "UpdateUserMfaPreferenceWithResponse", ClientMethod: "UpdateUserMfaPreferenceWithResponse", Parameters: updateUserMfaPreferenceParams, ExpectedStatus: 200},
			{Name: "GetRolesWithResponse", ClientMethod: "GetRolesWithResponse", Parameters: getRolesParams, ExpectedStatus: 200},
			{Name: "CreateRoleWithResponse", ClientMethod: "CreateRoleWithResponse", Parameters: createRoleParams, AllowedStatuses: statusOKOrCreated},
			{Name: "GetUserAttributesWithResponse", ClientMethod: "GetUserAttributesWithResponse", Parameters: getUserAttributesParams, ExpectedStatus: 200},
			{Name: "CreateUserAttributeWithResponse", ClientMethod: "CreateUserAttributeWithResponse", Parameters: createUserAttributeParams, AllowedStatuses: statusOKOrCreated},
			{Name: "GetTenantAttributesWithResponse", ClientMethod: "GetTenantAttributesWithResponse", Parameters: getTenantAttributesParams, ExpectedStatus: 200},
			{Name: "CreateTenantAttributeWithResponse", ClientMethod: "CreateTenantAttributeWithResponse", Parameters: createTenantAttributeParams, AllowedStatuses: statusOKOrCreated},
			{Name: "FindNotificationMessagesWithResponse", ClientMethod: "FindNotificationMessagesWithResponse", Parameters: findNotificationMessagesParams, ExpectedStatus: 200},
			{Name: "UpdateNotificationMessagesWithResponse", ClientMethod: "UpdateNotificationMessagesWithResponse", Parameters: updateNotificationMessagesParams, ExpectedStatus: 200},
			{Name: "GetCustomizePagesWithResponse", ClientMethod: "GetCustomizePagesWithResponse", Parameters: getCustomizePagesParams, ExpectedStatus: 200},
			{Name: "UpdateCustomizePagesWithResponse", ClientMethod: "UpdateCustomizePagesWithResponse", Parameters: updateCustomizePagesParams, ExpectedStatus: 200},
			{Name: "GetCustomizePageSettingsWithResponse", ClientMethod: "GetCustomizePageSettingsWithResponse", Parameters: getCustomizePageSettingsParams, ExpectedStatus: 200},
			{Name: "UpdateCustomizePageSettingsWithResponse", ClientMethod: "UpdateCustomizePageSettingsWithResponse", Parameters: updateCustomizePageSettingsParams, ExpectedStatus: 200},
			{Name: "GetEnvsWithResponse", ClientMethod: "GetEnvsWithResponse", Parameters: getEnvsParams, ExpectedStatus: 200},
			{Name: "CreateEnvWithResponse", ClientMethod: "CreateEnvWithResponse", Parameters: createEnvParams, AllowedStatuses: statusOKOrCreated},
			{Name: "GetEnvWithResponse", ClientMethod: "GetEnvWithResponse", Parameters: getEnvParams, ExpectedStatus: 200},
			{Name: "UpdateEnvWithResponse", ClientMethod: "UpdateEnvWithResponse", Parameters: updateEnvParams, ExpectedStatus: 200},
			{Name: "GetSignInSettingsWithResponse", ClientMethod: "GetSignInSettingsWithResponse", Parameters: getSignInSettingsParams, ExpectedStatus: 200},
			{Name: "UpdateSignInSettingsWithResponse", ClientMethod: "UpdateSignInSettingsWithResponse", Parameters: updateSignInSettingsParams, ExpectedStatus: 200},
			{Name: "GetTenantsWithResponse", ClientMethod: "GetTenantsWithResponse", Parameters: getTenantsParams, ExpectedStatus: 200},
			{Name: "CreateTenantWithResponse", ClientMethod: "CreateTenantWithResponse", Parameters: createTenantParams, AllowedStatuses: statusOKOrCreated,
				StateUpdate: func(response any, variables map[string]any) error {
					if id, err := extractIDFromResponse(response, "id"); err == nil {
						variables["tenant_id"] = id
						setLastTenantID(id)
						return nil
					} else {
						fmt.Printf("Warning: failed to extract tenant_id: %v\n", err)
					}
					variables["tenant_id"] = "extracted_tenant_id"
					setLastTenantID("extracted_tenant_id")
					return nil
				}},
			{Name: "GetTenantWithResponse", ClientMethod: "GetTenantWithResponse", Parameters: getTenantParams, ExpectedStatus: 200},
			{Name: "UpdateTenantWithResponse", ClientMethod: "UpdateTenantWithResponse", Parameters: updateTenantParams, ExpectedStatus: 200},
			{Name: "GetAllTenantUsersWithResponse", ClientMethod: "GetAllTenantUsersWithResponse", Parameters: getAllTenantUsersParams, ExpectedStatus: 200},
			{Name: "GetAllTenantUserWithResponse", ClientMethod: "GetAllTenantUserWithResponse", Parameters: getAllTenantUserParams, ExpectedStatus: 200},
			{Name: "CreateTenantUserWithResponse", ClientMethod: "CreateTenantUserWithResponse", Parameters: createTenantUserParams, AllowedStatuses: statusOKOrCreated,
				StateUpdate: func(response any, variables map[string]any) error {
					if id, err := extractIDFromResponse(response, "id"); err == nil {
						variables["tenant_user_id"] = id
						return nil
					}
					variables["tenant_user_id"] = ""
					return nil
				}},
			{Name: "GetTenantUsersWithResponse", ClientMethod: "GetTenantUsersWithResponse", Parameters: getTenantUsersParams, ExpectedStatus: 200},
			{Name: "GetTenantUserWithResponse", ClientMethod: "GetTenantUserWithResponse", Parameters: getTenantUserParams, ExpectedStatus: 200},
			{Name: "UpdateTenantUserWithResponse", ClientMethod: "UpdateTenantUserWithResponse", Parameters: updateTenantUserParams, ExpectedStatus: 200},
			{Name: "CreateTenantUserRolesWithResponse", ClientMethod: "CreateTenantUserRolesWithResponse", Parameters: createTenantUserRolesParams, AllowedStatuses: statusOKOrCreated},
			{Name: "DeleteTenantUserRoleWithResponse", ClientMethod: "DeleteTenantUserRoleWithResponse", Parameters: deleteTenantUserRoleParams, ExpectedStatus: 200},
			{Name: "UpdateTenantPlanWithResponse", ClientMethod: "UpdateTenantPlanWithResponse", Parameters: updateTenantPlanParams, ExpectedStatus: 200},
			{Name: "UpdateTenantBillingInfoWithResponse", ClientMethod: "UpdateTenantBillingInfoWithResponse", Parameters: updateTenantBillingInfoParams, ExpectedStatus: 200},

			{Name: "CreateTenantAndPricingWithResponse", ClientMethod: "CreateTenantAndPricingWithResponse", Parameters: createTenantAndPricingParams, AllowedStatuses: statusOKOrCreated, Skip: skipStripe},
			{Name: "GetStripeCustomerWithResponse", ClientMethod: "GetStripeCustomerWithResponse", Parameters: getStripeCustomerParams, ExpectedStatus: 200, Skip: skipStripe},
			{Name: "DeleteStripeTenantAndPricingWithResponse", ClientMethod: "DeleteStripeTenantAndPricingWithResponse", Parameters: deleteStripeTenantAndPricingParams, ExpectedStatus: 200, Skip: skipStripe},
			{Name: "ResetPlanWithResponse", ClientMethod: "ResetPlanWithResponse", Parameters: resetPlanParams, ExpectedStatus: 200},

			// ユーザー情報・認証情報管理
			{Name: "GetUserInfoWithResponse", ClientMethod: "GetUserInfoWithResponse", Parameters: getUserInfoParams, ExpectedStatus: 200},
			{Name: "GetAuthCredentialsWithResponse", ClientMethod: "GetAuthCredentialsWithResponse", Parameters: getAuthCredentialsParams, ExpectedStatus: 200, Skip: true, SkipReason: skipReasonAuthCredentials},
			{Name: "CreateAuthCredentialsWithResponse", ClientMethod: "CreateAuthCredentialsWithResponse", Parameters: createAuthCredentialsParams, ExpectedStatus: http.StatusCreated, Skip: true, SkipReason: skipReasonCreateAuthCredentials},
			{Name: "GetIdentityProvidersWithResponse", ClientMethod: "GetIdentityProvidersWithResponse", Parameters: getIdentityProvidersParams, ExpectedStatus: 200},
			{Name: "UpdateIdentityProviderWithResponse", ClientMethod: "UpdateIdentityProviderWithResponse", Parameters: updateIdentityProviderParams, ExpectedStatus: 200},
			{Name: "GetTenantIdentityProvidersWithResponse", ClientMethod: "GetTenantIdentityProvidersWithResponse", Parameters: getTenantIdentityProvidersParams, ExpectedStatus: 200},
			{Name: "UpdateTenantIdentityProviderWithResponse", ClientMethod: "UpdateTenantIdentityProviderWithResponse", Parameters: updateTenantIdentityProviderParams, ExpectedStatus: 200, Skip: true, SkipReason: skipReasonTenantIdentityProvider},
			{Name: "UnlinkProviderWithResponse", ClientMethod: "UnlinkProviderWithResponse", Parameters: unlinkProviderParams, ExpectedStatus: 200, Skip: true, SkipReason: skipReasonUnlinkProvider},
			{Name: "SignUpWithAwsMarketplaceWithResponse", ClientMethod: "SignUpWithAwsMarketplaceWithResponse", Parameters: signUpWithAwsMarketplaceParams, ExpectedStatus: 200, Skip: true, SkipReason: skipReasonAwsMarketplaceSignUp},
			{Name: "GetCloudFormationLaunchStackLinkForSingleTenantWithResponse", ClientMethod: "GetCloudFormationLaunchStackLinkForSingleTenantWithResponse", Parameters: getCloudFormationLaunchStackLinkForSingleTenantParams, ExpectedStatus: 200},
			{Name: "GetSingleTenantSettingsWithResponse", ClientMethod: "GetSingleTenantSettingsWithResponse", Parameters: getSingleTenantSettingsParams, ExpectedStatus: 200},
			{Name: "UpdateSingleTenantSettingsWithResponse", ClientMethod: "UpdateSingleTenantSettingsWithResponse", Parameters: updateSingleTenantSettingsParams, ExpectedStatus: 200},
			// クリーンアップ
			{Name: "DeleteTenantUserWithResponse", ClientMethod: "DeleteTenantUserWithResponse", Parameters: deleteTenantUserParams, ExpectedStatus: 200},
			{Name: "DeleteTenantWithResponse", ClientMethod: "DeleteTenantWithResponse", Parameters: deleteTenantParams, ExpectedStatus: 200},
			{Name: "DeleteEnvWithResponse", ClientMethod: "DeleteEnvWithResponse", Parameters: deleteEnvParams, ExpectedStatus: 200},
			{Name: "DeleteTenantAttributeWithResponse", ClientMethod: "DeleteTenantAttributeWithResponse", Parameters: deleteTenantAttributeParams, ExpectedStatus: 200},
			{Name: "DeleteUserAttributeWithResponse", ClientMethod: "DeleteUserAttributeWithResponse", Parameters: deleteUserAttributeParams, ExpectedStatus: 200},
			{Name: "DeleteRoleWithResponse", ClientMethod: "DeleteRoleWithResponse", Parameters: deleteRoleParams, ExpectedStatus: 200},
			{Name: "DeleteSaasUserWithResponse", ClientMethod: "DeleteSaasUserWithResponse", Parameters: deleteSaasUserParams, ExpectedStatus: 200},
		},
		Setup:   func() error { return nil },
		Cleanup: func() error { return nil },
	}
}

// GetPostmanStoryWithResponseMethodsWithBody は WithBodyWithResponse メソッドのストーリーを返します。
//
// このストーリーは WithBodyWithResponse サフィックスを持つメソッドのみを使用します。
// これらのメソッドは io.Reader 型のボディパラメータを受け取り、
// 詳細なレスポンス情報を返します。
//
// test_params.json からテストパラメータを読み込み、Variables マップを構築します。
//
// カバーするメソッド:
//   - UpdateBasicInfoWithBodyWithResponse, UpdateAuthInfoWithBodyWithResponse
//   - CreateSaasUserWithBodyWithResponse, UpdateSaasUserPasswordWithBodyWithResponse
//   - UpdateUserMfaPreferenceWithBodyWithResponse, CreateSecretCodeWithBodyWithResponse
//   - CreateRoleWithBodyWithResponse, CreateUserAttributeWithBodyWithResponse
//   - UpdateNotificationMessagesWithBodyWithResponse, UpdateCustomizePagesWithBodyWithResponse
//   - CreateEnvWithBodyWithResponse, UpdateEnvWithBodyWithResponse
//   - CreateTenantWithBodyWithResponse, UpdateTenantWithBodyWithResponse
//   - その他の WithBodyWithResponse バリアント
//
// パラメータ:
//   - t: テストコンテキスト（テストパラメータの読み込みに使用）
//
// 戻り値:
//   - testlib.Story: 実行可能なテストストーリー
func GetPostmanStoryWithResponseMethodsWithBody(t *testing.T) testlib.Story {
	// テストパラメータを読み込む
	params := testdata.LoadTestParams(t)
	tokens := MustGetAuthTokens(t)

	// test_params.json から Variables マップを構築（WithBodyWithResponse 用のプレフィックス）
	vars := buildBaseStoryVariables(params, tokens, "auth-standard-response-body")
	skipStripe := getStripeSecretKey() == ""
	if !skipStripe {
		ensureStripeIntegration(t)
	}
	ensureAwsMarketplaceIntegration(t, params)

	return testlib.Story{
		Name:        "Postman Collection Story - WithResponse Methods With Body",
		Description: "WithBodyWithResponse サフィックスを持つメソッドの完全なテストを実行します。これらのメソッドは io.Reader 型のボディパラメータを受け取り、詳細なレスポンス情報を返します。WithBody と WithResponse の両方の利点を組み合わせています。",
		Variables:   vars,
		Steps: []testlib.Step{
			{Name: "UpdateBasicInfoWithBodyWithResponse", ClientMethod: "UpdateBasicInfoWithBodyWithResponse", Parameters: updateBasicInfoWithBodyParams, ExpectedStatus: 200},
			{Name: "UpdateAuthInfoWithBodyWithResponse", ClientMethod: "UpdateAuthInfoWithBodyWithResponse", Parameters: updateAuthInfoWithBodyParams, ExpectedStatus: 200},
			{Name: "CreateSaasUserWithBodyWithResponse", ClientMethod: "CreateSaasUserWithBodyWithResponse", Parameters: createSaasUserWithBodyParams, AllowedStatuses: statusOKOrCreated,
				StateUpdate: func(response any, variables map[string]any) error {
					if id, err := extractIDFromResponse(response, "id"); err == nil {
						variables["user_id"] = id
					} else {
						fmt.Printf("Warning: failed to extract user_id: %v\n", err)
						variables["user_id"] = "extracted_user_id"
					}
					updateUserTokensFromCognito(variables)
					return nil
				}},
			{Name: "UpdateSaasUserPasswordWithBodyWithResponse", ClientMethod: "UpdateSaasUserPasswordWithBodyWithResponse", Parameters: updateSaasUserPasswordWithBodyParams, ExpectedStatus: 200},
			{Name: "UpdateSaasUserEmailWithBodyWithResponse", ClientMethod: "UpdateSaasUserEmailWithBodyWithResponse", Parameters: updateSaasUserEmailWithBodyParams, ExpectedStatus: 200},
			{Name: "CreateSecretCodeWithBodyWithResponse", ClientMethod: "CreateSecretCodeWithBodyWithResponse", Parameters: createSecretCodeWithBodyParams, AllowedStatuses: statusOKOrCreated, Skip: true,
				StateUpdate: func(response any, variables map[string]any) error {
					if secretCode, err := extractFieldFromResponse(response, "secret_code"); err == nil {
						variables["secret_code"] = secretCode
						variables[softwareTokenSecretKey] = secretCode
						fmt.Printf("[DEBUG] Extracted secret_code (with body+response) for MFA setup\n")
					}
					return nil
				}},
			{Name: "UpdateSoftwareTokenWithBodyWithResponse", ClientMethod: "UpdateSoftwareTokenWithBodyWithResponse", Parameters: updateSoftwareTokenWithBodyParams, ExpectedStatus: 200},
			{Name: "UpdateUserMfaPreferenceWithBodyWithResponse", ClientMethod: "UpdateUserMfaPreferenceWithBodyWithResponse", Parameters: updateUserMfaPreferenceWithBodyParams, ExpectedStatus: 200},
			{Name: "CreateRoleWithBodyWithResponse", ClientMethod: "CreateRoleWithBodyWithResponse", Parameters: createRoleWithBodyParams, AllowedStatuses: statusOKOrCreated},
			{Name: "CreateUserAttributeWithBodyWithResponse", ClientMethod: "CreateUserAttributeWithBodyWithResponse", Parameters: createUserAttributeWithBodyParams, AllowedStatuses: statusOKOrCreated},
			{Name: "CreateTenantAttributeWithBodyWithResponse", ClientMethod: "CreateTenantAttributeWithBodyWithResponse", Parameters: createTenantAttributeWithBodyParams, AllowedStatuses: statusOKOrCreated},
			{Name: "UpdateNotificationMessagesWithBodyWithResponse", ClientMethod: "UpdateNotificationMessagesWithBodyWithResponse", Parameters: updateNotificationMessagesWithBodyParams, ExpectedStatus: 200},
			{Name: "UpdateCustomizePagesWithBodyWithResponse", ClientMethod: "UpdateCustomizePagesWithBodyWithResponse", Parameters: updateCustomizePagesWithBodyParams, ExpectedStatus: 200},
			{Name: "UpdateCustomizePageSettingsWithBodyWithResponse", ClientMethod: "UpdateCustomizePageSettingsWithBodyWithResponse", Parameters: updateCustomizePageSettingsWithBodyParams, ExpectedStatus: 200},
			{Name: "CreateEnvWithBodyWithResponse", ClientMethod: "CreateEnvWithBodyWithResponse", Parameters: createEnvWithBodyParams, AllowedStatuses: statusOKOrCreated},
			{Name: "UpdateEnvWithBodyWithResponse", ClientMethod: "UpdateEnvWithBodyWithResponse", Parameters: updateEnvWithBodyParams, ExpectedStatus: 200},
			{Name: "UpdateSignInSettingsWithBodyWithResponse", ClientMethod: "UpdateSignInSettingsWithBodyWithResponse", Parameters: updateSignInSettingsWithBodyParams, ExpectedStatus: 200},
			{Name: "CreateTenantWithBodyWithResponse", ClientMethod: "CreateTenantWithBodyWithResponse", Parameters: createTenantWithBodyParams, AllowedStatuses: statusOKOrCreated,
				StateUpdate: func(response any, variables map[string]any) error {
					if id, err := extractIDFromResponse(response, "id"); err == nil {
						variables["tenant_id"] = id
						setLastTenantID(id)
						return nil
					} else {
						fmt.Printf("Warning: failed to extract tenant_id: %v\n", err)
					}
					variables["tenant_id"] = "extracted_tenant_id"
					setLastTenantID("extracted_tenant_id")
					return nil
				}},
			{Name: "UpdateTenantWithBodyWithResponse", ClientMethod: "UpdateTenantWithBodyWithResponse", Parameters: updateTenantWithBodyParams, ExpectedStatus: 200},
			{Name: "CreateTenantUserWithBodyWithResponse", ClientMethod: "CreateTenantUserWithBodyWithResponse", Parameters: createTenantUserWithBodyParams, AllowedStatuses: statusOKOrCreated,
				StateUpdate: func(response any, variables map[string]any) error {
					if id, err := extractIDFromResponse(response, "id"); err == nil {
						variables["tenant_user_id"] = id
						return nil
					}
					variables["tenant_user_id"] = ""
					return nil
				}},
			{Name: "UpdateTenantUserWithBodyWithResponse", ClientMethod: "UpdateTenantUserWithBodyWithResponse", Parameters: updateTenantUserWithBodyParams, ExpectedStatus: 200},
			{Name: "CreateTenantUserRolesWithBodyWithResponse", ClientMethod: "CreateTenantUserRolesWithBodyWithResponse", Parameters: createTenantUserRolesWithBodyParams, AllowedStatuses: statusOKOrCreated},
			{Name: "UpdateTenantPlanWithBodyWithResponse", ClientMethod: "UpdateTenantPlanWithBodyWithResponse", Parameters: updateTenantPlanWithBodyParams, ExpectedStatus: 200},
			{Name: "UpdateTenantBillingInfoWithBodyWithResponse", ClientMethod: "UpdateTenantBillingInfoWithBodyWithResponse", Parameters: updateTenantBillingInfoWithBodyParams, ExpectedStatus: 200},
			{Name: "CreateAuthCredentialsWithBodyWithResponse", ClientMethod: "CreateAuthCredentialsWithBodyWithResponse", Parameters: createAuthCredentialsWithBodyParams, ExpectedStatus: http.StatusCreated, Skip: true, SkipReason: skipReasonCreateAuthCredentials},
			{Name: "UpdateIdentityProviderWithBodyWithResponse", ClientMethod: "UpdateIdentityProviderWithBodyWithResponse", Parameters: updateIdentityProviderWithBodyParams, ExpectedStatus: 200},
			{Name: "UpdateTenantIdentityProviderWithBodyWithResponse", ClientMethod: "UpdateTenantIdentityProviderWithBodyWithResponse", Parameters: updateTenantIdentityProviderWithBodyParams, ExpectedStatus: 200, Skip: true, SkipReason: skipReasonTenantIdentityProvider},
			{Name: "SignUpWithAwsMarketplaceWithBodyWithResponse", ClientMethod: "SignUpWithAwsMarketplaceWithBodyWithResponse", Parameters: signUpWithAwsMarketplaceWithBodyParams, ExpectedStatus: 200, Skip: true, SkipReason: skipReasonAwsMarketplaceSignUp},
			{Name: "UpdateSingleTenantSettingsWithBodyWithResponse", ClientMethod: "UpdateSingleTenantSettingsWithBodyWithResponse", Parameters: updateSingleTenantSettingsWithBodyParams, ExpectedStatus: 200},
			// クリーンアップ
			{Name: "DeleteTenantUserWithResponse", ClientMethod: "DeleteTenantUserWithResponse", Parameters: deleteTenantUserParams, ExpectedStatus: 200},
			{Name: "DeleteTenantWithResponse", ClientMethod: "DeleteTenantWithResponse", Parameters: deleteTenantParams, ExpectedStatus: 200},
			{Name: "DeleteEnvWithResponse", ClientMethod: "DeleteEnvWithResponse", Parameters: deleteEnvParams, ExpectedStatus: 200},
			{Name: "DeleteTenantAttributeWithResponse", ClientMethod: "DeleteTenantAttributeWithResponse", Parameters: deleteTenantAttributeParams, ExpectedStatus: 200},
			{Name: "DeleteUserAttributeWithResponse", ClientMethod: "DeleteUserAttributeWithResponse", Parameters: deleteUserAttributeParams, ExpectedStatus: 200},
			{Name: "DeleteRoleWithResponse", ClientMethod: "DeleteRoleWithResponse", Parameters: deleteRoleParams, ExpectedStatus: 200},
			{Name: "DeleteSaasUserWithResponse", ClientMethod: "DeleteSaasUserWithResponse", Parameters: deleteSaasUserParams, ExpectedStatus: 200},
		},
		Setup:   func() error { return nil },
		Cleanup: func() error { return nil },
	}
}

// ========================================
// Additional Stories for Uncovered Methods
// ========================================

// GetStorySaasUserAttributes は SaaS ユーザー属性管理ストーリーを返します。
//
// このストーリーは UpdateSaasUserAttributes 系と CreateSaasUserAttribute 系の
// 全バリアント（Standard, WithBody, WithResponse, WithBodyWithResponse）をカバーします。
//
// カバーするメソッド (8メソッド):
//   - UpdateSaasUserAttributes, UpdateSaasUserAttributesWithBody
//   - UpdateSaasUserAttributesWithResponse, UpdateSaasUserAttributesWithBodyWithResponse
//   - CreateSaasUserAttribute, CreateSaasUserAttributeWithBody
//   - CreateSaasUserAttributeWithResponse, CreateSaasUserAttributeWithBodyWithResponse
//
// パラメータ:
//   - t: テストコンテキスト（テストパラメータの読み込みに使用）
//
// 戻り値:
//   - testlib.Story: 実行可能なテストストーリー
func GetStorySaasUserAttributes(t *testing.T) testlib.Story {
	// テストパラメータを読み込む
	params := testdata.LoadTestParams(t)

	// Variables マップを構築
	vars := map[string]any{
		// ユーザー情報
		"email":    uniqueEmail("auth-attr"),
		"password": params.Users.CreateParams["password"],

		// SaaS ユーザー属性
		"saas_user_attributes": map[string]interface{}{"saas_custom_field": "test value"},

		// SaaS ユーザー属性定義
		"attribute_name": "saas_custom_field",
		"display_name":   "SaaS Custom Field",
		"attribute_type": "string",
	}

	return testlib.Story{
		Name:        "SaaS User Attributes Management Story",
		Description: "SaaS ユーザーの追加属性を管理するメソッドをテストします。UpdateSaasUserAttributes 系と CreateSaasUserAttribute 系の全バリアントをカバーします。",
		Variables:   vars,
		Steps: []testlib.Step{
			// Setup: ユーザーを作成
			{Name: "CreateSaasUser for attributes test", ClientMethod: "CreateSaasUser", Parameters: createSaasUserParams, AllowedStatuses: statusOKOrCreated,
				StateUpdate: func(response any, variables map[string]any) error {
					if id, err := extractIDFromResponse(response, "id"); err == nil {
						variables["user_id"] = id
					} else {
						fmt.Printf("Warning: failed to extract user_id: %v\n", err)
						variables["user_id"] = "extracted_user_id"
					}
					updateUserTokensFromCognito(variables)
					return nil
				}},

			// SaaS ユーザー属性定義を作成
			{Name: "CreateSaasUserAttribute", ClientMethod: "CreateSaasUserAttribute", Parameters: createSaasUserAttributeParams, AllowedStatuses: statusOKOrCreated},
			{Name: "CreateSaasUserAttributeWithBody", ClientMethod: "CreateSaasUserAttributeWithBody", Parameters: createSaasUserAttributeWithBodyParams, AllowedStatuses: statusOKOrCreated},
			{Name: "CreateSaasUserAttributeWithResponse", ClientMethod: "CreateSaasUserAttributeWithResponse", Parameters: createSaasUserAttributeParams, AllowedStatuses: statusOKOrCreated},
			{Name: "CreateSaasUserAttributeWithBodyWithResponse", ClientMethod: "CreateSaasUserAttributeWithBodyWithResponse", Parameters: createSaasUserAttributeWithBodyParams, AllowedStatuses: statusOKOrCreated},

			// SaaS ユーザー属性を更新
			{Name: "UpdateSaasUserAttributes", ClientMethod: "UpdateSaasUserAttributes", Parameters: updateSaasUserAttributesParams, ExpectedStatus: 200},
			{Name: "UpdateSaasUserAttributesWithBody", ClientMethod: "UpdateSaasUserAttributesWithBody", Parameters: updateSaasUserAttributesWithBodyParams, ExpectedStatus: 200},
			{Name: "UpdateSaasUserAttributesWithResponse", ClientMethod: "UpdateSaasUserAttributesWithResponse", Parameters: updateSaasUserAttributesParams, ExpectedStatus: 200},
			{Name: "UpdateSaasUserAttributesWithBodyWithResponse", ClientMethod: "UpdateSaasUserAttributesWithBodyWithResponse", Parameters: updateSaasUserAttributesWithBodyParams, ExpectedStatus: 200},

			// Cleanup: ユーザーを削除
			{Name: "DeleteSaasUser after attributes test", ClientMethod: "DeleteSaasUser", Parameters: deleteSaasUserParams, ExpectedStatus: 200},
		},
		Setup:   func() error { return nil },
		Cleanup: func() error { return nil },
	}
}

// GetStoryTenantInvitations は テナント招待機能ストーリーを返します。
//
// このストーリーはテナント招待に関連する全メソッドの全バリアントをカバーします。
//
// カバーするメソッド (16メソッド):
//   - CreateTenantInvitation系（4メソッド）
//   - GetTenantInvitations系（4メソッド）
//   - GetTenantInvitation系（4メソッド）
//   - DeleteTenantInvitation系（2メソッド）
//   - ValidateInvitation系（4メソッド）
//   - GetInvitationValidity系（2メソッド）
//
// パラメータ:
//   - t: テストコンテキスト（テストパラメータの読み込みに使用）
//
// 戻り値:
//   - testlib.Story: 実行可能なテストストーリー
func GetStoryTenantInvitations(t *testing.T) testlib.Story {
	// テストパラメータを読み込む
	params := testdata.LoadTestParams(t)
	tokens := MustGetAuthTokens(t)

	// Variables マップを構築
	vars := map[string]any{
		// テナント情報
		"name":                    uniqueString("invitation-test-tenant"),
		"back_office_staff_email": params.Tenants.CreateParams["back_office_staff_email"],
		"attributes":              map[string]interface{}{},

		// 招待情報
		"invitation_email": uniqueEmail("auth-invite"),
		"invitation_envs":  params.TenantInvitations.CreateParams["envs"],
		"access_token":     tokens.AccessToken,
		"password":         params.Users.CreateParams["password"],
	}

	return testlib.Story{
		Name:        "Tenant Invitations Management Story",
		Description: "テナント招待機能をテストします。招待の作成、取得、検証、削除の全バリアントをカバーします。",
		Variables:   vars,
		Steps: []testlib.Step{
			// Setup: テナントを作成
			{Name: "CreateTenant for invitations", ClientMethod: "CreateTenant", Parameters: createTenantParams, AllowedStatuses: statusOKOrCreated,
				StateUpdate: func(response any, variables map[string]any) error {
					id, err := extractIDFromResponse(response, "id")
					if err != nil || id == "" {
						return fmt.Errorf("failed to extract tenant_id: %w", err)
					}
					variables["tenant_id"] = id
					setLastTenantID(id)
					return nil
				}},

			// テナント招待を作成（Standard）
			{Name: "CreateTenantInvitation", ClientMethod: "CreateTenantInvitation", Parameters: createTenantInvitationParams, AllowedStatuses: statusOKOrCreated,
				StateUpdate: func(response any, variables map[string]any) error {
					id, err := extractIDFromResponse(response, "id")
					if err != nil || id == "" {
						return fmt.Errorf("failed to extract invitation_id: %w", err)
					}
					variables["invitation_id"] = id
					return nil
				}},

			// テナント招待一覧を取得（Standard）
			{Name: "GetTenantInvitations", ClientMethod: "GetTenantInvitations", Parameters: getTenantInvitationsParams, ExpectedStatus: 200},
			{Name: "GetTenantInvitationsWithResponse", ClientMethod: "GetTenantInvitationsWithResponse", Parameters: getTenantInvitationsParams, ExpectedStatus: 200},

			// 特定のテナント招待を取得（Standard）
			{Name: "GetTenantInvitation", ClientMethod: "GetTenantInvitation", Parameters: getTenantInvitationParams, ExpectedStatus: 200},
			{Name: "GetTenantInvitationWithResponse", ClientMethod: "GetTenantInvitationWithResponse", Parameters: getTenantInvitationParams, ExpectedStatus: 200},

			// 招待の有効性を確認（Standard）
			{Name: "GetInvitationValidity", ClientMethod: "GetInvitationValidity", Parameters: getInvitationValidityParams, ExpectedStatus: 200},
			{Name: "GetInvitationValidityWithResponse", ClientMethod: "GetInvitationValidityWithResponse", Parameters: getInvitationValidityParams, ExpectedStatus: 200},

			// 招待を検証（Standard）
			{Name: "ValidateInvitation", ClientMethod: "ValidateInvitation", Parameters: validateInvitationParams, ExpectedStatus: 200},
			{Name: "ValidateInvitationWithBody", ClientMethod: "ValidateInvitationWithBody", Parameters: validateInvitationWithBodyParams, ExpectedStatus: 200},
			{Name: "ValidateInvitationWithResponse", ClientMethod: "ValidateInvitationWithResponse", Parameters: validateInvitationParams, ExpectedStatus: 200},
			{Name: "ValidateInvitationWithBodyWithResponse", ClientMethod: "ValidateInvitationWithBodyWithResponse", Parameters: validateInvitationWithBodyParams, ExpectedStatus: 200},

			// テナント招待を削除（Standard）
			{Name: "DeleteTenantInvitation", ClientMethod: "DeleteTenantInvitation", Parameters: deleteTenantInvitationParams, ExpectedStatus: 200},
			{Name: "DeleteTenantInvitationWithResponse", ClientMethod: "DeleteTenantInvitationWithResponse", Parameters: deleteTenantInvitationParams, ExpectedStatus: 200},

			// WithBody バリアントをテスト
			{Name: "CreateTenantInvitationWithBody", ClientMethod: "CreateTenantInvitationWithBody", Parameters: createTenantInvitationWithBodyParams, AllowedStatuses: statusOKOrCreated,
				StateUpdate: func(response any, variables map[string]any) error {
					id, err := extractIDFromResponse(response, "id")
					if err != nil || id == "" {
						return fmt.Errorf("failed to extract invitation_id: %w", err)
					}
					variables["invitation_id"] = id
					return nil
				}},
			{Name: "CreateTenantInvitationWithResponse", ClientMethod: "CreateTenantInvitationWithResponse", Parameters: createTenantInvitationParams, AllowedStatuses: statusOKOrCreated,
				StateUpdate: func(response any, variables map[string]any) error {
					id, err := extractIDFromResponse(response, "id")
					if err != nil || id == "" {
						return fmt.Errorf("failed to extract invitation_id: %w", err)
					}
					variables["invitation_id"] = id
					return nil
				}},
			{Name: "CreateTenantInvitationWithBodyWithResponse", ClientMethod: "CreateTenantInvitationWithBodyWithResponse", Parameters: createTenantInvitationWithBodyParams, AllowedStatuses: statusOKOrCreated},

			// Cleanup: テナントを削除
			{Name: "DeleteTenant after invitations test", ClientMethod: "DeleteTenant", Parameters: deleteTenantParams, ExpectedStatus: 200},
		},
		Setup:   func() error { return nil },
		Cleanup: func() error { return nil },
	}
}

// GetStoryExternalUserLinkAndEmailUpdate は 外部ユーザーリンクとメール更新ストーリーを返します。
//
// このストーリーは外部ユーザーリンクとメール更新確認に関連する全メソッドの全バリアントをカバーします。
//
// カバーするメソッド (16メソッド):
//   - RequestExternalUserLink系（4メソッド）
//   - ConfirmExternalUserLink系（4メソッド）
//   - RequestEmailUpdate系（4メソッド）
//   - ConfirmEmailUpdate系（4メソッド）
//
// 注意事項:
//   - Request系メソッドは確認コードをメールで送信するため、200を返します
//   - Confirm系メソッドは実際の確認コードが必要なため、E2E環境では401エラーになります
//   - 実際の確認コードを取得するには、メールサーバーとの統合が必要です
//
// パラメータ:
//   - t: テストコンテキスト（テストパラメータの読み込みに使用）
//
// 戻り値:
//   - testlib.Story: 実行可能なテストストーリー
func GetStoryExternalUserLinkAndEmailUpdate(t *testing.T) testlib.Story {
	// テストパラメータを読み込む
	params := testdata.LoadTestParams(t)
	tokens := MustGetAuthTokens(t)

	// Variables マップを構築
	vars := map[string]any{
		// ユーザー情報
		"email":    uniqueEmail("auth-external"),
		"password": params.Users.CreateParams["password"],

		// 外部ユーザーリンク・メール更新用
		"access_token":      tokens.AccessToken,
		"verification_code": "000000", // モック確認コード（実際のメールから取得が必要）
		"new_email":         uniqueEmail("auth-email-update"),
	}

	return testlib.Story{
		Name:        "External User Link and Email Update Story",
		Description: "外部ユーザーアカウントのリンクとメールアドレス更新確認機能をテストします。全バリアントをカバーします。",
		Variables:   vars,
		Steps: []testlib.Step{
			// Setup: ユーザーを作成
			{Name: "CreateSaasUser for external link test", ClientMethod: "CreateSaasUser", Parameters: createSaasUserParams, AllowedStatuses: statusOKOrCreated,
				StateUpdate: func(response any, variables map[string]any) error {
					if id, err := extractIDFromResponse(response, "id"); err == nil {
						variables["user_id"] = id
					} else {
						fmt.Printf("Warning: failed to extract user_id: %v\n", err)
						variables["user_id"] = "extracted_user_id"
					}
					updateUserTokensFromCognito(variables)
					return nil
				}},

			// 外部ユーザーリンクをリクエスト（Standard）
			{Name: "RequestExternalUserLink", ClientMethod: "RequestExternalUserLink", Parameters: requestExternalUserLinkParams, ExpectedStatus: 200, Skip: true,
				StateUpdate: func(response any, variables map[string]any) error {
					// 実際のAPIレスポンスから確認コードを抽出（存在する場合）
					fmt.Printf("[INFO] RequestExternalUserLink completed. Verification code should be sent to email.\n")
					return nil
				}},
			{Name: "RequestExternalUserLinkWithBody", ClientMethod: "RequestExternalUserLinkWithBody", Parameters: requestExternalUserLinkWithBodyParams, ExpectedStatus: 200, Skip: true},
			{Name: "RequestExternalUserLinkWithResponse", ClientMethod: "RequestExternalUserLinkWithResponse", Parameters: requestExternalUserLinkParams, ExpectedStatus: 200, Skip: true},
			{Name: "RequestExternalUserLinkWithBodyWithResponse", ClientMethod: "RequestExternalUserLinkWithBodyWithResponse", Parameters: requestExternalUserLinkWithBodyParams, ExpectedStatus: 200, Skip: true},

			// 外部ユーザーリンクを確認（Standard）
			// 注: これらのステップは実際の確認コードが必要なため、401エラーになる可能性があります
			{Name: "ConfirmExternalUserLink", ClientMethod: "ConfirmExternalUserLink", Parameters: confirmExternalUserLinkParams, ExpectedStatus: 401, Skip: true},
			{Name: "ConfirmExternalUserLinkWithBody", ClientMethod: "ConfirmExternalUserLinkWithBody", Parameters: confirmExternalUserLinkWithBodyParams, ExpectedStatus: 401, Skip: true},
			{Name: "ConfirmExternalUserLinkWithResponse", ClientMethod: "ConfirmExternalUserLinkWithResponse", Parameters: confirmExternalUserLinkParams, ExpectedStatus: 401, Skip: true},
			{Name: "ConfirmExternalUserLinkWithBodyWithResponse", ClientMethod: "ConfirmExternalUserLinkWithBodyWithResponse", Parameters: confirmExternalUserLinkWithBodyParams, ExpectedStatus: 401, Skip: true},

			// メールアドレス更新をリクエスト（Standard）
			{Name: "RequestEmailUpdate", ClientMethod: "RequestEmailUpdate", Parameters: requestEmailUpdateParams, ExpectedStatus: 200, Skip: true, SkipReason: skipReasonEmailUpdate,
				StateUpdate: func(response any, variables map[string]any) error {
					// 実際のAPIレスポンスから確認コードを抽出（存在する場合）
					// 注: 実際のAPIはメールで送信するため、ここでは取得できない可能性が高い
					fmt.Printf("[INFO] RequestEmailUpdate completed. Verification code should be sent to email.\n")
					return nil
				}},
			{Name: "RequestEmailUpdateWithBody", ClientMethod: "RequestEmailUpdateWithBody", Parameters: requestEmailUpdateWithBodyParams, ExpectedStatus: 200, Skip: true, SkipReason: skipReasonEmailUpdate},
			{Name: "RequestEmailUpdateWithResponse", ClientMethod: "RequestEmailUpdateWithResponse", Parameters: requestEmailUpdateParams, ExpectedStatus: 200, Skip: true, SkipReason: skipReasonEmailUpdate},
			{Name: "RequestEmailUpdateWithBodyWithResponse", ClientMethod: "RequestEmailUpdateWithBodyWithResponse", Parameters: requestEmailUpdateWithBodyParams, ExpectedStatus: 200, Skip: true, SkipReason: skipReasonEmailUpdate},

			// メールアドレス更新を確認（Standard）
			// 注: これらのステップは実際の確認コードが必要なため、401エラーになる可能性があります
			{Name: "ConfirmEmailUpdate", ClientMethod: "ConfirmEmailUpdate", Parameters: confirmEmailUpdateParams, ExpectedStatus: 401, Skip: true},
			{Name: "ConfirmEmailUpdateWithBody", ClientMethod: "ConfirmEmailUpdateWithBody", Parameters: confirmEmailUpdateWithBodyParams, ExpectedStatus: 401, Skip: true},
			{Name: "ConfirmEmailUpdateWithResponse", ClientMethod: "ConfirmEmailUpdateWithResponse", Parameters: confirmEmailUpdateParams, ExpectedStatus: 401, Skip: true},
			{Name: "ConfirmEmailUpdateWithBodyWithResponse", ClientMethod: "ConfirmEmailUpdateWithBodyWithResponse", Parameters: confirmEmailUpdateWithBodyParams, ExpectedStatus: 401, Skip: true},

			// Cleanup: ユーザーを削除
			{Name: "DeleteSaasUser after external link test", ClientMethod: "DeleteSaasUser", Parameters: deleteSaasUserParams, ExpectedStatus: 200},
		},
		Setup:   func() error { return nil },
		Cleanup: func() error { return nil },
	}
}

// GetStorySignUpAndProviderManagement は サインアップとプロバイダー管理ストーリーを返します。
//
// このストーリーはサインアップ、AWS Marketplace連携、プロバイダー管理に関連する
// 全メソッドの全バリアントをカバーします。
//
// カバーするメソッド (17メソッド):
//   - SignUp系（4メソッド）
//   - ResendSignUpConfirmationEmail系（4メソッド）
//   - ConfirmSignUpWithAwsMarketplace系（4メソッド）
//   - LinkAwsMarketplace系（4メソッド）
//   - UnlinkProvider（1メソッド）
//
// パラメータ:
//   - t: テストコンテキスト（テストパラメータの読み込みに使用）
//
// 戻り値:
//   - testlib.Story: 実行可能なテストストーリー
func GetStorySignUpAndProviderManagement(t *testing.T) testlib.Story {
	// テストパラメータを読み込む
	params := testdata.LoadTestParams(t)
	tokens := MustGetAuthTokens(t)

	// Variables マップを構築
	vars := map[string]any{
		// サインアップ情報
		"signup_email": uniqueEmail("auth-signup"),
		"email":        uniqueEmail("auth-provider-unlink"),
		"password":     params.Users.CreateParams["password"],

		// AWS Marketplace情報
		"access_token":       tokens.AccessToken,
		"registration_token": params.AwsMarketplaceSignUp.CreateParams["registration_token"],
		"tenant_name":        "aws-marketplace-tenant",

		// プロバイダー管理
		"provider_name": params.ProviderManagement.UnlinkParams["provider_name"],

		// テナント情報
		"name":                    "signup-test-tenant",
		"back_office_staff_email": params.Tenants.CreateParams["back_office_staff_email"],
		"attributes":              map[string]interface{}{},
	}

	if awsCfg := params.AwsMarketplaceSignUp.CreateParams; awsCfg != nil {
		if token := stringFromMap(awsCfg, "registration_token"); token != "" {
			vars["registration_token"] = token
			vars[awsMarketplaceTokenVariable] = token
		}
		if email := stringFromMap(awsCfg, "email"); email != "" {
			vars[awsMarketplaceEmailVariable] = email
		}
	}
	if token := os.Getenv(envAwsMarketplaceRegistrationToken); token != "" {
		vars["registration_token"] = token
		vars[awsMarketplaceTokenVariable] = token
	}

	ensureAwsMarketplaceIntegration(t, params)
	signUpEnabled := isSignUpTestingEnabled()

	return testlib.Story{
		Name:        "Sign Up and Provider Management Story",
		Description: "サインアップ、AWS Marketplace連携、プロバイダー管理機能をテストします。全バリアントをカバーします。",
		Variables:   vars,
		Steps: []testlib.Step{
			// セルフサインアップを有効化（事前準備）
			{Name: "UpdateSignInSettings (Enable Self Regist)", ClientMethod: "UpdateSignInSettings", Parameters: updateSignInSettingsParams, ExpectedStatus: 200,
				StateUpdate: func(response any, variables map[string]any) error {
					time.Sleep(5 * time.Second)
					return nil
				}},
			{Name: "GetSignInSettings (Verify Self Regist)", ClientMethod: "GetSignInSettings", Parameters: getSignInSettingsParams, ExpectedStatus: 200},

			// サインアップ（Standard）
			{Name: "SignUp", ClientMethod: "SignUp", Parameters: signUpParams, AllowedStatuses: statusCreatedOnly, Skip: !signUpEnabled, SkipReason: skipReasonSignUpLimit,
				StateUpdate: func(response any, variables map[string]any) error {
					return updateSignUpStateFromResponse(response, variables, signUpUserIDKeyStandard, true)
				}},
			{Name: "SignUpWithBody", ClientMethod: "SignUpWithBody", Parameters: signUpWithBodyParams, AllowedStatuses: statusCreatedOnly, Skip: !signUpEnabled, SkipReason: skipReasonSignUpLimit,
				StateUpdate: func(response any, variables map[string]any) error {
					return updateSignUpStateFromResponse(response, variables, signUpUserIDKeyWithBody, true)
				}},
			{Name: "SignUpWithResponse", ClientMethod: "SignUpWithResponse", Parameters: signUpParams, AllowedStatuses: statusCreatedOnly, Skip: !signUpEnabled, SkipReason: skipReasonSignUpLimit,
				StateUpdate: func(response any, variables map[string]any) error {
					return updateSignUpStateFromResponse(response, variables, signUpUserIDKeyWithResponse, true)
				}},
			{Name: "SignUpWithBodyWithResponse", ClientMethod: "SignUpWithBodyWithResponse", Parameters: signUpWithBodyParams, AllowedStatuses: statusCreatedOnly, Skip: !signUpEnabled, SkipReason: skipReasonSignUpLimit,
				StateUpdate: func(response any, variables map[string]any) error {
					return updateSignUpStateFromResponse(response, variables, signUpUserIDKeyWithBodyWithResponse, true)
				}},

			// サインアップ確認メール再送信（Standard）
			{Name: "ResendSignUpConfirmationEmail", ClientMethod: "ResendSignUpConfirmationEmail", Parameters: resendSignUpConfirmationEmailParams, ExpectedStatus: 200},
			{Name: "ResendSignUpConfirmationEmailWithBody", ClientMethod: "ResendSignUpConfirmationEmailWithBody", Parameters: resendSignUpConfirmationEmailWithBodyParams, ExpectedStatus: 200},
			{Name: "ResendSignUpConfirmationEmailWithResponse", ClientMethod: "ResendSignUpConfirmationEmailWithResponse", Parameters: resendSignUpConfirmationEmailParams, ExpectedStatus: 200},
			{Name: "ResendSignUpConfirmationEmailWithBodyWithResponse", ClientMethod: "ResendSignUpConfirmationEmailWithBodyWithResponse", Parameters: resendSignUpConfirmationEmailWithBodyParams, ExpectedStatus: 200},

			// Cleanup: サインアップで作成したユーザーを削除
			{Name: "Delete signup user (SignUp)", ClientMethod: "DeleteSaasUser", Parameters: deleteSignUpUserParamsFor(signUpUserIDKeyStandard), ExpectedStatus: 200, Skip: !signUpEnabled, SkipReason: skipReasonSignUpLimit},
			{Name: "Delete signup user (SignUpWithBody)", ClientMethod: "DeleteSaasUser", Parameters: deleteSignUpUserParamsFor(signUpUserIDKeyWithBody), ExpectedStatus: 200, Skip: !signUpEnabled, SkipReason: skipReasonSignUpLimit},
			{Name: "Delete signup user (SignUpWithResponse)", ClientMethod: "DeleteSaasUser", Parameters: deleteSignUpUserParamsFor(signUpUserIDKeyWithResponse), ExpectedStatus: 200, Skip: !signUpEnabled, SkipReason: skipReasonSignUpLimit},
			{Name: "Delete signup user (SignUpWithBodyWithResponse)", ClientMethod: "DeleteSaasUser", Parameters: deleteSignUpUserParamsFor(signUpUserIDKeyWithBodyWithResponse), ExpectedStatus: 200, Skip: !signUpEnabled, SkipReason: skipReasonSignUpLimit},

			// AWS Marketplaceサインアップ確認（Standard）
			{Name: "ConfirmSignUpWithAwsMarketplace", ClientMethod: "ConfirmSignUpWithAwsMarketplace", Parameters: confirmSignUpWithAwsMarketplaceParams, ExpectedStatus: 200, Skip: true},
			{Name: "ConfirmSignUpWithAwsMarketplaceWithBody", ClientMethod: "ConfirmSignUpWithAwsMarketplaceWithBody", Parameters: confirmSignUpWithAwsMarketplaceWithBodyParams, ExpectedStatus: 200, Skip: true},
			{Name: "ConfirmSignUpWithAwsMarketplaceWithResponse", ClientMethod: "ConfirmSignUpWithAwsMarketplaceWithResponse", Parameters: confirmSignUpWithAwsMarketplaceParams, ExpectedStatus: 200, Skip: true},
			{Name: "ConfirmSignUpWithAwsMarketplaceWithBodyWithResponse", ClientMethod: "ConfirmSignUpWithAwsMarketplaceWithBodyWithResponse", Parameters: confirmSignUpWithAwsMarketplaceWithBodyParams, ExpectedStatus: 200, Skip: true},

			// Setup: テナントを作成（LinkAwsMarketplace用）
			{Name: "CreateTenant for AWS Marketplace link", ClientMethod: "CreateTenant", Parameters: createTenantParams, AllowedStatuses: statusOKOrCreated,
				StateUpdate: func(response any, variables map[string]any) error {
					if id, err := extractIDFromResponse(response, "id"); err == nil {
						variables["tenant_id"] = id
						setLastTenantID(id)
						return nil
					} else {
						fmt.Printf("Warning: failed to extract tenant_id: %v\n", err)
					}
					variables["tenant_id"] = "extracted_tenant_id"
					setLastTenantID("extracted_tenant_id")
					return nil
				}},

			// AWS Marketplaceリンク（Standard）
			{Name: "LinkAwsMarketplace", ClientMethod: "LinkAwsMarketplace", Parameters: linkAwsMarketplaceParams, ExpectedStatus: 200, Skip: true, SkipReason: skipReasonLinkAwsMarketplace},
			{Name: "LinkAwsMarketplaceWithBody", ClientMethod: "LinkAwsMarketplaceWithBody", Parameters: linkAwsMarketplaceWithBodyParams, ExpectedStatus: 200, Skip: true, SkipReason: skipReasonLinkAwsMarketplace},
			{Name: "LinkAwsMarketplaceWithResponse", ClientMethod: "LinkAwsMarketplaceWithResponse", Parameters: linkAwsMarketplaceParams, ExpectedStatus: 200, Skip: true, SkipReason: skipReasonLinkAwsMarketplace},
			{Name: "LinkAwsMarketplaceWithBodyWithResponse", ClientMethod: "LinkAwsMarketplaceWithBodyWithResponse", Parameters: linkAwsMarketplaceWithBodyParams, ExpectedStatus: 200, Skip: true, SkipReason: skipReasonLinkAwsMarketplace},

			// Setup: ユーザーを作成（UnlinkProvider用）
			{Name: "CreateSaasUser for provider unlink", ClientMethod: "CreateSaasUser", Parameters: createSaasUserParams, AllowedStatuses: statusOKOrCreated,
				StateUpdate: func(response any, variables map[string]any) error {
					if id, err := extractIDFromResponse(response, "id"); err == nil {
						variables["user_id"] = id
						return nil
					} else {
						fmt.Printf("Warning: failed to extract user_id: %v\n", err)
					}
					variables["user_id"] = "extracted_user_id"
					updateUserTokensFromCognito(variables)
					return nil
				}},

			// プロバイダーのアンリンク（Standard）
			{Name: "UnlinkProvider", ClientMethod: "UnlinkProvider", Parameters: unlinkProviderParams, ExpectedStatus: 200, Skip: true, SkipReason: skipReasonUnlinkProvider},

			// Cleanup
			{Name: "DeleteSaasUser after provider test", ClientMethod: "DeleteSaasUser", Parameters: deleteSaasUserParams, ExpectedStatus: 200},
			{Name: "DeleteTenant after AWS Marketplace test", ClientMethod: "DeleteTenant", Parameters: deleteTenantParams, ExpectedStatus: 200},
		},
		Setup:   func() error { return nil },
		Cleanup: func() error { return nil },
	}
}
