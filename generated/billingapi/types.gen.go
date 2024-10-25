// Package billingapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.12.4 DO NOT EDIT.
package billingapi

const (
	BearerScopes = "Bearer.Scopes"
)

// Error defines model for Error.
type Error struct {
	// Message Error message
	Message string `json:"message"`

	// Type permission_denied
	Type string `json:"type"`
}

// StripeInfo defines model for StripeInfo.
type StripeInfo struct {
	IsRegistered bool `json:"is_registered"`
}

// UpdateStripeInfoParam defines model for UpdateStripeInfoParam.
type UpdateStripeInfoParam struct {
	// SecretKey secret key
	SecretKey string `json:"secret_key"`
}

// UpdateStripeInfoJSONRequestBody defines body for UpdateStripeInfo for application/json ContentType.
type UpdateStripeInfoJSONRequestBody = UpdateStripeInfoParam
