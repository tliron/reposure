package admin

import (
	"github.com/tliron/kutil/logging"
	directclient "github.com/tliron/reposure/client/direct"
	registryclient "github.com/tliron/reposure/client/registry"
	commandclient "github.com/tliron/reposure/client/surrogate/command"
	spoolerclient "github.com/tliron/reposure/client/surrogate/spooler"
	resources "github.com/tliron/reposure/resources/reposure.puccini.cloud/v1alpha1"
)

func (self *Client) RegistryClient() *registryclient.Client {
	return registryclient.NewClient(
		self.Kubernetes,
		self.Reposure,
		self.Context,
		logging.GetLoggerf("%s.registry", self.LogName),
		self.Namespace,
		tlsMountPath,
	)
}

func (self *Client) DirectClient(registry *resources.Registry) (*directclient.Client, error) {
	return self.RegistryClient().DirectClient(registry)
}

func (self *Client) SurrogateSpoolerClient(registry *resources.Registry) *spoolerclient.Client {
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

func (self *Client) SurrogateCommandClient(registry *resources.Registry) (*commandclient.Client, error) {
	appName := self.GetRegistrySurrogateAppName(registry.Name)
	registryClient := self.RegistryClient()

	if _, username, password, token, err := registryClient.GetAuthorization(registry); err == nil {
		if host, err := registryClient.GetHost(registry); err == nil {
			return commandclient.NewClient(
				self.Kubernetes,
				self.REST,
				self.Config,
				self.Context,
				nil,
				self.Namespace,
				appName,
				surrogateContainerName,
				host,
				registryClient.GetCertificatePath(registry),
				username,
				password,
				token,
				logging.GetLoggerf("%s.command", self.LogName),
			), nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}
