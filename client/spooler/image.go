package spooler

import (
	"io"
	"os"
	"strings"
)

func (self *Client) PushImageFromFile(imageName string, path string) error {
	name := imageName
	if strings.HasSuffix(path, ".tar.gz") {
		name += ".tar.gz"
	} else if strings.HasSuffix(path, ".tgz") {
		name += ".tgz"
	} else if strings.HasSuffix(path, ".tar") {
		name += ".tar"
	}

	if file, err := os.Open(path); err == nil {
		defer file.Close()
		return self.PushImage(name, file)
	} else {
		return err
	}
}

func (self *Client) PushImage(name string, reader io.Reader) error {
	if podName, err := self.getFirstPodName(); err == nil {
		path := self.getPath(name)
		tempPath := path + "~"
		if err := self.writeToContainer(podName, reader, tempPath); err == nil {
			return self.mv(podName, tempPath, path)
		} else {
			return err
		}
	} else {
		return err
	}
}

func (self *Client) DeleteImage(imageName string) error {
	if podName, err := self.getFirstPodName(); err == nil {
		path := self.getPath(imageName) + "!"
		return self.touch(podName, path)
	} else {
		return err
	}
}
