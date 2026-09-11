package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	openapiresolver "github.com/saasus-platform/saasus-sdk-go/tests/tools/openapi-resolver-lib"
)

// CLIConfig represents CLI configuration
type CLIConfig struct {
	SpecPath   string // OpenAPI仕様ファイルのパス
	Module     string // モジュール名（billing, auth等）
	Method     string // HTTPメソッド（GET, POST等）
	Path       string // エンドポイントパス
	ListAll    bool   // すべてのエンドポイントをリスト表示
	JSONOutput bool   // JSON形式で出力
}

// CLIHandler handles command-line interface
type CLIHandler struct {
	config   *CLIConfig
	resolver *openapiresolver.OperationResolver
}

// NewCLIHandler creates a new CLI handler
func NewCLIHandler() *CLIHandler {
	return &CLIHandler{
		config: &CLIConfig{},
	}
}

// ParseFlags parses command-line flags
func (h *CLIHandler) ParseFlags() error {
	// フラグの定義
	flag.StringVar(&h.config.SpecPath, "spec", "", "OpenAPI仕様ファイルのパス")
	flag.StringVar(&h.config.Module, "module", "", "モジュール名（billing, auth等）")
	flag.StringVar(&h.config.Method, "method", "", "HTTPメソッド（GET, POST, PUT, DELETE, PATCH）")
	flag.StringVar(&h.config.Path, "path", "", "エンドポイントパス（例: /stripe/info）")
	flag.BoolVar(&h.config.ListAll, "list", false, "すべてのエンドポイントをリスト表示")
	flag.BoolVar(&h.config.JSONOutput, "json", false, "JSON形式で出力")

	flag.Parse()

	// バリデーション
	if err := h.validateConfig(); err != nil {
		return err
	}

	return nil
}

// validateConfig validates CLI configuration
func (h *CLIHandler) validateConfig() error {
	// specPathとmoduleの両方が指定されていない場合はエラー
	if h.config.SpecPath == "" && h.config.Module == "" {
		return fmt.Errorf("either -spec or -module must be specified")
	}

	// moduleが指定されている場合、specPathを自動生成
	if h.config.Module != "" && h.config.SpecPath == "" {
		h.config.SpecPath = h.moduleToSpecPath(h.config.Module)
	}

	// listフラグが指定されていない場合、methodとpathが必要
	if !h.config.ListAll {
		if h.config.Method == "" {
			return fmt.Errorf("-method is required when -list is not specified")
		}

		if h.config.Path == "" {
			return fmt.Errorf("-path is required when -list is not specified")
		}

		// パスの基本的な検証
		if !strings.HasPrefix(h.config.Path, "/") {
			return fmt.Errorf("path must start with '/': %s", h.config.Path)
		}

		// HTTPメソッドの検証
		validMethods := map[string]bool{
			"GET": true, "POST": true, "PUT": true, "DELETE": true, "PATCH": true,
		}
		upperMethod := strings.ToUpper(h.config.Method)
		if !validMethods[upperMethod] {
			return fmt.Errorf("invalid HTTP method: %s (must be GET, POST, PUT, DELETE, or PATCH)", h.config.Method)
		}
		h.config.Method = upperMethod
	}

	return nil
}

// moduleToSpecPath converts module name to OpenAPI spec file path
func (h *CLIHandler) moduleToSpecPath(module string) string {
	// モジュール名から仕様ファイルパスを生成
	// 例: "billing" -> "../../../api/billingapi.yml"
	return filepath.Join("../../../api", fmt.Sprintf("%sapi.yml", module))
}

// InitResolver initializes the operation resolver
func (h *CLIHandler) InitResolver() error {
	resolver, err := openapiresolver.NewOperationResolver(h.config.SpecPath)
	if err != nil {
		return fmt.Errorf("failed to initialize resolver: %w", err)
	}
	h.resolver = resolver
	return nil
}

// PrintResult outputs a single operation mapping result
func (h *CLIHandler) PrintResult(mapping *openapiresolver.OperationMapping) error {
	if h.config.JSONOutput {
		return h.printResultJSON(mapping)
	}
	return h.printResultText(mapping)
}

// printResultJSON outputs result in JSON format
func (h *CLIHandler) printResultJSON(mapping *openapiresolver.OperationMapping) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(mapping); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}
	return nil
}

// printResultText outputs result in text format
func (h *CLIHandler) printResultText(mapping *openapiresolver.OperationMapping) error {
	fmt.Printf("HTTPメソッド: %s\n", mapping.HTTPMethod)
	fmt.Printf("パス: %s\n", mapping.Path)
	fmt.Printf("OperationID: %s\n", mapping.OperationID)
	fmt.Printf("リクエストボディ: %v\n", mapping.HasBody)

	if mapping.Summary != "" {
		fmt.Printf("概要: %s\n", mapping.Summary)
	}

	if mapping.Description != "" {
		fmt.Printf("説明: %s\n", mapping.Description)
	}

	if len(mapping.Variants) > 0 {
		fmt.Printf("\nメソッドバリアント:\n")
		for _, variant := range mapping.Variants {
			fmt.Printf("  - %s\n", variant)
		}
	}

	return nil
}

