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

package gnmi

import (
	"context"
	"fmt"

	"github.com/openconfig/lemming/gnmi/gnmit"
	"github.com/openconfig/lemming/gnmi/oc"
	"google.golang.org/grpc"
)

// Server is a fake gNMI implementation.
type Server struct {
	*gnmit.GNMIServer[*oc.Root]
}

func New(srv *grpc.Server, targetName string) (*Server, error) {
	configSchema, err := oc.Schema()
	if err != nil {
		return nil, fmt.Errorf("cannot create ygot schema object: %v", err)
	}
	stateSchema, err := oc.Schema()
	if err != nil {
		return nil, fmt.Errorf("cannot create ygot schema object: %v", err)
	}
	_, gnmiServer, err := gnmit.New[*oc.Root](context.Background(), configSchema, stateSchema, targetName, false)
	if err != nil {
		return nil, fmt.Errorf("failed to create gNMI server: %v", err)
	}
	gnmiServer.RegisterGRPCServer(srv)
	return &Server{gnmiServer}, nil
}
