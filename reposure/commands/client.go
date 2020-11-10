package commands

import (
	contextpkg "context"

	kubernetesutil "github.com/tliron/kutil/kubernetes"
	"github.com/tliron/kutil/util"
	reposurepkg "github.com/tliron/reposure/apis/clientset/versioned"
	adminclient "github.com/tliron/reposure/client/admin"
	"github.com/tliron/reposure/controller"
	apiextensionspkg "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	kubernetespkg "k8s.io/client-go/kubernetes"
	restpkg "k8s.io/client-go/rest"
)

//
// Client
//

type Client struct {
	Config     *restpkg.Config
	Kubernetes kubernetespkg.Interface
	REST       restpkg.Interface
	Context    contextpkg.Context
	Namespace  string
}

func NewClient() *Client {
	config, err := kubernetesutil.NewConfigFromFlags(masterUrl, kubeconfigPath, context, log)
	util.FailOnError(err)

	kubernetes, err := kubernetespkg.NewForConfig(config)
	util.FailOnError(err)

	namespace_ := namespace
	if cluster {
		namespace_ = ""
	} else if namespace_ == "" {
		if namespace__, ok := kubernetesutil.GetConfiguredNamespace(kubeconfigPath, context); ok {
			namespace_ = namespace__
		}
		if namespace_ == "" {
			util.Fail("could not discover namespace and \"--namespace\" not provided")
		}
	}

	return &Client{
		Config:     config,
		Kubernetes: kubernetes,
		REST:       kubernetes.CoreV1().RESTClient(),
		Context:    contextpkg.TODO(),
		Namespace:  namespace_,
	}
}

func (self *Client) AdminClient() *adminclient.Client {
	apiExtensions, err := apiextensionspkg.NewForConfig(self.Config)
	util.FailOnError(err)

	reposure, err := reposurepkg.NewForConfig(self.Config)
	util.FailOnError(err)

	return adminclient.NewClient(
		self.Kubernetes,
		apiExtensions,
		reposure,
		self.REST,
		self.Config,
		self.Context,
		cluster,
		self.Namespace,
		controller.NamePrefix,
		controller.PartOf,
		controller.ManagedBy,
		controller.OperatorImageReference,
		controller.RepositorySurrogateImageReference,
		controller.RegistryImageReference,
		"reposure.client",
	)
}
