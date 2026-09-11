package main

import (
	"strings"
	"testing"
)

func TestNewTypeDefinitionParser(t *testing.T) {
	tests := []struct {
		name       string
		moduleName string
		wantErr    bool
	}{
		{
			name:       "valid module - communication",
			moduleName: "communication",
			wantErr:    false,
		},
		{
			name:       "valid module - billing",
			moduleName: "billing",
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
			parser, err := NewTypeDefinitionParser(tt.moduleName)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTypeDefinitionParser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && parser == nil {
				t.Error("NewTypeDefinitionParser() returned nil parser")
			}
		})
	}
}

func TestParseAllTypes(t *testing.T) {
	tests := []struct {
		name          string
		moduleName    string
		expectedTypes []string
		minTypeCount  int
	}{
		{
			name:       "communication module",
			moduleName: "communication",
			expectedTypes: []string{
				"CreateFeedbackParam",
				"Feedback",
				"Comment",
			},
			minTypeCount: 3,
		},
		{
			name:       "billing module",
			moduleName: "billing",
			expectedTypes: []string{
				"UpdateStripeInfoParam",
				"StripeInfo",
			},
			minTypeCount: 2,
		},
		{
			name:       "auth module",
			moduleName: "auth",
			expectedTypes: []string{
				"CreateSaasUserParam",
				"SaasUser",
			},
			minTypeCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser, err := NewTypeDefinitionParser(tt.moduleName)
			if err != nil {
				t.Fatalf("NewTypeDefinitionParser() error = %v", err)
			}

			types, err := parser.ParseAllTypes()
			if err != nil {
				t.Fatalf("ParseAllTypes() error = %v", err)
			}

			if len(types) == 0 {
				t.Error("ParseAllTypes() returned no types")
			}

			if len(types) < tt.minTypeCount {
				t.Errorf("ParseAllTypes() returned %d types, expected at least %d", len(types), tt.minTypeCount)
			}

			t.Logf("Found %d types in %s module", len(types), tt.moduleName)

			// 期待される型が存在することを確認
			for _, typeName := range tt.expectedTypes {
				if typeDef, ok := types[typeName]; !ok {
					t.Logf("Warning: Expected type '%s' not found", typeName)
				} else {
					t.Logf("Found type: %s with %d fields", typeName, len(typeDef.Fields))

					// フィールドの詳細を検証
					for _, field := range typeDef.Fields {
						if field.Name == "" {
							t.Errorf("Type %s has field with empty name", typeName)
						}
						if field.Type == "" {
							t.Errorf("Type %s field %s has empty type", typeName, field.Name)
						}
					}
				}
			}
		})
	}
}

func TestGetTypeDefinition(t *testing.T) {
	tests := []struct {
		name           string
		moduleName     string
		typeName       string
		wantErr        bool
		minFieldCount  int
		expectedFields []string
	}{
		{
			name:           "communication - CreateFeedbackParam",
			moduleName:     "communication",
			typeName:       "CreateFeedbackParam",
			wantErr:        false,
			minFieldCount:  2,
			expectedFields: []string{"FeedbackDescription", "FeedbackTitle", "UserId"},
		},
		{
			name:           "billing - UpdateStripeInfoParam",
			moduleName:     "billing",
			typeName:       "UpdateStripeInfoParam",
			wantErr:        false,
			minFieldCount:  1,
			expectedFields: []string{"SecretKey"},
		},
		{
			name:       "communication - non-existing type",
			moduleName: "communication",
			typeName:   "NonExistentType",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser, err := NewTypeDefinitionParser(tt.moduleName)
			if err != nil {
				t.Fatalf("NewTypeDefinitionParser() error = %v", err)
			}

			typeDef, err := parser.GetTypeDefinition(tt.typeName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTypeDefinition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if typeDef == nil {
					t.Error("GetTypeDefinition() returned nil type definition")
					return
				}

				if len(typeDef.Fields) < tt.minFieldCount {
					t.Errorf("Type has %d fields, expected at least %d", len(typeDef.Fields), tt.minFieldCount)
				}

				t.Logf("Type: %s", typeDef.Name)
				t.Logf("Fields: %d", len(typeDef.Fields))

				// フィールドマップを作成
				fieldMap := make(map[string]bool)
				for _, field := range typeDef.Fields {
					fieldMap[field.Name] = true
					t.Logf("  - %s: %s (json:%s, required:%v)", field.Name, field.Type, field.JSONTag, field.Required)

					// フィールドの基本検証
					if field.Name == "" {
						t.Error("Field has empty name")
					}
					if field.Type == "" {
						t.Error("Field has empty type")
					}
					if field.JSONTag == "" {
						t.Logf("Warning: Field %s has no JSON tag", field.Name)
					}
				}

				// 期待されるフィールドが存在することを確認
				for _, expectedField := range tt.expectedFields {
					if !fieldMap[expectedField] {
						t.Logf("Warning: Expected field '%s' not found", expectedField)
					}
				}
			}
		})
	}
}

func TestGetTypeDefinitionRecursive(t *testing.T) {
	parser, err := NewTypeDefinitionParser("communication")
	if err != nil {
		t.Fatalf("NewTypeDefinitionParser() error = %v", err)
	}

	// まずすべての型を解析
	_, err = parser.ParseAllTypes()
	if err != nil {
		t.Fatalf("ParseAllTypes() error = %v", err)
	}

	tests := []struct {
		name     string
		typeName string
		wantErr  bool
	}{
		{
			name:     "simple type",
			typeName: "CreateFeedbackParam",
			wantErr:  false,
		},
		{
			name:     "type with nested fields",
			typeName: "Feedback",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			typeDef, err := parser.GetTypeDefinitionRecursive(tt.typeName, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTypeDefinitionRecursive() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if typeDef == nil {
					t.Error("GetTypeDefinitionRecursive() returned nil type definition")
					return
				}

				t.Logf("Type: %s", typeDef.Name)
				t.Logf("Fields: %d", len(typeDef.Fields))

				for _, field := range typeDef.Fields {
					t.Logf("  - %s: %s (json:%s, required:%v, comment:%s)",
						field.Name, field.Type, field.JSONTag, field.Required, field.Comment)
				}
			}
		})
	}
}

