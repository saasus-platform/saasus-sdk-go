package auth

import (
	"github.com/saasus-platform/saasus-sdk-go/client"
	"github.com/saasus-platform/saasus-sdk-go/generated/authapi"
)

var (
	server = "https://api.saasus.io/v1/auth"
)

func withRequestEditorFns(c *authapi.Client) error {
	c.RequestEditors = []authapi.RequestEditorFn{
		client.WithSaasusSigV1(),
	}

	return nil
}

func AuthWithResponse() (*authapi.ClientWithResponses, error) {
	authClientWithResponse, err := authapi.NewClientWithResponses(server, withRequestEditorFns)
	if err != nil {
		return nil, err
	}

	return authClientWithResponse, nil
}
