package main

import (
	"strings"
	"testing"
)

func TestNewClientInterfaceParser(t *testing.T) {
	tests := []struct {
		name       string
		moduleName string
		wantErr    bool
	}{
		{
			name:       "valid module - billing",
			moduleName: "billing",
			wantErr:    false,
		},
		{
			name:       "valid module - auth",
			moduleName: "auth",
			wantErr:    false,
		},
		{
			name:       "valid module - communication",
			moduleName: "communication",
			wantErr:    false,
		},
		{
			name:       "valid module - pricing",
			moduleName: "pricing",
			wantErr:    false,
		},
		{
			name:       "invalid module",
			moduleName: "nonexistent",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser, err := NewClientInterfaceParser(tt.moduleName)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClientInterfaceParser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && parser == nil {
				t.Error("NewClientInterfaceParser() returned nil parser")
			}
			if !tt.wantErr && parser.moduleName != tt.moduleName {
				t.Errorf("NewClientInterfaceParser() moduleName = %v, want %v", parser.moduleName, tt.moduleName)
			}
		})
	}
}

func TestParseAllMethods(t *testing.T) {
	tests := []struct {
		name            string
		moduleName      string
		expectedMethods []string
		minMethodCount  int
	}{
		{
			name:       "billing module",
			moduleName: "billing",
			expectedMethods: []string{
				"GetStripeInfo",
				"UpdateStripeInfo",
				"DeleteStripeInfo",
			},
			minMethodCount: 3,
		},
		{
			name:       "communication module",
			moduleName: "communication",
			expectedMethods: []string{
				"CreateFeedback",
				"GetFeedbacks",
				"GetFeedback",
				"UpdateFeedback",
				"DeleteFeedback",
			},
			minMethodCount: 5,
		},
		{
			name:       "auth module",
			moduleName: "auth",
			expectedMethods: []string{
				"GetUserInfo",
				"GetSaasUsers",
				"CreateSaasUser",
			},
			minMethodCount: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser, err := NewClientInterfaceParser(tt.moduleName)
			if err != nil {
				t.Fatalf("NewClientInterfaceParser() error = %v", err)
			}

			methods, err := parser.ParseAllMethods()
			if err != nil {
				t.Fatalf("ParseAllMethods() error = %v", err)
			}

			if len(methods) == 0 {
				t.Error("ParseAllMethods() returned no methods")
			}

			if len(methods) < tt.minMethodCount {
				t.Errorf("ParseAllMethods() returned %d methods, expected at least %d", len(methods), tt.minMethodCount)
			}

			methodMap := make(map[string]bool)
			for _, method := range methods {
				methodMap[method.MethodName] = true

				// メソッドシグネチャの基本的な検証
				if method.MethodName == "" {
					t.Error("Method has empty name")
				}
				if method.OperationID == "" {
					t.Error("Method has empty operationID")
				}

				// パラメータの検証
				if len(method.ParameterTypes) == 0 {
					t.Errorf("Method %s has no parameters (expected at least context)", method.MethodName)
				}
			}

			for _, expected := range tt.expectedMethods {
				if !methodMap[expected] {
					t.Errorf("Expected method %s not found in parsed methods", expected)
				}
			}

			t.Logf("Successfully parsed %d methods from %s module", len(methods), tt.moduleName)
		})
	}
}

func TestFindMethodByOperationID(t *testing.T) {
	tests := []struct {
		name        string
		moduleName  string
		operationID string
		wantErr     bool
	}{
		{
			name:        "billing - existing operation",
			moduleName:  "billing",
			operationID: "GetStripeInfo",
			wantErr:     false,
		},
		{
			name:        "billing - non-existing operation",
			moduleName:  "billing",
			operationID: "NonExistentMethod",
			wantErr:     true,
		},
		{
			name:        "communication - CreateFeedback",
			moduleName:  "communication",
			operationID: "CreateFeedback",
			wantErr:     false,
		},
		{
			name:        "communication - GetFeedbacks",
			moduleName:  "communication",
			operationID: "GetFeedbacks",
			wantErr:     false,
		},
		{
			name:        "auth - GetUserInfo",
			moduleName:  "auth",
			operationID: "GetUserInfo",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser, err := NewClientInterfaceParser(tt.moduleName)
			if err != nil {
				t.Fatalf("NewClientInterfaceParser() error = %v", err)
			}

			method, err := parser.FindMethodByOperationID(tt.operationID)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindMethodByOperationID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if method == nil {
					t.Error("FindMethodByOperationID() returned nil method")
					return
				}
				if method.MethodName == "" {
					t.Error("Method has empty name")
				}
				if len(method.ParameterTypes) == 0 {
					t.Error("Method has no parameters (expected at least context)")
				}

				t.Logf("Found method: %s with %d parameters", method.MethodName, len(method.ParameterTypes))
			}
		})
	}
}

