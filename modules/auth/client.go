package auth

import (
	"context"
	"net/http"

	"github.com/saasus-platform/saasus-sdk-go/client"
	"github.com/saasus-platform/saasus-sdk-go/generated/authapi"
)

var (
	server = "https://api.saasus.io/v1/auth"
)

func withRequestEditorFns(c *authapi.Client) error {
	c.RequestEditors = []authapi.RequestEditorFn{
		func(ctx context.Context, req *http.Request) error {
			client.SetReferer(ctx, req)
			return client.SetSigV1(req)
		},
	}

	return nil
}

// AuthWithResponse returns a ClientWithResponses with RequestEditorFn that generates signatures.
func AuthWithResponse() (*authapi.ClientWithResponses, error) {
	authClientWithResponse, err := authapi.NewClientWithResponses(server, withRequestEditorFns)
	if err != nil {
		return nil, err
	}

	return authClientWithResponse, nil
}
