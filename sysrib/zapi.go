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
	"fmt"
	"net"
	"os"
	"sync"

	"github.com/golang/glog"
	"github.com/openconfig/gribigo/afthelper"
	"github.com/wenovus/gobgp/v3/pkg/log"
	"github.com/wenovus/gobgp/v3/pkg/zebra"
)

// TODO(wenbli): Consider putting this into a different, internal package once
// the API boundary between ZAPI and the RIB manager becomes more stable. An
// interface is a way to break the circular dependency.

// TODO(wenbli): Unify logging calls. "Topic=" is actually a good way of
// filtering through logs rather than simply using verbosity.

// Client is a ZAPI client.
type Client struct {
	conn      net.Conn
	version   uint8
	allVrf    bool
	vrfId     uint32
	routeType zebra.RouteType
	zserver   *ZServer
}

var (
	// ClientMutex protects the ZAPI client map.
	ClientMutex sync.RWMutex
	// ClientMap stores all connected ZAPI clients.
	ClientMap = map[net.Conn]*Client{}
)

// ClientRegister creates a new ZAPI client connection.
func ClientRegister(conn net.Conn) *Client {
	ClientMutex.Lock()
	defer ClientMutex.Unlock()

	glog.Info("zapi:ClientRegister", conn)
	client := &Client{conn: conn}
	ClientMap[conn] = client
	return client
}

// ClientUnregister deletes a ZAPI client connection.
func ClientUnregister(conn net.Conn) {
	ClientMutex.Lock()
	defer ClientMutex.Unlock()

	glog.Info("zapi:ClientUnregister", conn)
	delete(ClientMap, conn)
}

// convertZebraRoute converts a zebra route to a Sysrib route.
func convertZebraRoute(niName string, zroute *zebra.IPRouteBody) *Route {
	var nexthops []*afthelper.NextHopSummary
	for _, znh := range zroute.Nexthops {
		nexthops = append(nexthops, &afthelper.NextHopSummary{
			Weight:          1,
			Address:         znh.Gate.String(),
			NetworkInstance: niName,
		})
	}
	return &Route{
		Prefix: fmt.Sprintf("%s/%d", zroute.Prefix.Prefix.String(), zroute.Prefix.PrefixLen),
		// NextHops is the set of IP nexthops that the route uses if
		// it is not a connected route.
		NextHops: nexthops,
		RoutePref: RoutePreference{
			AdminDistance: zroute.Distance,
			Metric:        zroute.Metric,
		},
	}
}

