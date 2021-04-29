package spooler

func (self *Client) DeleteImage(imageName string) error {
	if podName, err := self.getFirstPodName(); err == nil {
		path := self.getPath(imageName) + "!"
		return self.touch(podName, path)
	} else {
		return err
	}
}
