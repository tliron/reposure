*This is an early release. Some features are not yet fully implemented.*

Reposure
========

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Latest Release](https://img.shields.io/github/release/tliron/reposure.svg)](https://github.com/tliron/reposure/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/tliron/reposure)](https://goreportcard.com/report/github.com/tliron/reposure)

Manage and access cloud-native container image registries for [Kubernetes](https://kubernetes.io/).

Reposure has baked-in support for the built-in registries of OpenShift and Minikube. Additionally, it
can deploy its own "simple" cloud-native registry, based on the reference Docker registry, which can
be useful for testing and development. (Note that by default it uses self-signed certificates, which
may be inaccessible to your container runtime.) For more robust implementations, see
[Harbor](https://goharbor.io/) and [Quay](https://www.projectquay.io/).


Get It
------

[![Download](assets/media/download.png "Download")](https://github.com/tliron/reposure/releases)


Rationale
---------

### What are cloud-native registries?

Kubernetes's container runtime, whether it's CRI-O or Docker or something else, pulls its container
images from any OCI-compliant container image registry. Though there are publicly hosted registries,
such as [Docker Hub](https://hub.docker.com/) and [Quay](https://quay.io/), it is often desirable and
necessary to host a private registry.

With that in mind, it can makes sense to use Kubernetes to host the private registry. Moreover, though
the private registry can run on any Kubernetes cluster, it may also make sense to deploy it in *the same*
Kubernetes cluster that will be using it. This setup is what we're here calling a "cloud-native registry".

This setup can significantly simplify the deployment as there is no need to create or access an additional
cluster. Moreover, it may be possible to have applications deploy their own custom, private registries.
An obvious use case is cloud-native CI/CD pipelines, for which building and packaging code, publishing
images, and deploying containers all happen within the cluster. Indeed, OpenShift/OKD and Minikube come
with built-in cloud-native registries exactly for this purpose.

### The challenge

Setting up cloud-native registries can be quite challenging. Your private TLS certificates, certificate
authorities, and authorization credentials, must all be configured into the cluster's container runtime.
This is in addition to the requirement that the container runtime, which runs on the host, can route to
the registry's IP address, which in this case sits on Kubernetes's control plane. The container runtime
might also need access to the registry's domain name for TLS verification, which could require a DNS
relay, and `/etc/hosts` update, or similar.

Additionally, beyond just getting the cloud-native registry to work with the container runtime, it can
be challenging to access it from outside the cluster using tools like Buildah and Skopeo. Would an
ingress to the cluster be required? An ingress is non-trivial and sometimes impossible to set up. And
what about TLS authentication and authorization from the outside?


Reposure's Features
-------------------

Reposure can assist with many of the challenges mentioned abvove.

The Reposure operator manages registry "surrogates", which run as pods in the same cluster as the
registries. They are configured with the necessary TLS authentication and authorization, if required.
The surrogates can fetch and push container images for you. No ingress is required, just normal
`kubectl` access to the cluster's API server. The `reposure` CLI tool (it's also a `kubectl` plugin)
simplifies access to these surrogates.

Furthermore, Reposure provides APIs with which your application can access these registries:

* The ["surrogate" client](client/surrogate/) allows access to the surrogate from outside the
  cluster, and is essentially what the `reposure` CLI tool uses.
* The ["direct" client](client/direct/) allows programmatic direct access to the registry from
  *inside* the cluster, *not* via the surrogate but rather via the
  [go-containerregistry](https://github.com/google/go-containerregistry) library. Reposure will
  handle configuring the client with the authentication and authorization.


How the Surrogate Works
-----------------------

The Reposure surrogate is deliberately *not* a proxy. A registry proxy would not in fact help you
access the registry from outside the cluster because it would just shift the challenge and might
even makes things more difficult. Your outside tools, such as [`buildah`](https://buildah.io/) and
[`skopeo`](https://github.com/containers/skopeo), would still need to securely connect to that
proxy. Terminating TLS would also be challenging.

The alternative we chose is to use Kubernetes's existing control plane, which allows for executing
commands in containers as well as streaming stdout and stdin (via the SPDY procotol). We can use
this to transfer files (tarballs of container images) to and from the cluster, similarly to how
`kubectl cp` works. Thus, the surrogate functions as a
[jump server](https://en.wikipedia.org/wiki/Jump_server).

The Reposure surrogate comprises of two components:

* A [file spooler](reposure-registry-spooler), which watches a directory for incoming tarballs and
  pushes them to the registry. The spooler also handles deleting images from the registry via
  special filenames.
* A [client utility](reposure-registry-client), which can pull tarballs from the registry and
  deliver them to stdout. The utility can also list images in the registry, again delivering the
  list to stdout.

Together these two components allow you to push, pull, delete, and list images using the basic
`kubectl` connectivity you already have.

### Downsides

The problem with this solution is that if you are working outside the cluster then you would need
to use the `reposure` tool instead of your usual tools. So, for example, you can't use `buildah`
to directly push an image to the registry. (Wouldn't it be nice if `buildah` had built-in support
for Reposure?)

The workaround is to export your image to a tarball and push that instead. Note, though, that you
would also need to re-tag it so that Kubernetes's container runtime can pull it. You can do all of
this using the `podman` tool. Here's an example workflow:

    # Build
    CONTAINER_ID=$(buildah from scratch)
    ...
    buildah commit $CONTAINER_ID localhost/myrepo/myimage

    # Re-tag
    HOST=$(reposure registry info myregistry host)
    podman tag localhost/myrepo/myimage $HOST/myrepo/myimage

    # Export
    podman save $HOST/myrepo/myimage --output myimage.tar

    # Push
    reposure image push myregistry myrepo/myimage myimage.tar

(Note that Reposure accepts `.tar` as well as `.tar.gz` or `.tgz`.)

### But, On the Other Hand...

This downside also comes with a very useful advantage.

A side effect of the fact that the surrogate works directly with files is that it makes it very
easy to store arbitrary files in the registry, not just container images. You don't even have to
package them as tarballs: the spooler will automatically wrap the file in a tarball for you if
it's not one already. (This is essentially a container image with a single layer). Likewise, when
you pull a tarball, the `reposure` tool can automatically unpack that single layer for you. For
example:

    echo 'hello world' > hello.txt
    reposure image push myregistry myrepo/hello hello.txt
    reposure image pull myregistry myrepo/hello --unpack


Installation
------------

Use Minikube's registry add-on (with "view" cluster role):

    kubectl config set-context --current --namespace=mynamespace
    reposure operator install --role=view --wait
    reposure registry create default --provider=minikube --wait

Use built-in registry in OpenShift (with "view" cluster role):

    kubectl config set-context --current --namespace=mynamespace
    reposure operator install --role=view --wait
    reposure registry create default --provider=openshift --wait

Install the simple registry (for low-security clusters only, e.g. Minikube):

    kubectl config set-context --current --namespace=mynamespace
    reposure operator install --wait
    reposure simple install --wait
    reposure registry create default --provider=simple --wait

For a fuller example that includes installing, pushing, and using an actual container image, as
well as installing the simple registry with authentication and authorization, see
[`lab/test`](lab/test).


FAQ
---

### What's the difference between a "registry" and a "repository"?

A **registry** is the backend implementation, the actual server.

The image reference structure comprises a **repository** name and an image name (as well as a
"tag", which is usually used as the version). This extra naming level allows for namespace
separation as well as permission management per repository. Note that if you do not specify a
repository name in the reference it internally defaults to "library" (and if you don't specify a
tag it will default to "latest").

So, it is correct to say that the image is stored in a "repository" and it is also correct to say
that it is stored in a "registry".

### Why is OpenShift giving me access errors for pods using images from the built-in registry?

OpenShift's added security requires the repository name and namespace of the pod to be identical.
This improves isolation between namespaces: a namespace can't pull images that belong to another
namespace.

### Why is it called "Reposure"?

"Reposure" is the state of being calm or relaxed. It is recomended to stay calm and relaxed when
dealing with the complexities of cloud-native registries...

Also, it's kinda short for "**repo**sitory **sur**rogate".
