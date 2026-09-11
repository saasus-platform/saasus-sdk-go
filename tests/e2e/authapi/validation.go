package authapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/saasus-platform/saasus-sdk-go/generated/authapi"
	pricingapi "github.com/saasus-platform/saasus-sdk-go/generated/pricingapi"
	"github.com/saasus-platform/saasus-sdk-go/tests/e2e/authapi/testdata"
)

// このファイルには Auth API のパラメータ関数とバリデーション関数を実装します
// 各 API メソッドに適切なパラメータを提供し、レスポンスを検証します
//
// パラメータ関数の使用方法:
//
// 1. ストーリー関数で testdata.LoadTestParams(t) を呼び出してテストパラメータを読み込む
// 2. 読み込んだパラメータを variables マップに格納する
// 3. パラメータ関数が variables マップから値を取得する
//
// 例:
//
//	func GetPostmanStoryStandardMethods(t *testing.T) testlib.Story {
//	    // テストパラメータを読み込む
//	    params := testdata.LoadTestParams(t)
//
//	    // variables マップを構築
//	    variables := map[string]interface{}{
//	        "email":    params.Users.CreateParams["email"],
//	        "password": params.Users.CreateParams["password"],
//	        // その他の変数...
//	    }
//
//	    return testlib.Story{
//	        Variables: variables,
//	        Steps: []testlib.Step{
//	            {
//	                ClientMethod: "CreateSaasUser",
//	                Parameters:   createSaasUserParams, // variables から値を取得
//	            },
//	        },
//	    }
//	}
//
// ポインタヘルパー関数の使用:
//
// オプショナルなフィールドには testdata.StringPtr(), testdata.IntPtr(), testdata.BoolPtr() を使用します。
//
// 例:
//
//	return authapi.UpdateBasicInfoParam{
//	    DomainName:        domainName,
//	    FromEmailAddress:  fromEmail,
//	    ReplyEmailAddress: testdata.StringPtr(replyEmail), // ポインタ
//	}

// =============================================================================
// 基本設定関連のパラメータ関数
// =============================================================================

// getBasicInfoParams は GetBasicInfo メソッドのパラメータを返します
// メソッドシグネチャ: GetBasicInfo(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
func getBasicInfoParams(variables map[string]interface{}) interface{} {
	// パラメータ不要
	return nil
}

func nextEmailForUpdate(variables map[string]interface{}) string {
	if queue, ok := variables["_email_updates"].([]string); ok && len(queue) > 0 {
		email := queue[0]
		variables["_email_updates"] = queue[1:]
		variables["email"] = email
		return email
	}

	if email, ok := variables["email"].(string); ok && email != "" {
		return email
	}

	return fmt.Sprintf("fallback+%d@example.com", time.Now().UnixNano())
}

func getUserAccessToken(variables map[string]interface{}) string {
	// デフォルト優先順位: access_token -> user_access_token -> cognito_access_token
	if token, ok := variables["access_token"].(string); ok && token != "" {
		return token
	}
	if token, ok := variables["user_access_token"].(string); ok && token != "" {
		return token
	}
	if token, ok := variables["cognito_access_token"].(string); ok && token != "" {
		return token
	}
	return ""
}

func tenantUserID(variables map[string]interface{}) string {
	if tenantUser, ok := variables["tenant_user_id"].(string); ok && tenantUser != "" {
		return tenantUser
	}
	if userID, ok := variables["user_id"].(string); ok {
		return userID
	}
	return ""
}

func tenantRoleName(variables map[string]interface{}) string {
	if role, ok := variables["tenant_role_name"].(string); ok && role != "" {
		return role
	}
	if role, ok := variables["role_name"].(string); ok && role != "" {
		return role
	}
	return "admin"
}

func getStringVariable(variables map[string]interface{}, key, fallback string) string {
	if variables == nil {
		return fallback
	}
	if v, ok := variables[key].(string); ok && v != "" {
		return v
	}
	return fallback
}

func buildTenantBillingInfoParam(data map[string]interface{}) authapi.UpdateTenantBillingInfoParam {
	if data == nil {
		return authapi.UpdateTenantBillingInfoParam{
			Name:            "test_billing",
			InvoiceLanguage: authapi.JaJP,
			Address: authapi.BillingAddress{
				City:       "Tokyo",
				Country:    "JP",
				PostalCode: "100-0000",
				State:      "Tokyo",
				Street:     "1-2-3",
			},
		}
	}
	addrMap, _ := data["address"].(map[string]interface{})
	getString := func(m map[string]interface{}, key, def string) string {
		if m == nil {
			return def
		}
		if v, ok := m[key].(string); ok && v != "" {
			return v
		}
		return def
	}
	address := authapi.BillingAddress{
		City:       getString(addrMap, "city", "Tokyo"),
		Country:    getString(addrMap, "country", "JP"),
		PostalCode: getString(addrMap, "postal_code", "100-0000"),
		State:      getString(addrMap, "state", "Tokyo"),
		Street:     getString(addrMap, "street", "1-2-3"),
		AdditionalAddressInfo: func() *string {
			if v, ok := addrMap["additional_address_info"].(string); ok && v != "" {
				return &v
			}
			return nil
		}(),
	}
	invoiceLangStr, _ := data["invoice_language"].(string)
	invoiceLang := authapi.InvoiceLanguage(invoiceLangStr)
	if invoiceLang == "" {
		invoiceLang = authapi.JaJP
	}
	name := getString(data, "name", "test_billing")
	return authapi.UpdateTenantBillingInfoParam{
		Name:            name,
		InvoiceLanguage: invoiceLang,
		Address:         address,
	}
}

func getCognitoAccessToken(variables map[string]interface{}) string {
	if token, ok := variables["cognito_access_token"].(string); ok && token != "" {
		return token
	}
	return getUserAccessToken(variables)
}

func getSoftwareTokenVerificationCode(variables map[string]interface{}) string {
	secret := ensureSoftwareTokenSecret(variables)
	if secret != "" {
		code, err := generateCurrentTotp(secret)
		if err == nil {
			variables["verification_code"] = code
			return code
		}
		fmt.Printf("Warning: failed to generate TOTP code from secret: %v\n", err)
	}
	if code, ok := variables["verification_code"].(string); ok && code != "" {
		return code
	}
	return ""
}

func ensureSoftwareTokenSecret(variables map[string]interface{}) string {
	if secret, ok := variables[softwareTokenSecretKey].(string); ok && secret != "" {
		return secret
	}
	if secretCode, ok := variables["secret_code"].(string); ok && secretCode != "" {
		variables[softwareTokenSecretKey] = secretCode
		return secretCode
	}
	ensureCognitoSoftwareTokenPrepared(variables)
	if secret, ok := variables[softwareTokenSecretKey].(string); ok && secret != "" {
		return secret
	}
	return ""
}

func buildNotificationMessagesParam(data map[string]interface{}) authapi.UpdateNotificationMessagesParam {
	if data == nil {
		return authapi.UpdateNotificationMessagesParam{}
	}

	return authapi.UpdateNotificationMessagesParam{
		AuthenticationMfa:   toMessageTemplate(data, "authentication_mfa"),
		CreateUser:          toMessageTemplate(data, "create_user"),
		ForgotPassword:      toMessageTemplate(data, "forgot_password"),
		InviteTenantUser:    toMessageTemplate(data, "invite_tenant_user"),
		ResendCode:          toMessageTemplate(data, "resend_code"),
		SignUp:              toMessageTemplate(data, "sign_up"),
		UpdateUserAttribute: toMessageTemplate(data, "update_user_attribute"),
		VerifyExternalUser:  toMessageTemplate(data, "verify_external_user"),
		VerifyUserAttribute: toMessageTemplate(data, "verify_user_attribute"),
	}
}

func toMessageTemplate(data map[string]interface{}, key string) *authapi.MessageTemplate {
	entry, ok := data[key].(map[string]interface{})
	if !ok || entry == nil {
		return nil
	}
	message, _ := entry["message"].(string)
	subject, _ := entry["subject"].(string)
	if message == "" && subject == "" {
		return nil
	}
	return &authapi.MessageTemplate{Message: message, Subject: subject}
}

func buildCustomizePagesParam(data map[string]interface{}) authapi.UpdateCustomizePagesParam {
	if data == nil {
		return authapi.UpdateCustomizePagesParam{}
	}
	return authapi.UpdateCustomizePagesParam{
		PasswordResetPage: toCustomizePageProps(data, "password_reset_page"),
		SignInPage:        toCustomizePageProps(data, "sign_in_page"),
		SignUpPage:        toCustomizePageProps(data, "sign_up_page"),
	}
}

func toCustomizePageProps(data map[string]interface{}, key string) *authapi.CustomizePageProps {
	entry, ok := data[key].(map[string]interface{})
	if !ok || entry == nil {
		return nil
	}
	html, _ := entry["html_contents"].(string)
	isPrivacy, _ := entry["is_privacy_policy"].(bool)
	isTerms, _ := entry["is_terms_of_service"].(bool)
	if html == "" && !isPrivacy && !isTerms {
		return nil
	}
	return &authapi.CustomizePageProps{
		HtmlContents:     html,
		IsPrivacyPolicy:  isPrivacy,
		IsTermsOfService: isTerms,
	}
}

// updateBasicInfoParams は UpdateBasicInfo メソッドのパラメータを返します
// メソッドシグネチャ: UpdateBasicInfo(ctx context.Context, body UpdateBasicInfoParam, reqEditors ...RequestEditorFn) (*http.Response, error)
func updateBasicInfoParams(variables map[string]interface{}) interface{} {
	domainName, _ := variables["domain_name"].(string)
	fromEmail, _ := variables["from_email_address"].(string)
	replyEmail, _ := variables["reply_email_address"].(string)

	return authapi.UpdateBasicInfoParam{
		DomainName:        domainName,
		FromEmailAddress:  fromEmail,
		ReplyEmailAddress: testdata.StringPtr(replyEmail),
	}
}

// updateBasicInfoWithBodyParams は UpdateBasicInfoWithBody メソッドのパラメータを返します
// メソッドシグネチャ: UpdateBasicInfoWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
func updateBasicInfoWithBodyParams(variables map[string]interface{}) interface{} {
	domainName, _ := variables["domain_name"].(string)
	fromEmail, _ := variables["from_email_address"].(string)
	replyEmail, _ := variables["reply_email_address"].(string)

	param := authapi.UpdateBasicInfoParam{
		DomainName:        domainName,
		FromEmailAddress:  fromEmail,
		ReplyEmailAddress: testdata.StringPtr(replyEmail),
	}
	bodyBytes, _ := json.Marshal(param)

	return struct {
		ContentType string
		Body        io.Reader
	}{
		ContentType: "application/json",
		Body:        bytes.NewReader(bodyBytes),
	}
}

// =============================================================================
// 認証情報関連のパラメータ関数
// =============================================================================

// getAuthInfoParams は GetAuthInfo メソッドのパラメータを返します
// メソッドシグネチャ: GetAuthInfo(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
func getAuthInfoParams(variables map[string]interface{}) interface{} {
	// パラメータ不要
	return nil
}

// updateAuthInfoParams は UpdateAuthInfo メソッドのパラメータを返します
// メソッドシグネチャ: UpdateAuthInfo(ctx context.Context, body UpdateAuthInfoParam, reqEditors ...RequestEditorFn) (*http.Response, error)
func updateAuthInfoParams(variables map[string]interface{}) interface{} {
	callbackURL, _ := variables["callback_url"].(string)

	return authapi.UpdateAuthInfoParam{
		CallbackUrl: callbackURL,
	}
}

// updateAuthInfoWithBodyParams は UpdateAuthInfoWithBody メソッドのパラメータを返します
// メソッドシグネチャ: UpdateAuthInfoWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
func updateAuthInfoWithBodyParams(variables map[string]interface{}) interface{} {
	callbackURL, _ := variables["callback_url"].(string)

	param := authapi.UpdateAuthInfoParam{
		CallbackUrl: callbackURL,
	}
	bodyBytes, _ := json.Marshal(param)

	return struct {
		ContentType string
		Body        io.Reader
	}{
		ContentType: "application/json",
		Body:        bytes.NewReader(bodyBytes),
	}
}

// =============================================================================
// ユーザー管理関連のパラメータ関数
// =============================================================================

// getSaasUsersParams は GetSaasUsers メソッドのパラメータを返します
// メソッドシグネチャ: GetSaasUsers(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
func getSaasUsersParams(variables map[string]interface{}) interface{} {
	// パラメータ不要
	return nil
}

// getSaasUserParams は GetSaasUser メソッドのパラメータを返します
// メソッドシグネチャ: GetSaasUser(ctx context.Context, userId string, reqEditors ...RequestEditorFn) (*http.Response, error)
func getSaasUserParams(variables map[string]interface{}) interface{} {
	userID, _ := variables["user_id"].(string)
	return userID
}

// createSaasUserParams は CreateSaasUser メソッドのパラメータを返します
// メソッドシグネチャ: CreateSaasUser(ctx context.Context, body CreateSaasUserParam, reqEditors ...RequestEditorFn) (*http.Response, error)
func createSaasUserParams(variables map[string]interface{}) interface{} {
	email, _ := variables["email"].(string)
	password, _ := variables["password"].(string)

	return authapi.CreateSaasUserParam{
		Email:    testdata.StringPtr(email),
		Password: testdata.StringPtr(password),
	}
}

// createSaasUserWithBodyParams は CreateSaasUserWithBody メソッドのパラメータを返します
// メソッドシグネチャ: CreateSaasUserWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
func createSaasUserWithBodyParams(variables map[string]interface{}) interface{} {
	email, _ := variables["email"].(string)
	password, _ := variables["password"].(string)

	param := authapi.CreateSaasUserParam{
		Email:    testdata.StringPtr(email),
		Password: testdata.StringPtr(password),
	}
	bodyBytes, _ := json.Marshal(param)

	return struct {
		ContentType string
		Body        io.Reader
	}{
		ContentType: "application/json",
		Body:        bytes.NewReader(bodyBytes),
	}
}

