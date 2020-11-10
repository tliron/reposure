package repository

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

func (self *Client) GetCertificatePath(repository *resources.Repository) string {
	if repository.Spec.TLSSecret != "" {
		secretDataKey := repository.Spec.TLSSecretDataKey
		if secretDataKey == "" {
			secretDataKey = core.TLSCertKey
		}
		return fmt.Sprintf("%s/%s", self.TLSMountPath, secretDataKey)
	} else {
		return ""
	}
}

func (self *Client) GetTLSCertPool(repository *resources.Repository) (*x509.CertPool, error) {
	if repository.Spec.TLSSecret != "" {
		if secret, err := self.Kubernetes.CoreV1().Secrets(repository.Namespace).Get(self.Context, repository.Spec.TLSSecret, meta.GetOptions{}); err == nil {
			return kubernetesutil.GetSecretTLSCertPool(secret, repository.Spec.TLSSecretDataKey)
		} else {
			return nil, err
		}
	} else {
		return nil, nil
	}
}

func (self *Client) GetHTTPRoundTripper(repository *resources.Repository) (string, http.RoundTripper, error) {
	if certPool, err := self.GetTLSCertPool(repository); err == nil {
		if certPool != nil {
			if host, err := self.GetHost(repository); err == nil {
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
