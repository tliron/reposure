package controller

import (
	resources "github.com/tliron/reposure/resources/reposure.puccini.cloud/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
)

func (self *Controller) UpdateRegistrySurrogatePod(registry *resources.Registry, surrogatePod string) (*resources.Registry, error) {
	self.Log.Infof("updating surrogate pod to %q for registry: %s/%s", surrogatePod, registry.Namespace, registry.Name)

	for {
		registry = registry.DeepCopy()
		registry.Status.SurrogatePod = surrogatePod

		registry_, err, retry := self.updateRegistryStatus(registry)
		if retry {
			registry = registry_
		} else {
			return registry_, err
		}
	}
}

func (self *Controller) updateRegistryStatus(registry *resources.Registry) (*resources.Registry, error, bool) {
	if registry_, err := self.Client.UpdateRegistryStatus(registry); err == nil {
		return registry_, nil, false
	} else if errors.IsConflict(err) {
		self.Log.Warningf("retrying status update for registry: %s/%s", registry.Namespace, registry.Name)
		if registry_, err := self.Client.GetRegistry(registry.Namespace, registry.Name); err == nil {
			return registry_, nil, true
		} else {
			return registry, err, false
		}
	} else {
		return registry, err, false
	}
}

func (self *Controller) processRegistry(registry *resources.Registry) (bool, error) {
	// Create surrogate
	if pod, err := self.Client.CreateRegistrySurrogate(registry); err == nil {
		if _, err := self.UpdateRegistrySurrogatePod(registry, pod.Name); err == nil {
			return true, nil
		} else {
			return false, err
		}
	} else {
		return false, err
	}
}
