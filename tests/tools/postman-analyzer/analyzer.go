package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	openapiresolver "github.com/saasus-platform/saasus-sdk-go/tests/tools/openapi-resolver-lib"
)

// PostmanCollection represents the structure of a Postman collection
type PostmanCollection struct {
	Info PostmanInfo   `json:"info"`
	Item []PostmanItem `json:"item"`
}

type PostmanInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type PostmanItem struct {
	Name    string          `json:"name"`
	Item    []PostmanItem   `json:"item,omitempty"`
	Request *PostmanRequest `json:"request,omitempty"`
	Event   []PostmanEvent  `json:"event,omitempty"`
}

type PostmanRequest struct {
	Method string       `json:"method"`
	URL    PostmanURL   `json:"url"`
	Body   *PostmanBody `json:"body,omitempty"`
}

type PostmanURL struct {
	Raw  string   `json:"raw"`
	Path []string `json:"path"`
}

type PostmanBody struct {
	Mode string `json:"mode"`
	Raw  string `json:"raw"`
}

type PostmanEvent struct {
	Listen string        `json:"listen"`
	Script PostmanScript `json:"script"`
}

type PostmanScript struct {
	Exec []string `json:"exec"`
}

// PostmanAnalysis represents analyzed collection data
type PostmanAnalysis struct {
	CollectionName string            `json:"collection_name"`
	TestFlow       []PostmanTestStep `json:"test_flow"`
	Variables      map[string]string `json:"variables"`
	ClientMethods  []ClientMethodDef `json:"client_methods"`
}

type PostmanTestStep struct {
	Name              string                 `json:"name"`
	Method            string                 `json:"method"`
	Endpoint          string                 `json:"endpoint"`
	Module            string                 `json:"module"` // SaaSus module name (auth, pricing, billing, etc.)
	Description       string                 `json:"description"`
	Tests             []string               `json:"tests"`
	ExpectedStatus    int                    `json:"expected_status"`
	Variables         map[string]interface{} `json:"variables"`          // Variables set by this step
	RequiredVariables []string               `json:"required_variables"` // Variables required by this step (from URL params and body)
	RequestBody       map[string]interface{} `json:"request_body"`       // Request body structure
	BodyFields        []string               `json:"body_fields"`        // List of fields in request body

	// 新規追加フィールド（要件10.1, 10.2, 10.3, 10.4）
	OperationID         string               `json:"operation_id,omitempty"`          // OpenAPIのoperationId
	ClientMethodName    string               `json:"client_method_name,omitempty"`    // client.gen.goのメソッド名
	ParameterTypeName   string               `json:"parameter_type_name,omitempty"`   // パラメータ型名
	ParameterTypeFields map[string]FieldInfo `json:"parameter_type_fields,omitempty"` // 型のフィールド情報
}

// FieldInfo represents field information from types.gen.go
type FieldInfo struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	JSONTag  string `json:"json_tag"`
	Required bool   `json:"required"`
}

type ClientMethodDef struct {
	Name        string   `json:"name"`
	Type        string   `json:"type"`         // "standard", "with_response", "with_body", "with_body_with_response"
	HTTPMethod  string   `json:"http_method"`  // "GET", "POST", "PUT", "DELETE", etc.
	Endpoint    string   `json:"endpoint"`     // API endpoint path
	Module      string   `json:"module"`       // SaaSus module name (auth, pricing, billing, etc.)
	HasBody     bool     `json:"has_body"`     // Whether the method requires a request body
	Description string   `json:"description"`  // Method description
	OperationID string   `json:"operation_id"` // OpenAPI operationId
	ParamTypes  []string `json:"param_types"`  // Parameter types from types.gen.go
}

// resolverManager provides memoized access to OpenAPI operation resolvers per module.
type resolverManager struct {
	defaultModule   string
	defaultSpecPath string
	resolvers       map[string]*openapiresolver.OperationResolver
	mu              sync.Mutex
}

func newResolverManager(defaultModule, defaultSpecPath string) *resolverManager {
	return &resolverManager{
		defaultModule:   defaultModule,
		defaultSpecPath: defaultSpecPath,
		resolvers:       make(map[string]*openapiresolver.OperationResolver),
	}
}

func (rm *resolverManager) resolve(module, method, endpoint string) (*openapiresolver.OperationMapping, error) {
	resolver, err := rm.getResolver(module)
	if err != nil {
		return nil, err
	}
	if resolver == nil {
		return nil, fmt.Errorf("operation resolver not available for module '%s'", module)
	}

	mapping, err := resolver.Resolve(method, endpoint)
	if err != nil {
		return nil, err
	}

	return mapping, nil
}

func (rm *resolverManager) getResolver(module string) (*openapiresolver.OperationResolver, error) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	targetModule := module
	if targetModule == "" {
		targetModule = rm.defaultModule
	}

	if resolver, ok := rm.resolvers[targetModule]; ok {
		return resolver, nil
	}

	specPath := rm.resolveSpecPath(targetModule)
	if specPath == "" {
		return nil, fmt.Errorf("could not determine spec path for module '%s'", targetModule)
	}

	resolver, err := openapiresolver.NewOperationResolver(specPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize resolver for module '%s': %w", targetModule, err)
	}

	rm.resolvers[targetModule] = resolver
	return resolver, nil
}

