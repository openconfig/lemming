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

package transport

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/viper"
)

func TestLoopbackAddress(t *testing.T) {
	tests := []struct {
		address string
		want    string
	}{
		{address: "127.0.0.1:50000", want: "127.0.0.1:50000"},
		{address: "127.255.255.254:1", want: "127.255.255.254:1"},
		{address: "[::1]:65535", want: "[::1]:65535"},
		{address: "[0:0:0:0:0:0:0:1]:05000", want: "[::1]:5000"},
		{address: "[::ffff:127.0.0.1]:50000", want: "127.0.0.1:50000"},
		{address: "localhost:50000"},
		{address: "localhost.:50000"},
		{address: "127.0.0.1.example:50000"},
		{address: "192.0.2.1:50000"},
		{address: "10.0.0.1:50000"},
		{address: "0.0.0.0:50000"},
		{address: "[::]:50000"},
		{address: "[2001:db8::1]:50000"},
		{address: "[::ffff:192.0.2.1]:50000"},
		{address: "[::1%lo]:50000"},
		{address: "dns:///127.0.0.1:50000"},
		{address: "passthrough:///127.0.0.1:50000"},
		{address: "unix:///tmp/lucius.sock"},
		{address: "http://127.0.0.1:50000"},
		{address: "127.1:50000"},
		{address: "2130706433:50000"},
		{address: "0177.0.0.1:50000"},
		{address: "127.0.0.1"},
		{address: "::1:50000"},
		{address: "127.0.0.1:0"},
		{address: "127.0.0.1:65536"},
		{address: "127.0.0.1:http"},
		{address: "127.0.0.1:+1"},
		{address: "127.0.0.1:-1"},
		{address: "127.0.0.1:"},
		{address: ""},
	}
	for _, test := range tests {
		t.Run(test.address, func(t *testing.T) {
			got, err := loopbackAddress(test.address)
			if test.want == "" {
				if err == nil {
					t.Fatalf("accepted prohibited plaintext target %q as %q", test.address, got)
				}
				return
			}
			if err != nil || got != test.want {
				t.Fatalf("loopbackAddress(%q) = %q, %v; want %q", test.address, got, err, test.want)
			}
		})
	}
}

func TestDialConfigurationErrors(t *testing.T) {
	dir := t.TempDir()
	invalidCA := filepath.Join(dir, "invalid.pem")
	if err := os.WriteFile(invalidCA, []byte("not a certificate"), 0600); err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name    string
		values  map[string]any
		wantErr string
	}{
		{name: "missing-address", wantErr: "--address is required"},
		{name: "legacy-bypass", values: map[string]any{"tls_skip_verify": true}, wantErr: "no longer supported"},
		{name: "both-options", values: map[string]any{"tls_skip_verify": true, "insecure": true}, wantErr: "no longer supported"},
		{name: "remote-plaintext", values: map[string]any{"address": "192.0.2.1:50000", "insecure": true}, wantErr: "literal loopback"},
		{name: "ca-with-plaintext", values: map[string]any{"address": "127.0.0.1:50000", "insecure": true, "tls_ca_cert": invalidCA}, wantErr: "cannot be used"},
		{name: "missing-ca", values: map[string]any{"address": "dataplane.example:50000", "tls_ca_cert": filepath.Join(dir, "missing.pem")}, wantErr: "read --tls_ca_cert"},
		{name: "invalid-ca", values: map[string]any{"address": "dataplane.example:50000", "tls_ca_cert": invalidCA}, wantErr: "no valid PEM certificates"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			viper.Reset()
			t.Cleanup(viper.Reset)
			for key, value := range test.values {
				viper.Set(key, value)
			}
			conn, err := Dial()
			if conn != nil {
				conn.Close()
				t.Fatal("invalid configuration created a connection")
			}
			if err == nil || !strings.Contains(err.Error(), test.wantErr) {
				t.Fatalf("Dial() error = %v; want %q", err, test.wantErr)
			}
		})
	}
}
