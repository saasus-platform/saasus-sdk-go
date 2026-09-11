package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// StorySnapshotFileManager handles file operations for story snapshots
type StorySnapshotFileManager struct {
	BaseDir string // Base directory for story snapshots (e.g., "tests/e2e/snapshot")
	logger  *SnapshotLogger
}

// NewStorySnapshotFileManager creates a new file manager instance
func NewStorySnapshotFileManager(baseDir string) *StorySnapshotFileManager {
	if baseDir == "" {
		baseDir = "tests/e2e/snapshot"
	}

	logger := NewSnapshotLogger(LogLevelInfo, false)

	return &StorySnapshotFileManager{
		BaseDir: baseDir,
		logger:  logger,
	}
}

// SaveStorySnapshot saves a story snapshot to file with Git tag-based naming
func (fm *StorySnapshotFileManager) SaveStorySnapshot(snapshot *StorySnapshot) (string, error) {
	if snapshot == nil {
		return "", fmt.Errorf("snapshot cannot be nil")
	}

	// Get current Git tag
	tag, err := GetGitTagForFileManager()
	if err != nil || tag == "" {
		tag = fmt.Sprintf("snapshot_%d", time.Now().Unix())
	}

	// Create filename: story_snapshot_{tag}_{story_name}.json
	filename := fmt.Sprintf("story_snapshot_%s_%s.json", tag, sanitizeStoryName(snapshot.StoryName))

	// Create directory structure
	snapshotsDir := filepath.Join(fm.BaseDir, "story_snapshots", "tags")
	if err := os.MkdirAll(snapshotsDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create snapshots directory: %w", err)
	}

	// Full file path
	filePath := filepath.Join(snapshotsDir, filename)

	// Convert to JSON
	data, err := json.MarshalIndent(snapshot, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal snapshot to JSON: %w", err)
	}

	// Write to file
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return "", fmt.Errorf("failed to write snapshot file: %w", err)
	}

	return filePath, nil
}

// LoadStorySnapshot loads a story snapshot from file
func (fm *StorySnapshotFileManager) LoadStorySnapshot(filePath string) (*StorySnapshot, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read snapshot file: %w", err)
	}

	var snapshot StorySnapshot
	if err := json.Unmarshal(data, &snapshot); err != nil {
		return nil, fmt.Errorf("failed to unmarshal snapshot JSON: %w", err)
	}

	return &snapshot, nil
}

// LoadStorySnapshotByTag loads a story snapshot by tag and story name
func (fm *StorySnapshotFileManager) LoadStorySnapshotByTag(tag, storyName string) (*StorySnapshot, error) {
	filename := fmt.Sprintf("story_snapshot_%s_%s.json", tag, sanitizeStoryName(storyName))
	filePath := filepath.Join(fm.BaseDir, "story_snapshots", "tags", filename)

	return fm.LoadStorySnapshot(filePath)
}

// GetSnapshotFilePath returns the file path for a snapshot with given tag and story name
func (fm *StorySnapshotFileManager) GetSnapshotFilePath(tag, storyName string) string {
	filename := fmt.Sprintf("story_snapshot_%s_%s.json", tag, sanitizeStoryName(storyName))
	return filepath.Join(fm.BaseDir, "story_snapshots", "tags", filename)
}

// ListAvailableSnapshots returns a list of available snapshot files
func (fm *StorySnapshotFileManager) ListAvailableSnapshots() ([]string, error) {
	snapshotsDir := filepath.Join(fm.BaseDir, "story_snapshots", "tags")

	files, err := os.ReadDir(snapshotsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil // No snapshots directory exists yet
		}
		return nil, fmt.Errorf("failed to read snapshots directory: %w", err)
	}

	var snapshots []string
	for _, file := range files {
		if !file.IsDir() && strings.HasPrefix(file.Name(), "story_snapshot_") && strings.HasSuffix(file.Name(), ".json") {
			snapshots = append(snapshots, file.Name())
		}
	}

	sort.Strings(snapshots)
	return snapshots, nil
}

