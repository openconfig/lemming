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
	"time"

	log "github.com/golang/glog"
	"google.golang.org/grpc"

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
	fwdpb.UnimplementedServiceServer
	fwdpb.UnimplementedInfoServer
	fwdpb.UnimplementedPacketSinkServer

	mu  sync.Mutex
	ctx map[string]*fwdcontext.Context // forwarding contexts indexed by name

	name         string                                     // name of the forwarding engine
	info         *InfoList                                  // list of info elements that can be queried
	conn         map[string]*grpc.ClientConn                // client connections indexed by address
	notification map[string]fwdcontext.NotificationCallback // notification callback indexed by address
	packet       map[string]fwdcontext.PacketCallback       // packet callback indexed by address
}

// New creates a new forwarding instance using the specified name.
func New(name string) *Engine {
	return &Engine{
		name: name,
		ctx:  make(map[string]*fwdcontext.Context),
		info: NewInfoList(),
	}
}

// client returns a gRPC client connected to the specified address.
// It is assumed that clients are looked up while the service is locked.
func (e *Engine) client(addr string) (*grpc.ClientConn, error) {
	if c, ok := e.conn[addr]; ok {
		return c, nil
	}
	c, err := grpc.Dial(addr)
	if err != nil {
		return nil, fmt.Errorf("service: dial to %v failed, err %v", addr, err)
	}
	e.conn[addr] = c
	return c, nil
}

