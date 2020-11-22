package admin

import (
	resources "github.com/tliron/reposure/resources/reposure.puccini.cloud/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (self *Client) GetRegistry(namespace string, registryName string) (*resources.Registry, error) {
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

func (self *Client) ListRegistries() (*resources.RegistryList, error) {
	// TODO: all registries in cluster mode
	return self.Reposure.ReposureV1alpha1().Registries(self.Namespace).List(self.Context, meta.ListOptions{})
}

func (self *Client) CreateRegistryDirect(namespace string, registryName string, host string, authenticationSecretName string, authenticationSecretDataKey string, authorizationSecretName string) (*resources.Registry, error) {
	// Default to same namespace as operator
	if namespace == "" {
		namespace = self.Namespace
	}

	registry := &resources.Registry{
		ObjectMeta: meta.ObjectMeta{
			Name:      registryName,
			Namespace: namespace,
		},
		Spec: resources.RegistrySpec{
			Type: resources.RegistryTypeOCI,
			Direct: &resources.RegistryDirect{
				Host: host,
			},
			AuthenticationSecret:        authenticationSecretName,
			AuthenticationSecretDataKey: authenticationSecretDataKey,
			AuthorizationSecret:         authorizationSecretName,
		},
	}

	return self.createRegistry(namespace, registryName, registry)
}

func (self *Client) CreateRegistryIndirect(namespace string, registryName string, serviceNamespace string, serviceName string, port uint64, authenticationSecretName string, authenticationSecretDataKey string, authorizationSecretName string) (*resources.Registry, error) {
	// Default to same namespace as operator
	if namespace == "" {
		namespace = self.Namespace
	}

	registry := &resources.Registry{
		ObjectMeta: meta.ObjectMeta{
			Name:      registryName,
			Namespace: namespace,
		},
		Spec: resources.RegistrySpec{
			Type: resources.RegistryTypeOCI,
			Indirect: &resources.RegistryIndirect{
				Namespace: serviceNamespace,
				Service:   serviceName,
				Port:      port,
			},
			AuthenticationSecret:        authenticationSecretName,
			AuthenticationSecretDataKey: authenticationSecretDataKey,
			AuthorizationSecret:         authorizationSecretName,
		},
	}

	return self.createRegistry(namespace, registryName, registry)
}

func (self *Client) createRegistry(namespace string, registryName string, registry *resources.Registry) (*resources.Registry, error) {
	if registry, err := self.Reposure.ReposureV1alpha1().Registries(namespace).Create(self.Context, registry, meta.CreateOptions{}); err == nil {
		return registry, nil
	} else if errors.IsAlreadyExists(err) {
		return self.Reposure.ReposureV1alpha1().Registries(namespace).Get(self.Context, registryName, meta.GetOptions{})
	} else {
		return nil, err
	}
}

func (self *Client) UpdateRegistryStatus(registry *resources.Registry) (*resources.Registry, error) {
	if registry_, err := self.Reposure.ReposureV1alpha1().Registries(registry.Namespace).UpdateStatus(self.Context, registry, meta.UpdateOptions{}); err == nil {
		// When retrieved from cache the GVK may be empty
		if registry_.Kind == "" {
			registry_ = registry_.DeepCopy()
			registry_.APIVersion, registry_.Kind = resources.RegistryGVK.ToAPIVersionAndKind()
		}
		return registry_, nil
	} else {
		return registry, err
	}
}

func (self *Client) DeleteRegistry(namespace string, registryName string) error {
	// Default to same namespace as operator
	if namespace == "" {
		namespace = self.Namespace
	}

	return self.Reposure.ReposureV1alpha1().Registries(namespace).Delete(self.Context, registryName, meta.DeleteOptions{})
}
