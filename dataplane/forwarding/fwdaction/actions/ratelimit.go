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
	"fmt"
	"sync"
	"time"

	"github.com/openconfig/lemming/dataplane/forwarding/fwdaction"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdattribute"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdcontext"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdobject"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

// AttrRatelimitAdvisory controls if the ratelimit is enforced or advisory
var AttrRatelimitAdvisory = fwdattribute.ID("RatelimitAdvisory")

func init() {
	fwdattribute.Register(AttrRatelimitAdvisory, "Does not enforce ratelimit if set to true")
}

// A ratelimit is an action that policies packets flowing through it.
//
// It implements a simplified token bucket configured with a rate and a burst.
// The rate determines how many tokens (i.e. bytes) are added to the bucket per
// second and the burst determines the maximum allowed size of the bucket.
// When a packet arrives, it is allowed to continue only if there are
// sufficent number of tokens in the bucket.
//
// The tokens in the bucket are evaluated when each packet arrives. Note that
// the number of tokens never exceeds the burst size. A packet that is allowed
// consumes tokens from the bucket equal to its length in bytes.
type ratelimit struct {
	burst uint64           // burst size in bytes
	rate  uint64           // rate in bytes per second
	clock func() time.Time // function used to access current Unix time

	mu     sync.Mutex
	last   time.Time // time of last update
	tokens uint64    // size of bucket in bytes
	update struct {
		tokens   uint64        // number of tokens added in the last evaluation
		interval time.Duration // interval used for the last evaluation
	}
	running bool // true if the ratelimit has processed at least one packet
}

// String formats the state of the action as a string.
func (r *ratelimit) String() string {
	return fmt.Sprintf("Type=%v;Rate=%v;Burst=%v;Tokens=%v;Last=%v;Update.Tokens=%v;Update.Interval=%v;Running=%v", fwdpb.ActionType_ACTION_TYPE_RATE, r.rate, r.burst, r.tokens, r.last, r.update.tokens, r.update.interval, r.running)
}

// Allowed evaluates the token bucket and returns true if a packet of the
// specified length is allowed.
func (r *ratelimit) Allowed(length uint64) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := r.clock() // current time

	// Setup the initial values. Update them if the ratelimit is running.
	tokens := r.burst                 // current number of tokens
	updateTokens := tokens            // change in tokens
	updateInterval := 0 * time.Second // time elapsed since last change in tokens
	if r.running {
		updateInterval = now.Sub(r.last)
		updateTokens = r.rate * uint64(updateInterval) / uint64(time.Second)
		tokens = r.tokens + updateTokens
		if tokens > r.burst {
			tokens = r.burst
		}
	}

	// Evaluate if the current packet has sufficient number of tokens.
	if length > tokens {
		return false
	}

	// Update the number of tokens after consuming the packet.
	r.tokens = tokens - length

	// Update the timestamp only if there was a change.
	if updateTokens != 0 {
		r.last = now
		r.update.tokens = updateTokens
		r.update.interval = updateInterval
	}
	r.running = true
	return true
}

// Process allows packet processing to continue if the bucket has sufficient
// tokens. If the attribute "RatelimitAdvisory" is set, the packet is not
// actually dropped, but only counted as ratelimited.
func (r *ratelimit) Process(packet fwdpacket.Packet, counters fwdobject.Counters) (fwdaction.Actions, fwdaction.State) {
	length := uint64(packet.Length())
	if !r.Allowed(length) {
		counters.Increment(fwdpb.CounterId_COUNTER_ID_RATELIMIT_PACKETS, 1)
		counters.Increment(fwdpb.CounterId_COUNTER_ID_RATELIMIT_OCTETS, uint32(length))

		a := packet.Attributes()
		if a != nil {
			if value, ok := a.Get(AttrRatelimitAdvisory); ok && value == "true" {
				return nil, fwdaction.CONTINUE
			}
		}
		return nil, fwdaction.DROP
	}
	return nil, fwdaction.CONTINUE
}

// A ratelimitBuilder builds ratelimit actions.
type ratelimitBuilder struct{}

// init registers a builder for the ratelimit action type.
func init() {
	fwdaction.Register(fwdpb.ActionType_ACTION_TYPE_RATE, ratelimitBuilder{})
}

// Build creates a new ratelimit action.
func (ratelimitBuilder) Build(desc *fwdpb.ActionDesc, ctx *fwdcontext.Context) (fwdaction.Action, error) {
	r, ok := desc.Action.(*fwdpb.ActionDesc_Rate)
	if !ok {
		return nil, fmt.Errorf("actions: Build for ratelimit action failed, missing desc")
	}

	return &ratelimit{
		last:    time.Now(),
		burst:   uint64(r.Rate.GetBurstBytes()),
		rate:    uint64(r.Rate.GetRateBps()),
		clock:   time.Now,
		running: false,
	}, nil
}
