package admin

import (
	"fmt"

	"github.com/op/go-logging"
	commandclient "github.com/tliron/reposure/client/command"
	repositoryclient "github.com/tliron/reposure/client/repository"
	spoolerclient "github.com/tliron/reposure/client/spooler"
	resources "github.com/tliron/reposure/resources/reposure.puccini.cloud/v1alpha1"
)

func (self *Client) RepositoryClient(repository *resources.Repository) *repositoryclient.Client {
	return repositoryclient.NewClient(
		self.Kubernetes,
		self.Reposure,
		self.Context,
		logging.MustGetLogger(fmt.Sprintf("%s.repository", self.LogName)),
		self.Namespace,
		tlsMountPath,
	)
}

func (self *Client) SpoolerClient(repository *resources.Repository) *spoolerclient.Client {
	appName := self.GetRepositorySurrogateAppName(repository.Name)

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

func (self *Client) CommandClient(repository *resources.Repository) (*commandclient.Client, error) {
	appName := self.GetRepositorySurrogateAppName(repository.Name)

	repositoryClient := self.RepositoryClient(repository)

	if _, username, password, token, err := repositoryClient.GetAuth(repository); err == nil {
		if address, err := repositoryClient.GetHost(repository); err == nil {
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
				repositoryClient.GetCertificatePath(repository),
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
