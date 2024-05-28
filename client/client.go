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

	"github.com/saasus-platform/saasus-sdk-go/ctxlib"
)

// SetSigV1 sets the signature to the request header.
//
// ref: https://docs.saasus.io/reference/getting-started-with-your-api
func SetSigV1(req *http.Request) error {
	secret := os.Getenv("SAASUS_SECRET_KEY")
	apiKey := os.Getenv("SAASUS_API_KEY")
	saasID := os.Getenv("SAASUS_SAAS_ID")
	literal := "SAASUSSIGV1"

	if secret == "" || apiKey == "" || saasID == "" {
		return fmt.Errorf("invalid request: SAASUS_SECRET_KEY, SAASUS_API_KEY, SAASUS_SAAS_ID are required")
	}

	sigV1 := &sigV1{
		nowStr:          time.Now().UTC().Format("200601021504"),
		secret:          secret,
		apiKey:          apiKey,
		literal:         literal,
		saasID:          saasID,
		httpMethod:      req.Method,
		httpURL:         req.URL.String(),
		httpRequestBody: req.Body,
	}

	sig, err := sigV1.generate()
	if err != nil {
		return err
	}
	req.Body = sigV1.httpRequestBody

	req.Header.Set("Authorization", sig)
	return nil
}

type sigV1 struct {
	nowStr          string
	secret          string
	apiKey          string
	literal         string
	saasID          string
	httpMethod      string
	httpURL         string
	httpRequestBody io.ReadCloser
}

func (s *sigV1) generate() (string, error) {
	signatureHmac := hmac.New(sha256.New, []byte(s.secret))
	signatureHmac.Write([]byte(s.nowStr))
	signatureHmac.Write([]byte(s.apiKey))
	signatureHmac.Write([]byte(strings.ToUpper(s.httpMethod)))

	hostURI := strings.Split(s.httpURL, "//")
	if len(hostURI) < 2 {
		return "", fmt.Errorf("invalid URL format")
	}

	signatureHmac.Write([]byte(hostURI[1]))

	if s.httpRequestBody != nil {
		bodyBytes, err := io.ReadAll(s.httpRequestBody)
		if err != nil {
			return "", err
		}
		signatureHmac.Write(bodyBytes)
		s.httpRequestBody = io.NopCloser(bytes.NewReader(bodyBytes))
	}

	sig := fmt.Sprintf("%s Sig=%s, SaaSID=%s, APIKey=%s", s.literal, hex.EncodeToString(signatureHmac.Sum(nil)), s.saasID, s.apiKey)
	return sig, nil
}

// SetReferer sets the referer to the request header if existed.
func SetReferer(ctx context.Context, req *http.Request) {
	if referer, ok := ctx.Value(ctxlib.RefererKey).(string); ok {
		req.Header.Set("Referer", referer)
	}

	if xSaaSusReferer, ok := ctx.Value(ctxlib.XSaaSusRefererKey).(string); ok {
		req.Header.Set("X-SaaSus-Referer", xSaaSusReferer)
	}
}
