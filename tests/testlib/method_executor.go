package testlib

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"reflect"
)

// MethodExecutor handles the execution of client methods with different signatures
type MethodExecutor struct {
	logger *Logger
	dryRun bool
}

// NewMethodExecutor creates a new method executor
func NewMethodExecutor(logger *Logger, dryRun bool) *MethodExecutor {
	return &MethodExecutor{
		logger: logger,
		dryRun: dryRun,
	}
}

// ExecutionResult contains the result of method execution
type ExecutionResult struct {
	Response   any
	StatusCode int
	Error      error
	BodyBytes  []byte
}

// Execute executes a client method and returns the result
func (m *MethodExecutor) Execute(client any, methodName string, parameters any, variables map[string]any) ExecutionResult {
	m.logger.LogDebug(fmt.Sprintf("Executing client method: %s", methodName))

	if m.dryRun {
		m.logger.LogDebug("DRY RUN: Simulating method execution")
		return ExecutionResult{Response: nil, StatusCode: 200, Error: nil}
	}

	// Resolve parameters if it's a function
	resolvedParams := m.resolveParameters(parameters, variables)

	// Get method via reflection
	method, err := m.getMethod(client, methodName)
	if err != nil {
		return ExecutionResult{Error: err}
	}

	// Build arguments based on method signature
	args, err := m.buildArguments(method, methodName, resolvedParams)
	if err != nil {
		return ExecutionResult{Error: err}
	}

	// Execute method
	results := method.Call(args)
	return m.processResults(results, methodName)
}

// resolveParameters resolves parameter functions to actual values
func (m *MethodExecutor) resolveParameters(parameters any, variables map[string]any) any {
	if parameters == nil {
		return nil
	}

	paramValue := reflect.ValueOf(parameters)
	if paramValue.Kind() == reflect.Func {
		results := paramValue.Call([]reflect.Value{reflect.ValueOf(variables)})
		if len(results) > 0 {
			return results[0].Interface()
		}
	}

	return parameters
}

// getMethod retrieves a method by name from the client
func (m *MethodExecutor) getMethod(client any, methodName string) (reflect.Value, error) {
	clientValue := reflect.ValueOf(client)
	method := clientValue.MethodByName(methodName)

	if !method.IsValid() {
		return reflect.Value{}, fmt.Errorf("method %s not found on client", methodName)
	}

	return method, nil
}

// buildArguments constructs the argument list for method invocation
func (m *MethodExecutor) buildArguments(method reflect.Value, methodName string, parameters any) ([]reflect.Value, error) {
	methodType := method.Type()
	numIn := methodType.NumIn()

	m.logger.LogDebug(fmt.Sprintf("Method %s has %d input parameters", methodName, numIn))

	// Start with context
	args := []reflect.Value{reflect.ValueOf(context.Background())}

	// Add method-specific parameters
	strategy := m.selectArgumentStrategy(methodName, methodType, parameters)
	additionalArgs, err := strategy.BuildArguments(parameters)
	if err != nil {
		return nil, err
	}

	args = append(args, additionalArgs...)
	return args, nil
}

// processResults extracts response and status code from method results
func (m *MethodExecutor) processResults(results []reflect.Value, methodName string) ExecutionResult {
	if len(results) < 2 {
		return ExecutionResult{
			Error: fmt.Errorf("unexpected number of return values from %s", methodName),
		}
	}

	// Check for error (last return value)
	if !results[len(results)-1].IsNil() {
		err := results[len(results)-1].Interface().(error)
		return ExecutionResult{Error: err}
	}

	response := results[0].Interface()
	bodyCopy := m.cloneHTTPResponseBody(response)
	statusCode := m.extractStatusCode(response)

	return ExecutionResult{
		Response:   response,
		StatusCode: statusCode,
		Error:      nil,
		BodyBytes:  bodyCopy,
	}
}

func (m *MethodExecutor) cloneHTTPResponseBody(response any) []byte {
	resp, ok := response.(*http.Response)
	if !ok || resp == nil || resp.Body == nil {
		return nil
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}
	resp.Body.Close()
	resp.Body = io.NopCloser(bytes.NewReader(body))
	return body
}

// extractStatusCode extracts status code from different response types
func (m *MethodExecutor) extractStatusCode(response any) int {
	if response == nil {
		m.logger.LogDebug("[extractStatusCode] response is nil, returning 0")
		return 0
	}

	// Handle *http.Response
	if httpResp, ok := response.(*http.Response); ok {
		m.logger.LogDebug(fmt.Sprintf("[extractStatusCode] Found *http.Response with status code: %d", httpResp.StatusCode))
		return httpResp.StatusCode
	}

	// Handle responses with StatusCode() method
	respValue := reflect.ValueOf(response)
	respType := respValue.Type()
	m.logger.LogDebug(fmt.Sprintf("[extractStatusCode] Response type: %v, Kind: %v", respType, respValue.Kind()))

	if respValue.Kind() == reflect.Ptr {
		respValue = respValue.Elem()
		m.logger.LogDebug(fmt.Sprintf("[extractStatusCode] After Elem(), type: %v, Kind: %v", respValue.Type(), respValue.Kind()))
	}

	statusCodeMethod := respValue.MethodByName("StatusCode")
	if statusCodeMethod.IsValid() {
		m.logger.LogDebug("[extractStatusCode] Found StatusCode() method, calling it...")
		results := statusCodeMethod.Call([]reflect.Value{})
		if len(results) > 0 {
			if statusCode, ok := results[0].Interface().(int); ok {
				m.logger.LogDebug(fmt.Sprintf("[extractStatusCode] StatusCode() returned: %d", statusCode))
				return statusCode
			}
			m.logger.LogDebug(fmt.Sprintf("[extractStatusCode] StatusCode() returned non-int: %v", results[0].Interface()))
		}
	} else {
		m.logger.LogDebug("[extractStatusCode] StatusCode() method not found")
	}

	m.logger.LogDebug("[extractStatusCode] WARNING: Returning default 200 - this may be incorrect!")
	return 200 // Default
}

