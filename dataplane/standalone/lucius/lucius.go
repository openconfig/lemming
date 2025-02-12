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
	"log"
	"log/slog"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"

	"github.com/openconfig/lemming/dataplane/dplaneopts"
	"github.com/openconfig/lemming/dataplane/saiserver"
	"github.com/openconfig/lemming/dataplane/saiserver/attrmgr"
	"github.com/openconfig/lemming/internal/telemetry"

	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

var (
	port = flag.Int("port", 50000, "Port for api server")
	// All below flags are not used, kept for now to prevent breaking exsting users, will be removed in future release.
	_              = flag.String("config_file", "", "Path to config file (deprecated, no-op will always to true in future release)")
	_              = flag.String("port_map", "", "Map of modeled port names to Linux interface to  as comma seperated list (eg Ethernet8:eth1,Ethernet10,eth2) (deprecated, no-op will be removed in future release)")
	_              = flag.Bool("eth_dev_as_lane", true, "If true, when creating ports, use ethX and hardware lane X (deprecated, no-op will always to true in future release)")
	_              = flag.Bool("remote_cpu_port", true, "If true, send all packets from/to the CPU port over gRPC (deprecated, no-op will always to true in future release)")
	instGrpc       = flag.Bool("instrument_grpc", false, "If true, adds intrumentation for all gRPCS servers.")
	gcpTraceExport = flag.Bool("gcp_trace_export", false, "If true, export OTEL traces to GCP")
	gcpMeterExport = flag.Bool("gcp_meter_export", false, "If true, export OTEL meters to GCP")
	gcpLogExport   = flag.Bool("gcp_log_export", false, "If true, export application logs to GCP")
	gcpProject     = flag.String("gcp_project", "", "GCP project to export to, by default it will use project where the GCE instance is running")
	hwProfile      = flag.String("hw_profile", "", "Path to hardware profile config file.")
)

func main() {
	flag.Parse()

	cancel, err := telemetry.Setup(context.Background(), telemetry.WithGCPProject(*gcpProject), telemetry.WithGCPLogExport(*gcpLogExport), telemetry.WithGCPTraceExport(*gcpTraceExport), telemetry.WithGCPMeterExport(*gcpMeterExport))
	if err != nil {
		log.Fatal(err)
	}
	defer cancel(context.Background())
	start(*port)
}

func getLogger() logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		slog.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}

var tracer = otel.Tracer("")

func traceHandler(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	ctx, span := tracer.Start(ctx, info.FullMethod)
	defer span.End()

	resp, err := handler(ctx, req)
	grpc.SetTrailer(ctx, metadata.Pairs("traceparent", span.SpanContext().TraceID().String()))

	return resp, err
}

func start(port int) {
	slog.Info("lucius initialized")

	lis, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	logOpts := []logging.Option{logging.WithLogOnEvents(logging.FinishCall), logging.WithLevels(func(code codes.Code) logging.Level {
		switch code {
		case codes.OK:
			return logging.LevelDebug
		}
		return logging.DefaultServerCodeToLevel(code)
	})}

	mgr := attrmgr.New()
	srvOpts := []grpc.ServerOption{
		grpc.Creds(insecure.NewCredentials()),
		grpc.ChainUnaryInterceptor(logging.UnaryServerInterceptor(getLogger(), logOpts...), mgr.Interceptor, traceHandler),
		grpc.ChainStreamInterceptor(logging.StreamServerInterceptor(getLogger(), logOpts...)),
	}
	if *instGrpc {
		srvOpts = append(srvOpts, grpc.StatsHandler(otelgrpc.NewServerHandler()))
	}

	srv := grpc.NewServer(srvOpts...)

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
