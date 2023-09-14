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

import "C"

import (
	"context"
	"flag"
	"fmt"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/openconfig/lemming/dataplane/internal/engine"
	"github.com/openconfig/lemming/dataplane/standalone/saiserver"
	"github.com/openconfig/lemming/dataplane/standalone/saiserver/attrmgr"

	log "github.com/golang/glog"

	dpb "github.com/openconfig/lemming/proto/dataplane"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

var fwdCtxID string

//export getForwardCtxID
func getForwardCtxID() *C.char {
	return C.CString(fwdCtxID)
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

//export initialize
func initialize(port int) {
	log.Info("lucius initialized")

	e, err := engine.New(context.Background())
	if err != nil {
		log.Fatalf("failed create engine: %v", err)
	}
	fwdCtxID = e.ID()

	lis, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	mgr := attrmgr.New()

	srv := grpc.NewServer(grpc.Creds(insecure.NewCredentials()),
		grpc.ChainUnaryInterceptor(logging.UnaryServerInterceptor(getLogger()), mgr.Interceptor),
		grpc.ChainStreamInterceptor(logging.StreamServerInterceptor(getLogger())))
	fwdpb.RegisterForwardingServer(srv, e)
	dpb.RegisterDataplaneServer(srv, e)
	saiserver.New(mgr, srv)

	go func() {
		if err := srv.Serve(lis); err != nil {
			log.Fatalf("failed to serve forwarding server: %v", err)
		}
	}()
}

func main() {
}

func init() {
	// TODO: Figure a better way to config this when used standalone.
	flag.Set("log_dir", "/var/log/syncd")
}
