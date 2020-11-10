package command

import (
	"io"
)

func (self *Client) PullTarball(imageReference string, writer io.Writer) error {
	return self.Command(writer, "pull", imageReference)
}
