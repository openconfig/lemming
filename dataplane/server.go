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

// Package dataplane is an implementation of the dataplane HAL API.
package dataplane

import (
	"context"
	"fmt"
	"net"

	"github.com/openconfig/lemming/dataplane/internal/kernel"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	dpb "github.com/openconfig/lemming/proto/dataplane"
)

const (
	Port = 6443
)

// Server is an implementation of Dataplane HAL API.
type Server struct {
	dpb.UnimplementedHALServer
}

var (
	// Stubs for testing
	setInterfaceHWAddr = kernel.SetInterfaceHWAddr
	setInterfaceIPs    = kernel.SetInterfaceIPs
	setInterfaceState  = kernel.SetInterfaceState
)

// UpdatePort updates an interface properties from the input request.
func (s *Server) UpdatePort(_ context.Context, req *dpb.UpdatePortRequest) (*dpb.UpdatePortResponse, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name not set")
	}
	tapName := fmt.Sprintf("%s-tap", req.Name)
	if req.Hwaddr != "" {
		addr, err := net.ParseMAC(req.Hwaddr)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid MAC address %q: parse err %v", req.Hwaddr, err)
		}
		if err := setInterfaceHWAddr(tapName, addr); err != nil {
			return nil, status.Errorf(codes.Internal, "failed to set HW addr: %v", err)
		}
	}
	var ips []*net.IPNet
	for _, ip := range req.Ipv4S {
		parsedIP := net.ParseIP(ip.GetIp())
		if parsedIP == nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid IP address: %q", ip.GetIp())
		}
		ips = append(ips, &net.IPNet{
			IP:   parsedIP,
			Mask: net.CIDRMask(int(ip.GetPrefixLen()), 32),
		})
	}
	for _, ip := range req.Ipv6S {
		parsedIP := net.ParseIP(ip.GetIp())
		if parsedIP == nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid IP address: %q", ip.GetIp())
		}
		ips = append(ips, &net.IPNet{
			IP:   parsedIP,
			Mask: net.CIDRMask(int(ip.GetPrefixLen()), 128),
		})
	}
	if err := setInterfaceIPs(tapName, ips); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to set IPs: %v", err)
	}

	var err error
	switch req.AdminState {
	case dpb.PortState_PORT_STATE_UP:
		err = setInterfaceState(tapName, req.AdminState == dpb.PortState_PORT_STATE_UP)
	case dpb.PortState_PORT_STATE_DOWN:
		err = setInterfaceState(tapName, req.AdminState == dpb.PortState_PORT_STATE_UP)
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to set interface set: %v", err)
	}

	return &dpb.UpdatePortResponse{}, nil
}
