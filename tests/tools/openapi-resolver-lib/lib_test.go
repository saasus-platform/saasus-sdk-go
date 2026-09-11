package openapiresolver

import (
	"testing"
)

func TestNormalizePath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "コロンパラメータ",
			input:    "/tenants/:tenant_id",
			expected: "/tenants/{tenant_id}",
		},
		{
			name:     "異なるパラメータ名",
			input:    "/tenants/:id",
			expected: "/tenants/{id}",
		},
		{
			name:     "複数のパラメータ",
			input:    "/tenants/:tenant_id/users/:user_id",
			expected: "/tenants/{tenant_id}/users/{user_id}",
		},
		{
			name:     "ダブルブレースパラメータ",
			input:    "/tenants/{{tenant_id}}",
			expected: "/tenants/{tenant_id}",
		},
		{
			name:     "シングルブレースパラメータ",
			input:    "/tenants/{tenant_id}",
			expected: "/tenants/{tenant_id}",
		},
		{
			name:     "パラメータなし",
			input:    "/tenants",
			expected: "/tenants",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizePath(tt.input)
			if result != tt.expected {
				t.Errorf("normalizePath(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestPathsMatch(t *testing.T) {
	tests := []struct {
		name     string
		path1    string
		path2    string
		expected bool
	}{
		{
			name:     "同じパラメータ名",
			path1:    "/tenants/{tenant_id}",
			path2:    "/tenants/{tenant_id}",
			expected: true,
		},
		{
			name:     "異なるパラメータ名でも一致",
			path1:    "/tenants/{tenant_id}",
			path2:    "/tenants/{id}",
			expected: true,
		},
		{
			name:     "複数パラメータで異なる名前",
			path1:    "/tenants/{tenant_id}/users/{user_id}",
			path2:    "/tenants/{id}/users/{uid}",
			expected: true,
		},
		{
			name:     "固定セグメントが異なる",
			path1:    "/tenants/{id}",
			path2:    "/users/{id}",
			expected: false,
		},
		{
			name:     "セグメント数が異なる",
			path1:    "/tenants/{id}",
			path2:    "/tenants/{id}/users",
			expected: false,
		},
		{
			name:     "パラメータと固定セグメントの位置が異なる",
			path1:    "/tenants/{id}/users",
			path2:    "/tenants/active/{id}",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := pathsMatch(tt.path1, tt.path2)
			if result != tt.expected {
				t.Errorf("pathsMatch(%q, %q) = %v, want %v", tt.path1, tt.path2, result, tt.expected)
			}
		})
	}
}

func TestPathMatcher_FindMatchingPath(t *testing.T) {
	// テスト用のOpenAPI仕様を作成
	spec := &OpenAPISpec{
		Paths: map[string]PathItem{
			"/tenants/{tenant_id}": {
				Get: &Operation{OperationID: "getTenant"},
			},
			"/tenants/{tenant_id}/users/{user_id}": {
				Get: &Operation{OperationID: "getTenantUser"},
			},
			"/feedbacks": {
				Get: &Operation{OperationID: "getFeedbacks"},
			},
		},
	}

	matcher := NewPathMatcher(spec)

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "同じパラメータ名",
			input:    "/tenants/{tenant_id}",
			expected: "/tenants/{tenant_id}",
		},
		{
			name:     "異なるパラメータ名（コロン形式）",
			input:    "/tenants/:id",
			expected: "/tenants/{tenant_id}",
		},
		{
			name:     "異なるパラメータ名（ブレース形式）",
			input:    "/tenants/{id}",
			expected: "/tenants/{tenant_id}",
		},
		{
			name:     "ダブルブレース形式",
			input:    "/tenants/{{id}}",
			expected: "/tenants/{tenant_id}",
		},
		{
			name:     "複数パラメータ",
			input:    "/tenants/:id/users/:uid",
			expected: "/tenants/{tenant_id}/users/{user_id}",
		},
		{
			name:     "パラメータなし",
			input:    "/feedbacks",
			expected: "/feedbacks",
		},
		{
			name:     "マッチしないパス",
			input:    "/unknown/path",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := matcher.FindMatchingPath(tt.input)
			if result != tt.expected {
				t.Errorf("FindMatchingPath(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
