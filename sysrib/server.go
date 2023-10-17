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

package sysrib

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/netip"
	"os"
	"reflect"
	"strconv"
	"sync"

	log "github.com/golang/glog"
	"github.com/google/go-cmp/cmp"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/openconfig/gribigo/afthelper"
	"github.com/openconfig/ygnmi/ygnmi"
	"golang.org/x/exp/maps"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/openconfig/lemming/dataplane/handlers"
	"github.com/openconfig/lemming/gnmi/fakedevice"
	"github.com/openconfig/lemming/gnmi/oc"
	"github.com/openconfig/lemming/gnmi/oc/ocpath"

	gpb "github.com/openconfig/gnmi/proto/gnmi"

	dpb "github.com/openconfig/lemming/proto/dataplane"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
	sysribpb "github.com/openconfig/lemming/proto/sysrib"
)

const (
	// SockAddr is the unix domain socket address for the Sysrib API.
	SockAddr = "/tmp/sysrib.api"

	// ZAPIAddr is the connection address for ZAPI, which is in the form
	// of "type:address", where type can either be a unix or tcp socket.
	ZAPIAddr = "unix:/var/run/zserv.api"
)

// AdminDistance is the admin-distance of a routing protocol. See
// https://docs.frrouting.org/en/latest/zebra.html#administrative-distance
const (
	AdminDistanceConnected = 0
	AdminDistanceStatic    = 1
	AdminDistanceBGP       = 20
)

// Server is the implementation of the Sysrib API.
//
// API:
// - SetRoute
// - addInterfacePrefix
// - setInterface
type Server struct {
	sysribpb.UnimplementedSysribServer // For forward-compatibility

	// rib contains the current set of all raw routes from the routing
	// clients as well as the configured connected prefixes from the
	// dataplane.
	rib *SysRIB

	interfacesMu sync.Mutex
	// interfaces contains the status of all existing interfaces as
	// indicated by the forwarding plane.
	interfaces map[Interface]bool

	// bgpGUEPolicies contains the current set of BGP GUE policies as
	// received from the configuration. Comparison against this prevents
	// unnecessary sysrib updates.
	bgpGUEPolicies map[string]GUEPolicy

	programmedRoutesMu sync.Mutex
	// programmedRoutes contain a map of resolved routes with which to do
	// diff for sending to the dataplane for programming.
	programmedRoutes map[RouteKey]*ResolvedRoute

	resolvedRoutesMu sync.Mutex
	// resolvedRoutes keeps track of all the original routes that have now
	// been programmed. This allows looking up the top-level routes that
	// have been resolved.
	resolvedRoutes map[RouteKey]*Route

	dataplane dplane

	zServer *ZServer
}

// dplane represents the dataplane API accessible to sysrib for programming
// routes.
type dplane struct {
	Client *ygnmi.Client
}

// programRoute programs the route in the dataplane, returning an error on failure.
func (d *dplane) programRoute(ctx context.Context, r *ResolvedRoute) error {
	log.V(1).Infof("sysrib: programming resolved route: %+v", r)
	rr, err := resolvedRouteToRouteRequest(r)
	if err != nil {
		return err
	}
	_, err = ygnmi.Replace(ctx, d.Client, handlers.RouteQuery(rr.GetPrefix().GetVrfId(), rr.GetPrefix().GetCidr()), rr, ygnmi.WithSetFallbackEncoding())
	return err
}

// deprogramRoute de-programs the route in the dataplane, returning an error on failure.
func (d *dplane) deprogramRoute(ctx context.Context, r *ResolvedRoute) error {
	log.V(1).Infof("sysrib: deprogramming newly unresolved route: %+v", r)
	rr, err := resolvedRouteToRouteRequest(r)
	if err != nil {
		return err
	}
	_, err = ygnmi.Delete(ctx, d.Client, handlers.RouteQuery(rr.GetPrefix().GetVrfId(), rr.GetPrefix().GetCidr()))
	return err
}

// New instantiates server to handle client queries.
//
// If dp is nil, then a connection attempt is made.
func New(cfg *oc.Root) (*Server, error) {
	rib, err := NewSysRIB(cfg)
	if err != nil {
		return nil, err
	}

	s := &Server{
		rib:              rib,
		interfaces:       map[Interface]bool{},
		bgpGUEPolicies:   map[string]GUEPolicy{},
		programmedRoutes: map[RouteKey]*ResolvedRoute{},
		resolvedRoutes:   map[RouteKey]*Route{},
	}
	return s, nil
}

