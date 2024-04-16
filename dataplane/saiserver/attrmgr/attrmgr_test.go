// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package attrmgr

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/openconfig/gnmi/errdiff"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/testing/protocmp"

	saipb "github.com/openconfig/lemming/dataplane/proto/sai"
)

func ptrToValue(v protoreflect.Value) *protoreflect.Value {
	return &v
}

func TestInterceptor(t *testing.T) {
	tests := []struct {
		desc        string
		req         any
		attrs       map[string]map[int32]*protoreflect.Value
		handlerResp any
		handlerErr  error
		info        *grpc.UnaryServerInfo
		want        any
		wantErr     string
	}{{
		desc:  "not sai",
		info:  &grpc.UnaryServerInfo{FullMethod: "foo"},
		attrs: map[string]map[int32]*protoreflect.Value{},
	}, {
		desc:       "handler error",
		info:       &grpc.UnaryServerInfo{FullMethod: "/lemming.dataplane.sai.Switch/CreateSwitch"},
		attrs:      map[string]map[int32]*protoreflect.Value{},
		handlerErr: fmt.Errorf("foo"),
		wantErr:    "foo",
	}, {
		desc:  "create request",
		info:  &grpc.UnaryServerInfo{FullMethod: "/lemming.dataplane.sai.Switch/CreateSwitch"},
		attrs: map[string]map[int32]*protoreflect.Value{},
		req: &saipb.CreateSwitchRequest{
			RestartWarm: proto.Bool(true),
		},
		handlerResp: &saipb.CreateSwitchResponse{},
		want: &saipb.CreateSwitchResponse{
			Oid: 1,
		},
	}, {
		desc:  "create request unimplemented",
		info:  &grpc.UnaryServerInfo{FullMethod: "/lemming.dataplane.sai.Switch/CreateSwitch"},
		attrs: map[string]map[int32]*protoreflect.Value{},
		req: &saipb.CreateSwitchRequest{
			RestartWarm: proto.Bool(true),
		},
		handlerErr: status.Error(codes.Unimplemented, "foo"),
		wantErr:    "Unimplemented",
	}, {
		desc:  "create request unimplemented typed nil",
		info:  &grpc.UnaryServerInfo{FullMethod: "/lemming.dataplane.sai.Switch/CreateSwitch"},
		attrs: map[string]map[int32]*protoreflect.Value{},
		req: &saipb.CreateSwitchRequest{
			RestartWarm: proto.Bool(true),
		},
		handlerResp: (*saipb.CreateSwitchResponse)(nil),
		handlerErr:  status.Error(codes.Unimplemented, "foo"),
		wantErr:     "Unimplemented",
	}, {
		desc:  "create request entry",
		info:  &grpc.UnaryServerInfo{FullMethod: "/lemming.dataplane.sai.Route/CreateRoute"},
		attrs: map[string]map[int32]*protoreflect.Value{},
		req: &saipb.CreateRouteEntryRequest{
			Entry: &saipb.RouteEntry{
				SwitchId: 12,
			},
		},
		handlerResp: &saipb.CreateRouteEntryResponse{},
		want:        &saipb.CreateRouteEntryResponse{},
	}, {
		desc: "get request",
		info: &grpc.UnaryServerInfo{FullMethod: "/lemming.dataplane.sai.Switch/GetSwitchAttribute"},
		attrs: map[string]map[int32]*protoreflect.Value{
			"10": {
				int32(saipb.SwitchAttr_SWITCH_ATTR_CPU_PORT):        ptrToValue(protoreflect.ValueOfUint64(100)),
				int32(saipb.SwitchAttr_SWITCH_ATTR_PRE_INGRESS_ACL): ptrToValue(protoreflect.ValueOfUint64(300)),
			},
		},
		req: &saipb.GetSwitchAttributeRequest{
			Oid: 10,
			AttrType: []saipb.SwitchAttr{
				saipb.SwitchAttr_SWITCH_ATTR_CPU_PORT,
				saipb.SwitchAttr_SWITCH_ATTR_PRE_INGRESS_ACL,
			},
		},
		want: &saipb.GetSwitchAttributeResponse{
			Attr: &saipb.SwitchAttribute{
				CpuPort:       proto.Uint64(100),
				PreIngressAcl: proto.Uint64(300),
			},
		},
	}, {
		desc: "stats request",
		info: &grpc.UnaryServerInfo{FullMethod: "/lemming.dataplane.sai.Port/GetPortStats"},
		attrs: map[string]map[int32]*protoreflect.Value{
			"10": {
				int32(saipb.SwitchAttr_SWITCH_ATTR_CPU_PORT): ptrToValue(protoreflect.ValueOfUint64(100)),
			},
		},
		req: &saipb.GetPortStatsRequest{
			Oid: 10,
		},
		want: &saipb.GetPortStatsResponse{},
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			mgr := New()
			mgr.attrs = tt.attrs
			got, gotErr := mgr.Interceptor(context.TODO(), tt.req, tt.info, func(context.Context, any) (any, error) {
				return tt.handlerResp, tt.handlerErr
			})
			if diff := errdiff.Check(gotErr, tt.wantErr); diff != "" {
				t.Fatalf("Interceptor() unexpected err: %s", diff)
			}
			if gotErr != nil {
				return
			}
			if d := cmp.Diff(got, tt.want, protocmp.Transform()); d != "" {
				t.Fatalf("Interceptor() failed: diff(-got,+want)\n:%s", d)
			}
		})
	}
}

