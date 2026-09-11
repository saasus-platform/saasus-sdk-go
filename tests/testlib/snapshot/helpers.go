package snapshot

import (
	"os"

	"github.com/joho/godotenv"
)

// getSDKVersionFromEnv gets SDK version from environment
func getSDKVersionFromEnv() string {
	if version := getRealStripeEnvWithDefault("SDK_VERSION", ""); version != "" {
		return version
	}
	return "unknown"
}

// getRealStripeEnvWithDefault returns environment variable value with default fallback
func getRealStripeEnvWithDefault(key, defaultValue string) string {
	// Try to load .env file (ignore errors as it may not exist in all environments)
	_ = godotenv.Load()

	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