// GetAvailableTagsForStory returns available tags for a specific story
func (fm *StorySnapshotFileManager) GetAvailableTagsForStory(storyName string) ([]string, error) {
	snapshots, err := fm.ListAvailableSnapshots()
	if err != nil {
		return nil, err
	}

	sanitizedStoryName := sanitizeStoryName(storyName)
	var tags []string

	for _, snapshot := range snapshots {
		// Parse filename: story_snapshot_{tag}_{story_name}.json
		if strings.HasSuffix(snapshot, "_"+sanitizedStoryName+".json") {
			// Extract tag
			prefix := "story_snapshot_"
			suffix := "_" + sanitizedStoryName + ".json"
			if len(snapshot) > len(prefix)+len(suffix) {
				tag := snapshot[len(prefix) : len(snapshot)-len(suffix)]
				tags = append(tags, tag)
			}
		}
	}

	sort.Strings(tags)
	return tags, nil
}

// SaveStoryComparison saves a story comparison result to file
func (fm *StorySnapshotFileManager) SaveStoryComparison(comparison *StoryComparison) (string, error) {
	if comparison == nil {
		return "", fmt.Errorf("comparison cannot be nil")
	}

	// Create filename: story_comparison_{story_name}_{old_tag}_vs_{new_tag}.json
	var oldTag, newTag string
	if comparison.OldSnapshot != nil && comparison.OldSnapshot.Metadata.GitTag != "" {
		oldTag = comparison.OldSnapshot.Metadata.GitTag
	} else {
		oldTag = "baseline"
	}
	if comparison.NewSnapshot != nil && comparison.NewSnapshot.Metadata.GitTag != "" {
		newTag = comparison.NewSnapshot.Metadata.GitTag
	} else {
		newTag = "current"
	}

	filename := fmt.Sprintf("story_comparison_%s_%s_vs_%s.json",
		sanitizeStoryName(comparison.StoryName),
		sanitizeTag(oldTag),
		sanitizeTag(newTag))

	// Create directory structure
	comparisonsDir := filepath.Join(fm.BaseDir, "story_comparisons")
	if err := os.MkdirAll(comparisonsDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create comparisons directory: %w", err)
	}

	// Full file path
	filePath := filepath.Join(comparisonsDir, filename)

	// Convert to JSON
	data, err := json.MarshalIndent(comparison, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal comparison to JSON: %w", err)
	}

	// Write to file
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return "", fmt.Errorf("failed to write comparison file: %w", err)
	}

	return filePath, nil
}

// LoadStoryComparison loads a story comparison from file
func (fm *StorySnapshotFileManager) LoadStoryComparison(filePath string) (*StoryComparison, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read comparison file: %w", err)
	}

	var comparison StoryComparison
	if err := json.Unmarshal(data, &comparison); err != nil {
		return nil, fmt.Errorf("failed to unmarshal comparison JSON: %w", err)
	}

	return &comparison, nil
}

// SaveStoryReport saves a story report to file (both JSON and HTML formats)
func (fm *StorySnapshotFileManager) SaveStoryReport(report *StoryReport) error {
	if report == nil {
		return fmt.Errorf("report cannot be nil")
	}

	// Get current Git tag for filename
	tag, err := GetGitTagForFileManager()
	if err != nil || tag == "" {
		tag = fmt.Sprintf("report_%d", time.Now().Unix())
	}

	// Create base filename: story_report_{story_name}_{tag}
	baseFilename := fmt.Sprintf("story_report_%s_%s", sanitizeStoryName(report.StoryName), sanitizeTag(tag))

	// Create directory structure
	reportsDir := filepath.Join(fm.BaseDir, "story_reports")
	if err := os.MkdirAll(reportsDir, 0755); err != nil {
		return fmt.Errorf("failed to create reports directory: %w", err)
	}

	// Save JSON format
	jsonFilePath := filepath.Join(reportsDir, baseFilename+".json")
	jsonData, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal report to JSON: %w", err)
	}

	if err := os.WriteFile(jsonFilePath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write JSON report file: %w", err)
	}

	// Save HTML format using StoryReporter
	htmlFilePath := filepath.Join(reportsDir, baseFilename+".html")
	reporter := NewStoryReporter(&StorySnapshotConfig{})
	if err := reporter.SaveReportAsHTML(report, htmlFilePath); err != nil {
		return fmt.Errorf("failed to write HTML report file: %w", err)
	}

	return nil
}

// LoadStoryReport loads a story report from file
func (fm *StorySnapshotFileManager) LoadStoryReport(filePath string) (*StoryReport, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read report file: %w", err)
	}

	var report StoryReport
	if err := json.Unmarshal(data, &report); err != nil {
		return nil, fmt.Errorf("failed to unmarshal report JSON: %w", err)
	}

	return &report, nil
}

