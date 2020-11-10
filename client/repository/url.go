package repository

import (
	"fmt"
	"io"

	namepkg "github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
	gzip "github.com/klauspost/pgzip"
	urlpkg "github.com/tliron/kutil/url"
	resources "github.com/tliron/reposure/resources/reposure.puccini.cloud/v1alpha1"
)

func (self *Client) UpdateURLContext(repository *resources.Repository, urlContext *urlpkg.Context) error {
	if host, roundTripper, err := self.GetHTTPRoundTripper(repository); err == nil {
		if roundTripper != nil {
			urlContext.SetHTTPRoundTripper(host, roundTripper)
		}
	} else {
		return err
	}

	if host, username, password, token, err := self.GetAuth(repository); err == nil {
		if (username != "") || (token != "") {
			urlContext.SetCredentials(host, username, password, token)
		}
	} else {
		return err
	}

	return nil
}

// TODO: what uses this?
func (self *Client) PushTarball(artifactName string, sourceUrl urlpkg.URL, repository *resources.Repository) (string, error) {
	if repositoryHost, err := self.GetHost(repository); err == nil {
		if options, err := self.GetRemoteOptions(repository); err == nil {
			self.Log.Infof("publishing image %q at %q on %q", artifactName, sourceUrl.String(), repositoryHost)

			opener := func() (io.ReadCloser, error) {
				if reader, err := sourceUrl.Open(); err == nil {
					return gzip.NewReader(reader)
				} else {
					return nil, err
				}
			}

			if contentTag, err := namepkg.NewTag("portable"); err == nil {
				tag := fmt.Sprintf("%s/%s", repositoryHost, artifactName)
				if tag_, err := namepkg.NewTag(tag); err == nil {
					if image, err := tarball.Image(opener, &contentTag); err == nil {
						if err := remote.Write(tag_, image, options...); err == nil {
							self.Log.Infof("published image %q at %q on %q", tag, sourceUrl.String(), repositoryHost)
							return tag, nil
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
		} else {
			return "", err
		}
	} else {
		return "", err
	}
}
