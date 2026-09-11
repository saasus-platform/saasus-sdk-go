package main

import "testing"

func TestMapToGeneratedTypes_UsesReferencedModules(t *testing.T) {
	paramGroups := map[string]*ParamGroup{
		"tenants": {
			Name:             "tenants",
			CreateParams:     map[string]any{"attributes": nil, "back_office_staff_email": nil, "name": nil},
			UpdateParams:     map[string]any{},
			CreateEndpoint:   "/v1/tenants",
			CreateMethod:     "POST",
			ModuleCandidates: []string{"auth"},
		},
	}

	moduleOrder := []string{"pricing", "auth"}

	if err := mapToGeneratedTypes("pricing", paramGroups, moduleOrder); err != nil {
		t.Fatalf("mapToGeneratedTypes returned error: %v", err)
	}

	group := paramGroups["tenants"]

	email, ok := group.CreateParams["back_office_staff_email"].(string)
	if !ok {
		t.Fatalf("expected email field to be mapped to string, got %#v", group.CreateParams["back_office_staff_email"])
	}
	if email != "test@example.com" {
		t.Fatalf("unexpected email sample value: %q", email)
	}

	name, ok := group.CreateParams["name"].(string)
	if !ok {
		t.Fatalf("expected name field to be mapped to string, got %#v", group.CreateParams["name"])
	}
	if name != "テスト用名前" {
		t.Fatalf("unexpected name sample value: %q", name)
	}

	attributes, ok := group.CreateParams["attributes"].(map[string]any)
	if !ok {
		t.Fatalf("expected attributes field to be mapped to map, got %#v", group.CreateParams["attributes"])
	}
	if len(attributes) != 0 {
		t.Fatalf("expected attributes map to be empty, got %#v", attributes)
	}
}
