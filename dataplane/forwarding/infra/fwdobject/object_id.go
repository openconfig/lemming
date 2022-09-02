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

package fwdobject

import (
	"fmt"
)

// NID is a numeric identifier of a forwarding object within a Lucius context.
type NID uint64

// InvalidNID represents an invalid NID.
const InvalidNID = NID(0)

// NIDPool is a pool used to generate object NIDs. It is not safe for
// simultaneous use by multiple goroutines.
type NIDPool struct {
	last uint64 // The last allocated id.
	free []NID  // List of free identifiers.
	max  uint64 // Largest NID that can be allocated.
}

// NewNIDPool creates a pool with the specified number of object NIDs.
func NewNIDPool(count uint64) *NIDPool {
	return &NIDPool{
		max: count,
	}
}

// Acquire obtains a new object id.
func (p *NIDPool) Acquire() (NID, error) {
	if n := len(p.free); n > 0 {
		id := p.free[n-1]
		p.free = p.free[:n-1]
		return id, nil
	}
	if p.last >= p.max {
		return 0, fmt.Errorf("fwdobject: acquire failed, no free ids")
	}
	p.last++
	return NID(p.last), nil
}

// Release releases an object id into the pool.
func (p *NIDPool) Release(id NID) {
	p.free = append(p.free, id)
}
