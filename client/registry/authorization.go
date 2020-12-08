package registry

import (
	"fmt"

	kubernetesutil "github.com/tliron/kutil/kubernetes"
	"github.com/tliron/kutil/util"
	resources "github.com/tliron/reposure/resources/reposure.puccini.cloud/v1alpha1"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (self *Client) GetAuthorization(registry *resources.Registry) (string, string, string, string, error) {
	if host, err := self.GetHost(registry); err == nil {
		if registry.Spec.AuthorizationSecret != "" {
			if secret, err := self.Kubernetes.CoreV1().Secrets(self.Namespace).Get(self.Context, registry.Spec.AuthorizationSecret, meta.GetOptions{}); err == nil {
				switch secret.Type {
				case core.SecretTypeServiceAccountToken:
					if data, ok := secret.Data[core.ServiceAccountTokenKey]; ok {
						// OpenShift: you can also get a valid token from "oc whoami -t"
						token := util.BytesToString(data)
						return host, "", "", token, nil
					} else {
						return "", "", "", "", fmt.Errorf("malformed %q secret: %s", core.ServiceAccountTokenKey, secret.Data)
					}

				case core.SecretTypeDockerConfigJson, core.SecretTypeDockercfg:
					if table, err := kubernetesutil.NewRegistryCredentialsTableFromSecret(secret); err == nil {
						if credentials, ok := table[host]; ok {
							return host, credentials.Username, credentials.Password, "", nil
						}
					} else {
						return "", "", "", "", err
					}

				default:
					return "", "", "", "", fmt.Errorf("unsupposed authorization secret type: %s", secret.Type)
				}
			} else {
				return "", "", "", "", err
			}
		}

		return host, "", "", "", nil
	} else {
		return "", "", "", "", err
	}
}
