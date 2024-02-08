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

package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"strings"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/openconfig/lemming/dataplane/dplaneopts"
	"github.com/openconfig/lemming/dataplane/saiserver"
	"github.com/openconfig/lemming/dataplane/saiserver/attrmgr"

	log "github.com/golang/glog"

	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

var (
	port          = flag.Int("port", 50000, "Port for api server")
	configFile    = flag.String("config_file", "", "Path to config file")
	portMapString = flag.String("port_map", "", "Map of modeled port names to Linux interface to  as comma seperated list (eg Ethernet8:eth1,Ethernet10,eth2)")
)

func main() {
	flag.Parse()
	start(*port)
}

func getLogger() logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, level logging.Level, msg string, fields ...any) {
		switch level {
		case logging.LevelDebug:
			log.V(1).Info(msg, fields)
		case logging.LevelInfo:
			log.Info(msg, fields)
		case logging.LevelWarn:
			log.Warning(msg, fields)
		case logging.LevelError:
			log.Error(msg, fields)
		}
	})
}

func start(port int) {
	log.Info("lucius initialized")

	lis, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	portMap := map[string]string{}
	for _, mapping := range strings.Split(*portMapString, ",") {
		ports := strings.Split(mapping, ":")
		if len(ports) != 2 {
			log.Fatal("invalid port map format")
		}
		portMap[ports[0]] = ports[1]
	}

	mgr := attrmgr.New()

	srv := grpc.NewServer(grpc.Creds(insecure.NewCredentials()),
		grpc.ChainUnaryInterceptor(logging.UnaryServerInterceptor(getLogger()), mgr.Interceptor),
		grpc.ChainStreamInterceptor(logging.StreamServerInterceptor(getLogger())))

	opts := dplaneopts.ResolveOpts(
		dplaneopts.WithHostifNetDevPortType(fwdpb.PortType_PORT_TYPE_KERNEL),
		dplaneopts.WithPortConfigFile(*configFile),
		dplaneopts.WithPortMap(portMap),
	)

	if _, err := saiserver.New(context.Background(), mgr, srv, opts); err != nil {
		log.Fatalf("failed to create server: %v", err)
	}

	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve forwarding server: %v", err)
	}
}
