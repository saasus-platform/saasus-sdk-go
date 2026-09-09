package snapshot

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// StoryReporter handles story-level report generation
type StoryReporter struct {
	config *StorySnapshotConfig
	logger *SnapshotLogger
}

// NewStoryReporter creates a new story reporter
func NewStoryReporter(config *StorySnapshotConfig) *StoryReporter {
	// Create logger with appropriate level
	logLevel := LogLevelInfo
	debugMode := false
	if config != nil && config.CaptureLevel == CaptureLevelFull {
		logLevel = LogLevelDebug
		debugMode = true
	}

	logger := NewSnapshotLogger(logLevel, debugMode)

	return &StoryReporter{
		config: config,
		logger: logger,
	}
}

// GenerateStoryReport generates a story report from comparison result
func (sr *StoryReporter) GenerateStoryReport(comparison *StoryComparison) (*StoryReport, error) {
	if comparison == nil {
		return nil, NewReportError(ErrorTypeValidation,
			"comparison cannot be nil", nil, "", "")
	}

	sr.logger.LogPhaseStart(PhaseReport, fmt.Sprintf("story '%s'", comparison.StoryName))
	sr.logger.SetContext("story", comparison.StoryName)
	defer sr.logger.ClearContext()

	report := &StoryReport{
		Title:       fmt.Sprintf("Story Snapshot Report: %s", comparison.StoryName),
		GeneratedAt: time.Now(),
		StoryName:   comparison.StoryName,
	}

	// Generate execution summary
	if comparison.NewSnapshot != nil {
		report.ExecutionSummary = comparison.NewSnapshot.Summary
	}

	// Copy compatibility analysis
	report.CompatibilityAnalysis = comparison.CompatibilityReport

	// Generate step analysis
	report.StepAnalysis = sr.generateStepAnalysis(comparison)

	// Generate recommendations
	report.Recommendations = sr.generateRecommendations(comparison)

	// Generate conclusion
	report.Conclusion = sr.generateConclusion(comparison)

	return report, nil
}

// GenerateReportFromFile generates a report from a comparison file
func (sr *StoryReporter) GenerateReportFromFile(comparisonPath string) (*StoryReport, error) {
	sr.logger.LogInfo("Generating report from comparison file: %s", comparisonPath)

	comparison, err := sr.LoadComparisonResult(comparisonPath)
	if err != nil {
		snapErr := NewReportError(ErrorTypeFileIO,
			"failed to load comparison result", err, "", comparisonPath)
		sr.logger.LogSnapshotError(snapErr)
		return nil, snapErr
	}

	return sr.GenerateStoryReport(comparison)
}

// LoadComparisonResult loads comparison result from file
func (sr *StoryReporter) LoadComparisonResult(filePath string) (*StoryComparison, error) {
	sr.logger.LogFileOperation("read", filePath)

	data, err := os.ReadFile(filePath)
	if err != nil {
		snapErr := NewFileIOError("failed to read comparison file", err, filePath)
		sr.logger.LogSnapshotError(snapErr)
		return nil, snapErr
	}

	var comparison StoryComparison
	if err := json.Unmarshal(data, &comparison); err != nil {
		snapErr := NewReportError(ErrorTypeParsing,
			"failed to unmarshal comparison", err, "", filePath)
		sr.logger.LogSnapshotError(snapErr)
		return nil, snapErr
	}

	sr.logger.LogDebug("Successfully loaded comparison for story: %s", comparison.StoryName)
	return &comparison, nil
}

// SaveReport saves report to file (JSON and HTML formats)
func (sr *StoryReporter) SaveReport(report *StoryReport, filePath string) error {
	if filePath != "" {
		// Use provided file path
		return sr.saveReportToPath(report, filePath)
	}

	// Use file manager scoped to the module output directory for automatic naming
	baseDir := ""
	if sr.config != nil {
		baseDir = sr.config.GetModuleOutputDirectory()
	}
	fileManager := NewStorySnapshotFileManager(baseDir)
	return fileManager.SaveStoryReport(report)
}

