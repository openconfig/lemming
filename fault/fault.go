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
	"regexp"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/google/uuid"

	faultpb "github.com/openconfig/lemming/proto/fault"
)

func NewInterceptor() *Interceptor {
	return &Interceptor{
		receivers: map[string]chan *faultMessage{},
		faultSubs: map[string]*faultSubscription{},
	}
}

type Interceptor struct {
	faultpb.UnimplementedFaultInjectServer
	faultSubsMu sync.Mutex
	faultSubs   map[string]*faultSubscription
	receiversMu sync.Mutex
	receivers   map[string]chan *faultMessage
}

type faultMessage struct {
	rpcID, msgID string
	method       string
	msg          *anypb.Any
	status       *status.Status
}

type faultSubscription struct {
	exp         *regexp.Regexp
	originMsgCh chan *faultMessage
}

func (i *Interceptor) sendRecvFault(ch chan *faultMessage, rpcID string, msg any, oErr error) (any, error) {
	msgID := uuid.New().String()

	mpb, ok := msg.(proto.Message)
	if !ok { // Do not intercept RPC where the response is not a protobuf message.
		return msg, nil
	}
	mAny, err := anypb.New(mpb)
	if err != nil {
		return msg, oErr
	}

	recvCh := make(chan *faultMessage, 1)
	i.receiversMu.Lock()
	i.receivers[msgID] = recvCh
	i.receiversMu.Unlock()

	stErr, _ := status.FromError(oErr)

	ch <- &faultMessage{ // Send the original req to the fault RPC.
		rpcID:  rpcID,
		msgID:  msgID,
		msg:    mAny,
		status: stErr,
	}

	var recv *faultMessage
	select { // Receive the potential modified req from the fault RPC.
	case recv = <-recvCh:
	case <-time.After(time.Second):
	}
	if recv == nil {
		return msg, oErr
	}

	i.receiversMu.Lock()
	delete(i.receivers, msgID)
	i.receiversMu.Unlock()
	res, err := recv.msg.UnmarshalNew()
	if err != nil {
		return msg, oErr
	}

	return res, recv.status.Err()
}

func (i *Interceptor) Unary(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	rpcID := uuid.New().String()

	i.faultSubsMu.Lock()
	var sub *faultSubscription
	for _, fs := range i.faultSubs {
		if fs.exp.FindString(info.FullMethod) != "" {
			sub = fs
			break
		}
	}
	i.faultSubsMu.Unlock()
	if sub == nil {
		return handler(ctx, req)
	}

	modReq, oErr := i.sendRecvFault(sub.originMsgCh, rpcID, req, nil)
	if oErr != nil { // If the fault client wants to return an error, don't run the handler and return.
		return modReq, oErr
	}
	res, hErr := handler(ctx, modReq) // Run the implementation of the RPC.

	modResp, oErr := i.sendRecvFault(sub.originMsgCh, rpcID, res, hErr)
	return modResp, oErr
}

type streamInt struct {
	grpc.ServerStream
	int   *Interceptor
	fs    *faultSubscription
	rpcID string
}

func (si *streamInt) RecvMsg(m any) error {
	err := si.ServerStream.RecvMsg(m)
	msg, err := si.int.sendRecvFault(si.fs.originMsgCh, si.rpcID, m, err)
	if pm, ok := m.(proto.Message); ok {
		proto.Merge(pm, msg.(proto.Message))
	}

	return err
}

func (si *streamInt) SendMsg(m any) error {
	msg, err := si.int.sendRecvFault(si.fs.originMsgCh, si.rpcID, m, nil)
	if err != nil {
		return err
	}
	return si.ServerStream.SendMsg(msg)
}

func (i *Interceptor) Stream(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	if info.FullMethod == "/lemming.fault.FaultInject/Intercept" { // Do not self-intercept
		return handler(srv, stream)
	}

	i.faultSubsMu.Lock()
	var sub *faultSubscription
	for _, fs := range i.faultSubs {
		if fs.exp.FindString(info.FullMethod) != "" {
			sub = fs
			break
		}
	}
	i.faultSubsMu.Unlock()

	si := &streamInt{
		ServerStream: stream,
		int:          i,
		fs:           sub,
		rpcID:        uuid.New().String(),
	}

	return handler(srv, si)
}

func (i *Interceptor) Intercept(srv faultpb.FaultInject_InterceptServer) error {
	req, err := srv.Recv()
	if err != nil {
		return err
	}
	if req.GetIntSub() == nil {
		return status.Errorf(codes.InvalidArgument, "expected first request to be rpc filter")
	}
	exp, oErr := regexp.Compile(req.GetIntSub().GetMethodRegex())
	if oErr != nil {
		return status.Errorf(codes.InvalidArgument, "invalid regex: %v", oErr)
	}
	fs := &faultSubscription{
		exp:         exp,
		originMsgCh: make(chan *faultMessage),
	}

	intID := uuid.New().String()
	i.faultSubsMu.Lock()
	i.faultSubs[intID] = fs
	i.faultSubsMu.Unlock()

	defer func() {
		i.faultSubsMu.Lock()
		delete(i.faultSubs, intID)
		i.faultSubsMu.Unlock()
	}()

	recvErr := make(chan error, 1)

	go func() {
		for {
			msg, err := srv.Recv()
			if err != nil {
				recvErr <- err
				return
			}
			i.receiversMu.Lock()
			ch, ok := i.receivers[msg.GetFaultMsg().GetMsgId()]
			i.receiversMu.Unlock()
			if ok {
				ch <- &faultMessage{
					msgID:  msg.GetFaultMsg().GetMsgId(),
					msg:    msg.GetFaultMsg().GetMsg(),
					status: status.FromProto(msg.GetFaultMsg().GetStatus()),
				}
			}
		}
	}()

	for {
		select {
		case <-srv.Context().Done():
			return nil
		case err := <-recvErr:
			return err
		case req := <-fs.originMsgCh:
			srv.Send(&faultpb.InterceptResponse{
				OriginalMsg: &faultpb.ServerMessage{
					RpcId:  req.rpcID,
					MsgId:  req.msgID,
					Method: req.method,
					Msg:    req.msg,
				},
			})
		}
	}
}
