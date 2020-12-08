package direct

import (
	"io"
	"os"

	namepkg "github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
	gzip "github.com/klauspost/pgzip"
	urlpkg "github.com/tliron/kutil/url"
)

func (self *Client) PushTarball(path string, imageName string) error {
	name := self.getName(imageName)
	if tag, err := namepkg.NewTag(name); err == nil {
		if image, err := tarball.ImageFromPath(path, &tag); err == nil {
			return remote.Write(tag, image, self.Options...)
		} else {
			return err
		}
	} else {
		return err
	}
}

func (self *Client) PushGzippedTarball(path string, imageName string) error {
	opener := func() (io.ReadCloser, error) {
		if reader, err := os.Open(path); err == nil {
			return gzip.NewReader(reader)
		} else {
			return nil, err
		}
	}

	name := self.getName(imageName)
	if tag, err := namepkg.NewTag(name); err == nil {
		if image, err := tarball.Image(opener, &tag); err == nil {
			return remote.Write(tag, image, self.Options...)
		} else {
			return err
		}
	} else {
		return err
	}
}

func (self *Client) PushGzippedTarballFromURL(url urlpkg.URL, imageName string) (string, error) {
	opener := func() (io.ReadCloser, error) {
		if reader, err := url.Open(); err == nil {
			return gzip.NewReader(reader)
		} else {
			return nil, err
		}
	}

	if contentTag, err := namepkg.NewTag("portable"); err == nil {
		name := self.getName(imageName)
		if tag, err := namepkg.NewTag(name); err == nil {
			if image, err := tarball.Image(opener, &contentTag); err == nil {
				if err := remote.Write(tag, image, self.Options...); err == nil {
					return name, nil
				} else {
					return "", err
				}
			} else {
				return "", err
			}
		} else {
			return "", err
		}
	} else {
		return "", err
	}
}

func (self *Client) PullTarball(imageName string, path string) error {
	name := self.getName(imageName)
	if tag, err := namepkg.NewTag(name); err == nil {
		if image, err := remote.Image(tag, self.Options...); err == nil {
			var writer io.Writer
			if path == "" {
				writer = os.Stdout
			} else {
				if file, err := os.Create(path); err == nil {
					defer file.Close()
					writer = file
				} else {
					return err
				}
			}

			return tarball.Write(tag, image, writer)
		} else {
			return err
		}
	} else {
		return err
	}
}
