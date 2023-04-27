// Copyright 2022 Google LLC
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

package kerneltest

import (
	"fmt"
	"net"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/openconfig/gnmi/errdiff"
	"github.com/vishvananda/netlink"
)

func TestSetHWAddr(t *testing.T) {
	tests := []struct {
		desc    string
		name    string
		addr    string
		want    net.HardwareAddr
		wantErr string
	}{{
		desc:    "non existing iface",
		name:    "eth0",
		wantErr: "link eth0 doesn't exist",
	}, {
		desc:    "invalid address",
		name:    "test",
		addr:    "hi",
		wantErr: "invalid MAC address",
	}, {
		desc: "success",
		name: "test",
		addr: "01:01:01:01:01:01",
		want: []byte{1, 1, 1, 1, 1, 1},
	}}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			fi := New(map[string]*Iface{
				"test": {},
			})
			gotErr := fi.SetHWAddr(tt.name, tt.addr)
			if diff := errdiff.Check(gotErr, tt.wantErr); diff != "" {
				t.Fatalf("SetHWAddr(%s, %s) unexpected err: %s", tt.name, tt.addr, diff)
			}
			if gotErr != nil {
				return
			}
			if d := cmp.Diff(fi.Links[tt.name].HWAddr, tt.want); d != "" {
				t.Fatalf("SetHWAddr(%s, %s) failed: diff (-got +want) %s", tt.name, tt.addr, d)
			}
		})
	}
}

func TestReplaceIP(t *testing.T) {
	tests := []struct {
		desc      string
		name      string
		addr      string
		prefixLen int
		wantErr   string
	}{{
		desc:    "non existing iface",
		name:    "eth0",
		wantErr: "link eth0 doesn't exist",
	}, {
		desc:    "invalid address",
		name:    "test",
		addr:    "hi",
		wantErr: "invalid CIDR address",
	}, {
		desc:      "success",
		name:      "test",
		addr:      "127.0.0.1",
		prefixLen: 32,
	}}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			fi := New(map[string]*Iface{
				"test": {},
			})
			gotErr := fi.ReplaceIP(tt.name, tt.addr, tt.prefixLen)
			if diff := errdiff.Check(gotErr, tt.wantErr); diff != "" {
				t.Fatalf("ReplaceIP(%s, %s, %d) unexpected err: %s", tt.name, tt.addr, tt.prefixLen, diff)
			}
			if gotErr != nil {
				return
			}
			if _, ok := fi.Links[tt.name].IPs[fmt.Sprintf("%s/%d", tt.addr, tt.prefixLen)]; !ok {
				t.Fatalf("ReplaceIP(%s, %s, %d) failed: values doesn't not exist", tt.name, tt.addr, tt.prefixLen)
			}
		})
	}
}

func TestDeleteIP(t *testing.T) {
	tests := []struct {
		desc      string
		name      string
		addr      string
		prefixLen int
		wantErr   string
	}{{
		desc:    "non existing iface",
		name:    "eth0",
		wantErr: "link eth0 doesn't exist",
	}, {
		desc:    "invalid address",
		name:    "test",
		addr:    "hi",
		wantErr: "invalid CIDR address",
	}, {
		desc:      "delete ip that doesn't exist",
		name:      "test",
		addr:      "127.0.0.2",
		prefixLen: 32,
		wantErr:   "not set in interface",
	}, {
		desc:      "success",
		name:      "test",
		addr:      "127.0.0.1",
		prefixLen: 32,
	}}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			fi := New(map[string]*Iface{
				"test": {},
			})
			fi.Links = map[string]*Iface{
				"test": {
					IPs: map[string]struct{}{"127.0.0.1/32": {}},
				},
			}
			gotErr := fi.DeleteIP(tt.name, tt.addr, tt.prefixLen)
			if diff := errdiff.Check(gotErr, tt.wantErr); diff != "" {
				t.Fatalf("DeleteIP(%s, %s, %d) unexpected err: %s", tt.name, tt.addr, tt.prefixLen, diff)
			}
			if gotErr != nil {
				return
			}
			if _, ok := fi.Links[tt.name].IPs[fmt.Sprintf("%s/%d", tt.addr, tt.prefixLen)]; ok {
				t.Fatalf("DeleteIP(%s, %s, %d) failed: value exists in map", tt.name, tt.addr, tt.prefixLen)
			}
		})
	}
}

