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

// Package fwdcontext contains routines for managing the context
// of the forwarding engine.
package fwdcontext

import (
	"fmt"
	"sync"

	log "github.com/golang/glog"
	"github.com/google/gopacket"

	"github.com/openconfig/lemming/dataplane/forwarding/infra/deadlock"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdattribute"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdobject"
	"github.com/openconfig/lemming/dataplane/forwarding/util/queue"
	pktiopb "github.com/openconfig/lemming/dataplane/proto/packetio"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

// An NotificationCallback generates events to a notification service.
type NotificationCallback func(*fwdpb.EventDesc)

// Port is an interface for reading and writing to a port.
type Port interface {
	gopacket.PacketDataSource
	WritePacketData(data []byte) error
}

// FakePortManager is an interface for creating custom ports for fwdpb FakePort.
type FakePortManager interface {
	CreatePort(string) (Port, error)
}

// A Context encapsulates the state of an instance of the forwarding engine.
//
// A context is the domain for synchronization. There are two users of context;
// provisioning and packet processing. Provisioning takes a rw lock on the
// context, and modifies objects within the context. Packet processing takes
// a read-lock on the context to process packets. Since provisioning and packet
// processing directly manipulate objects within the context, they must
// explicitly take the appropriate lock.
//
// The context provides mechanisms to send notifications (notify) and packets
// (punt). Notify is non-blocking and sends notifications via an unbounded
// queue. Punt is a blocking call, and the caller is responsible for ordering
// and blocking guarantees.
type Context struct {
	sync.RWMutex                  // Synchronization between provisioning and forwarding
	Objects      *fwdobject.Table // Set of all visible forwarding objects
	ID           string           // ID of the context
	Instance     string           // Name of the forwarding engine instance
	Attributes   fwdattribute.Set

	notifyMu sync.Mutex   // Mutex protecting notification queue
	notify   *queue.Queue // Notification service

	eventMu     sync.Mutex // Mutex protecting the event notification
	nextEventID uint64     // Id of the next event id

	// FakePortManager is the implementation of the port creator for the Fake port type.
	FakePortManager FakePortManager
	cpuPortSink     CPUPortSink
	cpuPortSinkDone func()
}

// New creates a new forwarding context with the specified id and fwd engine
// name. The id identifies the forwarding context in an forwarding engine
// instance, and the instance identifies the forwarding engine instance in the
// universe.
func New(id, instance string) *Context {
	return &Context{
		Objects:    fwdobject.NewTable(),
		ID:         id,
		Attributes: fwdattribute.NewSet(),
		Instance:   instance,
	}
}

// String returns a formatted string representing the context.
func (ctx *Context) String() string {
	str := fmt.Sprintf("Ctx=%v;Instance=%v;NextEvent=%v", ctx.ID, ctx.Instance, ctx.nextEventID)
	if ctx.GetNotificationQueue() != nil {
		str += fmt.Sprintf("<Queue=%v>;", ctx.GetNotificationQueue())
	}
	return str
}

// GetNotificationQueue returns a pointer to the queue of notifications in the context.
func (ctx *Context) GetNotificationQueue() *queue.Queue {
	ctx.notifyMu.Lock()
	defer ctx.notifyMu.Unlock()
	return ctx.notify
}

// SetNotificationQueue sets the notification queue.
func (ctx *Context) SetNotificationQueue(val *queue.Queue) {
	ctx.notifyMu.Lock()
	defer ctx.notifyMu.Unlock()
	ctx.notify = val
}

// SetNotification sets the notification service for the context. If the
// notification service is set to nil, notifications are disabled for the context.
func (ctx *Context) SetNotification(call NotificationCallback) error {
	if call == nil {
		if nq := ctx.GetNotificationQueue(); nq != nil {
			nq.Close()
		}
		ctx.SetNotificationQueue(nil)
		return nil
	}

	h := func(v interface{}) {
		if event, ok := v.(*fwdpb.EventDesc); ok {
			call(event)
		}
	}
	n, err := queue.NewUnbounded("notification")
	if err != nil {
		return err
	}
	ctx.SetNotificationQueue(n)
	n.Run()
	go func() {
		for {
			v, ok := <-n.Receive()
			if !ok {
				return
			}
			h(v)
		}
	}()
	return nil
}

// Notify enqueues a notification request if there is a notification service.
// This is a non-blocking call.
func (ctx *Context) Notify(event *fwdpb.EventDesc) error {
	nq := ctx.GetNotificationQueue()
	if nq == nil {
		return fmt.Errorf("fwdcontext: unable to send notification in context %v, nil queue", ctx)
	}

	timer := deadlock.NewTimer(deadlock.Timeout, fmt.Sprintf("Notifying event %+v in context %v", event, ctx))
	defer timer.Stop()

	// Update the event id.
	ctx.eventMu.Lock()
	event.SequenceNumber = ctx.nextEventID
	ctx.nextEventID++
	ctx.eventMu.Unlock()

	return nq.Write(event)
}

type CPUPortSink func(*pktiopb.PacketOut) error

// SetCPUPortSink sets the port control service for the context
func (ctx *Context) SetCPUPortSink(fn CPUPortSink, doneFn func()) error {
	ctx.cpuPortSink = fn
	ctx.cpuPortSinkDone = doneFn
	return nil
}

// PacketSink returns a handler to port control service
func (ctx *Context) CPUPortSink() CPUPortSink {
	return ctx.cpuPortSink
}

// Cleanup cleans up the context.
// It first cleans up the objects that satisfy isPort.
// Then it unblocks the caller by sending a message on the channel.
// Then it cleans up the rest of the objects.
func (ctx *Context) Cleanup(ch chan bool, isPort func(*fwdpb.ObjectId) bool) {
	ctx.SetNotification(nil)

	if ctx.cpuPortSinkDone != nil {
		ctx.cpuPortSinkDone()
	}
	ctx.SetCPUPortSink(nil, nil)

	ids := ctx.Objects.IDs()

	// First remove the ports.
	for _, id := range ids {
		obj := &fwdpb.ObjectId{
			Id: string(id),
		}
		if isPort(obj) {
			log.Infof("Clean up port %v.", id)
			ctx.Objects.Remove(obj, true /*forceCleanup*/)
		}
	}

	// Then unblock the caller.
	ch <- true

	// And finally remove the other objects.
	for _, id := range ids {
		obj := &fwdpb.ObjectId{
			Id: string(id),
		}
		if !isPort(obj) {
			ctx.Objects.Remove(obj, true /*forceCleanup*/)
		}
	}
}
