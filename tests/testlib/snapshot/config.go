package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// DefaultStorySnapshotConfig returns a default configuration for story snapshots
func DefaultStorySnapshotConfig() *StorySnapshotConfig {
	return &StorySnapshotConfig{
		EnableCapture:            true,
		EnableComparison:         false, // Disabled by default for capture-only mode
		EnableReporting:          false, // Disabled by default for capture-only mode
		EnableValidation:         true,
		EnableCompatibilityCheck: true, // 要件 4.1, 4.2: デフォルトで後方互換性チェックを有効化
		ComparisonMode:           "release",
		OutputDirectory:          "tests/e2e/snapshot",
		ModuleName:               "",
		SnapshotOnly:             false,
		FileNameFormat:           "story_snapshot_{tag}_{story_name}.json",
		CaptureLevel:             CaptureLevelFull,
		ValidationRules:          DefaultValidationRules(),
	}
}

// DefaultValidationRules returns default validation rules
func DefaultValidationRules() []ValidationRule {
	return []ValidationRule{
		{
			Type:     ValidationRuleCompletion,
			Enabled:  true,
			Severity: ValidationSeverityError,
		},
		{
			Type:     ValidationRuleSequence,
			Enabled:  true,
			Severity: ValidationSeverityError,
		},
		{
			Type:     ValidationRuleStateTransition,
			Enabled:  true,
			Severity: ValidationSeverityWarning,
		},
		{
			Type:     ValidationRuleTiming,
			Enabled:  false, // Disabled by default as timing can be variable
			Severity: ValidationSeverityInfo,
		},
	}
}

// LoadStorySnapshotConfig loads configuration from a file
func LoadStorySnapshotConfig(configPath string) (*StorySnapshotConfig, error) {
	// Start with default config
	config := DefaultStorySnapshotConfig()

	// If no config file specified, return default
	if configPath == "" {
		return config, nil
	}

	// Check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return config, nil // Return default if file doesn't exist
	}

	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", configPath, err)
	}

	// Parse JSON
	if err := json.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse config file %s: %w", configPath, err)
	}

	// Validate config
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config in %s: %w", configPath, err)
	}

	return config, nil
}

// SaveStorySnapshotConfig saves configuration to a file
func SaveStorySnapshotConfig(config *StorySnapshotConfig, configPath string) error {
	// Validate config before saving
	if err := config.Validate(); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file %s: %w", configPath, err)
	}

	return nil
}

// Validate validates the configuration
func (c *StorySnapshotConfig) Validate() error {
	// Validate comparison mode
	validModes := map[string]bool{
		"release": true,
		"manual":  true,
		"skip":    true,
	}
	if !validModes[c.ComparisonMode] {
		return fmt.Errorf("invalid comparison mode: %s", c.ComparisonMode)
	}

	// Validate capture level
	validLevels := map[CaptureLevel]bool{
		CaptureLevelFull:     true,
		CaptureLevelStory:    true,
		CaptureLevelStep:     true,
		CaptureLevelResponse: true,
	}
	if !validLevels[c.CaptureLevel] {
		return fmt.Errorf("invalid capture level: %s", c.CaptureLevel)
	}

	// Validate output directory
	if c.OutputDirectory == "" {
		return fmt.Errorf("output directory cannot be empty")
	}

	// Validate file name format
	if c.FileNameFormat == "" {
		c.FileNameFormat = "story_snapshot_{tag}_{story_name}.json"
	}

	return nil
}

func (c *StorySnapshotConfig) getEffectiveOutputDirectory() string {
	base := c.OutputDirectory
	if base == "" {
		base = "tests/e2e/snapshot"
	}

	if slug := c.ModuleSlug(); slug != "" {
		return filepath.Join(base, slug)
	}

	return base
}

// GetModuleOutputDirectory resolves the effective output directory including module subdirectory if configured.
func (c *StorySnapshotConfig) GetModuleOutputDirectory() string {
	return c.getEffectiveOutputDirectory()
}

// ModuleSlug returns the normalized module identifier based on ModuleName.
func (c *StorySnapshotConfig) ModuleSlug() string {
	return slugifyModuleName(c.ModuleName)
}

// GetSnapshotDirectory returns the directory for storing snapshots
func (c *StorySnapshotConfig) GetSnapshotDirectory() string {
	return filepath.Join(c.getEffectiveOutputDirectory(), "story_snapshots", "tags")
}

// GetComparisonDirectory returns the directory for storing comparison results
func (c *StorySnapshotConfig) GetComparisonDirectory() string {
	return filepath.Join(c.getEffectiveOutputDirectory(), "story_comparisons")
}

// GetValidationDirectory returns the directory for storing validation results
func (c *StorySnapshotConfig) GetValidationDirectory() string {
	return filepath.Join(c.getEffectiveOutputDirectory(), "story_validations")
}

