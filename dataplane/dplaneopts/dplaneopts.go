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

package dplaneopts

import fwdpb "github.com/openconfig/lemming/proto/forwarding"

// Options configures the dataplane
type Options struct {
	// AddrPort is the address of the gRPC server.
	AddrPort string
	// Reconcilation enabes gNMI reconcilation.
	Reconcilation bool
	// HostifNetDevType is the fwdpb type for the saipb hostif netdev types.
	HostifNetDevType fwdpb.PortType
	// PortType is the fwdpb type for the port type.
	PortType fwdpb.PortType
	// PortConfigFile is the path of the port config.
	PortConfigFile string
	// PortMap maps the modeled port name (Ethernet1/1/1) to Linux port name (eth1).
	PortMap map[string]string
}

// Option exposes additional configuration for the dataplane.
type Option func(*Options)

// WithAddrPort sets the address of the dataplane gRPC server
// Default: 127.0.0.1:0
func WithAddrPort(addrPort string) Option {
	return func(o *Options) {
		o.AddrPort = addrPort
	}
}

// WithReconcilation enables the gNMI reconcilation.
// Default: true
func WithReconcilation(rec bool) Option {
	return func(o *Options) {
		o.Reconcilation = rec
	}
}

// WithHostifNetDevPortType sets the lucius port type for saipb hostif NETDEV.
// Default: fwdpb.PortType_PORT_TYPE_TAP
func WithHostifNetDevPortType(t fwdpb.PortType) Option {
	return func(o *Options) {
		o.HostifNetDevType = t
	}
}

// WithPortType sets the lucius port type for saipb ports.
// Default: fwdpb.PortType_PORT_TYPE_KERNEL
func WithPortType(t fwdpb.PortType) Option {
	return func(o *Options) {
		o.PortType = t
	}
}

// WithPortConfigFile sets the path of the port config file.
// Default: none
func WithPortConfigFile(file string) Option {
	return func(o *Options) {
		o.PortConfigFile = file
	}
}

// WithPortMap configure a map from port name to Linux network device to allow flexible port naming. (eg Ethernet8 -> eth1)
// Default: none
func WithPortMap(m map[string]string) Option {
	return func(o *Options) {
		o.PortMap = m
	}
}

// Port contains configuration data for a single port.
type Port struct {
	Lanes string `json:"lanes"`
}

// PortConfig contains configuration data for the dataplane ports.
type PortConfig struct {
	Ports map[string]*Port `json:"PORT"`
}

// ResolveOpts creates an option struct from the opts.
func ResolveOpts(opts ...Option) *Options {
	resolved := &Options{
		AddrPort:         "127.0.0.1:0",
		Reconcilation:    true,
		HostifNetDevType: fwdpb.PortType_PORT_TYPE_TAP,
		PortType:         fwdpb.PortType_PORT_TYPE_KERNEL,
		PortMap:          map[string]string{},
	}

	for _, opt := range opts {
		opt(resolved)
	}
	return resolved
}
