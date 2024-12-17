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
	"github.com/openconfig/gnmi/errdiff"
	"github.com/openconfig/gnmi/value"
	"github.com/openconfig/ygnmi/schemaless"
	"github.com/openconfig/ygnmi/ygnmi"
	"github.com/openconfig/ygot/ygot"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/local"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/prototext"

	"github.com/openconfig/lemming/gnmi/gnmiclient"
	"github.com/openconfig/lemming/gnmi/oc"
	"github.com/openconfig/lemming/gnmi/oc/ocpath"

	gpb "github.com/openconfig/gnmi/proto/gnmi"
)

const (
	targetName = "local"
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
	T      updateType
	TS     int64
	Target string
	Path   string
	Val    interface{}
}

func (u *upd) String() string {
	b := &bytes.Buffer{}
	b.WriteString("<")
	switch u.T {
	case VAL:
		b.WriteString(fmt.Sprintf("value, @%d (%s) %s=%v", u.TS, u.Target, u.Path, u.Val))
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
					T:      VAL,
					TS:     v.Update.GetTimestamp(),
					Target: v.Update.GetPrefix().GetTarget(),
					Path:   mustPathToString(u.Path),
					Val:    mustToScalar(u.Val),
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
func startServer(s *Server) (string, error) {
	// Start gNMI server.
	srv := grpc.NewServer(grpc.StreamInterceptor(NewSubscribeTargetUpdateInterceptor(targetName)))
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

func TestSetAndSubscribeOnce(t *testing.T) {
	tests := []struct {
		desc     string
		pathStr  string
		useSet   bool
		config   bool
		internal bool
	}{{
		desc:     "subscribe-once",
		pathStr:  "/hello",
		useSet:   false,
		config:   false,
		internal: false,
	}, {
		desc:     "set-config",
		pathStr:  "/interfaces/interface[name=foo]/config/description",
		useSet:   true,
		config:   true,
		internal: false,
	}, {
		desc:     "set-state",
		pathStr:  "/interfaces/interface[name=foo]/state/description",
		useSet:   true,
		config:   false,
		internal: false,
	}, {
		desc:     "set-internal",
		pathStr:  "/test/foo",
		useSet:   true,
		config:   true,
		internal: true,
	}}

	for _, tt := range tests {
		path := mustPath(tt.pathStr)
		if tt.internal {
			path.Origin = InternalOrigin
		}
		client, cleanup := testSetSubSetup(t, mustTargetPath(targetName, "", tt.useSet && !tt.internal), path, tt.useSet, tt.config)
		defer cleanup()
		t.Run(tt.desc+"-with-target", func(t *testing.T) {
			testSetSub(t, client, path, tt.config, targetName, tt.useSet, tt.internal)
		})
		t.Run(tt.desc+"-no-target", func(t *testing.T) {
			testSetSub(t, client, path, tt.config, "", tt.useSet, tt.internal)
		})
		// Run this again for repeatability (e.g. make sure Target in
		// the notification didn't get overwritten).
		t.Run(tt.desc+"-with-target-2", func(t *testing.T) {
			testSetSub(t, client, path, tt.config, targetName, tt.useSet, tt.internal)
		})
		t.Run(tt.desc+"-no-target-2", func(t *testing.T) {
			testSetSub(t, client, path, tt.config, "", tt.useSet, tt.internal)
		})
	}
}

// testSetSubSetup tests gnmi.Set and/or gnmi.Subscribe/ONCE on a config or state value.
//
// It purposely avoids using ygnmi in order to test lower-level details
// (e.g. timestamp metadata)
func testSetSubSetup(t *testing.T, prefix, path *gpb.Path, useSet, config bool) (gpb.GNMIClient, func()) {
	ctx := context.Background()
	gnmiServer, err := newServer(ctx, targetName, useSet)
	if err != nil {
		t.Fatalf("cannot create server, got err: %v", err)
	}
	addr, err := startServer(gnmiServer)
	if err != nil {
		t.Fatalf("cannot start server, got err: %v", err)
	}

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(local.NewCredentials()))
	if err != nil {
		t.Fatalf("cannot dial gNMI server, %v", err)
	}

	client := gpb.NewGNMIClient(conn)

	if useSet {
		if !config {
			ctx = metadata.AppendToOutgoingContext(ctx, GNMIModeMetadataKey, string(StateMode))
		}
		if _, err := client.Set(metadata.AppendToOutgoingContext(ctx, TimestampMetadataKey, strconv.FormatInt(42, 10)), &gpb.SetRequest{
			Prefix: prefix,
			Replace: []*gpb.Update{{
				Path: path,
				Val:  mustTypedValue("world"),
			}},
		}); err != nil {
			t.Fatalf("set request failed: %v", err)
		}
	} else {
		gnmiServer.c.TargetUpdate(&gpb.SubscribeResponse{
			Response: &gpb.SubscribeResponse_Update{
				Update: &gpb.Notification{
					Prefix:    prefix,
					Timestamp: 42,
					Update: []*gpb.Update{{
						Path: path,
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
	}

	return client, func() {
		gnmiServer.c.Stop()
	}
}

// testSetSub tests gnmi.Set and/or gnmi.Subscribe/ONCE on a config or state value.
//
// It purposely avoids using ygnmi in order to test lower-level details
// (e.g. timestamp metadata)
func testSetSub(t *testing.T, client gpb.GNMIClient, path *gpb.Path, config bool, wantTarget string, useSet, internal bool) {
	prefix := mustTargetPath(wantTarget, "", useSet && !internal)
	pathStr := mustPathToString(path)

	got := []*upd{}
	clientCtx, cancel := context.WithCancel(context.Background())
	var sendErr, recvErr error
	go func(ctx context.Context) {
		defer cancel()
		subc, err := client.Subscribe(ctx)
		if err != nil {
			sendErr = err
			return
		}
		sr := &gpb.SubscribeRequest{
			Request: &gpb.SubscribeRequest_Subscribe{
				Subscribe: &gpb.SubscriptionList{
					Prefix: prefix,
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

	if sendErr != nil {
		t.Errorf("got unexpected send error, %v", sendErr)
	}

	if recvErr != nil {
		t.Errorf("got unexpected recv error, %v", recvErr)
	}

	var cmpOptions []cmp.Option
	if config {
		cmpOptions = append(cmpOptions, cmpopts.IgnoreFields(upd{}, "TS"))
	}

	if diff := cmp.Diff(got, []*upd{{
		T:      VAL,
		TS:     42,
		Target: wantTarget,
		Path:   pathStr,
		Val:    "world",
	}, {
		T: SYNC,
	}}, cmpOptions...); diff != "" {
		t.Fatalf("did not get expected updates, diff(-got,+want)\n:%s", diff)
	}

	if config {
		// Test that timestamp is not 42: we don't want the timestamp metadata to affect config values.
		if cmp.Equal(got, []*upd{{
			T:      VAL,
			TS:     42,
			Target: wantTarget,
			Path:   pathStr,
			Val:    "world",
		}, {
			T: SYNC,
		}}) {
			t.Fatalf("Expected error -- timestamp metadata should be ignored but it is not ignored.")
		}
	}
}

func TestSetYGNMI(t *testing.T) {
	type testSpec struct {
		desc      string
		isState   bool
		inOp      func(c *ygnmi.Client) error
		getOp     func(t *testing.T, c *ygnmi.Client) (interface{}, bool)
		wantValue any
		wantErr   string
	}

	prefixSetName := "accept-route"
	policyStmts := &oc.RoutingPolicy_PolicyDefinition_Statement_OrderedMap{}
	policyName := "one"
	stmt, err := policyStmts.AppendNew("stmt1")
	if err != nil {
		t.Fatal(err)
	}
	stmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetPrefixSet(prefixSetName)
	stmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetMatchSetOptions(oc.PolicyTypes_MatchSetOptionsRestrictedType_ANY)
	stmt.GetOrCreateActions().GetOrCreateBgpActions().GetOrCreateSetCommunity().SetOptions(oc.BgpPolicy_BgpSetCommunityOptionType_REPLACE)
	stmt.GetOrCreateActions().GetOrCreateBgpActions().GetOrCreateSetCommunity().GetOrCreateInline().SetCommunities(
		[]oc.RoutingPolicy_PolicyDefinition_Statement_Actions_BgpActions_SetCommunity_Inline_Communities_Union{
			oc.UnionString("10000:10000"),
		},
	)
	stmt, err = policyStmts.AppendNew("stmt2")
	if err != nil {
		t.Fatal(err)
	}
	stmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetPrefixSet(prefixSetName)
	stmt.GetOrCreateConditions().GetOrCreateMatchPrefixSet().SetMatchSetOptions(oc.PolicyTypes_MatchSetOptionsRestrictedType_INVERT)
	stmt.GetOrCreateActions().GetOrCreateBgpActions().GetOrCreateSetCommunity().SetOptions(oc.BgpPolicy_BgpSetCommunityOptionType_ADD)
	stmt.GetOrCreateActions().GetOrCreateBgpActions().GetOrCreateSetCommunity().GetOrCreateInline().SetCommunities(
		[]oc.RoutingPolicy_PolicyDefinition_Statement_Actions_BgpActions_SetCommunity_Inline_Communities_Union{
			oc.UnionString("20000:20000"),
		},
	)
	policy := &oc.RoutingPolicy_PolicyDefinition{Name: ygot.String(policyName), Statement: policyStmts}
	policy.PopulateDefaults()

	schemalessQuery, err := schemaless.NewConfig[string](fmt.Sprintf("/dataplane/routes/route[prefix=%s][vrf=%d]", "1", 2), InternalOrigin)
	if err != nil {
		t.Fatal(err)
	}

	passingTests := []testSpec{{
		desc: "leaf config update",
		inOp: func(c *ygnmi.Client) error {
			_, err := ygnmi.Update(context.Background(), c, ocpath.Root().System().Hostname().Config(), "foo")
			return err
		},
		wantValue: "foo",
		getOp: func(t *testing.T, c *ygnmi.Client) (interface{}, bool) {
			v, err := ygnmi.Lookup[string](context.Background(), c, ocpath.Root().System().Hostname().Config())
			if err != nil {
				t.Fatal(err)
			}
			return v.Val()
		},
	}, {
		desc: "leaf config replace",
		inOp: func(c *ygnmi.Client) error {
			_, err := ygnmi.Replace(context.Background(), c, ocpath.Root().System().Hostname().Config(), "foo")
			return err
		},
		wantValue: "foo",
		getOp: func(t *testing.T, c *ygnmi.Client) (interface{}, bool) {
			v, err := ygnmi.Lookup[string](context.Background(), c, ocpath.Root().System().Hostname().Config())
			if err != nil {
				t.Fatal(err)
			}
			return v.Val()
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
		wantValue: nil,
		getOp: func(t *testing.T, c *ygnmi.Client) (interface{}, bool) {
			v, err := ygnmi.Lookup[string](context.Background(), c, ocpath.Root().System().Hostname().Config())
			if err != nil {
				t.Fatal(err)
			}
			return v.Val()
		},
	}, {
		desc: "leaf config update internal",
		inOp: func(c *ygnmi.Client) error {
			_, err = ygnmi.Update(context.Background(), c, schemalessQuery, "foo")
			return err
		},
		wantValue: "foo",
		getOp: func(t *testing.T, c *ygnmi.Client) (interface{}, bool) {
			v, err := ygnmi.Lookup[string](context.Background(), c, schemalessQuery)
			if err != nil {
				t.Fatal(err)
			}
			return v.Val()
		},
	}, {
		desc: "leaf config replace internal",
		inOp: func(c *ygnmi.Client) error {
			_, err = ygnmi.Replace(context.Background(), c, schemalessQuery, "foo")
			return err
		},
		wantValue: "foo",
		getOp: func(t *testing.T, c *ygnmi.Client) (interface{}, bool) {
			v, err := ygnmi.Lookup[string](context.Background(), c, schemalessQuery)
			if err != nil {
				t.Fatal(err)
			}
			return v.Val()
		},
	}, {
		desc: "leaf config delete internal",
		inOp: func(c *ygnmi.Client) error {
			if _, err := ygnmi.Update(context.Background(), c, schemalessQuery, "foo"); err != nil {
				return err
			}
			_, err = ygnmi.Delete(context.Background(), c, schemalessQuery)
			return err
		},
		wantValue: nil,
		getOp: func(t *testing.T, c *ygnmi.Client) (interface{}, bool) {
			v, err := ygnmi.Lookup[string](context.Background(), c, schemalessQuery)
			if err != nil {
				t.Fatal(err)
			}
			return v.Val()
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
		wantValue: func() interface{} {
			want := &oc.System{Hostname: ygot.String("foo"), MotdBanner: ygot.String("bar")}
			want.PopulateDefaults()
			return want
		}(),
		getOp: func(t *testing.T, c *ygnmi.Client) (interface{}, bool) {
			v, err := ygnmi.Lookup[*oc.System](context.Background(), c, ocpath.Root().System().Config())
			if err != nil {
				t.Fatal(err)
			}
			return v.Val()
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
		wantValue: func() interface{} {
			want := &oc.System{MotdBanner: ygot.String("foo")}
			want.PopulateDefaults()
			return want
		}(),
		getOp: func(t *testing.T, c *ygnmi.Client) (interface{}, bool) {
			v, err := ygnmi.Lookup[*oc.System](context.Background(), c, ocpath.Root().System().Config())
			if err != nil {
				t.Fatal(err)
			}
			return v.Val()
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
		wantValue: func() interface{} {
			want := &oc.System{}
			want.PopulateDefaults()
			return want
		}(),
		getOp: func(t *testing.T, c *ygnmi.Client) (interface{}, bool) {
			v, err := ygnmi.Lookup[*oc.System](context.Background(), c, ocpath.Root().System().Config())
			if err != nil {
				t.Fatal(err)
			}
			return v.Val()
		},
	}, {
		desc:    "leaf state update",
		isState: true,
		inOp: func(c *ygnmi.Client) error {
			_, err := gnmiclient.Update(context.Background(), c, ocpath.Root().System().Hostname().State(), "foo")
			return err
		},
		wantValue: "foo",
		getOp: func(t *testing.T, c *ygnmi.Client) (interface{}, bool) {
			v, err := ygnmi.Lookup(context.Background(), c, ocpath.Root().System().Hostname().State())
			if err != nil {
				t.Fatal(err)
			}
			return v.Val()
		},
	}, {
		desc:    "leaf state replace",
		isState: true,
		inOp: func(c *ygnmi.Client) error {
			_, err := gnmiclient.Replace(context.Background(), c, ocpath.Root().System().Hostname().State(), "foo")
			return err
		},
		wantValue: "foo",
		getOp: func(t *testing.T, c *ygnmi.Client) (interface{}, bool) {
			v, err := ygnmi.Lookup(context.Background(), c, ocpath.Root().System().Hostname().State())
			if err != nil {
				t.Fatal(err)
			}
			return v.Val()
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
		wantValue: nil,
		getOp: func(t *testing.T, c *ygnmi.Client) (interface{}, bool) {
			v, err := ygnmi.Lookup(context.Background(), c, ocpath.Root().System().Hostname().State())
			if err != nil {
				t.Fatal(err)
			}
			return v.Val()
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
		wantValue: &oc.System{Hostname: ygot.String("foo"), MotdBanner: ygot.String("bar")},
		getOp: func(t *testing.T, c *ygnmi.Client) (interface{}, bool) {
			v, err := ygnmi.Lookup(context.Background(), c, ocpath.Root().System().State())
			if err != nil {
				t.Fatal(err)
			}
			return v.Val()
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
		wantValue: &oc.System{MotdBanner: ygot.String("foo")},
		getOp: func(t *testing.T, c *ygnmi.Client) (interface{}, bool) {
			v, err := ygnmi.Lookup(context.Background(), c, ocpath.Root().System().State())
			if err != nil {
				t.Fatal(err)
			}
			return v.Val()
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
		wantValue: &oc.System{},
		getOp: func(t *testing.T, c *ygnmi.Client) (interface{}, bool) {
			v, err := ygnmi.Lookup(context.Background(), c, ocpath.Root().System().State())
			if err != nil {
				t.Fatal(err)
			}
			return v.Val()
		},
	}, {
		desc: "telemetry-atomic-config-replace",
		inOp: func(c *ygnmi.Client) error {
			route := "1.1.0.0/16"
			prefixPath := ocpath.Root().RoutingPolicy().DefinedSets().PrefixSet(prefixSetName).Prefix(route, "exact").IpPrefix()
			_, err = ygnmi.Replace(context.Background(), c, prefixPath.Config(), route)
			if err != nil {
				return err
			}
			_, err = ygnmi.Replace(context.Background(), c, ocpath.Root().RoutingPolicy().PolicyDefinition(policyName).Config(), policy)
			return err
		},
		wantValue: policy,
		getOp: func(t *testing.T, c *ygnmi.Client) (interface{}, bool) {
			v, err := ygnmi.Lookup[*oc.RoutingPolicy_PolicyDefinition](context.Background(), c, ocpath.Root().RoutingPolicy().PolicyDefinition(policyName).Config())
			if err != nil {
				t.Fatal(err)
			}
			return v.Val()
		},
	}}

	failingTests := []testSpec{{
		desc: "fail due to missing leafref",
		inOp: func(c *ygnmi.Client) error {
			_, err := ygnmi.Update(context.Background(), c, ocpath.Root().Lldp().Interface("eth1").Name().Config(), "eth1")
			return err
		},
		wantErr: "pointed-to value with path /interfaces/interface/name",
	}, {
		desc: "fail due to non-matching key names",
		inOp: func(c *ygnmi.Client) error {
			_, err := ygnmi.Update(context.Background(), c, ocpath.Root().RoutingPolicy().DefinedSets().PrefixSet("test").Prefix("1.1.1.1", "exact").IpPrefix().Config(), "2.2.2.2/32")
			return err
		},
		wantErr: "key value 2.2.2.2/32 for key field IpPrefix has different value from map key 1.1.1.1",
	}, {
		desc: "fail due to bad regex",
		inOp: func(c *ygnmi.Client) error {
			_, err := ygnmi.Update(context.Background(), c, ocpath.Root().RoutingPolicy().DefinedSets().PrefixSet("test").Prefix("1.1.1.1", "24").IpPrefix().Config(), "1.1.1.1")
			return err
		},
		wantErr: `"24" does not match regular expression pattern`,
	}}

	gnmiServer, err := newServer(context.Background(), targetName, true)
	if err != nil {
		t.Fatalf("cannot create server, got err: %v", err)
	}
	addr, err := startServer(gnmiServer)
	if err != nil {
		t.Fatalf("cannot start server, got err: %v", err)
	}
	t.Logf("Running gNMI server on %s", addr)
	defer gnmiServer.c.Stop()
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(local.NewCredentials()))
	if err != nil {
		t.Fatalf("cannot dial gNMI server, %v", err)
	}
	configClient, err := ygnmi.NewClient(gpb.NewGNMIClient(conn), ygnmi.WithTarget(targetName))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	stateClient, err := ygnmi.NewClient(gnmiServer.LocalClient(), ygnmi.WithTarget(targetName))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	for _, tt := range append(passingTests, failingTests...) {
		c := configClient
		if tt.isState {
			c = stateClient
		}
		t.Run(tt.desc, func(t *testing.T) {
			err := tt.inOp(c)
			if d := errdiff.Check(err, tt.wantErr); d != "" {
				t.Errorf("Set() unexpected err: %s", d)
			}
			if tt.wantErr != "" {
				return
			}

			got, ok := tt.getOp(t, c)
			switch want := tt.wantValue.(type) {
			case nil:
				if ok {
					t.Errorf("Got present, want not present")
				}
			case ygot.GoStruct:
				if !ok {
					t.Fatalf("Got not present, want present")
				}
				gotGS, ok := got.(ygot.GoStruct)
				if !ok {
					t.Fatalf("Got object not a GoStruct")
				}
				// Diffs between empty structs and nil should be ignored, so populate defaults.
				want.(populateDefaultser).PopulateDefaults()
				gotGS.(populateDefaultser).PopulateDefaults()
				nos, err := ygot.Diff(want, gotGS)
				if err != nil {
					t.Fatal(err)
				}
				if len(nos.Update)+len(nos.Delete) != 0 {
					t.Errorf("Got diff:\n%s\n(-want, +got):\n%s", nos, cmp.Diff(want, gotGS, cmp.AllowUnexported(oc.RoutingPolicy_PolicyDefinition_Statement_OrderedMap{})))
				}
			default:
				if diff := cmp.Diff(want, got); diff != "" {
					t.Errorf("(-want, +got):\n%s", diff)
				}
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
			gnmiServer, err := newServer(context.Background(), targetName, true)
			if err != nil {
				t.Fatalf("cannot create server, got err: %v", err)
			}

			addr, err := startServer(gnmiServer)
			if err != nil {
				t.Fatalf("cannot start server, got err: %v", err)
			}
			defer gnmiServer.c.Stop()
			conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(local.NewCredentials()))
			if err != nil {
				t.Fatalf("cannot dial gNMI server, %v", err)
			}
			c, err := ygnmi.NewClient(gpb.NewGNMIClient(conn), ygnmi.WithTarget(targetName))
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
	gnmiServer, err := newServer(ctx, targetName, false)
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
				Prefix:    mustTargetPath(targetName, "", false),
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
		conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(local.NewCredentials()))
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
					Prefix: mustTargetPath(targetName, "", false),
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
						Prefix:    mustTargetPath(targetName, "", false),
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
		T:      VAL,
		TS:     42,
		Target: targetName,
		Path:   "/hello",
		Val:    "world",
	}, {
		T: SYNC,
	}, {
		T:      VAL,
		TS:     43,
		Target: targetName,
		Path:   "/hello",
		Val:    "mercury",
	}, {
		T:      VAL,
		TS:     44,
		Target: targetName,
		Path:   "/hello",
		Val:    "venus",
	}, {
		T:      VAL,
		TS:     45,
		Target: targetName,
		Path:   "/hello",
		Val:    "earth",
	}, {
		T:      VAL,
		TS:     46,
		Target: targetName,
		Path:   "/hello",
		Val:    "mars",
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
		return false
	})); diff != "" {
		t.Fatalf("did not get expected updates, diff(-got,+want)\n:%s", diff)
	}
}

type testAuth struct {
	allow bool
}

func (t testAuth) CheckPermit(*gpb.Path, string, bool) bool {
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
			gnmiServer, err := newServer(context.Background(), targetName, false)
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
						Prefix:    mustTargetPath(targetName, "", true),
						Timestamp: 1,
						Update: []*gpb.Update{{
							Path: mustPath("/interfaces/interface[name=eth0]/state/oper-status"),
							Val:  mustTypedValue("UP"),
						}},
					},
				},
			})
			defer gnmiServer.c.Stop()
			conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(local.NewCredentials()))
			if err != nil {
				t.Fatalf("cannot dial gNMI server, %v", err)
			}
			c, err := ygnmi.NewClient(gpb.NewGNMIClient(conn), ygnmi.WithTarget(targetName))
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
