package contract

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/saasus-platform/saasus-sdk-go/generated/pricingapi"
)

// TestCheckUnions covers the one assertion whose effect is not visible from a
// passing round trip: a oneOf value deserializes even when its variant is
// unknown, so only checkUnions can tell the two apart.
func TestCheckUnions(t *testing.T) {
	tests := []struct {
		name      string
		body      string
		wantError bool
	}{
		{"known variant", `{"units":[{"type":"fixed","name":"n","display_name":"d","unit_amount":1,"currency":"USD","recurring_interval":"month"}]}`, false},
		{"unknown variant", `{"units":[{"type":"metered"}]}`, true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var units pricingapi.PricingUnits
			if err := json.Unmarshal([]byte(test.body), &units); err != nil {
				t.Fatalf("deserialization failed before the union check: %v", err)
			}
			err := checkUnions(reflect.ValueOf(units))
			if test.wantError && err == nil {
				t.Error("an unknown oneOf variant was not detected")
			}
			if !test.wantError && err != nil {
				t.Errorf("a known oneOf variant was rejected: %v", err)
			}
		})
	}
}