func (rm *resolverManager) resolveSpecPath(module string) string {
	if module == rm.defaultModule && rm.defaultSpecPath != "" {
		return rm.defaultSpecPath
	}

	return deriveSpecPathFromModule(module)
}

func main() {
	var collectionFile, outputDir, module, specPath string
	flag.StringVar(&collectionFile, "collection", "", "Path to Postman collection JSON file")
	flag.StringVar(&outputDir, "output", "", "Output directory for analysis files")
	flag.StringVar(&module, "module", "", "SaaSus module name (e.g., billing)")
	flag.StringVar(&specPath, "spec", "", "Path to OpenAPI spec file (optional, auto-detected from module)")
	flag.Parse()

	var moduleName string

	if module != "" && collectionFile == "" {
		// moduleが指定されていて、collectionFileが指定されていない場合のみ解決
		resolvedCollection, resolvedModuleName, err := resolveModuleCollection(module)
		if err != nil {
			fmt.Printf("Error resolving module '%s': %v\n", module, err)
			os.Exit(1)
		}
		moduleName = resolvedModuleName
		collectionFile = resolvedCollection

		if outputDir == "" {
			fallback := moduleName
			if fallback == "" {
				fallback = module
			}
			outputDir = deriveOutputDirFromModule(fallback)
		}
	} else if module != "" {
		// moduleが指定されているが、collectionFileも指定されている場合
		moduleName = module
		if outputDir == "" {
			outputDir = deriveOutputDirFromModule(module)
		}
	}

	if moduleName == "" {
		moduleName = detectModuleFromCollectionPath(collectionFile)
	}

	// Auto-detect OpenAPI spec path from module name if not provided
	if specPath == "" && moduleName != "" {
		specPath = deriveSpecPathFromModule(moduleName)
	}

	if collectionFile == "" {
		fmt.Println("Usage: go run analyzer.go -collection <path-to-collection.json> [-output <output-dir>] [-module <module-name>] [-spec <openapi-spec-path>]")
		os.Exit(1)
	}

	if outputDir == "" {
		if moduleName != "" {
			outputDir = deriveOutputDirFromModule(moduleName)
		} else {
			outputDir = "tests/tools/postman-analyzer/reports"
		}
	}

	analysis, err := AnalyzePostmanCollection(collectionFile, specPath, moduleName)
	if err != nil {
		fmt.Printf("Error analyzing collection: %v\n", err)
		os.Exit(1)
	}

	// Generate and export report using ReportGenerator
	generator := NewReportGenerator(analysis)
	generator.GenerateConsoleReport()
	if err := generator.ExportFiles(outputDir); err != nil {
		fmt.Printf("Error exporting files: %v\n", err)
		os.Exit(1)
	}
}

// AnalyzePostmanCollection analyzes a Postman collection file
func AnalyzePostmanCollection(filePath string, specPath string, defaultModule string) (*PostmanAnalysis, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read collection file: %w", err)
	}

	var collection PostmanCollection
	if err := json.Unmarshal(data, &collection); err != nil {
		return nil, fmt.Errorf("failed to parse collection JSON: %w", err)
	}

	analysis := &PostmanAnalysis{
		CollectionName: collection.Info.Name,
		TestFlow:       []PostmanTestStep{},
		Variables:      make(map[string]string),
		ClientMethods:  []ClientMethodDef{},
	}

	moduleName := defaultModule
	if moduleName == "" {
		moduleName = detectModuleFromCollectionPath(filePath)
	}

	// OpenAPI resolverをモジュール単位で管理
	var resolverMgr *resolverManager
	if specPath != "" || moduleName != "" {
		resolverMgr = newResolverManager(moduleName, specPath)
		if specPath != "" {
			if _, err := resolverMgr.getResolver(moduleName); err != nil {
				fmt.Printf("Warning: Failed to initialize OpenAPI resolver for module '%s': %v\n", moduleName, err)
				fmt.Println("Falling back to best-effort analysis; operationId resolution may be partial")
			}
		}
	}

	// Extract test steps from collection items
	extractTestSteps(collection.Item, analysis, resolverMgr)

	// 各テストステップに対してclient.gen.goとtypes.gen.goから型情報を取得（要件10.2, 10.3, 10.4）
	enrichTestStepsWithTypeInfo(analysis)

	// Generate client methods based on endpoints
	generateClientMethods(analysis, resolverMgr)

	return analysis, nil
}