// Start starts the sysrib gRPC service at a unix domain socket. This
// should be started prior to routing services to allow them to connect to
// sysrib during their initialization.
//
// - If zapiURL is not specified, then the ZAPI server will not be started.
func (s *Server) Start(ctx context.Context, gClient gpb.GNMIClient, target, zapiURL string) error {
	if s == nil {
		return errors.New("cannot start nil sysrib server")
	}

	yclient, err := ygnmi.NewClient(gClient, ygnmi.WithTarget(target), ygnmi.WithRequestLogLevel(2))
	if err != nil {
		return err
	}
	s.dataplane = dplane{Client: yclient}

	if err := s.monitorConnectedIntfs(ctx, yclient); err != nil {
		return err
	}

	if err := s.monitorBGPGUEPolicies(ctx, yclient); err != nil {
		return err
	}

	if err := s.monitorStaticRoutes(ctx, yclient); err != nil {
		return err
	}

	if err := os.RemoveAll(SockAddr); err != nil {
		return err
	}
	lis, err := net.Listen("unix", SockAddr)
	if err != nil {
		log.Fatalf("listen error: %v", err)
	}

	grpcServer := grpc.NewServer()
	sysribpb.RegisterSysribServer(grpcServer, s)

	go grpcServer.Serve(lis)

	// BEGIN Start ZAPI server.
	if zapiURL != "" {
		if s.zServer, err = StartZServer(zapiURL, 0, s); err != nil {
			return err
		}
	}
	// END Start ZAPI server.

	return nil
}

// Stop stops the sysrib server.
func (s *Server) Stop() {
	s.zServer.Stop()
}

// monitorConnectedIntfs starts a gothread to check for connected prefixes from
// connected interfaces and adds them to the sysrib. It returns an error if
// there is an error before monitoring can begin.
//
// - gnmiServerAddr is the address of the central gNMI datastore.
// - target is the name of the gNMI target.
func (s *Server) monitorConnectedIntfs(ctx context.Context, yclient *ygnmi.Client) error {
	b := &ocpath.Batch{}
	b.AddPaths(
		ocpath.Root().InterfaceAny().Enabled().State().PathStruct(),
		ocpath.Root().InterfaceAny().Ifindex().State().PathStruct(),
		ocpath.Root().InterfaceAny().Subinterface(0).Ipv4().AddressAny().Ip().State().PathStruct(),
		ocpath.Root().InterfaceAny().Subinterface(0).Ipv4().AddressAny().PrefixLength().State().PathStruct(),
		ocpath.Root().InterfaceAny().Subinterface(0).Ipv6().AddressAny().Ip().State().PathStruct(),
		ocpath.Root().InterfaceAny().Subinterface(0).Ipv6().AddressAny().PrefixLength().State().PathStruct(),
	)

	interfaceWatcher := ygnmi.Watch(
		ctx,
		yclient,
		b.State(),
		func(root *ygnmi.Value[*oc.Root]) error {
			rootVal, ok := root.Val()
			if !ok {
				return ygnmi.Continue
			}
			for name, intf := range rootVal.Interface {
				if intf.Enabled != nil {
					if intf.Ifindex != nil {
						ifindex := intf.GetIfindex()
						s.setInterface(ctx, name, int32(ifindex), intf.GetEnabled())
						// TODO(wenbli): Support other VRFs.
						if subintf := intf.GetSubinterface(0); subintf != nil {
							for _, addr := range subintf.GetOrCreateIpv4().Address {
								if addr.Ip != nil && addr.PrefixLength != nil {
									if err := s.addInterfacePrefix(ctx, name, int32(ifindex), fmt.Sprintf("%s/%d", addr.GetIp(), addr.GetPrefixLength()), fakedevice.DefaultNetworkInstance); err != nil {
										log.Warningf("adding interface prefix failed: %v", err)
									}
								}
							}
							for _, addr := range subintf.GetOrCreateIpv6().Address {
								if addr.Ip != nil && addr.PrefixLength != nil {
									if err := s.addInterfacePrefix(ctx, name, int32(ifindex), fmt.Sprintf("%s/%d", addr.GetIp(), addr.GetPrefixLength()), fakedevice.DefaultNetworkInstance); err != nil {
										log.Warningf("adding interface prefix failed: %v", err)
									}
								}
							}
						}
					}
				}
			}
			return ygnmi.Continue
		},
	)

	// TODO(wenbli): Support interface removal.
	go func() {
		if _, err := interfaceWatcher.Await(); err != nil {
			log.Warningf("Sysrib interface watcher has stopped: %v", err)
		}
	}()
	return nil
}

