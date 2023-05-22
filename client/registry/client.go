package registry

import (
	contextpkg "context"

	"github.com/tliron/commonlog"
	reposurepkg "github.com/tliron/reposure/apis/clientset/versioned"
	directclient "github.com/tliron/reposure/client/direct"
	resources "github.com/tliron/reposure/resources/reposure.puccini.cloud/v1alpha1"
	kubernetespkg "k8s.io/client-go/kubernetes"
)

//
// Client
//

type Client struct {
	Kubernetes kubernetespkg.Interface
	Reposure   reposurepkg.Interface
	Context    contextpkg.Context
	Log        commonlog.Logger

	Namespace    string
	TLSMountPath string
}

func NewClient(kubernetes kubernetespkg.Interface, reposure reposurepkg.Interface, context contextpkg.Context, log commonlog.Logger, namespace string, tlsMountPath string) *Client {
	return &Client{
		Kubernetes:   kubernetes,
		Reposure:     reposure,
		Context:      context,
		Log:          log,
		Namespace:    namespace,
		TLSMountPath: tlsMountPath,
	}
}

func (self *Client) DirectClient(registry *resources.Registry) (*directclient.Client, error) {
	if host, transport, err := self.GetHTTPRoundTripper(registry); err == nil {
		if _, username, password, token, err := self.GetAuthorization(registry); err == nil {
			return directclient.NewClient(
				host,
				transport,
				username,
				password,
				token,
				self.Context,
			), nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}