func extractTestSteps(items []PostmanItem, analysis *PostmanAnalysis, resolverMgr *resolverManager) {
	for _, item := range items {
		if item.Request != nil {
			// This is a request item
			step := PostmanTestStep{
				Name:                item.Name,
				Method:              item.Request.Method,
				Endpoint:            extractEndpoint(item.Request.URL),
				ExpectedStatus:      200, // Default
				Tests:               []string{},
				Variables:           make(map[string]interface{}),
				RequiredVariables:   []string{},
				RequestBody:         make(map[string]interface{}),
				BodyFields:          []string{},
				ParameterTypeFields: make(map[string]FieldInfo),
			}

			// Extract module information from URL
			step.Module = extractModuleFromURL(item.Request.URL)

			// OpenAPI resolverを使用してoperationIdを取得（要件10.1）
			if resolverMgr != nil {
				mapping, err := resolverMgr.resolve(step.Module, step.Method, step.Endpoint)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Warning: failed to resolve operationId (module=%s method=%s endpoint=%s): %v\n", step.Module, step.Method, step.Endpoint, err)
				} else if mapping != nil && mapping.OperationID != "" {
					step.OperationID = mapping.OperationID
				}
			}

			// Extract test scripts and variables
			for _, event := range item.Event {
				if event.Listen == "test" {
					step.Tests = append(step.Tests, event.Script.Exec...)
					// Extract variables set by this step
					extractVariablesFromScript(event.Script.Exec, &step)
				}
			}

			// Extract required variables from URL parameters
			extractRequiredVariablesFromURL(item.Request.URL, &step)

			// Extract request body information
			if item.Request.Body != nil && item.Request.Body.Mode == "raw" {
				extractRequestBodyInfo(item.Request.Body, &step)
			}

			analysis.TestFlow = append(analysis.TestFlow, step)
		} else if len(item.Item) > 0 {
			// This is a folder, recurse into it
			extractTestSteps(item.Item, analysis, resolverMgr)
		}
	}
}

// extractVariablesFromScript extracts variables set by pm.collectionVariables.set()
func extractVariablesFromScript(scriptLines []string, step *PostmanTestStep) {
	for _, line := range scriptLines {
		// Look for pm.collectionVariables.set("variable_name", ...)
		if strings.Contains(line, "pm.collectionVariables.set") {
			// Extract variable name
			start := strings.Index(line, "\"")
			if start != -1 {
				end := strings.Index(line[start+1:], "\"")
				if end != -1 {
					varName := line[start+1 : start+1+end]
					step.Variables[varName] = "extracted_from_response"
				}
			}
		}
	}
}

// extractRequiredVariablesFromURL extracts variables required from URL parameters
func extractRequiredVariablesFromURL(url PostmanURL, step *PostmanTestStep) {
	// Extract from path (e.g., :feedback_id)
	for _, pathSegment := range url.Path {
		if strings.HasPrefix(pathSegment, ":") {
			varName := strings.TrimPrefix(pathSegment, ":")
			step.RequiredVariables = append(step.RequiredVariables, varName)
		}
	}

	// Extract from raw URL (e.g., {{feedback_id}})
	if strings.Contains(url.Raw, "{{") {
		// Simple extraction of {{variable_name}} patterns
		raw := url.Raw
		for {
			start := strings.Index(raw, "{{")
			if start == -1 {
				break
			}
			end := strings.Index(raw[start:], "}}")
			if end == -1 {
				break
			}
			varName := raw[start+2 : start+end]
			// Skip URL variables like UrlCommunication
			if !strings.HasPrefix(varName, "Url") {
				if !contains(step.RequiredVariables, varName) {
					step.RequiredVariables = append(step.RequiredVariables, varName)
				}
			}
			raw = raw[start+end+2:]
		}
	}
}

// extractRequestBodyInfo extracts request body structure and fields
func extractRequestBodyInfo(body *PostmanBody, step *PostmanTestStep) {
	if body.Mode != "raw" {
		return
	}

	// Try to parse as JSON
	var bodyData map[string]interface{}
	if err := json.Unmarshal([]byte(body.Raw), &bodyData); err == nil {
		step.RequestBody = bodyData

		// Extract field names and required variables
		for key, value := range bodyData {
			step.BodyFields = append(step.BodyFields, key)

			// Check if value is a variable reference (e.g., "{{user_id}}")
			if strValue, ok := value.(string); ok {
				if strings.HasPrefix(strValue, "{{") && strings.HasSuffix(strValue, "}}") {
					varName := strings.TrimSuffix(strings.TrimPrefix(strValue, "{{"), "}}")
					if !contains(step.RequiredVariables, varName) {
						step.RequiredVariables = append(step.RequiredVariables, varName)
					}
				}
			}
		}
	}
}

// contains checks if a string slice contains a specific string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// extractModuleFromURL extracts the SaaSus module name from a Postman URL
func extractModuleFromURL(url PostmanURL) string {
	variableName := firstUrlVariable(url.Raw)
	if variableName == "" {
		return ""
	}

	if module, ok := urlVariableModuleMap[variableName]; ok {
		return module
	}

	normalized := normalizeModuleCandidate(variableName)
	if module, ok := urlVariableModuleMap[normalized]; ok {
		return module
	}

	return normalized
}

var urlVariableModuleMap = map[string]string{
	"UrlAuth":           "auth",
	"UrlPricing":        "pricing",
	"UrlBilling":        "billing",
	"UrlCommunication":  "communication",
	"UrlIntegration":    "integration",
	"UrlAwsMarketplace": "awsmarketplace",
	"UrlApiLog":         "apilog",
	"UrlApiLogs":        "apilog",
	"UrlApiGateway":     "api-gateway",
	"UrlApiGateways":    "api-gateway",
	"auth":              "auth",
	"pricing":           "pricing",
	"billing":           "billing",
	"communication":     "communication",
	"integration":       "integration",
	"awsmarketplace":    "awsmarketplace",
	"apilog":            "apilog",
	"apilogs":           "apilog",
	"api-gateway":       "api-gateway",
	"apigateway":        "api-gateway",
}

var urlVariablePattern = regexp.MustCompile(`\{\{([^}]+)\}\}`)

