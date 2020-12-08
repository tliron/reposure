package main

import (
	contextpkg "context"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofrs/flock"
	"github.com/op/go-logging"
	directclient "github.com/tliron/reposure/client/direct"
)

type Publisher struct {
	client *directclient.Client
	work   chan string
	log    *logging.Logger
}

func NewPublisher(context contextpkg.Context, host string, roundTripper http.RoundTripper, username string, password string, token string, queue int) *Publisher {
	return &Publisher{
		client: directclient.NewClient(host, roundTripper, username, password, token, context),
		work:   make(chan string, queue),
		log:    logging.MustGetLogger("publisher"),
	}
}

func (self *Publisher) Enqueue(path string) {
	self.log.Debugf("enqueuing: %s", path)
	self.work <- path
}

func (self *Publisher) Close() {
	close(self.work)
}

func (self *Publisher) Run() {
	defer self.Close()
	for self.Process() {
	}
}

func (self *Publisher) Process() bool {
	if path, ok := <-self.work; ok {
		// Lock file
		lock := flock.New(path)
		if err := lock.Lock(); err == nil {
			defer lock.Unlock()
		} else {
			self.log.Errorf("could not lock file %q: %s", path, err.Error())
			return true
		}

		// File may have already been deleted by another process
		if _, err := os.Stat(path); err != nil {
			if os.IsNotExist(err) {
				self.log.Infof("file %q already deleted", path)
			} else {
				self.log.Errorf("could not access file %q: %s", path, err.Error())
			}
			return true
		}

		/// Process
		if strings.HasSuffix(path, "!") {
			self.Delete(path[:len(path)-1])
		} else {
			self.Publish(path)
		}

		// Delete file
		if err := os.Remove(path); err == nil {
			self.log.Infof("deleted file %q", path)
		} else {
			self.log.Errorf("could not delete file %q: %s", path, err.Error())
		}

		return true
	} else {
		self.log.Warning("no more work")
		return false
	}
}

func (self *Publisher) Publish(path string) {
	imageName := self.getImageName(path)

	var err error
	if strings.HasSuffix(path, ".tar.gz") || strings.HasSuffix(path, ".tgz") {
		self.log.Infof("publishing gzipped tarball %q to image %q", path, imageName)
		err = self.client.PushGzippedTarball(path, imageName)
	} else if strings.HasSuffix(path, ".tar") {
		self.log.Infof("publishing tarball %q to image %q", path, imageName)
		err = self.client.PushTarball(path, imageName)
	} else {
		self.log.Infof("publishing layer %q to image %q", path, imageName)
		if file, err2 := os.Open(path); err2 == nil {
			err = self.client.PushLayer(file, imageName)
		} else {
			self.log.Errorf("could not read file %q: %s", path, err2.Error())
		}
	}

	if err == nil {
		self.log.Infof("published image %q", imageName)
	} else {
		self.log.Errorf("could not publish image %q: %s", imageName, err.Error())
	}
}

func (self *Publisher) Delete(path string) {
	imageName := self.getImageName(path)
	self.log.Infof("deleting image %q", imageName)
	if err := self.client.DeleteImage(imageName); err == nil {
		self.log.Infof("deleted image %q", imageName)
	} else {
		self.log.Errorf("could not delete image %q: %s", imageName, err.Error())
	}
}

func (self *Publisher) getImageName(path string) string {
	name := filepath.Base(path)
	if dot := strings.Index(name, "."); dot != -1 {
		// Note: filepath.Ext will return the last extension only
		name = name[:dot]
	}
	return strings.ReplaceAll(name, "\\", "/")
}
