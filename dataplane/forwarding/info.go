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

package forwarding

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	log "github.com/golang/glog"

	"github.com/openconfig/lemming/dataplane/forwarding/fwdport"
	"github.com/openconfig/lemming/dataplane/forwarding/fwdtable"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdattribute"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdcontext"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdobject"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	"github.com/openconfig/lemming/dataplane/forwarding/protocol"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

// indent creates a string with the specified number of spaces
func indent(count int) []byte {
	str := make([]byte, count)
	for pos := 0; pos < count; pos++ {
		str[pos] = ' '
	}
	return str
}

// format formats the description of various Lucius objects by replacing ';'
// with ' ', and replacing '<' and '>' with suitable indentation.
// TODO: Find a better approach. This hack allows us to use the same
// implementation for debug string (single line description) and display
// (formatted multiline description)
func format(str string) string {
	t := []byte(str)
	result := make([]byte, len(t))
	count := 0
	for _, c := range str {
		c := byte(c)
		switch c {
		case ';':
			result = append(result, ' ')
		case '<':
			count++
			result = append(result, '\n')
			t := indent(count)
			result = append(result, t...)
		case '>':
			t := indent(count)
			result = append(result, t...)
			count--
		default:
			result = append(result, c)
		}
	}
	return string(result)
}

// AttributeInfo returns the attributes from the global, context and arg as specified..
func AttributeInfo(ctx *fwdcontext.Context, arg interface{}) string {
	var buffer []string
	if attr := fwdattribute.Global; len(attr) != 0 {
		buffer = append(buffer, "Global Attributes:")
		buffer = append(buffer, attr.String())
	}
	if ctx != nil {
		if attr := ctx.Attributes; len(attr) != 0 {
			buffer = append(buffer, "Context Attributes: ")
			buffer = append(buffer, attr.String())
		}
	}
	if obj, ok := arg.(fwdobject.Object); ok {
		if attr := obj.Attributes(); len(attr) != 0 {
			buffer = append(buffer, "Object Attributes:")
			buffer = append(buffer, attr.String())
		}
	}
	return strings.Join(buffer, "\n")
}

// CounterInfo returns the counters of an object formatted as string.
func CounterInfo(obj fwdobject.Object) string {
	var names []string
	counters := obj.Counters()
	for _, c := range counters {
		names = append(names, c.ID.String())
	}
	sort.Strings(names)
	buffer := []string{"Counters:"}
	for _, n := range names {
		c := counters[fwdpb.CounterId(fwdpb.CounterId_value[n])]
		buffer = append(buffer, c.String())
	}
	buffer = append(buffer, "")
	return strings.Join(buffer, "\n")
}

// PortInfo returns the details of the specified Port.
// TODO: Use PortInfo, probably needs to return object info string alongside the PortElementInfo return value.
func PortInfo(ctx *fwdcontext.Context, arg interface{}) (*fwdpb.PortElementInfo, error) {
	port, ok := arg.(fwdport.Port)
	if !ok {
		return nil, fmt.Errorf("arg %v is not a Port", arg)
	}
	desc := &fwdpb.PortDesc{
		PortType: port.Type(),
		desc.PortId: fwdport.GetID(port),
	}
	switch desc.PortType {
		default:
			return nil, fmt.Errorf("unknown PortType %v", desc.PortType)
		case fwdpb.PortType_PORT_TYPE_CPU_PORT:
			cpu, ok := port.(*ports.CPUPort)
			if !ok {
				return nil, fmt.Errorf("ports: Unable to create cpu port")
			}
			desc.Port = &fwdpb.CPUPortDesc{
				QueueId: cpu.queueID,
				QueueLength: int32(cpu.queue.max),
				ExportFieldIds: cpu.export,
				Remote:  cpu.remote,
			}
		// TODO: cases for other port types
	}
	var counters []*fwdpb.Counter
	for _, c := range port.Counters() {
		counters = append(counters, &fwdpb.Counter{Id: c.ID, Value: c.Value})
	}
	return &fwdpb.PortElementInfo{Desc: desc, Counters: counters}, nil
}

