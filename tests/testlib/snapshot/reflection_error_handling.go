package snapshot

import (
	"fmt"
	"reflect"
	"strings"
)

// ArgumentMismatchError represents an error when method arguments don't match expectations
type ArgumentMismatchError struct {
	MethodName    string
	Expected      int
	Actual        int
	ExpectedTypes []reflect.Type
	ActualTypes   []reflect.Type
	Details       string
}

// Error implements the error interface for ArgumentMismatchError
func (e *ArgumentMismatchError) Error() string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("引数数不一致エラー: メソッド '%s' は %d 個の引数を期待していますが、%d 個が渡されました",
		e.MethodName, e.Expected, e.Actual))

	if len(e.ExpectedTypes) > 0 {
		builder.WriteString("\n期待される引数タイプ: ")
		for i, typ := range e.ExpectedTypes {
			if i > 0 {
				builder.WriteString(", ")
			}
			builder.WriteString(typ.String())
		}
	}

	if len(e.ActualTypes) > 0 {
		builder.WriteString("\n実際の引数タイプ: ")
		for i, typ := range e.ActualTypes {
			if i > 0 {
				builder.WriteString(", ")
			}
			builder.WriteString(typ.String())
		}
	}

	if e.Details != "" {
		builder.WriteString("\n詳細: ")
		builder.WriteString(e.Details)
	}

	return builder.String()
}

// ParameterTypeError represents an error when parameter types don't match expectations
type ParameterTypeError struct {
	MethodName    string
	ParameterName string
	ExpectedType  reflect.Type
	ActualType    reflect.Type
	Details       string
}

// Error implements the error interface for ParameterTypeError
func (e *ParameterTypeError) Error() string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("パラメータ型不適合エラー: メソッド '%s' のパラメータ '%s'",
		e.MethodName, e.ParameterName))

	if e.ExpectedType != nil {
		builder.WriteString(fmt.Sprintf("\n期待される型: %s", e.ExpectedType.String()))
	}

	if e.ActualType != nil {
		builder.WriteString(fmt.Sprintf("\n実際の型: %s", e.ActualType.String()))
	}

	if e.Details != "" {
		builder.WriteString("\n詳細: ")
		builder.WriteString(e.Details)
	}

	return builder.String()
}

// ReflectionErrorHandler handles reflection-related errors with detailed logging
type ReflectionErrorHandler struct {
	logger *SnapshotLogger
}

// NewReflectionErrorHandler creates a new reflection error handler
func NewReflectionErrorHandler(logger *SnapshotLogger) *ReflectionErrorHandler {
	return &ReflectionErrorHandler{
		logger: logger,
	}
}

// HandleArgumentMismatch handles argument mismatch errors with detailed logging
func (h *ReflectionErrorHandler) HandleArgumentMismatch(
	methodName string,
	methodType reflect.Type,
	actualArgs []reflect.Value,
) *ArgumentMismatchError {

	expected := methodType.NumIn()
	actual := len(actualArgs)

	// Extract expected types
	expectedTypes := make([]reflect.Type, expected)
	for i := 0; i < expected; i++ {
		expectedTypes[i] = methodType.In(i)
	}

	// Extract actual types
	actualTypes := make([]reflect.Type, actual)
	for i, arg := range actualArgs {
		actualTypes[i] = arg.Type()
	}

	// Create detailed error information
	details := h.generateArgumentMismatchDetails(methodName, methodType, actualArgs)

	err := &ArgumentMismatchError{
		MethodName:    methodName,
		Expected:      expected,
		Actual:        actual,
		ExpectedTypes: expectedTypes,
		ActualTypes:   actualTypes,
		Details:       details,
	}

	// Log the error with detailed information
	h.logArgumentMismatchError(err)

	return err
}

// HandleParameterTypeError handles parameter type errors with detailed logging
func (h *ReflectionErrorHandler) HandleParameterTypeError(
	methodName string,
	parameterName string,
	expectedType reflect.Type,
	actualType reflect.Type,
	details string,
) *ParameterTypeError {

	err := &ParameterTypeError{
		MethodName:    methodName,
		ParameterName: parameterName,
		ExpectedType:  expectedType,
		ActualType:    actualType,
		Details:       details,
	}

	// Log the error with detailed information
	h.logParameterTypeError(err)

	return err
}

