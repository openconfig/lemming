// Copyright 2023 Google LLC
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

package pathz

import (
	"context"
	"net"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/testing/protocmp"

	"github.com/google/go-cmp/cmp"
	"github.com/h-fam/errdiff"
	gpb "github.com/openconfig/gnmi/proto/gnmi"
	pathzpb "github.com/openconfig/gnsi/pathz"
)

func TestRotate(t *testing.T) {
	tests := []struct {
		desc     string
		reqs     []*pathzpb.RotateRequest
		wantErrs []string
	}{{
		desc: "invalid policy",
		reqs: []*pathzpb.RotateRequest{{
			RotateRequest: &pathzpb.RotateRequest_UploadRequest{
				UploadRequest: &pathzpb.UploadRequest{
					Policy: &pathzpb.AuthorizationPolicy{
						Rules: []*pathzpb.AuthorizationRule{{}},
					},
				},
			},
		}},
		wantErrs: []string{"invalid policy"},
	}, {
		desc: "multiple uploads",
		reqs: []*pathzpb.RotateRequest{{
			RotateRequest: &pathzpb.RotateRequest_UploadRequest{
				UploadRequest: &pathzpb.UploadRequest{
					Policy: &pathzpb.AuthorizationPolicy{
						Rules: []*pathzpb.AuthorizationRule{},
					},
				},
			},
		}, {
			RotateRequest: &pathzpb.RotateRequest_UploadRequest{
				UploadRequest: &pathzpb.UploadRequest{
					Policy: &pathzpb.AuthorizationPolicy{
						Rules: []*pathzpb.AuthorizationRule{},
					},
				},
			},
		}},
		wantErrs: []string{"", "single upload request"},
	}, {
		desc: "finalize before upload",
		reqs: []*pathzpb.RotateRequest{{
			RotateRequest: &pathzpb.RotateRequest_FinalizeRotation{},
		}},
		wantErrs: []string{"finalize rotation called before upload request"},
	}, {
		desc: "multiple uploads",
		reqs: []*pathzpb.RotateRequest{{
			RotateRequest: &pathzpb.RotateRequest_UploadRequest{
				UploadRequest: &pathzpb.UploadRequest{
					Policy: &pathzpb.AuthorizationPolicy{
						Rules: []*pathzpb.AuthorizationRule{},
					},
				},
			},
		}, {
			RotateRequest: &pathzpb.RotateRequest_FinalizeRotation{},
		}},
		wantErrs: []string{"", ""},
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			client, closeFn := start(t)
			defer closeFn()
			rot, err := client.Rotate(context.Background())
			if err != nil {
				t.Fatal(err)
			}
			for i, req := range tt.reqs {
				if err := rot.Send(req); err != nil {
					t.Fatal(err)
				}
				_, err := rot.Recv()
				if d := errdiff.Check(err, tt.wantErrs[i]); d != "" {
					t.Errorf("Rotate() unexpected err: %s", d)
				}
			}
		})
	}
	t.Run("concurrent rotation", func(t *testing.T) {
		client, closeFn := start(t)
		defer closeFn()
		if _, err := client.Rotate(context.Background()); err != nil {
			t.Fatal(err)
		}
		time.Sleep(10 * time.Millisecond)
		c, err := client.Rotate(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		_, err = c.Recv()
		if d := errdiff.Check(err, "another rotation is already in progress"); d != "" {
			t.Errorf("Rotate() unexpected err: %s", d)
		}
	})
}

