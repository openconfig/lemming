// Copyright 2026 Microsoft Corporation
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

package luciusctl

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/reflection"

	saipb "github.com/openconfig/lemming/dataplane/proto/sai"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

type diagnosticServer struct {
	fwdpb.UnimplementedInfoServer
	saipb.UnimplementedSwitchServer
	calls    atomic.Int32
	tlsCalls atomic.Int32
	streams  atomic.Int32
}

func (s *diagnosticServer) record(ctx context.Context) {
	s.calls.Add(1)
	if p, ok := peer.FromContext(ctx); ok {
		if _, ok := p.AuthInfo.(credentials.TLSInfo); ok {
			s.tlsCalls.Add(1)
		}
	}
}

func (s *diagnosticServer) InfoList(ctx context.Context, _ *fwdpb.InfoListRequest) (*fwdpb.InfoListReply, error) {
	s.record(ctx)
	return &fwdpb.InfoListReply{Names: []string{"synthetic-object"}}, nil
}

func (s *diagnosticServer) InfoElement(ctx context.Context, _ *fwdpb.InfoElementRequest) (*fwdpb.InfoElementReply, error) {
	s.record(ctx)
	return &fwdpb.InfoElementReply{Content: "synthetic-content"}, nil
}

func (s *diagnosticServer) GetSwitchAttribute(ctx context.Context, _ *saipb.GetSwitchAttributeRequest) (*saipb.GetSwitchAttributeResponse, error) {
	s.record(ctx)
	return &saipb.GetSwitchAttributeResponse{}, nil
}

func serveDiagnostics(t *testing.T, host string, cert *tls.Certificate) (string, *diagnosticServer) {
	t.Helper()
	listener, err := net.Listen("tcp", net.JoinHostPort(host, "0"))
	if err != nil {
		t.Fatal(err)
	}
	service := &diagnosticServer{}
	opts := []grpc.ServerOption{grpc.StreamInterceptor(func(srv any, stream grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		service.streams.Add(1)
		return handler(srv, stream)
	})}
	if cert != nil {
		opts = append(opts, grpc.Creds(credentials.NewTLS(&tls.Config{
			Certificates: []tls.Certificate{*cert},
			MinVersion:   tls.VersionTLS12,
		})))
	}
	server := grpc.NewServer(opts...)
	fwdpb.RegisterInfoServer(server, service)
	saipb.RegisterSwitchServer(server, service)
	reflection.Register(server)
	done := make(chan error, 1)
	go func() { done <- server.Serve(listener) }()
	t.Cleanup(func() {
		server.Stop()
		if err := <-done; err != nil && err != grpc.ErrServerStopped {
			t.Errorf("test server: %v", err)
		}
	})
	return listener.Addr().String(), service
}

func diagnosticCertificates(t *testing.T) (string, map[string]tls.Certificate) {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now()
	root := &x509.Certificate{
		SerialNumber: big.NewInt(1), Subject: pkix.Name{CommonName: "Disposable diagnostic test CA"},
		NotBefore: now.Add(-48 * time.Hour), NotAfter: now.Add(24 * time.Hour),
		IsCA: true, BasicConstraintsValid: true, KeyUsage: x509.KeyUsageCertSign,
	}
	der, err := x509.CreateCertificate(rand.Reader, root, root, &key.PublicKey, key)
	if err != nil {
		t.Fatal(err)
	}
	caFile := filepath.Join(t.TempDir(), "ca.pem")
	if err := os.WriteFile(caFile, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}), 0600); err != nil {
		t.Fatal(err)
	}
	certs := make(map[string]tls.Certificate)
	for index, name := range []string{"valid", "self-signed", "expired", "wrong-name"} {
		leafKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			t.Fatal(err)
		}
		leaf := &x509.Certificate{
			SerialNumber: big.NewInt(int64(index + 2)),
			NotBefore:    now.Add(-time.Hour), NotAfter: now.Add(time.Hour),
			ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
			KeyUsage:    x509.KeyUsageDigitalSignature,
			IPAddresses: []net.IP{net.ParseIP("127.0.0.1"), net.ParseIP("::1")},
		}
		parent, signer := root, key
		switch name {
		case "self-signed":
			parent, signer = leaf, leafKey
		case "expired":
			leaf.NotBefore, leaf.NotAfter = now.Add(-24*time.Hour), now.Add(-time.Hour)
		case "wrong-name":
			leaf.IPAddresses = nil
			leaf.DNSNames = []string{"wrong-name.invalid"}
		}
		der, err := x509.CreateCertificate(rand.Reader, leaf, parent, &leafKey.PublicKey, signer)
		if err != nil {
			t.Fatal(err)
		}
		certs[name] = tls.Certificate{Certificate: [][]byte{der}, PrivateKey: leafKey}
	}
	return caFile, certs
}

