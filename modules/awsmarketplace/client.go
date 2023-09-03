package awsmarketplace

import (
	"context"
	"net/http"

	"github.com/saasus-platform/saasus-sdk-go/client"
	"github.com/saasus-platform/saasus-sdk-go/generated/awsmarketplaceapi"
)

var (
	server = "https://api.saasus.io/v1/awsmarketplace"
)

func withRequestEditorFns(c *awsmarketplaceapi.Client) error {
	c.RequestEditors = []awsmarketplaceapi.RequestEditorFn{
		func(ctx context.Context, req *http.Request) error {
			client.SetReferer(ctx, req)
			return client.SetSigV1(req)
		},
	}

	return nil
}

// AwsMarketplaceWithResponse returns a ClientWithResponses with RequestEditorFn that generates signatures.
func AwsMarketplaceWithResponse() (*awsmarketplaceapi.ClientWithResponses, error) {
	clientWithResponse, err := awsmarketplaceapi.NewClientWithResponses(server, withRequestEditorFns)
	if err != nil {
		return nil, err
	}

	return clientWithResponse, nil
}
