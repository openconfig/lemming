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
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/wenovus/gobgp/v3/pkg/zebra"

	log "github.com/golang/glog"
	bgplog "github.com/wenovus/gobgp/v3/pkg/log"
)

// Part of this file was adapted from
// https://github.com/coreswitch/zebra/blob/master/rib/zapi.go

// TODO(wenbli): Consider putting ZAPI logic into a different, internal package
// once the API boundary between ZAPI and the RIB manager becomes more stable.
// An interface is a way to break the circular dependency.

// TODO(wenbli): Consider unifying logging calls. Although klog is the
// currently-desired logger, "Topic=" is actually a good way of filtering
// through logs rather than simply using verbosity.

// ZServer is a ZAPI server.
type ZServer struct {
	socketType string
	path       string
	vrfID      uint32
	sysrib     *Server
	lis        net.Listener

	// ClientMutex protects the ZAPI client map.
	ClientMutex sync.RWMutex
	// ClientMap stores all connected ZAPI clients.
	ClientMap map[net.Conn]*Client
}

// StartZServer starts a ZAPI server on the given connection type and path,
//
// e.g.
// - "unix", "/var/run/zapi.serv"
//
// It also resquires the sysrib server in order to send and receive
// redistributed routes.
//
// TODO: vrfID is not well-integrated with the sysrib.
func StartZServer(ctx context.Context, address string, vrfID uint32, sysrib *Server) (*ZServer, error) {
	l := strings.SplitN(address, ":", 2)
	if len(l) != 2 {
		return nil, fmt.Errorf("unsupported ZAPI url, has to be \"protocol:address\", got: %s", address)
	}
	socketType, path := l[0], l[1]

	if err := os.RemoveAll(path); err != nil {
		return nil, err
	}

	var err error
	var lis net.Listener

	switch socketType {
	case "unix", "unix-writable":
		lis, err = net.Listen("unix", path)
		if err != nil {
			return nil, fmt.Errorf("cannot start ZAPI server: %v", err)
		}
		if socketType == "unix-writable" {
			if err = os.Chmod(path, 0777); err != nil {
				return nil, fmt.Errorf("cannot start ZAPI server: %v", err)
			}
		}
	default:
		return nil, fmt.Errorf("zebra server socket type must be unix or unix-writable")
	}

	zServer := &ZServer{
		socketType: l[0],
		path:       l[1],
		vrfID:      vrfID,
		sysrib:     sysrib,
		lis:        lis,
		ClientMap:  map[net.Conn]*Client{},
	}

	go func() {
		log.Infof("ZAPI Server started at %s", path)
		for {
			// Listen for an incoming connection.
			conn, err := lis.Accept()
			if err != nil {
				log.Infof("Stopping ZAPI server: %v", err)
				return
			}

			// Register client.
			client := zServer.ClientRegister(conn)
			client.zServer = zServer

			// Handle connections in a new go routine.
			go client.HandleRequest(ctx, conn, vrfID)
		}
	}()

	return zServer, nil
}

// ClientRegister creates a new ZAPI client connection.
func (s *ZServer) ClientRegister(conn net.Conn) *Client {
	s.ClientMutex.Lock()
	defer s.ClientMutex.Unlock()

	log.Info("zapi:ClientRegister", conn)
	client := &Client{conn: conn}
	s.ClientMap[conn] = client
	return client
}

// ClientUnregister deletes a ZAPI client connection.
func (s *ZServer) ClientUnregister(conn net.Conn) {
	s.ClientMutex.Lock()
	defer s.ClientMutex.Unlock()

	log.Info("zapi:ClientUnregister", conn)
	delete(s.ClientMap, conn)
}

// Stop stops the ZAPI server.
func (s *ZServer) Stop() {
	if s != nil {
		if s.lis != nil {
			s.lis.Close()
		}
		os.Remove(s.path)
	}
}

// Client is a ZAPI client.
type Client struct {
	conn    net.Conn
	zServer *ZServer
}

