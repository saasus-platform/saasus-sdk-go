package snapshot

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// SDKMethodDifference represents a difference in SDK method snapshots
type SDKMethodDifference struct {
	Type        string      `json:"type"`
	MethodName  string      `json:"method_name"`
	Field       string      `json:"field,omitempty"`
	Description string      `json:"description"`
	OldValue    interface{} `json:"old_value,omitempty"`
	NewValue    interface{} `json:"new_value,omitempty"`
}

// compareSDKReturnValues compares SDK return values
func compareSDKReturnValues(methodName string, oldValue, newValue *SDKReturnValue) []*SDKMethodDifference {
	var differences []*SDKMethodDifference

	if oldValue == nil && newValue == nil {
		return differences
	}

	if oldValue == nil {
		differences = append(differences, &SDKMethodDifference{
			Type:        "RETURN_VALUE",
			MethodName:  methodName,
			Description: "Return value added",
			NewValue:    newValue,
		})
		return differences
	}

	if newValue == nil {
		differences = append(differences, &SDKMethodDifference{
			Type:        "RETURN_VALUE",
			MethodName:  methodName,
			Description: "Return value removed",
			OldValue:    oldValue,
		})
		return differences
	}

	// Compare status code
	if oldValue.StatusCode != newValue.StatusCode {
		differences = append(differences, &SDKMethodDifference{
			Type:        "STATUS_CODE",
			MethodName:  methodName,
			Field:       "StatusCode",
			Description: fmt.Sprintf("Status code changed from %d to %d", oldValue.StatusCode, newValue.StatusCode),
			OldValue:    oldValue.StatusCode,
			NewValue:    newValue.StatusCode,
		})
	}

	// Compare type
	if oldValue.Type != newValue.Type {
		differences = append(differences, &SDKMethodDifference{
			Type:        "TYPE",
			MethodName:  methodName,
			Field:       "Type",
			Description: fmt.Sprintf("Return type changed from %s to %s", oldValue.Type, newValue.Type),
			OldValue:    oldValue.Type,
			NewValue:    newValue.Type,
		})
	}

	// Compare JSON data (normalized)
	if !compareNormalizedJSONData(oldValue.JSONData, newValue.JSONData) {
		differences = append(differences, &SDKMethodDifference{
			Type:        "JSON_DATA",
			MethodName:  methodName,
			Field:       "JSONData",
			Description: "JSON response data changed",
			OldValue:    oldValue.JSONData,
			NewValue:    newValue.JSONData,
		})
	}

	// Compare body (normalized)
	if !compareNormalizedBody(oldValue.Body, newValue.Body) {
		differences = append(differences, &SDKMethodDifference{
			Type:        "BODY",
			MethodName:  methodName,
			Field:       "Body",
			Description: "Response body changed",
			OldValue:    oldValue.Body,
			NewValue:    newValue.Body,
		})
	}

	// Compare headers (excluding dynamic ones)
	if !compareNormalizedHeaders(oldValue.Headers, newValue.Headers) {
		differences = append(differences, &SDKMethodDifference{
			Type:        "HEADERS",
			MethodName:  methodName,
			Field:       "Headers",
			Description: "Response headers changed",
			OldValue:    oldValue.Headers,
			NewValue:    newValue.Headers,
		})
	}

	return differences
}

// compareNormalizedJSONData compares JSON data while ignoring dynamic fields
func compareNormalizedJSONData(oldData, newData interface{}) bool {
	// If both are nil, they're equal
	if oldData == nil && newData == nil {
		return true
	}

	// If one is nil and the other isn't, they're different
	if oldData == nil || newData == nil {
		return false
	}

	// Convert to normalized JSON strings for comparison
	oldNormalized := normalizeJSONForComparison(oldData)
	newNormalized := normalizeJSONForComparison(newData)

	return oldNormalized == newNormalized
}

// compareNormalizedBody compares response bodies while ignoring dynamic content
func compareNormalizedBody(oldBody, newBody string) bool {
	// Parse and normalize JSON bodies
	oldNormalized := normalizeJSONStringForComparison(oldBody)
	newNormalized := normalizeJSONStringForComparison(newBody)

	return oldNormalized == newNormalized
}

// compareNormalizedHeaders compares headers while ignoring dynamic ones
func compareNormalizedHeaders(oldHeaders, newHeaders map[string]string) bool {
	// Dynamic headers to ignore
	dynamicHeaders := map[string]bool{
		"Date":              true,
		"X-Saasus-Trace-Id": true,
		"X-Request-Id":      true,
		"X-Correlation-Id":  true,
		"Server":            true,
		"X-Runtime":         true,
	}

	// Create normalized header maps
	oldNormalized := make(map[string]string)
	newNormalized := make(map[string]string)

	for key, value := range oldHeaders {
		if !dynamicHeaders[key] {
			oldNormalized[key] = value
		}
	}

	for key, value := range newHeaders {
		if !dynamicHeaders[key] {
			newNormalized[key] = value
		}
	}

	return reflect.DeepEqual(oldNormalized, newNormalized)
}

// normalizeJSONForComparison normalizes JSON data for comparison
func normalizeJSONForComparison(data interface{}) string {
	if data == nil {
		return ""
	}

	// Convert to JSON bytes
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}

	// Parse back to normalize structure
	var normalized interface{}
	if err := json.Unmarshal(jsonBytes, &normalized); err != nil {
		return string(jsonBytes)
	}

	// Convert back to string with consistent formatting
	normalizedBytes, err := json.Marshal(normalized)
	if err != nil {
		return string(jsonBytes)
	}

	return string(normalizedBytes)
}

// normalizeJSONStringForComparison normalizes JSON string for comparison
func normalizeJSONStringForComparison(jsonStr string) string {
	if jsonStr == "" {
		return ""
	}

	// Try to parse as JSON
	var data interface{}
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		// If not valid JSON, return as-is
		return jsonStr
	}

	// Re-marshal with consistent formatting
	normalized, err := json.Marshal(data)
	if err != nil {
		return jsonStr
	}

	return string(normalized)
}
