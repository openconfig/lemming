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
	"fmt"
	"net"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/google/go-cmp/cmp"
	"github.com/openconfig/gnmi/errdiff"

	testpb "github.com/openconfig/lemming/fault/proto/test"
	faultpb "github.com/openconfig/lemming/proto/fault"
)

type pingServer struct {
	testpb.UnimplementedPingServer
	gotReq []*testpb.PingRequest
}

const returnErr = "return err"

func (s *pingServer) Unary(_ context.Context, r *testpb.PingRequest) (*testpb.PingResponse, error) {
	s.gotReq = append(s.gotReq, r)
	if r.Msg == returnErr {
		return nil, fmt.Errorf("return err")
	}
	return &testpb.PingResponse{
		Msg: r.GetMsg(),
	}, nil
}

func (s *pingServer) Stream(srv testpb.Ping_StreamServer) error {
	for {
		m, err := srv.Recv()
		s.gotReq = append(s.gotReq, m)
		if err != nil {
			return err
		}
		if m.Msg == returnErr {
			return fmt.Errorf("return err")
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

	unaryPing := func(pc testpb.PingClient, msg string) (string, error) {
		res, err := pc.Unary(context.Background(), &testpb.PingRequest{Msg: msg})
		return res.GetMsg(), err
	}
	streamPing := func(pc testpb.PingClient, msg string) (string, error) {
		sc, err := pc.Stream(context.Background())
		if err != nil {
			return "", err
		}
		if err := sc.Send(&testpb.PingRequest{Msg: msg}); err != nil {
			return "", err
		}
		res, err := sc.Recv()
		return res.GetMsg(), err
	}

	tests := []struct {
		desc      string
		want      string
		wantErr   string
		msg       string
		subReq    *faultpb.InterceptSubRequest
		testRPC   func(pc testpb.PingClient, msg string) (string, error)
		wantReq   []*testpb.PingRequest
		faultMsgs []*faultpb.FaultMessage
	}{{
		desc:    "no match",
		testRPC: unaryPing,
		want:    "hello",
		msg:     "hello",
		wantReq: []*testpb.PingRequest{{Msg: "hello"}},
	}, {
		desc:    "unary modify req and resp",
		testRPC: unaryPing,
		subReq:  &faultpb.InterceptSubRequest{Method: "/test.ping.Ping/Unary"},
		want:    "test2",
		wantReq: []*testpb.PingRequest{{Msg: "test1"}},
		faultMsgs: []*faultpb.FaultMessage{{
			Msg: mustMarshalAny(t, &testpb.PingRequest{Msg: "test1"}),
		}, {
			Msg: mustMarshalAny(t, &testpb.PingResponse{Msg: "test2"}),
		}},
	}, {
		desc:    "unary modify req with err",
		testRPC: unaryPing,
		subReq:  &faultpb.InterceptSubRequest{Method: "/test.ping.Ping/Unary"},
		wantErr: "error message",
		msg:     "hello",
		faultMsgs: []*faultpb.FaultMessage{{
			Msg:    mustMarshalAny(t, &testpb.PingResponse{Msg: "test"}),
			Status: status.New(codes.Internal, "error message").Proto(),
		}},
	}, {
		desc:    "unary modify resp with err",
		testRPC: unaryPing,
		subReq:  &faultpb.InterceptSubRequest{Method: "/test.ping.Ping/Unary"},
		wantErr: "error message",
		wantReq: []*testpb.PingRequest{{Msg: "test"}},
		faultMsgs: []*faultpb.FaultMessage{{
			Msg: mustMarshalAny(t, &testpb.PingRequest{Msg: "test"}),
		}, {
			Msg:    mustMarshalAny(t, &testpb.PingResponse{Msg: "test"}),
			Status: status.New(codes.Internal, "error message").Proto(),
		}},
	}, {
		desc:    "unary override handler err",
		testRPC: unaryPing,
		subReq:  &faultpb.InterceptSubRequest{Method: "/test.ping.Ping/Unary"},
		wantErr: "error message",
		wantReq: []*testpb.PingRequest{{Msg: returnErr}},
		faultMsgs: []*faultpb.FaultMessage{{
			Msg: mustMarshalAny(t, &testpb.PingRequest{Msg: returnErr}),
		}, {
			Msg:    mustMarshalAny(t, &testpb.PingResponse{Msg: "test"}),
			Status: status.New(codes.Internal, "error message").Proto(),
		}},
	}, {
		desc:    "stream modify req and resp",
		testRPC: streamPing,
		subReq:  &faultpb.InterceptSubRequest{Method: "/test.ping.Ping/Stream"},
		want:    "test2",
		msg:     "hello",
		wantReq: []*testpb.PingRequest{{Msg: "test1"}},
		faultMsgs: []*faultpb.FaultMessage{{
			Msg: mustMarshalAny(t, &testpb.PingRequest{Msg: "test1"}),
		}, {
			Msg: mustMarshalAny(t, &testpb.PingResponse{Msg: "test2"}),
		}, {}},
	}, {
		desc:    "stream modify req with err",
		testRPC: streamPing,
		subReq:  &faultpb.InterceptSubRequest{Method: "/test.ping.Ping/Stream"},
		wantErr: "error message",
		wantReq: []*testpb.PingRequest{nil}, // The server calls Recv, which returns the injected error and no value.
		faultMsgs: []*faultpb.FaultMessage{{
			Msg:    mustMarshalAny(t, &testpb.PingRequest{Msg: "test"}),
			Status: status.New(codes.Internal, "error message").Proto(),
		}, {}},
	}, {
		desc:    "stream modify resp with err",
		testRPC: streamPing,
		subReq:  &faultpb.InterceptSubRequest{Method: "/test.ping.Ping/Stream"},
		wantErr: "error message",
		wantReq: []*testpb.PingRequest{{Msg: "test"}},
		faultMsgs: []*faultpb.FaultMessage{{
			Msg: mustMarshalAny(t, &testpb.PingRequest{Msg: "test"}),
		}, {
			Msg:    mustMarshalAny(t, &testpb.PingResponse{Msg: "test"}),
			Status: status.New(codes.Internal, "error message").Proto(),
		}, {}},
	}, {
		desc:    "stream override handler err",
		testRPC: streamPing,
		subReq:  &faultpb.InterceptSubRequest{Method: "/test.ping.Ping/Stream"},
		wantErr: "error message",
		wantReq: []*testpb.PingRequest{{Msg: returnErr}},
		faultMsgs: []*faultpb.FaultMessage{{
			Msg: mustMarshalAny(t, &testpb.PingRequest{Msg: returnErr}),
		}, {
			Msg:    mustMarshalAny(t, &testpb.PingResponse{Msg: "test"}),
			Status: status.New(codes.Internal, "error message").Proto(),
		}},
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			ps.gotReq = nil
			ctx, cancelFn := context.WithCancel(context.Background())
			defer cancelFn()
			time.Sleep(10 * time.Millisecond) // Sleep long enough for previous Intercept to close.
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
			got, gotErr := tt.testRPC(pc, tt.msg)
			if diff := errdiff.Check(gotErr, tt.wantErr); diff != "" {
				t.Errorf("unexpected err: %s", diff)
			}
			if d := cmp.Diff(got, tt.want); d != "" {
				t.Errorf("unexpected result: diff(-got,+want)\n:%s", d)
			}
			if d := cmp.Diff(ps.gotReq, tt.wantReq, protocmp.Transform()); d != "" {
				t.Errorf("unexpected input requests: diff(-got,+want)\n:%s", d)
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