// SaveReportAsHTML saves report as HTML file
func (sr *StoryReporter) SaveReportAsHTML(report *StoryReport, filePath string) error {
	htmlContent, err := sr.generateHTMLReport(report)
	if err != nil {
		return fmt.Errorf("failed to generate HTML report: %w", err)
	}

	if err := os.WriteFile(filePath, []byte(htmlContent), 0644); err != nil {
		return fmt.Errorf("failed to write HTML report: %w", err)
	}

	return nil
}

// SaveReportAsJSON saves report as JSON file
func (sr *StoryReporter) SaveReportAsJSON(report *StoryReport, filePath string) error {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal report: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write JSON report: %w", err)
	}

	return nil
}

// generateStepAnalysis generates analysis for individual steps
func (sr *StoryReporter) generateStepAnalysis(comparison *StoryComparison) []StepAnalysis {
	var stepAnalyses []StepAnalysis

	// Create a map of step differences for easy lookup
	stepDifferences := make(map[string][]StoryDifference)
	for _, diff := range comparison.Differences {
		if diff.StepName != "" {
			stepDifferences[diff.StepName] = append(stepDifferences[diff.StepName], diff)
		}
	}

	// Analyze steps from new snapshot
	if comparison.NewSnapshot != nil {
		for _, step := range comparison.NewSnapshot.Steps {
			analysis := StepAnalysis{
				StepName:    step.StepName,
				Method:      step.Method,
				Status:      sr.determineStepStatus(step.Success),
				Differences: stepDifferences[step.StepName],
				Impact:      sr.determineStepImpact(stepDifferences[step.StepName]),
				Description: sr.generateStepDescription(step, stepDifferences[step.StepName]),
			}
			stepAnalyses = append(stepAnalyses, analysis)
		}
	}

	// Add analysis for removed steps
	if comparison.OldSnapshot != nil {
		for _, step := range comparison.OldSnapshot.Steps {
			// Check if step exists in new snapshot
			exists := false
			if comparison.NewSnapshot != nil {
				for _, newStep := range comparison.NewSnapshot.Steps {
					if newStep.StepName == step.StepName {
						exists = true
						break
					}
				}
			}

			if !exists {
				analysis := StepAnalysis{
					StepName:    step.StepName,
					Method:      step.Method,
					Status:      TestStatusFailed, // Removed step is considered failed
					Differences: stepDifferences[step.StepName],
					Impact:      Breaking,
					Description: fmt.Sprintf("Step %s was removed from the story", step.StepName),
				}
				stepAnalyses = append(stepAnalyses, analysis)
			}
		}
	}

	return stepAnalyses
}

// generateRecommendations generates actionable recommendations
func (sr *StoryReporter) generateRecommendations(comparison *StoryComparison) []Recommendation {
	var recommendations []Recommendation

	// Analyze compatibility issues
	if comparison.CompatibilityReport.Level == Breaking {
		recommendations = append(recommendations, Recommendation{
			Type:        RecommendationTypeCompatibility,
			Priority:    RecommendationPriorityHigh,
			Title:       "Breaking Changes Detected",
			Description: "This story contains breaking changes that may affect SDK users",
			Action:      "Review breaking changes and consider providing migration guide or deprecation notices",
		})
	}

	// Analyze step failures
	failedSteps := 0
	if comparison.NewSnapshot != nil {
		for _, step := range comparison.NewSnapshot.Steps {
			if !step.Success {
				failedSteps++
			}
		}
	}

	if failedSteps > 0 {
		recommendations = append(recommendations, Recommendation{
			Type:        RecommendationTypeCompatibility,
			Priority:    RecommendationPriorityHigh,
			Title:       "Step Failures Detected",
			Description: fmt.Sprintf("%d step(s) failed in this story execution", failedSteps),
			Action:      "Investigate failed steps and fix underlying issues before release",
		})
	}

	// Analyze performance issues
	if comparison.NewSnapshot != nil && comparison.OldSnapshot != nil {
		oldDuration := comparison.OldSnapshot.Duration.Seconds()
		newDuration := comparison.NewSnapshot.Duration.Seconds()
		if oldDuration > 0 && newDuration > oldDuration*1.5 { // 50% slower
			recommendations = append(recommendations, Recommendation{
				Type:        RecommendationTypePerformance,
				Priority:    RecommendationPriorityMedium,
				Title:       "Performance Regression Detected",
				Description: fmt.Sprintf("Story execution time increased from %.2fs to %.2fs", oldDuration, newDuration),
				Action:      "Investigate performance regression and optimize slow operations",
			})
		}
	}

	// Analyze step sequence changes
	for _, change := range comparison.CompatibilityReport.StepSequenceChanges {
		if change.Impact == Breaking {
			recommendations = append(recommendations, Recommendation{
				Type:        RecommendationTypeCompatibility,
				Priority:    RecommendationPriorityHigh,
				Title:       "Step Sequence Changes",
				Description: change.Description,
				Action:      "Ensure step sequence changes are intentional and document any workflow modifications",
			})
		}
	}

	// Add general best practice recommendations
	if len(recommendations) == 0 {
		recommendations = append(recommendations, Recommendation{
			Type:        RecommendationTypeBestPractice,
			Priority:    RecommendationPriorityLow,
			Title:       "Story Execution Successful",
			Description: "No critical issues detected in this story execution",
			Action:      "Continue monitoring story execution in future releases",
		})
	}

	return recommendations
}

