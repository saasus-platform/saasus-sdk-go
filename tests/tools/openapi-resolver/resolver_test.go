//go:build test
// +build test

package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewOperationResolver(t *testing.T) {
	// テスト用の一時的なOpenAPI仕様ファイルを作成
	tempDir := t.TempDir()
	specPath := filepath.Join(tempDir, "test_spec.yml")

	specContent := `openapi: 3.0.0
info:
  title: Test API
  version: 1.0.0
paths:
  /test:
    get:
      operationId: GetTest
      summary: Get test
    post:
      operationId: CreateTest
      summary: Create test
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
`

	if err := os.WriteFile(specPath, []byte(specContent), 0644); err != nil {
		t.Fatalf("Failed to create test spec file: %v", err)
	}

	// リゾルバーを作成
	resolver, err := NewOperationResolver(specPath)
	if err != nil {
		t.Fatalf("NewOperationResolver failed: %v", err)
	}

	if resolver == nil {
		t.Fatal("Expected resolver to be non-nil")
	}

	if resolver.spec == nil {
		t.Fatal("Expected spec to be non-nil")
	}

	if resolver.pathNormalizer == nil {
		t.Fatal("Expected pathNormalizer to be non-nil")
	}

	if resolver.variantGenerator == nil {
		t.Fatal("Expected variantGenerator to be non-nil")
	}
}

