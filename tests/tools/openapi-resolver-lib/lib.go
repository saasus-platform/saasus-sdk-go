package openapiresolver

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// エラー定義
var (
	ErrFileNotFound       = errors.New("OpenAPI spec file not found")
	ErrInvalidYAML        = errors.New("invalid YAML format")
	ErrOperationNotFound  = errors.New("operation not found for given method and path")
	ErrMissingOperationID = errors.New("operationId is missing in OpenAPI spec")
	ErrInvalidPath        = errors.New("invalid path format")
)

// ResolverError はコンテキスト情報を含むエラーをラップする
type ResolverError struct {
	Op   string
	Path string
	Err  error
}

func (e *ResolverError) Error() string {
	return fmt.Sprintf("%s: %s: %v", e.Op, e.Path, e.Err)
}

func (e *ResolverError) Unwrap() error {
	return e.Err
}

// OpenAPISpec represents parsed OpenAPI specification
type OpenAPISpec struct {
	OpenAPI string              `yaml:"openapi"`
	Info    OpenAPIInfo         `yaml:"info"`
	Paths   map[string]PathItem `yaml:"paths"`
}

// OpenAPIInfo represents OpenAPI info section
type OpenAPIInfo struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	Version     string `yaml:"version"`
}

// PathItem represents operations available on a single path
type PathItem struct {
	Get    *Operation `yaml:"get,omitempty"`
	Post   *Operation `yaml:"post,omitempty"`
	Put    *Operation `yaml:"put,omitempty"`
	Delete *Operation `yaml:"delete,omitempty"`
	Patch  *Operation `yaml:"patch,omitempty"`
}

// Operation represents a single API operation
type Operation struct {
	OperationID string       `yaml:"operationId"`
	Summary     string       `yaml:"summary"`
	Description string       `yaml:"description"`
	Tags        []string     `yaml:"tags"`
	RequestBody *RequestBody `yaml:"requestBody,omitempty"`
	Parameters  []Parameter  `yaml:"parameters,omitempty"`
}

// RequestBody represents request body definition
type RequestBody struct {
	Required bool                       `yaml:"required"`
	Content  map[string]MediaTypeObject `yaml:"content"`
}

// MediaTypeObject represents media type object
type MediaTypeObject struct {
	Schema SchemaObject `yaml:"schema"`
}

// SchemaObject represents schema object
type SchemaObject struct {
	Ref  string `yaml:"$ref"`
	Type string `yaml:"type"`
}

// Parameter represents operation parameter
type Parameter struct {
	Name     string `yaml:"name"`
	In       string `yaml:"in"`
	Required bool   `yaml:"required"`
}

// OperationMapping represents mapping between endpoint and operationId
type OperationMapping struct {
	HTTPMethod      string   `json:"http_method"`
	Path            string   `json:"path"`
	OperationID     string   `json:"operation_id"`
	HasBody         bool     `json:"has_body"`
	RequestBodyType string   `json:"request_body_type"` // requestBodyの型名（$refから抽出）
	Summary         string   `json:"summary"`
	Description     string   `json:"description"`
	Variants        []string `json:"variants"`
}

// OperationResolver resolves operationId from HTTP method and path
type OperationResolver struct {
	spec             *OpenAPISpec
	pathMatcher      *PathMatcher
	variantGenerator *MethodVariantGenerator
}

// PathMatcher matches API paths with different parameter names
type PathMatcher struct {
	spec *OpenAPISpec
}

// MethodVariantGenerator generates method variants
type MethodVariantGenerator struct{}

// NewOperationResolver creates a new resolver by loading OpenAPI spec
func NewOperationResolver(specPath string) (*OperationResolver, error) {
	spec, err := ParseOpenAPISpec(specPath)
	if err != nil {
		return nil, err
	}

	return &OperationResolver{
		spec:             spec,
		pathMatcher:      NewPathMatcher(spec),
		variantGenerator: NewMethodVariantGenerator(),
	}, nil
}

