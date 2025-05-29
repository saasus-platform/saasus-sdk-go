type GetAuthInfoResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *AuthInfo
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r GetAuthInfoResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetAuthInfoResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type UpdateAuthInfoResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r UpdateAuthInfoResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r UpdateAuthInfoResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type LinkAwsMarketplaceResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r LinkAwsMarketplaceResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r LinkAwsMarketplaceResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type SignUpWithAwsMarketplaceResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON201      *SaasUser
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r SignUpWithAwsMarketplaceResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r SignUpWithAwsMarketplaceResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type ConfirmSignUpWithAwsMarketplaceResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON201      *Tenant
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r ConfirmSignUpWithAwsMarketplaceResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r ConfirmSignUpWithAwsMarketplaceResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetBasicInfoResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *BasicInfo
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r GetBasicInfoResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetBasicInfoResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type UpdateBasicInfoResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r UpdateBasicInfoResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r UpdateBasicInfoResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetAuthCredentialsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *Credentials
	JSON404      *Error
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r GetAuthCredentialsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetAuthCredentialsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type CreateAuthCredentialsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON201      *AuthorizationTempCode
	JSON400      *Error
	JSON401      *Error
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r CreateAuthCredentialsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r CreateAuthCredentialsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetCustomizePageSettingsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *CustomizePageSettings
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r GetCustomizePageSettingsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetCustomizePageSettingsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type UpdateCustomizePageSettingsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r UpdateCustomizePageSettingsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r UpdateCustomizePageSettingsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetCustomizePagesResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *CustomizePages
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r GetCustomizePagesResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetCustomizePagesResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type UpdateCustomizePagesResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r UpdateCustomizePagesResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r UpdateCustomizePagesResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetEnvsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *Envs
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r GetEnvsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetEnvsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type CreateEnvResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON201      *Env
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r CreateEnvResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r CreateEnvResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type DeleteEnvResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON404      *Error
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r DeleteEnvResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r DeleteEnvResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetEnvResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *Env
	JSON404      *Error
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r GetEnvResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetEnvResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type UpdateEnvResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON404      *Error
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r UpdateEnvResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r UpdateEnvResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type ReturnInternalServerErrorResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r ReturnInternalServerErrorResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r ReturnInternalServerErrorResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type ConfirmExternalUserLinkResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON400      *Error
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r ConfirmExternalUserLinkResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r ConfirmExternalUserLinkResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type RequestExternalUserLinkResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON400      *Error
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r RequestExternalUserLinkResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r RequestExternalUserLinkResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetIdentityProvidersResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *IdentityProviders
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r GetIdentityProvidersResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetIdentityProvidersResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type UpdateIdentityProviderResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r UpdateIdentityProviderResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r UpdateIdentityProviderResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type ValidateInvitationResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON400      *Error
	JSON404      *Error
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r ValidateInvitationResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r ValidateInvitationResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetInvitationValidityResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *InvitationValidity
	JSON404      *Error
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r GetInvitationValidityResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetInvitationValidityResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type FindNotificationMessagesResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *NotificationMessages
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r FindNotificationMessagesResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r FindNotificationMessagesResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type UpdateNotificationMessagesResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r UpdateNotificationMessagesResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r UpdateNotificationMessagesResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type ResetPlanResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r ResetPlanResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r ResetPlanResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetRolesResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *Roles
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r GetRolesResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetRolesResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type CreateRoleResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON201      *Role
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r CreateRoleResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r CreateRoleResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type DeleteRoleResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON400      *Error
	JSON404      *Error
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r DeleteRoleResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r DeleteRoleResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type CreateSaasUserAttributeResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON201      *Attribute
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r CreateSaasUserAttributeResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r CreateSaasUserAttributeResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetSignInSettingsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *SignInSettings
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r GetSignInSettingsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetSignInSettingsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type UpdateSignInSettingsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r UpdateSignInSettingsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r UpdateSignInSettingsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type SignUpResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON201      *SaasUser
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r SignUpResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r SignUpResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type ResendSignUpConfirmationEmailResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r ResendSignUpConfirmationEmailResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r ResendSignUpConfirmationEmailResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetCloudFormationLaunchStackLinkForSingleTenantResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *CloudFormationLaunchStackLink
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r GetCloudFormationLaunchStackLinkForSingleTenantResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetCloudFormationLaunchStackLinkForSingleTenantResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetSingleTenantSettingsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *SingleTenantSettings
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r GetSingleTenantSettingsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetSingleTenantSettingsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type UpdateSingleTenantSettingsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON400      *Error
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r UpdateSingleTenantSettingsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r UpdateSingleTenantSettingsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type DeleteStripeTenantAndPricingResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r DeleteStripeTenantAndPricingResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r DeleteStripeTenantAndPricingResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type CreateTenantAndPricingResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r CreateTenantAndPricingResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r CreateTenantAndPricingResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetTenantAttributesResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *TenantAttributes
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r GetTenantAttributesResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetTenantAttributesResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type CreateTenantAttributeResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON201      *Attribute
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r CreateTenantAttributeResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r CreateTenantAttributeResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type DeleteTenantAttributeResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON404      *Error
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r DeleteTenantAttributeResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r DeleteTenantAttributeResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetTenantsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *Tenants
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r GetTenantsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetTenantsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type CreateTenantResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON201      *Tenant
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r CreateTenantResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r CreateTenantResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetAllTenantUsersResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *Users
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r GetAllTenantUsersResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetAllTenantUsersResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetAllTenantUserResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *Users
	JSON404      *Error
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r GetAllTenantUserResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetAllTenantUserResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type DeleteTenantResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r DeleteTenantResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r DeleteTenantResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetTenantResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *TenantDetail
	JSON404      *Error
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r GetTenantResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetTenantResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type UpdateTenantResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON404      *Error
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r UpdateTenantResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r UpdateTenantResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type UpdateTenantBillingInfoResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON400      *Error
	JSON404      *Error
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r UpdateTenantBillingInfoResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r UpdateTenantBillingInfoResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetTenantIdentityProvidersResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *TenantIdentityProviders
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r GetTenantIdentityProvidersResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetTenantIdentityProvidersResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type UpdateTenantIdentityProviderResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r UpdateTenantIdentityProviderResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r UpdateTenantIdentityProviderResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetTenantInvitationsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *Invitations
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r GetTenantInvitationsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetTenantInvitationsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type CreateTenantInvitationResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON201      *Invitation
	JSON400      *Error
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r CreateTenantInvitationResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r CreateTenantInvitationResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type DeleteTenantInvitationResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON404      *Error
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r DeleteTenantInvitationResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r DeleteTenantInvitationResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetTenantInvitationResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *Invitation
	JSON404      *Error
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r GetTenantInvitationResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetTenantInvitationResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type UpdateTenantPlanResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON400      *Error
	JSON404      *Error
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r UpdateTenantPlanResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r UpdateTenantPlanResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetStripeCustomerResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *StripeCustomer
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r GetStripeCustomerResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetStripeCustomerResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetTenantUsersResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *Users
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r GetTenantUsersResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetTenantUsersResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type CreateTenantUserResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON201      *User
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r CreateTenantUserResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r CreateTenantUserResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type DeleteTenantUserResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON404      *Error
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r DeleteTenantUserResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r DeleteTenantUserResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetTenantUserResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *User
	JSON404      *Error
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r GetTenantUserResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetTenantUserResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type UpdateTenantUserResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON404      *Error
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r UpdateTenantUserResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r UpdateTenantUserResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type CreateTenantUserRolesResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r CreateTenantUserRolesResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r CreateTenantUserRolesResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type DeleteTenantUserRoleResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON404      *Error
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r DeleteTenantUserRoleResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r DeleteTenantUserRoleResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetUserAttributesResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *UserAttributes
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r GetUserAttributesResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetUserAttributesResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type CreateUserAttributeResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON201      *Attribute
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r CreateUserAttributeResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r CreateUserAttributeResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type DeleteUserAttributeResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON404      *Error
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r DeleteUserAttributeResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r DeleteUserAttributeResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetUserInfoResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *UserInfo
	JSON401      *Error
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r GetUserInfoResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetUserInfoResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetSaasUsersResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *SaasUsers
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r GetSaasUsersResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetSaasUsersResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type CreateSaasUserResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON201      *SaasUser
	JSON400      *Error
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r CreateSaasUserResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r CreateSaasUserResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type DeleteSaasUserResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON404      *Error
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r DeleteSaasUserResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r DeleteSaasUserResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetSaasUserResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *SaasUser
	JSON404      *Error
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r GetSaasUserResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetSaasUserResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type UpdateSaasUserAttributesResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r UpdateSaasUserAttributesResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r UpdateSaasUserAttributesResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type UpdateSaasUserEmailResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r UpdateSaasUserEmailResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r UpdateSaasUserEmailResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type ConfirmEmailUpdateResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r ConfirmEmailUpdateResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r ConfirmEmailUpdateResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type RequestEmailUpdateResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r RequestEmailUpdateResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r RequestEmailUpdateResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetUserMfaPreferenceResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *MfaPreference
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r GetUserMfaPreferenceResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetUserMfaPreferenceResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type UpdateUserMfaPreferenceResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r UpdateUserMfaPreferenceResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r UpdateUserMfaPreferenceResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type UpdateSoftwareTokenResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r UpdateSoftwareTokenResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r UpdateSoftwareTokenResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type CreateSecretCodeResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON201      *SoftwareTokenSecretCode
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r CreateSecretCodeResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r CreateSecretCodeResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type UpdateSaasUserPasswordResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r UpdateSaasUserPasswordResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r UpdateSaasUserPasswordResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type UnlinkProviderResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r UnlinkProviderResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r UnlinkProviderResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// GetAuthInfoWithResponse request returning *GetAuthInfoResponse
