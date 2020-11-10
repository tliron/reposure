package main

import (
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/op/go-logging"
)

type Watcher struct {
	watcher  *fsnotify.Watcher
	handlers []*Handler
	log      *logging.Logger
}

type HandlerFunc func(path string)

type Handler struct {
	op     fsnotify.Op
	handle HandlerFunc
}

func NewWatcher() (*Watcher, error) {
	if watcher, err := fsnotify.NewWatcher(); err == nil {
		return &Watcher{
			watcher: watcher,
			log:     logging.MustGetLogger("watcher"),
		}, nil
	} else {
		return nil, err
	}
}

func (self *Watcher) Add(directoryPath string, op fsnotify.Op, handle HandlerFunc) error {
	self.handlers = append(self.handlers, &Handler{
		op:     op,
		handle: handle,
	})
	return self.watcher.Add(directoryPath)
}

func (self *Watcher) Close() error {
	return self.watcher.Close()
}

func (self *Watcher) Run() {
	defer self.Close()
	for self.Process() {
	}
}

func (self *Watcher) Process() bool {
	select {
	case event, ok := <-self.watcher.Events:
		if !ok {
			self.log.Warning("no more events")
			return false
		}

		// Ignore temporary files
		if strings.HasSuffix(event.Name, "~") {
			break
		}

		self.log.Debugf("event: %s", event)

		for _, handler := range self.handlers {
			if event.Op&handler.op == handler.op {
				handler.handle(event.Name)
			}
		}

	case err, ok := <-self.watcher.Errors:
		if !ok {
			self.log.Warning("no more errors")
			return false
		}

		self.log.Errorf("error: %s", err.Error())
	}

	return true
}
