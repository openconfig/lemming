// Copyright 2021 Google LLC
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

package gnmit

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"sync"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/prototext"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	gpb "github.com/openconfig/gnmi/proto/gnmi"
	"github.com/openconfig/gnmi/value"
	"github.com/openconfig/lemming/gnmi/internal/oc"
	"github.com/openconfig/ygot/ygot"
)

func mustPath(s string) *gpb.Path {
	p, err := ygot.StringToStructuredPath(s)
	if err != nil {
		panic(fmt.Sprintf("cannot parse subscription path %s, %v", s, err))
	}
	return p
}

// Disable linter for this helper function.
//
//nolint:unparam
func mustTargetPath(t, s string, addOpenConfigOrigin bool) *gpb.Path {
	p := mustPath(s)
	p.Target = t
	if addOpenConfigOrigin {
		p.Origin = "openconfig"
	}
	return p
}

func mustTypedValue(i interface{}) *gpb.TypedValue {
	v, err := value.FromScalar(i)
	if err != nil {
		panic(fmt.Sprintf("cannot parse %v into a TypedValue, %v", i, err))
	}
	return v
}

func mustPathToString(p *gpb.Path) string {
	s, err := ygot.PathToString(p)
	if err != nil {
		panic(fmt.Sprintf("cannot convert %s to path, %v", p, err))
	}
	return s
}

func mustToScalar(t *gpb.TypedValue) interface{} {
	v, err := value.ToScalar(t)
	if err != nil {
		panic(fmt.Sprintf("cannot convert %s to scalar, %v", t, err))
	}
	return v
}

type updateType int64

const (
	_ updateType = iota
	VAL
	SYNC
	DEL
	METACONNECTED
	METASYNC
)

type upd struct {
	T    updateType
	TS   int64
	Path string
	Val  interface{}
}

func (u *upd) String() string {
	b := &bytes.Buffer{}
	b.WriteString("<")
	switch u.T {
	case VAL:
		b.WriteString(fmt.Sprintf("value, @%d %s=%v", u.TS, u.Path, u.Val))
	case METACONNECTED:
		b.WriteString("meta/connected=true")
	case METASYNC:
		b.WriteString("meta/sync=true")
	case SYNC:
		b.WriteString("syncresponse")
	case DEL:
		b.WriteString(fmt.Sprintf("delete @%d %s", u.TS, u.Path))
	}
	b.WriteString(">")
	return b.String()
}

func toUpd(r *gpb.SubscribeResponse) []*upd {
	switch v := r.Response.(type) {
	case *gpb.SubscribeResponse_SyncResponse:
		return append([]*upd{}, &upd{T: SYNC})
	case *gpb.SubscribeResponse_Update:
		ret := []*upd{}
		for _, u := range v.Update.GetUpdate() {
			switch mustPathToString(u.Path) {
			case "/meta/connected":
				ret = append(ret, &upd{
					T: METACONNECTED,
				})
			case "/meta/sync":
				ret = append(ret, &upd{
					T: METASYNC,
				})
			default:
				ret = append(ret, &upd{
					T:    VAL,
					TS:   v.Update.GetTimestamp(),
					Path: mustPathToString(u.Path),
					Val:  mustToScalar(u.Val),
				})
			}
		}
		return ret
	}
	return nil
}

