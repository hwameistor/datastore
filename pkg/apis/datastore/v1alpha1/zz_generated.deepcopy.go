//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// Code generated by deepcopy-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StorageBackend) DeepCopyInto(out *StorageBackend) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StorageBackend.
func (in *StorageBackend) DeepCopy() *StorageBackend {
	if in == nil {
		return nil
	}
	out := new(StorageBackend)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *StorageBackend) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StorageBackendList) DeepCopyInto(out *StorageBackendList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]StorageBackend, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StorageBackendList.
func (in *StorageBackendList) DeepCopy() *StorageBackendList {
	if in == nil {
		return nil
	}
	out := new(StorageBackendList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *StorageBackendList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StorageBackendSpec) DeepCopyInto(out *StorageBackendSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StorageBackendSpec.
func (in *StorageBackendSpec) DeepCopy() *StorageBackendSpec {
	if in == nil {
		return nil
	}
	out := new(StorageBackendSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StorageBackendStatus) DeepCopyInto(out *StorageBackendStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StorageBackendStatus.
func (in *StorageBackendStatus) DeepCopy() *StorageBackendStatus {
	if in == nil {
		return nil
	}
	out := new(StorageBackendStatus)
	in.DeepCopyInto(out)
	return out
}
