package snapshot

import (
	"fmt"
	"time"
)

// Story represents a test story with multiple steps (snapshot version)
type Story struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Steps       []Step                 `json:"steps"`
	Variables   map[string]interface{} `json:"variables"`
	Setup       func() error           `json:"-"`
	Cleanup     func() error           `json:"-"`
}

// Step represents a single test step (snapshot version)
type Step struct {
	Name           string                                                             `json:"name"`
	ClientMethod   string                                                             `json:"client_method"`
	Parameters     interface{}                                                        `json:"parameters"`
	ExpectedStatus int                                                                `json:"expected_status"`
	Skip           bool                                                               `json:"skip"`
	SkipReason     string                                                             `json:"skip_reason,omitempty"`
	Validation     func(response interface{}) error                                   `json:"-"`
	StateUpdate    func(response interface{}, variables map[string]interface{}) error `json:"-"`
}

// SDKReturnValue captures the complete return value structure
type SDKReturnValue struct {
	Type         string                 `json:"type"`          // e.g., "*billingapi.GetStripeInfoResponse"
	StatusCode   int                    `json:"status_code"`   // resp.StatusCode()
	Status       string                 `json:"status"`        // resp.Status()
	HTTPResponse *HTTPResponseSnapshot  `json:"http_response"` // resp.HTTPResponse
	JSONData     map[string]interface{} `json:"json_data"`     // resp.JSON200, resp.JSON400, etc.
	Body         string                 `json:"body"`          // string(resp.Body)
	Headers      map[string]string      `json:"headers"`       // resp.HTTPResponse.Header
}

// HTTPResponseSnapshot captures HTTP response details
type HTTPResponseSnapshot struct {
	StatusCode    int               `json:"status_code"`
	Status        string            `json:"status"`
	Headers       map[string]string `json:"headers"`
	ContentLength int64             `json:"content_length"`
	TraceID       string            `json:"trace_id,omitempty"`
}

// SDKMethodError captures error information
type SDKMethodError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// CompatibilityLevel represents the level of compatibility impact
type CompatibilityLevel int

const (
	Compatible CompatibilityLevel = iota // No breaking changes
	Warning                              // Minor changes that might need attention
	Breaking                             // Breaking changes that require action
)

// String returns the string representation of CompatibilityLevel
func (c CompatibilityLevel) String() string {
	switch c {
	case Compatible:
		return "Compatible"
	case Warning:
		return "Warning"
	case Breaking:
		return "Breaking"
	default:
		return "Unknown"
	}
}

// CompatibilityIssue represents a compatibility issue
type CompatibilityIssue struct {
	Type        string             `json:"type"`
	Description string             `json:"description"`
	Impact      CompatibilityLevel `json:"impact"`
	Details     interface{}        `json:"details,omitempty"`
}

// StorySnapshot represents a snapshot of a complete story execution
type StorySnapshot struct {
	StoryName   string                 `json:"story_name"`
	Description string                 `json:"description"`
	Timestamp   time.Time              `json:"timestamp"`
	Duration    time.Duration          `json:"duration"`
	Status      TestStatus             `json:"status"`
	Variables   map[string]interface{} `json:"variables"`
	Steps       []StepSnapshot         `json:"steps"`
	Summary     StoryExecutionSummary  `json:"summary"`
	Metadata    SnapshotMetadata       `json:"metadata"`
}

// StepSnapshot represents a snapshot of a single step execution
type StepSnapshot struct {
	StepName     string                 `json:"step_name"`
	Method       string                 `json:"method"`
	Parameters   map[string]interface{} `json:"parameters"`
	ReturnValue  *SDKReturnValue        `json:"return_value"` // Reusing existing structure
	Duration     time.Duration          `json:"duration"`
	StatusCode   int                    `json:"status_code"`
	Success      bool                   `json:"success"`
	Status       TestStatus             `json:"status"`
	SkipReason   string                 `json:"skip_reason,omitempty"`
	Error        *SDKMethodError        `json:"error,omitempty"` // Reusing existing structure
	Timestamp    time.Time              `json:"timestamp"`
	StateChanges map[string]interface{} `json:"state_changes,omitempty"`
}

