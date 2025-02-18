// Copyright 2024 Google LLC
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

package telemetry

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	metricapi "cloud.google.com/go/monitoring/apiv3/v2"
	"cloud.google.com/go/monitoring/apiv3/v2/monitoringpb"
	traceapi "cloud.google.com/go/trace/apiv2"
	"cloud.google.com/go/trace/apiv2/tracepb"
	mexporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/metric"
	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"github.com/google/uuid"
	"github.com/googleapis/gax-go/v2/apierror"
	"go.opentelemetry.io/contrib/detectors/gcp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"golang.org/x/oauth2/google"
	"google.golang.org/genproto/googleapis/api/metric"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/openconfig/lemming/internal/cloudlog"
)

type opts struct {
	gcpProject string
	gcpTrace   bool
	gcpMeter   bool
	gcpLog     bool
}

type Option func(*opts)

// WithGCPProject sets to project for export OTEL traces and metrics as well as applications.
// If unset and telemetry epxort is enabled, the project is determined based the Application Default Credentials.
func WithGCPProject(proj string) Option {
	return func(o *opts) {
		o.gcpProject = proj
	}
}

func WithGCPMeterExport(enable bool) Option {
	return func(o *opts) {
		o.gcpMeter = enable
	}
}

// WithGCPTraceExport enables exporting OpenTelemetry traces and metrics to GCP.
func WithGCPTraceExport(enable bool) Option {
	return func(o *opts) {
		o.gcpTrace = enable
	}
}

// WithGCPLogExport enables streaming applications logs to GCP.
func WithGCPLogExport(enable bool) Option {
	return func(o *opts) {
		o.gcpLog = enable
	}
}

// Setup configures OpenTelemetry with option exporter.
func Setup(ctx context.Context, options ...Option) (func(context.Context) error, error) {
	o := &opts{}
	for _, opt := range options {
		opt(o)
	}

	if o.gcpProject == "" && (o.gcpLog || o.gcpTrace) {
		cred, err := google.FindDefaultCredentials(ctx, "https://www.googleapis.com/auth/cloud-platform")
		if err != nil {
			return nil, err
		}
		o.gcpProject = cred.ProjectID
	}

	var shutdownFuncs []func(context.Context) error

	var err error
	if o.gcpLog {
		shutdownFuncs = append(shutdownFuncs, cloudlog.SetGlobalLogger(ctx, o.gcpProject, "lucius", slog.LevelWarn))
	}

	host, err := os.Hostname()
	if err != nil {
		host = "unknown"
	}
	ns, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err != nil {
		ns = []byte{}
	}
	nodeID := os.Getenv("KNE_NODE_ID")
	if nodeID == "" {
		nodeID = uuid.New().String()
	}

	res, err := resource.New(ctx, resource.WithDetectors(gcp.NewDetector()),
		resource.WithHost(), resource.WithTelemetrySDK(), resource.WithAttributes(semconv.ServiceName(host), semconv.ServiceNamespace(string(ns)), semconv.ServiceInstanceID(nodeID)))
	if err != nil {
		return nil, err
	}

	shutdownFuncs = append(shutdownFuncs, setupTrace(ctx, res, o))
	shutdownFuncs = append(shutdownFuncs, setupMeter(ctx, res, o))

	shutdown := func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}

	// Set up logger provider.
	loggerProvider, err := newLoggerProvider(res)
	if err != nil {
		return nil, errors.Join(err, shutdown(ctx))
	}
	shutdownFuncs = append(shutdownFuncs, loggerProvider.Shutdown)
	global.SetLoggerProvider(loggerProvider)

	return shutdown, nil
}

