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
	"context"
	"errors"
	"fmt"
	"io"
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
	_ "github.com/openconfig/lemming/dataplane/forwarding/fwdtable/action"
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

// A Server is an instance of the forwarding server. It contains a set of
// forwarding contexts, each of which contain forwarding objects such as
// tables, ports and actions.
type Server struct {
	fwdpb.UnimplementedForwardingServer
	fwdpb.UnimplementedInfoServer

	mu  sync.Mutex
	ctx map[string]*fwdcontext.Context // forwarding contexts indexed by name

	name string    // name of the forwarding engine
	info *InfoList // list of info elements that can be queried
}

// New creates a new forwarding instance using the specified name.
func New(name string) *Server {
	return &Server{
		name: name,
		ctx:  make(map[string]*fwdcontext.Context),
		info: NewInfoList(),
	}
}

// UpdateNotification updates the notification service for a context. If the
// notification is set to nil, no notifications are generated for the context.
// The address is the address of the notification service (used in queries)
// in the host:port format.
func (e *Server) UpdateNotification(contextID *fwdpb.ContextId, notification fwdcontext.NotificationCallback) error {
	if contextID == nil {
		return errors.New("fwd: UpdateNotification failed, No context ID")
	}

	ctx, err := e.FindContext(contextID)
	if err != nil {
		return fmt.Errorf("fwd: UpdateNotification failed, err %v", err)
	}

	ctx.Lock()
	ctx.SetNotification(notification)
	ctx.Unlock()
	return nil
}

// ContextCreate creates a new context. Note that if the packet sink and/or
// notification services are specified but not reachable, the context creation
// fails.
func (e *Server) ContextCreate(_ context.Context, request *fwdpb.ContextCreateRequest) (*fwdpb.ContextCreateReply, error) {
	cid := request.GetContextId()
	if err := e.contextCreateByID(cid); err != nil {
		return nil, err
	}
	return &fwdpb.ContextCreateReply{}, nil
}

