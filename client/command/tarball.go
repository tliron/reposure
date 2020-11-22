package command

import (
	"io"
)

func (self *Client) PullTarball(imageName string, writer io.Writer) error {
	return self.Command(writer, "pull", imageName)
}
