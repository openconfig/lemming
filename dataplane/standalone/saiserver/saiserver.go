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

package saiserver

import (
	"google.golang.org/grpc"

	saipb "github.com/openconfig/lemming/dataplane/standalone/proto"
)

type acl struct {
	saipb.UnimplementedAclServer
}

type bfd struct {
	saipb.UnimplementedBfdServer
}

type buffer struct {
	saipb.UnimplementedBufferServer
}

type bridge struct {
	saipb.UnimplementedBridgeServer
}

type counter struct {
	saipb.UnimplementedCounterServer
}

type debugCounter struct {
	saipb.UnimplementedDebugCounterServer
}

type fdb struct {
	saipb.UnimplementedFdbServer
}

type hash struct {
	saipb.UnimplementedHashServer
}

type hostif struct {
	saipb.UnimplementedHostifServer
}

type ipmcGroup struct {
	saipb.UnimplementedIpmcGroupServer
}

type ipmc struct {
	saipb.UnimplementedIpmcServer
}

type ipsec struct {
	saipb.UnimplementedIpsecServer
}

type isolationGroup struct {
	saipb.UnimplementedIsolationGroupServer
}

type l2mcGroup struct {
	saipb.UnimplementedL2McGroupServer
}

type l2mc struct {
	saipb.UnimplementedL2McServer
}

type lag struct {
	saipb.UnimplementedLagServer
}

type macsec struct {
	saipb.UnimplementedMacsecServer
}

type mcastFdb struct {
	saipb.UnimplementedMcastFdbServer
}

type mirror struct {
	saipb.UnimplementedMirrorServer
}

type mpls struct {
	saipb.UnimplementedMplsServer
}

type myMac struct {
	saipb.UnimplementedMyMacServer
}

type nat struct {
	saipb.UnimplementedNatServer
}

type neighbor struct {
	saipb.UnimplementedNeighborServer
}

type nextHopGroup struct {
	saipb.UnimplementedNextHopGroupServer
}

type nextHop struct {
	saipb.UnimplementedNextHopServer
}

type policer struct {
	saipb.UnimplementedPolicerServer
}

type port struct {
	saipb.UnimplementedPortServer
}

type qosMap struct {
	saipb.UnimplementedQosMapServer
}

type queue struct {
	saipb.UnimplementedQueueServer
}

type route struct {
	saipb.UnimplementedRouteServer
}

type routerInterface struct {
	saipb.UnimplementedRouterInterfaceServer
}

type rpfGroup struct {
	saipb.UnimplementedRpfGroupServer
}

type samplePacket struct {
	saipb.UnimplementedSamplepacketServer
}

type schedulerGroup struct {
	saipb.UnimplementedSchedulerGroupServer
}

type scheduler struct {
	saipb.UnimplementedSchedulerServer
}

type srv6 struct {
	saipb.UnimplementedSrv6Server
}

type stp struct {
	saipb.UnimplementedStpServer
}

type saiSwitch struct {
	saipb.UnimplementedSwitchServer
}

type systemPort struct {
	saipb.UnimplementedSystemPortServer
}

type tam struct {
	saipb.UnimplementedTamServer
}

type tunnel struct {
	saipb.UnimplementedTunnelServer
}

type udf struct {
	saipb.UnimplementedUdfServer
}

type virtualRouter struct {
	saipb.UnimplementedVirtualRouterServer
}

type vlan struct {
	saipb.UnimplementedVlanServer
}

type wred struct {
	saipb.UnimplementedWredServer
}

type Server struct {
	acl             *acl
	bfd             *bfd
	buffer          *buffer
	bridge          *bridge
	counter         *counter
	debugCounter    *debugCounter
	fdb             *fdb
	hash            *hash
	hostif          *hostif
	ipmcGroup       *ipmcGroup
	ipmc            *ipmc
	ipsec           *ipsec
	isolationGroup  *isolationGroup
	l2mcGroup       *l2mcGroup
	l2mc            *l2mc
	lag             *lag
	macsec          *macsec
	mcastFdb        *mcastFdb
	mirror          *mirror
	mpls            *mpls
	myMac           *myMac
	nat             *nat
	neighbor        *neighbor
	nextHopGroup    *nextHopGroup
	nextHop         *nextHop
	policer         *policer
	port            *port
	qosMap          *qosMap
	queue           *queue
	route           *route
	routerInterface *routerInterface
	rpfGroup        *rpfGroup
	samplePacket    *samplePacket
	schedulerGroup  *schedulerGroup
	scheduler       *scheduler
	srv6            *srv6
	stp             *stp
	saiSwitch       *saiSwitch
	systemPort      *systemPort
	tam             *tam
	tunnel          *tunnel
	udf             *udf
	virtualRouter   *virtualRouter
	vlan            *vlan
	wred            *wred
}