// TableInfo returns the details of the specified Table as a string.
func TableInfo(ctx *fwdcontext.Context, arg interface{}) string {
	table, ok := arg.(fwdtable.Table)
	if !ok {
		return ""
	}
	var buffer []string
	buffer = append(buffer, AttributeInfo(ctx, table))
	buffer = append(buffer, "Description: ")
	buffer = append(buffer, format(table.String()))
	buffer = append(buffer, "Entries:")
	var list []string
	for _, entry := range table.Entries() {
		list = append(list, entry)
	}
	buffer = append(buffer, format(strings.Join(list, "")))
	buffer = append(buffer, CounterInfo(table))
	return strings.Join(buffer, "\n")
}

// ObjectInfo returns the details of the specified object as a string.
func ObjectInfo(ctx *fwdcontext.Context, arg interface{}) string {
	obj, ok := arg.(fwdobject.Object)
	if !ok {
		return ""
	}
	var buffer []string
	buffer = append(buffer, AttributeInfo(ctx, obj))
	buffer = append(buffer, "Description: ")
	buffer = append(buffer, format(obj.String()))
	buffer = append(buffer, CounterInfo(obj))
	return strings.Join(buffer, "\n")
}

// ContextInfo returns the details of the specified context as a string.
func ContextInfo(ctx *fwdcontext.Context, _ interface{}) string {
	var buffer []string
	buffer = append(buffer, AttributeInfo(ctx, nil))
	buffer = append(buffer, format(ctx.String()))
	return strings.Join(buffer, "\n")
}

// A Handler returns the contents of an info element.
type Handler func(*fwdcontext.Context, interface{}) string

// An entry represents a Handler and its argument.
type entry struct {
	handler Handler
	arg     interface{}
	ctx     *fwdcontext.Context
	name    string
}

// An InfoList stores info elements for each forwarding object. Each info
// element is identified by a string name and is associated with a handler and
// reference to the actual object itself.
type InfoList struct {
	mu       sync.Mutex
	handlers map[string]*entry // Map of handlers indexed by the info element name
}

// NewInfoList returns a new information list.
func NewInfoList() *InfoList {
	l := &InfoList{
		handlers: make(map[string]*entry),
	}

	l.add(nil, nil, "global", AttributeInfo)
	return l
}

// objectKey returns a key formed by the object's context and ID.
func objectKey(contextID string, id fwdobject.ID) string {
	return fmt.Sprintf("%v,%v", contextID, id)
}

// contextKey returns a key formed by the context ID.
func contextKey(contextID string) string {
	return fmt.Sprintf("%v", contextID)
}

// add adds a handler by key.
func (l *InfoList) add(ctx *fwdcontext.Context, arg interface{}, key string, h Handler) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.handlers[key] = &entry{
		handler: h,
		arg:     arg,
		ctx:     ctx,
		name:    key,
	}
}

// AddContext adds an info element for a forwarding context.
func (l *InfoList) AddContext(ctx *fwdcontext.Context) {
	if ctx == nil {
		log.Errorf("Invalid context.")
		return
	}
	l.add(ctx, nil, contextKey(ctx.ID), ContextInfo)
}

// AddObject adds an info element for an object.
func (l *InfoList) AddObject(ctx *fwdcontext.Context, obj fwdobject.Object) {
	if ctx == nil || obj == nil {
		log.Errorf("Invalidcontext and/or object.")
		return
	}

	// Use a special handler for tables.
	h := ObjectInfo
	if _, ok := obj.(fwdtable.Table); ok {
		h = TableInfo
	}
	l.add(ctx, obj, objectKey(ctx.ID, obj.ID()), h)
}

// remove removes an info element corresponding to the specified key.
func (l *InfoList) remove(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	delete(l.handlers, key)
}

// RemoveContext removes an info element corresponding to the specified context.
func (l *InfoList) RemoveContext(ctx *fwdcontext.Context) {
	if ctx == nil {
		log.Errorf("Invalid context.")
		return
	}
	l.remove(contextKey(ctx.ID))
}

