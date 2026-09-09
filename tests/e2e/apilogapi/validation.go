// Package apilogapi provides validation and parameter generation functions for Apilog API E2E tests.
package apilogapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"time"

	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/saasus-platform/saasus-sdk-go/generated/apilogapi"
)

// getLogsParams creates parameters for GetLogs (no query parameters)
func getLogsParams(variables map[string]any) any {
	// Return empty params struct (not nil) for GetLogs
	return &apilogapi.GetLogsParams{}
}

// getLogsWithQueryParams creates parameters for GetLogs with query parameters
func getLogsWithQueryParams(variables map[string]any) any {
	params := &apilogapi.GetLogsParams{}

	// CreatedDate expects *openapi_types.Date
	if createdDateStr, ok := variables["created_date"].(string); ok && createdDateStr != "" {
		date, err := time.Parse("2006-01-02", createdDateStr)
		if err != nil {
			// Log parse error but continue - this is a test helper
			fmt.Printf("Warning: failed to parse created_date '%s': %v\n", createdDateStr, err)
		} else {
			openapiDate := openapi_types.Date{Time: date}
			params.CreatedDate = &openapiDate
		}
	}

	// CreatedAt expects *time.Time
	if createdAtStr, ok := variables["created_at"].(string); ok && createdAtStr != "" {
		t, err := time.Parse(time.RFC3339, createdAtStr)
		if err != nil {
			// Log parse error but continue - this is a test helper
			fmt.Printf("Warning: failed to parse created_at '%s': %v\n", createdAtStr, err)
		} else {
			params.CreatedAt = &t
		}
	}

	if cursor, ok := variables["cursor"].(string); ok && cursor != "" {
		params.Cursor = &cursor
	}

	return params
}

// getLogParams creates parameters for GetLog
// Returns empty string if api_log_id is not found or not a string
func getLogParams(variables map[string]any) any {
	apiLogID, ok := variables["api_log_id"].(string)
	if !ok || apiLogID == "" {
		// Return empty string as fallback - testlib will handle validation
		return ""
	}
	return apiLogID
}

const (
	httpStatusOK              = 200
	httpStatusMultipleChoices = 300
	httpStatusUnauthorized    = 401
)

// validateHTTPStatusCode validates HTTP status code from standard response
func validateHTTPStatusCode(statusCode int) error {
	if statusCode == httpStatusUnauthorized {
		return fmt.Errorf("authentication failed: received 401 Unauthorized - check SAASUS_SECRET_KEY, SAASUS_API_KEY, and SAASUS_SAAS_ID")
	}
	if statusCode < httpStatusOK || statusCode >= httpStatusMultipleChoices {
		return fmt.Errorf("unexpected HTTP status: %d", statusCode)
	}
	return nil
}

// validateStandardResponse validates standard method response (*http.Response)
func validateStandardResponse(response any) error {
	if httpResp, ok := response.(*http.Response); ok {
		return validateHTTPStatusCode(httpResp.StatusCode)
	}
	return nil
}

// validateWithResponseMethod validates WithResponse method response using reflection
func validateWithResponseMethod(response any) error {
	respValue := reflect.ValueOf(response)
	if respValue.Kind() == reflect.Ptr && !respValue.IsNil() {
		statusCodeMethod := respValue.MethodByName("StatusCode")
		if statusCodeMethod.IsValid() {
			results := statusCodeMethod.Call([]reflect.Value{})
			if len(results) > 0 {
				if statusCode, ok := results[0].Interface().(int); ok {
					return validateHTTPStatusCode(statusCode)
				}
			}
		}
	}
	return nil
}

// validateGetLogsResponse validates GetLogs response (Standard Method)
func validateGetLogsResponse(response any) error {
	return validateStandardResponse(response)
}

// validateGetLogsWithResponse validates GetLogsWithResponse response (WithResponse Method)
func validateGetLogsWithResponse(response any) error {
	return validateWithResponseMethod(response)
}

// validateGetLogResponse validates GetLog response (Standard Method)
func validateGetLogResponse(response any) error {
	return validateStandardResponse(response)
}

// validateGetLogWithResponse validates GetLogWithResponse response (WithResponse Method)
func validateGetLogWithResponse(response any) error {
	return validateWithResponseMethod(response)
}

// convertUnixToISO8601 converts UNIX timestamp to ISO 8601 format
func convertUnixToISO8601(timestamp int64) string {
	t := time.Unix(timestamp, 0)
	return t.Format(time.RFC3339)
}

