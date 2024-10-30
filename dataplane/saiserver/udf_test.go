// Copyright 2024 Google LLC
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

package saiserver

import (
	"context"
	"fmt"
	"slices"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"

	"github.com/openconfig/gnmi/errdiff"

	saipb "github.com/openconfig/lemming/dataplane/proto/sai"
	"github.com/openconfig/lemming/dataplane/saiserver/attrmgr"
)

func TestCreateUdf(t *testing.T) {
	tests := []struct {
		desc      string
		wantErr   string
		req       *saipb.CreateUdfRequest
		wantGroup []uint64
	}{{
		desc: "add to empty",
		req: &saipb.CreateUdfRequest{
			GroupId: proto.Uint64(10),
		},
		wantGroup: []uint64{1},
	}, {
		desc: "add to existing",
		req: &saipb.CreateUdfRequest{
			GroupId: proto.Uint64(11),
		},
		wantGroup: []uint64{20, 1},
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			dplane := &fakeSwitchDataplane{}
			c, mgr, stopFn := newTestUdf(t, dplane)
			defer stopFn()
			mgr.StoreAttributes(10, &saipb.UdfGroupAttribute{})
			mgr.StoreAttributes(11, &saipb.UdfGroupAttribute{
				UdfList: []uint64{20},
			})
			_, gotErr := c.CreateUdf(context.Background(), tt.req)
			if diff := errdiff.Check(gotErr, tt.wantErr); diff != "" {
				t.Fatalf("CreateNeighborEntry() unexpected err: %s", diff)
			}
			if gotErr != nil {
				return
			}
			gr := &saipb.UdfGroupAttribute{}
			mgr.PopulateAllAttributes(fmt.Sprint(tt.req.GetGroupId()), gr)
			if slices.Compare(gr.UdfList, tt.wantGroup) != 0 {
				t.Errorf("invalid group got %v, want: %v", gr.UdfList, tt.wantGroup)
			}
		})
	}
}

func newTestUdf(t testing.TB, api switchDataplaneAPI) (saipb.UdfClient, *attrmgr.AttrMgr, func()) {
	conn, mgr, stopFn := newTestServer(t, func(mgr *attrmgr.AttrMgr, srv *grpc.Server) {
		newUdf(mgr, api, srv)
	})
	return saipb.NewUdfClient(conn), mgr, stopFn
}
