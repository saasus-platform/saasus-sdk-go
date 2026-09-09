package snapshot

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/saasus-platform/saasus-sdk-go/tests/testlib"
)

// E2EEngineWithSnapshot extends the existing E2EEngine with snapshot capabilities
// This provides non-destructive extension of the base E2EEngine
type E2EEngineWithSnapshot struct {
	*testlib.E2EEngine
	SnapshotCapture      *StorySnapshotCapture
	Comparator           *StoryComparator
	Reporter             *StoryReporter
	FileManager          *StorySnapshotFileManager
	Validator            *SimpleStoryValidator
	Config               *StorySnapshotConfig
	Logger               *SnapshotLogger
	ErrorHandler         *ReflectionErrorHandler
	FallbackMechanism    *FallbackMechanism    // 要件 4.1, 4.2: 後方互換性保証のためのフォールバック機構
	CompatibilityChecker *CompatibilityChecker // 要件 4.1, 4.2: 既存テストケースが影響を受けないことを保証する機能
}

// NewE2EEngineWithSnapshot creates a new E2E engine with snapshot capabilities
func NewE2EEngineWithSnapshot(client interface{}, methods []string, config *StorySnapshotConfig) *E2EEngineWithSnapshot {
	// Create base E2E engine
	baseEngine := testlib.NewE2EEngine(client, methods)

	// Use default config if none provided
	if config == nil {
		config = DefaultStorySnapshotConfig()
	}

	// Create snapshot capture
	snapshotCapture := NewStorySnapshotCapture(config)

	// Create comparator
	comparator := NewStoryComparator(config)

	// Create reporter
	reporter := NewStoryReporter(config)

	// Create file manager scoped to the module output directory
	fileManager := NewStorySnapshotFileManager(config.GetModuleOutputDirectory())

	// Create validator
	validator := NewSimpleStoryValidator(config)

	// Create logger
	logLevel := LogLevelInfo
	debugMode := false
	if config.CaptureLevel == CaptureLevelFull {
		logLevel = LogLevelDebug
		debugMode = true
	}
	logger := NewSnapshotLogger(logLevel, debugMode)

	// Create error handler
	errorHandler := NewReflectionErrorHandler(logger)

	// Create fallback mechanism
	// 要件 4.1, 4.2: 後方互換性保証のためのフォールバック機構
	fallbackMechanism := NewFallbackMechanism(logger, errorHandler)

	// Create compatibility checker
	// 要件 4.1, 4.2: 既存テストケースが影響を受けないことを保証する機能
	compatibilityChecker := NewCompatibilityChecker(logger, fallbackMechanism, client)

	return &E2EEngineWithSnapshot{
		E2EEngine:            baseEngine,
		SnapshotCapture:      snapshotCapture,
		Comparator:           comparator,
		Reporter:             reporter,
		FileManager:          fileManager,
		Validator:            validator,
		Config:               config,
		Logger:               logger,
		ErrorHandler:         errorHandler,
		FallbackMechanism:    fallbackMechanism,
		CompatibilityChecker: compatibilityChecker,
	}
}

// ExecuteStoriesWithSnapshot executes stories and captures snapshots (integrated mode)
// This method implements the full 3-phase workflow: capture, compare, and report
// Each phase can be enabled/disabled via configuration options
// 要件 4.1, 4.2: 既存テストケースが影響を受けないことを保証する機能
func (e *E2EEngineWithSnapshot) ExecuteStoriesWithSnapshot(stories []testlib.Story) []StorySnapshotResult {
	var results []StorySnapshotResult

	e.Logger.LogInfo(fmt.Sprintf("Starting integrated story snapshot execution with %d stories", len(stories)))
	e.Logger.LogInfo(fmt.Sprintf("Configuration: Capture=%v, Comparison=%v, Reporting=%v, Validation=%v",
		e.Config.EnableCapture, e.Config.EnableComparison, e.Config.EnableReporting, e.Config.EnableValidation))

	// Perform backward compatibility check before execution
	// 要件 4.1: 既存のメソッド呼び出しパターンを維持するフォールバック処理
	if e.Config.EnableCompatibilityCheck {
		e.Logger.LogInfo("Performing backward compatibility check...")
		compatibilityReport := e.performCompatibilityCheck(stories)
		e.logCompatibilityReport(compatibilityReport)

		// Log warnings if fallback is required for any methods
		if len(compatibilityReport.FallbackRequired) > 0 {
			e.Logger.LogWarning(fmt.Sprintf("Fallback mechanism will be used for %d methods: %v",
				len(compatibilityReport.FallbackRequired), compatibilityReport.FallbackRequired))
		}

	}

	startTime := time.Now()

	for i, story := range stories {
		e.Logger.LogInfo(fmt.Sprintf("Processing story %d/%d: %s", i+1, len(stories), story.Name))

		result := e.executeStoryWithIntegratedSnapshot(story)
		results = append(results, result)

		// Log result summary
		e.logStoryResultSummary(result)
	}

	totalDuration := time.Since(startTime)
	e.Logger.LogInfo(fmt.Sprintf("Completed integrated story snapshot execution in %v", totalDuration))
	e.logExecutionSummary(results)

	return results
}

// ExecuteStoryWithSnapshot executes a single story and captures its snapshot (legacy method)
func (e *E2EEngineWithSnapshot) ExecuteStoryWithSnapshot(story testlib.Story) StorySnapshotResult {
	return e.executeStoryWithIntegratedSnapshot(story)
}

// executeStoryWithIntegratedSnapshot executes a single story with full 3-phase integration
func (e *E2EEngineWithSnapshot) executeStoryWithIntegratedSnapshot(story testlib.Story) StorySnapshotResult {
	startTime := time.Now()
	storyStartTime := time.Now()

	e.Logger.LogInfo(fmt.Sprintf("Starting integrated execution for story: %s", story.Name))

	// Initialize result
	result := StorySnapshotResult{}

	// Skip E2E execution if only comparison or reporting is enabled
	if !e.Config.EnableCapture && (e.Config.EnableComparison || e.Config.EnableReporting) {
		e.Logger.LogInfo("Skipping E2E execution (compare/report only mode)")
		// Load existing snapshot for comparison/reporting
		snapshot, err := e.loadExistingSnapshot(story.Name)
		if err != nil {
			e.Logger.LogError(fmt.Sprintf("Failed to load existing snapshot for story '%s'", story.Name), err)
			return result
		}
		result.Snapshot = *snapshot
		// Set story result status to passed for compare/report only modes
		result.StoryResult = StoryResult{
			StoryName: story.Name,
			Status:    TestStatusPassed,
			Duration:  0,
		}

		// Phase 2: Comparison (if enabled)
		if e.Config.EnableComparison {
			e.Logger.LogInfo("Phase 2: Comparing with previous release")
			comparison, err := e.compareWithPreviousRelease(snapshot)
			if err != nil {
				e.Logger.LogError("Failed to compare with previous release", err)
			} else {
				result.Comparison = comparison
				if err := e.saveComparisonToFile(comparison); err != nil {
					e.Logger.LogError("Failed to save comparison result", err)
				}

				// Phase 3: Report Generation (if enabled)
				if e.Config.EnableReporting && comparison != nil {
					e.Logger.LogInfo("Phase 3: Generating story report")
					report, err := e.generateStoryReport(comparison)
					if err != nil {
						e.Logger.LogError("Failed to generate story report", err)
					} else {
						result.Report = report
						if err := e.saveReportToFile(report); err != nil {
							e.Logger.LogError("Failed to save story report", err)
						}
					}
				}
			}
		} else if e.Config.EnableReporting {
			// Report only mode: load existing comparison
			e.Logger.LogInfo("Phase 3: Generating story report (report only mode)")
			comparison, err := e.loadExistingComparison(story.Name)
			if err != nil {
				e.Logger.LogError(fmt.Sprintf("Failed to load existing comparison for story '%s'", story.Name), err)
			} else {
				result.Comparison = comparison
				report, err := e.generateStoryReport(comparison)
				if err != nil {
					e.Logger.LogError("Failed to generate story report", err)
				} else {
					result.Report = report
					if err := e.saveReportToFile(report); err != nil {
						e.Logger.LogError("Failed to save story report", err)
					}
				}
			}
		}
		result.StoryResult.Duration = time.Since(storyStartTime)
		return result
	}

	// Phase 1: E2E Execution + Snapshot Capture
	e.Logger.LogInfo("Phase 1: Executing E2E test and capturing snapshot")
	storyResult := e.executeStoryWithResponseCapture(story)
	result.StoryResult = storyResult

	// Only proceed with snapshot phases if story execution was successful or if configured to capture failures
	if storyResult.Status == TestStatusFailed && !e.shouldCaptureFailedStories() {
		e.Logger.LogInfo(fmt.Sprintf("Skipping snapshot capture for failed story: %s", story.Name))
		result.StoryResult.Duration = time.Since(storyStartTime)
		return result
	}

	// Capture snapshot if enabled
	if e.Config.EnableCapture {
		snapshot, err := e.captureStorySnapshot(story, storyResult)
		if err != nil {
			e.Logger.LogError("Failed to capture story snapshot", err)
			result.StoryResult.Duration = time.Since(storyStartTime)
			return result
		}
		result.Snapshot = *snapshot

		// Save snapshot to file
		if err := e.saveSnapshotToFile(snapshot); err != nil {
			e.Logger.LogError("Failed to save story snapshot", err)
		}

		// Phase 2: Comparison (if enabled)
		if e.Config.EnableComparison {
			e.Logger.LogInfo("Phase 2: Comparing with previous release")
			comparison, err := e.compareWithPreviousRelease(snapshot)
			if err != nil {
				e.Logger.LogError("Failed to compare with previous release", err)
			} else {
				result.Comparison = comparison

				// Save comparison result
				if err := e.saveComparisonToFile(comparison); err != nil {
					e.Logger.LogError("Failed to save comparison result", err)
				}

				// Phase 3: Report Generation (if enabled)
				if e.Config.EnableReporting && comparison != nil {
					e.Logger.LogInfo("Phase 3: Generating story report")
					report, err := e.generateStoryReport(comparison)
					if err != nil {
						e.Logger.LogError("Failed to generate story report", err)
					} else {
						result.Report = report

						// Save report
						if err := e.saveReportToFile(report); err != nil {
							e.Logger.LogError("Failed to save story report", err)
						}
					}
				}
			}
		}

		// Story Validation (if enabled)
		if e.Config.EnableValidation {
			e.Logger.LogInfo("Validating story snapshot")
			validation, err := e.validateStorySnapshot(snapshot)
			if err != nil {
				e.Logger.LogError("Failed to validate story snapshot", err)
			} else {
				result.Validation = validation

				// Save validation result
				if err := e.saveValidationToFile(validation); err != nil {
					e.Logger.LogError("Failed to save validation result", err)
				}
			}
		}
	}

	result.StoryResult.Duration = time.Since(storyStartTime)
	totalDuration := time.Since(startTime)
	e.Logger.LogInfo(fmt.Sprintf("Completed integrated execution for story '%s' in %v", story.Name, totalDuration))

	return result
}

