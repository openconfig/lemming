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

// Package fault implements the fault RPC service and provides client libraries.
package fault

import (
	"context"
	"sync"
	"time"

	log "github.com/golang/glog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/google/uuid"

	configpb "github.com/openconfig/lemming/proto/config"
	faultpb "github.com/openconfig/lemming/proto/fault"
)

func NewInterceptor() *Interceptor {
	return &Interceptor{
		receivers:        map[string]chan *faultMessage{},
		faultSubs:        map[string]*faultSubscription{},
		configuredFaults: map[string][]*faultpb.FaultMessage{},
	}
}

// NewInterceptorFromConfig creates a new interceptor configured with static faults from config
func NewInterceptorFromConfig(faultConfig *configpb.FaultServiceConfiguration) *Interceptor {
	interceptor := NewInterceptor()

	if faultConfig != nil {
		for _, gnoiFault := range faultConfig.GetGnoiFaults() {
			if err := interceptor.configureFaults(gnoiFault.GetRpcMethod(), gnoiFault.GetFaults()); err != nil {
				log.Warningf("Failed to configure faults for RPC method %s: %v", gnoiFault.GetRpcMethod(), err)
			}
		}
	}

	return interceptor
}

type Interceptor struct {
	faultpb.UnimplementedFaultInjectServer
	faultSubsMu sync.Mutex
	faultSubs   map[string]*faultSubscription
	receiversMu sync.Mutex
	receivers   map[string]chan *faultMessage
	// configuredFaults is a map of RPC methods to a slice of faults.
	configuredFaults map[string][]*faultpb.FaultMessage
	// Mutex for configured faults.
	configMu sync.Mutex
}

type faultMessage struct {
	rpcID, msgID string
	msgType      faultpb.MessageType
	msg          *anypb.Any
	status       *status.Status
}

type faultSubscription struct {
	originMsgCh chan *faultMessage
}

// sendRecvFault sends the msg to fault client and waits for the response.
// msg: the original request received or response from the handler
// oErr: the original error returned from the handler
// returns the (optionally) modified msg and err from the client.
func (i *Interceptor) sendRecvFault(ch chan *faultMessage, rpcID string, msg any, msgType faultpb.MessageType, oErr error) (any, error) {
	var mAny *anypb.Any
	if msg != nil {
		mpb, ok := msg.(proto.Message)
		if !ok { // Do not intercept RPC where the response is not a protobuf message.
			return msg, oErr
		}
		var err error
		mAny, err = anypb.New(mpb)
		if err != nil {
			return msg, oErr
		}
	}

	msgID := uuid.New().String()
	recvCh := make(chan *faultMessage, 1)
	i.receiversMu.Lock()
	i.receivers[msgID] = recvCh
	i.receiversMu.Unlock()

	stErr, _ := status.FromError(oErr)

	ch <- &faultMessage{ // Send the original req to the fault RPC.
		rpcID:   rpcID,
		msgID:   msgID,
		msg:     mAny,
		status:  stErr,
		msgType: msgType,
	}

	var recv *faultMessage
	select { // Receive the potential modified req from the fault RPC.
	case recv = <-recvCh:
	case <-time.After(5 * time.Second):
		log.Infof("timeout waiting for msg %v", msgID)
	}

	i.receiversMu.Lock()
	delete(i.receivers, msgID)
	i.receiversMu.Unlock()

	if recv == nil {
		return msg, oErr
	}
	if recv.msg == nil {
		return nil, recv.status.Err()
	}
	res, err := recv.msg.UnmarshalNew()
	if err != nil {
		return msg, oErr
	}

	return res, recv.status.Err()
}

// configureFaults sets pre-configured faults for a specific RPC method
func (i *Interceptor) configureFaults(rpcMethod string, faults []*faultpb.FaultMessage) error {
	i.configMu.Lock()
	defer i.configMu.Unlock()

	if rpcMethod == "" {
		return status.Errorf(codes.InvalidArgument, "rpc_method cannot be empty")
	}

	if len(faults) == 0 {
		delete(i.configuredFaults, rpcMethod)
		log.Infof("Removed fault configuration for RPC method: %s", rpcMethod)
		return nil
	}

	i.configuredFaults[rpcMethod] = faults
	log.Infof("Configured %d faults for RPC method: %s", len(faults), rpcMethod)
	return nil
}

