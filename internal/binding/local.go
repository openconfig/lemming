// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package binding

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/openconfig/magna/lwotg"
	"github.com/openconfig/ondatra/binding"

	"github.com/openconfig/lemming"

	opb "github.com/openconfig/ondatra/proto"

	saipb "github.com/openconfig/lemming/dataplane/proto"
)

type LocalBind struct {
	binding.Binding
}

type localLemming struct {
	binding.AbstractDUT
	l *lemming.Device
}

type localMagna struct {
	binding.AbstractATE
	l *lwotg.Server
}

type linkMgr struct{}

type port interface {
	io.Reader
	io.Writer
}

// Reserve creates a new local binding.
func (lb *LocalBind) Reserve(ctx context.Context, tb *opb.Testbed, runTime, waitTime time.Duration, partial map[string]string) (*binding.Reservation, error) {
	resv := binding.Reservation{}
	for _, dut := range tb.Duts {
		l, err := lemming.New(dut.Id, "")
		if err != nil {
			return nil, err
		}
		boundLemming := &localLemming{
			l: l,
			AbstractDUT: binding.AbstractDUT{
				Dims: &binding.Dims{
					Name:          dut.Id,
					Vendor:        opb.Device_OPENCONFIG,
					HardwareModel: "LEMMING",
				},
			},
		}

		dplane, err := boundLemming.l.Dataplane().Conn()
		pc := saipb.NewPortClient(dplane)
		for _, port := range dut.Ports {
			pc.CreatePort(ctx, &saipb.CreatePortRequest{})

			boundLemming.AbstractDUT.Dims.Ports[port.Id] = &binding.Port{
				Name: "1",
			}
		}
		resv.DUTs[dut.Id] = boundLemming
	}
	// for _, ate := range tb.Ates {
	// 	otgSrv := lwotg.New()
	// 	resv.ATEs[ate.Id] = &localMagna{
	// 		l: otgSrv,
	// 	}
	// }
	// for _, l := range tb.Links {
	// 	aDev, aPort, err := lb.checkLink(l.A)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	bDev, bPort, err := lb.checkLink(l.B)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }

	return nil, nil
}

func (lb *LocalBind) checkLink(link string) (string, string, error) {
	split := strings.Split(link, ":")
	if len(split) != 2 {
		return "", "", fmt.Errorf("invalid link format %q, expected dut:port", link)
	}

	return split[0], split[1], nil
}

// Release releases the reserved testbed.
func (lb *LocalBind) Release(ctx context.Context) error {
	return nil
}
