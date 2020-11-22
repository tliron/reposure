package registry

import (
	"fmt"

	resources "github.com/tliron/reposure/resources/reposure.puccini.cloud/v1alpha1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (self *Client) GetHost(registry *resources.Registry) (string, error) {
	if (registry.Spec.Direct != nil) && (registry.Spec.Direct.Host != "") {
		return registry.Spec.Direct.Host, nil
	} else if (registry.Spec.Indirect != nil) && (registry.Spec.Indirect.Service != "") {
		serviceNamespace := registry.Spec.Indirect.Namespace
		if serviceNamespace == "" {
			// Default to registry namespace
			serviceNamespace = registry.Namespace
		}

		if service, err := self.Kubernetes.CoreV1().Services(serviceNamespace).Get(self.Context, registry.Spec.Indirect.Service, meta.GetOptions{}); err == nil {
			return fmt.Sprintf("%s:%d", service.Spec.ClusterIP, registry.Spec.Indirect.Port), nil
		} else {
			return "", err
		}
	} else {
		return "", fmt.Errorf("malformed registry: %s", registry.Name)
	}
}
