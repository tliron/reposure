package spooler

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
	SpoolDirectory       string
}

func NewClient(kubernetes kubernetes.Interface, rest rest.Interface, config *rest.Config, context contextpkg.Context, stderr io.Writer, namespace string, surrogateAppName string, spoolerContainerName string, spoolDirectory string) *Client {
	return &Client{
		Kubernetes: kubernetes,
		REST:       rest,
		Config:     config,
		Context:    context,
		Stderr:     stderr,

		Namespace:            namespace,
		SurrogateAppName:     surrogateAppName,
		SpoolerContainerName: spoolerContainerName,
		SpoolDirectory:       spoolDirectory,
	}
}
