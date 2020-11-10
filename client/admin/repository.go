package admin

import (
	resources "github.com/tliron/reposure/resources/reposure.puccini.cloud/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (self *Client) GetRepository(namespace string, repositoryName string) (*resources.Repository, error) {
	// Default to same namespace as operator
	if namespace == "" {
		namespace = self.Namespace
	}

	if repository, err := self.Reposure.ReposureV1alpha1().Repositories(namespace).Get(self.Context, repositoryName, meta.GetOptions{}); err == nil {
		// When retrieved from cache the GVK may be empty
		if repository.Kind == "" {
			repository = repository.DeepCopy()
			repository.APIVersion, repository.Kind = resources.RepositoryGVK.ToAPIVersionAndKind()
		}
		return repository, nil
	} else {
		return nil, err
	}
}

func (self *Client) ListRepositories() (*resources.RepositoryList, error) {
	// TODO: all repositories in cluster mode
	return self.Reposure.ReposureV1alpha1().Repositories(self.Namespace).List(self.Context, meta.ListOptions{})
}

func (self *Client) CreateRepositoryDirect(namespace string, repositoryName string, host string, tlsSecretName string, tlsSecretDataKey string, authSecretName string) (*resources.Repository, error) {
	// Default to same namespace as operator
	if namespace == "" {
		namespace = self.Namespace
	}

	repository := &resources.Repository{
		ObjectMeta: meta.ObjectMeta{
			Name:      repositoryName,
			Namespace: namespace,
		},
		Spec: resources.RepositorySpec{
			Type: resources.RepositoryTypeRegistry,
			Direct: &resources.RepositoryDirect{
				Host: host,
			},
			TLSSecret:        tlsSecretName,
			TLSSecretDataKey: tlsSecretDataKey,
			AuthSecret:       authSecretName,
		},
	}

	return self.createRepository(namespace, repositoryName, repository)
}

func (self *Client) CreateRepositoryIndirect(namespace string, repositoryName string, serviceNamespace string, serviceName string, port uint64, tlsSecretName string, tlsSecretDataKey string, authSecretName string) (*resources.Repository, error) {
	// Default to same namespace as operator
	if namespace == "" {
		namespace = self.Namespace
	}

	repository := &resources.Repository{
		ObjectMeta: meta.ObjectMeta{
			Name:      repositoryName,
			Namespace: namespace,
		},
		Spec: resources.RepositorySpec{
			Type: resources.RepositoryTypeRegistry,
			Indirect: &resources.RepositoryIndirect{
				Namespace: serviceNamespace,
				Service:   serviceName,
				Port:      port,
			},
			TLSSecret:        tlsSecretName,
			TLSSecretDataKey: tlsSecretDataKey,
			AuthSecret:       authSecretName,
		},
	}

	return self.createRepository(namespace, repositoryName, repository)
}

func (self *Client) createRepository(namespace string, repositoryName string, repository *resources.Repository) (*resources.Repository, error) {
	if repository, err := self.Reposure.ReposureV1alpha1().Repositories(namespace).Create(self.Context, repository, meta.CreateOptions{}); err == nil {
		return repository, nil
	} else if errors.IsAlreadyExists(err) {
		return self.Reposure.ReposureV1alpha1().Repositories(namespace).Get(self.Context, repositoryName, meta.GetOptions{})
	} else {
		return nil, err
	}
}

func (self *Client) UpdateRepositoryStatus(repository *resources.Repository) (*resources.Repository, error) {
	if repository_, err := self.Reposure.ReposureV1alpha1().Repositories(repository.Namespace).UpdateStatus(self.Context, repository, meta.UpdateOptions{}); err == nil {
		// When retrieved from cache the GVK may be empty
		if repository_.Kind == "" {
			repository_ = repository_.DeepCopy()
			repository_.APIVersion, repository_.Kind = resources.RepositoryGVK.ToAPIVersionAndKind()
		}
		return repository_, nil
	} else {
		return repository, err
	}
}

func (self *Client) DeleteRepository(namespace string, repositoryName string) error {
	// Default to same namespace as operator
	if namespace == "" {
		namespace = self.Namespace
	}

	return self.Reposure.ReposureV1alpha1().Repositories(namespace).Delete(self.Context, repositoryName, meta.DeleteOptions{})
}
