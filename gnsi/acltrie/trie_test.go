// Copyright 2023 Google LLC
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

package acltrie

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/openconfig/gnmi/errdiff"
	"github.com/openconfig/ygot/ygot"

	gpb "github.com/openconfig/gnmi/proto/gnmi"
	pathzpb "github.com/openconfig/gnsi/pathz"
)

func mustPath(s string) *gpb.Path {
	p, err := ygot.StringToStructuredPath(s)
	if err != nil {
		panic(fmt.Sprintf("cannot parse subscription path %s, %v", s, err))
	}
	return p
}

func TestInsert(t *testing.T) {
	tests := []struct {
		desc        string
		wantErr     string
		initialRule *pathzpb.AuthorizationRule
		testRule    *pathzpb.AuthorizationRule
	}{{
		desc: "success empty trie",
		testRule: &pathzpb.AuthorizationRule{
			Path:      mustPath("/bar"),
			Principal: &pathzpb.AuthorizationRule_User{User: "bob"},
			Mode:      pathzpb.Mode_MODE_READ,
			Action:    pathzpb.Action_ACTION_PERMIT,
		},
	}, {
		desc: "success different path no list",
		initialRule: &pathzpb.AuthorizationRule{
			Path:      mustPath("/foo"),
			Principal: &pathzpb.AuthorizationRule_User{User: "bob"},
			Mode:      pathzpb.Mode_MODE_READ,
			Action:    pathzpb.Action_ACTION_PERMIT,
		},
		testRule: &pathzpb.AuthorizationRule{
			Path:      mustPath("/bar"),
			Principal: &pathzpb.AuthorizationRule_User{User: "bob"},
			Mode:      pathzpb.Mode_MODE_READ,
			Action:    pathzpb.Action_ACTION_PERMIT,
		},
	}, {
		desc: "success path with different list keys",
		initialRule: &pathzpb.AuthorizationRule{
			Path:      mustPath("/foo[a=1]"),
			Principal: &pathzpb.AuthorizationRule_User{User: "bob"},
			Mode:      pathzpb.Mode_MODE_READ,
			Action:    pathzpb.Action_ACTION_PERMIT,
		},
		testRule: &pathzpb.AuthorizationRule{
			Path:      mustPath("/foo[a=2]"),
			Principal: &pathzpb.AuthorizationRule_User{User: "bob"},
			Mode:      pathzpb.Mode_MODE_READ,
			Action:    pathzpb.Action_ACTION_PERMIT,
		},
	}, {
		desc: "success path with wildcard and non-wildcard list keys",
		initialRule: &pathzpb.AuthorizationRule{
			Path:      mustPath("/foo[a=1]"),
			Principal: &pathzpb.AuthorizationRule_User{User: "bob"},
			Mode:      pathzpb.Mode_MODE_READ,
			Action:    pathzpb.Action_ACTION_PERMIT,
		},
		testRule: &pathzpb.AuthorizationRule{
			Path:      mustPath("/foo[a=*]"),
			Principal: &pathzpb.AuthorizationRule_User{User: "bob"},
			Mode:      pathzpb.Mode_MODE_READ,
			Action:    pathzpb.Action_ACTION_PERMIT,
		},
	}, {
		desc: "success path with wildcard and non-wildcard list keys",
		initialRule: &pathzpb.AuthorizationRule{
			Path:      mustPath("/foo[a=*]"),
			Principal: &pathzpb.AuthorizationRule_User{User: "bob"},
			Mode:      pathzpb.Mode_MODE_READ,
			Action:    pathzpb.Action_ACTION_PERMIT,
		},
		testRule: &pathzpb.AuthorizationRule{
			Path:      mustPath("/foo[a=1]"),
			Principal: &pathzpb.AuthorizationRule_User{User: "bob"},
			Mode:      pathzpb.Mode_MODE_READ,
			Action:    pathzpb.Action_ACTION_PERMIT,
		},
	}, {
		desc: "success path with the same prefix",
		initialRule: &pathzpb.AuthorizationRule{
			Path:      mustPath("/foo[a=*]"),
			Principal: &pathzpb.AuthorizationRule_User{User: "bob"},
			Mode:      pathzpb.Mode_MODE_READ,
			Action:    pathzpb.Action_ACTION_PERMIT,
		},
		testRule: &pathzpb.AuthorizationRule{
			Path:      mustPath("/foo[a=*]/bar"),
			Principal: &pathzpb.AuthorizationRule_User{User: "bob"},
			Mode:      pathzpb.Mode_MODE_READ,
			Action:    pathzpb.Action_ACTION_PERMIT,
		},
	}, {
		desc: "success list variable amount of wildcards",
		initialRule: &pathzpb.AuthorizationRule{
			Path:      mustPath("/foo[a=*][b=2]"),
			Principal: &pathzpb.AuthorizationRule_User{User: "bob"},
			Mode:      pathzpb.Mode_MODE_READ,
			Action:    pathzpb.Action_ACTION_PERMIT,
		},
		testRule: &pathzpb.AuthorizationRule{
			Path:      mustPath("/foo[a=1][b=2]"),
			Principal: &pathzpb.AuthorizationRule_User{User: "bob"},
			Mode:      pathzpb.Mode_MODE_READ,
			Action:    pathzpb.Action_ACTION_PERMIT,
		},
	}, {
		desc: "success same path, different user",
		initialRule: &pathzpb.AuthorizationRule{
			Path:      mustPath("/foo"),
			Principal: &pathzpb.AuthorizationRule_User{User: "bob"},
			Mode:      pathzpb.Mode_MODE_READ,
			Action:    pathzpb.Action_ACTION_PERMIT,
		},
		testRule: &pathzpb.AuthorizationRule{
			Path:      mustPath("/foo"),
			Principal: &pathzpb.AuthorizationRule_User{User: "testuser"},
			Mode:      pathzpb.Mode_MODE_READ,
			Action:    pathzpb.Action_ACTION_PERMIT,
		},
	}, {
		desc: "success same path, different mode",
		initialRule: &pathzpb.AuthorizationRule{
			Path:      mustPath("/foo"),
			Principal: &pathzpb.AuthorizationRule_User{User: "bob"},
			Mode:      pathzpb.Mode_MODE_READ,
			Action:    pathzpb.Action_ACTION_PERMIT,
		},
		testRule: &pathzpb.AuthorizationRule{
			Path:      mustPath("/foo"),
			Principal: &pathzpb.AuthorizationRule_User{User: "bob"},
			Mode:      pathzpb.Mode_MODE_WRITE,
			Action:    pathzpb.Action_ACTION_PERMIT,
		},
	}, {
		desc: "success same path, group",
		initialRule: &pathzpb.AuthorizationRule{
			Path:      mustPath("/foo"),
			Principal: &pathzpb.AuthorizationRule_User{User: "bob"},
			Mode:      pathzpb.Mode_MODE_READ,
			Action:    pathzpb.Action_ACTION_PERMIT,
		},
		testRule: &pathzpb.AuthorizationRule{
			Path:      mustPath("/foo"),
			Principal: &pathzpb.AuthorizationRule_Group{Group: "bob"},
			Mode:      pathzpb.Mode_MODE_READ,
			Action:    pathzpb.Action_ACTION_PERMIT,
		},
	}, {
		desc: "success ambiguous list keys of different length",
		initialRule: &pathzpb.AuthorizationRule{
			Path:      mustPath("/foo[a=*][b=2]/bar"),
			Principal: &pathzpb.AuthorizationRule_User{User: "bob"},
			Mode:      pathzpb.Mode_MODE_READ,
			Action:    pathzpb.Action_ACTION_PERMIT,
		},
		testRule: &pathzpb.AuthorizationRule{
			Path:      mustPath("/foo[a=1][b=*]"),
			Principal: &pathzpb.AuthorizationRule_User{User: "bob"},
			Mode:      pathzpb.Mode_MODE_READ,
			Action:    pathzpb.Action_ACTION_PERMIT,
		},
	}, {
		desc: "failure no action",
		testRule: &pathzpb.AuthorizationRule{
			Path:      mustPath("/bar"),
			Principal: &pathzpb.AuthorizationRule_User{User: "bob"},
			Mode:      pathzpb.Mode_MODE_READ,
		},
		wantErr: "action unspecified",
	}, {
		desc: "failure no mode",
		testRule: &pathzpb.AuthorizationRule{
			Path:      mustPath("/bar"),
			Principal: &pathzpb.AuthorizationRule_User{User: "bob"},
			Action:    pathzpb.Action_ACTION_DENY,
		},
		wantErr: "mode unspecified",
	}, {
		desc: "failure no principal",
		testRule: &pathzpb.AuthorizationRule{
			Path:   mustPath("/bar"),
			Action: pathzpb.Action_ACTION_DENY,
			Mode:   pathzpb.Mode_MODE_READ,
		},
		wantErr: "principal unset",
	}, {
		desc: "failure bad path",
		testRule: &pathzpb.AuthorizationRule{
			Path:      mustPath("//bar"),
			Principal: &pathzpb.AuthorizationRule_User{User: "bob"},
			Mode:      pathzpb.Mode_MODE_READ,
			Action:    pathzpb.Action_ACTION_PERMIT,
		},
		wantErr: "empty name",
	}, {
		desc: "failure wildcard Name",
		initialRule: &pathzpb.AuthorizationRule{
			Path:      mustPath("/foo"),
			Principal: &pathzpb.AuthorizationRule_User{User: "bob"},
			Mode:      pathzpb.Mode_MODE_READ,
			Action:    pathzpb.Action_ACTION_PERMIT,
		},
		testRule: &pathzpb.AuthorizationRule{
			Path:      mustPath("/*/bar"),
			Principal: &pathzpb.AuthorizationRule_User{User: "bob"},
			Mode:      pathzpb.Mode_MODE_READ,
			Action:    pathzpb.Action_ACTION_PERMIT,
		},
		wantErr: "wildcard path names",
	}, {
		desc: "failure duplicate rule",
		initialRule: &pathzpb.AuthorizationRule{
			Path:      mustPath("/foo/bar"),
			Principal: &pathzpb.AuthorizationRule_User{User: "bob"},
			Mode:      pathzpb.Mode_MODE_READ,
			Action:    pathzpb.Action_ACTION_PERMIT,
		},
		testRule: &pathzpb.AuthorizationRule{
			Path:      mustPath("/foo/bar"),
			Principal: &pathzpb.AuthorizationRule_User{User: "bob"},
			Mode:      pathzpb.Mode_MODE_READ,
			Action:    pathzpb.Action_ACTION_PERMIT,
		},
		wantErr: "policy already contains action for principal",
	}, {
		desc: "failure duplicate rule list keys",
		initialRule: &pathzpb.AuthorizationRule{
			Path:      mustPath("/foo[a=1]/bar"),
			Principal: &pathzpb.AuthorizationRule_User{User: "bob"},
			Mode:      pathzpb.Mode_MODE_READ,
			Action:    pathzpb.Action_ACTION_PERMIT,
		},
		testRule: &pathzpb.AuthorizationRule{
			Path:      mustPath("/foo[a=1]/bar"),
			Principal: &pathzpb.AuthorizationRule_User{User: "bob"},
			Mode:      pathzpb.Mode_MODE_READ,
			Action:    pathzpb.Action_ACTION_PERMIT,
		},
		wantErr: "policy already contains action for principal",
	}, {
		desc: "failure ambiguous list keys",
		initialRule: &pathzpb.AuthorizationRule{
			Path:      mustPath("/foo[a=*][b=2]"),
			Principal: &pathzpb.AuthorizationRule_User{User: "bob"},
			Mode:      pathzpb.Mode_MODE_READ,
			Action:    pathzpb.Action_ACTION_PERMIT,
		},
		testRule: &pathzpb.AuthorizationRule{
			Path:      mustPath("/foo[a=1][b=*]"),
			Principal: &pathzpb.AuthorizationRule_User{User: "bob"},
			Mode:      pathzpb.Mode_MODE_READ,
			Action:    pathzpb.Action_ACTION_PERMIT,
		},
		wantErr: "policy path conflict",
	}, {
		desc: "failure ambiguous list keys implicit wildcard",
		initialRule: &pathzpb.AuthorizationRule{
			Path:      mustPath("/foo[b=2]"),
			Principal: &pathzpb.AuthorizationRule_User{User: "bob"},
			Mode:      pathzpb.Mode_MODE_READ,
			Action:    pathzpb.Action_ACTION_PERMIT,
		},
		testRule: &pathzpb.AuthorizationRule{
			Path:      mustPath("/foo[a=1]"),
			Principal: &pathzpb.AuthorizationRule_User{User: "bob"},
			Mode:      pathzpb.Mode_MODE_READ,
			Action:    pathzpb.Action_ACTION_PERMIT,
		},
		wantErr: "policy path conflict",
	}, {
		desc: "failure ambiguous list keys implicit wildcard",
		initialRule: &pathzpb.AuthorizationRule{
			Path:      mustPath("/foo[a=2]"),
			Principal: &pathzpb.AuthorizationRule_User{User: "bob"},
			Mode:      pathzpb.Mode_MODE_READ,
			Action:    pathzpb.Action_ACTION_PERMIT,
		},
		testRule: &pathzpb.AuthorizationRule{
			Path:      mustPath("/foo[b=1]"),
			Principal: &pathzpb.AuthorizationRule_User{User: "bob"},
			Mode:      pathzpb.Mode_MODE_READ,
			Action:    pathzpb.Action_ACTION_PERMIT,
		},
		wantErr: "policy path conflict",
	}, {
		desc: "failure ambiguous list keys with matching suffix",
		initialRule: &pathzpb.AuthorizationRule{
			Path:      mustPath("/foo[a=*][b=2]/c/d"),
			Principal: &pathzpb.AuthorizationRule_User{User: "bob"},
			Mode:      pathzpb.Mode_MODE_READ,
			Action:    pathzpb.Action_ACTION_PERMIT,
		},
		testRule: &pathzpb.AuthorizationRule{
			Path:      mustPath("/foo[a=1][b=*]/c/d"),
			Principal: &pathzpb.AuthorizationRule_User{User: "bob"},
			Mode:      pathzpb.Mode_MODE_READ,
			Action:    pathzpb.Action_ACTION_PERMIT,
		},
		wantErr: "policy path conflict",
	}, {
		desc: "failure ambiguous list keys across multiple lists",
		initialRule: &pathzpb.AuthorizationRule{
			Path:      mustPath("/foo[a=1][b=*]/bar[c=1]"),
			Principal: &pathzpb.AuthorizationRule_User{User: "bob"},
			Mode:      pathzpb.Mode_MODE_READ,
			Action:    pathzpb.Action_ACTION_PERMIT,
		},
		testRule: &pathzpb.AuthorizationRule{
			Path:      mustPath("/foo[a=1][b=2]/bar[c=*]"),
			Principal: &pathzpb.AuthorizationRule_User{User: "bob"},
			Mode:      pathzpb.Mode_MODE_READ,
			Action:    pathzpb.Action_ACTION_PERMIT,
		},
		wantErr: "policy path conflict",
	}, {
		desc: "failure ambiguous list keys across multiple lists",
		initialRule: &pathzpb.AuthorizationRule{
			Path:      mustPath("/foo[a=*][b=*]/bar[c=1]"),
			Principal: &pathzpb.AuthorizationRule_User{User: "bob"},
			Mode:      pathzpb.Mode_MODE_READ,
			Action:    pathzpb.Action_ACTION_PERMIT,
		},
		testRule: &pathzpb.AuthorizationRule{
			Path:      mustPath("/foo[a=1][b=*]/bar"),
			Principal: &pathzpb.AuthorizationRule_User{User: "bob"},
			Mode:      pathzpb.Mode_MODE_READ,
			Action:    pathzpb.Action_ACTION_PERMIT,
		},
		wantErr: "policy path conflict",
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			trie := &Trie{}
			if tt.initialRule != nil {
				err := trie.Insert(tt.initialRule)
				if err != nil {
					t.Fatalf("Insert() failed to setup initial trie: %v", err)
				}
			}
			gotErr := trie.Insert(tt.testRule)
			if d := errdiff.Check(gotErr, tt.wantErr); d != "" {
				t.Errorf("Insert() unexpected err: %s", d)
			}
		})
	}
}