// updateSaasUserPasswordParams は UpdateSaasUserPassword メソッドのパラメータを返します
// メソッドシグネチャ: UpdateSaasUserPassword(ctx context.Context, userId string, body UpdateSaasUserPasswordParam, reqEditors ...RequestEditorFn) (*http.Response, error)
func updateSaasUserPasswordParams(variables map[string]interface{}) interface{} {
	userID, _ := variables["user_id"].(string)
	password, _ := variables["password"].(string)

	return struct {
		UserId string
		Body   authapi.UpdateSaasUserPasswordParam
	}{
		UserId: userID,
		Body: authapi.UpdateSaasUserPasswordParam{
			Password: password,
		},
	}
}

// updateSaasUserPasswordWithBodyParams は UpdateSaasUserPasswordWithBody メソッドのパラメータを返します
// メソッドシグネチャ: UpdateSaasUserPasswordWithBody(ctx context.Context, userId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
func updateSaasUserPasswordWithBodyParams(variables map[string]interface{}) interface{} {
	userID, _ := variables["user_id"].(string)
	password, _ := variables["password"].(string)

	param := authapi.UpdateSaasUserPasswordParam{
		Password: password,
	}
	bodyBytes, _ := json.Marshal(param)

	return struct {
		UserId      string
		ContentType string
		Body        io.Reader
	}{
		UserId:      userID,
		ContentType: "application/json",
		Body:        bytes.NewReader(bodyBytes),
	}
}

// updateSaasUserEmailParams は UpdateSaasUserEmail メソッドのパラメータを返します
// メソッドシグネチャ: UpdateSaasUserEmail(ctx context.Context, userId string, body UpdateSaasUserEmailParam, reqEditors ...RequestEditorFn) (*http.Response, error)
func updateSaasUserEmailParams(variables map[string]interface{}) interface{} {
	userID, _ := variables["user_id"].(string)
	email := nextEmailForUpdate(variables)

	return struct {
		UserId string
		Body   authapi.UpdateSaasUserEmailParam
	}{
		UserId: userID,
		Body: authapi.UpdateSaasUserEmailParam{
			Email: email,
		},
	}
}

// updateSaasUserEmailWithBodyParams は UpdateSaasUserEmailWithBody メソッドのパラメータを返します
// メソッドシグネチャ: UpdateSaasUserEmailWithBody(ctx context.Context, userId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
func updateSaasUserEmailWithBodyParams(variables map[string]interface{}) interface{} {
	userID, _ := variables["user_id"].(string)
	email := nextEmailForUpdate(variables)

	param := authapi.UpdateSaasUserEmailParam{
		Email: email,
	}
	bodyBytes, _ := json.Marshal(param)

	return struct {
		UserId      string
		ContentType string
		Body        io.Reader
	}{
		UserId:      userID,
		ContentType: "application/json",
		Body:        bytes.NewReader(bodyBytes),
	}
}

// deleteSaasUserParams は DeleteSaasUser メソッドのパラメータを返します
// メソッドシグネチャ: DeleteSaasUser(ctx context.Context, userId string, reqEditors ...RequestEditorFn) (*http.Response, error)
func deleteSaasUserParams(variables map[string]interface{}) interface{} {
	userID, _ := variables["user_id"].(string)
	return userID
}

// =============================================================================
// MFA 設定関連のパラメータ関数
// =============================================================================

// getUserMfaPreferenceParams は GetUserMfaPreference メソッドのパラメータを返します
// メソッドシグネチャ: GetUserMfaPreference(ctx context.Context, userId string, reqEditors ...RequestEditorFn) (*http.Response, error)
func getUserMfaPreferenceParams(variables map[string]interface{}) interface{} {
	userID, _ := variables["user_id"].(string)
	return userID
}

// updateUserMfaPreferenceParams は UpdateUserMfaPreference メソッドのパラメータを返します
// メソッドシグネチャ: UpdateUserMfaPreference(ctx context.Context, userId string, body UpdateUserMfaPreferenceParam, reqEditors ...RequestEditorFn) (*http.Response, error)
func updateUserMfaPreferenceParams(variables map[string]interface{}) interface{} {
	userID, _ := variables["user_id"].(string)
	enabled, _ := variables["enabled"].(bool)

	method := authapi.MfaPreferenceMethodSoftwareToken
	return struct {
		UserId string
		Body   authapi.UpdateUserMfaPreferenceParam
	}{
		UserId: userID,
		Body: authapi.UpdateUserMfaPreferenceParam{
			Enabled: enabled,
			Method:  &method,
		},
	}
}

// updateUserMfaPreferenceWithBodyParams は UpdateUserMfaPreferenceWithBody メソッドのパラメータを返します
// メソッドシグネチャ: UpdateUserMfaPreferenceWithBody(ctx context.Context, userId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
func updateUserMfaPreferenceWithBodyParams(variables map[string]interface{}) interface{} {
	userID, _ := variables["user_id"].(string)
	enabled, _ := variables["enabled"].(bool)

	method := authapi.MfaPreferenceMethodSoftwareToken
	param := authapi.UpdateUserMfaPreferenceParam{
		Enabled: enabled,
		Method:  &method,
	}
	bodyBytes, _ := json.Marshal(param)

	return struct {
		UserId      string
		ContentType string
		Body        io.Reader
	}{
		UserId:      userID,
		ContentType: "application/json",
		Body:        bytes.NewReader(bodyBytes),
	}
}

// updateSoftwareTokenParams は UpdateSoftwareToken メソッドのパラメータを返します
// メソッドシグネチャ: UpdateSoftwareToken(ctx context.Context, userId string, body UpdateSoftwareTokenParam, reqEditors ...RequestEditorFn) (*http.Response, error)
func updateSoftwareTokenParams(variables map[string]interface{}) interface{} {
	userID, _ := variables["user_id"].(string)
	accessToken := getCognitoAccessToken(variables)
	verificationCode := getSoftwareTokenVerificationCode(variables)

	return struct {
		UserId string
		Body   authapi.UpdateSoftwareTokenParam
	}{
		UserId: userID,
		Body: authapi.UpdateSoftwareTokenParam{
			AccessToken:      accessToken,
			VerificationCode: verificationCode,
		},
	}
}

// updateSoftwareTokenWithBodyParams は UpdateSoftwareTokenWithBody メソッドのパラメータを返します
// メソッドシグネチャ: UpdateSoftwareTokenWithBody(ctx context.Context, userId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
func updateSoftwareTokenWithBodyParams(variables map[string]interface{}) interface{} {
	userID, _ := variables["user_id"].(string)
	accessToken := getCognitoAccessToken(variables)
	verificationCode := getSoftwareTokenVerificationCode(variables)

	param := authapi.UpdateSoftwareTokenParam{
		AccessToken:      accessToken,
		VerificationCode: verificationCode,
	}
	bodyBytes, _ := json.Marshal(param)

	return struct {
		UserId      string
		ContentType string
		Body        io.Reader
	}{
		UserId:      userID,
		ContentType: "application/json",
		Body:        bytes.NewReader(bodyBytes),
	}
}

// createSecretCodeParams は CreateSecretCode メソッドのパラメータを返します
// メソッドシグネチャ: CreateSecretCode(ctx context.Context, userId string, body CreateSecretCodeParam, reqEditors ...RequestEditorFn) (*http.Response, error)
func createSecretCodeParams(variables map[string]interface{}) interface{} {
	userID, _ := variables["user_id"].(string)
	accessToken := getCognitoAccessToken(variables)

	// デバッグ: どのトークンが使われているか確認
	fmt.Printf("[DEBUG] CreateSecretCode: user_id=%s\n", userID)
	fmt.Printf("[DEBUG] Token sources: cognito=%v, user=%v, default=%v\n",
		variables["cognito_access_token"] != nil && variables["cognito_access_token"] != "",
		variables["user_access_token"] != nil && variables["user_access_token"] != "",
		variables["access_token"] != nil && variables["access_token"] != "")
	if accessToken != "" && len(accessToken) > 20 {
		fmt.Printf("[DEBUG] Using access_token prefix: %s...\n", accessToken[:20])
	}

	return struct {
		UserId string
		Body   authapi.CreateSecretCodeParam
	}{
		UserId: userID,
		Body: authapi.CreateSecretCodeParam{
			AccessToken: accessToken,
		},
	}
}

// createSecretCodeWithBodyParams は CreateSecretCodeWithBody メソッドのパラメータを返します
// メソッドシグネチャ: CreateSecretCodeWithBody(ctx context.Context, userId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
func createSecretCodeWithBodyParams(variables map[string]interface{}) interface{} {
	userID, _ := variables["user_id"].(string)
	accessToken := getCognitoAccessToken(variables)

	param := authapi.CreateSecretCodeParam{
		AccessToken: accessToken,
	}
	bodyBytes, _ := json.Marshal(param)

	return struct {
		UserId      string
		ContentType string
		Body        io.Reader
	}{
		UserId:      userID,
		ContentType: "application/json",
		Body:        bytes.NewReader(bodyBytes),
	}
}

// =============================================================================
// ロール管理関連のパラメータ関数
// =============================================================================

// getRolesParams は GetRoles メソッドのパラメータを返します
// メソッドシグネチャ: GetRoles(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
func getRolesParams(variables map[string]interface{}) interface{} {
	// パラメータ不要
	return nil
}

// createRoleParams は CreateRole メソッドのパラメータを返します
// メソッドシグネチャ: CreateRole(ctx context.Context, body CreateRoleParam, reqEditors ...RequestEditorFn) (*http.Response, error)
func createRoleParams(variables map[string]interface{}) interface{} {
	roleName := tenantRoleName(variables)
	displayName, _ := variables["display_name"].(string)

	return authapi.CreateRoleParam{
		RoleName:    roleName,
		DisplayName: displayName,
	}
}

// createRoleWithBodyParams は CreateRoleWithBody メソッドのパラメータを返します
// メソッドシグネチャ: CreateRoleWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
func createRoleWithBodyParams(variables map[string]interface{}) interface{} {
	roleName, _ := variables["role_name"].(string)
	displayName, _ := variables["display_name"].(string)

	param := authapi.CreateRoleParam{
		RoleName:    roleName,
		DisplayName: displayName,
	}
	bodyBytes, _ := json.Marshal(param)

	return struct {
		ContentType string
		Body        io.Reader
	}{
		ContentType: "application/json",
		Body:        bytes.NewReader(bodyBytes),
	}
}

// deleteRoleParams は DeleteRole メソッドのパラメータを返します
// メソッドシグネチャ: DeleteRole(ctx context.Context, roleName string, reqEditors ...RequestEditorFn) (*http.Response, error)
func deleteRoleParams(variables map[string]interface{}) interface{} {
	roleName, _ := variables["role_name"].(string)
	return roleName
}

// =============================================================================
// 属性管理関連のパラメータ関数
// =============================================================================

// getUserAttributesParams は GetUserAttributes メソッドのパラメータを返します
// メソッドシグネチャ: GetUserAttributes(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
func getUserAttributesParams(variables map[string]interface{}) interface{} {
	// パラメータ不要
	return nil
}

// createUserAttributeParams は CreateUserAttribute メソッドのパラメータを返します
// メソッドシグネチャ: CreateUserAttribute(ctx context.Context, body CreateUserAttributeParam, reqEditors ...RequestEditorFn) (*http.Response, error)
func createUserAttributeParams(variables map[string]interface{}) interface{} {
	attributeName, _ := variables["attribute_name"].(string)
	displayName, _ := variables["display_name"].(string)
	attributeType, _ := variables["attribute_type"].(string)

	return authapi.CreateUserAttributeParam{
		AttributeName: attributeName,
		DisplayName:   displayName,
		AttributeType: authapi.AttributeType(attributeType),
	}
}

// createUserAttributeWithBodyParams は CreateUserAttributeWithBody メソッドのパラメータを返します
// メソッドシグネチャ: CreateUserAttributeWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
func createUserAttributeWithBodyParams(variables map[string]interface{}) interface{} {
	attributeName, _ := variables["attribute_name"].(string)
	displayName, _ := variables["display_name"].(string)
	attributeType, _ := variables["attribute_type"].(string)

	param := authapi.CreateUserAttributeParam{
		AttributeName: attributeName,
		DisplayName:   displayName,
		AttributeType: authapi.AttributeType(attributeType),
	}
	bodyBytes, _ := json.Marshal(param)

	return struct {
		ContentType string
		Body        io.Reader
	}{
		ContentType: "application/json",
		Body:        bytes.NewReader(bodyBytes),
	}
}

// deleteUserAttributeParams は DeleteUserAttribute メソッドのパラメータを返します
// メソッドシグネチャ: DeleteUserAttribute(ctx context.Context, attributeName string, reqEditors ...RequestEditorFn) (*http.Response, error)
func deleteUserAttributeParams(variables map[string]interface{}) interface{} {
	attributeName, _ := variables["attribute_name"].(string)
	return attributeName
}

// getTenantAttributesParams は GetTenantAttributes メソッドのパラメータを返します
// メソッドシグネチャ: GetTenantAttributes(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
func getTenantAttributesParams(variables map[string]interface{}) interface{} {
	// パラメータ不要
	return nil
}

// createTenantAttributeParams は CreateTenantAttribute メソッドのパラメータを返します
// メソッドシグネチャ: CreateTenantAttribute(ctx context.Context, body CreateTenantAttributeParam, reqEditors ...RequestEditorFn) (*http.Response, error)
func createTenantAttributeParams(variables map[string]interface{}) interface{} {
	attributeName, _ := variables["attribute_name"].(string)
	displayName, _ := variables["display_name"].(string)
	attributeType, _ := variables["attribute_type"].(string)

	return authapi.CreateTenantAttributeParam{
		AttributeName: attributeName,
		DisplayName:   displayName,
		AttributeType: authapi.AttributeType(attributeType),
	}
}