// executeStoryWithResponseCapture executes a story with enhanced response capture
func (e *E2EEngineWithSnapshot) executeStoryWithResponseCapture(story testlib.Story) StoryResult {
	startTime := time.Now()

	e.Logger.LogStoryStart(story.Name)

	result := StoryResult{
		StoryName: story.Name,
		Status:    TestStatusPassed,
		Variables: make(map[string]interface{}),
	}

	// Copy initial variables
	for k, v := range story.Variables {
		result.Variables[k] = v
	}

	// Execute setup
	if story.Setup != nil {
		if err := story.Setup(); err != nil {
			result.Status = TestStatusFailed
			result.Error = fmt.Errorf("setup failed: %w", err)
			result.Duration = time.Since(startTime)
			return result
		}
	}

	// Execute steps with response capture
	for i, step := range story.Steps {
		stepResult := e.executeStepWithResponseCapture(step, i+1, result.Variables)
		result.Steps = append(result.Steps, stepResult)

		if stepResult.Status == TestStatusFailed {
			result.Status = TestStatusFailed
			result.Error = stepResult.Error
			break
		}
	}

	// Execute cleanup
	if story.Cleanup != nil {
		if err := story.Cleanup(); err != nil {
			e.Logger.LogError("Cleanup failed", err)
		}
	}

	result.Duration = time.Since(startTime)
	e.Logger.LogStoryEnd(story.Name, result.Duration)

	return result
}

// executeStepWithResponseCapture executes a step with enhanced response capture
func (e *E2EEngineWithSnapshot) executeStepWithResponseCapture(step testlib.Step, stepNum int, variables map[string]interface{}) StepResult {
	startTime := time.Now()

	e.Logger.LogStepStart(stepNum, step.Name, step.ClientMethod)

	result := StepResult{
		StepName: step.Name,
		Method:   step.ClientMethod,
		Status:   TestStatusPassed,
	}

	if step.Skip {
		result.Status = TestStatusSkipped
		result.Duration = 0
		result.StatusCode = 0
		result.SkipReason = step.SkipReason
		result.ReturnValue = nil // Set empty return value for skipped steps
		e.Logger.LogInfo(fmt.Sprintf("⏭️  Skipping step '%s' (%s)", step.Name, step.ClientMethod))
		e.Logger.LogStepResult(step.Name, step.ClientMethod, result.Status, result.StatusCode, result.Duration)
		return result
	}

	// Log variables before execution
	e.Logger.LogDebug(fmt.Sprintf("Variables before step %s: %d variables", step.Name, len(variables)))
	for k := range variables {
		e.Logger.LogDebug(fmt.Sprintf("  - %s", k))
	}

	// Execute the client method with response capture
	response, statusCode, returnValue, err := e.ExecuteClientMethodWithResponseCapture(step.ClientMethod, step.Parameters, variables)

	result.StatusCode = statusCode
	result.Duration = time.Since(startTime)

	if err != nil {
		result.Status = TestStatusFailed
		result.Error = err
		e.Logger.LogError(fmt.Sprintf("Step '%s' failed", step.Name), err)
	} else {
		statusErr := testlib.ValidateExpectedStatus(step, statusCode)
		success := statusErr == nil
		errorMsg := ""
		if !success {
			errorMsg = statusErr.Error()
		}

		e.Coverage.RecordExecution(
			step.ClientMethod,
			"", // Story name will be set by caller
			step.Name,
			statusCode,
			result.Duration,
			success,
			errorMsg,
		)

		if statusErr != nil {
			result.Status = TestStatusFailed
			result.Error = statusErr
			e.Logger.LogError(fmt.Sprintf("Step '%s' failed", step.Name), statusErr)
			return result
		}

		// Store return value in step result for snapshot capture
		if e.Config.EnableCapture && returnValue != nil {
			result.ReturnValue = returnValue
		}

		// Execute validation if provided
		if step.Validation != nil {
			if validationErr := step.Validation(response); validationErr != nil {
				result.Status = TestStatusFailed
				result.Error = validationErr
				e.Logger.LogError(fmt.Sprintf("Validation failed for step '%s'", step.Name), validationErr)
				e.Logger.LogValidation(step.Name, false, 1)
			} else {
				e.Logger.LogValidation(step.Name, true, 0)
			}
		}

		// Execute state update if provided
		if step.StateUpdate != nil {
			// Try to use captured JSON data if available, otherwise use original response
			var stateUpdateResponse interface{} = response
			if returnValue != nil && returnValue.JSONData != nil && len(returnValue.JSONData) > 0 {
				// If there's only one JSON field (e.g., JSON200, JSON201), use its value directly
				if len(returnValue.JSONData) == 1 {
					for _, value := range returnValue.JSONData {
						stateUpdateResponse = value
						e.Logger.LogDebug("Using single JSON field value for StateUpdate")
						break
					}
				} else {
					// Use the entire JSONData map
					stateUpdateResponse = returnValue.JSONData
					e.Logger.LogDebug(fmt.Sprintf("Using captured JSON data for StateUpdate: %d keys", len(returnValue.JSONData)))
				}
			}

			if updateErr := step.StateUpdate(stateUpdateResponse, variables); updateErr != nil {
				e.Logger.LogError("State update failed", updateErr)
			} else {
				e.Logger.LogStateUpdate(step.Name, variables)
				// Log variables after StateUpdate
				e.Logger.LogDebug(fmt.Sprintf("Variables after StateUpdate in step %s: %d variables", step.Name, len(variables)))
				for k := range variables {
					e.Logger.LogDebug(fmt.Sprintf("  - %s", k))
				}
			}
		}
	}

	e.Logger.LogStepResult(step.Name, step.ClientMethod, result.Status, result.StatusCode, result.Duration)

	return result
}

// storeResponseForStep stores response data for later snapshot capture
func (e *E2EEngineWithSnapshot) storeResponseForStep(stepName string, returnValue *SDKReturnValue) {
	// For now, we'll store this in the engine for later retrieval
	// In a more sophisticated implementation, this could be stored in a map
	// keyed by step name for retrieval during snapshot capture
}

// CaptureResponseForStep captures response data for a specific step
func (e *E2EEngineWithSnapshot) CaptureResponseForStep(stepName string, response interface{}) (*SDKReturnValue, error) {
	if !e.Config.EnableCapture {
		return nil, nil
	}

	return e.SnapshotCapture.CaptureStepResponse(stepName, response)
}

// UpdateStepSnapshotWithResponse updates a step snapshot with captured response data
func (e *E2EEngineWithSnapshot) UpdateStepSnapshotWithResponse(stepSnapshot *StepSnapshot, response interface{}) error {
	if !e.Config.EnableCapture || response == nil {
		return nil
	}

	returnValue, err := e.SnapshotCapture.CaptureStepResponse(stepSnapshot.StepName, response)
	if err != nil {
		return fmt.Errorf("failed to capture response for step %s: %w", stepSnapshot.StepName, err)
	}

	stepSnapshot.ReturnValue = returnValue
	return nil
}

// TrackStateChangeForStep tracks a state change for a specific step
func (e *E2EEngineWithSnapshot) TrackStateChangeForStep(stepSnapshot *StepSnapshot, key string, oldValue, newValue interface{}) {
	if e.Config.EnableCapture {
		e.SnapshotCapture.TrackStateChange(stepSnapshot, key, oldValue, newValue)
	}
}