// contextCreateByID creates a new context with the given ID, erroring if it already exists.
func (e *Server) contextCreateByID(contextID *fwdpb.ContextId) error {
	if contextID == nil {
		return errors.New("No context ID")
	}

	if _, err := e.FindContext(contextID); err == nil {
		return fmt.Errorf("fwd: ContextCreate failed, %v already exists", contextID)
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	id := contextID.GetId()
	ctx := fwdcontext.New(id, e.name)
	e.ctx[id] = ctx
	e.info.AddContext(ctx)
	return nil
}

// ContextDelete deletes a context if it exists.
func (e *Server) ContextDelete(_ context.Context, request *fwdpb.ContextDeleteRequest) (*fwdpb.ContextDeleteReply, error) {
	if request.GetContextId() == nil {
		return nil, errors.New("No context ID")
	}

	ctx, err := e.FindContext(request.GetContextId())
	if err != nil {
		return nil, err
	}

	id := request.GetContextId().GetId()
	e.mu.Lock()
	delete(e.ctx, id)
	e.info.RemoveContext(ctx)
	e.mu.Unlock()

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
	return &fwdpb.ContextDeleteReply{}, nil
}

// ContextList lists all the contexts in the system.
func (e *Server) ContextList(_ context.Context, request *fwdpb.ContextListRequest) (*fwdpb.ContextListReply, error) {
	timer := deadlock.NewTimer(deadlock.Timeout, fmt.Sprintf("Processing %+v", request))
	defer timer.Stop()

	e.mu.Lock()
	defer e.mu.Unlock()

	reply := &fwdpb.ContextListReply{}
	for _, ctx := range e.ctx {
		reply.Contexts = append(reply.Contexts, &fwdpb.ContextAttr{
			ContextId: &fwdpb.ContextId{
				Id: ctx.ID,
			},
		})
	}
	return reply, nil
}

// FindContext finds the specified context.
func (e *Server) FindContext(contextID *fwdpb.ContextId) (*fwdcontext.Context, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if contextID == nil {
		return nil, errors.New("No context ID")
	}

	id := contextID.GetId()
	ctx, ok := e.ctx[id]
	if !ok {
		return nil, fmt.Errorf("Unknown context %v", id)
	}
	return ctx, nil
}

// ObjectDelete deletes an object.
func (e *Server) ObjectDelete(_ context.Context, request *fwdpb.ObjectDeleteRequest) (*fwdpb.ObjectDeleteReply, error) {
	timer := deadlock.NewTimer(deadlock.Timeout, fmt.Sprintf("Processing %+v", request))
	defer timer.Stop()

	ctx, err := e.FindContext(request.GetContextId())
	if err != nil {
		return nil, fmt.Errorf("fwd: ObjectDelete failed, err %v", err)
	}

	ctx.Lock()
	defer ctx.Unlock()

	e.info.RemoveObject(ctx, request.GetObjectId())
	if err := ctx.Objects.Remove(request.GetObjectId(), false /*forceCleanup*/); err != nil {
		return nil, fmt.Errorf("fwd: ObjectDelete failed, err %v", err)
	}
	return &fwdpb.ObjectDeleteReply{}, nil
}

// ObjectList lists all the objects in the system.
func (e *Server) ObjectList(_ context.Context, request *fwdpb.ObjectListRequest) (*fwdpb.ObjectListReply, error) {
	timer := deadlock.NewTimer(deadlock.Timeout, fmt.Sprintf("Processing %+v", request))
	defer timer.Stop()

	ctx, err := e.FindContext(request.GetContextId())
	if err != nil {
		return nil, fmt.Errorf("fwd: ObjectList failed, err %v", err)
	}

	ctx.Lock()
	defer ctx.Unlock()

	reply := &fwdpb.ObjectListReply{}
	for _, id := range ctx.Objects.IDs() {
		reply.Objects = append(reply.Objects, &fwdpb.ObjectId{
			Id: string(id),
		})
	}
	return reply, nil
}

// ObjectCounters retrieves all the counters associated on the object.
func (e *Server) ObjectCounters(_ context.Context, request *fwdpb.ObjectCountersRequest) (*fwdpb.ObjectCountersReply, error) {
	timer := deadlock.NewTimer(deadlock.Timeout, fmt.Sprintf("Processing %+v", request))
	defer timer.Stop()

	ctx, err := e.FindContext(request.GetContextId())
	if err != nil {
		return nil, fmt.Errorf("fwd: ObjectCounters failed, err %v", err)
	}

	ctx.RLock()
	defer ctx.RUnlock()

	object, err := ctx.Objects.FindID(request.GetObjectId())
	if err != nil {
		return nil, fmt.Errorf("fwd: ObjectCounters failed, err %v", err)
	}
	reply := &fwdpb.ObjectCountersReply{}
	for _, counter := range object.Counters() {
		counter := counter
		reply.Counters = append(reply.Counters, &fwdpb.Counter{
			Id:    counter.ID,
			Value: counter.Value,
		})
	}
	return reply, nil
}

// AttributeList lists all attributes.
func (*Server) AttributeList(context.Context, *fwdpb.AttributeListRequest) (*fwdpb.AttributeListReply, error) {
	reply := &fwdpb.AttributeListReply{}
	for id, help := range fwdattribute.List {
		reply.Attrs = append(reply.Attrs, &fwdpb.AttributeDesc{
			Name: string(id),
			Help: help,
		})
	}
	return reply, nil
}

// AttributeUpdate updates attributes in a forwarding context.
func (e *Server) AttributeUpdate(_ context.Context, request *fwdpb.AttributeUpdateRequest) (*fwdpb.AttributeUpdateReply, error) {
	timer := deadlock.NewTimer(deadlock.Timeout, fmt.Sprintf("Processing %+v", request))
	defer timer.Stop()

	// Find the attributes
	attributes, err := func() (fwdattribute.Set, error) {
		cid := request.GetContextId()
		if cid == nil {
			return fwdattribute.Global, nil
		}

		ctx, err := e.FindContext(request.GetContextId())
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
		return nil, err
	}
	if attributes == nil {
		return nil, errors.New("fwd: AttributeUpdate failed, no attribute on object")
	}

	// Set a value if it is specified, else unset it.
	if request.AttrId == "" {
		return nil, errors.New("fwd: AttributeUpdate failed, no attribute specified")
	}
	if request.AttrValue != "" {
		attributes.Add(fwdattribute.ID(request.AttrId), request.AttrValue)
	} else {
		attributes.Delete(fwdattribute.ID(request.AttrId))
	}
	return &fwdpb.AttributeUpdateReply{}, nil
}

// AttributeQuery queries attributes in a forwarding context.
func (e *Server) AttributeQuery(_ context.Context, request *fwdpb.AttributeQueryRequest) (*fwdpb.AttributeQueryReply, error) {
	timer := deadlock.NewTimer(deadlock.Timeout, fmt.Sprintf("Processing %+v", request))
	defer timer.Stop()

	// Find the attributes
	attributes, err := func() (fwdattribute.Set, error) {
		cid := request.GetContextId()
		if cid == nil {
			return fwdattribute.Global, nil
		}

		ctx, err := e.FindContext(request.GetContextId())
		if err != nil {
			return nil, fmt.Errorf("fwd: AttributeQuery failed, err %v", err)
		}

		ctx.Lock()
		defer ctx.Unlock()

		if request.ObjectId == nil {
			return ctx.Attributes, nil
		}

		object, err := ctx.Objects.FindID(request.GetObjectId())
		if err != nil {
			return nil, fmt.Errorf("fwd: AttributeQuery failed, err %v", err)
		}
		return object.Attributes(), nil
	}()
	if err != nil {
		return nil, err
	}
	if attributes == nil {
		return nil, errors.New("fwd: AttributeQuery failed, no attribute on object")
	}

	// Set a value if it is specified, else unset it.
	if request.AttrId == "" {
		return nil, errors.New("fwd: AttributeQuery failed, no attribute specified")
	}
	val, ok := attributes.Get(fwdattribute.ID(request.AttrId))
	if !ok {
		return nil, errors.New("fwd: AttributeQuery failed, attribute not set")
	}
	return &fwdpb.AttributeQueryReply{
		AttrValue: val,
	}, nil
}

// PortCreate creates a port.
func (e *Server) PortCreate(_ context.Context, request *fwdpb.PortCreateRequest) (*fwdpb.PortCreateReply, error) {
	timer := deadlock.NewTimer(deadlock.Timeout, fmt.Sprintf("Processing %+v", request))
	defer timer.Stop()

	ctx, err := e.FindContext(request.GetContextId())
	if err != nil {
		return nil, fmt.Errorf("fwd: PortCreate failed, err %v", err)
	}

	ctx.Lock()
	defer ctx.Unlock()

	object, err := fwdport.New(request.GetPort(), ctx)
	if err != nil {
		return nil, fmt.Errorf("fwd: PortCreate failed, err %v", err)
	}
	e.info.AddObject(ctx, object)
	reply := &fwdpb.PortCreateReply{
		ObjectIndex: &fwdpb.ObjectIndex{
			Index: uint64(object.NID()),
		},
	}
	// Notify with the current state of the port.
	resp, err := object.State(nil)
	if err != nil {
		log.Warningf("skipping initial state notification due to error: %v", err)
		return reply, nil
	}
	err = ctx.Notify(&fwdpb.EventDesc{
		Event: fwdpb.Event_EVENT_PORT,
		Desc: &fwdpb.EventDesc_Port{
			Port: &fwdpb.PortEventDesc{
				Context:  request.GetContextId(),
				PortId:   request.GetPort().GetPortId(),
				PortInfo: resp.GetStatus(),
			},
		},
	})
	if err != nil {
		log.Warningf("error sending initial state notification: %v", err)
	}

	return reply, nil
}

// PortUpdate updates a port.
func (e *Server) PortUpdate(_ context.Context, request *fwdpb.PortUpdateRequest) (*fwdpb.PortUpdateReply, error) {
	timer := deadlock.NewTimer(deadlock.Timeout, fmt.Sprintf("Processing %+v", request))
	defer timer.Stop()

	ctx, err := e.FindContext(request.GetContextId())
	if err != nil {
		return nil, fmt.Errorf("fwd: PortUpdate failed, err %v", err)
	}

	ctx.Lock()
	defer ctx.Unlock()

	port, err := fwdport.Find(request.GetPortId(), ctx)
	if err != nil {
		return nil, fmt.Errorf("fwd: PortUpdate failed, err %v", err)
	}
	upd := request.GetUpdate()
	if upd == nil {
		return nil, errors.New("fwd: PortUpdate failed, no update")
	}
	if err = port.Update(request.GetUpdate()); err != nil {
		return nil, fmt.Errorf("fwd: PortUpdate failed, err %v", err)
	}
	return &fwdpb.PortUpdateReply{}, nil
}

// PortState controls and queries the port state.
func (e *Server) PortState(_ context.Context, request *fwdpb.PortStateRequest) (*fwdpb.PortStateReply, error) {
	timer := deadlock.NewTimer(deadlock.Timeout, fmt.Sprintf("Processing %+v", request))
	defer timer.Stop()

	ctx, err := e.FindContext(request.GetContextId())
	if err != nil {
		return nil, fmt.Errorf("fwd: PortState failed, err %v", err)
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
		return nil, fmt.Errorf("fwd: PortState failed, err %v", err)
	}
	reply, err := port.State(request.Operation)
	if err != nil {
		return nil, fmt.Errorf("fwd: PortState failed, err %v", err)
	}
	return reply, nil
}

// TableCreate creates an empty table.
func (e *Server) TableCreate(_ context.Context, request *fwdpb.TableCreateRequest) (*fwdpb.TableCreateReply, error) {
	timer := deadlock.NewTimer(deadlock.Timeout, fmt.Sprintf("Processing %+v", request))
	defer timer.Stop()

	ctx, err := e.FindContext(request.GetContextId())
	if err != nil {
		return nil, fmt.Errorf("fwd: TableCreate failed, err %v", err)
	}

	ctx.Lock()
	defer ctx.Unlock()

	object, err := fwdtable.New(ctx, request.GetDesc())
	if err != nil {
		return nil, fmt.Errorf("fwd: TableCreate failed, table %q: err %v", request.GetDesc().GetTableId().GetObjectId().GetId(), err)
	}
	e.info.AddObject(ctx, object)
	reply := &fwdpb.TableCreateReply{
		ObjectIndex: &fwdpb.ObjectIndex{
			Index: uint64(object.NID()),
		},
	}
	return reply, nil
}

// TableEntryAdd adds (or updates) an entry in a table.
func (e *Server) TableEntryAdd(_ context.Context, request *fwdpb.TableEntryAddRequest) (*fwdpb.TableEntryAddReply, error) {
	timer := deadlock.NewTimer(deadlock.Timeout, fmt.Sprintf("Processing %+v", request))
	defer timer.Stop()

	ctx, err := e.FindContext(request.GetContextId())
	if err != nil {
		return nil, fmt.Errorf("fwd: TableEntryAdd failed, err %v", err)
	}

	ctx.Lock()
	defer ctx.Unlock()

	table, err := fwdtable.Find(ctx, request.GetTableId())
	if err != nil {
		return nil, fmt.Errorf("fwd: TableEntryAdd failed, err %v", err)
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
			return nil, err
		}
	}

	// Add the batched entries.
	for _, entry := range request.Entries {
		if err := add(entry.EntryDesc, entry.Actions); err != nil {
			return nil, err
		}
	}
	return &fwdpb.TableEntryAddReply{}, nil
}

