package billingapi

import (
	"path/filepath"

	"github.com/saasus-platform/saasus-sdk-go/tests/testlib/snapshot"
)

// GetBillingStorySnapshotConfig returns the default story snapshot configuration for billing API
func GetBillingStorySnapshotConfig() *snapshot.StorySnapshotConfig {
	// Get absolute path to project root
	projectRoot := filepath.Join("..", "..", "..")
	outputDir := filepath.Join(projectRoot, "tests", "e2e", "snapshot")

	return &snapshot.StorySnapshotConfig{
		EnableCapture:    false, // Will be set by flag parser
		EnableComparison: false, // Will be set by flag parser
		EnableReporting:  false, // Will be set by flag parser
		EnableValidation: true,
		ComparisonMode:   "release",
		OutputDirectory:  outputDir,
		ModuleName:       "billing",
		SnapshotOnly:     false,
		FileNameFormat:   "story_snapshot_{tag}_{story_name}.json",
		CaptureLevel:     snapshot.CaptureLevelFull,
		ValidationRules:  GetBillingValidationRules(),
	}
}

// GetBillingValidationRules returns validation rules specific to billing API
func GetBillingValidationRules() []snapshot.ValidationRule {
	return []snapshot.ValidationRule{
		{
			Type:     snapshot.ValidationRuleCompletion,
			Enabled:  true,
			Severity: snapshot.ValidationSeverityError,
			Parameters: map[string]interface{}{
				"description": "All billing API calls must complete successfully",
			},
		},
		{
			Type:     snapshot.ValidationRuleSequence,
			Enabled:  true,
			Severity: snapshot.ValidationSeverityError,
			Parameters: map[string]interface{}{
				"description": "All billing API calls must have valid step sequence",
			},
		},
		{
			Type:     snapshot.ValidationRuleStateTransition,
			Enabled:  true,
			Severity: snapshot.ValidationSeverityWarning,
			Parameters: map[string]interface{}{
				"description": "All billing API calls must have valid state transitions",
			},
		},
		{
			Type:     snapshot.ValidationRuleTiming,
			Enabled:  false, // Disabled by default as timing can be variable
			Severity: snapshot.ValidationSeverityInfo,
			Parameters: map[string]interface{}{
				"description": "Billing API call timing validation",
			},
		},
	}
}

// GetBillingSnapshotConfigForCapture returns configuration optimized for capture-only mode
func GetBillingSnapshotConfigForCapture() *snapshot.StorySnapshotConfig {
	config := GetBillingStorySnapshotConfig()
	config.EnableCapture = true
	config.EnableComparison = false
	config.EnableReporting = false
	config.EnableValidation = false
	config.SnapshotOnly = true
	return config
}

// GetBillingSnapshotConfigForComparison returns configuration optimized for comparison-only mode
func GetBillingSnapshotConfigForComparison() *snapshot.StorySnapshotConfig {
	config := GetBillingStorySnapshotConfig()
	config.EnableCapture = false
	config.EnableComparison = true
	config.EnableReporting = false
	config.EnableValidation = false
	return config
}

// GetBillingSnapshotConfigForReporting returns configuration optimized for reporting-only mode
func GetBillingSnapshotConfigForReporting() *snapshot.StorySnapshotConfig {
	config := GetBillingStorySnapshotConfig()
	config.EnableCapture = false
	config.EnableComparison = false
	config.EnableReporting = true
	config.EnableValidation = false
	return config
}

// GetBillingSnapshotConfigForIntegrated returns configuration for integrated execution (all phases)
func GetBillingSnapshotConfigForIntegrated() *snapshot.StorySnapshotConfig {
	config := GetBillingStorySnapshotConfig()
	config.EnableCapture = true
	config.EnableComparison = true
	config.EnableReporting = true
	config.EnableValidation = true
	config.SnapshotOnly = false
	return config
}
