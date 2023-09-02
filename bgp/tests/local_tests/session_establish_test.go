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
	"net"
	"net/netip"
	"sync/atomic"
	"testing"

	"github.com/openconfig/lemming"
	"github.com/openconfig/lemming/bgp"
	"github.com/openconfig/lemming/gnmi"
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

type Device struct {
	yc       *ygnmi.Client
	ID       uint
	AS       uint32
	bgpPort  uint16
	RouterID string
}

func nextLocalHostAddr() string {
	// Start at 127.0.0.2
	return netip.AddrFrom4([4]byte{127, 0, 0, byte(nextLocalHostAddrIndex.Add(1) + 1)}).String()
}

var (
	nextLocalHostAddrIndex = &atomic.Uint32{}
)

func ygnmiClient(t *testing.T, target, dialTarget string) *ygnmi.Client {
	t.Helper()
	conn, err := grpc.Dial(dialTarget, grpc.WithTransportCredentials(local.NewCredentials()))
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

func newLemming(t *testing.T, id uint, as uint32, connectedIntfs []*AddIntfAction) (*Device, func()) {
	routerID := nextLocalHostAddr()
	gnmiTarget := net.JoinHostPort(routerID, "7339")
	gribiTarget := net.JoinHostPort(routerID, "7340")
	opts := []lemming.Option{lemming.WithTransportCreds(insecure.NewCredentials()), lemming.WithGRIBIAddr(gribiTarget), lemming.WithGNMIAddr(gnmiTarget), lemming.WithBGPPort(1111)}

	target := fmt.Sprintf("dut%d", id)

	l, err := lemming.New(target, fmt.Sprintf("unix:/tmp/zserv-test%d.api", id), opts...)
	if err != nil {
		t.Fatal(err)
	}

	for _, intf := range connectedIntfs {
		configureInterface(t, intf, l.GNMI())
	}

	return &Device{
		yc:       ygnmiClient(t, target, gnmiTarget),
		ID:       id,
		AS:       as,
		bgpPort:  1111,
		RouterID: routerID,
	}, func() { l.Stop() }
}

type DevicePair struct {
	first  *Device
	second *Device
}

func establishSessionPairs(t *testing.T, dutPairs ...DevicePair) {
	t.Helper()
	for _, pair := range dutPairs {
		dut1, dut2 := pair.first, pair.second
		dutConf := bgpWithNbr(dut1.AS, dut1.RouterID, &oc.NetworkInstance_Protocol_Bgp_Neighbor{
			PeerAs:          ygot.Uint32(dut2.AS),
			NeighborAddress: ygot.String(dut2.RouterID),
			NeighborPort:    ygot.Uint16(dut2.bgpPort),
			Transport: &oc.NetworkInstance_Protocol_Bgp_Neighbor_Transport{
				LocalAddress: ygot.String(dut1.RouterID),
			},
		})
		dut2Conf := bgpWithNbr(dut2.AS, dut2.RouterID, &oc.NetworkInstance_Protocol_Bgp_Neighbor{
			PeerAs:          ygot.Uint32(dut1.AS),
			NeighborAddress: ygot.String(dut1.RouterID),
			NeighborPort:    ygot.Uint16(dut1.bgpPort),
			Transport: &oc.NetworkInstance_Protocol_Bgp_Neighbor_Transport{
				LocalAddress: ygot.String(dut2.RouterID),
			},
		})
		Update(t, dut1, bgp.BGPPath.Config(), dutConf)
		Update(t, dut2, bgp.BGPPath.Config(), dut2Conf)
	}

	for _, pair := range dutPairs {
		dut1, dut2 := pair.first, pair.second
		awaitSessionEstablished(t, dut1, dut2)
	}
}

func awaitSessionEstablished(t *testing.T, dut1, dut2 *Device) {
	t.Helper()
	Await(t, dut1, bgp.BGPPath.Neighbor(dut2.RouterID).SessionState().State(), oc.Bgp_Neighbor_SessionState_ESTABLISHED)
	Await(t, dut2, bgp.BGPPath.Neighbor(dut1.RouterID).SessionState().State(), oc.Bgp_Neighbor_SessionState_ESTABLISHED)
}

func TestSessionEstablish(t *testing.T) {
	dut1, stop1 := newLemming(t, 1, 64500, nil)
	defer stop1()
	dut2, stop2 := newLemming(t, 2, 64501, nil)
	defer stop2()

	establishSessionPairs(t, DevicePair{first: dut1, second: dut2})
}