// TableEntryRemove removes an entry from the table.
func (e *Server) TableEntryRemove(_ context.Context, request *fwdpb.TableEntryRemoveRequest) (*fwdpb.TableEntryRemoveReply, error) {
	timer := deadlock.NewTimer(deadlock.Timeout, fmt.Sprintf("Processing %+v", request))
	defer timer.Stop()

	ctx, err := e.FindContext(request.GetContextId())
	if err != nil {
		return nil, fmt.Errorf("fwd: TableEntryRemove failed, err %v", err)
	}

	ctx.Lock()
	defer ctx.Unlock()

	table, err := fwdtable.Find(ctx, request.GetTableId())
	if err != nil {
		return nil, fmt.Errorf("fwd: TableEntryRemove failed, err %v", err)
	}

	if desc := request.GetEntryDesc(); desc != nil {
		if err = table.RemoveEntry(desc); err != nil {
			return nil, fmt.Errorf("fwd: TableEntryRemove failed, err %v", err)
		}
	}
	for _, entry := range request.Entries {
		if err = table.RemoveEntry(entry); err != nil {
			return nil, fmt.Errorf("fwd: TableEntryRemove failed, err %v", err)
		}
	}
	return &fwdpb.TableEntryRemoveReply{}, nil
}

// TableList lists all the entries in the specified table.
func (e *Server) TableList(_ context.Context, request *fwdpb.TableListRequest) (*fwdpb.TableListReply, error) {
	timer := deadlock.NewTimer(deadlock.Timeout, fmt.Sprintf("Processing %+v", request))
	defer timer.Stop()

	ctx, err := e.FindContext(request.GetContextId())
	if err != nil {
		return nil, fmt.Errorf("fwd: TableEntryList failed, err %v", err)
	}

	ctx.RLock()
	defer ctx.RUnlock()

	table, err := fwdtable.Find(ctx, request.GetTableId())
	if err != nil {
		return nil, fmt.Errorf("fwd: TableEntryList failed, err %v", err)
	}
	reply := &fwdpb.TableListReply{
		Entries: table.Entries(),
	}
	return reply, nil
}

