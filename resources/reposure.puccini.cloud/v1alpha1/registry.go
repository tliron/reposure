package v1alpha1

import (
	"fmt"

	"github.com/tliron/kutil/ard"
	"github.com/tliron/kutil/kubernetes"
	group "github.com/tliron/reposure/resources/reposure.puccini.cloud"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var RegistryGVK = SchemeGroupVersion.WithKind(RegistryKind)

type RegistryType string

const (
	RegistryKind     = "Registry"
	RegistryListKind = "RegistryList"

	RegistrySingular  = "registry"
	RegistryPlural    = "registries"
	RegistryShortName = "reg"

	RegistryTypeOCI RegistryType = "oci"
)

//
// Registry
//

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Registry struct {
	meta.TypeMeta   `json:",inline"`
	meta.ObjectMeta `json:"metadata,omitempty"`

	Spec   RegistrySpec   `json:"spec"`
	Status RegistryStatus `json:"status"`
}

type RegistrySpec struct {
	Type                        RegistryType      `json:"type"`                                  // Registry type
	Direct                      *RegistryDirect   `json:"direct,omitempty"`                      // Direct reference to registry
	Indirect                    *RegistryIndirect `json:"indirect,omitempty"`                    // Indirect reference to registry
	AuthenticationSecret        string            `json:"authenticationSecret,omitempty"`        // Name of authentication Secret required for connecting to the registry (optional)
	AuthenticationSecretDataKey string            `json:"authenticationSecretDataKey,omitempty"` // Name of key within the authentication Secret data required for connecting to the registry (optional)
	AuthorizationSecret         string            `json:"authorizationSecret,omitempty"`         // Name of authorization Secret required for connecting to the registry (optional)
}

type RegistryDirect struct {
	Host string `json:"host"` // Registry host (either "host:port" or "host")
}

type RegistryIndirect struct {
	Namespace string `json:"namespace,omitempty"` // Namespace for service resource (optional; defaults to same namespace as this registry)
	Service   string `json:"service"`             // Name of service resource
	Port      uint64 `json:"port"`                // TCP port to use with service
}

type RegistryStatus struct {
	SurrogatePod string `json:"surrogatePod"` // Name of surrogate pod resource (in the same namespace as this registry)
}

//
// RegistryList
//

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type RegistryList struct {
	meta.TypeMeta `json:",inline"`
	meta.ListMeta `json:"metadata"`

	Items []Registry `json:"items"`
}

//
// RegistryCustomResourceDefinition
//

// See: assets/custom-resource-definitions.yaml

var RegistryResourcesName = fmt.Sprintf("%s.%s", RegistryPlural, group.GroupName)

var RegistryCustomResourceDefinition = apiextensions.CustomResourceDefinition{
	ObjectMeta: meta.ObjectMeta{
		Name: RegistryResourcesName,
	},
	Spec: apiextensions.CustomResourceDefinitionSpec{
		Group: group.GroupName,
		Names: apiextensions.CustomResourceDefinitionNames{
			Singular: RegistrySingular,
			Plural:   RegistryPlural,
			Kind:     RegistryKind,
			ListKind: RegistryListKind,
			ShortNames: []string{
				RegistryShortName,
			},
			Categories: []string{
				"all", // will appear in "kubectl get all"
			},
		},
		Scope: apiextensions.NamespaceScoped,
		Versions: []apiextensions.CustomResourceDefinitionVersion{
			{
				Name:    Version,
				Served:  true,
				Storage: true, // one and only one version must be marked with storage=true
				Subresources: &apiextensions.CustomResourceSubresources{ // requires CustomResourceSubresources feature gate enabled
					Status: &apiextensions.CustomResourceSubresourceStatus{},
				},
				Schema: &apiextensions.CustomResourceValidation{
					OpenAPIV3Schema: &apiextensions.JSONSchemaProps{
						Description: "Reposure registry",
						Type:        "object",
						Required:    []string{"spec"},
						Properties: map[string]apiextensions.JSONSchemaProps{
							"spec": {
								Type:     "object",
								Required: []string{"type"},
								Properties: map[string]apiextensions.JSONSchemaProps{
									"type": {
										Description: "Registry type",
										Type:        "string",
										Enum: []apiextensions.JSON{
											kubernetes.JSONString(RegistryTypeOCI),
										},
									},
									"direct": {
										Description: "Direct reference to registry",
										Type:        "object",
										Required:    []string{"host"},
										Properties: map[string]apiextensions.JSONSchemaProps{
											"host": {
												Description: "Registry host (either \"host:port\" or \"host\")",
												Type:        "string",
											},
										},
									},
									"indirect": {
										Description: "Indirect reference to registry",
										Type:        "object",
										Required:    []string{"service", "port"},
										Properties: map[string]apiextensions.JSONSchemaProps{
											"namespace": {
												Description: "Namespace for service resource (optional; defaults to same namespace as this registry)",
												Type:        "string",
											},
											"service": {
												Description: "Name of service resource",
												Type:        "string",
											},
											"port": {
												Description: "TCP port to use with service",
												Type:        "integer",
											},
										},
									},
									"authenticationSecret": {
										Description: "Name of authentication Secret required for connecting to the registry (optional)",
										Type:        "string",
									},
									"authenticationSecretDataKey": {
										Description: "Name of key within the authentication Secret data required for connecting to the registry (optional)",
										Type:        "string",
									},
									"authorizationSecret": {
										Description: "Name of authorization Secret required for connecting to the registry (optional)",
										Type:        "string",
									},
								},
								OneOf: []apiextensions.JSONSchemaProps{
									{
										Required: []string{"direct"},
									},
									{
										Required: []string{"indirect"},
									},
								},
							},
							"status": {
								Type: "object",
								Properties: map[string]apiextensions.JSONSchemaProps{
									"surrogatePod": {
										Description: "Name of surrogate pod resource (in the same namespace as this registry)",
										Type:        "string",
									},
								},
							},
						},
					},
				},
				AdditionalPrinterColumns: []apiextensions.CustomResourceColumnDefinition{
					{
						Name:     "Type",
						Type:     "string",
						JSONPath: ".spec.type",
					},
					{
						Name:     "SurrogatePod",
						Type:     "string",
						JSONPath: ".status.surrogatePod",
					},
				},
			},
		},
	},
}

func RegistryToARD(registry *Registry) ard.StringMap {
	map_ := make(ard.StringMap)
	map_["Name"] = registry.Name
	map_["Type"] = registry.Spec.Type
	if registry.Spec.Direct != nil {
		map_["Direct"] = ard.StringMap{
			"Host": registry.Spec.Direct.Host,
		}
	} else if registry.Spec.Indirect != nil {
		map_["Indirect"] = ard.StringMap{
			"Namespace": registry.Spec.Indirect.Namespace,
			"Service":   registry.Spec.Indirect.Service,
			"Port":      registry.Spec.Indirect.Port,
		}
	}
	map_["AuthenticationSecret"] = registry.Spec.AuthenticationSecret
	map_["AuthenticationSecretDataKey"] = registry.Spec.AuthenticationSecretDataKey
	map_["AuthSecret"] = registry.Spec.AuthorizationSecret
	map_["SurrogatePod"] = registry.Status.SurrogatePod
	return map_
}
