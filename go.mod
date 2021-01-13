module github.com/tliron/reposure

go 1.15

// replace github.com/tliron/kutil => /Depot/Projects/RedHat/kutil

require (
	github.com/fsnotify/fsnotify v1.4.9
	github.com/gofrs/flock v0.8.0
	github.com/google/go-containerregistry v0.3.0
	github.com/heptiolabs/healthcheck v0.0.0-20180807145615-6ff867650f40
	github.com/jetstack/cert-manager v1.1.0
	github.com/klauspost/pgzip v1.2.5
	github.com/op/go-logging v0.0.0-20160315200505-970db520ece7
	github.com/spf13/cobra v1.1.1
	github.com/tebeka/atexit v0.3.0
	github.com/tliron/kutil v0.1.12
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad
	k8s.io/api v0.20.1
	k8s.io/apiextensions-apiserver v0.20.1
	k8s.io/apimachinery v0.20.1
	k8s.io/client-go v0.20.1
)