// TestONCE tests the subscribe mode of gnmit.
func TestONCE(t *testing.T) {
	ctx := context.Background()
	c, addr, err := New(ctx, "localhost:0", "local", false, nil)
	if err != nil {
		t.Fatalf("cannot start server, got err: %v", err)
	}

	c.TargetUpdate(&gpb.SubscribeResponse{
		Response: &gpb.SubscribeResponse_Update{
			Update: &gpb.Notification{
				Prefix:    mustTargetPath("local", "", false),
				Timestamp: 42,
				Update: []*gpb.Update{{
					Path: mustPath("/hello"),
					Val:  mustTypedValue("world"),
				}},
			},
		},
	})
	c.TargetUpdate(&gpb.SubscribeResponse{
		Response: &gpb.SubscribeResponse_SyncResponse{
			SyncResponse: true,
		},
	})

	got := []*upd{}
	clientCtx, cancel := context.WithCancel(context.Background())
	var sendErr, recvErr error
	go func(ctx context.Context) {
		defer cancel()
		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			sendErr = fmt.Errorf("cannot dial gNMI server, %v", err)
			return
		}

		client := gpb.NewGNMIClient(conn)
		subc, err := client.Subscribe(ctx)
		if err != nil {
			sendErr = err
			return
		}
		sr := &gpb.SubscribeRequest{
			Request: &gpb.SubscribeRequest_Subscribe{
				Subscribe: &gpb.SubscriptionList{
					Prefix: mustTargetPath("local", "/", false),
					Mode:   gpb.SubscriptionList_ONCE,
					Subscription: []*gpb.Subscription{{
						Path: mustPath("/hello"),
					}},
				},
			},
		}

		if err := subc.Send(sr); err != nil {
			sendErr = fmt.Errorf("cannot send subscribe request %s, %v", prototext.Format(sr), err)
			return
		}

		for {
			in, err := subc.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				recvErr = err
				return
			}
			got = append(got, toUpd(in)...)
		}
	}(clientCtx)

	<-clientCtx.Done()

	c.Stop()

	if sendErr != nil {
		t.Errorf("got unexpected send error, %v", sendErr)
	}

	if recvErr != nil {
		t.Errorf("got unexpected recv error, %v", recvErr)
	}

	if diff := cmp.Diff(got, []*upd{{
		T:    VAL,
		TS:   42,
		Path: "/hello",
		Val:  "world",
	}, {
		T: SYNC,
	}}); diff != "" {
		t.Fatalf("did not get expected updates, diff(-got,+want)\n:%s", diff)
	}
}

// TestONCESet tests the subscribe mode of gnmit.
func TestONCESet(t *testing.T) {
	schema, err := oc.Schema()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	c, addr, err := NewSettable(ctx, "localhost:0", "local", false, schema, nil)
	if err != nil {
		t.Fatalf("cannot start server, got err: %v", err)
	}

	pathStr := "/interfaces/interface[name=foo]/config/description"
	path := mustPath(pathStr)

	got := []*upd{}
	clientCtx, cancel := context.WithCancel(context.Background())
	var sendErr, recvErr error
	go func(ctx context.Context) {
		defer cancel()
		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			sendErr = fmt.Errorf("cannot dial gNMI server, %v", err)
			return
		}

		client := gpb.NewGNMIClient(conn)

		if _, err := client.Set(ctx, &gpb.SetRequest{
			Prefix: mustTargetPath("local", "", true),
			Replace: []*gpb.Update{{
				Path: path,
				Val:  mustTypedValue("world"),
			}},
		}); err != nil {
			sendErr = fmt.Errorf("set request failed: %v", err)
			return
		}

		subc, err := client.Subscribe(ctx)
		if err != nil {
			sendErr = err
			return
		}
		sr := &gpb.SubscribeRequest{
			Request: &gpb.SubscribeRequest_Subscribe{
				Subscribe: &gpb.SubscriptionList{
					Prefix: mustTargetPath("local", "", true),
					Mode:   gpb.SubscriptionList_ONCE,
					Subscription: []*gpb.Subscription{{
						Path: path,
					}},
				},
			},
		}

		if err := subc.Send(sr); err != nil {
			sendErr = fmt.Errorf("cannot send subscribe request %s, %v", prototext.Format(sr), err)
			return
		}

		for {
			in, err := subc.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				recvErr = err
				return
			}
			got = append(got, toUpd(in)...)
		}
	}(clientCtx)

	<-clientCtx.Done()

	c.Stop()

	if sendErr != nil {
		t.Errorf("got unexpected send error, %v", sendErr)
	}

	if recvErr != nil {
		t.Errorf("got unexpected recv error, %v", recvErr)
	}

	if diff := cmp.Diff(got, []*upd{{
		T:    VAL,
		TS:   42,
		Path: pathStr,
		Val:  "world",
	}, {
		T: SYNC,
	}}, cmpopts.IgnoreFields(upd{}, "TS")); diff != "" {
		t.Fatalf("did not get expected updates, diff(-got,+want)\n:%s", diff)
	}
}

