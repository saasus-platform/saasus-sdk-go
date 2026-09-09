package snapshot

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

// SnapshotError represents errors that occur during snapshot operations
type SnapshotError struct {
	Phase     SnapshotPhase     `json:"phase"`
	Type      SnapshotErrorType `json:"type"`
	Message   string            `json:"message"`
	Cause     error             `json:"cause,omitempty"`
	Context   ErrorContext      `json:"context,omitempty"`
	Timestamp time.Time         `json:"timestamp"`
	Severity  ErrorSeverity     `json:"severity"`
}

// SnapshotPhase represents the phase where the error occurred
type SnapshotPhase string

const (
	PhaseCapture    SnapshotPhase = "capture"
	PhaseComparison SnapshotPhase = "compare"
	PhaseReport     SnapshotPhase = "report"
	PhaseValidation SnapshotPhase = "validation"
	PhaseFileIO     SnapshotPhase = "file_io"
	PhaseConfig     SnapshotPhase = "config"
)

// SnapshotErrorType represents the type of error
type SnapshotErrorType string

const (
	ErrorTypeFileIO      SnapshotErrorType = "file_io"
	ErrorTypeParsing     SnapshotErrorType = "parsing"
	ErrorTypeValidation  SnapshotErrorType = "validation"
	ErrorTypeComparison  SnapshotErrorType = "comparison"
	ErrorTypeGeneration  SnapshotErrorType = "generation"
	ErrorTypeConfig      SnapshotErrorType = "config"
	ErrorTypeNetwork     SnapshotErrorType = "network"
	ErrorTypeTimeout     SnapshotErrorType = "timeout"
	ErrorTypeMemory      SnapshotErrorType = "memory"
	ErrorTypePermission  SnapshotErrorType = "permission"
	ErrorTypeFormat      SnapshotErrorType = "format"
	ErrorTypeIntegration SnapshotErrorType = "integration"
)

// ErrorSeverity represents the severity of the error
type ErrorSeverity string

const (
	SeverityCritical ErrorSeverity = "critical"
	SeverityError    ErrorSeverity = "error"
	SeverityWarning  ErrorSeverity = "warning"
	SeverityInfo     ErrorSeverity = "info"
)

