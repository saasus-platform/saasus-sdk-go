package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	openapiresolver "github.com/saasus-platform/saasus-sdk-go/tests/tools/openapi-resolver-lib"
)

// OpenAPIMapper はOpenAPI定義から型名を取得します
type OpenAPIMapper struct {
	moduleName string
	resolver   *openapiresolver.OperationResolver
	cache      map[string]string // endpoint:method -> typeName
}

// NewOpenAPIMapper は新しいOpenAPIMapperを作成します
func NewOpenAPIMapper(moduleName string) (*OpenAPIMapper, error) {
	// OpenAPI仕様ファイルのパスを構築
	specPath, err := findOpenAPISpec(moduleName)
	if err != nil {
		return nil, err
	}

	// リゾルバーを初期化
	resolver, err := openapiresolver.NewOperationResolver(specPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create OpenAPI resolver: %w", err)
	}

	return &OpenAPIMapper{
		moduleName: moduleName,
		resolver:   resolver,
		cache:      make(map[string]string),
	}, nil
}

// findOpenAPISpec はモジュール名からOpenAPI仕様ファイルを探します
func findOpenAPISpec(moduleName string) (string, error) {
	// 候補パスのリスト
	candidates := []string{
		fmt.Sprintf("api/%sapi.yml", moduleName),
		fmt.Sprintf("api/%sapi.yaml", moduleName),
		fmt.Sprintf("../../../api/%sapi.yml", moduleName),
		fmt.Sprintf("../../../api/%sapi.yaml", moduleName),
	}

	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			absPath, _ := filepath.Abs(path)
			debugLog("OpenAPI仕様ファイル発見: %s\n", absPath)
			return path, nil
		}
	}

	return "", fmt.Errorf("OpenAPI spec file not found for module: %s", moduleName)
}

// GetTypeNameForEndpoint はエンドポイントとメソッドから型名を取得します
func (m *OpenAPIMapper) GetTypeNameForEndpoint(method, endpoint string) (string, bool) {
	// キャッシュをチェック
	cacheKey := fmt.Sprintf("%s:%s", endpoint, method)
	if typeName, found := m.cache[cacheKey]; found {
		return typeName, true
	}

	// HTTPメソッドを正規化（大文字に）
	method = strings.ToUpper(method)

	// リゾルバーで解決
	mapping, err := m.resolver.Resolve(method, endpoint)
	if err != nil {
		debugLog("OpenAPIマッピング失敗: %s %s - %v\n", method, endpoint, err)
		return "", false
	}

	// requestBodyTypeが空の場合は失敗
	if mapping.RequestBodyType == "" {
		debugLog("OpenAPIマッピング: requestBodyTypeが空 - %s %s\n", method, endpoint)
		return "", false
	}

	// キャッシュに保存
	m.cache[cacheKey] = mapping.RequestBodyType

	debugLog("OpenAPIマッピング成功: %s %s -> %s\n", method, endpoint, mapping.RequestBodyType)
	return mapping.RequestBodyType, true
}

// GetAllMappings はすべてのエンドポイントマッピングを取得します
func (m *OpenAPIMapper) GetAllMappings() (map[string]string, error) {
	mappings, err := m.resolver.ResolveAll()
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, mapping := range mappings {
		if mapping.RequestBodyType != "" {
			key := fmt.Sprintf("%s:%s", mapping.Path, mapping.HTTPMethod)
			result[key] = mapping.RequestBodyType
		}
	}

	return result, nil
}