// GetReportFilePath returns the file path for a report with given tag and story name
func (fm *StorySnapshotFileManager) GetReportFilePath(tag, storyName, format string) string {
	baseFilename := fmt.Sprintf("story_report_%s_%s", sanitizeStoryName(storyName), sanitizeTag(tag))

	var filename string
	switch strings.ToLower(format) {
	case "html", "htm":
		filename = baseFilename + ".html"
	case "json":
		filename = baseFilename + ".json"
	default:
		filename = baseFilename + ".json" // Default to JSON
	}

	return filepath.Join(fm.BaseDir, "story_reports", filename)
}

// ListAvailableReports returns a list of available report files
func (fm *StorySnapshotFileManager) ListAvailableReports() ([]string, error) {
	reportsDir := filepath.Join(fm.BaseDir, "story_reports")

	files, err := os.ReadDir(reportsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil // No reports directory exists yet
		}
		return nil, fmt.Errorf("failed to read reports directory: %w", err)
	}

	var reports []string
	for _, file := range files {
		if !file.IsDir() && strings.HasPrefix(file.Name(), "story_report_") {
			reports = append(reports, file.Name())
		}
	}

	sort.Strings(reports)
	return reports, nil
}

// Git utility functions (reused from existing implementation)

// GetGitTagForFileManager gets the current Git tag using the same logic as existing implementation
func GetGitTagForFileManager() (string, error) {
	// Try to get the current tag
	cmd := exec.Command("git", "describe", "--tags", "--exact-match", "HEAD")
	output, err := cmd.Output()
	if err == nil {
		return strings.TrimSpace(string(output)), nil
	}

	// If no exact tag, get the latest tag with commit info
	cmd = exec.Command("git", "describe", "--tags", "--always")
	output, err = cmd.Output()
	if err == nil {
		return strings.TrimSpace(string(output)), nil
	}

	// Fallback to commit hash
	cmd = exec.Command("git", "rev-parse", "--short", "HEAD")
	output, err = cmd.Output()
	if err != nil {
		return "unknown", fmt.Errorf("failed to get git information: %w", err)
	}

	return "dev-" + time.Now().Format("20060102-150405"), nil
}

// GetPreviousReleaseTag gets the previous release tag for story snapshots
func (fm *StorySnapshotFileManager) GetPreviousReleaseTag(currentTag, storyName string) (string, error) {
	tags, err := fm.GetAvailableTagsForStory(storyName)
	if err != nil {
		return "", err
	}

	// Filter out current tag and find the most recent previous tag
	var previousTags []string
	for _, tag := range tags {
		if tag != currentTag {
			previousTags = append(previousTags, tag)
		}
	}

	if len(previousTags) == 0 {
		return "", nil // No previous releases
	}

	// Sort and return the most recent (last in sorted order)
	sort.Strings(previousTags)
	return previousTags[len(previousTags)-1], nil
}

// Helper functions

// sanitizeStoryName sanitizes story name for use in filenames
func sanitizeStoryName(storyName string) string {
	// Replace spaces and special characters with underscores
	sanitized := strings.ReplaceAll(storyName, " ", "_")
	sanitized = strings.ReplaceAll(sanitized, "-", "_")
	sanitized = strings.ReplaceAll(sanitized, ".", "_")
	sanitized = strings.ReplaceAll(sanitized, "/", "_")
	sanitized = strings.ReplaceAll(sanitized, "\\", "_")
	sanitized = strings.ToLower(sanitized)

	// Remove multiple consecutive underscores
	for strings.Contains(sanitized, "__") {
		sanitized = strings.ReplaceAll(sanitized, "__", "_")
	}

	// Trim underscores from start and end
	sanitized = strings.Trim(sanitized, "_")

	return sanitized
}

