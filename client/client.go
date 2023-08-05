package client

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Anti-Pattern-Inc/saasus-sdk-go/generated/authapi"
)

func WithSaasusSigV1() authapi.RequestEditorFn {
	return func(ctx context.Context, req *http.Request) error {
		secret := os.Getenv("SAASUS_SECRET_KEY")
		apiKey := os.Getenv("SAASUS_API_KEY")
		saasID := os.Getenv("SAASUS_SAAS_ID")
		literal := "SAASUSSIGV1"

		if secret == "" || apiKey == "" || saasID == "" {
			return fmt.Errorf("invalid request: SAASUS_SECRET_KEY, SAASUS_API_KEY, SAASUS_SAAS_ID are required")
		}

		signatureHmac := hmac.New(sha256.New, []byte(secret))

		now := time.Now().UTC().Format("200601021504")
		signatureHmac.Write([]byte(now))
		signatureHmac.Write([]byte(apiKey))
		signatureHmac.Write([]byte(strings.ToUpper(req.Method))) // HTTP method

		hostURI := strings.Split(req.URL.String(), "//")
		if len(hostURI) < 2 {
			return fmt.Errorf("invalid URL format")
		}

		signatureHmac.Write([]byte(hostURI[1])) // host+URI

		if req.Body != nil {
			bodyBytes, err := io.ReadAll(req.Body)
			if err != nil {
				return err
			}
			signatureHmac.Write(bodyBytes)
			// Restore the io.ReadCloser to its original state
			req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		}

		authorization := fmt.Sprintf("%s Sig=%s, SaaSID=%s, APIKey=%s", literal, hex.EncodeToString(signatureHmac.Sum(nil)), saasID, apiKey)
		req.Header.Set("Authorization", authorization)

		return nil
	}
}
