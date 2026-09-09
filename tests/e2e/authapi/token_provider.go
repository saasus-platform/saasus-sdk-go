package authapi

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/joho/godotenv"
)

const (
	envCognitoUserPoolID = "E2E_COGNITO_USER_POOL_ID"
	envCognitoClientID   = "E2E_COGNITO_CLIENT_ID"
	envCognitoUsername   = "E2E_COGNITO_USERNAME"
	envCognitoPassword   = "E2E_COGNITO_PASSWORD"
	envCognitoRegion     = "E2E_COGNITO_REGION"
	envCognitoEndpoint   = "E2E_COGNITO_ENDPOINT"
	defaultCognitoRegion = "ap-northeast-1"
)

// AuthTokens keeps Cognito issued tokens that are reused across stories.
type AuthTokens struct {
	AccessToken  string
	IDToken      string
	RefreshToken string
}

var (
	tokenOnce sync.Once
	tokenSet  *AuthTokens
	tokenErr  error
)

// MustGetAuthTokens fetches Cognito tokens once per test run and fails the test if it cannot obtain them.
func MustGetAuthTokens(t *testing.T) *AuthTokens {
	t.Helper()

	tokenOnce.Do(func() {
		// Load .env file if it exists
		_ = godotenv.Load("../../../.env")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		cfg, err := LoadCognitoConfigFromEnv()
		if err != nil {
			tokenErr = err
			return
		}

		tokenSet, tokenErr = fetchCognitoTokens(ctx, cfg)
	})

	if tokenErr != nil {
		t.Skipf("Skipping test due to Cognito configuration issue: %v", tokenErr)
	}

	return tokenSet
}

// CognitoConfig holds configuration for Cognito interaction.
type CognitoConfig struct {
	UserPoolID string
	ClientID   string
	Username   string
	Password   string
	Region     string
	Endpoint   string
}

// LoadCognitoConfigFromEnv loads Cognito configuration from environment variables.
func LoadCognitoConfigFromEnv() (*CognitoConfig, error) {
	missing := make([]string, 0, 4)

	cfg := &CognitoConfig{
		UserPoolID: os.Getenv(envCognitoUserPoolID),
		ClientID:   os.Getenv(envCognitoClientID),
		Username:   os.Getenv(envCognitoUsername),
		Password:   os.Getenv(envCognitoPassword),
		Region:     os.Getenv(envCognitoRegion),
		Endpoint:   os.Getenv(envCognitoEndpoint),
	}

	if cfg.Region == "" {
		cfg.Region = defaultCognitoRegion
	}

	if cfg.UserPoolID == "" {
		missing = append(missing, envCognitoUserPoolID)
	}
	if cfg.ClientID == "" {
		missing = append(missing, envCognitoClientID)
	}
	if cfg.Username == "" {
		missing = append(missing, envCognitoUsername)
	}
	if cfg.Password == "" {
		missing = append(missing, envCognitoPassword)
	}

	if len(missing) > 0 {
		return nil, fmt.Errorf("missing required Cognito env vars: %v", missing)
	}

	return cfg, nil
}

func fetchCognitoTokens(ctx context.Context, cfg *CognitoConfig) (*AuthTokens, error) {
	awsCfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(cfg.Region))
	if err != nil {
		return nil, fmt.Errorf("load aws config: %w", err)
	}

	clientOpts := func(o *cognitoidentityprovider.Options) {
		if cfg.Endpoint != "" {
			o.BaseEndpoint = aws.String(cfg.Endpoint)
		}
	}

	client := cognitoidentityprovider.NewFromConfig(awsCfg, clientOpts)

	input := &cognitoidentityprovider.AdminInitiateAuthInput{
		AuthFlow: types.AuthFlowTypeAdminUserPasswordAuth,
		AuthParameters: map[string]string{
			"USERNAME": cfg.Username,
			"PASSWORD": cfg.Password,
		},
		ClientId:   aws.String(cfg.ClientID),
		UserPoolId: aws.String(cfg.UserPoolID),
	}

	output, err := client.AdminInitiateAuth(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("admin initiate auth: %w", err)
	}

	if output.AuthenticationResult == nil {
		return nil, fmt.Errorf("authentication result is empty")
	}

	return &AuthTokens{
		AccessToken:  aws.ToString(output.AuthenticationResult.AccessToken),
		IDToken:      aws.ToString(output.AuthenticationResult.IdToken),
		RefreshToken: aws.ToString(output.AuthenticationResult.RefreshToken),
	}, nil
}
