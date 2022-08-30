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

// Package forwarding manages forwarding contexts and forwarding objects.
package forwarding

import (
	"errors"
	"fmt"
	"sync"

	log "github.com/golang/glog"

	"github.com/openconfig/lemming/dataplane/forwarding/fwdaction"
	"github.com/openconfig/lemming/dataplane/forwarding/fwdport"
	"github.com/openconfig/lemming/dataplane/forwarding/fwdtable"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/deadlock"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdattribute"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdcontext"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdflowcounter"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdset"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"

	// The packages below are required to use fwd package. As all these are
	// essential, and there can be more to come, importing here is more
	// beneficial.
	_ "github.com/openconfig/lemming/dataplane/forwarding/fwdaction/actions"
	_ "github.com/openconfig/lemming/dataplane/forwarding/fwdport/ports"
	_ "github.com/openconfig/lemming/dataplane/forwarding/fwdtable/bridge"
	_ "github.com/openconfig/lemming/dataplane/forwarding/fwdtable/exact"
	_ "github.com/openconfig/lemming/dataplane/forwarding/fwdtable/flow"
	_ "github.com/openconfig/lemming/dataplane/forwarding/fwdtable/prefix"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/arp"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/ethernet"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/icmp"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/ip"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/metadata"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/opaque"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/tcp"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/udp"
)

// A Engine is an instance of the forwarding engine. It contains a set of
// forwarding contexts, each of which contain forwarding objects such as
// tables, ports and actions.
type Engine struct {
	mu  sync.Mutex
	ctx map[string]*fwdcontext.Context // forwarding contexts indexed by name

	name string    // name of the forwarding engine
	info *InfoList // list of info elements that can be queried
}

// New creates a new forwarding instance using the specified name.
func New(name string) *Engine {
	return &Engine{
		name: name,
		ctx:  make(map[string]*fwdcontext.Context),
		info: NewInfoList(),
	}
}

// UpdateNotification updates the notification service for a context. If the
// notification is set to nil, no notifications are generated for the context.
// The address is the address of the notification service (used in queries)
// in the host:port format.
func (f *Engine) UpdateNotification(contextID *fwdpb.ContextId, notification fwdcontext.NotificationCallback, address string) (err error) {
	defer func() {
		if err != nil {
			log.Errorf("%v: %v, UpdateNotification failed", f.name, err)
		}
	}()

	if contextID == nil {
		return errors.New("fwd: UpdateNotification failed, No context ID")
	}

	ctx, err := f.FindContext(contextID)
	if err != nil {
		return fmt.Errorf("fwd: UpdateNotification failed, err %v", err)
	}

	ctx.Lock()
	ctx.SetNotification(notification, address)
	ctx.Unlock()
	return nil
}

// UpdatePacketSink updates the packet sink service for a context. If the
// service is set to nil, no packets are delivered externally for the context.
// The address is the address of the packet service (used in queries)
// in the host:port format.
func (f *Engine) UpdatePacketSink(contextID *fwdpb.ContextId, packet fwdcontext.PacketCallback, address string) (err error) {
	defer func() {
		if err != nil {
			log.Errorf("%v: %v, UpdatePacketSink failed", f.name, err)
		}
	}()

	if contextID == nil {
		return errors.New("fwd: UpdatePacketSink failed, No context ID")
	}

	ctx, err := f.FindContext(contextID)
	if err != nil {
		return fmt.Errorf("fwd: UpdatePacketSink failed, err %v", err)
	}

	ctx.Lock()
	defer ctx.Unlock()
	ctx.SetPacketSink(packet, address)
	return nil
}

