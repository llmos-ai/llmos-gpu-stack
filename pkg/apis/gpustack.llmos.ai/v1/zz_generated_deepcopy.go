//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2024 llmos.ai.

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
// Code generated by main. DO NOT EDIT.

package v1

import (
	common "github.com/llmos-ai/llmos-gpu-stack/pkg/apis/common"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GPUDevice) DeepCopyInto(out *GPUDevice) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GPUDevice.
func (in *GPUDevice) DeepCopy() *GPUDevice {
	if in == nil {
		return nil
	}
	out := new(GPUDevice)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *GPUDevice) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GPUDeviceInfo) DeepCopyInto(out *GPUDeviceInfo) {
	*out = *in
	if in.Index != nil {
		in, out := &in.Index, &out.Index
		*out = new(int)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GPUDeviceInfo.
func (in *GPUDeviceInfo) DeepCopy() *GPUDeviceInfo {
	if in == nil {
		return nil
	}
	out := new(GPUDeviceInfo)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GPUDeviceList) DeepCopyInto(out *GPUDeviceList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]GPUDevice, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GPUDeviceList.
func (in *GPUDeviceList) DeepCopy() *GPUDeviceList {
	if in == nil {
		return nil
	}
	out := new(GPUDeviceList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *GPUDeviceList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GPUDeviceSpec) DeepCopyInto(out *GPUDeviceSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GPUDeviceSpec.
func (in *GPUDeviceSpec) DeepCopy() *GPUDeviceSpec {
	if in == nil {
		return nil
	}
	out := new(GPUDeviceSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GPUDeviceStatus) DeepCopyInto(out *GPUDeviceStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]common.Condition, len(*in))
		copy(*out, *in)
	}
	in.GPUDeviceInfo.DeepCopyInto(&out.GPUDeviceInfo)
	if in.Pods != nil {
		in, out := &in.Pods, &out.Pods
		*out = make([]GPUPod, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GPUDeviceStatus.
func (in *GPUDeviceStatus) DeepCopy() *GPUDeviceStatus {
	if in == nil {
		return nil
	}
	out := new(GPUDeviceStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GPUPod) DeepCopyInto(out *GPUPod) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GPUPod.
func (in *GPUPod) DeepCopy() *GPUPod {
	if in == nil {
		return nil
	}
	out := new(GPUPod)
	in.DeepCopyInto(out)
	return out
}
