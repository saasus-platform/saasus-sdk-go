package authapi

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

const (
	softwareTokenSecretKey  = "software_token_secret"
	softwareTokenSessionKey = "software_token_session"
)

var totpOpts = totp.ValidateOpts{
	Period:    30,
	Skew:      1,
	Digits:    otp.DigitsSix,
	Algorithm: otp.AlgorithmSHA1,
}

// ensureCognitoSoftwareTokenPrepared fetches a TOTP secret directly from Cognito so that
// tests can generate valid verification codes before hitting SaaSus API endpoints.
func ensureCognitoSoftwareTokenPrepared(variables map[string]any) {
	if secret, _ := variables[softwareTokenSecretKey].(string); secret != "" {
		return
	}

	accessToken, _ := variables["cognito_access_token"].(string)
	if accessToken == "" {
		return
	}

	cfg, err := LoadCognitoConfigFromEnv()
	if err != nil {
		fmt.Printf("Warning: unable to load Cognito config for MFA preparation: %v\n", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := newCognitoIDPClient(ctx, cfg)
	if err != nil {
		fmt.Printf("Warning: failed to initialize Cognito client for MFA preparation: %v\n", err)
		return
	}

	resp, err := client.AssociateSoftwareToken(ctx, &cognitoidentityprovider.AssociateSoftwareTokenInput{
		AccessToken: aws.String(accessToken),
	})
	if err != nil {
		fmt.Printf("Warning: AssociateSoftwareToken pre-step failed: %v\n", err)
		return
	}

	secret := aws.ToString(resp.SecretCode)
	if secret == "" {
		fmt.Printf("Warning: AssociateSoftwareToken returned empty secret\n")
		return
	}

	variables[softwareTokenSecretKey] = secret
	if resp.Session != nil && len(*resp.Session) > 0 {
		variables[softwareTokenSessionKey] = aws.ToString(resp.Session)
	}

	fmt.Printf("[DEBUG] Prepared Cognito software token secret via SDK\n")
}

func newCognitoIDPClient(ctx context.Context, cfg *CognitoConfig) (*cognitoidentityprovider.Client, error) {
	awsCfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(cfg.Region))
	if err != nil {
		return nil, fmt.Errorf("load aws config: %w", err)
	}

	client := cognitoidentityprovider.NewFromConfig(awsCfg, func(o *cognitoidentityprovider.Options) {
		if cfg.Endpoint != "" {
			o.BaseEndpoint = aws.String(cfg.Endpoint)
		}
	})
	return client, nil
}

func generateCurrentTotp(secret string) (string, error) {
	if secret == "" {
		return "", fmt.Errorf("software token secret is empty")
	}

	return totp.GenerateCodeCustom(secret, time.Now(), totpOpts)
}
