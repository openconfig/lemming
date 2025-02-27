// Copyright 2025 Google LLC
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

package fault

import (
	"context"

	"google.golang.org/grpc"

	gpb "github.com/openconfig/gnmi/proto/gnmi"

	faultpb "github.com/openconfig/lemming/proto/fault"
)

func NewClient(conn grpc.ClientConnInterface) *Client {
	return &Client{
		fc: faultpb.NewFaultInjectClient(conn),
	}
}

type Client struct {
	fc faultpb.FaultInjectClient
}

func (c *Client) GNMISubscribe() error {
	s, err := c.fc.Intercept(context.TODO())
	if err != nil {
		return err
	}
	s.Send(&faultpb.InterceptRequest{Msg: &faultpb.InterceptRequest_IntSub{
		IntSub: &faultpb.InterceptSubRequest{
			Method: "/gnmi.gNMI/Subscribe",
		},
	}})

	go func() {
		for {
			msg, err := s.Recv()
			if err != nil {
				return
			}
			switch msg.GetOriginalMsg().MsgType {
			case faultpb.MessageType_MESSAGE_TYPE_REQUEST:
				req := &gpb.SubscribeRequest{}
				err := msg.GetOriginalMsg().GetMsg().UnmarshalTo(req)
				if err != nil {
					return
				}
			case faultpb.MessageType_MESSAGE_TYPE_RESPONSE:
				req := &gpb.SubscribeResponse{}
				err := msg.GetOriginalMsg().GetMsg().UnmarshalTo(req)
				if err != nil {
					return
				}
			}
		}
	}()

	return nil
}

func (c *Client) GNMISet() error {
	s, err := c.fc.Intercept(context.TODO())
	if err != nil {
		return err
	}
	s.Send(&faultpb.InterceptRequest{Msg: &faultpb.InterceptRequest_IntSub{
		IntSub: &faultpb.InterceptSubRequest{
			Method: "/gnmi.gNMI/Set",
		},
	}})
	return nil
}
