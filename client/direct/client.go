package direct

import (
	contextpkg "context"
	"net/http"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

//
// Client
//

type Client struct {
	Secure  bool
	Options []remote.Option
	Context contextpkg.Context
}

func NewClient(context contextpkg.Context, transport http.RoundTripper, username string, password string, token string) *Client {
	var options []remote.Option

	if transport != nil {
		options = append(options, remote.WithTransport(transport))
	}

	if (username != "") || (token != "") {
		authenticator := authn.FromConfig(authn.AuthConfig{
			Username:      username,
			Password:      password,
			RegistryToken: token,
		})
		options = append(options, remote.WithAuth(authenticator))
	}

	return &Client{
		Secure:  transport != nil,
		Options: options,
		Context: context,
	}
}