// StoryExecutionSummary provides execution summary for a story
type StoryExecutionSummary struct {
	TotalSteps          int           `json:"total_steps"`
	SuccessfulSteps     int           `json:"successful_steps"`
	FailedSteps         int           `json:"failed_steps"`
	SkippedSteps        int           `json:"skipped_steps"`
	TotalDuration       time.Duration `json:"total_duration"`
	AverageStepDuration time.Duration `json:"average_step_duration"`
}

// SnapshotMetadata contains metadata about the snapshot
type SnapshotMetadata struct {
	SDKVersion      string       `json:"sdk_version"`
	TestEnvironment string       `json:"test_environment"`
	CaptureLevel    CaptureLevel `json:"capture_level"`
	GitTag          string       `json:"git_tag,omitempty"`
	GitCommit       string       `json:"git_commit,omitempty"`
}

// TestStatus represents the status of a test execution
type TestStatus string

const (
	TestStatusPassed  TestStatus = "passed"
	TestStatusFailed  TestStatus = "failed"
	TestStatusSkipped TestStatus = "skipped"
)

// CaptureLevel defines the level of detail to capture
type CaptureLevel string

const (
	CaptureLevelFull     CaptureLevel = "FULL"
	CaptureLevelStory    CaptureLevel = "STORY"
	CaptureLevelStep     CaptureLevel = "STEP"
	CaptureLevelResponse CaptureLevel = "RESPONSE"
)

// StorySnapshotConfig configures story snapshot behavior
type StorySnapshotConfig struct {
	EnableCapture            bool             `json:"enable_capture"`
	EnableComparison         bool             `json:"enable_comparison"`
	EnableReporting          bool             `json:"enable_reporting"`
	EnableValidation         bool             `json:"enable_validation"`
	EnableCompatibilityCheck bool             `json:"enable_compatibility_check"` // 要件 4.1, 4.2: 後方互換性チェック機能
	ComparisonMode           string           `json:"comparison_mode"`            // "release", "manual", "skip"
	OutputDirectory          string           `json:"output_directory"`
	ModuleName               string           `json:"module_name,omitempty"`
	SnapshotOnly             bool             `json:"snapshot_only"`
	FileNameFormat           string           `json:"file_name_format"`
	CaptureLevel             CaptureLevel     `json:"capture_level"`
	ValidationRules          []ValidationRule `json:"validation_rules"`
}