// Conversion methods between testlib types and snapshot types

func (e *E2EEngineWithSnapshot) convertToSnapshotStory(story testlib.Story) Story {
	snapshotStory := Story{
		Name:        story.Name,
		Description: story.Description,
		Variables:   story.Variables,
		Setup:       story.Setup,
		Cleanup:     story.Cleanup,
		Steps:       make([]Step, len(story.Steps)),
	}

	for i, step := range story.Steps {
		snapshotStory.Steps[i] = e.convertToSnapshotStep(step)
	}

	return snapshotStory
}

func (e *E2EEngineWithSnapshot) convertToSnapshotStep(step testlib.Step) Step {
	return Step{
		Name:           step.Name,
		ClientMethod:   step.ClientMethod,
		Parameters:     step.Parameters,
		ExpectedStatus: step.ExpectedStatus,
		Skip:           step.Skip,
		SkipReason:     step.SkipReason,
		Validation:     step.Validation,
		StateUpdate:    step.StateUpdate,
	}
}

// Legacy conversion methods for backward compatibility (used by tests)
func (e *E2EEngineWithSnapshot) convertToTestlibStory(story Story) testlib.Story {
	testlibStory := testlib.Story{
		Name:        story.Name,
		Description: story.Description,
		Variables:   story.Variables,
		Setup:       story.Setup,
		Cleanup:     story.Cleanup,
		Steps:       make([]testlib.Step, len(story.Steps)),
	}

	for i, step := range story.Steps {
		testlibStory.Steps[i] = e.convertToTestlibStep(step)
	}

	return testlibStory
}

func (e *E2EEngineWithSnapshot) convertToTestlibStep(step Step) testlib.Step {
	return testlib.Step{
		Name:           step.Name,
		ClientMethod:   step.ClientMethod,
		Parameters:     step.Parameters,
		ExpectedStatus: step.ExpectedStatus,
		Skip:           step.Skip,
		SkipReason:     step.SkipReason,
		Validation:     step.Validation,
		StateUpdate:    step.StateUpdate,
	}
}

// Enhanced execution methods with response capture

// ExecuteClientMethodWithResponseCapture executes a client method and captures the response
// This method implements enhanced reflection-based method calling with dynamic argument building
// 要件 1.1, 1.2: リフレクションエラーなしでのAPIメソッド呼び出しと適切な引数数での実行
// 要件 4.1, 4.2: 後方互換性保証のためのフォールバック機構
func (e *E2EEngineWithSnapshot) ExecuteClientMethodWithResponseCapture(methodName string, parameters interface{}, variables map[string]interface{}) (interface{}, int, *SDKReturnValue, error) {
	e.Logger.LogDebug(fmt.Sprintf("Executing client method with response capture: %s", methodName))

	if e.E2EEngine.Config.DryRun {
		e.Logger.LogDebug("DRY RUN: Simulating method execution with response capture")
		return nil, 200, nil, nil
	}

	// If parameters is a function, call it to get the actual parameters
	if parameters != nil {
		paramValue := reflect.ValueOf(parameters)
		if paramValue.Kind() == reflect.Func {
			// Call the function with variables to get actual parameters
			results := paramValue.Call([]reflect.Value{reflect.ValueOf(variables)})
			if len(results) > 0 {
				parameters = results[0].Interface()
			}
		}
	}

	// Use reflection to call the method on the client
	clientValue := reflect.ValueOf(e.Client)
	method := clientValue.MethodByName(methodName)
	if !method.IsValid() {
		return nil, 0, nil, fmt.Errorf("method %s not found on client", methodName)
	}

	// Get method type for analysis
	methodType := method.Type()

	// メソッドシグネチャ解析を統合：詳細な解析情報をログ出力
	methodAnalyzer := NewMethodSignatureAnalyzer(methodType, methodName)
	strategy := methodAnalyzer.AnalyzeSignature()

	e.Logger.LogInfo(fmt.Sprintf("Method analysis for %s:", methodName))
	e.Logger.LogInfo(fmt.Sprintf("  Strategy: %s", strategy.String()))
	e.Logger.LogInfo(fmt.Sprintf("  Signature: %s", methodAnalyzer.GetMethodSignatureString()))
	e.Logger.LogInfo(fmt.Sprintf("  Expected args: %d", methodAnalyzer.GetExpectedArgumentCount()))
	e.Logger.LogInfo(fmt.Sprintf("  Is variadic: %t", methodAnalyzer.IsVariadic()))

	// Log method call attempt with detailed information
	e.ErrorHandler.LogMethodCallAttempt(methodName, methodType, parameters)

	// Try enhanced method calling with fallback to legacy approach
	// 要件 4.1: 既存のメソッド呼び出しパターンを維持するフォールバック処理
	response, statusCode, returnValue, err := e.executeMethodWithFallback(method, methodName, methodType, parameters)
	if err != nil {
		e.ErrorHandler.LogMethodCallFailure(methodName, err)
		return nil, 0, nil, err
	}

	// Log successful method call with detailed information
	responseType := reflect.TypeOf(response)
	e.ErrorHandler.LogMethodCallSuccess(methodName, statusCode, responseType)
	e.Logger.LogDebug(fmt.Sprintf("Method %s executed successfully with status code %d", methodName, statusCode))

	return response, statusCode, returnValue, nil
}

// executeMethodWithFallback implements the fallback mechanism for method execution
// 要件 4.1, 4.2: 後方互換性保証のためのフォールバック機構実装
func (e *E2EEngineWithSnapshot) executeMethodWithFallback(method reflect.Value, methodName string, methodType reflect.Type, parameters interface{}) (interface{}, int, *SDKReturnValue, error) {
	// Build context
	ctx := context.Background()
	if e.E2EEngine.Config.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(e.E2EEngine.Config.Timeout)*time.Second)
		defer cancel()
	}

	// Use the new fallback mechanism
	// 要件 4.1, 4.2: 後方互換性保証のためのフォールバック機構
	e.Logger.LogDebug(fmt.Sprintf("Using fallback mechanism for method: %s", methodName))

	results, err := e.FallbackMechanism.ExecuteWithFallback(method, methodName, methodType, ctx, parameters)
	if err != nil {
		e.Logger.LogError(fmt.Sprintf("Fallback mechanism failed for method %s", methodName), err)
		return nil, 0, nil, err
	}

	// Validate return values
	if len(results) < 2 {
		return nil, 0, nil, fmt.Errorf("unexpected number of return values from %s: expected at least 2, got %d", methodName, len(results))
	}

	// Check for error in return values (last return value should be error)
	errorResult := results[len(results)-1]
	if !errorResult.IsNil() {
		err := errorResult.Interface().(error)
		e.Logger.LogDebug(fmt.Sprintf("Method %s returned error: %v", methodName, err))
		return nil, 0, nil, err
	}

	// Extract response and status code
	response := results[0].Interface()
	statusCode := e.extractStatusCode(response)

	e.Logger.LogDebug(fmt.Sprintf("Method %s completed successfully with status code %d", methodName, statusCode))

	// Capture response data if enabled
	var returnValue *SDKReturnValue
	if e.Config.EnableCapture {
		var captureErr error
		returnValue, captureErr = e.SnapshotCapture.CaptureStepResponse(methodName, response)
		if captureErr != nil {
			e.Logger.LogError("Failed to capture response", captureErr)
			// Don't fail the entire method call if capture fails
		} else {
			e.Logger.LogDebug(fmt.Sprintf("Successfully captured response for method %s", methodName))
		}
	}

	return response, statusCode, returnValue, nil
}

