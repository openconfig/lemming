// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package fwdconfig

import fwdpb "github.com/openconfig/lemming/proto/forwarding"

// TableEntryAddRequest builds TableEntryAddRequests.
type TableEntryAddRequestBuilder struct {
	contextID string
	tableID   string
	entries   []*EntryDescBuilder
	actions   [][]*ActionBuilder
}

// TableEntryAddRequest creates a new TableEntryAddRequestBuilder.
func TableEntryAddRequest(ctxID, tableID string) *TableEntryAddRequestBuilder {
	return &TableEntryAddRequestBuilder{
		contextID: ctxID,
		tableID:   tableID,
	}
}

// AppendEntry adds an entry and the actions to the requests.
func (b *TableEntryAddRequestBuilder) AppendEntry(entry *EntryDescBuilder, actions ...*ActionBuilder) *TableEntryAddRequestBuilder {
	b.entries = append(b.entries, entry)
	b.actions = append(b.actions, actions)
	return b
}

// AppendEntry adds the actions to the requests.
func (b *TableEntryAddRequestBuilder) AppendActions(actions ...*ActionBuilder) *TableEntryAddRequestBuilder {
	b.actions = append(b.actions, actions)
	return b
}

// Build returns a new TableEntryAddRequest.
func (b TableEntryAddRequestBuilder) Build() *fwdpb.TableEntryAddRequest {
	req := &fwdpb.TableEntryAddRequest{
		ContextId: &fwdpb.ContextId{Id: b.contextID},
		TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: b.tableID}},
		Entries:   []*fwdpb.TableEntryAddRequest_Entry{},
	}
	for i, entry := range b.entries {
		var acts []*fwdpb.ActionDesc
		for _, act := range b.actions[i] {
			acts = append(acts, act.Build())
		}
		req.Entries = append(req.Entries, &fwdpb.TableEntryAddRequest_Entry{
			EntryDesc: entry.Build(),
			Actions:   acts,
		})
	}
	return req
}

type entryDescBuilder interface {
	set(*fwdpb.EntryDesc)
}

// EntryDesc returns a new EntryDescBuilder.
type EntryDescBuilder struct {
	entryBuild entryDescBuilder
}

// EntryDesc returns a new EntryDescBuilder.
func EntryDesc(e entryDescBuilder) *EntryDescBuilder {
	return &EntryDescBuilder{
		entryBuild: e,
	}
}

// Build returns a new EntryDesc.
func (edb EntryDescBuilder) Build() *fwdpb.EntryDesc {
	entry := &fwdpb.EntryDesc{}
	edb.entryBuild.set(entry)
	return entry
}

// ExactEntryBuilder builds exact table entries.
type ExactEntryBuilder struct {
	fields []*PacketFieldBytesBuilder
}

// ExtactEntry creates a new exact entry builder.
func ExtactEntry(fields ...*PacketFieldBytesBuilder) *ExactEntryBuilder {
	return &ExactEntryBuilder{
		fields: fields,
	}
}

func (eeb ExactEntryBuilder) set(ed *fwdpb.EntryDesc) {
	exact := &fwdpb.ExactEntryDesc{}
	for _, b := range eeb.fields {
		exact.Fields = append(exact.Fields, b.Build())
	}

	ed.Entry = &fwdpb.EntryDesc_Exact{
		Exact: exact,
	}
}

// FlowEntryBuilder builds flow table entries.
type FlowEntryBuilder struct {
	fields []*PacketFieldMaskedBytesBuilder
}

// FlowEntry creates a new flow entry builder.
func FlowEntry(fields ...*PacketFieldMaskedBytesBuilder) *FlowEntryBuilder {
	return &FlowEntryBuilder{
		fields: fields,
	}
}

func (eeb FlowEntryBuilder) set(ed *fwdpb.EntryDesc) {
	flow := &fwdpb.FlowEntryDesc{}
	for _, b := range eeb.fields {
		flow.Fields = append(flow.Fields, b.Build())
	}

	ed.Entry = &fwdpb.EntryDesc_Flow{
		Flow: flow,
	}
}
