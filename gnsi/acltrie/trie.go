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
	root        *trieNode
	memberships map[string]map[string]bool
}

type trieNode struct {
	children map[string]*trieNode
	elem     *gpb.PathElem
	ruleID   string
	users    policies
	groups   policies
	parent   *trieNode
}

// getAction returns the action associated for the given user and mode.
// It returns UNSPECIFIED, if the user does not have an action in the node.
// It returns true, if the policy was for the user, and false if it was for the group.
func (tn *trieNode) getAction(user string, mode pathzpb.Mode, memberships map[string]map[string]bool) (pathzpb.Action, bool) {
	if m, ok := tn.users[mode]; ok {
		if act, ok := m[user]; ok {
			return act, true
		}
	}

	var act pathzpb.Action
	for group, action := range tn.groups[mode] {
		if _, ok := memberships[group][user]; ok {
			if action == pathzpb.Action_ACTION_DENY { // DENY action take precedence over PERMIT.
				return action, false
			}
			act = action
		}
	}
	return act, false
}

// Policies are map of mode to user to action.
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

// Insert inserts a new policy into the trie.
func (t *Trie) Insert(r *pathzpb.AuthorizationRule) error {
	if t.root == nil {
		t.root = &trieNode{
			children: make(map[string]*trieNode),
		}
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

	for _, elem := range path.Elem {
		if elem.Name == "*" {
			return fmt.Errorf("wildcard path names are not permitted")
		}

		// Normalize path string by pruning wildcard keys.
		pathStr, err := elemToString(elem.Name, elem.Key)
		if err != nil {
			return fmt.Errorf("invalid path element: %v", err)
		}

		// If the node already exists, keep going.
		if _, ok := node.children[pathStr]; ok {
			node = node.children[pathStr]
			continue
		}

		node.children[pathStr] = &trieNode{
			elem:     elem,
			groups:   policies{},
			users:    policies{},
			parent:   node,
			children: make(map[string]*trieNode),
		}

		node = node.children[pathStr]
	}
	// Validate the path by ensuring that compared to all other paths with the same length and same name,
	// all list keys do not overlap.
	err := t.walk(func(node *trieNode, depth int) (bool, error) {
		if depth >= len(path.Elem) {
			return false, nil
		}
		if node.elem.Name != path.Elem[depth].Name {
			return false, nil
		}
		if depth == len(path.Elem)-1 && node.ruleID != "" {
			pathCmp := other
			n := node
			for i := depth; i >= 0; i-- {
				elemCmp, err := comparePathElem(path.Elem[i], n.elem)
				if err != nil {
					return false, err
				}
				if pathCmp != other && elemCmp != pathCmp {
					return false, fmt.Errorf("path is not consistently subset or superset of other rules")
				}
				if elemCmp != other {
					pathCmp = elemCmp
				}

				n = n.parent
			}
		}
		return true, nil
	})
	if err != nil {
		return fmt.Errorf("policy path conflict: %v", err)
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
	node.ruleID = r.GetId()

	return nil
}

// Probe returns the for the given path, user, and mode; if there is no match, deny is returned
func (t *Trie) Probe(path *gpb.Path, user string, mode pathzpb.Mode) (string, pathzpb.Action) {
	potentialPolicies := []*trieNode{}
	maxDepth := 0

	// Walk the matching rules, keeping only the longest ones that contain either the user or a group that the user belongs too.
	t.walk(func(node *trieNode, depth int) (bool, error) {
		if !pathElemsMatch(node.elem, path.Elem[depth]) {
			return false, nil
		}
		if node.ruleID == "" {
			return true, nil
		}
		if act, _ := node.getAction(user, mode, t.memberships); act != pathzpb.Action_ACTION_UNSPECIFIED {
			if depth > maxDepth {
				potentialPolicies = nil
				maxDepth = depth
			}
			potentialPolicies = append(potentialPolicies, node)
		}

		return true, nil
	})
	if len(potentialPolicies) == 0 {
		return "", pathzpb.Action_ACTION_DENY
	}

	// Pick the policies with the largest number of definite keys.
	var bestPaths []*trieNode
	maxSetKeys := -1
	for _, p := range potentialPolicies {
		path := &gpb.Path{
			Elem: make([]*gpb.PathElem, maxDepth),
		}
		n := p
		setKey := 0
		for i := maxDepth - 1; i >= 0; i-- {
			path.Elem[i] = n.elem
			setKey += setKeys(path.Elem[i].Key)
			n = p.parent
		}
		if setKey >= maxSetKeys {
			maxSetKeys = setKey
			bestPaths = append(bestPaths, p)
		}
	}

	var finalAction pathzpb.Action
	var finalID string
	// Prefer user over groups and DENY over permit.
	for _, n := range bestPaths {
		act, isUser := n.getAction(user, mode, t.memberships)
		switch {
		case isUser && act == pathzpb.Action_ACTION_DENY: // Prefer user and deny, so return immediately.
			return n.ruleID, act
		case isUser: // Prefer a user over group.
			finalAction = act
			finalID = n.ruleID
		case finalAction != pathzpb.Action_ACTION_DENY: // Prefer deny over allow.
			finalAction = act
			finalID = n.ruleID
		}
	}

	return finalID, finalAction
}

// walk explores the trie in breadth first order and invokes walkFn on every node.
// To continue exploring children of the node, the walkFn must return true.
func (t *Trie) walk(walkFn func(node *trieNode, depth int) (bool, error)) error {
	type traversalNode struct {
		node  *trieNode
		depth int
	}

	queue := []*traversalNode{{node: t.root}}
	for len(queue) > 0 {
		front := queue[0]
		queue = queue[1:]
		for _, c := range front.node.children {
			cont, err := walkFn(c, front.depth)
			if err != nil {
				return err
			}
			if cont {
				queue = append(queue, &traversalNode{node: c, depth: front.depth + 1})
			}
		}
	}
	return nil
}

type compareResult int

const (
	other compareResult = iota
	subset
	superset
)

// comparePathElem compare two path elements a, b and returns:
// subset: if every definite key in a is wildcard in b.
// superset: if every wildcard key in b is non-wildcard in b.
// other: all keys are the same or all keys are different.
// error: not two keys are both subset and superset.
// TODO: Change to more broadly useful compare func that outputs subset, superset, equal, disjoint and upstream to ygot.
func comparePathElem(a, b *gpb.PathElem) (compareResult, error) {
	setRelation := other
	for k, aVal := range a.Key {
		bVal, ok := b.Key[k]
		switch {
		case aVal == bVal:
			continue
		case aVal == "*" && !ok: // b is implicitly wildcarded.
			continue
		case aVal == "*":
			if setRelation == subset {
				return other, fmt.Errorf("path %v is not consistently a superset of %v", a, b)
			}
			setRelation = superset
		case bVal == "*" || !ok:
			if setRelation == superset {
				return other, fmt.Errorf("path %v is not consistently a subset of %v", a, b)
			}
			setRelation = subset
		}
	}
	for k, bVal := range b.Key {
		_, ok := a.Key[k]
		if !ok && bVal != "*" { // If a contains an implicit wildcard.
			if setRelation == subset {
				return other, fmt.Errorf("path %v is not consistently a superset of %v", a, b)
			}
			setRelation = superset
		}
	}

	return setRelation, nil
}

// elemToString returns a formatted string representation of a single path elem.
// wildcard keys are pruned from the resulting string.
// TODO: upstream to ygot.
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

// pathElemsMatch returns true if the a PathElem matches the b PathElem.
// A match is when both the name match and all keys match (where * or unset matches with anything).
func pathElemsMatch(a, b *gpb.PathElem) bool {
	if a.Name != b.Name {
		return false
	}

	for k, bVal := range b.Key {
		aVal, ok := a.Key[k]
		if !ok || aVal == "*" { // a key matches against wildcard.
			continue
		}
		if bVal == "*" { // b is wildcard, matches against anything.
			continue
		}
		if bVal != aVal {
			return false
		}
	}

	return true
}

// setKeys returns the number on non-wildcard keys in the map.
func setKeys(in map[string]string) int {
	c := 0
	for _, v := range in {
		if v != "*" {
			c++
		}
	}
	return c
}
