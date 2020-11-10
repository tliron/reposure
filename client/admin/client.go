package admin

import (
	contextpkg "context"
	"fmt"

	certmanagerpkg "github.com/jetstack/cert-manager/pkg/client/clientset/versioned"
	"github.com/op/go-logging"
	reposurepkg "github.com/tliron/reposure/apis/clientset/versioned"
	apiextensionspkg "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	kubernetespkg "k8s.io/client-go/kubernetes"
	restpkg "k8s.io/client-go/rest"
)

//
// Client
//

type Client struct {
	Kubernetes    kubernetespkg.Interface
	APIExtensions apiextensionspkg.Interface
	Reposure      reposurepkg.Interface
	REST          restpkg.Interface
	CertManager   certmanagerpkg.Interface
	Config        *restpkg.Config
	Context       contextpkg.Context

	Cluster                           bool
	Namespace                         string
	NamePrefix                        string
	PartOf                            string
	ManagedBy                         string
	OperatorImageReference            string
	RepositorySurrogateImageReference string
	RegistryImageReference            string

	LogName string
	Log     *logging.Logger
}

func NewClient(kubernetes kubernetespkg.Interface, apiExtensions apiextensionspkg.Interface, reposure reposurepkg.Interface, rest restpkg.Interface, config *restpkg.Config, context contextpkg.Context, cluster bool, namespace string, namePrefix string, partOf string, managedBy string, operatorImageReference string, repositorySurrogateImageReference string, registryImageReference string, logName string) *Client {
	return &Client{
		Kubernetes:                        kubernetes,
		APIExtensions:                     apiExtensions,
		Reposure:                          reposure,
		REST:                              rest,
		Config:                            config,
		Context:                           context,
		Cluster:                           cluster,
		Namespace:                         namespace,
		NamePrefix:                        namePrefix,
		PartOf:                            partOf,
		ManagedBy:                         managedBy,
		OperatorImageReference:            operatorImageReference,
		RepositorySurrogateImageReference: repositorySurrogateImageReference,
		RegistryImageReference:            registryImageReference,
		LogName:                           logName,
		Log:                               logging.MustGetLogger(fmt.Sprintf("%s.admin", logName)),
	}
}
