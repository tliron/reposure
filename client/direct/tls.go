package direct

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"

	"github.com/tliron/kutil/util"
)

func TLSRoundTripper(certificatePath string) (http.RoundTripper, error) {
	if certPool, err := CertPool(certificatePath); err == nil {
		if certPool != nil {
			// We need to force HTTPS because go-containerregistry will attempt to drop down to HTTP for local addresses
			return util.NewForceHTTPSRoundTripper(&http.Transport{
				TLSClientConfig: &tls.Config{
					RootCAs: certPool,
				},
			}), nil
		} else {
			return nil, nil
		}
	} else {
		return nil, err
	}
}

func CertPool(certificatePath string) (*x509.CertPool, error) {
	if bytes, err := ioutil.ReadFile(certificatePath); err == nil {
		return util.ParseX509CertPool(bytes)
	} else {
		return nil, err
	}
}
