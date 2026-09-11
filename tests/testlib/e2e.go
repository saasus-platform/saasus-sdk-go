package testlib

import (
	"fmt"
	"reflect"
	"time"
)

// E2EEngine provides the core E2E test execution engine
type E2EEngine struct {
	Client   any
	Config   *Config
	Logger   *Logger
	Coverage *CoverageTracker
	Reporter *Reporter
}

// NewE2EEngine creates a new E2E test engine
func NewE2EEngine(client any, methods []string) *E2EEngine {
	config := NewConfig()
	config.LoadFromEnvironment()

	logger := NewLogger(config.LogLevel)
	coverage := NewCoverageTracker(methods)
	reporter := NewReporter(coverage)

	return &E2EEngine{
		Client:   client,
		Config:   config,
		Logger:   logger,
		Coverage: coverage,
		Reporter: reporter,
	}
}

// ExecuteStories executes a list of test stories
func (e *E2EEngine) ExecuteStories(stories []Story) []StoryResult {
	var results []StoryResult

	e.Logger.LogInfo("Starting E2E test execution")

	for _, story := range stories {
		if e.Config.DryRun {
			e.Logger.LogInfo(fmt.Sprintf("DRY RUN: Would execute story '%s'", story.Name))
			// In dry run, still record method coverage for analysis
			for _, step := range story.Steps {
				e.Coverage.RecordExecution(step.ClientMethod, story.Name, step.Name, 200, 0, true, "")
			}
			continue
		}

		result := e.ExecuteStory(story)
		results = append(results, result)
	}

	return results
}

// ExecuteStory executes a single test story
func (e *E2EEngine) ExecuteStory(story Story) StoryResult {
	startTime := time.Now()

	e.Logger.LogStoryStart(story.Name)

	result := StoryResult{
		StoryName: story.Name,
		Status:    TestStatusPassed,
		Variables: make(map[string]any),
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

	// Execute steps
	for i, step := range story.Steps {
		stepResult := e.ExecuteStep(step, i+1, result.Variables)
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

// ExecuteStep executes a single test step
func (e *E2EEngine) ExecuteStep(step Step, stepNum int, variables map[string]any) StepResult {
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
		e.Logger.LogInfo(fmt.Sprintf("⏭️  Skipping step '%s' (%s)", step.Name, step.ClientMethod))
		e.Logger.LogStepResult(step.Name, step.ClientMethod, result.Status, result.StatusCode, result.Duration)
		return result
	}

	// Execute the client method and get both response and status code
	execResult := e.executeClientMethodWithResponse(step.ClientMethod, step.Parameters, variables)
	response := execResult.Response
	statusCode := execResult.StatusCode
	err := execResult.Error
	bodyBytes := execResult.BodyBytes

	result.StatusCode = statusCode
	result.Duration = time.Since(startTime)

	success := false
	errMsg := ""

	if err != nil {
		result.Status = TestStatusFailed
		result.Error = err
		errMsg = err.Error()
		e.Logger.LogError(fmt.Sprintf("Step '%s' failed", step.Name), err)
		// Log response body if available for debugging
		e.logResponseBody(response)
	} else {
		statusErr := ValidateExpectedStatus(step, statusCode)
		success = statusErr == nil // 外側のsuccess変数を更新
		errorMsg := ""
		if statusErr != nil {
			errorMsg = statusErr.Error()
		}

		// Record execution outcome before potential early return
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
			errMsg = statusErr.Error()
			e.Logger.LogError(fmt.Sprintf("Step '%s' failed", step.Name), statusErr)
			if len(bodyBytes) > 0 {
				preview := bodyBytes
				if len(preview) > 2048 {
					preview = preview[:2048]
				}
				e.Logger.LogError("HTTP response body", fmt.Errorf("%s", string(preview)))
			} else {
				e.logResponseBody(response)
			}
			return result
		}

		// Execute validation if provided
		if step.Validation != nil {
			if validationErr := step.Validation(response); validationErr != nil {
				success = false
				result.Status = TestStatusFailed
				result.Error = validationErr
				errMsg = validationErr.Error()
				e.Logger.LogValidation(step.Name, false, validationErr.Error())
			} else {
				e.Logger.LogValidation(step.Name, true, "All validations passed")
			}
		}

		// Execute state update if provided (only when still successful)
		if success && step.StateUpdate != nil {
			if updateErr := step.StateUpdate(response, variables); updateErr != nil {
				success = false
				result.Status = TestStatusFailed
				result.Error = updateErr
				errMsg = updateErr.Error()
				e.Logger.LogError("State update failed", updateErr)
			} else {
				e.Logger.LogStateUpdate(step.Name, variables)
			}
		}
	}

	// 最終的なステータス判定
	if !success && result.Status != TestStatusFailed {
		result.Status = TestStatusFailed
	}
	if success && result.Status != TestStatusPassed {
		result.Status = TestStatusPassed
	}

	e.Coverage.RecordExecution(
		step.ClientMethod,
		"", // Story name will be set by caller
		step.Name,
		statusCode,
		result.Duration,
		success,
		errMsg,
	)

	e.Logger.LogStepResult(step.Name, step.ClientMethod, result.Status, result.StatusCode, result.Duration)

	return result
}

// logResponseBody logs the response body for debugging purposes
func (e *E2EEngine) logResponseBody(response any) {
	if response == nil {
		return
	}

	// Try to extract body from different response types
	var bodyBytes []byte

	// Check if it's a struct with a Body field
	responseValue := reflect.ValueOf(response)
	if responseValue.Kind() == reflect.Ptr {
		responseValue = responseValue.Elem()
	}

	if responseValue.Kind() == reflect.Struct {
		bodyField := responseValue.FieldByName("Body")
		if bodyField.IsValid() && bodyField.Kind() == reflect.Slice && bodyField.Type().Elem().Kind() == reflect.Uint8 {
			bodyBytes = bodyField.Bytes()
		}
	}

	if len(bodyBytes) > 0 && len(bodyBytes) < 10000 { // Only log if reasonable size
		e.Logger.LogDebug(fmt.Sprintf("Response body: %s", string(bodyBytes)))
	}
}

// executeClientMethodWithResponse executes the actual client method using reflection and returns response, status code, and error
func (e *E2EEngine) executeClientMethodWithResponse(methodName string, parameters any, variables map[string]any) ExecutionResult {
	executor := NewMethodExecutor(e.Logger, e.Config.DryRun)
	return executor.Execute(e.Client, methodName, parameters, variables)
}

// PrintResults prints test execution results
func (e *E2EEngine) PrintResults(results []StoryResult) {
	e.Reporter.PrintSummary(results)
	e.Coverage.PrintSummary()
}
