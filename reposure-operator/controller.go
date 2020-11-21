package main

import (
	"fmt"
	"net/http"

	"github.com/heptiolabs/healthcheck"
	"github.com/tebeka/atexit"
	"github.com/tliron/kutil/kubernetes"
	"github.com/tliron/kutil/util"
	versionpkg "github.com/tliron/kutil/version"
	reposurepkg "github.com/tliron/reposure/apis/clientset/versioned"
	controllerpkg "github.com/tliron/reposure/controller"
	apiextensionspkg "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/dynamic"
	kubernetespkg "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	// Load all auth plugins:
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

func Controller() {
	if version {
		versionpkg.Print()
		atexit.Exit(0)
		return
	}

	log.Noticef("%s version=%s revision=%s", toolName, versionpkg.GitVersion, versionpkg.GitRevision)

	// Config

	config, err := clientcmd.BuildConfigFromFlags(masterUrl, kubeconfigPath)
	util.FailOnError(err)

	if cluster {
		namespace = ""
	} else if namespace == "" {
		if namespace_, ok := kubernetes.GetConfiguredNamespace(kubeconfigPath, context); ok {
			namespace = namespace_
		}
		if namespace == "" {
			namespace = kubernetes.GetServiceAccountNamespace()
		}
		if namespace == "" {
			log.Fatal("could not discover namespace and namespace not provided")
		}
	}

	// Clients

	kubernetesClient, err := kubernetespkg.NewForConfig(config)
	util.FailOnError(err)

	apiExtensionsClient, err := apiextensionspkg.NewForConfig(config)
	util.FailOnError(err)

	dynamicClient, err := dynamic.NewForConfig(config)
	util.FailOnError(err)

	reposureClient, err := reposurepkg.NewForConfig(config)
	util.FailOnError(err)

	// Controller

	controller := controllerpkg.NewController(
		toolName,
		cluster,
		namespace,
		dynamicClient,
		kubernetesClient,
		apiExtensionsClient,
		reposureClient,
		config,
		resyncPeriod,
		util.SetupSignalHandler(),
	)

	// Run

	err = controller.Run(concurrency, func() {
		log.Info("starting health monitor")
		health := healthcheck.NewHandler()
		err := http.ListenAndServe(fmt.Sprintf(":%d", healthPort), health)
		util.FailOnError(err)
	})
	util.FailOnError(err)
}
