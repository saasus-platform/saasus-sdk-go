package snapshot

import "fmt"

type validationCounts struct {
	errors   int
	warnings int
	infos    int
}

// generateValidationComparison builds a comparison summary between previous and current validations.
func generateValidationComparison(previous, current *StoryValidation, previousFile string) *ValidationComparison {
	if previous == nil || current == nil {
		return nil
	}

	comparison := &ValidationComparison{
		PreviousFile: previousFile,
	}
	if !previous.ValidationTime.IsZero() {
		prevTime := previous.ValidationTime
		comparison.PreviousValidationTime = &prevTime
	}

	prevMap, prevCounts := flattenValidationErrors(previous)
	currMap, currCounts := flattenValidationErrors(current)

	for key, entry := range currMap {
		if _, exists := prevMap[key]; !exists {
			comparison.NewFindings = append(comparison.NewFindings, entry)
		}
	}
	for key, entry := range prevMap {
		if _, exists := currMap[key]; !exists {
			comparison.ResolvedFindings = append(comparison.ResolvedFindings, entry)
		}
	}

	comparison.ErrorCountDelta = currCounts.errors - prevCounts.errors
	comparison.WarningCountDelta = currCounts.warnings - prevCounts.warnings
	comparison.InfoCountDelta = currCounts.infos - prevCounts.infos

	return comparison
}

func flattenValidationErrors(validation *StoryValidation) (map[string]ValidationDelta, validationCounts) {
	errors := make(map[string]ValidationDelta)
	counts := validationCounts{}

	collect := func(list []ValidationError) {
		for _, err := range list {
			key := fmt.Sprintf("%s|%s|%s|%s", err.Type, err.Severity, err.StepName, err.Message)
			errors[key] = ValidationDelta{
				Type:     err.Type,
				StepName: err.StepName,
				Message:  err.Message,
				Severity: err.Severity,
			}
			switch err.Severity {
			case ValidationSeverityError:
				counts.errors++
			case ValidationSeverityWarning:
				counts.warnings++
			case ValidationSeverityInfo:
				counts.infos++
			}
		}
	}

	collect(validation.SequenceErrors)
	collect(validation.StateTransitionErrors)
	collect(validation.TimingErrors)

	return errors, counts
}
