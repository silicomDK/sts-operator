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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// StsConfigSpec defines the desired state of StsConfig
type StsConfigSpec struct {
	Interfaces   []StsInterfaceSpec `json:"interfaces"`
	NodeSelector map[string]string  `json:"nodeSelector,omitempty"`

	// +kubebuilder:default:string=quay.io/silicom
	ImageRegistry string `json:"imageRegistry,omitempty"`

	// +kubebuilder:default:string="2.0.0.0"
	StsVersion string `json:"stsVersion,omitempty"`

	// +kubebuilder:validation:Enum=T-GM.8275.1;T-BC-8275.1;T-TSC.8275.1
	// +kubebuilder:default:="T-GM.8275.1"
	Mode string `json:"mode,omitempty"`

	// +kubebuilder:default:="sts-silicom"
	Namespace string `json:"namespace,omitempty"`

	// +kubebuilder:default:=24
	// +kubebuilder:validation:Minimum=24
	// +kubebuilder:validation:Maximum=48
	DomainNumber int `json:"domainNumber,omitempty"`

	// +kubebuilder:default:=1
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=3
	SrcPPS int `json:"srcPPS,omitempty"`

	// +kubebuilder:default:=1
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=3
	Src10MHz int `json:"src10MHz,omitempty"`

	// +kubebuilder:default:=0
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=7
	SynceRecClkPort int `json:"synceRecClkPort,omitempty"`

	// +kubebuilder:default:=0
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=1
	PhyLedsCtl int `json:"phyLedsCtl,omitempty"`

	// +kubebuilder:default:=1
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=2
	SyncOption int `json:"synceOption,omitempty"`

	// +kubebuilder:default:=10
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=40
	SynceCpu int `json:"synceCpu,omitempty"`
}

type STSNodeStatus struct {
	Name        string      `json:"name,omitempty"`
	TsyncStatus TsyncStatus `json:"tsyncStatus,omitempty"`
	GpsStatus   GPSStatus   `json:"gpsStatus,omitempty"`
}

type TsyncStatus struct {
	Mode   string `json:"mode"`
	Status string `json:"status"`
}

type GPSStatus struct {
	Time string `json:"time"`
	Lat  int    `json:"lat"`
	Lon  int    `json:"lon"`
}

type StsInterfaceSpec struct {
	EthName string `json:"ethName"`
	EthPort int    `json:"ethPort"`

	// +kubebuilder:default:=0
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=1
	SyncE int `json:"synce,omitempty"`

	// +kubebuilder:default:=500
	// +kubebuilder:validation:Minimum=300
	// +kubebuilder:validation:Maximum=1800
	HoldOff int `json:"holdoff,omitempty"`

	// +kubebuilder:validation:Enum=Master;Slave
	// +kubebuilder:default:=Master
	Mode string `json:"mode,omitempty"`

	// +kubebuilder:default:=0
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=1
	QlEnable int `json:"qlEnable,omitempty"`
}

// StsConfigStatus defines the observed state of StsConfig
type StsConfigStatus struct {
	NodeStatus []STSNodeStatus `json:"nodeStatus,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// StsConfig is the Schema for the stsconfigs API
type StsConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	// +kubebuilder:validation:Required
	Spec   StsConfigSpec   `json:"spec,omitempty"`
	Status StsConfigStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// StsConfigList contains a list of StsConfig
type StsConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []StsConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&StsConfig{}, &StsConfigList{})
}
