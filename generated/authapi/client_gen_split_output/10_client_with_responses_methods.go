func (c *ClientWithResponses) GetAuthInfoWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetAuthInfoResponse, error) {
	rsp, err := c.GetAuthInfo(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetAuthInfoResponse(rsp)
}

// UpdateAuthInfoWithBodyWithResponse request with arbitrary body returning *UpdateAuthInfoResponse
func (c *ClientWithResponses) UpdateAuthInfoWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateAuthInfoResponse, error) {
	rsp, err := c.UpdateAuthInfoWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateAuthInfoResponse(rsp)
}

func (c *ClientWithResponses) UpdateAuthInfoWithResponse(ctx context.Context, body UpdateAuthInfoJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateAuthInfoResponse, error) {
	rsp, err := c.UpdateAuthInfo(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateAuthInfoResponse(rsp)
}

// LinkAwsMarketplaceWithBodyWithResponse request with arbitrary body returning *LinkAwsMarketplaceResponse
func (c *ClientWithResponses) LinkAwsMarketplaceWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*LinkAwsMarketplaceResponse, error) {
	rsp, err := c.LinkAwsMarketplaceWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseLinkAwsMarketplaceResponse(rsp)
}

func (c *ClientWithResponses) LinkAwsMarketplaceWithResponse(ctx context.Context, body LinkAwsMarketplaceJSONRequestBody, reqEditors ...RequestEditorFn) (*LinkAwsMarketplaceResponse, error) {
	rsp, err := c.LinkAwsMarketplace(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseLinkAwsMarketplaceResponse(rsp)
}

// SignUpWithAwsMarketplaceWithBodyWithResponse request with arbitrary body returning *SignUpWithAwsMarketplaceResponse
func (c *ClientWithResponses) SignUpWithAwsMarketplaceWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*SignUpWithAwsMarketplaceResponse, error) {
	rsp, err := c.SignUpWithAwsMarketplaceWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseSignUpWithAwsMarketplaceResponse(rsp)
}

func (c *ClientWithResponses) SignUpWithAwsMarketplaceWithResponse(ctx context.Context, body SignUpWithAwsMarketplaceJSONRequestBody, reqEditors ...RequestEditorFn) (*SignUpWithAwsMarketplaceResponse, error) {
	rsp, err := c.SignUpWithAwsMarketplace(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseSignUpWithAwsMarketplaceResponse(rsp)
}

// ConfirmSignUpWithAwsMarketplaceWithBodyWithResponse request with arbitrary body returning *ConfirmSignUpWithAwsMarketplaceResponse
func (c *ClientWithResponses) ConfirmSignUpWithAwsMarketplaceWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*ConfirmSignUpWithAwsMarketplaceResponse, error) {
	rsp, err := c.ConfirmSignUpWithAwsMarketplaceWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseConfirmSignUpWithAwsMarketplaceResponse(rsp)
}

func (c *ClientWithResponses) ConfirmSignUpWithAwsMarketplaceWithResponse(ctx context.Context, body ConfirmSignUpWithAwsMarketplaceJSONRequestBody, reqEditors ...RequestEditorFn) (*ConfirmSignUpWithAwsMarketplaceResponse, error) {
	rsp, err := c.ConfirmSignUpWithAwsMarketplace(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseConfirmSignUpWithAwsMarketplaceResponse(rsp)
}

// GetBasicInfoWithResponse request returning *GetBasicInfoResponse
func (c *ClientWithResponses) GetBasicInfoWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetBasicInfoResponse, error) {
	rsp, err := c.GetBasicInfo(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetBasicInfoResponse(rsp)
}

// UpdateBasicInfoWithBodyWithResponse request with arbitrary body returning *UpdateBasicInfoResponse
func (c *ClientWithResponses) UpdateBasicInfoWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateBasicInfoResponse, error) {
	rsp, err := c.UpdateBasicInfoWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateBasicInfoResponse(rsp)
}

