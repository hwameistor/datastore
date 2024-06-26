//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// Code generated by deepcopy-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BaseModel) DeepCopyInto(out *BaseModel) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BaseModel.
func (in *BaseModel) DeepCopy() *BaseModel {
	if in == nil {
		return nil
	}
	out := new(BaseModel)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *BaseModel) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BaseModelList) DeepCopyInto(out *BaseModelList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]BaseModel, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BaseModelList.
func (in *BaseModelList) DeepCopy() *BaseModelList {
	if in == nil {
		return nil
	}
	out := new(BaseModelList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *BaseModelList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BaseModelSpec) DeepCopyInto(out *BaseModelSpec) {
	*out = *in
	if in.MinIO != nil {
		in, out := &in.MinIO, &out.MinIO
		*out = new(MinIOSpec)
		**out = **in
	}
	if in.NFS != nil {
		in, out := &in.NFS, &out.NFS
		*out = new(NFSSpec)
		**out = **in
	}
	if in.FTP != nil {
		in, out := &in.FTP, &out.FTP
		*out = new(FTPSpec)
		**out = **in
	}
	if in.HTTP != nil {
		in, out := &in.HTTP, &out.HTTP
		*out = new(HTTPSpec)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BaseModelSpec.
func (in *BaseModelSpec) DeepCopy() *BaseModelSpec {
	if in == nil {
		return nil
	}
	out := new(BaseModelSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BaseModelStatus) DeepCopyInto(out *BaseModelStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BaseModelStatus.
func (in *BaseModelStatus) DeepCopy() *BaseModelStatus {
	if in == nil {
		return nil
	}
	out := new(BaseModelStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Checkpoint) DeepCopyInto(out *Checkpoint) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Checkpoint.
func (in *Checkpoint) DeepCopy() *Checkpoint {
	if in == nil {
		return nil
	}
	out := new(Checkpoint)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Checkpoint) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CheckpointBackup) DeepCopyInto(out *CheckpointBackup) {
	*out = *in
	if in.MinIO != nil {
		in, out := &in.MinIO, &out.MinIO
		*out = new(MinIOSpec)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CheckpointBackup.
func (in *CheckpointBackup) DeepCopy() *CheckpointBackup {
	if in == nil {
		return nil
	}
	out := new(CheckpointBackup)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CheckpointList) DeepCopyInto(out *CheckpointList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Checkpoint, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CheckpointList.
func (in *CheckpointList) DeepCopy() *CheckpointList {
	if in == nil {
		return nil
	}
	out := new(CheckpointList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *CheckpointList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CheckpointRecord) DeepCopyInto(out *CheckpointRecord) {
	*out = *in
	if in.CreateTime != nil {
		in, out := &in.CreateTime, &out.CreateTime
		*out = (*in).DeepCopy()
	}
	if in.ExpiredTime != nil {
		in, out := &in.ExpiredTime, &out.ExpiredTime
		*out = (*in).DeepCopy()
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CheckpointRecord.
func (in *CheckpointRecord) DeepCopy() *CheckpointRecord {
	if in == nil {
		return nil
	}
	out := new(CheckpointRecord)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CheckpointSpec) DeepCopyInto(out *CheckpointSpec) {
	*out = *in
	if in.Backup != nil {
		in, out := &in.Backup, &out.Backup
		*out = new(CheckpointBackup)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CheckpointSpec.
func (in *CheckpointSpec) DeepCopy() *CheckpointSpec {
	if in == nil {
		return nil
	}
	out := new(CheckpointSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CheckpointStatus) DeepCopyInto(out *CheckpointStatus) {
	*out = *in
	if in.Records != nil {
		in, out := &in.Records, &out.Records
		*out = make([]*CheckpointRecord, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(CheckpointRecord)
				(*in).DeepCopyInto(*out)
			}
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CheckpointStatus.
func (in *CheckpointStatus) DeepCopy() *CheckpointStatus {
	if in == nil {
		return nil
	}
	out := new(CheckpointStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DataLoadRequest) DeepCopyInto(out *DataLoadRequest) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DataLoadRequest.
func (in *DataLoadRequest) DeepCopy() *DataLoadRequest {
	if in == nil {
		return nil
	}
	out := new(DataLoadRequest)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DataLoadRequest) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DataLoadRequestList) DeepCopyInto(out *DataLoadRequestList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]DataLoadRequest, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DataLoadRequestList.
func (in *DataLoadRequestList) DeepCopy() *DataLoadRequestList {
	if in == nil {
		return nil
	}
	out := new(DataLoadRequestList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DataLoadRequestList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DataLoadRequestSpec) DeepCopyInto(out *DataLoadRequestSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DataLoadRequestSpec.
func (in *DataLoadRequestSpec) DeepCopy() *DataLoadRequestSpec {
	if in == nil {
		return nil
	}
	out := new(DataLoadRequestSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DataLoadRequestStatus) DeepCopyInto(out *DataLoadRequestStatus) {
	*out = *in
	if in.ReadyNodes != nil {
		in, out := &in.ReadyNodes, &out.ReadyNodes
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DataLoadRequestStatus.
func (in *DataLoadRequestStatus) DeepCopy() *DataLoadRequestStatus {
	if in == nil {
		return nil
	}
	out := new(DataLoadRequestStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DataSet) DeepCopyInto(out *DataSet) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DataSet.
func (in *DataSet) DeepCopy() *DataSet {
	if in == nil {
		return nil
	}
	out := new(DataSet)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DataSet) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DataSetList) DeepCopyInto(out *DataSetList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]DataSet, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DataSetList.
func (in *DataSetList) DeepCopy() *DataSetList {
	if in == nil {
		return nil
	}
	out := new(DataSetList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DataSetList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DataSetSpec) DeepCopyInto(out *DataSetSpec) {
	*out = *in
	if in.MinIO != nil {
		in, out := &in.MinIO, &out.MinIO
		*out = new(MinIOSpec)
		**out = **in
	}
	if in.NFS != nil {
		in, out := &in.NFS, &out.NFS
		*out = new(NFSSpec)
		**out = **in
	}
	if in.FTP != nil {
		in, out := &in.FTP, &out.FTP
		*out = new(FTPSpec)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DataSetSpec.
func (in *DataSetSpec) DeepCopy() *DataSetSpec {
	if in == nil {
		return nil
	}
	out := new(DataSetSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DataSetStatus) DeepCopyInto(out *DataSetStatus) {
	*out = *in
	if in.LastRefreshTimestamp != nil {
		in, out := &in.LastRefreshTimestamp, &out.LastRefreshTimestamp
		*out = (*in).DeepCopy()
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DataSetStatus.
func (in *DataSetStatus) DeepCopy() *DataSetStatus {
	if in == nil {
		return nil
	}
	out := new(DataSetStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FTPSpec) DeepCopyInto(out *FTPSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FTPSpec.
func (in *FTPSpec) DeepCopy() *FTPSpec {
	if in == nil {
		return nil
	}
	out := new(FTPSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HTTPSpec) DeepCopyInto(out *HTTPSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HTTPSpec.
func (in *HTTPSpec) DeepCopy() *HTTPSpec {
	if in == nil {
		return nil
	}
	out := new(HTTPSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MinIOSpec) DeepCopyInto(out *MinIOSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MinIOSpec.
func (in *MinIOSpec) DeepCopy() *MinIOSpec {
	if in == nil {
		return nil
	}
	out := new(MinIOSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NFSSpec) DeepCopyInto(out *NFSSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NFSSpec.
func (in *NFSSpec) DeepCopy() *NFSSpec {
	if in == nil {
		return nil
	}
	out := new(NFSSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SSHSpec) DeepCopyInto(out *SSHSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SSHSpec.
func (in *SSHSpec) DeepCopy() *SSHSpec {
	if in == nil {
		return nil
	}
	out := new(SSHSpec)
	in.DeepCopyInto(out)
	return out
}