func TestInvokeAndSave(t *testing.T) {
	tests := []struct {
		desc      string
		req       proto.Message
		rpc       func(context.Context, proto.Message) (proto.Message, error)
		wantAttrs map[string]map[int32]*protoreflect.Value
		wantErr   string
	}{{
		desc: "rpc error",
		rpc: func(context.Context, proto.Message) (proto.Message, error) {
			return nil, fmt.Errorf("foo")
		},
		wantErr: "foo",
	}, {
		desc: "success",
		rpc: func(context.Context, proto.Message) (proto.Message, error) {
			return nil, nil
		},
		req: &saipb.SetPortAttributeRequest{
			Oid:             1,
			AdminState:      proto.Bool(true),
			AdvertisedSpeed: []uint32{},
		},
		wantAttrs: map[string]map[int32]*protoreflect.Value{
			"1": {
				int32(saipb.PortAttr_PORT_ATTR_ADMIN_STATE):      ptrToValue(protoreflect.ValueOfBool(true)),
				int32(saipb.PortAttr_PORT_ATTR_ADVERTISED_SPEED): nil,
			},
		},
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			mgr := New()
			_, gotErr := InvokeAndSave(context.Background(), mgr, tt.rpc, tt.req)
			if diff := errdiff.Check(gotErr, tt.wantErr); diff != "" {
				t.Fatalf("InvokeAndSave() unexpected err: %s", diff)
			}
			if gotErr != nil {
				return
			}
			if d := cmp.Diff(mgr.attrs, tt.wantAttrs, protocmp.Transform()); d != "" {
				t.Fatalf("InvokeAndSave() failed: diff(-got,+want)\n:%s", d)
			}
		})
	}
}

func TestGetType(t *testing.T) {
	tests := []struct {
		desc  string
		types map[string]saipb.ObjectType
		id    string
		want  saipb.ObjectType
	}{{
		desc:  "unknown val",
		types: make(map[string]saipb.ObjectType),
		want:  saipb.ObjectType_OBJECT_TYPE_NULL,
	}, {
		desc:  "existing val",
		id:    "1",
		types: map[string]saipb.ObjectType{"1": saipb.ObjectType_OBJECT_TYPE_SWITCH},
		want:  saipb.ObjectType_OBJECT_TYPE_SWITCH,
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			mgr := New()
			mgr.idToType = tt.types
			got := mgr.GetType(tt.id)
			if d := cmp.Diff(got, tt.want, protocmp.Transform()); d != "" {
				t.Fatalf("Interceptor() failed: diff(-got,+want)\n:%s", d)
			}
		})
	}
}
