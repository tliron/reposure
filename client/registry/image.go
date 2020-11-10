package registry

import (
	"fmt"
	"os"

	namepkg "github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

func (self *Client) ListImages(registryHost string) ([]string, error) {
	if registry, err := self.newRegistry(registryHost); err == nil {
		if repositoryNames, err := remote.Catalog(self.Context, registry, self.Options...); err == nil {
			var imageReferences []string

			for _, repositoryName := range repositoryNames {
				repositoryName_ := fmt.Sprintf("%s/%s", registryHost, repositoryName)
				if repository, err := namepkg.NewRepository(repositoryName_); err == nil {
					if repositoryImageReferences, err := remote.List(repository, self.Options...); err == nil {
						for _, repositoryImageReference := range repositoryImageReferences {
							repositoryImageReference = fmt.Sprintf("%s:%s", repositoryName, repositoryImageReference)
							imageReferences = append(imageReferences, repositoryImageReference)
						}
					} else {
						fmt.Fprintln(os.Stderr, "1")
						return nil, err
					}
				} else {
					fmt.Fprintln(os.Stderr, "2")
					return nil, err
				}
			}

			return imageReferences, nil
		} else {
			fmt.Fprintln(os.Stderr, "3")
			return nil, err
		}
	} else {
		fmt.Fprintln(os.Stderr, "4")
		return nil, err
	}
}

func (self *Client) DeleteImage(name string) error {
	// Note: the registry repository in which the image is located
	// will *not* be deleted, even if this deletion causes it to be
	// empty

	if tag, err := namepkg.NewTag(name); err == nil {
		if image, err := remote.Image(tag, self.Options...); err == nil {
			if hash, err := image.Digest(); err == nil {
				digest := tag.Digest(hash.String())

				return remote.Delete(digest, self.Options...)
			} else {
				return err
			}
		} else {
			return err
		}
	} else {
		return err
	}
}