// nextConfiguredFault returns the next configured fault for an RPC method
func (i *Interceptor) nextConfiguredFault(rpcMethod string) *faultpb.FaultMessage {
	i.configMu.Lock()
	defer i.configMu.Unlock()

	faults, exists := i.configuredFaults[rpcMethod]
	if !exists || len(faults) == 0 {
		return nil
	}

	// Consume the first fault from the slice
	fault := faults[0]
	i.configuredFaults[rpcMethod] = faults[1:]

	return fault
}

// applyConfiguredFault applies a configured fault to the message
func (i *Interceptor) applyConfiguredFault(rpcMethod string, fault *faultpb.FaultMessage, originalMsg any, originalErr error) (any, error) {
	if fault == nil {
		return originalMsg, originalErr
	}

	log.Infof("Applying configured fault for RPC %s: msg_id=%s", rpcMethod, fault.GetMsgId())

	// If fault has a message, use it; otherwise use original.
	if fault.GetMsg() != nil {
		faultMsg, err := fault.GetMsg().UnmarshalNew()
		if err == nil {
			if fault.GetStatus() != nil {
				// If status is OK, return modified message with no error
				if fault.GetStatus().Code == 0 {
					return faultMsg, nil
				}
				// Otherwise return the fault status error
				return faultMsg, status.FromProto(fault.GetStatus()).Err()
			}
			// If no status provided, return modified message with default error
			return faultMsg, status.Errorf(codes.Internal, "configured fault without status for %s", rpcMethod)
		}
		log.Warningf("Failed to unmarshal fault message for %s: %v", rpcMethod, err)
	}

	// If fault only has status, return original message with fault status
	if fault.GetStatus() != nil {
		return originalMsg, status.FromProto(fault.GetStatus()).Err()
	}
	return originalMsg, status.Errorf(codes.Internal, "configured fault without status for %s", rpcMethod)
}

// Unary implements the grpc unary server interceptor interface, adding fault injection to unary RPCs.
func (i *Interceptor) Unary(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	// Check for configured faults first
	if fault := i.nextConfiguredFault(info.FullMethod); fault != nil {
		modifiedReq, faultErr := i.applyConfiguredFault(info.FullMethod, fault, req, nil)

		// If fault has an error status, return immediately without calling handler.
		if faultErr != nil {
			return modifiedReq, faultErr
		}
		req = modifiedReq
	}

	// Check for live fault subscriptions
	i.faultSubsMu.Lock()
	sub, ok := i.faultSubs[info.FullMethod]
	i.faultSubsMu.Unlock()
	if !ok {
		return handler(ctx, req)
	}

	rpcID := uuid.New().String()
	modReq, oErr := i.sendRecvFault(sub.originMsgCh, rpcID, req, faultpb.MessageType_MESSAGE_TYPE_REQUEST, nil)
	if oErr != nil { // If the fault client wants to return an error, don't run the handler and return.
		return modReq, oErr
	}
	res, hErr := handler(ctx, modReq) // Run the implementation of the RPC.

	modResp, err := i.sendRecvFault(sub.originMsgCh, rpcID, res, faultpb.MessageType_MESSAGE_TYPE_RESPONSE, hErr)
	return modResp, err
}

type streamInt struct {
	grpc.ServerStream
	int   *Interceptor
	fs    *faultSubscription
	rpcID string
}

func (si *streamInt) RecvMsg(m any) error {
	err := si.ServerStream.RecvMsg(m)
	msg, err := si.int.sendRecvFault(si.fs.originMsgCh, si.rpcID, m, faultpb.MessageType_MESSAGE_TYPE_REQUEST, err)
	log.Infof("fault stream recv, msg %v, err %v", msg, err)
	if msg != nil {
		if pm, ok := m.(proto.Message); ok {
			proto.Merge(pm, msg.(proto.Message))
		}
	}
	return err
}

