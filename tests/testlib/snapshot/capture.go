package snapshot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"
)

// StorySnapshotCapture handles capturing story snapshots during E2E test execution
type StorySnapshotCapture struct {
	config *StorySnapshotConfig
	logger *SnapshotLogger
}

// NewStorySnapshotCapture creates a new story snapshot capture instance
func NewStorySnapshotCapture(config *StorySnapshotConfig) *StorySnapshotCapture {
	if config == nil {
		config = DefaultStorySnapshotConfig()
	}

	// Create logger with appropriate level
	logLevel := LogLevelInfo
	debugMode := false
	if config.CaptureLevel == CaptureLevelFull {
		logLevel = LogLevelDebug
		debugMode = true
	}

	logger := NewSnapshotLogger(logLevel, debugMode)

	return &StorySnapshotCapture{
		config: config,
		logger: logger,
	}
}

// CaptureStorySnapshot captures a complete story execution as a snapshot
func (ssc *StorySnapshotCapture) CaptureStorySnapshot(story Story, storyResult StoryResult) (*StorySnapshot, error) {
	if !ssc.config.EnableCapture {
		ssc.logger.LogDebug("Snapshot capture is disabled, skipping story: %s", story.Name)
		return nil, nil
	}

	ssc.logger.LogPhaseStart(PhaseCapture, fmt.Sprintf("story '%s'", story.Name))
	ssc.logger.SetContext("story", story.Name)
	defer ssc.logger.ClearContext()

	snapshot := &StorySnapshot{
		StoryName:   story.Name,
		Description: story.Description,
		Timestamp:   time.Now(),
		Duration:    storyResult.Duration,
		Status:      storyResult.Status,
		Variables:   make(map[string]interface{}),
		Steps:       []StepSnapshot{},
		Metadata: SnapshotMetadata{
			SDKVersion:      getSDKVersionFromEnv(),
			TestEnvironment: getTestEnvironment(),
			CaptureLevel:    ssc.config.CaptureLevel,
			GitTag:          getGitTagSafe(),
			GitCommit:       getGitCommitSafe(),
		},
	}

	// Copy story variables (mask sensitive data)
	for k, v := range storyResult.Variables {
		snapshot.Variables[k] = ssc.maskSensitiveData(k, v)
	}

	// Capture each step
	for i, stepResult := range storyResult.Steps {
		var step Step
		if i < len(story.Steps) {
			step = story.Steps[i]
		}

		ssc.logger.LogDebug("Capturing step %d: %s", i+1, stepResult.StepName)
		stepSnapshot, err := ssc.captureStepSnapshot(step, stepResult)
		if err != nil {
			snapErr := NewCaptureError(ErrorTypeParsing,
				fmt.Sprintf("failed to capture step %s", stepResult.StepName),
				err, story.Name, stepResult.StepName)
			ssc.logger.LogSnapshotError(snapErr)
			return nil, snapErr
		}

		snapshot.Steps = append(snapshot.Steps, *stepSnapshot)
	}

	// Generate execution summary
	snapshot.Summary = ssc.generateExecutionSummary(snapshot.Steps)

	ssc.logger.LogInfo("Successfully captured story snapshot: %d steps, %v duration",
		len(snapshot.Steps), snapshot.Duration)

	return snapshot, nil
}

// captureStepSnapshot captures a single step execution as a snapshot
func (ssc *StorySnapshotCapture) captureStepSnapshot(step Step, stepResult StepResult) (*StepSnapshot, error) {
	stepSnapshot := &StepSnapshot{
		StepName:     stepResult.StepName,
		Method:       stepResult.Method,
		Parameters:   make(map[string]interface{}),
		Duration:     stepResult.Duration,
		StatusCode:   stepResult.StatusCode,
		Success:      stepResult.Status == TestStatusPassed,
		Status:       stepResult.Status,
		SkipReason:   stepResult.SkipReason,
		Timestamp:    time.Now(),
		StateChanges: make(map[string]interface{}),
	}

	// Capture parameters (mask sensitive data)
	// Skip function parameters as they cannot be serialized to JSON
	if step.Parameters != nil {
		paramValue := reflect.ValueOf(step.Parameters)
		if paramValue.Kind() != reflect.Func {
			stepSnapshot.Parameters = ssc.maskSensitiveParameters(step.Parameters)
		}
	}

	// Capture error if present
	if stepResult.Error != nil {
		stepSnapshot.Error = &SDKMethodError{
			Type:    "StepExecutionError",
			Message: stepResult.Error.Error(),
		}
	}

	// Use ReturnValue from StepResult if available
	stepSnapshot.ReturnValue = stepResult.ReturnValue

	// Ensure ReturnValue.StatusCode matches StepSnapshot.StatusCode
	// This handles cases where CaptureStepResponse couldn't extract the StatusCode
	if stepSnapshot.ReturnValue != nil && stepSnapshot.ReturnValue.StatusCode == 0 && stepSnapshot.StatusCode != 0 {
		stepSnapshot.ReturnValue.StatusCode = stepSnapshot.StatusCode
	}

	return stepSnapshot, nil
}

