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

	"github.com/openconfig/ygnmi/ygnmi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/local"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/proto"

	"github.com/openconfig/lemming/dataplane/dplaneopts"
	"github.com/openconfig/lemming/dataplane/saiserver"
	"github.com/openconfig/lemming/dataplane/saiserver/attrmgr"
	"github.com/openconfig/lemming/dataplane/standalone/pkthandler/pktiohandler"
	"github.com/openconfig/lemming/gnmi/oc"
	"github.com/openconfig/lemming/gnmi/reconciler"

	gpb "github.com/openconfig/gnmi/proto/gnmi"

	pktiopb "github.com/openconfig/lemming/dataplane/proto/packetio"
	saipb "github.com/openconfig/lemming/dataplane/proto/sai"
)

// Dataplane is an implementation of Dataplane HAL API.
type Dataplane struct {
	saiserv     *saiserver.Server
	srv         *grpc.Server
	lis         net.Listener
	reconcilers []reconciler.Reconciler
	opt         *dplaneopts.Options
	cancelFn    func()
}

// New create a new dataplane instance.
func New(ctx context.Context, opts ...dplaneopts.Option) (*Dataplane, error) {
	data := &Dataplane{
		opt: dplaneopts.ResolveOpts(append(opts, dplaneopts.WithEthDevAsLane(true), dplaneopts.WithRemoteCPUPort(true))...),
	}

	lis, err := (&net.ListenConfig{}).Listen(ctx, "tcp", data.opt.AddrPort)
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %w", err)
	}

	mgr := attrmgr.New()
	srv := grpc.NewServer(grpc.Creds(local.NewCredentials()), grpc.ChainUnaryInterceptor(mgr.Interceptor))
	reflection.Register(srv)

	saiserv, err := saiserver.New(ctx, mgr, srv, data.opt)
	if err != nil {
		return nil, fmt.Errorf("failed to create: %w", err)
	}
	data.saiserv = saiserv
	data.lis = lis

	go srv.Serve(data.lis)

	return data, nil
}

// ID returns the ID of the dataplane reconciler.
func (d *Dataplane) ID() string {
	return "dataplane"
}

// Start starts the HAL gRPC server and packet forwarding engine.
func (d *Dataplane) Start(ctx context.Context, c gpb.GNMIClient, target string) error {
	if d.srv != nil {
		return fmt.Errorf("dataplane already started")
	}

	conn, err := d.Conn()
	if err != nil {
		return err
	}
	sw := saipb.NewSwitchClient(conn)
	hostif := saipb.NewHostifClient(conn)
	swResp, err := sw.CreateSwitch(ctx, &saipb.CreateSwitchRequest{})
	if err != nil {
		return err
	}
	swAttrs, err := sw.GetSwitchAttribute(ctx, &saipb.GetSwitchAttributeRequest{
		Oid: swResp.Oid,
		AttrType: []saipb.SwitchAttr{
			saipb.SwitchAttr_SWITCH_ATTR_CPU_PORT,
		},
	})
	if err != nil {
		return err
	}

	// Allow all traffic to L3 processing.
	hash := saipb.NewHashClient(conn)
	hashResp, err := hash.CreateHash(context.Background(), &saipb.CreateHashRequest{
		Switch: swResp.Oid,
		NativeHashFieldList: []saipb.NativeHashField{
			saipb.NativeHashField_NATIVE_HASH_FIELD_SRC_IP,
			saipb.NativeHashField_NATIVE_HASH_FIELD_DST_IP,
			saipb.NativeHashField_NATIVE_HASH_FIELD_L4_SRC_PORT,
			saipb.NativeHashField_NATIVE_HASH_FIELD_L4_DST_PORT,
		},
	})
	if err != nil {
		return err
	}
	_, err = sw.SetSwitchAttribute(ctx, &saipb.SetSwitchAttributeRequest{
		Oid:          swResp.Oid,
		EcmpHashIpv4: proto.Uint64(hashResp.Oid),
		EcmpHashIpv6: proto.Uint64(hashResp.Oid),
	})
	if err != nil {
		return err
	}

	// Allow all traffic to L3 processing.
	mmc := saipb.NewMyMacClient(conn)
	_, err = mmc.CreateMyMac(context.Background(), &saipb.CreateMyMacRequest{
		Switch:         swResp.Oid,
		Priority:       proto.Uint32(1),
		MacAddress:     []byte{0, 0, 0, 0, 0, 0},
		MacAddressMask: []byte{0, 0, 0, 0, 0, 0},
	})
	if err != nil {
		return err
	}

	_, err = hostif.CreateHostifTrap(ctx, &saipb.CreateHostifTrapRequest{
		Switch:       swResp.Oid,
		TrapType:     saipb.HostifTrapType_HOSTIF_TRAP_TYPE_ARP_REQUEST.Enum(),
		PacketAction: saipb.PacketAction_PACKET_ACTION_TRAP.Enum(),
	})
	if err != nil {
		return err
	}
	_, err = hostif.CreateHostifTrap(ctx, &saipb.CreateHostifTrapRequest{
		Switch:       swResp.Oid,
		TrapType:     saipb.HostifTrapType_HOSTIF_TRAP_TYPE_IPV6_NEIGHBOR_DISCOVERY.Enum(),
		PacketAction: saipb.PacketAction_PACKET_ACTION_TRAP.Enum(),
	})
	if err != nil {
		return err
	}

	h, err := pktiohandler.New("")
	if err != nil {
		return err
	}

	pktio := pktiopb.NewPacketIOClient(conn)

	ctx, d.cancelFn = context.WithCancel(ctx)
	portCtl, err := pktio.HostPortControl(ctx)
	if err != nil {
		return err
	}
	packet, err := pktio.CPUPacketStream(ctx)
	if err != nil {
		return err
	}

	go h.ManagePorts(portCtl)
	go h.StreamPackets(packet)

	if d.opt.Reconcilation {
		d.reconcilers = append(d.reconcilers, getReconcilers(conn, swResp.Oid, *swAttrs.GetAttr().CpuPort, "lucius")...)

		for _, rec := range d.reconcilers {
			if err := rec.Start(ctx, c, target); err != nil {
				return fmt.Errorf("failed to stop handler %q: %v", rec.ID(), err)
			}
		}
	}

	return nil
}

// FwdClient gets a gRPC client to the packet forwarding engine.
func (d *Dataplane) Conn() (grpc.ClientConnInterface, error) {
	conn, err := grpc.Dial(d.lis.Addr().String(), grpc.WithTransportCredentials(local.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to dial server: %w", err)
	}
	return conn, nil
}

func (d *Dataplane) SaiServer() *saiserver.Server {
	return d.saiserv
}

// Stop gracefully stops the server.
func (d *Dataplane) Stop(ctx context.Context) error {
	d.cancelFn()
	for _, rec := range d.reconcilers {
		if err := rec.Stop(ctx); err != nil {
			return fmt.Errorf("failed to stop handler %q: %v", rec.ID(), err)
		}
	}
	d.srv.GracefulStop()
	return nil
}

// Validate is a noop to implement to the reconciler interface.
func (d *Dataplane) Validate(*oc.Root) error {
	return nil
}

// ValidationPaths is a noop to implement to the reconciler interface.
func (d *Dataplane) ValidationPaths() []ygnmi.PathStruct {
	return nil
}