// monitorBGPGUEPolicies starts a gothread to check for BGP GUE policies being
// added or deleted from the config, and informs the sysrib server accordingly
// to update programmed routes.
func (s *Server) monitorBGPGUEPolicies(ctx context.Context, yclient *ygnmi.Client) error {
	b := &ocpath.Batch{}
	b.AddPaths(
		ocpath.Root().BgpGueIpv4GlobalPolicyAny().Prefix().Config().PathStruct(),
		ocpath.Root().BgpGueIpv4GlobalPolicyAny().DstPortIpv4().Config().PathStruct(),
		ocpath.Root().BgpGueIpv4GlobalPolicyAny().DstPortIpv6().Config().PathStruct(),
		ocpath.Root().BgpGueIpv4GlobalPolicyAny().SrcIp().Config().PathStruct(),
		ocpath.Root().BgpGueIpv6GlobalPolicyAny().Prefix().Config().PathStruct(),
		ocpath.Root().BgpGueIpv6GlobalPolicyAny().DstPortIpv6().Config().PathStruct(),
		ocpath.Root().BgpGueIpv6GlobalPolicyAny().SrcIp().Config().PathStruct(),
	)

	bgpGUEPolicyWatcher := ygnmi.Watch(
		ctx,
		yclient,
		b.Config(),
		func(root *ygnmi.Value[*oc.Root]) error {
			rootVal, ok := root.Val()
			if !ok {
				return ygnmi.Continue
			}

			policiesFound := map[string]bool{}
			updatePolicy := func(prefix string, policy GUEPolicy) {
				policiesFound[prefix] = true
				if existingPolicy := s.bgpGUEPolicies[prefix]; policy != existingPolicy {
					log.V(1).Infof("Adding new/updated BGP GUE policy: %s: %v", prefix, policy)
					if err := s.setGUEPolicy(ctx, prefix, policy); err != nil {
						log.Errorf("Failed while setting BGP GUE Policy: %v", err)
					} else {
						s.bgpGUEPolicies[prefix] = policy
					}
				}
			}

			// Add new/updated policies.
			for pfx, ocPolicy := range rootVal.BgpGueIpv4GlobalPolicy {
				// TODO(wenbli): Support other VRFs.
				prefix, err := canonicalPrefix(pfx)
				if err != nil {
					// TODO(wenbli): This should be a Reconciler.Validate checker function.
					log.Errorf("BGP GUE Policy prefix cannot be parsed: %v", err)
					continue
				}
				if ocPolicy.DstPortIpv4 == nil || ocPolicy.DstPortIpv6 == nil || ocPolicy.SrcIp == nil {
					continue // Wait for complete configuration to arrive.
				}
				addr, err := netip.ParseAddr(*ocPolicy.SrcIp)
				if err != nil {
					log.Errorf("BGP GUE Policy source IP cannot be parsed: %v", err)
					continue
				}
				updatePolicy(prefix.String(), GUEPolicy{
					dstPortv4: *ocPolicy.DstPortIpv4,
					dstPortv6: *ocPolicy.DstPortIpv6,
					srcIP4:    addr.As4(),
				})
			}
			for pfx, ocPolicy := range rootVal.BgpGueIpv6GlobalPolicy {
				// TODO(wenbli): Support other VRFs.
				prefix, err := canonicalPrefix(pfx)
				if err != nil {
					// TODO(wenbli): This should be a Reconciler.Validate checker function.
					log.Errorf("BGP GUE Policy prefix cannot be parsed: %v", err)
					continue
				}
				if ocPolicy.DstPortIpv6 == nil || ocPolicy.SrcIp == nil {
					continue // Wait for complete configuration to arrive.
				}
				addr, err := netip.ParseAddr(*ocPolicy.SrcIp)
				if err != nil {
					log.Errorf("BGP GUE Policy source IP cannot be parsed: %v", err)
					continue
				}
				updatePolicy(prefix.String(), GUEPolicy{
					dstPortv6: *ocPolicy.DstPortIpv6,
					srcIP6:    addr.As16(),
				})
			}

			// Delete incomplete/non-existent policies.
			for prefix := range s.bgpGUEPolicies {
				if _, ok := policiesFound[prefix]; !ok {
					log.Infof("Deleting incomplete/non-existent policy: %s", prefix)
					if err := s.deleteGUEPolicy(ctx, prefix); err != nil {
						log.Errorf("Failed while deleting BGP GUE Policy: %v", err)
					} else {
						delete(s.bgpGUEPolicies, prefix)
					}
				}
			}
			return ygnmi.Continue
		},
	)

	go func() {
		if _, err := bgpGUEPolicyWatcher.Await(); err != nil {
			log.Warningf("Sysrib BGP GUE policy watcher has stopped: %v", err)
		}
	}()
	return nil
}

