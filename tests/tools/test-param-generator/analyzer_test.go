package main

import (
	"os"
	"path/filepath"
	"testing"
)

// TestExtractResourceName はextractResourceName関数のテスト
func TestExtractResourceName(t *testing.T) {
	tests := []struct {
		name     string
		endpoint string
		want     string
	}{
		{
			name:     "simple endpoint",
			endpoint: "/feedbacks",
			want:     "feedbacks",
		},
		{
			name:     "endpoint with parameter",
			endpoint: "/feedbacks/:feedback_id",
			want:     "feedbacks",
		},
		{
			name:     "nested endpoint",
			endpoint: "/feedbacks/:feedback_id/comments",
			want:     "feedbacks-comments",
		},
		{
			name:     "nested endpoint with multiple parameters",
			endpoint: "/feedbacks/:feedback_id/comments/:comment_id",
			want:     "feedbacks-comments",
		},
		{
			name:     "endpoint with trailing slash",
			endpoint: "/feedbacks/",
			want:     "feedbacks",
		},
		{
			name:     "empty endpoint",
			endpoint: "",
			want:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractResourceName(tt.endpoint)
			if got != tt.want {
				t.Errorf("extractResourceName(%q) = %q, want %q", tt.endpoint, got, tt.want)
			}
		})
	}
}

// TestExtractParamGroups はextractParamGroups関数のテスト
func TestExtractParamGroups(t *testing.T) {
	testFlow := []TestStep{
		{
			Name:     "CreateFeedback",
			Method:   "POST",
			Endpoint: "/feedbacks",
			RequestBody: map[string]any{
				"user_id":              "{{user_id}}",
				"feedback_title":       "{{feedback_title}}",
				"feedback_description": "{{feedback_description}}",
			},
		},
		{
			Name:     "UpdateFeedback",
			Method:   "PATCH",
			Endpoint: "/feedbacks/:feedback_id",
			RequestBody: map[string]any{
				"feedback_title":       "{{feedback_title_update}}",
				"feedback_description": "{{feedback_description_update}}",
			},
		},
		{
			Name:              "GetFeedback",
			Method:            "GET",
			Endpoint:          "/feedbacks/:feedback_id",
			RequiredVariables: []string{"feedback_id"},
		},
		{
			Name:              "DeleteFeedback",
			Method:            "DELETE",
			Endpoint:          "/feedbacks/:feedback_id",
			RequiredVariables: []string{"feedback_id"},
		},
	}

	groups := extractParamGroups(testFlow)

	// feedbacksグループが存在することを確認
	if _, exists := groups["feedbacks"]; !exists {
		t.Fatal("feedbacksグループが見つかりません")
	}

	feedbackGroup := groups["feedbacks"]

	// CreateParamsが存在することを確認
	if len(feedbackGroup.CreateParams) == 0 {
		t.Error("CreateParamsが空です")
	}

	// UpdateParamsが存在することを確認
	if len(feedbackGroup.UpdateParams) == 0 {
		t.Error("UpdateParamsが空です")
	}

	t.Logf("feedbacksグループ: CreateParams=%d, UpdateParams=%d",
		len(feedbackGroup.CreateParams),
		len(feedbackGroup.UpdateParams))
}

