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

//go:build linux

package ports

import (
	"fmt"
	"net"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/openconfig/gnmi/errdiff"
	"github.com/vishvananda/netlink"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"

	"github.com/openconfig/lemming/dataplane/forwarding/fwdport/porttestutil"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdcontext"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

type fakeStateIfaceMgr struct {
	link        netlink.Link
	linkErr     error
	setStateErr error
	gotName     string
	gotUp       *bool
}

func (f *fakeStateIfaceMgr) LinkByName(string) (netlink.Link, error) {
	return f.link, f.linkErr
}

func (f *fakeStateIfaceMgr) SetState(name string, up bool) error {
	f.gotName = name
	f.gotUp = &up
	return f.setStateErr
}

func TestGetAndSetState(t *testing.T) {
	tests := []struct {
		desc      string
		ifMgr     *fakeStateIfaceMgr
		pi        *fwdpb.PortInfo
		want      *fwdpb.PortStateReply
		wantSetUp *bool
		wantErr   string
	}{{
		desc: "link by name error",
		ifMgr: &fakeStateIfaceMgr{
			linkErr: fmt.Errorf("no such link"),
		},
		wantErr: "no such link",
	}, {
		desc: "read state admin down oper down",
		ifMgr: &fakeStateIfaceMgr{
			link: &netlink.Dummy{
				LinkAttrs: netlink.LinkAttrs{
					Flags:     0,
					OperState: netlink.OperDown,
				},
			},
		},
		want: &fwdpb.PortStateReply{
			Status: &fwdpb.PortInfo{
				AdminStatus: fwdpb.PortState_PORT_STATE_DISABLED_DOWN,
				OperStatus:  fwdpb.PortState_PORT_STATE_DISABLED_DOWN,
			},
		},
	}, {
		desc: "read state admin up oper lower layer down",
		ifMgr: &fakeStateIfaceMgr{
			link: &netlink.Dummy{
				LinkAttrs: netlink.LinkAttrs{
					Flags:     net.FlagUp,
					OperState: netlink.OperLowerLayerDown,
				},
			},
		},
		want: &fwdpb.PortStateReply{
			Status: &fwdpb.PortInfo{
				AdminStatus: fwdpb.PortState_PORT_STATE_ENABLED_UP,
				OperStatus:  fwdpb.PortState_PORT_STATE_DISABLED_DOWN,
			},
		},
	}, {
		desc: "read state admin up oper unknown",
		ifMgr: &fakeStateIfaceMgr{
			link: &netlink.Dummy{
				LinkAttrs: netlink.LinkAttrs{
					Flags:     net.FlagUp,
					OperState: netlink.OperUnknown,
				},
			},
		},
		want: &fwdpb.PortStateReply{
			Status: &fwdpb.PortInfo{
				AdminStatus: fwdpb.PortState_PORT_STATE_ENABLED_UP,
				OperStatus:  fwdpb.PortState_PORT_STATE_ENABLED_UP,
			},
		},
	}, {
		desc: "read state admin up oper up",
		ifMgr: &fakeStateIfaceMgr{
			link: &netlink.Dummy{
				LinkAttrs: netlink.LinkAttrs{
					Flags:     net.FlagUp,
					OperState: netlink.OperUp,
				},
			},
		},
		want: &fwdpb.PortStateReply{
			Status: &fwdpb.PortInfo{
				AdminStatus: fwdpb.PortState_PORT_STATE_ENABLED_UP,
				OperStatus:  fwdpb.PortState_PORT_STATE_ENABLED_UP,
			},
		},
	}, {
		desc: "set state unchanged admin state",
		ifMgr: &fakeStateIfaceMgr{
			link: &netlink.Dummy{
				LinkAttrs: netlink.LinkAttrs{
					Flags:     net.FlagUp,
					OperState: netlink.OperUp,
				},
			},
		},
		pi: &fwdpb.PortInfo{
			AdminStatus: fwdpb.PortState_PORT_STATE_ENABLED_UP,
		},
		want: &fwdpb.PortStateReply{
			Status: &fwdpb.PortInfo{
				AdminStatus: fwdpb.PortState_PORT_STATE_ENABLED_UP,
				OperStatus:  fwdpb.PortState_PORT_STATE_ENABLED_UP,
			},
		},
	}, {
		desc: "set state enable up",
		ifMgr: &fakeStateIfaceMgr{
			link: &netlink.Dummy{
				LinkAttrs: netlink.LinkAttrs{
					Flags:     0,
					OperState: netlink.OperDown,
				},
			},
		},
		pi: &fwdpb.PortInfo{
			AdminStatus: fwdpb.PortState_PORT_STATE_ENABLED_UP,
		},
		want: &fwdpb.PortStateReply{
			Status: &fwdpb.PortInfo{
				AdminStatus: fwdpb.PortState_PORT_STATE_ENABLED_UP,
				OperStatus:  fwdpb.PortState_PORT_STATE_DISABLED_DOWN,
			},
		},
		wantSetUp: proto.Bool(true),
	}, {
		desc: "set state disable down",
		ifMgr: &fakeStateIfaceMgr{
			link: &netlink.Dummy{
				LinkAttrs: netlink.LinkAttrs{
					Flags:     net.FlagUp,
					OperState: netlink.OperUp,
				},
			},
		},
		pi: &fwdpb.PortInfo{
			AdminStatus: fwdpb.PortState_PORT_STATE_DISABLED_DOWN,
		},
		want: &fwdpb.PortStateReply{
			Status: &fwdpb.PortInfo{
				AdminStatus: fwdpb.PortState_PORT_STATE_DISABLED_DOWN,
				OperStatus:  fwdpb.PortState_PORT_STATE_ENABLED_UP,
			},
		},
		wantSetUp: proto.Bool(false),
	}, {
		desc: "set state unspecified admin status",
		ifMgr: &fakeStateIfaceMgr{
			link: &netlink.Dummy{
				LinkAttrs: netlink.LinkAttrs{
					Flags:     net.FlagUp,
					OperState: netlink.OperUp,
				},
			},
		},
		pi: &fwdpb.PortInfo{
			AdminStatus: fwdpb.PortState_PORT_STATE_UNSPECIFIED,
		},
		want: &fwdpb.PortStateReply{
			Status: &fwdpb.PortInfo{
				AdminStatus: fwdpb.PortState_PORT_STATE_ENABLED_UP,
				OperStatus:  fwdpb.PortState_PORT_STATE_ENABLED_UP,
			},
		},
	}, {
		desc: "set state error",
		ifMgr: &fakeStateIfaceMgr{
			link: &netlink.Dummy{
				LinkAttrs: netlink.LinkAttrs{
					Flags:     0,
					OperState: netlink.OperDown,
				},
			},
			setStateErr: fmt.Errorf("set state failed"),
		},
		pi: &fwdpb.PortInfo{
			AdminStatus: fwdpb.PortState_PORT_STATE_ENABLED_UP,
		},
		wantSetUp: proto.Bool(true),
		wantErr:   "set state failed",
	}}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			got, err := getAndSetState("eth0", tt.ifMgr, tt.pi)
			if d := errdiff.Check(err, tt.wantErr); d != "" {
				t.Fatalf("getAndSetState() unexpected error: %s", d)
			}
			if d := cmp.Diff(got, tt.want, protocmp.Transform()); d != "" {
				t.Errorf("getAndSetState() unexpected diff (-got, +want):\n%s", d)
			}
			if d := cmp.Diff(tt.ifMgr.gotUp, tt.wantSetUp); d != "" {
				t.Errorf("SetState() unexpected up value (-got, +want):\n%s", d)
			}
		})
	}
}