// RedistributeResolvedRoutes sends RedistributeRouteAdd messages to the client
// connection for all currently-resolved routes.
func (c *Client) RedistributeResolvedRoutes(conn net.Conn) {
	resolvableRoutes := c.zServer.sysrib.ResolvedRoutes()
	programmedRoutes := c.zServer.sysrib.ProgrammedRoutes()
	topicLogger.Info(fmt.Sprintf("Sending %d resolved routes to client", len(programmedRoutes)),
		bgplog.Fields{
			"Topic": "Sysrib",
		})
	for routeKey, rr := range resolvableRoutes {
		route, ok := programmedRoutes[routeKey]
		if !ok {
			continue
		}
		zrouteBody, err := convertToZAPIRoute(routeKey, rr, route)
		if err != nil {
			topicLogger.Warn(fmt.Sprintf("failed to convert resolved route to zebra BGP route: %v", err),
				bgplog.Fields{
					"Topic": "Sysrib",
				})
		}
		topicLogger.Info("Sending resolved route",
			bgplog.Fields{
				"Topic":   "Sysrib",
				"Message": zrouteBody,
			})
		if zrouteBody != nil {
			if err := serverSendMessage(conn, zebra.RedistributeRouteAdd, zrouteBody); err != nil {
				topicLogger.Error(fmt.Sprintf("Cannot send RedistributeRouteAdd message: %v", err),
					bgplog.Fields{
						"Topic": "Sysrib",
					})
			}
		}
	}
}

// HandleRequest handles an incoming ZAPI client connection.
func (c *Client) HandleRequest(ctx context.Context, conn net.Conn, vrfID uint32) {
	version := zebra.MaxZapiVer
	software := zebra.MaxSoftware
	defer func() {
		err := conn.Close()
		if err != nil {
			topicLogger.Error("error while closing connection to client, stopping client handling thread.",
				bgplog.Fields{
					"Topic": "Sysrib",
					"Error": err,
				})
		}
		log.Infof("[zapi] disconnected, vrf %d, version %v", vrfID, version)
		c.zServer.ClientUnregister(conn)
	}()

	for {
		m, err := zebra.ReceiveSingleMsg(topicLogger, conn, version, software, "Sysrib")
		switch {
		case err == io.EOF:
			log.Warningf("ZAPI server stopping after receiving EOF")
			return
		case err != nil:
			log.Errorf("ZAPI server stopping, HandleRequest error: %v", err)
			return
		case m == nil:
			continue
		}

		command := m.Header.Command.ToCommon(version, software)
		switch command {
		case zebra.Hello:
			topicLogger.Info("Received Zebra Hello from client:",
				bgplog.Fields{
					"Topic":   "Sysrib",
					"Message": m,
				})
			// TODO(wenbli): A capabilities message should be sent instead.
			// This doesn't matter right now because it appears no
			// client (isisd nor GoBGP) actually looks at this message.
			if err := serverSendMessage(conn, zebra.Hello, &zebra.HelloBody{}); err != nil {
				topicLogger.Error(fmt.Sprintf("Cannot send hello message: %v", err),
					bgplog.Fields{
						"Topic":   "Sysrib",
						"Message": m,
					})
				return
			}
			c.RedistributeResolvedRoutes(conn)
		case zebra.RouteAdd:
			topicLogger.Info("Received Zebra RouteAdd from client:",
				bgplog.Fields{
					"Topic":   "Sysrib",
					"Message": m,
				})
			if err := c.zServer.sysrib.setZebraRoute(ctx, vrfIDToNiName(vrfID), m.Body.(*zebra.IPRouteBody)); err != nil {
				topicLogger.Warn(fmt.Sprintf("Could not add route to sysrib: %v", err),
					bgplog.Fields{
						"Topic":   "Sysrib",
						"Message": m,
					})
			}
		case zebra.RouteDelete:
			// TODO(wenbli): Implement RouteDelete.
			topicLogger.Warn("Received Zebra RouteDelete from client which is not handled:",
				bgplog.Fields{
					"Topic":   "Sysrib",
					"Message": m,
				})
		default:
			topicLogger.Warn(fmt.Sprintf("Received unhandled Zebra message %v from client:", command),
				bgplog.Fields{
					"Topic":   "Sysrib",
					"Message": m,
				})
		}
	}
}

// serverSendMessage sends a message and returns a bool indicating whether a
// fatal error was encountered and logged.
func serverSendMessage(conn net.Conn, command zebra.APIType, body zebra.Body) error {
	serverVersion := zebra.MaxZapiVer
	serverSoftware := zebra.MaxSoftware
	m := &zebra.Message{
		Header: zebra.Header{
			Len:     zebra.HeaderSize(serverVersion),
			Marker:  zebra.HeaderMarker(serverVersion),
			Version: serverVersion,
			VrfID:   zebra.DefaultVrf,
			Command: command.ToEach(serverVersion, serverSoftware),
		},
		Body: body,
	}
	topicLogger.Info(fmt.Sprintf("sending message: %v", command),
		bgplog.Fields{
			"Topic":   "Sysrib",
			"Message": m,
		})
	b, err := m.Serialize(serverSoftware)
	if err != nil {
		return err
	}

	if _, err := conn.Write(b); err != nil {
		return err
	}
	return nil
}
