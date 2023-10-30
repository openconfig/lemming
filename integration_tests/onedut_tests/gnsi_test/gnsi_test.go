// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gnsi_test

import (
	"context"
	"testing"

	"github.com/openconfig/gnmi/errdiff"
	"github.com/openconfig/ondatra"
	"github.com/openconfig/ondatra/gnmi"
	"github.com/openconfig/ondatra/gnmi/oc"
	"github.com/openconfig/ygnmi/ygnmi"
	"github.com/openconfig/ygot/ygot"
	"google.golang.org/grpc/metadata"

	"github.com/openconfig/lemming/internal/binding"

	gpb "github.com/openconfig/gnmi/proto/gnmi"
	pathzpb "github.com/openconfig/gnsi/pathz"
)

func TestMain(m *testing.M) {
	ondatra.RunTests(m, binding.Get(".."))
}

func TestPathz(t *testing.T) {
	// TODO: Remove when KNEbind supports DialGNSI.
	t.Skip("Ondatra/KNEbind currently doesn't support dialing gNSI")
	dut := ondatra.DUT(t, "dut")
	yc, err := ygnmi.NewClient(dut.RawAPIs().GNMI(t), ygnmi.WithTarget(dut.Name()))
	if err != nil {
		t.Fatal(err)
	}
	defer reset(t, dut)

	tests := []struct {
		desc       string
		policy     *pathzpb.AuthorizationPolicy
		op         func(ctx context.Context, c *ygnmi.Client) (*ygnmi.Result, error)
		wantOpErr  string
		wantGetErr string
		want       *oc.System
	}{{
		desc: "success",
		policy: &pathzpb.AuthorizationPolicy{
			Rules: []*pathzpb.AuthorizationRule{{
				Path:      mustPath(t, "/system"),
				Principal: &pathzpb.AuthorizationRule_User{User: "testuser"},
				Action:    pathzpb.Action_ACTION_PERMIT,
				Mode:      pathzpb.Mode_MODE_WRITE,
			}, {
				Path:      mustPath(t, "/system"),
				Principal: &pathzpb.AuthorizationRule_User{User: "testuser"},
				Action:    pathzpb.Action_ACTION_PERMIT,
				Mode:      pathzpb.Mode_MODE_READ,
			}},
		},
		op: func(ctx context.Context, c *ygnmi.Client) (*ygnmi.Result, error) {
			return ygnmi.Update(ctx, c, gnmi.OC().System().Config(), &oc.System{Hostname: ygot.String("test")})
		},
		want: &oc.System{
			Hostname:   ygot.String("test"),
			DomainName: ygot.String("lemming.example.com"),
		},
	}, {
		desc: "deny match inside struct",
		policy: &pathzpb.AuthorizationPolicy{
			Rules: []*pathzpb.AuthorizationRule{{
				Path:      mustPath(t, "/system"),
				Principal: &pathzpb.AuthorizationRule_User{User: "testuser"},
				Action:    pathzpb.Action_ACTION_PERMIT,
				Mode:      pathzpb.Mode_MODE_WRITE,
			}, {
				Path:      mustPath(t, "/system/config/hostname"),
				Principal: &pathzpb.AuthorizationRule_User{User: "testuser"},
				Action:    pathzpb.Action_ACTION_DENY,
				Mode:      pathzpb.Mode_MODE_WRITE,
			}, {
				Path:      mustPath(t, "/system"),
				Principal: &pathzpb.AuthorizationRule_User{User: "testuser"},
				Action:    pathzpb.Action_ACTION_PERMIT,
				Mode:      pathzpb.Mode_MODE_READ,
			}},
		},
		op: func(ctx context.Context, c *ygnmi.Client) (*ygnmi.Result, error) {
			return ygnmi.Update(ctx, c, gnmi.OC().System().Config(), &oc.System{Hostname: ygot.String("test")})
		},
		want: &oc.System{
			Hostname:   ygot.String("lemming"),
			DomainName: ygot.String("lemming.example.com"),
		},
		wantOpErr: "PermissionDenied",
	}, {
		desc: "deny implicit delete",
		policy: &pathzpb.AuthorizationPolicy{
			Rules: []*pathzpb.AuthorizationRule{{
				Path:      mustPath(t, "/system"),
				Principal: &pathzpb.AuthorizationRule_User{User: "testuser"},
				Action:    pathzpb.Action_ACTION_PERMIT,
				Mode:      pathzpb.Mode_MODE_WRITE,
			}, {
				Path:      mustPath(t, "/system/config/domain-name"),
				Principal: &pathzpb.AuthorizationRule_User{User: "testuser"},
				Action:    pathzpb.Action_ACTION_DENY,
				Mode:      pathzpb.Mode_MODE_WRITE,
			}, {
				Path:      mustPath(t, "/system"),
				Principal: &pathzpb.AuthorizationRule_User{User: "testuser"},
				Action:    pathzpb.Action_ACTION_PERMIT,
				Mode:      pathzpb.Mode_MODE_READ,
			}},
		},
		op: func(ctx context.Context, c *ygnmi.Client) (*ygnmi.Result, error) {
			return ygnmi.Replace(ctx, c, gnmi.OC().System().Config(), &oc.System{Hostname: ygot.String("test")})
		},
		want: &oc.System{
			Hostname:   ygot.String("lemming"),
			DomainName: ygot.String("lemming.example.com"),
		},
		wantOpErr: "PermissionDenied",
	}, {
		desc: "deny read all paths",
		policy: &pathzpb.AuthorizationPolicy{
			Rules: []*pathzpb.AuthorizationRule{{
				Path:      mustPath(t, "/system"),
				Principal: &pathzpb.AuthorizationRule_User{User: "testuser"},
				Action:    pathzpb.Action_ACTION_PERMIT,
				Mode:      pathzpb.Mode_MODE_WRITE,
			}, {
				Path:      mustPath(t, "/system"),
				Principal: &pathzpb.AuthorizationRule_User{User: "testuser"},
				Action:    pathzpb.Action_ACTION_DENY,
				Mode:      pathzpb.Mode_MODE_READ,
			}},
		},
		op: func(ctx context.Context, c *ygnmi.Client) (*ygnmi.Result, error) {
			return ygnmi.Replace(ctx, c, gnmi.OC().System().Config(), &oc.System{Hostname: ygot.String("test")})
		},
		want: &oc.System{
			Hostname:   ygot.String("lemming"),
			DomainName: ygot.String("lemming.example.com"),
		},
		wantGetErr: "value not present",
	}, {
		desc: "deny read sub path",
		policy: &pathzpb.AuthorizationPolicy{
			Rules: []*pathzpb.AuthorizationRule{{
				Path:      mustPath(t, "/system"),
				Principal: &pathzpb.AuthorizationRule_User{User: "testuser"},
				Action:    pathzpb.Action_ACTION_PERMIT,
				Mode:      pathzpb.Mode_MODE_WRITE,
			}, {
				Path:      mustPath(t, "/system"),
				Principal: &pathzpb.AuthorizationRule_User{User: "testuser"},
				Action:    pathzpb.Action_ACTION_PERMIT,
				Mode:      pathzpb.Mode_MODE_READ,
			}, {
				Path:      mustPath(t, "/system/config/hostname"),
				Principal: &pathzpb.AuthorizationRule_User{User: "testuser"},
				Action:    pathzpb.Action_ACTION_DENY,
				Mode:      pathzpb.Mode_MODE_READ,
			}},
		},
		op: func(ctx context.Context, c *ygnmi.Client) (*ygnmi.Result, error) {
			return ygnmi.Update(ctx, c, gnmi.OC().System().Config(), &oc.System{Hostname: ygot.String("test")})
		},
		want: &oc.System{
			DomainName: ygot.String("lemming.example.com"),
		},
	}}

	ctx := metadata.AppendToOutgoingContext(context.Background(), "username", "testuser")

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			reset(t, dut)
			installPolicy(t, dut.RawAPIs().GNSI(t).Pathz(), tt.policy)

			_, err := tt.op(ctx, yc)
			if d := errdiff.Check(err, tt.wantOpErr); d != "" {
				t.Errorf("Replace() unexpected err: %s", d)
			}
			got, err := ygnmi.Get[*oc.System](ctx, yc, gnmi.OC().System().Config())
			if d := errdiff.Check(err, tt.wantGetErr); d != "" {
				t.Errorf("Get() unexpected err: %s", d)
			}
			if err != nil {
				return
			}
			if got.GetHostname() != tt.want.GetHostname() || got.GetDomainName() != tt.want.GetDomainName() {
				t.Errorf("Get() unexpected result got (%s, %s), want (%s, %s)", got.GetHostname(), got.GetDomainName(), tt.want.GetHostname(), tt.want.GetDomainName())
			}
		})
	}
}

