package admin

import (
	"fmt"

	certmanager "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	certmanagermeta "github.com/cert-manager/cert-manager/pkg/apis/meta/v1"
	"github.com/tliron/kutil/kubernetes"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const simpleHtpasswdMountPath = "/htpasswd"

var simpleTlsCertificatePath = fmt.Sprintf("%s/%s", tlsMountPath, core.TLSCertKey)
var simpleTlsKeyPath = fmt.Sprintf("%s/%s", tlsMountPath, core.TLSPrivateKeyKey)
var simpleHtpasswdPath = fmt.Sprintf("%s/htpasswd", simpleHtpasswdMountPath)

func (self *Client) InstallSimple(sourceRegistryHost string, authentication bool, authorization bool, wait bool) error {
	// Authentication requires Cert-Manager: https://github.com/cert-manager/cert-manager

	// Authorization expects a generic secret named "reposure-simple-htpasswd",
	// which is a file named htpasswd in htpasswd format that uses bcrypt for passwords
	// E.g.: htpasswd -cbB
	// See: https://docs.docker.com/registry/configuration/#htpasswd

	var err error

	if authentication {
		if err = self.GetCertManager(); err != nil {
			return err
		}
	}

	if sourceRegistryHost, err = self.GetSourceRegistryHost(sourceRegistryHost); err != nil {
		return err
	}

	var serviceAccount *core.ServiceAccount
	if serviceAccount, err = self.GetOperatorServiceAccount(); err != nil {
		return err
	}

	var registryDeployment *apps.Deployment
	if registryDeployment, err = self.createSimpleDeployment(sourceRegistryHost, serviceAccount, 1, authentication, authorization); err != nil {
		return err
	}

	var service *core.Service
	if service, err = self.createSimpleService(); err != nil {
		return err
	}

	if authentication {
		var issuer *certmanager.Issuer
		if issuer, err = self.createSimpleCertificateIssuer(); err != nil {
			return err
		}

		if _, err = self.createSimpleCertificate(issuer, service); err != nil {
			return err
		}
	}

	if wait {
		if _, err := kubernetes.WaitForDeployment(self.Context, self.Kubernetes, self.Log, self.Namespace, registryDeployment.Name); err != nil {
			return err
		}
	}

	return nil
}

func (self *Client) UninstallSimple(wait bool) {
	var gracePeriodSeconds int64 = 0
	deleteOptions := meta.DeleteOptions{
		GracePeriodSeconds: &gracePeriodSeconds,
	}

	appName := fmt.Sprintf("%s-simple", self.NamePrefix)
	secretName := fmt.Sprintf("%s-authentication", appName)

	// Service
	if err := self.Kubernetes.CoreV1().Services(self.Namespace).Delete(self.Context, appName, deleteOptions); err != nil {
		self.Log.Warningf("%s", err)
	}

	// Deployment
	if err := self.Kubernetes.AppsV1().Deployments(self.Namespace).Delete(self.Context, appName, deleteOptions); err != nil {
		self.Log.Warningf("%s", err)
	}

	if err := self.GetCertManager(); err != nil {
		self.Log.Warningf("%s", err.Error())
	}

	// Certificate
	if err := self.CertManager.CertmanagerV1().Certificates(self.Namespace).Delete(self.Context, appName, deleteOptions); err != nil {
		self.Log.Warningf("%s", err)
	}

	// Issuer
	if err := self.CertManager.CertmanagerV1().Issuers(self.Namespace).Delete(self.Context, appName, deleteOptions); err != nil {
		self.Log.Warningf("%s", err)
	}

	// Secret (deleting the Certificate will not delete the Secret!)
	if err := self.Kubernetes.CoreV1().Secrets(self.Namespace).Delete(self.Context, secretName, deleteOptions); err != nil {
		self.Log.Warningf("%s", err)
	}

	if wait {
		getOptions := meta.GetOptions{}
		kubernetes.WaitForDeletion(self.Log, "simple service", func() bool {
			_, err := self.Kubernetes.CoreV1().Services(self.Namespace).Get(self.Context, appName, getOptions)
			return err == nil
		})
		kubernetes.WaitForDeletion(self.Log, "simple deployment", func() bool {
			_, err := self.Kubernetes.AppsV1().Deployments(self.Namespace).Get(self.Context, appName, getOptions)
			return err == nil
		})
		kubernetes.WaitForDeletion(self.Log, "simple certificate", func() bool {
			_, err := self.CertManager.CertmanagerV1().Certificates(self.Namespace).Get(self.Context, appName, getOptions)
			return err == nil
		})
		kubernetes.WaitForDeletion(self.Log, "simple issuer", func() bool {
			_, err := self.CertManager.CertmanagerV1().Issuers(self.Namespace).Get(self.Context, appName, getOptions)
			return err == nil
		})
		kubernetes.WaitForDeletion(self.Log, "simple authentication secret", func() bool {
			_, err := self.Kubernetes.CoreV1().Secrets(self.Namespace).Get(self.Context, secretName, getOptions)
			return err == nil
		})
	}
}

func (self *Client) SimpleService() (*core.Service, error) {
	appName := fmt.Sprintf("%s-simple", self.NamePrefix)

	return self.Kubernetes.CoreV1().Services(self.Namespace).Get(self.Context, appName, meta.GetOptions{})
}

func (self *Client) SimpleHost() (string, error) {
	if service, err := self.SimpleService(); err == nil {
		return fmt.Sprintf("%s:5000", service.Spec.ClusterIP), nil
	} else {
		return "", err
	}
}

func (self *Client) createSimpleDeployment(registryAddress string, serviceAccount *core.ServiceAccount, replicas int32, authentication bool, authorization bool) (*apps.Deployment, error) {
	// https://hub.docker.com/_/registry
	// https://github.com/ContainerSolutions/trow
	// https://github.com/google/go-containerregistry

	appName := fmt.Sprintf("%s-simple", self.NamePrefix)
	labels := self.Labels(appName, "simple", self.Namespace)

	deployment := &apps.Deployment{
		ObjectMeta: meta.ObjectMeta{
			Name:   appName,
			Labels: labels,
		},
		Spec: apps.DeploymentSpec{
			Replicas: &replicas,
			Selector: &meta.LabelSelector{
				MatchLabels: labels,
			},
			Template: core.PodTemplateSpec{
				ObjectMeta: meta.ObjectMeta{
					Labels: labels,
				},
				Spec: core.PodSpec{
					Containers: []core.Container{
						{
							Name:            "registry",
							Image:           fmt.Sprintf("%s/%s", registryAddress, self.SimpleImageReference),
							ImagePullPolicy: core.PullAlways,
							VolumeMounts: []core.VolumeMount{
								{
									Name:      "registry",
									MountPath: "/var/lib/registry",
								},
							},
							Env: []core.EnvVar{
								{
									// necessary!
									Name:  "REGISTRY_STORAGE_DELETE_ENABLED",
									Value: "true",
								},
								{
									// For kutil's kubernetes.GetConfiguredNamespace
									Name: "KUBERNETES_NAMESPACE",
									ValueFrom: &core.EnvVarSource{
										FieldRef: &core.ObjectFieldSelector{
											FieldPath: "metadata.namespace",
										},
									},
								},
							},
							// Note: Probes skip certificate validation for HTTPS
							LivenessProbe: &core.Probe{
								ProbeHandler: core.ProbeHandler{
									HTTPGet: &core.HTTPGetAction{
										Port: intstr.FromInt(5000),
									},
								},
							},
							ReadinessProbe: &core.Probe{
								ProbeHandler: core.ProbeHandler{
									HTTPGet: &core.HTTPGetAction{
										Port: intstr.FromInt(5000),
									},
								},
							},
						},
					},
					Volumes: []core.Volume{
						{
							Name:         "registry",
							VolumeSource: self.VolumeSource("1Gi"),
						},
					},
				},
			},
		},
	}

	if authentication {
		secretName := fmt.Sprintf("%s-authentication", appName)

		deployment.Spec.Template.Spec.Containers[0].VolumeMounts = append(deployment.Spec.Template.Spec.Containers[0].VolumeMounts, core.VolumeMount{
			Name:      "tls",
			MountPath: tlsMountPath,
			ReadOnly:  true,
		})

		deployment.Spec.Template.Spec.Containers[0].Env = append(deployment.Spec.Template.Spec.Containers[0].Env,
			core.EnvVar{
				Name:  "REGISTRY_HTTP_TLS_CERTIFICATE",
				Value: simpleTlsCertificatePath,
			},
			core.EnvVar{
				Name:  "REGISTRY_HTTP_TLS_KEY",
				Value: simpleTlsKeyPath,
			},
		)

		deployment.Spec.Template.Spec.Containers[0].LivenessProbe.ProbeHandler.HTTPGet.Scheme = core.URISchemeHTTPS
		deployment.Spec.Template.Spec.Containers[0].ReadinessProbe.ProbeHandler.HTTPGet.Scheme = core.URISchemeHTTPS

		deployment.Spec.Template.Spec.Volumes = append(deployment.Spec.Template.Spec.Volumes, core.Volume{
			Name: "tls",
			VolumeSource: core.VolumeSource{
				Secret: &core.SecretVolumeSource{
					SecretName: secretName,
				},
			},
		})
	}

	if authorization {
		deployment.Spec.Template.Spec.Containers[0].VolumeMounts = append(deployment.Spec.Template.Spec.Containers[0].VolumeMounts, core.VolumeMount{
			Name:      "htpasswd",
			MountPath: simpleHtpasswdMountPath,
			ReadOnly:  true,
		})

		deployment.Spec.Template.Spec.Containers[0].Env = append(deployment.Spec.Template.Spec.Containers[0].Env,
			core.EnvVar{
				Name:  "REGISTRY_AUTH",
				Value: "htpasswd",
			},
			core.EnvVar{
				Name:  "REGISTRY_AUTH_HTPASSWD_PATH",
				Value: simpleHtpasswdPath,
			},
			core.EnvVar{
				Name:  "REGISTRY_AUTH_HTPASSWD_REALM",
				Value: "Registry",
			},
		)

		deployment.Spec.Template.Spec.Volumes = append(deployment.Spec.Template.Spec.Volumes, core.Volume{
			Name: "htpasswd",
			VolumeSource: core.VolumeSource{
				Secret: &core.SecretVolumeSource{
					SecretName: "reposure-simple-htpasswd",
				},
			},
		})
	}

	return self.CreateDeployment(deployment)
}

func (self *Client) createSimpleService() (*core.Service, error) {
	appName := fmt.Sprintf("%s-simple", self.NamePrefix)
	labels := self.Labels(appName, "simple", self.Namespace)

	service := &core.Service{
		ObjectMeta: meta.ObjectMeta{
			Name:   appName,
			Labels: labels,
		},
		Spec: core.ServiceSpec{
			Type:     core.ServiceTypeClusterIP,
			Selector: labels,
			Ports: []core.ServicePort{
				{
					Name:       "registry",
					Protocol:   "TCP",
					TargetPort: intstr.FromInt(5000),
					Port:       5000,
				},
			},
		},
	}

	return self.CreateService(service)
}

func (self *Client) createSimpleCertificateIssuer() (*certmanager.Issuer, error) {
	appName := fmt.Sprintf("%s-simple", self.NamePrefix)

	issuer := &certmanager.Issuer{
		ObjectMeta: meta.ObjectMeta{
			Name:   appName,
			Labels: self.Labels(appName, "simple", self.Namespace),
		},
		Spec: certmanager.IssuerSpec{
			IssuerConfig: certmanager.IssuerConfig{
				SelfSigned: &certmanager.SelfSignedIssuer{},
			},
		},
	}

	return self.CreateCertificateIssuer(issuer)
}

func (self *Client) createSimpleCertificate(issuer *certmanager.Issuer, service *core.Service) (*certmanager.Certificate, error) {
	appName := fmt.Sprintf("%s-simple", self.NamePrefix)
	secretName := fmt.Sprintf("%s-authentication", appName)
	ipAddress := service.Spec.ClusterIP

	certificate := &certmanager.Certificate{
		ObjectMeta: meta.ObjectMeta{
			Name:   appName,
			Labels: self.Labels(appName, "simple", self.Namespace),
		},
		Spec: certmanager.CertificateSpec{
			SecretName:  secretName,
			IPAddresses: []string{ipAddress},
			URIs:        []string{"https://reposure.puccini.cloud"},
			IssuerRef: certmanagermeta.ObjectReference{
				Name: issuer.Name,
			},
		},
	}

	return self.CreateCertificate(certificate)
}

/*
func (self *Client) createSimpleConfigMap() (*core.ConfigMap, error) {
	appName := fmt.Sprintf("%s-simple", self.NamePrefix)
	instanceName := fmt.Sprintf("%s-%s", appName, self.Namespace)

	configMap := &core.ConfigMap{
		ObjectMeta: meta.ObjectMeta{
			Name: appName,
			Labels: map[string]string{
				"app.kubernetes.io/name":       appName,
				"app.kubernetes.io/instance":   instanceName,
				"app.kubernetes.io/version":    version.GitVersion,
				"app.kubernetes.io/component":  "registry",
				"app.kubernetes.io/part-of":    self.PartOf,
				"app.kubernetes.io/managed-by": self.ManagedBy,
			},
		},
	}

	if configMap, err := self.Kubernetes.CoreV1().ConfigMaps(self.Namespace).Create(self.Context, configMap, meta.CreateOptions{}); err == nil {
		return configMap, nil
	} else if errors.IsAlreadyExists(err) {
		self.Log.Infof("%s", err.Error())
		return self.Kubernetes.CoreV1().ConfigMaps(self.Namespace).Get(self.Context, appName, meta.GetOptions{})
	} else {
		return nil, err
	}
}

func (self *Client) createRegistryImagePullSecret(server string, username string, password string) (*core.Secret, error) {
	// See: https://kubernetes.io/docs/concepts/containers/images/#specifying-imagepullsecrets-on-a-pod
	//      https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/
	//      https://docs.docker.com/engine/reference/commandline/cli/#configjson-properties

	appName := fmt.Sprintf("%s-simple", self.NamePrefix)
	instanceName := fmt.Sprintf("%s-%s", appName, self.Namespace)

	secret := &core.Secret{
		ObjectMeta: meta.ObjectMeta{
			Name: appName,
			Labels: map[string]string{
				"app.kubernetes.io/name":       appName,
				"app.kubernetes.io/instance":   instanceName,
				"app.kubernetes.io/version":    version.GitVersion,
				"app.kubernetes.io/component":  "registry",
				"app.kubernetes.io/part-of":    self.PartOf,
				"app.kubernetes.io/managed-by": self.ManagedBy,
			},
		},
	}

	if err := kubernetes.SetSecretDockerConfigJson(secret, server, username, password); err != nil {
		return nil, err
	}

	if secret, err := self.Kubernetes.CoreV1().Secrets(self.Namespace).Create(self.Context, secret, meta.CreateOptions{}); err == nil {
		return secret, nil
	} else if errors.IsAlreadyExists(err) {
		self.Log.Infof("%s", err.Error())
		return self.Kubernetes.CoreV1().Secrets(self.Namespace).Get(self.Context, appName, meta.GetOptions{})
	} else {
		return nil, err
	}
}

// See: https://nip.io/
//      https://cert-manager.io/docs/

func (self *Client) createRegistryTlsSecret() (*core.Secret, error) {
	appName := fmt.Sprintf("%s-simple", self.NamePrefix)
	instanceName := fmt.Sprintf("%s-%s", appName, self.Namespace)

	var crt []byte
	var key []byte

	secret := &core.Secret{
		ObjectMeta: meta.ObjectMeta{
			Name: appName,
			Labels: map[string]string{
				"app.kubernetes.io/name":       appName,
				"app.kubernetes.io/instance":   instanceName,
				"app.kubernetes.io/version":    version.GitVersion,
				"app.kubernetes.io/component":  "registry",
				"app.kubernetes.io/part-of":    self.PartOf,
				"app.kubernetes.io/managed-by": self.ManagedBy,
			},
		},
		Type: core.SecretTypeTLS,
		Data: map[string][]byte{
			core.TLSCertKey:       crt,
			core.TLSPrivateKeyKey: key,
		},
	}

	if secret, err := self.Kubernetes.CoreV1().Secrets(self.Namespace).Create(self.Context, secret, meta.CreateOptions{}); err == nil {
		return secret, nil
	} else if errors.IsAlreadyExists(err) {
		self.Log.Infof("%s", err.Error())
		return self.Kubernetes.CoreV1().Secrets(self.Namespace).Get(self.Context, appName, meta.GetOptions{})
	} else {
		return nil, err
	}
}
*/