func firstUrlVariable(raw string) string {
	if raw == "" {
		return ""
	}

	matches := urlVariablePattern.FindAllStringSubmatch(raw, -1)
	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		candidate := match[1]
		if strings.HasPrefix(candidate, "Url") || strings.HasPrefix(strings.ToLower(candidate), "url") {
			return candidate
		}
	}

	return ""
}

func normalizeModuleCandidate(variable string) string {
	if variable == "" {
		return ""
	}

	candidate := strings.TrimSpace(variable)
	if strings.HasPrefix(candidate, "Url") {
		candidate = candidate[3:]
	} else if strings.HasPrefix(strings.ToLower(candidate), "url") {
		candidate = candidate[3:]
	}

	lower := strings.ToLower(candidate)
	switch lower {
	case "apilogs", "api-logs":
		return "apilog"
	case "apigateway", "api-gateway", "api_gateway":
		return "api-gateway"
	case "awsmarketplace", "aws-marketplace", "aws_marketplace":
		return "awsmarketplace"
	default:
		return lower
	}
}

func extractEndpoint(url PostmanURL) string {
	if len(url.Path) > 0 {
		return "/" + strings.Join(url.Path, "/")
	}
	// Fallback to parsing from raw URL
	if strings.Contains(url.Raw, "{{") {
		parts := strings.Split(url.Raw, "/")
		if len(parts) > 3 {
			return "/" + strings.Join(parts[3:], "/")
		}
	}
	return url.Raw
}

func generateClientMethods(analysis *PostmanAnalysis, resolverMgr *resolverManager) {
	methodMap := make(map[string]bool)

	// エンドポイントとモジュールのマッピング用マップを作成
	endpointToModule := make(map[string]string)

	// テストステップをループしてモジュール情報を収集
	for _, step := range analysis.TestFlow {
		key := step.Method + " " + step.Endpoint
		// 最初に見つかったモジュール名を使用（空でない場合のみ）
		if step.Module != "" && endpointToModule[key] == "" {
			endpointToModule[key] = step.Module
		}
	}

	for _, step := range analysis.TestFlow {
		endpoint := step.Endpoint
		method := step.Method

		key := method + " " + endpoint
		moduleName := endpointToModule[key]
		if moduleName == "" {
			moduleName = step.Module
		}

		var methodName string
		var operationID string
		var variants []string
		var hasBody bool

		// OpenAPI resolverを使用してoperationIdを取得
		if resolverMgr != nil && moduleName != "" {
			mapping, err := resolverMgr.resolve(moduleName, method, endpoint)
			if err == nil && mapping != nil && mapping.OperationID != "" {
				operationID = mapping.OperationID
				methodName = operationID
				variants = mapping.Variants
				hasBody = mapping.HasBody
			} else {
				methodName = generateMethodName(method, endpoint)
				hasBody = method == "POST" || method == "PUT" || method == "PATCH"
				variants = generateFallbackVariants(methodName, hasBody)
			}
		} else {
			methodName = generateMethodName(method, endpoint)
			hasBody = method == "POST" || method == "PUT" || method == "PATCH"
			variants = generateFallbackVariants(methodName, hasBody)
		}

		if !methodMap[methodName] {
			// Extract parameter types from types.gen.go
			paramTypes := extractParamTypesFromTypesGen(moduleName, operationID)

			// Extract body type if applicable
			var bodyType string
			if hasBody {
				bodyType = extractBodyTypeFromTypesGen(moduleName, operationID)
			}

			// バリアントごとにClientMethodDefを生成
			for _, variant := range variants {
				variantType := determineVariantType(variant, methodName)

				// Determine which param types to include based on variant
				variantParamTypes := []string{}
				if len(paramTypes) > 0 {
					variantParamTypes = append(variantParamTypes, fmt.Sprintf("%sParams", operationID))
				}
				if bodyType != "" && (variantType == "with_body" || variantType == "with_body_with_response") {
					variantParamTypes = append(variantParamTypes, bodyType)
				}

				analysis.ClientMethods = append(analysis.ClientMethods, ClientMethodDef{
					Name:        variant,
					Type:        variantType,
					HTTPMethod:  method,
					Endpoint:    endpoint,
					Module:      moduleName,
					HasBody:     hasBody,
					Description: fmt.Sprintf("%s %s", method, endpoint),
					OperationID: operationID,
					ParamTypes:  variantParamTypes,
				})
			}

			methodMap[methodName] = true
		}
	}
}

// generateFallbackVariants generates method variants using fallback logic
func generateFallbackVariants(methodName string, hasBody bool) []string {
	variants := []string{
		methodName,
		methodName + "WithResponse",
	}

	if hasBody {
		variants = append(variants,
			methodName+"WithBody",
			methodName+"WithBodyWithResponse",
		)
	}

	return variants
}

// determineVariantType determines the variant type from the variant name
func determineVariantType(variant, baseName string) string {
	if variant == baseName {
		return "standard"
	} else if variant == baseName+"WithResponse" {
		return "with_response"
	} else if variant == baseName+"WithBody" {
		return "with_body"
	} else if variant == baseName+"WithBodyWithResponse" {
		return "with_body_with_response"
	}
	return "standard"
}

