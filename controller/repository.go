package controller

import (
	resources "github.com/tliron/reposure/resources/reposure.puccini.cloud/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
)

func (self *Controller) UpdateRepositorySurrogatePod(repository *resources.Repository, spoolerPod string) (*resources.Repository, error) {
	self.Log.Infof("updating spooler pod to %q for repository: %s/%s", spoolerPod, repository.Namespace, repository.Name)

	for {
		repository = repository.DeepCopy()
		repository.Status.SpoolerPod = spoolerPod

		service_, err, retry := self.updateRepositoryStatus(repository)
		if retry {
			repository = service_
		} else {
			return service_, err
		}
	}
}
func (self *Controller) updateRepositoryStatus(repository *resources.Repository) (*resources.Repository, error, bool) {
	if repository_, err := self.Client.UpdateRepositoryStatus(repository); err == nil {
		return repository_, nil, false
	} else if errors.IsConflict(err) {
		self.Log.Warningf("retrying status update for repository: %s/%s", repository.Namespace, repository.Name)
		if repository_, err := self.Client.GetRepository(repository.Namespace, repository.Name); err == nil {
			return repository_, nil, true
		} else {
			return repository, err, false
		}
	} else {
		return repository, err, false
	}
}

func (self *Controller) processRepository(repository *resources.Repository) (bool, error) {
	// Create surrogate
	if pod, err := self.Client.CreateRepositorySurrogate(repository); err == nil {
		if _, err := self.UpdateRepositorySurrogatePod(repository, pod.Name); err == nil {
			return true, nil
		} else {
			return false, err
		}
	} else {
		return false, err
	}
}
