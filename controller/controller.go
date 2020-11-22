package controller

import (
	contextpkg "context"
	"fmt"
	"time"

	"github.com/op/go-logging"
	kubernetesutil "github.com/tliron/kutil/kubernetes"
	reposureclientset "github.com/tliron/reposure/apis/clientset/versioned"
	reposureinformers "github.com/tliron/reposure/apis/informers/externalversions"
	reposurelisters "github.com/tliron/reposure/apis/listers/reposure.puccini.cloud/v1alpha1"
	adminclient "github.com/tliron/reposure/client/admin"
	reposureresources "github.com/tliron/reposure/resources/reposure.puccini.cloud/v1alpha1"
	apiextensionspkg "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	dynamicpkg "k8s.io/client-go/dynamic"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	restpkg "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
)

//
// Controller
//

type Controller struct {
	Config      *restpkg.Config
	Dynamic     *kubernetesutil.Dynamic
	Kubernetes  kubernetes.Interface
	Reposure    reposureclientset.Interface
	Client      *adminclient.Client
	CachePath   string
	StopChannel <-chan struct{}

	Processors *kubernetesutil.Processors
	Events     record.EventRecorder

	KubernetesInformerFactory informers.SharedInformerFactory
	ReposureInformerFactory   reposureinformers.SharedInformerFactory

	Registries reposurelisters.RegistryLister

	Context contextpkg.Context
	Log     *logging.Logger
}

func NewController(toolName string, clusterMode bool, clusterRole string, namespace string, dynamic dynamicpkg.Interface, kubernetes kubernetes.Interface, apiExtensions apiextensionspkg.Interface, reposure reposureclientset.Interface, config *restpkg.Config, informerResyncPeriod time.Duration, stopChannel <-chan struct{}) *Controller {
	context := contextpkg.TODO()

	if clusterMode {
		namespace = ""
		if clusterRole != "" {
			clusterRole = "cluster-admin"
		}
	}

	log := logging.MustGetLogger(fmt.Sprintf("%s.controller", toolName))

	self := Controller{
		Config:      config,
		Dynamic:     kubernetesutil.NewDynamic(toolName, dynamic, kubernetes.Discovery(), namespace, context),
		Kubernetes:  kubernetes,
		Reposure:    reposure,
		StopChannel: stopChannel,
		Processors:  kubernetesutil.NewProcessors(toolName),
		Events:      kubernetesutil.CreateEventRecorder(kubernetes, "Reposure", log),
		Context:     context,
		Log:         log,
	}

	self.Client = adminclient.NewClient(
		kubernetes,
		apiExtensions,
		reposure,
		kubernetes.CoreV1().RESTClient(),
		config,
		context,
		clusterMode,
		clusterRole,
		namespace,
		NamePrefix,
		PartOf,
		ManagedBy,
		OperatorImageReference,
		SurrogateImageReference,
		SimpleImageReference,
		fmt.Sprintf("%s.client", toolName),
	)

	if clusterMode {
		self.KubernetesInformerFactory = informers.NewSharedInformerFactory(kubernetes, informerResyncPeriod)
		self.ReposureInformerFactory = reposureinformers.NewSharedInformerFactory(reposure, informerResyncPeriod)
	} else {
		self.KubernetesInformerFactory = informers.NewSharedInformerFactoryWithOptions(kubernetes, informerResyncPeriod, informers.WithNamespace(namespace))
		self.ReposureInformerFactory = reposureinformers.NewSharedInformerFactoryWithOptions(reposure, informerResyncPeriod, reposureinformers.WithNamespace(namespace))
	}

	// Informers
	registryInformer := self.ReposureInformerFactory.Reposure().V1alpha1().Registries()

	// Listers
	self.Registries = registryInformer.Lister()

	// Processors

	processorPeriod := 5 * time.Second

	self.Processors.Add(reposureresources.RegistryGVK, kubernetesutil.NewProcessor(
		toolName,
		"registries",
		registryInformer.Informer(),
		processorPeriod,
		func(name string, namespace string) (interface{}, error) {
			return self.Client.GetRegistry(namespace, name)
		},
		func(object interface{}) (bool, error) {
			return self.processRegistry(object.(*reposureresources.Registry))
		},
	))

	return &self
}

func (self *Controller) Run(concurrency uint, startup func()) error {
	defer utilruntime.HandleCrash()

	self.Log.Info("starting informer factories")
	self.KubernetesInformerFactory.Start(self.StopChannel)
	self.ReposureInformerFactory.Start(self.StopChannel)

	self.Log.Info("waiting for processor informer caches to sync")
	utilruntime.HandleError(self.Processors.WaitForCacheSync(self.StopChannel))

	self.Log.Infof("starting processors (concurrency=%d)", concurrency)
	self.Processors.Start(concurrency, self.StopChannel)
	defer self.Processors.ShutDown()

	if startup != nil {
		go startup()
	}

	<-self.StopChannel

	self.Log.Info("shutting down")

	return nil
}
