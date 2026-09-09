package snapshot

import (
	"fmt"
	"time"
)

// StoryValidator interface defines the contract for story validation
type StoryValidator interface {
	ValidateStorySnapshot(snapshot *StorySnapshot) (*StoryValidation, error)
	ValidateStoryCompletion(snapshot *StorySnapshot) bool
	ValidateStepSequence(steps []StepSnapshot) []ValidationError
	ValidateStateTransitions(steps []StepSnapshot) []ValidationError
	DetectIncompleteExecution(snapshot *StorySnapshot) []ValidationError
}

// SimpleStoryValidator provides basic story validation functionality
type SimpleStoryValidator struct {
	config *StorySnapshotConfig
	logger *SnapshotLogger
}

// NewSimpleStoryValidator creates a new simple story validator
func NewSimpleStoryValidator(config *StorySnapshotConfig) *SimpleStoryValidator {
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

	return &SimpleStoryValidator{
		config: config,
		logger: logger,
	}
}

// ValidateStorySnapshot validates a complete story snapshot
func (v *SimpleStoryValidator) ValidateStorySnapshot(snapshot *StorySnapshot) (*StoryValidation, error) {
	if snapshot == nil {
		return nil, NewValidationError(ErrorTypeValidation,
			"snapshot cannot be nil", nil, "", "")
	}

	v.logger.LogPhaseStart(PhaseValidation, fmt.Sprintf("story '%s'", snapshot.StoryName))
	v.logger.SetContext("story", snapshot.StoryName)
	defer v.logger.ClearContext()

	validation := &StoryValidation{
		StoryName:      snapshot.StoryName,
		ValidationTime: time.Now(),
		IsValid:        true,
		Summary: ValidationSummary{
			IsValid: true,
		},
	}
	validation.SkippedSteps = v.collectSkippedSteps(snapshot)

	// Determine completion status
	validation.CompletionStatus = v.determineCompletionStatus(snapshot)

	// Validate story completion if rule is enabled
	if v.config.IsValidationRuleEnabled(ValidationRuleCompletion) {
		if !v.ValidateStoryCompletion(snapshot) {
			severity := v.config.GetValidationRuleSeverity(ValidationRuleCompletion)
			error := ValidationError{
				Type:     ValidationErrorSequence,
				StepName: "",
				Message:  "Story execution is not complete",
				Severity: severity,
			}
			validation.SequenceErrors = append(validation.SequenceErrors, error)
			validation.IsValid = false
		}
	}

	// Validate step sequence if rule is enabled
	if v.config.IsValidationRuleEnabled(ValidationRuleSequence) {
		sequenceErrors := v.ValidateStepSequence(snapshot.Steps)
		validation.SequenceErrors = append(validation.SequenceErrors, sequenceErrors...)
		if len(sequenceErrors) > 0 {
			validation.IsValid = false
		}
	}

	// Validate state transitions if rule is enabled
	if v.config.IsValidationRuleEnabled(ValidationRuleStateTransition) {
		stateErrors := v.ValidateStateTransitions(snapshot.Steps)
		validation.StateTransitionErrors = append(validation.StateTransitionErrors, stateErrors...)
		if len(stateErrors) > 0 && v.hasErrorSeverity(stateErrors) {
			validation.IsValid = false
		}
	}

	// Validate timing if rule is enabled
	if v.config.IsValidationRuleEnabled(ValidationRuleTiming) {
		timingErrors := v.ValidateStepTiming(snapshot.Steps)
		validation.TimingErrors = append(validation.TimingErrors, timingErrors...)
		if len(timingErrors) > 0 && v.hasErrorSeverity(timingErrors) {
			validation.IsValid = false
		}
	}

	// Detect incomplete execution
	incompleteErrors := v.DetectIncompleteExecution(snapshot)
	validation.SequenceErrors = append(validation.SequenceErrors, incompleteErrors...)
	if len(incompleteErrors) > 0 {
		validation.IsValid = false
	}

	// Generate summary
	v.generateValidationSummary(validation)

	return validation, nil
}