// ValidationRule defines a validation rule for story snapshots
type ValidationRule struct {
	Type       ValidationRuleType     `json:"type"`
	Enabled    bool                   `json:"enabled"`
	Severity   ValidationSeverity     `json:"severity"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
}

// ValidationRuleType defines types of validation rules
type ValidationRuleType string

const (
	ValidationRuleSequence        ValidationRuleType = "sequence"
	ValidationRuleStateTransition ValidationRuleType = "state_transition"
	ValidationRuleTiming          ValidationRuleType = "timing"
	ValidationRuleDependency      ValidationRuleType = "dependency"
	ValidationRuleCompletion      ValidationRuleType = "completion"
)

// StoryValidation represents the validation result of a story snapshot
type StoryValidation struct {
	StoryName             string                `json:"story_name"`
	ValidationTime        time.Time             `json:"validation_time"`
	IsValid               bool                  `json:"is_valid"`
	CompletionStatus      CompletionStatus      `json:"completion_status"`
	SequenceErrors        []ValidationError     `json:"sequence_errors"`
	StateTransitionErrors []ValidationError     `json:"state_transition_errors"`
	TimingErrors          []ValidationError     `json:"timing_errors"`
	SkippedSteps          []SkippedStepInfo     `json:"skipped_steps,omitempty"`
	Comparison            *ValidationComparison `json:"comparison,omitempty"`
	Summary               ValidationSummary     `json:"summary"`
}

// SkippedStepInfo captures information about skipped steps during execution
type SkippedStepInfo struct {
	StepName string `json:"step_name"`
	Method   string `json:"method"`
	Reason   string `json:"reason,omitempty"`
}

// ValidationError represents a validation error
type ValidationError struct {
	Type          ValidationErrorType `json:"type"`
	StepName      string              `json:"step_name"`
	Message       string              `json:"message"`
	Severity      ValidationSeverity  `json:"severity"`
	ExpectedValue interface{}         `json:"expected_value,omitempty"`
	ActualValue   interface{}         `json:"actual_value,omitempty"`
}

// CompletionStatus represents the completion status of a story
type CompletionStatus string

const (
	CompletionStatusComplete   CompletionStatus = "complete"
	CompletionStatusIncomplete CompletionStatus = "incomplete"
	CompletionStatusFailed     CompletionStatus = "failed"
	CompletionStatusPartial    CompletionStatus = "partial"
)

// ValidationErrorType defines types of validation errors
type ValidationErrorType string

const (
	ValidationErrorSequence        ValidationErrorType = "sequence"
	ValidationErrorStateTransition ValidationErrorType = "state_transition"
	ValidationErrorTiming          ValidationErrorType = "timing"
	ValidationErrorDependency      ValidationErrorType = "dependency"
)

// ValidationSeverity defines the severity of validation errors
type ValidationSeverity string

const (
	ValidationSeverityError   ValidationSeverity = "error"
	ValidationSeverityWarning ValidationSeverity = "warning"
	ValidationSeverityInfo    ValidationSeverity = "info"
)

// ValidationSummary provides a summary of validation results
type ValidationSummary struct {
	TotalErrors   int  `json:"total_errors"`
	TotalWarnings int  `json:"total_warnings"`
	TotalInfo     int  `json:"total_info"`
	IsValid       bool `json:"is_valid"`
}

// ValidationComparison summarizes differences from the previous validation
type ValidationComparison struct {
	PreviousFile           string            `json:"previous_file,omitempty"`
	PreviousValidationTime *time.Time        `json:"previous_validation_time,omitempty"`
	NewFindings            []ValidationDelta `json:"new_findings,omitempty"`
	ResolvedFindings       []ValidationDelta `json:"resolved_findings,omitempty"`
	ErrorCountDelta        int               `json:"error_count_delta"`
	WarningCountDelta      int               `json:"warning_count_delta"`
	InfoCountDelta         int               `json:"info_count_delta"`
}

// ValidationDelta describes a single validation change between runs
type ValidationDelta struct {
	Type     ValidationErrorType `json:"type"`
	StepName string              `json:"step_name"`
	Message  string              `json:"message"`
	Severity ValidationSeverity  `json:"severity"`
}

// StoryComparison represents the comparison result between two story snapshots
type StoryComparison struct {
	StoryName           string                   `json:"story_name"`
	ComparisonType      string                   `json:"comparison_type"`
	Timestamp           time.Time                `json:"timestamp"`
	OldSnapshot         *StorySnapshot           `json:"old_snapshot"`
	NewSnapshot         *StorySnapshot           `json:"new_snapshot"`
	Differences         []StoryDifference        `json:"differences"`
	CompatibilityReport StoryCompatibilityReport `json:"compatibility_report"`
	Summary             StoryComparisonSummary   `json:"summary"`
}

// StoryDifference represents a difference between two story snapshots
type StoryDifference struct {
	Type        string             `json:"type"`
	StepName    string             `json:"step_name"`
	Field       string             `json:"field"`
	Description string             `json:"description"`
	OldValue    interface{}        `json:"old_value"`
	NewValue    interface{}        `json:"new_value"`
	Impact      CompatibilityLevel `json:"impact"`
}

// Note: CompatibilityLevel is now defined above

// StoryFlowChange represents changes in overall story flow
type StoryFlowChange struct {
	Type        string             `json:"type"`
	Description string             `json:"description"`
	Impact      CompatibilityLevel `json:"impact"`
	Details     interface{}        `json:"details,omitempty"`
}

// StepSequenceChange represents changes in step sequence
type StepSequenceChange struct {
	Type        string             `json:"type"`
	StepName    string             `json:"step_name"`
	Position    int                `json:"position,omitempty"`
	Description string             `json:"description"`
	Impact      CompatibilityLevel `json:"impact"`
}

// StateTransitionChange represents changes in state transitions
type StateTransitionChange struct {
	Type        string             `json:"type"`
	FromStep    string             `json:"from_step"`
	ToStep      string             `json:"to_step"`
	Description string             `json:"description"`
	Impact      CompatibilityLevel `json:"impact"`
	OldState    interface{}        `json:"old_state,omitempty"`
	NewState    interface{}        `json:"new_state,omitempty"`
}

// StoryCompatibilityReport provides compatibility analysis
type StoryCompatibilityReport struct {
	Level                  CompatibilityLevel      `json:"level"`
	Summary                string                  `json:"summary"`
	OverallCompatibility   CompatibilityLevel      `json:"overall_compatibility"`
	BreakingChanges        int                     `json:"breaking_changes"`
	MajorChanges           int                     `json:"major_changes"`
	MinorChanges           int                     `json:"minor_changes"`
	Recommendations        []string                `json:"recommendations"`
	StoryFlowChanges       []StoryFlowChange       `json:"story_flow_changes"`
	StepSequenceChanges    []StepSequenceChange    `json:"step_sequence_changes"`
	StateTransitionChanges []StateTransitionChange `json:"state_transition_changes"`
	Issues                 []CompatibilityIssue    `json:"issues"`
	Passed                 bool                    `json:"passed"`
}

// StoryComparisonSummary provides a summary of the comparison
type StoryComparisonSummary struct {
	TotalDifferences     int                `json:"total_differences"`
	BreakingChanges      int                `json:"breaking_changes"`
	WarningChanges       int                `json:"warning_changes"`
	InfoChanges          int                `json:"info_changes"`
	CompatibleChanges    int                `json:"compatible_changes"`
	StepsCompared        int                `json:"steps_compared"`
	StepsAdded           int                `json:"steps_added"`
	StepsRemoved         int                `json:"steps_removed"`
	StepDifferences      int                `json:"step_differences"`
	ResponseDifferences  int                `json:"response_differences"`
	OverallCompatibility CompatibilityLevel `json:"overall_compatibility"`
}

// StoryReport represents a generated report for story comparison
type StoryReport struct {
	Title                 string                   `json:"title"`
	GeneratedAt           time.Time                `json:"generated_at"`
	StoryName             string                   `json:"story_name"`
	ExecutionSummary      StoryExecutionSummary    `json:"execution_summary"`
	CompatibilityAnalysis StoryCompatibilityReport `json:"compatibility_analysis"`
	StepAnalysis          []StepAnalysis           `json:"step_analysis"`
	Recommendations       []Recommendation         `json:"recommendations"`
	Conclusion            string                   `json:"conclusion"`
}

// StepAnalysis provides analysis for individual steps
type StepAnalysis struct {
	StepName    string             `json:"step_name"`
	Method      string             `json:"method"`
	Status      TestStatus         `json:"status"`
	Differences []StoryDifference  `json:"differences"`
	Impact      CompatibilityLevel `json:"impact"`
	Description string             `json:"description"`
}

// Recommendation provides actionable recommendations
type Recommendation struct {
	Type        RecommendationType     `json:"type"`
	Priority    RecommendationPriority `json:"priority"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Action      string                 `json:"action"`
}