// ErrorContext provides additional context about the error
type ErrorContext struct {
	StoryName   string                 `json:"story_name,omitempty"`
	StepName    string                 `json:"step_name,omitempty"`
	FilePath    string                 `json:"file_path,omitempty"`
	Method      string                 `json:"method,omitempty"`
	Operation   string                 `json:"operation,omitempty"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
	StackTrace  string                 `json:"stack_trace,omitempty"`
	Environment string                 `json:"environment,omitempty"`
	GitTag      string                 `json:"git_tag,omitempty"`
	Timestamp   time.Time              `json:"timestamp,omitempty"`
}

// Error implements the error interface
func (se *SnapshotError) Error() string {
	if se.Cause != nil {
		return fmt.Sprintf("[%s:%s] %s: %v", se.Phase, se.Type, se.Message, se.Cause)
	}
	return fmt.Sprintf("[%s:%s] %s", se.Phase, se.Type, se.Message)
}

// Unwrap returns the underlying error for error unwrapping
func (se *SnapshotError) Unwrap() error {
	return se.Cause
}

// IsRetryable returns whether the error is retryable
func (se *SnapshotError) IsRetryable() bool {
	switch se.Type {
	case ErrorTypeNetwork, ErrorTypeTimeout, ErrorTypeFileIO:
		return true
	case ErrorTypeMemory:
		return se.Severity != SeverityCritical
	default:
		return false
	}
}

// IsCritical returns whether the error is critical
func (se *SnapshotError) IsCritical() bool {
	return se.Severity == SeverityCritical
}

// GetSuggestion returns a suggestion for resolving the error
func (se *SnapshotError) GetSuggestion() string {
	switch se.Type {
	case ErrorTypeFileIO:
		return "Check file permissions and disk space. Ensure the output directory exists and is writable."
	case ErrorTypeParsing:
		return "Verify the file format is valid JSON. Check for syntax errors or corrupted data."
	case ErrorTypeValidation:
		return "Review the validation rules and ensure the snapshot data meets the expected format."
	case ErrorTypeComparison:
		return "Check that both snapshots are from compatible versions and have the same story structure."
	case ErrorTypeConfig:
		return "Verify the configuration file syntax and ensure all required fields are present."
	case ErrorTypeNetwork:
		return "Check network connectivity and retry the operation. Consider increasing timeout values."
	case ErrorTypeTimeout:
		return "Increase timeout values in configuration or reduce the scope of the operation."
	case ErrorTypeMemory:
		return "Reduce the number of concurrent operations or increase available memory."
	case ErrorTypePermission:
		return "Check file and directory permissions. Ensure the process has read/write access."
	case ErrorTypeFormat:
		return "Verify the data format matches the expected schema. Check for version compatibility."
	default:
		return "Review the error details and check the documentation for troubleshooting steps."
	}
}

// SnapshotLogger provides logging functionality for snapshot operations
type SnapshotLogger struct {
	level     LogLevel
	output    *log.Logger
	errorLog  *log.Logger
	debugMode bool
	context   map[string]interface{}
}

// LogLevel represents different logging levels
type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarning
	LogLevelError
	LogLevelCritical
)

// String returns the string representation of the log level
func (l LogLevel) String() string {
	switch l {
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelInfo:
		return "INFO"
	case LogLevelWarning:
		return "WARNING"
	case LogLevelError:
		return "ERROR"
	case LogLevelCritical:
		return "CRITICAL"
	default:
		return "UNKNOWN"
	}
}

// NewSnapshotLogger creates a new snapshot logger
func NewSnapshotLogger(level LogLevel, debugMode bool) *SnapshotLogger {
	return &SnapshotLogger{
		level:     level,
		output:    log.New(os.Stdout, "", log.LstdFlags),
		errorLog:  log.New(os.Stderr, "", log.LstdFlags),
		debugMode: debugMode,
		context:   make(map[string]interface{}),
	}
}

// SetContext sets context information for logging
func (sl *SnapshotLogger) SetContext(key string, value interface{}) {
	sl.context[key] = value
}

// ClearContext clears all context information
func (sl *SnapshotLogger) ClearContext() {
	sl.context = make(map[string]interface{})
}

// LogDebug logs debug messages
func (sl *SnapshotLogger) LogDebug(message string, args ...interface{}) {
	if sl.level <= LogLevelDebug && sl.debugMode {
		sl.log(LogLevelDebug, message, args...)
	}
}

// LogInfo logs info messages
func (sl *SnapshotLogger) LogInfo(message string, args ...interface{}) {
	if sl.level <= LogLevelInfo {
		sl.log(LogLevelInfo, message, args...)
	}
}

// LogWarning logs warning messages
func (sl *SnapshotLogger) LogWarning(message string, args ...interface{}) {
	if sl.level <= LogLevelWarning {
		sl.log(LogLevelWarning, message, args...)
	}
}

// LogError logs error messages
func (sl *SnapshotLogger) LogError(message string, err error, args ...interface{}) {
	if sl.level <= LogLevelError {
		if err != nil {
			message = fmt.Sprintf("%s: %v", message, err)
		}
		sl.logToError(LogLevelError, message, args...)
	}
}

// LogCritical logs critical messages
func (sl *SnapshotLogger) LogCritical(message string, err error, args ...interface{}) {
	if err != nil {
		message = fmt.Sprintf("%s: %v", message, err)
	}
	sl.logToError(LogLevelCritical, message, args...)
}

// LogSnapshotError logs a SnapshotError with full context
func (sl *SnapshotLogger) LogSnapshotError(err *SnapshotError) {
	contextStr := sl.formatContext(err.Context)
	message := fmt.Sprintf("[%s] %s %s", err.Phase, err.Error(), contextStr)

	switch err.Severity {
	case SeverityCritical:
		sl.LogCritical(message, err.Cause)
	case SeverityError:
		sl.LogError(message, err.Cause)
	case SeverityWarning:
		sl.LogWarning(message)
	default:
		sl.LogInfo(message)
	}

	// Log suggestion if available
	if suggestion := err.GetSuggestion(); suggestion != "" {
		sl.LogInfo("Suggestion: %s", suggestion)
	}
}

// LogPhaseStart logs the start of a snapshot phase
func (sl *SnapshotLogger) LogPhaseStart(phase SnapshotPhase, operation string) {
	sl.LogInfo("Starting %s phase: %s", phase, operation)
}

// LogPhaseEnd logs the end of a snapshot phase
func (sl *SnapshotLogger) LogPhaseEnd(phase SnapshotPhase, operation string, duration time.Duration) {
	sl.LogInfo("Completed %s phase: %s (took %v)", phase, operation, duration)
}

// LogStoryProcessing logs story processing information
func (sl *SnapshotLogger) LogStoryProcessing(storyName string, operation string) {
	sl.LogInfo("Processing story '%s': %s", storyName, operation)
}

// LogFileOperation logs file operations
func (sl *SnapshotLogger) LogFileOperation(operation, filePath string) {
	sl.LogDebug("File operation: %s - %s", operation, filePath)
}

// LogComparison logs comparison operations
func (sl *SnapshotLogger) LogComparison(storyName string, oldTag, newTag string, differences int) {
	sl.LogInfo("Comparison for story '%s' (%s vs %s): %d differences found",
		storyName, oldTag, newTag, differences)
}

// LogValidation logs validation results
func (sl *SnapshotLogger) LogValidation(storyName string, isValid bool, errors int) {
	if isValid {
		sl.LogInfo("Validation passed for story '%s'", storyName)
	} else {
		sl.LogWarning("Validation failed for story '%s': %d errors", storyName, errors)
	}
}

// LogStoryStart logs the start of story execution
func (sl *SnapshotLogger) LogStoryStart(storyName string) {
	sl.LogInfo("Starting story execution: %s", storyName)
}

// LogStoryEnd logs the end of story execution
func (sl *SnapshotLogger) LogStoryEnd(storyName string, duration time.Duration) {
	sl.LogInfo("Completed story execution: %s (took %v)", storyName, duration)
}

// LogStepStart logs the start of step execution
func (sl *SnapshotLogger) LogStepStart(stepNum int, stepName, method string) {
	sl.LogDebug("Starting step %d: %s (%s)", stepNum, stepName, method)
}

// LogStepResult logs the result of step execution
func (sl *SnapshotLogger) LogStepResult(stepName, method string, status TestStatus, statusCode int, duration time.Duration) {
	sl.LogDebug("Step result: %s (%s) - %s [%d] (took %v)", stepName, method, status, statusCode, duration)
}

// LogStateUpdate logs state updates
func (sl *SnapshotLogger) LogStateUpdate(stepName string, variables map[string]interface{}) {
	sl.LogDebug("State updated in step %s: %d variables", stepName, len(variables))
}

// log logs a message with the specified level
func (sl *SnapshotLogger) log(level LogLevel, message string, args ...interface{}) {
	contextStr := sl.formatCurrentContext()
	fullMessage := fmt.Sprintf("[%s] %s %s", level.String(), fmt.Sprintf(message, args...), contextStr)
	sl.output.Println(fullMessage)
}

// logToError logs a message to stderr
func (sl *SnapshotLogger) logToError(level LogLevel, message string, args ...interface{}) {
	contextStr := sl.formatCurrentContext()
	fullMessage := fmt.Sprintf("[%s] %s %s", level.String(), fmt.Sprintf(message, args...), contextStr)
	sl.errorLog.Println(fullMessage)
}

// formatCurrentContext formats the current context for logging
func (sl *SnapshotLogger) formatCurrentContext() string {
	if len(sl.context) == 0 {
		return ""
	}

	var parts []string
	for key, value := range sl.context {
		parts = append(parts, fmt.Sprintf("%s=%v", key, value))
	}
	return fmt.Sprintf("[%s]", strings.Join(parts, ", "))
}

// formatContext formats an ErrorContext for logging
func (sl *SnapshotLogger) formatContext(ctx ErrorContext) string {
	var parts []string

	if ctx.StoryName != "" {
		parts = append(parts, fmt.Sprintf("story=%s", ctx.StoryName))
	}
	if ctx.StepName != "" {
		parts = append(parts, fmt.Sprintf("step=%s", ctx.StepName))
	}
	if ctx.FilePath != "" {
		parts = append(parts, fmt.Sprintf("file=%s", ctx.FilePath))
	}
	if ctx.Method != "" {
		parts = append(parts, fmt.Sprintf("method=%s", ctx.Method))
	}
	if ctx.Operation != "" {
		parts = append(parts, fmt.Sprintf("op=%s", ctx.Operation))
	}

	if len(parts) == 0 {
		return ""
	}
	return fmt.Sprintf("[%s]", strings.Join(parts, ", "))
}

// Error creation helper functions

// NewSnapshotError creates a new SnapshotError
func NewSnapshotError(phase SnapshotPhase, errorType SnapshotErrorType, message string, cause error) *SnapshotError {
	return &SnapshotError{
		Phase:     phase,
		Type:      errorType,
		Message:   message,
		Cause:     cause,
		Timestamp: time.Now(),
		Severity:  determineSeverity(errorType, cause),
	}
}

// NewSnapshotErrorWithContext creates a new SnapshotError with context
func NewSnapshotErrorWithContext(phase SnapshotPhase, errorType SnapshotErrorType, message string, cause error, context ErrorContext) *SnapshotError {
	return &SnapshotError{
		Phase:     phase,
		Type:      errorType,
		Message:   message,
		Cause:     cause,
		Context:   context,
		Timestamp: time.Now(),
		Severity:  determineSeverity(errorType, cause),
	}
}

// WrapError wraps an existing error as a SnapshotError
func WrapError(phase SnapshotPhase, errorType SnapshotErrorType, message string, err error) *SnapshotError {
	if err == nil {
		return nil
	}

	// If it's already a SnapshotError, preserve the original context
	if snapErr, ok := err.(*SnapshotError); ok {
		return &SnapshotError{
			Phase:     phase,
			Type:      errorType,
			Message:   message,
			Cause:     snapErr,
			Context:   snapErr.Context,
			Timestamp: time.Now(),
			Severity:  snapErr.Severity,
		}
	}

	return NewSnapshotError(phase, errorType, message, err)
}

// determineSeverity determines the severity based on error type and cause
func determineSeverity(errorType SnapshotErrorType, cause error) ErrorSeverity {
	switch errorType {
	case ErrorTypeMemory:
		return SeverityCritical
	case ErrorTypePermission, ErrorTypeConfig:
		return SeverityError
	case ErrorTypeFileIO, ErrorTypeParsing, ErrorTypeValidation:
		return SeverityError
	case ErrorTypeComparison, ErrorTypeGeneration:
		return SeverityWarning
	case ErrorTypeNetwork, ErrorTypeTimeout:
		return SeverityWarning
	case ErrorTypeFormat, ErrorTypeIntegration:
		return SeverityWarning
	default:
		return SeverityError
	}
}

// Specific error creation functions for common scenarios

// NewCaptureError creates an error for the capture phase
func NewCaptureError(errorType SnapshotErrorType, message string, cause error, storyName, stepName string) *SnapshotError {
	context := ErrorContext{
		StoryName: storyName,
		StepName:  stepName,
		Operation: "capture",
		Timestamp: time.Now(),
	}
	return NewSnapshotErrorWithContext(PhaseCapture, errorType, message, cause, context)
}

// NewComparisonError creates an error for the comparison phase
func NewComparisonError(errorType SnapshotErrorType, message string, cause error, storyName string) *SnapshotError {
	context := ErrorContext{
		StoryName: storyName,
		Operation: "compare",
		Timestamp: time.Now(),
	}
	return NewSnapshotErrorWithContext(PhaseComparison, errorType, message, cause, context)
}

// NewReportError creates an error for the report phase
func NewReportError(errorType SnapshotErrorType, message string, cause error, storyName, filePath string) *SnapshotError {
	context := ErrorContext{
		StoryName: storyName,
		FilePath:  filePath,
		Operation: "report",
		Timestamp: time.Now(),
	}
	return NewSnapshotErrorWithContext(PhaseReport, errorType, message, cause, context)
}

// NewValidationError creates an error for the validation phase
func NewValidationError(errorType SnapshotErrorType, message string, cause error, storyName, stepName string) *SnapshotError {
	context := ErrorContext{
		StoryName: storyName,
		StepName:  stepName,
		Operation: "validate",
		Timestamp: time.Now(),
	}
	return NewSnapshotErrorWithContext(PhaseValidation, errorType, message, cause, context)
}

// NewFileIOError creates an error for file I/O operations
func NewFileIOError(message string, cause error, filePath string) *SnapshotError {
	context := ErrorContext{
		FilePath:  filePath,
		Operation: "file_io",
		Timestamp: time.Now(),
	}
	return NewSnapshotErrorWithContext(PhaseFileIO, ErrorTypeFileIO, message, cause, context)
}

// NewConfigError creates an error for configuration issues
func NewConfigError(message string, cause error, filePath string) *SnapshotError {
	context := ErrorContext{
		FilePath:  filePath,
		Operation: "config",
		Timestamp: time.Now(),
	}
	return NewSnapshotErrorWithContext(PhaseConfig, ErrorTypeConfig, message, cause, context)
}

// IsSnapshotError checks if an error is a SnapshotError
func IsSnapshotError(err error) bool {
	_, ok := err.(*SnapshotError)
	return ok
}

// GetSnapshotError extracts a SnapshotError from an error chain
func GetSnapshotError(err error) *SnapshotError {
	if snapErr, ok := err.(*SnapshotError); ok {
		return snapErr
	}
	return nil
}
