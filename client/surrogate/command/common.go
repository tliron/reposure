package command

import (
	"io"
	"strings"

	kubernetesutil "github.com/tliron/kutil/kubernetes"
)

func (self *Client) Command(writer io.Writer, arguments ...string) error {
	if podName, err := self.getFirstPodName(); err == nil {
		arguments = append([]string{"reposure-registry-client"}, arguments...)

		arguments = append(arguments, "--host", self.RegistryHost)

		if self.RegistryCertificate != "" {
			arguments = append(arguments, "--certificate", self.RegistryCertificate)
		}
		if self.RegistryUsername != "" {
			arguments = append(arguments, "--username", self.RegistryUsername)
		}
		if self.RegistryPassword != "" {
			arguments = append(arguments, "--password", self.RegistryPassword)
		}
		if self.RegistryToken != "" {
			arguments = append(arguments, "--token", self.RegistryToken)
		}

		self.Log.Debug(strings.Join(arguments, " "))

		return self.exec(podName, nil, writer, arguments...)
	} else {
		return err
	}
}

func (self *Client) getFirstPodName() (string, error) {
	return kubernetesutil.GetFirstPodName(self.Context, self.Kubernetes, self.Namespace, self.SurrogateAppName)
}

func (self *Client) exec(podName string, stdin io.Reader, stdout io.Writer, command ...string) error {
	return kubernetesutil.Exec(self.REST, self.Config, self.Namespace, podName, self.SpoolerContainerName, stdin, stdout, self.Stderr, false, command...)
}
