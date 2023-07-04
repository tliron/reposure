// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	json "encoding/json"
	"fmt"
	"time"

	reposurepuccinicloudv1alpha1 "github.com/tliron/reposure/apis/applyconfiguration/reposure.puccini.cloud/v1alpha1"
	scheme "github.com/tliron/reposure/apis/clientset/versioned/scheme"
	v1alpha1 "github.com/tliron/reposure/resources/reposure.puccini.cloud/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// RegistriesGetter has a method to return a RegistryInterface.
// A group's client should implement this interface.
type RegistriesGetter interface {
	Registries(namespace string) RegistryInterface
}

// RegistryInterface has methods to work with Registry resources.
type RegistryInterface interface {
	Create(ctx context.Context, registry *v1alpha1.Registry, opts v1.CreateOptions) (*v1alpha1.Registry, error)
	Update(ctx context.Context, registry *v1alpha1.Registry, opts v1.UpdateOptions) (*v1alpha1.Registry, error)
	UpdateStatus(ctx context.Context, registry *v1alpha1.Registry, opts v1.UpdateOptions) (*v1alpha1.Registry, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.Registry, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.RegistryList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.Registry, err error)
	Apply(ctx context.Context, registry *reposurepuccinicloudv1alpha1.RegistryApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.Registry, err error)
	ApplyStatus(ctx context.Context, registry *reposurepuccinicloudv1alpha1.RegistryApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.Registry, err error)
	RegistryExpansion
}

// registries implements RegistryInterface
type registries struct {
	client rest.Interface
	ns     string
}

// newRegistries returns a Registries
func newRegistries(c *ReposureV1alpha1Client, namespace string) *registries {
	return &registries{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the registry, and returns the corresponding registry object, and an error if there is any.
func (c *registries) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.Registry, err error) {
	result = &v1alpha1.Registry{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("registries").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Registries that match those selectors.
func (c *registries) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.RegistryList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.RegistryList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("registries").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested registries.
func (c *registries) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("registries").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a registry and creates it.  Returns the server's representation of the registry, and an error, if there is any.
func (c *registries) Create(ctx context.Context, registry *v1alpha1.Registry, opts v1.CreateOptions) (result *v1alpha1.Registry, err error) {
	result = &v1alpha1.Registry{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("registries").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(registry).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a registry and updates it. Returns the server's representation of the registry, and an error, if there is any.
func (c *registries) Update(ctx context.Context, registry *v1alpha1.Registry, opts v1.UpdateOptions) (result *v1alpha1.Registry, err error) {
	result = &v1alpha1.Registry{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("registries").
		Name(registry.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(registry).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *registries) UpdateStatus(ctx context.Context, registry *v1alpha1.Registry, opts v1.UpdateOptions) (result *v1alpha1.Registry, err error) {
	result = &v1alpha1.Registry{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("registries").
		Name(registry.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(registry).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the registry and deletes it. Returns an error if one occurs.
func (c *registries) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("registries").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *registries) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("registries").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched registry.
func (c *registries) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.Registry, err error) {
	result = &v1alpha1.Registry{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("registries").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}

// Apply takes the given apply declarative configuration, applies it and returns the applied registry.
func (c *registries) Apply(ctx context.Context, registry *reposurepuccinicloudv1alpha1.RegistryApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.Registry, err error) {
	if registry == nil {
		return nil, fmt.Errorf("registry provided to Apply must not be nil")
	}
	patchOpts := opts.ToPatchOptions()
	data, err := json.Marshal(registry)
	if err != nil {
		return nil, err
	}
	name := registry.Name
	if name == nil {
		return nil, fmt.Errorf("registry.Name must be provided to Apply")
	}
	result = &v1alpha1.Registry{}
	err = c.client.Patch(types.ApplyPatchType).
		Namespace(c.ns).
		Resource("registries").
		Name(*name).
		VersionedParams(&patchOpts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}

// ApplyStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating ApplyStatus().
func (c *registries) ApplyStatus(ctx context.Context, registry *reposurepuccinicloudv1alpha1.RegistryApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.Registry, err error) {
	if registry == nil {
		return nil, fmt.Errorf("registry provided to Apply must not be nil")
	}
	patchOpts := opts.ToPatchOptions()
	data, err := json.Marshal(registry)
	if err != nil {
		return nil, err
	}

	name := registry.Name
	if name == nil {
		return nil, fmt.Errorf("registry.Name must be provided to Apply")
	}

	result = &v1alpha1.Registry{}
	err = c.client.Patch(types.ApplyPatchType).
		Namespace(c.ns).
		Resource("registries").
		Name(*name).
		SubResource("status").
		VersionedParams(&patchOpts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
