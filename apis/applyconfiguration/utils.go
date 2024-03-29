// Code generated by applyconfiguration-gen. DO NOT EDIT.

package applyconfiguration

import (
	reposurepuccinicloudv1alpha1 "github.com/tliron/reposure/apis/applyconfiguration/reposure.puccini.cloud/v1alpha1"
	v1alpha1 "github.com/tliron/reposure/resources/reposure.puccini.cloud/v1alpha1"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
)

// ForKind returns an apply configuration type for the given GroupVersionKind, or nil if no
// apply configuration type exists for the given GroupVersionKind.
func ForKind(kind schema.GroupVersionKind) interface{} {
	switch kind {
	// Group=reposure.puccini.cloud, Version=v1alpha1
	case v1alpha1.SchemeGroupVersion.WithKind("Registry"):
		return &reposurepuccinicloudv1alpha1.RegistryApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("RegistryDirect"):
		return &reposurepuccinicloudv1alpha1.RegistryDirectApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("RegistryIndirect"):
		return &reposurepuccinicloudv1alpha1.RegistryIndirectApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("RegistrySpec"):
		return &reposurepuccinicloudv1alpha1.RegistrySpecApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("RegistryStatus"):
		return &reposurepuccinicloudv1alpha1.RegistryStatusApplyConfiguration{}

	}
	return nil
}
