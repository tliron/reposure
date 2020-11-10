package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofrs/flock"
	"github.com/op/go-logging"
	registryclient "github.com/tliron/reposure/client/registry"
)

type Publisher struct {
	registry string
	client   *registryclient.Client
	work     chan string
	log      *logging.Logger
}

func NewPublisher(registry string, roundTripper http.RoundTripper, username string, password string, token string, queue int) *Publisher {
	return &Publisher{
		registry: registry,
		client:   registryclient.NewClient(roundTripper, username, password, token),
		work:     make(chan string, queue),
		log:      logging.MustGetLogger("publisher"),
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
	name := self.getImageReference(path)

	var err error
	if strings.HasSuffix(path, ".tar.gz") || strings.HasSuffix(path, ".tgz") {
		self.log.Infof("publishing gzipped tarball %q to image %q", path, name)
		err = self.client.PushGzippedTarball(path, name)
	} else if strings.HasSuffix(path, ".tar") {
		self.log.Infof("publishing tarball %q to image %q", path, name)
		err = self.client.PushTarball(path, name)
	} else {
		self.log.Infof("publishing layer %q to image %q", path, name)
		if file, err2 := os.Open(path); err2 == nil {
			err = self.client.PushLayer(file, name)
		} else {
			self.log.Errorf("could not read file %q: %s", path, err2.Error())
		}
	}

	if err == nil {
		self.log.Infof("published image %q", name)
	} else {
		self.log.Errorf("could not publish image %q: %s", name, err.Error())
	}
}

func (self *Publisher) Delete(path string) {
	name := self.getImageReference(path)
	self.log.Infof("deleting image %q", name)
	if err := self.client.DeleteImage(name); err == nil {
		self.log.Infof("deleted image %q", name)
	} else {
		self.log.Errorf("could not delete image %q: %s", name, err.Error())
	}
}

func (self *Publisher) getImageReference(path string) string {
	name := filepath.Base(path)
	if dot := strings.Index(name, "."); dot != -1 {
		// Note: filepath.Ext will return the last extension only
		name = name[:dot]
	}
	name = strings.ReplaceAll(name, "\\", "/")
	return fmt.Sprintf("%s/%s", self.registry, name)
}
