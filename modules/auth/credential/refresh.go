package credential

import (
	"context"
	"fmt"
	"net/http"

	"github.com/saasus-platform/saasus-sdk-go/generated/authapi"
	"github.com/saasus-platform/saasus-sdk-go/modules/auth"
)

// GetAuthCredentialsByRefreshTokenAuth executes GetAuthCredentials API with refreshTokenAuth.
func GetAuthCredentialsWithRefreshTokenAuth(ctx context.Context, w http.ResponseWriter, r *http.Request, refreshToken string) (*authapi.Credentials, error) {
	if refreshToken == "" {
		return nil, fmt.Errorf("saasus_refresh_token cookie is required")
	}

	authClientWithResponse, err := auth.AuthWithResponse()
	if err != nil {
		return nil, err
	}

	authFlow := authapi.GetAuthCredentialsParamsAuthFlow("refreshTokenAuth")
	req := &authapi.GetAuthCredentialsParams{
		Code:         nil,
		AuthFlow:     &authFlow,
		RefreshToken: &refreshToken,
	}
	res, err := authClientWithResponse.GetAuthCredentialsWithResponse(ctx, req)
	if err != nil {
		return nil, err
	}
	if res.JSON500 != nil {
		return nil, fmt.Errorf("data = %s, message = %s, type = %s", res.JSON500.Data, res.JSON500.Message, res.JSON500.Type)
	}
	if res.JSON404 != nil {
		return nil, fmt.Errorf("data = %s, message = %s, type = %s", res.JSON404.Data, res.JSON404.Message, res.JSON404.Type)
	}
	if res.JSON200 == nil {
		return nil, fmt.Errorf("response failed")
	}

	return res.JSON200, nil
}
