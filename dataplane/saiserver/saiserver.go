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

	"github.com/openconfig/lemming/dataplane/dplaneopts"
	"github.com/openconfig/lemming/dataplane/forwarding"
	"github.com/openconfig/lemming/dataplane/saiserver/attrmgr"

	"google.golang.org/grpc"

	log "github.com/golang/glog"

	saipb "github.com/openconfig/lemming/dataplane/proto/sai"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

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

type ipmcGroup struct {
	saipb.UnimplementedIpmcGroupServer
}

type ipmc struct {
	saipb.UnimplementedIpmcServer
}

type ipsec struct {
	saipb.UnimplementedIpsecServer
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

type nat struct {
	saipb.UnimplementedNatServer
}

type samplePacket struct {
	saipb.UnimplementedSamplepacketServer
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

type wred struct {
	saipb.UnimplementedWredServer
}

type forwardingContext struct {
	*forwarding.Server
	id string
}

func (fc *forwardingContext) ID() string {
	return fc.id
}

type Server struct {
	saipb.UnimplementedEntrypointServer
	*forwardingContext
	mgr          *attrmgr.AttrMgr
	initialized  bool
	bfd          *bfd
	buffer       *buffer
	counter      *counter
	debugCounter *debugCounter
	dtel         *dtel
	fdb          *fdb
	ipmcGroup    *ipmcGroup
	ipmc         *ipmc
	ipsec        *ipsec
	macsec       *macsec
	mcastFdb     *mcastFdb
	mirror       *mirror
	mpls         *mpls
	nat          *nat
	samplePacket *samplePacket
	srv6         *srv6
	saiSwitch    *saiSwitch
	systemPort   *systemPort
	tam          *tam
	wred         *wred
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

func (s *Server) Initialize(ctx context.Context, _ *saipb.InitializeRequest) (*saipb.InitializeResponse, error) {
	if s.initialized {
		log.Info("dataplane already intialized, reseting")
		s.mgr.Reset()
		s.saiSwitch.Reset()
		if err := s.Reset(ctx); err != nil {
			return nil, err
		}
	}
	s.initialized = true

	return &saipb.InitializeResponse{}, nil
}

func (s *Server) Reset(ctx context.Context) error {
	_, err := s.forwardingContext.ContextDelete(ctx, &fwdpb.ContextDeleteRequest{
		ContextId: &fwdpb.ContextId{Id: s.forwardingContext.id},
	})
	if err != nil {
		return err
	}

	_, err = s.forwardingContext.ContextCreate(ctx, &fwdpb.ContextCreateRequest{
		ContextId: &fwdpb.ContextId{Id: s.forwardingContext.id},
	})
	if err != nil {
		return err
	}

	return nil
}

func New(ctx context.Context, mgr *attrmgr.AttrMgr, s *grpc.Server, opts *dplaneopts.Options) (*Server, error) {
	fwdCtx := &forwardingContext{Server: forwarding.New("engine"), id: "lucius"}
	_, err := fwdCtx.ContextCreate(ctx, &fwdpb.ContextCreateRequest{
		ContextId: &fwdpb.ContextId{Id: fwdCtx.id},
	})
	if err != nil {
		return nil, err
	}
	sw, err := newSwitch(mgr, fwdCtx, s, opts)
	if err != nil {
		return nil, err
	}

	srv := &Server{
		mgr:               mgr,
		forwardingContext: fwdCtx,
		bfd:               &bfd{},
		buffer:            &buffer{},
		counter:           &counter{},
		debugCounter:      &debugCounter{},
		dtel:              &dtel{},
		fdb:               &fdb{},
		ipmcGroup:         &ipmcGroup{},
		ipmc:              &ipmc{},
		ipsec:             &ipsec{},
		macsec:            &macsec{},
		mcastFdb:          &mcastFdb{},
		mirror:            &mirror{},
		mpls:              &mpls{},
		nat:               &nat{},
		samplePacket:      &samplePacket{},
		srv6:              &srv6{},
		saiSwitch:         sw,
		systemPort:        &systemPort{},
		tam:               &tam{},
		wred:              &wred{},
	}
	fwdpb.RegisterForwardingServer(s, fwdCtx)
	fwdpb.RegisterInfoServer(s, fwdCtx)
	saipb.RegisterEntrypointServer(s, srv)
	saipb.RegisterBfdServer(s, srv.bfd)
	saipb.RegisterCounterServer(s, srv.counter)
	saipb.RegisterDebugCounterServer(s, srv.debugCounter)
	saipb.RegisterDtelServer(s, srv.dtel)
	saipb.RegisterFdbServer(s, srv.fdb)
	saipb.RegisterIpmcGroupServer(s, srv.ipmcGroup)
	saipb.RegisterIpmcServer(s, srv.ipmc)
	saipb.RegisterIpsecServer(s, srv.ipsec)
	saipb.RegisterMacsecServer(s, srv.macsec)
	saipb.RegisterMcastFdbServer(s, srv.mcastFdb)
	saipb.RegisterMirrorServer(s, srv.mirror)
	saipb.RegisterMplsServer(s, srv.mpls)
	saipb.RegisterNatServer(s, srv.nat)
	saipb.RegisterSamplepacketServer(s, srv.samplePacket)
	saipb.RegisterSrv6Server(s, srv.srv6)
	saipb.RegisterSystemPortServer(s, srv.systemPort)
	saipb.RegisterTamServer(s, srv.tam)
	saipb.RegisterWredServer(s, srv.wred)

	return srv, nil
}