func TestDiagnosticTransports(t *testing.T) {
	caFile, certs := diagnosticCertificates(t)
	tests := []struct {
		name    string
		cert    string
		flags   []string
		target  string
		wantErr string
	}{
		{name: "trusted-tls", cert: "valid", flags: []string{"--tls_ca_cert", caFile}},
		{name: "resolver-trusted-tls", cert: "valid", flags: []string{"--tls_ca_cert", caFile}, target: "dns:///"},
		{name: "unknown-ca", cert: "valid", wantErr: "certificate signed by unknown authority"},
		{name: "self-signed", cert: "self-signed", flags: []string{"--tls_ca_cert", caFile}, wantErr: "certificate signed by unknown authority"},
		{name: "expired", cert: "expired", flags: []string{"--tls_ca_cert", caFile}, wantErr: "certificate has expired"},
		{name: "wrong-name", cert: "wrong-name", flags: []string{"--tls_ca_cert", caFile}, wantErr: "IP SANs"},
		{name: "self-signed-bypass-rejected", cert: "self-signed", flags: []string{"--tls_skip_verify"}, wantErr: "no longer supported"},
		{name: "expired-bypass-rejected", cert: "expired", flags: []string{"--tls_skip_verify"}, wantErr: "no longer supported"},
		{name: "wrong-name-bypass-rejected", cert: "wrong-name", flags: []string{"--tls_skip_verify"}, wantErr: "no longer supported"},
		{name: "loopback-plaintext", flags: []string{"--insecure"}},
		{name: "plaintext-with-default-tls", wantErr: "TLS handshake"},
		{name: "both-options", flags: []string{"--insecure", "--tls_skip_verify"}, wantErr: "no longer supported"},
		{name: "resolver-plaintext", flags: []string{"--insecure"}, target: "dns:///", wantErr: "literal loopback"},
		{name: "hostname-plaintext", flags: []string{"--insecure"}, target: "localhost", wantErr: "literal loopback"},
		{name: "remote-plaintext", flags: []string{"--insecure"}, target: "192.0.2.1", wantErr: "literal loopback"},
		{name: "ca-with-plaintext", flags: []string{"--insecure", "--tls_ca_cert", caFile}, wantErr: "cannot be used"},
	}
	commands := [][]string{
		{"info", "list"},
		{"info", "get", "synthetic-object"},
		{"info", "lookup", "synthetic-object", "00"},
		{"info", "port-input", "synthetic-object", "00"},
		{"info", "port-output", "synthetic-object", "00"},
		{"sai", "get-attr", "switch", "1", "cpu_port"},
	}
	for _, command := range commands {
		for _, test := range tests {
			t.Run(strings.Join(command[:2], "-")+"/"+test.name, func(t *testing.T) {
				viper.Reset()
				t.Cleanup(viper.Reset)
				var cert *tls.Certificate
				if test.cert != "" {
					value := certs[test.cert]
					cert = &value
				}
				address, server := serveDiagnostics(t, "127.0.0.1", cert)
				switch test.target {
				case "dns:///":
					address = test.target + address
				case "localhost", "192.0.2.1":
					_, port, err := net.SplitHostPort(address)
					if err != nil {
						t.Fatal(err)
					}
					address = net.JoinHostPort(test.target, port)
				}
				cmd := New()
				cmd.SilenceUsage, cmd.SilenceErrors = true, true
				cmd.SetOut(io.Discard)
				cmd.SetErr(io.Discard)
				args := append([]string{"--address", address}, test.flags...)
				cmd.SetArgs(append(args, command...))
				ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
				defer cancel()
				err := cmd.ExecuteContext(ctx)
				if test.wantErr != "" {
					if err == nil || !strings.Contains(err.Error(), test.wantErr) || server.calls.Load() != 0 || server.streams.Load() != 0 {
						t.Fatalf("expected %q before RPC; err=%v calls=%d streams=%d", test.wantErr, err, server.calls.Load(), server.streams.Load())
					}

					return
				}
				if err != nil || server.calls.Load() != 1 {
					t.Fatalf("expected completed diagnostic RPC; err=%v calls=%d", err, server.calls.Load())
				}
				wantTLS := int32(0)
				if test.cert != "" {
					wantTLS = 1
				}
				if server.tlsCalls.Load() != wantTLS {
					t.Fatalf("TLS RPC count = %d, want %d", server.tlsCalls.Load(), wantTLS)
				}
			})
		}
	}
}