func generateMethodName(httpMethod, endpoint string) string {
	// Convert HTTP method to action
	var action string
	switch httpMethod {
	case "GET":
		action = "Get"
	case "POST":
		action = "Create"
	case "PUT":
		action = "Update"
	case "DELETE":
		action = "Delete"
	case "PATCH":
		action = "Update"
	default:
		action = strings.Title(strings.ToLower(httpMethod))
	}

	// Convert endpoint to resource name
	parts := strings.Split(strings.Trim(endpoint, "/"), "/")
	var resource strings.Builder

	for _, part := range parts {
		if part != "" && !strings.Contains(part, "{") && !strings.Contains(part, "}") {
			// Convert to PascalCase
			resource.WriteString(strings.Title(strings.ToLower(part)))
		}
	}

	return action + resource.String()
}

func getCollectionsDir() string {
	possiblePaths := []string{
		"collections",
		"tests/tools/postman-analyzer/collections",
		filepath.Join("tests", "tools", "postman-analyzer", "collections"),
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return "collections"
}

func resolveModuleCollection(module string) (string, string, error) {
	resolver := NewModuleResolver(getCollectionsDir())
	return resolver.ResolveCollection(module)
}

func detectModuleFromCollectionPath(path string) string {
	resolver := NewModuleResolver(getCollectionsDir())
	return resolver.DetectModuleFromPath(path)
}

func deriveOutputDirFromModule(name string) string {
	resolver := NewModuleResolver(getCollectionsDir())
	return resolver.DeriveOutputDir(name)
}

// extractParamTypesFromTypesGen extracts parameter types from types.gen.go for a given operation ID
func extractParamTypesFromTypesGen(moduleName, operationID string) []string {
	if moduleName == "" || operationID == "" {
		return nil
	}

	// Construct path to types.gen.go
	typesGenPath := filepath.Join("generated", moduleName+"api", "types.gen.go")

	// Read the file
	content, err := os.ReadFile(typesGenPath)
	if err != nil {
		return nil
	}

	// Convert operationID to expected param type name
	// e.g., "GetFeedbacks" -> "GetFeedbacksParams"
	paramTypeName := operationID + "Params"

	// Search for the param type definition
	// Pattern: type GetFeedbacksParams struct {
	pattern := regexp.MustCompile(`type\s+` + regexp.QuoteMeta(paramTypeName) + `\s+struct\s*\{`)
	if !pattern.Match(content) {
		return nil
	}

	// Extract the struct definition
	lines := strings.Split(string(content), "\n")
	var paramTypes []string
	inStruct := false

	for _, line := range lines {
		if pattern.MatchString(line) {
			inStruct = true
			continue
		}

		if inStruct {
			// End of struct
			if strings.Contains(line, "}") {
				break
			}

			// Extract field name and type
			// Pattern: FieldName Type `json:"field_name"`
			trimmed := strings.TrimSpace(line)
			if trimmed == "" || strings.HasPrefix(trimmed, "//") {
				continue
			}

			// Parse field definition
			parts := strings.Fields(trimmed)
			if len(parts) >= 2 {
				fieldName := parts[0]
				fieldType := parts[1]
				paramTypes = append(paramTypes, fmt.Sprintf("%s %s", fieldName, fieldType))
			}
		}
	}

	return paramTypes
}

// extractBodyTypeFromTypesGen extracts request body type from types.gen.go for a given operation ID
func extractBodyTypeFromTypesGen(moduleName, operationID string) string {
	if moduleName == "" || operationID == "" {
		return ""
	}

	// Construct path to types.gen.go
	typesGenPath := filepath.Join("generated", moduleName+"api", "types.gen.go")

	// Read the file
	content, err := os.ReadFile(typesGenPath)
	if err != nil {
		return ""
	}

	// Convert operationID to expected body type name
	// e.g., "CreateFeedback" -> "CreateFeedbackJSONRequestBody"
	bodyTypeName := operationID + "JSONRequestBody"

	// Search for the body type definition
	// Pattern: type CreateFeedbackJSONRequestBody = ...
	pattern := regexp.MustCompile(`type\s+` + regexp.QuoteMeta(bodyTypeName) + `\s*=\s*(\w+)`)
	matches := pattern.FindStringSubmatch(string(content))

	if len(matches) > 1 {
		return matches[1]
	}

	return ""
}

// enrichTestStepsWithTypeInfo enriches test steps with type information from client.gen.go and types.gen.go
func enrichTestStepsWithTypeInfo(analysis *PostmanAnalysis) {
	// モジュールごとにパーサーをキャッシュ
	clientParsers := make(map[string]*ClientInterfaceParserLocal)
	typeParsers := make(map[string]*TypeDefinitionParserLocal)

	for i := range analysis.TestFlow {
		step := &analysis.TestFlow[i]

		// operationIdが取得できていない場合はスキップ
		if step.OperationID == "" || step.Module == "" {
			continue
		}

		// ClientInterfaceParserを初期化（キャッシュから取得または新規作成）
		clientParser, ok := clientParsers[step.Module]
		if !ok {
			var err error
			clientParser, err = NewClientInterfaceParserLocal(step.Module)
			if err != nil {
				// パーサーの初期化に失敗した場合は警告を表示してスキップ
				fmt.Printf("Warning: Failed to initialize client parser for module '%s': %v\n", step.Module, err)
				continue
			}
			clientParsers[step.Module] = clientParser
		}

		// operationIdからメソッド名を取得（要件10.2）
		methodSig, err := clientParser.FindMethodByOperationID(step.OperationID)
		if err != nil {
			// メソッドが見つからない場合は警告を表示してスキップ
			fmt.Printf("Warning: Method not found for operationId '%s' in module '%s': %v\n", step.OperationID, step.Module, err)
			continue
		}

		step.ClientMethodName = methodSig.MethodName

		// メソッドシグネチャからパラメータ型名を抽出（要件10.3）
		for _, param := range methodSig.ParameterTypes {
			// JSONRequestBody型またはParam型を探す
			if strings.Contains(param.Type, "JSONRequestBody") || strings.HasSuffix(param.Type, "Param") {
				step.ParameterTypeName = param.Type

				// TypeDefinitionParserを初期化（キャッシュから取得または新規作成）
				typeParser, ok := typeParsers[step.Module]
				if !ok {
					var err error
					typeParser, err = NewTypeDefinitionParserLocal(step.Module)
					if err != nil {
						// パーサーの初期化に失敗した場合は警告を表示してスキップ
						fmt.Printf("Warning: Failed to initialize type parser for module '%s': %v\n", step.Module, err)
						break
					}
					typeParsers[step.Module] = typeParser
				}

				// パラメータ型名から型定義を取得（要件10.4）
				typeDef, err := typeParser.GetTypeDefinition(param.Type)
				if err != nil {
					// 型定義が見つからない場合は警告を表示してスキップ
					fmt.Printf("Warning: Type definition not found for type '%s' in module '%s': %v\n", param.Type, step.Module, err)
					break
				}

				// フィールド情報をPostmanTestStepに追加
				if step.ParameterTypeFields == nil {
					step.ParameterTypeFields = make(map[string]FieldInfo)
				}

				for _, field := range typeDef.Fields {
					step.ParameterTypeFields[field.JSONTag] = FieldInfo{
						Name:     field.Name,
						Type:     field.Type,
						JSONTag:  field.JSONTag,
						Required: field.Required,
					}
				}

				break // 最初のパラメータ型を見つけたら終了
			}
		}
	}
}

// ClientInterfaceParserLocal parses client.gen.go to extract method information
type ClientInterfaceParserLocal struct {
	moduleName string
	clientPath string
	methods    map[string]*MethodSignatureLocal
}

// MethodSignatureLocal represents a client method signature
type MethodSignatureLocal struct {
	MethodName     string
	OperationID    string
	ParameterTypes []ParameterTypeLocal
	ReturnType     string
}

// ParameterTypeLocal represents a method parameter
type ParameterTypeLocal struct {
	Name string
	Type string
}

// NewClientInterfaceParserLocal creates a new parser
func NewClientInterfaceParserLocal(moduleName string) (*ClientInterfaceParserLocal, error) {
	// generated/{module}api/client.gen.go のパスを構築
	clientPath := filepath.Join("generated", moduleName+"api", "client.gen.go")

	// ファイルが存在するか確認
	if _, err := os.Stat(clientPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("client.gen.go file not found for module: %s", moduleName)
	}

	parser := &ClientInterfaceParserLocal{
		moduleName: moduleName,
		clientPath: clientPath,
		methods:    make(map[string]*MethodSignatureLocal),
	}

	return parser, nil
}

// ParseAllMethods parses all methods in ClientInterface
func (p *ClientInterfaceParserLocal) ParseAllMethods() ([]MethodSignatureLocal, error) {
	// ファイルを読み込んでASTを解析
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, p.clientPath, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file %s: %w", p.clientPath, err)
	}

	// ClientInterfaceを検索
	var methods []MethodSignatureLocal
	ast.Inspect(node, func(n ast.Node) bool {
		// インターフェース型を検索
		typeSpec, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}

		// ClientInterfaceを特定
		if typeSpec.Name.Name != "ClientInterface" {
			return true
		}

		// インターフェース型を取得
		interfaceType, ok := typeSpec.Type.(*ast.InterfaceType)
		if !ok {
			return true
		}

		// 各メソッドを解析
		for _, method := range interfaceType.Methods.List {
			if len(method.Names) == 0 {
				continue
			}

			methodName := method.Names[0].Name
			funcType, ok := method.Type.(*ast.FuncType)
			if !ok {
				continue
			}

			// メソッドシグネチャを抽出
			sig := p.extractMethodSignature(methodName, funcType)
			if sig != nil {
				methods = append(methods, *sig)
				p.methods[sig.OperationID] = sig
			}
		}

		return false // ClientInterfaceを見つけたので探索終了
	})

	if len(methods) == 0 {
		return nil, fmt.Errorf("no methods found in ClientInterface in %s", p.clientPath)
	}

	return methods, nil
}

