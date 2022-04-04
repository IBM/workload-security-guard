//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright The Kubernetes Authors.

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

// Code generated by deepcopy-gen. DO NOT EDIT.

package v1

import (
	net "net"

	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CidrSet.
func (in CidrSet) DeepCopy() CidrSet {
	if in == nil {
		return nil
	}
	out := new(CidrSet)
	in.DeepCopyInto(out)
	return *out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Critiria) DeepCopyInto(out *Critiria) {
	*out = *in
	in.Req.DeepCopyInto(&out.Req)
	in.Resp.DeepCopyInto(&out.Resp)
	in.Process.DeepCopyInto(&out.Process)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Critiria.
func (in *Critiria) DeepCopy() *Critiria {
	if in == nil {
		return nil
	}
	out := new(Critiria)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Ctrl) DeepCopyInto(out *Ctrl) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Ctrl.
func (in *Ctrl) DeepCopy() *Ctrl {
	if in == nil {
		return nil
	}
	out := new(Ctrl)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Guardian) DeepCopyInto(out *Guardian) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	if in.Spec != nil {
		in, out := &in.Spec, &out.Spec
		*out = new(GuardianSpec)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Guardian.
func (in *Guardian) DeepCopy() *Guardian {
	if in == nil {
		return nil
	}
	out := new(Guardian)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Guardian) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GuardianList) DeepCopyInto(out *GuardianList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Guardian, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GuardianList.
func (in *GuardianList) DeepCopy() *GuardianList {
	if in == nil {
		return nil
	}
	out := new(GuardianList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *GuardianList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GuardianSpec) DeepCopyInto(out *GuardianSpec) {
	*out = *in
	if in.Configured != nil {
		in, out := &in.Configured, &out.Configured
		*out = new(Critiria)
		(*in).DeepCopyInto(*out)
	}
	if in.Learned != nil {
		in, out := &in.Learned, &out.Learned
		*out = new(Critiria)
		(*in).DeepCopyInto(*out)
	}
	out.Control = in.Control
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GuardianSpec.
func (in *GuardianSpec) DeepCopy() *GuardianSpec {
	if in == nil {
		return nil
	}
	out := new(GuardianSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HeadersConfig) DeepCopyInto(out *HeadersConfig) {
	*out = *in
	in.Kv.DeepCopyInto(&out.Kv)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HeadersConfig.
func (in *HeadersConfig) DeepCopy() *HeadersConfig {
	if in == nil {
		return nil
	}
	out := new(HeadersConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HeadersPile) DeepCopyInto(out *HeadersPile) {
	*out = *in
	if in.Kv != nil {
		in, out := &in.Kv, &out.Kv
		*out = new(KeyValPile)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HeadersPile.
func (in *HeadersPile) DeepCopy() *HeadersPile {
	if in == nil {
		return nil
	}
	out := new(HeadersPile)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HeadersProfile) DeepCopyInto(out *HeadersProfile) {
	*out = *in
	if in.Kv != nil {
		in, out := &in.Kv, &out.Kv
		*out = new(KeyValProfile)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HeadersProfile.
func (in *HeadersProfile) DeepCopy() *HeadersProfile {
	if in == nil {
		return nil
	}
	out := new(HeadersProfile)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IpPile) DeepCopyInto(out *IpPile) {
	*out = *in
	if in.List != nil {
		in, out := &in.List, &out.List
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.m != nil {
		in, out := &in.m, &out.m
		*out = make(map[string]bool, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IpPile.
func (in *IpPile) DeepCopy() *IpPile {
	if in == nil {
		return nil
	}
	out := new(IpPile)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IpSet) DeepCopyInto(out *IpSet) {
	*out = *in
	if in.list != nil {
		in, out := &in.list, &out.list
		*out = make([]net.IP, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = make(net.IP, len(*in))
				copy(*out, *in)
			}
		}
	}
	if in.m != nil {
		in, out := &in.m, &out.m
		*out = make(map[string]bool, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IpSet.
func (in *IpSet) DeepCopy() *IpSet {
	if in == nil {
		return nil
	}
	out := new(IpSet)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KeyValConfig) DeepCopyInto(out *KeyValConfig) {
	*out = *in
	if in.Vals != nil {
		in, out := &in.Vals, &out.Vals
		*out = make(map[string]*SimpleValConfig, len(*in))
		for key, val := range *in {
			var outVal *SimpleValConfig
			if val == nil {
				(*out)[key] = nil
			} else {
				in, out := &val, &outVal
				*out = new(SimpleValConfig)
				(*in).DeepCopyInto(*out)
			}
			(*out)[key] = outVal
		}
	}
	if in.OtherVals != nil {
		in, out := &in.OtherVals, &out.OtherVals
		*out = new(SimpleValConfig)
		(*in).DeepCopyInto(*out)
	}
	if in.OtherKeynames != nil {
		in, out := &in.OtherKeynames, &out.OtherKeynames
		*out = new(SimpleValConfig)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KeyValConfig.
func (in *KeyValConfig) DeepCopy() *KeyValConfig {
	if in == nil {
		return nil
	}
	out := new(KeyValConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KeyValPile) DeepCopyInto(out *KeyValPile) {
	*out = *in
	if in.Vals != nil {
		in, out := &in.Vals, &out.Vals
		*out = make(map[string]*SimpleValPile, len(*in))
		for key, val := range *in {
			var outVal *SimpleValPile
			if val == nil {
				(*out)[key] = nil
			} else {
				in, out := &val, &outVal
				*out = new(SimpleValPile)
				(*in).DeepCopyInto(*out)
			}
			(*out)[key] = outVal
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KeyValPile.
func (in *KeyValPile) DeepCopy() *KeyValPile {
	if in == nil {
		return nil
	}
	out := new(KeyValPile)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KeyValProfile) DeepCopyInto(out *KeyValProfile) {
	*out = *in
	if in.Vals != nil {
		in, out := &in.Vals, &out.Vals
		*out = make(map[string]*SimpleValProfile, len(*in))
		for key, val := range *in {
			var outVal *SimpleValProfile
			if val == nil {
				(*out)[key] = nil
			} else {
				in, out := &val, &outVal
				*out = new(SimpleValProfile)
				(*in).DeepCopyInto(*out)
			}
			(*out)[key] = outVal
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KeyValProfile.
func (in *KeyValProfile) DeepCopy() *KeyValProfile {
	if in == nil {
		return nil
	}
	out := new(KeyValProfile)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Pile) DeepCopyInto(out *Pile) {
	*out = *in
	in.Req.DeepCopyInto(&out.Req)
	in.Resp.DeepCopyInto(&out.Resp)
	in.Process.DeepCopyInto(&out.Process)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Pile.
func (in *Pile) DeepCopy() *Pile {
	if in == nil {
		return nil
	}
	out := new(Pile)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ProcessConfig) DeepCopyInto(out *ProcessConfig) {
	*out = *in
	out.tcp4Peers = in.tcp4Peers.DeepCopy()
	out.udp4Peers = in.udp4Peers.DeepCopy()
	out.udplite4Peers = in.udplite4Peers.DeepCopy()
	out.tcp6Peers = in.tcp6Peers.DeepCopy()
	out.udp6Peers = in.udp6Peers.DeepCopy()
	out.udplite6Peers = in.udplite6Peers.DeepCopy()
	if in.ResponseTime != nil {
		in, out := &in.ResponseTime, &out.ResponseTime
		*out = make(U8MinmaxSlice, len(*in))
		copy(*out, *in)
	}
	if in.CompletionTime != nil {
		in, out := &in.CompletionTime, &out.CompletionTime
		*out = make(U8MinmaxSlice, len(*in))
		copy(*out, *in)
	}
	if in.Tcp4Peers != nil {
		in, out := &in.Tcp4Peers, &out.Tcp4Peers
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Udp4Peers != nil {
		in, out := &in.Udp4Peers, &out.Udp4Peers
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Udplite4Peers != nil {
		in, out := &in.Udplite4Peers, &out.Udplite4Peers
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Tcp6Peers != nil {
		in, out := &in.Tcp6Peers, &out.Tcp6Peers
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Udp6Peers != nil {
		in, out := &in.Udp6Peers, &out.Udp6Peers
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Udplite6Peers != nil {
		in, out := &in.Udplite6Peers, &out.Udplite6Peers
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ProcessConfig.
func (in *ProcessConfig) DeepCopy() *ProcessConfig {
	if in == nil {
		return nil
	}
	out := new(ProcessConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ProcessPile) DeepCopyInto(out *ProcessPile) {
	*out = *in
	if in.ResponseTime != nil {
		in, out := &in.ResponseTime, &out.ResponseTime
		*out = make([]byte, len(*in))
		copy(*out, *in)
	}
	if in.CompletionTime != nil {
		in, out := &in.CompletionTime, &out.CompletionTime
		*out = make([]byte, len(*in))
		copy(*out, *in)
	}
	in.Tcp4Peers.DeepCopyInto(&out.Tcp4Peers)
	in.Udp4Peers.DeepCopyInto(&out.Udp4Peers)
	in.Udplite4Peers.DeepCopyInto(&out.Udplite4Peers)
	in.Tcp6Peers.DeepCopyInto(&out.Tcp6Peers)
	in.Udp6Peers.DeepCopyInto(&out.Udp6Peers)
	in.Udplite6Peers.DeepCopyInto(&out.Udplite6Peers)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ProcessPile.
func (in *ProcessPile) DeepCopy() *ProcessPile {
	if in == nil {
		return nil
	}
	out := new(ProcessPile)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ProcessProfile) DeepCopyInto(out *ProcessProfile) {
	*out = *in
	if in.Tcp4Peers != nil {
		in, out := &in.Tcp4Peers, &out.Tcp4Peers
		*out = new(IpSet)
		(*in).DeepCopyInto(*out)
	}
	if in.Udp4Peers != nil {
		in, out := &in.Udp4Peers, &out.Udp4Peers
		*out = new(IpSet)
		(*in).DeepCopyInto(*out)
	}
	if in.Udplite4Peers != nil {
		in, out := &in.Udplite4Peers, &out.Udplite4Peers
		*out = new(IpSet)
		(*in).DeepCopyInto(*out)
	}
	if in.Tcp6Peers != nil {
		in, out := &in.Tcp6Peers, &out.Tcp6Peers
		*out = new(IpSet)
		(*in).DeepCopyInto(*out)
	}
	if in.Udp6Peers != nil {
		in, out := &in.Udp6Peers, &out.Udp6Peers
		*out = new(IpSet)
		(*in).DeepCopyInto(*out)
	}
	if in.Udplite6Peers != nil {
		in, out := &in.Udplite6Peers, &out.Udplite6Peers
		*out = new(IpSet)
		(*in).DeepCopyInto(*out)
	}
	if in.Processes != nil {
		in, out := &in.Processes, &out.Processes
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	in.IoRchar.DeepCopyInto(&out.IoRchar)
	in.IoWchar.DeepCopyInto(&out.IoWchar)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ProcessProfile.
func (in *ProcessProfile) DeepCopy() *ProcessProfile {
	if in == nil {
		return nil
	}
	out := new(ProcessProfile)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *QueryConfig) DeepCopyInto(out *QueryConfig) {
	*out = *in
	in.Kv.DeepCopyInto(&out.Kv)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new QueryConfig.
func (in *QueryConfig) DeepCopy() *QueryConfig {
	if in == nil {
		return nil
	}
	out := new(QueryConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *QueryPile) DeepCopyInto(out *QueryPile) {
	*out = *in
	if in.Kv != nil {
		in, out := &in.Kv, &out.Kv
		*out = new(KeyValPile)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new QueryPile.
func (in *QueryPile) DeepCopy() *QueryPile {
	if in == nil {
		return nil
	}
	out := new(QueryPile)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *QueryProfile) DeepCopyInto(out *QueryProfile) {
	*out = *in
	if in.Kv != nil {
		in, out := &in.Kv, &out.Kv
		*out = new(KeyValProfile)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new QueryProfile.
func (in *QueryProfile) DeepCopy() *QueryProfile {
	if in == nil {
		return nil
	}
	out := new(QueryProfile)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ReqConfig) DeepCopyInto(out *ReqConfig) {
	*out = *in
	out.clientIp = in.clientIp.DeepCopy()
	out.hopIp = in.hopIp.DeepCopy()
	if in.ClientIp != nil {
		in, out := &in.ClientIp, &out.ClientIp
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.HopIp != nil {
		in, out := &in.HopIp, &out.HopIp
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Method != nil {
		in, out := &in.Method, &out.Method
		*out = make(Set, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Proto != nil {
		in, out := &in.Proto, &out.Proto
		*out = make(Set, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.ContentLength != nil {
		in, out := &in.ContentLength, &out.ContentLength
		*out = make(U8MinmaxSlice, len(*in))
		copy(*out, *in)
	}
	in.Url.DeepCopyInto(&out.Url)
	in.Qs.DeepCopyInto(&out.Qs)
	in.Headers.DeepCopyInto(&out.Headers)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ReqConfig.
func (in *ReqConfig) DeepCopy() *ReqConfig {
	if in == nil {
		return nil
	}
	out := new(ReqConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ReqPile) DeepCopyInto(out *ReqPile) {
	*out = *in
	if in.ClientIp != nil {
		in, out := &in.ClientIp, &out.ClientIp
		*out = make([]net.IP, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = make(net.IP, len(*in))
				copy(*out, *in)
			}
		}
	}
	if in.HopIp != nil {
		in, out := &in.HopIp, &out.HopIp
		*out = make([]net.IP, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = make(net.IP, len(*in))
				copy(*out, *in)
			}
		}
	}
	if in.Method != nil {
		in, out := &in.Method, &out.Method
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Proto != nil {
		in, out := &in.Proto, &out.Proto
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.ContentLength != nil {
		in, out := &in.ContentLength, &out.ContentLength
		*out = make([]byte, len(*in))
		copy(*out, *in)
	}
	in.Url.DeepCopyInto(&out.Url)
	in.Qs.DeepCopyInto(&out.Qs)
	in.Headers.DeepCopyInto(&out.Headers)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ReqPile.
func (in *ReqPile) DeepCopy() *ReqPile {
	if in == nil {
		return nil
	}
	out := new(ReqPile)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ReqProfile) DeepCopyInto(out *ReqProfile) {
	*out = *in
	if in.ClientIp != nil {
		in, out := &in.ClientIp, &out.ClientIp
		*out = make(net.IP, len(*in))
		copy(*out, *in)
	}
	if in.HopIp != nil {
		in, out := &in.HopIp, &out.HopIp
		*out = make(net.IP, len(*in))
		copy(*out, *in)
	}
	if in.Url != nil {
		in, out := &in.Url, &out.Url
		*out = new(UrlProfile)
		(*in).DeepCopyInto(*out)
	}
	if in.Qs != nil {
		in, out := &in.Qs, &out.Qs
		*out = new(QueryProfile)
		(*in).DeepCopyInto(*out)
	}
	if in.Headers != nil {
		in, out := &in.Headers, &out.Headers
		*out = new(HeadersProfile)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ReqProfile.
func (in *ReqProfile) DeepCopy() *ReqProfile {
	if in == nil {
		return nil
	}
	out := new(ReqProfile)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RespConfig) DeepCopyInto(out *RespConfig) {
	*out = *in
	in.Headers.DeepCopyInto(&out.Headers)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RespConfig.
func (in *RespConfig) DeepCopy() *RespConfig {
	if in == nil {
		return nil
	}
	out := new(RespConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RespPile) DeepCopyInto(out *RespPile) {
	*out = *in
	in.Headers.DeepCopyInto(&out.Headers)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RespPile.
func (in *RespPile) DeepCopy() *RespPile {
	if in == nil {
		return nil
	}
	out := new(RespPile)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RespProfile) DeepCopyInto(out *RespProfile) {
	*out = *in
	if in.Headers != nil {
		in, out := &in.Headers, &out.Headers
		*out = new(HeadersProfile)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RespProfile.
func (in *RespProfile) DeepCopy() *RespProfile {
	if in == nil {
		return nil
	}
	out := new(RespProfile)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in Set) DeepCopyInto(out *Set) {
	{
		in := &in
		*out = make(Set, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
		return
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Set.
func (in Set) DeepCopy() Set {
	if in == nil {
		return nil
	}
	out := new(Set)
	in.DeepCopyInto(out)
	return *out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SimpleValConfig) DeepCopyInto(out *SimpleValConfig) {
	*out = *in
	if in.NonReadables != nil {
		in, out := &in.NonReadables, &out.NonReadables
		*out = make(U8MinmaxSlice, len(*in))
		copy(*out, *in)
	}
	if in.Spaces != nil {
		in, out := &in.Spaces, &out.Spaces
		*out = make(U8MinmaxSlice, len(*in))
		copy(*out, *in)
	}
	if in.Unicodes != nil {
		in, out := &in.Unicodes, &out.Unicodes
		*out = make(U8MinmaxSlice, len(*in))
		copy(*out, *in)
	}
	if in.Digits != nil {
		in, out := &in.Digits, &out.Digits
		*out = make(U8MinmaxSlice, len(*in))
		copy(*out, *in)
	}
	if in.Letters != nil {
		in, out := &in.Letters, &out.Letters
		*out = make(U8MinmaxSlice, len(*in))
		copy(*out, *in)
	}
	if in.SpecialChars != nil {
		in, out := &in.SpecialChars, &out.SpecialChars
		*out = make(U8MinmaxSlice, len(*in))
		copy(*out, *in)
	}
	if in.Sequences != nil {
		in, out := &in.Sequences, &out.Sequences
		*out = make(U8MinmaxSlice, len(*in))
		copy(*out, *in)
	}
	if in.UnicodeFlags != nil {
		in, out := &in.UnicodeFlags, &out.UnicodeFlags
		*out = make(Uint32Slice, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SimpleValConfig.
func (in *SimpleValConfig) DeepCopy() *SimpleValConfig {
	if in == nil {
		return nil
	}
	out := new(SimpleValConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SimpleValPile) DeepCopyInto(out *SimpleValPile) {
	*out = *in
	if in.NonReadables != nil {
		in, out := &in.NonReadables, &out.NonReadables
		*out = make([]byte, len(*in))
		copy(*out, *in)
	}
	if in.Spaces != nil {
		in, out := &in.Spaces, &out.Spaces
		*out = make([]byte, len(*in))
		copy(*out, *in)
	}
	if in.Unicodes != nil {
		in, out := &in.Unicodes, &out.Unicodes
		*out = make([]byte, len(*in))
		copy(*out, *in)
	}
	if in.Digits != nil {
		in, out := &in.Digits, &out.Digits
		*out = make([]byte, len(*in))
		copy(*out, *in)
	}
	if in.Letters != nil {
		in, out := &in.Letters, &out.Letters
		*out = make([]byte, len(*in))
		copy(*out, *in)
	}
	if in.SpecialChars != nil {
		in, out := &in.SpecialChars, &out.SpecialChars
		*out = make([]byte, len(*in))
		copy(*out, *in)
	}
	if in.Sequences != nil {
		in, out := &in.Sequences, &out.Sequences
		*out = make([]byte, len(*in))
		copy(*out, *in)
	}
	if in.UnicodeFlags != nil {
		in, out := &in.UnicodeFlags, &out.UnicodeFlags
		*out = make(Uint32Slice, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SimpleValPile.
func (in *SimpleValPile) DeepCopy() *SimpleValPile {
	if in == nil {
		return nil
	}
	out := new(SimpleValPile)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SimpleValProfile) DeepCopyInto(out *SimpleValProfile) {
	*out = *in
	if in.UnicodeFlags != nil {
		in, out := &in.UnicodeFlags, &out.UnicodeFlags
		*out = make(Uint32Slice, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SimpleValProfile.
func (in *SimpleValProfile) DeepCopy() *SimpleValProfile {
	if in == nil {
		return nil
	}
	out := new(SimpleValProfile)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *U8Minmax) DeepCopyInto(out *U8Minmax) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new U8Minmax.
func (in *U8Minmax) DeepCopy() *U8Minmax {
	if in == nil {
		return nil
	}
	out := new(U8Minmax)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in U8MinmaxSlice) DeepCopyInto(out *U8MinmaxSlice) {
	{
		in := &in
		*out = make(U8MinmaxSlice, len(*in))
		copy(*out, *in)
		return
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new U8MinmaxSlice.
func (in U8MinmaxSlice) DeepCopy() U8MinmaxSlice {
	if in == nil {
		return nil
	}
	out := new(U8MinmaxSlice)
	in.DeepCopyInto(out)
	return *out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in Uint32Slice) DeepCopyInto(out *Uint32Slice) {
	{
		in := &in
		*out = make(Uint32Slice, len(*in))
		copy(*out, *in)
		return
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Uint32Slice.
func (in Uint32Slice) DeepCopy() Uint32Slice {
	if in == nil {
		return nil
	}
	out := new(Uint32Slice)
	in.DeepCopyInto(out)
	return *out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *UrlConfig) DeepCopyInto(out *UrlConfig) {
	*out = *in
	in.Val.DeepCopyInto(&out.Val)
	if in.Segments != nil {
		in, out := &in.Segments, &out.Segments
		*out = make(U8MinmaxSlice, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new UrlConfig.
func (in *UrlConfig) DeepCopy() *UrlConfig {
	if in == nil {
		return nil
	}
	out := new(UrlConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *UrlPile) DeepCopyInto(out *UrlPile) {
	*out = *in
	if in.Val != nil {
		in, out := &in.Val, &out.Val
		*out = new(SimpleValPile)
		(*in).DeepCopyInto(*out)
	}
	if in.Segments != nil {
		in, out := &in.Segments, &out.Segments
		*out = make([]byte, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new UrlPile.
func (in *UrlPile) DeepCopy() *UrlPile {
	if in == nil {
		return nil
	}
	out := new(UrlPile)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *UrlProfile) DeepCopyInto(out *UrlProfile) {
	*out = *in
	if in.Val != nil {
		in, out := &in.Val, &out.Val
		*out = new(SimpleValProfile)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new UrlProfile.
func (in *UrlProfile) DeepCopy() *UrlProfile {
	if in == nil {
		return nil
	}
	out := new(UrlProfile)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WsGate) DeepCopyInto(out *WsGate) {
	*out = *in
	if in.Configured != nil {
		in, out := &in.Configured, &out.Configured
		*out = new(Critiria)
		(*in).DeepCopyInto(*out)
	}
	if in.Learned != nil {
		in, out := &in.Learned, &out.Learned
		*out = new(Critiria)
		(*in).DeepCopyInto(*out)
	}
	out.Control = in.Control
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WsGate.
func (in *WsGate) DeepCopy() *WsGate {
	if in == nil {
		return nil
	}
	out := new(WsGate)
	in.DeepCopyInto(out)
	return out
}