// CaptureStepResponse captures the response data from a step execution
func (ssc *StorySnapshotCapture) CaptureStepResponse(stepName string, response interface{}) (*SDKReturnValue, error) {
	if response == nil {
		return nil, nil
	}

	returnValue := &SDKReturnValue{
		Type:     reflect.TypeOf(response).String(),
		JSONData: make(map[string]interface{}),
		Headers:  make(map[string]string),
	}

	// Use reflection to extract response fields
	respValue := reflect.ValueOf(response)

	// Extract StatusCode using method call
	if statusCodeMethod := respValue.MethodByName("StatusCode"); statusCodeMethod.IsValid() {
		if result := statusCodeMethod.Call(nil); len(result) > 0 {
			if statusCode, ok := result[0].Interface().(int); ok {
				returnValue.StatusCode = statusCode
			}
		}
	}

	// Extract Status using method call
	if statusMethod := respValue.MethodByName("Status"); statusMethod.IsValid() {
		if result := statusMethod.Call(nil); len(result) > 0 {
			if status, ok := result[0].Interface().(string); ok {
				returnValue.Status = status
			}
		}
	}

	// Dereference pointer for field access
	if respValue.Kind() == reflect.Ptr {
		respValue = respValue.Elem()
	}

	// Extract Body
	if bodyField := respValue.FieldByName("Body"); bodyField.IsValid() && bodyField.CanInterface() {
		if bodyBytes, ok := bodyField.Interface().([]byte); ok {
			returnValue.Body = string(bodyBytes)
		}
	}

	// Handle *http.Response specially to parse JSON body
	if httpResp, ok := response.(*http.Response); ok {
		httpSnapshot, err := ssc.captureHTTPResponse(httpResp)
		if err == nil {
			returnValue.HTTPResponse = httpSnapshot
			// Copy headers to top level
			for k, v := range httpSnapshot.Headers {
				returnValue.Headers[k] = v
			}
			returnValue.StatusCode = httpSnapshot.StatusCode
			returnValue.Status = httpSnapshot.Status
		}

		// Parse JSON body if available
		if httpResp.Body != nil {
			bodyBytes, err := io.ReadAll(httpResp.Body)
			if err == nil {
				_ = httpResp.Body.Close()

				// Replace body with a new reader so it can be read again by StateUpdate
				httpResp.Body = io.NopCloser(bytes.NewReader(bodyBytes))

				// Try to parse as JSON
				var jsonData map[string]interface{}
				if err := json.Unmarshal(bodyBytes, &jsonData); err == nil {
					if masked := ssc.maskSensitiveInterface(jsonData); masked != nil {
						if maskedMap, ok := masked.(map[string]interface{}); ok {
							returnValue.JSONData = maskedMap
						}
					}
					if ssc.logger != nil {
						ssc.logger.LogDebug(fmt.Sprintf("Parsed HTTP response body as JSON: %d keys", len(jsonData)))
					}
				}
				returnValue.Body = ssc.maskBodyString(string(bodyBytes))
			}
		}
	} else {
		// Extract HTTPResponse or create from response itself
		if httpRespField := respValue.FieldByName("HTTPResponse"); httpRespField.IsValid() && httpRespField.CanInterface() {
			if httpResp := httpRespField.Interface(); httpResp != nil {
				httpSnapshot, err := ssc.captureHTTPResponse(httpResp)
				if err == nil {
					returnValue.HTTPResponse = httpSnapshot
					// Copy headers to top level
					for k, v := range httpSnapshot.Headers {
						returnValue.Headers[k] = v
					}
					// If StatusCode wasn't captured from method, use HTTPResponse.StatusCode
					if returnValue.StatusCode == 0 && httpSnapshot.StatusCode != 0 {
						returnValue.StatusCode = httpSnapshot.StatusCode
						returnValue.Status = httpSnapshot.Status
					}
				}
			}
		} else {
			// Try to create HTTPResponse from the response itself
			httpSnapshot, err := ssc.captureHTTPResponse(response)
			if err == nil && httpSnapshot != nil {
				returnValue.HTTPResponse = httpSnapshot
				// Copy headers to top level
				for k, v := range httpSnapshot.Headers {
					returnValue.Headers[k] = v
				}
				// If StatusCode wasn't captured from method, use HTTPResponse.StatusCode
				if returnValue.StatusCode == 0 && httpSnapshot.StatusCode != 0 {
					returnValue.StatusCode = httpSnapshot.StatusCode
					returnValue.Status = httpSnapshot.Status
				}
			}
		}
	}

	// Extract JSON response fields (JSON200, JSON400, etc.)
	respType := respValue.Type()
	for i := 0; i < respValue.NumField(); i++ {
		field := respValue.Field(i)
		fieldType := respType.Field(i)

		// Look for JSON response fields
		if len(fieldType.Name) >= 4 && fieldType.Name[:4] == "JSON" && field.CanInterface() && !field.IsNil() {
			jsonData, err := json.Marshal(field.Interface())
			if err == nil {
				var jsonObj interface{}
				if err := json.Unmarshal(jsonData, &jsonObj); err == nil {
					returnValue.JSONData[fieldType.Name] = ssc.maskSensitiveInterface(jsonObj)
				}
			}
		}
	}

	returnValue.Body = ssc.maskBodyString(returnValue.Body)

	return returnValue, nil
}

