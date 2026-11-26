package sdkuseragent_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/saasus-platform/saasus-sdk-go/generated/authapi"
	"github.com/saasus-platform/saasus-sdk-go/internal/sdkuseragent"
	"github.com/saasus-platform/saasus-sdk-go/modules/auth"
)

type doerFunc func(*http.Request) (*http.Response, error)

func (f doerFunc) Do(req *http.Request) (*http.Response, error) {
	return f(req)
}

const product = "SaaSus-SDK-go/" + sdkuseragent.Version

func TestSet(t *testing.T) {
	for _, test := range []struct {
		name, userAgent, want string
	}{
		{"default", "", product},
		{"custom", "Example/1.0", "Example/1.0 " + product},
		{"already present", "Example/1.0 " + product, "Example/1.0 " + product},
	} {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/v1/auth/auth-info", nil)
			req.Header.Set("User-Agent", test.userAgent)
			sdkuseragent.Set(req)
			if got := req.UserAgent(); got != test.want {
				t.Errorf("User-Agent = %q, want %q", got, test.want)
			}
		})
	}
}

func TestGetAuthInfoUserAgent(t *testing.T) {
	t.Setenv("SAASUS_SAAS_ID", "test-saas")
	t.Setenv("SAASUS_API_KEY", "test-key")
	t.Setenv("SAASUS_SECRET_KEY", "test-secret")

	client, err := auth.AuthWithResponse()
	if err != nil {
		t.Fatal(err)
	}
	rawClient := client.ClientInterface.(*authapi.Client)
	rawClient.Client = doerFunc(func(req *http.Request) (*http.Response, error) {
		if got := req.UserAgent(); got != product {
			t.Errorf("sent User-Agent = %q, want %q", got, product)
		}
		if !strings.HasPrefix(req.Header.Get("Authorization"), "SAASUSSIGV1 ") {
			t.Error("request is missing the SaaSus signature")
		}
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader("{}"))}, nil
	})
	response, err := client.GetAuthInfo(context.Background(), func(_ context.Context, req *http.Request) error {
		if got := req.UserAgent(); got != product {
			t.Errorf("User-Agent before the per-call editor = %q, want %q", got, product)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
}
