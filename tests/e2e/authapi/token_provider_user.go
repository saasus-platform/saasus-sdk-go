package authapi

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

// UserAuthTokens keeps per-user Cognito tokens.
type UserAuthTokens struct {
	AccessToken  string
	IdToken      string
	RefreshToken string
}

// ObtainUserTokens authenticates a SaaS user via Cognito and returns tokens.
func ObtainUserTokens(email, password string) (*UserAuthTokens, error) {
	cfg, err := LoadCognitoConfigFromEnv()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	awsCfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(cfg.Region))
	if err != nil {
		return nil, fmt.Errorf("load aws config: %w", err)
	}

	client := cognitoidentityprovider.NewFromConfig(awsCfg, func(o *cognitoidentityprovider.Options) {
		if cfg.Endpoint != "" {
			o.BaseEndpoint = aws.String(cfg.Endpoint)
		}
	})

	input := &cognitoidentityprovider.AdminInitiateAuthInput{
		AuthFlow:   types.AuthFlowTypeAdminUserPasswordAuth,
		ClientId:   aws.String(cfg.ClientID),
		UserPoolId: aws.String(cfg.UserPoolID),
		AuthParameters: map[string]string{
			"USERNAME": email,
			"PASSWORD": password,
		},
	}

	resp, err := client.AdminInitiateAuth(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("admin initiate auth: %w", err)
	}

	if resp.ChallengeName == types.ChallengeNameTypeNewPasswordRequired {
		challengeInput := &cognitoidentityprovider.RespondToAuthChallengeInput{
			ChallengeName: resp.ChallengeName,
			ClientId:      aws.String(cfg.ClientID),
			ChallengeResponses: map[string]string{
				"USERNAME":     email,
				"NEW_PASSWORD": password,
			},
			Session: resp.Session,
		}
		challengeResp, err := client.RespondToAuthChallenge(ctx, challengeInput)
		if err != nil {
			return nil, fmt.Errorf("respond to auth challenge: %w", err)
		}
		resp.AuthenticationResult = challengeResp.AuthenticationResult
	}

	if resp.AuthenticationResult == nil {
		return nil, fmt.Errorf("empty authentication result")
	}

	return &UserAuthTokens{
		AccessToken:  aws.ToString(resp.AuthenticationResult.AccessToken),
		IdToken:      aws.ToString(resp.AuthenticationResult.IdToken),
		RefreshToken: aws.ToString(resp.AuthenticationResult.RefreshToken),
	}, nil
}