// extractLogsVariables extracts variables from GetLogs response (Standard Method)
func extractLogsVariables(response any, variables map[string]any) error {
	// Handle case where snapshot framework extracts single JSON field value (api_logs array)
	// This happens when returnValue.JSONData has only one key
	if apiLogsArray, ok := response.([]any); ok {
		// Reconstruct the expected JSON structure
		jsonData := map[string]any{
			"api_logs": apiLogsArray,
		}
		return extractLogsVariablesFromJSON(jsonData, variables)
	}

	// Try to handle captured JSON data from snapshot framework
	if jsonData, ok := response.(map[string]any); ok {
		// Snapshot framework provides parsed JSON data
		return extractLogsVariablesFromJSON(jsonData, variables)
	}

	// For Standard Methods, response is *http.Response, need to decode body
	if httpResp, ok := response.(*http.Response); ok {
		if httpResp.Body == nil {
			return fmt.Errorf("response body is nil")
		}

		// Read the entire body into memory first
		bodyBytes, err := io.ReadAll(httpResp.Body)
		if err != nil {
			_ = httpResp.Body.Close()
			return fmt.Errorf("failed to read response body: %w", err)
		}

		// Close the original body
		_ = httpResp.Body.Close()

		// Replace body with a new reader so it can be read again by snapshot framework
		// This is critical - without this, the snapshot framework cannot read the body
		httpResp.Body = io.NopCloser(bytes.NewReader(bodyBytes))

		// Decode the body bytes into ApiLogs struct
		var apiLogs apilogapi.ApiLogs
		if err := json.Unmarshal(bodyBytes, &apiLogs); err != nil {
			return fmt.Errorf("failed to decode response body: %w", err)
		}

		// Check if there are any logs in the response
		if len(apiLogs.ApiLogs) == 0 {
			// No logs available, but this is not an error - just no data to extract
			return nil
		}

		// Extract variables from the decoded response
		return extractLogsVariablesFromStruct(&apiLogs, variables)
	}
	return nil
}

// extractLogsVariablesFromJSON extracts variables from JSON map (used by snapshot framework)
func extractLogsVariablesFromJSON(jsonData map[string]any, variables map[string]any) error {
	// Extract cursor if present
	if cursor, ok := jsonData["cursor"].(string); ok {
		variables["cursor"] = cursor
	}

	// Extract first log entry if present
	if apiLogsArray, ok := jsonData["api_logs"].([]any); ok && len(apiLogsArray) > 0 {
		if firstLog, ok := apiLogsArray[0].(map[string]any); ok {
			if apiLogID, ok := firstLog["api_log_id"].(string); ok {
				variables["api_log_id"] = apiLogID
			}
			if createdDate, ok := firstLog["created_date"].(string); ok {
				variables["created_date"] = createdDate
			}
			if createdAt, ok := firstLog["created_at"].(float64); ok {
				timestamp := int64(createdAt)
				variables["timestamp"] = timestamp
				variables["created_at"] = convertUnixToISO8601(timestamp)
			}
		}
	}

	return nil
}

// extractLogsVariablesFromStruct extracts variables from ApiLogs struct
func extractLogsVariablesFromStruct(apiLogs *apilogapi.ApiLogs, variables map[string]any) error {
	if apiLogs.Cursor != nil {
		variables["cursor"] = *apiLogs.Cursor
	}
	if len(apiLogs.ApiLogs) > 0 {
		firstLog := apiLogs.ApiLogs[0]
		variables["api_log_id"] = firstLog.ApiLogId
		variables["created_date"] = firstLog.CreatedDate
		variables["timestamp"] = firstLog.CreatedAt
		variables["created_at"] = convertUnixToISO8601(firstLog.CreatedAt)
	}
	return nil
}

// extractLogsVariablesWithResponse extracts variables from GetLogsWithResponse response (WithResponse Method)
func extractLogsVariablesWithResponse(response any, variables map[string]any) error {
	// Handle case where snapshot framework extracts single JSON field value (api_logs array)
	// This happens when returnValue.JSONData has only one key
	if apiLogsArray, ok := response.([]any); ok {
		// Reconstruct the expected JSON structure
		jsonData := map[string]any{
			"api_logs": apiLogsArray,
		}
		return extractLogsVariablesFromJSON(jsonData, variables)
	}

	// Try to handle captured JSON data from snapshot framework
	if jsonData, ok := response.(map[string]any); ok {
		// Snapshot framework provides parsed JSON data
		return extractLogsVariablesFromJSON(jsonData, variables)
	}

	// For WithResponse Methods, we need to extract the body using reflection
	respValue := reflect.ValueOf(response)
	if respValue.Kind() == reflect.Ptr && !respValue.IsNil() {
		// Try to get JSON200 field (successful response)
		json200Field := respValue.Elem().FieldByName("JSON200")
		if json200Field.IsValid() && !json200Field.IsNil() {
			if apiLogs, ok := json200Field.Interface().(*apilogapi.ApiLogs); ok {
				return extractLogsVariablesFromStruct(apiLogs, variables)
			}
		}
	}
	return nil
}
