package direct

import (
	"fmt"

	namepkg "github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

func (self *Client) ListImages() ([]string, error) {
	if registry, err := self.newRegistry(self.Host); err == nil {
		if repositoryNames, err := remote.Catalog(self.Context, registry, self.Options...); err == nil {
			var imageReferences []string

			for _, repositoryName := range repositoryNames {
				repositoryName_ := self.getName(repositoryName)
				if repository, err := namepkg.NewRepository(repositoryName_); err == nil {
					if repositoryImageReferences, err := remote.List(repository, self.Options...); err == nil {
						for _, repositoryImageReference := range repositoryImageReferences {
							repositoryImageReference = fmt.Sprintf("%s:%s", repositoryName, repositoryImageReference)
							imageReferences = append(imageReferences, repositoryImageReference)
						}
					} else {
						return nil, err
					}
				} else {
					return nil, err
				}
			}

			return imageReferences, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func (self *Client) DeleteImage(imageName string) error {
	// Note: the registry repository in which the image is located
	// will *not* be deleted, even if this deletion causes it to be
	// empty

	name := self.getName(imageName)
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
