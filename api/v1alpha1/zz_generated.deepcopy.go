//go:build !ignore_autogenerated
// +build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GPSStatus) DeepCopyInto(out *GPSStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GPSStatus.
func (in *GPSStatus) DeepCopy() *GPSStatus {
	if in == nil {
		return nil
	}
	out := new(GPSStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SroCfg) DeepCopyInto(out *SroCfg) {
	*out = *in
	in.Chart.DeepCopyInto(&out.Chart)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SroCfg.
func (in *SroCfg) DeepCopy() *SroCfg {
	if in == nil {
		return nil
	}
	out := new(SroCfg)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StsConfig) DeepCopyInto(out *StsConfig) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StsConfig.
func (in *StsConfig) DeepCopy() *StsConfig {
	if in == nil {
		return nil
	}
	out := new(StsConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *StsConfig) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StsConfigList) DeepCopyInto(out *StsConfigList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]StsConfig, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StsConfigList.
func (in *StsConfigList) DeepCopy() *StsConfigList {
	if in == nil {
		return nil
	}
	out := new(StsConfigList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *StsConfigList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StsConfigSpec) DeepCopyInto(out *StsConfigSpec) {
	*out = *in
	if in.Interfaces != nil {
		in, out := &in.Interfaces, &out.Interfaces
		*out = make([]StsInterfaceSpec, len(*in))
		copy(*out, *in)
	}
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StsConfigSpec.
func (in *StsConfigSpec) DeepCopy() *StsConfigSpec {
	if in == nil {
		return nil
	}
	out := new(StsConfigSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StsConfigStatus) DeepCopyInto(out *StsConfigStatus) {
	*out = *in
	if in.Nodes != nil {
		in, out := &in.Nodes, &out.Nodes
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StsConfigStatus.
func (in *StsConfigStatus) DeepCopy() *StsConfigStatus {
	if in == nil {
		return nil
	}
	out := new(StsConfigStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StsImages) DeepCopyInto(out *StsImages) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StsImages.
func (in *StsImages) DeepCopy() *StsImages {
	if in == nil {
		return nil
	}
	out := new(StsImages)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StsInterfaceSpec) DeepCopyInto(out *StsInterfaceSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StsInterfaceSpec.
func (in *StsInterfaceSpec) DeepCopy() *StsInterfaceSpec {
	if in == nil {
		return nil
	}
	out := new(StsInterfaceSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StsNode) DeepCopyInto(out *StsNode) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StsNode.
func (in *StsNode) DeepCopy() *StsNode {
	if in == nil {
		return nil
	}
	out := new(StsNode)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *StsNode) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StsNodeInterfaceStatus) DeepCopyInto(out *StsNodeInterfaceStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StsNodeInterfaceStatus.
func (in *StsNodeInterfaceStatus) DeepCopy() *StsNodeInterfaceStatus {
	if in == nil {
		return nil
	}
	out := new(StsNodeInterfaceStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StsNodeList) DeepCopyInto(out *StsNodeList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]StsNode, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StsNodeList.
func (in *StsNodeList) DeepCopy() *StsNodeList {
	if in == nil {
		return nil
	}
	out := new(StsNodeList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *StsNodeList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StsNodeSpec) DeepCopyInto(out *StsNodeSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StsNodeSpec.
func (in *StsNodeSpec) DeepCopy() *StsNodeSpec {
	if in == nil {
		return nil
	}
	out := new(StsNodeSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StsNodeStatus) DeepCopyInto(out *StsNodeStatus) {
	*out = *in
	out.TsyncStatus = in.TsyncStatus
	out.GpsStatus = in.GpsStatus
	if in.EthInterfaces != nil {
		in, out := &in.EthInterfaces, &out.EthInterfaces
		*out = make([]StsNodeInterfaceStatus, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StsNodeStatus.
func (in *StsNodeStatus) DeepCopy() *StsNodeStatus {
	if in == nil {
		return nil
	}
	out := new(StsNodeStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StsOperatorConfig) DeepCopyInto(out *StsOperatorConfig) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StsOperatorConfig.
func (in *StsOperatorConfig) DeepCopy() *StsOperatorConfig {
	if in == nil {
		return nil
	}
	out := new(StsOperatorConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *StsOperatorConfig) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StsOperatorConfigList) DeepCopyInto(out *StsOperatorConfigList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]StsOperatorConfig, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StsOperatorConfigList.
func (in *StsOperatorConfigList) DeepCopy() *StsOperatorConfigList {
	if in == nil {
		return nil
	}
	out := new(StsOperatorConfigList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *StsOperatorConfigList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StsOperatorConfigSpec) DeepCopyInto(out *StsOperatorConfigSpec) {
	*out = *in
	out.Images = in.Images
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StsOperatorConfigSpec.
func (in *StsOperatorConfigSpec) DeepCopy() *StsOperatorConfigSpec {
	if in == nil {
		return nil
	}
	out := new(StsOperatorConfigSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StsOperatorConfigStatus) DeepCopyInto(out *StsOperatorConfigStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StsOperatorConfigStatus.
func (in *StsOperatorConfigStatus) DeepCopy() *StsOperatorConfigStatus {
	if in == nil {
		return nil
	}
	out := new(StsOperatorConfigStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TsyncStatus) DeepCopyInto(out *TsyncStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TsyncStatus.
func (in *TsyncStatus) DeepCopy() *TsyncStatus {
	if in == nil {
		return nil
	}
	out := new(TsyncStatus)
	in.DeepCopyInto(out)
	return out
}
