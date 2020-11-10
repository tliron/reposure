package repository

import (
	"fmt"

	resources "github.com/tliron/reposure/resources/reposure.puccini.cloud/v1alpha1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (self *Client) GetHost(repository *resources.Repository) (string, error) {
	if (repository.Spec.Direct != nil) && (repository.Spec.Direct.Host != "") {
		return repository.Spec.Direct.Host, nil
	} else if (repository.Spec.Indirect != nil) && (repository.Spec.Indirect.Service != "") {
		serviceNamespace := repository.Spec.Indirect.Namespace
		if serviceNamespace == "" {
			// Default to repository namespace
			serviceNamespace = repository.Namespace
		}

		if service, err := self.Kubernetes.CoreV1().Services(serviceNamespace).Get(self.Context, repository.Spec.Indirect.Service, meta.GetOptions{}); err == nil {
			return fmt.Sprintf("%s:%d", service.Spec.ClusterIP, repository.Spec.Indirect.Port), nil
		} else {
			return "", err
		}
	} else {
		return "", fmt.Errorf("malformed repository: %s", repository.Name)
	}
}
