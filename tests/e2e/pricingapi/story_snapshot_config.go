package pricingapi

import (
	"path/filepath"

	"github.com/saasus-platform/saasus-sdk-go/tests/testlib/snapshot"
)

// GetPricingStorySnapshotConfig returns the default story snapshot configuration for pricing API
func GetPricingStorySnapshotConfig() *snapshot.StorySnapshotConfig {
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
		ModuleName:       "pricing",
		SnapshotOnly:     false,
		FileNameFormat:   "story_snapshot_{tag}_{story_name}.json",
		CaptureLevel:     snapshot.CaptureLevelFull,
		ValidationRules:  GetPricingValidationRules(),
	}
}

// GetPricingValidationRules returns validation rules specific to pricing API
func GetPricingValidationRules() []snapshot.ValidationRule {
	return []snapshot.ValidationRule{
		{
			Type:     snapshot.ValidationRuleCompletion,
			Enabled:  true,
			Severity: snapshot.ValidationSeverityError,
			Parameters: map[string]interface{}{
				"description": "All pricing API calls must complete successfully",
			},
		},
		{
			Type:     snapshot.ValidationRuleSequence,
			Enabled:  true,
			Severity: snapshot.ValidationSeverityError,
			Parameters: map[string]interface{}{
				"description": "All pricing API calls must have valid step sequence",
			},
		},
		{
			Type:     snapshot.ValidationRuleStateTransition,
			Enabled:  true,
			Severity: snapshot.ValidationSeverityWarning,
			Parameters: map[string]interface{}{
				"description": "All pricing API calls must have valid state transitions",
			},
		},
		{
			Type:     snapshot.ValidationRuleTiming,
			Enabled:  false, // Disabled by default as timing can be variable
			Severity: snapshot.ValidationSeverityInfo,
			Parameters: map[string]interface{}{
				"description": "Pricing API call timing validation",
			},
		},
	}
}

// GetPricingSnapshotConfigForCapture returns configuration optimized for capture-only mode
func GetPricingSnapshotConfigForCapture() *snapshot.StorySnapshotConfig {
	config := GetPricingStorySnapshotConfig()
	config.EnableCapture = true
	config.EnableComparison = false
	config.EnableReporting = false
	config.EnableValidation = false
	config.SnapshotOnly = true
	return config
}

// GetPricingSnapshotConfigForComparison returns configuration optimized for comparison-only mode
func GetPricingSnapshotConfigForComparison() *snapshot.StorySnapshotConfig {
	config := GetPricingStorySnapshotConfig()
	config.EnableCapture = false
	config.EnableComparison = true
	config.EnableReporting = false
	config.EnableValidation = false
	return config
}

// GetPricingSnapshotConfigForReporting returns configuration optimized for reporting-only mode
func GetPricingSnapshotConfigForReporting() *snapshot.StorySnapshotConfig {
	config := GetPricingStorySnapshotConfig()
	config.EnableCapture = false
	config.EnableComparison = false
	config.EnableReporting = true
	config.EnableValidation = false
	return config
}

// GetPricingSnapshotConfigForIntegrated returns configuration for integrated execution (all phases)
func GetPricingSnapshotConfigForIntegrated() *snapshot.StorySnapshotConfig {
	config := GetPricingStorySnapshotConfig()
	config.EnableCapture = true
	config.EnableComparison = true
	config.EnableReporting = true
	config.EnableValidation = true
	config.SnapshotOnly = false
	return config
}
