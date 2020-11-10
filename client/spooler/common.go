package spooler

import (
	"io"
	"path/filepath"
	"strings"

	kubernetesutil "github.com/tliron/kutil/kubernetes"
)

func (self *Client) getFirstPodName() (string, error) {
	return kubernetesutil.GetFirstPodName(self.Context, self.Kubernetes, self.Namespace, self.SurrogateAppName)
}

func (self *Client) getPath(imageReference string) string {
	return filepath.Join(self.SpoolDirectory, strings.ReplaceAll(imageReference, "/", "\\"))
}

func (self *Client) writeToContainer(podName string, reader io.Reader, targetPath string) error {
	return kubernetesutil.WriteToContainer(self.REST, self.Config, self.Namespace, podName, self.SpoolerContainerName, reader, targetPath, nil)
}

func (self *Client) readFromContainer(podName string, writer io.Writer, sourcePath string) error {
	return kubernetesutil.ReadFromContainer(self.REST, self.Config, self.Namespace, podName, self.SpoolerContainerName, writer, sourcePath)
}

func (self *Client) mv(podName string, fromPath string, toPath string) error {
	return self.exec(podName, nil, nil, "mv", fromPath, toPath)
}

func (self *Client) touch(podName string, path string) error {
	return self.exec(podName, nil, nil, "touch", path)
}

func (self *Client) exec(podName string, stdin io.Reader, stdout io.Writer, command ...string) error {
	return kubernetesutil.Exec(self.REST, self.Config, self.Namespace, podName, self.SpoolerContainerName, stdin, stdout, self.Stderr, false, command...)
}
