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
	"flag"

	"github.com/openconfig/lemming"
	"github.com/openconfig/lemming/sysrib"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	log "github.com/golang/glog"
)

var (
	gnmiAddr    = flag.String("gnmi", ":9339", "gNMI listen address")
	gribiAddr   = flag.String("gribi", ":9340", "gRIBI listen address")
	target      = pflag.String("target", "fakedut", "name of the fake target")
	tlsKeyFile  = pflag.String("tls_key_file", "", "Controls whether to enable TLS for gNXI services. If unspecified, insecure credentials are used.")
	tlsCertFile = pflag.String("tls_cert_file", "", "Controls whether to enable TLS for gNXI services. If unspecified, insecure credentials are used.")
	zapiAddr    = pflag.String("zapi_addr", sysrib.ZAPIAddr, "Custom ZAPI address: use unix:/tmp/zserv.api for a temp.")
)

func main() {
	pflag.Bool("enable_dataplane", false, "Controls whether to enable dataplane")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	creds := insecure.NewCredentials()
	if *tlsCertFile != "" && *tlsKeyFile != "" {
		var err error
		creds, err = credentials.NewServerTLSFromFile(*tlsCertFile, *tlsKeyFile)
		if err != nil {
			log.Fatalf("failed to create tls credentials: %v", err)
		}
	}

	f, err := lemming.New(*target, *zapiAddr, lemming.WithTransportCreds(creds), lemming.WithGRIBIAddr(*gribiAddr), lemming.WithGNMIAddr(*gnmiAddr))
	if err != nil {
		log.Fatalf("Failed to start lemming: %v", err)
	}
	defer f.Stop()

	log.Info("lemming initialization complete")
	select {}
}