// extractMethodSignature extracts method signature from AST function type
func (p *ClientInterfaceParserLocal) extractMethodSignature(methodName string, funcType *ast.FuncType) *MethodSignatureLocal {
	sig := &MethodSignatureLocal{
		MethodName:     methodName,
		OperationID:    methodName, // デフォルトではメソッド名をoperationIDとして使用
		ParameterTypes: []ParameterTypeLocal{},
	}

	// パラメータを抽出
	if funcType.Params != nil {
		for _, field := range funcType.Params.List {
			paramType := p.exprToString(field.Type)

			// パラメータ名がある場合
			if len(field.Names) > 0 {
				for _, name := range field.Names {
					sig.ParameterTypes = append(sig.ParameterTypes, ParameterTypeLocal{
						Name: name.Name,
						Type: paramType,
					})
				}
			} else {
				// 無名パラメータの場合
				sig.ParameterTypes = append(sig.ParameterTypes, ParameterTypeLocal{
					Name: "",
					Type: paramType,
				})
			}
		}
	}

	// 戻り値の型を抽出
	if funcType.Results != nil && len(funcType.Results.List) > 0 {
		sig.ReturnType = p.exprToString(funcType.Results.List[0].Type)
	}

	return sig
}

// exprToString converts AST expression to string representation
func (p *ClientInterfaceParserLocal) exprToString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + p.exprToString(t.X)
	case *ast.SelectorExpr:
		return p.exprToString(t.X) + "." + t.Sel.Name
	case *ast.ArrayType:
		return "[]" + p.exprToString(t.Elt)
	case *ast.MapType:
		return "map[" + p.exprToString(t.Key) + "]" + p.exprToString(t.Value)
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.Ellipsis:
		return "..." + p.exprToString(t.Elt)
	default:
		return ""
	}
}

