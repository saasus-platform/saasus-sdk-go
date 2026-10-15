// Package contract round-trips the generated SDK clients against a Prism mock
// started from the latest OpenAPI spec.
//
// It checks that requests the SDK builds are accepted by the mock and that the
// responses deserialize into the SDK types. Operations are enumerated from the
// spec at run time, so an SDK that does not have an operation simply has no
// case for it — which is what makes the same suite usable for a previously
// released SDK (forward compatibility).
//
// Only the typed <Operation>WithResponse form is exercised. The generated
// alternatives are strict sub-paths of it:
//
//	OpWithResponse        -> Client.Op -> NewOpRequest -> NewOpRequestWithBody -> Do -> ParseOpResponse
//	OpWithBodyWithResponse-> Client.OpWithBody ------------^                          ^
//	Client.Op / OpWithBody-> the same request construction, without deserialization
//
// so calling them as well would execute no additional logic.
//
// Synthetic data only: the bearer token is a constant, request values come from
// the spec, and no response body, token or PII is logged.
package contract

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/saasus-platform/saasus-sdk-go/generated/apilogapi"
	"github.com/saasus-platform/saasus-sdk-go/generated/authapi"
	"github.com/saasus-platform/saasus-sdk-go/generated/awsmarketplaceapi"
	"github.com/saasus-platform/saasus-sdk-go/generated/billingapi"
	"github.com/saasus-platform/saasus-sdk-go/generated/communicationapi"
	"github.com/saasus-platform/saasus-sdk-go/generated/integrationapi"
	"github.com/saasus-platform/saasus-sdk-go/generated/pricingapi"
)

// The specs declare a bearer scheme and the mock rejects requests without
// credentials, but it does not validate the value.
const syntheticToken = "synthetic-bearer-token"

// bearerTransport supplies the credentials the mock requires. The hand-written
// authentication wrapper is deliberately bypassed: this suite verifies the
// generated client against the spec.
type bearerTransport struct{}

func (bearerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req = req.Clone(req.Context())
	req.Header.Set("Authorization", "Bearer "+syntheticToken)
	return http.DefaultTransport.RoundTrip(req)
}

var httpClient = &http.Client{Transport: bearerTransport{}}

// modules lists the generated clients. Go cannot enumerate packages at run
// time, so this is the only place that names them.
var modules = []struct {
	name   string
	client func(baseURL string) (any, error)
}{
	{"apilog", func(u string) (any, error) {
		return apilogapi.NewClientWithResponses(u, apilogapi.WithHTTPClient(httpClient))
	}},
	{"auth", func(u string) (any, error) {
		return authapi.NewClientWithResponses(u, authapi.WithHTTPClient(httpClient))
	}},
	{"awsmarketplace", func(u string) (any, error) {
		return awsmarketplaceapi.NewClientWithResponses(u, awsmarketplaceapi.WithHTTPClient(httpClient))
	}},
	{"billing", func(u string) (any, error) {
		return billingapi.NewClientWithResponses(u, billingapi.WithHTTPClient(httpClient))
	}},
	{"communication", func(u string) (any, error) {
		return communicationapi.NewClientWithResponses(u, communicationapi.WithHTTPClient(httpClient))
	}},
	{"integration", func(u string) (any, error) {
		return integrationapi.NewClientWithResponses(u, integrationapi.WithHTTPClient(httpClient))
	}},
	{"pricing", func(u string) (any, error) {
		return pricingapi.NewClientWithResponses(u, pricingapi.WithHTTPClient(httpClient))
	}},
}

// TestContract runs one case per operation (and per request-body variant).
//
// CONTRACT_SPEC_DIR points at the <module>.yml files; PRISM_URL_<MODULE> gives
// each module's mock. A module without a mock contributes no case, so the suite
// can be run against a single module locally.
func TestContract(t *testing.T) {
	specDir := os.Getenv("CONTRACT_SPEC_DIR")
	if specDir == "" {
		t.Skip("CONTRACT_SPEC_DIR is not set")
	}

	var missing []string
	ran := 0
	for _, module := range modules {
		baseURL := os.Getenv("PRISM_URL_" + strings.ToUpper(module.name))
		if baseURL == "" {
			continue
		}
		s, err := loadSpec(specDir, module.name)
		if err != nil {
			t.Fatalf("%s: cannot load spec: %v", module.name, err)
		}
		client, err := module.client(baseURL)
		if err != nil {
			t.Fatalf("%s: cannot build client: %v", module.name, err)
		}
		value := reflect.ValueOf(client)

		for _, id := range sortedOperationIDs(s) {
			op := s.ops[id]
			method := value.MethodByName(id + "WithResponse")
			if !method.IsValid() {
				missing = append(missing, module.name+"."+id)
				continue
			}
			labels := s.bodyVariants(op)
			if labels == nil {
				labels = []string{""}
			}
			for i, label := range labels {
				name := module.name + "." + id
				if label != "" {
					name += "(" + label + ")"
				}
				variantIndex := i
				t.Run(name, func(t *testing.T) {
					roundTrip(t, s, op, method, variantIndex)
				})
				ran++
			}
		}
	}

	if ran == 0 {
		t.Fatal("no case ran: no module has a mock configured")
	}
	// The latest SDK must cover every operation the spec declares. A previously
	// released SDK legitimately lacks newer operations, so this is only enforced
	// where the caller asks for it.
	if os.Getenv("CONTRACT_REQUIRE_COMPLETE") == "true" && len(missing) > 0 {
		t.Errorf("spec operations not present in the SDK: %v", missing)
	}
}

