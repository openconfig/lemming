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

	"github.com/openconfig/lemming/gnmi/fakedevice"
	"github.com/openconfig/lemming/gnmi/gnmit"
	"github.com/openconfig/lemming/gnmi/internal/oc"

	gnmipb "github.com/openconfig/gnmi/proto/gnmi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server is a fake gNMI implementation.
type Server struct {
	*gnmit.GNMIServer
	s            *grpc.Server
	Responses    [][]*gnmipb.SubscribeResponse
	GetResponses []interface{}
	Errs         []error
}

func createGNMIServer(targetName string) (*gnmit.GNMIServer, error) {
	schema, err := oc.Schema()
	if err != nil {
		return nil, fmt.Errorf("cannot create ygot schema object: %v", err)
	}
	if err := gnmit.SetupSchema(schema); err != nil {
		return nil, fmt.Errorf("gnmi: cannot setup ygot schema object: %v", err)
	}
	_, gnmiServer, err := gnmit.NewServer(context.Background(), schema, targetName, false, fakedevice.Tasks(targetName))
	if err != nil {
		return nil, fmt.Errorf("failed to create gNMI server: %v", err)
	}
	if _, err := gnmit.StartDatastoreServer(gnmiServer); err != nil {
		return nil, fmt.Errorf("failed to start datastore server: %v", err)
	}
	return gnmiServer, nil
}

// New returns a new fake gNMI server.
func New(s *grpc.Server, targetName string) (*Server, error) {
	gnmiServer, err := createGNMIServer(targetName)
	if err != nil {
		return nil, fmt.Errorf("gnmi: cannot create gNMI server: %v", err)
	}

	srv := &Server{
		GNMIServer: gnmiServer,
		s:          s,
	}
	gnmipb.RegisterGNMIServer(s, srv)
	return srv, nil
}

func (s *Server) Capabilities(ctx context.Context, req *gnmipb.CapabilityRequest) (*gnmipb.CapabilityResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "Fake Unimplemented")
}

func (s *Server) Get(ctx context.Context, req *gnmipb.GetRequest) (*gnmipb.GetResponse, error) {
	if len(s.GetResponses) == 0 {
		return nil, status.Errorf(codes.Unimplemented, "Fake Unimplemented")
	}
	resp := s.GetResponses[0]
	s.GetResponses = s.GetResponses[1:]
	switch v := resp.(type) {
	case error:
		return nil, v
	case *gnmipb.GetResponse:
		return v, nil
	default:
		return nil, status.Errorf(codes.DataLoss, "Unknown message type: %T", resp)
	}
}
