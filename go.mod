module github.com/tliron/reposure

go 1.16

// replace github.com/tliron/kutil => /Depot/Projects/RedHat/kutil

require (
	github.com/fsnotify/fsnotify v1.4.9
	github.com/gofrs/flock v0.8.0
	github.com/google/go-containerregistry v0.4.1
	github.com/heptiolabs/healthcheck v0.0.0-20180807145615-6ff867650f40
	github.com/jetstack/cert-manager v1.2.0
	github.com/klauspost/pgzip v1.2.5
	github.com/spf13/cobra v1.1.3
	github.com/tliron/kutil v0.1.22
	golang.org/x/crypto v0.0.0-20210220033148-5ea612d1eb83
	gopkg.in/DATA-DOG/go-sqlmock.v1 v1.3.0 // indirect
	k8s.io/api v0.20.4
	k8s.io/apiextensions-apiserver v0.20.4
	k8s.io/apimachinery v0.20.4
	k8s.io/client-go v0.20.4
	k8s.io/klog/v2 v2.6.0
)