// SetCreate creates a new set.
func (e *Server) SetCreate(_ context.Context, request *fwdpb.SetCreateRequest) (*fwdpb.SetCreateReply, error) {
	timer := deadlock.NewTimer(deadlock.Timeout, fmt.Sprintf("Processing %+v", request))
	defer timer.Stop()

	ctx, err := e.FindContext(request.GetContextId())
	if err != nil {
		return nil, fmt.Errorf("fwd: SetCreate failed, err %v", err)
	}

	ctx.Lock()
	defer ctx.Unlock()

	c, err := fwdset.New(ctx, request.GetSetId())
	if err != nil {
		return nil, fmt.Errorf("fwd: SetCreate failed, err %v", err)
	}
	e.info.AddObject(ctx, c)
	reply := &fwdpb.SetCreateReply{
		ObjectIndex: &fwdpb.ObjectIndex{
			Index: uint64(c.NID()),
		},
	}

	return reply, nil
}

// SetUpdate updates a set.
func (e *Server) SetUpdate(_ context.Context, request *fwdpb.SetUpdateRequest) (*fwdpb.SetUpdateReply, error) {
	timer := deadlock.NewTimer(deadlock.Timeout, fmt.Sprintf("Processing %+v", request))
	defer timer.Stop()

	ctx, err := e.FindContext(request.GetContextId())
	if err != nil {
		return nil, fmt.Errorf("fwd: SetUpdate failed, err %v", err)
	}

	ctx.Lock()
	defer ctx.Unlock()

	c, err := fwdset.Find(ctx, request.GetSetId())
	if err != nil {
		return nil, fmt.Errorf("fwd: SetUpdate failed, err %v", err)
	}
	c.Update(request.GetBytes())
	return &fwdpb.SetUpdateReply{}, nil
}

