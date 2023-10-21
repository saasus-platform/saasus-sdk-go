package billing

import (
	"context"
	"net/http"

	"github.com/saasus-platform/saasus-sdk-go/client"
	"github.com/saasus-platform/saasus-sdk-go/generated/billingapi"
)

var (
	server = "https://api.saasus.io/v1/billing"
)

func withRequestEditorFns(c *billingapi.Client) error {
	c.RequestEditors = []billingapi.RequestEditorFn{
		func(ctx context.Context, req *http.Request) error {
			client.SetReferer(ctx, req)
			return client.SetSigV1(req)
		},
	}

	return nil
}

// BillingWithResponse returns a ClientWithResponses with RequestEditorFn that generates signatures.
func BillingWithResponse() (*billingapi.ClientWithResponses, error) {
	clientWithResponse, err := billingapi.NewClientWithResponses(server, withRequestEditorFns)
	if err != nil {
		return nil, err
	}

	return clientWithResponse, nil
}