// ParseOpenAPISpec parses OpenAPI YAML file
func ParseOpenAPISpec(filePath string) (*OpenAPISpec, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, &ResolverError{
			Op:   "ParseOpenAPISpec",
			Path: filePath,
			Err:  ErrFileNotFound,
		}
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, &ResolverError{
			Op:   "ParseOpenAPISpec",
			Path: filePath,
			Err:  err,
		}
	}

	var spec OpenAPISpec
	if err := yaml.Unmarshal(data, &spec); err != nil {
		return nil, &ResolverError{
			Op:   "ParseOpenAPISpec",
			Path: filePath,
			Err:  ErrInvalidYAML,
		}
	}

	return &spec, nil
}

// Resolve resolves operationId from HTTP method and path
func (r *OperationResolver) Resolve(method, path string) (*OperationMapping, error) {
	// OpenAPI仕様のパスとマッチング（パラメータ名が異なっても構造が同じならマッチ）
	matchedPath := r.pathMatcher.FindMatchingPath(path)
	if matchedPath == "" {
		return nil, &ResolverError{
			Op:   "Resolve",
			Path: path,
			Err:  ErrOperationNotFound,
		}
	}

	operation, err := r.spec.GetOperation(method, matchedPath)
	if err != nil {
		return nil, &ResolverError{
			Op:   "Resolve",
			Path: path,
			Err:  ErrOperationNotFound,
		}
	}

	return r.createOperationMapping(method, matchedPath, operation), nil
}

// GetOperation retrieves operation by HTTP method and path
func (spec *OpenAPISpec) GetOperation(method, path string) (*Operation, error) {
	pathItem, exists := spec.Paths[path]
	if !exists {
		return nil, &ResolverError{
			Op:   "GetOperation",
			Path: path,
			Err:  ErrOperationNotFound,
		}
	}

	var operation *Operation
	switch method {
	case "GET":
		operation = pathItem.Get
	case "POST":
		operation = pathItem.Post
	case "PUT":
		operation = pathItem.Put
	case "DELETE":
		operation = pathItem.Delete
	case "PATCH":
		operation = pathItem.Patch
	}

	if operation == nil {
		return nil, &ResolverError{
			Op:   "GetOperation",
			Path: path,
			Err:  ErrOperationNotFound,
		}
	}

	return operation, nil
}

// createOperationMapping creates an OperationMapping from an Operation
func (r *OperationResolver) createOperationMapping(method, path string, operation *Operation) *OperationMapping {
	hasBody := operation.RequestBody != nil
	variants := r.variantGenerator.GenerateVariants(operation.OperationID, method, hasBody)

	// requestBodyの型名を抽出
	requestBodyType := ""
	if hasBody {
		requestBodyType = extractRequestBodyType(operation.RequestBody)
	}

	return &OperationMapping{
		HTTPMethod:      method,
		Path:            path,
		OperationID:     operation.OperationID,
		HasBody:         hasBody,
		RequestBodyType: requestBodyType,
		Summary:         operation.Summary,
		Description:     operation.Description,
		Variants:        variants,
	}
}

// extractRequestBodyType extracts type name from requestBody schema reference
// 例: "#/components/schemas/CreateUserJSONRequestBody" -> "CreateUserJSONRequestBody"
func extractRequestBodyType(requestBody *RequestBody) string {
	if requestBody == nil {
		return ""
	}

	// application/jsonのスキーマを取得
	mediaType, ok := requestBody.Content["application/json"]
	if !ok {
		return ""
	}

	// $refから型名を抽出
	ref := mediaType.Schema.Ref
	if ref == "" {
		return ""
	}

	// "#/components/schemas/TypeName" から "TypeName" を抽出
	parts := strings.Split(ref, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}

	return ""
}

// NewPathMatcher creates a new PathMatcher
func NewPathMatcher(spec *OpenAPISpec) *PathMatcher {
	return &PathMatcher{spec: spec}
}

// FindMatchingPath finds an OpenAPI path that matches the given path
// パラメータ名が異なっても、パスの構造が同じであればマッチする
// マッチした場合はOpenAPI仕様のパス（正しいパラメータ名を含む）を返す
func (m *PathMatcher) FindMatchingPath(inputPath string) string {
	// 入力パスを正規化（:param や {{param}} を {param} に変換）
	normalizedInput := normalizePath(inputPath)

	// OpenAPI仕様の各パスと比較
	for specPath := range m.spec.Paths {
		normalizedSpec := normalizePath(specPath)
		if pathsMatch(normalizedInput, normalizedSpec) {
			// マッチした場合はOpenAPI仕様のパス（元のパラメータ名を含む）を返す
			return specPath
		}
	}

	return ""
}

