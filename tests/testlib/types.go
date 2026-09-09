package testlib

import "time"

// TestStatus represents the status of a test execution
type TestStatus string

const (
	TestStatusPassed  TestStatus = "passed"
	TestStatusFailed  TestStatus = "failed"
	TestStatusSkipped TestStatus = "skipped"
)

// Story represents a test story with multiple steps
type Story struct {
	Name        string
	Description string
	Steps       []Step
	Variables   map[string]any
	Setup       func() error
	Cleanup     func() error
}

// Step represents a single test step
type Step struct {
	Name            string
	ClientMethod    string
	Parameters      any
	ExpectedStatus  int
	AllowedStatuses []int
	Skip            bool
	SkipReason      string
	Validation      func(response any) error
	StateUpdate     func(response any, variables map[string]any) error
}

// StoryResult represents the result of executing a story
type StoryResult struct {
	StoryName string
	Status    TestStatus
	Duration  time.Duration
	Steps     []StepResult
	Error     error
	Variables map[string]any
}

// StepResult represents the result of executing a step
type StepResult struct {
	StepName   string
	Method     string
	Status     TestStatus
	Duration   time.Duration
	StatusCode int
	Error      error
	SkipReason string
}

// MethodExecution represents a method execution record
type MethodExecution struct {
	MethodName string
	StoryName  string
	StepName   string
	StatusCode int
	Duration   time.Duration
	Success    bool
	Error      string
	Timestamp  time.Time
}