// ContextCreate creates a new context.
func (f *Engine) ContextCreate(contextID *fwdpb.ContextId) (err error) {
	defer func() {
		if err != nil {
			log.Errorf("%v: %v, ContextCreate failed", f.name, err)
		}
	}()

	if contextID == nil {
		return errors.New("No context ID")
	}

	if _, err := f.FindContext(contextID); err == nil {
		return fmt.Errorf("fwd: ContextCreate failed, %v already exists", contextID)
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	id := contextID.GetId()
	ctx := fwdcontext.New(id, f.name)
	f.ctx[id] = ctx
	f.info.AddContext(ctx)
	return nil
}

// ContextDelete deletes a context if it exists.
func (f *Engine) ContextDelete(contextID *fwdpb.ContextId) (err error) {
	defer func() {
		if err != nil {
			log.Errorf("%v: %v, ContextDelete failed", f.name, err)
		}
	}()

	if contextID == nil {
		return errors.New("No context ID")
	}

	ctx, err := f.FindContext(contextID)
	if err != nil {
		return err
	}

	id := contextID.GetId()
	f.mu.Lock()
	delete(f.ctx, id)
	f.info.RemoveContext(ctx)
	f.mu.Unlock()

	ch := make(chan bool)
	go func() {
		ctx.Lock()
		defer ctx.Unlock()

		isPort := func(oid *fwdpb.ObjectId) bool {
			obj, err := ctx.Objects.FindID(oid)
			if err != nil {
				return false
			}
			_, ok := obj.(fwdport.Port)
			return ok
		}

		// Clean up remaining context objects.
		ctx.Cleanup(ch, isPort)
	}()
	<-ch
	return nil
}

// ContextList lists all the contexts in the system.
func (f *Engine) ContextList(request *fwdpb.ContextListRequest, reply *fwdpb.ContextListReply) (err error) {
	defer func() {
		if err != nil {
			log.Errorf("%v: %v, request %v", f.name, err, request)
		}
	}()

	timer := deadlock.NewTimer(deadlock.Timeout, fmt.Sprintf("Processing %+v", request))
	defer timer.Stop()

	f.mu.Lock()
	defer f.mu.Unlock()

	for _, ctx := range f.ctx {
		reply.Contexts = append(reply.Contexts, &fwdpb.ContextAttr{
			ContextId: &fwdpb.ContextId{
				Id: string(ctx.ID),
			},
			PacketAddress:       ctx.PacketAddress,
			NotificationAddress: ctx.NotificationAddress,
		})
	}
	return nil
}

// FindContext finds the specified context.
func (f *Engine) FindContext(contextID *fwdpb.ContextId) (*fwdcontext.Context, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if contextID == nil {
		return nil, errors.New("No context ID")
	}

	id := contextID.GetId()
	ctx, ok := f.ctx[id]
	if !ok {
		return nil, fmt.Errorf("Unknown context %v", id)
	}
	return ctx, nil
}

// ObjectDelete deletes an object.
func (f *Engine) ObjectDelete(request *fwdpb.ObjectDeleteRequest, _ *fwdpb.ObjectDeleteReply) (err error) {
	defer func() {
		if err != nil {
			log.Errorf("%v: %v, request %v", f.name, err, request)
		}
	}()

	timer := deadlock.NewTimer(deadlock.Timeout, fmt.Sprintf("Processing %+v", request))
	defer timer.Stop()

	ctx, err := f.FindContext(request.GetContextId())
	if err != nil {
		return fmt.Errorf("fwd: ObjectDelete failed, err %v", err)
	}

	ctx.Lock()
	defer ctx.Unlock()

	f.info.RemoveObject(ctx, request.GetObjectId())
	if err := ctx.Objects.Remove(request.GetObjectId(), false /*forceCleanup*/); err != nil {
		return fmt.Errorf("fwd: ObjectDelete failed, err %v", err)
	}
	return nil
}

// ObjectList lists all the objects in the system.
func (f *Engine) ObjectList(request *fwdpb.ObjectListRequest, reply *fwdpb.ObjectListReply) (err error) {
	defer func() {
		if err != nil {
			log.Errorf("%v: %v, request %v", f.name, err, request)
		}
	}()

	timer := deadlock.NewTimer(deadlock.Timeout, fmt.Sprintf("Processing %+v", request))
	defer timer.Stop()

	ctx, err := f.FindContext(request.GetContextId())
	if err != nil {
		return fmt.Errorf("fwd: ObjectList failed, err %v", err)
	}

	ctx.Lock()
	defer ctx.Unlock()

	for _, id := range ctx.Objects.IDs() {
		reply.Objects = append(reply.Objects, &fwdpb.ObjectId{
			Id: string(id),
		})
	}
	return nil
}

// ObjectCounters retrieves all the counters associated on the object.
func (f *Engine) ObjectCounters(request *fwdpb.ObjectCountersRequest, reply *fwdpb.ObjectCountersReply) (err error) {
	defer func() {
		if err != nil {
			log.Errorf("%v: %v, request %v", f.name, err, request)
		}
	}()

	timer := deadlock.NewTimer(deadlock.Timeout, fmt.Sprintf("Processing %+v", request))
	defer timer.Stop()

	ctx, err := f.FindContext(request.GetContextId())
	if err != nil {
		return fmt.Errorf("fwd: ObjectCounters failed, err %v", err)
	}

	ctx.RLock()
	defer ctx.RUnlock()

	object, err := ctx.Objects.FindID(request.GetObjectId())
	if err != nil {
		return fmt.Errorf("fwd: ObjectCounters failed, err %v", err)
	}
	for _, counter := range object.Counters() {
		counter := counter
		reply.Counters = append(reply.Counters, &fwdpb.Counter{
			Id:    counter.ID,
			Value: counter.Value,
		})
	}
	return nil
}

// AttributeList lists all attributes.
func (*Engine) AttributeList(request *fwdpb.AttributeListRequest, reply *fwdpb.AttributeListReply) error {
	for id, help := range fwdattribute.List {
		reply.Attrs = append(reply.Attrs, &fwdpb.AttributeDesc{
			Name: string(id),
			Help: help,
		})
	}
	return nil
}

// AttributeUpdate updates attributes in a forwarding context.
func (f *Engine) AttributeUpdate(request *fwdpb.AttributeUpdateRequest, _ *fwdpb.AttributeUpdateReply) (err error) {
	defer func() {
		if err != nil {
			log.Errorf("%v: %v, request %v", f.name, err, request)
		}
	}()

	timer := deadlock.NewTimer(deadlock.Timeout, fmt.Sprintf("Processing %+v", request))
	defer timer.Stop()

	// Find the attributes
	attributes, err := func() (fwdattribute.Set, error) {
		cid := request.GetContextId()
		if cid == nil {
			return fwdattribute.Global, nil
		}

		ctx, err := f.FindContext(request.GetContextId())
		if err != nil {
			return nil, fmt.Errorf("fwd: AttributeUpdate failed, err %v", err)
		}

		ctx.Lock()
		defer ctx.Unlock()

		if request.ObjectId == nil {
			return ctx.Attributes, nil
		}

		object, err := ctx.Objects.FindID(request.GetObjectId())
		if err != nil {
			return nil, fmt.Errorf("fwd: AttributeUpdate failed, err %v", err)
		}
		return object.Attributes(), nil
	}()
	if err != nil {
		return err
	}
	if attributes == nil {
		return errors.New("fwd: AttributeUpdate failed, no attribute on object")
	}

	// Set a value if it is specified, else unset it.
	if request.AttrId == "" {
		return errors.New("fwd: AttributeUpdate failed, no attribute specified")
	}
	if request.AttrValue != "" {
		attributes.Add(fwdattribute.ID(request.AttrId), request.AttrValue)
	} else {
		attributes.Delete(fwdattribute.ID(request.AttrId))
	}
	return nil
}

// PortCreate creates a port.
func (f *Engine) PortCreate(request *fwdpb.PortCreateRequest, reply *fwdpb.PortCreateReply) (err error) {
	defer func() {
		if err != nil {
			log.Errorf("%v: %v, request %v", f.name, err, request)
		}
	}()

	timer := deadlock.NewTimer(deadlock.Timeout, fmt.Sprintf("Processing %+v", request))
	defer timer.Stop()

	ctx, err := f.FindContext(request.GetContextId())
	if err != nil {
		return fmt.Errorf("fwd: PortCreate failed, err %v", err)
	}

	ctx.Lock()
	defer ctx.Unlock()

	object, err := fwdport.New(request.GetPort(), ctx)
	if err != nil {
		return fmt.Errorf("fwd: PortCreate failed, err %v", err)
	}
	f.info.AddObject(ctx, object)
	reply.ObjectIndex = &fwdpb.ObjectIndex{
		Index: uint64(object.NID()),
	}
	return nil
}

// PortUpdate updates a port.
func (f *Engine) PortUpdate(request *fwdpb.PortUpdateRequest, _ *fwdpb.PortUpdateReply) (err error) {
	defer func() {
		if err != nil {
			log.Errorf("%v: %v request %v", f.name, err, request)
		}
	}()

	timer := deadlock.NewTimer(deadlock.Timeout, fmt.Sprintf("Processing %+v", request))
	defer timer.Stop()

	ctx, err := f.FindContext(request.GetContextId())
	if err != nil {
		return fmt.Errorf("fwd: PortUpdate failed, err %v", err)
	}

	ctx.Lock()
	defer ctx.Unlock()

	port, err := fwdport.Find(request.GetPortId(), ctx)
	if err != nil {
		return fmt.Errorf("fwd: PortUpdate failed, err %v", err)
	}
	upd := request.GetUpdate()
	if upd == nil {
		return errors.New("fwd: PortUpdate failed, no update")
	}
	if err = port.Update(request.GetUpdate()); err != nil {
		return fmt.Errorf("fwd: PortUpdate failed, err %v", err)
	}
	return nil
}

// PortState controls and queries the port state.
func (f *Engine) PortState(request *fwdpb.PortStateRequest, reply *fwdpb.PortStateReply) (err error) {
	defer func() {
		if err != nil {
			log.Errorf("%v: %v, request %v", f.name, err, request)
		}
	}()

	timer := deadlock.NewTimer(deadlock.Timeout, fmt.Sprintf("Processing %+v", request))
	defer timer.Stop()

	ctx, err := f.FindContext(request.GetContextId())
	if err != nil {
		return fmt.Errorf("fwd: PortState failed, err %v", err)
	}

	// If the request has no specified "operation", it is a request which can
	// be satisified with a read lock.
	if request.Operation == nil {
		ctx.RLock()
		defer ctx.RUnlock()
	} else {
		ctx.Lock()
		defer ctx.Unlock()
	}

	port, err := fwdport.Find(request.GetPortId(), ctx)
	if err != nil {
		return fmt.Errorf("fwd: PortState failed, err %v", err)
	}
	reply, err = port.State(request.Operation)
	if err != nil {
		return fmt.Errorf("fwd: PortState failed, err %v", err)
	}
	return nil
}

// TableCreate creates an empty table.
func (f *Engine) TableCreate(request *fwdpb.TableCreateRequest, reply *fwdpb.TableCreateReply) (err error) {
	defer func() {
		if err != nil {
			log.Errorf("%v: %v, request %v", f.name, err, request)
		}
	}()

	timer := deadlock.NewTimer(deadlock.Timeout, fmt.Sprintf("Processing %+v", request))
	defer timer.Stop()

	ctx, err := f.FindContext(request.GetContextId())
	if err != nil {
		return fmt.Errorf("fwd: TableCreate failed, err %v", err)
	}

	ctx.Lock()
	defer ctx.Unlock()

	object, err := fwdtable.New(ctx, request.GetDesc())
	if err != nil {
		return fmt.Errorf("fwd: TableCreate failed, err %v", err)
	}
	f.info.AddObject(ctx, object)
	reply.ObjectIndex = &fwdpb.ObjectIndex{
		Index: uint64(object.NID()),
	}
	return nil
}

// TableEntryAdd adds (or updates) an entry in a table.
func (f *Engine) TableEntryAdd(request *fwdpb.TableEntryAddRequest, _ *fwdpb.TableEntryAddReply) (err error) {
	defer func() {
		if err != nil {
			log.Errorf("%v: %v, request %v", f.name, err, request)
		}
	}()

	timer := deadlock.NewTimer(deadlock.Timeout, fmt.Sprintf("Processing %+v", request))
	defer timer.Stop()

	ctx, err := f.FindContext(request.GetContextId())
	if err != nil {
		return fmt.Errorf("fwd: TableEntryAdd failed, err %v", err)
	}

	ctx.Lock()
	defer ctx.Unlock()

	table, err := fwdtable.Find(ctx, request.GetTableId())
	if err != nil {
		return fmt.Errorf("fwd: TableEntryAdd failed, err %v", err)
	}

	// If requested, clear the table before adding entries.
	if request.GetClearBeforeAdd() {
		table.Clear()
	}

	add := func(desc *fwdpb.EntryDesc, actions []*fwdpb.ActionDesc) error {
		if err := table.AddEntry(desc, actions); err != nil {
			return fmt.Errorf("fwd: TableEntryAdd failed, err %v", err)
		}
		return nil
	}

	// Add the singleton entry if it exists.
	if request.EntryDesc != nil {
		if err := add(request.EntryDesc, request.Actions); err != nil {
			return err
		}
	}

	// Add the batched entries.
	for _, entry := range request.Entries {
		if err := add(entry.EntryDesc, entry.Actions); err != nil {
			return err
		}
	}
	return nil
}

// TableEntryRemove removes an entry from the table.
func (f *Engine) TableEntryRemove(request *fwdpb.TableEntryRemoveRequest, _ *fwdpb.TableEntryRemoveReply) (err error) {
	defer func() {
		if err != nil {
			log.Errorf("%v: %v, request %v", f.name, err, request)
		}
	}()

	timer := deadlock.NewTimer(deadlock.Timeout, fmt.Sprintf("Processing %+v", request))
	defer timer.Stop()

	ctx, err := f.FindContext(request.GetContextId())
	if err != nil {
		return fmt.Errorf("fwd: TableEntryRemove failed, err %v", err)
	}

	ctx.Lock()
	defer ctx.Unlock()

	table, err := fwdtable.Find(ctx, request.GetTableId())
	if err != nil {
		return fmt.Errorf("fwd: TableEntryRemove failed, err %v", err)
	}

	if desc := request.GetEntryDesc(); desc != nil {
		if err = table.RemoveEntry(desc); err != nil {
			return fmt.Errorf("fwd: TableEntryRemove failed, err %v", err)
		}
	}
	for _, entry := range request.Entries {
		if err = table.RemoveEntry(entry); err != nil {
			return fmt.Errorf("fwd: TableEntryRemove failed, err %v", err)
		}
	}
	return nil
}

// TableList lists all the entries in the specifed table.
func (f *Engine) TableList(request *fwdpb.TableListRequest, reply *fwdpb.TableListReply) (err error) {
	defer func() {
		if err != nil {
			log.Errorf("%v: %v, request %v", f.name, err, request)
		}
	}()

	timer := deadlock.NewTimer(deadlock.Timeout, fmt.Sprintf("Processing %+v", request))
	defer timer.Stop()

	ctx, err := f.FindContext(request.GetContextId())
	if err != nil {
		return fmt.Errorf("fwd: TableEntryList failed, err %v", err)
	}

	ctx.RLock()
	defer ctx.RUnlock()

	table, err := fwdtable.Find(ctx, request.GetTableId())
	if err != nil {
		return fmt.Errorf("fwd: TableEntryList failed, err %v", err)
	}
	reply.Entries = table.Entries()
	return nil
}

// SetCreate creates a new set.
func (f *Engine) SetCreate(request *fwdpb.SetCreateRequest, reply *fwdpb.SetCreateReply) (err error) {
	defer func() {
		if err != nil {
			log.Errorf("%v: %v, request %v", f.name, err, request)
		}
	}()

	timer := deadlock.NewTimer(deadlock.Timeout, fmt.Sprintf("Processing %+v", request))
	defer timer.Stop()

	ctx, err := f.FindContext(request.GetContextId())
	if err != nil {
		return fmt.Errorf("fwd: SetCreate failed, err %v", err)
	}

	ctx.Lock()
	defer ctx.Unlock()

	c, err := fwdset.New(ctx, request.GetSetId())
	if err != nil {
		return fmt.Errorf("fwd: SetCreate failed, err %v", err)
	}
	f.info.AddObject(ctx, c)
	reply.ObjectIndex = &fwdpb.ObjectIndex{
		Index: uint64(c.NID()),
	}
	return nil
}

// SetUpdate updates a set.
func (f *Engine) SetUpdate(request *fwdpb.SetUpdateRequest, _ *fwdpb.SetUpdateReply) (err error) {
	defer func() {
		if err != nil {
			log.Errorf("%v: %v, request %v", f.name, err, request)
		}
	}()

	timer := deadlock.NewTimer(deadlock.Timeout, fmt.Sprintf("Processing %+v", request))
	defer timer.Stop()

	ctx, err := f.FindContext(request.GetContextId())
	if err != nil {
		return fmt.Errorf("fwd: SetUpdate failed, err %v", err)
	}

	ctx.Lock()
	defer ctx.Unlock()

	c, err := fwdset.Find(ctx, request.GetSetId())
	if err != nil {
		return fmt.Errorf("fwd: SetUpdate failed, err %v", err)
	}
	c.Update(request.GetBytes())
	return nil
}

// FlowCounterCreate creates a flow counter with specified ObjectId in the specified context.
func (f *Engine) FlowCounterCreate(request *fwdpb.FlowCounterCreateRequest, reply *fwdpb.FlowCounterCreateReply) (err error) {
	defer func() {
		if err != nil {
			log.Errorf("%v: %v, request %v", f.name, err, request)
		}
	}()

	timer := deadlock.NewTimer(deadlock.Timeout, fmt.Sprintf("Processing %+v", request))
	defer timer.Stop()

	ctx, err := f.FindContext(request.GetContextId())
	if err != nil {
		return fmt.Errorf("fwd: FlowCounterCreate failed, err %v", err)
	}

	ctx.Lock()
	defer ctx.Unlock()

	// Check to see if the FlowCounter object already exists.
	if _, err := ctx.Objects.FindID(request.GetId().GetObjectId()); err == nil {
		return fmt.Errorf("fwd: FlowCounterCreate failed, it already exists: err %v", err)
	}

	fc, err := fwdflowcounter.New(ctx, request)
	if err != nil {
		return fmt.Errorf("fwd: FlowCounterCreate failed, err %v", err)
	}
	f.info.AddObject(ctx, fc)

	// FlowCounterCreateReply is empty.

	return nil
}

// FlowCounterQuery queries for the values of specified counters.
func (f *Engine) FlowCounterQuery(request *fwdpb.FlowCounterQueryRequest, reply *fwdpb.FlowCounterQueryReply) (err error) {
	defer func() {
		if err != nil {
			log.Errorf("%v: %v, request %v", f.name, err, request)
		}
	}()

	timer := deadlock.NewTimer(deadlock.Timeout, fmt.Sprintf("Processing %+v", request))
	defer timer.Stop()

	ctx, err := f.FindContext(request.GetContextId())
	if err != nil {
		return fmt.Errorf("fwd: FlowCounterQuery failed, err %v", err)
	}

	ctx.RLock()
	defer ctx.RUnlock()

	for _, fcid := range request.Ids {
		obj, err := ctx.Objects.FindID(fcid.GetObjectId())
		if err != nil {
			continue
		}
		fc, ok := obj.(*fwdflowcounter.FlowCounter)
		if ok != true {
			continue
		}
		fcval, fcerr := fc.Query()
		if fcerr != nil {
			continue
		}
		reply.Counters = append(reply.Counters, fcval)
	}

	return nil
}

// PacketInject injects a packet in the specified forwarding context and port.
func (f *Engine) PacketInject(request *fwdpb.PacketInjectRequest) (err error) {
	timer := deadlock.NewTimer(deadlock.Timeout, fmt.Sprintf("Processing %+v", request))
	defer timer.Stop()

	ctx, err := f.FindContext(request.GetContextId())
	if err != nil {
		return fmt.Errorf("fwd: PacketInject failed, err %v", err)
	}

	// In a goroutine, acquire a RWLock on the context, create the preprocessing
	// actions and process the packet. The RPC does not wait for the complete
	// packet processing. It returns once it has validated the input parameters.
	status := make(chan error)

	go func() {
		ctx.RLock()
		defer ctx.RUnlock()

		port, err := fwdport.Find(request.GetPortId(), ctx)
		if err != nil {
			status <- fmt.Errorf("fwd: PortInject failed, err %v", err)
			return
		}

		packet, err := fwdpacket.New(request.GetStartHeader(), request.GetBytes())
		if err != nil {
			status <- fmt.Errorf("fwd: PortInject failed, err %v", err)
			return
		}

		pre, err := fwdaction.NewActions(request.GetPreprocesses(), ctx)
		if err != nil {
			status <- fmt.Errorf("fwd: PortInject failed to create preprocessing actions %v, err %v", request.GetPreprocesses(), err)
			return
		}

		// Apply the preprocessing actions on the packet and inject it into the
		// port while holding the context's RLock. After packet processing,
		// cleanup the port and actions. Publish a "no error" on the status
		// channel so that the RPC can return. Note that this also serializes
		// the packets arriving from the CPU.
		status <- nil

		func() {
			defer func() {
				if pre != nil {
					pre.Cleanup()
				}
			}()

			packet.Debug(request.GetDebug())
			if len(pre) != 0 {
				packet.Logf(fwdpacket.LogDesc, fmt.Sprintf("%v: Preprocess %v", ctx.ID, port.ID()))
				state, err := fwdaction.ProcessPacket(packet, pre, port)
				if state != fwdaction.CONTINUE || err != nil {
					log.Errorf("%v: preprocessing failed, state %v, err %v", port.ID(), state, err)
					return
				}
				packet.Logf(fwdpacket.LogDebugFrame, "injecting packet")
			}
			fwdport.Process(port, packet, request.GetAction(), ctx, "Control")
		}()
	}()

	if err = <-status; err != nil {
		return fmt.Errorf("fwd: PacketInject failed, err %v", err)
	}
	return nil
}

// InfoList retrieves a list of all information elements.
func (f *Engine) InfoList(_ *fwdpb.InfoListRequest, reply *fwdpb.InfoListReply) error {
	reply.Names = f.info.List()
	return nil
}

// InfoElement retrieves the contents of a specific information element.
func (f *Engine) InfoElement(request *fwdpb.InfoElementRequest, reply *fwdpb.InfoElementReply) error {
	content, err := f.info.Element(request.GetName(), request.GetType(), request.GetFrame(), request.GetStartHeader())
	if err != nil {
		return err
	}
	reply.Content = content
	return nil
}
