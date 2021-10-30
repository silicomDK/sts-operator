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
	// +kubebuilder:validation:Pattern=[a-z0-9\.\-]+
	Name string `json:"name"`

	Interfaces   []StsInterfaceSpec `json:"interfaces"`
	NodeSelector map[string]string  `json:"nodeSelector,omitempty"`

	ImageRegistry string `json:"imageRegistry"`
	Mode          string `json:"mode"`
	Namespace     string `json:"namespace"`
}

type STSNodeStatus struct {
	Name        string      `json:"name"`
	TsyncStatus TsyncStatus `json:"tsyncStatus"`
	GpsStatus   GPSStatus   `json:"gpsStatus,omitempty"`
}

type TsyncStatus struct {
	Mode   string `json:"mode"`
	Status string `json:"status"`
}

type GPSStatus struct {
	Status string `json:"status"`
}

type StsInterfaceSpec struct {
	EthName string `json:"ethName"`
	SyncE   bool   `json:"synce"`
	HoldOff int    `json:"holdoff"`
}

// StsConfigStatus defines the observed state of StsConfig
type StsConfigStatus struct {
	NodeStatus []STSNodeStatus `json:"nodeStatus"`
	State      string          `json:"state"`
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