// createTenantAttributeWithBodyParams は CreateTenantAttributeWithBody メソッドのパラメータを返します
// メソッドシグネチャ: CreateTenantAttributeWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
func createTenantAttributeWithBodyParams(variables map[string]interface{}) interface{} {
	attributeName, _ := variables["attribute_name"].(string)
	displayName, _ := variables["display_name"].(string)
	attributeType, _ := variables["attribute_type"].(string)

	param := authapi.CreateTenantAttributeParam{
		AttributeName: attributeName,
		DisplayName:   displayName,
		AttributeType: authapi.AttributeType(attributeType),
	}
	bodyBytes, _ := json.Marshal(param)

	return struct {
		ContentType string
		Body        io.Reader
	}{
		ContentType: "application/json",
		Body:        bytes.NewReader(bodyBytes),
	}
}

// deleteTenantAttributeParams は DeleteTenantAttribute メソッドのパラメータを返します
// メソッドシグネチャ: DeleteTenantAttribute(ctx context.Context, attributeName string, reqEditors ...RequestEditorFn) (*http.Response, error)
func deleteTenantAttributeParams(variables map[string]interface{}) interface{} {
	attributeName, _ := variables["attribute_name"].(string)
	return attributeName
}

// =============================================================================
// 通知・カスタマイズ関連のパラメータ関数
// =============================================================================

// findNotificationMessagesParams は FindNotificationMessages メソッドのパラメータを返します
// メソッドシグネチャ: FindNotificationMessages(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
func findNotificationMessagesParams(variables map[string]interface{}) interface{} {
	// パラメータ不要
	return nil
}

// updateNotificationMessagesParams は UpdateNotificationMessages メソッドのパラメータを返します
// メソッドシグネチャ: UpdateNotificationMessages(ctx context.Context, body UpdateNotificationMessagesParam, reqEditors ...RequestEditorFn) (*http.Response, error)
func updateNotificationMessagesParams(variables map[string]interface{}) interface{} {
	messages, _ := variables["notification_messages"].(map[string]interface{})
	return buildNotificationMessagesParam(messages)
}

// updateNotificationMessagesWithBodyParams は UpdateNotificationMessagesWithBody メソッドのパラメータを返します
// メソッドシグネチャ: UpdateNotificationMessagesWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
func updateNotificationMessagesWithBodyParams(variables map[string]interface{}) interface{} {
	messages, _ := variables["notification_messages"].(map[string]interface{})
	param := buildNotificationMessagesParam(messages)
	bodyBytes, _ := json.Marshal(param)

	return struct {
		ContentType string
		Body        io.Reader
	}{
		ContentType: "application/json",
		Body:        bytes.NewReader(bodyBytes),
	}
}

// getCustomizePagesParams は GetCustomizePages メソッドのパラメータを返します
// メソッドシグネチャ: GetCustomizePages(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
func getCustomizePagesParams(variables map[string]interface{}) interface{} {
	// パラメータ不要
	return nil
}

// updateCustomizePagesParams は UpdateCustomizePages メソッドのパラメータを返します
// メソッドシグネチャ: UpdateCustomizePages(ctx context.Context, body UpdateCustomizePagesParam, reqEditors ...RequestEditorFn) (*http.Response, error)
func updateCustomizePagesParams(variables map[string]interface{}) interface{} {
	pages, _ := variables["customize_pages"].(map[string]interface{})
	return buildCustomizePagesParam(pages)
}

// updateCustomizePagesWithBodyParams は UpdateCustomizePagesWithBody メソッドのパラメータを返します
// メソッドシグネチャ: UpdateCustomizePagesWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
func updateCustomizePagesWithBodyParams(variables map[string]interface{}) interface{} {
	pages, _ := variables["customize_pages"].(map[string]interface{})
	param := buildCustomizePagesParam(pages)
	bodyBytes, _ := json.Marshal(param)

	return struct {
		ContentType string
		Body        io.Reader
	}{
		ContentType: "application/json",
		Body:        bytes.NewReader(bodyBytes),
	}
}

// getCustomizePageSettingsParams は GetCustomizePageSettings メソッドのパラメータを返します
// メソッドシグネチャ: GetCustomizePageSettings(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
func getCustomizePageSettingsParams(variables map[string]interface{}) interface{} {
	// パラメータ不要
	return nil
}

// updateCustomizePageSettingsParams は UpdateCustomizePageSettings メソッドのパラメータを返します
// メソッドシグネチャ: UpdateCustomizePageSettings(ctx context.Context, body UpdateCustomizePageSettingsParam, reqEditors ...RequestEditorFn) (*http.Response, error)

func updateCustomizePageSettingsParams(variables map[string]interface{}) interface{} {
	favicon, _ := variables["favicon"].(string)
	icon, _ := variables["icon"].(string)
	title, _ := variables["title"].(string)
	termsURL, _ := variables["terms_of_service_url"].(string)
	privacyURL, _ := variables["privacy_policy_url"].(string)
	gtmID, _ := variables["google_tag_manager_container_id"].(string)

	return authapi.UpdateCustomizePageSettingsParam{
		Favicon:                     favicon,
		Icon:                        icon,
		Title:                       title,
		TermsOfServiceUrl:           termsURL,
		PrivacyPolicyUrl:            privacyURL,
		GoogleTagManagerContainerId: gtmID,
	}
}

// updateCustomizePageSettingsWithBodyParams は UpdateCustomizePageSettingsWithBody メソッドのパラメータを返します
// メソッドシグネチャ: UpdateCustomizePageSettingsWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
func updateCustomizePageSettingsWithBodyParams(variables map[string]interface{}) interface{} {
	favicon, _ := variables["favicon"].(string)
	icon, _ := variables["icon"].(string)
	title, _ := variables["title"].(string)
	termsURL, _ := variables["terms_of_service_url"].(string)
	privacyURL, _ := variables["privacy_policy_url"].(string)
	gtmID, _ := variables["google_tag_manager_container_id"].(string)

	param := authapi.UpdateCustomizePageSettingsParam{
		Favicon:                     favicon,
		Icon:                        icon,
		Title:                       title,
		TermsOfServiceUrl:           termsURL,
		PrivacyPolicyUrl:            privacyURL,
		GoogleTagManagerContainerId: gtmID,
	}
	bodyBytes, _ := json.Marshal(param)

	return struct {
		ContentType string
		Body        io.Reader
	}{
		ContentType: "application/json",
		Body:        bytes.NewReader(bodyBytes),
	}
}

// =============================================================================
// 環境・サインイン設定関連のパラメータ関数
// =============================================================================

// getEnvsParams は GetEnvs メソッドのパラメータを返します
// メソッドシグネチャ: GetEnvs(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
func getEnvsParams(variables map[string]interface{}) interface{} {
	// パラメータ不要
	return nil
}

// createEnvParams は CreateEnv メソッドのパラメータを返します
// メソッドシグネチャ: CreateEnv(ctx context.Context, body Env, reqEditors ...RequestEditorFn) (*http.Response, error)
func createEnvParams(variables map[string]interface{}) interface{} {
	envID, ok := variables["env_id"].(int)
	if !ok {
		if envIDFloat, ok := variables["env_id"].(float64); ok {
			envID = int(envIDFloat)
		} else {
			envID = 10 // デフォルト値
		}
	}
	name, _ := variables["name"].(string)
	displayName, _ := variables["display_name"].(string)

	return authapi.Env{
		Id:          authapi.Id(uint64(envID)),
		Name:        name,
		DisplayName: testdata.StringPtr(displayName),
	}
}

// createEnvWithBodyParams は CreateEnvWithBody メソッドのパラメータを返します
// メソッドシグネチャ: CreateEnvWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
func createEnvWithBodyParams(variables map[string]interface{}) interface{} {
	envID, ok := variables["env_id"].(int)
	if !ok {
		if envIDFloat, ok := variables["env_id"].(float64); ok {
			envID = int(envIDFloat)
		} else {
			envID = 10 // デフォルト値
		}
	}
	name, _ := variables["name"].(string)
	displayName, _ := variables["display_name"].(string)

	param := authapi.Env{
		Id:          authapi.Id(uint64(envID)),
		Name:        name,
		DisplayName: testdata.StringPtr(displayName),
	}
	bodyBytes, _ := json.Marshal(param)

	return struct {
		ContentType string
		Body        io.Reader
	}{
		ContentType: "application/json",
		Body:        bytes.NewReader(bodyBytes),
	}
}

// getEnvParams は GetEnv メソッドのパラメータを返します
// メソッドシグネチャ: GetEnv(ctx context.Context, envId uint64, reqEditors ...RequestEditorFn) (*http.Response, error)
func getEnvParams(variables map[string]interface{}) interface{} {
	envID, ok := variables["env_id"].(int)
	if !ok {
		// Try float64 (JSON unmarshaling default for numbers)
		if envIDFloat, ok := variables["env_id"].(float64); ok {
			if envIDFloat >= 0 && envIDFloat <= float64(^uint64(0)) {
				return uint64(envIDFloat)
			}
		}
		return uint64(10) // デフォルト値
	}
	if envID >= 0 {
		return uint64(envID)
	}
	return uint64(10) // デフォルト値
}

// updateEnvParams は UpdateEnv メソッドのパラメータを返します
// メソッドシグネチャ: UpdateEnv(ctx context.Context, envId uint64, body UpdateEnvParam, reqEditors ...RequestEditorFn) (*http.Response, error)
func updateEnvParams(variables map[string]interface{}) interface{} {
	envID, ok := variables["env_id"].(int)
	if !ok {
		if envIDFloat, ok := variables["env_id"].(float64); ok {
			envID = int(envIDFloat)
		} else {
			envID = 10 // デフォルト値
		}
	}
	name, _ := variables["name"].(string)
	displayName, _ := variables["display_name"].(string)

	return struct {
		EnvId uint64
		Body  authapi.UpdateEnvParam
	}{
		EnvId: uint64(envID),
		Body: authapi.UpdateEnvParam{
			Name:        name,
			DisplayName: testdata.StringPtr(displayName),
		},
	}
}

// updateEnvWithBodyParams は UpdateEnvWithBody メソッドのパラメータを返します
// メソッドシグネチャ: UpdateEnvWithBody(ctx context.Context, envId uint64, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
func updateEnvWithBodyParams(variables map[string]interface{}) interface{} {
	envID, ok := variables["env_id"].(int)
	if !ok {
		if envIDFloat, ok := variables["env_id"].(float64); ok {
			envID = int(envIDFloat)
		} else {
			envID = 10 // デフォルト値
		}
	}
	name, _ := variables["name"].(string)
	displayName, _ := variables["display_name"].(string)

	param := authapi.UpdateEnvParam{
		Name:        name,
		DisplayName: testdata.StringPtr(displayName),
	}
	bodyBytes, _ := json.Marshal(param)

	return struct {
		EnvId       uint64
		ContentType string
		Body        io.Reader
	}{
		EnvId:       uint64(envID),
		ContentType: "application/json",
		Body:        bytes.NewReader(bodyBytes),
	}
}

// deleteEnvParams は DeleteEnv メソッドのパラメータを返します
// メソッドシグネチャ: DeleteEnv(ctx context.Context, envId uint64, reqEditors ...RequestEditorFn) (*http.Response, error)
func deleteEnvParams(variables map[string]interface{}) interface{} {
	envID, ok := variables["env_id"].(int)
	if !ok {
		if envIDFloat, ok := variables["env_id"].(float64); ok {
			if envIDFloat >= 0 && envIDFloat <= float64(^uint64(0)) {
				return uint64(envIDFloat)
			}
		}
		return uint64(10) // デフォルト値
	}
	if envID >= 0 {
		return uint64(envID)
	}
	return uint64(10) // デフォルト値
}

// getSignInSettingsParams は GetSignInSettings メソッドのパラメータを返します
// メソッドシグネチャ: GetSignInSettings(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
func getSignInSettingsParams(variables map[string]interface{}) interface{} {
	// パラメータ不要
	return nil
}

// updateSignInSettingsParams は UpdateSignInSettings メソッドのパラメータを返します
// メソッドシグネチャ: UpdateSignInSettings(ctx context.Context, body UpdateSignInSettingsParam, reqEditors ...RequestEditorFn) (*http.Response, error)
func updateSignInSettingsParams(variables map[string]interface{}) interface{} {
	// セルフサインアップを有効化
	return authapi.UpdateSignInSettingsParam{
		SelfRegist: &authapi.SelfRegist{
			Enable: true,
		},
	}
}

// updateSignInSettingsWithBodyParams は UpdateSignInSettingsWithBody メソッドのパラメータを返します
// メソッドシグネチャ: UpdateSignInSettingsWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
func updateSignInSettingsWithBodyParams(variables map[string]interface{}) interface{} {
	// セルフサインアップを有効化
	param := authapi.UpdateSignInSettingsParam{
		SelfRegist: &authapi.SelfRegist{
			Enable: true,
		},
	}
	bodyBytes, _ := json.Marshal(param)

	return struct {
		ContentType string
		Body        io.Reader
	}{
		ContentType: "application/json",
		Body:        bytes.NewReader(bodyBytes),
	}
}

// =============================================================================
// テナント管理関連のパラメータ関数
// =============================================================================

// getTenantsParams は GetTenants メソッドのパラメータを返します
// メソッドシグネチャ: GetTenants(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
func getTenantsParams(variables map[string]interface{}) interface{} {
	// パラメータ不要
	return nil
}

// createTenantParams は CreateTenant メソッドのパラメータを返します
// メソッドシグネチャ: CreateTenant(ctx context.Context, body CreateTenantParam, reqEditors ...RequestEditorFn) (*http.Response, error)
func createTenantParams(variables map[string]interface{}) interface{} {
	name, _ := variables["name"].(string)
	backOfficeEmail, _ := variables["back_office_staff_email"].(string)
	attributes, _ := variables["attributes"].(map[string]interface{})

	return authapi.CreateTenantParam{
		Name:                 name,
		BackOfficeStaffEmail: backOfficeEmail,
		Attributes:           attributes,
	}
}

