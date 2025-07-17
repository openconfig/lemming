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

package main

import (
	"context"
	"flag"
	"os"
	"os/signal"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/openconfig/lemming"
	"github.com/openconfig/lemming/dataplane/dplaneopts"
	"github.com/openconfig/lemming/internal/telemetry"

	log "github.com/golang/glog"
)

var (
	gnmiAddr       = pflag.String("gnmi", ":9339", "gNMI listen address")
	gribiAddr      = pflag.String("gribi", ":9340", "gRIBI listen address")
	p4rtAddr       = pflag.String("p4rt_addr", ":9559", "p4rt listen address")
	bgpPort        = pflag.Uint("bgp_port", 179, "BGP listening port")
	target         = pflag.String("target", "fakedut", "name of the fake target")
	tlsKeyFile     = pflag.String("tls_key_file", "", "Controls whether to enable TLS for gNXI services. If unspecified, insecure credentials are used.")
	tlsCertFile    = pflag.String("tls_cert_file", "", "Controls whether to enable TLS for gNXI services. If unspecified, insecure credentials are used.")
	zapiAddr       = pflag.String("zapi_addr", "unix:/var/run/zserv.api", "Custom ZAPI address: use unix:/tmp/zserv.api for a temp.")
	dplane         = pflag.Bool("enable_dataplane", false, "Controls whether to enable dataplane")
	gcpTraceExport = pflag.Bool("gcp_trace_export", false, "If true, export OTEL traces to GCP")
	gcpMeterExport = pflag.Bool("gcp_meter_export", false, "If true, export OTEL meters to GCP")
	gcpLogExport   = pflag.Bool("gcp_log_export", false, "If true, export application logs to GCP")
	gcpProject     = pflag.String("gcp_project", "", "GCP project to export to, by default it will use project where the GCE instance is running")
	faultAddr      = pflag.String("fault_addr", ":9399", "fault server listen address")
	faultEnable    = pflag.Bool("enable_fault", true, "Enable fault service")
	configFile     = pflag.String("config_file", "", "Path to configuration file or vendor preset (e.g., 'arista'). If not specified, checks LEMMING_CONFIG_FILE, then uses defaults.")
)

func main() {
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	cancel, err := telemetry.Setup(context.Background(), telemetry.WithGCPProject(*gcpProject), telemetry.WithGCPLogExport(*gcpLogExport), telemetry.WithGCPTraceExport(*gcpTraceExport), telemetry.WithGCPMeterExport(*gcpMeterExport))
	if err != nil {
		log.Fatal(err)
	}
	defer cancel(context.Background())

	creds := insecure.NewCredentials()
	if *tlsCertFile != "" && *tlsKeyFile != "" {
		var err error
		creds, err = credentials.NewServerTLSFromFile(*tlsCertFile, *tlsKeyFile)
		if err != nil {
			log.Exitf("failed to create tls credentials: %v", err)
		}
	}

	f, err := lemming.New(*target, *zapiAddr,
		lemming.WithConfigFile(*configFile),
		lemming.WithTransportCreds(creds),
		lemming.WithGRIBIAddr(*gribiAddr),
		lemming.WithGNMIAddr(*gnmiAddr),
		lemming.WithBGPPort(uint16(*bgpPort)),
		lemming.WithDataplane(*dplane),
		lemming.WithFaultAddr(*faultAddr),
		lemming.WithFaultInjection(*faultEnable),
		lemming.WithP4RTAddr(*p4rtAddr),
		lemming.WithDataplaneOpts(dplaneopts.WithSkipIPValidation()),
	)
	if err != nil {
		log.Exitf("Failed to start lemming: %v", err)
	}
	defer f.Stop()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	log.Info("lemming initialization complete")
	select {
	case <-c:
		log.Info("received sigint")
		return
	}
}
