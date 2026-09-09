package snapshot

import (
	"flag"
	"fmt"
	"strings"
)

// TestSnapshotConfig represents configuration for test-based snapshot execution
// This is specifically designed for integration with go test -v command flags
type TestSnapshotConfig struct {
	Mode          string   `json:"mode"`           // "full", "capture", "compare", "report"
	Stories       []string `json:"stories"`        // Comma-separated story names to execute
	OldTag        string   `json:"old_tag"`        // Old tag for comparison
	NewTag        string   `json:"new_tag"`        // New tag for comparison
	ConfigFile    string   `json:"config_file"`    // Snapshot configuration file path
	OutputDir     string   `json:"output_dir"`     // Output directory for snapshots
	EnableCapture bool     `json:"enable_capture"` // キャプチャ有効化
	EnableCompare bool     `json:"enable_compare"` // 比較有効化
	EnableReport  bool     `json:"enable_report"`  // レポート有効化
	Verbose       bool     `json:"verbose"`        // 詳細ログ出力
}

// Test flags for snapshot functionality
// These flags are parsed when go test -v is executed with snapshot options
var (
	snapshotMode    = flag.String("snapshot-mode", "", "Snapshot mode: capture, compare, report, full")
	snapshotStories = flag.String("snapshot-stories", "", "Comma-separated story names to execute")
	snapshotOldTag  = flag.String("snapshot-old-tag", "", "Old tag for comparison")
	snapshotNewTag  = flag.String("snapshot-new-tag", "", "New tag for comparison")
	snapshotConfig  = flag.String("snapshot-config", "", "Snapshot configuration file path")
	snapshotOutput  = flag.String("snapshot-output", "", "Output directory for snapshots")
	snapshotVerbose = flag.Bool("snapshot-verbose", false, "Enable verbose snapshot logging")
)

// ParseSnapshotFlags parses command line flags and returns TestSnapshotConfig
// This function is called during test execution to determine snapshot behavior
func ParseSnapshotFlags() *TestSnapshotConfig {
	// Parse flags if not already parsed
	if !flag.Parsed() {
		flag.Parse()
	}

	config := &TestSnapshotConfig{
		Mode:       *snapshotMode,
		OldTag:     *snapshotOldTag,
		NewTag:     *snapshotNewTag,
		ConfigFile: *snapshotConfig,
		OutputDir:  *snapshotOutput,
		Verbose:    *snapshotVerbose,
	}

	// Parse comma-separated stories
	if *snapshotStories != "" {
		stories := strings.Split(*snapshotStories, ",")
		// Trim whitespace from each story name and filter out empty strings
		config.Stories = make([]string, 0, len(stories))
		for _, story := range stories {
			trimmed := strings.TrimSpace(story)
			if trimmed != "" {
				config.Stories = append(config.Stories, trimmed)
			}
		}
	}

	// Set feature flags based on mode
	switch config.Mode {
	case "capture":
		config.EnableCapture = true
		config.EnableCompare = false
		config.EnableReport = false
	case "compare":
		config.EnableCapture = false
		config.EnableCompare = true
		config.EnableReport = false
	case "report":
		config.EnableCapture = false
		config.EnableCompare = false
		config.EnableReport = true
	case "full":
		config.EnableCapture = true
		config.EnableCompare = true
		config.EnableReport = true
	default:
		// No snapshot mode specified - all features disabled
		config.EnableCapture = false
		config.EnableCompare = false
		config.EnableReport = false
	}

	return config
}

// IsSnapshotModeEnabled returns true if any snapshot mode is enabled
func (c *TestSnapshotConfig) IsSnapshotModeEnabled() bool {
	return c.Mode != "" && (c.EnableCapture || c.EnableCompare || c.EnableReport)
}

// IsNormalTestMode returns true if no snapshot mode is specified
// In this case, tests should run normally without snapshot functionality
func (c *TestSnapshotConfig) IsNormalTestMode() bool {
	return c.Mode == ""
}

// Validate validates the test snapshot configuration
func (c *TestSnapshotConfig) Validate() error {
	// If no mode specified, it's valid (normal test mode)
	if c.Mode == "" {
		return nil
	}

	// Validate mode
	validModes := map[string]bool{
		"capture": true,
		"compare": true,
		"report":  true,
		"full":    true,
	}
	if !validModes[c.Mode] {
		return fmt.Errorf("invalid snapshot mode: %s (valid modes: capture, compare, report, full)", c.Mode)
	}

	// Mode-specific validations
	switch c.Mode {
	case "compare":
		// Compare mode should have either old/new tags or stories
		if c.OldTag == "" && c.NewTag == "" && len(c.Stories) == 0 {
			return fmt.Errorf("compare mode requires either -snapshot-old-tag/-snapshot-new-tag or -snapshot-stories")
		}
	case "report":
		// Report mode should have stories or tags for report generation
		if len(c.Stories) == 0 && c.OldTag == "" && c.NewTag == "" {
			return fmt.Errorf("report mode requires either -snapshot-stories or -snapshot-old-tag/-snapshot-new-tag")
		}
	}

	return nil
}

