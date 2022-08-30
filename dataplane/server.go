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

package dataplane

import (
	"context"

	dpb "github.com/openconfig/lemming/proto/dataplane"
)

const (
	Port = 6443
)

func New() *Server {
	srv := &Server{}

	return srv
}

// Server is an implementation of Dataplane HAL API.
type Server struct {
	dpb.UnimplementedHALServer
}

func (s *Server) UpdatePort(context.Context, *dpb.UpdatePortRequest) (*dpb.UpdatePortResponse, error) {

}