// RouteKey is the unique identifier of an IP route.
type RouteKey struct {
	Prefix string
	NIName string
}

// ResolvedRoute represents a route that is ready to be programmed into the forwarding plane.
type ResolvedRoute struct {
	RouteKey

	// NOTE: The order of the nexthops should not matter when being programmed into the forwarding plane. As such, the forwarding plane should sort these nexthops before assigning the hash output for ECMP.
	Nexthops map[ResolvedNexthop]bool
	// TODO(wenbli): backup nexthops.
}

// ResolvedNexthop contains the information required to forward an IP packet.
//
// This type must be hashable, and uniquely identifies nexthops.
type ResolvedNexthop struct {
	afthelper.NextHopSummary

	Port       Interface
	GUEHeaders GUEHeaders
}

// HasGUE returns a bool indicating whether the resolved nexthop contains GUE
// information.
func (nh ResolvedNexthop) HasGUE() bool {
	return nh.GUEHeaders != GUEHeaders{}
}

func vrfIDToNiName(vrfID uint32) string {
	switch vrfID {
	case 0:
		return fakedevice.DefaultNetworkInstance
	default:
		return strconv.Itoa(int(vrfID))
	}
}

func niNameToVrfID(niName string) (uint32, error) {
	switch niName {
	case fakedevice.DefaultNetworkInstance:
		return 0, nil
	default:
		// TODO(wenbli): This mapping should probably be stored in a map.
		return 1, fmt.Errorf("sysrib: only %s VRF is recognized", fakedevice.DefaultNetworkInstance)
	}
}

func prefixString(prefix *sysribpb.Prefix) (string, error) {
	switch fam := prefix.GetFamily(); fam {
	case sysribpb.Prefix_FAMILY_IPV4, sysribpb.Prefix_FAMILY_IPV6:
		// TODO(wenbli): Handle invalid input values.
		return fmt.Sprintf("%s/%d", prefix.GetAddress(), prefix.GetMaskLength()), nil
	default:
		return "", fmt.Errorf("unrecognized prefix family: %v", fam)
	}
}

