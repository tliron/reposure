package admin

import (
	contextpkg "context"

	certmanagerpkg "github.com/cert-manager/cert-manager/pkg/client/clientset/versioned"
	"github.com/tliron/commonlog"
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

	ClusterMode             bool
	ClusterRole             string
	Namespace               string
	NamePrefix              string
	PartOf                  string
	ManagedBy               string
	OperatorImageReference  string
	SurrogateImageReference string
	SimpleImageReference    string

	LogName string
	Log     commonlog.Logger
}

func NewClient(kubernetes kubernetespkg.Interface, apiExtensions apiextensionspkg.Interface, reposure reposurepkg.Interface, rest restpkg.Interface, config *restpkg.Config, context contextpkg.Context, clusterMode bool, clusterRole string, namespace string, namePrefix string, partOf string, managedBy string, operatorImageReference string, surrogateImageReference string, simpleImageReference string, logName string) *Client {
	return &Client{
		Kubernetes:              kubernetes,
		APIExtensions:           apiExtensions,
		Reposure:                reposure,
		REST:                    rest,
		Config:                  config,
		Context:                 context,
		ClusterMode:             clusterMode,
		ClusterRole:             clusterRole,
		Namespace:               namespace,
		NamePrefix:              namePrefix,
		PartOf:                  partOf,
		ManagedBy:               managedBy,
		OperatorImageReference:  operatorImageReference,
		SurrogateImageReference: surrogateImageReference,
		SimpleImageReference:    simpleImageReference,
		LogName:                 logName,
		Log:                     commonlog.GetLoggerf("%s.admin", logName),
	}
}
