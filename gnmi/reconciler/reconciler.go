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

// Package reconciler contains a common interface for gNMI reconciler.
package reconciler

import (
	"context"

	gpb "github.com/openconfig/gnmi/proto/gnmi"
	"github.com/openconfig/lemming/gnmi/oc"
)

// Reconciler is a common interface for gNMI reconciler.
// Reconcilers are responsible for watching config and state and performing any necessary reconciliation.
type Reconciler interface {
	// ID the id for the reconciler.
	ID() string
	// Start start the reconciliation loop.
	// The client and target are connected to the local gNMI cache.
	// An error returned during Start will cause lemming to exit, so Start should return non-retriable errors.
	Start(ctx context.Context, client gpb.GNMIClient, target string) error
	// Stop stops the reconciliation loop.
	Stop(context.Context) error
	// Validate is called after a SetRequest is checked for schema compliance, but before data is written to the cache.
	// Reconcilers can validate a SetRequest for semantic correctness, config that can never be reconciled.
	// The Validate func is only called if the SetRequest contains paths which match the ValidationPaths.
	Validate(intendedConfig *oc.Root) error
	// ValidationPaths returns the set of path prefixes that the reconciler can validate.
	ValidationPaths() []*gpb.Path
}
