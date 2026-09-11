package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	// DefaultReportsDir is the default directory for analysis reports
	DefaultReportsDir = "tests/tools/postman-analyzer/reports"
)

// ModuleResolver handles module name resolution and collection file discovery.
// It provides methods to:
//   - Resolve module names to collection file paths
//   - Detect module names from file paths
//   - Derive output directories from module names
type ModuleResolver struct {
	collectionsDir string
}

// NewModuleResolver creates a new module resolver
func NewModuleResolver(collectionsDir string) *ModuleResolver {
	return &ModuleResolver{
		collectionsDir: collectionsDir,
	}
}

var moduleSlugAliases = map[string][]string{
	"awsmarketplace": {"aws-marketplace"},
	"apilog":         {"apilogs", "api-logs"},
	"integration":    {"eventbridge"},
	"analysis":       {"analysis"},
	"stripe":         {"stripe"},
}

// ResolveCollection resolves a module name to a collection file path.
// Returns the full path to the collection file and the extracted module name.
func (r *ModuleResolver) ResolveCollection(module string) (collectionPath string, moduleName string, err error) {
	moduleSlug := slugifyModuleName(module)
	candidateSlugs := []string{}
	if moduleSlug != "" {
		candidateSlugs = append(candidateSlugs, moduleSlug)
		if aliases, ok := moduleSlugAliases[moduleSlug]; ok {
			candidateSlugs = append(candidateSlugs, aliases...)
		}
	}

	entries, err := os.ReadDir(r.collectionsDir)
	if err != nil {
		return "", "", fmt.Errorf("failed to read collections directory: %w", err)
	}

	// Pre-allocate with estimated capacity
	available := make(map[string]struct{}, len(entries))

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := extractModuleNameFromFilename(entry.Name())
		if name == "" {
			continue
		}

		slug := slugifyModuleName(name)
		if slug == "" {
			continue
		}

		available[slug] = struct{}{}
		if slugMatchesCandidate(slug, candidateSlugs) {
			canonical := moduleSlug
			if canonical == "" {
				canonical = slug
			}
			return filepath.Join(r.collectionsDir, entry.Name()), canonical, nil
		}
	}

	return "", "", r.buildNotFoundError(moduleSlug, available)
}

func slugMatchesCandidate(slug string, candidates []string) bool {
	for _, candidate := range candidates {
		if candidate == slug {
			return true
		}
	}
	return false
}

// DetectModuleFromPath detects module name from collection file path
func (r *ModuleResolver) DetectModuleFromPath(path string) string {
	if path == "" {
		return ""
	}
	return extractModuleNameFromFilename(filepath.Base(path))
}

// DeriveOutputDir derives output directory from module name.
// Returns DefaultReportsDir if moduleName is empty or invalid.
func (r *ModuleResolver) DeriveOutputDir(moduleName string) string {
	slug := slugifyModuleName(moduleName)
	if slug == "" {
		return DefaultReportsDir
	}
	return filepath.Join(DefaultReportsDir, slug)
}

// buildNotFoundError creates a detailed error message with available modules
func (r *ModuleResolver) buildNotFoundError(moduleSlug string, available map[string]struct{}) error {
	if moduleSlug == "" {
		moduleSlug = "unknown"
	}

	var availableList []string
	for slug := range available {
		availableList = append(availableList, slug)
	}
	sort.Strings(availableList)

	if len(availableList) == 0 {
		return fmt.Errorf("module %q not found and no SaaSus collections are available", moduleSlug)
	}

	return fmt.Errorf("module %q not found; available modules: %s", moduleSlug, strings.Join(availableList, ", "))
}

// extractModuleNameFromFilename extracts module name from SaaSus collection filename.
// Expected format: "SaaSus {ModuleName} API Schema.postman_collection.json"
// Example: "SaaSus Billing API Schema.postman_collection.json" -> "Billing"
func extractModuleNameFromFilename(filename string) string {
	const prefix = "SaaSus "
	const suffix = " API Schema.postman_collection.json"

	name, ok := strings.CutPrefix(filename, prefix)
	if !ok {
		return ""
	}

	name, ok = strings.CutSuffix(name, suffix)
	if !ok {
		return ""
	}

	return name
}

// slugifyModuleName converts module name to URL-safe slug.
// Examples:
//
//	"Billing API" -> "billing-api"
//	"AWS Marketplace" -> "aws-marketplace"
//	"Auth  API" -> "auth-api" (multiple spaces collapsed)
func slugifyModuleName(name string) string {
	name = strings.ToLower(strings.TrimSpace(name))
	if name == "" {
		return ""
	}

	var builder strings.Builder
	builder.Grow(len(name)) // Pre-allocate memory
	lastWasHyphen := false

	for _, r := range name {
		switch {
		case isAlphanumeric(r):
			builder.WriteRune(r)
			lastWasHyphen = false
		case isSeparator(r):
			// Convert consecutive separators to single hyphen
			if !lastWasHyphen && builder.Len() > 0 {
				builder.WriteRune('-')
				lastWasHyphen = true
			}
		}
	}

	return strings.Trim(builder.String(), "-")
}

// isAlphanumeric checks if rune is alphanumeric (lowercase letter or digit)
func isAlphanumeric(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9')
}

// isSeparator checks if rune is a separator character
func isSeparator(r rune) bool {
	return r == ' ' || r == '-' || r == '_' || r == '/'
}

var moduleSpecOverrides = map[string]string{
	"api-gateway": "apigatewayapi.yml",
	"apigateway":  "apigatewayapi.yml",
	"apilog":      "apilogapi.yml",
	"apilogs":     "apilogapi.yml",
	"integration": "integration.yml",
	"stripe":      "stripeapi.yml",
}

// deriveSpecPathFromModule derives OpenAPI spec path from module name
func deriveSpecPathFromModule(moduleName string) string {
	if moduleName == "" {
		return ""
	}

	if specPath := resolveSpecWithOverride(moduleName); specPath != "" {
		return specPath
	}

	slug := slugifyModuleName(moduleName)
	if specPath := resolveSpecWithOverride(slug); specPath != "" {
		return specPath
	}

	if slug == "" {
		return ""
	}

	return trySpecSearch(slug + "api.yml")
}

func resolveSpecWithOverride(key string) string {
	fileName, ok := moduleSpecOverrides[key]
	if !ok {
		return ""
	}

	return trySpecSearch(fileName)
}

func trySpecSearch(fileName string) string {
	if fileName == "" {
		return ""
	}

	possiblePaths := []string{
		filepath.Join("api", fileName),
		filepath.Join("..", "..", "..", "api", fileName),
		filepath.Join("..", "..", "api", fileName),
	}

	for _, specPath := range possiblePaths {
		if _, err := os.Stat(specPath); err == nil {
			return specPath
		}
	}

	return ""
}
