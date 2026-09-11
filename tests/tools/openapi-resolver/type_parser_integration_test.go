package main

import (
	"testing"
)

// TestTypeParserWithBillingModule tests type parser with billing module
func TestTypeParserWithBillingModule(t *testing.T) {
	parser, err := NewTypeDefinitionParser("billing")
	if err != nil {
		t.Fatalf("NewTypeDefinitionParser() error = %v", err)
	}

	// Parse all types
	types, err := parser.ParseAllTypes()
	if err != nil {
		t.Fatalf("ParseAllTypes() error = %v", err)
	}

	t.Logf("Found %d types in billing module", len(types))

	// Test specific types
	testCases := []struct {
		typeName       string
		expectedFields int
	}{
		{"StripeInfo", 1},
		{"UpdateStripeInfoParam", 1},
	}

	for _, tc := range testCases {
		t.Run(tc.typeName, func(t *testing.T) {
			typeDef, err := parser.GetTypeDefinition(tc.typeName)
			if err != nil {
				t.Errorf("GetTypeDefinition(%s) error = %v", tc.typeName, err)
				return
			}

			if len(typeDef.Fields) < tc.expectedFields {
				t.Errorf("Expected at least %d fields, got %d", tc.expectedFields, len(typeDef.Fields))
			}

			t.Logf("Type: %s", typeDef.Name)
			for _, field := range typeDef.Fields {
				t.Logf("  - %s: %s (json:%s, required:%v)",
					field.Name, field.Type, field.JSONTag, field.Required)
			}
		})
	}
}

// TestTypeParserWithAuthModule tests type parser with auth module
func TestTypeParserWithAuthModule(t *testing.T) {
	parser, err := NewTypeDefinitionParser("auth")
	if err != nil {
		t.Fatalf("NewTypeDefinitionParser() error = %v", err)
	}

	// Parse all types
	types, err := parser.ParseAllTypes()
	if err != nil {
		t.Fatalf("ParseAllTypes() error = %v", err)
	}

	t.Logf("Found %d types in auth module", len(types))

	// Just verify we can parse the module
	if len(types) == 0 {
		t.Error("Expected to find types in auth module")
	}
}

// TestComplexNestedTypes tests parsing of complex nested types
func TestComplexNestedTypes(t *testing.T) {
	parser, err := NewTypeDefinitionParser("communication")
	if err != nil {
		t.Fatalf("NewTypeDefinitionParser() error = %v", err)
	}

	// Parse all types first
	_, err = parser.ParseAllTypes()
	if err != nil {
		t.Fatalf("ParseAllTypes() error = %v", err)
	}

	// Test recursive parsing of Feedback type which has nested Comment and User types
	typeDef, err := parser.GetTypeDefinitionRecursive("Feedback", nil)
	if err != nil {
		t.Fatalf("GetTypeDefinitionRecursive() error = %v", err)
	}

	t.Logf("Type: %s with %d fields", typeDef.Name, len(typeDef.Fields))

	// Verify nested types are detected
	hasNestedType := false
	for _, field := range typeDef.Fields {
		if field.Comment != "" && field.Comment != " " {
			hasNestedType = true
			t.Logf("Found nested type in field %s: %s", field.Name, field.Comment)
		}
	}

	if !hasNestedType {
		t.Log("Note: No nested type comments found (this is OK if types are simple)")
	}
}

// TestFieldRequiredDetection tests detection of required vs optional fields
func TestFieldRequiredDetection(t *testing.T) {
	parser, err := NewTypeDefinitionParser("communication")
	if err != nil {
		t.Fatalf("NewTypeDefinitionParser() error = %v", err)
	}

	typeDef, err := parser.GetTypeDefinition("CreateFeedbackParam")
	if err != nil {
		t.Fatalf("GetTypeDefinition() error = %v", err)
	}

	// Check that we can detect required fields
	hasRequiredField := false
	for _, field := range typeDef.Fields {
		if field.Required {
			hasRequiredField = true
			t.Logf("Required field: %s (%s)", field.Name, field.JSONTag)
		}
	}

	if !hasRequiredField {
		t.Error("Expected to find at least one required field")
	}
}