// FlowCounterCreate creates a flow counter with specified ObjectId in the specified context.
func (e *Server) FlowCounterCreate(_ context.Context, request *fwdpb.FlowCounterCreateRequest) (*fwdpb.FlowCounterCreateReply, error) {
	timer := deadlock.NewTimer(deadlock.Timeout, fmt.Sprintf("Processing %+v", request))
	defer timer.Stop()

	ctx, err := e.FindContext(request.GetContextId())
	if err != nil {
		return nil, fmt.Errorf("fwd: FlowCounterCreate failed, err %v", err)
	}

	ctx.Lock()
	defer ctx.Unlock()

	// Check to see if the FlowCounter object already exists.
	if _, err := ctx.Objects.FindID(request.GetId().GetObjectId()); err == nil {
		return nil, fmt.Errorf("fwd: FlowCounterCreate failed, it already exists: err %v", err)
	}

	fc, err := fwdflowcounter.New(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("fwd: FlowCounterCreate failed, err %v", err)
	}
	e.info.AddObject(ctx, fc)

	return &fwdpb.FlowCounterCreateReply{}, nil
}

// FlowCounterQuery queries for the values of specified counters.
func (e *Server) FlowCounterQuery(_ context.Context, request *fwdpb.FlowCounterQueryRequest) (*fwdpb.FlowCounterQueryReply, error) {
	timer := deadlock.NewTimer(deadlock.Timeout, fmt.Sprintf("Processing %+v", request))
	defer timer.Stop()

	ctx, err := e.FindContext(request.GetContextId())
	if err != nil {
		return nil, fmt.Errorf("fwd: FlowCounterQuery failed, err %v", err)
	}

	ctx.RLock()
	defer ctx.RUnlock()

	reply := &fwdpb.FlowCounterQueryReply{}
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

	return reply, nil
}

