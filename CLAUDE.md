# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

This is the SaaSus SDK for Go, providing client libraries for the SaaSus platform APIs. The SDK is primarily generated from OpenAPI specifications and includes authentication, middleware, and comprehensive E2E testing.

## Key Architecture

### Generated Code Structure
- **`generated/`** - Contains auto-generated client code from OpenAPI specs
  - `authapi/` - Authentication API client (most comprehensive with E2E tests)
  - `apilogapi/`, `billingapi/`, `communicationapi/`, `integrationapi/`, `pricingapi/`, `awsmarketplaceapi/` - Other API clients
- **`modules/`** - High-level SDK modules that wrap generated clients with authentication
- **`client/`** - Core client functionality including SaaSus signature authentication (SigV1)
- **`middleware/`** - Request middleware for authentication and referer handling

### Authentication System
The SDK uses SaaSus SigV1 authentication implemented in `client/client.go:22`. All API calls require:
- `SAASUS_SECRET_KEY` - Secret key for HMAC signature
- `SAASUS_API_KEY` - API key identifier  
- `SAASUS_SAAS_ID` - SaaS application identifier

Each module in `modules/` provides a `*WithResponse()` function that returns an authenticated client.

## Common Commands

### Code Generation
```bash
# Generate client from OpenAPI spec
./generate.sh <package_name> <openapi_file.yml>
# Example: ./generate.sh auth authapi.yml
```

### Testing
```bash
# Run unit tests
go test ./...

# Run E2E tests (requires environment variables)
cd generated/authapi/e2e
go test -v

# Run specific E2E story test
go test -v -run TestE2EStoryExecution

# Run with verbose logging
go test -v -verbose
```

### Working with Generated Code
```bash
# Split large generated client file for easier editing
cd generated/authapi
./split_client_gen.sh client.gen.go
```

## Development Patterns

### E2E Test Structure
The authapi has comprehensive E2E tests organized as "stories" in `generated/authapi/e2e/stories/`. Each story represents a complete user workflow:
- Stories run in dependency order (01→10)
- Global setup/cleanup manages test environment
- Individual story setup/cleanup per test
- Test data organized in `testdata/` with YAML configurations

### Adding New API Support
1. Add OpenAPI spec file (`.yml`) to root
2. Run `./generate.sh <api_name> <spec_file>`
3. Create module wrapper in `modules/<api_name>/client.go`
4. Add authentication integration following existing patterns

### Environment Variables Required for E2E Tests
```bash
export SAASUS_SAAS_ID="your-saas-id"
export SAASUS_API_KEY="your-api-key"
export SAASUS_SECRET_KEY="your-secret-key"
```

## File Generation Notes

- Generated files should not be manually edited
- Use `split_client_gen.sh` to break down large generated files for review
- The `client_gen_split_output/` directory contains organized sections of generated code
- Test parameters are stored in `test_params.json` files alongside generated code