// RecommendationType defines types of recommendations
type RecommendationType string

const (
	RecommendationTypeCompatibility RecommendationType = "compatibility"
	RecommendationTypePerformance   RecommendationType = "performance"
	RecommendationTypeSecurity      RecommendationType = "security"
	RecommendationTypeBestPractice  RecommendationType = "best_practice"
)

// RecommendationPriority defines the priority of recommendations
type RecommendationPriority string

const (
	RecommendationPriorityHigh   RecommendationPriority = "high"
	RecommendationPriorityMedium RecommendationPriority = "medium"
	RecommendationPriorityLow    RecommendationPriority = "low"
)

// StorySnapshotResult represents the complete result of story snapshot execution
type StorySnapshotResult struct {
	StoryResult StoryResult      `json:"story_result"`
	Snapshot    StorySnapshot    `json:"snapshot"`
	Comparison  *StoryComparison `json:"comparison,omitempty"`
	Report      *StoryReport     `json:"report,omitempty"`
	Validation  *StoryValidation `json:"validation,omitempty"`
}

// CLIConfig represents configuration for CLI-based execution control
// This enables flexible execution of different phases from main.go
type CLIConfig struct {
	Mode           string   `json:"mode"`            // "full", "capture", "compare", "report"
	Service        string   `json:"service"`         // "billing", "auth", "pricing", etc.
	Stories        []string `json:"stories"`         // Specific story names to execute
	OldTag         string   `json:"old_tag"`         // Tag for comparison (old version)
	NewTag         string   `json:"new_tag"`         // Tag for comparison (new version)
	ComparisonFile string   `json:"comparison_file"` // Path to comparison result file
	SnapshotFile   string   `json:"snapshot_file"`   // Path to snapshot file
	ConfigFile     string   `json:"config_file"`     // Path to configuration file
	OutputDir      string   `json:"output_dir"`      // Output directory override
	Verbose        bool     `json:"verbose"`         // Enable verbose logging
	Parallel       bool     `json:"parallel"`        // Enable parallel execution
	DryRun         bool     `json:"dry_run"`         // Dry run mode
}

