package apilogapi

import (
	"path/filepath"

	"github.com/saasus-platform/saasus-sdk-go/tests/testlib/snapshot"
)

// GetApilogStorySnapshotConfig returns the default story snapshot configuration for apilog API
func GetApilogStorySnapshotConfig() *snapshot.StorySnapshotConfig {
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
		ModuleName:       "apilog",
		SnapshotOnly:     false,
		FileNameFormat:   "story_snapshot_{tag}_{story_name}.json",
		CaptureLevel:     snapshot.CaptureLevelFull,
		ValidationRules:  GetApilogValidationRules(),
	}
}

// GetApilogValidationRules returns validation rules specific to apilog API
func GetApilogValidationRules() []snapshot.ValidationRule {
	return []snapshot.ValidationRule{
		{
			Type:     snapshot.ValidationRuleCompletion,
			Enabled:  true,
			Severity: snapshot.ValidationSeverityError,
			Parameters: map[string]any{
				"description": "All apilog API calls must complete successfully",
			},
		},
		{
			Type:     snapshot.ValidationRuleSequence,
			Enabled:  true,
			Severity: snapshot.ValidationSeverityError,
			Parameters: map[string]any{
				"description": "All apilog API calls must have valid step sequence",
			},
		},
		{
			Type:     snapshot.ValidationRuleStateTransition,
			Enabled:  true,
			Severity: snapshot.ValidationSeverityWarning,
			Parameters: map[string]any{
				"description": "All apilog API calls must have valid state transitions",
			},
		},
		{
			Type:     snapshot.ValidationRuleTiming,
			Enabled:  false, // Disabled by default as timing can be variable
			Severity: snapshot.ValidationSeverityInfo,
			Parameters: map[string]any{
				"description": "Apilog API call timing validation",
			},
		},
	}
}

// GetApilogSnapshotConfigForCapture returns configuration optimized for capture-only mode
func GetApilogSnapshotConfigForCapture() *snapshot.StorySnapshotConfig {
	config := GetApilogStorySnapshotConfig()
	config.EnableCapture = true
	config.EnableComparison = false
	config.EnableReporting = false
	config.EnableValidation = false
	config.SnapshotOnly = true
	return config
}

// GetApilogSnapshotConfigForComparison returns configuration optimized for comparison-only mode
func GetApilogSnapshotConfigForComparison() *snapshot.StorySnapshotConfig {
	config := GetApilogStorySnapshotConfig()
	config.EnableCapture = false
	config.EnableComparison = true
	config.EnableReporting = false
	config.EnableValidation = false
	return config
}

// GetApilogSnapshotConfigForReporting returns configuration optimized for reporting-only mode
func GetApilogSnapshotConfigForReporting() *snapshot.StorySnapshotConfig {
	config := GetApilogStorySnapshotConfig()
	config.EnableCapture = false
	config.EnableComparison = false
	config.EnableReporting = true
	config.EnableValidation = false
	return config
}

// GetApilogSnapshotConfigForIntegrated returns configuration for integrated execution (all phases)
func GetApilogSnapshotConfigForIntegrated() *snapshot.StorySnapshotConfig {
	config := GetApilogStorySnapshotConfig()
	config.EnableCapture = true
	config.EnableComparison = true
	config.EnableReporting = true
	config.EnableValidation = true
	config.SnapshotOnly = false
	return config
}
