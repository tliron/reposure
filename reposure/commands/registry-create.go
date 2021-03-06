package commands

import (
	contextpkg "context"

	"github.com/spf13/cobra"
	"github.com/tliron/kutil/util"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var host string
var serviceNamespace string
var service string
var port uint64
var authenticationSecret string
var authenticationSecretDataKey string
var authorizationSecret string
var provider string

func init() {
	registryCommand.AddCommand(registryCreateCommand)
	registryCreateCommand.Flags().StringVarP(&host, "host", "", "", "registry host (\"host\" or \"host:port\")")
	registryCreateCommand.Flags().StringVarP(&serviceNamespace, "service-namespace", "", "", "registry service namespace name (defaults to registry namespace)")
	registryCreateCommand.Flags().StringVarP(&service, "service", "", "", "registry service name")
	registryCreateCommand.Flags().Uint64VarP(&port, "port", "", 5000, "registry service port")
	registryCreateCommand.Flags().StringVarP(&authenticationSecret, "authentication-secret", "", "", "registry authentication secret name")
	registryCreateCommand.Flags().StringVarP(&authenticationSecretDataKey, "authentication-secret-data-key", "", "", "registry authentication secret data key name")
	registryCreateCommand.Flags().StringVarP(&authorizationSecret, "authorization-secret", "", "", "registry authorization secret name")
	registryCreateCommand.Flags().StringVarP(&provider, "provider", "", "", "registry provider (\"simple\", \"minikube\", or \"openshift\")")
	registryCreateCommand.Flags().BoolVarP(&wait, "wait", "w", false, "wait for registry surrogate to come up")
}

var registryCreateCommand = &cobra.Command{
	Use:   "create [REGISTRY NAME]",
	Short: "Create a registry",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		CreateRegistry(args[0])
	},
}

func CreateRegistry(registryName string) {
	if (host == "") && (service == "") && (provider == "") {
		failRegistryCreate()
	}

	client := NewClient()

	if host != "" {
		if (service != "") || (provider != "") {
			failRegistryCreate()
		}
	} else if service != "" {
		if (host != "") || (provider != "") {
			failRegistryCreate()
		}
	} else if provider != "" {
		if (host != "") || (service != "") {
			failRegistryCreate()
		}

		switch provider {
		case "simple":
			service = "reposure-simple"
			if authenticationSecret == "" {
				authenticationSecret = "reposure-simple-authentication"
			}

		case "minikube":
			// Make sure to install the operator with "--role=view" so it can access "kube-system"

			// Note: The Docker container runtime always treats the registry at "127.0.0.1" as insecure
			// However CRI-O does not, thus the most compatible approach is to use the service
			// See: https://github.com/kubernetes/minikube/issues/6982
			serviceNamespace = "kube-system"
			service = "registry"
			// Insecure on port 80
			port = 80

		case "openshift":
			host = "image-registry.openshift-image-registry.svc:5000"
			if (authenticationSecret == "") || (authorizationSecret == "") {
				// We will use the "builder" service account's service-ca certificate and auth token
				serviceAccount, err := client.Kubernetes.CoreV1().ServiceAccounts(client.Namespace).Get(contextpkg.TODO(), "builder", meta.GetOptions{})
				util.FailOnError(err)
				for _, secretName := range serviceAccount.Secrets {
					secret, err := client.Kubernetes.CoreV1().Secrets(client.Namespace).Get(contextpkg.TODO(), secretName.Name, meta.GetOptions{})
					util.FailOnError(err)
					if secret.Type == core.SecretTypeServiceAccountToken {
						if authenticationSecret == "" {
							authenticationSecret = secret.Name
						}
						if authenticationSecretDataKey == "" {
							authenticationSecretDataKey = "service-ca.crt"
						}
						if authorizationSecret == "" {
							authorizationSecret = secret.Name
						}
						break
					}
				}
			}

		default:
			util.Fail("unsupported \"--provider\": must be \"simple\", \"minikube\", or \"openshift\"")
		}
	}

	adminClient := client.AdminClient()

	var err error
	if service != "" {
		_, err = adminClient.CreateRegistryIndirect(namespace, registryName, serviceNamespace, service, port, authenticationSecret, authenticationSecretDataKey, authorizationSecret)
	} else {
		_, err = adminClient.CreateRegistryDirect(namespace, registryName, host, authenticationSecret, authenticationSecretDataKey, authorizationSecret)
	}
	util.FailOnError(err)

	if wait {
		_, err = adminClient.WaitForRegistrySurrogate(namespace, registryName)
		util.FailOnError(err)
	}
}

func failRegistryCreate() {
	util.Fail("must specify only one of \"--host\", \"--service\", or \"--provider\"")
}
