package authapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/saasus-platform/saasus-sdk-go/generated/authapi"
	"github.com/saasus-platform/saasus-sdk-go/generated/billingapi"
	"github.com/saasus-platform/saasus-sdk-go/modules/billing"
)

type authTenantTracker struct {
	mu       sync.Mutex
	tenantID string
}

var lastAuthTenant authTenantTracker

func setLastTenantID(id string) {
	lastAuthTenant.mu.Lock()
	defer lastAuthTenant.mu.Unlock()
	lastAuthTenant.tenantID = id
}

func getLastTenantID() string {
	lastAuthTenant.mu.Lock()
	defer lastAuthTenant.mu.Unlock()
	return lastAuthTenant.tenantID
}

// AuthClientWithStripeFallback wraps the auth client to add Stripe recovery
// logic for CreateTenantAndPricing calls.
type AuthClientWithStripeFallback struct {
	*authapi.ClientWithResponses
}

func NewAuthClientWithStripeFallback(client *authapi.ClientWithResponses) *AuthClientWithStripeFallback {
	return &AuthClientWithStripeFallback{ClientWithResponses: client}
}

func (c *AuthClientWithStripeFallback) CreateTenantAndPricing(ctx context.Context, reqEditors ...authapi.RequestEditorFn) (*http.Response, error) {
	resp, err := c.ClientWithResponses.CreateTenantAndPricing(ctx, reqEditors...)
	if err != nil || resp == nil {
		return resp, err
	}
	if isSuccessStatus(resp.StatusCode) {
		return resp, nil
	}

	bodyStr := readAndRestoreBody(resp)
	if shouldRetryAfterStripeKeyError(resp.StatusCode, bodyStr) {
		if err := updateStripeInfo(ctx); err == nil {
			retryResp, retryErr := c.ClientWithResponses.CreateTenantAndPricing(ctx, reqEditors...)
			if retryErr == nil && retryResp != nil && isSuccessStatus(retryResp.StatusCode) {
				return retryResp, nil
			}
			if retryErr != nil {
				return retryResp, retryErr
			}
			resp = retryResp
			bodyStr = readAndRestoreBody(resp)
		}
	}

	if shouldFallbackToStripeCustomer(resp.StatusCode, bodyStr) {
		if okResp := c.tryGetStripeCustomer(ctx, reqEditors...); okResp != nil {
			return okResp, nil
		}
	}

	return resp, nil
}

func (c *AuthClientWithStripeFallback) CreateTenantAndPricingWithResponse(ctx context.Context, reqEditors ...authapi.RequestEditorFn) (*authapi.CreateTenantAndPricingResponse, error) {
	resp, err := c.ClientWithResponses.CreateTenantAndPricingWithResponse(ctx, reqEditors...)
	if err != nil || resp == nil {
		return resp, err
	}
	if isSuccessStatus(resp.StatusCode()) {
		return resp, nil
	}

	bodyStr := string(resp.Body)
	if shouldRetryAfterStripeKeyError(resp.StatusCode(), bodyStr) {
		if err := updateStripeInfo(ctx); err == nil {
			retryResp, retryErr := c.ClientWithResponses.CreateTenantAndPricingWithResponse(ctx, reqEditors...)
			if retryErr == nil && retryResp != nil && isSuccessStatus(retryResp.StatusCode()) {
				return retryResp, nil
			}
			if retryErr != nil {
				return retryResp, retryErr
			}
			resp = retryResp
			bodyStr = string(resp.Body)
		}
	}

	if shouldFallbackToStripeCustomer(resp.StatusCode(), bodyStr) {
		if okResp := c.tryGetStripeCustomer(ctx, reqEditors...); okResp != nil {
			return &authapi.CreateTenantAndPricingResponse{
				Body:         []byte{},
				HTTPResponse: okResp,
			}, nil
		}
	}

	return resp, nil
}

func (c *AuthClientWithStripeFallback) tryGetStripeCustomer(ctx context.Context, reqEditors ...authapi.RequestEditorFn) *http.Response {
	tenantID := getLastTenantID()
	if tenantID == "" {
		return nil
	}
	resp, err := c.ClientWithResponses.GetStripeCustomer(ctx, tenantID, reqEditors...)
	if err != nil || resp == nil {
		return nil
	}
	defer func() {
		if resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()
	if !isSuccessStatus(resp.StatusCode) {
		return nil
	}
	return newEmptyOKResponse(resp.Header)
}

func updateStripeInfo(ctx context.Context) error {
	key := getStripeSecretKey()
	if key == "" {
		return fmt.Errorf("STRIPE_SECRET_KEY is not set")
	}

	client, err := billing.BillingWithResponse()
	if err != nil {
		return fmt.Errorf("failed to create billing client: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	payload := billingapi.UpdateStripeInfoParam{SecretKey: key}
	resp, err := client.UpdateStripeInfoWithResponse(ctx, payload)
	if err != nil {
		return fmt.Errorf("failed to update stripe info: %w", err)
	}

	statusCode := resp.StatusCode()
	bodyStr := string(resp.Body)
	if statusCode == http.StatusBadRequest {
		if isStripeAlreadyLinkedMessage(bodyStr) {
			return nil
		}
		return fmt.Errorf("update stripe info returned status 400: %s", bodyStr)
	}
	if statusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("update stripe info returned status %d: %s", statusCode, bodyStr)
	}
	return nil
}

func readAndRestoreBody(resp *http.Response) string {
	if resp == nil || resp.Body == nil {
		return ""
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	_ = resp.Body.Close()
	resp.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	return string(bodyBytes)
}

func extractMessage(bodyStr string) string {
	if bodyStr == "" {
		return ""
	}
	var payload struct {
		Message string `json:"message"`
	}
	if err := json.Unmarshal([]byte(bodyStr), &payload); err == nil {
		return payload.Message
	}
	return ""
}

func shouldRetryAfterStripeKeyError(statusCode int, bodyStr string) bool {
	if statusCode < http.StatusBadRequest {
		return false
	}
	msg := strings.ToLower(extractMessage(bodyStr))
	if msg == "" {
		return false
	}
	if strings.Contains(msg, "stripe key is not registered") {
		return true
	}
	return strings.Contains(msg, "stripe key") && strings.Contains(msg, "not registered")
}

func shouldFallbackToStripeCustomer(statusCode int, bodyStr string) bool {
	return statusCode >= http.StatusBadRequest
}

func isStripeAlreadyLinkedMessage(bodyStr string) bool {
	msg := strings.ToLower(extractMessage(bodyStr))
	return strings.Contains(msg, "already") || strings.Contains(msg, "exist") || strings.Contains(msg, "linked")
}

func isSuccessStatus(statusCode int) bool {
	return statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices
}

func newEmptyOKResponse(header http.Header) *http.Response {
	clone := make(http.Header, len(header))
	for k, v := range header {
		clone[k] = append([]string(nil), v...)
	}
	if clone.Get("Content-Length") == "" {
		clone.Set("Content-Length", "0")
	}
	return &http.Response{
		StatusCode:    http.StatusOK,
		Status:        "200 OK",
		Header:        clone,
		ContentLength: 0,
		Body:          io.NopCloser(bytes.NewReader(nil)),
	}
}
