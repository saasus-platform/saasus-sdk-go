package billingapi

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/saasus-platform/saasus-sdk-go/modules/billing"
	"github.com/saasus-platform/saasus-sdk-go/tests/testlib"
)

func TestBillingAPIE2E(t *testing.T) {
	fmt.Println("🚀 Starting SaaSus Billing API E2E Tests")
	fmt.Println("========================================")

	// Load environment variables from .env file
	if err := godotenv.Load("../../../.env"); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	// Create actual billing client
	client, err := billing.BillingWithResponse()
	if err != nil {
		t.Fatalf("❌ Failed to create billing client: %v", err)
	}

	engine := testlib.NewE2EEngine(client, GetBillingMethods())

	// Parse command line arguments
	if err := engine.Config.ParseArgs(os.Args[1:]); err != nil {
		t.Fatalf("❌ Error parsing arguments: %v", err)
	}

	// Validate configuration
	if err := engine.Config.Validate(); err != nil {
		t.Fatalf("❌ Configuration error: %v", err)
	}

	// Get billing stories
	stories := GetBillingStories()

	// Execute stories
	results := engine.ExecuteStories(stories)

	// Print results
	engine.PrintResults(results)

	// Check if all tests passed
	allPassed := true
	for _, result := range results {
		if result.Status != testlib.TestStatusPassed {
			allPassed = false
			break
		}
	}

	if allPassed && engine.Coverage.IsFullyCovered() {
		fmt.Println("\n🎉 ALL TESTS PASSED - FULL METHOD COVERAGE ACHIEVED!")
	} else {
		t.Fatal("\n❌ Some tests failed or coverage is incomplete")
	}
}
