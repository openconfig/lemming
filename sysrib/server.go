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
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/openconfig/gribigo/afthelper"
	"github.com/openconfig/ygnmi/ygnmi"
	"github.com/wenovus/gobgp/v3/pkg/zebra"
	"golang.org/x/exp/maps"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/openconfig/lemming/dataplane/handlers"
	"github.com/openconfig/lemming/gnmi/oc"
	"github.com/openconfig/lemming/gnmi/oc/ocpath"

	gpb "github.com/openconfig/gnmi/proto/gnmi"
	dpb "github.com/openconfig/lemming/proto/dataplane"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
	pb "github.com/openconfig/lemming/proto/sysrib"
)

const (
	// SockAddr is the unix domain socket address for the Sysrib API.
	SockAddr  = "/tmp/sysrib.api"
	defaultNI = "DEFAULT"

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
	pb.UnimplementedSysribServer // For forward-compatibility

	// rib contains the current set of all raw routes from the routing
	// clients as well as the configured connected prefixes from the
	// dataplane.
	rib *SysRIB

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

	dataplane dataplaneAPI

	zServer *ZServer
}

type dataplaneAPI interface {
	ProgramRoute(*ResolvedRoute) error
}

// Dataplane is a wrapper around dpb.HALClient to enable testing before
// resolved route translation.
//
// TODO(wenbli): This is a temporary workaround due to the instability of the
// API. Once the dataplane API is stable, then we'll want to test at the API
// layer instead.
type Dataplane struct {
	Client *ygnmi.Client
}

// programRoute programs the route in the dataplane, returning an error on failure.
func (d *Dataplane) ProgramRoute(r *ResolvedRoute) error {
	log.V(1).Infof("sysrib: programming resolved route: %+v", r)
	rr, err := resolvedRouteToRouteRequest(r)
	if err != nil {
		return err
	}
	_, err = ygnmi.Replace(context.TODO(), d.Client, handlers.RouteQuery(rr.GetVrf(), rr.GetPrefix()), rr, ygnmi.WithSetFallbackEncoding())
	return err
}

