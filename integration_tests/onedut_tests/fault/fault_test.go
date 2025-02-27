// Copyright 2025 Google LLC
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

package fault

import (
	"context"
	"testing"

	"github.com/openconfig/ondatra"
	"google.golang.org/grpc"

	"github.com/openconfig/lemming/internal/binding"
	faultpb "github.com/openconfig/lemming/proto/fault"
)

func TestMain(m *testing.M) {
	ondatra.RunTests(m, binding.KNE(".."))
}

type grpcDialer interface {
	DialGRPC(ctx context.Context, serviceName string, opts ...grpc.DialOption) (*grpc.ClientConn, error)
}

func TestFault(t *testing.T) {
	dut := ondatra.DUT(t, "dut")
	ld := dut.RawAPIs().BindingDUT().(grpcDialer)
	conn, err := ld.DialGRPC(context.Background(), "fault")
	if err != nil {
		t.Fatal(err)
	}
	faultpb.NewFaultInjectClient(conn)
}
