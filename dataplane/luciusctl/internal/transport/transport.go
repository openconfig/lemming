// Copyright 2026 Google LLC
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

// Package transport provides verified TLS or local-only plaintext connections.
package transport

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"net/netip"
	"os"
	"strconv"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// Dial creates a channel using the configured diagnostic transport policy.
func Dial() (*grpc.ClientConn, error) {
	if viper.GetBool("tls_skip_verify") {
		return nil, fmt.Errorf("--tls_skip_verify is no longer supported; use --tls_ca_cert to trust the server's CA")
	}
	address := viper.GetString("address")
	if address == "" {
		return nil, fmt.Errorf("--address is required")
	}
	caFile := viper.GetString("tls_ca_cert")
	if viper.GetBool("insecure") {
		if caFile != "" {
			return nil, fmt.Errorf("--tls_ca_cert cannot be used with --insecure")
		}
		address, err := loopbackAddress(address)
		if err != nil {
			return nil, err
		}
		// Pin every connection to the validated IP, bypassing resolvers and proxies.
		dialer := &net.Dialer{}
		return grpc.NewClient("passthrough:///"+address,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
				return dialer.DialContext(ctx, "tcp", address)
			}),
		)
	}

	config := &tls.Config{MinVersion: tls.VersionTLS12}
	if caFile != "" {
		certs, err := os.ReadFile(caFile)
		if err != nil {
			return nil, fmt.Errorf("read --tls_ca_cert: %w", err)
		}
		roots, err := x509.SystemCertPool()
		if err != nil {
			return nil, fmt.Errorf("load system certificate roots: %w", err)
		}
		if !roots.AppendCertsFromPEM(certs) {
			return nil, fmt.Errorf("--tls_ca_cert contains no valid PEM certificates")
		}
		config.RootCAs = roots
	}
	return grpc.NewClient(address, grpc.WithTransportCredentials(credentials.NewTLS(config)))
}

func loopbackAddress(address string) (string, error) {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return "", fmt.Errorf("--insecure requires a literal loopback IP and port: %w", err)
	}
	ip, err := netip.ParseAddr(host)
	if err != nil || ip.Zone() != "" || !ip.Unmap().IsLoopback() {
		return "", fmt.Errorf("--insecure requires a literal loopback IP; hostnames, remote IPs and resolver targets are not supported")
	}
	number, err := strconv.ParseUint(port, 10, 16)
	if err != nil || number == 0 {
		return "", fmt.Errorf("--insecure requires a port between 1 and 65535")
	}
	return net.JoinHostPort(ip.Unmap().String(), strconv.FormatUint(number, 10)), nil
}