// InjectPacket inserts a packet in the forwarding pipeline, orginating from the specified port.
func (e *Server) InjectPacket(contextID *fwdpb.ContextId, id *fwdpb.PortId, hid fwdpb.PacketHeaderId, frame []byte, preActions []*fwdpb.ActionDesc, debug bool, dir fwdpb.PortAction) error {
	timer := deadlock.NewTimer(deadlock.Timeout, fmt.Sprintf("Processing packet"))
	defer timer.Stop()

	ctx, err := e.FindContext(contextID)
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

		port, err := fwdport.Find(id, ctx)
		if err != nil {
			status <- fmt.Errorf("fwd: PortInject failed, err %v", err)
			return
		}

		packet, err := fwdpacket.New(hid, frame)
		if err != nil {
			status <- fmt.Errorf("fwd: PortInject failed, err %v", err)
			return
		}

		pre, err := fwdaction.NewActions(preActions, ctx)
		if err != nil {
			status <- fmt.Errorf("fwd: PortInject failed to create preprocessing actions %v, err %v", preActions, err)
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

			packet.Debug(debug)
			if len(pre) != 0 {
				packet.Log().WithValues("context", ctx.ID, "port", port.ID())
				state, err := fwdaction.ProcessPacket(packet, pre, port)
				if state != fwdaction.CONTINUE || err != nil {
					log.Errorf("%v: preprocessing failed, state %v, err %v", port.ID(), state, err)
					return
				}
				packet.Log().V(1).Info("injecting packet", "frame", fwdpacket.IncludeFrameInLog)
			}
			fwdport.Process(port, packet, dir, ctx, "Control")
		}()
	}()

	if err = <-status; err != nil {
		return fmt.Errorf("fwd: PacketInject failed, err %v", err)
	}
	return nil
}