// generateConclusion generates conclusion based on analysis
func (sr *StoryReporter) generateConclusion(comparison *StoryComparison) string {
	if comparison.CompatibilityReport.Level == Breaking {
		return "FAILED: Breaking changes detected that require immediate attention before release."
	}

	if comparison.CompatibilityReport.Level == Warning {
		return "WARNING: Some compatibility issues detected. Review recommended actions before release."
	}

	if comparison.NewSnapshot != nil && comparison.NewSnapshot.Status != TestStatusPassed {
		return "FAILED: Story execution failed. Fix issues before proceeding with release."
	}

	return "PASSED: Story execution completed successfully with no critical compatibility issues."
}

// generateHTMLReport generates HTML report content
func (sr *StoryReporter) generateHTMLReport(report *StoryReport) (string, error) {
	tmpl := template.New("story_report").Funcs(template.FuncMap{
		"formatDuration":   sr.formatDuration,
		"formatTime":       sr.formatTime,
		"getStatusClass":   sr.getStatusClass,
		"getImpactClass":   sr.getImpactClass,
		"getPriorityClass": sr.getPriorityClass,
		"formatJSON":       sr.formatJSON,
		"contains":         strings.Contains,
	})

	htmlTemplate := sr.getHTMLTemplate()
	tmpl, err := tmpl.Parse(htmlTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML template: %w", err)
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, report); err != nil {
		return "", fmt.Errorf("failed to execute HTML template: %w", err)
	}

	return buf.String(), nil
}

// Helper methods

// saveReportToPath saves report to specified path (determines format by extension)
func (sr *StoryReporter) saveReportToPath(report *StoryReport, filePath string) error {
	// Ensure directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".html", ".htm":
		return sr.SaveReportAsHTML(report, filePath)
	case ".json":
		return sr.SaveReportAsJSON(report, filePath)
	default:
		// Default to JSON if no extension or unknown extension
		jsonPath := filePath + ".json"
		htmlPath := filePath + ".html"

		if err := sr.SaveReportAsJSON(report, jsonPath); err != nil {
			return err
		}
		return sr.SaveReportAsHTML(report, htmlPath)
	}
}

// determineStepStatus determines test status from step success
func (sr *StoryReporter) determineStepStatus(success bool) TestStatus {
	if success {
		return TestStatusPassed
	}
	return TestStatusFailed
}

// determineStepImpact determines impact level from step differences
func (sr *StoryReporter) determineStepImpact(differences []StoryDifference) CompatibilityLevel {
	maxImpact := Compatible
	for _, diff := range differences {
		if diff.Impact > maxImpact {
			maxImpact = diff.Impact
		}
	}
	return maxImpact
}

// generateStepDescription generates description for step analysis
func (sr *StoryReporter) generateStepDescription(step StepSnapshot, differences []StoryDifference) string {
	if len(differences) == 0 {
		if step.Success {
			return fmt.Sprintf("Step executed successfully with status code %d", step.StatusCode)
		}
		return fmt.Sprintf("Step failed with status code %d", step.StatusCode)
	}

	var descriptions []string
	for _, diff := range differences {
		descriptions = append(descriptions, diff.Description)
	}

	return strings.Join(descriptions, "; ")
}

// Template helper functions