// createTenantWithBodyParams は CreateTenantWithBody メソッドのパラメータを返します
// メソッドシグネチャ: CreateTenantWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
func createTenantWithBodyParams(variables map[string]interface{}) interface{} {
	name, _ := variables["name"].(string)
	backOfficeEmail, _ := variables["back_office_staff_email"].(string)
	attributes, _ := variables["attributes"].(map[string]interface{})

	param := authapi.CreateTenantParam{
		Name:                 name,
		BackOfficeStaffEmail: backOfficeEmail,
		Attributes:           attributes,
	}
	bodyBytes, _ := json.Marshal(param)

	return struct {
		ContentType string
		Body        io.Reader
	}{
		ContentType: "application/json",
		Body:        bytes.NewReader(bodyBytes),
	}
}

// getTenantParams は GetTenant メソッドのパラメータを返します
// メソッドシグネチャ: GetTenant(ctx context.Context, tenantId string, reqEditors ...RequestEditorFn) (*http.Response, error)
func getTenantParams(variables map[string]interface{}) interface{} {
	tenantID, _ := variables["tenant_id"].(string)
	return tenantID
}

// updateTenantParams は UpdateTenant メソッドのパラメータを返します
// メソッドシグネチャ: UpdateTenant(ctx context.Context, tenantId string, body UpdateTenantParam, reqEditors ...RequestEditorFn) (*http.Response, error)
func updateTenantParams(variables map[string]interface{}) interface{} {
	tenantID, _ := variables["tenant_id"].(string)
	name, _ := variables["name"].(string)
	backOfficeEmail, _ := variables["back_office_staff_email"].(string)
	attributes, _ := variables["attributes"].(map[string]interface{})

	return struct {
		TenantId string
		Body     authapi.UpdateTenantParam
	}{
		TenantId: tenantID,
		Body: authapi.UpdateTenantParam{
			Name:                 name,
			BackOfficeStaffEmail: backOfficeEmail,
			Attributes:           attributes,
		},
	}
}

// updateTenantWithBodyParams は UpdateTenantWithBody メソッドのパラメータを返します
// メソッドシグネチャ: UpdateTenantWithBody(ctx context.Context, tenantId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
func updateTenantWithBodyParams(variables map[string]interface{}) interface{} {
	tenantID, _ := variables["tenant_id"].(string)
	name, _ := variables["name"].(string)
	backOfficeEmail, _ := variables["back_office_staff_email"].(string)
	attributes, _ := variables["attributes"].(map[string]interface{})

	param := authapi.UpdateTenantParam{
		Name:                 name,
		BackOfficeStaffEmail: backOfficeEmail,
		Attributes:           attributes,
	}
	bodyBytes, _ := json.Marshal(param)

	return struct {
		TenantId    string
		ContentType string
		Body        io.Reader
	}{
		TenantId:    tenantID,
		ContentType: "application/json",
		Body:        bytes.NewReader(bodyBytes),
	}
}

// deleteTenantParams は DeleteTenant メソッドのパラメータを返します
// メソッドシグネチャ: DeleteTenant(ctx context.Context, tenantId string, reqEditors ...RequestEditorFn) (*http.Response, error)
func deleteTenantParams(variables map[string]interface{}) interface{} {
	tenantID, _ := variables["tenant_id"].(string)
	return tenantID
}

// createTenantUserParams は CreateTenantUser メソッドのパラメータを返します
// メソッドシグネチャ: CreateTenantUser(ctx context.Context, tenantId string, body CreateTenantUserParam, reqEditors ...RequestEditorFn) (*http.Response, error)
func createTenantUserParams(variables map[string]interface{}) interface{} {
	tenantID, _ := variables["tenant_id"].(string)
	email, _ := variables["email"].(string)
	attributes, _ := variables["attributes"].(map[string]interface{})

	return struct {
		TenantId string
		Body     authapi.CreateTenantUserParam
	}{
		TenantId: tenantID,
		Body: authapi.CreateTenantUserParam{
			Email:      testdata.StringPtr(email),
			Attributes: attributes,
		},
	}
}

// createTenantUserWithBodyParams は CreateTenantUserWithBody メソッドのパラメータを返します
// メソッドシグネチャ: CreateTenantUserWithBody(ctx context.Context, tenantId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
func createTenantUserWithBodyParams(variables map[string]interface{}) interface{} {
	tenantID, _ := variables["tenant_id"].(string)
	email, _ := variables["email"].(string)
	attributes, _ := variables["attributes"].(map[string]interface{})

	param := authapi.CreateTenantUserParam{
		Email:      testdata.StringPtr(email),
		Attributes: attributes,
	}
	bodyBytes, _ := json.Marshal(param)

	return struct {
		TenantId    string
		ContentType string
		Body        io.Reader
	}{
		TenantId:    tenantID,
		ContentType: "application/json",
		Body:        bytes.NewReader(bodyBytes),
	}
}

// getAllTenantUsersParams は GetAllTenantUsers メソッドのパラメータを返します
// メソッドシグネチャ: GetAllTenantUsers(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
func getAllTenantUsersParams(variables map[string]interface{}) interface{} {
	// パラメータ不要
	return nil
}

// getAllTenantUserParams は GetAllTenantUser メソッドのパラメータを返します
// メソッドシグネチャ: GetAllTenantUser(ctx context.Context, userId string, reqEditors ...RequestEditorFn) (*http.Response, error)
func getAllTenantUserParams(variables map[string]interface{}) interface{} {
	userID, _ := variables["user_id"].(string)
	return userID
}

// =============================================================================
// testdata からのパラメータ読み込み
// =============================================================================

// パラメータ関数は variables マップから値を取得します。
// variables マップは以下のソースから値を取得できます：
//
// 1. testdata/test_params.json - 自動生成されたテストパラメータ
//    LoadTestParams() 関数を使用して読み込みます
//
// 2. ストーリーの Variables フィールド - 初期変数
//    各ストーリーで定義された初期値
//
// 3. StateUpdate 関数 - レスポンスから抽出された値
//    API レスポンスから動的に抽出された値（例: user_id, tenant_id）
//
// 使用例:
//   params := testdata.LoadTestParams(t)
//   variables := map[string]interface{}{
//       "email":    params.Users.CreateParams["email"],
//       "password": params.Users.CreateParams["password"],
//   }
//   result := createSaasUserParams(variables)

// =============================================================================
// テナントユーザー管理関連のパラメータ関数
// =============================================================================

// getTenantUsersParams は GetTenantUsers メソッドのパラメータを返します
// メソッドシグネチャ: GetTenantUsers(ctx context.Context, tenantId string, reqEditors ...RequestEditorFn) (*http.Response, error)
func getTenantUsersParams(variables map[string]interface{}) interface{} {
	tenantID, _ := variables["tenant_id"].(string)
	return tenantID
}

// getTenantUserParams は GetTenantUser メソッドのパラメータを返します
// メソッドシグネチャ: GetTenantUser(ctx context.Context, tenantId string, userId string, reqEditors ...RequestEditorFn) (*http.Response, error)
func getTenantUserParams(variables map[string]interface{}) interface{} {
	tenantID, _ := variables["tenant_id"].(string)
	userID := tenantUserID(variables)
	return struct {
		TenantId string
		UserId   string
	}{
		TenantId: tenantID,
		UserId:   userID,
	}
}

// updateTenantUserParams は UpdateTenantUser メソッドのパラメータを返します
// メソッドシグネチャ: UpdateTenantUser(ctx context.Context, tenantId string, userId string, body UpdateTenantUserParam, reqEditors ...RequestEditorFn) (*http.Response, error)
func updateTenantUserParams(variables map[string]interface{}) interface{} {
	tenantID, _ := variables["tenant_id"].(string)
	userID := tenantUserID(variables)
	attributes, _ := variables["attributes"].(map[string]interface{})

	return struct {
		TenantId string
		UserId   string
		Body     authapi.UpdateTenantUserParam
	}{
		TenantId: tenantID,
		UserId:   userID,
		Body: authapi.UpdateTenantUserParam{
			Attributes: attributes,
		},
	}
}

// updateTenantUserWithBodyParams は UpdateTenantUserWithBody メソッドのパラメータを返します
// メソッドシグネチャ: UpdateTenantUserWithBody(ctx context.Context, tenantId string, userId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
func updateTenantUserWithBodyParams(variables map[string]interface{}) interface{} {
	tenantID, _ := variables["tenant_id"].(string)
	userID := tenantUserID(variables)
	attributes, _ := variables["attributes"].(map[string]interface{})

	param := authapi.UpdateTenantUserParam{
		Attributes: attributes,
	}
	bodyBytes, _ := json.Marshal(param)

	return struct {
		TenantId    string
		UserId      string
		ContentType string
		Body        io.Reader
	}{
		TenantId:    tenantID,
		UserId:      userID,
		ContentType: "application/json",
		Body:        bytes.NewReader(bodyBytes),
	}
}

// deleteTenantUserParams は DeleteTenantUser メソッドのパラメータを返します
// メソッドシグネチャ: DeleteTenantUser(ctx context.Context, tenantId string, userId string, reqEditors ...RequestEditorFn) (*http.Response, error)
func deleteTenantUserParams(variables map[string]interface{}) interface{} {
	tenantID, _ := variables["tenant_id"].(string)
	userID := tenantUserID(variables)
	return struct {
		TenantId string
		UserId   string
	}{
		TenantId: tenantID,
		UserId:   userID,
	}
}

// =============================================================================
// テナントユーザーロール関連のパラメータ関数
// =============================================================================

// createTenantUserRolesParams は CreateTenantUserRoles メソッドのパラメータを返します
// メソッドシグネチャ: CreateTenantUserRoles(ctx context.Context, tenantId string, userId string, envId uint64, body CreateTenantUserRolesParam, reqEditors ...RequestEditorFn) (*http.Response, error)
func createTenantUserRolesParams(variables map[string]interface{}) interface{} {
	tenantID, _ := variables["tenant_id"].(string)
	userID := tenantUserID(variables)
	envID, ok := variables["env_id"].(int)
	if !ok {
		if envIDFloat, ok := variables["env_id"].(float64); ok {
			envID = int(envIDFloat)
		} else {
			envID = 10 // デフォルト値
		}
	}
	roleName := tenantRoleName(variables)
	roleNames := []string{roleName}

	return struct {
		TenantId string
		UserId   string
		EnvId    uint64
		Body     authapi.CreateTenantUserRolesParam
	}{
		TenantId: tenantID,
		UserId:   userID,
		EnvId:    uint64(envID),
		Body: authapi.CreateTenantUserRolesParam{
			RoleNames: roleNames,
		},
	}
}

// createTenantUserRolesWithBodyParams は CreateTenantUserRolesWithBody メソッドのパラメータを返します
// メソッドシグネチャ: CreateTenantUserRolesWithBody(ctx context.Context, tenantId string, userId string, envId uint64, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
func createTenantUserRolesWithBodyParams(variables map[string]interface{}) interface{} {
	tenantID, _ := variables["tenant_id"].(string)
	userID := tenantUserID(variables)
	envID, ok := variables["env_id"].(int)
	if !ok {
		if envIDFloat, ok := variables["env_id"].(float64); ok {
			envID = int(envIDFloat)
		} else {
			envID = 10 // デフォルト値
		}
	}
	roleName := tenantRoleName(variables)
	roleNames := []string{roleName}

	param := authapi.CreateTenantUserRolesParam{
		RoleNames: roleNames,
	}
	bodyBytes, _ := json.Marshal(param)

	return struct {
		TenantId    string
		UserId      string
		EnvId       uint64
		ContentType string
		Body        io.Reader
	}{
		TenantId:    tenantID,
		UserId:      userID,
		EnvId:       uint64(envID),
		ContentType: "application/json",
		Body:        bytes.NewReader(bodyBytes),
	}
}

// deleteTenantUserRoleParams は DeleteTenantUserRole メソッドのパラメータを返します
// メソッドシグネチャ: DeleteTenantUserRole(ctx context.Context, tenantId string, userId string, envId uint64, roleName string, reqEditors ...RequestEditorFn) (*http.Response, error)
func deleteTenantUserRoleParams(variables map[string]interface{}) interface{} {
	tenantID, _ := variables["tenant_id"].(string)
	userID := tenantUserID(variables)
	envID, ok := variables["env_id"].(int)
	if !ok {
		if envIDFloat, ok := variables["env_id"].(float64); ok {
			envID = int(envIDFloat)
		} else {
			envID = 10 // デフォルト値
		}
	}
	roleName, _ := variables["role_name"].(string)

	return struct {
		TenantId string
		UserId   string
		EnvId    uint64
		RoleName string
	}{
		TenantId: tenantID,
		UserId:   userID,
		EnvId:    uint64(envID),
		RoleName: roleName,
	}
}

// =============================================================================
// テナントプラン・請求関連のパラメータ関数
// =============================================================================

// updateTenantPlanParams は UpdateTenantPlan メソッドのパラメータを返します
// メソッドシグネチャ: UpdateTenantPlan(ctx context.Context, tenantId string, body UpdateTenantPlanParam, reqEditors ...RequestEditorFn) (*http.Response, error)
func updateTenantPlanParams(variables map[string]interface{}) interface{} {
	tenantID, _ := variables["tenant_id"].(string)

	return struct {
		TenantId string
		Body     authapi.UpdateTenantPlanParam
	}{
		TenantId: tenantID,
		Body:     authapi.UpdateTenantPlanParam{},
	}
}

// updateTenantPlanWithBodyParams は UpdateTenantPlanWithBody メソッドのパラメータを返します
// メソッドシグネチャ: UpdateTenantPlanWithBody(ctx context.Context, tenantId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
func updateTenantPlanWithBodyParams(variables map[string]interface{}) interface{} {
	tenantID, _ := variables["tenant_id"].(string)
	param := authapi.UpdateTenantPlanParam{}
	bodyBytes, _ := json.Marshal(param)

	return struct {
		TenantId    string
		ContentType string
		Body        io.Reader
	}{
		TenantId:    tenantID,
		ContentType: "application/json",
		Body:        bytes.NewReader(bodyBytes),
	}
}