// tryEnhancedMethodCall attempts to call the method using enhanced signature analysis
// 要件 1.1, 1.2: メソッドシグネチャ解析を統合し、動的な引数構築ロジックを組み込む
func (e *E2EEngineWithSnapshot) tryEnhancedMethodCall(method reflect.Value, methodName string, methodType reflect.Type, parameters interface{}) (interface{}, int, *SDKReturnValue, error) {
	// Validate method signature using the analyzer
	methodAnalyzer := NewMethodSignatureAnalyzer(methodType, methodName)
	if err := methodAnalyzer.ValidateMethodSignature(); err != nil {
		return nil, 0, nil, fmt.Errorf("invalid method signature for %s: %w", methodName, err)
	}

	// Analyze method signature to determine call strategy
	strategy := methodAnalyzer.AnalyzeSignature()

	e.Logger.LogDebug(fmt.Sprintf("Method %s analyzed as strategy: %s", methodName, strategy.String()))
	e.Logger.LogDebug(fmt.Sprintf("Method signature: %s", methodAnalyzer.GetMethodSignatureString()))

	// Build arguments using the appropriate strategy with method type integration
	var argumentBuilder ArgumentBuilder
	switch strategy {
	case WithBodyMethod:
		argumentBuilder = &WithBodyArgumentBuilder{}
		e.Logger.LogDebug(fmt.Sprintf("Using WithBodyArgumentBuilder for method %s", methodName))
	case NoParameterMethod:
		argumentBuilder = &NoParameterArgumentBuilder{}
		e.Logger.LogDebug(fmt.Sprintf("Using NoParameterArgumentBuilder for method %s", methodName))
	default:
		argumentBuilder = &StandardArgumentBuilder{}
		e.Logger.LogDebug(fmt.Sprintf("Using StandardArgumentBuilder for method %s", methodName))
	}

	// Build arguments with context from the E2E engine
	ctx := context.Background()
	if e.E2EEngine.Config.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(e.E2EEngine.Config.Timeout)*time.Second)
		defer cancel()
	}

	// 動的な引数構築ロジック：メソッド型を考慮した引数構築を優先的に試行
	var args []reflect.Value
	var err error

	// TypedArgumentBuilderインターフェースをサポートしているかチェック
	if typedBuilder, ok := argumentBuilder.(TypedArgumentBuilder); ok {
		e.Logger.LogDebug(fmt.Sprintf("Using typed argument builder for method %s", methodName))
		args, err = typedBuilder.BuildArgumentsWithMethodType(ctx, parameters, methodType)
		if err != nil {
			e.Logger.LogDebug(fmt.Sprintf("Typed argument building failed for %s: %v", methodName, err))
			// フォールバック：従来の方法を試行
			args, err = argumentBuilder.BuildArguments(ctx, parameters)
			if err != nil {
				return nil, 0, nil, fmt.Errorf("failed to build arguments (both typed and standard): %w", err)
			}
			e.Logger.LogDebug(fmt.Sprintf("Fallback to standard argument building succeeded for %s", methodName))
		} else {
			e.Logger.LogDebug(fmt.Sprintf("Typed argument building succeeded for %s", methodName))
		}
	} else {
		// 従来の方法のみサポート
		e.Logger.LogDebug(fmt.Sprintf("Using standard argument builder for method %s", methodName))
		args, err = argumentBuilder.BuildArguments(ctx, parameters)
		if err != nil {
			return nil, 0, nil, fmt.Errorf("failed to build arguments: %w", err)
		}
	}

	// 引数数の詳細ログ出力
	expectedArgCount := methodType.NumIn()
	actualArgCount := len(args)
	e.Logger.LogDebug(fmt.Sprintf("Method %s: expected %d arguments, built %d arguments",
		methodName, expectedArgCount, actualArgCount))

	// 各引数の型情報をログ出力
	for i, arg := range args {
		expectedType := methodType.In(i)
		actualType := arg.Type()
		e.Logger.LogDebug(fmt.Sprintf("  Arg %d: expected %s, actual %s",
			i, expectedType.String(), actualType.String()))
	}

	// Validate argument count matches method signature
	if actualArgCount != expectedArgCount {
		argMismatchErr := e.ErrorHandler.HandleArgumentMismatch(methodName, methodType, args)
		return nil, 0, nil, argMismatchErr
	}

	// Execute the method call
	return e.executeMethodCall(method, methodName, args)
}

// tryLegacyMethodCall attempts to call the method using the legacy pattern
// 要件 4.1: 既存のメソッド呼び出しパターンを維持するフォールバック処理
func (e *E2EEngineWithSnapshot) tryLegacyMethodCall(method reflect.Value, methodName string, methodType reflect.Type, parameters interface{}) (interface{}, int, *SDKReturnValue, error) {
	e.Logger.LogDebug(fmt.Sprintf("Trying legacy method call pattern for: %s", methodName))

	// Build context
	ctx := context.Background()
	if e.E2EEngine.Config.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(e.E2EEngine.Config.Timeout)*time.Second)
		defer cancel()
	}

	// Legacy pattern: ctx + parameters (most common pattern in existing tests)
	args := []reflect.Value{
		reflect.ValueOf(ctx),
	}

	// Add parameters if provided and method expects them
	if parameters != nil && methodType.NumIn() > 1 {
		args = append(args, reflect.ValueOf(parameters))
	}

	// Check if argument count matches
	if len(args) != methodType.NumIn() {
		return nil, 0, nil, fmt.Errorf("legacy pattern argument count mismatch: expected %d, got %d", methodType.NumIn(), len(args))
	}

	// Execute the method call
	return e.executeMethodCall(method, methodName, args)
}

// trySafeMethodCall attempts to call the method using safe patterns for unknown methods
// 要件 4.2: 未知のメソッドパターンに対する安全な処理
func (e *E2EEngineWithSnapshot) trySafeMethodCall(method reflect.Value, methodName string, methodType reflect.Type, parameters interface{}) (interface{}, int, *SDKReturnValue, error) {
	e.Logger.LogDebug(fmt.Sprintf("Trying safe method call pattern for: %s", methodName))

	// Build context
	ctx := context.Background()
	if e.E2EEngine.Config.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(e.E2EEngine.Config.Timeout)*time.Second)
		defer cancel()
	}

	// Safe pattern: Try different argument combinations based on method signature
	numIn := methodType.NumIn()

	// Pattern 1: Context only (for methods with no parameters)
	if numIn == 1 {
		args := []reflect.Value{reflect.ValueOf(ctx)}
		return e.executeMethodCall(method, methodName, args)
	}

	// Pattern 2: Context + parameters (standard pattern)
	if numIn == 2 && parameters != nil {
		args := []reflect.Value{
			reflect.ValueOf(ctx),
			reflect.ValueOf(parameters),
		}
		return e.executeMethodCall(method, methodName, args)
	}

	// Pattern 3: Try to match method signature with available arguments
	args := []reflect.Value{reflect.ValueOf(ctx)}

	// Add parameters if available and method expects more arguments
	if parameters != nil && numIn > 1 {
		// Try to add parameters as-is
		args = append(args, reflect.ValueOf(parameters))

		// If still not enough arguments, try to fill with zero values
		for len(args) < numIn {
			paramType := methodType.In(len(args))
			zeroValue := reflect.Zero(paramType)
			args = append(args, zeroValue)
		}
	} else {
		// Fill remaining arguments with zero values
		for len(args) < numIn {
			paramType := methodType.In(len(args))
			zeroValue := reflect.Zero(paramType)
			args = append(args, zeroValue)
		}
	}

	// Ensure we don't exceed the expected argument count
	if len(args) > numIn {
		args = args[:numIn]
	}

	return e.executeMethodCall(method, methodName, args)
}

// executeMethodCall executes the actual method call with proper error handling
// 動的な引数構築ロジックを組み込んだメソッド実行処理
func (e *E2EEngineWithSnapshot) executeMethodCall(method reflect.Value, methodName string, args []reflect.Value) (interface{}, int, *SDKReturnValue, error) {
	// Log detailed argument information
	e.Logger.LogDebug(fmt.Sprintf("Calling method %s with %d arguments", methodName, len(args)))
	for i, arg := range args {
		argType := arg.Type()
		argValue := "nil"
		if arg.IsValid() && arg.CanInterface() {
			argValue = fmt.Sprintf("%v", arg.Interface())
			// 長い値は切り詰める
			if len(argValue) > 100 {
				argValue = argValue[:100] + "..."
			}
		}
		e.Logger.LogDebug(fmt.Sprintf("  Arg %d: %s = %s", i, argType.String(), argValue))
	}

	// Call the method with proper error handling and panic recovery
	var results []reflect.Value
	var callErr error
	func() {
		defer func() {
			if r := recover(); r != nil {
				// Handle panic during method call with detailed information
				callErr = fmt.Errorf("panic during method call %s: %v", methodName, r)
				e.Logger.LogError(fmt.Sprintf("Panic occurred during method call %s", methodName), callErr)
			}
		}()

		// 動的な引数構築ロジック：メソッド型に基づいて適切な呼び出し方法を選択
		methodType := method.Type()

		// 可変長引数メソッドの処理を改善
		if methodType.IsVariadic() && len(args) > 0 {
			e.Logger.LogDebug(fmt.Sprintf("Method %s is variadic, using appropriate call method", methodName))

			// 最後の引数が可変長引数に対応するスライスかチェック
			lastArgIndex := len(args) - 1
			lastArg := args[lastArgIndex]
			variadicType := methodType.In(methodType.NumIn() - 1)

			e.Logger.LogDebug(fmt.Sprintf("Variadic type: %s, last arg type: %s",
				variadicType.String(), lastArg.Type().String()))

			// スライス型で、可変長引数の型と互換性がある場合はCallSliceを使用
			if lastArg.Kind() == reflect.Slice && lastArg.Type().AssignableTo(variadicType) {
				e.Logger.LogDebug(fmt.Sprintf("Using CallSlice for method %s", methodName))
				results = method.CallSlice(args)
				return
			}

			// 可変長引数を展開してCallを使用
			e.Logger.LogDebug(fmt.Sprintf("Using Call for variadic method %s", methodName))
		}

		// 通常のCall（非可変長引数メソッドまたはCallSliceが適用できない場合）
		e.Logger.LogDebug(fmt.Sprintf("Using regular Call for method %s", methodName))
		results = method.Call(args)
	}()

	// Check if panic occurred during method call
	if callErr != nil {
		return nil, 0, nil, callErr
	}

	// Validate return values
	if len(results) < 2 {
		return nil, 0, nil, fmt.Errorf("unexpected number of return values from %s: expected at least 2, got %d", methodName, len(results))
	}

	// Log return value information
	e.Logger.LogDebug(fmt.Sprintf("Method %s returned %d values", methodName, len(results)))
	for i, result := range results {
		resultType := result.Type()
		e.Logger.LogDebug(fmt.Sprintf("  Return %d: %s", i, resultType.String()))
	}

	// Check for error in return values (last return value should be error)
	errorResult := results[len(results)-1]
	if !errorResult.IsNil() {
		err := errorResult.Interface().(error)
		e.Logger.LogDebug(fmt.Sprintf("Method %s returned error: %v", methodName, err))
		return nil, 0, nil, err
	}

	// Extract response and status code
	response := results[0].Interface()
	statusCode := e.extractStatusCode(response)

	e.Logger.LogDebug(fmt.Sprintf("Method %s completed successfully with status code %d", methodName, statusCode))

	// Capture response data if enabled
	var returnValue *SDKReturnValue
	if e.Config.EnableCapture {
		var captureErr error
		returnValue, captureErr = e.SnapshotCapture.CaptureStepResponse(methodName, response)
		if captureErr != nil {
			e.Logger.LogError("Failed to capture response", captureErr)
			// Don't fail the entire method call if capture fails
		} else {
			e.Logger.LogDebug(fmt.Sprintf("Successfully captured response for method %s", methodName))
		}
	}

	return response, statusCode, returnValue, nil
}