// captureHTTPResponse captures HTTP response details
func (ssc *StorySnapshotCapture) captureHTTPResponse(httpResp interface{}) (*HTTPResponseSnapshot, error) {
	httpSnapshot := &HTTPResponseSnapshot{
		Headers: make(map[string]string),
	}

	// Handle *http.Response directly
	if httpResponse, ok := httpResp.(*http.Response); ok {
		httpSnapshot.StatusCode = httpResponse.StatusCode
		httpSnapshot.Status = httpResponse.Status
		httpSnapshot.ContentLength = httpResponse.ContentLength

		// Extract headers
		for key, values := range httpResponse.Header {
			if len(values) > 0 {
				httpSnapshot.Headers[key] = ssc.maskHeaderValue(key, values[0])
				if key == "X-Saasus-Trace-Id" {
					httpSnapshot.TraceID = values[0]
				}
			}
		}

		return httpSnapshot, nil
	}

	// Use reflection to extract HTTP response fields
	respValue := reflect.ValueOf(httpResp)
	if respValue.Kind() == reflect.Ptr {
		respValue = respValue.Elem()
	}

	// Extract StatusCode - try field first, then method
	if statusCodeField := respValue.FieldByName("StatusCode"); statusCodeField.IsValid() && statusCodeField.CanInterface() {
		if statusCode, ok := statusCodeField.Interface().(int); ok {
			httpSnapshot.StatusCode = statusCode
		}
	} else {
		// Try method call
		respValuePtr := reflect.ValueOf(httpResp)
		if statusCodeMethod := respValuePtr.MethodByName("StatusCode"); statusCodeMethod.IsValid() {
			if result := statusCodeMethod.Call(nil); len(result) > 0 {
				if statusCode, ok := result[0].Interface().(int); ok {
					httpSnapshot.StatusCode = statusCode
				}
			}
		}
	}

	// Extract Status - try field first, then method
	if statusField := respValue.FieldByName("Status"); statusField.IsValid() && statusField.CanInterface() {
		if status, ok := statusField.Interface().(string); ok {
			httpSnapshot.Status = status
		}
	} else {
		// Try method call
		respValuePtr := reflect.ValueOf(httpResp)
		if statusMethod := respValuePtr.MethodByName("Status"); statusMethod.IsValid() {
			if result := statusMethod.Call(nil); len(result) > 0 {
				if status, ok := result[0].Interface().(string); ok {
					httpSnapshot.Status = status
				}
			}
		}
	}

	// Extract ContentLength
	if contentLengthField := respValue.FieldByName("ContentLength"); contentLengthField.IsValid() && contentLengthField.CanInterface() {
		if contentLength, ok := contentLengthField.Interface().(int64); ok {
			httpSnapshot.ContentLength = contentLength
		}
	}

	// Extract Headers
	if headerField := respValue.FieldByName("Header"); headerField.IsValid() && headerField.CanInterface() {
		if headers := headerField.Interface(); headers != nil {
			headerValue := reflect.ValueOf(headers)
			if headerValue.Kind() == reflect.Map {
				for _, key := range headerValue.MapKeys() {
					keyStr := key.String()
					values := headerValue.MapIndex(key)
					if values.IsValid() && values.Len() > 0 {
						firstValue := values.Index(0)
						if firstValue.IsValid() {
							valueStr := firstValue.String()
							httpSnapshot.Headers[keyStr] = ssc.maskHeaderValue(keyStr, valueStr)

							// Extract trace ID
							if keyStr == "X-Saasus-Trace-Id" {
								httpSnapshot.TraceID = valueStr
							}
						}
					}
				}
			}
		}
	}

	return httpSnapshot, nil
}