func TestSetState(t *testing.T) {
	tests := []struct {
		desc    string
		name    string
		state   bool
		want    bool
		wantErr string
	}{{
		desc:    "non existing iface",
		name:    "eth0",
		wantErr: "link eth0 doesn't exist",
	}, {
		desc:  "success",
		name:  "test",
		state: true,
		want:  true,
	}}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			fi := New(map[string]*Iface{
				"test": {},
			})
			gotErr := fi.SetState(tt.name, tt.state)
			if diff := errdiff.Check(gotErr, tt.wantErr); diff != "" {
				t.Fatalf("SetState(%s, %v) unexpected err: %s", tt.name, tt.state, diff)
			}
			if gotErr != nil {
				return
			}
			if d := cmp.Diff(fi.Links[tt.name].up, tt.want); d != "" {
				t.Fatalf("SetState(%s, %v) failed: diff (-got +want) %s", tt.name, tt.state, d)
			}
		})
	}
}

func TestCreateTAP(t *testing.T) {
	tests := []struct {
		desc    string
		name    string
		want    *Iface
		wantFd  int
		wantErr string
	}{{
		desc:    "existing iface",
		name:    "test",
		wantErr: "link test already exist",
	}, {
		desc: "success",
		name: "tap0",
		want: &Iface{
			Idx:    2,
			HWAddr: []byte{1, 1, 0, 0, 1, 1},
		},
	}}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			fi := New(map[string]*Iface{
				"test": {},
			})
			gotFd, gotErr := fi.CreateTAP(tt.name)
			if diff := errdiff.Check(gotErr, tt.wantErr); diff != "" {
				t.Fatalf("CreateTAP(%s) unexpected err: %s", tt.name, diff)
			}
			if gotErr != nil {
				return
			}
			if d := cmp.Diff(fi.Links[tt.name], tt.want, cmp.AllowUnexported(Iface{})); d != "" {
				t.Errorf("CreateTAP(%s) failed: diff (-got +want) %s", tt.name, d)
			}
			if d := cmp.Diff(gotFd, tt.wantFd); d != "" {
				t.Errorf("CreateTAP(%s) failed: diff (-got +want) %s", tt.name, d)
			}
		})
	}
}

func TestLinkSubscribe(t *testing.T) {
	fi := New(map[string]*Iface{
		"test": {
			up: true,
		},
	})
	ch := make(chan netlink.LinkUpdate, 1)
	done := make(chan struct{})
	defer close(done)

	fi.LinkSubscribe(ch, done)
	fi.SetHWAddr("test", "01:01:01:01:01:01")
	got := <-ch
	want := netlink.LinkUpdate{
		Link: &netlink.Dummy{
			LinkAttrs: netlink.LinkAttrs{
				Name:         "test",
				HardwareAddr: []byte{01, 01, 01, 01, 01, 01},
				Flags:        net.FlagUp,
				OperState:    netlink.OperUp,
			},
		},
	}

	if d := cmp.Diff(got, want); d != "" {
		t.Errorf("LinkSubscribe() failed: diff (-got +want) %s", d)
	}
}

func TestAddrSubscribe(t *testing.T) {
	fi := New(map[string]*Iface{
		"test": {
			Idx: 1,
		},
	})
	ch := make(chan netlink.AddrUpdate, 1)
	done := make(chan struct{})

	fi.AddrSubscribe(ch, done)
	fi.ReplaceIP("test", "127.0.0.1", 32)
	got := <-ch
	want := netlink.AddrUpdate{
		LinkAddress: net.IPNet{
			IP:   []byte{127, 0, 0, 1},
			Mask: []byte{255, 255, 255, 255},
		},
		NewAddr:   true,
		LinkIndex: 1,
	}

	if d := cmp.Diff(got, want); d != "" {
		t.Errorf("AddrSubscribe() failed: diff (-got +want) %s", d)
	}
	close(done)
}
