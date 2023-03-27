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

package stats

import (
	"testing"
)

const (
	testCounter1 = 1 + iota
	testCounter2
)

func TestStat(t *testing.T) {
	counters := []EntryDesc{
		{testCounter1, "testCounter1"},
		{testCounter2, "testCounter2"},
	}
	desc := "test metric"
	st, err := New(desc, counters...)
	if err != nil {
		t.Fatalf("New(%v,%v) = %v, %v, want _, nil", desc, counters, st, err)
	}
	all := st.GetAll()
	if len(all) != len(st.statMap) {
		t.Errorf("Check error: GetAll() got %v items, want %v items", len(all), len(st.statMap))
	}
	for k, v := range all {
		if want, err := st.Get(k); want != v || err != nil {
			t.Errorf("Get(%v) = %v, want %v", k, v, want)
		}
	}

	regTests := []struct {
		in      []EntryDesc
		wantErr bool
	}{
		{
			in: counters,
		},
		{
			in: []EntryDesc{
				{testCounter1, "testCounter1"},
				{testCounter1, "testCounter1"},
			},
			wantErr: true,
		},
	}
	for _, tt := range regTests {
		if got, err := New(desc, tt.in...); tt.wantErr == (err == nil) {
			t.Fatalf("New(%v, %v) = %v, want error(%v)", desc, tt.in, got, tt.wantErr)
		}
	}

	testCounters := []struct {
		id      int
		name    string
		wantErr bool
	}{
		{
			id:   testCounter1,
			name: "testCounter1",
		},
		{
			id:   testCounter2,
			name: "testCounter2",
		},
		{
			id:      3,
			name:    "",
			wantErr: true,
		},
	}
	for _, tt := range testCounters {
		before, err := st.Get(tt.id)
		if tt.wantErr == (err == nil) {
			t.Errorf("stat.Get(%v) got _, %v, want _, error(%v)", tt.id, err, tt.wantErr)
		}
		want := before + 1
		st.Add(tt.id, 1)
		if got, _ := st.Get(tt.id); !tt.wantErr && got != want {
			t.Errorf("stat.Get(%v) got %v, nil, want %v, nil", tt.id, got, want)
		}
		if got, _ := st.GetStatName(tt.id); !tt.wantErr && got != tt.name {
			t.Errorf("stat.GetStatName(%v) got %v, nil, want %v, nil", tt.id, got, tt.name)
		}
	}

	functionTest := []struct {
		id      int
		in      int64
		wantErr bool
	}{
		{1, 100, false},
		{100, 200, true}, // this is to test illegal id won't crash
	}
	for _, tt := range functionTest {
		st.Update(tt.id, tt.in)
		got, err := st.Get(tt.id)
		if tt.wantErr == (err == nil) {
			t.Errorf("st.Get(%v) got _, %v, want _, error(%v)", tt.id, err, tt.wantErr)
		}
		if !tt.wantErr && got != tt.in {
			t.Errorf("st.Get(%v) got %v, nil, want %v, nil", tt.id, got, tt.in)
		}
	}
}