// gueActions generates the forwarding actions that encapsulates a packet with
// a UDP and then an IP header using the information from gueHeaders.
//
// - isRouteV4 indicates whether the route is a v4/v6-mapped v4 route or a v6 route.
func gueActions(isRouteV4 bool, gueHeaders GUEHeaders) ([]*fwdpb.ActionDesc, error) {
	var ip gopacket.SerializableLayer
	var headerID fwdpb.PacketHeaderId
	if !gueHeaders.isV6 {
		ip = &layers.IPv4{
			Version:  4,
			IHL:      5,
			Protocol: layers.IPProtocolUDP,
			SrcIP:    gueHeaders.srcIP4[:],
			DstIP:    gueHeaders.dstIP4[:],
		}
		headerID = fwdpb.PacketHeaderId_PACKET_HEADER_ID_IP4
	} else {
		ip = &layers.IPv6{
			Version:    6,
			NextHeader: layers.IPProtocolUDP,
			SrcIP:      gueHeaders.srcIP6[:],
			DstIP:      gueHeaders.dstIP6[:],
		}
		headerID = fwdpb.PacketHeaderId_PACKET_HEADER_ID_IP6
	}

	udp := &layers.UDP{
		SrcPort: 0,  // TODO(wenbli): Implement hashing for srcPort.
		Length:  34, // TODO(wenbli): Figure out how to not make this hardcoded.
	}
	if isRouteV4 {
		udp.DstPort = layers.UDPPort(gueHeaders.dstPortv4)
	} else {
		udp.DstPort = layers.UDPPort(gueHeaders.dstPortv6)
	}

	buf := gopacket.NewSerializeBuffer()
	if err := gopacket.SerializeLayers(buf, gopacket.SerializeOptions{}, ip, udp); err != nil {
		return nil, fmt.Errorf("failed to serialize GUE headers: %v", err)
	}

	return []*fwdpb.ActionDesc{{
		ActionType: fwdpb.ActionType_ACTION_TYPE_REPARSE,
		Action: &fwdpb.ActionDesc_Reparse{
			Reparse: &fwdpb.ReparseActionDesc{
				HeaderId: headerID,
				FieldIds: []*fwdpb.PacketFieldId{
					{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_NEXT_HOP_IP}},
					{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_PORT_INPUT}},
					{Field: &fwdpb.PacketField{FieldNum: fwdpb.PacketFieldNum_PACKET_FIELD_NUM_PACKET_PORT_OUTPUT}},
				},
				// After the UDP header, the rest of the packet (original packet) will be classified as payload.
				Prepend: buf.Bytes(),
			},
		},
	}}, nil
}

func resolvedRouteToRouteRequest(r *ResolvedRoute) (*dpb.Route, error) {
	vrfID, err := niNameToVrfID(r.NIName)
	if err != nil {
		return nil, err
	}

	pfx, err := netip.ParsePrefix(r.Prefix)
	if err != nil {
		log.Errorf("Route prefix cannot be parsed: %v", err)
		return nil, err
	}

	var nexthops []*dpb.NextHop
	for nh := range r.Nexthops {
		var actions []*fwdpb.ActionDesc
		if nh.HasGUE() {
			if actions, err = gueActions(pfx.Addr().Is4() || pfx.Addr().Is4In6(), nh.GUEHeaders); err != nil {
				return nil, fmt.Errorf("error retrieving GUE actions: %v", err)
			}
		}
		dnh := &dpb.NextHop{
			Dev: &dpb.NextHop_Port{
				Port: nh.Port.Name,
			},
			Weight:             nh.Weight,
			PreTransmitActions: actions,
		}
		if nh.Address != "" {
			dnh.Ip = &dpb.NextHop_IpStr{IpStr: nh.Address}
		}
		nexthops = append(nexthops, dnh)
	}

	return &dpb.Route{
		Prefix: &dpb.RoutePrefix{
			VrfId: uint64(vrfID),
			Prefix: &dpb.RoutePrefix_Cidr{
				Cidr: r.Prefix,
			},
		},
		Hop: &dpb.Route_NextHops{
			NextHops: &dpb.NextHopList{Hops: nexthops},
		},
	}, nil
}

// ResolveAndProgramDiff walks through each prefix in the RIB, resolving it and
// programs the forwarding plane.
func (s *Server) ResolveAndProgramDiff(ctx context.Context) error {
	log.Info("Recalculating resolved RIB")
	if debugRIB {
		defer s.rib.PrintRIB()
	}
	s.rib.mu.RLock()
	defer s.rib.mu.RUnlock()

	// Program new/changed resolved routes.
	newResolvedRoutes := map[RouteKey]*Route{}
	for niName, ni := range s.rib.NI {
		for it := ni.IPV4.Iterate(); it.Next(); {
			s.resolveAndProgramDiffAux(ctx, niName, ni, it.Address().String(), newResolvedRoutes)
		}
		for it := ni.IPV6.Iterate(); it.Next(); {
			s.resolveAndProgramDiffAux(ctx, niName, ni, it.Address().String(), newResolvedRoutes)
		}
	}

	// Deprogram newly unresolved routes.
	for routeKey, rr := range s.ProgrammedRoutes() {
		if _, ok := newResolvedRoutes[routeKey]; !ok {
			log.V(1).Infof("Deleting route %s", rr.RouteKey)
			if err := s.dataplane.deprogramRoute(ctx, rr); err != nil {
				log.Warningf("failed to deprogram route %+v: %v", rr, err)
				continue
			}
			s.programmedRoutesMu.Lock()
			delete(s.programmedRoutes, rr.RouteKey)
			s.programmedRoutesMu.Unlock()
			// TODO(wenbli): Add ZAPI-handling for sending a route deletion message.
		}
	}
	if debugRIB {
		s.PrintProgrammedRoutes()
	}

	s.resolvedRoutesMu.Lock()
	s.resolvedRoutes = newResolvedRoutes
	s.resolvedRoutesMu.Unlock()
	return nil
}