// RemoveObject removes an info element corresponding to the specified objectID.
func (l *InfoList) RemoveObject(ctx *fwdcontext.Context, oid *fwdpb.ObjectId) {
	if ctx == nil || oid == nil {
		log.Errorf("Invalid context and/or object.")
		return
	}
	l.remove(objectKey(ctx.ID, fwdobject.ID(oid.GetId())))
}

// List retrieves a list of all information elements.
func (l *InfoList) List() []string {
	l.mu.Lock()
	defer l.mu.Unlock()

	var names sort.StringSlice
	for n := range l.handlers {
		names = append(names, n)
	}
	names.Sort()
	return names
}

// allInfo returns all the information available in the specified information element.
func (l *InfoList) allInfo(e *entry) (string, error) {
	return e.handler(e.ctx, e.arg), nil
}

// lookupInfo returns the information related to the lookup of the given packet
// in the specified table.
func (l *InfoList) lookupInfo(e *entry, arg []byte, l2 fwdpb.PacketHeaderId) (string, error) {
	if e.arg == nil || e.ctx == nil {
		return "", fmt.Errorf("infolist: %v is not associated with an object %v and/or context %v", e.name, e.arg, e.ctx)
	}

	e.ctx.RLock()
	defer e.ctx.RUnlock()

	table, ok := e.arg.(fwdtable.Table)
	if !ok {
		return "", fmt.Errorf("infolist: %v is not a table", e.name)
	}

	packet, err := fwdpacket.New(l2, arg)
	if err != nil {
		return "", fmt.Errorf("infolist: Unable to create packet, err %v", err)
	}
	// By default, the packet logging is dependent on the --v flag. Override this behavior
	if pkt, ok := packet.(*protocol.Packet); ok {
		pkt.OverrideGlobalLogLevel()
	}

	packet.Debug(true)
	table.Process(packet, table)
	m := packet.LogMsgs()
	return strings.Join(m, "\n"), nil
}

// lookupPacket returns the information related to the processing of the given packet
// in the specified port. Note that the packet is actually injected into the network.
func (l *InfoList) lookupPacket(e *entry, arg []byte, l2 fwdpb.PacketHeaderId, dir fwdpb.PortAction) (string, error) {
	if e.arg == nil || e.ctx == nil {
		return "", fmt.Errorf("infolist: %v is not associated with an object %v and/or context %v", e.name, e.arg, e.ctx)
	}

	e.ctx.RLock()
	defer e.ctx.RUnlock()

	port, ok := e.arg.(fwdport.Port)
	if !ok {
		return "", fmt.Errorf("infolist: %v is not a port", e.name)
	}

	packet, err := fwdpacket.New(l2, arg)
	if err != nil {
		return "", fmt.Errorf("infolist: Unable to create packet, err %v", err)
	}
	// By default, the packet logging is dependent on the --v flag. Override this behavior
	if pkt, ok := packet.(*protocol.Packet); ok {
		pkt.OverrideGlobalLogLevel()
	}

	packet.Debug(true)
	fwdport.Process(port, packet, dir, e.ctx, "Info")
	m := packet.LogMsgs()
	return strings.Join(m, "\n"), nil
}

// Element retrieves the contents of a specific information element.
func (l *InfoList) Element(name string, infoType fwdpb.InfoType, frame []byte, start fwdpb.PacketHeaderId) (string, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	e, ok := l.handlers[name]
	if !ok {
		return "", fmt.Errorf("infolist: Unable to find handler for %v", name)
	}

	switch infoType {
	case fwdpb.InfoType_INFO_TYPE_ALL:
		return l.allInfo(e)

	case fwdpb.InfoType_INFO_TYPE_LOOKUP:
		return l.lookupInfo(e, frame, start)

	case fwdpb.InfoType_INFO_TYPE_PORT_INPUT:
		return l.lookupPacket(e, frame, start, fwdpb.PortAction_PORT_ACTION_INPUT)

	case fwdpb.InfoType_INFO_TYPE_PORT_OUTPUT:
		return l.lookupPacket(e, frame, start, fwdpb.PortAction_PORT_ACTION_OUTPUT)

	default:
		return "", fmt.Errorf("infolist: Unable to handle infoType %v for %v", infoType, name)
	}
}
