package integration

import (
	"context"
	"net/http"

	"github.com/saasus-platform/saasus-sdk-go/client"
	"github.com/saasus-platform/saasus-sdk-go/generated/integrationapi"
)

var (
	server = "https://api.saasus.io/v1/integration"
)

func withRequestEditorFns(c *integrationapi.Client) error {
	c.RequestEditors = []integrationapi.RequestEditorFn{
		func(ctx context.Context, req *http.Request) error {
			client.SetReferer(ctx, req)
			return client.SetSigV1(req)
		}}

	return nil
}

// IntegrationWithResponse returns a ClientWithResponses with RequestEditorFn that generates signatures.
func IntegrationWithResponse() (*integrationapi.ClientWithResponses, error) {
	clientWithResponse, err := integrationapi.NewClientWithResponses(server, withRequestEditorFns)
	if err != nil {
		return nil, err
	}

	return clientWithResponse, nil
}
