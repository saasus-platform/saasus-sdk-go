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

// TestAuthAPISingleStory は特定のストーリーのみを実行するテストです
// 使用例: go test -v -run TestAuthAPISingleStory
func TestAuthAPISingleStory(t *testing.T) {
	fmt.Println("🚀 Starting SaaSus Auth API Single Story Test")
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

	// Get all auth stories (passing testing.T for test parameter loading)
	allStories := GetAuthStories(t)

	// Filter to only execute the first story (Standard Methods)
	stories := []testlib.Story{allStories[0]}
	stories = attachStripeSetupToStories(t, stories)

	fmt.Printf("📋 Executing single story: %s\n", stories[0].Name)

	// Execute stories
	results := engine.ExecuteStories(stories)

	// Print results
	engine.PrintResults(results)

	// Check if all tests passed
	allPassed := true
	for _, result := range results {
		if result.Status != testlib.TestStatusPassed {
			allPassed = false
			t.Errorf("❌ Story '%s' failed", result.StoryName)
		}
	}

	if allPassed {
		fmt.Println("\n✅ Single story test completed!")
	} else {
		t.Fatal("\n❌ Story test failed")
	}
}
