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

package dataplane

import (
	"context"
	"errors"
	"net"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	dpb "github.com/openconfig/lemming/proto/dataplane"
)

func TestUpdatePort(t *testing.T) {
	tests := []struct {
		desc          string
		inReq         *dpb.UpdatePortRequest
		inMACErr      error
		inSetIPErr    error
		inSetStateErr error
		wantErrCode   codes.Code
	}{{
		desc:        "no name",
		inReq:       &dpb.UpdatePortRequest{},
		wantErrCode: codes.InvalidArgument,
	}, {
		desc: "invalid MAC",
		inReq: &dpb.UpdatePortRequest{
			Name:   "eth0",
			Hwaddr: "hi",
		},
		wantErrCode: codes.InvalidArgument,
	}, {
		desc: "error setting MAC",
		inReq: &dpb.UpdatePortRequest{
			Name:   "eth0",
			Hwaddr: "11:11:11:11:11:11",
		},
		inMACErr:    errors.New("fake"),
		wantErrCode: codes.Internal,
	}, {
		desc: "invalid IPv4",
		inReq: &dpb.UpdatePortRequest{
			Name: "eth0",
			Ipv4S: []*dpb.IPNetwork{{
				Ip: "hi",
			}},
		},
		wantErrCode: codes.InvalidArgument,
	}, {
		desc: "invalid IPv6",
		inReq: &dpb.UpdatePortRequest{
			Name: "eth0",
			Ipv6S: []*dpb.IPNetwork{{
				Ip: "hi",
			}},
		},
		wantErrCode: codes.InvalidArgument,
	}, {
		desc: "error setting IPs",
		inReq: &dpb.UpdatePortRequest{
			Name: "eth0",
			Ipv4S: []*dpb.IPNetwork{{
				Ip:        "127.0.0.1",
				PrefixLen: 24,
			}},
			Ipv6S: []*dpb.IPNetwork{{
				Ip:        "::1",
				PrefixLen: 24,
			}},
		},
		inSetIPErr:  errors.New("fake"),
		wantErrCode: codes.Internal,
	}, {
		desc: "error setting state",
		inReq: &dpb.UpdatePortRequest{
			Name:       "eth0",
			AdminState: dpb.PortState_PORT_STATE_UP,
		},
		inSetStateErr: errors.New("fake"),
		wantErrCode:   codes.Internal,
	}, {
		desc: "success with all fields set",
		inReq: &dpb.UpdatePortRequest{
			Name:   "eth0",
			Hwaddr: "11:11:11:11:11:11",
			Ipv4S: []*dpb.IPNetwork{{
				Ip:        "127.0.0.1",
				PrefixLen: 24,
			}},
			Ipv6S: []*dpb.IPNetwork{{
				Ip:        "::1",
				PrefixLen: 24,
			}},
			AdminState: dpb.PortState_PORT_STATE_UP,
		},
		wantErrCode: codes.OK,
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			srv := &Server{}
			setInterfaceHWAddr = func(string, net.HardwareAddr) error {
				return tt.inMACErr
			}
			setInterfaceIPs = func(string, []*net.IPNet) error {
				return tt.inSetIPErr
			}
			setInterfaceState = func(string, bool) error {
				return tt.inSetStateErr
			}
			_, err := srv.UpdatePort(context.Background(), tt.inReq)
			if gotCode := status.Code(err); gotCode != tt.wantErrCode {
				t.Fatalf("UpdatePort(%v) returned unexpected code: got %v, want %v", tt.inReq, gotCode, tt.wantErrCode)
			}
		})
	}
}
