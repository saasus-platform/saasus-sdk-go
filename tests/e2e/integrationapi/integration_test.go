package integrationapi

import (
	"flag"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/saasus-platform/saasus-sdk-go/modules/integration"
	"github.com/saasus-platform/saasus-sdk-go/tests/testlib"
	"github.com/saasus-platform/saasus-sdk-go/tests/testlib/snapshot"
)

func TestIntegrationAPIE2E(t *testing.T) {
	// Load environment variables from .env file
	if err := godotenv.Load("../../../.env"); err != nil {
		t.Logf("Warning: .env file not found: %v", err)
	}

	// Verify required environment variables
	requiredEnvVars := []string{"SAASUS_SAAS_ID", "SAASUS_API_KEY", "SAASUS_SECRET_KEY", "TEST_AWS_ACCOUNT_ID"}
	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			t.Fatalf("Required environment variable %s is not set", envVar)
		}
	}

	// Create integration API client
	client, err := integration.IntegrationWithResponse()
	if err != nil {
		t.Fatalf("Failed to create integration client: %v", err)
	}

	// Parse command line flags (only if not already parsed)
	var storyFilter string
	var snapshotMode string
	if !flag.Parsed() {
		flag.StringVar(&storyFilter, "stories", "", "Comma-separated list of story names to run")
		flag.StringVar(&snapshotMode, "snapshot-mode", "", "Snapshot mode: capture, comparison, reporting, integrated")
		flag.Parse()
	} else {
		// Flags already parsed, get values from existing flags
		if f := flag.Lookup("stories"); f != nil {
			storyFilter = f.Value.String()
		}
		if f := flag.Lookup("snapshot-mode"); f != nil {
			snapshotMode = f.Value.String()
		}
	}

	// Get all stories
	stories := GetIntegrationStories()

	// Branch based on snapshot mode
	if snapshotMode != "" {
		// Snapshot test mode
		t.Logf("Running in snapshot mode: %s", snapshotMode)
		config := getSnapshotConfigForMode(snapshotMode)
		engine := snapshot.NewE2EEngineWithSnapshot(client, GetIntegrationMethods(), config)

		results := engine.ExecuteStoriesWithSnapshot(stories)

		// Convert snapshot results to testlib results for printing
		testlibResults := make([]testlib.StoryResult, len(results))
		for i, result := range results {
			testlibResults[i] = testlib.StoryResult{
				StoryName: result.StoryResult.StoryName,
				Status:    testlib.TestStatus(result.StoryResult.Status),
				Error:     result.StoryResult.Error,
				Duration:  result.StoryResult.Duration,
				Steps:     convertSnapshotStepsToTestlib(result.StoryResult.Steps),
				Variables: result.StoryResult.Variables,
			}
		}
		engine.PrintResults(testlibResults)

		// Verify all stories passed
		for _, result := range results {
			if result.StoryResult.Status != snapshot.TestStatusPassed {
				t.Errorf("Story '%s' failed: %v", result.StoryResult.StoryName, result.StoryResult.Error)
			}
		}
	} else {
		// Normal E2E test mode
		t.Logf("Running in normal E2E test mode")
		engine := testlib.NewE2EEngine(client, GetIntegrationMethods())

		// Execute stories
		t.Logf("Executing %d integration API stories...", len(stories))
		results := engine.ExecuteStories(stories)

		// Print results
		engine.PrintResults(results)

		// Verify all stories passed
		for _, result := range results {
			if result.Status != testlib.TestStatusPassed {
				t.Errorf("Story '%s' failed: %v", result.StoryName, result.Error)
			}
		}

		// Print summary
		passedCount := 0
		for _, result := range results {
			if result.Status == testlib.TestStatusPassed {
				passedCount++
			}
		}
		t.Logf("Test Summary: %d/%d stories passed", passedCount, len(results))
	}
}

// getSnapshotConfigForMode returns the appropriate snapshot config for the given mode
func getSnapshotConfigForMode(mode string) *snapshot.StorySnapshotConfig {
	switch mode {
	case "capture":
		return GetIntegrationSnapshotConfigForCapture()
	case "comparison":
		return GetIntegrationSnapshotConfigForComparison()
	case "reporting":
		return GetIntegrationSnapshotConfigForReporting()
	case "integrated":
		return GetIntegrationSnapshotConfigForIntegrated()
	default:
		return GetIntegrationSnapshotConfigForIntegrated()
	}
}

// convertSnapshotStepsToTestlib converts snapshot step results to testlib step results
func convertSnapshotStepsToTestlib(snapshotSteps []snapshot.StepResult) []testlib.StepResult {
	testlibSteps := make([]testlib.StepResult, len(snapshotSteps))
	for i, step := range snapshotSteps {
		testlibSteps[i] = testlib.StepResult{
			StepName:   step.StepName,
			Method:     step.Method,
			StatusCode: step.StatusCode,
			Status:     testlib.TestStatus(step.Status),
			Error:      step.Error,
			Duration:   step.Duration,
		}
	}
	return testlibSteps
}
