package testlib

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// LogLevel represents the logging level
type LogLevel int

const (
	LogLevelError LogLevel = iota
	LogLevelWarn
	LogLevelInfo
	LogLevelDebug
)

// String returns the string representation of LogLevel
func (l LogLevel) String() string {
	switch l {
	case LogLevelError:
		return "error"
	case LogLevelWarn:
		return "warn"
	case LogLevelInfo:
		return "info"
	case LogLevelDebug:
		return "debug"
	default:
		return "info"
	}
}

// parseLogLevel converts a string to LogLevel with default handling and warnings
func parseLogLevel(level string) LogLevel {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "debug":
		return LogLevelDebug
	case "info":
		return LogLevelInfo
	case "warn", "warning":
		return LogLevelWarn
	case "error":
		return LogLevelError
	default:
		// Invalid value - use default and show warning
		if level != "" {
			fmt.Fprintf(os.Stderr, "Warning: Invalid LOG_LEVEL '%s', using 'info'\n", level)
		}
		return LogLevelInfo
	}
}

// getLogLevelFromEnv gets log level from environment variables with proper priority
func getLogLevelFromEnv() LogLevel {
	// E2E-specific setting takes priority
	if level := os.Getenv("E2E_LOG_LEVEL"); level != "" {
		return parseLogLevel(level)
	}

	// General LOG_LEVEL setting
	if level := os.Getenv("LOG_LEVEL"); level != "" {
		return parseLogLevel(level)
	}

	// Default to info level
	return LogLevelInfo
}

// Config holds test execution configuration
type Config struct {
	SaaSID     string
	APIKey     string
	SecretKey  string
	BaseURL    string
	LogLevel   LogLevel
	StripeKey  string
	Timeout    int
	MaxRetries int
	DryRun     bool
	// Snapshot configuration (optional)
	Snapshot *SnapshotConfig `json:"snapshot,omitempty"`
}

// SnapshotConfig holds snapshot-specific configuration
type SnapshotConfig struct {
	EnableCapture    bool   `json:"enable_capture"`
	EnableComparison bool   `json:"enable_comparison"`
	EnableReporting  bool   `json:"enable_reporting"`
	OutputDirectory  string `json:"output_directory"`
	CaptureLevel     string `json:"capture_level"`
}

// NewConfig creates a new configuration with defaults
func NewConfig() *Config {
	return &Config{
		LogLevel:   getLogLevelFromEnv(),
		Timeout:    300,
		MaxRetries: 3,
		DryRun:     false,
	}
}

// LoadFromEnvironment loads configuration from environment variables
func (c *Config) LoadFromEnvironment() {
	// Load .env file if it exists
	c.loadEnvFile()

	c.SaaSID = os.Getenv("SAASUS_SAAS_ID")
	c.APIKey = os.Getenv("SAASUS_API_KEY")
	c.SecretKey = os.Getenv("SAASUS_SECRET_KEY")
	c.BaseURL = os.Getenv("SAASUS_BASE_URL")
	c.StripeKey = os.Getenv("STRIPE_SECRET_KEY")

	if timeout := os.Getenv("E2E_TIMEOUT"); timeout != "" {
		if t, err := strconv.Atoi(timeout); err == nil && t > 0 {
			c.Timeout = t
		}
	}

	if retries := os.Getenv("E2E_MAX_RETRIES"); retries != "" {
		if r, err := strconv.Atoi(retries); err == nil && r >= 0 {
			c.MaxRetries = r
		}
	}

	c.DryRun = getEnvWithDefault("E2E_DRY_RUN", "false") == "true"

	// Update log level from environment using priority logic
	c.LogLevel = getLogLevelFromEnv()

	// Load snapshot configuration from environment
	c.loadSnapshotConfigFromEnv()
}

// loadSnapshotConfigFromEnv loads snapshot configuration from environment variables
func (c *Config) loadSnapshotConfigFromEnv() {
	// Initialize snapshot config if any snapshot-related env vars are set
	if os.Getenv("E2E_SNAPSHOT_ENABLE") != "" ||
		os.Getenv("E2E_SNAPSHOT_OUTPUT_DIR") != "" ||
		os.Getenv("E2E_SNAPSHOT_CAPTURE_LEVEL") != "" {

		if c.Snapshot == nil {
			c.Snapshot = &SnapshotConfig{}
		}

		c.Snapshot.EnableCapture = getEnvWithDefault("E2E_SNAPSHOT_ENABLE", "false") == "true"
		c.Snapshot.EnableComparison = getEnvWithDefault("E2E_SNAPSHOT_COMPARISON", "false") == "true"
		c.Snapshot.EnableReporting = getEnvWithDefault("E2E_SNAPSHOT_REPORTING", "false") == "true"
		c.Snapshot.OutputDirectory = getEnvWithDefault("E2E_SNAPSHOT_OUTPUT_DIR", "tests/e2e/snapshot")
		c.Snapshot.CaptureLevel = getEnvWithDefault("E2E_SNAPSHOT_CAPTURE_LEVEL", "FULL")
	}
}

// ParseArgs parses command line arguments
func (c *Config) ParseArgs(args []string) error {
	for i, arg := range args {
		switch arg {
		case "-v", "--verbose":
			c.LogLevel = LogLevelDebug
		case "--timeout":
			if i+1 < len(args) {
				if timeout, err := strconv.Atoi(args[i+1]); err == nil && timeout > 0 {
					c.Timeout = timeout
				} else {
					return fmt.Errorf("--timeout requires a positive integer")
				}
			}
		case "--dry-run":
			c.DryRun = true
		}
	}
	return nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.SaaSID == "" {
		return fmt.Errorf("missing required environment variable: SAASUS_SAAS_ID")
	}
	if c.APIKey == "" {
		return fmt.Errorf("missing required environment variable: SAASUS_API_KEY")
	}
	if c.SecretKey == "" {
		return fmt.Errorf("missing required environment variable: SAASUS_SECRET_KEY")
	}
	return nil
}

// loadEnvFile loads environment variables from .env file
func (c *Config) loadEnvFile() {
	paths := []string{".env", "../.env", "../../.env", "../../../.env", "../../../../.env"}
	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			godotenv.Load(path)
			break
		}
	}
}

// getEnvWithDefault returns environment variable or default value
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
