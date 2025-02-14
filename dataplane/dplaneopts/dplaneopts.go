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

import (
	"os"

	"gopkg.in/yaml.v3"

	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

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
	// HardwareProfile is the "hardware" like config options.
	HardwareProfile *HardwareProfile
	// SkipIPValidation skips droping packets with invalid src or dst IPs.
	SkipIPValidation bool
}

// Option exposes additional configuration for the dataplane.
type Option func(*Options) error

// WithAddrPort sets the address of the dataplane gRPC server
// Default: 127.0.0.1:0
func WithAddrPort(addrPort string) Option {
	return func(o *Options) error {
		o.AddrPort = addrPort
		return nil
	}
}

// WithReconcilation enables the gNMI reconcilation.
// Default: true
func WithReconcilation(rec bool) Option {
	return func(o *Options) error {
		o.Reconcilation = rec
		return nil
	}
}

// WithHostifNetDevPortType sets the lucius port type for saipb hostif NETDEV.
// Default: fwdpb.PortType_PORT_TYPE_TAP
func WithHostifNetDevPortType(t fwdpb.PortType) Option {
	return func(o *Options) error {
		o.HostifNetDevType = t
		return nil
	}
}

// WithPortType sets the lucius port type for saipb ports.
// Default: fwdpb.PortType_PORT_TYPE_KERNEL
func WithPortType(t fwdpb.PortType) Option {
	return func(o *Options) error {
		o.PortType = t
		return nil
	}
}

// WithHardwareProfile sets location of the hardware profile config
func WithHardwareProfile(file string) Option {
	return func(o *Options) error {
		if file == "" {
			return nil
		}
		data, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		return yaml.Unmarshal(data, o.HardwareProfile)
	}
}

func WithSkipIPValidation() Option {
	return func(o *Options) error {
		o.SkipIPValidation = true
		return nil
	}
}

// ResolveOpts creates an option struct from the opts.
func ResolveOpts(opts ...Option) *Options {
	resolved := &Options{
		AddrPort:         "127.0.0.1:0",
		Reconcilation:    true,
		HostifNetDevType: fwdpb.PortType_PORT_TYPE_TAP,
		PortType:         fwdpb.PortType_PORT_TYPE_KERNEL,
		HardwareProfile:  &HardwareProfile{},
	}

	for _, opt := range opts {
		opt(resolved)
	}
	return resolved
}

type FECMode struct {
	Speed int      // Speed in Gbps.
	Lanes int      // Number of lanes.
	Modes []string // Supported modes.
}

// HardwareProfile is the "hardware" like config options.
type HardwareProfile struct {
	// FECModes configures the support FECMode for a given speed and lanes combination.
	FECModes []*FECMode `yaml:"fec_modes"`
}
