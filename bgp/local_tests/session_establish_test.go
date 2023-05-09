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

type BGPDevice struct {
	AS       uint32
	RouterID string
}

var (
	dut1BGP = BGPDevice{
		AS:       64500,
		RouterID: "127.0.0.1",
	}
	dut2BGP = BGPDevice{
		AS:       64501,
		RouterID: "127.0.0.2",
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

func establishSession(t *testing.T, dut1BgpPort, dut2BgpPort uint16, lemming1Intfs, lemming2Intfs []*AddIntfAction) (*ygnmi.Client, *ygnmi.Client, func()) {
	lemming1Opts := []lemming.Option{lemming.WithTransportCreds(insecure.NewCredentials()), lemming.WithGRIBIAddr(":7340"), lemming.WithGNMIAddr(":7339"), lemming.WithBGPPort(dut1BgpPort)}

	l1, err := lemming.New("dut1", "unix:/tmp/zserv-test1.api", lemming1Opts...)
	if err != nil {
		t.Fatal(err)
	}

	lemming2Opts := []lemming.Option{lemming.WithTransportCreds(insecure.NewCredentials()), lemming.WithGRIBIAddr(":8340"), lemming.WithGNMIAddr(":8339"), lemming.WithBGPPort(dut2BgpPort)}

	l2, err := lemming.New("dut2", "unix:/tmp/zserv-test2.api", lemming2Opts...)
	if err != nil {
		t.Fatal(err)
	}

	for _, intf := range lemming1Intfs {
		configureInterface(t, intf, l1.GNMI())
	}
	for _, intf := range lemming2Intfs {
		configureInterface(t, intf, l2.GNMI())
	}

	dut1 := ygnmiClient(t, "dut1", 7339)
	dut2 := ygnmiClient(t, "dut2", 8339)

	bgpPath := ocpath.Root().NetworkInstance(fakedevice.DefaultNetworkInstance).Protocol(oc.PolicyTypes_INSTALL_PROTOCOL_TYPE_BGP, fakedevice.BGPRoutingProtocol).Bgp()

	// Remove any existing BGP config
	Delete(t, dut1, bgpPath.Config())
	Delete(t, dut2, bgpPath.Config())

	// Start a new session
	dutConf := bgpWithNbr(dut1BGP.AS, dut1BGP.RouterID, &oc.NetworkInstance_Protocol_Bgp_Neighbor{
		PeerAs:          ygot.Uint32(dut2BGP.AS),
		NeighborAddress: ygot.String(dut2BGP.RouterID),
		NeighborPort:    ygot.Uint16(dut2BgpPort),
	})
	dut2Conf := bgpWithNbr(dut2BGP.AS, dut2BGP.RouterID, &oc.NetworkInstance_Protocol_Bgp_Neighbor{
		PeerAs:          ygot.Uint32(dut1BGP.AS),
		NeighborAddress: ygot.String(dut1BGP.RouterID),
		NeighborPort:    ygot.Uint16(dut1BgpPort),
	})
	Replace(t, dut1, bgpPath.Config(), dutConf)
	Replace(t, dut2, bgpPath.Config(), dut2Conf)

	nbrPath := bgpPath.Neighbor(dut2BGP.RouterID)
	Await(t, dut1, nbrPath.SessionState().State(), oc.Bgp_Neighbor_SessionState_ESTABLISHED)

	return dut1, dut2, func() { l1.Stop(); l2.Stop() }
}

func TestSessionEstablish(t *testing.T) {
	_, _, stop := establishSession(t, 1111, 1112, nil, nil)
	stop()
}