// New instantiates server to handle client queries.
//
// If dp is nil, then a connection attempt is made.
func New() (*Server, error) {
	rib, err := NewSysRIB(nil)
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
func (s *Server) Start(gClient gpb.GNMIClient, target, zapiURL string) error {
	if s == nil {
		return errors.New("cannot start nil sysrib server")
	}

	yclient, err := ygnmi.NewClient(gClient, ygnmi.WithTarget(target))
	if err != nil {
		return err
	}
	s.dataplane = &Dataplane{Client: yclient}

	if err := s.monitorConnectedIntfs(yclient); err != nil {
		return err
	}

	if err := s.monitorBGPGUEPolicies(yclient); err != nil {
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
	pb.RegisterSysribServer(grpcServer, s)

	go grpcServer.Serve(lis)

	// Start ZAPI server.
	if zapiURL != "" {
		if s.zServer, err = StartZServer(zapiURL, 0, s); err != nil {
			return err
		}
	}

	return nil
}

// Stop stops the sysrib server.
func (s *Server) Stop() {
	if s.zServer != nil {
		s.zServer.Stop()
	}
}

// monitorConnectedIntfs starts a gothread to check for connected prefixes from
// connected interfaces and adds them to the sysrib. It returns an error if
// there is an error before monitoring can begin.
//
// - gnmiServerAddr is the address of the central gNMI datastore.
// - target is the name of the gNMI target.
func (s *Server) monitorConnectedIntfs(yclient *ygnmi.Client) error {
	b := &ocpath.Batch{}
	b.AddPaths(
		ocpath.Root().InterfaceAny().State().PathStruct(),
	)

	interfaceWatcher := ygnmi.Watch(
		context.Background(),
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
						s.setInterface(name, int32(ifindex), intf.GetEnabled())
						// TODO(wenbli): Support other VRFs.
						if subintf := intf.GetSubinterface(0); subintf != nil {
							for _, addr := range subintf.GetOrCreateIpv4().Address {
								if addr.Ip != nil && addr.PrefixLength != nil {
									if err := s.addInterfacePrefix(name, int32(ifindex), fmt.Sprintf("%s/%d", addr.GetIp(), addr.GetPrefixLength()), defaultNI); err != nil {
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

	// TODO(wenbli): Ideally, this is implemented by watching more fine-grained paths.
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
func (s *Server) monitorBGPGUEPolicies(yclient *ygnmi.Client) error {
	b := &ocpath.Batch{}
	b.AddPaths(
		ocpath.Root().BgpGuePolicyAny().Prefix().Config().PathStruct(),
		ocpath.Root().BgpGuePolicyAny().DstPort().Config().PathStruct(),
		ocpath.Root().BgpGuePolicyAny().SrcIp().Config().PathStruct(),
	)

	bgpGUEPolicyWatcher := ygnmi.Watch(
		context.Background(),
		yclient,
		b.Config(),
		func(root *ygnmi.Value[*oc.Root]) error {
			rootVal, ok := root.Val()
			if !ok {
				return ygnmi.Continue
			}

			policiesFound := map[string]bool{}
			// Add new/updated policies.
			for nonCanonicalPrefix, ocPolicy := range rootVal.BgpGuePolicy {
				// TODO(wenbli): Support other VRFs.
				pfx, err := netip.ParsePrefix(nonCanonicalPrefix)
				if err != nil {
					// TODO(wenbli): This should be a Reconciler.Validate checker function.
					log.Errorf("BGP GUE Policy prefix cannot be parsed: %v", err)
					continue
				}
				pfx = pfx.Masked() // make canonical
				prefix := pfx.String()
				if ocPolicy.DstPort == nil || ocPolicy.SrcIp == nil {
					continue // Wait for complete configuration to arrive.
				}
				policy := GUEPolicy{
					isV6:    pfx.Addr().Is6(),
					dstPort: *ocPolicy.DstPort,
				}
				addr, err := netip.ParseAddr(*ocPolicy.SrcIp)
				if err != nil {
					log.Errorf("BGP GUE Policy source IP cannot be parsed: %v", err)
					continue
				}
				if policy.isV6 {
					if addr.Is4() {
						// TODO(wenbli): This should be a Reconciler.Validate checker function.
						log.Errorf("BGP GUE Policy prefix is IPv6 (%s) but source IP is IPv4: %s", prefix, *ocPolicy.SrcIp)
						continue
					}
					policy.srcIP6 = addr.As16()
				} else {
					if addr.Is6() {
						// TODO(wenbli): This should be a Reconciler.Validate checker function.
						log.Errorf("BGP GUE Policy prefix is IPv4 (%s) but source IP is IPv6: %s", prefix, *ocPolicy.SrcIp)
						continue
					}
					policy.srcIP4 = addr.As4()
				}
				policiesFound[prefix] = true
				if existingPolicy := s.bgpGUEPolicies[prefix]; policy != existingPolicy {
					log.V(1).Infof("Adding new/updated BGP GUE policy: %s: %v", prefix, policy)
					if err := s.setGUEPolicy(prefix, policy); err != nil {
						log.Errorf("Failed while setting BGP GUE Policy: %v", err)
					} else {
						s.bgpGUEPolicies[prefix] = policy
					}
				}
			}

			// Delete incomplete/non-existent policies.
			for prefix := range s.bgpGUEPolicies {
				if _, ok := policiesFound[prefix]; !ok {
					log.Infof("Deleting incomplete/non-existent policy: %s", prefix)
					if err := s.deleteGUEPolicy(prefix); err != nil {
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
		return defaultNI
	default:
		return strconv.Itoa(int(vrfID))
	}
}

func niNameToVrfID(niName string) (uint32, error) {
	switch niName {
	case defaultNI:
		return 0, nil
	default:
		// TODO(wenbli): This mapping should probably be stored in a map.
		return 1, fmt.Errorf("sysrib: only %s VRF is recognized", defaultNI)
	}
}

func prefixString(prefix *pb.Prefix) (string, error) {
	switch fam := prefix.GetFamily(); fam {
	case pb.Prefix_FAMILY_IPV4:
		// TODO(wenbli): Handle invalid input values.
		return fmt.Sprintf("%s/%d", prefix.GetAddress(), prefix.GetMaskLength()), nil
	default:
		return "", fmt.Errorf("unrecognized prefix family: %v", fam)
	}
}

// gueActions generates the forwarding actions that encapsulates a packet with
// a UDP and then an IP header using the information from gueHeaders.
func gueActions(gueHeaders GUEHeaders) ([]*fwdpb.ActionDesc, error) {
	var ip gopacket.SerializableLayer
	ip = &layers.IPv4{
		Version:  4,
		IHL:      5,
		Protocol: layers.IPProtocolUDP,
		SrcIP:    gueHeaders.srcIP4[:],
		DstIP:    gueHeaders.dstIP4[:],
	}
	if gueHeaders.isV6 {
		ip = &layers.IPv6{
			Version:    6,
			NextHeader: layers.IPProtocolUDP,
			SrcIP:      gueHeaders.srcIP6[:],
			DstIP:      gueHeaders.dstIP6[:],
		}
	}
	udp := &layers.UDP{
		SrcPort: 0, // TODO(wenbli): Implement hashing for srcPort.
		DstPort: layers.UDPPort(gueHeaders.dstPort),
		Length:  34,
	}

	buf := gopacket.NewSerializeBuffer()
	if err := gopacket.SerializeLayers(buf, gopacket.SerializeOptions{}, ip, udp); err != nil {
		return nil, fmt.Errorf("failed to serialize GUE headers: %v", err)
	}

	return []*fwdpb.ActionDesc{{
		ActionType: fwdpb.ActionType_ACTION_TYPE_REPARSE,
		Action: &fwdpb.ActionDesc_Reparse{
			Reparse: &fwdpb.ReparseActionDesc{
				HeaderId: fwdpb.PacketHeaderId_PACKET_HEADER_ID_IP4,
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

	var nexthops []*dpb.NextHop
	for nh := range r.Nexthops {
		var actions []*fwdpb.ActionDesc
		if nh.HasGUE() {
			if actions, err = gueActions(nh.GUEHeaders); err != nil {
				return nil, fmt.Errorf("error retrieving GUE actions: %v", err)
			}
		}
		nexthops = append(nexthops, &dpb.NextHop{
			Port:               nh.Port.Name,
			Ip:                 nh.Address,
			Weight:             nh.Weight,
			PreTransmitActions: actions,
		})
	}

	return &dpb.Route{
		Vrf:      uint64(vrfID),
		Prefix:   r.Prefix,
		NextHops: nexthops,
	}, nil
}

// programRoute programs the route in the dataplane, returning an error on failure.
func (s *Server) programRoute(r *ResolvedRoute) error {
	return s.dataplane.ProgramRoute(r)
}

// convertToZAPIRoute converts a route to a ZAPI route for redistributing to
// other protocols (e.g. BGP).
//
// TODO(wenbli): Add support for converting an IPv6 route.
func convertToZAPIRoute(routeKey RouteKey, route *Route) (*zebra.IPRouteBody, error) {
	if route.Connected != nil {
		// TODO(wenbli): Connected routes not supported. This is not
		// needed right now since only need to redistribute
		// non-connected routes.
		return nil, nil
	}
	vrfID, err := niNameToVrfID(routeKey.NIName)
	if err != nil {
		return nil, err
	}

	_, ipv4Net, err := net.ParseCIDR(route.Prefix)
	if err != nil {
		return nil, fmt.Errorf("gribigo/zapi: %v", err)
	}
	prefixLen, _ := ipv4Net.Mask.Size()

	var nexthops []zebra.Nexthop
	for _, nh := range route.NextHops {
		nexthops = append(nexthops, zebra.Nexthop{
			VrfID:  vrfID,
			Gate:   net.ParseIP(nh.Address),
			Weight: uint32(nh.Weight),
		})
	}

	return &zebra.IPRouteBody{
		Flags:   zebra.FlagAllowRecursion,
		Type:    zebra.RouteStatic,
		Safi:    zebra.SafiUnicast,
		Message: zebra.MessageNexthop,
		Prefix: zebra.Prefix{
			Prefix:    ipv4Net.IP.To4(),
			PrefixLen: uint8(prefixLen),
		},
		Nexthops: nexthops,
		Distance: route.RoutePref.AdminDistance,
	}, nil
}

// ResolveAndProgramDiff walks through each prefix in the RIB, resolving it and
// programs the forwarding plane.
//
// TODO(wenbli): handle route deletion.
func (s *Server) ResolveAndProgramDiff() error {
	log.Info("Recalculating resolved RIB")
	s.rib.mu.RLock()
	defer s.rib.mu.RUnlock()
	// newResolvedRoutes keeps track of the new set of top-level resolved
	// routes after re-processing.
	newResolvedRoutes := map[RouteKey]*Route{}
	for niName, ni := range s.rib.NI {
		for it := ni.IPV4.Iterate(); it.Next(); {
			log.V(1).Infof("Iterating at prefix %v out of %d tags", it.Address().String(), ni.IPV4.CountTags())
			_, prefix, err := net.ParseCIDR(it.Address().String())
			if err != nil {
				log.Errorf("sysrib: %v", err)
				continue
			}
			nhs, route, err := s.rib.EgressNexthops(niName, prefix, s.interfaces)
			if err != nil {
				log.Errorf("sysrib: %v", err)
				continue
			}
			routeIsResolved := len(nhs) > 0

			rr := &ResolvedRoute{
				RouteKey: RouteKey{
					// TODO(wenbli): Could it.Address() be different from prefix.String()?
					Prefix: prefix.String(),
					NIName: niName,
				},
				Nexthops: nhs,
			}
			if routeIsResolved {
				newResolvedRoutes[rr.RouteKey] = route
			}

			s.programmedRoutesMu.Lock()
			currentRoute, ok := s.programmedRoutes[rr.RouteKey]
			s.programmedRoutesMu.Unlock()
			switch {
			case !ok && routeIsResolved, ok && !reflect.DeepEqual(currentRoute, rr):
				if err := s.programRoute(rr); err != nil {
					log.Warningf("failed to program route %+v: %v", rr, err)
					continue
				}
				s.programmedRoutesMu.Lock()
				s.programmedRoutes[rr.RouteKey] = rr
				s.programmedRoutesMu.Unlock()
				// ZAPI: If a new/updated route is programmed, redistribute it to clients.
				// TODO(wenbli): RedistributeRouteDel
				zrouteBody, err := convertToZAPIRoute(rr.RouteKey, route)
				if err != nil {
					log.Warningf("failed to convert resolved route to zebra BGP route: %v", err)
				}
				if zrouteBody != nil && s.zServer != nil {
					log.V(1).Info("Sending new route to ZAPI clients: ", zrouteBody)
					s.zServer.ClientMutex.RLock()
					for conn := range s.zServer.ClientMap {
						serverSendMessage(conn, zebra.RedistributeRouteAdd, zrouteBody)
					}
					s.zServer.ClientMutex.RUnlock()
				}
			default:
				// No diff, so don't do anything.
			}
		}
	}
	s.resolvedRoutesMu.Lock()
	s.resolvedRoutes = newResolvedRoutes
	s.resolvedRoutesMu.Unlock()
	return nil
}

// ResolvedRoutes returns the shallow copy of the resolved routes of the RIB
// manager.
func (s *Server) ResolvedRoutes() map[RouteKey]*Route {
	s.resolvedRoutesMu.Lock()
	defer s.resolvedRoutesMu.Unlock()
	return maps.Clone(s.resolvedRoutes)
}

// SetRoute implements ROUTE_ADD and ROUTE_DELETE
func (s *Server) SetRoute(_ context.Context, req *pb.SetRouteRequest) (*pb.SetRouteResponse, error) {
	// TODO(wenbli): Handle route deletion.
	if req.Delete {
		return nil, status.Errorf(codes.Unimplemented, "route delete not yet supported")
	}

	pfx, err := prefixString(req.Prefix)
	if err != nil {
		return nil, err
	}

	nexthops := []*afthelper.NextHopSummary{}
	for _, nh := range req.GetNexthops() {
		if nh.GetType() != pb.Nexthop_TYPE_IPV4 {
			return nil, status.Errorf(codes.Unimplemented, "non-IPV4 nexthop not supported")
		}
		nexthops = append(nexthops, &afthelper.NextHopSummary{
			Weight:          nh.GetWeight(),
			Address:         nh.GetAddress(),
			NetworkInstance: vrfIDToNiName(nh.GetVrfId()),
		})
	}

	// TODO(wenbli): Check if recursive resolution is an infinite recursion. This happens if there is a cycle.

	niName := vrfIDToNiName(req.GetVrfId())
	if err := s.setRoute(niName, &Route{
		// TODO(wenbli): check if pfx has to be canonical or does it tolerate it: i.e. 1.1.1.0/24 instead of 1.1.1.1/24
		Prefix:   pfx,
		NextHops: nexthops,
		RoutePref: RoutePreference{
			AdminDistance: uint8(req.GetAdminDistance()),
			Metric:        req.GetMetric(),
		},
	}); err != nil {
		return nil, status.Error(codes.Aborted, fmt.Sprint(err))
	}

	// There could be operations carried out by ResolveAndProgramDiff() other than the input route, so we look up our particular prefix.
	status := pb.SetRouteResponse_STATUS_FAIL
	s.programmedRoutesMu.Lock()
	if _, ok := s.programmedRoutes[RouteKey{Prefix: pfx, NIName: niName}]; ok {
		status = pb.SetRouteResponse_STATUS_SUCCESS
	}
	s.programmedRoutesMu.Unlock()
	return &pb.SetRouteResponse{
		Status: status,
	}, nil
}

// setRoute adds/deletes a route from the RIB manager.
func (s *Server) setRoute(niName string, route *Route) error {
	if _, err := s.rib.AddRoute(niName, route); err != nil {
		return fmt.Errorf("error while adding route to sysrib: %v", err)
	}

	if err := s.ResolveAndProgramDiff(); err != nil {
		return fmt.Errorf("error while resolving sysrib: %v", err)
	}
	return nil
}

// setZebraRoute calls setRoute after reformatting a zebra-formatted input route.
func (s *Server) setZebraRoute(niName string, zroute *zebra.IPRouteBody) error {
	if s == nil {
		return fmt.Errorf("cannot add route to nil sysrib server")
	}
	log.V(1).Infof("setZebraRoute: %+v", *zroute)
	route := convertZebraRoute(niName, zroute)
	if err := s.setRoute(niName, route); err != nil {
		return err
	}
	return nil
}

// convertZebraRoute converts a zebra route to a Sysrib route.
func convertZebraRoute(niName string, zroute *zebra.IPRouteBody) *Route {
	var nexthops []*afthelper.NextHopSummary
	for _, znh := range zroute.Nexthops {
		nexthops = append(nexthops, &afthelper.NextHopSummary{
			Weight:          1,
			Address:         znh.Gate.String(),
			NetworkInstance: niName,
		})
	}
	var routePref RoutePreference
	switch zroute.Type {
	case zebra.RouteBGP:
		routePref.AdminDistance = AdminDistanceBGP
	}
	routePref.Metric = zroute.Metric
	return &Route{
		Prefix: fmt.Sprintf("%s/%d", zroute.Prefix.Prefix.String(), zroute.Prefix.PrefixLen),
		// NextHops is the set of IP nexthops that the route uses if
		// it is not a connected route.
		NextHops:  nexthops,
		RoutePref: routePref,
	}
}

// addInterfacePrefix adds a prefix to the sysrib as a connected route.
func (s *Server) addInterfacePrefix(name string, ifindex int32, prefix string, niName string) error {
	log.V(1).Infof("Adding interface prefix: intf %s, idx %d, prefix %s, ni %s", name, ifindex, prefix, niName)
	connectedRoute := &Route{
		Prefix: prefix,
		Connected: &Interface{
			Name:  name,
			Index: ifindex,
		},
		RoutePref: RoutePreference{
			// Connected routes have admin-distance of 0.
			AdminDistance: 0,
		},
	}

	if _, err := s.rib.AddRoute(niName, connectedRoute); err != nil {
		return fmt.Errorf("failed to add route to Sysrib: %v", err)
	}
	return s.ResolveAndProgramDiff()
}

// setInterface responds to INTERFACE_UP/INTERFACE_DOWN messages from the dataplane.
func (s *Server) setInterface(name string, ifindex int32, enabled bool) error {
	log.V(1).Infof("Setting interface %q(%d) to enabled=%v", name, ifindex, enabled)
	s.interfaces[Interface{
		Name:  name,
		Index: ifindex,
	}] = enabled

	return s.ResolveAndProgramDiff()
}

// TODO(wenbli): Do we need to handle interface deletion?
// This is not required in the MVP since basic tests will just need to enable/disable interfaces.

// setGUEPolicy adds a new GUE policy and triggers resolved route
// computation and programming.
func (s *Server) setGUEPolicy(prefix string, policy GUEPolicy) error {
	if err := s.rib.SetGUEPolicy(prefix, policy); err != nil {
		return fmt.Errorf("error while adding route to sysrib: %v", err)
	}

	if err := s.ResolveAndProgramDiff(); err != nil {
		return fmt.Errorf("error while resolving sysrib: %v", err)
	}
	return nil
}

// deleteGUEPolicy adds a new GUE policy and triggers resolved route
// computation and programming.
func (s *Server) deleteGUEPolicy(prefix string) error {
	if _, err := s.rib.DeleteGUEPolicy(prefix); err != nil {
		return fmt.Errorf("error while adding route to sysrib: %v", err)
	}

	if err := s.ResolveAndProgramDiff(); err != nil {
		return fmt.Errorf("error while resolving sysrib: %v", err)
	}
	return nil
}
