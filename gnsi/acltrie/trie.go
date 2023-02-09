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

// Package acltrie contains a trie for gNMI path with a corresponding policy.
package acltrie

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	gpb "github.com/openconfig/gnmi/proto/gnmi"
	pathzpb "github.com/openconfig/gnsi/pathz"
)

// Trie is the root of the ACL trie.
type Trie struct {
	root *trieNode
}

type trieNode struct {
	// childrenByRank orders the children of the node by number the set leaves for the element.
	// This allows fast lookup of the best match: the one with the most number of set keys.
	childrenByRank []map[string]*trieNode
	elem           *gpb.PathElem
	hasPolicy      bool
	users          policies
	groups         policies
}

// Insert inserts a new policy into the trie.
func (t *Trie) Insert(r *pathzpb.AuthorizationRule) error {
	if t.root == nil {
		t.root = &trieNode{}
	}
	if r.Action == pathzpb.Action_ACTION_UNSPECIFIED {
		return fmt.Errorf("action unspecified")
	}
	if r.Mode == pathzpb.Mode_MODE_UNSPECIFIED {
		return fmt.Errorf("mode unspecified")
	}
	if r.GetGroup() == "" && r.GetUser() == "" {
		return fmt.Errorf("principal unset")
	}

	path := r.GetPath()
	node := t.root

	for i, elem := range path.Elem {
		if elem.Name == "*" {
			return fmt.Errorf("wildcard path names are not permitted")
		}

		// Normalize path string by pruning wildcard keys.
		pathStr, err := elemToString(elem.Name, elem.Key)
		if err != nil {
			return fmt.Errorf("invalid path element: %v", err)
		}

		// The rank of the child is the number of non-wildcards key (if any).
		rank := len(setKeys(elem.Key))
		if rank >= len(node.childrenByRank) {
			node.childrenByRank = append(node.childrenByRank, make([]map[string]*trieNode, rank-len(node.childrenByRank)+1)...)
		}
		if node.childrenByRank[rank] == nil {
			node.childrenByRank[rank] = make(map[string]*trieNode)
		}
		children := node.childrenByRank[rank]

		// If the node already exists, keep going.
		if _, ok := children[pathStr]; ok {
			node = children[pathStr]
			continue
		}

		// Before adding a new node, check if the path conflicts with another policy.
		for _, child := range children {
			if checkPathConflict(child.elem, elem) {
				return fmt.Errorf("policy path conflict with %v/%v", path.Elem[0:i], child.elem)
			}
		}

		children[pathStr] = &trieNode{
			elem:   elem,
			groups: policies{},
			users:  policies{},
		}

		node = children[pathStr]
	}
	principal := r.GetUser()
	policy := node.users
	if _, isUser := r.GetPrincipal().(*pathzpb.AuthorizationRule_User); !isUser {
		principal = r.GetGroup()
		policy = node.groups
	}

	if err := policy.insert(principal, r.GetMode(), r.GetAction()); err != nil {
		return fmt.Errorf("error inserting policy at %v: %v", r.Path, err)
	}
	node.hasPolicy = true

	return nil
}

// policies are map of mode to user to action.
type policies map[pathzpb.Mode]map[string]pathzpb.Action

// insert adds a principal for a given mode and action to policies, returning an error on duplicate entry.
func (p policies) insert(principal string, mode pathzpb.Mode, action pathzpb.Action) error {
	m, ok := p[mode]
	if !ok {
		p[mode] = make(map[string]pathzpb.Action)
		m = p[mode]
	}
	if _, ok := m[principal]; ok {
		return fmt.Errorf("policy already contains action for principal")
	}
	m[principal] = action

	return nil
}

// checkPathConflict returns an error if the candidate path overlaps with the accepted path.
func checkPathConflict(accepted, candidate *gpb.PathElem) bool {
	if accepted.Name != candidate.Name {
		return false
	}

	for k, v := range candidate.Key {
		acceptedV, ok := accepted.Key[k]
		if !ok || acceptedV == "*" { // Candidate key matches against wildcard accepted.
			continue
		}
		if v == "*" { // Candidate is wildcard, matches against
			continue
		}
		if v != acceptedV {
			return false
		}
	}

	return true
}

// elemToString returns a formatted string representation of a single path elem.
// wildcard keys are pruned from the resulting string
func elemToString(name string, kv map[string]string) (string, error) {
	if name == "" {
		return "", errors.New("empty name for PathElem")
	}
	if len(kv) == 0 {
		return name, nil
	}

	var keys []string
	for k, v := range kv {
		if k == "" {
			return "", fmt.Errorf("empty key name (value: %s) in element %s", v, name)
		}
		if v != "*" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := strings.ReplaceAll(kv[k], `=`, `\=`)
		v = strings.ReplaceAll(v, `]`, `\]`)
		name = fmt.Sprintf("%s[%s=%s]", name, k, v)
	}

	return name, nil
}

// setKeys returns a copy of the input map containing only the keys that not wildcards.
func setKeys(in map[string]string) map[string]string {
	m := make(map[string]string)
	for k, v := range in {
		if v != "*" {
			m[k] = v
		}
	}
	return m
}