// ArgumentStrategy defines how to build arguments for different method types
type ArgumentStrategy interface {
	BuildArguments(parameters any) ([]reflect.Value, error)
}

// selectArgumentStrategy selects the appropriate strategy based on method signature
func (m *MethodExecutor) selectArgumentStrategy(methodName string, methodType reflect.Type, parameters any) ArgumentStrategy {
	numIn := methodType.NumIn()

	// Check for WithBody methods
	if isWithBodyMethod(methodName) {
		return &WithBodyArgumentStrategy{logger: m.logger}
	}

	// Check for methods with body parameter
	if hasBodyParameter(methodType) {
		return &BodyArgumentStrategy{logger: m.logger}
	}

	// Check for methods with multiple parameters
	if numIn > 1 && parameters != nil {
		return &MultiParameterArgumentStrategy{logger: m.logger}
	}

	// Default: no additional parameters
	return &NoParameterArgumentStrategy{}
}

func isWithBodyMethod(methodName string) bool {
	return contains(methodName, "WithBody")
}

func hasBodyParameter(methodType reflect.Type) bool {
	if methodType.NumIn() < 2 {
		return false
	}
	// Check if second parameter (after context) looks like a body type
	paramType := methodType.In(1)
	return paramType.Kind() == reflect.Struct || paramType.Kind() == reflect.Interface
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[len(s)-len(substr):] == substr
}

// NoParameterArgumentStrategy handles methods with no additional parameters
type NoParameterArgumentStrategy struct{}

func (s *NoParameterArgumentStrategy) BuildArguments(parameters any) ([]reflect.Value, error) {
	return []reflect.Value{}, nil
}

// BodyArgumentStrategy handles methods with a body parameter
type BodyArgumentStrategy struct {
	logger *Logger
}

func (s *BodyArgumentStrategy) BuildArguments(parameters any) ([]reflect.Value, error) {
	if parameters == nil {
		return nil, fmt.Errorf("body parameter required but not provided")
	}
	return []reflect.Value{reflect.ValueOf(parameters)}, nil
}

// WithBodyArgumentStrategy handles methods with ContentType and Body
type WithBodyArgumentStrategy struct {
	logger *Logger
}

func (s *WithBodyArgumentStrategy) BuildArguments(parameters any) ([]reflect.Value, error) {
	if parameters == nil {
		return nil, fmt.Errorf("WithBody method requires parameters")
	}

	paramValue := reflect.ValueOf(parameters)
	if paramValue.Kind() != reflect.Struct {
		return nil, fmt.Errorf("WithBody method requires struct with ContentType and Body fields")
	}

	contentTypeField := paramValue.FieldByName("ContentType")
	bodyField := paramValue.FieldByName("Body")

	if !contentTypeField.IsValid() || !bodyField.IsValid() {
		return nil, fmt.Errorf("WithBody method requires ContentType and Body fields")
	}

	// Build arguments list - first add any path parameters, then ContentType and Body
	var args []reflect.Value
	paramType := paramValue.Type()

	// Add all fields except ContentType and Body as path parameters
	for i := 0; i < paramValue.NumField(); i++ {
		fieldName := paramType.Field(i).Name
		if fieldName != "ContentType" && fieldName != "Body" {
			field := paramValue.Field(i)
			if field.CanInterface() {
				s.logger.LogDebug(fmt.Sprintf("  Adding path parameter %s", fieldName))
				args = append(args, field)
			}
		}
	}

	// Add ContentType and Body at the end
	args = append(args, contentTypeField, bodyField)

	return args, nil
}

// MultiParameterArgumentStrategy handles methods with multiple parameters
type MultiParameterArgumentStrategy struct {
	logger *Logger
}

func (s *MultiParameterArgumentStrategy) BuildArguments(parameters any) ([]reflect.Value, error) {
	if parameters == nil {
		return []reflect.Value{}, nil
	}

	paramValue := reflect.ValueOf(parameters)

	// If it's a struct, extract fields as separate arguments
	if paramValue.Kind() == reflect.Struct {
		var args []reflect.Value
		numFields := paramValue.NumField()
		paramType := paramValue.Type()

		for i := 0; i < numFields; i++ {
			field := paramValue.Field(i)
			if field.CanInterface() {
				s.logger.LogDebug(fmt.Sprintf("  Adding field %s as argument", paramType.Field(i).Name))
				args = append(args, field)
			}
		}
		return args, nil
	}

	// For non-struct parameters, add directly
	return []reflect.Value{paramValue}, nil
}