// TrackStateChange tracks a state change during story execution
func (ssc *StorySnapshotCapture) TrackStateChange(stepSnapshot *StepSnapshot, key string, oldValue, newValue interface{}) {
	if stepSnapshot.StateChanges == nil {
		stepSnapshot.StateChanges = make(map[string]interface{})
	}

	stepSnapshot.StateChanges[key] = map[string]interface{}{
		"old_value": ssc.maskSensitiveData(key, oldValue),
		"new_value": ssc.maskSensitiveData(key, newValue),
		"timestamp": time.Now(),
	}
}

// generateExecutionSummary generates execution summary for a story
func (ssc *StorySnapshotCapture) generateExecutionSummary(steps []StepSnapshot) StoryExecutionSummary {
	summary := StoryExecutionSummary{
		TotalSteps: len(steps),
	}

	var totalDuration time.Duration
	for _, step := range steps {
		totalDuration += step.Duration

		if step.Status == TestStatusSkipped {
			summary.SkippedSteps++
			continue
		}
		if step.Success {
			summary.SuccessfulSteps++
		} else {
			summary.FailedSteps++
		}
	}

	summary.TotalDuration = totalDuration
	if summary.TotalSteps > 0 {
		summary.AverageStepDuration = totalDuration / time.Duration(summary.TotalSteps)
	}

	return summary
}

// maskSensitiveData masks sensitive data in variables and parameters
var sensitiveKeySubstrings = []string{
	"secret",
	"password",
	"token",
	"credential",
	"api_key",
	"api-key",
	"secret_key",
	"secret-key",
	"saasus_secret",
	"saasus_api_key",
	"stripe_key",
	"stripe_secret",
	"access_key",
	"access-key",
	"client_secret",
	"client-secret",
	"refresh_token",
	"refresh-token",
	"set-cookie",
	"cookie",
	"authorization",
	"proxy-authorization",
	"x-api-key",
}

func (ssc *StorySnapshotCapture) maskSensitiveData(key string, value interface{}) interface{} {
	keyLower := strings.ToLower(fmt.Sprintf("%v", key))

	if value == nil {
		return nil
	}

	if isSensitiveKey(keyLower) {
		return maskValue(value)
	}

	switch v := value.(type) {
	case map[string]interface{}:
		return ssc.maskSensitiveInterface(v)
	case []interface{}:
		masked := make([]interface{}, len(v))
		for i, item := range v {
			masked[i] = ssc.maskSensitiveData(key, item)
		}
		return masked
	case map[string]string:
		masked := make(map[string]interface{}, len(v))
		for k, item := range v {
			masked[k] = ssc.maskSensitiveData(k, item)
		}
		return masked
	case []string:
		masked := make([]string, len(v))
		for i, item := range v {
			maskedValue := ssc.maskSensitiveData(key, item)
			if maskedVal, ok := maskedValue.(string); ok {
				masked[i] = maskedVal
			} else {
				masked[i] = fmt.Sprintf("%v", maskedValue)
			}
		}
		return masked
	case []byte:
		return string(v)
	}

	rv := reflect.ValueOf(value)
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return nil
		}
		return ssc.maskSensitiveData(key, rv.Elem().Interface())
	}

	switch rv.Kind() {
	case reflect.Map:
		masked := make(map[string]interface{})
		for _, mapKey := range rv.MapKeys() {
			mapKeyStr := fmt.Sprintf("%v", mapKey.Interface())
			masked[mapKeyStr] = ssc.maskSensitiveData(mapKeyStr, rv.MapIndex(mapKey).Interface())
		}
		return masked
	case reflect.Struct:
		masked := make(map[string]interface{})
		typeOfVal := rv.Type()
		for i := 0; i < rv.NumField(); i++ {
			field := rv.Field(i)
			if !field.CanInterface() {
				continue
			}
			fieldName := typeOfVal.Field(i).Name
			masked[fieldName] = ssc.maskSensitiveData(fieldName, field.Interface())
		}
		return masked
	case reflect.Slice, reflect.Array:
		if rv.Type().Elem().Kind() == reflect.Uint8 {
			return value
		}
		masked := make([]interface{}, rv.Len())
		for i := 0; i < rv.Len(); i++ {
			masked[i] = ssc.maskSensitiveData(key, rv.Index(i).Interface())
		}
		return masked
	}

	return value
}

