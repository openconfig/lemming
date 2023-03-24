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

// Package stats implements Stats that collects data.
package stats

import (
	"fmt"
	"sync/atomic"
)

// statEntry is a structure to manage and export information like stats
// and other numeric metrics in transport.
type statEntry struct {
	name  string
	value int64
}

// Stats is a map of stats for the object.
type Stats struct {
	statMap map[int]*statEntry
}

// EntryDesc defines a pair of ID and Name for a stat field.
type EntryDesc struct {
	// ID is unique in a Stats instance, or Register() will reject.
	ID int
	// Name will be part of the variable name.
	Name string
}

// New returns a new Stats.
func New(desc string, stats ...EntryDesc) (*Stats, error) {
	s := &Stats{
		statMap: make(map[int]*statEntry),
	}
	for _, stat := range stats {
		if _, ok := s.statMap[stat.ID]; ok {
			return nil, fmt.Errorf("stats: %v is already in the Stats", stat.ID)
		}
		s.statMap[stat.ID] = &statEntry{value: 0, name: stat.Name}
	}
	return s, nil
}

// Get reads the value of the given stat atomically. Non-existing ID will
// return an error.
func (s *Stats) Get(id int) (int64, error) {
	if stat, ok := s.statMap[id]; ok {
		return atomic.LoadInt64(&stat.value), nil
	}
	return -1, fmt.Errorf("stats: non-existing stat %v value requested", id)
}

// GetStatName returns the string name of the stat ID. Non-existing ID will
// return an error.
func (s *Stats) GetStatName(id int) (string, error) {
	if stat, ok := s.statMap[id]; ok {
		return stat.name, nil
	}
	return "", fmt.Errorf("stats: non-existing stat %v name requested", id)
}

// Add adds the given delta to the stat.
func (s *Stats) Add(id int, delta int64) error {
	if stat, ok := s.statMap[id]; ok {
		atomic.AddInt64(&stat.value, delta)
		return nil
	}
	return fmt.Errorf("stats: non-existing stat %v increment requested", id)
}

// Update updates the stat with the value provided.
func (s *Stats) Update(id int, value int64) error {
	if stat, ok := s.statMap[id]; ok {
		atomic.StoreInt64(&stat.value, value)
		return nil
	}
	return fmt.Errorf("stats: non-existing stat %v update requested", id)
}

// GetAll returns all stats as a map of stat values, indexed by id.
func (s *Stats) GetAll() map[int]int64 {
	r := make(map[int]int64)
	for id, stat := range s.statMap {
		r[id] = atomic.LoadInt64(&stat.value)
	}
	return r
}