func setupTrace(ctx context.Context, res *resource.Resource, o *opts) func(context.Context) error {
	prop := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
	otel.SetTextMapPropagator(prop)

	opts := []sdktrace.TracerProviderOption{sdktrace.WithSampler(sdktrace.AlwaysSample()), sdktrace.WithResource(res)}
	tracerProvider := sdktrace.NewTracerProvider(opts...)
	otel.SetTracerProvider(tracerProvider)
	if !o.gcpTrace {
		return tracerProvider.Shutdown
	}

	// Create a test span and try exporting it to verify authentication is successful.
	c, err := traceapi.NewClient(ctx)
	if err != nil {
		return tracerProvider.Shutdown
	}

	tracer := otel.Tracer("lucius")

	_, span := tracer.Start(context.Background(), "ping")
	span.AddEvent("test")
	span.End()

	req := &tracepb.BatchWriteSpansRequest{
		Name: fmt.Sprintf("projects/%s", o.gcpProject),
		Spans: []*tracepb.Span{{
			Name:        fmt.Sprintf("projects/%s/traces/%s/spans/%s", o.gcpProject, span.SpanContext().TraceID().String(), span.SpanContext().SpanID().String()),
			SpanId:      span.SpanContext().SpanID().String(),
			DisplayName: &tracepb.TruncatableString{Value: "test"},
			StartTime:   timestamppb.Now(),
			EndTime:     timestamppb.Now(),
		}},
	}

	exportOpts := []texporter.Option{texporter.WithProjectID(o.gcpProject)}

	err = c.BatchWriteSpans(ctx, req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to export test trace: %v\n", err)
		batchCtx := metadata.NewOutgoingContext(ctx, metadata.New(map[string]string{"x-goog-user-project": o.gcpProject}))
		err = c.BatchWriteSpans(batchCtx, req)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to export test trace with destination quota, disabling exporter: %v\n", err)
			return tracerProvider.Shutdown
		}
		exportOpts = append(exportOpts, texporter.WithDestinationProjectQuota())
	}

	traceExp, err := texporter.New(exportOpts...)
	if err != nil {
		return tracerProvider.Shutdown
	}
	opts = append(opts, sdktrace.WithSpanProcessor(sdktrace.NewBatchSpanProcessor(traceExp)))
	tracerProvider = sdktrace.NewTracerProvider(opts...)
	otel.SetTracerProvider(tracerProvider)
	return tracerProvider.Shutdown
}

func setupMeter(ctx context.Context, res *resource.Resource, o *opts) func(context.Context) error {
	mopts := []sdkmetric.Option{sdkmetric.WithResource(res)}
	meterProvider := sdkmetric.NewMeterProvider(mopts...)
	otel.SetMeterProvider(meterProvider)
	if !o.gcpMeter {
		return meterProvider.Shutdown
	}

	now := time.Now()
	mc, err := metricapi.NewMetricClient(ctx)
	if err != nil {
		return meterProvider.Shutdown
	}
	req := &monitoringpb.CreateTimeSeriesRequest{
		Name: fmt.Sprintf("projects/%s", o.gcpProject),
		TimeSeries: []*monitoringpb.TimeSeries{{
			Metric: &metric.Metric{
				Type: "custom.googleapis.com/lemming/test",
			},
			Points: []*monitoringpb.Point{{
				Interval: &monitoringpb.TimeInterval{
					StartTime: timestamppb.New(now),
					EndTime:   timestamppb.New(now),
				},
				Value: &monitoringpb.TypedValue{
					Value: &monitoringpb.TypedValue_Int64Value{
						Int64Value: 1,
					},
				},
			}},
		}},
	}

	exportOpts := []mexporter.Option{
		mexporter.WithProjectID(o.gcpProject),
		mexporter.WithFilteredResourceAttributes(attribute.NewAllowKeysFilter(semconv.ServiceNameKey, semconv.ServiceNamespaceKey, semconv.ServiceInstanceIDKey, semconv.CloudAccountIDKey, semconv.GCPGceInstanceNameKey)),
	}

	err = mc.CreateTimeSeries(ctx, req)
	var apiErr *apierror.APIError
	if errors.As(err, &apiErr) && apiErr.GRPCStatus().Code() == codes.InvalidArgument {
		fmt.Fprintf(os.Stderr, "test metric got invalid argument, likely ratelimit, continuing: %v\n", err)
	} else if err != nil {
		fmt.Fprintf(os.Stderr, "failed to export test metric: %v\n", err)
		batchCtx := metadata.NewOutgoingContext(ctx, metadata.New(map[string]string{"x-goog-user-project": o.gcpProject}))
		err = mc.CreateTimeSeries(batchCtx, req)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to export test metric with destination quota, disabling exporter: %v\n", err)
			return meterProvider.Shutdown
		}
		exportOpts = append(exportOpts, mexporter.WithDestinationProjectQuota())
	}

	metricExp, err := mexporter.New(exportOpts...)
	if err != nil {
		return meterProvider.Shutdown
	}

	mopts = append(mopts, sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExp)))
	meterProvider = sdkmetric.NewMeterProvider(mopts...)
	otel.SetMeterProvider(meterProvider)

	return meterProvider.Shutdown
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