// updateTenantBillingInfoParams は UpdateTenantBillingInfo メソッドのパラメータを返します
// メソッドシグネチャ: UpdateTenantBillingInfo(ctx context.Context, tenantId string, body UpdateTenantBillingInfoParam, reqEditors ...RequestEditorFn) (*http.Response, error)
func updateTenantBillingInfoParams(variables map[string]interface{}) interface{} {
	tenantID, _ := variables["tenant_id"].(string)
	info, _ := variables["tenant_billing_info"].(map[string]interface{})

	return struct {
		TenantId string
		Body     authapi.UpdateTenantBillingInfoParam
	}{
		TenantId: tenantID,
		Body:     buildTenantBillingInfoParam(info),
	}
}

// updateTenantBillingInfoWithBodyParams は UpdateTenantBillingInfoWithBody メソッドのパラメータを返します
// メソッドシグネチャ: UpdateTenantBillingInfoWithBody(ctx context.Context, tenantId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
func updateTenantBillingInfoWithBodyParams(variables map[string]interface{}) interface{} {
	tenantID, _ := variables["tenant_id"].(string)
	info, _ := variables["tenant_billing_info"].(map[string]interface{})
	param := buildTenantBillingInfoParam(info)
	bodyBytes, _ := json.Marshal(param)

	return struct {
		TenantId    string
		ContentType string
		Body        io.Reader
	}{
		TenantId:    tenantID,
		ContentType: "application/json",
		Body:        bytes.NewReader(bodyBytes),
	}
}

// =============================================================================
// Stripe連携関連のパラメータ関数
// =============================================================================

// getStripeInfoParams は GetStripeInfo メソッドのパラメータを返します
// メソッドシグネチャ: GetStripeInfo(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
func getStripeInfoParams(variables map[string]interface{}) interface{} {
	return nil
}

// updateStripeInfoParams は UpdateStripeInfo メソッドのパラメータを返します
// メソッドシグネチャ: UpdateStripeInfo(ctx context.Context, body UpdateStripeInfoParam, reqEditors ...RequestEditorFn) (*http.Response, error)
func updateStripeInfoParams(variables map[string]interface{}) interface{} {
	// Stripe関連の型はBilling APIで定義されているため、nilを返す
	return nil
}

// deleteStripeInfoParams は DeleteStripeInfo メソッドのパラメータを返します
// メソッドシグネチャ: DeleteStripeInfo(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
func deleteStripeInfoParams(variables map[string]interface{}) interface{} {
	return nil
}

// createTenantAndPricingParams は CreateTenantAndPricing メソッドのパラメータを返します
// メソッドシグネチャ: CreateTenantAndPricing(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
func createTenantAndPricingParams(variables map[string]interface{}) interface{} {
	return nil
}

// deleteStripeTenantAndPricingParams は DeleteStripeTenantAndPricing メソッドのパラメータを返します
// メソッドシグネチャ: DeleteStripeTenantAndPricing(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
func deleteStripeTenantAndPricingParams(variables map[string]interface{}) interface{} {
	return nil
}

// getStripeCustomerParams は GetStripeCustomer メソッドのパラメータを返します
// メソッドシグネチャ: GetStripeCustomer(ctx context.Context, tenantId string, reqEditors ...RequestEditorFn) (*http.Response, error)
func getStripeCustomerParams(variables map[string]interface{}) interface{} {
	tenantID, _ := variables["tenant_id"].(string)
	return tenantID
}

// resetPlanParams は ResetPlan メソッドのパラメータを返します
// メソッドシグネチャ: ResetPlan(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
func resetPlanParams(variables map[string]interface{}) interface{} {
	return nil
}

// =============================================================================
// 認証情報関連のパラメータ関数
// =============================================================================

// getUserInfoParams は GetUserInfo メソッドのパラメータを返します
// メソッドシグネチャ: GetUserInfo(ctx context.Context, params *GetUserInfoParams, reqEditors ...RequestEditorFn) (*http.Response, error)
func getUserInfoParams(variables map[string]interface{}) interface{} {
	token, _ := variables["token"].(string)
	return &authapi.GetUserInfoParams{
		Token: token,
	}
}

// getAuthCredentialsParams は GetAuthCredentials メソッドのパラメータを返します
// メソッドシグネチャ: GetAuthCredentials(ctx context.Context, params *GetAuthCredentialsParams, reqEditors ...RequestEditorFn) (*http.Response, error)
func getAuthCredentialsParams(variables map[string]interface{}) interface{} {
	return &authapi.GetAuthCredentialsParams{}
}

// createAuthCredentialsParams は CreateAuthCredentials メソッドのパラメータを返します
// メソッドシグネチャ: CreateAuthCredentials(ctx context.Context, body CreateAuthCredentialsParam, reqEditors ...RequestEditorFn) (*http.Response, error)
func createAuthCredentialsParams(variables map[string]interface{}) interface{} {
	return authapi.CreateAuthCredentialsParam{
		IdToken:      "test_id_token",
		AccessToken:  "test_access_token",
		RefreshToken: testdata.StringPtr("test_refresh_token"),
	}
}

// createAuthCredentialsWithBodyParams は CreateAuthCredentialsWithBody メソッドのパラメータを返します
// メソッドシグネチャ: CreateAuthCredentialsWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
func createAuthCredentialsWithBodyParams(variables map[string]interface{}) interface{} {
	param := authapi.CreateAuthCredentialsParam{
		IdToken:      "test_id_token",
		AccessToken:  "test_access_token",
		RefreshToken: testdata.StringPtr("test_refresh_token"),
	}
	bodyBytes, _ := json.Marshal(param)

	return struct {
		ContentType string
		Body        io.Reader
	}{
		ContentType: "application/json",
		Body:        bytes.NewReader(bodyBytes),
	}
}

// =============================================================================
// IDプロバイダー関連のパラメータ関数
// =============================================================================

// getIdentityProvidersParams は GetIdentityProviders メソッドのパラメータを返します
// メソッドシグネチャ: GetIdentityProviders(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
func getIdentityProvidersParams(variables map[string]interface{}) interface{} {
	return nil
}

// updateIdentityProviderParams は UpdateIdentityProvider メソッドのパラメータを返します
// メソッドシグネチャ: UpdateIdentityProvider(ctx context.Context, body UpdateIdentityProviderParam, reqEditors ...RequestEditorFn) (*http.Response, error)
func updateIdentityProviderParams(variables map[string]interface{}) interface{} {
	return buildUpdateIdentityProviderParam(variables)
}

// updateIdentityProviderWithBodyParams は UpdateIdentityProviderWithBody メソッドのパラメータを返します
// メソッドシグネチャ: UpdateIdentityProviderWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
func updateIdentityProviderWithBodyParams(variables map[string]interface{}) interface{} {
	param := buildUpdateIdentityProviderParam(variables)
	bodyBytes, _ := json.Marshal(param)

	return struct {
		ContentType string
		Body        io.Reader
	}{
		ContentType: "application/json",
		Body:        bytes.NewReader(bodyBytes),
	}
}

// getTenantIdentityProvidersParams は GetTenantIdentityProviders メソッドのパラメータを返します
// メソッドシグネチャ: GetTenantIdentityProviders(ctx context.Context, tenantId string, reqEditors ...RequestEditorFn) (*http.Response, error)
func getTenantIdentityProvidersParams(variables map[string]interface{}) interface{} {
	tenantID, _ := variables["tenant_id"].(string)
	return tenantID
}

// updateTenantIdentityProviderParams は UpdateTenantIdentityProvider メソッドのパラメータを返します
// メソッドシグネチャ: UpdateTenantIdentityProvider(ctx context.Context, tenantId string, body UpdateTenantIdentityProviderParam, reqEditors ...RequestEditorFn) (*http.Response, error)
func updateTenantIdentityProviderParams(variables map[string]interface{}) interface{} {
	tenantID, _ := variables["tenant_id"].(string)

	return struct {
		TenantId string
		Body     authapi.UpdateTenantIdentityProviderParam
	}{
		TenantId: tenantID,
		Body:     buildUpdateTenantIdentityProviderParam(variables),
	}
}

func buildUpdateTenantIdentityProviderParam(variables map[string]interface{}) authapi.UpdateTenantIdentityProviderParam {
	if raw, ok := variables["tenant_identity_provider_update_params"].(map[string]any); ok {
		param := authapi.UpdateTenantIdentityProviderParam{}
		param.ProviderType = normalizeProviderType(getStringFromMap(raw, "provider_type", ""))
		if propsRaw, ok := raw["identity_provider_props"].(map[string]any); ok {
			var tenantProps authapi.TenantIdentityProviderProps
			saml := authapi.IdentityProviderSaml{
				EmailAttribute: getStringFromMap(propsRaw, "email_attribute", "email"),
				MetadataUrl:    getStringFromMap(propsRaw, "metadata_url", "https://example.com/metadata.xml"),
			}
			if err := tenantProps.FromIdentityProviderSaml(saml); err == nil {
				param.IdentityProviderProps = &tenantProps
			} else {
				fmt.Printf("Warning: failed to encode tenant identity provider props: %v\n", err)
			}
		}
		return param
	}
	return authapi.UpdateTenantIdentityProviderParam{
		ProviderType: authapi.SAML,
	}
}

func buildUpdateIdentityProviderParam(variables map[string]interface{}) authapi.UpdateIdentityProviderParam {
	if raw, ok := variables["identity_provider_update_params"].(map[string]any); ok {
		var param authapi.UpdateIdentityProviderParam
		if err := decodeMapToStruct(raw, &param); err == nil {
			return param
		} else {
			fmt.Printf("Warning: failed to decode identity provider params: %v\n", err)
		}
	}
	return authapi.UpdateIdentityProviderParam{
		Provider: authapi.Google,
	}
}

func decodeMapToStruct(src map[string]any, dest interface{}) error {
	data, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

func getStringFromMap(m map[string]any, key, def string) string {
	if m == nil {
		return def
	}
	if v, ok := m[key].(string); ok && v != "" {
		return v
	}
	return def
}

func normalizeProviderType(v string) authapi.ProviderType {
	switch strings.ToUpper(v) {
	case "SAML", "GOOGLE":
		return authapi.SAML
	default:
		return authapi.SAML
	}
}

// updateTenantIdentityProviderWithBodyParams は UpdateTenantIdentityProviderWithBody メソッドのパラメータを返します
// メソッドシグネチャ: UpdateTenantIdentityProviderWithBody(ctx context.Context, tenantId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
func updateTenantIdentityProviderWithBodyParams(variables map[string]interface{}) interface{} {
	tenantID, _ := variables["tenant_id"].(string)
	param := buildUpdateTenantIdentityProviderParam(variables)
	bodyBytes, _ := json.Marshal(param)

	return struct {
		TenantId    string
		ContentType string
		Body        io.Reader
	}{
		TenantId:    tenantID,
		ContentType: "application/json",
		Body:        bytes.NewReader(bodyBytes),
	}
}

// unlinkProviderParams は UnlinkProvider メソッドのパラメータを返します
// メソッドシグネチャ: UnlinkProvider(ctx context.Context, providerName string, userId string, reqEditors ...RequestEditorFn) (*http.Response, error)
func unlinkProviderParams(variables map[string]interface{}) interface{} {
	userID, _ := variables["user_id"].(string)
	return struct {
		ProviderName string
		UserId       string
	}{
		ProviderName: "Google",
		UserId:       userID,
	}
}

// =============================================================================
// AWS Marketplace関連のパラメータ関数
// =============================================================================

// signUpWithAwsMarketplaceParams は SignUpWithAwsMarketplace メソッドのパラメータを返します
// メソッドシグネチャ: SignUpWithAwsMarketplace(ctx context.Context, body SignUpWithAwsMarketplaceParam, reqEditors ...RequestEditorFn) (*http.Response, error)
func signUpWithAwsMarketplaceParams(variables map[string]interface{}) interface{} {
	email := getStringVariable(variables, awsMarketplaceEmailVariable, "aws-marketplace@example.com")
	token := getStringVariable(variables, awsMarketplaceTokenVariable, "test_registration_token")

	return authapi.SignUpWithAwsMarketplaceParam{
		Email:             email,
		RegistrationToken: token,
	}
}

// signUpWithAwsMarketplaceWithBodyParams は SignUpWithAwsMarketplaceWithBody メソッドのパラメータを返します
// メソッドシグネチャ: SignUpWithAwsMarketplaceWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
func signUpWithAwsMarketplaceWithBodyParams(variables map[string]interface{}) interface{} {
	email := getStringVariable(variables, awsMarketplaceEmailVariable, "aws-marketplace@example.com")
	token := getStringVariable(variables, awsMarketplaceTokenVariable, "test_registration_token")

	param := authapi.SignUpWithAwsMarketplaceParam{
		Email:             email,
		RegistrationToken: token,
	}
	bodyBytes, _ := json.Marshal(param)

	return struct {
		ContentType string
		Body        io.Reader
	}{
		ContentType: "application/json",
		Body:        bytes.NewReader(bodyBytes),
	}
}

// getListingStatusParams は GetListingStatus メソッドのパラメータを返します
// メソッドシグネチャ: GetListingStatus(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
func getListingStatusParams(variables map[string]interface{}) interface{} {
	return nil
}

// updateListingStatusParams は UpdateListingStatus メソッドのパラメータを返します
// メソッドシグネチャ: UpdateListingStatus(ctx context.Context, body UpdateListingStatusParam, reqEditors ...RequestEditorFn) (*http.Response, error)
func updateListingStatusParams(variables map[string]interface{}) interface{} {
	// AWS Marketplace関連の型はAWS Marketplace APIで定義されているため、nilを返す
	return nil
}

// updateSettingsParams は UpdateSettings メソッドのパラメータを返します
// メソッドシグネチャ: UpdateSettings(ctx context.Context, body UpdateSettingsParam, reqEditors ...RequestEditorFn) (*http.Response, error)
func updateSettingsParams(variables map[string]interface{}) interface{} {
	// AWS Marketplace関連の型はAWS Marketplace APIで定義されているため、nilを返す
	return nil
}

// =============================================================================
// シングルテナント関連のパラメータ関数
// =============================================================================

// getCloudFormationLaunchStackLinkForSingleTenantParams は GetCloudFormationLaunchStackLinkForSingleTenant メソッドのパラメータを返します
// メソッドシグネチャ: GetCloudFormationLaunchStackLinkForSingleTenant(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
func getCloudFormationLaunchStackLinkForSingleTenantParams(variables map[string]interface{}) interface{} {
	return nil
}

// getSingleTenantSettingsParams は GetSingleTenantSettings メソッドのパラメータを返します
// メソッドシグネチャ: GetSingleTenantSettings(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
func getSingleTenantSettingsParams(variables map[string]interface{}) interface{} {
	return nil
}

// updateSingleTenantSettingsParams は UpdateSingleTenantSettings メソッドのパラメータを返します
// メソッドシグネチャ: UpdateSingleTenantSettings(ctx context.Context, body UpdateSingleTenantSettingsParam, reqEditors ...RequestEditorFn) (*http.Response, error)
func updateSingleTenantSettingsParams(variables map[string]interface{}) interface{} {
	return authapi.UpdateSingleTenantSettingsParam{
		Enabled: testdata.BoolPtr(false),
	}
}

// updateSingleTenantSettingsWithBodyParams は UpdateSingleTenantSettingsWithBody メソッドのパラメータを返します
// メソッドシグネチャ: UpdateSingleTenantSettingsWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
func updateSingleTenantSettingsWithBodyParams(variables map[string]interface{}) interface{} {
	param := authapi.UpdateSingleTenantSettingsParam{
		Enabled: testdata.BoolPtr(false),
	}
	bodyBytes, _ := json.Marshal(param)

	return struct {
		ContentType string
		Body        io.Reader
	}{
		ContentType: "application/json",
		Body:        bytes.NewReader(bodyBytes),
	}
}

// =============================================================================
// Pricing/Billing API メソッド（Auth APIから呼び出し可能）
// =============================================================================

func ensurePricingResourceName(variables map[string]interface{}, key, prefix string) string {
	if value, ok := variables[key].(string); ok && value != "" {
		return value
	}
	generated := fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano())
	variables[key] = generated
	return generated
}

