# Billing API Snapshot Flag Usage Guide

This guide explains how to use the flag-based snapshot functionality with the Billing API tests.

## Overview

The Billing API tests now support automatic flag detection and execution mode switching. You can run tests in different modes using command-line flags without modifying code.

## Basic Usage

### Normal Test Execution (No Snapshot)
```bash
# Run standard E2E tests without snapshot functionality
go test -v ./tests/e2e/billingapi -run TestBillingAPIE2E
```

### Automatic Flag-Based Execution
```bash
# The main test function that automatically detects flags
go test -v ./tests/e2e/billingapi -run TestBillingAPIWithSnapshot
```

## Snapshot Modes

### Capture Mode
Captures snapshots during test execution:
```bash
go test -v ./tests/e2e/billingapi -run TestBillingAPIWithSnapshot -snapshot-mode=capture
```

### Compare Mode
Compares current execution with previous snapshots:
```bash
go test -v ./tests/e2e/billingapi -run TestBillingAPIWithSnapshot -snapshot-mode=compare
```

### Report Mode
Generates reports from comparison results:
```bash
go test -v ./tests/e2e/billingapi -run TestBillingAPIWithSnapshot -snapshot-mode=report
```

### Full Mode
Executes all phases (capture, compare, and report):
```bash
go test -v ./tests/e2e/billingapi -run TestBillingAPIWithSnapshot -snapshot-mode=full
```

## Advanced Options

### Specific Stories
Execute only specific stories:
```bash
go test -v ./tests/e2e/billingapi -run TestBillingAPIWithSnapshot \
  -snapshot-mode=capture \
  -snapshot-stories="Postman Collection Story - Standard Methods"
```

### Multiple Stories
Execute multiple specific stories:
```bash
go test -v ./tests/e2e/billingapi -run TestBillingAPIWithSnapshot \
  -snapshot-mode=capture \
  -snapshot-stories="Postman Collection Story - Standard Methods,Postman Collection Story - WithResponse Methods"
```

### Verbose Output
Enable detailed logging:
```bash
go test -v ./tests/e2e/billingapi -run TestBillingAPIWithSnapshot \
  -snapshot-mode=capture \
  -snapshot-verbose
```

### Tag Comparison
Compare specific tags:
```bash
go test -v ./tests/e2e/billingapi -run TestBillingAPIWithSnapshot \
  -snapshot-mode=compare \
  -snapshot-old-tag=v1.2.3 \
  -snapshot-new-tag=v1.2.4
```

### Custom Output Directory
Specify custom output directory:
```bash
go test -v ./tests/e2e/billingapi -run TestBillingAPIWithSnapshot \
  -snapshot-mode=capture \
  -snapshot-output=custom_snapshots
```

### Custom Configuration File
Use custom configuration:
```bash
go test -v ./tests/e2e/billingapi -run TestBillingAPIWithSnapshot \
  -snapshot-mode=capture \
  -snapshot-config=configs/billing_snapshot.json
```

## Specific Test Functions

### Forced Mode Tests
These tests run in specific modes regardless of flags:

```bash
# Force capture mode
go test -v ./tests/e2e/billingapi -run TestBillingAPISnapshotModeCapture

# Force compare mode
go test -v ./tests/e2e/billingapi -run TestBillingAPISnapshotModeCompare

# Force full mode
go test -v ./tests/e2e/billingapi -run TestBillingAPISnapshotModeFull
```

### Conditional Tests
These tests only run when specific flags are provided:

```bash
# Only runs when -snapshot-stories is provided
go test -v ./tests/e2e/billingapi -run TestBillingAPIWithSpecificStories \
  -snapshot-stories="Postman Collection Story - Standard Methods"

# Only runs when -snapshot-verbose is provided
go test -v ./tests/e2e/billingapi -run TestBillingAPIVerboseMode \
  -snapshot-verbose
```

### Configuration Validation
Test configuration validation:
```bash
go test -v ./tests/e2e/billingapi -run TestBillingAPIConfigurationValidation
```

## Flag Integration Tests

Test the flag integration functionality itself:
```bash
# Test basic integration
go test -v ./tests/e2e/billingapi -run TestBillingFlagIntegration

# Test mode detection
go test -v ./tests/e2e/billingapi -run TestBillingFlagModeDetection

# Test with flags
go test -v ./tests/e2e/billingapi -run TestBillingFlagModeDetection \
  -snapshot-mode=capture -snapshot-verbose
```

## Available Flags

| Flag | Description | Example |
|------|-------------|---------|
| `-snapshot-mode` | Snapshot execution mode | `capture`, `compare`, `report`, `full` |
| `-snapshot-stories` | Comma-separated story names | `"story1,story2"` |
| `-snapshot-old-tag` | Old tag for comparison | `v1.2.3` |
| `-snapshot-new-tag` | New tag for comparison | `v1.2.4` |
| `-snapshot-config` | Configuration file path | `configs/snapshot.json` |
| `-snapshot-output` | Output directory | `custom_output` |
| `-snapshot-verbose` | Enable verbose logging | (no value needed) |

## Execution Flow

1. **Flag Detection**: Tests automatically detect snapshot flags
2. **Mode Selection**: Based on flags, appropriate execution mode is selected
3. **Configuration**: Test configuration is built from flags and defaults
4. **Validation**: Configuration is validated before execution
5. **Execution**: Tests run in the selected mode with specified options
6. **Results**: Detailed results and logs are provided

## Error Handling

The tests include comprehensive error handling:

- **Invalid Modes**: Clear error messages for unsupported modes
- **Missing Requirements**: Validation errors for incomplete configurations
- **File Operations**: Graceful handling of file system errors
- **API Failures**: Proper error reporting for API call failures

## Integration with CI/CD

The flag-based approach integrates well with CI/CD pipelines:

```yaml
# Example GitHub Actions workflow
- name: Run Snapshot Capture
  run: go test -v ./tests/e2e/billingapi -run TestBillingAPIWithSnapshot -snapshot-mode=capture

- name: Run Snapshot Comparison
  run: go test -v ./tests/e2e/billingapi -run TestBillingAPIWithSnapshot -snapshot-mode=compare

- name: Generate Reports
  run: go test -v ./tests/e2e/billingapi -run TestBillingAPIWithSnapshot -snapshot-mode=report
```

## Troubleshooting

### Common Issues

1. **No snapshots captured**: Ensure `-snapshot-mode=capture` is specified
2. **Comparison failures**: Check that previous snapshots exist
3. **Configuration errors**: Use `-snapshot-verbose` for detailed logging
4. **Story not found**: Verify story names match exactly (case-sensitive)

### Debug Commands

```bash
# Check configuration
go test -v ./tests/e2e/billingapi -run TestBillingAPIConfigurationValidation -snapshot-verbose

# Test mode detection
go test -v ./tests/e2e/billingapi -run TestBillingFlagModeDetection -snapshot-mode=capture -snapshot-verbose

# List available stories
go test -v ./tests/e2e/billingapi -run TestBillingFlagStoryFiltering -snapshot-verbose
```

## Best Practices

1. **Use specific modes**: Choose the appropriate mode for your use case
2. **Enable verbose logging**: Use `-snapshot-verbose` for debugging
3. **Filter stories**: Use `-snapshot-stories` to focus on specific tests
4. **Validate configuration**: Run configuration validation tests regularly
5. **Check results**: Review test output and generated files
6. **Handle errors**: Pay attention to error messages and suggestions

This flag-based approach provides a flexible and powerful way to control snapshot testing without modifying test code, making it ideal for both development and CI/CD environments.