func TestNewOperationResolver_FileNotFound(t *testing.T) {
	tests := []struct {
		name     string
		specPath string
		wantErr  bool
	}{
		{
			name:     "Non-existent file",
			specPath: "/nonexistent/path/spec.yml",
			wantErr:  true,
		},
		{
			name:     "Empty path",
			specPath: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewOperationResolver(tt.specPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewOperationResolver() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestResolve(t *testing.T) {
	// テスト用の一時的なOpenAPI仕様ファイルを作成
	tempDir := t.TempDir()
	specPath := filepath.Join(tempDir, "test_spec.yml")

	specContent := `openapi: 3.0.0
info:
  title: Test API
  version: 1.0.0
paths:
  /test:
    get:
      operationId: GetTest
      summary: Get test
    post:
      operationId: CreateTest
      summary: Create test
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
  /test/{id}:
    get:
      operationId: GetTestById
      summary: Get test by ID
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
`

	if err := os.WriteFile(specPath, []byte(specContent), 0644); err != nil {
		t.Fatalf("Failed to create test spec file: %v", err)
	}

	resolver, err := NewOperationResolver(specPath)
	if err != nil {
		t.Fatalf("NewOperationResolver failed: %v", err)
	}

	tests := []struct {
		name            string
		method          string
		path            string
		wantOperationID string
		wantHasBody     bool
		wantVariants    int
		wantErr         bool
	}{
		{
			name:            "GET /test",
			method:          "GET",
			path:            "/test",
			wantOperationID: "GetTest",
			wantHasBody:     false,
			wantVariants:    2, // GetTest, GetTestWithResponse
			wantErr:         false,
		},
		{
			name:            "POST /test",
			method:          "POST",
			path:            "/test",
			wantOperationID: "CreateTest",
			wantHasBody:     true,
			wantVariants:    4, // CreateTest, CreateTestWithResponse, CreateTestWithBody, CreateTestWithBodyWithResponse
			wantErr:         false,
		},
		{
			name:            "GET /test/{id}",
			method:          "GET",
			path:            "/test/{id}",
			wantOperationID: "GetTestById",
			wantHasBody:     false,
			wantVariants:    2,
			wantErr:         false,
		},
		{
			name:    "Non-existent endpoint",
			method:  "GET",
			path:    "/nonexistent",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mapping, err := resolver.Resolve(tt.method, tt.path)

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if mapping.OperationID != tt.wantOperationID {
				t.Errorf("OperationID = %v, want %v", mapping.OperationID, tt.wantOperationID)
			}

			if mapping.HasBody != tt.wantHasBody {
				t.Errorf("HasBody = %v, want %v", mapping.HasBody, tt.wantHasBody)
			}

			if len(mapping.Variants) != tt.wantVariants {
				t.Errorf("Variants count = %v, want %v", len(mapping.Variants), tt.wantVariants)
			}

			if mapping.HTTPMethod != tt.method {
				t.Errorf("HTTPMethod = %v, want %v", mapping.HTTPMethod, tt.method)
			}

			if mapping.Path != tt.path {
				t.Errorf("Path = %v, want %v", mapping.Path, tt.path)
			}
		})
	}
}

func TestResolve_PathNormalization(t *testing.T) {
	// テスト用の一時的なOpenAPI仕様ファイルを作成
	tempDir := t.TempDir()
	specPath := filepath.Join(tempDir, "test_spec.yml")

	specContent := `openapi: 3.0.0
info:
  title: Test API
  version: 1.0.0
paths:
  /test/{id}:
    get:
      operationId: GetTestById
      summary: Get test by ID
`

	if err := os.WriteFile(specPath, []byte(specContent), 0644); err != nil {
		t.Fatalf("Failed to create test spec file: %v", err)
	}

	resolver, err := NewOperationResolver(specPath)
	if err != nil {
		t.Fatalf("NewOperationResolver failed: %v", err)
	}

	// 異なるパスパラメータ形式でテスト
	tests := []struct {
		name string
		path string
	}{
		{"Standard format", "/test/{id}"},
		{"Colon format", "/test/:id"},
		{"Double brace format", "/test/{{id}}"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mapping, err := resolver.Resolve("GET", tt.path)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if mapping.OperationID != "GetTestById" {
				t.Errorf("OperationID = %v, want GetTestById", mapping.OperationID)
			}
		})
	}
}

func TestResolveAll(t *testing.T) {
	// テスト用の一時的なOpenAPI仕様ファイルを作成
	tempDir := t.TempDir()
	specPath := filepath.Join(tempDir, "test_spec.yml")

	specContent := `openapi: 3.0.0
info:
  title: Test API
  version: 1.0.0
paths:
  /test:
    get:
      operationId: GetTest
      summary: Get test
    post:
      operationId: CreateTest
      summary: Create test
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
  /test/{id}:
    get:
      operationId: GetTestById
      summary: Get test by ID
    delete:
      operationId: DeleteTest
      summary: Delete test
`

	if err := os.WriteFile(specPath, []byte(specContent), 0644); err != nil {
		t.Fatalf("Failed to create test spec file: %v", err)
	}

	resolver, err := NewOperationResolver(specPath)
	if err != nil {
		t.Fatalf("NewOperationResolver failed: %v", err)
	}

	mappings, err := resolver.ResolveAll()
	if err != nil {
		t.Fatalf("ResolveAll failed: %v", err)
	}

	// 4つのエンドポイントが期待される
	expectedCount := 4
	if len(mappings) != expectedCount {
		t.Errorf("Expected %d mappings, got %d", expectedCount, len(mappings))
	}

	// すべてのマッピングにoperationIDが設定されているか確認
	for _, mapping := range mappings {
		if mapping.OperationID == "" {
			t.Errorf("Empty operationID for %s %s", mapping.HTTPMethod, mapping.Path)
		}

		if len(mapping.Variants) == 0 {
			t.Errorf("No variants for %s %s", mapping.HTTPMethod, mapping.Path)
		}
	}

	// 特定のエンドポイントが含まれているか確認
	foundGetTest := false
	foundCreateTest := false
	for _, mapping := range mappings {
		if mapping.HTTPMethod == "GET" && mapping.Path == "/test" && mapping.OperationID == "GetTest" {
			foundGetTest = true
		}
		if mapping.HTTPMethod == "POST" && mapping.Path == "/test" && mapping.OperationID == "CreateTest" {
			foundCreateTest = true
			// POSTでbodyがある場合は4つのバリアントが期待される
			if len(mapping.Variants) != 4 {
				t.Errorf("Expected 4 variants for CreateTest, got %d", len(mapping.Variants))
			}
		}
	}

	if !foundGetTest {
		t.Error("GetTest endpoint not found in mappings")
	}

	if !foundCreateTest {
		t.Error("CreateTest endpoint not found in mappings")
	}
}