func requirePricingVariable(variables map[string]interface{}, key, prefix string) string {
	if value, ok := variables[key].(string); ok && value != "" {
		return value
	}
	generated := ensurePricingResourceName(variables, key, prefix)
	fmt.Printf("[pricing] Warning: variable %s was empty; generated placeholder %s. Ensure dependent resources exist before calling downstream APIs.\n", key, generated)
	return generated
}

func pricingStringOrDefault(variables map[string]interface{}, key, fallback string) string {
	if value, ok := variables[key].(string); ok && value != "" {
		return value
	}
	return fallback
}

func pricingIntOrDefault(variables map[string]interface{}, key string, fallback int) int {
	maxInt := func() int64 {
		if strconv.IntSize == 32 {
			return int64(math.MaxInt32)
		}
		return math.MaxInt64
	}()
	if raw, ok := variables[key]; ok {
		switch v := raw.(type) {
		case int:
			if v >= 0 {
				return v
			}
		case int64:
			if v >= 0 && v <= maxInt {
				return int(v)
			}
		case float64:
			if v >= 0 && v <= float64(maxInt) {
				return int(v)
			}
		case json.Number:
			if parsed, err := strconv.ParseInt(string(v), 10, strconv.IntSize); err == nil && parsed >= 0 {
				return int(parsed)
			}
		case string:
			if parsed, err := strconv.ParseInt(v, 10, strconv.IntSize); err == nil && parsed >= 0 {
				return int(parsed)
			}
		}
	}
	if fallback < 0 {
		return 0
	}
	return fallback
}

func pricingUint64OrDefault(variables map[string]interface{}, key string, fallback uint64) uint64 {
	if raw, ok := variables[key]; ok {
		switch v := raw.(type) {
		case int:
			if v >= 0 {
				return uint64(v)
			}
		case int64:
			if v >= 0 {
				return uint64(v)
			}
		case float64:
			if v >= 0 && v <= float64(^uint64(0)) {
				return uint64(v)
			}
		case json.Number:
			// ParseUintを使用して直接uint64に変換（より安全）
			if parsed, err := strconv.ParseUint(string(v), 10, 64); err == nil {
				return parsed
			}
		case string:
			if parsed, err := strconv.ParseUint(v, 10, 64); err == nil {
				return parsed
			}
		}
	}
	return fallback
}

func pricingAggregatePtr(value pricingapi.AggregateUsage) *pricingapi.AggregateUsage {
	copy := value
	return &copy
}

func collectUUIDs(variables map[string]interface{}, pluralKey, singleKey string) []pricingapi.Uuid {
	var ids []pricingapi.Uuid
	if raw, ok := variables[pluralKey]; ok {
		switch v := raw.(type) {
		case []string:
			for _, id := range v {
				if id != "" {
					ids = append(ids, pricingapi.Uuid(id))
				}
			}
		case []interface{}:
			for _, val := range v {
				if str, ok := val.(string); ok && str != "" {
					ids = append(ids, pricingapi.Uuid(str))
				}
			}
		}
	}
	if id, ok := variables[singleKey].(string); ok && id != "" {
		ids = append(ids, pricingapi.Uuid(id))
	}
	return ids
}

func defaultPricingTiers() []pricingapi.PricingTier {
	return []pricingapi.PricingTier{
		{UpTo: 100, UnitAmount: 300, FlatAmount: 0},
		{UpTo: 500, UnitAmount: 250, FlatAmount: 0},
		{Inf: true, UnitAmount: 200, FlatAmount: 0},
	}
}

func buildFixedPricingUnit(name, displayName, description string, currency pricingapi.Currency, variables map[string]interface{}) pricingapi.PricingUnitForSave {
	interval := pricingapi.RecurringInterval(pricingStringOrDefault(variables, "pricing_unit_recurring_interval", string(pricingapi.Month)))
	unitAmount := pricingUint64OrDefault(variables, "pricing_unit_unit_amount", 3000)
	payload := pricingapi.PricingFixedUnitForSave{
		Type:              pricingapi.Fixed,
		Name:              name,
		DisplayName:       displayName,
		Description:       description,
		Currency:          currency,
		RecurringInterval: interval,
		UnitAmount:        unitAmount,
	}
	var unit pricingapi.PricingUnitForSave
	if err := unit.FromPricingFixedUnitForSave(payload); err != nil {
		panic(fmt.Sprintf("failed to encode fixed pricing unit payload: %v", err))
	}
	return unit
}

func buildUsagePricingUnit(name, displayName, description string, currency pricingapi.Currency, meteringUnitName string, variables map[string]interface{}) pricingapi.PricingUnitForSave {
	aggregate := pricingapi.AggregateUsage(pricingStringOrDefault(variables, "pricing_unit_aggregate_usage", string(pricingapi.Sum)))
	unitAmount := pricingUint64OrDefault(variables, "pricing_unit_unit_amount", 1200)
	upperCount := pricingUint64OrDefault(variables, "pricing_unit_upper_count", 10000)
	payload := pricingapi.PricingUsageUnitForSave{
		Type:             pricingapi.Usage,
		Name:             name,
		DisplayName:      displayName,
		Description:      description,
		Currency:         currency,
		MeteringUnitName: meteringUnitName,
		UnitAmount:       unitAmount,
		UpperCount:       upperCount,
		AggregateUsage:   pricingAggregatePtr(aggregate),
	}
	var unit pricingapi.PricingUnitForSave
	if err := unit.FromPricingUsageUnitForSave(payload); err != nil {
		panic(fmt.Sprintf("failed to encode usage pricing unit payload: %v", err))
	}
	return unit
}

func buildTieredPricingUnit(name, displayName, description string, currency pricingapi.Currency, meteringUnitName string, variables map[string]interface{}) pricingapi.PricingUnitForSave {
	aggregate := pricingapi.AggregateUsage(pricingStringOrDefault(variables, "pricing_unit_aggregate_usage", string(pricingapi.Sum)))
	upperCount := pricingUint64OrDefault(variables, "pricing_unit_upper_count", 10000)
	payload := pricingapi.PricingTieredUnitForSave{
		Type:             pricingapi.Tiered,
		Name:             name,
		DisplayName:      displayName,
		Description:      description,
		Currency:         currency,
		MeteringUnitName: meteringUnitName,
		Tiers:            defaultPricingTiers(),
		UpperCount:       upperCount,
		AggregateUsage:   pricingAggregatePtr(aggregate),
	}
	var unit pricingapi.PricingUnitForSave
	if err := unit.FromPricingTieredUnitForSave(payload); err != nil {
		panic(fmt.Sprintf("failed to encode tiered pricing unit payload: %v", err))
	}
	return unit
}

func buildTieredUsagePricingUnit(name, displayName, description string, currency pricingapi.Currency, meteringUnitName string, variables map[string]interface{}) pricingapi.PricingUnitForSave {
	aggregate := pricingapi.AggregateUsage(pricingStringOrDefault(variables, "pricing_unit_aggregate_usage", string(pricingapi.Sum)))
	upperCount := pricingUint64OrDefault(variables, "pricing_unit_upper_count", 10000)
	payload := pricingapi.PricingTieredUsageUnitForSave{
		Type:             pricingapi.TieredUsage,
		Name:             name,
		DisplayName:      displayName,
		Description:      description,
		Currency:         currency,
		MeteringUnitName: meteringUnitName,
		Tiers:            defaultPricingTiers(),
		UpperCount:       upperCount,
		AggregateUsage:   pricingAggregatePtr(aggregate),
	}
	var unit pricingapi.PricingUnitForSave
	if err := unit.FromPricingTieredUsageUnitForSave(payload); err != nil {
		panic(fmt.Sprintf("failed to encode tiered usage pricing unit payload: %v", err))
	}
	return unit
}

func resolveBillingTenantID(variables map[string]interface{}) string {
	if tenantID, ok := variables["tenant_id"].(string); ok && tenantID != "" {
		return tenantID
	}
	if tenantID, ok := variables["tenantId"].(string); ok && tenantID != "" {
		return tenantID
	}
	generated := fmt.Sprintf("tenant-%d", time.Now().Unix())
	fmt.Printf("[pricing] Warning: tenant_id was not set; using placeholder %s.\n", generated)
	variables["tenant_id"] = generated
	return generated
}

// createPricingUnitParams は CreatePricingUnit メソッドのパラメータを返します
func createPricingUnitParams(variables map[string]interface{}) interface{} {
	unitType := pricingapi.UnitType(pricingStringOrDefault(variables, "pricing_unit_type", string(pricingapi.TieredUsage)))
	name := ensurePricingResourceName(variables, "pricing_unit_name", "pricing-unit")
	displayName := pricingStringOrDefault(variables, "pricing_unit_display_name", "Tiered Usage Unit")
	description := pricingStringOrDefault(variables, "pricing_unit_description", "Tiered usage billing for Stripe integration tests")
	currency := pricingapi.Currency(pricingStringOrDefault(variables, "pricing_unit_currency", string(pricingapi.JPY)))

	var unit pricingapi.PricingUnitForSave
	switch unitType {
	case pricingapi.Fixed:
		unit = buildFixedPricingUnit(name, displayName, description, currency, variables)
	case pricingapi.Usage:
		meteringUnitName := requirePricingVariable(variables, "metering_unit_name", "meter")
		unit = buildUsagePricingUnit(name, displayName, description, currency, meteringUnitName, variables)
	case pricingapi.Tiered:
		meteringUnitName := requirePricingVariable(variables, "metering_unit_name", "meter")
		unit = buildTieredPricingUnit(name, displayName, description, currency, meteringUnitName, variables)
	default:
		meteringUnitName := requirePricingVariable(variables, "metering_unit_name", "meter")
		unitType = pricingapi.TieredUsage
		unit = buildTieredUsagePricingUnit(name, displayName, description, currency, meteringUnitName, variables)
	}

	variables["pricing_unit_type"] = string(unitType)
	return pricingapi.CreatePricingUnitParam(unit)
}

// getPricingUnitsParams は GetPricingUnits メソッドのパラメータを返します
func getPricingUnitsParams(variables map[string]interface{}) interface{} {
	return nil
}

// createPricingMenuParams は CreatePricingMenu メソッドのパラメータを返します
func createPricingMenuParams(variables map[string]interface{}) interface{} {
	name := ensurePricingResourceName(variables, "pricing_menu_name", "pricing-menu")
	displayName := pricingStringOrDefault(variables, "pricing_menu_display_name", "Standard Pricing Menu")
	description := pricingStringOrDefault(variables, "pricing_menu_description", "Menu that bundles pricing units for Stripe plan creation")
	unitIDs := collectUUIDs(variables, "pricing_unit_ids", "pricing_unit_id")
	if len(unitIDs) == 0 {
		id := requirePricingVariable(variables, "pricing_unit_id", "pricing-unit")
		unitIDs = []pricingapi.Uuid{pricingapi.Uuid(id)}
	}
	return pricingapi.CreatePricingMenuParam{
		Name:        name,
		DisplayName: displayName,
		Description: description,
		UnitIds:     unitIDs,
	}
}

// getPricingMenusParams は GetPricingMenus メソッドのパラメータを返します
func getPricingMenusParams(variables map[string]interface{}) interface{} {
	return nil
}

// createPricingPlanParams は CreatePricingPlan メソッドのパラメータを返します
func createPricingPlanParams(variables map[string]interface{}) interface{} {
	name := ensurePricingResourceName(variables, "pricing_plan_name", "pricing-plan")
	displayName := pricingStringOrDefault(variables, "pricing_plan_display_name", "Sample Pricing Plan")
	description := pricingStringOrDefault(variables, "pricing_plan_description", "Plan that links pricing menus for tenant onboarding")
	menuIDs := collectUUIDs(variables, "pricing_menu_ids", "pricing_menu_id")
	if len(menuIDs) == 0 {
		id := requirePricingVariable(variables, "pricing_menu_id", "pricing-menu")
		menuIDs = []pricingapi.Uuid{pricingapi.Uuid(id)}
	}
	return pricingapi.CreatePricingPlanParam{
		Name:        name,
		DisplayName: displayName,
		Description: description,
		MenuIds:     menuIDs,
	}
}

