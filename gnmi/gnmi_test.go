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

package gnmi

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/h-fam/errdiff"
	"github.com/openconfig/gnmi/value"
	"github.com/openconfig/lemming/gnmi/gnmiclient"
	"github.com/openconfig/lemming/gnmi/oc"
	"github.com/openconfig/lemming/gnmi/oc/ocpath"
	"github.com/openconfig/lemming/internal/creds"
	"github.com/openconfig/ygnmi/ygnmi"
	"github.com/openconfig/ygot/ygot"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/local"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/prototext"

	gpb "github.com/openconfig/gnmi/proto/gnmi"
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
		p.Origin = OpenConfigOrigin
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

// startServer starts the collector-backed gNMI server that listens on the specified
// addr (in the form host:port).
//
// It returns the address it is listening on in the form hostname:port or any
// errors encounted whilst setting it up.
func startServer(s *Server, opts ...grpc.ServerOption) (string, error) {
	// Start gNMI server.
	srv := grpc.NewServer(opts...)
	gpb.RegisterGNMIServer(srv, s)
	// Forward streaming updates to clients.
	// Register listening port and start serving.
	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return "", fmt.Errorf("failed to listen: %v", err)
	}

	go func() {
		if err := srv.Serve(lis); err != nil {
			log.Fatalf("Error while serving gnmi target: %v", err)
		}
	}()
	s.c.stopFn = func() {
		srv.GracefulStop()
	}
	return lis.Addr().String(), nil
}