func TestProbe(t *testing.T) {
	tests := []struct {
		desc                string
		req                 *pathzpb.RotateRequest
		probeBeforeFinalize bool
		probeReq            *pathzpb.ProbeRequest
		want                *pathzpb.ProbeResponse
		wantErr             string
	}{{
		desc: "error no mode",
		probeReq: &pathzpb.ProbeRequest{
			User:           "me",
			Path:           &gpb.Path{},
			PolicyInstance: pathzpb.PolicyInstance_POLICY_INSTANCE_ACTIVE,
		},
		wantErr: "mode not specified",
	}, {
		desc: "error no user",
		probeReq: &pathzpb.ProbeRequest{
			Mode:           pathzpb.Mode_MODE_READ,
			Path:           &gpb.Path{},
			PolicyInstance: pathzpb.PolicyInstance_POLICY_INSTANCE_ACTIVE,
		},
		wantErr: "user not specified",
	}, {
		desc: "error no path",
		probeReq: &pathzpb.ProbeRequest{
			Mode:           pathzpb.Mode_MODE_READ,
			User:           "me",
			PolicyInstance: pathzpb.PolicyInstance_POLICY_INSTANCE_ACTIVE,
		},
		wantErr: "path not specified",
	}, {
		desc: "error no instance in req",
		probeReq: &pathzpb.ProbeRequest{
			Mode: pathzpb.Mode_MODE_READ,
			User: "me",
			Path: &gpb.Path{},
		},
		wantErr: "unknown instance type",
	}, {
		desc: "error instance not ready",
		probeReq: &pathzpb.ProbeRequest{
			Mode:           pathzpb.Mode_MODE_READ,
			User:           "me",
			Path:           &gpb.Path{},
			PolicyInstance: pathzpb.PolicyInstance_POLICY_INSTANCE_ACTIVE,
		},
		wantErr: "requested policy instance is nil",
	}, {
		desc: "success sandbox",
		probeReq: &pathzpb.ProbeRequest{
			Mode:           pathzpb.Mode_MODE_READ,
			User:           "me",
			Path:           &gpb.Path{Elem: []*gpb.PathElem{{Name: "test"}}},
			PolicyInstance: pathzpb.PolicyInstance_POLICY_INSTANCE_SANDBOX,
		},
		req: &pathzpb.RotateRequest{
			RotateRequest: &pathzpb.RotateRequest_UploadRequest{
				UploadRequest: &pathzpb.UploadRequest{
					Version: "1",
					Policy: &pathzpb.AuthorizationPolicy{
						Rules: []*pathzpb.AuthorizationRule{{
							Path:      &gpb.Path{Elem: []*gpb.PathElem{{Name: "test"}}},
							Principal: &pathzpb.AuthorizationRule_User{User: "me"},
							Mode:      pathzpb.Mode_MODE_READ,
							Action:    pathzpb.Action_ACTION_PERMIT,
						}},
					},
				},
			},
		},
		probeBeforeFinalize: true,
		want: &pathzpb.ProbeResponse{
			Version: "1",
			Action:  pathzpb.Action_ACTION_PERMIT,
		},
	}, {
		desc: "success active",
		probeReq: &pathzpb.ProbeRequest{
			Mode:           pathzpb.Mode_MODE_READ,
			User:           "me",
			Path:           &gpb.Path{Elem: []*gpb.PathElem{{Name: "test"}}},
			PolicyInstance: pathzpb.PolicyInstance_POLICY_INSTANCE_ACTIVE,
		},
		req: &pathzpb.RotateRequest{
			RotateRequest: &pathzpb.RotateRequest_UploadRequest{
				UploadRequest: &pathzpb.UploadRequest{
					Version: "1",
					Policy: &pathzpb.AuthorizationPolicy{
						Rules: []*pathzpb.AuthorizationRule{{
							Path:      &gpb.Path{Elem: []*gpb.PathElem{{Name: "test"}}},
							Principal: &pathzpb.AuthorizationRule_User{User: "me"},
							Mode:      pathzpb.Mode_MODE_READ,
							Action:    pathzpb.Action_ACTION_PERMIT,
						}},
					},
				},
			},
		},
		want: &pathzpb.ProbeResponse{
			Version: "1",
			Action:  pathzpb.Action_ACTION_PERMIT,
		},
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			client, closeFn := start(t)
			defer closeFn()
			rc, err := client.Rotate(context.Background())
			if err != nil {
				t.Fatal(err)
			}
			if tt.req != nil {
				mustSendAndRecv(t, rc, tt.req)
				if !tt.probeBeforeFinalize {
					mustFinalize(t, rc)
				}
			}

			got, err := client.Probe(context.Background(), tt.probeReq)
			if d := errdiff.Check(err, tt.wantErr); d != "" {
				t.Errorf("Probe() unexpected err: %s", d)
			}
			if err != nil {
				return
			}
			if d := cmp.Diff(tt.want, got, protocmp.Transform()); d != "" {
				t.Errorf("Probe() unexpected diff: %s", d)
			}
		})
	}
}