// formatDuration formats duration for display
func (sr *StoryReporter) formatDuration(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%.0fms", float64(d.Nanoseconds())/1e6)
	}
	return fmt.Sprintf("%.2fs", d.Seconds())
}

// formatTime formats time for display
func (sr *StoryReporter) formatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05 MST")
}

// getStatusClass returns CSS class for test status
func (sr *StoryReporter) getStatusClass(status TestStatus) string {
	switch status {
	case TestStatusPassed:
		return "status-passed"
	case TestStatusFailed:
		return "status-failed"
	case TestStatusSkipped:
		return "status-skipped"
	default:
		return "status-unknown"
	}
}

// getImpactClass returns CSS class for compatibility impact
func (sr *StoryReporter) getImpactClass(impact CompatibilityLevel) string {
	switch impact {
	case Compatible:
		return "impact-compatible"
	case Warning:
		return "impact-warning"
	case Breaking:
		return "impact-breaking"
	default:
		return "impact-unknown"
	}
}

// getPriorityClass returns CSS class for recommendation priority
func (sr *StoryReporter) getPriorityClass(priority RecommendationPriority) string {
	switch priority {
	case RecommendationPriorityHigh:
		return "priority-high"
	case RecommendationPriorityMedium:
		return "priority-medium"
	case RecommendationPriorityLow:
		return "priority-low"
	default:
		return "priority-unknown"
	}
}

// formatJSON formats JSON for display
func (sr *StoryReporter) formatJSON(data interface{}) string {
	if data == nil {
		return "null"
	}

	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Sprintf("Error formatting JSON: %v", err)
	}

	return string(jsonBytes)
}