// PacketInject is a streaming RPC to inject a packets in the specified forwarding context and port.
func (e *Server) PacketInject(srv fwdpb.Forwarding_PacketInjectServer) error {
	for {
		request, err := srv.Recv()
		if err != nil {
			return err
		}
		err = e.InjectPacket(request.GetContextId(), request.GetPortId(), request.GetStartHeader(),
			request.GetBytes(), request.GetPreprocesses(), request.GetDebug(), request.GetAction())
		if err != nil {
			return err
		}
	}
}

// InfoList retrieves a list of all information elements.
func (e *Server) InfoList(context.Context, *fwdpb.InfoListRequest) (*fwdpb.InfoListReply, error) {
	return &fwdpb.InfoListReply{
		Names: e.info.List(),
	}, nil
}

// InfoElement retrieves the contents of a specific information element.
func (e *Server) InfoElement(_ context.Context, request *fwdpb.InfoElementRequest) (*fwdpb.InfoElementReply, error) {
	return e.info.Element(request.GetName(), request.GetType(), request.GetFrame(), request.GetStartHeader())
}

// Operation processes incoming OperationRequests and extracts the contained Request.
func (e *Server) Operation(stream fwdpb.Forwarding_OperationServer) error {
	for {
		operationRequest, err := stream.Recv()

		// Client is done. Clean up the bindings and return nil.
		if err == io.EOF {
			return nil
		}

		// Return the error.
		if err != nil {
			return err
		}

		var operationReply fwdpb.OperationReply
		switch request := operationRequest.Request.(type) {
		case *fwdpb.OperationRequest_TableEntryAdd:
			reply, err := e.TableEntryAdd(stream.Context(), request.TableEntryAdd)
			if err != nil {
				return err
			}
			operationReply.Reply = &fwdpb.OperationReply_TableEntryAdd{TableEntryAdd: reply}
		case *fwdpb.OperationRequest_TableEntryRemove:
			reply, err := e.TableEntryRemove(stream.Context(), request.TableEntryRemove)
			if err != nil {
				return err
			}
			operationReply.Reply = &fwdpb.OperationReply_TableEntryRemove{TableEntryRemove: reply}
		case nil:
			// The field is not set.
			return fmt.Errorf("OperationRequest.Request was not set")
		default:
			return fmt.Errorf("OperationRequest.Request has unexpected type %T", request)
		}
		if err := stream.Send(&operationReply); err != nil {
			return err
		}
	}
}

// NotifySubscribe subscribes to notification for a forwarding context.
// TODO: Only one notification client is supported per forwarding context.
func (e *Server) NotifySubscribe(sub *fwdpb.NotifySubscribeRequest, srv fwdpb.Forwarding_NotifySubscribeServer) error {
	eventCh := make(chan *fwdpb.EventDesc)
	fn := func(ed *fwdpb.EventDesc) {
		eventCh <- ed
	}

	if err := e.UpdateNotification(sub.GetContext(), fn); err != nil {
		return err
	}

	for {
		e := <-eventCh
		if err := srv.Send(e); err != nil {
			return err
		}
	}
}

// ObjectNID returns the numeric id for the object id if it exists.
func (e *Server) ObjectNID(_ context.Context, request *fwdpb.ObjectNIDRequest) (*fwdpb.ObjectNIDReply, error) {
	timer := deadlock.NewTimer(deadlock.Timeout, fmt.Sprintf("Processing %+v", request))
	defer timer.Stop()

	ctx, err := e.FindContext(request.GetContextId())
	if err != nil {
		return nil, fmt.Errorf("fwd: ObjectNID failed, err %v", err)
	}

	ctx.Lock()
	defer ctx.Unlock()
	obj, err := ctx.Objects.FindID(request.GetObjectId())
	if err != nil {
		return nil, fmt.Errorf("fwd: ObjectNID failed, err %v", err)
	}
	return &fwdpb.ObjectNIDReply{Nid: uint64(obj.NID())}, nil
}