// TestONCE tests the subscribe mode of gnmit.
func TestONCE(t *testing.T) {
	ctx := context.Background()
	gnmiServer, err := newServer(ctx, "local", false)
	if err != nil {
		t.Fatalf("cannot create server, got err: %v", err)
	}
	addr, err := startServer(gnmiServer)
	if err != nil {
		t.Fatalf("cannot start server, got err: %v", err)
	}

	gnmiServer.c.TargetUpdate(&gpb.SubscribeResponse{
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
	gnmiServer.c.TargetUpdate(&gpb.SubscribeResponse{
		Response: &gpb.SubscribeResponse_SyncResponse{
			SyncResponse: true,
		},
	})

	got := []*upd{}
	clientCtx, cancel := context.WithCancel(context.Background())
	var sendErr, recvErr error
	go func(ctx context.Context) {
		defer cancel()
		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(local.NewCredentials()))
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

	gnmiServer.c.Stop()

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

// TestSetConfig tests gnmi.Set on a config value.
//
// It purposely avoids using ygnmi in order to test lower-level details
// (e.g. timestamp metadata)
func TestSetConfig(t *testing.T) {
	ctx := context.Background()
	gnmiServer, err := newServer(ctx, "local", true)
	if err != nil {
		t.Fatalf("cannot create server, got err: %v", err)
	}
	addr, err := startServer(gnmiServer)
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
		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(local.NewCredentials()))
		if err != nil {
			sendErr = fmt.Errorf("cannot dial gNMI server, %v", err)
			return
		}

		client := gpb.NewGNMIClient(conn)

		if _, err := client.Set(metadata.AppendToOutgoingContext(ctx, TimestampMetadataKey, strconv.FormatInt(42, 10)), &gpb.SetRequest{
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

	gnmiServer.c.Stop()

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

	// Test that timestamp is not 42: we don't want the timestamp metadata to affect config values.
	if cmp.Equal(got, []*upd{{
		T:    VAL,
		TS:   42,
		Path: pathStr,
		Val:  "world",
	}, {
		T: SYNC,
	}}) {
		t.Fatalf("Expected error -- timestamp metadata should be ignored but it is not ignored.")
	}
}

// TestSetState tests gnmi.Set on a state value.
//
// It purposely avoids using ygnmi in order to test lower-level details
// (e.g. timestamp metadata)
func TestSetState(t *testing.T) {
	ctx := context.Background()
	gnmiServer, err := newServer(ctx, "local", true)
	if err != nil {
		t.Fatalf("cannot create server, got err: %v", err)
	}
	addr, err := startServer(gnmiServer)
	if err != nil {
		t.Fatalf("cannot start server, got err: %v", err)
	}

	pathStr := "/interfaces/interface[name=foo]/state/description"
	path := mustPath(pathStr)

	got := []*upd{}
	clientCtx, cancel := context.WithCancel(context.Background())
	var sendErr, recvErr error
	go func(ctx context.Context) {
		defer cancel()
		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(local.NewCredentials()))
		if err != nil {
			sendErr = fmt.Errorf("cannot dial gNMI server, %v", err)
			return
		}

		client := gpb.NewGNMIClient(conn)

		ctx = metadata.AppendToOutgoingContext(ctx, GNMIModeMetadataKey, string(StateMode))
		if _, err := client.Set(metadata.AppendToOutgoingContext(ctx, TimestampMetadataKey, strconv.FormatInt(42, 10)), &gpb.SetRequest{
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

	gnmiServer.c.Stop()

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
	}}); diff != "" {
		t.Fatalf("did not get expected updates, diff(-got,+want)\n:%s", diff)
	}
}

// TestSetInternal tests that the server is able to handle schemaless queries.
func TestSetInternal(t *testing.T) {
	ctx := context.Background()
	gnmiServer, err := newServer(ctx, "local", true)
	if err != nil {
		t.Fatalf("cannot create server, got err: %v", err)
	}
	addr, err := startServer(gnmiServer)
	if err != nil {
		t.Fatalf("cannot start server, got err: %v", err)
	}

	pathStr := "/test/foo"
	path := mustPath(pathStr)
	path.Origin = InternalOrigin

	got := []*upd{}
	clientCtx, cancel := context.WithCancel(context.Background())
	var sendErr, recvErr error
	go func(ctx context.Context) {
		defer cancel()
		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(local.NewCredentials()))
		if err != nil {
			sendErr = fmt.Errorf("cannot dial gNMI server, %v", err)
			return
		}

		client := gpb.NewGNMIClient(conn)

		if _, err := client.Set(metadata.AppendToOutgoingContext(ctx, TimestampMetadataKey, strconv.FormatInt(42, 10)), &gpb.SetRequest{
			Prefix: mustTargetPath("local", "", false),
			Replace: []*gpb.Update{{
				Path: path,
				Val:  &gpb.TypedValue{Value: &gpb.TypedValue_StringVal{StringVal: "test"}},
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
					Prefix: mustTargetPath("local", "", false),
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

	gnmiServer.c.Stop()

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
		Val:  "test",
	}, {
		T: SYNC,
	}}, cmpopts.IgnoreFields(upd{}, "TS")); diff != "" {
		t.Fatalf("did not get expected updates, diff(-got,+want)\n:%s", diff)
	}

	// Test that timestamp is not 42: we don't want the timestamp metadata to affect config values.
	if cmp.Equal(got, []*upd{{
		T:    VAL,
		TS:   42,
		Path: pathStr,
		Val:  "world",
	}, {
		T: SYNC,
	}}) {
		t.Fatalf("Expected error -- timestamp metadata should be ignored but it is not ignored.")
	}
}

func TestSetYGNMI(t *testing.T) {
	tests := []struct {
		desc    string
		isState bool
		inOp    func(c *ygnmi.Client) error
		checkOp func(t *testing.T, c *ygnmi.Client)
		wantErr string
	}{{
		desc: "leaf config update",
		inOp: func(c *ygnmi.Client) error {
			_, err := ygnmi.Update(context.Background(), c, ocpath.Root().System().Hostname().Config(), "foo")
			return err
		},
		checkOp: func(t *testing.T, c *ygnmi.Client) {
			v, err := ygnmi.Get[string](context.Background(), c, ocpath.Root().System().Hostname().Config())
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff("foo", v); diff != "" {
				t.Errorf("(-want, +got):\n%s", diff)
			}
		},
	}, {
		desc: "leaf config replace",
		inOp: func(c *ygnmi.Client) error {
			_, err := ygnmi.Replace(context.Background(), c, ocpath.Root().System().Hostname().Config(), "foo")
			return err
		},
		checkOp: func(t *testing.T, c *ygnmi.Client) {
			v, err := ygnmi.Get[string](context.Background(), c, ocpath.Root().System().Hostname().Config())
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff("foo", v); diff != "" {
				t.Errorf("(-want, +got):\n%s", diff)
			}
		},
	}, {
		desc: "leaf config delete",
		inOp: func(c *ygnmi.Client) error {
			if _, err := ygnmi.Update(context.Background(), c, ocpath.Root().System().Hostname().Config(), "foo"); err != nil {
				return err
			}
			_, err := ygnmi.Delete(context.Background(), c, ocpath.Root().System().Hostname().Config())
			return err
		},
		checkOp: func(t *testing.T, c *ygnmi.Client) {
			v, err := ygnmi.Lookup[string](context.Background(), c, ocpath.Root().System().Hostname().Config())
			if err != nil {
				t.Fatal(err)
			}
			if v.IsPresent() {
				t.Errorf("Got present, want not present")
			}
		},
	}, {
		desc: "non-leaf config update",
		inOp: func(c *ygnmi.Client) error {
			if _, err := ygnmi.Update(context.Background(), c, ocpath.Root().System().Config(), &oc.System{Hostname: ygot.String("foo")}); err != nil {
				return err
			}
			_, err := gnmiclient.Update[*oc.System](context.Background(), c, ocpath.Root().System().Config(), &oc.System{MotdBanner: ygot.String("bar")})
			return err
		},
		checkOp: func(t *testing.T, c *ygnmi.Client) {
			v, err := ygnmi.Get[*oc.System](context.Background(), c, ocpath.Root().System().Config())
			if err != nil {
				t.Fatal(err)
			}
			nos, err := ygot.Diff(&oc.System{Hostname: ygot.String("foo"), MotdBanner: ygot.String("bar")}, v)
			if err != nil {
				t.Fatal(err)
			}
			if len(nos.Update)+len(nos.Delete) != 0 {
				t.Errorf("Got diff:\n%s", nos)
			}
		},
	}, {
		desc: "non-leaf config replace",
		inOp: func(c *ygnmi.Client) error {
			if _, err := ygnmi.Update(context.Background(), c, ocpath.Root().System().Config(), &oc.System{Hostname: ygot.String("foo")}); err != nil {
				return err
			}
			_, err := ygnmi.Replace(context.Background(), c, ocpath.Root().System().Config(), &oc.System{MotdBanner: ygot.String("foo")})
			return err
		},
		checkOp: func(t *testing.T, c *ygnmi.Client) {
			v, err := ygnmi.Get[*oc.System](context.Background(), c, ocpath.Root().System().Config())
			if err != nil {
				t.Fatal(err)
			}
			nos, err := ygot.Diff(&oc.System{MotdBanner: ygot.String("foo")}, v)
			if err != nil {
				t.Fatal(err)
			}
			if len(nos.Update)+len(nos.Delete) != 0 {
				t.Errorf("Got diff:\n%s", nos)
			}
		},
	}, {
		desc: "non-leaf config delete",
		inOp: func(c *ygnmi.Client) error {
			if _, err := ygnmi.Update(context.Background(), c, ocpath.Root().System().Config(), &oc.System{Hostname: ygot.String("foo")}); err != nil {
				return err
			}
			_, err := ygnmi.Delete(context.Background(), c, ocpath.Root().System().Config())
			return err
		},
		checkOp: func(t *testing.T, c *ygnmi.Client) {
			v, err := ygnmi.Lookup[*oc.System](context.Background(), c, ocpath.Root().System().Config())
			if err != nil {
				t.Fatal(err)
			}
			val, ok := v.Val()
			if !ok {
				return
			}
			nos, err := ygot.Diff(&oc.System{}, val)
			if err != nil {
				t.Fatal(err)
			}
			if len(nos.Update)+len(nos.Delete) != 0 {
				t.Errorf("Got diff:\n%s", nos)
			}
		},
	}, {
		desc: "fail due to missing leafref",
		inOp: func(c *ygnmi.Client) error {
			_, err := ygnmi.Update(context.Background(), c, ocpath.Root().Lldp().Interface("eth1").Name().Config(), "eth1")
			return err
		},
		wantErr: "pointed-to value with path /interfaces/interface/name",
	}, {
		desc:    "leaf state update",
		isState: true,
		inOp: func(c *ygnmi.Client) error {
			_, err := gnmiclient.Update(context.Background(), c, ocpath.Root().System().Hostname().State(), "foo")
			return err
		},
		checkOp: func(t *testing.T, c *ygnmi.Client) {
			v, err := ygnmi.Get(context.Background(), c, ocpath.Root().System().Hostname().State())
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff("foo", v); diff != "" {
				t.Errorf("(-want, +got):\n%s", diff)
			}
		},
	}, {
		desc:    "leaf state replace",
		isState: true,
		inOp: func(c *ygnmi.Client) error {
			_, err := gnmiclient.Replace(context.Background(), c, ocpath.Root().System().Hostname().State(), "foo")
			return err
		},
		checkOp: func(t *testing.T, c *ygnmi.Client) {
			v, err := ygnmi.Get(context.Background(), c, ocpath.Root().System().Hostname().State())
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff("foo", v); diff != "" {
				t.Errorf("(-want, +got):\n%s", diff)
			}
		},
	}, {
		desc:    "leaf state delete",
		isState: true,
		inOp: func(c *ygnmi.Client) error {
			if _, err := gnmiclient.Update(context.Background(), c, ocpath.Root().System().Hostname().State(), "foo"); err != nil {
				return err
			}
			_, err := gnmiclient.Delete(context.Background(), c, ocpath.Root().System().Hostname().State())
			return err
		},
		checkOp: func(t *testing.T, c *ygnmi.Client) {
			v, err := ygnmi.Lookup(context.Background(), c, ocpath.Root().System().Hostname().State())
			if err != nil {
				t.Fatal(err)
			}
			if v.IsPresent() {
				t.Errorf("Got present, want not present")
			}
		},
	}, {
		desc:    "non-leaf state update",
		isState: true,
		inOp: func(c *ygnmi.Client) error {
			if _, err := gnmiclient.Update(context.Background(), c, ocpath.Root().System().State(), &oc.System{Hostname: ygot.String("foo")}); err != nil {
				return err
			}
			_, err := gnmiclient.Update(context.Background(), c, ocpath.Root().System().State(), &oc.System{MotdBanner: ygot.String("bar")})
			return err
		},
		checkOp: func(t *testing.T, c *ygnmi.Client) {
			v, err := ygnmi.Get(context.Background(), c, ocpath.Root().System().State())
			if err != nil {
				t.Fatal(err)
			}
			want := &oc.System{Hostname: ygot.String("foo"), MotdBanner: ygot.String("bar")}
			want.PopulateDefaults()
			nos, err := ygot.Diff(want, v)
			if err != nil {
				t.Fatal(err)
			}
			if len(nos.Update)+len(nos.Delete) != 0 {
				t.Errorf("Got diff:\n%s", nos)
			}
		},
	}, {
		desc:    "non-leaf state replace",
		isState: true,
		inOp: func(c *ygnmi.Client) error {
			if _, err := gnmiclient.Update(context.Background(), c, ocpath.Root().System().State(), &oc.System{Hostname: ygot.String("foo")}); err != nil {
				return err
			}
			_, err := gnmiclient.Replace(context.Background(), c, ocpath.Root().System().State(), &oc.System{MotdBanner: ygot.String("foo")})
			return err
		},
		checkOp: func(t *testing.T, c *ygnmi.Client) {
			v, err := ygnmi.Get(context.Background(), c, ocpath.Root().System().State())
			if err != nil {
				t.Fatal(err)
			}
			want := &oc.System{MotdBanner: ygot.String("foo")}
			want.PopulateDefaults()
			nos, err := ygot.Diff(want, v)
			if err != nil {
				t.Fatal(err)
			}
			if len(nos.Update)+len(nos.Delete) != 0 {
				t.Errorf("Got diff:\n%s", nos)
			}
		},
	}, {
		desc:    "non-leaf state delete",
		isState: true,
		inOp: func(c *ygnmi.Client) error {
			if _, err := gnmiclient.Update(context.Background(), c, ocpath.Root().System().State(), &oc.System{Hostname: ygot.String("foo")}); err != nil {
				return err
			}
			_, err := gnmiclient.Delete(context.Background(), c, ocpath.Root().System().State())
			return err
		},
		checkOp: func(t *testing.T, c *ygnmi.Client) {
			v, err := ygnmi.Lookup(context.Background(), c, ocpath.Root().System().State())
			if err != nil {
				t.Fatal(err)
			}
			val, ok := v.Val()
			if !ok {
				return
			}
			want := &oc.System{}
			want.PopulateDefaults()
			nos, err := ygot.Diff(want, val)
			if err != nil {
				t.Fatal(err)
			}
			if len(nos.Update)+len(nos.Delete) != 0 {
				t.Errorf("Got diff:\n%s", nos)
			}
		},
	}}

	gnmiServer, err := newServer(context.Background(), "local", true)
	if err != nil {
		t.Fatalf("cannot create server, got err: %v", err)
	}
	addr, err := startServer(gnmiServer)
	if err != nil {
		t.Fatalf("cannot start server, got err: %v", err)
	}
	defer gnmiServer.c.Stop()
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(local.NewCredentials()))
	if err != nil {
		t.Fatalf("cannot dial gNMI server, %v", err)
	}
	configClient, err := ygnmi.NewClient(gpb.NewGNMIClient(conn), ygnmi.WithTarget("local"))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	stateClient, err := ygnmi.NewClient(gnmiServer.LocalClient(), ygnmi.WithTarget("local"))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	for _, tt := range tests {
		c := configClient
		if tt.isState {
			c = stateClient
		}
		t.Run(tt.desc, func(t *testing.T) {
			err := tt.inOp(c)
			if d := errdiff.Check(err, tt.wantErr); d != "" {
				t.Errorf("Set() unexpected err: %s", d)
			}
			if tt.wantErr == "" {
				tt.checkOp(t, c)
			}
		})
	}
}

func TestSetWithAuth(t *testing.T) {
	tests := []struct {
		desc      string
		authAllow bool
		user      string
		wantErr   string
	}{{
		desc:      "allowed",
		authAllow: true,
		user:      "test",
	}, {
		desc:      "denied",
		authAllow: false,
		user:      "test",
		wantErr:   "PermissionDenied",
	}, {
		desc:      "error no user",
		authAllow: false,
		wantErr:   "no username set in metadata",
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			gnmiServer, err := newServer(context.Background(), "local", true)
			if err != nil {
				t.Fatalf("cannot create server, got err: %v", err)
			}

			addr, err := startServer(gnmiServer)
			if err != nil {
				t.Fatalf("cannot start server, got err: %v", err)
			}
			defer gnmiServer.c.Stop()
			conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(local.NewCredentials()))
			if err != nil {
				t.Fatalf("cannot dial gNMI server, %v", err)
			}
			c, err := ygnmi.NewClient(gpb.NewGNMIClient(conn), ygnmi.WithTarget("local"))
			if err != nil {
				t.Fatalf("failed to create client: %v", err)
			}

			gnmiServer.pathAuth = &testAuth{allow: tt.authAllow}

			ctx := metadata.NewOutgoingContext(context.Background(), metadata.New(map[string]string{"username": tt.user}))
			_, err = ygnmi.Update(ctx, c, ocpath.Root().System().Hostname().Config(), "test")
			if d := errdiff.Check(err, tt.wantErr); d != "" {
				t.Errorf("Set() unexpected err: %s", d)
			}
		})
	}
}

// TestSTREAM tests the STREAM mode of gnmit.
func TestSTREAM(t *testing.T) {
	ctx := context.Background()
	gnmiServer, err := newServer(ctx, "local", false)
	if err != nil {
		t.Fatalf("cannot create server, got err: %v", err)
	}
	addr, err := startServer(gnmiServer)
	if err != nil {
		t.Fatalf("cannot start server, got err: %v", err)
	}

	gnmiServer.c.TargetUpdate(&gpb.SubscribeResponse{
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
	gnmiServer.c.TargetUpdate(&gpb.SubscribeResponse{
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
		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(local.NewCredentials()))
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
			gnmiServer.c.TargetUpdate(&gpb.SubscribeResponse{
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
	gnmiServer.c.Stop()

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

type testAuth struct {
	allow bool
}

func (t testAuth) CheckPermit(path *gpb.Path, user string, write bool) bool {
	return t.allow
}

func (t testAuth) IsInitialized() bool {
	return true
}

func TestSubscribeWithAuth(t *testing.T) {
	tests := []struct {
		desc      string
		authAllow bool
		user      string
		want      oc.E_Interface_OperStatus
		wantErr   string
	}{{
		desc:      "allowed",
		authAllow: true,
		user:      "test",
		want:      oc.Interface_OperStatus_UP,
	}, {
		desc:      "denied",
		authAllow: false,
		user:      "test",
		wantErr:   "not present",
	}, {
		desc:      "error no user",
		authAllow: false,
		wantErr:   "no username set in metadata",
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			gnmiServer, err := newServer(context.Background(), "local", false)
			if err != nil {
				t.Fatalf("cannot create server, got err: %v", err)
			}
			gnmiServer.pathAuth = &testAuth{allow: tt.authAllow}
			addr, err := startServer(gnmiServer)
			if err != nil {
				t.Fatalf("cannot start server, got err: %v", err)
			}
			gnmiServer.c.TargetUpdate(&gpb.SubscribeResponse{
				Response: &gpb.SubscribeResponse_Update{
					Update: &gpb.Notification{
						Prefix:    mustTargetPath("local", "", true),
						Timestamp: 1,
						Update: []*gpb.Update{{
							Path: mustPath("/interfaces/interface[name=eth0]/state/oper-status"),
							Val:  mustTypedValue("UP"),
						}},
					},
				},
			})
			defer gnmiServer.c.Stop()
			conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(local.NewCredentials()))
			if err != nil {
				t.Fatalf("cannot dial gNMI server, %v", err)
			}
			c, err := ygnmi.NewClient(gpb.NewGNMIClient(conn), ygnmi.WithTarget("local"))
			if err != nil {
				t.Fatalf("failed to create client: %v", err)
			}
			ctx := metadata.NewOutgoingContext(context.Background(), metadata.New(map[string]string{"username": tt.user}))
			got, err := ygnmi.Get(ctx, c, ocpath.Root().Interface("eth0").OperStatus().State())
			if d := errdiff.Check(err, tt.wantErr); d != "" {
				t.Errorf("Subscribe() unexpected err: %s", d)
			}
			if err != nil {
				return
			}
			if d := cmp.Diff(tt.want, got); d != "" {
				t.Errorf("Subscribe() unexpected diff: %s", d)
			}
		})
	}
}

func TestWithCreds(t *testing.T) {
	creds, err := creds.NewCreds()
	if err != nil {
		t.Fatalf("cannot open test TLS credentials, %v", err)
	}

	gnmiServer, err := newServer(context.Background(), "local", true)
	if err != nil {
		t.Fatalf("cannot create server, got err: %v", err)
	}
	addr, err := startServer(gnmiServer, grpc.Creds(creds))
	if err != nil {
		t.Fatalf("cannot start server, got err: %v", err)
	}

	cctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if _, err := grpc.DialContext(cctx, addr,
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: true,
		})),
		grpc.WithBlock()); err != nil {
		t.Fatalf("cannot dial server with TLS credentials, err: %v", err)
	}
}