// ValidateStoryCompletion validates that the story completed successfully
func (v *SimpleStoryValidator) ValidateStoryCompletion(snapshot *StorySnapshot) bool {
	if snapshot == nil {
		return false
	}

	// Check overall story status
	if snapshot.Status != TestStatusPassed {
		return false
	}

	// Check that we have steps
	if len(snapshot.Steps) == 0 {
		return false
	}

	// Check summary
	if snapshot.Summary.FailedSteps > 0 {
		return false
	}

	return true
}

// ValidateStepSequence validates the sequence of steps
func (v *SimpleStoryValidator) ValidateStepSequence(steps []StepSnapshot) []ValidationError {
	var errors []ValidationError
	severity := v.config.GetValidationRuleSeverity(ValidationRuleSequence)

	for i, step := range steps {
		// Check if step has return value
		if step.ReturnValue == nil {
			errors = append(errors, ValidationError{
				Type:     ValidationErrorSequence,
				StepName: step.StepName,
				Message:  "Missing return_value",
				Severity: severity,
			})
		}

		// Check if step success flag matches return value
		if step.ReturnValue != nil && step.Success {
			// For successful steps, check if status code indicates success
			if step.StatusCode >= 400 {
				errors = append(errors, ValidationError{
					Type:          ValidationErrorSequence,
					StepName:      step.StepName,
					Message:       "Step marked as successful but has error status code",
					Severity:      severity,
					ExpectedValue: "status code < 400",
					ActualValue:   step.StatusCode,
				})
			}
		}

		// Check timestamp ordering
		if i > 0 {
			prevStep := steps[i-1]
			if step.Timestamp.Before(prevStep.Timestamp) {
				errors = append(errors, ValidationError{
					Type:     ValidationErrorSequence,
					StepName: step.StepName,
					Message:  "Step timestamp is before previous step",
					Severity: ValidationSeverityWarning, // Usually not critical
				})
			}
		}
	}

	return errors
}

// ValidateStateTransitions validates state transitions between steps
func (v *SimpleStoryValidator) ValidateStateTransitions(steps []StepSnapshot) []ValidationError {
	var errors []ValidationError
	severity := v.config.GetValidationRuleSeverity(ValidationRuleStateTransition)

	for _, step := range steps {
		// Basic validation: check if failed steps have error information
		if !step.Success && step.Error == nil {
			errors = append(errors, ValidationError{
				Type:     ValidationErrorStateTransition,
				StepName: step.StepName,
				Message:  "Step marked as failed but no error information provided",
				Severity: severity,
			})
		}

		// Check if successful steps don't have error information
		if step.Success && step.Error != nil {
			errors = append(errors, ValidationError{
				Type:     ValidationErrorStateTransition,
				StepName: step.StepName,
				Message:  "Step marked as successful but has error information",
				Severity: ValidationSeverityWarning,
			})
		}

		// Validate return value consistency
		if step.ReturnValue != nil {
			// Check if status code matches success flag
			isSuccessStatusCode := step.ReturnValue.StatusCode >= 200 && step.ReturnValue.StatusCode < 400
			if step.Success != isSuccessStatusCode {
				errors = append(errors, ValidationError{
					Type:          ValidationErrorStateTransition,
					StepName:      step.StepName,
					Message:       "Step success flag doesn't match return value status code",
					Severity:      severity,
					ExpectedValue: isSuccessStatusCode,
					ActualValue:   step.Success,
				})
			}
		}
	}

	return errors
}

