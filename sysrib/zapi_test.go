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
	"fmt"
	"net"
	"testing"

	"github.com/coreswitch/netutil"
	"github.com/wenovus/gobgp/v3/pkg/zebra"
)

func SendMessage(conn net.Conn, msg *zebra.Message) {
	s, _ := msg.Serialize(zebra.MaxSoftware)
	conn.Write(s)
}

func Dial() (net.Conn, error) {
	// Connect to server
	conn, err := net.Dial("tcp", ":9000")
	if err != nil {
		return nil, err
	}
	return conn, err
}

// ZAPI TCP server start.
func ZAPIServerStart() *ZServer {
	// Start ZAPI server at port 9000 for VRF 0.
	return ZServerStart("tcp", ":9000", 0, nil)
}

func TestHello(t *testing.T) {
	fmt.Println("Hello")
	ZAPIServerStart()

	conn, err := Dial()
	if err != nil {
		t.Errorf("Dial failed %v\n", err)
		return
	}
	defer conn.Close()
	msg := zebra.Message{
		Header: zebra.Header{Command: zebra.Hello},
	}
	SendMessage(conn, &msg)
}

func TestRouteAdd(t *testing.T) {
	fmt.Println("Nexthop IPv4 Lookup")

	s := ZAPIServerStart()
	conn, err := Dial()
	if err != nil {
		t.Errorf("Dial failed %v\n", err)
		return
	}
	defer conn.Close()

	addr := netutil.ParseIPv4("10.0.0.0")
	body := &zebra.IPRouteBody{Prefix: zebra.Prefix{
		Family:    4,
		PrefixLen: 24,
		Prefix:    addr,
	}}

	msg := zebra.Message{
		Header: zebra.Header{Command: zebra.RouteAdd},
		Body:   body,
	}
	SendMessage(conn, &msg)
	s.Stop()
}
