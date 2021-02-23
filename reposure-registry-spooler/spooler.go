package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/heptiolabs/healthcheck"
	"github.com/tliron/kutil/logging"
	"github.com/tliron/kutil/util"
	directclient "github.com/tliron/reposure/client/direct"
)

var log = logging.GetLogger("reposure-registry-spooler")

func RunSpooler(registryUrl string, path string) {
	stopChannel := util.SetupSignalHandler()

	var roundTripper http.RoundTripper
	if certificatePath != "" {
		log.Infof("certificate path: %s", certificatePath)
		var err error
		roundTripper, err = directclient.TLSRoundTripper(certificatePath)
		util.FailOnError(err)
	}

	/*if username != "" {
		log.Infof("username: %s", username)
		log.Infof("password: %s", password)
	} else if token != "" {
		log.Infof("token: %s", token)
	}*/

	publisher := NewPublisher(context, registryUrl, roundTripper, username, password, token, queue)
	log.Info("starting publisher")
	go publisher.Run()
	defer publisher.Close()

	fileInfos, err := ioutil.ReadDir(path)
	util.FailOnError(err)
	for _, fileInfo := range fileInfos {
		publisher.Enqueue(filepath.Join(path, fileInfo.Name()))
	}

	watcher, err := NewWatcher()
	util.FailOnError(err)

	err = watcher.Add(path, fsnotify.Create, func(path string) {
		publisher.Enqueue(path)
	})
	util.FailOnError(err)

	log.Info("starting watcher")
	go watcher.Run()

	go func() {
		log.Info("starting health monitor")
		health := healthcheck.NewHandler()
		err := http.ListenAndServe(fmt.Sprintf(":%d", healthPort), health)
		util.FailOnError(err)
	}()

	<-stopChannel
}