// SaveStoryValidation saves a story validation result to file
func (fm *StorySnapshotFileManager) SaveStoryValidation(validation *StoryValidation) error {
	if validation == nil {
		return fmt.Errorf("validation cannot be nil")
	}

	validationsDir := filepath.Join(fm.BaseDir, "story_validations")
	if err := os.MkdirAll(validationsDir, 0755); err != nil {
		return fmt.Errorf("failed to create validations directory: %w", err)
	}

	tag, err := GetGitTagForFileManager()
	if err != nil || tag == "" {
		tag = fmt.Sprintf("validation_%d", time.Now().Unix())
	}
	fileName := fmt.Sprintf("story_validation_%s_%s.json",
		sanitizeStoryName(validation.StoryName), sanitizeTag(tag))
	currentPath := filepath.Join(validationsDir, fileName)

	previousPath, previousValidation, err := fm.loadLatestValidation(validationsDir, validation.StoryName)
	if err != nil {
		return err
	}

	if previousValidation != nil {
		validation.Comparison = generateValidationComparison(previousValidation, validation, filepath.Base(previousPath))
	}

	data, err := json.MarshalIndent(validation, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal validation to JSON: %w", err)
	}

	if err := os.WriteFile(currentPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write validation file: %w", err)
	}

	if err := fm.retainLatestValidations(validationsDir, validation.StoryName, 2); err != nil {
		return err
	}

	return nil
}

// LoadStoryValidation loads a story validation from file
func (fm *StorySnapshotFileManager) LoadStoryValidation(filePath string) (*StoryValidation, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read validation file: %w", err)
	}

	var validation StoryValidation
	if err := json.Unmarshal(data, &validation); err != nil {
		return nil, fmt.Errorf("failed to unmarshal validation JSON: %w", err)
	}

	return &validation, nil
}

// GetValidationFilePath returns the file path for a validation with given tag and story name
func (fm *StorySnapshotFileManager) GetValidationFilePath(tag, storyName string) string {
	return filepath.Join(fm.BaseDir, "story_validations",
		fmt.Sprintf("story_validation_%s_%s.json", sanitizeStoryName(storyName), sanitizeTag(tag)))
}

// ListAvailableValidations returns a list of available validation files
func (fm *StorySnapshotFileManager) ListAvailableValidations() ([]string, error) {
	validationsDir := filepath.Join(fm.BaseDir, "story_validations")

	files, err := os.ReadDir(validationsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil // No validations directory exists yet
		}
		return nil, fmt.Errorf("failed to read validations directory: %w", err)
	}

	var validations []string
	for _, file := range files {
		if !file.IsDir() && strings.HasPrefix(file.Name(), "story_validation_") {
			validations = append(validations, file.Name())
		}
	}

	sort.Strings(validations)
	return validations, nil
}

// sanitizeTag sanitizes Git tag for use in filenames
func sanitizeTag(tag string) string {
	// Replace problematic characters
	sanitized := strings.ReplaceAll(tag, "/", "_")
	sanitized = strings.ReplaceAll(sanitized, "\\", "_")
	sanitized = strings.ReplaceAll(sanitized, ":", "_")
	return sanitized
}

func (fm *StorySnapshotFileManager) loadLatestValidation(dir, storyName string) (string, *StoryValidation, error) {
	files, err := fm.sortedValidationFiles(dir, storyName)
	if err != nil {
		return "", nil, err
	}
	if len(files) == 0 {
		return "", nil, nil
	}
	latest := files[0]
	data, err := os.ReadFile(filepath.Join(dir, latest.name))
	if err != nil {
		return "", nil, fmt.Errorf("failed to read previous validation: %w", err)
	}
	var prev StoryValidation
	if err := json.Unmarshal(data, &prev); err != nil {
		return "", nil, fmt.Errorf("failed to unmarshal previous validation: %w", err)
	}
	return filepath.Join(dir, latest.name), &prev, nil
}

func (fm *StorySnapshotFileManager) retainLatestValidations(dir, storyName string, keep int) error {
	files, err := fm.sortedValidationFiles(dir, storyName)
	if err != nil {
		return err
	}
	for idx := keep; idx < len(files); idx++ {
		path := filepath.Join(dir, files[idx].name)
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to remove old validation '%s': %w", files[idx].name, err)
		}
	}
	return nil
}

type validationFileInfo struct {
	name    string
	modTime time.Time
}

func (fm *StorySnapshotFileManager) sortedValidationFiles(dir, storyName string) ([]validationFileInfo, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read validations directory: %w", err)
	}
	prefix := fmt.Sprintf("story_validation_%s_", sanitizeStoryName(storyName))
	var files []validationFileInfo
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasPrefix(entry.Name(), prefix) {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		files = append(files, validationFileInfo{name: entry.Name(), modTime: info.ModTime()})
	}
	sort.Slice(files, func(i, j int) bool {
		return files[i].modTime.After(files[j].modTime)
	})
	return files, nil
}
