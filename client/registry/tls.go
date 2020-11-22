package registry

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"

	kubernetesutil "github.com/tliron/kutil/kubernetes"
	"github.com/tliron/kutil/util"
	resources "github.com/tliron/reposure/resources/reposure.puccini.cloud/v1alpha1"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (self *Client) GetCertificatePath(registry *resources.Registry) string {
	if registry.Spec.AuthenticationSecret != "" {
		secretDataKey := registry.Spec.AuthenticationSecretDataKey
		if secretDataKey == "" {
			secretDataKey = core.TLSCertKey
		}
		return fmt.Sprintf("%s/%s", self.TLSMountPath, secretDataKey)
	} else {
		return ""
	}
}

func (self *Client) GetTLSCertBytes(registry *resources.Registry) ([]byte, error) {
	if registry.Spec.AuthenticationSecret != "" {
		if secret, err := self.Kubernetes.CoreV1().Secrets(registry.Namespace).Get(self.Context, registry.Spec.AuthenticationSecret, meta.GetOptions{}); err == nil {
			return kubernetesutil.GetSecretTLSCertBytes(secret, registry.Spec.AuthenticationSecretDataKey)

		} else {
			return nil, err
		}
	} else {
		return nil, nil
	}
}

func (self *Client) GetTLSCertPool(registry *resources.Registry) (*x509.CertPool, error) {
	if bytes, err := self.GetTLSCertBytes(registry); err == nil {
		return util.ParseX509CertPool(bytes)
	} else {
		return nil, err
	}
}

func (self *Client) GetHTTPRoundTripper(registry *resources.Registry) (string, http.RoundTripper, error) {
	if certPool, err := self.GetTLSCertPool(registry); err == nil {
		if certPool != nil {
			if host, err := self.GetHost(registry); err == nil {
				roundTripper := util.NewForceHTTPSRoundTripper(&http.Transport{
					TLSClientConfig: &tls.Config{
						RootCAs: certPool,
					},
				})
				return host, roundTripper, nil
			} else {
				return "", nil, err
			}
		} else {
			return "", nil, nil
		}
	} else {
		return "", nil, err
	}
}
