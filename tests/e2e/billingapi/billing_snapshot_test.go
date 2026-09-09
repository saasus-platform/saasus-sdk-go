package billingapi

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/saasus-platform/saasus-sdk-go/modules/billing"
	"github.com/saasus-platform/saasus-sdk-go/tests/testlib/snapshot"
)

// TestBillingAPIWithSnapshot tests billing API with snapshot functionality
func TestBillingAPIWithSnapshot(t *testing.T) {
	fmt.Println("🚀 Starting SaaSus Billing API Snapshot Tests")
	fmt.Println("==============================================")

	// Load environment variables from .env file
	if err := godotenv.Load("../../../.env"); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	// Create billing client
	client, err := billing.BillingWithResponse()
	if err != nil {
		t.Fatalf("❌ Failed to create billing client: %v", err)
	}

	// Get snapshot config
	config := GetBillingStorySnapshotConfig()

	// Parse command line flags for snapshot mode
	parseSnapshotFlags(config)

	// Create snapshot engine
	engine := snapshot.NewE2EEngineWithSnapshot(client, GetBillingMethods(), config)

	// Get billing stories
	stories := GetBillingStories()

	fmt.Printf("\n📋 Snapshot Configuration:\n")
	fmt.Printf("   Capture: %v\n", config.EnableCapture)
	fmt.Printf("   Compare: %v\n", config.EnableComparison)
	fmt.Printf("   Report: %v\n", config.EnableReporting)
	fmt.Printf("   Module: %s\n", config.ModuleSlug())
	fmt.Printf("   Output Directory: %s\n\n", config.GetModuleOutputDirectory())

	// Execute stories with snapshot
	results := engine.ExecuteStoriesWithSnapshot(stories)

	// Print results summary
	fmt.Println("\n📊 Snapshot Test Execution Summary:")
	fmt.Println("═══════════════════════════════════")

	allPassed := true
	for _, result := range results {
		status := result.StoryResult.Status
		fmt.Printf("Story: %s - Status: %s\n", result.StoryResult.StoryName, status)

		if status != "passed" {
			allPassed = false
			if result.StoryResult.Error != nil {
				fmt.Printf("  Error: %v\n", result.StoryResult.Error)
			}
		}

		// Print snapshot information
		if config.EnableCapture && result.Snapshot.StoryName != "" {
			fmt.Printf("  ✅ Snapshot captured: %d steps\n", len(result.Snapshot.Steps))
		}

		if config.EnableComparison && result.Comparison != nil {
			fmt.Printf("  📊 Comparison: %d differences found\n", len(result.Comparison.Differences))
		}

		if config.EnableReporting && result.Report != nil {
			fmt.Printf("  📄 Report generated\n")
		}
	}

	if !allPassed {
		t.Fatal("\n❌ Some snapshot tests failed")
	}

	fmt.Println("\n🎉 ALL SNAPSHOT TESTS PASSED!")
}

// parseSnapshotFlags parses command line flags for snapshot mode
func parseSnapshotFlags(config *snapshot.StorySnapshotConfig) {
	modeFound := false
	// Check for snapshot flags
	for i, arg := range os.Args {
		// Handle -snapshot-mode=value format
		if len(arg) > 15 && arg[:15] == "-snapshot-mode=" {
			mode := arg[15:]
			modeFound = true
			switch mode {
			case "capture":
				config.EnableCapture = true
				config.EnableComparison = false
				config.EnableReporting = false
				fmt.Println("📸 Mode: Capture only")
			case "compare":
				config.EnableCapture = false
				config.EnableComparison = true
				config.EnableReporting = false
				fmt.Println("🔍 Mode: Compare only")
			case "report":
				config.EnableCapture = false
				config.EnableComparison = false
				config.EnableReporting = true
				fmt.Println("📄 Mode: Report only")
			case "full":
				config.EnableCapture = true
				config.EnableComparison = true
				config.EnableReporting = true
				fmt.Println("🔄 Mode: Full (Capture + Compare + Report)")
			}
		}

		// Handle -snapshot-mode value format
		if arg == "-snapshot-mode" && i+1 < len(os.Args) {
			mode := os.Args[i+1]
			modeFound = true
			switch mode {
			case "capture":
				config.EnableCapture = true
				config.EnableComparison = false
				config.EnableReporting = false
				fmt.Println("📸 Mode: Capture only")
			case "compare":
				config.EnableCapture = false
				config.EnableComparison = true
				config.EnableReporting = false
				fmt.Println("🔍 Mode: Compare only")
			case "report":
				config.EnableCapture = false
				config.EnableComparison = false
				config.EnableReporting = true
				fmt.Println("📄 Mode: Report only")
			case "full":
				config.EnableCapture = true
				config.EnableComparison = true
				config.EnableReporting = true
				fmt.Println("🔄 Mode: Full (Capture + Compare + Report)")
			}
		}

		if arg == "-snapshot-verbose" {
			fmt.Println("📝 Verbose mode enabled")
		}

		if arg == "-snapshot-stories" && i+1 < len(os.Args) {
			stories := os.Args[i+1]
			fmt.Printf("📋 Filtering stories: %s\n", stories)
		}
	}

	// Default to capture mode if no mode specified
	if !modeFound {
		config.EnableCapture = true
		fmt.Println("📸 Mode: Capture (default)")
	}
}
