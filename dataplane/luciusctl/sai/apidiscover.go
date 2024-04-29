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
	"strings"

	reflectpb "google.golang.org/grpc/reflection/grpc_reflection_v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

type args struct {
	fullName string
	enumVal  protoreflect.EnumNumber
}

type method struct {
	fullName string
	input    protoreflect.MessageDescriptor
	output   protoreflect.MessageDescriptor
	args     map[string]*args
	hasOID   bool
}

// discoverRPCs uses gRPC reflection to find all the SAI APIs on the server.
// It returns a CLI-friendly names for the RPCs and their attributes.
func discoverRPCs(refl reflectpb.ServerReflectionClient) (map[string]map[string]*method, error) {
	methods := map[string]map[string]*method{
		getAttr: {},
	}

	info, err := refl.ServerReflectionInfo(context.Background())
	if err != nil {
		return methods, err
	}
	err = info.Send(&reflectpb.ServerReflectionRequest{
		MessageRequest: &reflectpb.ServerReflectionRequest_ListServices{
			ListServices: "foo",
		},
	})
	if err != nil {
		return methods, err
	}
	resp, err := info.Recv()
	if err != nil {
		return methods, err
	}
	set := &descriptorpb.FileDescriptorSet{}
	seenFiles := map[string]struct{}{}
	for _, svc := range resp.GetListServicesResponse().GetService() {
		if !strings.HasPrefix(svc.GetName(), saiPrefix) {
			continue
		}
		err := info.Send(&reflectpb.ServerReflectionRequest{
			MessageRequest: &reflectpb.ServerReflectionRequest_FileContainingSymbol{
				FileContainingSymbol: svc.GetName(),
			},
		})
		if err != nil {
			return methods, err
		}
		resp, err := info.Recv()
		if err != nil {
			return methods, err
		}

		for _, descByte := range resp.GetFileDescriptorResponse().GetFileDescriptorProto() {
			file := &descriptorpb.FileDescriptorProto{}
			if err := proto.Unmarshal(descByte, file); err != nil {
				return methods, err
			}
			if _, ok := seenFiles[file.GetName()]; !ok {
				set.File = append(set.File, file)
				seenFiles[file.GetName()] = struct{}{}
			}
		}
	}
	reg, err := protodesc.NewFiles(set)
	if err != nil {
		return nil, err
	}
	reg.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		for svc := 0; svc < fd.Services().Len(); svc++ {
			for m := 0; m < fd.Services().Get(svc).Methods().Len(); m++ {
				methodpb := fd.Services().Get(svc).Methods().Get(m)
				switch { // TODO: support other methods
				case strings.HasPrefix(string(methodpb.Name()), "Get") && strings.HasSuffix(string(methodpb.Name()), "Attribute"):
					shortName := strings.TrimPrefix(strings.TrimSuffix(string(methodpb.Name()), "Attribute"), "Get")
					arg := map[string]*args{}
					for i := 0; i < methodpb.Input().Fields().ByNumber(2).Enum().Values().Len(); i++ {
						name := methodpb.Input().Fields().ByNumber(2).Enum().Values().Get(i).Name()
						_, argShortName, _ := strings.Cut(string(name), "ATTR_")
						arg[strings.ToLower(argShortName)] = &args{
							fullName: string(name),
							enumVal:  methodpb.Input().Fields().ByNumber(2).Enum().Values().Get(i).Number(),
						}
					}
					methods[getAttr][strings.ToLower(shortName)] = &method{
						fullName: fmt.Sprintf("/%s.%s/%s", fd.Package(), fd.Services().Get(svc).Name(), methodpb.Name()),
						input:    methodpb.Input(),
						output:   methodpb.Output(),
						args:     arg,
					}
					if methodpb.Input().Fields().ByNumber(1).Kind() == protoreflect.Uint64Kind {
						methods[getAttr][strings.ToLower(shortName)].hasOID = true
					}
				}
			}
		}

		return true
	})
	return methods, err
}
