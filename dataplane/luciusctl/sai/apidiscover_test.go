// Copyright 2024 Google LLC
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

package sai

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc"
	reflectpb "google.golang.org/grpc/reflection/grpc_reflection_v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

func TestDiscoverRPCs(t *testing.T) {
	tests := []struct {
		desc       string
		fakeClient *fakeClient
		want       map[string]map[string]*method
		wantErr    bool
	}{{
		desc: "success",
		fakeClient: &fakeClient{
			methods: []string{"lemming.dataplane.sai.Switch"},
			descriptors: map[string]*descriptorpb.FileDescriptorProto{
				"lemming.dataplane.sai.Switch": {
					Name:    proto.String("foo.proto"),
					Package: proto.String("lemming.dataplane.sai"),
					EnumType: []*descriptorpb.EnumDescriptorProto{{
						Name: proto.String("SwitchAttr"),
						Value: []*descriptorpb.EnumValueDescriptorProto{{
							Name:   proto.String("SWITCH_ATTR_CPU_PORT"),
							Number: proto.Int32(1),
						}},
					}},
					MessageType: []*descriptorpb.DescriptorProto{{
						Name: proto.String("GetSwitchAttributeRequest"),
						Field: []*descriptorpb.FieldDescriptorProto{{
							Name:   proto.String("oid"),
							Number: proto.Int32(1),
							Type:   descriptorpb.FieldDescriptorProto_TYPE_UINT64.Enum(),
						}, {
							Name:     proto.String("attr_type"),
							Number:   proto.Int32(2),
							TypeName: proto.String(".lemming.dataplane.sai.SwitchAttr"),
							Type:     descriptorpb.FieldDescriptorProto_TYPE_ENUM.Enum(),
							Label:    descriptorpb.FieldDescriptorProto_LABEL_REPEATED.Enum(),
						}},
					}, {
						Name: proto.String("GetSwitchAttributeResponse"),
					}},
					Service: []*descriptorpb.ServiceDescriptorProto{{
						Name: proto.String("Switch"),
						Method: []*descriptorpb.MethodDescriptorProto{{
							Name:       proto.String("GetSwitchAttribute"),
							InputType:  proto.String(".lemming.dataplane.sai.GetSwitchAttributeRequest"),
							OutputType: proto.String(".lemming.dataplane.sai.GetSwitchAttributeResponse"),
						}},
					}},
				},
			},
		},
		want: map[string]map[string]*method{
			"get-attribute": {
				"switch": {
					fullName: "/lemming.dataplane.sai.Switch/GetSwitchAttribute",
					hasOID:   true,
					args:     map[string]*args{"cpu_port": {fullName: "SWITCH_ATTR_CPU_PORT", enumVal: 1}},
				},
			},
		},
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			got, err := discoverRPCs(tt.fakeClient)
			if (err != nil) != tt.wantErr {
				t.Errorf("discoverRPCs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			opts := cmp.Comparer(func(a, b *method) bool {
				if a.fullName != b.fullName {
					return false
				}
				for k, v := range a.args {
					if !reflect.DeepEqual(v, b.args[k]) {
						return false
					}
				}
				for k, v := range b.args {
					if !reflect.DeepEqual(v, a.args[k]) {
						return false
					}
				}
				return true
			})
			if diff := cmp.Diff(tt.want, got, opts); diff != "" {
				t.Errorf("discoverRPCs() (-want, +got):\n%s", diff)
			}
		})
	}
}

type fakeClient struct {
	reflectpb.ServerReflectionClient
	reflectpb.ServerReflection_ServerReflectionInfoClient

	lastRequest *reflectpb.ServerReflectionRequest
	methods     []string
	descriptors map[string]*descriptorpb.FileDescriptorProto
}

func (f *fakeClient) ServerReflectionInfo(context.Context, ...grpc.CallOption) (reflectpb.ServerReflection_ServerReflectionInfoClient, error) {
	return f, nil
}

func (f *fakeClient) Send(req *reflectpb.ServerReflectionRequest) error {
	f.lastRequest = req
	return nil
}

func (f *fakeClient) Recv() (*reflectpb.ServerReflectionResponse, error) {
	switch f.lastRequest.GetMessageRequest().(type) {
	case *reflectpb.ServerReflectionRequest_ListServices:
		svcs := []*reflectpb.ServiceResponse{}
		for _, m := range f.methods {
			svcs = append(svcs, &reflectpb.ServiceResponse{
				Name: m,
			})
		}
		resp := &reflectpb.ServerReflectionResponse{
			MessageResponse: &reflectpb.ServerReflectionResponse_ListServicesResponse{
				ListServicesResponse: &reflectpb.ListServiceResponse{
					Service: svcs,
				},
			},
		}

		return resp, nil
	case *reflectpb.ServerReflectionRequest_FileContainingSymbol:
		if f.descriptors[f.lastRequest.GetFileContainingSymbol()] == nil {
			return nil, fmt.Errorf("not found")
		}
		data, err := proto.Marshal(f.descriptors[f.lastRequest.GetFileContainingSymbol()])
		if err != nil {
			return nil, err
		}
		return &reflectpb.ServerReflectionResponse{
			MessageResponse: &reflectpb.ServerReflectionResponse_FileDescriptorResponse{
				FileDescriptorResponse: &reflectpb.FileDescriptorResponse{
					FileDescriptorProto: [][]byte{data},
				},
			},
		}, nil
	default:
		return nil, fmt.Errorf("unknown type %v", f.lastRequest.GetMessageRequest())
	}
}