// extractStatusCode extracts status code from response (copied from base engine)
func (e *E2EEngineWithSnapshot) extractStatusCode(response interface{}) int {
	if response == nil {
		return 0
	}

	if statusProvider, ok := response.(interface{ StatusCode() int }); ok {
		return statusProvider.StatusCode()
	}
	if httpResp, ok := response.(*http.Response); ok && httpResp != nil {
		return httpResp.StatusCode
	}

	respValue := reflect.ValueOf(response)
	if respValue.Kind() == reflect.Ptr {
		respValue = respValue.Elem()
	}
	statusCodeMethod := respValue.MethodByName("StatusCode")
	if statusCodeMethod.IsValid() {
		results := statusCodeMethod.Call(nil)
		if len(results) > 0 {
			if statusCode, ok := results[0].Interface().(int); ok {
				return statusCode
			}
		}
	}
	return http.StatusOK
}

// Separated execution modes (Requirements 2.1, 4.1)

// ExecuteSnapshotOnly executes stories and captures snapshots only (Phase 1)
// This method implements isolated snapshot capture without comparison or reporting
func (e *E2EEngineWithSnapshot) ExecuteSnapshotOnly(stories []testlib.Story) []StorySnapshot {
	e.Logger.LogInfo(fmt.Sprintf("Starting snapshot capture only mode with %d stories", len(stories)))

	// Temporarily enable only capture mode
	originalConfig := *e.Config
	e.Config.EnableCapture = true
	e.Config.EnableComparison = false
	e.Config.EnableReporting = false
	e.Config.EnableValidation = false
	defer func() { e.Config = &originalConfig }()

	var snapshots []StorySnapshot
	startTime := time.Now()

	for i, story := range stories {
		e.Logger.LogInfo(fmt.Sprintf("Capturing snapshot for story %d/%d: %s", i+1, len(stories), story.Name))

		// Execute story and capture snapshot
		storyResult := e.executeStoryWithResponseCapture(story)

		// Only capture snapshot if story execution was successful or if configured to capture failures
		if storyResult.Status == TestStatusFailed && !e.shouldCaptureFailedStories() {
			e.Logger.LogInfo(fmt.Sprintf("Skipping snapshot capture for failed story: %s", story.Name))
			continue
		}

		// Capture snapshot
		snapshot, err := e.captureStorySnapshot(story, storyResult)
		if err != nil {
			e.Logger.LogError(fmt.Sprintf("Failed to capture snapshot for story '%s'", story.Name), err)
			continue
		}

		// Save snapshot to file
		if err := e.saveSnapshotToFile(snapshot); err != nil {
			e.Logger.LogError(fmt.Sprintf("Failed to save snapshot for story '%s'", story.Name), err)
			continue
		}

		snapshots = append(snapshots, *snapshot)
		e.Logger.LogInfo(fmt.Sprintf("Successfully captured and saved snapshot for story: %s", story.Name))
	}

	totalDuration := time.Since(startTime)
	e.Logger.LogInfo(fmt.Sprintf("Completed snapshot capture only mode: %d snapshots captured in %v",
		len(snapshots), totalDuration))

	return snapshots
}

// CompareSnapshots compares existing snapshots (Phase 2)
// This method implements isolated snapshot comparison without E2E execution
func (e *E2EEngineWithSnapshot) CompareSnapshots(snapshotPaths []string) ([]StoryComparison, error) {
	e.Logger.LogInfo(fmt.Sprintf("Starting snapshot comparison mode with %d snapshot files", len(snapshotPaths)))

	var comparisons []StoryComparison
	startTime := time.Now()

	for i, snapshotPath := range snapshotPaths {
		e.Logger.LogInfo(fmt.Sprintf("Comparing snapshot %d/%d: %s", i+1, len(snapshotPaths), snapshotPath))

		// Load current snapshot
		currentSnapshot, err := e.FileManager.LoadStorySnapshot(snapshotPath)
		if err != nil {
			e.Logger.LogError(fmt.Sprintf("Failed to load snapshot from %s", snapshotPath), err)
			continue
		}

		// Compare with previous release
		comparison, err := e.Comparator.CompareWithPreviousRelease(currentSnapshot)
		if err != nil {
			e.Logger.LogError(fmt.Sprintf("Failed to compare snapshot for story '%s'", currentSnapshot.StoryName), err)
			continue
		}

		// Save comparison result
		if err := e.saveComparisonToFile(comparison); err != nil {
			e.Logger.LogError(fmt.Sprintf("Failed to save comparison for story '%s'", comparison.StoryName), err)
			continue
		}

		comparisons = append(comparisons, *comparison)
		e.Logger.LogInfo(fmt.Sprintf("Successfully compared and saved comparison for story: %s", comparison.StoryName))
	}

	totalDuration := time.Since(startTime)
	e.Logger.LogInfo(fmt.Sprintf("Completed snapshot comparison mode: %d comparisons completed in %v",
		len(comparisons), totalDuration))

	return comparisons, nil
}

// GenerateReports generates reports from comparison results (Phase 3)
// This method implements isolated report generation from existing comparison files
func (e *E2EEngineWithSnapshot) GenerateReports(comparisonPaths []string) ([]StoryReport, error) {
	e.Logger.LogInfo(fmt.Sprintf("Starting report generation mode with %d comparison files", len(comparisonPaths)))

	var reports []StoryReport
	startTime := time.Now()

	for i, comparisonPath := range comparisonPaths {
		e.Logger.LogInfo(fmt.Sprintf("Generating report %d/%d: %s", i+1, len(comparisonPaths), comparisonPath))

		// Load comparison result
		comparison, err := e.FileManager.LoadStoryComparison(comparisonPath)
		if err != nil {
			e.Logger.LogError(fmt.Sprintf("Failed to load comparison from %s", comparisonPath), err)
			continue
		}

		// Generate report
		report, err := e.Reporter.GenerateStoryReport(comparison)
		if err != nil {
			e.Logger.LogError(fmt.Sprintf("Failed to generate report for story '%s'", comparison.StoryName), err)
			continue
		}

		// Save report
		if err := e.saveReportToFile(report); err != nil {
			e.Logger.LogError(fmt.Sprintf("Failed to save report for story '%s'", report.StoryName), err)
			continue
		}

		reports = append(reports, *report)
		e.Logger.LogInfo(fmt.Sprintf("Successfully generated and saved report for story: %s", report.StoryName))
	}

	totalDuration := time.Since(startTime)
	e.Logger.LogInfo(fmt.Sprintf("Completed report generation mode: %d reports generated in %v",
		len(reports), totalDuration))

	return reports, nil
}

// Configuration management methods

// GetSnapshotConfig returns the current snapshot configuration
func (e *E2EEngineWithSnapshot) GetSnapshotConfig() *StorySnapshotConfig {
	return e.Config
}

// UpdateSnapshotConfig updates the snapshot configuration
func (e *E2EEngineWithSnapshot) UpdateSnapshotConfig(config *StorySnapshotConfig) {
	e.Config = config
	e.SnapshotCapture = NewStorySnapshotCapture(config)
}

// EnableSnapshotCapture enables or disables snapshot capture
func (e *E2EEngineWithSnapshot) EnableSnapshotCapture(enabled bool) {
	e.Config.EnableCapture = enabled
}

// EnableSnapshotComparison enables or disables snapshot comparison
func (e *E2EEngineWithSnapshot) EnableSnapshotComparison(enabled bool) {
	e.Config.EnableComparison = enabled
}

// EnableSnapshotReporting enables or disables snapshot reporting
func (e *E2EEngineWithSnapshot) EnableSnapshotReporting(enabled bool) {
	e.Config.EnableReporting = enabled
}

// SetSnapshotOnly sets the engine to snapshot-only mode
func (e *E2EEngineWithSnapshot) SetSnapshotOnly(snapshotOnly bool) {
	e.Config.SnapshotOnly = snapshotOnly
	if snapshotOnly {
		e.Config.EnableCapture = true
		e.Config.EnableComparison = false
		e.Config.EnableReporting = false
	}
}