// LogMethodCallAttempt logs detailed information about a method call attempt
func (h *ReflectionErrorHandler) LogMethodCallAttempt(
	methodName string,
	methodType reflect.Type,
	parameters interface{},
) {
	if h.logger == nil {
		return
	}

	h.logger.LogDebug(fmt.Sprintf("メソッド呼び出し試行: %s", methodName))
	h.logger.LogDebug(fmt.Sprintf("メソッドシグネチャ: %s", methodType.String()))

	// Log expected arguments
	numIn := methodType.NumIn()
	h.logger.LogDebug(fmt.Sprintf("期待される引数数: %d", numIn))

	for i := 0; i < numIn; i++ {
		argType := methodType.In(i)
		h.logger.LogDebug(fmt.Sprintf("  引数 %d: %s", i+1, argType.String()))
	}

	// Log provided parameters
	if parameters != nil {
		paramType := reflect.TypeOf(parameters)
		h.logger.LogDebug(fmt.Sprintf("提供されたパラメータ型: %s", paramType.String()))

		// If it's a struct, log field information
		if paramType.Kind() == reflect.Struct {
			h.logStructFields(paramType)
		}
	} else {
		h.logger.LogDebug("提供されたパラメータ: nil")
	}
}

// LogMethodCallSuccess logs successful method call information
func (h *ReflectionErrorHandler) LogMethodCallSuccess(
	methodName string,
	statusCode int,
	responseType reflect.Type,
) {
	if h.logger == nil {
		return
	}

	h.logger.LogDebug(fmt.Sprintf("メソッド呼び出し成功: %s", methodName))
	h.logger.LogDebug(fmt.Sprintf("ステータスコード: %d", statusCode))

	if responseType != nil {
		h.logger.LogDebug(fmt.Sprintf("レスポンス型: %s", responseType.String()))
	}
}

// LogMethodCallFailure logs method call failure information
func (h *ReflectionErrorHandler) LogMethodCallFailure(
	methodName string,
	err error,
) {
	if h.logger == nil {
		return
	}

	h.logger.LogError(fmt.Sprintf("メソッド呼び出し失敗: %s", methodName), err)

	// Log additional context based on error type
	switch err.(type) {
	case *ArgumentMismatchError:
		h.logger.LogError("引数数不一致が検出されました", nil)
	case *ParameterTypeError:
		h.logger.LogError("パラメータ型不適合が検出されました", nil)
	default:
		h.logger.LogError("予期しないエラーが発生しました", nil)
	}
}

// generateArgumentMismatchDetails generates detailed information about argument mismatch
func (h *ReflectionErrorHandler) generateArgumentMismatchDetails(
	methodName string,
	methodType reflect.Type,
	actualArgs []reflect.Value,
) string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("メソッド '%s' の詳細分析:\n", methodName))

	// Method signature analysis
	builder.WriteString(fmt.Sprintf("メソッドシグネチャ: %s\n", methodType.String()))

	// Expected arguments analysis
	numIn := methodType.NumIn()
	builder.WriteString(fmt.Sprintf("期待される引数 (%d個):\n", numIn))
	for i := 0; i < numIn; i++ {
		argType := methodType.In(i)
		builder.WriteString(fmt.Sprintf("  %d: %s\n", i+1, argType.String()))
	}

	// Actual arguments analysis
	builder.WriteString(fmt.Sprintf("実際の引数 (%d個):\n", len(actualArgs)))
	for i, arg := range actualArgs {
		builder.WriteString(fmt.Sprintf("  %d: %s (値: %v)\n", i+1, arg.Type().String(), arg.Interface()))
	}

	// Suggest possible solutions
	builder.WriteString("推奨される解決策:\n")
	if len(actualArgs) < numIn {
		builder.WriteString("- 不足している引数を追加してください\n")
		for i := len(actualArgs); i < numIn; i++ {
			argType := methodType.In(i)
			builder.WriteString(fmt.Sprintf("  - 引数 %d: %s\n", i+1, argType.String()))
		}
	} else if len(actualArgs) > numIn {
		builder.WriteString("- 余分な引数を削除してください\n")
	}

	return builder.String()
}

// logArgumentMismatchError logs argument mismatch error details
func (h *ReflectionErrorHandler) logArgumentMismatchError(err *ArgumentMismatchError) {
	if h.logger == nil {
		return
	}

	h.logger.LogError("引数数不一致エラーが発生しました", err)
	h.logger.LogError(fmt.Sprintf("メソッド: %s", err.MethodName), nil)
	h.logger.LogError(fmt.Sprintf("期待: %d引数, 実際: %d引数", err.Expected, err.Actual), nil)

	if len(err.ExpectedTypes) > 0 {
		h.logger.LogDebug("期待される引数タイプ:")
		for i, typ := range err.ExpectedTypes {
			h.logger.LogDebug(fmt.Sprintf("  %d: %s", i+1, typ.String()))
		}
	}

	if len(err.ActualTypes) > 0 {
		h.logger.LogDebug("実際の引数タイプ:")
		for i, typ := range err.ActualTypes {
			h.logger.LogDebug(fmt.Sprintf("  %d: %s", i+1, typ.String()))
		}
	}
}

