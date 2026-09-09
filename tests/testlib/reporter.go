package testlib

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Reporter provides test result reporting functionality
type Reporter struct {
	coverage *CoverageTracker
}

// NewReporter creates a new reporter
func NewReporter(coverage *CoverageTracker) *Reporter {
	return &Reporter{
		coverage: coverage,
	}
}

// PrintSummary prints a summary of test results
func (r *Reporter) PrintSummary(results []StoryResult) {
	fmt.Println("\n📊 Test Execution Summary:")

	totalStories := len(results)
	passedStories := 0
	totalSteps := 0
	passedSteps := 0
	totalDuration := time.Duration(0)

	for _, result := range results {
		if result.Status == TestStatusPassed {
			passedStories++
		}
		totalSteps += len(result.Steps)
		totalDuration += result.Duration

		for _, step := range result.Steps {
			if step.Status == TestStatusPassed {
				passedSteps++
			}
		}
	}

	covered, total, coverage := r.coverage.GetCoverage()

	fmt.Printf("Stories: %d/%d passed (%.1f%%)\n",
		passedStories, totalStories,
		func() float64 {
			if totalStories > 0 {
				return float64(passedStories) / float64(totalStories) * 100
			}
			return 0
		}())

	fmt.Printf("Steps: %d/%d passed (%.1f%%)\n",
		passedSteps, totalSteps,
		func() float64 {
			if totalSteps > 0 {
				return float64(passedSteps) / float64(totalSteps) * 100
			}
			return 0
		}())

	fmt.Printf("Method Coverage: %d/%d methods (%.1f%%)\n", covered, total, coverage)
	fmt.Printf("Execution Time: %v\n", totalDuration)

	if passedStories == totalStories && r.coverage.IsFullyCovered() {
		fmt.Println("🎉 All tests passed successfully!")
	} else if passedStories < totalStories {
		fmt.Printf("❌ %d/%d stories failed\n", totalStories-passedStories, totalStories)
	}
}

// PrintDetailed prints detailed test results
func (r *Reporter) PrintDetailed(results []StoryResult) {
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("                    DETAILED TEST EXECUTION REPORT")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("Generated at: %s\n\n", time.Now().Format("2006-01-02 15:04:05"))

	for _, result := range results {
		statusIcon := "✅"
		if result.Status == TestStatusFailed {
			statusIcon = "❌"
		}

		fmt.Printf("%s %s (Duration: %v)\n", statusIcon, result.StoryName, result.Duration)

		if result.Error != nil {
			fmt.Printf("   Error: %v\n", result.Error)
		}

		fmt.Printf("   Steps: %d/%d passed\n",
			func() int {
				passed := 0
				for _, step := range result.Steps {
					if step.Status == TestStatusPassed {
						passed++
					}
				}
				return passed
			}(), len(result.Steps))

		for i, step := range result.Steps {
			stepIcon := "✅"
			if step.Status == TestStatusFailed {
				stepIcon = "❌"
			}
			fmt.Printf("   %d. %s %s -> %s (%d) in %v\n",
				i+1, stepIcon, step.StepName, step.Method, step.StatusCode, step.Duration)

			if step.Error != nil {
				fmt.Printf("      Error: %v\n", step.Error)
			}
		}
		fmt.Println()
	}
}

// ExportJSON exports results as JSON
func (r *Reporter) ExportJSON(results []StoryResult) ([]byte, error) {
	covered, total, coverage := r.coverage.GetCoverage()

	report := map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"summary": map[string]interface{}{
			"total_stories":   len(results),
			"passed_stories":  r.countPassedStories(results),
			"total_steps":     r.countTotalSteps(results),
			"passed_steps":    r.countPassedSteps(results),
			"coverage":        coverage,
			"covered_methods": covered,
			"total_methods":   total,
		},
		"stories": results,
	}

	return json.MarshalIndent(report, "", "  ")
}

// Helper methods
func (r *Reporter) countPassedStories(results []StoryResult) int {
	count := 0
	for _, result := range results {
		if result.Status == TestStatusPassed {
			count++
		}
	}
	return count
}

func (r *Reporter) countTotalSteps(results []StoryResult) int {
	count := 0
	for _, result := range results {
		count += len(result.Steps)
	}
	return count
}

func (r *Reporter) countPassedSteps(results []StoryResult) int {
	count := 0
	for _, result := range results {
		for _, step := range result.Steps {
			if step.Status == TestStatusPassed {
				count++
			}
		}
	}
	return count
}
