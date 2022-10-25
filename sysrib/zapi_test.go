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

// Copyright 2016, 2017 zebra project.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sysrib

import (
	"context"
	"net"
	"syscall"
	"testing"
	"time"

	"github.com/coreswitch/netutil"
	"github.com/openconfig/lemming/gnmi"
	pb "github.com/openconfig/lemming/proto/sysrib"
	"github.com/wenovus/gobgp/v3/pkg/zebra"
	"google.golang.org/grpc"
)

func SendMessage(t *testing.T, conn net.Conn, msg *zebra.Message) error {
	s, err := msg.Serialize(zebra.MaxSoftware)
	if err != nil {
		return err
	}
	n, err := conn.Write(s)
	t.Logf("SendMessage: %d bytes written", n)
	return err
}

func Dial() (net.Conn, error) {
	conn, err := net.Dial("unix", "/tmp/zserv.api")
	if err != nil {
		return nil, err
	}
	return conn, err
}

func ZAPIServerStart(t *testing.T) *ZServer {
	t.Helper()
	dp := NewFakeDataplane()
	sysribServer, err := New(dp)
	if err != nil {
		t.Fatal(err)
	}

	s, err := ZServerStart("unix", "/tmp/zserv.api", 0, sysribServer)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func TestZServer(t *testing.T) {
	// Need to test these serially since they all use the same UDS address.
	testHello(t)
	testRouteAdd(t)
	testRouteRedistribution(t, true)
	testRouteRedistribution(t, false)
}

func testHello(t *testing.T) {
	s := ZAPIServerStart(t)
	defer s.Stop()

	conn, err := Dial()
	if err != nil {
		t.Errorf("Dial failed %v\n", err)
		return
	}
	defer conn.Close()

	serverVersion := zebra.MaxZapiVer
	serverSoftware := zebra.MaxSoftware
	msg := zebra.Message{
		Header: zebra.Header{
			Len:     zebra.HeaderSize(serverVersion),
			Marker:  zebra.HeaderMarker(serverVersion),
			Version: serverVersion,
			Command: zebra.Hello.ToEach(serverVersion, serverSoftware),
		},
		Body: &zebra.HelloBody{},
	}
	if err := SendMessage(t, conn, &msg); err != nil {
		t.Error(err)
	}
}

func testRouteAdd(t *testing.T) {
	s := ZAPIServerStart(t)
	defer s.Stop()

	conn, err := Dial()
	if err != nil {
		t.Errorf("Dial failed %v\n", err)
		return
	}
	defer conn.Close()

	addr := netutil.ParseIPv4("10.0.0.0")
	body := &zebra.IPRouteBody{Prefix: zebra.Prefix{
		Family:    syscall.AF_INET,
		PrefixLen: 24,
		Prefix:    addr,
	}}

	serverVersion := zebra.MaxZapiVer
	serverSoftware := zebra.MaxSoftware
	helloMsg := zebra.Message{
		Header: zebra.Header{
			Len:     zebra.HeaderSize(serverVersion),
			Marker:  zebra.HeaderMarker(serverVersion),
			Version: serverVersion,
			Command: zebra.Hello.ToEach(serverVersion, serverSoftware),
		},
		Body: &zebra.HelloBody{},
	}
	if err := SendMessage(t, conn, &helloMsg); err != nil {
		t.Error(err)
	}

	msg := zebra.Message{
		Header: zebra.Header{
			Len:     zebra.HeaderSize(serverVersion),
			Marker:  zebra.HeaderMarker(serverVersion),
			Version: serverVersion,
			Command: zebra.RouteAdd.ToEach(serverVersion, serverSoftware),
		},
		Body: body,
	}
	if err := SendMessage(t, conn, &msg); err != nil {
		t.Error(err)
	}
}

func testRouteRedistribution(t *testing.T, routeReadyBeforeDial bool) {
	dp := NewFakeDataplane()
	s, err := New(dp)
	if err != nil {
		t.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	gnmiServer, err := gnmi.New(grpcServer, "local")
	if err != nil {
		t.Fatal(err)
	}
	client := gnmiServer.LocalClient()
	if err := s.Start(client, "local", "unix:/tmp/zserv.api"); err != nil {
		t.Fatalf("cannot start sysrib server, %v", err)
	}
	defer s.Stop()

	configureInterface(t, AddIntfAction{
		name:    "eth0",
		ifindex: 0,
		enabled: true,
		prefix:  "192.168.1.1/24",
		niName:  "DEFAULT",
	}, client)

	routeReq := &pb.SetRouteRequest{
		AdminDistance: 10,
		Metric:        10,
		Prefix: &pb.Prefix{
			Family:     pb.Prefix_FAMILY_IPV4,
			Address:    "10.0.0.0",
			MaskLength: 8,
		},
		Nexthops: []*pb.Nexthop{{
			Type:    pb.Nexthop_TYPE_IPV4,
			Address: "192.168.1.42",
			Weight:  1,
		}},
	}

	if routeReadyBeforeDial {
		if _, err := s.SetRoute(context.Background(), routeReq); err != nil {
			t.Fatalf("Got unexpected error during call to SetRoute: %v", err)
		}
		for i := 0; i != maxGNMIWaitQuanta; i++ {
			if len(dp.GetRoutes()) != 0 {
				break
			}
			time.Sleep(100 * time.Millisecond)
		}
		if len(dp.GetRoutes()) == 0 {
			t.Fatalf("Route not resolved in time limit.")
		}
	}

	conn, err := Dial()
	if err != nil {
		t.Errorf("Dial failed %v\n", err)
		return
	}
	defer conn.Close()

	version := zebra.MaxZapiVer
	software := zebra.MaxSoftware
	logger := NewLogger()

	SendMessage(t, conn, &zebra.Message{
		Header: zebra.Header{
			Len:     zebra.HeaderSize(version),
			Marker:  zebra.HeaderMarker(version),
			Version: version,
			Command: zebra.Hello.ToEach(version, software),
		},
		Body: &zebra.HelloBody{},
	})

	_, err = zebra.ReceiveSingleMsg(logger, conn, version, software, "test-client")

	if !routeReadyBeforeDial {
		if _, err := s.SetRoute(context.Background(), routeReq); err != nil {
			t.Fatalf("Got unexpected error during call to SetRoute: %v", err)
		}
	}

	m, err := zebra.ReceiveSingleMsg(logger, conn, version, software, "test-client")
	if err != nil {
		t.Fatal(err)
	} else if m == nil {
		t.Fatal("got empty message")
	}
	if got, want := m.Header.Command.ToCommon(version, software), zebra.RedistributeRouteAdd; got != want {
		t.Errorf("Got message %s, want %s", got, want)
	}
}
