package pricingapi

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/saasus-platform/saasus-sdk-go/modules/auth"
	"github.com/saasus-platform/saasus-sdk-go/modules/pricing"
	"github.com/saasus-platform/saasus-sdk-go/tests/testlib"
)

func TestPricingAPIE2E(t *testing.T) {
	fmt.Println("🚀 Starting SaaSus Pricing API E2E Tests")
	fmt.Println("========================================")

	// Load environment variables from .env file
	if err := godotenv.Load("../../../.env"); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	// Create pricing client
	pricingClient, err := pricing.PricingWithResponse()
	if err != nil {
		t.Fatalf("❌ Failed to create pricing client: %v", err)
	}

	// Create auth client for prerequisite steps
	authClient, err := auth.AuthWithResponse()
	if err != nil {
		t.Fatalf("❌ Failed to create auth client: %v", err)
	}

	// Combine methods from both APIs
	allMethods := append(GetPricingMethods(), GetAuthMethods()...)

	// Initialize E2E engine with pricing client
	// Note: Auth client will be used separately for prerequisite steps
	engine := testlib.NewE2EEngine(pricingClient, allMethods)

	// Store auth client in engine for prerequisite steps
	// We need to handle auth API calls separately
	_ = authClient // Will be used in prerequisite steps

	// Parse command line arguments
	if err := engine.Config.ParseArgs(os.Args[1:]); err != nil {
		t.Fatalf("❌ Error parsing arguments: %v", err)
	}

	// Validate configuration
	if err := engine.Config.Validate(); err != nil {
		t.Fatalf("❌ Configuration error: %v", err)
	}

	// Get stories
	stories := GetPricingStories()

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

	// Calculate coverage percentage
	coveredMethods, totalMethods, coveragePercent := engine.Coverage.GetCoverage()

	if allPassed {
		if engine.Coverage.IsFullyCovered() {
			fmt.Println("\n🎉 ALL TESTS PASSED - FULL METHOD COVERAGE ACHIEVED!")
		} else {
			fmt.Printf("\n✅ ALL TESTS PASSED - Coverage: %.1f%% (%d/%d methods)\n",
				coveragePercent, coveredMethods, totalMethods)
		}
	} else {
		t.Fatal("\n❌ Some tests failed")
	}
}
