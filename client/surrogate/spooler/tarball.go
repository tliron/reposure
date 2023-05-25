package spooler

import (
	contextpkg "context"
	"io"
	"os"
	"strings"

	"github.com/tliron/exturl"
	"github.com/tliron/kutil/util"
)

func (self *Client) PushTarballFromURL(context contextpkg.Context, imageName string, url exturl.URL) error {
	reader, err := url.Open(context)
	if err != nil {
		return err
	}

	reader, err = url.Open(context)
	if err != nil {
		return err
	}

	reader = util.NewContextualReadCloser(context, reader)
	defer reader.Close()
	if err = self.PushTarball(imageName, reader); err == nil {
		return nil
	} else {
		return err
	}
}

func (self *Client) PushTarballFromFile(context contextpkg.Context, imageName string, path string) error {
	fileName := imageName
	if strings.HasSuffix(path, ".tar.gz") {
		fileName += ".tar.gz"
	} else if strings.HasSuffix(path, ".tgz") {
		fileName += ".tgz"
	} else if strings.HasSuffix(path, ".tar") {
		fileName += ".tar"
	}

	if file, err := os.Open(path); err == nil {
		reader := util.NewContextualReadCloser(context, file)
		defer reader.Close()
		return self.PushTarball(fileName, reader)
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
