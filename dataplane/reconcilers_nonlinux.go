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

//go:build !linux

// Package dataplane is an implementation of the dataplane HAL API.
package dataplane

import (
	"google.golang.org/grpc"

	"github.com/openconfig/lemming/dataplane/dplanerc"
	"github.com/openconfig/lemming/gnmi/reconciler"
)

func getReconcilers(conn grpc.ClientConnInterface, switchID uint64, cpuPortID uint64, contextID string) []reconciler.Reconciler {
	r := dplanerc.New(conn, switchID, cpuPortID, contextID)

	return []reconciler.Reconciler{
		reconciler.NewBuilder("inferface").WithStart(r.StartInterface).Build(),
	}
}
