// Copyright 2023 Google LLC
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

// Package attrmgr contains a SAI attribute key/value store.
// Each object has a set of attributes: a map of attribute id (enum value) to value (a number of different types).
// Attributes are set using Create and Set RPCs and retrieved using Get RPCs.
package attrmgr

import (
	"context"
	"fmt"
	"log/slog"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"

	saipb "github.com/openconfig/lemming/dataplane/proto/sai"
)

// AttrMgr stores and retrieve the SAI attributes.
type AttrMgr struct {
	mu sync.Mutex
	// attrs is a map of object id (string) to a map of attributes (key: attr id, some enum value).
	attrs   map[string]map[int32]*protoreflect.Value
	nextOid atomic.Uint64
	// idToType maps an object id to its SAI type.
	idToType map[string]saipb.ObjectType
	// msgEnumToFieldNum maps a proto message name to a map of an attribute enum to its corresponding proto field.
	// For example, for SwitchAttribute SWITCH_ATTR_MAX_SYSTEM_CORES (enum val 182) -> field max_system_cores (num 172)
	msgEnumToFieldNum map[string]map[int32]int
	switchID          string
}

func deleteOID(mgr *AttrMgr, oid string) error {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()
	if _, ok := mgr.attrs[oid]; !ok {
		return fmt.Errorf("OID not found: %s", oid)
	}
	delete(mgr.attrs, oid)
	if oid == mgr.switchID {
		mgr.switchID = ""
	}
	return nil
}

// New returns a new AttrMgr.
func New() *AttrMgr {
	mgr := &AttrMgr{
		attrs:             make(map[string]map[int32]*protoreflect.Value),
		idToType:          make(map[string]saipb.ObjectType),
		msgEnumToFieldNum: make(map[string]map[int32]int),
	}
	return mgr
}

func (mgr *AttrMgr) set(id string, attr int32, val *protoreflect.Value) {
	if _, ok := mgr.attrs[id]; !ok {
		mgr.attrs[id] = make(map[int32]*protoreflect.Value)
	}
	mgr.attrs[id][attr] = val
}

const protoNS = "lemming.dataplane.sai"

// InvokeAndSave calls the RPC method and saves the attributes in the request and returns the RPC response.
// This is the same behavior as the Interceptor, except for invoking server methods directly (not using gRPC).
func InvokeAndSave[T proto.Message, S proto.Message](ctx context.Context, mgr *AttrMgr, rpc func(context.Context, T) (S, error), req T) (S, error) {
	resp, err := rpc(ctx, req)
	// Ignore unimplemented errors, so that we don't have to implement APIs we don't support.
	if st, _ := status.FromError(err); err != nil && (st.Code() != codes.Unimplemented) {
		return resp, err
	}
	respMsg, err := createResponse(req, resp)
	if err != nil {
		return resp, err
	}
	id, err := mgr.getID(req, respMsg)
	if err != nil {
		slog.WarnContext(ctx, "failed to get id", "err", err)
		return respMsg.(S), nil
	}
	mgr.storeAttributes(id, req)

	return respMsg.(S), nil
}

// Interceptor returns a gRPC interceptor that automatically store values set in requests and fills in responses with the stored values.
func (mgr *AttrMgr) Interceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	if !strings.Contains(info.FullMethod, protoNS) {
		return handler(ctx, req)
	}
	resp, err := handler(ctx, req)
	// Ignore unimplemented error for Get*Attribute.
	if err != nil {
		if st, _ := status.FromError(err); st.Code() != codes.Unimplemented || !strings.Contains(info.FullMethod, "Get") {
			return resp, err
		}
	}
	reqMsg := req.(proto.Message)
	// If the resp is nil, then create a new response type.
	respMsg, err := createResponse(reqMsg, resp)
	if err != nil {
		return resp, err
	}

	switch {
	case strings.Contains(info.FullMethod, "Create") || strings.Contains(info.FullMethod, "Set"):
		id, err := mgr.getID(reqMsg, respMsg)
		if err != nil {
			slog.WarnContext(ctx, "failed to get id", "err", err)
			return respMsg, nil
		}
		mgr.storeAttributes(id, reqMsg)
	case strings.Contains(info.FullMethod, "Get") && strings.Contains(info.FullMethod, "Attribute"):
		if err := mgr.PopulateAttributes(reqMsg, respMsg); err != nil {
			return nil, err
		}
	case strings.Contains(info.FullMethod, "Remove"):
		switchID, ok := mgr.GetSwitchID()
		if ok {
			val := mgr.GetAttribute(switchID, int32(saipb.SwitchAttr_SWITCH_ATTR_RESTART_WARM))
			if isWarm, ok := val.(bool); ok && isWarm {
				slog.InfoContext(ctx, "attrmgr.Interceptor: Skipping attribute deletion for RemoveSwitch RPC during warm restart.")
				return respMsg, nil
			}
		}
		id, err := mgr.getID(reqMsg, respMsg)
		if err != nil {
			slog.WarnContext(ctx, "failed to get id", "err", err)
			return respMsg, nil
		}
		if err := deleteOID(mgr, id); err != nil {
			return nil, err
		}
	}
	return respMsg, nil
}

