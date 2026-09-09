package snapshot

import (
	"strings"
	"testing"
)

type dummyHTTPResponse struct {
	StatusCode    int
	Status        string
	ContentLength int64
	Header        map[string][]string
}

type dummyJSONPayload struct {
	ClientSecret string `json:"client_secret"`
	Message      string `json:"message"`
}

type dummyResponse struct {
	Body         []byte
	HTTPResponse *dummyHTTPResponse
	JSON200      *dummyJSONPayload
}

func (d *dummyResponse) StatusCode() int {
	if d.HTTPResponse == nil {
		return 0
	}
	return d.HTTPResponse.StatusCode
}

func (d *dummyResponse) Status() string {
	if d.HTTPResponse == nil {
		return ""
	}
	return d.HTTPResponse.Status
}

func TestMaskSensitiveDataMasksKnownKeys(t *testing.T) {
	capture := NewStorySnapshotCapture(nil)

	masked := capture.maskSensitiveData("access_token", "super-secret-token")
	if maskedStr, ok := masked.(string); !ok || !strings.HasPrefix(maskedStr, "[MASKED") {
		t.Fatalf("expected masked string, got %#v", masked)
	}

	nested := capture.maskSensitiveData("metadata", map[string]interface{}{
		"client_secret": "abcdef",
		"normal":        "value",
	})

	nestedMap, ok := nested.(map[string]interface{})
	if !ok {
		t.Fatalf("expected map result, got %#v", nested)
	}

	if _, ok := nestedMap["normal"].(string); !ok {
		t.Fatalf("expected normal field to remain a string, got %#v", nestedMap["normal"])
	}

	if secretVal, ok := nestedMap["client_secret"].(string); !ok || !strings.HasPrefix(secretVal, "[MASKED") {
		t.Fatalf("expected masked client_secret, got %#v", nestedMap["client_secret"])
	}
}

func TestCaptureStepResponseMasksSecrets(t *testing.T) {
	capture := NewStorySnapshotCapture(nil)

	response := &dummyResponse{
		Body: []byte("{\n  \"client_secret\": \"top-secret\",\n  \"message\": \"ok\"\n}\n"),
		HTTPResponse: &dummyHTTPResponse{
			StatusCode:    200,
			Status:        "200 OK",
			ContentLength: 123,
			Header: map[string][]string{
				"Authorization": []string{"Bearer abcdef"},
				"Content-Type":  []string{"application/json"},
			},
		},
		JSON200: &dummyJSONPayload{
			ClientSecret: "top-secret",
			Message:      "ok",
		},
	}

	result, err := capture.CaptureStepResponse("dummy", response)
	if err != nil {
		t.Fatalf("CaptureStepResponse returned error: %v", err)
	}

	if headerVal := result.Headers["Authorization"]; !strings.HasPrefix(headerVal, "[MASKED") {
		t.Fatalf("expected masked Authorization header, got %q", headerVal)
	}

	jsonData, ok := result.JSONData["JSON200"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected JSON200 entry to be map, got %#v", result.JSONData["JSON200"])
	}

	if secretVal, ok := jsonData["client_secret"].(string); !ok || !strings.HasPrefix(secretVal, "[MASKED") {
		t.Fatalf("expected masked client_secret in JSON, got %#v", jsonData["client_secret"])
	}

	if strings.Contains(result.Body, "top-secret") {
		t.Fatalf("expected masked body, got %q", result.Body)
	}
}