// GenerateReportsFromComparisons generates reports from provided comparison results
func (e *E2EEngineWithSnapshot) GenerateReportsFromComparisons(comparisons []StoryComparison) []StoryReport {
	var reports []StoryReport

	e.Logger.LogInfo("Generating reports from comparison results")

	for _, comparison := range comparisons {
		report, err := e.Reporter.GenerateStoryReport(&comparison)
		if err != nil {
			e.Logger.LogError(fmt.Sprintf("Failed to generate report for '%s'", comparison.StoryName), err)
			continue
		}

		// Save report
		err = e.FileManager.SaveStoryReport(report)
		if err != nil {
			e.Logger.LogError(fmt.Sprintf("Failed to save report for '%s'", comparison.StoryName), err)
		} else {
			e.Logger.LogInfo(fmt.Sprintf("Saved report for '%s'", comparison.StoryName))
		}

		reports = append(reports, *report)
	}

	return reports
}

// CLI-based execution control methods (Requirements 2.1, 4.1)

// ExecuteWithCLIConfig executes based on CLI configuration
// This method provides unified execution control for different phases
func (e *E2EEngineWithSnapshot) ExecuteWithCLIConfig(cliConfig *CLIConfig, stories []testlib.Story) error {
	e.Logger.LogInfo(fmt.Sprintf("Starting execution with CLI config: mode=%s, service=%s",
		cliConfig.Mode, cliConfig.Service))

	// Apply CLI config overrides
	if err := e.applyCLIConfigOverrides(cliConfig); err != nil {
		return fmt.Errorf("failed to apply CLI config overrides: %w", err)
	}

	// Execute based on mode
	switch ExecutionMode(cliConfig.Mode) {
	case ExecutionModeFull:
		return e.executeFullMode(cliConfig, stories)
	case ExecutionModeCapture:
		return e.executeCaptureMode(cliConfig, stories)
	case ExecutionModeCompare:
		return e.executeCompareMode(cliConfig)
	case ExecutionModeReport:
		return e.executeReportMode(cliConfig)
	default:
		return fmt.Errorf("unknown execution mode: %s", cliConfig.Mode)
	}
}

// executeFullMode executes all phases (capture, compare, report)
func (e *E2EEngineWithSnapshot) executeFullMode(cliConfig *CLIConfig, stories []testlib.Story) error {
	e.Logger.LogInfo("Executing full mode (all phases)")

	// Filter stories if specific stories are requested
	filteredStories := e.filterStoriesByNames(stories, cliConfig.Stories)

	// Execute integrated mode
	results := e.ExecuteStoriesWithSnapshot(filteredStories)

	// Log summary
	e.Logger.LogInfo(fmt.Sprintf("Full mode completed: %d stories processed", len(results)))
	return nil
}

// executeCaptureMode executes only snapshot capture phase
func (e *E2EEngineWithSnapshot) executeCaptureMode(cliConfig *CLIConfig, stories []testlib.Story) error {
	e.Logger.LogInfo("Executing capture mode (snapshot capture only)")

	// Filter stories if specific stories are requested
	filteredStories := e.filterStoriesByNames(stories, cliConfig.Stories)

	// Execute snapshot capture only
	snapshots := e.ExecuteSnapshotOnly(filteredStories)

	// Log summary
	e.Logger.LogInfo(fmt.Sprintf("Capture mode completed: %d snapshots captured", len(snapshots)))
	return nil
}

// executeCompareMode executes only snapshot comparison phase
func (e *E2EEngineWithSnapshot) executeCompareMode(cliConfig *CLIConfig) error {
	e.Logger.LogInfo("Executing compare mode (snapshot comparison only)")

	// Determine snapshot files to compare
	snapshotPaths, err := e.determineSnapshotPaths(cliConfig)
	if err != nil {
		return fmt.Errorf("failed to determine snapshot paths: %w", err)
	}

	// Execute comparison
	comparisons, err := e.CompareSnapshots(snapshotPaths)
	if err != nil {
		return fmt.Errorf("failed to compare snapshots: %w", err)
	}

	// Log summary
	e.Logger.LogInfo(fmt.Sprintf("Compare mode completed: %d comparisons completed", len(comparisons)))
	return nil
}

// executeReportMode executes only report generation phase
func (e *E2EEngineWithSnapshot) executeReportMode(cliConfig *CLIConfig) error {
	e.Logger.LogInfo("Executing report mode (report generation only)")

	// Determine comparison files to process
	comparisonPaths, err := e.determineComparisonPaths(cliConfig)
	if err != nil {
		return fmt.Errorf("failed to determine comparison paths: %w", err)
	}

	// Execute report generation
	reports, err := e.GenerateReports(comparisonPaths)
	if err != nil {
		return fmt.Errorf("failed to generate reports: %w", err)
	}

	// Log summary
	e.Logger.LogInfo(fmt.Sprintf("Report mode completed: %d reports generated", len(reports)))
	return nil
}

// applyCLIConfigOverrides applies CLI configuration overrides to the engine config
func (e *E2EEngineWithSnapshot) applyCLIConfigOverrides(cliConfig *CLIConfig) error {
	// Override output directory if specified
	if cliConfig.OutputDir != "" {
		e.Config.OutputDirectory = cliConfig.OutputDir
		// Recreate file manager with module-scoped output directory
		e.FileManager = NewStorySnapshotFileManager(e.Config.GetModuleOutputDirectory())
	}

	// Apply dry run mode
	if cliConfig.DryRun {
		e.E2EEngine.Config.DryRun = true
	}

	// Apply verbose logging
	if cliConfig.Verbose {
		// Enable debug logging if available
		e.Logger.LogInfo("Verbose logging enabled")
	}

	return nil
}

// filterStoriesByNames filters stories by specified names
func (e *E2EEngineWithSnapshot) filterStoriesByNames(stories []testlib.Story, storyNames []string) []testlib.Story {
	if len(storyNames) == 0 {
		return stories // Return all stories if no filter specified
	}

	nameSet := make(map[string]bool)
	for _, name := range storyNames {
		nameSet[name] = true
	}

	var filteredStories []testlib.Story
	for _, story := range stories {
		if nameSet[story.Name] {
			filteredStories = append(filteredStories, story)
		}
	}

	e.Logger.LogInfo(fmt.Sprintf("Filtered stories: %d out of %d stories selected",
		len(filteredStories), len(stories)))

	return filteredStories
}

// determineSnapshotPaths determines which snapshot files to compare
func (e *E2EEngineWithSnapshot) determineSnapshotPaths(cliConfig *CLIConfig) ([]string, error) {
	var snapshotPaths []string

	// If specific snapshot file is provided
	if cliConfig.SnapshotFile != "" {
		snapshotPaths = append(snapshotPaths, cliConfig.SnapshotFile)
		return snapshotPaths, nil
	}

	// If specific stories are requested, find their snapshot files
	if len(cliConfig.Stories) > 0 {
		for _, storyName := range cliConfig.Stories {
			// Determine current tag
			currentTag, err := GetGitTagForFileManager()
			if err != nil {
				return nil, fmt.Errorf("failed to get current git tag: %w", err)
			}

			// Build snapshot file path
			fileName := e.Config.GetSnapshotFileName(currentTag, storyName)
			filePath := filepath.Join(e.Config.GetSnapshotDirectory(), fileName)
			snapshotPaths = append(snapshotPaths, filePath)
		}
		return snapshotPaths, nil
	}

	// Otherwise, find all snapshot files for the current tag
	currentTag, err := GetGitTagForFileManager()
	if err != nil {
		return nil, fmt.Errorf("failed to get current git tag: %w", err)
	}

	// Find all snapshot files with current tag
	snapshotDir := e.Config.GetSnapshotDirectory()
	files, err := filepath.Glob(filepath.Join(snapshotDir, fmt.Sprintf("story_snapshot_%s_*.json", currentTag)))
	if err != nil {
		return nil, fmt.Errorf("failed to find snapshot files: %w", err)
	}

	return files, nil
}

// determineComparisonPaths determines which comparison files to process for reports
func (e *E2EEngineWithSnapshot) determineComparisonPaths(cliConfig *CLIConfig) ([]string, error) {
	var comparisonPaths []string

	// If specific comparison file is provided
	if cliConfig.ComparisonFile != "" {
		comparisonPaths = append(comparisonPaths, cliConfig.ComparisonFile)
		return comparisonPaths, nil
	}

	// If specific stories are requested, find their comparison files
	if len(cliConfig.Stories) > 0 {
		for _, storyName := range cliConfig.Stories {
			// Build comparison file path (using release comparison by default)
			fileName := e.Config.GetComparisonFileName(storyName, "", "")
			filePath := filepath.Join(e.Config.GetComparisonDirectory(), fileName)
			comparisonPaths = append(comparisonPaths, filePath)
		}
		return comparisonPaths, nil
	}

	// Otherwise, find all comparison files
	comparisonDir := e.Config.GetComparisonDirectory()
	files, err := filepath.Glob(filepath.Join(comparisonDir, "story_comparison_*.json"))
	if err != nil {
		return nil, fmt.Errorf("failed to find comparison files: %w", err)
	}

	return files, nil
}