func TestGet(t *testing.T) {
	tests := []struct {
		desc                string
		req                 *pathzpb.RotateRequest
		probeBeforeFinalize bool
		getReq              *pathzpb.GetRequest
		want                *pathzpb.GetResponse
		wantErr             string
	}{{
		desc:    "error instance not set",
		getReq:  &pathzpb.GetRequest{},
		wantErr: "unknown instance type",
	}, {
		desc: "error instance not ready",
		getReq: &pathzpb.GetRequest{
			PolicyInstance: pathzpb.PolicyInstance_POLICY_INSTANCE_ACTIVE,
		},
		wantErr: "requested policy instance is nil",
	}, {
		desc: "success sandbox",
		getReq: &pathzpb.GetRequest{
			PolicyInstance: pathzpb.PolicyInstance_POLICY_INSTANCE_SANDBOX,
		},
		req: &pathzpb.RotateRequest{
			RotateRequest: &pathzpb.RotateRequest_UploadRequest{
				UploadRequest: &pathzpb.UploadRequest{
					Version: "1",
					Policy: &pathzpb.AuthorizationPolicy{
						Rules: []*pathzpb.AuthorizationRule{{
							Path:      &gpb.Path{Elem: []*gpb.PathElem{{Name: "test"}}},
							Principal: &pathzpb.AuthorizationRule_User{User: "me"},
							Mode:      pathzpb.Mode_MODE_READ,
							Action:    pathzpb.Action_ACTION_PERMIT,
						}},
					},
				},
			},
		},
		probeBeforeFinalize: true,
		want: &pathzpb.GetResponse{
			Version: "1",
			Policy: &pathzpb.AuthorizationPolicy{
				Rules: []*pathzpb.AuthorizationRule{{
					Path:      &gpb.Path{Elem: []*gpb.PathElem{{Name: "test"}}},
					Principal: &pathzpb.AuthorizationRule_User{User: "me"},
					Mode:      pathzpb.Mode_MODE_READ,
					Action:    pathzpb.Action_ACTION_PERMIT,
				}},
			},
		},
	}, {
		desc: "success active",
		getReq: &pathzpb.GetRequest{
			PolicyInstance: pathzpb.PolicyInstance_POLICY_INSTANCE_ACTIVE,
		},
		req: &pathzpb.RotateRequest{
			RotateRequest: &pathzpb.RotateRequest_UploadRequest{
				UploadRequest: &pathzpb.UploadRequest{
					Version: "1",
					Policy: &pathzpb.AuthorizationPolicy{
						Rules: []*pathzpb.AuthorizationRule{{
							Path:      &gpb.Path{Elem: []*gpb.PathElem{{Name: "test"}}},
							Principal: &pathzpb.AuthorizationRule_User{User: "me"},
							Mode:      pathzpb.Mode_MODE_READ,
							Action:    pathzpb.Action_ACTION_PERMIT,
						}},
					},
				},
			},
		},
		want: &pathzpb.GetResponse{
			Version: "1",
			Policy: &pathzpb.AuthorizationPolicy{
				Rules: []*pathzpb.AuthorizationRule{{
					Path:      &gpb.Path{Elem: []*gpb.PathElem{{Name: "test"}}},
					Principal: &pathzpb.AuthorizationRule_User{User: "me"},
					Mode:      pathzpb.Mode_MODE_READ,
					Action:    pathzpb.Action_ACTION_PERMIT,
				}},
			},
		},
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			client, closeFn := start(t)
			defer closeFn()
			rc, err := client.Rotate(context.Background())
			if err != nil {
				t.Fatal(err)
			}
			if tt.req != nil {
				mustSendAndRecv(t, rc, tt.req)
				if !tt.probeBeforeFinalize {
					mustFinalize(t, rc)
				}
			}

			got, err := client.Get(context.Background(), tt.getReq)
			if d := errdiff.Check(err, tt.wantErr); d != "" {
				t.Errorf("Probe() unexpected err: %s", d)
			}
			if err != nil {
				return
			}
			if d := cmp.Diff(tt.want, got, protocmp.Transform()); d != "" {
				t.Errorf("Probe() unexpected diff: %s", d)
			}
		})
	}
}

func mustSendAndRecv(t testing.TB, rc pathzpb.Pathz_RotateClient, req *pathzpb.RotateRequest) {
	t.Helper()
	if err := rc.Send(req); err != nil {
		t.Fatalf("failed to send: %v", err)
	}
	if _, err := rc.Recv(); err != nil {
		t.Fatalf("failed to send: %v", err)
	}
}

func mustFinalize(t testing.TB, rc pathzpb.Pathz_RotateClient) {
	t.Helper()
	mustSendAndRecv(t, rc, &pathzpb.RotateRequest{RotateRequest: &pathzpb.RotateRequest_FinalizeRotation{}})
}

func start(t testing.TB) (pathzpb.PathzClient, func()) {
	t.Helper()
	pathzServer := &Server{}

	s := grpc.NewServer()
	pathzpb.RegisterPathzServer(s, pathzServer)

	l, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatalf("failed to start listener: %v", err)
	}

	go s.Serve(l)

	conn, err := grpc.Dial(l.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("failed dial server: %v", err)
	}
	return pathzpb.NewPathzClient(conn), func() { s.Stop() }
}