func (c *ClientWithResponses) UpdateBasicInfoWithResponse(ctx context.Context, body UpdateBasicInfoJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateBasicInfoResponse, error) {
	rsp, err := c.UpdateBasicInfo(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateBasicInfoResponse(rsp)
}

// GetAuthCredentialsWithResponse request returning *GetAuthCredentialsResponse
func (c *ClientWithResponses) GetAuthCredentialsWithResponse(ctx context.Context, params *GetAuthCredentialsParams, reqEditors ...RequestEditorFn) (*GetAuthCredentialsResponse, error) {
	rsp, err := c.GetAuthCredentials(ctx, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetAuthCredentialsResponse(rsp)
}

// CreateAuthCredentialsWithBodyWithResponse request with arbitrary body returning *CreateAuthCredentialsResponse
func (c *ClientWithResponses) CreateAuthCredentialsWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateAuthCredentialsResponse, error) {
	rsp, err := c.CreateAuthCredentialsWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateAuthCredentialsResponse(rsp)
}

func (c *ClientWithResponses) CreateAuthCredentialsWithResponse(ctx context.Context, body CreateAuthCredentialsJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateAuthCredentialsResponse, error) {
	rsp, err := c.CreateAuthCredentials(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateAuthCredentialsResponse(rsp)
}

// GetCustomizePageSettingsWithResponse request returning *GetCustomizePageSettingsResponse
func (c *ClientWithResponses) GetCustomizePageSettingsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetCustomizePageSettingsResponse, error) {
	rsp, err := c.GetCustomizePageSettings(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetCustomizePageSettingsResponse(rsp)
}

// UpdateCustomizePageSettingsWithBodyWithResponse request with arbitrary body returning *UpdateCustomizePageSettingsResponse
func (c *ClientWithResponses) UpdateCustomizePageSettingsWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateCustomizePageSettingsResponse, error) {
	rsp, err := c.UpdateCustomizePageSettingsWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateCustomizePageSettingsResponse(rsp)
}

func (c *ClientWithResponses) UpdateCustomizePageSettingsWithResponse(ctx context.Context, body UpdateCustomizePageSettingsJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateCustomizePageSettingsResponse, error) {
	rsp, err := c.UpdateCustomizePageSettings(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateCustomizePageSettingsResponse(rsp)
}

// GetCustomizePagesWithResponse request returning *GetCustomizePagesResponse
func (c *ClientWithResponses) GetCustomizePagesWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetCustomizePagesResponse, error) {
	rsp, err := c.GetCustomizePages(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetCustomizePagesResponse(rsp)
}

// UpdateCustomizePagesWithBodyWithResponse request with arbitrary body returning *UpdateCustomizePagesResponse
func (c *ClientWithResponses) UpdateCustomizePagesWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateCustomizePagesResponse, error) {
	rsp, err := c.UpdateCustomizePagesWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateCustomizePagesResponse(rsp)
}

func (c *ClientWithResponses) UpdateCustomizePagesWithResponse(ctx context.Context, body UpdateCustomizePagesJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateCustomizePagesResponse, error) {
	rsp, err := c.UpdateCustomizePages(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateCustomizePagesResponse(rsp)
}

// GetEnvsWithResponse request returning *GetEnvsResponse
func (c *ClientWithResponses) GetEnvsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetEnvsResponse, error) {
	rsp, err := c.GetEnvs(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetEnvsResponse(rsp)
}

// CreateEnvWithBodyWithResponse request with arbitrary body returning *CreateEnvResponse
func (c *ClientWithResponses) CreateEnvWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateEnvResponse, error) {
	rsp, err := c.CreateEnvWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateEnvResponse(rsp)
}

func (c *ClientWithResponses) CreateEnvWithResponse(ctx context.Context, body CreateEnvJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateEnvResponse, error) {
	rsp, err := c.CreateEnv(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateEnvResponse(rsp)
}

// DeleteEnvWithResponse request returning *DeleteEnvResponse
func (c *ClientWithResponses) DeleteEnvWithResponse(ctx context.Context, envId EnvId, reqEditors ...RequestEditorFn) (*DeleteEnvResponse, error) {
	rsp, err := c.DeleteEnv(ctx, envId, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseDeleteEnvResponse(rsp)
}

// GetEnvWithResponse request returning *GetEnvResponse
func (c *ClientWithResponses) GetEnvWithResponse(ctx context.Context, envId EnvId, reqEditors ...RequestEditorFn) (*GetEnvResponse, error) {
	rsp, err := c.GetEnv(ctx, envId, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetEnvResponse(rsp)
}

// UpdateEnvWithBodyWithResponse request with arbitrary body returning *UpdateEnvResponse
func (c *ClientWithResponses) UpdateEnvWithBodyWithResponse(ctx context.Context, envId EnvId, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateEnvResponse, error) {
	rsp, err := c.UpdateEnvWithBody(ctx, envId, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateEnvResponse(rsp)
}

func (c *ClientWithResponses) UpdateEnvWithResponse(ctx context.Context, envId EnvId, body UpdateEnvJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateEnvResponse, error) {
	rsp, err := c.UpdateEnv(ctx, envId, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateEnvResponse(rsp)
}

// ReturnInternalServerErrorWithResponse request returning *ReturnInternalServerErrorResponse
func (c *ClientWithResponses) ReturnInternalServerErrorWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*ReturnInternalServerErrorResponse, error) {
	rsp, err := c.ReturnInternalServerError(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseReturnInternalServerErrorResponse(rsp)
}

// ConfirmExternalUserLinkWithBodyWithResponse request with arbitrary body returning *ConfirmExternalUserLinkResponse
func (c *ClientWithResponses) ConfirmExternalUserLinkWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*ConfirmExternalUserLinkResponse, error) {
	rsp, err := c.ConfirmExternalUserLinkWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseConfirmExternalUserLinkResponse(rsp)
}

func (c *ClientWithResponses) ConfirmExternalUserLinkWithResponse(ctx context.Context, body ConfirmExternalUserLinkJSONRequestBody, reqEditors ...RequestEditorFn) (*ConfirmExternalUserLinkResponse, error) {
	rsp, err := c.ConfirmExternalUserLink(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseConfirmExternalUserLinkResponse(rsp)
}

// RequestExternalUserLinkWithBodyWithResponse request with arbitrary body returning *RequestExternalUserLinkResponse
func (c *ClientWithResponses) RequestExternalUserLinkWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*RequestExternalUserLinkResponse, error) {
	rsp, err := c.RequestExternalUserLinkWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseRequestExternalUserLinkResponse(rsp)
}

func (c *ClientWithResponses) RequestExternalUserLinkWithResponse(ctx context.Context, body RequestExternalUserLinkJSONRequestBody, reqEditors ...RequestEditorFn) (*RequestExternalUserLinkResponse, error) {
	rsp, err := c.RequestExternalUserLink(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseRequestExternalUserLinkResponse(rsp)
}

// GetIdentityProvidersWithResponse request returning *GetIdentityProvidersResponse
func (c *ClientWithResponses) GetIdentityProvidersWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetIdentityProvidersResponse, error) {
	rsp, err := c.GetIdentityProviders(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetIdentityProvidersResponse(rsp)
}

// UpdateIdentityProviderWithBodyWithResponse request with arbitrary body returning *UpdateIdentityProviderResponse
func (c *ClientWithResponses) UpdateIdentityProviderWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateIdentityProviderResponse, error) {
	rsp, err := c.UpdateIdentityProviderWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateIdentityProviderResponse(rsp)
}

func (c *ClientWithResponses) UpdateIdentityProviderWithResponse(ctx context.Context, body UpdateIdentityProviderJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateIdentityProviderResponse, error) {
	rsp, err := c.UpdateIdentityProvider(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateIdentityProviderResponse(rsp)
}

// ValidateInvitationWithBodyWithResponse request with arbitrary body returning *ValidateInvitationResponse
func (c *ClientWithResponses) ValidateInvitationWithBodyWithResponse(ctx context.Context, invitationId InvitationId, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*ValidateInvitationResponse, error) {
	rsp, err := c.ValidateInvitationWithBody(ctx, invitationId, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseValidateInvitationResponse(rsp)
}

func (c *ClientWithResponses) ValidateInvitationWithResponse(ctx context.Context, invitationId InvitationId, body ValidateInvitationJSONRequestBody, reqEditors ...RequestEditorFn) (*ValidateInvitationResponse, error) {
	rsp, err := c.ValidateInvitation(ctx, invitationId, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseValidateInvitationResponse(rsp)
}

// GetInvitationValidityWithResponse request returning *GetInvitationValidityResponse
func (c *ClientWithResponses) GetInvitationValidityWithResponse(ctx context.Context, invitationId InvitationId, reqEditors ...RequestEditorFn) (*GetInvitationValidityResponse, error) {
	rsp, err := c.GetInvitationValidity(ctx, invitationId, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetInvitationValidityResponse(rsp)
}

// FindNotificationMessagesWithResponse request returning *FindNotificationMessagesResponse
func (c *ClientWithResponses) FindNotificationMessagesWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*FindNotificationMessagesResponse, error) {
	rsp, err := c.FindNotificationMessages(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseFindNotificationMessagesResponse(rsp)
}

// UpdateNotificationMessagesWithBodyWithResponse request with arbitrary body returning *UpdateNotificationMessagesResponse
func (c *ClientWithResponses) UpdateNotificationMessagesWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateNotificationMessagesResponse, error) {
	rsp, err := c.UpdateNotificationMessagesWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateNotificationMessagesResponse(rsp)
}

func (c *ClientWithResponses) UpdateNotificationMessagesWithResponse(ctx context.Context, body UpdateNotificationMessagesJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateNotificationMessagesResponse, error) {
	rsp, err := c.UpdateNotificationMessages(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateNotificationMessagesResponse(rsp)
}

// ResetPlanWithResponse request returning *ResetPlanResponse
func (c *ClientWithResponses) ResetPlanWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*ResetPlanResponse, error) {
	rsp, err := c.ResetPlan(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseResetPlanResponse(rsp)
}

// GetRolesWithResponse request returning *GetRolesResponse
func (c *ClientWithResponses) GetRolesWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetRolesResponse, error) {
	rsp, err := c.GetRoles(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetRolesResponse(rsp)
}

// CreateRoleWithBodyWithResponse request with arbitrary body returning *CreateRoleResponse
func (c *ClientWithResponses) CreateRoleWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateRoleResponse, error) {
	rsp, err := c.CreateRoleWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateRoleResponse(rsp)
}

func (c *ClientWithResponses) CreateRoleWithResponse(ctx context.Context, body CreateRoleJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateRoleResponse, error) {
	rsp, err := c.CreateRole(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateRoleResponse(rsp)
}

// DeleteRoleWithResponse request returning *DeleteRoleResponse
func (c *ClientWithResponses) DeleteRoleWithResponse(ctx context.Context, roleName RoleName, reqEditors ...RequestEditorFn) (*DeleteRoleResponse, error) {
	rsp, err := c.DeleteRole(ctx, roleName, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseDeleteRoleResponse(rsp)
}

// CreateSaasUserAttributeWithBodyWithResponse request with arbitrary body returning *CreateSaasUserAttributeResponse
func (c *ClientWithResponses) CreateSaasUserAttributeWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateSaasUserAttributeResponse, error) {
	rsp, err := c.CreateSaasUserAttributeWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateSaasUserAttributeResponse(rsp)
}

func (c *ClientWithResponses) CreateSaasUserAttributeWithResponse(ctx context.Context, body CreateSaasUserAttributeJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateSaasUserAttributeResponse, error) {
	rsp, err := c.CreateSaasUserAttribute(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateSaasUserAttributeResponse(rsp)
}

// GetSignInSettingsWithResponse request returning *GetSignInSettingsResponse
func (c *ClientWithResponses) GetSignInSettingsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetSignInSettingsResponse, error) {
	rsp, err := c.GetSignInSettings(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetSignInSettingsResponse(rsp)
}

// UpdateSignInSettingsWithBodyWithResponse request with arbitrary body returning *UpdateSignInSettingsResponse
func (c *ClientWithResponses) UpdateSignInSettingsWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateSignInSettingsResponse, error) {
	rsp, err := c.UpdateSignInSettingsWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateSignInSettingsResponse(rsp)
}

func (c *ClientWithResponses) UpdateSignInSettingsWithResponse(ctx context.Context, body UpdateSignInSettingsJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateSignInSettingsResponse, error) {
	rsp, err := c.UpdateSignInSettings(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateSignInSettingsResponse(rsp)
}

// SignUpWithBodyWithResponse request with arbitrary body returning *SignUpResponse
func (c *ClientWithResponses) SignUpWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*SignUpResponse, error) {
	rsp, err := c.SignUpWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseSignUpResponse(rsp)
}

func (c *ClientWithResponses) SignUpWithResponse(ctx context.Context, body SignUpJSONRequestBody, reqEditors ...RequestEditorFn) (*SignUpResponse, error) {
	rsp, err := c.SignUp(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseSignUpResponse(rsp)
}

// ResendSignUpConfirmationEmailWithBodyWithResponse request with arbitrary body returning *ResendSignUpConfirmationEmailResponse
func (c *ClientWithResponses) ResendSignUpConfirmationEmailWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*ResendSignUpConfirmationEmailResponse, error) {
	rsp, err := c.ResendSignUpConfirmationEmailWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseResendSignUpConfirmationEmailResponse(rsp)
}

func (c *ClientWithResponses) ResendSignUpConfirmationEmailWithResponse(ctx context.Context, body ResendSignUpConfirmationEmailJSONRequestBody, reqEditors ...RequestEditorFn) (*ResendSignUpConfirmationEmailResponse, error) {
	rsp, err := c.ResendSignUpConfirmationEmail(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseResendSignUpConfirmationEmailResponse(rsp)
}

// GetCloudFormationLaunchStackLinkForSingleTenantWithResponse request returning *GetCloudFormationLaunchStackLinkForSingleTenantResponse
func (c *ClientWithResponses) GetCloudFormationLaunchStackLinkForSingleTenantWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetCloudFormationLaunchStackLinkForSingleTenantResponse, error) {
	rsp, err := c.GetCloudFormationLaunchStackLinkForSingleTenant(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetCloudFormationLaunchStackLinkForSingleTenantResponse(rsp)
}

// GetSingleTenantSettingsWithResponse request returning *GetSingleTenantSettingsResponse
func (c *ClientWithResponses) GetSingleTenantSettingsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetSingleTenantSettingsResponse, error) {
	rsp, err := c.GetSingleTenantSettings(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetSingleTenantSettingsResponse(rsp)
}

// UpdateSingleTenantSettingsWithBodyWithResponse request with arbitrary body returning *UpdateSingleTenantSettingsResponse
func (c *ClientWithResponses) UpdateSingleTenantSettingsWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateSingleTenantSettingsResponse, error) {
	rsp, err := c.UpdateSingleTenantSettingsWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateSingleTenantSettingsResponse(rsp)
}

func (c *ClientWithResponses) UpdateSingleTenantSettingsWithResponse(ctx context.Context, body UpdateSingleTenantSettingsJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateSingleTenantSettingsResponse, error) {
	rsp, err := c.UpdateSingleTenantSettings(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateSingleTenantSettingsResponse(rsp)
}

// DeleteStripeTenantAndPricingWithResponse request returning *DeleteStripeTenantAndPricingResponse
func (c *ClientWithResponses) DeleteStripeTenantAndPricingWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*DeleteStripeTenantAndPricingResponse, error) {
	rsp, err := c.DeleteStripeTenantAndPricing(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseDeleteStripeTenantAndPricingResponse(rsp)
}

// CreateTenantAndPricingWithResponse request returning *CreateTenantAndPricingResponse
func (c *ClientWithResponses) CreateTenantAndPricingWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*CreateTenantAndPricingResponse, error) {
	rsp, err := c.CreateTenantAndPricing(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateTenantAndPricingResponse(rsp)
}

// GetTenantAttributesWithResponse request returning *GetTenantAttributesResponse
func (c *ClientWithResponses) GetTenantAttributesWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetTenantAttributesResponse, error) {
	rsp, err := c.GetTenantAttributes(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetTenantAttributesResponse(rsp)
}

// CreateTenantAttributeWithBodyWithResponse request with arbitrary body returning *CreateTenantAttributeResponse
func (c *ClientWithResponses) CreateTenantAttributeWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateTenantAttributeResponse, error) {
	rsp, err := c.CreateTenantAttributeWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateTenantAttributeResponse(rsp)
}

func (c *ClientWithResponses) CreateTenantAttributeWithResponse(ctx context.Context, body CreateTenantAttributeJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateTenantAttributeResponse, error) {
	rsp, err := c.CreateTenantAttribute(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateTenantAttributeResponse(rsp)
}

// DeleteTenantAttributeWithResponse request returning *DeleteTenantAttributeResponse
func (c *ClientWithResponses) DeleteTenantAttributeWithResponse(ctx context.Context, attributeName string, reqEditors ...RequestEditorFn) (*DeleteTenantAttributeResponse, error) {
	rsp, err := c.DeleteTenantAttribute(ctx, attributeName, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseDeleteTenantAttributeResponse(rsp)
}

// GetTenantsWithResponse request returning *GetTenantsResponse
func (c *ClientWithResponses) GetTenantsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetTenantsResponse, error) {
	rsp, err := c.GetTenants(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetTenantsResponse(rsp)
}

// CreateTenantWithBodyWithResponse request with arbitrary body returning *CreateTenantResponse
func (c *ClientWithResponses) CreateTenantWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateTenantResponse, error) {
	rsp, err := c.CreateTenantWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateTenantResponse(rsp)
}

func (c *ClientWithResponses) CreateTenantWithResponse(ctx context.Context, body CreateTenantJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateTenantResponse, error) {
	rsp, err := c.CreateTenant(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateTenantResponse(rsp)
}

// GetAllTenantUsersWithResponse request returning *GetAllTenantUsersResponse
func (c *ClientWithResponses) GetAllTenantUsersWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetAllTenantUsersResponse, error) {
	rsp, err := c.GetAllTenantUsers(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetAllTenantUsersResponse(rsp)
}

// GetAllTenantUserWithResponse request returning *GetAllTenantUserResponse
func (c *ClientWithResponses) GetAllTenantUserWithResponse(ctx context.Context, userId UserId, reqEditors ...RequestEditorFn) (*GetAllTenantUserResponse, error) {
	rsp, err := c.GetAllTenantUser(ctx, userId, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetAllTenantUserResponse(rsp)
}

// DeleteTenantWithResponse request returning *DeleteTenantResponse
func (c *ClientWithResponses) DeleteTenantWithResponse(ctx context.Context, tenantId TenantId, reqEditors ...RequestEditorFn) (*DeleteTenantResponse, error) {
	rsp, err := c.DeleteTenant(ctx, tenantId, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseDeleteTenantResponse(rsp)
}

// GetTenantWithResponse request returning *GetTenantResponse
func (c *ClientWithResponses) GetTenantWithResponse(ctx context.Context, tenantId TenantId, reqEditors ...RequestEditorFn) (*GetTenantResponse, error) {
	rsp, err := c.GetTenant(ctx, tenantId, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetTenantResponse(rsp)
}

// UpdateTenantWithBodyWithResponse request with arbitrary body returning *UpdateTenantResponse
func (c *ClientWithResponses) UpdateTenantWithBodyWithResponse(ctx context.Context, tenantId TenantId, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateTenantResponse, error) {
	rsp, err := c.UpdateTenantWithBody(ctx, tenantId, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateTenantResponse(rsp)
}

func (c *ClientWithResponses) UpdateTenantWithResponse(ctx context.Context, tenantId TenantId, body UpdateTenantJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateTenantResponse, error) {
	rsp, err := c.UpdateTenant(ctx, tenantId, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateTenantResponse(rsp)
}

// UpdateTenantBillingInfoWithBodyWithResponse request with arbitrary body returning *UpdateTenantBillingInfoResponse
func (c *ClientWithResponses) UpdateTenantBillingInfoWithBodyWithResponse(ctx context.Context, tenantId TenantId, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateTenantBillingInfoResponse, error) {
	rsp, err := c.UpdateTenantBillingInfoWithBody(ctx, tenantId, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateTenantBillingInfoResponse(rsp)
}

func (c *ClientWithResponses) UpdateTenantBillingInfoWithResponse(ctx context.Context, tenantId TenantId, body UpdateTenantBillingInfoJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateTenantBillingInfoResponse, error) {
	rsp, err := c.UpdateTenantBillingInfo(ctx, tenantId, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateTenantBillingInfoResponse(rsp)
}

// GetTenantIdentityProvidersWithResponse request returning *GetTenantIdentityProvidersResponse
func (c *ClientWithResponses) GetTenantIdentityProvidersWithResponse(ctx context.Context, tenantId TenantId, reqEditors ...RequestEditorFn) (*GetTenantIdentityProvidersResponse, error) {
	rsp, err := c.GetTenantIdentityProviders(ctx, tenantId, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetTenantIdentityProvidersResponse(rsp)
}

// UpdateTenantIdentityProviderWithBodyWithResponse request with arbitrary body returning *UpdateTenantIdentityProviderResponse
func (c *ClientWithResponses) UpdateTenantIdentityProviderWithBodyWithResponse(ctx context.Context, tenantId TenantId, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateTenantIdentityProviderResponse, error) {
	rsp, err := c.UpdateTenantIdentityProviderWithBody(ctx, tenantId, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateTenantIdentityProviderResponse(rsp)
}

func (c *ClientWithResponses) UpdateTenantIdentityProviderWithResponse(ctx context.Context, tenantId TenantId, body UpdateTenantIdentityProviderJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateTenantIdentityProviderResponse, error) {
	rsp, err := c.UpdateTenantIdentityProvider(ctx, tenantId, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateTenantIdentityProviderResponse(rsp)
}

// GetTenantInvitationsWithResponse request returning *GetTenantInvitationsResponse
func (c *ClientWithResponses) GetTenantInvitationsWithResponse(ctx context.Context, tenantId TenantId, reqEditors ...RequestEditorFn) (*GetTenantInvitationsResponse, error) {
	rsp, err := c.GetTenantInvitations(ctx, tenantId, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetTenantInvitationsResponse(rsp)
}

// CreateTenantInvitationWithBodyWithResponse request with arbitrary body returning *CreateTenantInvitationResponse
func (c *ClientWithResponses) CreateTenantInvitationWithBodyWithResponse(ctx context.Context, tenantId TenantId, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateTenantInvitationResponse, error) {
	rsp, err := c.CreateTenantInvitationWithBody(ctx, tenantId, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateTenantInvitationResponse(rsp)
}

func (c *ClientWithResponses) CreateTenantInvitationWithResponse(ctx context.Context, tenantId TenantId, body CreateTenantInvitationJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateTenantInvitationResponse, error) {
	rsp, err := c.CreateTenantInvitation(ctx, tenantId, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateTenantInvitationResponse(rsp)
}

// DeleteTenantInvitationWithResponse request returning *DeleteTenantInvitationResponse
func (c *ClientWithResponses) DeleteTenantInvitationWithResponse(ctx context.Context, tenantId TenantId, invitationId InvitationId, reqEditors ...RequestEditorFn) (*DeleteTenantInvitationResponse, error) {
	rsp, err := c.DeleteTenantInvitation(ctx, tenantId, invitationId, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseDeleteTenantInvitationResponse(rsp)
}

// GetTenantInvitationWithResponse request returning *GetTenantInvitationResponse
func (c *ClientWithResponses) GetTenantInvitationWithResponse(ctx context.Context, tenantId TenantId, invitationId InvitationId, reqEditors ...RequestEditorFn) (*GetTenantInvitationResponse, error) {
	rsp, err := c.GetTenantInvitation(ctx, tenantId, invitationId, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetTenantInvitationResponse(rsp)
}

// UpdateTenantPlanWithBodyWithResponse request with arbitrary body returning *UpdateTenantPlanResponse
func (c *ClientWithResponses) UpdateTenantPlanWithBodyWithResponse(ctx context.Context, tenantId TenantId, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateTenantPlanResponse, error) {
	rsp, err := c.UpdateTenantPlanWithBody(ctx, tenantId, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateTenantPlanResponse(rsp)
}

func (c *ClientWithResponses) UpdateTenantPlanWithResponse(ctx context.Context, tenantId TenantId, body UpdateTenantPlanJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateTenantPlanResponse, error) {
	rsp, err := c.UpdateTenantPlan(ctx, tenantId, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateTenantPlanResponse(rsp)
}

// GetStripeCustomerWithResponse request returning *GetStripeCustomerResponse
func (c *ClientWithResponses) GetStripeCustomerWithResponse(ctx context.Context, tenantId TenantId, reqEditors ...RequestEditorFn) (*GetStripeCustomerResponse, error) {
	rsp, err := c.GetStripeCustomer(ctx, tenantId, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetStripeCustomerResponse(rsp)
}

// GetTenantUsersWithResponse request returning *GetTenantUsersResponse
func (c *ClientWithResponses) GetTenantUsersWithResponse(ctx context.Context, tenantId TenantId, reqEditors ...RequestEditorFn) (*GetTenantUsersResponse, error) {
	rsp, err := c.GetTenantUsers(ctx, tenantId, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetTenantUsersResponse(rsp)
}

// CreateTenantUserWithBodyWithResponse request with arbitrary body returning *CreateTenantUserResponse
func (c *ClientWithResponses) CreateTenantUserWithBodyWithResponse(ctx context.Context, tenantId TenantId, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateTenantUserResponse, error) {
	rsp, err := c.CreateTenantUserWithBody(ctx, tenantId, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateTenantUserResponse(rsp)
}

func (c *ClientWithResponses) CreateTenantUserWithResponse(ctx context.Context, tenantId TenantId, body CreateTenantUserJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateTenantUserResponse, error) {
	rsp, err := c.CreateTenantUser(ctx, tenantId, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateTenantUserResponse(rsp)
}

// DeleteTenantUserWithResponse request returning *DeleteTenantUserResponse
func (c *ClientWithResponses) DeleteTenantUserWithResponse(ctx context.Context, tenantId TenantId, userId UserId, reqEditors ...RequestEditorFn) (*DeleteTenantUserResponse, error) {
	rsp, err := c.DeleteTenantUser(ctx, tenantId, userId, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseDeleteTenantUserResponse(rsp)
}

// GetTenantUserWithResponse request returning *GetTenantUserResponse
func (c *ClientWithResponses) GetTenantUserWithResponse(ctx context.Context, tenantId TenantId, userId UserId, reqEditors ...RequestEditorFn) (*GetTenantUserResponse, error) {
	rsp, err := c.GetTenantUser(ctx, tenantId, userId, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetTenantUserResponse(rsp)
}

// UpdateTenantUserWithBodyWithResponse request with arbitrary body returning *UpdateTenantUserResponse
func (c *ClientWithResponses) UpdateTenantUserWithBodyWithResponse(ctx context.Context, tenantId TenantId, userId UserId, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateTenantUserResponse, error) {
	rsp, err := c.UpdateTenantUserWithBody(ctx, tenantId, userId, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateTenantUserResponse(rsp)
}

func (c *ClientWithResponses) UpdateTenantUserWithResponse(ctx context.Context, tenantId TenantId, userId UserId, body UpdateTenantUserJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateTenantUserResponse, error) {
	rsp, err := c.UpdateTenantUser(ctx, tenantId, userId, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateTenantUserResponse(rsp)
}

// CreateTenantUserRolesWithBodyWithResponse request with arbitrary body returning *CreateTenantUserRolesResponse
func (c *ClientWithResponses) CreateTenantUserRolesWithBodyWithResponse(ctx context.Context, tenantId TenantId, userId UserId, envId EnvId, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateTenantUserRolesResponse, error) {
	rsp, err := c.CreateTenantUserRolesWithBody(ctx, tenantId, userId, envId, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateTenantUserRolesResponse(rsp)
}

func (c *ClientWithResponses) CreateTenantUserRolesWithResponse(ctx context.Context, tenantId TenantId, userId UserId, envId EnvId, body CreateTenantUserRolesJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateTenantUserRolesResponse, error) {
	rsp, err := c.CreateTenantUserRoles(ctx, tenantId, userId, envId, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateTenantUserRolesResponse(rsp)
}

// DeleteTenantUserRoleWithResponse request returning *DeleteTenantUserRoleResponse
func (c *ClientWithResponses) DeleteTenantUserRoleWithResponse(ctx context.Context, tenantId TenantId, userId UserId, envId EnvId, roleName RoleName, reqEditors ...RequestEditorFn) (*DeleteTenantUserRoleResponse, error) {
	rsp, err := c.DeleteTenantUserRole(ctx, tenantId, userId, envId, roleName, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseDeleteTenantUserRoleResponse(rsp)
}

// GetUserAttributesWithResponse request returning *GetUserAttributesResponse
func (c *ClientWithResponses) GetUserAttributesWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetUserAttributesResponse, error) {
	rsp, err := c.GetUserAttributes(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetUserAttributesResponse(rsp)
}

// CreateUserAttributeWithBodyWithResponse request with arbitrary body returning *CreateUserAttributeResponse
func (c *ClientWithResponses) CreateUserAttributeWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateUserAttributeResponse, error) {
	rsp, err := c.CreateUserAttributeWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateUserAttributeResponse(rsp)
}

func (c *ClientWithResponses) CreateUserAttributeWithResponse(ctx context.Context, body CreateUserAttributeJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateUserAttributeResponse, error) {
	rsp, err := c.CreateUserAttribute(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateUserAttributeResponse(rsp)
}

// DeleteUserAttributeWithResponse request returning *DeleteUserAttributeResponse
func (c *ClientWithResponses) DeleteUserAttributeWithResponse(ctx context.Context, attributeName string, reqEditors ...RequestEditorFn) (*DeleteUserAttributeResponse, error) {
	rsp, err := c.DeleteUserAttribute(ctx, attributeName, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseDeleteUserAttributeResponse(rsp)
}

// GetUserInfoWithResponse request returning *GetUserInfoResponse
func (c *ClientWithResponses) GetUserInfoWithResponse(ctx context.Context, params *GetUserInfoParams, reqEditors ...RequestEditorFn) (*GetUserInfoResponse, error) {
	rsp, err := c.GetUserInfo(ctx, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetUserInfoResponse(rsp)
}

// GetSaasUsersWithResponse request returning *GetSaasUsersResponse
func (c *ClientWithResponses) GetSaasUsersWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetSaasUsersResponse, error) {
	rsp, err := c.GetSaasUsers(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetSaasUsersResponse(rsp)
}

// CreateSaasUserWithBodyWithResponse request with arbitrary body returning *CreateSaasUserResponse
func (c *ClientWithResponses) CreateSaasUserWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateSaasUserResponse, error) {
	rsp, err := c.CreateSaasUserWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateSaasUserResponse(rsp)
}

func (c *ClientWithResponses) CreateSaasUserWithResponse(ctx context.Context, body CreateSaasUserJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateSaasUserResponse, error) {
	rsp, err := c.CreateSaasUser(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateSaasUserResponse(rsp)
}

// DeleteSaasUserWithResponse request returning *DeleteSaasUserResponse
func (c *ClientWithResponses) DeleteSaasUserWithResponse(ctx context.Context, userId UserId, reqEditors ...RequestEditorFn) (*DeleteSaasUserResponse, error) {
	rsp, err := c.DeleteSaasUser(ctx, userId, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseDeleteSaasUserResponse(rsp)
}

// GetSaasUserWithResponse request returning *GetSaasUserResponse
func (c *ClientWithResponses) GetSaasUserWithResponse(ctx context.Context, userId UserId, reqEditors ...RequestEditorFn) (*GetSaasUserResponse, error) {
	rsp, err := c.GetSaasUser(ctx, userId, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetSaasUserResponse(rsp)
}

// UpdateSaasUserAttributesWithBodyWithResponse request with arbitrary body returning *UpdateSaasUserAttributesResponse
func (c *ClientWithResponses) UpdateSaasUserAttributesWithBodyWithResponse(ctx context.Context, userId UserId, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateSaasUserAttributesResponse, error) {
	rsp, err := c.UpdateSaasUserAttributesWithBody(ctx, userId, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateSaasUserAttributesResponse(rsp)
}

func (c *ClientWithResponses) UpdateSaasUserAttributesWithResponse(ctx context.Context, userId UserId, body UpdateSaasUserAttributesJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateSaasUserAttributesResponse, error) {
	rsp, err := c.UpdateSaasUserAttributes(ctx, userId, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateSaasUserAttributesResponse(rsp)
}

// UpdateSaasUserEmailWithBodyWithResponse request with arbitrary body returning *UpdateSaasUserEmailResponse
func (c *ClientWithResponses) UpdateSaasUserEmailWithBodyWithResponse(ctx context.Context, userId UserId, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateSaasUserEmailResponse, error) {
	rsp, err := c.UpdateSaasUserEmailWithBody(ctx, userId, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateSaasUserEmailResponse(rsp)
}

func (c *ClientWithResponses) UpdateSaasUserEmailWithResponse(ctx context.Context, userId UserId, body UpdateSaasUserEmailJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateSaasUserEmailResponse, error) {
	rsp, err := c.UpdateSaasUserEmail(ctx, userId, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateSaasUserEmailResponse(rsp)
}

// ConfirmEmailUpdateWithBodyWithResponse request with arbitrary body returning *ConfirmEmailUpdateResponse
func (c *ClientWithResponses) ConfirmEmailUpdateWithBodyWithResponse(ctx context.Context, userId UserId, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*ConfirmEmailUpdateResponse, error) {
	rsp, err := c.ConfirmEmailUpdateWithBody(ctx, userId, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseConfirmEmailUpdateResponse(rsp)
}

func (c *ClientWithResponses) ConfirmEmailUpdateWithResponse(ctx context.Context, userId UserId, body ConfirmEmailUpdateJSONRequestBody, reqEditors ...RequestEditorFn) (*ConfirmEmailUpdateResponse, error) {
	rsp, err := c.ConfirmEmailUpdate(ctx, userId, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseConfirmEmailUpdateResponse(rsp)
}

// RequestEmailUpdateWithBodyWithResponse request with arbitrary body returning *RequestEmailUpdateResponse
func (c *ClientWithResponses) RequestEmailUpdateWithBodyWithResponse(ctx context.Context, userId UserId, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*RequestEmailUpdateResponse, error) {
	rsp, err := c.RequestEmailUpdateWithBody(ctx, userId, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseRequestEmailUpdateResponse(rsp)
}

func (c *ClientWithResponses) RequestEmailUpdateWithResponse(ctx context.Context, userId UserId, body RequestEmailUpdateJSONRequestBody, reqEditors ...RequestEditorFn) (*RequestEmailUpdateResponse, error) {
	rsp, err := c.RequestEmailUpdate(ctx, userId, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseRequestEmailUpdateResponse(rsp)
}

// GetUserMfaPreferenceWithResponse request returning *GetUserMfaPreferenceResponse
func (c *ClientWithResponses) GetUserMfaPreferenceWithResponse(ctx context.Context, userId UserId, reqEditors ...RequestEditorFn) (*GetUserMfaPreferenceResponse, error) {
	rsp, err := c.GetUserMfaPreference(ctx, userId, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetUserMfaPreferenceResponse(rsp)
}

// UpdateUserMfaPreferenceWithBodyWithResponse request with arbitrary body returning *UpdateUserMfaPreferenceResponse
func (c *ClientWithResponses) UpdateUserMfaPreferenceWithBodyWithResponse(ctx context.Context, userId UserId, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateUserMfaPreferenceResponse, error) {
	rsp, err := c.UpdateUserMfaPreferenceWithBody(ctx, userId, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateUserMfaPreferenceResponse(rsp)
}

func (c *ClientWithResponses) UpdateUserMfaPreferenceWithResponse(ctx context.Context, userId UserId, body UpdateUserMfaPreferenceJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateUserMfaPreferenceResponse, error) {
	rsp, err := c.UpdateUserMfaPreference(ctx, userId, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateUserMfaPreferenceResponse(rsp)
}

// UpdateSoftwareTokenWithBodyWithResponse request with arbitrary body returning *UpdateSoftwareTokenResponse
func (c *ClientWithResponses) UpdateSoftwareTokenWithBodyWithResponse(ctx context.Context, userId UserId, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateSoftwareTokenResponse, error) {
	rsp, err := c.UpdateSoftwareTokenWithBody(ctx, userId, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateSoftwareTokenResponse(rsp)
}

func (c *ClientWithResponses) UpdateSoftwareTokenWithResponse(ctx context.Context, userId UserId, body UpdateSoftwareTokenJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateSoftwareTokenResponse, error) {
	rsp, err := c.UpdateSoftwareToken(ctx, userId, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateSoftwareTokenResponse(rsp)
}

// CreateSecretCodeWithBodyWithResponse request with arbitrary body returning *CreateSecretCodeResponse
func (c *ClientWithResponses) CreateSecretCodeWithBodyWithResponse(ctx context.Context, userId UserId, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateSecretCodeResponse, error) {
	rsp, err := c.CreateSecretCodeWithBody(ctx, userId, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateSecretCodeResponse(rsp)
}

func (c *ClientWithResponses) CreateSecretCodeWithResponse(ctx context.Context, userId UserId, body CreateSecretCodeJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateSecretCodeResponse, error) {
	rsp, err := c.CreateSecretCode(ctx, userId, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateSecretCodeResponse(rsp)
}

// UpdateSaasUserPasswordWithBodyWithResponse request with arbitrary body returning *UpdateSaasUserPasswordResponse
func (c *ClientWithResponses) UpdateSaasUserPasswordWithBodyWithResponse(ctx context.Context, userId UserId, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateSaasUserPasswordResponse, error) {
	rsp, err := c.UpdateSaasUserPasswordWithBody(ctx, userId, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateSaasUserPasswordResponse(rsp)
}

func (c *ClientWithResponses) UpdateSaasUserPasswordWithResponse(ctx context.Context, userId UserId, body UpdateSaasUserPasswordJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateSaasUserPasswordResponse, error) {
	rsp, err := c.UpdateSaasUserPassword(ctx, userId, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateSaasUserPasswordResponse(rsp)
}

// UnlinkProviderWithResponse request returning *UnlinkProviderResponse
func (c *ClientWithResponses) UnlinkProviderWithResponse(ctx context.Context, userId UserId, providerName string, reqEditors ...RequestEditorFn) (*UnlinkProviderResponse, error) {
	rsp, err := c.UnlinkProvider(ctx, userId, providerName, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUnlinkProviderResponse(rsp)
}

// ParseGetAuthInfoResponse parses an HTTP response from a GetAuthInfoWithResponse call
