package main

import (
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/tliron/commonlog"
)

type Watcher struct {
	Watcher  *fsnotify.Watcher
	Handlers []*Handler
	Log      commonlog.Logger
}

type HandlerFunc func(path string)

type Handler struct {
	op     fsnotify.Op
	handle HandlerFunc
}

func NewWatcher() (*Watcher, error) {
	if watcher, err := fsnotify.NewWatcher(); err == nil {
		return &Watcher{
			Watcher: watcher,
			Log:     commonlog.GetLogger("watcher"),
		}, nil
	} else {
		return nil, err
	}
}

func (self *Watcher) Add(directoryPath string, op fsnotify.Op, handle HandlerFunc) error {
	self.Handlers = append(self.Handlers, &Handler{
		op:     op,
		handle: handle,
	})
	return self.Watcher.Add(directoryPath)
}

func (self *Watcher) Close() error {
	return self.Watcher.Close()
}

func (self *Watcher) Run() {
	defer self.Close()
	for self.Process() {
	}
}

func (self *Watcher) Process() bool {
	select {
	case event, ok := <-self.Watcher.Events:
		if !ok {
			self.Log.Warning("no more events")
			return false
		}

		// Ignore temporary files
		if strings.HasSuffix(event.Name, "~") {
			break
		}

		self.Log.Debugf("event: %s", event)

		for _, handler := range self.Handlers {
			if event.Op&handler.op == handler.op {
				handler.handle(event.Name)
			}
		}

	case err, ok := <-self.Watcher.Errors:
		if !ok {
			self.Log.Warning("no more errors")
			return false
		}

		self.Log.Errorf("error: %s", err.Error())
	}

	return true
}