// TestGenerateSampleValue はgenerateSampleValue関数のテスト
func TestGenerateSampleValue(t *testing.T) {
	tests := []struct {
		name      string
		fieldType string
		jsonTag   string
		wantType  string // 期待される値の型
	}{
		{
			name:      "UUID field",
			fieldType: "string",
			jsonTag:   "user_id",
			wantType:  "string",
		},
		{
			name:      "title field",
			fieldType: "string",
			jsonTag:   "feedback_title",
			wantType:  "string",
		},
		{
			name:      "description field",
			fieldType: "string",
			jsonTag:   "feedback_description",
			wantType:  "string",
		},
		{
			name:      "body field",
			fieldType: "string",
			jsonTag:   "body",
			wantType:  "string",
		},
		{
			name:      "status field",
			fieldType: "int",
			jsonTag:   "status",
			wantType:  "int",
		},
		{
			name:      "count field",
			fieldType: "int",
			jsonTag:   "count",
			wantType:  "int",
		},
		{
			name:      "boolean field",
			fieldType: "bool",
			jsonTag:   "enabled",
			wantType:  "bool",
		},
		{
			name:      "email field",
			fieldType: "string",
			jsonTag:   "email",
			wantType:  "string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := generateSampleValue(tt.fieldType, tt.jsonTag)
			if got == nil {
				t.Errorf("generateSampleValue(%q, %q) = nil, want non-nil", tt.fieldType, tt.jsonTag)
				return
			}

			// 型チェック
			switch tt.wantType {
			case "string":
				if _, ok := got.(string); !ok {
					t.Errorf("generateSampleValue(%q, %q) = %T, want string", tt.fieldType, tt.jsonTag, got)
				}
			case "int":
				if _, ok := got.(int); !ok {
					t.Errorf("generateSampleValue(%q, %q) = %T, want int", tt.fieldType, tt.jsonTag, got)
				}
			case "bool":
				if _, ok := got.(bool); !ok {
					t.Errorf("generateSampleValue(%q, %q) = %T, want bool", tt.fieldType, tt.jsonTag, got)
				}
			}

			t.Logf("generateSampleValue(%q, %q) = %v", tt.fieldType, tt.jsonTag, got)
		})
	}
}

func TestParseGeneratedTypesEnhancedResolvesAliasChain(t *testing.T) {
	tmpDir := t.TempDir()
	typesDir := filepath.Join(tmpDir, "generated", "integrationapi")
	if err := os.MkdirAll(typesDir, 0o755); err != nil {
		t.Fatalf("failed to create temp types directory: %v", err)
	}

	typesFile := filepath.Join(typesDir, "types.gen.go")
	source := `package integrationapi

type EventBridgeSettings struct {
	AWSAccountID string ` + "`json:\"aws_account_id\"`" + `
	AWSRegion    string ` + "`json:\"aws_region\"`" + `
}

type SaveEventBridgeSettingsParam = EventBridgeSettings
type SaveEventBridgeSettingsJSONRequestBody = SaveEventBridgeSettingsParam
`

	if err := os.WriteFile(typesFile, []byte(source), 0o644); err != nil {
		t.Fatalf("failed to write temp types file: %v", err)
	}

	GetGlobalCache().Clear()
	typeInfo, err := parseGeneratedTypesEnhanced(typesFile)
	if err != nil {
		t.Fatalf("parseGeneratedTypesEnhanced returned error: %v", err)
	}

	info, ok := typeInfo["SaveEventBridgeSettingsJSONRequestBody"]
	if !ok {
		t.Fatalf("alias SaveEventBridgeSettingsJSONRequestBody not resolved")
	}

	if _, ok := info.Fields["AWSAccountID"]; !ok {
		t.Errorf("alias missing field AWSAccountID")
	}

	if _, ok := info.Fields["AWSRegion"]; !ok {
		t.Errorf("alias missing field AWSRegion")
	}
}

// TestNormalizeResourceName はnormalizeResourceName関数のテスト
func TestNormalizeResourceName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "plural to singular",
			input: "feedbacks",
			want:  "Feedback",
		},
		{
			name:  "already singular",
			input: "feedback",
			want:  "Feedback",
		},
		{
			name:  "comments",
			input: "comments",
			want:  "Comment",
		},
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizeResourceName(tt.input)
			if got != tt.want {
				t.Errorf("normalizeResourceName(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// TestNormalizeFieldName はnormalizeFieldName関数のテスト
func TestNormalizeFieldName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "camelCase to snake_case",
			input: "userId",
			want:  "user_id",
		},
		{
			name:  "already snake_case",
			input: "user_id",
			want:  "user_id",
		},
		{
			name:  "PascalCase to snake_case",
			input: "FeedbackTitle",
			want:  "feedback_title",
		},
		{
			name:  "single word",
			input: "name",
			want:  "name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizeFieldName(tt.input)
			if got != tt.want {
				t.Errorf("normalizeFieldName(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
