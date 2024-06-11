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
	"net/netip"

	log "github.com/golang/glog"
	"github.com/openconfig/gribigo/aft"
	"github.com/openconfig/gribigo/afthelper"
	"github.com/openconfig/gribigo/constants"
	"github.com/openconfig/gribigo/rib"
	"github.com/openconfig/gribigo/server"
	"github.com/openconfig/ygnmi/ygnmi"
	"github.com/openconfig/ygot/ygot"
	"github.com/openconfig/ygot/ytypes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/openconfig/lemming/gnmi/gnmiclient"
	"github.com/openconfig/lemming/gnmi/oc"
	"github.com/openconfig/lemming/gnmi/oc/ocpath"

	gpb "github.com/openconfig/gnmi/proto/gnmi"
	gribipb "github.com/openconfig/gribi/v1/proto/service"

	sysribpb "github.com/openconfig/lemming/proto/sysrib"
)

// Server is a fake gRIBI implementation.
type Server struct {
	*server.Server
	s *grpc.Server
}

// New returns a new fake gRIBI server.
//
// - s is the gRPC server on which the reference gRIBI service will be
// installed.
// - root, if specified, will be used to populate connected routes into the RIB
// manager. Note this is intended to be used for unit/standalone device testing.
func New(s *grpc.Server, gClient gpb.GNMIClient, target string, root *oc.Root, sysribAddr string) (*Server, error) {
	gs, err := createGRIBIServer(gClient, target, root, sysribAddr)
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
//
// - root, if specified, will be used to populate connected routes into the RIB
// manager. Note this is intended to be used for unit/standalone device testing.
func createGRIBIServer(gClient gpb.GNMIClient, target string, root *oc.Root, sysribAddr string) (*server.Server, error) {
	gzebraConn, err := grpc.NewClient(sysribAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("cannot dial to sysrib, %v", err)
	}
	gzebraClient := sysribpb.NewSysribClient(gzebraConn)

	networkInstances := []string{}
	for name, ni := range root.NetworkInstance {
		if ni.Type == oc.NetworkInstanceTypes_NETWORK_INSTANCE_TYPE_L3VRF {
			networkInstances = append(networkInstances, name)
		}
	}

	yclient, err := ygnmi.NewClient(gClient, ygnmi.WithTarget(target), ygnmi.WithRequestLogLevel(2))
	if err != nil {
		return nil, err
	}

	ribHookfn := func(o constants.OpType, _ int64, ni string, data ygot.ValidatedGoStruct) {
		if o != constants.Add {
			// TODO(wenbli): handle replace and delete :-)
			return
		}
		// write gNMI notifications
		if err := updateAft(yclient, o, ni, data); err != nil {
			log.Errorf("invalid notifications, %v", err)
		}

		// TODO(wenbli): Check if this is needed with @robshakir.
		// server.WithFIBProgrammedCheck()
		//   -> gives us a function that checks whether an ID is a tristate (ok, failed, pending)
		//   -> plumb this through to the rib - and have a fib pending queue for responding.
		//
		// here we just write to something that the server has access to.
	}

	ribAddfn := func(ribs map[string]*aft.RIB, optype constants.OpType, netinst string, aft constants.AFT, key any, _ ...rib.ResolvedDetails) {
		prefix, ok := key.(string)
		if !ok {
			log.Errorf("Key is not a string type: (%T, %v)", key, key)
		}
		switch aft {
		case constants.IPv4, constants.IPv6:
		default:
			log.Errorf("Incompatible type of route receive, type: %s, key: %v", aft, key)
		}
		nhSum := []*afthelper.NextHopSummary{}
		switch optype {
		case constants.Add, constants.Replace:
			nhs, err := afthelper.NextHopAddrsForPrefix(ribs, netinst, prefix)
			if err != nil {
				log.Errorf("cannot add netinst:prefix %s:%s to the RIB, %v", netinst, prefix, err)
				return
			}
			for _, nh := range nhs {
				nhSum = append(nhSum, nh)
			}
		case constants.Delete:
		default:
			return
		}

		routeReq, err := createSetRouteRequest(prefix, nhSum)
		if err != nil {
			log.Errorf("Cannot create SetRouteRequest: %v", err)
			return
		}
		if optype == constants.Delete {
			routeReq.Delete = true
		}

		resp, err := gzebraClient.SetRoute(context.Background(), routeReq)
		if err != nil {
			log.Errorf("Error sending route to sysrib: %v", err)
			return
		}
		log.Infof("Sent route %v with response %v", routeReq, resp)
	}

	return server.New(
		server.WithPostChangeRIBHook(ribHookfn),
		server.WithRIBResolvedEntryHook(ribAddfn),
		server.WithVRFs(networkInstances),
	)
}

// createSetRouteRequest converts a Route to a sysrib SetRouteRequest
func createSetRouteRequest(prefix string, nexthops []*afthelper.NextHopSummary) (*sysribpb.SetRouteRequest, error) {
	pfx, err := netip.ParsePrefix(prefix)
	if err != nil {
		log.Errorf("Cannot parse prefix %q as CIDR for calling sysrib", prefix)
	}

	if err != nil {
		return nil, fmt.Errorf("gribigo/sysrib: %v", err)
	}

	var zNexthops []*sysribpb.Nexthop
	for _, nhs := range nexthops {
		zNexthops = append(zNexthops, &sysribpb.Nexthop{
			Type:    sysribpb.Nexthop_TYPE_IPV4,
			Address: nhs.Address,
			Weight:  nhs.Weight,
		})
	}

	family := sysribpb.Prefix_FAMILY_IPV4
	if pfx.Addr().Is6() {
		family = sysribpb.Prefix_FAMILY_IPV6
	}

	return &sysribpb.SetRouteRequest{
		AdminDistance: 5,
		ProtocolName:  "gRIBI",
		Safi:          sysribpb.SetRouteRequest_SAFI_UNICAST,
		Prefix: &sysribpb.Prefix{
			Family:     family,
			Address:    pfx.Addr().String(),
			MaskLength: uint32(pfx.Bits()),
		},
		Nexthops: zNexthops,
	}, nil
}

// convertGoStruct converts GoStruct a to GoStruct b.
//
// - unmarshal is the generated Unmarshal function of b's generated package.
func convertGoStruct(a, b ygot.GoStruct, unmarshal func(data []byte, destStruct ygot.GoStruct, opts ...ytypes.UnmarshalOpt) error) error {
	data, err := ygot.Marshal7951(a)
	if err != nil {
		return err
	}
	return unmarshal(data, b)
}

// updateAft creates the corresponding ygnmi PathStruct from a RIB operation.
func updateAft(yclient *ygnmi.Client, _ constants.OpType, ni string, e ygot.GoStruct) error {
	var err error
	switch t := e.(type) {
	case *aft.Afts_Ipv4Entry:
		dst := &oc.NetworkInstance_Afts_Ipv4Entry{}
		if err = convertGoStruct(t, dst, oc.Unmarshal); err != nil {
			break
		}
		path := ocpath.Root().NetworkInstance(ni).Afts().Ipv4Entry(t.GetPrefix()).State()
		if _, err := gnmiclient.Update(context.Background(), yclient, path, dst); err != nil {
			log.Warningf("unable to update gRIBI data: %v", err)
		}
	case *aft.Afts_Ipv6Entry:
		dst := &oc.NetworkInstance_Afts_Ipv6Entry{}
		if err = convertGoStruct(t, dst, oc.Unmarshal); err != nil {
			break
		}
		path := ocpath.Root().NetworkInstance(ni).Afts().Ipv6Entry(t.GetPrefix()).State()
		if _, err := gnmiclient.Update(context.Background(), yclient, path, dst); err != nil {
			log.Warningf("unable to update gRIBI data: %v", err)
		}
	case *aft.Afts_NextHopGroup:
		dst := &oc.NetworkInstance_Afts_NextHopGroup{}
		if err = convertGoStruct(t, dst, oc.Unmarshal); err != nil {
			break
		}
		path := ocpath.Root().NetworkInstance(ni).Afts().NextHopGroup(t.GetId()).State()
		if _, err := gnmiclient.Update(context.Background(), yclient, path, dst); err != nil {
			log.Warningf("unable to update gRIBI data: %v", err)
		}
	case *aft.Afts_NextHop:
		dst := &oc.NetworkInstance_Afts_NextHop{}
		if err = convertGoStruct(t, dst, oc.Unmarshal); err != nil {
			break
		}
		path := ocpath.Root().NetworkInstance(ni).Afts().NextHop(t.GetIndex()).State()
		if _, err := gnmiclient.Update(context.Background(), yclient, path, dst); err != nil {
			log.Warningf("unable to update gRIBI data: %v", err)
		}
	case *aft.Afts_LabelEntry:
		dst := &oc.NetworkInstance_Afts_LabelEntry{}
		if err = convertGoStruct(t, dst, oc.Unmarshal); err != nil {
			break
		}
		var dstLabel oc.NetworkInstance_Afts_LabelEntry_Label_Union
		switch l := t.GetLabel().(type) {
		case aft.E_MplsTypes_MplsLabel_Enum:
			dstLabel = oc.E_MplsTypes_MplsLabel_Enum(l)
		case aft.UnionUint32:
			dstLabel = oc.UnionUint32(l)
		default:
			return fmt.Errorf("Unhandled Label entry type")
		}
		path := ocpath.Root().NetworkInstance(ni).Afts().LabelEntry(dstLabel).State()
		if _, err := gnmiclient.Update(context.Background(), yclient, path, dst); err != nil {
			log.Warningf("unable to update gRIBI data: %v", err)
		}
	default:
		return fmt.Errorf("unrecognized GoStruct type: %T", e)
	}
	return nil
}
