package apilog

import (
	"context"
	"net/http"

	"github.com/saasus-platform/saasus-sdk-go/client"
	"github.com/saasus-platform/saasus-sdk-go/generated/apilogapi"
)

var (
	server = "https://api.saasus.io/v1/apilog"
)

func withRequestEditorFns(c *apilogapi.Client) error {
	c.RequestEditors = []apilogapi.RequestEditorFn{
		func(ctx context.Context, req *http.Request) error {
			client.SetReferer(ctx, req)
			return client.SetSigV1(req)
		},
	}

	return nil
}

// ApiLogWithResponse returns a ClientWithResponses with RequestEditorFn that generates signatures.
func ApiLogWithResponse() (*apilogapi.ClientWithResponses, error) {
	clientWithResponse, err := apilogapi.NewClientWithResponses(server, withRequestEditorFns)
	if err != nil {
		return nil, err
	}

	return clientWithResponse, nil
}
