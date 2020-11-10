package admin

import (
	"fmt"

	core "k8s.io/api/core/v1"
)

const tlsMountPath = "/tls"

// TODO: these depend on the dataKey!
var tlsCertificatePath = fmt.Sprintf("%s/%s", tlsMountPath, core.TLSCertKey)
var tlsKeyPath = fmt.Sprintf("%s/%s", tlsMountPath, core.TLSPrivateKeyKey)
