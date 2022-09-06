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

package actions

import (
	"github.com/openconfig/lemming/dataplane/forwarding/fwdaction"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdcontext"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

// noOPBuilder builds no operation actions.
type noOPBuilder struct{}

func init() {
	// Register a builder for the continue action type.
	fwdaction.Register(fwdpb.ActionType_ACTION_TYPE_CONTINUE, &noOPBuilder{})
}

// Build creates a new noOP action. The builder just returns nil, which is then
// ignored by the action processing framework.
func (noOPBuilder) Build(*fwdpb.ActionDesc, *fwdcontext.Context) (fwdaction.Action, error) {
	return nil, nil
}
