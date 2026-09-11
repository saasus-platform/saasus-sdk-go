package integrationapi

import (
	"path/filepath"

	"github.com/saasus-platform/saasus-sdk-go/tests/testlib/snapshot"
)

// GetIntegrationSnapshotConfigForCapture returns snapshot configuration for capture mode
func GetIntegrationSnapshotConfigForCapture() *snapshot.StorySnapshotConfig {
	// Get absolute path to project root
	projectRoot := filepath.Join("..", "..", "..")
	outputDir := filepath.Join(projectRoot, "tests", "e2e", "snapshot")

	return &snapshot.StorySnapshotConfig{
		EnableCapture:            true,
		EnableComparison:         false,
		EnableReporting:          false,
		EnableValidation:         true,
		EnableCompatibilityCheck: false,
		ComparisonMode:           "skip",
		OutputDirectory:          outputDir,
		ModuleName:               "integration",
		SnapshotOnly:             false,
		FileNameFormat:           "story_snapshot_{tag}_{story_name}.json",
		CaptureLevel:             snapshot.CaptureLevelFull,
		ValidationRules: []snapshot.ValidationRule{
			{
				Type:     snapshot.ValidationRuleCompletion,
				Enabled:  true,
				Severity: snapshot.ValidationSeverityError,
			},
			{
				Type:     snapshot.ValidationRuleSequence,
				Enabled:  true,
				Severity: snapshot.ValidationSeverityError,
			},
			{
				Type:     snapshot.ValidationRuleStateTransition,
				Enabled:  true,
				Severity: snapshot.ValidationSeverityWarning,
			},
			{
				Type:     snapshot.ValidationRuleTiming,
				Enabled:  false, // Disabled by default
				Severity: snapshot.ValidationSeverityInfo,
			},
		},
	}
}

// GetIntegrationSnapshotConfigForComparison returns snapshot configuration for comparison mode
func GetIntegrationSnapshotConfigForComparison() *snapshot.StorySnapshotConfig {
	config := GetIntegrationSnapshotConfigForCapture()
	config.EnableCapture = false
	config.EnableComparison = true
	config.ComparisonMode = "release"
	return config
}

// GetIntegrationSnapshotConfigForReporting returns snapshot configuration for reporting mode
func GetIntegrationSnapshotConfigForReporting() *snapshot.StorySnapshotConfig {
	config := GetIntegrationSnapshotConfigForCapture()
	config.EnableCapture = false
	config.EnableComparison = false
	config.EnableReporting = true
	return config
}

// GetIntegrationSnapshotConfigForIntegrated returns snapshot configuration for integrated mode
func GetIntegrationSnapshotConfigForIntegrated() *snapshot.StorySnapshotConfig {
	config := GetIntegrationSnapshotConfigForCapture()
	config.EnableCapture = true
	config.EnableComparison = true
	config.EnableReporting = true
	config.EnableCompatibilityCheck = true
	config.ComparisonMode = "release"
	return config
}
