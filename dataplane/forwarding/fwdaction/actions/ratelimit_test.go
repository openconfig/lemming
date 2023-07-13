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

package actions

import (
	"testing"
	"time"

	"go.uber.org/mock/gomock"

	"github.com/openconfig/lemming/dataplane/forwarding/fwdaction"
	"github.com/openconfig/lemming/dataplane/forwarding/fwdaction/mock_fwdpacket"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdcontext"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdobject"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"

	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/arp"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/ethernet"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/metadata"
	_ "github.com/openconfig/lemming/dataplane/forwarding/protocol/opaque"
)

// TestRate tests the ratelimit action and builder.
func TestRate(t *testing.T) {
	// Create a controller for creating mock packets.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a forwarding context.
	ctx := fwdcontext.New("test", "fwd")

	type packet struct {
		length   int           // Packet length in bytes
		duration time.Duration // Time elapsed since the first packet
		drop     bool          // Indicates if the packet should be dropped
	}

	tests := []struct {
		rate    int32    // Rate in bytes per second
		burst   int32    // Burst in bytes
		packets []packet // sequence of packets processed
	}{
		{
			rate:  1000,
			burst: 100,
			packets: []packet{
				{
					length:   1000,
					duration: 0 * time.Second,
					drop:     true,
				},
				{
					length:   50,
					duration: 0 * time.Second,
					drop:     false,
				},
				{
					length:   50,
					duration: 0 * time.Second,
					drop:     false,
				},
				{
					length:   50,
					duration: 0 * time.Second,
					drop:     true,
				},
			},
		},
		{
			rate:  1000,
			burst: 100,
			packets: []packet{
				{
					length:   100,
					duration: 0 * time.Second,
					drop:     false,
				},
				{
					length:   100,
					duration: 99 * time.Millisecond,
					drop:     true,
				},
				{
					length:   100,
					duration: 100 * time.Millisecond,
					drop:     false,
				},
				{
					length:   101,
					duration: 500 * time.Millisecond,
					drop:     true,
				},
			},
		},
	}
	for pos, test := range tests {
		desc := fwdpb.ActionDesc{
			ActionType: fwdpb.ActionType_ACTION_TYPE_RATE,
		}
		rateDesc := fwdpb.RateActionDesc{
			RateBps:    test.rate,
			BurstBytes: test.burst,
		}
		desc.Action = &fwdpb.ActionDesc_Rate{
			Rate: &rateDesc,
		}
		action, err := fwdaction.New(&desc, ctx)
		if err != nil {
			t.Errorf("%d: NewAction failed, desc %v failed, err %v.", pos, &desc, err)
		}

		// Change the ratelimit action's clock function.
		rate := action.(*ratelimit)
		now := time.Now()

		for _, p := range test.packets {
			rate.clock = func() time.Time {
				return now.Add(p.duration)
			}
			var base fwdobject.Base
			if err := base.InitCounters("desc", fwdpb.CounterId_COUNTER_ID_RATELIMIT_PACKETS, fwdpb.CounterId_COUNTER_ID_RATELIMIT_OCTETS); err != nil {
				t.Fatalf("InitCounters failed, %v", err)
			}

			// Process a mock packet of the specified length.
			packet := mock_fwdpacket.NewMockPacket(ctrl)
			packet.EXPECT().Length().Return(p.length).AnyTimes()
			packet.EXPECT().Attributes().Return(nil).AnyTimes()
			next, state := action.Process(packet, &base)

			// Verify the result and the updated counters.
			var wantPkts, wantBytes uint64
			wantState := fwdaction.CONTINUE
			if p.drop {
				wantPkts = 1
				wantBytes = uint64(p.length)
				wantState = fwdaction.DROP
			}
			if next != nil {
				t.Errorf("%d: %v processing returned next actions. Got %v, want <nil>.", pos, action, next)
			}
			if state != wantState {
				t.Errorf("%d: %v processing returned bad state. Got %v, want %v.", pos, action, state, wantState)
			}
			counters := base.Counters()
			if counter, ok := counters[fwdpb.CounterId_COUNTER_ID_RATELIMIT_PACKETS]; !ok || counter.Value != wantPkts {
				t.Errorf("%d: %v processing returned invalid counter %v. Got %v, want %v.", pos, action, counter, counter.Value, wantPkts)
			}
			if counter, ok := counters[fwdpb.CounterId_COUNTER_ID_RATELIMIT_OCTETS]; !ok || counter.Value != wantBytes {
				t.Errorf("%d: %v processing returned invalid counter %v.  Got %v, want %v.", pos, action, counter, counter.Value, wantBytes)
			}
		}
	}
}
