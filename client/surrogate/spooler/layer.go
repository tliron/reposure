package spooler

import (
	"io"

	urlpkg "github.com/tliron/exturl"
	"github.com/tliron/kutil/util"
)

func (self *Client) PushLayerFromURL(fileName string, url urlpkg.URL) error {
	if size, err := urlpkg.Size(url); err == nil {
		if reader, err := url.Open(); err == nil {
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
