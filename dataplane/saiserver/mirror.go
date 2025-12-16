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

package saiserver

import (
	"context"

	saipb "github.com/openconfig/lemming/dataplane/proto/sai"
)

func (m *mirror) CreateMirrorSession(ctx context.Context, req *saipb.CreateMirrorSessionRequest) (*saipb.CreateMirrorSessionResponse, error) {
	return &saipb.CreateMirrorSessionResponse{}, nil
}

func (m *mirror) RemoveMirrorSession(ctx context.Context, req *saipb.RemoveMirrorSessionRequest) (*saipb.RemoveMirrorSessionResponse, error) {
	return &saipb.RemoveMirrorSessionResponse{}, nil
}

func (m *mirror) SetMirrorSessionAttribute(ctx context.Context, req *saipb.SetMirrorSessionAttributeRequest) (*saipb.SetMirrorSessionAttributeResponse, error) {
	return &saipb.SetMirrorSessionAttributeResponse{}, nil
}

func (m *mirror) GetMirrorSessionAttribute(ctx context.Context, req *saipb.GetMirrorSessionAttributeRequest) (*saipb.GetMirrorSessionAttributeResponse, error) {
	return &saipb.GetMirrorSessionAttributeResponse{}, nil
}