// getHTMLTemplate returns the HTML template for story reports
func (sr *StoryReporter) getHTMLTemplate() string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 1200px;
            margin: 0 auto;
            padding: 20px;
            background-color: #f5f5f5;
        }

        .container {
            background: white;
            border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
            padding: 30px;
            margin-bottom: 20px;
        }

        h1, h2, h3 {
            color: #2c3e50;
            margin-top: 0;
        }

        h1 {
            border-bottom: 3px solid #3498db;
            padding-bottom: 10px;
        }

        h2 {
            border-bottom: 2px solid #ecf0f1;
            padding-bottom: 8px;
            margin-top: 30px;
        }

        .header-info {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 20px;
            margin: 20px 0;
            padding: 20px;
            background-color: #f8f9fa;
            border-radius: 6px;
        }

        .info-item {
            text-align: center;
        }

        .info-label {
            font-weight: bold;
            color: #7f8c8d;
            font-size: 0.9em;
            text-transform: uppercase;
        }

        .info-value {
            font-size: 1.2em;
            margin-top: 5px;
        }

        .status-passed { color: #27ae60; font-weight: bold; }
        .status-failed { color: #e74c3c; font-weight: bold; }
        .status-skipped { color: #f39c12; font-weight: bold; }
        .status-unknown { color: #95a5a6; font-weight: bold; }

        .impact-compatible { color: #27ae60; }
        .impact-warning { color: #f39c12; }
        .impact-breaking { color: #e74c3c; }
        .impact-unknown { color: #95a5a6; }

        .priority-high {
            background-color: #e74c3c;
            color: white;
            padding: 2px 8px;
            border-radius: 4px;
            font-size: 0.8em;
        }
        .priority-medium {
            background-color: #f39c12;
            color: white;
            padding: 2px 8px;
            border-radius: 4px;
            font-size: 0.8em;
        }
        .priority-low {
            background-color: #27ae60;
            color: white;
            padding: 2px 8px;
            border-radius: 4px;
            font-size: 0.8em;
        }

        .summary-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
            gap: 15px;
            margin: 20px 0;
        }

        .summary-card {
            background: #f8f9fa;
            padding: 15px;
            border-radius: 6px;
            text-align: center;
            border-left: 4px solid #3498db;
        }

        .summary-number {
            font-size: 2em;
            font-weight: bold;
            color: #2c3e50;
        }

        .summary-label {
            color: #7f8c8d;
            font-size: 0.9em;
            text-transform: uppercase;
        }

        .step-table {
            width: 100%;
            border-collapse: collapse;
            margin: 20px 0;
        }

        .step-table th,
        .step-table td {
            padding: 12px;
            text-align: left;
            border-bottom: 1px solid #ecf0f1;
        }

        .step-table th {
            background-color: #f8f9fa;
            font-weight: bold;
            color: #2c3e50;
        }

        .step-table tr:hover {
            background-color: #f8f9fa;
        }

        .recommendation {
            margin: 15px 0;
            padding: 15px;
            border-radius: 6px;
            border-left: 4px solid #3498db;
        }

        .recommendation.high {
            border-left-color: #e74c3c;
            background-color: #fdf2f2;
        }

        .recommendation.medium {
            border-left-color: #f39c12;
            background-color: #fef9e7;
        }

        .recommendation.low {
            border-left-color: #27ae60;
            background-color: #f0f9f0;
        }

        .recommendation-title {
            font-weight: bold;
            margin-bottom: 8px;
        }

        .recommendation-action {
            font-style: italic;
            color: #7f8c8d;
            margin-top: 8px;
        }

        .conclusion {
            padding: 20px;
            border-radius: 6px;
            font-weight: bold;
            text-align: center;
            margin: 20px 0;
        }

        .conclusion.passed {
            background-color: #d4edda;
            color: #155724;
            border: 1px solid #c3e6cb;
        }

        .conclusion.failed {
            background-color: #f8d7da;
            color: #721c24;
            border: 1px solid #f5c6cb;
        }

        .conclusion.warning {
            background-color: #fff3cd;
            color: #856404;
            border: 1px solid #ffeaa7;
        }

        .code-block {
            background-color: #f4f4f4;
            border: 1px solid #ddd;
            border-radius: 4px;
            padding: 10px;
            font-family: 'Courier New', monospace;
            font-size: 0.9em;
            overflow-x: auto;
            white-space: pre-wrap;
        }

        .expandable {
            cursor: pointer;
            user-select: none;
        }

        .expandable:hover {
            background-color: #f0f0f0;
        }

        .expandable-content {
            display: none;
            margin-top: 10px;
        }

        .expandable.expanded .expandable-content {
            display: block;
        }

        .badge {
            display: inline-block;
            padding: 2px 8px;
            border-radius: 12px;
            font-size: 0.8em;
            font-weight: bold;
        }

        .footer {
            text-align: center;
            color: #7f8c8d;
            font-size: 0.9em;
            margin-top: 40px;
            padding-top: 20px;
            border-top: 1px solid #ecf0f1;
        }
    </style>
    <script>
        function toggleExpandable(element) {
            element.classList.toggle('expanded');
        }

        document.addEventListener('DOMContentLoaded', function() {
            // Add click handlers for expandable sections
            document.querySelectorAll('.expandable').forEach(function(element) {
                element.addEventListener('click', function() {
                    toggleExpandable(this);
                });
            });
        });
    </script>
</head>
<body>
    <div class="container">
        <h1>{{.Title}}</h1>

        <div class="header-info">
            <div class="info-item">
                <div class="info-label">Story Name</div>
                <div class="info-value">{{.StoryName}}</div>
            </div>
            <div class="info-item">
                <div class="info-label">Generated At</div>
                <div class="info-value">{{formatTime .GeneratedAt}}</div>
            </div>
            <div class="info-item">
                <div class="info-label">Compatibility Level</div>
                <div class="info-value {{getImpactClass .CompatibilityAnalysis.Level}}">
                    {{.CompatibilityAnalysis.Level}}
                </div>
            </div>
            <div class="info-item">
                <div class="info-label">Overall Status</div>
                <div class="info-value">
                    {{if .CompatibilityAnalysis.Passed}}
                        <span class="status-passed">PASSED</span>
                    {{else}}
                        <span class="status-failed">FAILED</span>
                    {{end}}
                </div>
            </div>
        </div>
    </div>

    <div class="container">
        <h2>Execution Summary</h2>
        <div class="summary-grid">
            <div class="summary-card">
                <div class="summary-number">{{.ExecutionSummary.TotalSteps}}</div>
                <div class="summary-label">Total Steps</div>
            </div>
            <div class="summary-card">
                <div class="summary-number">{{.ExecutionSummary.SuccessfulSteps}}</div>
                <div class="summary-label">Successful</div>
            </div>
            <div class="summary-card">
                <div class="summary-number">{{.ExecutionSummary.FailedSteps}}</div>
                <div class="summary-label">Failed</div>
            </div>
            <div class="summary-card">
                <div class="summary-number">{{formatDuration .ExecutionSummary.TotalDuration}}</div>
                <div class="summary-label">Total Duration</div>
            </div>
            <div class="summary-card">
                <div class="summary-number">{{formatDuration .ExecutionSummary.AverageStepDuration}}</div>
                <div class="summary-label">Avg Step Duration</div>
            </div>
        </div>
    </div>

    <div class="container">
        <h2>Compatibility Analysis</h2>
        <p><strong>Summary:</strong> {{.CompatibilityAnalysis.Summary}}</p>

        {{if .CompatibilityAnalysis.StoryFlowChanges}}
        <h3>Story Flow Changes</h3>
        {{range .CompatibilityAnalysis.StoryFlowChanges}}
        <div class="recommendation {{if eq .Impact 2}}high{{else if eq .Impact 1}}medium{{else}}low{{end}}">
            <div class="recommendation-title">{{.Type}}: {{.Description}}</div>
            <div class="{{getImpactClass .Impact}}">Impact: {{.Impact}}</div>
        </div>
        {{end}}
        {{end}}

        {{if .CompatibilityAnalysis.StepSequenceChanges}}
        <h3>Step Sequence Changes</h3>
        {{range .CompatibilityAnalysis.StepSequenceChanges}}
        <div class="recommendation {{if eq .Impact 2}}high{{else if eq .Impact 1}}medium{{else}}low{{end}}">
            <div class="recommendation-title">{{.Type}}: {{.Description}}</div>
            {{if .StepName}}<div><strong>Step:</strong> {{.StepName}}</div>{{end}}
            <div class="{{getImpactClass .Impact}}">Impact: {{.Impact}}</div>
        </div>
        {{end}}
        {{end}}

        {{if .CompatibilityAnalysis.StateTransitionChanges}}
        <h3>State Transition Changes</h3>
        {{range .CompatibilityAnalysis.StateTransitionChanges}}
        <div class="recommendation {{if eq .Impact 2}}high{{else if eq .Impact 1}}medium{{else}}low{{end}}">
            <div class="recommendation-title">{{.Type}}: {{.Description}}</div>
            {{if .FromStep}}<div><strong>From:</strong> {{.FromStep}} → <strong>To:</strong> {{.ToStep}}</div>{{end}}
            <div class="{{getImpactClass .Impact}}">Impact: {{.Impact}}</div>
        </div>
        {{end}}
        {{end}}
    </div>

    <div class="container">
        <h2>Step Analysis</h2>
        <table class="step-table">
            <thead>
                <tr>
                    <th>Step Name</th>
                    <th>Method</th>
                    <th>Status</th>
                    <th>Impact</th>
                    <th>Description</th>
                </tr>
            </thead>
            <tbody>
                {{range .StepAnalysis}}
                <tr>
                    <td>{{.StepName}}</td>
                    <td><code>{{.Method}}</code></td>
                    <td><span class="{{getStatusClass .Status}}">{{.Status}}</span></td>
                    <td><span class="{{getImpactClass .Impact}}">{{.Impact}}</span></td>
                    <td>{{.Description}}</td>
                </tr>
                {{end}}
            </tbody>
        </table>
    </div>

    <div class="container">
        <h2>Recommendations</h2>
        {{range .Recommendations}}
        <div class="recommendation {{if eq .Priority "high"}}high{{else if eq .Priority "medium"}}medium{{else}}low{{end}}">
            <div class="recommendation-title">
                {{.Title}}
                <span class="{{getPriorityClass .Priority}}">{{.Priority}}</span>
            </div>
            <div>{{.Description}}</div>
            <div class="recommendation-action"><strong>Action:</strong> {{.Action}}</div>
        </div>
        {{end}}
    </div>

    <div class="container">
        <h2>Conclusion</h2>
        <div class="conclusion {{if contains .Conclusion "PASSED"}}passed{{else if contains .Conclusion "WARNING"}}warning{{else}}failed{{end}}">
            {{.Conclusion}}
        </div>
    </div>

    <div class="footer">
        <p>Generated by SaaSus SDK Story Snapshot Testing Framework</p>
        <p>Report generated at {{formatTime .GeneratedAt}}</p>
    </div>
</body>
</html>`
}