// getPricingPlansParams は GetPricingPlans メソッドのパラメータを返します
func getPricingPlansParams(variables map[string]interface{}) interface{} {
	return nil
}

// getPricingPlanParams は GetPricingPlan メソッドのパラメータを返します
func getPricingPlanParams(variables map[string]interface{}) interface{} {
	planID, _ := variables["plan_id"].(string)
	return planID
}

// createTaxRateParams は CreateTaxRate メソッドのパラメータを返します
func createTaxRateParams(variables map[string]interface{}) interface{} {
	return pricingapi.CreateTaxRateParam{
		Name:        pricingStringOrDefault(variables, "tax_rate_name", "japanese_consumption_tax_inclusive"),
		DisplayName: pricingStringOrDefault(variables, "tax_rate_display_name", "日本の消費税(内税)"),
		Description: pricingStringOrDefault(variables, "tax_rate_description", "日本国内向けの10%内税設定"),
		Country:     pricingStringOrDefault(variables, "tax_rate_country", "JP"),
		Inclusive:   true,
		Percentage:  10,
	}
}

// getTaxRatesParams は GetTaxRates メソッドのパラメータを返します
func getTaxRatesParams(variables map[string]interface{}) interface{} {
	return nil
}

// createMeteringUnitParams は CreateMeteringUnit メソッドのパラメータを返します
func createMeteringUnitParams(variables map[string]interface{}) interface{} {
	name := ensurePricingResourceName(variables, "metering_unit_name", "meter")
	displayName := pricingStringOrDefault(variables, "metering_unit_display_name", "Active Users")
	description := pricingStringOrDefault(variables, "metering_unit_description", "Counts active users per tenant")
	aggregate := pricingapi.AggregateUsage(pricingStringOrDefault(variables, "metering_unit_aggregate_usage", string(pricingapi.Sum)))
	return pricingapi.CreateMeteringUnitParam{
		UnitName:       name,
		DisplayName:    displayName,
		Description:    description,
		AggregateUsage: pricingAggregatePtr(aggregate),
	}
}

// updateMeteringUnitTimestampCountParams は UpdateMeteringUnitTimestampCount メソッドのパラメータを返します
func updateMeteringUnitTimestampCountParams(variables map[string]interface{}) interface{} {
	tenantID := resolveBillingTenantID(variables)
	meteringUnitName := requirePricingVariable(variables, "metering_unit_name", "meter")
	timestamp := pricingIntOrDefault(variables, "metering_timestamp", int(time.Now().Add(-24*time.Hour).Unix()))
	count := pricingIntOrDefault(variables, "metering_count_delta", 5)
	if count < 0 {
		count = 0
	}

	body := pricingapi.UpdateMeteringUnitTimestampCountJSONRequestBody{
		Method: pricingapi.UpdateMeteringUnitTimestampCountMethod("add"),
		Count:  count,
	}
	return struct {
		TenantId         string
		MeteringUnitName string
		Timestamp        int
		Body             pricingapi.UpdateMeteringUnitTimestampCountJSONRequestBody
	}{
		TenantId:         tenantID,
		MeteringUnitName: meteringUnitName,
		Timestamp:        timestamp,
		Body:             body,
	}
}

// getMeteringUnitDateCountsByTenantIdAndDateParams は GetMeteringUnitDateCountsByTenantIdAndDate メソッドのパラメータを返します
func getMeteringUnitDateCountsByTenantIdAndDateParams(variables map[string]interface{}) interface{} {
	date := pricingStringOrDefault(variables, "metering_unit_count_date", time.Now().UTC().Format("2006-01-02"))
	return struct {
		TenantId string
		Date     string
	}{
		TenantId: resolveBillingTenantID(variables),
		Date:     date,
	}
}

// deleteAllPlansAndMenusAndUnitsAndMetersAndTaxRatesParams は DeleteAllPlansAndMenusAndUnitsAndMetersAndTaxRates メソッドのパラメータを返します
func deleteAllPlansAndMenusAndUnitsAndMetersAndTaxRatesParams(variables map[string]interface{}) interface{} {
	return nil
}

// ========================================
// SaaS User Attributes Management
// ========================================

// updateSaasUserAttributesParams は UpdateSaasUserAttributes メソッドのパラメータを返します
// メソッドシグネチャ: UpdateSaasUserAttributes(ctx context.Context, userId string, body UpdateSaasUserAttributesParam, reqEditors ...RequestEditorFn) (*http.Response, error)
func updateSaasUserAttributesParams(variables map[string]interface{}) interface{} {
	userID, _ := variables["user_id"].(string)
	attributes, _ := variables["saas_user_attributes"].(map[string]interface{})

	return struct {
		UserId string
		Body   authapi.UpdateSaasUserAttributesJSONRequestBody
	}{
		UserId: userID,
		Body: authapi.UpdateSaasUserAttributesJSONRequestBody{
			Attributes: attributes,
		},
	}
}

// updateSaasUserAttributesWithBodyParams は UpdateSaasUserAttributesWithBody メソッドのパラメータを返します
// メソッドシグネチャ: UpdateSaasUserAttributesWithBody(ctx context.Context, userId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
func updateSaasUserAttributesWithBodyParams(variables map[string]interface{}) interface{} {
	userID, _ := variables["user_id"].(string)
	attributes, _ := variables["saas_user_attributes"].(map[string]interface{})

	param := authapi.UpdateSaasUserAttributesParam{
		Attributes: attributes,
	}
	body, _ := json.Marshal(param)

	return struct {
		UserId      string
		ContentType string
		Body        io.Reader
	}{
		UserId:      userID,
		ContentType: "application/json",
		Body:        bytes.NewReader(body),
	}
}

// createSaasUserAttributeParams は CreateSaasUserAttribute メソッドのパラメータを返します
// メソッドシグネチャ: CreateSaasUserAttribute(ctx context.Context, body CreateSaasUserAttributeParam, reqEditors ...RequestEditorFn) (*http.Response, error)
func createSaasUserAttributeParams(variables map[string]interface{}) interface{} {
	attributeName, _ := variables["attribute_name"].(string)
	displayName, _ := variables["display_name"].(string)
	attributeType, _ := variables["attribute_type"].(string)

	return authapi.CreateSaasUserAttributeJSONRequestBody{
		AttributeName: attributeName,
		DisplayName:   displayName,
		AttributeType: authapi.AttributeType(attributeType),
	}
}

// createSaasUserAttributeWithBodyParams は CreateSaasUserAttributeWithBody メソッドのパラメータを返します
// メソッドシグネチャ: CreateSaasUserAttributeWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
func createSaasUserAttributeWithBodyParams(variables map[string]interface{}) interface{} {
	attributeName, _ := variables["attribute_name"].(string)
	displayName, _ := variables["display_name"].(string)
	attributeType, _ := variables["attribute_type"].(string)

	param := authapi.CreateSaasUserAttributeParam{
		AttributeName: attributeName,
		DisplayName:   displayName,
		AttributeType: authapi.AttributeType(attributeType),
	}
	body, _ := json.Marshal(param)

	return struct {
		ContentType string
		Body        io.Reader
	}{
		ContentType: "application/json",
		Body:        bytes.NewReader(body),
	}
}

// ========================================
// Tenant Invitations Management
// ========================================

// getTenantInvitationsParams は GetTenantInvitations メソッドのパラメータを返します
// メソッドシグネチャ: GetTenantInvitations(ctx context.Context, tenantId string, reqEditors ...RequestEditorFn) (*http.Response, error)
func getTenantInvitationsParams(variables map[string]interface{}) interface{} {
	tenantID, _ := variables["tenant_id"].(string)
	return tenantID
}

// createTenantInvitationParams は CreateTenantInvitation メソッドのパラメータを返します
// メソッドシグネチャ: CreateTenantInvitation(ctx context.Context, tenantId string, body CreateTenantInvitationParam, reqEditors ...RequestEditorFn) (*http.Response, error)
func createTenantInvitationParams(variables map[string]interface{}) interface{} {
	tenantID, _ := variables["tenant_id"].(string)
	email, _ := variables["invitation_email"].(string)
	accessToken, _ := variables["access_token"].(string)
	envs, _ := variables["invitation_envs"].([]interface{})

	// envs を適切な型に変換
	var envsParam []struct {
		Id        authapi.Id `json:"id"`
		RoleNames []string   `json:"role_names"`
	}
	for _, env := range envs {
		envMap, _ := env.(map[string]interface{})
		envID, _ := envMap["id"].(float64)
		roleNames, _ := envMap["role_names"].([]interface{})

		var roles []string
		for _, role := range roleNames {
			roles = append(roles, role.(string))
		}

		envsParam = append(envsParam, struct {
			Id        authapi.Id `json:"id"`
			RoleNames []string   `json:"role_names"`
		}{
			Id:        authapi.Id(envID),
			RoleNames: roles,
		})
	}

	return struct {
		TenantId string
		Body     authapi.CreateTenantInvitationJSONRequestBody
	}{
		TenantId: tenantID,
		Body: authapi.CreateTenantInvitationJSONRequestBody{
			AccessToken: accessToken,
			Email:       email,
			Envs:        envsParam,
		},
	}
}

// createTenantInvitationWithBodyParams は CreateTenantInvitationWithBody メソッドのパラメータを返します
// メソッドシグネチャ: CreateTenantInvitationWithBody(ctx context.Context, tenantId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
func createTenantInvitationWithBodyParams(variables map[string]interface{}) interface{} {
	tenantID, _ := variables["tenant_id"].(string)
	email, _ := variables["invitation_email"].(string)
	accessToken, _ := variables["access_token"].(string)
	envs, _ := variables["invitation_envs"].([]interface{})

	// envs を適切な型に変換
	var envsParam []struct {
		Id        authapi.Id `json:"id"`
		RoleNames []string   `json:"role_names"`
	}
	for _, env := range envs {
		envMap, _ := env.(map[string]interface{})
		envID, _ := envMap["id"].(float64)
		roleNames, _ := envMap["role_names"].([]interface{})

		var roles []string
		for _, role := range roleNames {
			roles = append(roles, role.(string))
		}

		envsParam = append(envsParam, struct {
			Id        authapi.Id `json:"id"`
			RoleNames []string   `json:"role_names"`
		}{
			Id:        authapi.Id(envID),
			RoleNames: roles,
		})
	}

	param := authapi.CreateTenantInvitationParam{
		AccessToken: accessToken,
		Email:       email,
		Envs:        envsParam,
	}
	bodyBytes, _ := json.Marshal(param)

	return struct {
		TenantId    string
		ContentType string
		Body        io.Reader
	}{
		TenantId:    tenantID,
		ContentType: "application/json",
		Body:        bytes.NewReader(bodyBytes),
	}
}

// getTenantInvitationParams は GetTenantInvitation メソッドのパラメータを返します
// メソッドシグネチャ: GetTenantInvitation(ctx context.Context, tenantId string, invitationId string, reqEditors ...RequestEditorFn) (*http.Response, error)
func getTenantInvitationParams(variables map[string]interface{}) interface{} {
	tenantID, _ := variables["tenant_id"].(string)
	invitationID, _ := variables["invitation_id"].(string)

	return struct {
		TenantId     string
		InvitationId string
	}{
		TenantId:     tenantID,
		InvitationId: invitationID,
	}
}

// deleteTenantInvitationParams は DeleteTenantInvitation メソッドのパラメータを返します
// メソッドシグネチャ: DeleteTenantInvitation(ctx context.Context, tenantId string, invitationId string, reqEditors ...RequestEditorFn) (*http.Response, error)
func deleteTenantInvitationParams(variables map[string]interface{}) interface{} {
	tenantID, _ := variables["tenant_id"].(string)
	invitationID, _ := variables["invitation_id"].(string)

	return struct {
		TenantId     string
		InvitationId string
	}{
		TenantId:     tenantID,
		InvitationId: invitationID,
	}
}

// getInvitationValidityParams は GetInvitationValidity メソッドのパラメータを返します
// メソッドシグネチャ: GetInvitationValidity(ctx context.Context, invitationId string, reqEditors ...RequestEditorFn) (*http.Response, error)
func getInvitationValidityParams(variables map[string]interface{}) interface{} {
	invitationID, _ := variables["invitation_id"].(string)
	return invitationID
}

// validateInvitationParams は ValidateInvitation メソッドのパラメータを返します
// メソッドシグネチャ: ValidateInvitation(ctx context.Context, invitationId string, body ValidateInvitationParam, reqEditors ...RequestEditorFn) (*http.Response, error)
func validateInvitationParams(variables map[string]interface{}) interface{} {
	invitationID, _ := variables["invitation_id"].(string)
	accessToken, _ := variables["access_token"].(string)
	email, _ := variables["invitation_email"].(string)
	password, _ := variables["password"].(string)

	return struct {
		InvitationId string
		Body         authapi.ValidateInvitationJSONRequestBody
	}{
		InvitationId: invitationID,
		Body: authapi.ValidateInvitationJSONRequestBody{
			AccessToken: &accessToken,
			Email:       &email,
			Password:    &password,
		},
	}
}

// validateInvitationWithBodyParams は ValidateInvitationWithBody メソッドのパラメータを返します
// メソッドシグネチャ: ValidateInvitationWithBody(ctx context.Context, invitationId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
func validateInvitationWithBodyParams(variables map[string]interface{}) interface{} {
	invitationID, _ := variables["invitation_id"].(string)
	accessToken, _ := variables["access_token"].(string)
	email, _ := variables["invitation_email"].(string)
	password, _ := variables["password"].(string)

	param := authapi.ValidateInvitationParam{
		AccessToken: &accessToken,
		Email:       &email,
		Password:    &password,
	}
	bodyBytes, _ := json.Marshal(param)

	return struct {
		InvitationId string
		ContentType  string
		Body         io.Reader
	}{
		InvitationId: invitationID,
		ContentType:  "application/json",
		Body:         bytes.NewReader(bodyBytes),
	}
}

// ========================================
// External User Link and Email Update
// ========================================

