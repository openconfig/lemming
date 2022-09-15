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

package gribi

import (
	"context"
	"fmt"
	"net"

	log "github.com/golang/glog"
	"github.com/openconfig/gribigo/aft"
	"github.com/openconfig/gribigo/afthelper"
	"github.com/openconfig/gribigo/constants"
	"github.com/openconfig/gribigo/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	gribipb "github.com/openconfig/gribi/v1/proto/service"
	zpb "github.com/openconfig/lemming/proto/sysrib"
	"github.com/openconfig/lemming/sysrib"
)

// Server is a fake gRIBI implementation.
type Server struct {
	*server.Server
	s *grpc.Server
}

// New returns a new fake gRIBI server.
func New(s *grpc.Server) (*Server, error) {
	gs, err := createGRIBIServer()
	if err != nil {
		return nil, fmt.Errorf("cannot create gRIBI server, %v", err)
	}

	srv := &Server{
		Server: gs,
		s:      s,
	}
	gribipb.RegisterGRIBIServer(s, srv)

	return srv, nil
}

// createGRIBIServer creates and returns a gRIBI server that is ready be
// registered by a gRPC server.
func createGRIBIServer() (*server.Server, error) {
	gzebraConn, err := grpc.DialContext(context.Background(), fmt.Sprintf("unix:%s", sysrib.SockAddr), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("cannot dial to sysrib, %v", err)
	}
	gzebraClient := zpb.NewSysribClient(gzebraConn)

	ribAddfn := func(ribs map[string]*aft.RIB, optype constants.OpType, netinst, prefix string) {
		if optype != constants.Add {
			// TODO(wenbli): handle replace and delete :-)
			// For replace, just need to ensure Sysrib's gRPC supports it.
			return
		}
		nhs, err := afthelper.NextHopAddrsForPrefix(ribs, netinst, prefix)
		if err != nil {
			log.Errorf("cannot add netinst:prefix %s:%s to the RIB, %v", netinst, prefix, err)
			return
		}
		nhSum := []*afthelper.NextHopSummary{}
		for _, nh := range nhs {
			nhSum = append(nhSum, nh)
		}

		routeReq, err := createSetRouteRequest(prefix, nhSum)
		if err != nil {
			log.Errorf("Cannot create SetRouteRequest: %v", err)
		}

		resp, err := gzebraClient.SetRoute(context.Background(), routeReq)
		if err != nil {
			log.Errorf("Error sending route to sysrib: %v", err)
		}
		log.Infof("Sent route %v with response %v", routeReq, resp)
	}

	return server.New(
		server.WithRIBResolvedEntryHook(ribAddfn),
	)
}

// createSetRouteRequest converts a Route to a sysrib SetRouteRequest
func createSetRouteRequest(prefix string, nexthops []*afthelper.NextHopSummary) (*zpb.SetRouteRequest, error) {
	ip, ipnet, err := net.ParseCIDR(prefix)
	if err != nil {
		log.Errorf("Cannot parse prefix %q as CIDR for calling sysrib", prefix)
	}

	if err != nil {
		return nil, fmt.Errorf("gribigo/sysrib: %v", err)
	}
	maskLength, _ := ipnet.Mask.Size()

	var zNexthops []*zpb.Nexthop
	for _, nhs := range nexthops {
		zNexthops = append(zNexthops, &zpb.Nexthop{
			Type:    zpb.Nexthop_TYPE_IPV4,
			Address: nhs.Address,
			Weight:  nhs.Weight,
		})
	}

	return &zpb.SetRouteRequest{
		AdminDistance: 5,
		ProtocolName:  "gRIBI",
		Safi:          zpb.SetRouteRequest_SAFI_UNICAST,
		Prefix: &zpb.Prefix{
			Family:     zpb.Prefix_FAMILY_IPV4,
			Address:    ip.String(),
			MaskLength: uint32(maskLength),
		},
		Nexthops: zNexthops,
	}, nil
}
