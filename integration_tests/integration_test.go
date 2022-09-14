// Copyright 2022 Google LLC
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

package integration_test

import (
	"context"
	"testing"
	"time"

	"github.com/openconfig/gribigo/chk"
	"github.com/openconfig/gribigo/constants"
	"github.com/openconfig/gribigo/fluent"
	"github.com/openconfig/ondatra"
	kinit "github.com/openconfig/ondatra/knebind/init"
)

func TestMain(m *testing.M) {
	ondatra.RunTests(m, kinit.Init)
}

// awaitTimeout calls a fluent client Await, adding a timeout to the context.
func awaitTimeout(ctx context.Context, c *fluent.GRIBIClient, t testing.TB, timeout time.Duration) error {
	subctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return c.Await(subctx, t)
}

func TestGRIBI(t *testing.T) {
	dut := ondatra.DUT(t, "dut")

	// Dial gRIBI
	ctx := context.Background()
	gribic := dut.RawAPIs().GRIBI().Default(t)

	// Configure the gRIBI client c with election ID of 10.
	c := fluent.NewClient()

	c.Connection().WithStub(gribic).WithInitialElectionID(10, 0).
		WithRedundancyMode(fluent.ElectedPrimaryClient).
		WithPersistence()

	c.Start(context.Background(), t)
	defer c.Stop(t)
	c.StartSending(context.Background(), t)
	if err := awaitTimeout(ctx, c, t, time.Minute); err != nil {
		t.Fatalf("Await got error during session negotiation for c: %v", err)
	}

	// Add an IPv4Entry for 1.1.1.1/30 pointing to 192.168.0.1.
	// This is similar to the first subtest in TE2.1.
	c.Modify().AddEntry(t,
		fluent.NextHopEntry().
			WithNetworkInstance("DEFAULT").
			WithIndex(1).
			WithIPAddress("192.168.0.1"))

	c.Modify().AddEntry(t,
		fluent.NextHopGroupEntry().
			WithNetworkInstance("DEFAULT").
			WithID(1).
			AddNextHop(1, 1))

	c.Modify().AddEntry(t,
		fluent.IPv4Entry().
			WithPrefix("1.1.1.1/30").
			WithNetworkInstance("DEFAULT").
			WithNextHopGroup(1))

	if err := awaitTimeout(ctx, c, t, time.Minute); err != nil {
		t.Fatalf("Could not program entries via c, got err: %v", err)
	}

	// This check assumes that 192.168.0.1 is reachable.
	chk.HasResult(t, c.Results(t),
		fluent.OperationResult().
			WithIPv4Operation("1.1.1.1/30").
			WithOperationType(constants.Add).
			WithProgrammingResult(fluent.InstalledInRIB).
			AsResult(),
		chk.IgnoreOperationID(),
	)
}