// logParameterTypeError logs parameter type error details
func (h *ReflectionErrorHandler) logParameterTypeError(err *ParameterTypeError) {
	if h.logger == nil {
		return
	}

	h.logger.LogError("パラメータ型不適合エラーが発生しました", err)
	h.logger.LogError(fmt.Sprintf("メソッド: %s", err.MethodName), nil)
	h.logger.LogError(fmt.Sprintf("パラメータ: %s", err.ParameterName), nil)

	if err.ExpectedType != nil {
		h.logger.LogError(fmt.Sprintf("期待される型: %s", err.ExpectedType.String()), nil)
	}

	if err.ActualType != nil {
		h.logger.LogError(fmt.Sprintf("実際の型: %s", err.ActualType.String()), nil)
	}
}

// logStructFields logs struct field information for debugging
func (h *ReflectionErrorHandler) logStructFields(structType reflect.Type) {
	if h.logger == nil || structType.Kind() != reflect.Struct {
		return
	}

	numFields := structType.NumField()
	h.logger.LogDebug(fmt.Sprintf("構造体フィールド (%d個):", numFields))

	for i := 0; i < numFields; i++ {
		field := structType.Field(i)
		h.logger.LogDebug(fmt.Sprintf("  %s: %s", field.Name, field.Type.String()))
	}
}

// ValidateMethodSignature validates if a method signature matches expected patterns
func (h *ReflectionErrorHandler) ValidateMethodSignature(
	methodName string,
	methodType reflect.Type,
) error {
	if methodType == nil {
		return fmt.Errorf("メソッド '%s' のタイプが nil です", methodName)
	}

	if methodType.Kind() != reflect.Func {
		return fmt.Errorf("メソッド '%s' は関数型ではありません: %s", methodName, methodType.Kind())
	}

	// Check minimum requirements (at least context parameter)
	numIn := methodType.NumIn()
	if numIn < 1 {
		return fmt.Errorf("メソッド '%s' は最低1つの引数（context）が必要です", methodName)
	}

	// Check if first parameter is context
	firstParam := methodType.In(0)
	if !h.isContextType(firstParam) {
		h.logger.LogWarning(fmt.Sprintf("メソッド '%s' の最初の引数がcontextではない可能性があります: %s",
			methodName, firstParam.String()))
	}

	// Check return values (should have at least error)
	numOut := methodType.NumOut()
	if numOut < 1 {
		return fmt.Errorf("メソッド '%s' は最低1つの戻り値が必要です", methodName)
	}

	// Check if last return value is error
	lastReturn := methodType.Out(numOut - 1)
	if !h.isErrorType(lastReturn) {
		h.logger.LogWarning(fmt.Sprintf("メソッド '%s' の最後の戻り値がerrorではない可能性があります: %s",
			methodName, lastReturn.String()))
	}

	return nil
}

// isContextType checks if a type is context.Context
func (h *ReflectionErrorHandler) isContextType(typ reflect.Type) bool {
	return typ.String() == "context.Context"
}

// isErrorType checks if a type is error interface
func (h *ReflectionErrorHandler) isErrorType(typ reflect.Type) bool {
	return typ.String() == "error"
}

// CreateDetailedErrorMessage creates a detailed error message for debugging
func (h *ReflectionErrorHandler) CreateDetailedErrorMessage(
	methodName string,
	originalError error,
	context map[string]interface{},
) string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("メソッド '%s' でエラーが発生しました\n", methodName))
	builder.WriteString(fmt.Sprintf("元のエラー: %s\n", originalError.Error()))

	if len(context) > 0 {
		builder.WriteString("コンテキスト情報:\n")
		for key, value := range context {
			builder.WriteString(fmt.Sprintf("  %s: %v\n", key, value))
		}
	}

	builder.WriteString("デバッグのヒント:\n")
	builder.WriteString("1. メソッドシグネチャを確認してください\n")
	builder.WriteString("2. 引数の数と型が正しいか確認してください\n")
	builder.WriteString("3. パラメータ構造体のフィールドが正しいか確認してください\n")

	return builder.String()
}
