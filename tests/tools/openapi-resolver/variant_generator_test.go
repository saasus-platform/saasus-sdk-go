//go:build test
// +build test

package main

import (
	"reflect"
	"testing"
)

func TestMethodVariantGenerator_GenerateVariants(t *testing.T) {
	generator := NewMethodVariantGenerator()

	tests := []struct {
		name        string
		operationID string
		httpMethod  string
		hasBody     bool
		want        []string
	}{
		{
			name:        "GET method without body",
			operationID: "GetStripeInfo",
			httpMethod:  "GET",
			hasBody:     false,
			want: []string{
				"GetStripeInfo",
				"GetStripeInfoWithResponse",
			},
		},
		{
			name:        "DELETE method without body",
			operationID: "DeleteStripeInfo",
			httpMethod:  "DELETE",
			hasBody:     false,
			want: []string{
				"DeleteStripeInfo",
				"DeleteStripeInfoWithResponse",
			},
		},
		{
			name:        "POST method with body",
			operationID: "CreateTenant",
			httpMethod:  "POST",
			hasBody:     true,
			want: []string{
				"CreateTenant",
				"CreateTenantWithResponse",
				"CreateTenantWithBody",
				"CreateTenantWithBodyWithResponse",
			},
		},
		{
			name:        "PUT method with body",
			operationID: "UpdateStripeInfo",
			httpMethod:  "PUT",
			hasBody:     true,
			want: []string{
				"UpdateStripeInfo",
				"UpdateStripeInfoWithResponse",
				"UpdateStripeInfoWithBody",
				"UpdateStripeInfoWithBodyWithResponse",
			},
		},
		{
			name:        "PATCH method with body",
			operationID: "PatchResource",
			httpMethod:  "PATCH",
			hasBody:     true,
			want: []string{
				"PatchResource",
				"PatchResourceWithResponse",
				"PatchResourceWithBody",
				"PatchResourceWithBodyWithResponse",
			},
		},
		{
			name:        "POST method without body",
			operationID: "TriggerAction",
			httpMethod:  "POST",
			hasBody:     false,
			want: []string{
				"TriggerAction",
				"TriggerActionWithResponse",
			},
		},
		{
			name:        "Empty operationID",
			operationID: "",
			httpMethod:  "GET",
			hasBody:     false,
			want:        []string{},
		},
		{
			name:        "GET method with body (edge case)",
			operationID: "GetWithBody",
			httpMethod:  "GET",
			hasBody:     true,
			want: []string{
				"GetWithBody",
				"GetWithBodyWithResponse",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := generator.GenerateVariants(tt.operationID, tt.httpMethod, tt.hasBody)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateVariants() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsBodyMethod(t *testing.T) {
	tests := []struct {
		name   string
		method string
		want   bool
	}{
		{"POST supports body", "POST", true},
		{"PUT supports body", "PUT", true},
		{"PATCH supports body", "PATCH", true},
		{"GET does not support body", "GET", false},
		{"DELETE does not support body", "DELETE", false},
		{"HEAD does not support body", "HEAD", false},
		{"OPTIONS does not support body", "OPTIONS", false},
		{"Lowercase post", "post", false}, // Case sensitive
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isBodyMethod(tt.method)
			if got != tt.want {
				t.Errorf("isBodyMethod(%q) = %v, want %v", tt.method, got, tt.want)
			}
		})
	}
}