// FindMethodByOperationID finds method name by operationId
func (p *ClientInterfaceParserLocal) FindMethodByOperationID(operationID string) (*MethodSignatureLocal, error) {
	// メソッドがまだ解析されていない場合は解析を実行
	if len(p.methods) == 0 {
		if _, err := p.ParseAllMethods(); err != nil {
			return nil, err
		}
	}

	// operationIDで検索
	if method, ok := p.methods[operationID]; ok {
		return method, nil
	}

	// WithResponseサフィックスを除いて検索
	baseOperationID := strings.TrimSuffix(operationID, "WithResponse")
	if method, ok := p.methods[baseOperationID]; ok {
		return method, nil
	}

	// WithBodyサフィックスを除いて検索
	baseOperationID = strings.TrimSuffix(operationID, "WithBody")
	baseOperationID = strings.TrimSuffix(baseOperationID, "WithResponse")
	if method, ok := p.methods[baseOperationID]; ok {
		return method, nil
	}

	// 利用可能なメソッドのリストを作成
	availableMethods := make([]string, 0, len(p.methods))
	for methodName := range p.methods {
		availableMethods = append(availableMethods, methodName)
	}

	return nil, fmt.Errorf("method not found for operationId '%s'. Available methods: %v", operationID, availableMethods)
}

// TypeDefinitionParserLocal parses types.gen.go to extract type definitions
type TypeDefinitionParserLocal struct {
	moduleName string
	typesPath  string
	typeCache  map[string]*TypeDefinitionLocal
	aliasCache map[string]string
}

// TypeDefinitionLocal represents a complete type definition
type TypeDefinitionLocal struct {
	Name   string
	Fields []FieldDefinitionLocal
}

// FieldDefinitionLocal represents a struct field
type FieldDefinitionLocal struct {
	Name     string
	Type     string
	JSONTag  string
	Required bool
	Comment  string
}

// NewTypeDefinitionParserLocal creates a new parser
func NewTypeDefinitionParserLocal(moduleName string) (*TypeDefinitionParserLocal, error) {
	// generated/{module}api/types.gen.go のパスを構築
	typesPath := filepath.Join("generated", moduleName+"api", "types.gen.go")

	// ファイルが存在するか確認
	if _, err := os.Stat(typesPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("types.gen.go file not found for module: %s", moduleName)
	}

	parser := &TypeDefinitionParserLocal{
		moduleName: moduleName,
		typesPath:  typesPath,
		typeCache:  make(map[string]*TypeDefinitionLocal),
		aliasCache: make(map[string]string),
	}

	return parser, nil
}

// ParseAllTypes parses all type definitions
func (p *TypeDefinitionParserLocal) ParseAllTypes() (map[string]*TypeDefinitionLocal, error) {
	// ファイルを読み込んでASTを解析
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, p.typesPath, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file %s: %w", p.typesPath, err)
	}

	// すべての型定義を解析
	ast.Inspect(node, func(n ast.Node) bool {
		// 型定義を検索
		typeSpec, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}

		switch specType := typeSpec.Type.(type) {
		case *ast.StructType:
			// 構造体型をキャッシュ
			typeDef := p.extractTypeDefinition(typeSpec.Name.Name, specType)
			if typeDef != nil {
				p.typeCache[typeDef.Name] = typeDef
			}
		default:
			// 構造体以外は型エイリアスとして保持
			aliasTarget := p.exprToString(typeSpec.Type)
			if aliasTarget != "" {
				p.aliasCache[typeSpec.Name.Name] = aliasTarget
			}
		}

		return true
	})

	if len(p.typeCache) == 0 {
		return nil, fmt.Errorf("no type definitions found in %s", p.typesPath)
	}

	return p.typeCache, nil
}

