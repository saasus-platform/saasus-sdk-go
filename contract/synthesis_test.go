package contract

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// These tests cover the input synthesis on its own, without a mock: the specs
// in use declare neither header parameters nor anyOf bodies, so TestContract
// would not exercise either path.

// writeSpec writes a spec document to a temporary directory and loads it
// through the same entry point the suite uses.
func writeSpec(t *testing.T, document string) *spec {
	t.Helper()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "test.yml"), []byte(document), 0o600); err != nil {
		t.Fatal(err)
	}
	s, err := loadSpec(dir, "test")
	if err != nil {
		t.Fatal(err)
	}
	return s
}

// TestParamsValuesIncludeHeaders checks that header parameters are synthesized
// alongside the query ones: the generator gathers both into the single Params
// struct argument, so dropping the headers would leave that argument
// incomplete, and an operation whose only parameters are headers would get no
// argument at all.
func TestParamsValuesIncludeHeaders(t *testing.T) {
	s := writeSpec(t, `
openapi: 3.0.0
paths:
  /tenants/{tenant_id}/users:
    parameters:
      - name: tenant_id
        in: path
        schema:
          type: string
        example: tenant-1
    get:
      operationId: GetUsers
      parameters:
        - name: limit
          in: query
          schema:
            type: integer
          example: 10
        - name: X-Saasus-Referer
          in: header
          schema:
            type: string
          example: referer-1
      responses:
        200:
          description: ok
  /ping:
    get:
      operationId: Ping
      parameters:
        - name: X-Request-Id
          in: header
          schema:
            type: string
          example: req-1
      responses:
        200:
          description: ok
`)

	tests := []struct {
		id   string
		want map[string]any
	}{
		{"GetUsers", map[string]any{"limit": 10, "X-Saasus-Referer": "referer-1"}},
		{"Ping", map[string]any{"X-Request-Id": "req-1"}},
	}
	for _, test := range tests {
		op, ok := s.ops[test.id]
		if !ok {
			t.Fatalf("%s: not indexed", test.id)
		}
		if got := s.paramsValues(op); !reflect.DeepEqual(got, test.want) {
			t.Errorf("%s: params = %v, want %v", test.id, got, test.want)
		}
	}

	// The path parameter must stay a positional argument rather than join the
	// Params struct.
	if got := len(s.ops["GetUsers"].pathParams); got != 1 {
		t.Errorf("GetUsers: %d path parameters, want 1", got)
	}
}

// TestBodyVariantsCoverOneOfAndAnyOf checks that every member of a composite
// request body is sent, for anyOf as well as oneOf. A body whose members are
// only reachable through anyOf would otherwise be synthesized as an empty
// object.
func TestBodyVariantsCoverOneOfAndAnyOf(t *testing.T) {
	const document = `
openapi: 3.0.0
paths:
  /units:
    post:
      operationId: SaveUnit
      requestBody:
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/UnitForSave'
      responses:
        200:
          description: ok
components:
  schemas:
    UnitForSave:
      %s:
        - $ref: '#/components/schemas/FixedUnit'
        - $ref: '#/components/schemas/UsageUnit'
      discriminator:
        propertyName: type
        mapping:
          fixed: '#/components/schemas/FixedUnit'
          usage: '#/components/schemas/UsageUnit'
    FixedUnit:
      type: object
      properties:
        name:
          type: string
          example: fixed-unit
    UsageUnit:
      type: object
      properties:
        name:
          type: string
          example: usage-unit
`

	for _, keyword := range []string{"oneOf", "anyOf"} {
		s := writeSpec(t, fmt.Sprintf(document, keyword))
		op := s.ops["SaveUnit"]

		wantLabels := []string{"fixed", "usage"}
		if got := s.bodyVariants(op); !reflect.DeepEqual(got, wantLabels) {
			t.Errorf("%s: variants = %v, want %v", keyword, got, wantLabels)
		}

		wantBody := map[string]any{"name": "usage-unit", "type": "usage"}
		if got := s.bodyValue(op, 1); !reflect.DeepEqual(got, wantBody) {
			t.Errorf("%s: body = %v, want %v", keyword, got, wantBody)
		}
	}
}

// TestBuildArgsToleratesMissingParamsArgument: when the spec adds an
// operation's first optional query/header parameter, an SDK generated before it
// has no Params argument. That backward-compatible change must not fail;
// buildArgs attributes the missing argument to Params and omits it.
func TestBuildArgsToleratesMissingParamsArgument(t *testing.T) {
	s := writeSpec(t, `
openapi: 3.0.0
paths:
  /tenants/{tenant_id}/users:
    parameters:
      - name: tenant_id
        in: path
        schema:
          type: string
        example: tenant-1
    get:
      operationId: GetUsers
      parameters:
        - name: limit
          in: query
          schema:
            type: integer
          example: 10
      responses:
        200:
          description: ok
`)
	op := s.ops["GetUsers"]
	if len(op.params) != 1 {
		t.Fatalf("GetUsers: %d params, want 1", len(op.params))
	}

	// Old SDK signature (context, tenant_id, ...editors): no Params argument,
	// which the current spec now implies.
	oldMethod := reflect.TypeOf(func(context.Context, string, ...func(context.Context, *http.Request) error) (*struct{}, error) {
		return nil, nil
	})

	args := buildArgs(t, s, op, oldMethod, 0)
	if len(args) != 2 { // context + path parameter; the Params argument is skipped
		t.Fatalf("got %d args, want 2 (context + path parameter, Params skipped)", len(args))
	}
	if _, ok := args[0].Interface().(context.Context); !ok {
		t.Errorf("first argument is %s, want context.Context", args[0].Type())
	}
	if got := args[1].Interface(); got != "tenant-1" {
		t.Errorf("path argument = %v, want tenant-1", got)
	}
}
