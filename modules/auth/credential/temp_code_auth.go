package credential

import (
	"context"
	"fmt"

	"github.com/saasus-platform/saasus-sdk-go/generated/authapi"
	"github.com/saasus-platform/saasus-sdk-go/modules/auth"
)

// GetAuthCredentialsWithTempCodeAuth executes GetAuthCredentials API with tempCodeAuth.
func GetAuthCredentialsWithTempCodeAuth(ctx context.Context, code string) (*authapi.Credentials, error) {
	if code == "" {
		return nil, fmt.Errorf("code is not provided by query parameter")
	}

	authClientWithResponse, err := auth.AuthWithResponse()
	if err != nil {
		return nil, err
	}

	authFlow := authapi.GetAuthCredentialsParamsAuthFlow("tempCodeAuth")
	req := &authapi.GetAuthCredentialsParams{
		Code:         &code,
		AuthFlow:     &authFlow,
		RefreshToken: nil,
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
