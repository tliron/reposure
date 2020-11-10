package registry

import (
	namepkg "github.com/google/go-containerregistry/pkg/name"
)

func (self *Client) newRegistry(registryHost string) (namepkg.Registry, error) {
	if self.Secure {
		return namepkg.NewRegistry(registryHost)
	} else {
		return namepkg.NewRegistry(registryHost, namepkg.Insecure)
	}
}