// DetectIncompleteExecution detects incomplete or failed story execution
func (v *SimpleStoryValidator) DetectIncompleteExecution(snapshot *StorySnapshot) []ValidationError {
	var errors []ValidationError

	// Check if any steps are missing return values (except skipped steps)
	// Skipped steps have: ReturnValue = nil, Success = false, Error = nil, StatusCode = 0
	for _, step := range snapshot.Steps {
		isSkipped := step.ReturnValue == nil && !step.Success && step.Error == nil && step.StatusCode == 0
		if step.ReturnValue == nil && !isSkipped {
			errors = append(errors, ValidationError{
				Type:     ValidationErrorSequence,
				StepName: step.StepName,
				Message:  "Step execution incomplete - missing return value",
				Severity: ValidationSeverityError,
			})
		}
	}

	// Check if story status indicates failure but no step errors are recorded
	if snapshot.Status == TestStatusFailed {
		hasStepErrors := false
		for _, step := range snapshot.Steps {
			if !step.Success || step.Error != nil {
				hasStepErrors = true
				break
			}
		}

		if !hasStepErrors {
			errors = append(errors, ValidationError{
				Type:     ValidationErrorSequence,
				StepName: "",
				Message:  "Story marked as failed but no step failures detected",
				Severity: ValidationSeverityError,
			})
		}
	}

	return errors
}

// ValidateStepTiming validates timing aspects of step execution
func (v *SimpleStoryValidator) ValidateStepTiming(steps []StepSnapshot) []ValidationError {
	var errors []ValidationError
	severity := v.config.GetValidationRuleSeverity(ValidationRuleTiming)

	for _, step := range steps {
		// Check for unreasonably long durations (> 5 minutes)
		if step.Duration > 5*time.Minute {
			errors = append(errors, ValidationError{
				Type:          ValidationErrorTiming,
				StepName:      step.StepName,
				Message:       "Step duration is unusually long",
				Severity:      severity,
				ExpectedValue: "< 5 minutes",
				ActualValue:   step.Duration.String(),
			})
		}

		// Check for zero duration (might indicate timing issue)
		if step.Duration == 0 {
			errors = append(errors, ValidationError{
				Type:     ValidationErrorTiming,
				StepName: step.StepName,
				Message:  "Step duration is zero",
				Severity: ValidationSeverityInfo,
			})
		}
	}

	return errors
}

// Helper methods

// determineCompletionStatus determines the completion status of a story
func (v *SimpleStoryValidator) determineCompletionStatus(snapshot *StorySnapshot) CompletionStatus {
	if snapshot.Status == TestStatusPassed && snapshot.Summary.FailedSteps == 0 {
		return CompletionStatusComplete
	}

	if snapshot.Status == TestStatusFailed {
		if snapshot.Summary.SuccessfulSteps > 0 {
			return CompletionStatusPartial
		}
		return CompletionStatusFailed
	}

	if snapshot.Summary.SuccessfulSteps > 0 && snapshot.Summary.FailedSteps > 0 {
		return CompletionStatusPartial
	}

	return CompletionStatusIncomplete
}

// hasErrorSeverity checks if any errors have error severity
func (v *SimpleStoryValidator) hasErrorSeverity(errors []ValidationError) bool {
	for _, err := range errors {
		if err.Severity == ValidationSeverityError {
			return true
		}
	}
	return false
}

// generateValidationSummary generates the validation summary
func (v *SimpleStoryValidator) generateValidationSummary(validation *StoryValidation) {
	summary := &validation.Summary

	// Count errors by severity
	allErrors := append(validation.SequenceErrors, validation.StateTransitionErrors...)
	allErrors = append(allErrors, validation.TimingErrors...)

	for _, err := range allErrors {
		switch err.Severity {
		case ValidationSeverityError:
			summary.TotalErrors++
		case ValidationSeverityWarning:
			summary.TotalWarnings++
		case ValidationSeverityInfo:
			summary.TotalInfo++
		}
	}

	// Overall validity is based on error count
	summary.IsValid = summary.TotalErrors == 0
	validation.IsValid = summary.IsValid
}

func (v *SimpleStoryValidator) collectSkippedSteps(snapshot *StorySnapshot) []SkippedStepInfo {
	if snapshot == nil {
		return nil
	}
	var skipped []SkippedStepInfo
	for _, step := range snapshot.Steps {
		if step.Status == TestStatusSkipped {
			skipped = append(skipped, SkippedStepInfo{
				StepName: step.StepName,
				Method:   step.Method,
				Reason:   step.SkipReason,
			})
		}
	}
	return skipped
}
