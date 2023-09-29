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
	"context"
	"fmt"

	"github.com/openconfig/lemming/dataplane/standalone/saiserver/attrmgr"

	"google.golang.org/grpc"

	log "github.com/golang/glog"

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

type counter struct {
	saipb.UnimplementedCounterServer
}

type debugCounter struct {
	saipb.UnimplementedDebugCounterServer
}

type dtel struct {
	saipb.UnimplementedDtelServer
}

type fdb struct {
	saipb.UnimplementedFdbServer
}

type hash struct {
	saipb.UnimplementedHashServer
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

type policer struct {
	saipb.UnimplementedPolicerServer
}

type qosMap struct {
	saipb.UnimplementedQosMapServer
}

type queue struct {
	saipb.UnimplementedQueueServer
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

type wred struct {
	saipb.UnimplementedWredServer
}

type Server struct {
	saipb.UnimplementedEntrypointServer
	mgr            *attrmgr.AttrMgr
	acl            *acl
	bfd            *bfd
	buffer         *buffer
	counter        *counter
	debugCounter   *debugCounter
	dtel           *dtel
	fdb            *fdb
	ipmcGroup      *ipmcGroup
	ipmc           *ipmc
	ipsec          *ipsec
	isolationGroup *isolationGroup
	l2mcGroup      *l2mcGroup
	l2mc           *l2mc
	lag            *lag
	macsec         *macsec
	mcastFdb       *mcastFdb
	mirror         *mirror
	mpls           *mpls
	myMac          *myMac
	nat            *nat
	policer        *policer
	qosMap         *qosMap
	queue          *queue
	rpfGroup       *rpfGroup
	samplePacket   *samplePacket
	schedulerGroup *schedulerGroup
	scheduler      *scheduler
	srv6           *srv6
	saiSwitch      *saiSwitch
	systemPort     *systemPort
	tam            *tam
	tunnel         *tunnel
	udf            *udf
	wred           *wred
}

func (s *Server) ObjectTypeQuery(_ context.Context, req *saipb.ObjectTypeQueryRequest) (*saipb.ObjectTypeQueryResponse, error) {
	val := s.mgr.GetType(fmt.Sprint(req.GetObject()))
	if val == saipb.ObjectType_OBJECT_TYPE_NULL {
		log.Warningf("unknown object id %v, type %v", req.Object, val)
	}
	return &saipb.ObjectTypeQueryResponse{
		Type: val,
	}, nil
}

func New(mgr *attrmgr.AttrMgr, engine switchDataplaneAPI, s *grpc.Server) *Server {
	srv := &Server{
		mgr:            mgr,
		acl:            &acl{},
		bfd:            &bfd{},
		buffer:         &buffer{},
		counter:        &counter{},
		debugCounter:   &debugCounter{},
		dtel:           &dtel{},
		fdb:            &fdb{},
		ipmcGroup:      &ipmcGroup{},
		ipmc:           &ipmc{},
		ipsec:          &ipsec{},
		isolationGroup: &isolationGroup{},
		l2mcGroup:      &l2mcGroup{},
		l2mc:           &l2mc{},
		lag:            &lag{},
		macsec:         &macsec{},
		mcastFdb:       &mcastFdb{},
		mirror:         &mirror{},
		mpls:           &mpls{},
		myMac:          &myMac{},
		nat:            &nat{},
		policer:        &policer{},
		qosMap:         &qosMap{},
		queue:          &queue{},
		rpfGroup:       &rpfGroup{},
		samplePacket:   &samplePacket{},
		schedulerGroup: &schedulerGroup{},
		scheduler:      &scheduler{},
		srv6:           &srv6{},
		saiSwitch:      newSwitch(mgr, engine, s),
		systemPort:     &systemPort{},
		tam:            &tam{},
		tunnel:         &tunnel{},
		udf:            &udf{},
		wred:           &wred{},
	}
	saipb.RegisterEntrypointServer(s, srv)
	saipb.RegisterAclServer(s, srv.acl)
	saipb.RegisterBfdServer(s, srv.bfd)
	saipb.RegisterCounterServer(s, srv.counter)
	saipb.RegisterDebugCounterServer(s, srv.debugCounter)
	saipb.RegisterDtelServer(s, srv.dtel)
	saipb.RegisterFdbServer(s, srv.fdb)
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
	saipb.RegisterPolicerServer(s, srv.policer)
	saipb.RegisterQosMapServer(s, srv.qosMap)
	saipb.RegisterQueueServer(s, srv.queue)
	saipb.RegisterRpfGroupServer(s, srv.rpfGroup)
	saipb.RegisterSamplepacketServer(s, srv.samplePacket)
	saipb.RegisterSchedulerGroupServer(s, srv.schedulerGroup)
	saipb.RegisterSchedulerServer(s, srv.scheduler)
	saipb.RegisterSrv6Server(s, srv.srv6)
	saipb.RegisterSystemPortServer(s, srv.systemPort)
	saipb.RegisterTamServer(s, srv.tam)
	saipb.RegisterTunnelServer(s, srv.tunnel)
	saipb.RegisterUdfServer(s, srv.udf)
	saipb.RegisterWredServer(s, srv.wred)

	return srv
}