// resolveAndProgramDiffAux is a helper function for ResolveAndProgramDiff.
//
// It carries out the following functionalities:
// - Resolve a single route specified by prefix and program if it's different.
// - Populate the resolved route into newResolvedRoutes.
func (s *Server) resolveAndProgramDiffAux(ctx context.Context, niName string, ni *NIRIB, prefix string, newResolvedRoutes map[RouteKey]*Route) {
	log.V(1).Infof("Iterating at prefix %v (v4 has %d tags) (v6 has %d tags)", prefix, ni.IPV4.CountTags(), ni.IPV6.CountTags())
	_, pfx, err := net.ParseCIDR(prefix)
	if err != nil {
		log.Errorf("sysrib: %v", err)
		return
	}
	s.interfacesMu.Lock()
	nhs, route, err := s.rib.egressNexthops(niName, pfx, s.interfaces)
	s.interfacesMu.Unlock()
	if err != nil {
		log.Errorf("sysrib: %v", err)
		return
	}
	routeIsResolved := len(nhs) > 0

	cPfx, err := canonicalPrefix(prefix)
	if err != nil {
		log.Errorf("sysrib: %v", err)
	}

	rr := &ResolvedRoute{
		RouteKey: RouteKey{
			Prefix: cPfx.String(),
			NIName: niName,
		},
		Nexthops: nhs,
	}
	if routeIsResolved {
		newResolvedRoutes[rr.RouteKey] = route
	}

	s.programmedRoutesMu.Lock()
	currentRoute, _ := s.programmedRoutes[rr.RouteKey]
	s.programmedRoutesMu.Unlock()
	switch {
	case routeIsResolved && !reflect.DeepEqual(currentRoute, rr):
		log.V(1).Infof("(-currentRoute, +resolvedRoute):\n%s", cmp.Diff(currentRoute, rr))
		if err := s.dataplane.programRoute(ctx, rr); err != nil {
			log.Warningf("failed to program route %+v: %v", rr, err)
			return
		}
		s.programmedRoutesMu.Lock()
		s.programmedRoutes[rr.RouteKey] = rr
		s.programmedRoutesMu.Unlock()
		if debugRIB {
			s.PrintProgrammedRoutes()
		}
		// ZAPI: If a new/updated route is programmed, redistribute it to clients.
		distributeRoute(s.zServer, rr, route)
	}
}

// ResolvedRoutes returns the shallow copy of the resolved routes of the RIB
// manager.
func (s *Server) ResolvedRoutes() map[RouteKey]*Route {
	s.resolvedRoutesMu.Lock()
	defer s.resolvedRoutesMu.Unlock()
	return maps.Clone(s.resolvedRoutes)
}

// ProgrammedRoutes returns the shallow copy of the programmed routes of the RIB
// manager.
func (s *Server) ProgrammedRoutes() map[RouteKey]*ResolvedRoute {
	s.programmedRoutesMu.Lock()
	defer s.programmedRoutesMu.Unlock()
	return maps.Clone(s.programmedRoutes)
}

