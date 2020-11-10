package admin

import (
	"fmt"

	"github.com/tliron/kutil/kubernetes"
	resources "github.com/tliron/reposure/resources/reposure.puccini.cloud/v1alpha1"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const spoolPath = "/spool"

const surrogateContainerName = "surrogate"

func (self *Client) CreateRepositorySurrogate(repository *resources.Repository) (*core.Pod, error) {
	repositoryClient := self.RepositoryClient(repository)

	var repositoryHost string
	repositoryHost, err := repositoryClient.GetHost(repository)
	if err != nil {
		return nil, err
	}

	registryHost := "docker.io"
	appName := self.GetRepositorySurrogateAppName(repository.Name)

	pod := &core.Pod{
		ObjectMeta: meta.ObjectMeta{
			Name:      appName,
			Namespace: repository.Namespace,
			Labels:    self.Labels(appName, "surrogate", repository.Namespace),
		},
		Spec: core.PodSpec{
			Containers: []core.Container{
				{
					Name:            surrogateContainerName,
					Image:           fmt.Sprintf("%s/%s", registryHost, self.RepositorySurrogateImageReference),
					ImagePullPolicy: core.PullAlways,
					VolumeMounts: []core.VolumeMount{
						{
							Name:      "spool",
							MountPath: spoolPath,
						},
					},
					Env: []core.EnvVar{
						{
							Name:  "REPOSURE_REGISTRY_SPOOLER_registry",
							Value: repositoryHost,
						},
						{
							Name:  "REPOSURE_REGISTRY_SPOOLER_verbose",
							Value: "2",
						},
					},
					LivenessProbe: &core.Probe{
						Handler: core.Handler{
							HTTPGet: &core.HTTPGetAction{
								Port: intstr.FromInt(8086),
								Path: "/live",
							},
						},
					},
					ReadinessProbe: &core.Probe{
						Handler: core.Handler{
							HTTPGet: &core.HTTPGetAction{
								Port: intstr.FromInt(8086),
								Path: "/ready",
							},
						},
					},
				},
			},
			Volumes: []core.Volume{
				{
					Name:         "spool",
					VolumeSource: self.VolumeSource("1Gi"),
				},
			},
		},
	}

	if repository.Spec.TLSSecret != "" {
		pod.Spec.Containers[0].VolumeMounts = append(pod.Spec.Containers[0].VolumeMounts, core.VolumeMount{
			Name:      "tls",
			MountPath: tlsMountPath,
			ReadOnly:  true,
		})

		pod.Spec.Containers[0].Env = append(pod.Spec.Containers[0].Env, core.EnvVar{
			Name:  "REPOSURE_REGISTRY_SPOOLER_certificate",
			Value: repositoryClient.GetCertificatePath(repository),
		})

		pod.Spec.Volumes = append(pod.Spec.Volumes, core.Volume{
			Name: "tls",
			VolumeSource: core.VolumeSource{
				Secret: &core.SecretVolumeSource{
					SecretName: repository.Spec.TLSSecret,
				},
			},
		})
	}

	if _, username, password, token, err := repositoryClient.GetAuth(repository); err == nil {
		if username != "" {
			pod.Spec.Containers[0].Env = append(pod.Spec.Containers[0].Env, core.EnvVar{
				Name:  "REPOSURE_REGISTRY_SPOOLER_username",
				Value: username,
			})
		}
		if password != "" {
			pod.Spec.Containers[0].Env = append(pod.Spec.Containers[0].Env, core.EnvVar{
				Name:  "REPOSURE_REGISTRY_SPOOLER_password",
				Value: password,
			})
		}
		if token != "" {
			pod.Spec.Containers[0].Env = append(pod.Spec.Containers[0].Env, core.EnvVar{
				Name:  "REPOSURE_REGISTRY_SPOOLER_token",
				Value: token,
			})
		}
	} else {
		return nil, err
	}

	ownerReferences := pod.GetOwnerReferences()
	ownerReferences = append(ownerReferences, *meta.NewControllerRef(repository, repository.GroupVersionKind()))
	pod.SetOwnerReferences(ownerReferences)

	return self.CreatePod(pod)
}

func (self *Client) WaitForRepositorySurrogate(namespace string, repositoryName string) (*core.Pod, error) {
	appName := self.GetRepositorySurrogateAppName(repositoryName)
	return kubernetes.WaitForPod(self.Context, self.Kubernetes, self.Log, namespace, appName)
}

func (self *Client) GetRepositorySurrogateAppName(repositoryName string) string {
	return fmt.Sprintf("%s-surrogate-%s", self.NamePrefix, repositoryName)
}
