package repository

import (
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	resources "github.com/tliron/reposure/resources/reposure.puccini.cloud/v1alpha1"
)

func (self *Client) GetRemoteOptions(repository *resources.Repository) ([]remote.Option, error) {
	var options []remote.Option

	if _, roundTripper, err := self.GetHTTPRoundTripper(repository); err == nil {
		if roundTripper != nil {
			options = append(options, remote.WithTransport(roundTripper))
		}
	} else {
		return nil, err
	}

	if _, username, password, token, err := self.GetAuth(repository); err == nil {
		if (username != "") || (token != "") {
			authenticator := authn.FromConfig(authn.AuthConfig{
				Username:      username,
				Password:      password,
				RegistryToken: token,
			})
			options = append(options, remote.WithAuth(authenticator))
		}
	} else {
		return nil, err
	}

	return options, nil
}
