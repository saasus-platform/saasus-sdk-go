package authapi

import (
	"fmt"
	"testing"

	"github.com/saasus-platform/saasus-sdk-go/tests/testlib"
)

// attachStripeSetupToStories wraps each story's Setup to ensure Stripe integration
// is configured (via UpdateStripeInfo) right before the story executes.
func attachStripeSetupToStories(t *testing.T, stories []testlib.Story) []testlib.Story {
	prepared := make([]testlib.Story, len(stories))
	for i, s := range stories {
		story := s
		storyName := story.Name
		originalSetup := story.Setup
		story.Setup = func() error {
			prepareStripeIntegrationForStory(t, storyName)
			if originalSetup != nil {
				return originalSetup()
			}
			return nil
		}
		prepared[i] = story
	}
	return prepared
}

// prepareStripeIntegrationForStory logs and ensures the Stripe integration exists for a story.
func prepareStripeIntegrationForStory(t *testing.T, storyName string) {
	if getStripeSecretKey() == "" {
		fmt.Printf("⚠️  Story '%s': STRIPE_SECRET_KEY not set; Stripe steps will be skipped\n", storyName)
		return
	}

	fmt.Printf("🔧 Story '%s': ensuring Stripe integration via UpdateStripeInfo\n", storyName)
	if err := ensureStripeIntegrationFresh(t); err != nil {
		fmt.Printf("⚠️  Story '%s': failed to configure Stripe integration: %v\n", storyName, err)
	}
}
