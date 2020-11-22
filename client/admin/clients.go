package admin

import (
	"fmt"

	"github.com/op/go-logging"
	commandclient "github.com/tliron/reposure/client/command"
	registryclient "github.com/tliron/reposure/client/registry"
	spoolerclient "github.com/tliron/reposure/client/spooler"
	resources "github.com/tliron/reposure/resources/reposure.puccini.cloud/v1alpha1"
)

func (self *Client) RegistryClient() *registryclient.Client {
	return registryclient.NewClient(
		self.Kubernetes,
		self.Reposure,
		self.Context,
		logging.MustGetLogger(fmt.Sprintf("%s.registry", self.LogName)),
		self.Namespace,
		tlsMountPath,
	)
}

func (self *Client) SpoolerClient(registry *resources.Registry) *spoolerclient.Client {
	appName := self.GetRegistrySurrogateAppName(registry.Name)

	return spoolerclient.NewClient(
		self.Kubernetes,
		self.REST,
		self.Config,
		self.Context,
		nil,
		self.Namespace,
		appName,
		surrogateContainerName,
		spoolPath,
	)
}

func (self *Client) CommandClient(registry *resources.Registry) (*commandclient.Client, error) {
	appName := self.GetRegistrySurrogateAppName(registry.Name)
	registryClient := self.RegistryClient()

	if _, username, password, token, err := registryClient.GetAuthorization(registry); err == nil {
		if address, err := registryClient.GetHost(registry); err == nil {
			return commandclient.NewClient(
				self.Kubernetes,
				self.REST,
				self.Config,
				self.Context,
				nil,
				self.Namespace,
				appName,
				surrogateContainerName,
				address,
				registryClient.GetCertificatePath(registry),
				username,
				password,
				token,
			), nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}