func mustPath(t testing.TB, s string) *gpb.Path {
	p, err := ygot.StringToStructuredPath(s)
	if err != nil {
		t.Fatalf("cannot parse  path %s, %v", s, err)
	}
	return p
}

func reset(t testing.TB, dut *ondatra.DUTDevice) {
	t.Helper()
	installPolicy(t, dut.RawAPIs().GNSI(t).Pathz(), &pathzpb.AuthorizationPolicy{
		Rules: []*pathzpb.AuthorizationRule{{
			Path:      mustPath(t, "/"),
			Principal: &pathzpb.AuthorizationRule_User{User: "testuser"},
			Action:    pathzpb.Action_ACTION_PERMIT,
			Mode:      pathzpb.Mode_MODE_WRITE,
		}, {
			Path:      mustPath(t, "/"),
			Principal: &pathzpb.AuthorizationRule_User{User: "testuser"},
			Action:    pathzpb.Action_ACTION_PERMIT,
			Mode:      pathzpb.Mode_MODE_READ,
		}},
	})

	gnmi.Update(t, dut.GNMIOpts().WithMetadata(metadata.Pairs("username", "testuser")), gnmi.OC().System().Config(), &oc.System{
		DomainName: ygot.String("lemming.example.com"),
		Hostname:   ygot.String("lemming"),
	})
}

func installPolicy(t testing.TB, pathzClient pathzpb.PathzClient, req *pathzpb.AuthorizationPolicy) {
	t.Helper()

	rc, err := pathzClient.Rotate(context.Background())
	if err != nil {
		t.Fatalf("failed to start rotation: %v", err)
	}
	if err := rc.Send(&pathzpb.RotateRequest{
		RotateRequest: &pathzpb.RotateRequest_UploadRequest{UploadRequest: &pathzpb.UploadRequest{
			Policy: req,
		}},
	}); err != nil {
		t.Fatalf("failed to send upload req: %v", err)
	}
	if _, err := rc.Recv(); err != nil {
		t.Fatalf("failed to recv upload resp: %v", err)
	}
	if err := rc.Send(&pathzpb.RotateRequest{
		RotateRequest: &pathzpb.RotateRequest_FinalizeRotation{},
	}); err != nil {
		t.Fatalf("failed to send finalize req: %v", err)
	}
	if _, err := rc.Recv(); err != nil {
		t.Fatalf("failed to recv finalize resp: %v", err)
	}
}
