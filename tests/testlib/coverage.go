package testlib

import (
	"fmt"
	"sync"
	"time"
)

// CoverageTracker tracks method execution coverage
type CoverageTracker struct {
	mu         sync.RWMutex
	executions map[string][]MethodExecution
	methods    []string
}

// NewCoverageTracker creates a new coverage tracker
func NewCoverageTracker(methods []string) *CoverageTracker {
	return &CoverageTracker{
		executions: make(map[string][]MethodExecution),
		methods:    methods,
	}
}

// RecordExecution records a method execution
func (ct *CoverageTracker) RecordExecution(methodName, storyName, stepName string, statusCode int, duration time.Duration, success bool, errorMsg string) {
	ct.mu.Lock()
	defer ct.mu.Unlock()

	execution := MethodExecution{
		MethodName: methodName,
		StoryName:  storyName,
		StepName:   stepName,
		StatusCode: statusCode,
		Duration:   duration,
		Success:    success,
		Error:      errorMsg,
		Timestamp:  time.Now(),
	}

	ct.executions[methodName] = append(ct.executions[methodName], execution)
}

// GetCoverage returns coverage statistics
func (ct *CoverageTracker) GetCoverage() (int, int, float64) {
	ct.mu.RLock()
	defer ct.mu.RUnlock()

	totalMethods := len(ct.methods)
	coveredMethods := len(ct.executions)
	coverage := float64(coveredMethods) / float64(totalMethods) * 100

	return coveredMethods, totalMethods, coverage
}

// GetMethodStats returns statistics for a specific method
func (ct *CoverageTracker) GetMethodStats(methodName string) (int, float64, time.Duration) {
	ct.mu.RLock()
	defer ct.mu.RUnlock()

	executions := ct.executions[methodName]
	if len(executions) == 0 {
		return 0, 0, 0
	}

	successCount := 0
	totalDuration := time.Duration(0)

	for _, exec := range executions {
		if exec.Success {
			successCount++
		}
		totalDuration += exec.Duration
	}

	successRate := float64(successCount) / float64(len(executions)) * 100
	avgDuration := totalDuration / time.Duration(len(executions))

	return len(executions), successRate, avgDuration
}

// GetUntestedMethods returns methods that haven't been executed
func (ct *CoverageTracker) GetUntestedMethods() []string {
	ct.mu.RLock()
	defer ct.mu.RUnlock()

	var untested []string
	for _, method := range ct.methods {
		if _, exists := ct.executions[method]; !exists {
			untested = append(untested, method)
		}
	}
	return untested
}

// IsFullyCovered returns true if all methods have been executed
func (ct *CoverageTracker) IsFullyCovered() bool {
	return len(ct.GetUntestedMethods()) == 0
}

// PrintSummary prints coverage summary
func (ct *CoverageTracker) PrintSummary() {
	covered, total, coverage := ct.GetCoverage()

	fmt.Println("📊 Method Coverage Summary:")
	fmt.Printf("Total Coverage: %d/%d (%.1f%%)\n", covered, total, coverage)

	if ct.IsFullyCovered() {
		fmt.Println("🎉 All methods have been tested!")
	} else {
		untested := ct.GetUntestedMethods()
		fmt.Printf("⚠️  %d methods still need testing:\n", len(untested))
		for _, method := range untested {
			fmt.Printf("  - %s\n", method)
		}
	}
}
