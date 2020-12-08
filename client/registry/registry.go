package registry

import (
	resources "github.com/tliron/reposure/resources/reposure.puccini.cloud/v1alpha1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (self *Client) Get(namespace string, registryName string) (*resources.Registry, error) {
	// Default to same namespace as operator
	if namespace == "" {
		namespace = self.Namespace
	}

	if registry, err := self.Reposure.ReposureV1alpha1().Registries(namespace).Get(self.Context, registryName, meta.GetOptions{}); err == nil {
		// When retrieved from cache the GVK may be empty
		if registry.Kind == "" {
			registry = registry.DeepCopy()
			registry.APIVersion, registry.Kind = resources.RegistryGVK.ToAPIVersionAndKind()
		}
		return registry, nil
	} else {
		return nil, err
	}
}
