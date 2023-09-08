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

	saipb "github.com/openconfig/lemming/dataplane/standalone/proto"
)

func TestInterceptor(t *testing.T) {
	tests := []struct {
		desc        string
		req         any
		attrs       map[string]map[int32]protoreflect.Value
		handlerResp any
		handlerErr  error
		info        *grpc.UnaryServerInfo
		want        any
		wantErr     string
	}{{
		desc:  "not sai",
		info:  &grpc.UnaryServerInfo{FullMethod: "foo"},
		attrs: map[string]map[int32]protoreflect.Value{},
	}, {
		desc:       "handler error",
		info:       &grpc.UnaryServerInfo{FullMethod: "/lemming.dataplane.sai.Switch/CreateSwitch"},
		attrs:      map[string]map[int32]protoreflect.Value{},
		handlerErr: fmt.Errorf("foo"),
		wantErr:    "foo",
	}, {
		desc:  "create request",
		info:  &grpc.UnaryServerInfo{FullMethod: "/lemming.dataplane.sai.Switch/CreateSwitch"},
		attrs: map[string]map[int32]protoreflect.Value{},
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
		attrs: map[string]map[int32]protoreflect.Value{},
		req: &saipb.CreateSwitchRequest{
			RestartWarm: proto.Bool(true),
		},
		handlerErr: status.Error(codes.Unimplemented, "foo"),
		want: &saipb.CreateSwitchResponse{
			Oid: 1,
		},
	}, {
		desc:  "create request entry",
		info:  &grpc.UnaryServerInfo{FullMethod: "/lemming.dataplane.sai.Route/CreateRoute"},
		attrs: map[string]map[int32]protoreflect.Value{},
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
		attrs: map[string]map[int32]protoreflect.Value{
			"10": {int32(saipb.SwitchAttr_SWITCH_ATTR_CPU_PORT): protoreflect.ValueOfUint64(100)},
		},
		req: &saipb.GetSwitchAttributeRequest{
			Oid: 10,
			AttrType: []saipb.SwitchAttr{
				saipb.SwitchAttr_SWITCH_ATTR_CPU_PORT,
			},
		},
		want: &saipb.GetSwitchAttributeResponse{
			Attr: &saipb.SwitchAttribute{
				CpuPort: proto.Uint64(100),
			},
		},
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			mgr := AttrMgr{
				attrs: tt.attrs,
			}
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