// ToStorySnapshotConfig converts TestSnapshotConfig to StorySnapshotConfig
// This allows integration with existing snapshot infrastructure
func (c *TestSnapshotConfig) ToStorySnapshotConfig() *StorySnapshotConfig {
	config := DefaultStorySnapshotConfig()

	// Override with test-specific settings
	config.EnableCapture = c.EnableCapture
	config.EnableComparison = c.EnableCompare
	config.EnableReporting = c.EnableReport

	// Set output directory if specified
	if c.OutputDir != "" {
		config.OutputDirectory = c.OutputDir
	}

	// Set comparison mode based on tags
	if c.OldTag != "" && c.NewTag != "" {
		config.ComparisonMode = "manual"
	} else {
		config.ComparisonMode = "release"
	}

	return config
}

// GetExecutionMode returns the execution mode as a string
func (c *TestSnapshotConfig) GetExecutionMode() string {
	return c.Mode
}

// HasSpecificStories returns true if specific stories are requested
func (c *TestSnapshotConfig) HasSpecificStories() bool {
	return len(c.Stories) > 0
}

// HasTagComparison returns true if tag-based comparison is requested
func (c *TestSnapshotConfig) HasTagComparison() bool {
	return c.OldTag != "" && c.NewTag != ""
}

// GetStoryNames returns the list of story names to execute
func (c *TestSnapshotConfig) GetStoryNames() []string {
	return c.Stories
}

// GetOldTag returns the old tag for comparison
func (c *TestSnapshotConfig) GetOldTag() string {
	return c.OldTag
}

// GetNewTag returns the new tag for comparison
func (c *TestSnapshotConfig) GetNewTag() string {
	return c.NewTag
}

// GetConfigFile returns the configuration file path
func (c *TestSnapshotConfig) GetConfigFile() string {
	return c.ConfigFile
}

// GetOutputDir returns the output directory
func (c *TestSnapshotConfig) GetOutputDir() string {
	return c.OutputDir
}

// IsVerbose returns true if verbose logging is enabled
func (c *TestSnapshotConfig) IsVerbose() bool {
	return c.Verbose
}

// String returns a string representation of the configuration
func (c *TestSnapshotConfig) String() string {
	if c.IsNormalTestMode() {
		return "TestSnapshotConfig{mode: normal}"
	}

	parts := []string{
		fmt.Sprintf("mode: %s", c.Mode),
	}

	if len(c.Stories) > 0 {
		parts = append(parts, fmt.Sprintf("stories: [%s]", strings.Join(c.Stories, ", ")))
	}

	if c.OldTag != "" {
		parts = append(parts, fmt.Sprintf("old_tag: %s", c.OldTag))
	}

	if c.NewTag != "" {
		parts = append(parts, fmt.Sprintf("new_tag: %s", c.NewTag))
	}

	if c.ConfigFile != "" {
		parts = append(parts, fmt.Sprintf("config_file: %s", c.ConfigFile))
	}

	if c.OutputDir != "" {
		parts = append(parts, fmt.Sprintf("output_dir: %s", c.OutputDir))
	}

	if c.Verbose {
		parts = append(parts, "verbose: true")
	}

	return fmt.Sprintf("TestSnapshotConfig{%s}", strings.Join(parts, ", "))
}

// GetUsageExamples returns usage examples for the snapshot flags
func GetUsageExamples() string {
	return `
Snapshot Testing Usage Examples:

# Normal test execution (no snapshot functionality)
go test -v ./tests/e2e/billingapi

# Capture snapshots during test execution
go test -v ./tests/e2e/billingapi -snapshot-mode=capture

# Capture snapshots for specific stories
go test -v ./tests/e2e/billingapi -snapshot-mode=capture -snapshot-stories=postman_standard,postman_withresponse

# Compare snapshots with previous release
go test -v ./tests/e2e/billingapi -snapshot-mode=compare

# Compare specific tags
go test -v ./tests/e2e/billingapi -snapshot-mode=compare -snapshot-old-tag=v1.2.3 -snapshot-new-tag=v1.2.4

# Generate reports from comparison results
go test -v ./tests/e2e/billingapi -snapshot-mode=report

# Full execution (capture, compare, and report)
go test -v ./tests/e2e/billingapi -snapshot-mode=full

# Use custom configuration file
go test -v ./tests/e2e/billingapi -snapshot-mode=capture -snapshot-config=configs/billing_snapshot.json

# Use custom output directory
go test -v ./tests/e2e/billingapi -snapshot-mode=capture -snapshot-output=custom_output

# Enable verbose logging
go test -v ./tests/e2e/billingapi -snapshot-mode=capture -snapshot-verbose

Available Flags:
  -snapshot-mode string
        Snapshot mode: capture, compare, report, full
  -snapshot-stories string
        Comma-separated story names to execute
  -snapshot-old-tag string
        Old tag for comparison
  -snapshot-new-tag string
        New tag for comparison
  -snapshot-config string
        Snapshot configuration file path
  -snapshot-output string
        Output directory for snapshots
  -snapshot-verbose
        Enable verbose snapshot logging
`
}
