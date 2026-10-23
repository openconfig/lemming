// Copyright 2026 Google LLC
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

//go:build linux

package dplanerc

import (
	"fmt"
	"testing"
)

func TestHostifName(t *testing.T) {
	tests := []struct {
		desc    string
		name    string
		ifIndex int
		want    string
	}{{
		desc:    "short name uses full suffix",
		name:    "eth1",
		ifIndex: 3,
		want:    "eth1-internal",
	}, {
		desc:    "longest name that fits full suffix",
		name:    "eth999",
		ifIndex: 1001,
		want:    "eth999-internal",
	}, {
		desc:    "four digit port uses short suffix",
		name:    "eth1000",
		ifIndex: 1002,
		want:    "eth1000-int",
	}, {
		desc:    "longest name that fits short suffix",
		name:    "Ethernet100",
		ifIndex: 102,
		want:    "Ethernet100-int",
	}, {
		desc:    "long name uses ifindex",
		name:    "Ethernet1000",
		ifIndex: 1234,
		want:    "if1234-int",
	}, {
		desc:    "long name with large ifindex",
		name:    "Ethernet1/1/1000",
		ifIndex: 999999999,
		want:    "if999999999-int",
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			got := hostifName(tt.name, tt.ifIndex)
			if got != tt.want {
				t.Errorf("hostifName(%q, %d) got %q, want %q", tt.name, tt.ifIndex, got, tt.want)
			}
			if len(got) > maxIfNameLen {
				t.Errorf("hostifName(%q, %d) got %q (len %d), exceeds max len %d", tt.name, tt.ifIndex, got, len(got), maxIfNameLen)
			}
			if !isHostifName(got) {
				t.Errorf("isHostifName(%q) got false, want true", got)
			}
		})
	}
}

func TestHostifNameUnique(t *testing.T) {
	seen := map[string]string{}
	ifIndex := 2
	for _, format := range []string{"eth%d", "Ethernet%d", "Ethernet1/%d"} {
		for i := 0; i < 4096; i++ {
			name := fmt.Sprintf(format, i)
			ifIndex++
			got := hostifName(name, ifIndex)
			if len(got) > maxIfNameLen {
				t.Fatalf("hostifName(%q, %d) got %q (len %d), exceeds max len %d", name, ifIndex, got, len(got), maxIfNameLen)
			}
			if prev, ok := seen[got]; ok {
				t.Fatalf("hostifName(%q, %d) = %q, same as for %q; want unique names", name, ifIndex, got, prev)
			}
			seen[got] = name
		}
	}
}

func TestIsHostifName(t *testing.T) {
	tests := []struct {
		in   string
		want bool
	}{
		{in: "eth0", want: false},
		{in: "eth1", want: false},
		{in: "Ethernet1000", want: false},
		{in: "eth1-internal", want: true},
		{in: "eth1000-int", want: true},
		{in: "if1234-int", want: true},
	}
	for _, tt := range tests {
		if got := isHostifName(tt.in); got != tt.want {
			t.Errorf("isHostifName(%q) got %v, want %v", tt.in, got, tt.want)
		}
	}
}