// TestONCESet tests the subscribe mode of gnmit.
func TestONCESetJSON(t *testing.T) {
	schema, err := oc.Schema()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	c, addr, err := NewSettable(ctx, "localhost:0", "local", false, schema, nil)
	if err != nil {
		t.Fatalf("cannot start server, got err: %v", err)
	}

	pathStr := "/interfaces/interface[name=foo]/config/description"
	path := mustPath(pathStr)

	got := []*upd{}
	clientCtx, cancel := context.WithCancel(context.Background())
	var sendErr, recvErr error
	go func(ctx context.Context) {
		defer cancel()
		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			sendErr = fmt.Errorf("cannot dial gNMI server, %v", err)
			return
		}

		client := gpb.NewGNMIClient(conn)

		if _, err := client.Set(ctx, &gpb.SetRequest{
			Prefix: mustTargetPath("local", "", true),
			Replace: []*gpb.Update{{
				Path: path,
				Val: &gpb.TypedValue{
					Value: &gpb.TypedValue_JsonIetfVal{
						JsonIetfVal: []byte{0x22, 0x77, 0x6f, 0x72, 0x6c, 0x64, 0x22},
					},
				},
			}},
		}); err != nil {
			sendErr = fmt.Errorf("set request failed: %v", err)
			return
		}

		subc, err := client.Subscribe(ctx)
		if err != nil {
			sendErr = err
			return
		}
		sr := &gpb.SubscribeRequest{
			Request: &gpb.SubscribeRequest_Subscribe{
				Subscribe: &gpb.SubscriptionList{
					Prefix: mustTargetPath("local", "", true),
					Mode:   gpb.SubscriptionList_ONCE,
					Subscription: []*gpb.Subscription{{
						Path: path,
					}},
				},
			},
		}

		if err := subc.Send(sr); err != nil {
			sendErr = fmt.Errorf("cannot send subscribe request %s, %v", prototext.Format(sr), err)
			return
		}

		for {
			in, err := subc.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				recvErr = err
				return
			}
			got = append(got, toUpd(in)...)
		}
	}(clientCtx)

	<-clientCtx.Done()

	c.Stop()

	if sendErr != nil {
		t.Errorf("got unexpected send error, %v", sendErr)
	}

	if recvErr != nil {
		t.Errorf("got unexpected recv error, %v", recvErr)
	}

	if diff := cmp.Diff(got, []*upd{{
		T:    VAL,
		TS:   42,
		Path: pathStr,
		Val:  "world",
	}, {
		T: SYNC,
	}}, cmpopts.IgnoreFields(upd{}, "TS")); diff != "" {
		t.Fatalf("did not get expected updates, diff(-got,+want)\n:%s", diff)
	}
}

