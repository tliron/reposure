package direct

import (
	"io"
	"io/ioutil"

	namepkg "github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/empty"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/go-containerregistry/pkg/v1/stream"
)

func (self *Client) PushLayer(readCloser io.ReadCloser, name string) error {
	if tag, err := namepkg.NewTag(name); err == nil {
		// See: https://github.com/google/go-containerregistry/issues/707
		layer := stream.NewLayer(ioutil.NopCloser(readCloser))
		//layer = stream.NewLayer(readCloser)

		if image, err := mutate.AppendLayers(empty.Image, layer); err == nil {
			return remote.Write(tag, image, self.Options...)
		} else {
			return err
		}
	} else {
		return err
	}
}
