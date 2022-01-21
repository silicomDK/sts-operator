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

	// +kubebuilder:validation:Enum=T-GM.8275.1;T-BC-8275.1;T-TSC.8275.1
	// +kubebuilder:default:="T-GM.8275.1"
	//Telecom G8275 Profile
	//
	// T-BC-8275.1 (default)
	//
	// T-GM.8275.1
	//
	// T-TSC.8275.1
	Mode string `json:"mode,omitempty"`

	// +kubebuilder:default:="sts-silicom"
	// +kubebuilder:validation:Optional
	Namespace string `json:"namespace,omitempty"`

	// +kubebuilder:default:=24
	// +kubebuilder:validation:Minimum=24
	// +kubebuilder:validation:Maximum=48
	// +kubebuilder:validation:Optional
	// PTP domain number
	DomainNumber int `json:"domainNumber,omitempty"`

	// +kubebuilder:default:=2
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=2
	// +kubebuilder:validation:Optional
	//Set 1PPS Connector Mode
	//
	//1 - PPS IN
	//
	//2 - PPS OUT (default)
	//
	ModePPS int `json:"modePPS,omitempty"`

	// +kubebuilder:default:=1
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=3
	// +kubebuilder:validation:Optional
	// Set PPS OUT Source
	//
	// 1 - PLL (default)
	//
	// 2 - GPS
	//
	// 3 - IN
	SrcPPS int `json:"srcPPS,omitempty"`

	// +kubebuilder:default:=2
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=3
	// +kubebuilder:validation:Optional
	// Set 10MHz Connector Mode
	//
	// 1 - 10MHz IN
	//
	// 2 - 10MHz OUT (default)
	//
	// 3 - PPS OUT
	Mode10MHz int `json:"mode10MHz,omitempty"`

	// +kubebuilder:default:=1
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=3
	// +kubebuilder:validation:Optional
	// Set 10MHz OUT Source
	//
	// 1 - PLL (default)
	//
	// 2 - GPS
	//
	// 3 - IN
	Src10MHz int `json:"src10MHz,omitempty"`

	// +kubebuilder:default:=0
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=7
	// +kubebuilder:validation:Optional
	// Set SyncE Recovery Clock Port
	SynceRecClkPort int `json:"synceRecClkPort,omitempty"`

	// +kubebuilder:default:=0
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=1
	// +kubebuilder:validation:Optional
	//Enable/disable Phy Leds Control Switch
	//
	//0 - disable Phy Leds Control Switch (default)
	//
	//1 - enable Phy Leds Control Switch
	PhyLedsCtl int `json:"phyLedsCtl,omitempty"`

	// +kubebuilder:default:=1
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=2
	// +kubebuilder:validation:Optional
	// Configures the synchronization network
	//
	//1 - Option 1 refers to synchronization networks designed for Europe
	//
	//2 - Option 2 refers to synchronization networks designed for United States
	SyncOption int `json:"syncOption,omitempty"`

	// +kubebuilder:default:=10
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=40
	// +kubebuilder:validation:Optional
	// Set CPU Pin for SyncE ESMC thread
	SynceCpu int `json:"synceCpu,omitempty"`

	// +kubebuilder:default:=0
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=1
	// +kubebuilder:validation:Optional
	//Enable/disable two-step PTP Clock
	//
	//0 - Disable two-step clock, (set one-step clock) (default)
	//
	//1 - Enable two-step clock
	TwoStep int `json:"twoStep,omitempty"`

	// +kubebuilder:default:=128
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=255
	// +kubebuilder:validation:Optional
	//Set Priority 2 for GM PTP Clock
	//
	// Valid range 0-255, smaller values indicate higher priority
	Priority2 int `json:"priority2,omitempty"`

	// +kubebuilder:default:=0
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=1
	// +kubebuilder:validation:Optional
	//Forwardable/Non-Forwardable Multicast Address
	//
	//0 - Non-Forwardable (default)
	//
	//1 - Forwardable
	Forwardable int `json:"forwardable,omitempty"`

	// +kubebuilder:default:=-1
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=10
	// +kubebuilder:validation:Optional
	//Trace PTP Message
	//
	//Valid range -1-10
	//
	//-1 - Disable Trace log for PTP Messages (default)
	//
	//0 - Enable Trace for all types of PTP Messages
	//
	//1 - Enable Trace for SYNC Messages
	//
	//2 - Enable Trace for DELAY_REQ Messages
	//
	//3 - Enable Trace for PEER_DELAY_REQ Messages
	//
	//4 - Enable Trace for PEER_DELAY_RESP Messages
	//
	//5 - Enable Trace for FOLLOW_UP Messages
	//
	//6 - Enable Trace for DELAY_RESP Messages
	//
	//7 - Enable Trace for PEER_DELAY_FOLLOW_UP Messages
	//
	//8 - Enable Trace for ANNOUNCE Messages
	//
	//9 - Enable Trace for SIGNAL Messages
	//
	//10 - Enable Trace for MANAGEMENT Messages
	TracePtpMsg int `json:"tracePtpMsg,omitempty"`

	// +kubebuilder:default:=0
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=8
	// +kubebuilder:validation:Optional
	TraceLevel int `json:"traceLevel,omitempty"`

	// +kubebuilder:default:=2
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=2
	// +kubebuilder:validation:Optional
	//Configures the ESMC Mode
	//
	//1 - Manual
	//
	//2 - Auto (default)
	EsmcMode int `json:"esmcMode,omitempty"`

	// +kubebuilder:default:=1
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=2
	// +kubebuilder:validation:Optional
	//Configures the SSM Mode
	//
	//1 - SSM Code (default)
	//
	//2 - ESSM Code
	SsmMode int `json:"ssmMode,omitempty"`

	// +kubebuilder:default:=0
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=8
	// +kubebuilder:validation:Optional
	AprLevel int `json:"aprLevel,omitempty"`
}

type StsInterfaceSpec struct {
	EthName string `json:"ethName"`
	EthPort int    `json:"ethPort"`

	// +kubebuilder:default:=0
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=1
	// +kubebuilder:validation:Optional
	SyncE int `json:"synce,omitempty"`

	// +kubebuilder:default:=500
	// +kubebuilder:validation:Minimum=300
	// +kubebuilder:validation:Maximum=1800
	// +kubebuilder:validation:Optional
	HoldOff int `json:"holdoff,omitempty"`

	// +kubebuilder:validation:Enum=Master;Slave
	// +kubebuilder:default:=Master
	// +kubebuilder:validation:Optional
	Mode string `json:"mode,omitempty"`

	// +kubebuilder:default:=0
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=1
	// +kubebuilder:validation:Optional
	// Configures QL enable for the interface
	// 0 - Disable QL
	// 1 - Enable QL (default)
	QlEnable int `json:"qlEnable,omitempty"`

	// +kubebuilder:default:=4
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=7
	// +kubebuilder:validation:Optional
	//Configures QL value for the interface
	//
	//if syncOption is 1 (Europe)
	//
	//===========================
	//
	//1 - QL-PRC
	//
	//2 - QL-PRTC
	//
	//3 - QL-EEC1
	//
	//4 - QL-DNU (default)
	//
	//if syncOption is 2 (United States)
	//
	//==================================
	//
	//5 - QL-PRS
	//
	//2 - QL-PRTC
	//
	//6 - QL-EEC2
	//
	//7 - QL-DUS (default)
	Ql int `json:"ql,omitempty"`
}

// StsConfigStatus defines the observed state of StsConfig
type StsConfigStatus struct {
	Nodes []string `json:"nodes,omitempty"`
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