// normalizePath converts various parameter formats to {param_name} format
func normalizePath(path string) string {
	result := normalizeColonParams(path)
	result = normalizeDoubleBraces(result)
	return result
}

// pathsMatch checks if two normalized paths have the same structure
// パラメータ名は無視して、パスの構造だけを比較
func pathsMatch(path1, path2 string) bool {
	segments1 := strings.Split(strings.Trim(path1, "/"), "/")
	segments2 := strings.Split(strings.Trim(path2, "/"), "/")

	if len(segments1) != len(segments2) {
		return false
	}

	for i := range segments1 {
		seg1 := segments1[i]
		seg2 := segments2[i]

		// 両方がパラメータの場合は一致とみなす（パラメータ名は無視）
		isParam1 := strings.HasPrefix(seg1, "{") && strings.HasSuffix(seg1, "}")
		isParam2 := strings.HasPrefix(seg2, "{") && strings.HasSuffix(seg2, "}")

		if isParam1 && isParam2 {
			continue
		}

		// 両方が固定セグメントの場合は完全一致が必要
		if seg1 != seg2 {
			return false
		}
	}

	return true
}

func normalizeColonParams(path string) string {
	var result []rune
	inParam := false

	for i, ch := range path {
		if ch == ':' && (i == 0 || path[i-1] == '/') {
			inParam = true
			result = append(result, '{')
		} else if inParam && (ch == '/' || i == len(path)-1) {
			if i == len(path)-1 && ch != '/' {
				result = append(result, ch)
			}
			result = append(result, '}')
			if ch == '/' {
				result = append(result, ch)
			}
			inParam = false
		} else {
			result = append(result, ch)
		}
	}

	if inParam {
		result = append(result, '}')
	}

	return string(result)
}

func normalizeDoubleBraces(path string) string {
	var result []rune
	i := 0
	pathRunes := []rune(path)

	for i < len(pathRunes) {
		if i < len(pathRunes)-1 && pathRunes[i] == '{' && pathRunes[i+1] == '{' {
			result = append(result, '{')
			i += 2
			for i < len(pathRunes) {
				if i < len(pathRunes)-1 && pathRunes[i] == '}' && pathRunes[i+1] == '}' {
					result = append(result, '}')
					i += 2
					break
				}
				result = append(result, pathRunes[i])
				i++
			}
		} else {
			result = append(result, pathRunes[i])
			i++
		}
	}

	return string(result)
}

// NewMethodVariantGenerator creates a new MethodVariantGenerator
func NewMethodVariantGenerator() *MethodVariantGenerator {
	return &MethodVariantGenerator{}
}

// GenerateVariants generates all method variants for an operation
func (g *MethodVariantGenerator) GenerateVariants(operationID string, httpMethod string, hasBody bool) []string {
	if operationID == "" {
		return []string{}
	}

	variants := []string{
		operationID,
		operationID + "WithResponse",
	}

	if hasBody && (httpMethod == "POST" || httpMethod == "PUT" || httpMethod == "PATCH") {
		variants = append(variants,
			operationID+"WithBody",
			operationID+"WithBodyWithResponse",
		)
	}

	return variants
}

// ResolveAll returns all operation mappings from the OpenAPI spec
func (r *OperationResolver) ResolveAll() ([]OperationMapping, error) {
	estimatedSize := len(r.spec.Paths) * 5
	mappings := make([]OperationMapping, 0, estimatedSize)

	for path, pathItem := range r.spec.Paths {
		methods := map[string]*Operation{
			"GET":    pathItem.Get,
			"POST":   pathItem.Post,
			"PUT":    pathItem.Put,
			"DELETE": pathItem.Delete,
			"PATCH":  pathItem.Patch,
		}

		for method, operation := range methods {
			if operation != nil {
				mapping := r.createOperationMapping(method, path, operation)
				mappings = append(mappings, *mapping)
			}
		}
	}

	return mappings, nil
}
