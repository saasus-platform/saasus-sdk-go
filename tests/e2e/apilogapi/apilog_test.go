// Package apilogapi provides E2E tests for the SaaSus Apilog API module.
// These tests verify the functionality of API log retrieval operations using
// the testlib framework with story-based test scenarios.
package apilogapi

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/saasus-platform/saasus-sdk-go/modules/apilog"
	"github.com/saasus-platform/saasus-sdk-go/tests/testlib"
)

// TestApilogAPIE2E is the main E2E test function for the Apilog API.
// It executes all test stories defined in stories.go and verifies:
// - All API methods work correctly
// - Response validation passes
// - Full method coverage is achieved
//
// Environment variables required:
// - SAASUS_SAAS_ID: SaaS identifier
// - SAASUS_API_KEY: API authentication key
// - SAASUS_SECRET_KEY: Secret key for request signing
//
// Optional environment variables:
// - LOG_LEVEL: Logging level (debug, info, warn, error)
// - SAASUS_BASE_URL: Override default API base URL
func TestApilogAPIE2E(t *testing.T) {
	fmt.Println("🚀 Starting SaaSus Apilog API E2E Tests")
	fmt.Println("========================================")

	// Load environment variables from .env file
	// This is optional - variables can also be set in the system environment
	if err := godotenv.Load("../../../.env"); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	// Create actual apilog client with response wrapper
	// This client automatically handles SigV1 authentication
	client, err := apilog.ApiLogWithResponse()
	if err != nil {
		t.Fatalf("❌ Failed to create apilog client: %v", err)
	}

	// Initialize test engine with client and method list for coverage tracking
	engine := testlib.NewE2EEngine(client, GetApilogMethods())

	// Parse command line arguments for test configuration
	// Supports flags like -stories, -timeout, -snapshot, etc.
	if err := engine.Config.ParseArgs(os.Args[1:]); err != nil {
		t.Fatalf("❌ Error parsing arguments: %v", err)
	}

	// Validate configuration before running tests
	// Ensures all required settings are properly configured
	if err := engine.Config.Validate(); err != nil {
		t.Fatalf("❌ Configuration error: %v", err)
	}

	// Get all test stories (scenarios) to execute
	// Each story contains multiple steps that test different API operations
	stories := GetApilogStories()

	// Execute all stories and collect results
	// The engine handles step execution, validation, and state management
	results := engine.ExecuteStories(stories)

	// Print detailed results including coverage information
	engine.PrintResults(results)

	// Check if all tests passed and full coverage was achieved
	allPassed := true
	for _, result := range results {
		if result.Status != testlib.TestStatusPassed {
			allPassed = false
			break
		}
	}

	// Report final test status
	if allPassed && engine.Coverage.IsFullyCovered() {
		fmt.Println("\n🎉 ALL TESTS PASSED - FULL METHOD COVERAGE ACHIEVED!")
	} else {
		t.Fatal("\n❌ Some tests failed or coverage is incomplete")
	}
}
