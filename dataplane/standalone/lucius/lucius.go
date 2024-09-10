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

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/openconfig/lemming/dataplane/dplaneopts"
	"github.com/openconfig/lemming/dataplane/saiserver"
	"github.com/openconfig/lemming/dataplane/saiserver/attrmgr"

	log "github.com/golang/glog"

	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

var (
	port = flag.Int("port", 50000, "Port for api server")
	// All below flags are not used, kept for now to prevent breaking exsting users, will be removed in future release.
	_         = flag.String("config_file", "", "Path to config file (deprecated, no-op will always to true in future release)")
	_         = flag.String("port_map", "", "Map of modeled port names to Linux interface to  as comma seperated list (eg Ethernet8:eth1,Ethernet10,eth2) (deprecated, no-op will be removed in future release)")
	_         = flag.Bool("eth_dev_as_lane", true, "If true, when creating ports, use ethX and hardware lane X (deprecated, no-op will always to true in future release)")
	_         = flag.Bool("remote_cpu_port", true, "If true, send all packets from/to the CPU port over gRPC (deprecated, no-op will always to true in future release)")
	hwProfile = flag.String("hw_profile", "", "Path to hardware profile config file.")
)

func main() {
	flag.Parse()
	start(*port)
}

func getLogger() logging.Logger {
	return logging.LoggerFunc(func(_ context.Context, level logging.Level, msg string, fields ...any) {
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

	mgr := attrmgr.New()
	srv := grpc.NewServer(grpc.Creds(insecure.NewCredentials()),
		grpc.ChainUnaryInterceptor(logging.UnaryServerInterceptor(getLogger()), mgr.Interceptor),
		grpc.ChainStreamInterceptor(logging.StreamServerInterceptor(getLogger())))

	reflection.Register(srv)

	opts := dplaneopts.ResolveOpts(
		dplaneopts.WithHostifNetDevPortType(fwdpb.PortType_PORT_TYPE_KERNEL),
		dplaneopts.WithHardwareProfile(*hwProfile),
	)

	if _, err := saiserver.New(context.Background(), mgr, srv, opts); err != nil {
		log.Fatalf("failed to create server: %v", err)
	}

	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve forwarding server: %v", err)
	}
}