// GetReportDirectory returns the directory for storing reports
func (c *StorySnapshotConfig) GetReportDirectory() string {
	return filepath.Join(c.getEffectiveOutputDirectory(), "story_reports")
}

// GetSnapshotFileName generates a snapshot file name based on the format
func (c *StorySnapshotConfig) GetSnapshotFileName(tag, storyName string) string {
	fileName := c.FileNameFormat
	fileName = replaceToken(fileName, "{tag}", tag)
	fileName = replaceToken(fileName, "{story_name}", sanitizeFileName(storyName))
	return fileName
}

// GetComparisonFileName generates a comparison file name
func (c *StorySnapshotConfig) GetComparisonFileName(storyName, oldTag, newTag string) string {
	if oldTag != "" && newTag != "" {
		return fmt.Sprintf("story_comparison_%s_%s_vs_%s.json",
			sanitizeFileName(storyName), oldTag, newTag)
	}
	return fmt.Sprintf("story_comparison_%s_release.json", sanitizeFileName(storyName))
}

// GetValidationFileName generates a validation file name
func (c *StorySnapshotConfig) GetValidationFileName(storyName, tag string) string {
	return fmt.Sprintf("story_validation_%s_%s.json", sanitizeFileName(storyName), tag)
}

// GetReportFileName generates a report file name
func (c *StorySnapshotConfig) GetReportFileName(storyName, oldTag, newTag string) string {
	if oldTag != "" && newTag != "" {
		return fmt.Sprintf("story_report_%s_%s_vs_%s.html",
			sanitizeFileName(storyName), oldTag, newTag)
	}
	return fmt.Sprintf("story_report_%s_release.html", sanitizeFileName(storyName))
}

// IsValidationRuleEnabled checks if a validation rule is enabled
func (c *StorySnapshotConfig) IsValidationRuleEnabled(ruleType ValidationRuleType) bool {
	for _, rule := range c.ValidationRules {
		if rule.Type == ruleType {
			return rule.Enabled
		}
	}
	return false
}

// GetValidationRuleSeverity gets the severity for a validation rule
func (c *StorySnapshotConfig) GetValidationRuleSeverity(ruleType ValidationRuleType) ValidationSeverity {
	for _, rule := range c.ValidationRules {
		if rule.Type == ruleType {
			return rule.Severity
		}
	}
	return ValidationSeverityError // Default severity
}

func slugifyModuleName(name string) string {
	name = strings.ToLower(strings.TrimSpace(name))
	if name == "" {
		return ""
	}

	var builder strings.Builder
	lastHyphen := false

	for _, r := range name {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			builder.WriteRune(r)
			lastHyphen = false
		case r == ' ' || r == '-' || r == '_' || r == '/':
			if !lastHyphen && builder.Len() > 0 {
				builder.WriteRune('-')
				lastHyphen = true
			}
		default:
			// skip unsupported characters
		}
	}

	return strings.Trim(builder.String(), "-")
}

// Helper functions

// replaceToken replaces a token in a string
func replaceToken(str, token, replacement string) string {
	// Simple string replacement - could be enhanced with regex if needed
	result := ""
	tokenLen := len(token)
	for i := 0; i < len(str); {
		if i+tokenLen <= len(str) && str[i:i+tokenLen] == token {
			result += replacement
			i += tokenLen
		} else {
			result += string(str[i])
			i++
		}
	}
	return result
}

// sanitizeFileName sanitizes a string to be safe for use as a filename
func sanitizeFileName(name string) string {
	// Replace spaces and special characters with underscores
	result := ""
	for _, r := range name {
		switch r {
		case ' ', '-', '/', '\\', ':', '*', '?', '"', '<', '>', '|', '.':
			result += "_"
		default:
			result += string(r)
		}
	}

	// Convert to lowercase for consistency
	result = toLower(result)

	// Remove multiple consecutive underscores
	for i := 0; i < len(result)-1; {
		if result[i] == '_' && result[i+1] == '_' {
			result = result[:i] + result[i+1:]
		} else {
			i++
		}
	}

	// Trim underscores from start and end
	result = trimUnderscores(result)

	return result
}

// toLower converts a string to lowercase
func toLower(s string) string {
	result := ""
	for _, r := range s {
		if r >= 'A' && r <= 'Z' {
			result += string(r + 32)
		} else {
			result += string(r)
		}
	}
	return result
}

// trimUnderscores removes leading and trailing underscores
func trimUnderscores(s string) string {
	// Trim leading underscores
	for len(s) > 0 && s[0] == '_' {
		s = s[1:]
	}
	// Trim trailing underscores
	for len(s) > 0 && s[len(s)-1] == '_' {
		s = s[:len(s)-1]
	}
	return s
}
