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
	actions   [][]ActionDescBuilder
}

// TableEntryAddRequest creates a new TableEntryAddRequestBuilder.
func TableEntryAddRequest(ctxID, tableID string) *TableEntryAddRequestBuilder {
	return &TableEntryAddRequestBuilder{
		contextID: ctxID,
		tableID:   tableID,
	}
}

// AppendEntry adds an entry and the actions to the requests.
func (b *TableEntryAddRequestBuilder) AppendEntry(entry *EntryDescBuilder, actions ...ActionDescBuilder) *TableEntryAddRequestBuilder {
	b.entries = append(b.entries, entry)
	if len(actions) != 0 {
		b.actions = append(b.actions, actions)
	}
	return b
}

// AppendEntry adds the actions to the requests.
func (b *TableEntryAddRequestBuilder) AppendActions(actions ...ActionDescBuilder) *TableEntryAddRequestBuilder {
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
		tableEntry := &fwdpb.TableEntryAddRequest_Entry{
			EntryDesc: entry.Build(),
		}
		req.Entries = append(req.Entries, tableEntry)
		if i >= len(b.actions) {
			break
		}
		for _, act := range b.actions[i] {
			d := &fwdpb.ActionDesc{
				ActionType: act.actionType(),
			}
			act.set(d)
			tableEntry.Actions = append(tableEntry.Actions, d)
		}
	}
	return req
}

// TableEntryAddRequest builds TableEntryAddRequests.
type TableEntryRemoveRequestBuilder struct {
	contextID string
	tableID   string
	entries   []*EntryDescBuilder
}

// TableEntryRemoveRequest creates a new TableEntryRemoveRequestBuilder.
func TableEntryRemoveRequest(ctxID, tableID string) *TableEntryRemoveRequestBuilder {
	return &TableEntryRemoveRequestBuilder{
		contextID: ctxID,
		tableID:   tableID,
	}
}

// AppendEntry adds an entry to the requests.
func (b *TableEntryRemoveRequestBuilder) AppendEntry(entry *EntryDescBuilder) *TableEntryRemoveRequestBuilder {
	b.entries = append(b.entries, entry)
	return b
}

// Build returns a new TableEntryAddRequest.
func (b TableEntryRemoveRequestBuilder) Build() *fwdpb.TableEntryRemoveRequest {
	req := &fwdpb.TableEntryRemoveRequest{
		ContextId: &fwdpb.ContextId{Id: b.contextID},
		TableId:   &fwdpb.TableId{ObjectId: &fwdpb.ObjectId{Id: b.tableID}},
		Entries:   []*fwdpb.EntryDesc{},
	}
	for _, entry := range b.entries {
		req.Entries = append(req.Entries, entry.Build())
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

// ExactEntry creates a new exact entry builder.
func ExactEntry(fields ...*PacketFieldBytesBuilder) *ExactEntryBuilder {
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

// PrefixEntryBuilder builds prefix table entries.
type PrefixEntryBuilder struct {
	fields []*PacketFieldMaskedBytesBuilder
}

// PrefixEntry creates a new prefix entry builder.
func PrefixEntry(fields ...*PacketFieldMaskedBytesBuilder) *PrefixEntryBuilder {
	return &PrefixEntryBuilder{
		fields: fields,
	}
}

func (eeb PrefixEntryBuilder) set(ed *fwdpb.EntryDesc) {
	prefix := &fwdpb.PrefixEntryDesc{}
	for _, b := range eeb.fields {
		prefix.Fields = append(prefix.Fields, b.Build())
	}

	ed.Entry = &fwdpb.EntryDesc_Prefix{
		Prefix: prefix,
	}
}

// ActionEntryBuilder builds action table entries.
type ActionEntryBuilder struct {
	id           string
	insertMethod fwdpb.ActionEntryDesc_InsertMethod
}

// ActionEntry creates a new action entry builder.
func ActionEntry(id string, insertMethod fwdpb.ActionEntryDesc_InsertMethod) *ActionEntryBuilder {
	return &ActionEntryBuilder{
		id:           id,
		insertMethod: insertMethod,
	}
}

func (b ActionEntryBuilder) set(ed *fwdpb.EntryDesc) {
	action := &fwdpb.ActionEntryDesc{
		Id:           b.id,
		InsertMethod: b.insertMethod,
	}

	ed.Entry = &fwdpb.EntryDesc_Action{
		Action: action,
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