// Reset resets all fields attributes in the manager.
func (mgr *AttrMgr) Reset() {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()
	mgr.attrs = make(map[string]map[int32]*protoreflect.Value)
	mgr.idToType = make(map[string]saipb.ObjectType)
	mgr.msgEnumToFieldNum = make(map[string]map[int32]int)
	mgr.switchID = ""
	mgr.nextOid.Store(0)
}

func createResponse(req proto.Message, resp any) (proto.Message, error) {
	if resp == nil {
		respName := strings.ReplaceAll(string(req.ProtoReflect().Descriptor().FullName()), "Request", "Response")
		pType, err := protoregistry.GlobalTypes.FindMessageByName(protoreflect.FullName(respName))
		if err != nil {
			return nil, err
		}
		resp = pType.New().Interface()
	} else if val := reflect.ValueOf(resp); val.IsNil() { // If resp is a typed-nil value.
		return reflect.New(val.Type().Elem()).Interface().(proto.Message), nil
	}
	return resp.(proto.Message), nil
}

func (mgr *AttrMgr) getEnumToFields(message protoreflect.MessageDescriptor) map[int32]int {
	name := string(message.FullName())

	if enumToFields, ok := mgr.msgEnumToFieldNum[name]; ok {
		return enumToFields
	}
	mgr.msgEnumToFieldNum[name] = make(map[int32]int)
	for i := 0; i < message.Fields().Len(); i++ {
		opt, ok := message.Fields().Get(i).Options().(*descriptorpb.FieldOptions)
		if !ok {
			continue
		}
		mgr.msgEnumToFieldNum[name][proto.GetExtension(opt, saipb.E_AttrEnumValue).(int32)] = i
	}
	return mgr.msgEnumToFieldNum[name]
}

// populateAttributes fills the resp with the requests attributes.
// This must called with GetFooAttributeRequest and GetFooAttributeResponse message types.
func (mgr *AttrMgr) PopulateAttributes(req, resp proto.Message) error {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()
	id, err := mgr.getID(req, resp)
	if err != nil {
		return err
	}
	attrTypeFd := req.ProtoReflect().Descriptor().Fields().ByTextName("attr_type")
	attrFd := resp.ProtoReflect().Descriptor().Fields().ByTextName("attr")
	if attrFd == nil || attrTypeFd == nil {
		return fmt.Errorf("req and resp didn't have required attributes")
	}

	enumToFields := mgr.getEnumToFields(attrFd.Message())
	attrs := resp.ProtoReflect().Mutable(attrFd).Message()
	reqList := req.ProtoReflect().Get(attrTypeFd).List()

	for i := 0; i < reqList.Len(); i++ {
		enumVal := reqList.Get(i).Enum()
		val, ok := mgr.attrs[id][int32(enumVal)]
		if !ok {
			return fmt.Errorf("requested attribute not set: %v in OID: %v", attrTypeFd.Enum().Values().ByNumber(reqList.Get(i).Enum()).Name(), id)
		}
		// Empty lists exist so they are not errors, but are not settable.
		if val != nil {
			attrs.Set(attrs.Descriptor().Fields().Get(enumToFields[int32(enumVal)]), *val)
		}
	}
	return nil
}

// PopulateAllAttributes retrieves all attributes for an object.
// Supported message types FooAttribute, CreateFooRequest, SetFooRequest.
func (mgr *AttrMgr) PopulateAllAttributes(id string, msg proto.Message) (rerr error) {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()
	defer func() {
		if r := recover(); r != nil {
			rerr = fmt.Errorf("protoreflect error: %v", r)
		}
	}()

	desc := msg.ProtoReflect().Descriptor()
	for i := 0; i < desc.Fields().Len(); i++ {
		opt, ok := desc.Fields().Get(i).Options().(*descriptorpb.FieldOptions)
		if !ok {
			continue
		}
		enumVal := proto.GetExtension(opt, saipb.E_AttrEnumValue).(int32)
		if enumVal == 0 {
			continue
		}
		val, ok := mgr.attrs[id][enumVal]
		if !ok || val == nil {
			continue
		}

		msg.ProtoReflect().Set(desc.Fields().Get(i), *val)
	}
	return nil
}