// extractTypeDefinition extracts type definition from AST struct type
func (p *TypeDefinitionParserLocal) extractTypeDefinition(typeName string, structType *ast.StructType) *TypeDefinitionLocal {
	typeDef := &TypeDefinitionLocal{
		Name:   typeName,
		Fields: []FieldDefinitionLocal{},
	}

	// 各フィールドを解析
	for _, field := range structType.Fields.List {
		// フィールド名を取得
		var fieldName string
		if len(field.Names) > 0 {
			fieldName = field.Names[0].Name
		} else {
			// 埋め込みフィールドの場合はスキップ
			continue
		}

		// フィールド型を取得
		fieldType := p.exprToString(field.Type)

		// JSONタグを取得
		jsonTag := ""
		required := false
		if field.Tag != nil {
			tag := strings.Trim(field.Tag.Value, "`")
			jsonTag = p.extractJSONTag(tag)
			// requiredの判定（オプショナルでない場合はrequired）
			required = !strings.Contains(tag, "omitempty")
		}

		// コメントを取得
		comment := ""
		if field.Doc != nil && len(field.Doc.List) > 0 {
			comment = strings.TrimSpace(field.Doc.List[0].Text)
			comment = strings.TrimPrefix(comment, "//")
			comment = strings.TrimSpace(comment)
		}

		typeDef.Fields = append(typeDef.Fields, FieldDefinitionLocal{
			Name:     fieldName,
			Type:     fieldType,
			JSONTag:  jsonTag,
			Required: required,
			Comment:  comment,
		})
	}

	return typeDef
}

// extractJSONTag extracts JSON tag from struct tag
func (p *TypeDefinitionParserLocal) extractJSONTag(tag string) string {
	// json:"field_name" または json:"field_name,omitempty" の形式から抽出
	parts := strings.Split(tag, " ")
	for _, part := range parts {
		if strings.HasPrefix(part, "json:") {
			jsonPart := strings.TrimPrefix(part, "json:")
			jsonPart = strings.Trim(jsonPart, "\"")
			// カンマで分割して最初の部分を取得
			jsonFields := strings.Split(jsonPart, ",")
			if len(jsonFields) > 0 && jsonFields[0] != "-" {
				return jsonFields[0]
			}
		}
	}
	return ""
}

// exprToString converts AST expression to string representation
func (p *TypeDefinitionParserLocal) exprToString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + p.exprToString(t.X)
	case *ast.SelectorExpr:
		return p.exprToString(t.X) + "." + t.Sel.Name
	case *ast.ArrayType:
		return "[]" + p.exprToString(t.Elt)
	case *ast.MapType:
		return "map[" + p.exprToString(t.Key) + "]" + p.exprToString(t.Value)
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.Ellipsis:
		return "..." + p.exprToString(t.Elt)
	default:
		return ""
	}
}

// GetTypeDefinition gets type definition by type name
func (p *TypeDefinitionParserLocal) GetTypeDefinition(typeName string) (*TypeDefinitionLocal, error) {
	// キャッシュがまだ構築されていない場合は解析を実行
	if len(p.typeCache) == 0 {
		if _, err := p.ParseAllTypes(); err != nil {
			return nil, err
		}
	}

	resolvedName := p.normalizeTypeName(typeName)

	// 型名で検索
	if typeDef, ok := p.typeCache[resolvedName]; ok {
		return typeDef, nil
	}

	// 型エイリアスを解決（例: CreateFeedbackJSONRequestBody -> CreateFeedbackParam -> FeedbackParam）
	if aliasTarget, ok := p.resolveAliasChain(resolvedName); ok {
		if typeDef, exists := p.typeCache[aliasTarget]; exists {
			return typeDef, nil
		}
	}

	// 利用可能な型のリストを作成
	availableTypes := make([]string, 0, len(p.typeCache))
	for name := range p.typeCache {
		availableTypes = append(availableTypes, name)
	}

	return nil, fmt.Errorf("type definition not found for type '%s'. Available types: %v", typeName, availableTypes)
}

// resolveTypeAlias resolves type alias to actual type name
func (p *TypeDefinitionParserLocal) resolveTypeAlias(typeName string) (string, error) {
	// 既に解析済みのエイリアス情報を利用
	if len(p.typeCache) == 0 {
		if _, err := p.ParseAllTypes(); err != nil {
			return typeName, err
		}
	}

	if resolved, ok := p.resolveAliasChain(p.normalizeTypeName(typeName)); ok {
		return resolved, nil
	}

	return typeName, fmt.Errorf("type alias not found for type '%s'", typeName)
}

// resolveAliasChain follows alias definitions until reaching a concrete type name.
func (p *TypeDefinitionParserLocal) resolveAliasChain(typeName string) (string, bool) {
	if typeName == "" {
		return "", false
	}

	visited := make(map[string]struct{})
	current := typeName

	for {
		alias, ok := p.aliasCache[current]
		if !ok {
			return current, current != typeName
		}

		if _, seen := visited[current]; seen {
			return current, current != typeName
		}
		visited[current] = struct{}{}

		next := p.normalizeTypeName(alias)
		if next == "" || next == current {
			return current, current != typeName
		}
		current = next
	}
}

// normalizeTypeName removes pointer prefixes and package qualifiers.
func (p *TypeDefinitionParserLocal) normalizeTypeName(typeName string) string {
	if typeName == "" {
		return ""
	}

	trimmed := strings.TrimPrefix(typeName, "*")
	if idx := strings.LastIndex(trimmed, "."); idx != -1 {
		return trimmed[idx+1:]
	}
	return trimmed
}
