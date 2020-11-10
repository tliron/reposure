package repository

import (
	contextpkg "context"

	"github.com/op/go-logging"
	reposurepkg "github.com/tliron/reposure/apis/clientset/versioned"
	resources "github.com/tliron/reposure/resources/reposure.puccini.cloud/v1alpha1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubernetespkg "k8s.io/client-go/kubernetes"
)

//
// Client
//

type Client struct {
	Kubernetes kubernetespkg.Interface
	Reposure   reposurepkg.Interface
	Context    contextpkg.Context
	Log        *logging.Logger

	Namespace    string
	TLSMountPath string
}

func NewClient(kubernetes kubernetespkg.Interface, reposure reposurepkg.Interface, context contextpkg.Context, log *logging.Logger, namespace string, tlsMountPath string) *Client {
	return &Client{
		Kubernetes:   kubernetes,
		Reposure:     reposure,
		Context:      context,
		Log:          log,
		Namespace:    namespace,
		TLSMountPath: tlsMountPath,
	}
}

func (self *Client) Get(namespace string, repositoryName string) (*resources.Repository, error) {
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