// requestExternalUserLinkParams は RequestExternalUserLink メソッドのパラメータを返します
// メソッドシグネチャ: RequestExternalUserLink(ctx context.Context, body RequestExternalUserLinkParam, reqEditors ...RequestEditorFn) (*http.Response, error)
func requestExternalUserLinkParams(variables map[string]interface{}) interface{} {
	accessToken := getCognitoAccessToken(variables)

	return authapi.RequestExternalUserLinkJSONRequestBody{
		AccessToken: accessToken,
	}
}

// requestExternalUserLinkWithBodyParams は RequestExternalUserLinkWithBody メソッドのパラメータを返します
// メソッドシグネチャ: RequestExternalUserLinkWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
func requestExternalUserLinkWithBodyParams(variables map[string]interface{}) interface{} {
	accessToken := getCognitoAccessToken(variables)

	param := authapi.RequestExternalUserLinkParam{
		AccessToken: accessToken,
	}
	bodyBytes, _ := json.Marshal(param)

	return struct {
		ContentType string
		Body        io.Reader
	}{
		ContentType: "application/json",
		Body:        bytes.NewReader(bodyBytes),
	}
}

// confirmExternalUserLinkParams は ConfirmExternalUserLink メソッドのパラメータを返します
// メソッドシグネチャ: ConfirmExternalUserLink(ctx context.Context, body ConfirmExternalUserLinkParam, reqEditors ...RequestEditorFn) (*http.Response, error)
func confirmExternalUserLinkParams(variables map[string]interface{}) interface{} {
	accessToken := getCognitoAccessToken(variables)
	code, _ := variables["verification_code"].(string)

	return authapi.ConfirmExternalUserLinkJSONRequestBody{
		AccessToken: accessToken,
		Code:        code,
	}
}

// confirmExternalUserLinkWithBodyParams は ConfirmExternalUserLinkWithBody メソッドのパラメータを返します
// メソッドシグネチャ: ConfirmExternalUserLinkWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
func confirmExternalUserLinkWithBodyParams(variables map[string]interface{}) interface{} {
	accessToken := getCognitoAccessToken(variables)
	code, _ := variables["verification_code"].(string)

	param := authapi.ConfirmExternalUserLinkParam{
		AccessToken: accessToken,
		Code:        code,
	}
	bodyBytes, _ := json.Marshal(param)

	return struct {
		ContentType string
		Body        io.Reader
	}{
		ContentType: "application/json",
		Body:        bytes.NewReader(bodyBytes),
	}
}

// requestEmailUpdateParams は RequestEmailUpdate メソッドのパラメータを返します
// メソッドシグネチャ: RequestEmailUpdate(ctx context.Context, userId string, body RequestEmailUpdateParam, reqEditors ...RequestEditorFn) (*http.Response, error)
func requestEmailUpdateParams(variables map[string]interface{}) interface{} {
	userID, _ := variables["user_id"].(string)
	email, _ := variables["new_email"].(string)
	accessToken := getCognitoAccessToken(variables)

	return struct {
		UserId string
		Body   authapi.RequestEmailUpdateJSONRequestBody
	}{
		UserId: userID,
		Body: authapi.RequestEmailUpdateJSONRequestBody{
			AccessToken: accessToken,
			Email:       email,
		},
	}
}

// requestEmailUpdateWithBodyParams は RequestEmailUpdateWithBody メソッドのパラメータを返します
// メソッドシグネチャ: RequestEmailUpdateWithBody(ctx context.Context, userId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
func requestEmailUpdateWithBodyParams(variables map[string]interface{}) interface{} {
	userID, _ := variables["user_id"].(string)
	email, _ := variables["new_email"].(string)
	accessToken := getCognitoAccessToken(variables)

	param := authapi.RequestEmailUpdateParam{
		AccessToken: accessToken,
		Email:       email,
	}
	bodyBytes, _ := json.Marshal(param)

	return struct {
		UserId      string
		ContentType string
		Body        io.Reader
	}{
		UserId:      userID,
		ContentType: "application/json",
		Body:        bytes.NewReader(bodyBytes),
	}
}

// confirmEmailUpdateParams は ConfirmEmailUpdate メソッドのパラメータを返します
// メソッドシグネチャ: ConfirmEmailUpdate(ctx context.Context, userId string, body ConfirmEmailUpdateParam, reqEditors ...RequestEditorFn) (*http.Response, error)
func confirmEmailUpdateParams(variables map[string]interface{}) interface{} {
	userID, _ := variables["user_id"].(string)
	code, _ := variables["verification_code"].(string)
	accessToken := getCognitoAccessToken(variables)

	return struct {
		UserId string
		Body   authapi.ConfirmEmailUpdateJSONRequestBody
	}{
		UserId: userID,
		Body: authapi.ConfirmEmailUpdateJSONRequestBody{
			AccessToken: accessToken,
			Code:        code,
		},
	}
}

// confirmEmailUpdateWithBodyParams は ConfirmEmailUpdateWithBody メソッドのパラメータを返します
// メソッドシグネチャ: ConfirmEmailUpdateWithBody(ctx context.Context, userId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
func confirmEmailUpdateWithBodyParams(variables map[string]interface{}) interface{} {
	userID, _ := variables["user_id"].(string)
	code, _ := variables["verification_code"].(string)
	accessToken := getCognitoAccessToken(variables)

	param := authapi.ConfirmEmailUpdateParam{
		AccessToken: accessToken,
		Code:        code,
	}
	bodyBytes, _ := json.Marshal(param)

	return struct {
		UserId      string
		ContentType string
		Body        io.Reader
	}{
		UserId:      userID,
		ContentType: "application/json",
		Body:        bytes.NewReader(bodyBytes),
	}
}

// ========================================
// Sign Up and Provider Management
// ========================================

// signUpParams は SignUp メソッドのパラメータを返します
// メソッドシグネチャ: SignUp(ctx context.Context, body SignUpParam, reqEditors ...RequestEditorFn) (*http.Response, error)
func signUpParams(variables map[string]interface{}) interface{} {
	email, _ := variables["signup_email"].(string)

	return authapi.SignUpJSONRequestBody{
		Email: email,
	}
}

// signUpWithBodyParams は SignUpWithBody メソッドのパラメータを返します
// メソッドシグネチャ: SignUpWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
func signUpWithBodyParams(variables map[string]interface{}) interface{} {
	email, _ := variables["signup_email"].(string)

	param := authapi.SignUpParam{
		Email: email,
	}
	bodyBytes, _ := json.Marshal(param)

	return struct {
		ContentType string
		Body        io.Reader
	}{
		ContentType: "application/json",
		Body:        bytes.NewReader(bodyBytes),
	}
}

// resendSignUpConfirmationEmailParams は ResendSignUpConfirmationEmail メソッドのパラメータを返します
// メソッドシグネチャ: ResendSignUpConfirmationEmail(ctx context.Context, body ResendSignUpConfirmationEmailParam, reqEditors ...RequestEditorFn) (*http.Response, error)
func resendSignUpConfirmationEmailParams(variables map[string]interface{}) interface{} {
	email, _ := variables["signup_email"].(string)

	return authapi.ResendSignUpConfirmationEmailJSONRequestBody{
		Email: email,
	}
}

// resendSignUpConfirmationEmailWithBodyParams は ResendSignUpConfirmationEmailWithBody メソッドのパラメータを返します
// メソッドシグネチャ: ResendSignUpConfirmationEmailWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
func resendSignUpConfirmationEmailWithBodyParams(variables map[string]interface{}) interface{} {
	email, _ := variables["signup_email"].(string)

	param := authapi.ResendSignUpConfirmationEmailParam{
		Email: email,
	}
	bodyBytes, _ := json.Marshal(param)

	return struct {
		ContentType string
		Body        io.Reader
	}{
		ContentType: "application/json",
		Body:        bytes.NewReader(bodyBytes),
	}
}

// confirmSignUpWithAwsMarketplaceParams は ConfirmSignUpWithAwsMarketplace メソッドのパラメータを返します
// メソッドシグネチャ: ConfirmSignUpWithAwsMarketplace(ctx context.Context, body ConfirmSignUpWithAwsMarketplaceParam, reqEditors ...RequestEditorFn) (*http.Response, error)
func confirmSignUpWithAwsMarketplaceParams(variables map[string]interface{}) interface{} {
	accessToken, _ := variables["access_token"].(string)
	registrationToken, _ := variables["registration_token"].(string)
	tenantName, _ := variables["tenant_name"].(string)

	return authapi.ConfirmSignUpWithAwsMarketplaceJSONRequestBody{
		AccessToken:       accessToken,
		RegistrationToken: registrationToken,
		TenantName:        &tenantName,
	}
}

// confirmSignUpWithAwsMarketplaceWithBodyParams は ConfirmSignUpWithAwsMarketplaceWithBody メソッドのパラメータを返します
// メソッドシグネチャ: ConfirmSignUpWithAwsMarketplaceWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
func confirmSignUpWithAwsMarketplaceWithBodyParams(variables map[string]interface{}) interface{} {
	accessToken, _ := variables["access_token"].(string)
	registrationToken, _ := variables["registration_token"].(string)
	tenantName, _ := variables["tenant_name"].(string)

	param := authapi.ConfirmSignUpWithAwsMarketplaceParam{
		AccessToken:       accessToken,
		RegistrationToken: registrationToken,
		TenantName:        &tenantName,
	}
	bodyBytes, _ := json.Marshal(param)

	return struct {
		ContentType string
		Body        io.Reader
	}{
		ContentType: "application/json",
		Body:        bytes.NewReader(bodyBytes),
	}
}

// linkAwsMarketplaceParams は LinkAwsMarketplace メソッドのパラメータを返します
// メソッドシグネチャ: LinkAwsMarketplace(ctx context.Context, body LinkAwsMarketplaceParam, reqEditors ...RequestEditorFn) (*http.Response, error)
func linkAwsMarketplaceParams(variables map[string]interface{}) interface{} {
	tenantID, _ := variables["tenant_id"].(string)
	accessToken, _ := variables["access_token"].(string)
	registrationToken, _ := variables["registration_token"].(string)

	return authapi.LinkAwsMarketplaceJSONRequestBody{
		AccessToken:       accessToken,
		RegistrationToken: registrationToken,
		TenantId:          tenantID,
	}
}

// linkAwsMarketplaceWithBodyParams は LinkAwsMarketplaceWithBody メソッドのパラメータを返します
// メソッドシグネチャ: LinkAwsMarketplaceWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)
func linkAwsMarketplaceWithBodyParams(variables map[string]interface{}) interface{} {
	tenantID, _ := variables["tenant_id"].(string)
	accessToken, _ := variables["access_token"].(string)
	registrationToken, _ := variables["registration_token"].(string)

	param := authapi.LinkAwsMarketplaceParam{
		AccessToken:       accessToken,
		RegistrationToken: registrationToken,
		TenantId:          tenantID,
	}
	bodyBytes, _ := json.Marshal(param)

	return struct {
		ContentType string
		Body        io.Reader
	}{
		ContentType: "application/json",
		Body:        bytes.NewReader(bodyBytes),
	}
}

// extractResponseJSON は API レスポンスから JSON データを抽出します
// *http.Response と生成された型の両方に対応します
func extractResponseJSON(response any) (map[string]any, error) {
	if response == nil {
		return nil, fmt.Errorf("response is nil")
	}

	// すでに map[string]any の場合は直接返す
	if respMap, ok := response.(map[string]any); ok {
		return respMap, nil
	}

	// *http.Response の場合
	if httpResp, ok := response.(*http.Response); ok {
		if httpResp.Body == nil {
			return nil, fmt.Errorf("response body is nil")
		}

		// ボディを読み取る
		bodyBytes, err := io.ReadAll(httpResp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		// ボディを復元（他の処理で使えるように）
		httpResp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// JSON をパース
		var result map[string]any
		if err := json.Unmarshal(bodyBytes, &result); err != nil {
			return nil, fmt.Errorf("failed to parse JSON: %w", err)
		}

		return result, nil
	}

	// 生成された型の場合（例：*authapi.CreateTenantResponse）
	// JSON200, JSON201 などのフィールドを探す
	respValue := reflect.ValueOf(response)
	if respValue.Kind() == reflect.Ptr {
		respValue = respValue.Elem()
	}

	// 構造体でない場合はエラー
	if respValue.Kind() != reflect.Struct {
		return nil, fmt.Errorf("response is not a struct, got %v", respValue.Kind())
	}

	// JSON200, JSON201, JSON202 などのフィールドを探す
	for i := 0; i < respValue.NumField(); i++ {
		field := respValue.Field(i)
		fieldName := respValue.Type().Field(i).Name

		// JSON で始まるフィールドを探す
		if len(fieldName) > 4 && fieldName[:4] == "JSON" {
			if !field.IsNil() {
				// フィールドの値を map[string]any に変換
				jsonBytes, err := json.Marshal(field.Interface())
				if err != nil {
					continue
				}

				var result map[string]any
				if err := json.Unmarshal(jsonBytes, &result); err != nil {
					continue
				}

				return result, nil
			}
		}
	}

	// Body フィールドを探す
	bodyField := respValue.FieldByName("Body")
	if bodyField.IsValid() && bodyField.Kind() == reflect.Slice && bodyField.Type().Elem().Kind() == reflect.Uint8 {
		bodyBytes := bodyField.Bytes()
		if len(bodyBytes) > 0 {
			var result map[string]any
			if err := json.Unmarshal(bodyBytes, &result); err == nil {
				return result, nil
			}
		}
	}

	return nil, fmt.Errorf("unable to extract JSON from response type %T", response)
}

// extractIDFromResponse はレスポンスから ID を抽出します
func extractIDFromResponse(response any, idField string) (string, error) {
	jsonData, err := extractResponseJSON(response)
	if err != nil {
		return "", err
	}

	if id, ok := jsonData[idField].(string); ok && id != "" {
		return id, nil
	}

	return "", fmt.Errorf("field '%s' not found or empty in response", idField)
}

// extractFieldFromResponse はレスポンスから任意のフィールドを抽出します（extractIDFromResponseのエイリアス）
func extractFieldFromResponse(response any, fieldName string) (string, error) {
	return extractIDFromResponse(response, fieldName)
}
