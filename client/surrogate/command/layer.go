package command

import (
	"io"

	"github.com/tliron/kutil/compression"
)

func (self *Client) PullLayer(imageName string, writer io.Writer) error {
	pipeReader, pipeWriter := io.Pipe()

	go func() {
		if err := self.PullTarball(imageName, pipeWriter); err == nil {
			pipeWriter.Close()
		} else {
			pipeWriter.CloseWithError(err)
		}
	}()

	decoder := compression.NewFirstTarballInTarballDecoder(pipeReader)
	if _, err := io.Copy(writer, decoder.Decode()); err == nil {
		return nil
	} else {
		return err
	}
}
