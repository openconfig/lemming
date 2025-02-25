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
	"net"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/google/go-cmp/cmp"
	"github.com/openconfig/gnmi/errdiff"

	testpb "github.com/openconfig/lemming/fault/proto/test"
	faultpb "github.com/openconfig/lemming/proto/fault"
)

type pingServer struct {
	testpb.UnimplementedPingServer
}

func (s *pingServer) Unary(_ context.Context, r *testpb.PingRequest) (*testpb.PingResponse, error) {
	return &testpb.PingResponse{
		Msg: r.GetMsg(),
	}, nil
}

func (s *pingServer) Stream(srv testpb.Ping_StreamServer) error {
	for {
		m, err := srv.Recv()
		if err != nil {
			return err
		}
		err = srv.Send(&testpb.PingResponse{Msg: m.GetMsg()})
		if err != nil {
			return err
		}
	}
}

func TestIntercept(t *testing.T) {
	int := NewInterceptor()
	ps := &pingServer{}

	srv := grpc.NewServer(grpc.UnaryInterceptor(int.Unary), grpc.StreamInterceptor(int.Stream))
	testpb.RegisterPingServer(srv, ps)
	faultpb.RegisterFaultInjectServer(srv, int)

	l, err := net.Listen("tcp", "localhost:")
	if err != nil {
		t.Fatalf("listen failed: %v", err)
	}

	go srv.Serve(l)
	defer srv.Stop()

	c, err := grpc.NewClient(l.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}

	pc := testpb.NewPingClient(c)
	fc := faultpb.NewFaultInjectClient(c)

	unaryPing := func(pc testpb.PingClient) (string, error) {
		res, err := pc.Unary(context.Background(), &testpb.PingRequest{Msg: "hello"})
		return res.GetMsg(), err
	}
	streamPing := func(pc testpb.PingClient) (string, error) {
		sc, err := pc.Stream(context.Background())
		if err != nil {
			return "", err
		}
		if err := sc.Send(&testpb.PingRequest{Msg: "hello"}); err != nil {
			return "", err
		}
		res, err := sc.Recv()
		return res.GetMsg(), err
	}

	tests := []struct {
		desc      string
		want      string
		wantErr   string
		subReq    *faultpb.InterceptSubRequest
		testRPC   func(pc testpb.PingClient) (string, error)
		faultMsgs []*faultpb.FaultMessage
	}{{
		desc:    "no match",
		testRPC: unaryPing,
		want:    "hello",
	}, {
		desc:    "unary modify req and resp",
		testRPC: unaryPing,
		subReq:  &faultpb.InterceptSubRequest{MethodRegex: ".*"},
		want:    "test2",
		faultMsgs: []*faultpb.FaultMessage{{
			Msg: mustMarshalAny(t, &testpb.PingRequest{Msg: "test1"}),
		}, {
			Msg: mustMarshalAny(t, &testpb.PingResponse{Msg: "test2"}),
		}},
	}, {
		desc:    "unary modify req with err",
		testRPC: unaryPing,
		subReq:  &faultpb.InterceptSubRequest{MethodRegex: ".*"},
		wantErr: "error message",
		faultMsgs: []*faultpb.FaultMessage{{
			Msg:    mustMarshalAny(t, &testpb.PingResponse{Msg: "test"}),
			Status: status.New(codes.Internal, "error message").Proto(),
		}},
	}, {
		desc:    "unary modify resp with err",
		testRPC: unaryPing,
		subReq:  &faultpb.InterceptSubRequest{MethodRegex: ".*"},
		wantErr: "error message",
		faultMsgs: []*faultpb.FaultMessage{{
			Msg: mustMarshalAny(t, &testpb.PingRequest{Msg: "test"}),
		}, {
			Msg:    mustMarshalAny(t, &testpb.PingResponse{Msg: "test"}),
			Status: status.New(codes.Internal, "error message").Proto(),
		}},
	}, {
		desc:    "stream modify req and resp",
		testRPC: streamPing,
		subReq:  &faultpb.InterceptSubRequest{MethodRegex: ".*"},
		want:    "test2",
		faultMsgs: []*faultpb.FaultMessage{{
			Msg: mustMarshalAny(t, &testpb.PingRequest{Msg: "test1"}),
		}, {
			Msg: mustMarshalAny(t, &testpb.PingResponse{Msg: "test2"}),
		}},
	}, {
		desc:    "stream modify req with err",
		testRPC: streamPing,
		subReq:  &faultpb.InterceptSubRequest{MethodRegex: ".*"},
		wantErr: "error message",
		faultMsgs: []*faultpb.FaultMessage{{
			Msg:    mustMarshalAny(t, &testpb.PingRequest{Msg: "test"}),
			Status: status.New(codes.Internal, "error message").Proto(),
		}},
	}, {
		desc:    "stream modify resp with err",
		testRPC: streamPing,
		subReq:  &faultpb.InterceptSubRequest{MethodRegex: ".*"},
		wantErr: "error message",
		faultMsgs: []*faultpb.FaultMessage{{
			Msg: mustMarshalAny(t, &testpb.PingRequest{Msg: "test"}),
		}, {
			Msg:    mustMarshalAny(t, &testpb.PingResponse{Msg: "test"}),
			Status: status.New(codes.Internal, "error message").Proto(),
		}},
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			ctx, cancelFn := context.WithCancel(context.Background())
			defer cancelFn()
			if tt.subReq != nil {
				sub, err := fc.Intercept(ctx)
				if err != nil {
					t.Fatalf("intercept failed: %v", err)
				}
				if err := sub.Send(&faultpb.InterceptRequest{Msg: &faultpb.InterceptRequest_IntSub{IntSub: tt.subReq}}); err != nil {
					t.Fatalf("interceptor sub failed: %v", err)
				}
				go func() {
					for _, m := range tt.faultMsgs {
						oReq, _ := sub.Recv()
						sub.Send(&faultpb.InterceptRequest{Msg: &faultpb.InterceptRequest_FaultMsg{
							FaultMsg: &faultpb.FaultMessage{
								MsgId:  oReq.GetOriginalMsg().GetMsgId(),
								Msg:    m.GetMsg(),
								Status: m.GetStatus(),
							},
						}})
					}
				}()
			}
			time.Sleep(10 * time.Millisecond) // Sleep long enough for the 1st fault RPC to send.
			got, gotErr := tt.testRPC(pc)
			if diff := errdiff.Check(gotErr, tt.wantErr); diff != "" {
				t.Fatalf("unexpected err: %s", diff)
			}
			if gotErr != nil {
				return
			}
			if d := cmp.Diff(got, tt.want); d != "" {
				t.Errorf("unexpect result: diff(-got,+want)\n:%s", d)
			}
		})
	}
}

func mustMarshalAny(t testing.TB, m proto.Message) *anypb.Any {
	t.Helper()
	a, err := anypb.New(m)
	if err != nil {
		t.Fatal(err)
	}
	return a
}
