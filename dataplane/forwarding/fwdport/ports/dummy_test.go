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

package ports

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/openconfig/gnmi/errdiff"
	"google.golang.org/protobuf/testing/protocmp"

	"github.com/openconfig/lemming/dataplane/forwarding/fwdaction"
	"github.com/openconfig/lemming/dataplane/forwarding/fwdport"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdcontext"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdobject"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

func TestDummyBuild(t *testing.T) {
	tests := []struct {
		desc      string
		inDesc    *fwdpb.PortDesc
		wantState *fwdpb.PortStateReply
		wantDev   string
		wantErr   string
	}{{
		desc: "success admin down",
		inDesc: &fwdpb.PortDesc{
			PortType: fwdpb.PortType_PORT_TYPE_DUMMY,
			PortId:   fwdport.MakeID(fwdobject.NewID("dummy1")),
			Port: &fwdpb.PortDesc_Dummy{
				Dummy: &fwdpb.DummyPortDesc{
					DeviceName: "eth0",
					AdminState: false,
				},
			},
		},
		wantDev: "eth0",
		wantState: &fwdpb.PortStateReply{
			Status: &fwdpb.PortInfo{
				OperStatus:  fwdpb.PortState_PORT_STATE_DISABLED_DOWN,
				AdminStatus: fwdpb.PortState_PORT_STATE_DISABLED_DOWN,
			},
		},
	}, {
		desc: "success admin up",
		inDesc: &fwdpb.PortDesc{
			PortType: fwdpb.PortType_PORT_TYPE_DUMMY,
			PortId:   fwdport.MakeID(fwdobject.NewID("dummy2")),
			Port: &fwdpb.PortDesc_Dummy{
				Dummy: &fwdpb.DummyPortDesc{
					DeviceName: "eth1",
					AdminState: true,
				},
			},
		},
		wantDev: "eth1",
		wantState: &fwdpb.PortStateReply{
			Status: &fwdpb.PortInfo{
				OperStatus:  fwdpb.PortState_PORT_STATE_DISABLED_DOWN,
				AdminStatus: fwdpb.PortState_PORT_STATE_ENABLED_UP,
			},
		},
	}}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			ctx := fwdcontext.New("test", "fwd")
			p, err := (dummyBuilder{}).Build(tt.inDesc, ctx)
			if d := errdiff.Check(err, tt.wantErr); d != "" {
				t.Fatalf("Build() unexpected error: %s", d)
			}
			if tt.wantErr != "" {
				return
			}

			if got := p.Type(); got != fwdpb.PortType_PORT_TYPE_DUMMY {
				t.Errorf("Type() = %v, want %v", got, fwdpb.PortType_PORT_TYPE_DUMMY)
			}
			if d := cmp.Diff(p.Desc(), tt.inDesc, protocmp.Transform()); d != "" {
				t.Errorf("Desc() unexpected diff (-got, +want):\n%s", d)
			}

			gotState, err := p.State(nil)
			if err != nil {
				t.Fatalf("State(nil) unexpected error: %v", err)
			}
			if d := cmp.Diff(gotState, tt.wantState, protocmp.Transform()); d != "" {
				t.Errorf("State(nil) unexpected diff (-got, +want):\n%s", d)
			}
		})
	}
}

func TestDummyPortOperations(t *testing.T) {
	ctx := fwdcontext.New("test", "fwd")
	desc := &fwdpb.PortDesc{
		PortType: fwdpb.PortType_PORT_TYPE_DUMMY,
		PortId:   fwdport.MakeID(fwdobject.NewID("dummy-port")),
		Port: &fwdpb.PortDesc_Dummy{
			Dummy: &fwdpb.DummyPortDesc{
				DeviceName: "eth0",
				AdminState: false,
			},
		},
	}

	p, err := fwdport.New(desc, ctx)
	if err != nil {
		t.Fatalf("fwdport.New() unexpected error: %v", err)
	}

	if err := p.Update(&fwdpb.PortUpdateDesc{}); err != nil {
		t.Errorf("Update() unexpected error: %v", err)
	}

	if got := p.Actions(fwdpb.PortAction_PORT_ACTION_INPUT); got != nil {
		t.Errorf("Actions(INPUT) = %v, want nil", got)
	}
	if got := p.Actions(fwdpb.PortAction_PORT_ACTION_OUTPUT); got != nil {
		t.Errorf("Actions(OUTPUT) = %v, want nil", got)
	}

	state, err := p.Write(nil)
	if err != nil {
		t.Fatalf("Write() unexpected error: %v", err)
	}
	if state != fwdaction.CONSUME {
		t.Errorf("Write() state = %v, want %v", state, fwdaction.CONSUME)
	}

	// Update admin state via State().
	gotState, err := p.State(&fwdpb.PortInfo{
		AdminStatus: fwdpb.PortState_PORT_STATE_ENABLED_UP,
	})
	if err != nil {
		t.Fatalf("State() unexpected error: %v", err)
	}
	wantState := &fwdpb.PortStateReply{
		Status: &fwdpb.PortInfo{
			OperStatus:  fwdpb.PortState_PORT_STATE_DISABLED_DOWN,
			AdminStatus: fwdpb.PortState_PORT_STATE_ENABLED_UP,
		},
	}
	if d := cmp.Diff(gotState, wantState, protocmp.Transform()); d != "" {
		t.Errorf("State(enable) unexpected diff (-got, +want):\n%s", d)
	}

	// Verify the updated admin state persists when querying with nil.
	gotState, err = p.State(nil)
	if err != nil {
		t.Fatalf("State(nil) unexpected error: %v", err)
	}
	if d := cmp.Diff(gotState, wantState, protocmp.Transform()); d != "" {
		t.Errorf("State(nil) after update unexpected diff (-got, +want):\n%s", d)
	}
}
