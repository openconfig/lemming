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

// Package tabletestutil consists of routines used to test Lucius tables like
// prefix match, exact match and flow match tables.
package tabletestutil

import (
	"github.com/openconfig/lemming/dataplane/forwarding/fwdtable/tableutil"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

// KeyDescFields creates a key descriptor from a set of packet fields.
func KeyDescFields(fields []fwdpb.PacketFieldNum) tableutil.KeyDesc {
	var desc []*fwdpb.PacketFieldId
	for _, field := range fields {
		field := field
		desc = append(desc, &fwdpb.PacketFieldId{
			Field: &fwdpb.PacketField{
				FieldNum: field,
			}})
	}
	return tableutil.MakeKeyDesc(desc)
}

// KeyDescBytes creates a key descriptor from a set of packet bytes.
func KeyDescBytes(fields []fwdpb.PacketBytes) tableutil.KeyDesc {
	var desc []*fwdpb.PacketFieldId
	for _, f := range fields {
		copied := f
		desc = append(desc, &fwdpb.PacketFieldId{
			Bytes: &copied,
		})
	}
	return tableutil.MakeKeyDesc(desc)
}

// ActionDesc creates a desc for a set of actions.
func ActionDesc() []*fwdpb.ActionDesc {
	desc := &fwdpb.ActionDesc{
		ActionType: fwdpb.ActionType_ACTION_TYPE_DROP,
	}
	return []*fwdpb.ActionDesc{desc}
}
