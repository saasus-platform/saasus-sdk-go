package communication

import (
	"context"
	"net/http"

	"github.com/saasus-platform/saasus-sdk-go/client"
	"github.com/saasus-platform/saasus-sdk-go/generated/communicationapi"
)

var (
	server = "https://api.saasus.io/v1/communication"
)

func withRequestEditorFns(c *communicationapi.Client) error {
	c.RequestEditors = []communicationapi.RequestEditorFn{
		func(ctx context.Context, req *http.Request) error {
			client.SetReferer(ctx, req)
			return client.SetSigV1(req)
		},
	}

	return nil
}

// CommunicationWithResponse returns a ClientWithResponses with RequestEditorFn that generates signatures.
func CommunicationWithResponse() (*communicationapi.ClientWithResponses, error) {
	clientWithResponse, err := communicationapi.NewClientWithResponses(server, withRequestEditorFns)
	if err != nil {
		return nil, err
	}

	return clientWithResponse, nil
}
