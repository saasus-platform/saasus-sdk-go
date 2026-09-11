package testlib

import "fmt"

// ValidateExpectedStatus ensures the actual status code matches the step expectation when one is defined.
func ValidateExpectedStatus(step Step, statusCode int) error {
	if statusCode == 0 {
		return nil
	}

	if len(step.AllowedStatuses) > 0 {
		for _, allowed := range step.AllowedStatuses {
			if statusCode == allowed {
				return nil
			}
		}
		return fmt.Errorf("unexpected status code for step '%s' (%s): expected one of %v, got %d",
			step.Name, step.ClientMethod, step.AllowedStatuses, statusCode)
	}

	if step.ExpectedStatus == 0 {
		return nil
	}

	if step.ExpectedStatus != statusCode {
		return fmt.Errorf("unexpected status code for step '%s' (%s): expected %d, got %d",
			step.Name, step.ClientMethod, step.ExpectedStatus, statusCode)
	}

	return nil
}
