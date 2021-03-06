//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// Copyright (c) 2020 Red Hat, Inc.

// Code generated by controller-gen. DO NOT EDIT.

package v1beta1

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ManagedClusterView) DeepCopyInto(out *ManagedClusterView) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ManagedClusterView.
func (in *ManagedClusterView) DeepCopy() *ManagedClusterView {
	if in == nil {
		return nil
	}
	out := new(ManagedClusterView)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ManagedClusterView) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ManagedClusterViewList) DeepCopyInto(out *ManagedClusterViewList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ManagedClusterView, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ManagedClusterViewList.
func (in *ManagedClusterViewList) DeepCopy() *ManagedClusterViewList {
	if in == nil {
		return nil
	}
	out := new(ManagedClusterViewList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ManagedClusterViewList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ViewScope) DeepCopyInto(out *ViewScope) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ViewScope.
func (in *ViewScope) DeepCopy() *ViewScope {
	if in == nil {
		return nil
	}
	out := new(ViewScope)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ViewSpec) DeepCopyInto(out *ViewSpec) {
	*out = *in
	out.Scope = in.Scope
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ViewSpec.
func (in *ViewSpec) DeepCopy() *ViewSpec {
	if in == nil {
		return nil
	}
	out := new(ViewSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ViewStatus) DeepCopyInto(out *ViewStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]v1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	in.Result.DeepCopyInto(&out.Result)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ViewStatus.
func (in *ViewStatus) DeepCopy() *ViewStatus {
	if in == nil {
		return nil
	}
	out := new(ViewStatus)
	in.DeepCopyInto(out)
	return out
}