// HandleRequest handles an incoming ZAPI client connection.
func (c *Client) HandleRequest(conn net.Conn, vrfID uint32) {
	version := zebra.MaxZapiVer
	software := zebra.MaxSoftware
	logger := NewLogger()
	defer func() {
		err := conn.Close()
		if err != nil {
			logger.Error("error while closing connection to client, stopping client handling thread.",
				log.Fields{
					"Topic": "Sysrib",
					"Error": err,
				})
		}
		glog.Infof("[zapi] disconnected, vrf %d, version %v", vrfID, version)
		ClientUnregister(conn)
	}()

	for {
		m, err := zebra.ReceiveSingleMsg(logger, conn, version, software, "Sysrib")
		if err != nil {
			glog.Errorf("ZAPI server stopping, HandleRequest error: %v", err)
			return
		} else if m == nil {
			continue
		}

		command := m.Header.Command.ToCommon(version, software)
		switch command {
		case zebra.Hello:
			logger.Info("Received Zebra Hello from client:",
				log.Fields{
					"Topic":   "Sysrib",
					"Message": m,
				})
			// HACK: A capabilities message should be sent instead.
			// This doesn't matter right now because it appears no
			// client (isisd nor GoBGP) actually looks at this message.
			if serverSendMessage(logger, conn, zebra.Hello, &zebra.HelloBody{}) {
				logger.Error("Cannot send hello message",
					log.Fields{
						"Topic":   "Sysrib",
						"Message": m,
					})
				return
			}
			resolvedRoutes := c.zserver.sysrib.ResolvedRoutes()
			logger.Info(fmt.Sprintf("Sending %d resolved routes to client", len(resolvedRoutes)),
				log.Fields{
					"Topic":   "Sysrib",
					"Message": m,
				})
			for routeKey, route := range resolvedRoutes {
				zrouteBody, _, err := convertToZAPIRoute(routeKey, route)
				if err != nil {
					logger.Warn(fmt.Sprintf("failed to convert resolved route to zebra BGP route: %v", err),
						log.Fields{
							"Topic":   "Sysrib",
							"Message": m,
						})
				}
				logger.Info("Sending resolved route",
					log.Fields{
						"Topic":   "Sysrib",
						"Message": zrouteBody,
					})
				if zrouteBody != nil {
					if serverSendMessage(logger, conn, zebra.RedistributeRouteAdd, zrouteBody) {
						logger.Error("Cannot send RedistributeRouteAdd message",
							log.Fields{
								"Topic":   "Sysrib",
								"Message": m,
							})
					}
				}
			}
		case zebra.RouteAdd:
			logger.Info("Received Zebra RouteAdd from client:",
				log.Fields{
					"Topic":   "Sysrib",
					"Message": m,
				})
			zroute := m.Body.(*zebra.IPRouteBody)
			niName := vrfIDToNiName(vrfID)
			route := convertZebraRoute(niName, zroute)
			if c.zserver.sysrib != nil {
				if err := c.zserver.sysrib.setRoute(niName, route); err != nil {
					logger.Warn(fmt.Sprintf("Could not add route to sysrib: %v", route),
						log.Fields{
							"Topic":   "Sysrib",
							"Message": m,
						})
				}
			} else {
				logger.Warn("ZServer does not have reference to sysrib, cannot add route to RIB manager",
					log.Fields{
						"Topic":   "Sysrib",
						"Message": m,
					})
			}
		case zebra.RouteDelete:
			logger.Info("Received Zebra RouteDelete from client:",
				log.Fields{
					"Topic":   "Sysrib",
					"Message": m,
				})
		default:
			logger.Warn(fmt.Sprintf("Received unhandled Zebra message %v from client:", command),
				log.Fields{
					"Topic":   "Sysrib",
					"Message": m,
				})
		}
	}
}

// serverSendMessage sends a message and returns a bool indicating whether a
// fatal error was encountered and logged.
func serverSendMessage(logger log.Logger, conn net.Conn, command zebra.APIType, body zebra.Body) bool {
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
	logger.Info(fmt.Sprintf("sending message: %v", command),
		log.Fields{
			"Topic":   "Sysrib",
			"Message": m,
		})
	b, err := m.Serialize(serverSoftware)
	if err != nil {
		logger.Warn(fmt.Sprintf("failed to serialize: %v", m),
			log.Fields{
				"Topic": "Sysrib",
				"Error": err,
			})
		return false
	}

	_, err = conn.Write(b)
	if err != nil {
		logger.Error("failed to write to client, closing connection to client and stopping client handling thread.",
			log.Fields{
				"Topic": "Sysrib",
				"Error": err,
			})
		return true
	}
	return false
}

// ZServer is a ZAPI server.
type ZServer struct {
	path   string
	vrfID  uint32
	sysrib *Server
	lis    net.Listener
}

// ZServerStart starts a ZAPI server on the given connection type and path,
//
// e.g.
// - "unix", "/var/run/zapi.serv"
func ZServerStart(socketType string, path string, vrfID uint32, sysrib *Server) (*ZServer, error) {
	if err := os.RemoveAll(path); err != nil {
		return nil, err
	}

	var lis net.Listener
	var err error

	switch socketType {
	case "unix", "unix-writable":
		os.Remove(path)
		lis, err = net.Listen("unix", path)
		if err != nil {
			return nil, fmt.Errorf("ZAPI server: %v", err)
		}
		if socketType == "unix-writable" {
			if err = os.Chmod(path, 0777); err != nil {
				return nil, fmt.Errorf("ZAPI server: %v", err)
			}
		}
	default:
		return nil, fmt.Errorf("ZServer socket type is not unix.")
	}

	server := &ZServer{
		path:   path,
		lis:    lis,
		vrfID:  vrfID,
		sysrib: sysrib,
	}

	go func() {
		glog.Infof("ZAPI Server started at %s", path)
		for {
			// Listen for an incoming connection.
			conn, err := lis.Accept()
			if err != nil {
				glog.Infof("Stopping ZAPI server: %v", err)
				return
			}

			// Register client.
			client := ClientRegister(conn)
			client.zserver = server

			// Handle connections in a new go routine.
			go client.HandleRequest(conn, vrfID)
		}
	}()

	return server, nil
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
