package testlib

import (
	"fmt"
	"strings"
	"time"
)

// Logger provides enhanced logging functionality
type Logger struct {
	logLevel LogLevel
}

// NewLogger creates a new logger
func NewLogger(logLevel LogLevel) *Logger {
	return &Logger{
		logLevel: logLevel,
	}
}

// LogStoryStart logs the start of a story
func (l *Logger) LogStoryStart(storyName string) {
	if l.logLevel >= LogLevelDebug {
		fmt.Printf("\n🎬 STORY START: %s\n", storyName)
		fmt.Println(strings.Repeat("═", 63))
	}
}

// LogStoryEnd logs the end of a story
func (l *Logger) LogStoryEnd(storyName string, duration time.Duration) {
	if l.logLevel >= LogLevelDebug {
		fmt.Printf("\n✅ STORY COMPLETED: %s (Duration: %v)\n", storyName, duration)
		fmt.Println(strings.Repeat("═", 63))
	}
}

// LogStepStart logs the start of a step
func (l *Logger) LogStepStart(stepNum int, stepName, method string) {
	if l.logLevel >= LogLevelDebug {
		fmt.Printf("\n📋 STEP %d: %s\n", stepNum, stepName)
		fmt.Printf("   Method: %s\n", method)
		fmt.Println("   " + strings.Repeat("─", 59))
	}
}

// LogStepResult logs the result of a step
func (l *Logger) LogStepResult(stepName, method string, status TestStatus, statusCode int, duration time.Duration) {
	if l.logLevel >= LogLevelDebug {
		statusIcon := "✅"
		if status == TestStatusFailed {
			statusIcon = "❌"
		} else if status == TestStatusSkipped {
			statusIcon = "⏭️"
		}
		fmt.Printf("   %s Step Result: %s | Status Code: %d | Duration: %v\n",
			statusIcon, status, statusCode, duration)
		fmt.Printf("   📊 Step Summary: %s -> %s (%d) in %v\n",
			stepName, method, statusCode, duration)
	}
}

// LogValidation logs validation results
func (l *Logger) LogValidation(stepName string, success bool, details string) {
	if l.logLevel >= LogLevelDebug {
		icon := "✅"
		result := "PASSED"
		if !success {
			icon = "❌"
			result = "FAILED"
		}
		fmt.Printf("\n🔍 VALIDATION %s: %s\n", result, stepName)
		fmt.Printf("   %s Validation Result: %s\n", icon, result)
		if details != "" {
			fmt.Printf("   Details: %s\n", details)
		}
	}
}

// LogStateUpdate logs state updates
func (l *Logger) LogStateUpdate(stepName string, variables map[string]interface{}) {
	if l.logLevel >= LogLevelDebug && len(variables) > 0 {
		fmt.Printf("\n🔄 STATE UPDATE: %s\n", stepName)
		fmt.Println("   Updated Variables:")
		for key, value := range variables {
			// Mask sensitive values
			displayValue := l.maskSensitiveValue(key, value)
			fmt.Printf("     %s: %v\n", key, displayValue)
		}
	}
}

// LogCoverage logs coverage updates
func (l *Logger) LogCoverage(methodName string, executions int, successRate float64) {
	if l.logLevel >= LogLevelDebug {
		fmt.Printf("📊 COVERAGE UPDATE: %s | Executions: %d | Success Rate: %.1f%%\n",
			methodName, executions, successRate)
	}
}

// LogDebug logs debug information
func (l *Logger) LogDebug(message string) {
	if l.logLevel >= LogLevelDebug {
		fmt.Printf("🔍 DEBUG: %s\n", message)
	}
}

// LogError logs error information
func (l *Logger) LogError(message string, err error) {
	fmt.Printf("❌ ERROR: %s", message)
	if err != nil {
		fmt.Printf(": %v", err)
	}
	fmt.Println()
}

// LogInfo logs informational messages
func (l *Logger) LogInfo(message string) {
	fmt.Printf("ℹ️  %s\n", message)
}

// maskSensitiveValue masks sensitive information in logs
func (l *Logger) maskSensitiveValue(key string, value interface{}) interface{} {
	keyLower := strings.ToLower(key)
	if strings.Contains(keyLower, "secret") ||
		strings.Contains(keyLower, "key") ||
		strings.Contains(keyLower, "token") {
		if str, ok := value.(string); ok && len(str) > 8 {
			return str[:4] + "..." + str[len(str)-4:]
		}
	}
	return value
}