// maskSensitiveParameters masks sensitive data in parameters
func (ssc *StorySnapshotCapture) maskSensitiveParameters(params interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	if params == nil {
		return result
	}

	// Use reflection to extract parameters
	paramValue := reflect.ValueOf(params)
	if paramValue.Kind() == reflect.Ptr {
		if paramValue.IsNil() {
			return result
		}
		paramValue = paramValue.Elem()
	}

	switch paramValue.Kind() {
	case reflect.Struct:
		paramType := paramValue.Type()
		for i := 0; i < paramValue.NumField(); i++ {
			field := paramValue.Field(i)
			fieldType := paramType.Field(i)

			if field.CanInterface() {
				fieldName := fieldType.Name
				fieldValue := field.Interface()
				result[fieldName] = ssc.maskSensitiveData(fieldName, fieldValue)
			}
		}

	case reflect.Map:
		for _, key := range paramValue.MapKeys() {
			keyStr := fmt.Sprintf("%v", key.Interface())
			value := paramValue.MapIndex(key)
			if value.IsValid() && value.CanInterface() {
				result[keyStr] = ssc.maskSensitiveData(keyStr, value.Interface())
			}
		}

	default:
		// For other types, convert to string representation
		result["value"] = ssc.maskSensitiveData("value", params)
	}

	return result
}

// Helper functions

func (ssc *StorySnapshotCapture) maskSensitiveInterface(data interface{}) interface{} {
	switch v := data.(type) {
	case map[string]interface{}:
		masked := make(map[string]interface{}, len(v))
		for key, value := range v {
			masked[key] = ssc.maskSensitiveData(key, value)
		}
		return masked
	case []interface{}:
		masked := make([]interface{}, len(v))
		for i, item := range v {
			masked[i] = ssc.maskSensitiveInterface(item)
		}
		return masked
	default:
		return data
	}
}

func (ssc *StorySnapshotCapture) maskHeaderValue(key, value string) string {
	masked := ssc.maskSensitiveData(key, value)
	switch v := masked.(type) {
	case nil:
		return ""
	case string:
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
}

func (ssc *StorySnapshotCapture) maskBodyString(body string) string {
	trimmed := strings.TrimSpace(body)
	if trimmed == "" {
		return body
	}

	firstChar := trimmed[0]
	if firstChar != '{' && firstChar != '[' {
		return body
	}

	var payload interface{}
	if err := json.Unmarshal([]byte(trimmed), &payload); err != nil {
		return body
	}

	masked := ssc.maskSensitiveInterface(payload)
	maskedBytes, err := json.MarshalIndent(masked, "", "  ")
	if err != nil {
		return body
	}

	if strings.HasSuffix(body, "\n") {
		return string(maskedBytes) + "\n"
	}
	return string(maskedBytes)
}

func isSensitiveKey(keyLower string) bool {
	for _, keyword := range sensitiveKeySubstrings {
		if strings.Contains(keyLower, keyword) {
			return true
		}
	}
	return false
}

func maskValue(value interface{}) interface{} {
	switch v := value.(type) {
	case string:
		if v == "" {
			return ""
		}
		return maskString(v)
	case []byte:
		if len(v) == 0 {
			return ""
		}
		return maskString(string(v))
	default:
		return "[MASKED]"
	}
}

func maskString(str string) string {
	if str == "" {
		return ""
	}
	return fmt.Sprintf("[MASKED len=%d]", len(str))
}

// getTestEnvironment gets the test environment
func getTestEnvironment() string {
	if env := os.Getenv("TEST_ENVIRONMENT"); env != "" {
		return env
	}
	return "dev"
}

// Helper functions for git information

// getGitTagSafe gets the current git tag safely
func getGitTagSafe() string {
	if tag, err := GetGitTagForFileManager(); err == nil && tag != "" {
		return tag
	}
	if tag := os.Getenv("GIT_TAG"); tag != "" {
		return tag
	}
	return ""
}

// getGitCommitSafe gets the current git commit safely
func getGitCommitSafe() string {
	// getCurrentGitCommit doesn't exist in the existing system, so use environment
	if commit := os.Getenv("GIT_COMMIT"); commit != "" {
		return commit
	}
	return ""
}