// ExecutionMode represents different execution modes for CLI
type ExecutionMode string

const (
	ExecutionModeFull    ExecutionMode = "full"    // Execute all phases (capture, compare, report)
	ExecutionModeCapture ExecutionMode = "capture" // Execute only snapshot capture
	ExecutionModeCompare ExecutionMode = "compare" // Execute only snapshot comparison
	ExecutionModeReport  ExecutionMode = "report"  // Execute only report generation
)

// Validate validates the CLI configuration
func (c *CLIConfig) Validate() error {
	// Validate execution mode
	validModes := map[string]bool{
		"full":    true,
		"capture": true,
		"compare": true,
		"report":  true,
	}
	if !validModes[c.Mode] {
		return fmt.Errorf("invalid execution mode: %s", c.Mode)
	}

	// Mode-specific validations
	switch ExecutionMode(c.Mode) {
	case ExecutionModeCompare:
		// Compare mode requires either snapshot file or stories/tags
		if c.SnapshotFile == "" && len(c.Stories) == 0 && c.OldTag == "" && c.NewTag == "" {
			return fmt.Errorf("compare mode requires either snapshot_file, stories, or old_tag/new_tag")
		}
	case ExecutionModeReport:
		// Report mode requires either comparison file or stories
		if c.ComparisonFile == "" && len(c.Stories) == 0 {
			return fmt.Errorf("report mode requires either comparison_file or stories")
		}
	case ExecutionModeFull, ExecutionModeCapture:
		// Full and capture modes require service to be specified for story selection
		if c.Service == "" && len(c.Stories) == 0 {
			return fmt.Errorf("%s mode requires either service or specific stories", c.Mode)
		}
	}

	return nil
}

// GetExecutionMode returns the execution mode as an enum
func (c *CLIConfig) GetExecutionMode() ExecutionMode {
	return ExecutionMode(c.Mode)
}

// IsVerbose returns whether verbose logging is enabled
func (c *CLIConfig) IsVerbose() bool {
	return c.Verbose
}

// IsDryRun returns whether dry run mode is enabled
func (c *CLIConfig) IsDryRun() bool {
	return c.DryRun
}

// IsParallel returns whether parallel execution is enabled
func (c *CLIConfig) IsParallel() bool {
	return c.Parallel
}

// HasSpecificStories returns whether specific stories are requested
func (c *CLIConfig) HasSpecificStories() bool {
	return len(c.Stories) > 0
}

// HasTagComparison returns whether tag-based comparison is requested
func (c *CLIConfig) HasTagComparison() bool {
	return c.OldTag != "" && c.NewTag != ""
}

// GetOutputDirectory returns the output directory, using default if not specified
func (c *CLIConfig) GetOutputDirectory(defaultDir string) string {
	if c.OutputDir != "" {
		return c.OutputDir
	}
	return defaultDir
}

// StoryResult represents the result of executing a story (snapshot version)
type StoryResult struct {
	StoryName string                 `json:"story_name"`
	Status    TestStatus             `json:"status"`
	Duration  time.Duration          `json:"duration"`
	Steps     []StepResult           `json:"steps"`
	Error     error                  `json:"error,omitempty"`
	Variables map[string]interface{} `json:"variables"`
}

// StepResult represents the result of executing a step (snapshot version)
type StepResult struct {
	StepName    string          `json:"step_name"`
	Method      string          `json:"method"`
	Status      TestStatus      `json:"status"`
	Duration    time.Duration   `json:"duration"`
	StatusCode  int             `json:"status_code"`
	Error       error           `json:"error,omitempty"`
	ReturnValue *SDKReturnValue `json:"return_value,omitempty"`
	SkipReason  string          `json:"skip_reason,omitempty"`
}
