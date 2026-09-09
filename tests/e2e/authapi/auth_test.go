package authapi

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/saasus-platform/saasus-sdk-go/modules/auth"
	"github.com/saasus-platform/saasus-sdk-go/tests/testlib"
)

// TestAuthAPIE2E は Auth API の E2E テストを実行します
func TestAuthAPIE2E(t *testing.T) {
	fmt.Println("🚀 Starting SaaSus Auth API E2E Tests")
	fmt.Println("========================================")

	// Load environment variables from .env file
	if err := godotenv.Load("../../../.env"); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	// Create actual auth client
	rawClient, err := auth.AuthWithResponse()
	if err != nil {
		t.Fatalf("❌ Failed to create auth client: %v", err)
	}
	client := NewAuthClientWithStripeFallback(rawClient)

	// Initialize E2E engine with client and method list
	engine := testlib.NewE2EEngine(client, GetAuthMethods())

	// Parse command line arguments
	if err := engine.Config.ParseArgs(os.Args[1:]); err != nil {
		t.Fatalf("❌ Error parsing arguments: %v", err)
	}

	// Validate configuration
	if err := engine.Config.Validate(); err != nil {
		t.Fatalf("❌ Configuration error: %v", err)
	}

	// Get auth stories (passing testing.T for test parameter loading)
	stories := GetAuthStories(t)
	stories = attachStripeSetupToStories(t, stories)

	// Execute stories
	results := engine.ExecuteStories(stories)

	// Print results
	engine.PrintResults(results)

	// Verify method coverage (passing testing.T for test parameter loading)
	if err := VerifyMethodCoverage(t); err != nil {
		t.Errorf("❌ Method coverage verification failed: %v", err)
	}

	// Check if all tests passed
	allPassed := true
	for _, result := range results {
		if result.Status != testlib.TestStatusPassed {
			allPassed = false
			t.Errorf("❌ Story '%s' failed", result.StoryName)
		}
	}

	if allPassed {
		fmt.Println("\n✅ All Auth API E2E tests completed!")
	} else {
		t.Fatal("\n❌ Some tests failed")
	}
}