// StoreAttributes stores all the attributes in the message.
// Note: for lists, a nil lists is an unset attribute, but a non-nil empty list is set.
// so querying a nil list returns an error, even though they look the same on the wire.
func (mgr *AttrMgr) StoreAttributes(id uint64, msg proto.Message) {
	mgr.storeAttributes(fmt.Sprint(id), msg)
}

// GetType returns the SAI type for the object.
func (mgr *AttrMgr) GetType(id string) saipb.ObjectType {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()
	val, ok := mgr.idToType[id]
	if !ok {
		return saipb.ObjectType_OBJECT_TYPE_NULL
	}
	return val
}

// GetAttribute returns the value of the requested attribute.
func (mgr *AttrMgr) GetAttribute(id string, attr int32) any {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()
	if objAttrs, ok := mgr.attrs[id]; ok {
		if val := objAttrs[attr]; val != nil {
			return val.Interface()
		}
	}
	return nil
}

// GetSwitchID returns the ID of the switch if it exists.
func (mgr *AttrMgr) GetSwitchID() (string, bool) {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()
	if mgr.switchID != "" {
		return mgr.switchID, true
	}
	return "", false
}

// GetType returns the SAI type for the object.
func (mgr *AttrMgr) SetType(id string, t saipb.ObjectType) {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()
	mgr.idToType[id] = t
	if t == saipb.ObjectType_OBJECT_TYPE_SWITCH {
		mgr.switchID = id
	}
}

// storeAttributes stores all the attributes in the message.
func (mgr *AttrMgr) storeAttributes(id string, msg proto.Message) {
	ty := proto.GetExtension(msg.ProtoReflect().Descriptor().Options(), saipb.E_SaiType).(saipb.ObjectType)
	if ty != saipb.ObjectType_OBJECT_TYPE_UNSPECIFIED {
		mgr.SetType(id, ty)
	}

	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	// Protoreflect treats nil lists and empty lists as the same. However We want to store the value of empty lists, but not nil lists.
	// So use regular go reflect for that case.
	rv := reflect.ValueOf(msg).Elem()
	rt := reflect.TypeOf(msg).Elem()
	for i := 0; i < rt.NumField(); i++ {
		tag := rt.Field(i).Tag.Get("protobuf")
		if tag == "" {
			continue
		}
		var fName string
		for _, v := range strings.Split(tag, ",") {
			if strings.HasPrefix(v, "name=") {
				fName = strings.TrimPrefix(v, "name=")
			}
		}
		fd := msg.ProtoReflect().Descriptor().Fields().ByTextName(fName)
		opt, ok := fd.Options().(*descriptorpb.FieldOptions)
		if !ok {
			continue
		}
		enumValue := proto.GetExtension(opt, saipb.E_AttrEnumValue).(int32)
		if enumValue == 0 {
			continue
		}
		if v := msg.ProtoReflect().Get(fd); msg.ProtoReflect().Has(fd) {
			mgr.set(id, enumValue, &v)
		} else if fd.IsList() && !rv.Field(i).IsNil() {
			mgr.set(id, enumValue, nil)
		}
	}
}

// NextID returns the next available object id.
func (mgr *AttrMgr) NextID() uint64 {
	return mgr.nextOid.Add(1)
}

// getID returns the id from either the request or response. If the id is unset, it allocates a new one.
func (mgr *AttrMgr) getID(req, resp proto.Message) (string, error) {
	msgs := []proto.Message{req, resp}
	for _, msg := range msgs {
		if fd := msg.ProtoReflect().Descriptor().Fields().ByTextName("oid"); fd != nil {
			v := msg.ProtoReflect().Get(fd).Uint()
			if v == 0 {
				id := mgr.NextID()
				msg.ProtoReflect().Set(fd, protoreflect.ValueOfUint64(id))
				return fmt.Sprint(id), nil
			}
			return fmt.Sprint(v), nil
		}
	}
	fd := req.ProtoReflect().Descriptor().Fields().ByTextName("entry")
	if fd == nil {
		return "", fmt.Errorf("no id found in message")
	}
	entry := req.ProtoReflect().Get(fd).Message().Interface()
	pBytes, err := proto.Marshal(entry)
	if err != nil {
		return "", err
	}
	return string(pBytes), nil
}