// PrintAllResults outputs all operation mappings
func (h *CLIHandler) PrintAllResults(mappings []openapiresolver.OperationMapping) error {
	if h.config.JSONOutput {
		return h.printAllResultsJSON(mappings)
	}
	return h.printAllResultsText(mappings)
}

// printAllResultsJSON outputs all results in JSON format
func (h *CLIHandler) printAllResultsJSON(mappings []openapiresolver.OperationMapping) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(mappings); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}
	return nil
}

// printAllResultsText outputs all results in text format
func (h *CLIHandler) printAllResultsText(mappings []openapiresolver.OperationMapping) error {
	fmt.Printf("合計 %d 個のエンドポイント:\n\n", len(mappings))

	for i, mapping := range mappings {
		fmt.Printf("[%d] %s %s\n", i+1, mapping.HTTPMethod, mapping.Path)
		fmt.Printf("    OperationID: %s\n", mapping.OperationID)

		if mapping.Summary != "" {
			fmt.Printf("    概要: %s\n", mapping.Summary)
		}

		if len(mapping.Variants) > 0 {
			fmt.Printf("    バリアント: %s\n", strings.Join(mapping.Variants, ", "))
		}

		fmt.Println()
	}

	return nil
}

// Execute runs the CLI command based on configuration
func (h *CLIHandler) Execute() error {
	// リゾルバーを初期化
	if err := h.InitResolver(); err != nil {
		return err
	}

	// listフラグが指定されている場合、すべてのエンドポイントを表示
	if h.config.ListAll {
		return h.executeListAll()
	}

	// 特定のエンドポイントを解決
	return h.executeResolve()
}

// executeListAll lists all endpoints
func (h *CLIHandler) executeListAll() error {
	mappings, err := h.resolver.ResolveAll()
	if err != nil {
		return fmt.Errorf("failed to resolve all operations: %w", err)
	}

	if err := h.PrintAllResults(mappings); err != nil {
		return fmt.Errorf("failed to print results: %w", err)
	}

	return nil
}

// executeResolve resolves a specific endpoint
func (h *CLIHandler) executeResolve() error {
	mapping, err := h.resolver.Resolve(h.config.Method, h.config.Path)
	if err != nil {
		return fmt.Errorf("failed to resolve operation: %w", err)
	}

	if err := h.PrintResult(mapping); err != nil {
		return fmt.Errorf("failed to print result: %w", err)
	}

	return nil
}

// PrintUsage prints usage information
func (h *CLIHandler) PrintUsage() {
	fmt.Fprintf(os.Stderr, "使用方法:\n")
	fmt.Fprintf(os.Stderr, "  openapi-resolver [オプション]\n\n")
	fmt.Fprintf(os.Stderr, "オプション:\n")
	fmt.Fprintf(os.Stderr, "  -spec <path>      OpenAPI仕様ファイルのパス\n")
	fmt.Fprintf(os.Stderr, "  -module <name>    モジュール名（billing, auth等）\n")
	fmt.Fprintf(os.Stderr, "  -method <method>  HTTPメソッド（GET, POST, PUT, DELETE, PATCH）\n")
	fmt.Fprintf(os.Stderr, "  -path <path>      エンドポイントパス（例: /stripe/info）\n")
	fmt.Fprintf(os.Stderr, "  -list             すべてのエンドポイントをリスト表示\n")
	fmt.Fprintf(os.Stderr, "  -json             JSON形式で出力\n\n")
	fmt.Fprintf(os.Stderr, "使用例:\n")
	fmt.Fprintf(os.Stderr, "  # 特定のエンドポイントのoperationIDを取得\n")
	fmt.Fprintf(os.Stderr, "  openapi-resolver -spec ../../api/billingapi.yml -method PUT -path /stripe/info\n\n")
	fmt.Fprintf(os.Stderr, "  # モジュール名から自動的に仕様ファイルを読み込む\n")
	fmt.Fprintf(os.Stderr, "  openapi-resolver -module billing -method PUT -path /stripe/info\n\n")
	fmt.Fprintf(os.Stderr, "  # すべてのエンドポイントをリスト表示\n")
	fmt.Fprintf(os.Stderr, "  openapi-resolver -spec ../../api/billingapi.yml -list\n\n")
	fmt.Fprintf(os.Stderr, "  # JSON形式で出力\n")
	fmt.Fprintf(os.Stderr, "  openapi-resolver -module billing -method PUT -path /stripe/info -json\n")
}
