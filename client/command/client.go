package command

import (
	contextpkg "context"
	"io"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

//
// Client
//

type Client struct {
	Kubernetes kubernetes.Interface
	REST       rest.Interface
	Config     *rest.Config
	Context    contextpkg.Context
	Stderr     io.Writer

	Namespace            string
	SurrogateAppName     string
	SpoolerContainerName string
	Registry             string
	RegistryCertificate  string
	RegistryUsername     string
	RegistryPassword     string
	RegistryToken        string
}

func NewClient(kubernetes kubernetes.Interface, rest rest.Interface, config *rest.Config, context contextpkg.Context, stderr io.Writer, namespace string, surrogateAppName string, spoolerContainerName string, registry string, certificate string, username string, password string, token string) *Client {
	if registry == "" {
		// Default for sidecars
		registry = "localhost:5000"
	}

	return &Client{
		Kubernetes:           kubernetes,
		REST:                 rest,
		Config:               config,
		Context:              context,
		Stderr:               stderr,
		Namespace:            namespace,
		SurrogateAppName:     surrogateAppName,
		SpoolerContainerName: spoolerContainerName,
		Registry:             registry,
		RegistryCertificate:  certificate,
		RegistryUsername:     username,
		RegistryPassword:     password,
		RegistryToken:        token,
	}
}
