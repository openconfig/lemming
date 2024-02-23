// Copyright 2022 Google LLC
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

package reconciler

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/openconfig/gnmi/errdiff"
	"github.com/openconfig/ygnmi/ygnmi"
	"google.golang.org/protobuf/testing/protocmp"

	"github.com/openconfig/lemming/gnmi/oc"
	"github.com/openconfig/lemming/gnmi/oc/ocpath"

	gpb "github.com/openconfig/gnmi/proto/gnmi"
)

func TestWithStart(t *testing.T) {
	tests := []struct {
		desc string
		rec  Reconciler
		want string
	}{{
		desc: "single starter",
		rec: (&Builder{}).WithStart(func(context.Context, *ygnmi.Client) error {
			return fmt.Errorf("start 1 err")
		}).Build(),
		want: "start 1 err",
	}, {
		desc: "2 start funcs",
		rec: (&Builder{}).WithStart(func(context.Context, *ygnmi.Client) error {
			return fmt.Errorf("start 1 err")
		}).WithStart(func(context.Context, *ygnmi.Client) error {
			return fmt.Errorf("start 2 err")
		}).Build(),
		want: "start 1 err, start 2 err",
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			got := tt.rec.Start(context.Background(), nil, "")
			if diff := errdiff.Check(got, tt.want); diff != "" {
				t.Fatalf("WithStart failed: got %q want %q", got, tt.want)
			}
		})
	}
}

func TestWithValidator(t *testing.T) {
	tests := []struct {
		desc      string
		rec       Reconciler
		want      string
		wantPaths []ygnmi.PathStruct
	}{{
		desc: "single validator",
		rec: (&Builder{}).WithValidator([]ygnmi.PathStruct{ocpath.Root()}, func(*oc.Root) error {
			return fmt.Errorf("validator 1")
		}).Build(),
		want:      "validator 1",
		wantPaths: []ygnmi.PathStruct{ocpath.Root()},
	}, {
		desc: "2 validate funcs",
		rec: (&Builder{}).WithValidator([]ygnmi.PathStruct{ocpath.Root()}, func(*oc.Root) error {
			return fmt.Errorf("validator 1")
		}).WithValidator([]ygnmi.PathStruct{ocpath.Root().InterfaceAny()}, func(*oc.Root) error {
			return fmt.Errorf("validator 2")
		}).Build(),
		want:      "validator 1, validator 2",
		wantPaths: []ygnmi.PathStruct{ocpath.Root(), ocpath.Root().InterfaceAny()},
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			gotPaths := tt.rec.ValidationPaths()
			if diff := cmp.Diff(resolvePaths(t, gotPaths), resolvePaths(t, tt.wantPaths), protocmp.Transform()); diff != "" {
				t.Fatalf("WithValidate unexpected paths: %s", diff)
			}
			got := tt.rec.Validate(nil)
			if diff := errdiff.Check(got, tt.want); diff != "" {
				t.Fatalf("WithValidate unexpected error: %s", diff)
			}
		})
	}
}

func resolvePaths(t testing.TB, paths []ygnmi.PathStruct) []*gpb.Path {
	t.Helper()
	protoPaths := make([]*gpb.Path, len(paths))
	for i, p := range paths {
		p, _, err := ygnmi.ResolvePath(p)
		if err != nil {
			t.Fatal(err)
		}
		protoPaths[i] = p
	}
	return protoPaths
}