// Helper methods for integrated execution

// loadExistingSnapshot loads an existing snapshot for a story
func (e *E2EEngineWithSnapshot) loadExistingSnapshot(storyName string) (*StorySnapshot, error) {
	// Get current tag
	currentTag, err := GetGitTagForFileManager()
	if err != nil {
		return nil, fmt.Errorf("failed to get current git tag: %w", err)
	}

	// Build snapshot file path
	fileName := e.Config.GetSnapshotFileName(currentTag, storyName)
	filePath := filepath.Join(e.Config.GetSnapshotDirectory(), fileName)

	// Load snapshot
	snapshot, err := e.FileManager.LoadStorySnapshot(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load snapshot from %s: %w", filePath, err)
	}

	e.Logger.LogInfo(fmt.Sprintf("Loaded existing snapshot for story: %s", storyName))
	return snapshot, nil
}

// loadExistingComparison loads an existing comparison for a story
func (e *E2EEngineWithSnapshot) loadExistingComparison(storyName string) (*StoryComparison, error) {
	// Sanitize story name for file matching
	sanitizedName := sanitizeStoryNameForComparison(storyName)

	// Try to find the most recent comparison file for this story
	comparisonDir := e.Config.GetComparisonDirectory()
	pattern := filepath.Join(comparisonDir, fmt.Sprintf("story_comparison_%s_*.json", sanitizedName))

	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to search for comparison files: %w", err)
	}

	if len(files) == 0 {
		// Try release comparison as fallback
		fileName := e.Config.GetComparisonFileName(storyName, "", "")
		filePath := filepath.Join(comparisonDir, fileName)
		return nil, fmt.Errorf("no comparison files found for story '%s' (tried pattern: %s and %s)", storyName, pattern, filePath)
	}

	// Use the most recent file (last in sorted order)
	filePath := files[len(files)-1]

	// Load comparison
	comparison, err := e.FileManager.LoadStoryComparison(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load comparison from %s: %w", filePath, err)
	}

	e.Logger.LogInfo(fmt.Sprintf("Loaded existing comparison for story: %s from %s", storyName, filepath.Base(filePath)))
	return comparison, nil
}

// sanitizeStoryNameForComparison sanitizes story name for file matching
func sanitizeStoryNameForComparison(storyName string) string {
	// Use the same sanitization as file_manager.go
	sanitized := strings.ReplaceAll(storyName, " ", "_")
	sanitized = strings.ReplaceAll(sanitized, "-", "_")
	sanitized = strings.ReplaceAll(sanitized, ".", "_")
	sanitized = strings.ReplaceAll(sanitized, "/", "_")
	sanitized = strings.ReplaceAll(sanitized, "\\", "_")
	sanitized = strings.ToLower(sanitized)

	// Remove multiple consecutive underscores
	for strings.Contains(sanitized, "__") {
		sanitized = strings.ReplaceAll(sanitized, "__", "_")
	}

	// Trim underscores from start and end
	sanitized = strings.Trim(sanitized, "_")
	return sanitized
}

// captureStorySnapshot captures a snapshot for a story
func (e *E2EEngineWithSnapshot) captureStorySnapshot(story testlib.Story, storyResult StoryResult) (*StorySnapshot, error) {
	// Convert testlib.Story to snapshot.Story for snapshot capture
	snapshotStory := e.convertToSnapshotStory(story)

	snapshot, err := e.SnapshotCapture.CaptureStorySnapshot(snapshotStory, storyResult)
	if err != nil {
		return nil, fmt.Errorf("failed to capture story snapshot: %w", err)
	}

	e.Logger.LogInfo(fmt.Sprintf("Successfully captured snapshot for story: %s", story.Name))
	return snapshot, nil
}

// compareWithPreviousRelease compares snapshot with previous release
func (e *E2EEngineWithSnapshot) compareWithPreviousRelease(snapshot *StorySnapshot) (*StoryComparison, error) {
	comparison, err := e.Comparator.CompareWithPreviousRelease(snapshot)
	if err != nil {
		return nil, fmt.Errorf("failed to compare with previous release: %w", err)
	}

	e.Logger.LogInfo(fmt.Sprintf("Successfully compared snapshot for story: %s", snapshot.StoryName))
	return comparison, nil
}

// generateStoryReport generates a report from comparison results
func (e *E2EEngineWithSnapshot) generateStoryReport(comparison *StoryComparison) (*StoryReport, error) {
	report, err := e.Reporter.GenerateStoryReport(comparison)
	if err != nil {
		return nil, fmt.Errorf("failed to generate story report: %w", err)
	}

	e.Logger.LogInfo(fmt.Sprintf("Successfully generated report for story: %s", comparison.StoryName))
	return report, nil
}

// validateStorySnapshot validates a story snapshot
func (e *E2EEngineWithSnapshot) validateStorySnapshot(snapshot *StorySnapshot) (*StoryValidation, error) {
	// Create validator if not exists
	if e.Validator == nil {
		e.Validator = NewSimpleStoryValidator(e.Config)
	}

	validation, err := e.Validator.ValidateStorySnapshot(snapshot)
	if err != nil {
		return nil, fmt.Errorf("failed to validate story snapshot: %w", err)
	}

	e.Logger.LogInfo(fmt.Sprintf("Successfully validated snapshot for story: %s (Valid: %v)",
		snapshot.StoryName, validation.IsValid))
	return validation, nil
}

// File saving helper methods

// saveSnapshotToFile saves a snapshot to file
func (e *E2EEngineWithSnapshot) saveSnapshotToFile(snapshot *StorySnapshot) error {
	if e.E2EEngine.Config.DryRun {
		e.Logger.LogInfo("DRY RUN: Would save snapshot to file")
		return nil
	}

	filePath, err := e.FileManager.SaveStorySnapshot(snapshot)
	if err != nil {
		return fmt.Errorf("failed to save snapshot to file: %w", err)
	}
	e.Logger.LogInfo(fmt.Sprintf("Saved snapshot to: %s", filePath))
	return nil
}

// saveComparisonToFile saves a comparison result to file
func (e *E2EEngineWithSnapshot) saveComparisonToFile(comparison *StoryComparison) error {
	if e.E2EEngine.Config.DryRun {
		e.Logger.LogInfo("DRY RUN: Would save comparison to file")
		return nil
	}

	filePath, err := e.FileManager.SaveStoryComparison(comparison)
	if err != nil {
		return fmt.Errorf("failed to save comparison to file: %w", err)
	}
	e.Logger.LogInfo(fmt.Sprintf("Saved comparison to: %s", filePath))
	return nil
}

// saveReportToFile saves a report to file
func (e *E2EEngineWithSnapshot) saveReportToFile(report *StoryReport) error {
	if e.E2EEngine.Config.DryRun {
		e.Logger.LogInfo("DRY RUN: Would save report to file")
		return nil
	}

	err := e.FileManager.SaveStoryReport(report)
	if err != nil {
		return fmt.Errorf("failed to save report to file: %w", err)
	}
	e.Logger.LogInfo(fmt.Sprintf("Saved report for story: %s", report.StoryName))
	return nil
}

// saveValidationToFile saves a validation result to file
func (e *E2EEngineWithSnapshot) saveValidationToFile(validation *StoryValidation) error {
	if e.E2EEngine.Config.DryRun {
		e.Logger.LogInfo("DRY RUN: Would save validation to file")
		return nil
	}

	err := e.FileManager.SaveStoryValidation(validation)
	if err != nil {
		return fmt.Errorf("failed to save validation to file: %w", err)
	}
	e.Logger.LogInfo(fmt.Sprintf("Saved validation for story: %s", validation.StoryName))
	return nil
}

// Configuration helper methods

// shouldCaptureFailedStories determines if failed stories should be captured
func (e *E2EEngineWithSnapshot) shouldCaptureFailedStories() bool {
	// For now, always capture failed stories for debugging purposes
	// This could be made configurable in the future
	return true
}

// Logging helper methods

// logStoryResultSummary logs a summary of the story result
func (e *E2EEngineWithSnapshot) logStoryResultSummary(result StorySnapshotResult) {
	status := result.StoryResult.Status
	storyName := result.StoryResult.StoryName

	e.Logger.LogInfo(fmt.Sprintf("Story '%s' completed with status: %s", storyName, status))

	if result.Snapshot.StoryName != "" {
		e.Logger.LogInfo(fmt.Sprintf("  - Snapshot captured: %d steps", len(result.Snapshot.Steps)))
	}

	if result.Comparison != nil {
		e.Logger.LogInfo(fmt.Sprintf("  - Comparison completed: %d differences found",
			result.Comparison.Summary.TotalDifferences))
	}

	if result.Report != nil {
		e.Logger.LogInfo(fmt.Sprintf("  - Report generated: %s", result.Report.Title))
	}

	if result.Validation != nil {
		e.Logger.LogInfo(fmt.Sprintf("  - Validation completed: Valid=%v", result.Validation.IsValid))
	}
}