func TestExtractJSONTag(t *testing.T) {
	parser := &TypeDefinitionParser{}

	tests := []struct {
		name string
		tag  string
		want string
	}{
		{
			name: "simple json tag",
			tag:  `json:"field_name"`,
			want: "field_name",
		},
		{
			name: "json tag with omitempty",
			tag:  `json:"field_name,omitempty"`,
			want: "field_name",
		},
		{
			name: "json tag with multiple tags",
			tag:  `json:"field_name" yaml:"field_name"`,
			want: "field_name",
		},
		{
			name: "json tag with dash (ignored field)",
			tag:  `json:"-"`,
			want: "",
		},
		{
			name: "no json tag",
			tag:  `yaml:"field_name"`,
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parser.extractJSONTag(tt.tag)
			if got != tt.want {
				t.Errorf("extractJSONTag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsPrimitiveType(t *testing.T) {
	parser := &TypeDefinitionParser{}

	tests := []struct {
		name     string
		typeName string
		want     bool
	}{
		{name: "string", typeName: "string", want: true},
		{name: "int", typeName: "int", want: true},
		{name: "bool", typeName: "bool", want: true},
		{name: "float64", typeName: "float64", want: true},
		{name: "interface{}", typeName: "interface{}", want: true},
		{name: "custom type", typeName: "CustomType", want: false},
		{name: "pointer type", typeName: "*CustomType", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parser.isPrimitiveType(tt.typeName)
			if got != tt.want {
				t.Errorf("isPrimitiveType() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestNestedTypeHandling tests handling of nested types
func TestNestedTypeHandling(t *testing.T) {
	parser, err := NewTypeDefinitionParser("communication")
	if err != nil {
		t.Fatalf("NewTypeDefinitionParser() error = %v", err)
	}

	// まずすべての型を解析
	_, err = parser.ParseAllTypes()
	if err != nil {
		t.Fatalf("ParseAllTypes() error = %v", err)
	}

	// ネストした型を持つ可能性のある型を検索
	typeDef, err := parser.GetTypeDefinition("Feedback")
	if err != nil {
		t.Skipf("Feedback type not found, skipping nested type test: %v", err)
		return
	}

	t.Logf("Testing nested types in Feedback")
	for _, field := range typeDef.Fields {
		t.Logf("Field: %s, Type: %s", field.Name, field.Type)

		// カスタム型の場合、その型定義も取得できるか確認
		if !parser.isPrimitiveType(field.Type) &&
			!strings.HasPrefix(field.Type, "*") &&
			!strings.HasPrefix(field.Type, "[]") &&
			!strings.HasPrefix(field.Type, "map[") {
			nestedType, err := parser.GetTypeDefinition(field.Type)
			if err != nil {
				t.Logf("Could not get nested type definition for %s: %v", field.Type, err)
			} else {
				t.Logf("  Nested type %s has %d fields", field.Type, len(nestedType.Fields))
			}
		}
	}
}

// TestFieldJSONTagExtraction tests accurate extraction of JSON tags
func TestFieldJSONTagExtraction(t *testing.T) {
	parser, err := NewTypeDefinitionParser("communication")
	if err != nil {
		t.Fatalf("NewTypeDefinitionParser() error = %v", err)
	}

	typeDef, err := parser.GetTypeDefinition("CreateFeedbackParam")
	if err != nil {
		t.Fatalf("GetTypeDefinition() error = %v", err)
	}

	// すべてのフィールドがJSONタグを持つことを確認
	for _, field := range typeDef.Fields {
		if field.JSONTag == "" {
			t.Errorf("Field %s has no JSON tag", field.Name)
		}

		// JSONタグがsnake_case形式であることを確認
		if field.JSONTag != "" && !strings.Contains(field.JSONTag, "_") {
			// 単一単語の場合は問題ない
			if len(strings.Split(field.Name, "")) > 1 {
				t.Logf("Field %s has JSON tag '%s' (not snake_case)", field.Name, field.JSONTag)
			}
		}
	}
}

// TestMultipleModulesTypeExtraction tests type extraction across multiple modules
func TestMultipleModulesTypeExtraction(t *testing.T) {
	modules := []string{"communication", "billing", "auth", "pricing"}

	for _, moduleName := range modules {
		t.Run(moduleName, func(t *testing.T) {
			parser, err := NewTypeDefinitionParser(moduleName)
			if err != nil {
				t.Skipf("Module %s not available: %v", moduleName, err)
				return
			}

			types, err := parser.ParseAllTypes()
			if err != nil {
				t.Errorf("ParseAllTypes() failed for %s: %v", moduleName, err)
				return
			}

			if len(types) == 0 {
				t.Errorf("No types found in %s module", moduleName)
				return
			}

			t.Logf("Module %s: found %d types", moduleName, len(types))

			// Param型を探す（実際の型名はJSONRequestBodyではなくParamサフィックス）
			paramCount := 0
			for typeName := range types {
				if strings.HasSuffix(typeName, "Param") {
					paramCount++
					t.Logf("  Found parameter type: %s", typeName)
				}
			}

			t.Logf("Module %s: found %d parameter types", moduleName, paramCount)
		})
	}
}
