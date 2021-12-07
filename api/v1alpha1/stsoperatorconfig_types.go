/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	helmerv1beta1 "github.com/openshift-psap/special-resource-operator/pkg/helmer/api/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// StsOperatorConfigSpec defines the desired state of StsOperatorConfig
type StsOperatorConfigSpec struct {

	// +kubebuilder:default:int32=50051
	GrpcSvcPort int `json:"grpcSvcPort,omitempty"`

	// +kubebuilder:default:int32=2947
	GpsSvcPort int `json:"gpsSvcPort,omitempty"`

	// +kubebuilder:default:string="sts-silicom"
	Namespace string `json:"namespace,omitempty"`

	Images StsImages `json:"images,omitempty"`
}

type SroCfg struct {
	Build bool `json:"build,omitempty"`

	Chart helmerv1beta1.HelmChart `json:"chart,omitempty"`

	Namespace string `json:"namespace,omitempty"`
}

type StsImages struct {

	// +kubebuilder:default:string="quay.io/silicom/tsyncd:2.0.1.0"
	Tsyncd string `json:"tsyncd,omitempty"`

	// +kubebuilder:default:string="quay.io/silicom/grpc-tsyncd:2.0.1.0"
	GrpcTsyncd string `json:"grpcTsyncd,omitempty"`

	// +kubebuilder:default:string="quay.io/silicom/sts-plugin:2.0.1.0"
	StsPlugin string `json:"stsPlugin,omitempty"`

	// +kubebuilder:default:string="quay.io/silicom/gpsd:3.23.1"
	Gpsd string `json:"gpsd,omitempty"`

	// +kubebuilder:default:string="quay.io/silicom/tsync_extts:1.0.0"
	TsyncExtts string `json:"tsyncExtts,omitempty"`

	// +kubebuilder:default:string="quay.io/silicom/phc2sys:3.1.1"
	Phc2sys string `json:"phc2sys,omitempty"`
}

// StsOperatorConfigStatus defines the observed state of StsOperatorConfig
type StsOperatorConfigStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// StsOperatorConfig is the Schema for the stsoperatorconfigs API
type StsOperatorConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   StsOperatorConfigSpec   `json:"spec,omitempty"`
	Status StsOperatorConfigStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// StsOperatorConfigList contains a list of StsOperatorConfig
type StsOperatorConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []StsOperatorConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&StsOperatorConfig{}, &StsOperatorConfigList{})
}