// TestSTREAM tests the STREAM mode of gnmit.
func TestSTREAM(t *testing.T) {
	ctx := context.Background()
	c, addr, err := New(ctx, "localhost:0", "local", false, nil)
	if err != nil {
		t.Fatalf("cannot start server, got err: %v", err)
	}

	c.TargetUpdate(&gpb.SubscribeResponse{
		Response: &gpb.SubscribeResponse_Update{
			Update: &gpb.Notification{
				Prefix:    mustTargetPath("local", "", false),
				Timestamp: 42,
				Update: []*gpb.Update{{
					Path: mustPath("/hello"),
					Val:  mustTypedValue("world"),
				}},
			},
		},
	})
	c.TargetUpdate(&gpb.SubscribeResponse{
		Response: &gpb.SubscribeResponse_SyncResponse{
			SyncResponse: true,
		},
	})

	planets := []string{"mercury", "venus", "earth", "mars"}

	var gotMu sync.RWMutex
	got := []*upd{}

	addGot := func(in *gpb.SubscribeResponse) {
		gotMu.Lock()
		defer gotMu.Unlock()
		got = append(got, toUpd(in)...)
	}

	clientCtx, cancel := context.WithCancel(context.Background())
	var sendErr, recvErr error
	var wg sync.WaitGroup
	wg.Add(1)
	go func(ctx context.Context, cfn func()) {
		defer cfn()
		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			sendErr = fmt.Errorf("cannot dial gNMI server, %v", err)
			return
		}

		client := gpb.NewGNMIClient(conn)
		subc, err := client.Subscribe(ctx)
		if err != nil {
			sendErr = err
			return
		}
		sr := &gpb.SubscribeRequest{
			Request: &gpb.SubscribeRequest_Subscribe{
				Subscribe: &gpb.SubscriptionList{
					Prefix: mustTargetPath("local", "", false),
					Mode:   gpb.SubscriptionList_STREAM,
					Subscription: []*gpb.Subscription{{
						Path: mustPath("/"),
						Mode: gpb.SubscriptionMode_TARGET_DEFINED,
					}},
				},
			},
		}
		if err := subc.Send(sr); err != nil {
			sendErr = fmt.Errorf("cannot send subscribe request %s, %v", prototext.Format(sr), err)
			return
		}

		var j int
		for {
			in, err := subc.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				recvErr = err
				return
			}

			addGot(in)

			j++
			if j == len(planets)+4 { // we also get original update, meta/sync and meta/connected + sync_response
				wg.Done()
				return
			}
		}
	}(clientCtx, cancel)

	go func() {
		// time to connect
		time.Sleep(3 * time.Second)
		for i, p := range planets {
			// sleep enough to prevent the cache coalescing
			time.Sleep(1 * time.Second)
			c.TargetUpdate(&gpb.SubscribeResponse{
				Response: &gpb.SubscribeResponse_Update{
					Update: &gpb.Notification{
						Prefix:    mustTargetPath("local", "", false),
						Timestamp: int64(42 + 1 + i),
						Update: []*gpb.Update{{
							Path: mustPath("/hello"),
							Val:  mustTypedValue(p),
						}},
					},
				},
			})
		}
	}()

	<-clientCtx.Done()
	c.Stop()

	if sendErr != nil {
		t.Errorf("got unexpected send error, %v", sendErr)
	}

	if recvErr != nil {
		t.Errorf("got unexpected recv error, %v", recvErr)
	}

	// the semantics of what we need to see here are:
	//  - we need at least one /hello before SYNC
	//  - we need to see all the updates that we expect.
	seenVal := map[string]bool{}
	meta := 0

	wg.Wait()
	for _, s := range got {
		if s.T == SYNC {
			if len(seenVal) < 1 || meta < 1 { // seen hello, may see meta/sync, meta/connected
				t.Fatalf("did not get expected set of updates from client before sync, got: %d %s, want: 3 values, sync (updates %v, meta = %d)", len(got), got, seenVal, meta)
			}
		}
		switch s.T {
		case VAL:
			seenVal[s.Path] = true
		case METACONNECTED, METASYNC:
			meta++
		}
	}

	// now we can check whether we got all updates ignoring order.
	if diff := cmp.Diff(got, []*upd{{
		T: METACONNECTED,
	}, {
		T: METASYNC,
	}, {
		T:    VAL,
		TS:   42,
		Path: "/hello",
		Val:  "world",
	}, {
		T: SYNC,
	}, {
		T:    VAL,
		TS:   43,
		Path: "/hello",
		Val:  "mercury",
	}, {
		T:    VAL,
		TS:   44,
		Path: "/hello",
		Val:  "venus",
	}, {
		T:    VAL,
		TS:   45,
		Path: "/hello",
		Val:  "earth",
	}, {
		T:    VAL,
		TS:   46,
		Path: "/hello",
		Val:  "mars",
	}}, cmpopts.SortSlices(func(a, b *upd) bool {
		if a.T != b.T {
			return a.T < b.T
		}
		if a.TS != b.TS {
			return a.TS < b.TS
		}
		if a.Path != b.Path {
			return a.Path < b.Path
		}
		if a.Val != b.Val {
			return fmt.Sprintf("%v", a.Val) < fmt.Sprintf("%v", b.Val)
		}
		return true
	})); diff != "" {
		t.Fatalf("did not get expected updates, diff(-got,+want)\n:%s", diff)
	}
}