func TestMethodSignatureExtraction(t *testing.T) {
	tests := []struct {
		name            string
		moduleName      string
		operationID     string
		expectContext   bool
		expectBodyParam bool
		minParamCount   int
	}{
		{
			name:            "billing - GetStripeInfo (no body)",
			moduleName:      "billing",
			operationID:     "GetStripeInfo",
			expectContext:   true,
			expectBodyParam: false,
			minParamCount:   1,
		},
		{
			name:            "billing - UpdateStripeInfo (with body)",
			moduleName:      "billing",
			operationID:     "UpdateStripeInfo",
			expectContext:   true,
			expectBodyParam: true,
			minParamCount:   2,
		},
		{
			name:            "communication - CreateFeedback (with body)",
			moduleName:      "communication",
			operationID:     "CreateFeedback",
			expectContext:   true,
			expectBodyParam: true,
			minParamCount:   2,
		},
		{
			name:            "communication - GetFeedbacks (no body)",
			moduleName:      "communication",
			operationID:     "GetFeedbacks",
			expectContext:   true,
			expectBodyParam: false,
			minParamCount:   1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser, err := NewClientInterfaceParser(tt.moduleName)
			if err != nil {
				t.Fatalf("NewClientInterfaceParser() error = %v", err)
			}

			method, err := parser.FindMethodByOperationID(tt.operationID)
			if err != nil {
				t.Fatalf("FindMethodByOperationID() error = %v", err)
			}

			// パラメータ数の検証
			if len(method.ParameterTypes) < tt.minParamCount {
				t.Errorf("Method has %d parameters, expected at least %d", len(method.ParameterTypes), tt.minParamCount)
			}

			// 最初のパラメータはcontextであるべき
			if tt.expectContext && len(method.ParameterTypes) > 0 {
				firstParam := method.ParameterTypes[0]
				if firstParam.Type != "context.Context" {
					t.Errorf("First parameter type = %v, want context.Context", firstParam.Type)
				}
			}

			// bodyパラメータの検証
			if tt.expectBodyParam {
				hasBodyParam := false
				for _, param := range method.ParameterTypes {
					if strings.Contains(param.Type, "JSONRequestBody") ||
						strings.Contains(param.Type, "RequestBody") ||
						param.Name == "body" {
						hasBodyParam = true
						break
					}
				}
				if !hasBodyParam {
					t.Error("Expected body parameter but none found")
				}
			}

			// 戻り値の型が設定されているべき
			if method.ReturnType == "" {
				t.Error("Method has no return type")
			}

			// パラメータの詳細をログ出力
			t.Logf("Method: %s", method.MethodName)
			t.Logf("Return Type: %s", method.ReturnType)
			for i, param := range method.ParameterTypes {
				t.Logf("  Param %d: %s %s", i, param.Name, param.Type)
			}
		})
	}
}

// TestParameterTypeExtraction tests accurate extraction of parameter types
func TestParameterTypeExtraction(t *testing.T) {
	parser, err := NewClientInterfaceParser("communication")
	if err != nil {
		t.Fatalf("NewClientInterfaceParser() error = %v", err)
	}

	method, err := parser.FindMethodByOperationID("CreateFeedback")
	if err != nil {
		t.Fatalf("FindMethodByOperationID() error = %v", err)
	}

	// CreateFeedbackは少なくともcontext.ContextとCreateFeedbackJSONRequestBodyを持つべき
	foundContext := false
	foundBody := false
	var bodyTypeName string

	for _, param := range method.ParameterTypes {
		if param.Type == "context.Context" {
			foundContext = true
		}
		if strings.Contains(param.Type, "CreateFeedback") && strings.Contains(param.Type, "JSONRequestBody") {
			foundBody = true
			bodyTypeName = param.Type
		}
	}

	if !foundContext {
		t.Error("context.Context parameter not found")
	}

	if !foundBody {
		t.Error("CreateFeedbackJSONRequestBody parameter not found")
	}

	if bodyTypeName != "" {
		t.Logf("Found body parameter type: %s", bodyTypeName)
	}
}
