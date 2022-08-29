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

// Package fwdset implements a set of members.
package fwdset

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdcontext"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdobject"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

// A Set is a forwarding object that represents a set of members. Each
// member is a byte slice, and may represent things like packet field values.
// The set is used to perform one-of matches i.e. it answers the query
// if a given slice of bytes is one-of a given set of byte slices.
type Set struct {
	fwdobject.Base
	members map[string]bool
}

// String returns all the members in a ';' delimited string format. The bytes
// are printed in hex.
func (c *Set) String() string {
	list := make([]string, 0, 1+len(c.members))
	for m := range c.members {
		list = append(list, hex.EncodeToString([]byte(m)))
	}
	return strings.Join(append(list, ""), ";")
}

// Update updates the members in a set.
func (c *Set) Update(members [][]byte) {
	t := make(map[string]bool)
	for _, m := range members {
		b := string(m)
		t[b] = true
	}
	c.members = t
}

// Contains returns true if the specific string exists in the set.
func (c *Set) Contains(member []byte) bool {
	return c.members[string(member)]
}

// New creates a new set object with the specified set id.
func New(ctx *fwdcontext.Context, id *fwdpb.SetId) (*Set, error) {
	if id == nil {
		return nil, errors.New("fwdset: find failed, no set specified")
	}
	c := &Set{
		members: make(map[string]bool),
	}
	if err := ctx.Objects.Insert(c, id.GetObjectId()); err != nil {
		return nil, err
	}
	return c, nil
}

// Find finds a set.
func Find(ctx *fwdcontext.Context, id *fwdpb.SetId) (*Set, error) {
	if id == nil {
		return nil, errors.New("fwdset: find failed, no set specified")
	}
	object, err := ctx.Objects.FindID(id.GetObjectId())
	if err != nil {
		return nil, err
	}
	if c, ok := object.(*Set); ok {
		return c, nil
	}
	return nil, fmt.Errorf("fwdset: find failed, %v is not a set", id)
}