type fakeLinkSubscriber struct {
	updCh    chan<- netlink.LinkUpdate
	doneCh   <-chan struct{}
	subErr   error
	subCalls int
}

func (f *fakeLinkSubscriber) LinkSubscribe(ch chan<- netlink.LinkUpdate, done <-chan struct{}) error {
	f.subCalls++
	if f.subErr != nil {
		return f.subErr
	}
	f.updCh = ch
	f.doneCh = done
	return nil
}

func TestLinkStateWatcher(t *testing.T) {
	t.Run("watch notify and unwatch", func(t *testing.T) {
		sub := &fakeLinkSubscriber{}
		w := &linkStateWatcher{
			ifaceMgr: sub,
			ports:    make(map[string]portTarget),
		}
		ctx := fwdcontext.New("test-ctx", "fwd")
		eventCh := make(chan *fwdpb.EventDesc, 10)
		if err := ctx.SetNotification(func(event *fwdpb.EventDesc) {
			eventCh <- event
		}); err != nil {
			t.Fatalf("SetNotification() failed: %v", err)
		}

		p1 := porttestutil.CreateTestPort(t, ctx, "port1")
		p2 := porttestutil.CreateTestPort(t, ctx, "port2")

		if err := w.watch("eth0", p1, ctx); err != nil {
			t.Fatalf("watch(eth0) unexpected error: %v", err)
		}
		if err := w.watch("eth1", p2, ctx); err != nil {
			t.Fatalf("watch(eth1) unexpected error: %v", err)
		}
		if sub.subCalls != 1 {
			t.Fatalf("LinkSubscribe() calls = %d, want 1", sub.subCalls)
		}

		// Update for unwatched interface should be ignored.
		sub.updCh <- netlink.LinkUpdate{
			Link: &netlink.Dummy{
				LinkAttrs: netlink.LinkAttrs{
					Name:      "unwatched0",
					Flags:     net.FlagUp,
					OperState: netlink.OperUp,
				},
			},
		}

		// Update for eth0 (admin up, oper up).
		sub.updCh <- netlink.LinkUpdate{
			Link: &netlink.Dummy{
				LinkAttrs: netlink.LinkAttrs{
					Name:      "eth0",
					Flags:     net.FlagUp,
					OperState: netlink.OperUp,
				},
			},
		}

		got := <-eventCh
		want := &fwdpb.EventDesc{
			Event: fwdpb.Event_EVENT_PORT,
			Desc: &fwdpb.EventDesc_Port{
				Port: &fwdpb.PortEventDesc{
					Context: &fwdpb.ContextId{Id: "test-ctx"},
					PortId:  &fwdpb.PortId{ObjectId: &fwdpb.ObjectId{Id: "port1"}},
					PortInfo: &fwdpb.PortInfo{
						AdminStatus: fwdpb.PortState_PORT_STATE_ENABLED_UP,
						OperStatus:  fwdpb.PortState_PORT_STATE_ENABLED_UP,
					},
				},
			},
		}
		if d := cmp.Diff(got, want, protocmp.Transform(), protocmp.IgnoreFields(&fwdpb.EventDesc{}, "sequence_number")); d != "" {
			t.Errorf("received event diff (-got, +want):\n%s", d)
		}

		// Unwatch eth0; watcher should remain active because eth1 is still watched.
		w.unwatch("eth0")
		select {
		case <-sub.doneCh:
			t.Fatal("doneCh closed before all ports were unwatched")
		default:
		}

		// Subsequent update for eth0 should be ignored, while eth1 update should notify.
		sub.updCh <- netlink.LinkUpdate{
			Link: &netlink.Dummy{
				LinkAttrs: netlink.LinkAttrs{
					Name:      "eth0",
					Flags:     0,
					OperState: netlink.OperDown,
				},
			},
		}
		sub.updCh <- netlink.LinkUpdate{
			Link: &netlink.Dummy{
				LinkAttrs: netlink.LinkAttrs{
					Name:      "eth1",
					Flags:     net.FlagUp,
					OperState: netlink.OperLowerLayerDown,
				},
			},
		}

		got = <-eventCh
		want = &fwdpb.EventDesc{
			Event: fwdpb.Event_EVENT_PORT,
			Desc: &fwdpb.EventDesc_Port{
				Port: &fwdpb.PortEventDesc{
					Context: &fwdpb.ContextId{Id: "test-ctx"},
					PortId:  &fwdpb.PortId{ObjectId: &fwdpb.ObjectId{Id: "port2"}},
					PortInfo: &fwdpb.PortInfo{
						AdminStatus: fwdpb.PortState_PORT_STATE_ENABLED_UP,
						OperStatus:  fwdpb.PortState_PORT_STATE_DISABLED_DOWN,
					},
				},
			},
		}
		if d := cmp.Diff(got, want, protocmp.Transform(), protocmp.IgnoreFields(&fwdpb.EventDesc{}, "sequence_number")); d != "" {
			t.Errorf("received event diff (-got, +want):\n%s", d)
		}

		// Unwatch last port; doneCh is closed synchronously in unwatch.
		w.unwatch("eth1")
		select {
		case <-sub.doneCh:
		default:
			t.Fatal("doneCh was not closed after unwatching all ports")
		}
	})
}
