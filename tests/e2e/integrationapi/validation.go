package integrationapi

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

	"github.com/saasus-platform/saasus-sdk-go/generated/integrationapi"
	"github.com/saasus-platform/saasus-sdk-go/modules/integration"
)

// getTestAWSAccountID returns the AWS Account ID configured for the test environment.
func getTestAWSAccountID() string {
	return getEnvWithDefault("TEST_AWS_ACCOUNT_ID", "")
}

// getTestAWSRegion returns test AWS Region
func getTestAWSRegion() string {
	if region := getEnvWithDefault("TEST_AWS_REGION", ""); region != "" {
		return region
	}
	return "ap-northeast-1"
}

// createSaveEventBridgeSettingsParams creates parameters for SaveEventBridgeSettings
func createSaveEventBridgeSettingsParams() integrationapi.SaveEventBridgeSettingsJSONRequestBody {
	return integrationapi.SaveEventBridgeSettingsParam{
		AwsAccountId: getTestAWSAccountID(),
		AwsRegion:    integrationapi.ApNortheast1,
	}
}

// createSaveEventBridgeSettingsBodyParams creates parameters for SaveEventBridgeSettingsWithBody
func createSaveEventBridgeSettingsBodyParams() struct {
	ContentType string
	Body        io.Reader
} {
	param := integrationapi.SaveEventBridgeSettingsParam{
		AwsAccountId: getTestAWSAccountID(),
		AwsRegion:    integrationapi.ApNortheast1,
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

// createEventBridgeEventParams creates parameters for CreateEventBridgeEvent
func createEventBridgeEventParams() integrationapi.CreateEventBridgeEventJSONRequestBody {
	return integrationapi.CreateEventBridgeEventParam{
		EventMessages: []integrationapi.EventMessage{
			{
				EventType:       "api_call",
				EventDetailType: "create_user",
				Message:         "{id:8b79528a-ec3b-4f68-b7c4-d793e3894561,name:test222}",
			},
		},
	}
}

// createEventBridgeEventBodyParams creates parameters for CreateEventBridgeEventWithBody
func createEventBridgeEventBodyParams() struct {
	ContentType string
	Body        io.Reader
} {
	param := integrationapi.CreateEventBridgeEventParam{
		EventMessages: []integrationapi.EventMessage{
			{
				EventType:       "api_call",
				EventDetailType: "create_user",
				Message:         "{id:8b79528a-ec3b-4f68-b7c4-d793e3894561,name:test222}",
			},
		},
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

// validateEventBridgeSettingsResponse validates GetEventBridgeSettings response
func validateEventBridgeSettingsResponse(response interface{}, expectSettings bool) error {
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

// validateEventBridgeSettingsWithResponse validates GetEventBridgeSettingsWithResponse response
func validateEventBridgeSettingsWithResponse(response interface{}, expectSettings bool) error {
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

// cleanupEventBridgeSettings cleans up any existing EventBridge settings
func cleanupEventBridgeSettings() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create integration client
	client, err := integration.IntegrationWithResponse()
	if err != nil {
		// If client creation fails, just log and continue
		// This is cleanup, so we don't want to fail the test
		return nil
	}

	// Try to delete any existing EventBridge settings
	// Ignore errors as the settings might not exist
	_, _ = client.DeleteEventBridgeSettings(ctx)

	return nil
}

// getEnvWithDefault returns environment variable or default value
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