func (si *streamInt) SendMsg(m any) error {
	msg, err := si.int.sendRecvFault(si.fs.originMsgCh, si.rpcID, m, faultpb.MessageType_MESSAGE_TYPE_RESPONSE, nil)
	log.Infof("fault stream send, msg %v, err %v", msg, err)
	if err != nil {
		return err
	}
	return si.ServerStream.SendMsg(msg)
}

// Stream implements the grpc stream server interceptor interface, adding fault injection to streaming RPCs.
func (i *Interceptor) Stream(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	if info.FullMethod == "/lemming.fault.FaultInject/Intercept" { // Do not self-intercept
		return handler(srv, stream)
	}

	// Check for configured faults first - for streaming, apply fault immediately
	if fault := i.nextConfiguredFault(info.FullMethod); fault != nil {
		log.Infof("Applying configured stream fault for RPC %s: msg_id=%s", info.FullMethod, fault.GetMsgId())
		if fault.GetStatus() != nil {
			return status.FromProto(fault.GetStatus()).Err()
		}
		return status.Errorf(codes.Internal, "configured fault without status for %s", info.FullMethod)
	}

	// Check for live fault subscriptions
	i.faultSubsMu.Lock()
	sub, ok := i.faultSubs[info.FullMethod]
	i.faultSubsMu.Unlock()
	if !ok {
		return handler(srv, stream)
	}

	si := &streamInt{
		ServerStream: stream,
		int:          i,
		fs:           sub,
		rpcID:        uuid.New().String(),
	}
	hErr := handler(srv, si)
	log.Infof("original stream end, err %v", hErr)
	// After the handler exits, there may be an additional to should be injected.
	_, err := si.int.sendRecvFault(si.fs.originMsgCh, si.rpcID, nil, faultpb.MessageType_MESSAGE_TYPE_STREAM_END, hErr)
	log.Infof("fault stream end, err %v", err)
	return err
}

// Intercept streams RPC requests and responses to the fault client allowing injection of errors.
// Flow:
//  1. Client sends InterceptSubRequest matching the RPC to inject faults.
//  2. Server sends original request or response
//  3. Client replies with optional modified request and error.
//
// Unary RPCs:
//
//	When the client receives a request, it can either reply with a request message and OK status which
//	causes RPC processing to proceed normally, or it can reply with a response and an error bypassing the server implementatiom
//	When the client receives a response, it can reply with the same resposnse or a modified response. It can reply with a non-OK status
//	to inject an into the resposnse.
//
// Streaming RPC:
//
//	Streaming RPC works the same as unary, except if the client wants to inject an error in request,
//	it must supply the request type (NOT reponse)
func (i *Interceptor) Intercept(srv faultpb.FaultInject_InterceptServer) error {
	req, err := srv.Recv()
	if err != nil {
		return err
	}
	if req.GetIntSub() == nil {
		return status.Errorf(codes.InvalidArgument, "expected first request to be rpc filter")
	}
	fs := &faultSubscription{
		originMsgCh: make(chan *faultMessage),
	}

	i.faultSubsMu.Lock()
	if _, ok := i.faultSubs[req.GetIntSub().GetMethod()]; ok {
		i.faultSubsMu.Unlock()
		return status.Errorf(codes.FailedPrecondition, "interceptor already registered for this RPC")
	}
	i.faultSubs[req.GetIntSub().GetMethod()] = fs
	i.faultSubsMu.Unlock()

	defer func() {
		i.faultSubsMu.Lock()
		delete(i.faultSubs, req.GetIntSub().GetMethod())
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
			return srv.Context().Err()
		case err := <-recvErr:
			return err
		case req := <-fs.originMsgCh:
			err := srv.Send(&faultpb.InterceptResponse{
				OriginalMsg: &faultpb.ServerMessage{
					RpcId:   req.rpcID,
					MsgId:   req.msgID,
					MsgType: req.msgType,
					Msg:     req.msg,
					Status:  req.status.Proto(),
				},
			})
			if err != nil {
				return err
			}
		}
	}
}
