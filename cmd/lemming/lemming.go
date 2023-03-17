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
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"math/big"

	log "github.com/golang/glog"
	"github.com/openconfig/lemming"
	"github.com/openconfig/lemming/sysrib"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	gnmiAddr  = flag.String("gnmi", ":9339", "gNMI listen address")
	gribiAddr = flag.String("gribi", ":9340", "gRIBI listen address")
	target    = pflag.String("target", "fakedut", "name of the fake target")
	enableTLS = pflag.Bool("enable_tls", false, "Controls whether to enable TLS for gNXI services. If enabled and TLS key/cert path unspecified, a generated cert is used.")
	zapiAddr  = pflag.String("zapi_addr", sysrib.ZAPIAddr, "Custom ZAPI address: use unix:/tmp/zserv.api for a temp.")
)

func main() {
	pflag.Bool("enable_dataplane", false, "Controls whether to enable dataplane")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	creds, err := newCreds()
	if err != nil {
		log.Fatalf("failed to create credentials: %v", err)
	}

	f, err := lemming.New(*target, *zapiAddr, lemming.WithTransportCreds(creds), lemming.WithGRIBIAddr(*gribiAddr), lemming.WithGNMIAddr(*gnmiAddr))
	if err != nil {
		log.Fatalf("Failed to start lemming: %v", err)
	}
	defer f.Stop()

	log.Info("lemming initialization complete")
	select {}
}

// newCreds returns either insecure or tls credentials, depending the enable_tls flag.
// TODO: figure out long term plan for certs, this implementation is here to unblock using Ondatra KNEBind.
func newCreds() (credentials.TransportCredentials, error) {
	if !*enableTLS {
		return insecure.NewCredentials(), nil
	}
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(1234),
	}

	certDer, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		return nil, err
	}
	certPem := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDer})

	keyDer, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return nil, err
	}
	keyPem := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: keyDer})

	serverCert, err := tls.X509KeyPair(certPem, keyPem)
	if err != nil {
		return nil, err
	}
	return credentials.NewTLS(&tls.Config{
		MinVersion:   tls.VersionTLS13,
		Certificates: []tls.Certificate{serverCert},
	}), nil
}
