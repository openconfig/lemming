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

// Package local_test is an integration test of lemming's BGP functionality
// using lemmings instantiated on localhost.
package local_test

import (
	"context"
	"fmt"
	"net/netip"
	"testing"

	"github.com/openconfig/lemming"
	"github.com/openconfig/lemming/gnmi"
	"github.com/openconfig/lemming/gnmi/fakedevice"
	"github.com/openconfig/lemming/gnmi/gnmiclient"
	"github.com/openconfig/lemming/gnmi/oc"
	"github.com/openconfig/lemming/gnmi/oc/ocpath"
	"github.com/openconfig/ygnmi/ygnmi"
	"github.com/openconfig/ygot/ygot"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/credentials/local"

	gpb "github.com/openconfig/gnmi/proto/gnmi"
)

// TODO: Consolidate test helper code with integration and other unit tests.

type DeviceSpec struct {
	ID        uint
	gnmiPort  uint
	gribiPort uint
	bgpPort   uint16
	AS        uint32
	RouterID  string
}

var (
	dut1spec = DeviceSpec{
		ID:        1,
		gnmiPort:  7339,
		gribiPort: 7340,
		bgpPort:   1111,
		AS:        64500,
		RouterID:  "127.0.0.1",
	}
	dut2spec = DeviceSpec{
		ID:        2,
		gnmiPort:  8339,
		gribiPort: 8340,
		bgpPort:   1112,
		AS:        64501,
		RouterID:  "127.0.0.2",
	}
	dut3spec = DeviceSpec{
		ID:        3,
		gnmiPort:  9339,
		gribiPort: 9340,
		bgpPort:   1113,
		AS:        64502,
		RouterID:  "127.0.0.3",
	}
)

func ygnmiClient(t *testing.T, target string, port int) *ygnmi.Client {
	t.Helper()
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", port), grpc.WithTransportCredentials(local.NewCredentials()))
	if err != nil {
		t.Fatalf("cannot dial gNMI server, %v", err)
	}
	yc, err := ygnmi.NewClient(gpb.NewGNMIClient(conn), ygnmi.WithTarget(target))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	return yc
}

// AddIntfAction represents an action to add a connected interface so that the
// device is able to resolve routes.
type AddIntfAction struct {
	name    string
	ifindex int32
	enabled bool
	prefix  string
	niName  string
}

// configureInterface configures an interface on fake the lemming gNMI server.
func configureInterface(t *testing.T, intf *AddIntfAction, gnmiSv *gnmi.Server) {
	t.Helper()

	c, err := ygnmi.NewClient(gnmiSv.LocalClient(), ygnmi.WithTarget("local"))
	if err != nil {
		t.Fatalf("cannot create ygnmi client: %v", err)
	}

	ocintf := &oc.Interface{}
	ocintf.Name = ygot.String(intf.name)
	ocintf.Enabled = ygot.Bool(intf.enabled)
	ocintf.Ifindex = ygot.Uint32(uint32(intf.ifindex))
	prefix, err := netip.ParsePrefix(intf.prefix)
	if err != nil {
		t.Fatalf("Invalid prefix: %q", intf.prefix)
	}
	switch {
	case prefix.Addr().Is4():
		ocaddr := ocintf.GetOrCreateSubinterface(0).GetOrCreateIpv4().GetOrCreateAddress(prefix.Addr().String())
		ocaddr.PrefixLength = ygot.Uint8(uint8(prefix.Bits()))
	case prefix.Addr().Is6():
		ocaddr := ocintf.GetOrCreateSubinterface(0).GetOrCreateIpv6().GetOrCreateAddress(prefix.Addr().String())
		ocaddr.PrefixLength = ygot.Uint8(uint8(prefix.Bits()))
	default:
		t.Fatalf("Prefix is neither IPv4 nor IPv6: %q", intf.prefix)
	}

	if _, err := gnmiclient.Replace(context.Background(), c, ocpath.Root().Interface(intf.name).State(), ocintf); err != nil {
		t.Fatalf("Cannot configure interface: %v", err)
	}
}

func bgpWithNbr(as uint32, routerID string, nbr *oc.NetworkInstance_Protocol_Bgp_Neighbor) *oc.NetworkInstance_Protocol_Bgp {
	bgp := &oc.NetworkInstance_Protocol_Bgp{}
	bgp.GetOrCreateGlobal().As = ygot.Uint32(as)
	if routerID != "" {
		bgp.Global.RouterId = ygot.String(routerID)
	}
	bgp.AppendNeighbor(nbr)
	return bgp
}

func newLemming(t *testing.T, dev DeviceSpec, connectedIntfs []*AddIntfAction) (*ygnmi.Client, func()) {
	opts := []lemming.Option{lemming.WithTransportCreds(insecure.NewCredentials()), lemming.WithGRIBIAddr(fmt.Sprintf(":%d", dev.gribiPort)), lemming.WithGNMIAddr(fmt.Sprintf(":%d", dev.gnmiPort)), lemming.WithBGPPort(dev.bgpPort)}

	target := fmt.Sprintf("dut%d", dev.ID)

	l, err := lemming.New(target, fmt.Sprintf("unix:/tmp/zserv-test%d.api", dev.ID), opts...)
	if err != nil {
		t.Fatal(err)
	}

	for _, intf := range connectedIntfs {
		configureInterface(t, intf, l.GNMI())
	}

	return ygnmiClient(t, target, int(dev.gnmiPort)), func() { l.Stop() }
}

func establishSessionPair(t *testing.T, dut1, dut2 *ygnmi.Client, spec1, spec2 DeviceSpec) {
	bgpPath := ocpath.Root().NetworkInstance(fakedevice.DefaultNetworkInstance).Protocol(oc.PolicyTypes_INSTALL_PROTOCOL_TYPE_BGP, fakedevice.BGPRoutingProtocol).Bgp()

	dutConf := bgpWithNbr(spec1.AS, spec1.RouterID, &oc.NetworkInstance_Protocol_Bgp_Neighbor{
		PeerAs:          ygot.Uint32(spec2.AS),
		NeighborAddress: ygot.String(spec2.RouterID),
		NeighborPort:    ygot.Uint16(spec2.bgpPort),
	})
	dut2Conf := bgpWithNbr(spec2.AS, spec2.RouterID, &oc.NetworkInstance_Protocol_Bgp_Neighbor{
		PeerAs:          ygot.Uint32(spec1.AS),
		NeighborAddress: ygot.String(spec1.RouterID),
		NeighborPort:    ygot.Uint16(spec1.bgpPort),
	})
	Update(t, dut1, bgpPath.Config(), dutConf)
	Update(t, dut2, bgpPath.Config(), dut2Conf)

	nbrPath := bgpPath.Neighbor(spec2.RouterID)
	Await(t, dut1, nbrPath.SessionState().State(), oc.Bgp_Neighbor_SessionState_ESTABLISHED)
}

func TestSessionEstablish(t *testing.T) {
	dut1, stop1 := newLemming(t, dut1spec, nil)
	defer stop1()
	dut2, stop2 := newLemming(t, dut2spec, nil)
	defer stop2()

	establishSessionPair(t, dut1, dut2, dut1spec, dut2spec)
}

func TestEstablishDifferentIP(t *testing.T) {
	dut2, stop2 := newLemming(t, dut2spec, nil)
	defer stop2()
	dut3, stop3 := newLemming(t, dut3spec, nil)
	defer stop3()

	establishSessionPair(t, dut2, dut3, dut2spec, dut3spec)
}
