type ClientInterface interface {
	// GetAuthInfo request
	GetAuthInfo(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// UpdateAuthInfo request with any body
	UpdateAuthInfoWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	UpdateAuthInfo(ctx context.Context, body UpdateAuthInfoJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// LinkAwsMarketplace request with any body
	LinkAwsMarketplaceWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	LinkAwsMarketplace(ctx context.Context, body LinkAwsMarketplaceJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// SignUpWithAwsMarketplace request with any body
	SignUpWithAwsMarketplaceWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	SignUpWithAwsMarketplace(ctx context.Context, body SignUpWithAwsMarketplaceJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// ConfirmSignUpWithAwsMarketplace request with any body
	ConfirmSignUpWithAwsMarketplaceWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	ConfirmSignUpWithAwsMarketplace(ctx context.Context, body ConfirmSignUpWithAwsMarketplaceJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetBasicInfo request
	GetBasicInfo(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// UpdateBasicInfo request with any body
	UpdateBasicInfoWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	UpdateBasicInfo(ctx context.Context, body UpdateBasicInfoJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetAuthCredentials request
	GetAuthCredentials(ctx context.Context, params *GetAuthCredentialsParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// CreateAuthCredentials request with any body
	CreateAuthCredentialsWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	CreateAuthCredentials(ctx context.Context, body CreateAuthCredentialsJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetCustomizePageSettings request
	GetCustomizePageSettings(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// UpdateCustomizePageSettings request with any body
	UpdateCustomizePageSettingsWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	UpdateCustomizePageSettings(ctx context.Context, body UpdateCustomizePageSettingsJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetCustomizePages request
	GetCustomizePages(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// UpdateCustomizePages request with any body
	UpdateCustomizePagesWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	UpdateCustomizePages(ctx context.Context, body UpdateCustomizePagesJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetEnvs request
	GetEnvs(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// CreateEnv request with any body
	CreateEnvWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	CreateEnv(ctx context.Context, body CreateEnvJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// DeleteEnv request
	DeleteEnv(ctx context.Context, envId EnvId, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetEnv request
	GetEnv(ctx context.Context, envId EnvId, reqEditors ...RequestEditorFn) (*http.Response, error)

	// UpdateEnv request with any body
	UpdateEnvWithBody(ctx context.Context, envId EnvId, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	UpdateEnv(ctx context.Context, envId EnvId, body UpdateEnvJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// ReturnInternalServerError request
	ReturnInternalServerError(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// ConfirmExternalUserLink request with any body
	ConfirmExternalUserLinkWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	ConfirmExternalUserLink(ctx context.Context, body ConfirmExternalUserLinkJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// RequestExternalUserLink request with any body
	RequestExternalUserLinkWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	RequestExternalUserLink(ctx context.Context, body RequestExternalUserLinkJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetIdentityProviders request
	GetIdentityProviders(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// UpdateIdentityProvider request with any body
	UpdateIdentityProviderWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	UpdateIdentityProvider(ctx context.Context, body UpdateIdentityProviderJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// ValidateInvitation request with any body
	ValidateInvitationWithBody(ctx context.Context, invitationId InvitationId, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	ValidateInvitation(ctx context.Context, invitationId InvitationId, body ValidateInvitationJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetInvitationValidity request
	GetInvitationValidity(ctx context.Context, invitationId InvitationId, reqEditors ...RequestEditorFn) (*http.Response, error)

	// FindNotificationMessages request
	FindNotificationMessages(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// UpdateNotificationMessages request with any body
	UpdateNotificationMessagesWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	UpdateNotificationMessages(ctx context.Context, body UpdateNotificationMessagesJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// ResetPlan request
	ResetPlan(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetRoles request
	GetRoles(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// CreateRole request with any body
	CreateRoleWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	CreateRole(ctx context.Context, body CreateRoleJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// DeleteRole request
	DeleteRole(ctx context.Context, roleName RoleName, reqEditors ...RequestEditorFn) (*http.Response, error)

	// CreateSaasUserAttribute request with any body
	CreateSaasUserAttributeWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	CreateSaasUserAttribute(ctx context.Context, body CreateSaasUserAttributeJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetSignInSettings request
	GetSignInSettings(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// UpdateSignInSettings request with any body
	UpdateSignInSettingsWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	UpdateSignInSettings(ctx context.Context, body UpdateSignInSettingsJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// SignUp request with any body
	SignUpWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	SignUp(ctx context.Context, body SignUpJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// ResendSignUpConfirmationEmail request with any body
	ResendSignUpConfirmationEmailWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	ResendSignUpConfirmationEmail(ctx context.Context, body ResendSignUpConfirmationEmailJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetCloudFormationLaunchStackLinkForSingleTenant request
	GetCloudFormationLaunchStackLinkForSingleTenant(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetSingleTenantSettings request
	GetSingleTenantSettings(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// UpdateSingleTenantSettings request with any body
	UpdateSingleTenantSettingsWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	UpdateSingleTenantSettings(ctx context.Context, body UpdateSingleTenantSettingsJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// DeleteStripeTenantAndPricing request
	DeleteStripeTenantAndPricing(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// CreateTenantAndPricing request
	CreateTenantAndPricing(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetTenantAttributes request
	GetTenantAttributes(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// CreateTenantAttribute request with any body
	CreateTenantAttributeWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	CreateTenantAttribute(ctx context.Context, body CreateTenantAttributeJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// DeleteTenantAttribute request
	DeleteTenantAttribute(ctx context.Context, attributeName string, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetTenants request
	GetTenants(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// CreateTenant request with any body
	CreateTenantWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	CreateTenant(ctx context.Context, body CreateTenantJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetAllTenantUsers request
	GetAllTenantUsers(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetAllTenantUser request
	GetAllTenantUser(ctx context.Context, userId UserId, reqEditors ...RequestEditorFn) (*http.Response, error)

	// DeleteTenant request
	DeleteTenant(ctx context.Context, tenantId TenantId, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetTenant request
	GetTenant(ctx context.Context, tenantId TenantId, reqEditors ...RequestEditorFn) (*http.Response, error)

	// UpdateTenant request with any body
	UpdateTenantWithBody(ctx context.Context, tenantId TenantId, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	UpdateTenant(ctx context.Context, tenantId TenantId, body UpdateTenantJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// UpdateTenantBillingInfo request with any body
	UpdateTenantBillingInfoWithBody(ctx context.Context, tenantId TenantId, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	UpdateTenantBillingInfo(ctx context.Context, tenantId TenantId, body UpdateTenantBillingInfoJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetTenantIdentityProviders request
	GetTenantIdentityProviders(ctx context.Context, tenantId TenantId, reqEditors ...RequestEditorFn) (*http.Response, error)

	// UpdateTenantIdentityProvider request with any body
	UpdateTenantIdentityProviderWithBody(ctx context.Context, tenantId TenantId, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	UpdateTenantIdentityProvider(ctx context.Context, tenantId TenantId, body UpdateTenantIdentityProviderJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetTenantInvitations request
	GetTenantInvitations(ctx context.Context, tenantId TenantId, reqEditors ...RequestEditorFn) (*http.Response, error)

	// CreateTenantInvitation request with any body
	CreateTenantInvitationWithBody(ctx context.Context, tenantId TenantId, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	CreateTenantInvitation(ctx context.Context, tenantId TenantId, body CreateTenantInvitationJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// DeleteTenantInvitation request
	DeleteTenantInvitation(ctx context.Context, tenantId TenantId, invitationId InvitationId, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetTenantInvitation request
	GetTenantInvitation(ctx context.Context, tenantId TenantId, invitationId InvitationId, reqEditors ...RequestEditorFn) (*http.Response, error)

	// UpdateTenantPlan request with any body
	UpdateTenantPlanWithBody(ctx context.Context, tenantId TenantId, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	UpdateTenantPlan(ctx context.Context, tenantId TenantId, body UpdateTenantPlanJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetStripeCustomer request
	GetStripeCustomer(ctx context.Context, tenantId TenantId, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetTenantUsers request
	GetTenantUsers(ctx context.Context, tenantId TenantId, reqEditors ...RequestEditorFn) (*http.Response, error)

	// CreateTenantUser request with any body
	CreateTenantUserWithBody(ctx context.Context, tenantId TenantId, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	CreateTenantUser(ctx context.Context, tenantId TenantId, body CreateTenantUserJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// DeleteTenantUser request
	DeleteTenantUser(ctx context.Context, tenantId TenantId, userId UserId, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetTenantUser request
	GetTenantUser(ctx context.Context, tenantId TenantId, userId UserId, reqEditors ...RequestEditorFn) (*http.Response, error)

	// UpdateTenantUser request with any body
	UpdateTenantUserWithBody(ctx context.Context, tenantId TenantId, userId UserId, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	UpdateTenantUser(ctx context.Context, tenantId TenantId, userId UserId, body UpdateTenantUserJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// CreateTenantUserRoles request with any body
	CreateTenantUserRolesWithBody(ctx context.Context, tenantId TenantId, userId UserId, envId EnvId, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	CreateTenantUserRoles(ctx context.Context, tenantId TenantId, userId UserId, envId EnvId, body CreateTenantUserRolesJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// DeleteTenantUserRole request
	DeleteTenantUserRole(ctx context.Context, tenantId TenantId, userId UserId, envId EnvId, roleName RoleName, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetUserAttributes request
	GetUserAttributes(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// CreateUserAttribute request with any body
	CreateUserAttributeWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	CreateUserAttribute(ctx context.Context, body CreateUserAttributeJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// DeleteUserAttribute request
	DeleteUserAttribute(ctx context.Context, attributeName string, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetUserInfo request
	GetUserInfo(ctx context.Context, params *GetUserInfoParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetSaasUsers request
	GetSaasUsers(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// CreateSaasUser request with any body
	CreateSaasUserWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	CreateSaasUser(ctx context.Context, body CreateSaasUserJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// DeleteSaasUser request
	DeleteSaasUser(ctx context.Context, userId UserId, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetSaasUser request
	GetSaasUser(ctx context.Context, userId UserId, reqEditors ...RequestEditorFn) (*http.Response, error)

	// UpdateSaasUserAttributes request with any body
	UpdateSaasUserAttributesWithBody(ctx context.Context, userId UserId, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	UpdateSaasUserAttributes(ctx context.Context, userId UserId, body UpdateSaasUserAttributesJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// UpdateSaasUserEmail request with any body
	UpdateSaasUserEmailWithBody(ctx context.Context, userId UserId, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	UpdateSaasUserEmail(ctx context.Context, userId UserId, body UpdateSaasUserEmailJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// ConfirmEmailUpdate request with any body
	ConfirmEmailUpdateWithBody(ctx context.Context, userId UserId, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	ConfirmEmailUpdate(ctx context.Context, userId UserId, body ConfirmEmailUpdateJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// RequestEmailUpdate request with any body
	RequestEmailUpdateWithBody(ctx context.Context, userId UserId, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	RequestEmailUpdate(ctx context.Context, userId UserId, body RequestEmailUpdateJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetUserMfaPreference request
	GetUserMfaPreference(ctx context.Context, userId UserId, reqEditors ...RequestEditorFn) (*http.Response, error)

	// UpdateUserMfaPreference request with any body
	UpdateUserMfaPreferenceWithBody(ctx context.Context, userId UserId, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	UpdateUserMfaPreference(ctx context.Context, userId UserId, body UpdateUserMfaPreferenceJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// UpdateSoftwareToken request with any body
	UpdateSoftwareTokenWithBody(ctx context.Context, userId UserId, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	UpdateSoftwareToken(ctx context.Context, userId UserId, body UpdateSoftwareTokenJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// CreateSecretCode request with any body
	CreateSecretCodeWithBody(ctx context.Context, userId UserId, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	CreateSecretCode(ctx context.Context, userId UserId, body CreateSecretCodeJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// UpdateSaasUserPassword request with any body
	UpdateSaasUserPasswordWithBody(ctx context.Context, userId UserId, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	UpdateSaasUserPassword(ctx context.Context, userId UserId, body UpdateSaasUserPasswordJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// UnlinkProvider request
	UnlinkProvider(ctx context.Context, userId UserId, providerName string, reqEditors ...RequestEditorFn) (*http.Response, error)
}

