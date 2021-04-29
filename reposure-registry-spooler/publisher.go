package main

import (
	contextpkg "context"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofrs/flock"
	"github.com/tliron/kutil/logging"
	"github.com/tliron/reposure/client/direct"
)

type Publisher struct {
	Client *direct.Client
	Work   chan string
	Log    logging.Logger
}

func NewPublisher(context contextpkg.Context, host string, roundTripper http.RoundTripper, username string, password string, token string, queue int) *Publisher {
	return &Publisher{
		Client: direct.NewClient(host, roundTripper, username, password, token, context),
		Work:   make(chan string, queue),
		Log:    logging.GetLogger("publisher"),
	}
}

func (self *Publisher) Enqueue(path string) {
	self.Log.Debugf("enqueuing: %s", path)
	self.Work <- path
}

func (self *Publisher) Close() {
	close(self.Work)
}

func (self *Publisher) Run() {
	defer self.Close()
	for self.Process() {
	}
}

func (self *Publisher) Process() bool {
	if path, ok := <-self.Work; ok {
		// Lock file
		lock := flock.New(path)
		if err := lock.Lock(); err == nil {
			defer lock.Unlock()
		} else {
			self.Log.Errorf("could not lock file %q: %s", path, err.Error())
			return true
		}

		// File may have already been deleted by another process
		if _, err := os.Stat(path); err != nil {
			if os.IsNotExist(err) {
				self.Log.Infof("file %q already deleted", path)
			} else {
				self.Log.Errorf("could not access file %q: %s", path, err.Error())
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
			self.Log.Infof("deleted file %q", path)
		} else {
			self.Log.Errorf("could not delete file %q: %s", path, err.Error())
		}

		return true
	} else {
		self.Log.Warning("no more work")
		return false
	}
}

func (self *Publisher) Publish(path string) {
	imageName := self.getImageName(path)

	var err error
	if strings.HasSuffix(path, ".tar.gz") || strings.HasSuffix(path, ".tgz") {
		self.Log.Infof("publishing gzipped tarball %q to image %q", path, imageName)
		err = self.Client.PushGzippedTarball(path, imageName)
	} else if strings.HasSuffix(path, ".tar") {
		self.Log.Infof("publishing tarball %q to image %q", path, imageName)
		err = self.Client.PushTarball(path, imageName)
	} else {
		self.Log.Infof("publishing layer %q to image %q", path, imageName)
		if file, err2 := os.Open(path); err2 == nil {
			err = self.Client.PushLayer(file, imageName)
		} else {
			self.Log.Errorf("could not read file %q: %s", path, err2.Error())
		}
	}

	if err == nil {
		self.Log.Infof("published image %q", imageName)
	} else {
		self.Log.Errorf("could not publish image %q: %s", imageName, err.Error())
	}
}

func (self *Publisher) Delete(path string) {
	imageName := self.getImageName(path)
	self.Log.Infof("deleting image %q", imageName)
	if err := self.Client.DeleteImage(imageName); err == nil {
		self.Log.Infof("deleted image %q", imageName)
	} else {
		self.Log.Errorf("could not delete image %q: %s", imageName, err.Error())
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
