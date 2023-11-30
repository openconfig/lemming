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

package engine

import (
	"context"

	"github.com/openconfig/lemming/dataplane/forwarding"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdcontext"

	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

const (
	FIBV4Table            = "fib-v4"
	FIBV6Table            = "fib-v6"
	SRCMACTable           = "port-mac"
	FIBSelectorTable      = "fib-selector"
	NeighborTable         = "neighbor"
	NHGTable              = "nhg-table"
	NHTable               = "nh-table"
	layer2PuntTable       = "layer2-punt"
	layer3PuntTable       = "layer3-punt"
	arpPuntTable          = "arp-punt"
	PreIngressActionTable = "preingress-table"
	IngressActionTable    = "ingress-table"
	EgressActionTable     = "egress-action-table"
	NHActionTable         = "nh-action"
)

// Engine contains a routing context and methods to manage it.
type Engine struct {
	*forwarding.Server
	id       string
	cancelFn func()
}

// New creates a new engine and sets up the forwarding tables.
func New(ctx context.Context) (*Engine, error) {
	e := &Engine{
		id:     "lucius",
		Server: forwarding.New("engine"),
	}

	ctx, e.cancelFn = context.WithCancel(ctx)

	_, err := e.Server.ContextCreate(ctx, &fwdpb.ContextCreateRequest{
		ContextId: &fwdpb.ContextId{Id: e.id},
	})
	if err != nil {
		return nil, err
	}

	return e, nil
}

func (e *Engine) Reset(ctx context.Context) error {
	e.cancelFn()
	_, err := e.Server.ContextDelete(ctx, &fwdpb.ContextDeleteRequest{
		ContextId: &fwdpb.ContextId{Id: e.id},
	})
	if err != nil {
		return err
	}

	_, err = e.Server.ContextCreate(ctx, &fwdpb.ContextCreateRequest{
		ContextId: &fwdpb.ContextId{Id: e.id},
	})
	if err != nil {
		return err
	}

	_, e.cancelFn = context.WithCancel(ctx)
	return nil
}

// Context returns the forward context assoicated with the engine.
func (e *Engine) Context() (*fwdcontext.Context, error) {
	return e.Server.FindContext(&fwdpb.ContextId{Id: e.ID()})
}

// ID returns the engine's forwarding context id.
func (e *Engine) ID() string {
	return e.id
}
