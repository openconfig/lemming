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

package gnmi

import (
	"context"
	"fmt"
	"testing"

	gpb "github.com/openconfig/gnmi/proto/gnmi"
	"github.com/openconfig/lemming/gnmi/gnmiclient"
	"github.com/openconfig/lemming/gnmi/oc"
	"github.com/openconfig/lemming/gnmi/oc/ocpath"
	"github.com/openconfig/ygnmi/ygnmi"
	"github.com/openconfig/ygot/ygot"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/local"
)

func BenchmarkGNMISet(b *testing.B) {
	gnmiServer, err := newServer(context.Background(), "local", true)
	if err != nil {
		b.Fatalf("cannot create server, got err: %v", err)
	}
	addr, err := startServer(gnmiServer)
	if err != nil {
		b.Fatalf("cannot start server, got err: %v", err)
	}
	defer gnmiServer.c.Stop()
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(local.NewCredentials()))
	if err != nil {
		b.Fatalf("cannot dial gNMI server, %v", err)
	}
	configClient, err := ygnmi.NewClient(gpb.NewGNMIClient(conn), ygnmi.WithTarget("local"))
	if err != nil {
		b.Fatalf("failed to create client: %v", err)
	}
	stateClient, err := ygnmi.NewClient(gnmiServer.LocalClient(), ygnmi.WithTarget("local"))
	if err != nil {
		b.Fatalf("failed to create client: %v", err)
	}

	tests := []struct {
		desc string
		op   func(name string, i int) error
		op2  func(name string, i int) error
	}{{
		desc: "config replace",
		op: func(name string, i int) error {
			_, err := gnmiclient.Replace[*oc.Interface](context.Background(), configClient, ocpath.Root().Interface(name).Config(),
				&oc.Interface{
					Description:  ygot.String(fmt.Sprintf("iteration %d", i)),
					Enabled:      ygot.Bool(true),
					Type:         oc.IETFInterfaces_InterfaceType_fast,
					Mtu:          ygot.Uint16(42),
					LoopbackMode: oc.Interfaces_LoopbackModeType_FACILITY,
					Tpid:         oc.VlanTypes_TPID_TYPES_TPID_0X8100,
				},
			)
			return err
		},
	}, {
		desc: "state replace",
		op: func(name string, i int) error {
			_, err := gnmiclient.Replace(context.Background(), stateClient, ocpath.Root().Interface(name).State(),
				&oc.Interface{
					Description:  ygot.String(fmt.Sprintf("iteration %d", i)),
					Enabled:      ygot.Bool(true),
					Type:         oc.IETFInterfaces_InterfaceType_fast,
					Mtu:          ygot.Uint16(42),
					LoopbackMode: oc.Interfaces_LoopbackModeType_FACILITY,
					Tpid:         oc.VlanTypes_TPID_TYPES_TPID_0X8100,
				},
			)
			return err
		},
	}, {
		desc: "config update",
		op: func(name string, i int) error {
			if _, err := gnmiclient.Update[string](context.Background(), configClient, ocpath.Root().Interface(name).Description().Config(), fmt.Sprintf("iteration %d", i)); err != nil {
				return err
			}
			if _, err := gnmiclient.Update[bool](context.Background(), configClient, ocpath.Root().Interface(name).Enabled().Config(), true); err != nil {
				return err
			}
			if _, err := gnmiclient.Update[oc.E_IETFInterfaces_InterfaceType](context.Background(), configClient, ocpath.Root().Interface(name).Type().Config(), oc.IETFInterfaces_InterfaceType_fast); err != nil {
				return err
			}
			if _, err := gnmiclient.Update[uint16](context.Background(), configClient, ocpath.Root().Interface(name).Mtu().Config(), 42); err != nil {
				return err
			}
			if _, err := gnmiclient.Update[oc.E_Interfaces_LoopbackModeType](context.Background(), configClient, ocpath.Root().Interface(name).LoopbackMode().Config(), oc.Interfaces_LoopbackModeType_FACILITY); err != nil {
				return err
			}
			if _, err := gnmiclient.Update[oc.E_VlanTypes_TPID_TYPES](context.Background(), configClient, ocpath.Root().Interface(name).Tpid().Config(), oc.VlanTypes_TPID_TYPES_TPID_0X8100); err != nil {
				return err
			}
			return nil
		},
	}, {
		desc: "state update",
		op: func(name string, i int) error {
			if _, err := gnmiclient.Update(context.Background(), stateClient, ocpath.Root().Interface(name).Description().State(), fmt.Sprintf("iteration %d", i)); err != nil {
				return err
			}
			if _, err := gnmiclient.Update(context.Background(), stateClient, ocpath.Root().Interface(name).Enabled().State(), true); err != nil {
				return err
			}
			if _, err := gnmiclient.Update(context.Background(), stateClient, ocpath.Root().Interface(name).Type().State(), oc.IETFInterfaces_InterfaceType_fast); err != nil {
				return err
			}
			if _, err := gnmiclient.Update(context.Background(), stateClient, ocpath.Root().Interface(name).Mtu().State(), 42); err != nil {
				return err
			}
			if _, err := gnmiclient.Update(context.Background(), stateClient, ocpath.Root().Interface(name).LoopbackMode().State(), oc.Interfaces_LoopbackModeType_FACILITY); err != nil {
				return err
			}
			if _, err := gnmiclient.Update(context.Background(), stateClient, ocpath.Root().Interface(name).Tpid().State(), oc.VlanTypes_TPID_TYPES_TPID_0X8100); err != nil {
				return err
			}
			return nil
		},
	}, {
		desc: "config replace and delete",
		op: func(name string, i int) error {
			_, err := gnmiclient.Replace[*oc.Interface](context.Background(), configClient, ocpath.Root().Interface(name).Config(), &oc.Interface{Description: ygot.String(fmt.Sprintf("iteration %d", i))})
			return err
		},
		op2: func(name string, i int) error {
			_, err := gnmiclient.Delete[*oc.Interface](context.Background(), configClient, ocpath.Root().Interface(name).Config())
			return err
		},
	}, {
		desc: "state replace and delete",
		op: func(name string, i int) error {
			_, err := gnmiclient.Replace(context.Background(), configClient, ocpath.Root().Interface(name).State(), &oc.Interface{Description: ygot.String(fmt.Sprintf("iteration %d", i))})
			return err
		},
		op2: func(name string, i int) error {
			_, err := gnmiclient.Delete(context.Background(), configClient, ocpath.Root().Interface(name).State())
			return err
		},
	}}

	interfaceN := 10
	for _, bb := range tests {
		b.Run(bb.desc, func(b *testing.B) {
			for i := 0; i != b.N; i++ {
				for j := 0; j != interfaceN; j++ {
					if err := bb.op(fmt.Sprintf("eth%d", j), i); err != nil {
						b.Fatal(err)
					}
				}
				if bb.op2 != nil {
					for j := 0; j != interfaceN; j++ {
						if err := bb.op2(fmt.Sprintf("eth%d", j), i); err != nil {
							b.Fatal(err)
						}
					}
				}
			}
		})
	}
}
