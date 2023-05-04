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
	"fmt"
	"testing"

	"github.com/openconfig/lemming"
	"github.com/openconfig/lemming/gnmi/fakedevice"
	"github.com/openconfig/lemming/gnmi/oc"
	"github.com/openconfig/lemming/gnmi/oc/ocpath"
	"github.com/openconfig/ygnmi/ygnmi"
	"github.com/openconfig/ygot/ygot"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/credentials/local"

	gpb "github.com/openconfig/gnmi/proto/gnmi"
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

func bgpWithNbr(as uint32, routerID string, nbr *oc.NetworkInstance_Protocol_Bgp_Neighbor) *oc.NetworkInstance_Protocol_Bgp {
	bgp := &oc.NetworkInstance_Protocol_Bgp{}
	bgp.GetOrCreateGlobal().As = ygot.Uint32(as)
	if routerID != "" {
		bgp.Global.RouterId = ygot.String(routerID)
	}
	bgp.AppendNeighbor(nbr)
	return bgp
}

func TestSessionEstablish(t *testing.T) {
	const (
		dut1BgpPort uint16 = 1111
		dut2BgpPort uint16 = 1112
		dut1AS      uint32 = 64500
		dut2AS      uint32 = 64501
		addr1       string = "127.0.0.1"
		addr2       string = "127.0.0.2"
	)

	l1, err := lemming.New("dut1", "unix:/tmp/zserv-test1.api", lemming.WithTransportCreds(insecure.NewCredentials()), lemming.WithGRIBIAddr(":7340"), lemming.WithGNMIAddr(":7339"), lemming.WithBGPPort(dut1BgpPort))
	if err != nil {
		t.Fatal(err)
	}
	defer l1.Stop()

	l2, err := lemming.New("dut2", "unix:/tmp/zserv-test2.api", lemming.WithTransportCreds(insecure.NewCredentials()), lemming.WithGRIBIAddr(":8340"), lemming.WithGNMIAddr(":8339"), lemming.WithBGPPort(dut2BgpPort))
	if err != nil {
		t.Fatal(err)
	}
	defer l2.Stop()

	dut1 := ygnmiClient(t, "dut1", 7339)
	dut2 := ygnmiClient(t, "dut2", 8339)

	bgpPath := ocpath.Root().NetworkInstance(fakedevice.DefaultNetworkInstance).Protocol(oc.PolicyTypes_INSTALL_PROTOCOL_TYPE_BGP, "BGP").Bgp()

	// Remove any existing BGP config
	Delete(t, dut1, bgpPath.Config())
	Delete(t, dut2, bgpPath.Config())

	// Start a new session
	dutConf := bgpWithNbr(dut1AS, addr1, &oc.NetworkInstance_Protocol_Bgp_Neighbor{
		PeerAs:          ygot.Uint32(dut2AS),
		NeighborAddress: ygot.String(addr2),
		NeighborPort:    ygot.Uint16(dut2BgpPort),
	})
	dut2Conf := bgpWithNbr(dut2AS, addr2, &oc.NetworkInstance_Protocol_Bgp_Neighbor{
		PeerAs:          ygot.Uint32(dut1AS),
		NeighborAddress: ygot.String(addr1),
		NeighborPort:    ygot.Uint16(dut1BgpPort),
	})
	Replace(t, dut1, bgpPath.Config(), dutConf)
	Replace(t, dut2, bgpPath.Config(), dut2Conf)

	nbrPath := bgpPath.Neighbor(addr2)
	Await(t, dut1, nbrPath.SessionState().State(), oc.Bgp_Neighbor_SessionState_ESTABLISHED)
}
