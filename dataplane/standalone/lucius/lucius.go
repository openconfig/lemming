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
	"errors"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net"

	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"go.opentelemetry.io/contrib/detectors/gcp"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"

	"github.com/openconfig/lemming/dataplane/dplaneopts"
	"github.com/openconfig/lemming/dataplane/saiserver"
	"github.com/openconfig/lemming/dataplane/saiserver/attrmgr"
	"github.com/openconfig/lemming/internal/cloudlog"

	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

var (
	port = flag.Int("port", 50000, "Port for api server")
	// All below flags are not used, kept for now to prevent breaking exsting users, will be removed in future release.
	_              = flag.String("config_file", "", "Path to config file (deprecated, no-op will always to true in future release)")
	_              = flag.String("port_map", "", "Map of modeled port names to Linux interface to  as comma seperated list (eg Ethernet8:eth1,Ethernet10,eth2) (deprecated, no-op will be removed in future release)")
	_              = flag.Bool("eth_dev_as_lane", true, "If true, when creating ports, use ethX and hardware lane X (deprecated, no-op will always to true in future release)")
	_              = flag.Bool("remote_cpu_port", true, "If true, send all packets from/to the CPU port over gRPC (deprecated, no-op will always to true in future release)")
	gcpTelemExport = flag.Bool("gcp_telem_export", false, "If true, export OTEL telemetry and logs to GCP")
	gcpProject     = flag.String("gcp_project", "", "GCP project to export to, by default it will use project where the GCE instance is running")
	hwProfile      = flag.String("hw_profile", "", "Path to hardware profile config file.")
)

func main() {
	flag.Parse()
	cancel, err := setupOTelSDK(context.Background())
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

	mgr := attrmgr.New()
	srv := grpc.NewServer(grpc.Creds(insecure.NewCredentials()),
		grpc.ChainUnaryInterceptor(logging.UnaryServerInterceptor(getLogger()), mgr.Interceptor, traceHandler),
		grpc.ChainStreamInterceptor(logging.StreamServerInterceptor(getLogger())),
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)

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

// setupOTelSDK bootstraps the OpenTelemetry pipeline.
// If it does not return an error, make sure to call shutdown for proper cleanup.
func setupOTelSDK(ctx context.Context) (func(context.Context) error, error) {
	var shutdownFuncs []func(context.Context) error

	var exporter sdktrace.SpanExporter

	var err error
	if *gcpTelemExport {
		exporter, err = texporter.New(texporter.WithProjectID(*gcpProject), texporter.WithDestinationProjectQuota())
		if err != nil {
			return nil, err
		}
		cloudlog.SetGlobalLogger(ctx, *gcpProject, "lucius")
	}

	shutdown := func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}

	res, err := resource.New(ctx, resource.WithDetectors(gcp.NewDetector()), resource.WithHost(), resource.WithTelemetrySDK())
	if err != nil {
		return nil, err
	}

	// Set up propagator.
	prop := newPropagator()
	otel.SetTextMapPropagator(prop)

	bsp := sdktrace.NewBatchSpanProcessor(exporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
	)

	shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)
	otel.SetTracerProvider(tracerProvider)

	shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)
	otel.SetMeterProvider(meterProvider)

	// Set up logger provider.
	loggerProvider, err := newLoggerProvider(res)
	if err != nil {
		return nil, errors.Join(err, shutdown(ctx))
	}
	shutdownFuncs = append(shutdownFuncs, loggerProvider.Shutdown)
	global.SetLoggerProvider(loggerProvider)

	return shutdown, nil
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func newLoggerProvider(res *resource.Resource) (*sdklog.LoggerProvider, error) {
	logExporter, err := stdoutlog.New()
	if err != nil {
		return nil, err
	}

	loggerProvider := sdklog.NewLoggerProvider(
		sdklog.WithProcessor(sdklog.NewBatchProcessor(logExporter)),
		sdklog.WithResource(res),
	)
	return loggerProvider, nil
}