// SetRoute implements ROUTE_ADD and ROUTE_DELETE
func (s *Server) SetRoute(ctx context.Context, req *sysribpb.SetRouteRequest) (*sysribpb.SetRouteResponse, error) {
	pfx, err := prefixString(req.Prefix)
	if err != nil {
		return nil, err
	}

	nexthops := []*afthelper.NextHopSummary{}
	for _, nh := range req.GetNexthops() {
		if nh.GetType() != sysribpb.Nexthop_TYPE_IPV4 && nh.GetType() != sysribpb.Nexthop_TYPE_IPV6 {
			return nil, status.Errorf(codes.Unimplemented, "Unrecognized nexthop type: %s", nh.GetType())
		}
		nexthops = append(nexthops, &afthelper.NextHopSummary{
			Weight:          nh.GetWeight(),
			Address:         nh.GetAddress(),
			NetworkInstance: vrfIDToNiName(nh.GetVrfId()),
		})
	}

	// TODO(wenbli): Check if recursive resolution is an infinite recursion. This happens if there is a cycle.

	niName := vrfIDToNiName(req.GetVrfId())
	if err := s.setRoute(ctx, niName, &Route{
		// TODO(wenbli): check if pfx has to be canonical or does it tolerate it: i.e. 1.1.1.0/24 instead of 1.1.1.1/24
		Prefix:   pfx,
		NextHops: nexthops,
		RoutePref: RoutePreference{
			AdminDistance: uint8(req.GetAdminDistance()),
			Metric:        req.GetMetric(),
		},
	}, req.Delete); err != nil {
		return nil, status.Error(codes.Aborted, fmt.Sprint(err))
	}

	// There could be operations carried out by ResolveAndProgramDiff() other than the input route, so we look up our particular prefix.
	status := sysribpb.SetRouteResponse_STATUS_FAIL
	s.programmedRoutesMu.Lock()
	if _, ok := s.programmedRoutes[RouteKey{Prefix: pfx, NIName: niName}]; ok {
		status = sysribpb.SetRouteResponse_STATUS_SUCCESS
	}
	s.programmedRoutesMu.Unlock()
	return &sysribpb.SetRouteResponse{
		Status: status,
	}, nil
}

// setRoute adds/deletes a route from the RIB manager.
func (s *Server) setRoute(ctx context.Context, niName string, route *Route, isDelete bool) error {
	if err := s.rib.setRoute(niName, route, isDelete); err != nil {
		return fmt.Errorf("error while adding route to sysrib: %v", err)
	}

	if err := s.ResolveAndProgramDiff(ctx); err != nil {
		return fmt.Errorf("error while resolving sysrib: %v", err)
	}
	return nil
}

// addInterfacePrefix adds a prefix to the sysrib as a connected route.
func (s *Server) addInterfacePrefix(ctx context.Context, name string, ifindex int32, prefix string, niName string) error {
	log.V(1).Infof("Adding interface prefix: intf %s, idx %d, prefix %s, ni %s", name, ifindex, prefix, niName)
	return s.setRoute(ctx, niName, &Route{
		Prefix: prefix,
		Connected: &Interface{
			Name:  name,
			Index: ifindex,
		},
		RoutePref: RoutePreference{
			// Connected routes have admin-distance of 0.
			AdminDistance: 0,
		},
	}, false)
}

// setInterface responds to INTERFACE_UP/INTERFACE_DOWN messages from the dataplane.
func (s *Server) setInterface(ctx context.Context, name string, ifindex int32, enabled bool) error {
	log.V(1).Infof("Setting interface %q(%d) to enabled=%v", name, ifindex, enabled)
	s.interfacesMu.Lock()
	s.interfaces[Interface{
		Name:  name,
		Index: ifindex,
	}] = enabled
	s.interfacesMu.Unlock()

	return s.ResolveAndProgramDiff(ctx)
}

// TODO(wenbli): Do we need to handle interface deletion?
// This is not required in the MVP since basic tests will just need to enable/disable interfaces.

// setGUEPolicy adds a new GUE policy and triggers resolved route
// computation and programming.
func (s *Server) setGUEPolicy(ctx context.Context, prefix string, policy GUEPolicy) error {
	if err := s.rib.SetGUEPolicy(prefix, policy); err != nil {
		return fmt.Errorf("error while adding route to sysrib: %v", err)
	}

	if err := s.ResolveAndProgramDiff(ctx); err != nil {
		return fmt.Errorf("error while resolving sysrib: %v", err)
	}
	return nil
}

// deleteGUEPolicy adds a new GUE policy and triggers resolved route
// computation and programming.
func (s *Server) deleteGUEPolicy(ctx context.Context, prefix string) error {
	if _, err := s.rib.DeleteGUEPolicy(prefix); err != nil {
		return fmt.Errorf("error while adding route to sysrib: %v", err)
	}

	if err := s.ResolveAndProgramDiff(ctx); err != nil {
		return fmt.Errorf("error while resolving sysrib: %v", err)
	}
	return nil
}
