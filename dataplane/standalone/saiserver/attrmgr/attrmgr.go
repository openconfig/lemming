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

	saipb "github.com/openconfig/lemming/dataplane/standalone/proto"
)

// AttrMgr stores and retrieve the SAI attributes.
type AttrMgr struct {
	mu sync.Mutex
	// attrs is a map of object id (string) to a map of attributes (key: attr id, some enum value).
	attrs   map[string]map[int32]protoreflect.Value
	nextOid atomic.Uint64
	// idToType maps an object id to its SAI type.
	idToType map[string]saipb.ObjectType
	// msgEnumToFieldNum maps a proto message name to a map of an attribute enum to its corresponding proto field.
	// For example, for SwitchAttribute SWITCH_ATTR_MAX_SYSTEM_CORES (enum val 182) -> field max_system_cores (num 172)
	msgEnumToFieldNum map[string]map[int32]int
}

// New returns a new AttrMgr.
func New() *AttrMgr {
	mgr := &AttrMgr{
		attrs:             make(map[string]map[int32]protoreflect.Value),
		idToType:          make(map[string]saipb.ObjectType),
		msgEnumToFieldNum: make(map[string]map[int32]int),
	}
	return mgr
}

func (mgr *AttrMgr) set(id string, attr int32, val protoreflect.Value) {
	if _, ok := mgr.attrs[id]; !ok {
		mgr.attrs[id] = make(map[int32]protoreflect.Value)
	}
	mgr.attrs[id][attr] = val
}

const protoNS = "lemming.dataplane.sai"

// InvokeAndSave calls the RPC method and saves the attributes in the request.
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
		return resp, err
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
	// Ignore unimplemented errors, so that we don't have to implement APIs we don't support.
	if st, _ := status.FromError(err); err != nil && (st.Code() != codes.Unimplemented) {
		return resp, err
	}
	reqMsg := req.(proto.Message)
	// If the resp is nil, then create a new response type.
	respMsg, err := createResponse(reqMsg, resp)
	if err != nil {
		return resp, err
	}
	if strings.Contains(info.FullMethod, "Create") || strings.Contains(info.FullMethod, "Set") {
		id, err := mgr.getID(reqMsg, respMsg)
		if err != nil {
			return nil, err
		}
		mgr.storeAttributes(id, reqMsg)
	} else if strings.Contains(info.FullMethod, "Get") {
		id, err := mgr.getID(reqMsg, respMsg)
		if err != nil {
			return nil, err
		}
		if err := mgr.PopulateAttributes(id, reqMsg, respMsg); err != nil {
			return nil, err
		}
	}
	return respMsg, nil
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

// populateAttributes fills the resp with the requests attributes.
// This must called with GetFooAttributeRequest and GetFooAttributeResponse message types.
func (mgr *AttrMgr) PopulateAttributes(id string, req, resp proto.Message) error {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	attrTypeFd := req.ProtoReflect().Descriptor().Fields().ByTextName("attr_type")
	attrFd := resp.ProtoReflect().Descriptor().Fields().ByTextName("attr")
	if attrFd == nil || attrTypeFd == nil {
		return fmt.Errorf("req and resp didn't have required attributes")
	}

	// Populate the msgEnumToFieldNum if it doens't exist for this message.
	enumToFields, ok := mgr.msgEnumToFieldNum[string(attrFd.FullName())]
	if !ok {
		mgr.msgEnumToFieldNum[string(attrFd.FullName())] = make(map[int32]int)
		enumToFields = mgr.msgEnumToFieldNum[string(attrFd.FullName())]
		for i := 0; i < attrFd.Message().Fields().Len(); i++ {
			opt, ok := attrFd.Message().Fields().Get(i).Options().(*descriptorpb.FieldOptions)
			if !ok {
				continue
			}
			enumToFields[proto.GetExtension(opt, saipb.E_AttrEnumValue).(int32)] = i
		}
	}

	attrs := resp.ProtoReflect().Mutable(attrFd).Message()
	reqList := req.ProtoReflect().Get(attrTypeFd).List()

	for i := 0; i < reqList.Len(); i++ {
		enumVal := reqList.Get(i).Enum()
		val, ok := mgr.attrs[id][int32(enumVal)]
		if !ok {
			continue
		}
		attrs.Set(attrs.Descriptor().Fields().Get(enumToFields[int32(enumVal)]), val)
	}
	return nil
}

// PopulateAllAttributes fills the resp with the requests attributes.
// This must called with FooAttribute message type.
func (mgr *AttrMgr) PopulateAllAttributes(id string, msg proto.Message) error {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	attrFd := msg.ProtoReflect().Descriptor()
	// Populate the msgEnumToFieldNum if it doens't exist for this message.
	enumToFields, ok := mgr.msgEnumToFieldNum[string(attrFd.FullName())]
	if !ok {
		mgr.msgEnumToFieldNum[string(attrFd.FullName())] = make(map[int32]int)
		enumToFields = mgr.msgEnumToFieldNum[string(attrFd.FullName())]
		for i := 0; i < attrFd.Fields().Len(); i++ {
			opt, ok := attrFd.Fields().Get(i).Options().(*descriptorpb.FieldOptions)
			if !ok {
				continue
			}
			enumVal := proto.GetExtension(opt, saipb.E_AttrEnumValue).(int32)
			enumToFields[enumVal] = i
			val, ok := mgr.attrs[id][int32(enumVal)]
			if !ok {
				continue
			}
			msg.ProtoReflect().Set(attrFd.Fields().Get(i), val)
		}
	}

	return nil
}

// storeAttributes stores all the attributes in the message.
func (mgr *AttrMgr) StoreAttributes(id uint64, msg proto.Message) {
	mgr.storeAttributes(fmt.Sprint(id), msg)
}

// GetType returns the SAI type for the object.
func (mgr *AttrMgr) GetType(id string) saipb.ObjectType {
	val, ok := mgr.idToType[id]
	if !ok {
		return saipb.ObjectType_OBJECT_TYPE_NULL
	}
	return val
}

// GetType returns the SAI type for the object.
func (mgr *AttrMgr) SetType(id string, t saipb.ObjectType) {
	mgr.idToType[id] = t
}

// storeAttributes stores all the attributes in the message.
func (mgr *AttrMgr) storeAttributes(id string, msg proto.Message) {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	ty := proto.GetExtension(msg.ProtoReflect().Descriptor().Options(), saipb.E_SaiType).(saipb.ObjectType)
	if ty != saipb.ObjectType_OBJECT_TYPE_UNSPECIFIED {
		mgr.SetType(id, ty)
	}

	msg.ProtoReflect().Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		opt, ok := fd.Options().(*descriptorpb.FieldOptions)
		if !ok {
			return true
		}
		enumValue := proto.GetExtension(opt, saipb.E_AttrEnumValue).(int32)
		if enumValue != 0 {
			mgr.set(id, enumValue, v)
		}
		return true
	})
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
	entry := req.ProtoReflect().Get(req.ProtoReflect().Descriptor().Fields().ByTextName("entry")).Message().Interface()
	pBytes, err := proto.Marshal(entry)
	if err != nil {
		return "", err
	}
	return string(pBytes), nil
}
