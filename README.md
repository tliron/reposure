*This is an early release. Some features are not yet fully implemented.*

Reposure
========

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Latest Release](https://img.shields.io/github/release/tliron/reposure.svg)](https://github.com/tliron/reposure/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/tliron/reposure)](https://goreportcard.com/report/github.com/tliron/reposure)

Manage and access cloud-native container image registries for [Kubernetes](https://kubernetes.io/).


Get It
------

[![Download](assets/media/download.png "Download")](https://github.com/tliron/reposure/releases)


Rationale
---------

### What are cloud-native registries?

Kubernetes's container runtime, whether it's CRI-O or Docker or something else, can pull container
images from any OCI-compliant container image registry. While various public registries exist, such
as Docker Hub and Quay, it is often necessary to host a private registry.

It of course makes sense to use Kubernetes to deploy your private registry. Moreover, though your
private registry can run anywhere, it may also make sense to deploy it in the same Kubernetes
cluster that will be using it, a configuration that we're here calling a "cloud-native registry".
Doing so can significantly simplify the deployment of applications that require a private registry,
as indeed the registry would be just part of the workload. The most obvious use case is cloud-native
CI/CD pipelines, for which building and packaging code, publishing images, and deploying containers
all happen within the cluster. Indeed, platform as diverse as OpenShift, OKD, CodeReady Containers,
and Minikube come with built-in cloud-native registries exactly for this purpose.

### The challenge

Unfortunately, setting up cloud-native registries can be quite challenging. Your private TLS
certificates, certificate authorities, and authorization credentials, must all be configured into
the cluster's container runtime. This is in addition to the requirement that the container runtime,
which runs on the host, even has network connectivity to the registry, which uses Kubernetes's
control plane. The container runtime might also need access to the registry's domain name for TLS.

Additionally, beyond just getting the cloud-native registry to work with the container runtime, it
can be challenging to access it from outside the cluster using tools like Buildah and Skopeo. Would
an ingress to the cluster be required? Those are non-trivial and sometimes impossible to set up. And
what about TLS authentication and authorization?

### A solution

Reposure can assist with many of these challenges:

* The Reposure operator manages registry "surrogates" running in the cluster, thus they have network
  access to the registries. They are also configured with the necessary TLS authentication and
  authorization, if required. The surrogate can fetch and push container images for you. No ingress
  is required, just normal `kubectl` access to the cluster.
* The `reposure` CLI tool (it's also a `kubectl` plugin) can be used to access the registry via the
  surrogate. Likewise, if you need programmatic access to the registry from *outside* the cluster,
  Reposure provides a client API for the surrogate (Go language).
* If on the other hand you need programmatic access to the registry from *inside* the cluster
  (direct, without the surrogate), Reposure provides a client API based on
  [go-containerregistry](https://github.com/google/go-containerregistry) (Go language).
* Reposure has baked-in support for the built-in registries of OpenShift, OKD, CodeReady Containers,
  and Minikube.
* Reposure can deploy its own cloud-native registry, based on the Docker registry, which can useful
  for testing. (Note that by default it uses self-signed certificates, which may be inaccessible
  to your container runtime.)


How the Surrogate Works
-----------------------

The Reposure surrogate is emphatically *not* a proxy. A registry proxy would not in fact help you
access the registry from outside the cluster because it would just shift and might even multiply the
challenge. You would still need to connect to and secure the proxy. Tools like `buildah` and `skopeo`
can work in insecure mode, which might be fine for development, but otherwise you would still need
to set up TLS authentication and authorization.

Another challenge is that references within the image contain the registry in which they are stored,
e.g. `10.97.119.139:5000/myrepo/myimage:latest`. This means that the proxy would have to unpack the image,
rewrite the manifest, and repack it, and this would have to happen for both pushing and pulling.

The trick, then, is to use Kubernetes's existing control plane, which allows for executing commands
in containers as well as streaming stdout and stdin (via the SPDY procotol). Thus, instead of
dealing with images directly we will be using files, specifically tarball representations of the
images.

The surrogate is implemented as a combination of two tools:

* A file spooler, which watches a directory for incoming tarballs and pushes them to the registry.
  The spooler also handles deleting images from the registry via special filenames.
* A simple client utility, which can pull tarballs from the registry and deliver them to stdout.
  The utility can also list images in the registry, again delivering the list to stdout.

This combination allows you to push, pull, delete, and list images using nothing other than
`kubectl`. However, for convenience we provide the `reposure` CLI tool that makes these operations
even easier.

### Downsides

The problem with this approach is that if you are working outside the cluster then you would need to
use the `reposure` tool instead of your usual tools. So, for example, you can't use `buildah` to
directly push an image to the registry.

The workaround is often to export your image to a tarball and push that instead. Note that if you
want Kubernetes's container runtime to be able to pull it (say, for a Pod) then you would also need
to make sure to re-tag it accordingly. The `podman` tool can do this for you. For example, if we are
using `buildah` to build locally:

    # Build
    CONTAINER_ID=$(buildah from scratch)
    ...
    buildah commit $CONTAINER_ID localhost/myrepo/myimage
    
    # Re-tag
    HOST=$(reposure registry info myrepo host)
    podman tag localhost/myrepo/myimage $HOST/myrepo/myimage

    # Export and push
    podman save $HOST/myrepo/myimage --output myimage.tar
    reposure image push myregistry myrepo/myimage myimage.tar

(Note that spooler supports `.tar` as well as `.tar.gz` or `.tgz`.)

TODO: The `reposure` CLI currently does not block until the spooler succeeds or fails, and also there
is no forwarding of error messages.

### But, On the Other Hand...

The limitation of Reposure also comes with an advantage.

A side effect of the fact that the surrogate works directly with files is that it makes it very
easy to store arbitrary files on the registry. If you try to push a file that is not a tarball,
the spooler will create one for you (an image with a single layer). Likewise, when you pull a
tarball, the `reposure` tool can unpack that single layer for you. For example:

    echo 'hello world' > hello.txt
    reposure image push myregistry myrepo/hello hello.txt
    reposure image pull myregistry myrepo/hello --unpack


Terminology
===========

* *Registry*: This is the backend implementation, the actual server.
* *Repository*: The image reference structure comprises a repository name and an image name (as well
  as a "tag", whhich is usually used as the version). This extra naming level allows for namespace
  separation as well as permission management per repository. So, it is correct to say that the
  image is stored in a "repository" and it is also correct to say that it is stored in a "registry".
  Note that if you do not specify a repository name in the reference it internally defaults to
  "library" (and if you don't specify a tag it will default to "latest").
* *Simple*: Reposure can deploy a "simple" registry instance for you, based on the default Docker
  registry. This is intended for development and testing purposes but may be good enough for some
  production uses. Of course there exist more robust implementations, such as Quay and Harbor.


Basic Usage
-----------

Use Minikube's registry add-on (with "view" cluster role):

    kubectl config set-context --current --namespace=mynamespace
    reposure operator install --role=view --wait
    reposure registry create default --provider=minikube --wait

Use built-in registry in OpenShift or CodeReady Containers (with "view" cluster role):

    kubectl config set-context --current --namespace=mynamespace
    reposure operator install --role=view --wait
    reposure registry create default --provider=openshift --wait

Install the simple registry (for low-security clusters only, e.g. Minikube):

    kubectl config set-context --current --namespace=mynamespace
    reposure operator install --wait
    reposure simple install --wait
    reposure registry create default --provider=simple --wait

Quick test:

    echo 'Hello, world!' > hello.txt
    reposure image push default test hello.txt
    reposure image pull default test --unpack

For a fuller example that includes installing pushing and using an actual container image and also
installing the simple registry with authentication and authorization see [`lab/test`](lab/test).


FAQ
---

### Why is the surrogate implemented as a spooler? Why not just use a proxy?

An HTTPS proxy would be not make you life that much easier, as you would still need to handle TLS
certificates and domain names. The spooler allows you to forego all that and just use stdin/stdout
forwarding, which is built into the Kubernetes control plane.

### Why is it called "Reposure"?

"Reposure" is the state of being calm or relaxed, which is the ideal attitude for dealing with the
complexities of cloud-native registries.

Also, it's kinda short for "repository surrogate".