func TestLoopbackPlaintextBypassesProxy(t *testing.T) {
	for _, host := range []string{"127.0.0.1", "::1"} {
		t.Run(host, func(t *testing.T) {
			viper.Reset()
			t.Cleanup(viper.Reset)
			address, server := serveDiagnostics(t, host, nil)
			proxy, err := net.Listen("tcp", "127.0.0.1:0")
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { proxy.Close() })
			t.Setenv("HTTPS_PROXY", "http://"+proxy.Addr().String())
			t.Setenv("NO_PROXY", "")
			viper.Set("address", address)
			viper.Set("insecure", true)
			cmd := New()
			cmd.SetArgs([]string{"info", "list"})
			ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
			defer cancel()
			if err := cmd.ExecuteContext(ctx); err != nil || server.calls.Load() != 1 {
				t.Fatalf("direct local RPC failed: err=%v calls=%d", err, server.calls.Load())
			}
			if err := proxy.(*net.TCPListener).SetDeadline(time.Now().Add(20 * time.Millisecond)); err != nil {
				t.Fatal(err)
			}
			conn, err := proxy.Accept()
			if err == nil {
				conn.Close()
				t.Fatal("plaintext connection reached proxy")
			}
			if err, ok := err.(net.Error); !ok || !err.Timeout() {
				t.Fatalf("unexpected proxy listener error: %v", err)
			}
		})
	}
}

func TestSAICompletionTransportPolicy(t *testing.T) {
	tests := []struct {
		name       string
		target     string
		skipVerify bool
		wantError  bool
	}{
		{name: "local-plaintext"},
		{name: "legacy-bypass", skipVerify: true, wantError: true},
		{name: "remote-plaintext", target: "192.0.2.1:50000", wantError: true},
		{name: "resolver-plaintext", target: "dns:///127.0.0.1:50000", wantError: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			viper.Reset()
			t.Cleanup(viper.Reset)
			address, server := serveDiagnostics(t, "127.0.0.1", nil)
			if test.target != "" {
				address = test.target
			}
			viper.Set("address", address)
			viper.Set("insecure", true)
			viper.Set("tls_skip_verify", test.skipVerify)
			cmd, _, err := New().Find([]string{"sai", "get-attr"})
			if err != nil {
				t.Fatal(err)
			}
			values, directive := cmd.ValidArgsFunction(cmd, nil, "")
			if test.wantError {
				if directive&cobra.ShellCompDirectiveError == 0 || server.streams.Load() != 0 {
					t.Fatalf("expected completion error before reflection; directive=%v streams=%d", directive, server.streams.Load())
				}
			} else if directive&cobra.ShellCompDirectiveError != 0 || !slices.Contains(values, "switch") || server.streams.Load() != 1 {
				t.Fatalf("local completion failed: values=%v directive=%v streams=%d", values, directive, server.streams.Load())
			}
		})
	}
}
