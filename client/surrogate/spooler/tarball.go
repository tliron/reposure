package spooler

import (
	"io"
	"os"
	"strings"

	urlpkg "github.com/tliron/exturl"
)

func (self *Client) PushTarballFromURL(imageName string, url urlpkg.URL) error {
	reader, err := url.Open()
	if err != nil {
		return err
	}

	reader, err = url.Open()
	if err != nil {
		return err
	}

	if readCloser, ok := reader.(io.ReadCloser); ok {
		defer readCloser.Close()
	}

	if err = self.PushTarball(imageName, reader); err == nil {
		return nil
	} else {
		return err
	}
}

func (self *Client) PushTarballFromFile(imageName string, path string) error {
	fileName := imageName
	if strings.HasSuffix(path, ".tar.gz") {
		fileName += ".tar.gz"
	} else if strings.HasSuffix(path, ".tgz") {
		fileName += ".tgz"
	} else if strings.HasSuffix(path, ".tar") {
		fileName += ".tar"
	}

	if file, err := os.Open(path); err == nil {
		defer file.Close()
		return self.PushTarball(fileName, file)
	} else {
		return err
	}
}

func (self *Client) PushTarball(fileName string, reader io.Reader) error {
	if podName, err := self.getFirstPodName(); err == nil {
		path := self.getPath(fileName)
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
