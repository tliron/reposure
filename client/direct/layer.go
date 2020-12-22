package direct

import (
	"io"

	namepkg "github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/empty"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/go-containerregistry/pkg/v1/stream"
)

func (self *Client) PushLayer(readCloser io.ReadCloser, imageName string) error {
	name := self.getName(imageName)
	if tag, err := namepkg.NewTag(name); err == nil {
		layer := stream.NewLayer(readCloser)
		if image, err := mutate.AppendLayers(empty.Image, layer); err == nil {
			return remote.Write(tag, image, self.Options...)
		} else {
			return err
		}
	} else {
		return err
	}
}