func roundTrip(t *testing.T, s *spec, op *operation, method reflect.Value, variantIndex int) {
	out := method.Call(buildArgs(t, s, op, method.Type(), variantIndex))
	if err, _ := out[1].Interface().(error); err != nil {
		t.Fatalf("round trip failed: %v", err)
	}

	response := out[0].Elem()
	status := response.FieldByName("HTTPResponse").Interface().(*http.Response).StatusCode
	if status < 200 || status >= 300 {
		// An operation that declares no success response is expected to surface
		// the error it does declare.
		if op.hasSuccess || !op.statuses[status] {
			t.Fatalf("unexpected status %d", status)
		}
		return
	}

	body := response.FieldByName(fmt.Sprintf("JSON%d", status))
	if !body.IsValid() {
		return // the spec declares no JSON body for this status
	}
	if body.IsNil() {
		t.Fatalf("status %d: response was not deserialized into the SDK type", status)
	}
	if err := checkUnions(body); err != nil {
		t.Fatal(err)
	}
}

// buildArgs fills the generated argument list: context, the path parameters in
// path-template order, the query parameter struct, and the request body. That
// order and composition is what the generator emits, so a mismatch in the
// number of arguments means the spec and the SDK disagree about the operation.
func buildArgs(t *testing.T, s *spec, op *operation, method reflect.Type, variantIndex int) []reflect.Value {
	want := 1 + len(op.pathParams) + 1 // context + path parameters + variadic request editors
	if len(op.queryParams) > 0 {
		want++
	}
	if op.bodySchema != nil {
		want++
	}
	if method.NumIn() != want {
		t.Fatalf("%s: method takes %d arguments, the spec implies %d", op.id, method.NumIn(), want)
	}

	args := []reflect.Value{reflect.ValueOf(context.Background())}
	for _, p := range op.pathParams {
		args = append(args, decode(t, method.In(len(args)), s.paramValue(p)))
	}
	if len(op.queryParams) > 0 {
		args = append(args, decode(t, method.In(len(args)), s.queryValues(op)))
	}
	if op.bodySchema != nil {
		args = append(args, decode(t, method.In(len(args)), s.bodyValue(op, variantIndex)))
	}
	return args
}

// decode converts a spec-derived value into the generated argument type through
// JSON, which is the same conversion the SDK applies to a request body. The
// generated query parameter struct tags each field with its spec name, so a
// JSON object keyed by spec names binds by name rather than by position.
func decode(t *testing.T, target reflect.Type, value any) reflect.Value {
	raw, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("cannot encode the spec value for %s: %v", target, err)
	}
	out := reflect.New(target)
	if err := json.Unmarshal(raw, out.Interface()); err != nil {
		t.Fatalf("the spec value does not fit %s: %v", target, err)
	}
	return out.Elem()
}

// checkUnions resolves every oneOf value in the response.
//
// A oneOf is generated as a holder whose UnmarshalJSON copies the raw bytes
// without inspecting them, so deserializing it always succeeds and proves
// nothing. ValueByDiscriminator is where the variant is actually selected and
// parsed, and it fails on a discriminator the SDK does not know — which is the
// forward-compatibility break this detects.
func checkUnions(value reflect.Value) error {
	if method := value.MethodByName("ValueByDiscriminator"); method.IsValid() {
		out := method.Call(nil)
		if err, _ := out[1].Interface().(error); err != nil {
			return fmt.Errorf("oneOf value could not be resolved: %w", err)
		}
		return checkUnions(reflect.ValueOf(out[0].Interface()))
	}
	switch value.Kind() {
	case reflect.Pointer, reflect.Interface:
		if !value.IsNil() {
			return checkUnions(value.Elem())
		}
	case reflect.Struct:
		for i := 0; i < value.NumField(); i++ {
			if field := value.Field(i); field.CanInterface() {
				if err := checkUnions(field); err != nil {
					return err
				}
			}
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < value.Len(); i++ {
			if err := checkUnions(value.Index(i)); err != nil {
				return err
			}
		}
	case reflect.Map:
		for _, key := range value.MapKeys() {
			if err := checkUnions(value.MapIndex(key)); err != nil {
				return err
			}
		}
	}
	return nil
}

func sortedOperationIDs(s *spec) []string {
	ids := make([]string, 0, len(s.ops))
	for id := range s.ops {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}