// GetNotificationCallback returns a callback that posts notifications to the
// specified address.  If the address is "", the returned callback ignores all
// events.
func (e *Engine) GetNotificationCallback(address string) (fwdcontext.NotificationCallback, error) {
	if address == "" {
		return nil, nil
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	if h, ok := e.notification[address]; ok {
		return h, nil
	}
	c, err := e.client(address)
	if err != nil {
		return nil, fmt.Errorf("service: connection to notification service failed, err %v", err)
	}
	client := fwdpb.NewNotificationClient(c)
	h := func(ed *fwdpb.EventDesc) {
		client.Notify(context.TODO(), ed)
	}
	e.notification[address] = h
	return h, nil
}

// GetPacketSinkCallback returns a callback that posts packets to a packet sink
// service at the specified address. If the address is "", the packet sink
// service is disabled for the context.
func (e *Engine) GetPacketSinkCallback(address string) (fwdcontext.PacketCallback, error) {
	if address == "" {
		return nil, nil
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	if h, ok := e.packet[address]; ok {
		return h, nil
	}
	c, err := e.client(address)
	if err != nil {
		return nil, fmt.Errorf("service: connection to packet service failed, err %v", err)
	}
	client := fwdpb.NewPacketSinkClient(c)
	h := func(p *fwdpb.PacketInjectRequest) (*fwdpb.PacketInjectReply, error) {
		// Execute the RPC with a 1 min timeout.
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cancel()
		return client.PacketInject(ctx, p)
	}
	e.packet[address] = h
	return h, nil
}

// UpdateNotification updates the notification service for a context. If the
// notification is set to nil, no notifications are generated for the context.
// The address is the address of the notification service (used in queries)
// in the host:port format.
func (e *Engine) UpdateNotification(contextID *fwdpb.ContextId, notification fwdcontext.NotificationCallback, address string) error {
	if contextID == nil {
		return errors.New("fwd: UpdateNotification failed, No context ID")
	}

	ctx, err := e.FindContext(contextID)
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
func (e *Engine) UpdatePacketSink(contextID *fwdpb.ContextId, packet fwdcontext.PacketCallback, address string) error {
	if contextID == nil {
		return errors.New("fwd: UpdatePacketSink failed, No context ID")
	}

	ctx, err := e.FindContext(contextID)
	if err != nil {
		return fmt.Errorf("fwd: UpdatePacketSink failed, err %v", err)
	}

	ctx.Lock()
	defer ctx.Unlock()
	ctx.SetPacketSink(packet, address)
	return nil
}

// ContextCreate creates a new context. Note that if the packet sink and/or
// notification services are specified but not reachable, the context creation
// fails.
func (e *Engine) ContextCreate(_ context.Context, request *fwdpb.ContextCreateRequest) (*fwdpb.ContextCreateReply, error) {
	paddr := request.GetPacketAddress()
	naddr := request.GetNotificationAddress()
	ns, err := e.GetNotificationCallback(naddr)
	if err != nil {
		return nil, err
	}

	ps, err := e.GetPacketSinkCallback(paddr)
	if err != nil {
		return nil, err
	}

	cid := request.GetContextId()
	if err := e.contextCreateByID(cid); err != nil {
		return nil, err
	}

	if err := e.UpdateNotification(cid, ns, naddr); err != nil {
		return nil, err
	}
	return &fwdpb.ContextCreateReply{}, e.UpdatePacketSink(cid, ps, paddr)
}

// ContextUpdate updates a forwarding context. Note that if the packet sink and/or
// notification services are specified but not reachable, the context update
// fails.
func (e *Engine) ContextUpdate(_ context.Context, request *fwdpb.ContextUpdateRequest) (*fwdpb.ContextUpdateReply, error) {
	paddr := request.GetPacketAddress()
	naddr := request.GetNotificationAddress()
	cid := request.GetContextId()

	for _, op := range request.GetOperations() {
		switch op {
		case fwdpb.ContextUpdateRequest_OPERATION_UPDATE_PACKET_ADDRESS:
			ps, err := e.GetPacketSinkCallback(paddr)
			if err != nil {
				return nil, err
			}
			if err = e.UpdatePacketSink(cid, ps, paddr); err != nil {
				return nil, err
			}
		case fwdpb.ContextUpdateRequest_OPERATION_UPDATE_NOTIFICATION_ADDRESS:
			ns, err := e.GetNotificationCallback(naddr)
			if err != nil {
				return nil, err
			}
			if err = e.UpdateNotification(cid, ns, naddr); err != nil {
				return nil, err
			}
		}
	}
	return &fwdpb.ContextUpdateReply{}, nil
}

// contextCreateByID creates a new context with the given ID, erroring if it already exists.
func (e *Engine) contextCreateByID(contextID *fwdpb.ContextId) error {
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
func (e *Engine) ContextDelete(_ context.Context, request *fwdpb.ContextDeleteRequest) (*fwdpb.ContextDeleteReply, error) {
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
func (e *Engine) ContextList(_ context.Context, request *fwdpb.ContextListRequest) (*fwdpb.ContextListReply, error) {
	timer := deadlock.NewTimer(deadlock.Timeout, fmt.Sprintf("Processing %+v", request))
	defer timer.Stop()

	e.mu.Lock()
	defer e.mu.Unlock()

	reply := &fwdpb.ContextListReply{}
	for _, ctx := range e.ctx {
		reply.Contexts = append(reply.Contexts, &fwdpb.ContextAttr{
			ContextId: &fwdpb.ContextId{
				Id: string(ctx.ID),
			},
			PacketAddress:       ctx.PacketAddress,
			NotificationAddress: ctx.NotificationAddress,
		})
	}
	return reply, nil
}

// FindContext finds the specified context.
func (e *Engine) FindContext(contextID *fwdpb.ContextId) (*fwdcontext.Context, error) {
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
func (e *Engine) ObjectDelete(_ context.Context, request *fwdpb.ObjectDeleteRequest) (*fwdpb.ObjectDeleteReply, error) {
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
func (e *Engine) ObjectList(_ context.Context, request *fwdpb.ObjectListRequest) (*fwdpb.ObjectListReply, error) {
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
func (e *Engine) ObjectCounters(_ context.Context, request *fwdpb.ObjectCountersRequest) (*fwdpb.ObjectCountersReply, error) {
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
func (*Engine) AttributeList(_ context.Context, request *fwdpb.AttributeListRequest) (*fwdpb.AttributeListReply, error) {
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
func (e *Engine) AttributeUpdate(_ context.Context, request *fwdpb.AttributeUpdateRequest) (*fwdpb.AttributeUpdateReply, error) {
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

// PortCreate creates a port.
func (e *Engine) PortCreate(_ context.Context, request *fwdpb.PortCreateRequest) (*fwdpb.PortCreateReply, error) {
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
	return reply, nil
}

// PortUpdate updates a port.
func (e *Engine) PortUpdate(_ context.Context, request *fwdpb.PortUpdateRequest) (*fwdpb.PortUpdateReply, error) {
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
func (e *Engine) PortState(_ context.Context, request *fwdpb.PortStateRequest) (*fwdpb.PortStateReply, error) {
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
func (e *Engine) TableCreate(_ context.Context, request *fwdpb.TableCreateRequest) (*fwdpb.TableCreateReply, error) {
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
		return nil, fmt.Errorf("fwd: TableCreate failed, err %v", err)
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
func (e *Engine) TableEntryAdd(_ context.Context, request *fwdpb.TableEntryAddRequest) (*fwdpb.TableEntryAddReply, error) {
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
func (e *Engine) TableEntryRemove(_ context.Context, request *fwdpb.TableEntryRemoveRequest) (*fwdpb.TableEntryRemoveReply, error) {
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
func (e *Engine) TableList(_ context.Context, request *fwdpb.TableListRequest) (*fwdpb.TableListReply, error) {
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
func (e *Engine) SetCreate(_ context.Context, request *fwdpb.SetCreateRequest) (*fwdpb.SetCreateReply, error) {
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
func (e *Engine) SetUpdate(_ context.Context, request *fwdpb.SetUpdateRequest) (*fwdpb.SetUpdateReply, error) {
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
func (e *Engine) FlowCounterCreate(_ context.Context, request *fwdpb.FlowCounterCreateRequest) (*fwdpb.FlowCounterCreateReply, error) {
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
func (e *Engine) FlowCounterQuery(_ context.Context, request *fwdpb.FlowCounterQueryRequest) (*fwdpb.FlowCounterQueryReply, error) {
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

// PacketInject injects a packet in the specified forwarding context and port.
func (e *Engine) PacketInject(_ context.Context, request *fwdpb.PacketInjectRequest) (*fwdpb.PacketInjectReply, error) {
	timer := deadlock.NewTimer(deadlock.Timeout, fmt.Sprintf("Processing %+v", request))
	defer timer.Stop()

	ctx, err := e.FindContext(request.GetContextId())
	if err != nil {
		return nil, fmt.Errorf("fwd: PacketInject failed, err %v", err)
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
		return nil, fmt.Errorf("fwd: PacketInject failed, err %v", err)
	}
	return &fwdpb.PacketInjectReply{}, nil
}

// InfoList retrieves a list of all information elements.
func (e *Engine) InfoList(context.Context, *fwdpb.InfoListRequest) (*fwdpb.InfoListReply, error) {
	return &fwdpb.InfoListReply{
		Names: e.info.List(),
	}, nil
}

// InfoElement retrieves the contents of a specific information element.
func (e *Engine) InfoElement(_ context.Context, request *fwdpb.InfoElementRequest) (*fwdpb.InfoElementReply, error) {
	content, err := e.info.Element(request.GetName(), request.GetType(), request.GetFrame(), request.GetStartHeader())
	if err != nil {
		return nil, err
	}
	return &fwdpb.InfoElementReply{Content: content}, nil
}

// Operation processes incoming OperationRequests and extracts the contained Request.
func (e *Engine) Operation(stream fwdpb.Service_OperationServer) error {
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