func TestProbe(t *testing.T) {
	tests := []struct {
		desc         string
		initialRules []*pathzpb.AuthorizationRule
		path         *gpb.Path
		want         pathzpb.Action
	}{{
		desc: "simple user  match",
		initialRules: []*pathzpb.AuthorizationRule{{
			Path:      mustPath("/foo/bar"),
			Principal: &pathzpb.AuthorizationRule_User{User: "bob"},
			Mode:      pathzpb.Mode_MODE_READ,
			Action:    pathzpb.Action_ACTION_PERMIT,
		}},
		path: mustPath("/foo/bar"),
		want: pathzpb.Action_ACTION_PERMIT,
	}, {
		desc: "simple group match",
		initialRules: []*pathzpb.AuthorizationRule{{
			Path:      mustPath("/foo/bar"),
			Principal: &pathzpb.AuthorizationRule_Group{Group: "admin"},
			Mode:      pathzpb.Mode_MODE_READ,
			Action:    pathzpb.Action_ACTION_PERMIT,
		}},
		path: mustPath("/foo/bar"),
		want: pathzpb.Action_ACTION_PERMIT,
	}, {
		desc: "rule is a prefix match",
		initialRules: []*pathzpb.AuthorizationRule{{
			Path:      mustPath("/foo"),
			Principal: &pathzpb.AuthorizationRule_User{User: "bob"},
			Mode:      pathzpb.Mode_MODE_READ,
			Action:    pathzpb.Action_ACTION_PERMIT,
		}},
		path: mustPath("/foo/bar"),
		want: pathzpb.Action_ACTION_PERMIT,
	}, {
		desc: "explicit key against partial wildcard",
		initialRules: []*pathzpb.AuthorizationRule{{
			Path:      mustPath("/foo[a=1][b=*]"),
			Principal: &pathzpb.AuthorizationRule_User{User: "bob"},
			Mode:      pathzpb.Mode_MODE_READ,
			Action:    pathzpb.Action_ACTION_PERMIT,
		}},
		path: mustPath("/foo[a=1][b=2]"),
		want: pathzpb.Action_ACTION_PERMIT,
	}, {
		desc: "prefer longer path",
		initialRules: []*pathzpb.AuthorizationRule{{
			Path:      mustPath("/foo[a=*]/bar"),
			Principal: &pathzpb.AuthorizationRule_User{User: "bob"},
			Mode:      pathzpb.Mode_MODE_READ,
			Action:    pathzpb.Action_ACTION_PERMIT,
		}, {
			Path:      mustPath("/foo[a=*]"),
			Principal: &pathzpb.AuthorizationRule_User{User: "bob"},
			Mode:      pathzpb.Mode_MODE_READ,
			Action:    pathzpb.Action_ACTION_DENY,
		}},
		path: mustPath("/foo[a=1]/bar"),
		want: pathzpb.Action_ACTION_PERMIT,
	}, {
		desc: "prefer definite key",
		initialRules: []*pathzpb.AuthorizationRule{{
			Path:      mustPath("/foo[a=*]"),
			Principal: &pathzpb.AuthorizationRule_User{User: "bob"},
			Mode:      pathzpb.Mode_MODE_READ,
			Action:    pathzpb.Action_ACTION_PERMIT,
		}, {
			Path:      mustPath("/foo[a=1]"),
			Principal: &pathzpb.AuthorizationRule_User{User: "bob"},
			Mode:      pathzpb.Mode_MODE_READ,
			Action:    pathzpb.Action_ACTION_DENY,
		}},
		path: mustPath("/foo[a=1]/bar"),
		want: pathzpb.Action_ACTION_DENY,
	}, {
		desc: "prefer user over group",
		initialRules: []*pathzpb.AuthorizationRule{{
			Path:      mustPath("/foo[a=1]"),
			Principal: &pathzpb.AuthorizationRule_Group{Group: "admin"},
			Mode:      pathzpb.Mode_MODE_READ,
			Action:    pathzpb.Action_ACTION_DENY,
		}, {
			Path:      mustPath("/foo[a=1]"),
			Principal: &pathzpb.AuthorizationRule_User{User: "bob"},
			Mode:      pathzpb.Mode_MODE_READ,
			Action:    pathzpb.Action_ACTION_PERMIT,
		}},
		path: mustPath("/foo[a=1]/bar"),
		want: pathzpb.Action_ACTION_PERMIT,
	}, {
		desc: "prefer deny over permit",
		initialRules: []*pathzpb.AuthorizationRule{{
			Path:      mustPath("/foo[a=1]"),
			Principal: &pathzpb.AuthorizationRule_Group{Group: "admin"},
			Mode:      pathzpb.Mode_MODE_READ,
			Action:    pathzpb.Action_ACTION_DENY,
		}, {
			Path:      mustPath("/foo[a=1]"),
			Principal: &pathzpb.AuthorizationRule_Group{Group: "reader"},
			Mode:      pathzpb.Mode_MODE_READ,
			Action:    pathzpb.Action_ACTION_PERMIT,
		}},
		path: mustPath("/foo[a=1]/bar"),
		want: pathzpb.Action_ACTION_DENY,
	}, {
		desc: "default policy",
		initialRules: []*pathzpb.AuthorizationRule{{
			Path:      mustPath("/foo"),
			Principal: &pathzpb.AuthorizationRule_User{User: "bob"},
			Mode:      pathzpb.Mode_MODE_READ,
			Action:    pathzpb.Action_ACTION_DENY,
		}},
		path: mustPath("/bar"),
		want: pathzpb.Action_ACTION_DENY,
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			trie := &Trie{
				memberships: map[string]map[string]bool{
					"admin": {
						"bob": true,
					},
					"reader": {
						"bob": true,
					},
				},
			}
			for _, rule := range tt.initialRules {
				err := trie.Insert(rule)
				if err != nil {
					t.Fatalf("Get() failed to setup initial trie: %v", err)
				}
			}
			got := trie.Probe(tt.path, "bob", pathzpb.Mode_MODE_READ)
			if d := cmp.Diff(tt.want, got); d != "" {
				t.Errorf("Get() unexpected diff: %s", d)
			}
		})
	}
}
