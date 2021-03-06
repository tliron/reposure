package direct

import (
	"fmt"

	namepkg "github.com/google/go-containerregistry/pkg/name"
)

func (self *Client) getName(imageName string) string {
	return fmt.Sprintf("%s/%s", self.Host, imageName)
}

func (self *Client) newRegistry(registryHost string) (namepkg.Registry, error) {
	if self.Secure {
		return namepkg.NewRegistry(registryHost)
	} else {
		return namepkg.NewRegistry(registryHost, namepkg.Insecure)
	}
}
