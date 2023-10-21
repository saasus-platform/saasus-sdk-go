package pricing

import (
	"context"
	"net/http"

	"github.com/saasus-platform/saasus-sdk-go/client"
	"github.com/saasus-platform/saasus-sdk-go/generated/pricingapi"
)

var (
	server = "https://api.saasus.io/v1/pricing"
)

func withRequestEditorFns(c *pricingapi.Client) error {
	c.RequestEditors = []pricingapi.RequestEditorFn{
		func(ctx context.Context, req *http.Request) error {
			client.SetReferer(ctx, req)
			return client.SetSigV1(req)
		}}

	return nil
}

// PricingWithResponse returns a ClientWithResponses with RequestEditorFn that generates signatures.
func PricingWithResponse() (*pricingapi.ClientWithResponses, error) {
	clientWithResponse, err := pricingapi.NewClientWithResponses(server, withRequestEditorFns)
	if err != nil {
		return nil, err
	}

	return clientWithResponse, nil
}