func New(s *grpc.Server) *Server {
	srv := &Server{
		acl:             &acl{},
		bfd:             &bfd{},
		bridge:          &bridge{},
		buffer:          &buffer{},
		counter:         &counter{},
		debugCounter:    &debugCounter{},
		fdb:             &fdb{},
		hash:            &hash{},
		hostif:          &hostif{},
		ipmcGroup:       &ipmcGroup{},
		ipmc:            &ipmc{},
		ipsec:           &ipsec{},
		isolationGroup:  &isolationGroup{},
		l2mcGroup:       &l2mcGroup{},
		l2mc:            &l2mc{},
		lag:             &lag{},
		macsec:          &macsec{},
		mcastFdb:        &mcastFdb{},
		mirror:          &mirror{},
		mpls:            &mpls{},
		myMac:           &myMac{},
		nat:             &nat{},
		neighbor:        &neighbor{},
		nextHopGroup:    &nextHopGroup{},
		nextHop:         &nextHop{},
		policer:         &policer{},
		port:            &port{},
		qosMap:          &qosMap{},
		queue:           &queue{},
		route:           &route{},
		routerInterface: &routerInterface{},
		rpfGroup:        &rpfGroup{},
		samplePacket:    &samplePacket{},
		schedulerGroup:  &schedulerGroup{},
		scheduler:       &scheduler{},
		srv6:            &srv6{},
		stp:             &stp{},
		saiSwitch:       &saiSwitch{},
		systemPort:      &systemPort{},
		tam:             &tam{},
		tunnel:          &tunnel{},
		udf:             &udf{},
		virtualRouter:   &virtualRouter{},
		vlan:            &vlan{},
		wred:            &wred{},
	}
	saipb.RegisterAclServer(s, srv.acl)
	saipb.RegisterBfdServer(s, srv.bfd)
	saipb.RegisterBridgeServer(s, srv.bridge)
	saipb.RegisterCounterServer(s, srv.counter)
	saipb.RegisterDebugCounterServer(s, srv.debugCounter)
	saipb.RegisterFdbServer(s, srv.fdb)
	saipb.RegisterHashServer(s, srv.hash)
	saipb.RegisterHostifServer(s, srv.hostif)
	saipb.RegisterIpmcGroupServer(s, srv.ipmcGroup)
	saipb.RegisterIpmcServer(s, srv.ipmc)
	saipb.RegisterIpsecServer(s, srv.ipsec)
	saipb.RegisterIsolationGroupServer(s, srv.isolationGroup)
	saipb.RegisterL2McGroupServer(s, srv.l2mcGroup)
	saipb.RegisterL2McServer(s, srv.l2mc)
	saipb.RegisterLagServer(s, srv.lag)
	saipb.RegisterMacsecServer(s, srv.macsec)
	saipb.RegisterMcastFdbServer(s, srv.mcastFdb)
	saipb.RegisterMirrorServer(s, srv.mirror)
	saipb.RegisterMplsServer(s, srv.mpls)
	saipb.RegisterMyMacServer(s, srv.myMac)
	saipb.RegisterNatServer(s, srv.nat)
	saipb.RegisterNeighborServer(s, srv.neighbor)
	saipb.RegisterNextHopGroupServer(s, srv.nextHopGroup)
	saipb.RegisterNextHopServer(s, srv.nextHop)
	saipb.RegisterPolicerServer(s, srv.policer)
	saipb.RegisterPortServer(s, srv.port)
	saipb.RegisterQosMapServer(s, srv.qosMap)
	saipb.RegisterQueueServer(s, srv.queue)
	saipb.RegisterRouteServer(s, srv.route)
	saipb.RegisterRouterInterfaceServer(s, srv.routerInterface)
	saipb.RegisterRpfGroupServer(s, srv.rpfGroup)
	saipb.RegisterSamplepacketServer(s, srv.samplePacket)
	saipb.RegisterSchedulerGroupServer(s, srv.schedulerGroup)
	saipb.RegisterSchedulerServer(s, srv.scheduler)
	saipb.RegisterSrv6Server(s, srv.srv6)
	saipb.RegisterStpServer(s, srv.stp)
	saipb.RegisterSwitchServer(s, srv.saiSwitch)
	saipb.RegisterSystemPortServer(s, srv.systemPort)
	saipb.RegisterTamServer(s, srv.tam)
	saipb.RegisterTunnelServer(s, srv.tunnel)
	saipb.RegisterUdfServer(s, srv.udf)
	saipb.RegisterVirtualRouterServer(s, srv.virtualRouter)
	saipb.RegisterVlanServer(s, srv.vlan)
	saipb.RegisterWredServer(s, srv.wred)

	return srv
}
