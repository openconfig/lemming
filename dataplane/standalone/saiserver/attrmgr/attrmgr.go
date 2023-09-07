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

package attrmgr

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"

	saipb "github.com/openconfig/lemming/dataplane/standalone/proto"
)

// AttrMgr stores and retrieve the SAI attributes.
type AttrMgr struct {
	attrsMu sync.RWMutex
	attrs   map[string]map[int32]protoreflect.Value
	nextOid atomic.Uint64
}

// New returns a new AttrMgr.
func New() *AttrMgr {
	mgr := &AttrMgr{
		attrs: make(map[string]map[int32]protoreflect.Value),
	}
	return mgr
}

func (mgr *AttrMgr) set(id string, attr int32, val protoreflect.Value) {
	mgr.attrsMu.Lock()
	defer mgr.attrsMu.Unlock()
	if _, ok := mgr.attrs[id]; !ok {
		mgr.attrs[id] = make(map[int32]protoreflect.Value)
	}
	mgr.attrs[id][attr] = val
}

var (
	editRPC = regexp.MustCompile("Create|Set")
	readRPC = regexp.MustCompile("Get")
)

const protoNS = "lemming.dataplane.sai"

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
	if resp == nil {
		respName := strings.ReplaceAll(string(reqMsg.ProtoReflect().Descriptor().FullName()), "Request", "Response")
		pType, err := protoregistry.GlobalTypes.FindMessageByName(protoreflect.FullName(respName))
		if err != nil {
			return nil, err
		}
		resp = pType.New().Interface()
	}
	respMsg := resp.(proto.Message)
	if res := editRPC.FindString(info.FullMethod); res != "" {
		id, err := mgr.getID(reqMsg, respMsg)
		if err != nil {
			return nil, err
		}
		mgr.storeAttributes(id, reqMsg)
	} else if res := readRPC.FindString(info.FullMethod); res != "" {
		id, err := mgr.getID(reqMsg, respMsg)
		if err != nil {
			return nil, err
		}
		mgr.populateAttributes(id, reqMsg, respMsg)
	}
	return resp, nil
}

// populateAttributes fills the resp with the requests attributes.
func (mgr *AttrMgr) populateAttributes(id string, req, resp proto.Message) {
	mgr.attrsMu.RLock()
	defer mgr.attrsMu.RUnlock()

	attrs := resp.ProtoReflect().Mutable(resp.ProtoReflect().Descriptor().Fields().ByTextName("attr")).Message()
	reqList := req.ProtoReflect().Get(req.ProtoReflect().Descriptor().Fields().ByTextName("attr_type")).List()

	for i := 0; i < reqList.Len(); i++ {
		enumVal := reqList.Get(0).Enum()
		val, ok := mgr.attrs[id][int32(enumVal)]
		if !ok {
			continue
		}
		attrs.Set(attrs.Descriptor().Fields().ByNumber(protowire.Number(enumVal)), val)
	}
}

// storeAttributes stores all the attributes in the message.
func (mgr *AttrMgr) storeAttributes(id string, msg proto.Message) {
	msg.ProtoReflect().Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		opt, ok := fd.Options().(*descriptorpb.FieldOptions)
		if !ok {
			return true
		}
		enumValue := proto.GetExtension(opt, saipb.E_AttrEnumValue).(int32)
		mgr.set(id, enumValue, v)
		return true
	})
}

// getID returns the id from either the request or response. If the id is unset, it allocates a new one.
func (mgr *AttrMgr) getID(req, resp proto.Message) (string, error) {
	msgs := []proto.Message{req, resp}
	for _, msg := range msgs {
		if fd := msg.ProtoReflect().Descriptor().Fields().ByTextName("oid"); fd != nil {
			v := msg.ProtoReflect().Get(fd).Uint()
			if v == 0 {
				msg.ProtoReflect().Set(fd, protoreflect.ValueOfUint64(uint64(mgr.nextOid.Add(1))))
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
