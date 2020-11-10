package spooler

import (
	"io"
)

func (self *Client) PushImage(imageReference string, reader io.Reader) error {
	if podName, err := self.getFirstPodName(); err == nil {
		path := self.getPath(imageReference)
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

func (self *Client) DeleteImage(imageReference string) error {
	if podName, err := self.getFirstPodName(); err == nil {
		path := self.getPath(imageReference) + "!"
		return self.touch(podName, path)
	} else {
		return err
	}
}
