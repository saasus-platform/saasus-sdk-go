package billingapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"time"

	"github.com/saasus-platform/saasus-sdk-go/generated/billingapi"
)

// getTestStripeKey returns test Stripe key
func getTestStripeKey() string {
	if key := getEnvWithDefault("STRIPE_SECRET_KEY", ""); key != "" {
		return key
	}
	return "sk_test_example_key_for_testing"
}

// createUpdateStripeInfoParams creates parameters for UpdateStripeInfo
func createUpdateStripeInfoParams() billingapi.UpdateStripeInfoJSONRequestBody {
	return billingapi.UpdateStripeInfoParam{
		SecretKey: getTestStripeKey(),
	}
}

// createUpdateStripeInfoBodyParams creates parameters for UpdateStripeInfoWithBody
func createUpdateStripeInfoBodyParams() struct {
	ContentType string
	Body        io.Reader
} {
	param := billingapi.UpdateStripeInfoParam{
		SecretKey: getTestStripeKey(),
	}
	body, _ := json.Marshal(param)
	return struct {
		ContentType string
		Body        io.Reader
	}{
		ContentType: "application/json",
		Body:        bytes.NewReader(body),
	}
}

// validateStripeInfoResponse validates GetStripeInfo response
func validateStripeInfoResponse(response interface{}, expectedRegistered bool) error {
	// Check if response is an HTTP response with status code
	if httpResp, ok := response.(*http.Response); ok {
		if httpResp.StatusCode == 401 {
			return fmt.Errorf("authentication failed: received 401 Unauthorized - check SAASUS_SECRET_KEY, SAASUS_API_KEY, and SAASUS_SAAS_ID")
		}
		if httpResp.StatusCode < 200 || httpResp.StatusCode >= 300 {
			return fmt.Errorf("unexpected HTTP status: %d", httpResp.StatusCode)
		}
	}
	return nil
}

// validateStripeInfoWithResponse validates GetStripeInfoWithResponse response
func validateStripeInfoWithResponse(response interface{}, expectedRegistered bool) error {
	// Check if response has StatusCode method (WithResponse types)
	respValue := reflect.ValueOf(response)
	if respValue.Kind() == reflect.Ptr && !respValue.IsNil() {
		statusCodeMethod := respValue.MethodByName("StatusCode")
		if statusCodeMethod.IsValid() {
			results := statusCodeMethod.Call([]reflect.Value{})
			if len(results) > 0 {
				if statusCode, ok := results[0].Interface().(int); ok {
					if statusCode == 401 {
						return fmt.Errorf("authentication failed: received 401 Unauthorized - check SAASUS_SECRET_KEY, SAASUS_API_KEY, and SAASUS_SAAS_ID")
					}
					if statusCode < 200 || statusCode >= 300 {
						return fmt.Errorf("unexpected HTTP status: %d", statusCode)
					}
				}
			}
		}
	}
	return nil
}

// validateUpdateStripeInfoWithResponse validates update response
func validateUpdateStripeInfoWithResponse(response interface{}) error {
	// Check if response has StatusCode method (WithResponse types)
	respValue := reflect.ValueOf(response)
	if respValue.Kind() == reflect.Ptr && !respValue.IsNil() {
		statusCodeMethod := respValue.MethodByName("StatusCode")
		if statusCodeMethod.IsValid() {
			results := statusCodeMethod.Call([]reflect.Value{})
			if len(results) > 0 {
				if statusCode, ok := results[0].Interface().(int); ok {
					if statusCode == 401 {
						return fmt.Errorf("authentication failed: received 401 Unauthorized - check SAASUS_SECRET_KEY, SAASUS_API_KEY, and SAASUS_SAAS_ID")
					}
					if statusCode < 200 || statusCode >= 300 {
						return fmt.Errorf("unexpected HTTP status: %d", statusCode)
					}
				}
			}
		}
	}
	return nil
}

// validateDeleteStripeInfoWithResponse validates delete response
func validateDeleteStripeInfoWithResponse(response interface{}) error {
	// Check if response has StatusCode method (WithResponse types)
	respValue := reflect.ValueOf(response)
	if respValue.Kind() == reflect.Ptr && !respValue.IsNil() {
		statusCodeMethod := respValue.MethodByName("StatusCode")
		if statusCodeMethod.IsValid() {
			results := statusCodeMethod.Call([]reflect.Value{})
			if len(results) > 0 {
				if statusCode, ok := results[0].Interface().(int); ok {
					if statusCode == 401 {
						return fmt.Errorf("authentication failed: received 401 Unauthorized - check SAASUS_SECRET_KEY, SAASUS_API_KEY, and SAASUS_SAAS_ID")
					}
					if statusCode < 200 || statusCode >= 300 {
						return fmt.Errorf("unexpected HTTP status: %d", statusCode)
					}
				}
			}
		}
	}
	return nil
}

// cleanupStripeInfo cleans up any existing Stripe info
func cleanupStripeInfo() error {
	// Implementation will clean up any existing Stripe info
	// This is a placeholder for the actual cleanup logic
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Try to delete any existing Stripe info
	// This will be implemented in the actual test execution engine
	_ = ctx
	return nil
}

// getEnvWithDefault returns environment variable or default value
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
