package spooler

import (
	contextpkg "context"
	"io"

	"github.com/tliron/exturl"
	"github.com/tliron/kutil/util"
)

func (self *Client) PushLayerFromURL(context contextpkg.Context, fileName string, url exturl.URL) error {
	if size, err := exturl.Size(context, url); err == nil {
		if reader, err := url.Open(context); err == nil {
			defer reader.Close()
			return self.PushLayer(fileName, reader, size)
		} else {
			return err
		}
	} else {
		return err
	}
}

func (self *Client) PushLayer(fileName string, reader io.Reader, size int64) error {
	encoder := util.NewTarEncoder(reader, "portable", size)
	return self.PushTarball(fileName, encoder.Encode())
}