// logExecutionSummary logs a summary of the entire execution
func (e *E2EEngineWithSnapshot) logExecutionSummary(results []StorySnapshotResult) {
	totalStories := len(results)
	successfulStories := 0
	snapshotsCaptured := 0
	comparisonsCompleted := 0
	reportsGenerated := 0
	validationsCompleted := 0

	for _, result := range results {
		if result.StoryResult.Status == TestStatusPassed {
			successfulStories++
		}
		if result.Snapshot.StoryName != "" {
			snapshotsCaptured++
		}
		if result.Comparison != nil {
			comparisonsCompleted++
		}
		if result.Report != nil {
			reportsGenerated++
		}
		if result.Validation != nil {
			validationsCompleted++
		}
	}

	e.Logger.LogInfo("=== Execution Summary ===")
	e.Logger.LogInfo(fmt.Sprintf("Total stories: %d", totalStories))
	e.Logger.LogInfo(fmt.Sprintf("Successful stories: %d", successfulStories))
	e.Logger.LogInfo(fmt.Sprintf("Snapshots captured: %d", snapshotsCaptured))
	e.Logger.LogInfo(fmt.Sprintf("Comparisons completed: %d", comparisonsCompleted))
	e.Logger.LogInfo(fmt.Sprintf("Reports generated: %d", reportsGenerated))
	e.Logger.LogInfo(fmt.Sprintf("Validations completed: %d", validationsCompleted))
}

// requiresFallback determines if a method requires fallback mechanism
// 要件 4.1, 4.2: 既存のメソッド呼び出しパターンを維持するフォールバック処理
func (e *E2EEngineWithSnapshot) requiresFallback(methodName string, methodType reflect.Type, parameters interface{}) bool {
	// WithBody methods always require fallback
	if strings.Contains(methodName, "WithBody") {
		return true
	}

	// Methods with unusual parameter counts may require fallback
	numIn := methodType.NumIn()
	if numIn > 4 || numIn == 0 {
		return true
	}

	// Check if parameters match expected pattern
	if parameters != nil && numIn == 1 {
		// Method expects only context but parameters are provided
		return true
	}

	if parameters == nil && numIn > 1 {
		// Method expects parameters but none provided
		return true
	}

	return false
}

// BackwardCompatibilityReport represents a simple compatibility report
type BackwardCompatibilityReport struct {
	TotalMethods            int
	CompatibleMethods       int
	CompatibilityPercentage float64
	FallbackRequired        []string
}

// GetSummary returns a summary of the compatibility report
func (r *BackwardCompatibilityReport) GetSummary() string {
	return fmt.Sprintf("Compatibility: %.1f%% (%d/%d methods compatible), %d methods require fallback",
		r.CompatibilityPercentage, r.CompatibleMethods, r.TotalMethods, len(r.FallbackRequired))
}

// performCompatibilityCheck performs backward compatibility check on all stories
// 要件 4.1, 4.2: 既存テストケースが影響を受けないことを保証する機能
func (e *E2EEngineWithSnapshot) performCompatibilityCheck(stories []testlib.Story) *BackwardCompatibilityReport {
	e.Logger.LogInfo("Starting backward compatibility check")

	report := &BackwardCompatibilityReport{
		TotalMethods:      0,
		CompatibleMethods: 0,
		FallbackRequired:  []string{},
	}

	clientValue := reflect.ValueOf(e.Client)

	for _, story := range stories {
		for _, step := range story.Steps {
			report.TotalMethods++

			// Check if method exists and is compatible
			method := clientValue.MethodByName(step.ClientMethod)
			if !method.IsValid() {
				e.Logger.LogWarning(fmt.Sprintf("Method %s not found on client", step.ClientMethod))
				continue
			}

			methodType := method.Type()

			// Check if fallback is required
			if e.requiresFallback(step.ClientMethod, methodType, step.Parameters) {
				report.FallbackRequired = append(report.FallbackRequired, step.ClientMethod)
				e.Logger.LogInfo(fmt.Sprintf("Method %s will use fallback mechanism", step.ClientMethod))
			} else {
				report.CompatibleMethods++
			}
		}
	}

	if report.TotalMethods > 0 {
		report.CompatibilityPercentage = float64(report.CompatibleMethods) / float64(report.TotalMethods) * 100
	}

	e.Logger.LogInfo(fmt.Sprintf("Compatibility check completed: %s", report.GetSummary()))

	return report
}

// logCompatibilityReport logs the compatibility report details
func (e *E2EEngineWithSnapshot) logCompatibilityReport(report *BackwardCompatibilityReport) {
	e.Logger.LogInfo("=== Backward Compatibility Report ===")
	e.Logger.LogInfo(fmt.Sprintf("Total methods checked: %d", report.TotalMethods))
	e.Logger.LogInfo(fmt.Sprintf("Compatible methods: %d (%.1f%%)", report.CompatibleMethods, report.CompatibilityPercentage))

	if len(report.FallbackRequired) > 0 {
		e.Logger.LogInfo(fmt.Sprintf("Methods requiring fallback: %d", len(report.FallbackRequired)))
		for _, method := range report.FallbackRequired {
			e.Logger.LogInfo(fmt.Sprintf("  - %s", method))
		}
	}

	e.Logger.LogInfo("=== End Compatibility Report ===")
}

// EnableCompatibilityCheck enables or disables backward compatibility checking
// 要件 4.1: 既存のメソッド呼び出しパターンを維持するフォールバック処理
func (e *E2EEngineWithSnapshot) EnableCompatibilityCheck(enabled bool) {
	if e.Config.EnableCompatibilityCheck != enabled {
		e.Config.EnableCompatibilityCheck = enabled
		if enabled {
			e.Logger.LogInfo("Backward compatibility checking enabled")
		} else {
			e.Logger.LogInfo("Backward compatibility checking disabled")
		}
	}
}

// ValidateStoriesCompatibility は、ストーリーの互換性を検証します
// 要件 4.1, 4.2: 既存テストケースが影響を受けないことを保証する機能
func (e *E2EEngineWithSnapshot) ValidateStoriesCompatibility(stories []testlib.Story, minThreshold float64) error {
	if !e.Config.EnableCompatibilityCheck {
		e.Logger.LogInfo("Compatibility checking is disabled, skipping validation")
		return nil
	}

	e.Logger.LogInfo(fmt.Sprintf("Validating compatibility for %d stories (threshold: %.1f%%)", len(stories), minThreshold))

	return e.CompatibilityChecker.ValidateTestCaseCompatibility(stories, minThreshold)
}

// GetCompatibilityReport は、ストーリーの互換性レポートを取得します
// 要件 4.1, 4.2: 既存テストケースが影響を受けないことを保証する機能
func (e *E2EEngineWithSnapshot) GetCompatibilityReport(stories []testlib.Story) *CompatibilityCheckResult {
	e.Logger.LogInfo(fmt.Sprintf("Generating compatibility report for %d stories", len(stories)))

	result := e.CompatibilityChecker.CheckStoriesCompatibility(stories)

	// レポートを出力
	if e.Logger != nil {
		e.CompatibilityChecker.PrintCompatibilityReport(result)
	}

	return result
}

// GetIncompatibleMethods は、非互換メソッドのリストを取得します
func (e *E2EEngineWithSnapshot) GetIncompatibleMethods(stories []testlib.Story) []string {
	return e.CompatibilityChecker.GetIncompatibleMethods(stories)
}

// GetRequiredFallbackStrategies は、必要なフォールバック戦略のリストを取得します
func (e *E2EEngineWithSnapshot) GetRequiredFallbackStrategies(stories []testlib.Story) map[string]int {
	return e.CompatibilityChecker.GetRequiredStrategies(stories)
}

// TestFallbackMechanism は、フォールバック機構をテストします
// 要件 4.1, 4.2: 後方互換性保証のためのフォールバック機構のテスト
func (e *E2EEngineWithSnapshot) TestFallbackMechanism(methodNames []string) *FallbackCompatibilityReport {
	e.Logger.LogInfo(fmt.Sprintf("Testing fallback mechanism for %d methods", len(methodNames)))

	return e.FallbackMechanism.GetCompatibilityReport(e.Client, methodNames)
}

// ExecuteWithCompatibilityValidation は、互換性検証付きでストーリーを実行します
// 要件 4.1, 4.2: 既存テストケースが影響を受けないことを保証する機能
func (e *E2EEngineWithSnapshot) ExecuteWithCompatibilityValidation(stories []testlib.Story, minThreshold float64) ([]StorySnapshotResult, error) {
	// 実行前に互換性を検証
	if err := e.ValidateStoriesCompatibility(stories, minThreshold); err != nil {
		e.Logger.LogError("Compatibility validation failed before execution", err)
		return nil, fmt.Errorf("compatibility validation failed: %w", err)
	}

	e.Logger.LogInfo("Compatibility validation passed, proceeding with story execution")

	// 通常の実行を行う
	results := e.ExecuteStoriesWithSnapshot(stories)

	// 実行後の検証（オプション）
	e.Logger.LogInfo("Story execution completed, performing post-execution compatibility check")
	postExecutionReport := e.GetCompatibilityReport(stories)

	if postExecutionReport.CompatibilityPercentage < minThreshold {
		e.Logger.LogWarning(fmt.Sprintf("Post-execution compatibility dropped to %.1f%% (below threshold %.1f%%)",
			postExecutionReport.CompatibilityPercentage, minThreshold))
	}

	return results, nil
}
