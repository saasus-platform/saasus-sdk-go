type ClientWithResponsesInterface interface {
	// GetAuthInfo request
	GetAuthInfoWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetAuthInfoResponse, error)

	// UpdateAuthInfo request with any body
	UpdateAuthInfoWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateAuthInfoResponse, error)

	UpdateAuthInfoWithResponse(ctx context.Context, body UpdateAuthInfoJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateAuthInfoResponse, error)

	// LinkAwsMarketplace request with any body
	LinkAwsMarketplaceWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*LinkAwsMarketplaceResponse, error)

	LinkAwsMarketplaceWithResponse(ctx context.Context, body LinkAwsMarketplaceJSONRequestBody, reqEditors ...RequestEditorFn) (*LinkAwsMarketplaceResponse, error)

	// SignUpWithAwsMarketplace request with any body
	SignUpWithAwsMarketplaceWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*SignUpWithAwsMarketplaceResponse, error)

	SignUpWithAwsMarketplaceWithResponse(ctx context.Context, body SignUpWithAwsMarketplaceJSONRequestBody, reqEditors ...RequestEditorFn) (*SignUpWithAwsMarketplaceResponse, error)

	// ConfirmSignUpWithAwsMarketplace request with any body
	ConfirmSignUpWithAwsMarketplaceWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*ConfirmSignUpWithAwsMarketplaceResponse, error)

	ConfirmSignUpWithAwsMarketplaceWithResponse(ctx context.Context, body ConfirmSignUpWithAwsMarketplaceJSONRequestBody, reqEditors ...RequestEditorFn) (*ConfirmSignUpWithAwsMarketplaceResponse, error)

	// GetBasicInfo request
	GetBasicInfoWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetBasicInfoResponse, error)

	// UpdateBasicInfo request with any body
	UpdateBasicInfoWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateBasicInfoResponse, error)

	UpdateBasicInfoWithResponse(ctx context.Context, body UpdateBasicInfoJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateBasicInfoResponse, error)

	// GetAuthCredentials request
	GetAuthCredentialsWithResponse(ctx context.Context, params *GetAuthCredentialsParams, reqEditors ...RequestEditorFn) (*GetAuthCredentialsResponse, error)

	// CreateAuthCredentials request with any body
	CreateAuthCredentialsWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateAuthCredentialsResponse, error)

	CreateAuthCredentialsWithResponse(ctx context.Context, body CreateAuthCredentialsJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateAuthCredentialsResponse, error)

	// GetCustomizePageSettings request
	GetCustomizePageSettingsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetCustomizePageSettingsResponse, error)

	// UpdateCustomizePageSettings request with any body
	UpdateCustomizePageSettingsWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateCustomizePageSettingsResponse, error)

	UpdateCustomizePageSettingsWithResponse(ctx context.Context, body UpdateCustomizePageSettingsJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateCustomizePageSettingsResponse, error)

	// GetCustomizePages request
	GetCustomizePagesWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetCustomizePagesResponse, error)

	// UpdateCustomizePages request with any body
	UpdateCustomizePagesWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateCustomizePagesResponse, error)

	UpdateCustomizePagesWithResponse(ctx context.Context, body UpdateCustomizePagesJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateCustomizePagesResponse, error)

	// GetEnvs request
	GetEnvsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetEnvsResponse, error)

	// CreateEnv request with any body
	CreateEnvWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateEnvResponse, error)

	CreateEnvWithResponse(ctx context.Context, body CreateEnvJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateEnvResponse, error)

	// DeleteEnv request
	DeleteEnvWithResponse(ctx context.Context, envId EnvId, reqEditors ...RequestEditorFn) (*DeleteEnvResponse, error)

	// GetEnv request
	GetEnvWithResponse(ctx context.Context, envId EnvId, reqEditors ...RequestEditorFn) (*GetEnvResponse, error)

	// UpdateEnv request with any body
	UpdateEnvWithBodyWithResponse(ctx context.Context, envId EnvId, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateEnvResponse, error)

	UpdateEnvWithResponse(ctx context.Context, envId EnvId, body UpdateEnvJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateEnvResponse, error)

	// ReturnInternalServerError request
	ReturnInternalServerErrorWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*ReturnInternalServerErrorResponse, error)

	// ConfirmExternalUserLink request with any body
	ConfirmExternalUserLinkWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*ConfirmExternalUserLinkResponse, error)

	ConfirmExternalUserLinkWithResponse(ctx context.Context, body ConfirmExternalUserLinkJSONRequestBody, reqEditors ...RequestEditorFn) (*ConfirmExternalUserLinkResponse, error)

	// RequestExternalUserLink request with any body
	RequestExternalUserLinkWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*RequestExternalUserLinkResponse, error)

	RequestExternalUserLinkWithResponse(ctx context.Context, body RequestExternalUserLinkJSONRequestBody, reqEditors ...RequestEditorFn) (*RequestExternalUserLinkResponse, error)

	// GetIdentityProviders request
	GetIdentityProvidersWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetIdentityProvidersResponse, error)

	// UpdateIdentityProvider request with any body
	UpdateIdentityProviderWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateIdentityProviderResponse, error)

	UpdateIdentityProviderWithResponse(ctx context.Context, body UpdateIdentityProviderJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateIdentityProviderResponse, error)

	// ValidateInvitation request with any body
	ValidateInvitationWithBodyWithResponse(ctx context.Context, invitationId InvitationId, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*ValidateInvitationResponse, error)

	ValidateInvitationWithResponse(ctx context.Context, invitationId InvitationId, body ValidateInvitationJSONRequestBody, reqEditors ...RequestEditorFn) (*ValidateInvitationResponse, error)

	// GetInvitationValidity request
	GetInvitationValidityWithResponse(ctx context.Context, invitationId InvitationId, reqEditors ...RequestEditorFn) (*GetInvitationValidityResponse, error)

	// FindNotificationMessages request
	FindNotificationMessagesWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*FindNotificationMessagesResponse, error)

	// UpdateNotificationMessages request with any body
	UpdateNotificationMessagesWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateNotificationMessagesResponse, error)

	UpdateNotificationMessagesWithResponse(ctx context.Context, body UpdateNotificationMessagesJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateNotificationMessagesResponse, error)

	// ResetPlan request
	ResetPlanWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*ResetPlanResponse, error)

	// GetRoles request
	GetRolesWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetRolesResponse, error)

	// CreateRole request with any body
	CreateRoleWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateRoleResponse, error)

	CreateRoleWithResponse(ctx context.Context, body CreateRoleJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateRoleResponse, error)

	// DeleteRole request
	DeleteRoleWithResponse(ctx context.Context, roleName RoleName, reqEditors ...RequestEditorFn) (*DeleteRoleResponse, error)

	// CreateSaasUserAttribute request with any body
	CreateSaasUserAttributeWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateSaasUserAttributeResponse, error)

	CreateSaasUserAttributeWithResponse(ctx context.Context, body CreateSaasUserAttributeJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateSaasUserAttributeResponse, error)

	// GetSignInSettings request
	GetSignInSettingsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetSignInSettingsResponse, error)

	// UpdateSignInSettings request with any body
	UpdateSignInSettingsWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateSignInSettingsResponse, error)

	UpdateSignInSettingsWithResponse(ctx context.Context, body UpdateSignInSettingsJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateSignInSettingsResponse, error)

	// SignUp request with any body
	SignUpWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*SignUpResponse, error)

	SignUpWithResponse(ctx context.Context, body SignUpJSONRequestBody, reqEditors ...RequestEditorFn) (*SignUpResponse, error)

	// ResendSignUpConfirmationEmail request with any body
	ResendSignUpConfirmationEmailWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*ResendSignUpConfirmationEmailResponse, error)

	ResendSignUpConfirmationEmailWithResponse(ctx context.Context, body ResendSignUpConfirmationEmailJSONRequestBody, reqEditors ...RequestEditorFn) (*ResendSignUpConfirmationEmailResponse, error)

	// GetCloudFormationLaunchStackLinkForSingleTenant request
	GetCloudFormationLaunchStackLinkForSingleTenantWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetCloudFormationLaunchStackLinkForSingleTenantResponse, error)

	// GetSingleTenantSettings request
	GetSingleTenantSettingsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetSingleTenantSettingsResponse, error)

	// UpdateSingleTenantSettings request with any body
	UpdateSingleTenantSettingsWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateSingleTenantSettingsResponse, error)

	UpdateSingleTenantSettingsWithResponse(ctx context.Context, body UpdateSingleTenantSettingsJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateSingleTenantSettingsResponse, error)

	// DeleteStripeTenantAndPricing request
	DeleteStripeTenantAndPricingWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*DeleteStripeTenantAndPricingResponse, error)

	// CreateTenantAndPricing request
	CreateTenantAndPricingWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*CreateTenantAndPricingResponse, error)

	// GetTenantAttributes request
	GetTenantAttributesWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetTenantAttributesResponse, error)

	// CreateTenantAttribute request with any body
	CreateTenantAttributeWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateTenantAttributeResponse, error)

	CreateTenantAttributeWithResponse(ctx context.Context, body CreateTenantAttributeJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateTenantAttributeResponse, error)

	// DeleteTenantAttribute request
	DeleteTenantAttributeWithResponse(ctx context.Context, attributeName string, reqEditors ...RequestEditorFn) (*DeleteTenantAttributeResponse, error)

	// GetTenants request
	GetTenantsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetTenantsResponse, error)

	// CreateTenant request with any body
	CreateTenantWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateTenantResponse, error)

	CreateTenantWithResponse(ctx context.Context, body CreateTenantJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateTenantResponse, error)

	// GetAllTenantUsers request
	GetAllTenantUsersWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetAllTenantUsersResponse, error)

	// GetAllTenantUser request
	GetAllTenantUserWithResponse(ctx context.Context, userId UserId, reqEditors ...RequestEditorFn) (*GetAllTenantUserResponse, error)

	// DeleteTenant request
	DeleteTenantWithResponse(ctx context.Context, tenantId TenantId, reqEditors ...RequestEditorFn) (*DeleteTenantResponse, error)

	// GetTenant request
	GetTenantWithResponse(ctx context.Context, tenantId TenantId, reqEditors ...RequestEditorFn) (*GetTenantResponse, error)

	// UpdateTenant request with any body
	UpdateTenantWithBodyWithResponse(ctx context.Context, tenantId TenantId, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateTenantResponse, error)

	UpdateTenantWithResponse(ctx context.Context, tenantId TenantId, body UpdateTenantJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateTenantResponse, error)

	// UpdateTenantBillingInfo request with any body
	UpdateTenantBillingInfoWithBodyWithResponse(ctx context.Context, tenantId TenantId, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateTenantBillingInfoResponse, error)

	UpdateTenantBillingInfoWithResponse(ctx context.Context, tenantId TenantId, body UpdateTenantBillingInfoJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateTenantBillingInfoResponse, error)

	// GetTenantIdentityProviders request
	GetTenantIdentityProvidersWithResponse(ctx context.Context, tenantId TenantId, reqEditors ...RequestEditorFn) (*GetTenantIdentityProvidersResponse, error)

	// UpdateTenantIdentityProvider request with any body
	UpdateTenantIdentityProviderWithBodyWithResponse(ctx context.Context, tenantId TenantId, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateTenantIdentityProviderResponse, error)

	UpdateTenantIdentityProviderWithResponse(ctx context.Context, tenantId TenantId, body UpdateTenantIdentityProviderJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateTenantIdentityProviderResponse, error)

	// GetTenantInvitations request
	GetTenantInvitationsWithResponse(ctx context.Context, tenantId TenantId, reqEditors ...RequestEditorFn) (*GetTenantInvitationsResponse, error)

	// CreateTenantInvitation request with any body
	CreateTenantInvitationWithBodyWithResponse(ctx context.Context, tenantId TenantId, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateTenantInvitationResponse, error)

	CreateTenantInvitationWithResponse(ctx context.Context, tenantId TenantId, body CreateTenantInvitationJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateTenantInvitationResponse, error)

	// DeleteTenantInvitation request
	DeleteTenantInvitationWithResponse(ctx context.Context, tenantId TenantId, invitationId InvitationId, reqEditors ...RequestEditorFn) (*DeleteTenantInvitationResponse, error)

	// GetTenantInvitation request
	GetTenantInvitationWithResponse(ctx context.Context, tenantId TenantId, invitationId InvitationId, reqEditors ...RequestEditorFn) (*GetTenantInvitationResponse, error)

	// UpdateTenantPlan request with any body
	UpdateTenantPlanWithBodyWithResponse(ctx context.Context, tenantId TenantId, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateTenantPlanResponse, error)

	UpdateTenantPlanWithResponse(ctx context.Context, tenantId TenantId, body UpdateTenantPlanJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateTenantPlanResponse, error)

	// GetStripeCustomer request
	GetStripeCustomerWithResponse(ctx context.Context, tenantId TenantId, reqEditors ...RequestEditorFn) (*GetStripeCustomerResponse, error)

	// GetTenantUsers request
	GetTenantUsersWithResponse(ctx context.Context, tenantId TenantId, reqEditors ...RequestEditorFn) (*GetTenantUsersResponse, error)

	// CreateTenantUser request with any body
	CreateTenantUserWithBodyWithResponse(ctx context.Context, tenantId TenantId, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateTenantUserResponse, error)

	CreateTenantUserWithResponse(ctx context.Context, tenantId TenantId, body CreateTenantUserJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateTenantUserResponse, error)

	// DeleteTenantUser request
	DeleteTenantUserWithResponse(ctx context.Context, tenantId TenantId, userId UserId, reqEditors ...RequestEditorFn) (*DeleteTenantUserResponse, error)

	// GetTenantUser request
	GetTenantUserWithResponse(ctx context.Context, tenantId TenantId, userId UserId, reqEditors ...RequestEditorFn) (*GetTenantUserResponse, error)

	// UpdateTenantUser request with any body
	UpdateTenantUserWithBodyWithResponse(ctx context.Context, tenantId TenantId, userId UserId, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateTenantUserResponse, error)

	UpdateTenantUserWithResponse(ctx context.Context, tenantId TenantId, userId UserId, body UpdateTenantUserJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateTenantUserResponse, error)

	// CreateTenantUserRoles request with any body
	CreateTenantUserRolesWithBodyWithResponse(ctx context.Context, tenantId TenantId, userId UserId, envId EnvId, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateTenantUserRolesResponse, error)

	CreateTenantUserRolesWithResponse(ctx context.Context, tenantId TenantId, userId UserId, envId EnvId, body CreateTenantUserRolesJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateTenantUserRolesResponse, error)

	// DeleteTenantUserRole request
	DeleteTenantUserRoleWithResponse(ctx context.Context, tenantId TenantId, userId UserId, envId EnvId, roleName RoleName, reqEditors ...RequestEditorFn) (*DeleteTenantUserRoleResponse, error)

	// GetUserAttributes request
	GetUserAttributesWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetUserAttributesResponse, error)

	// CreateUserAttribute request with any body
	CreateUserAttributeWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateUserAttributeResponse, error)

	CreateUserAttributeWithResponse(ctx context.Context, body CreateUserAttributeJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateUserAttributeResponse, error)

	// DeleteUserAttribute request
	DeleteUserAttributeWithResponse(ctx context.Context, attributeName string, reqEditors ...RequestEditorFn) (*DeleteUserAttributeResponse, error)

	// GetUserInfo request
	GetUserInfoWithResponse(ctx context.Context, params *GetUserInfoParams, reqEditors ...RequestEditorFn) (*GetUserInfoResponse, error)

	// GetSaasUsers request
	GetSaasUsersWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetSaasUsersResponse, error)

	// CreateSaasUser request with any body
	CreateSaasUserWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateSaasUserResponse, error)

	CreateSaasUserWithResponse(ctx context.Context, body CreateSaasUserJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateSaasUserResponse, error)

	// DeleteSaasUser request
	DeleteSaasUserWithResponse(ctx context.Context, userId UserId, reqEditors ...RequestEditorFn) (*DeleteSaasUserResponse, error)

	// GetSaasUser request
	GetSaasUserWithResponse(ctx context.Context, userId UserId, reqEditors ...RequestEditorFn) (*GetSaasUserResponse, error)

	// UpdateSaasUserAttributes request with any body
	UpdateSaasUserAttributesWithBodyWithResponse(ctx context.Context, userId UserId, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateSaasUserAttributesResponse, error)

	UpdateSaasUserAttributesWithResponse(ctx context.Context, userId UserId, body UpdateSaasUserAttributesJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateSaasUserAttributesResponse, error)

	// UpdateSaasUserEmail request with any body
	UpdateSaasUserEmailWithBodyWithResponse(ctx context.Context, userId UserId, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateSaasUserEmailResponse, error)

	UpdateSaasUserEmailWithResponse(ctx context.Context, userId UserId, body UpdateSaasUserEmailJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateSaasUserEmailResponse, error)

	// ConfirmEmailUpdate request with any body
	ConfirmEmailUpdateWithBodyWithResponse(ctx context.Context, userId UserId, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*ConfirmEmailUpdateResponse, error)

	ConfirmEmailUpdateWithResponse(ctx context.Context, userId UserId, body ConfirmEmailUpdateJSONRequestBody, reqEditors ...RequestEditorFn) (*ConfirmEmailUpdateResponse, error)

	// RequestEmailUpdate request with any body
	RequestEmailUpdateWithBodyWithResponse(ctx context.Context, userId UserId, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*RequestEmailUpdateResponse, error)

	RequestEmailUpdateWithResponse(ctx context.Context, userId UserId, body RequestEmailUpdateJSONRequestBody, reqEditors ...RequestEditorFn) (*RequestEmailUpdateResponse, error)

	// GetUserMfaPreference request
	GetUserMfaPreferenceWithResponse(ctx context.Context, userId UserId, reqEditors ...RequestEditorFn) (*GetUserMfaPreferenceResponse, error)

	// UpdateUserMfaPreference request with any body
	UpdateUserMfaPreferenceWithBodyWithResponse(ctx context.Context, userId UserId, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateUserMfaPreferenceResponse, error)

	UpdateUserMfaPreferenceWithResponse(ctx context.Context, userId UserId, body UpdateUserMfaPreferenceJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateUserMfaPreferenceResponse, error)

	// UpdateSoftwareToken request with any body
	UpdateSoftwareTokenWithBodyWithResponse(ctx context.Context, userId UserId, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateSoftwareTokenResponse, error)

	UpdateSoftwareTokenWithResponse(ctx context.Context, userId UserId, body UpdateSoftwareTokenJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateSoftwareTokenResponse, error)

	// CreateSecretCode request with any body
	CreateSecretCodeWithBodyWithResponse(ctx context.Context, userId UserId, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateSecretCodeResponse, error)

	CreateSecretCodeWithResponse(ctx context.Context, userId UserId, body CreateSecretCodeJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateSecretCodeResponse, error)

	// UpdateSaasUserPassword request with any body
	UpdateSaasUserPasswordWithBodyWithResponse(ctx context.Context, userId UserId, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateSaasUserPasswordResponse, error)

	UpdateSaasUserPasswordWithResponse(ctx context.Context, userId UserId, body UpdateSaasUserPasswordJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateSaasUserPasswordResponse, error)

	// UnlinkProvider request
	UnlinkProviderWithResponse(ctx context.Context, userId UserId, providerName string, reqEditors ...RequestEditorFn) (*UnlinkProviderResponse, error)
}

