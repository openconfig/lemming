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

	"github.com/openconfig/ygot/util"

	gpb "github.com/openconfig/gnmi/proto/gnmi"
	pathzpb "github.com/openconfig/gnsi/pathz"
)

// Trie is the root of the ACL trie.
type Trie struct {
	root *trieNode
	// memberships is a map of group name to a set of users.
	memberships map[string]map[string]bool
}

type trieNode struct {
	children     map[string]*trieNode
	elem         *gpb.PathElem
	hasPolicy    bool
	users        policies
	groups       policies
	totalSetKeys int
	parent       *trieNode
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

// FromPolicy creates a new trie from a pathzpb.AuthorizationPolicy.
func FromPolicy(p *pathzpb.AuthorizationPolicy) (*Trie, error) {
	t := &Trie{
		memberships: map[string]map[string]bool{},
	}
	for _, group := range p.GetGroups() {
		if _, ok := t.memberships[group.GetName()]; !ok {
			t.memberships[group.GetName()] = make(map[string]bool)
		}
		for _, user := range group.GetUsers() {
			t.memberships[group.GetName()][user.GetName()] = true
		}
	}

	for _, rule := range p.GetRules() {
		if err := t.Insert(rule); err != nil {
			return nil, err
		}
	}
	return t, nil
}

// Insert inserts a new rule into the trie.
func (t *Trie) Insert(r *pathzpb.AuthorizationRule) error {
	if t.root == nil {
		t.root = &trieNode{
			children: make(map[string]*trieNode),
			users:    make(policies),
			groups:   make(policies),
		}
	}
	if r.Action == pathzpb.Action_ACTION_UNSPECIFIED {
		return fmt.Errorf("action unspecified")
	}
	if r.Action != pathzpb.Action_ACTION_DENY && r.Action != pathzpb.Action_ACTION_PERMIT {
		return fmt.Errorf("unknown action type")
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
		pathStr, keys, err := elemToString(elem.Name, elem.Key)
		if err != nil {
			return fmt.Errorf("invalid path element: %v", err)
		}

		// If the node already exists, keep going.
		if _, ok := node.children[pathStr]; ok {
			node = node.children[pathStr]
			continue
		}

		node.children[pathStr] = &trieNode{
			elem:         elem,
			groups:       policies{},
			users:        policies{},
			parent:       node,
			totalSetKeys: node.totalSetKeys + keys,
			children:     make(map[string]*trieNode),
		}

		node = node.children[pathStr]
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

// Probe returns the action for the given path, user, and mode; if there is no match, deny is returned
func (t *Trie) Probe(path *gpb.Path, user string, mode pathzpb.Mode) pathzpb.Action {
	matchingPolicies := []*trieNode{}
	longestPolicyLen := 0

	// Walk the matching rules, keeping only the longest ones that contain either the user or a group that the user belongs too.
	t.walk(func(node *trieNode, walkPath *gpb.Path) (bool, error) {
		if res := util.ComparePaths(path, walkPath); res != util.Equal && res != util.Subset {
			return false, nil
		}
		if !node.hasPolicy {
			return true, nil
		}
		if act, _ := node.getAction(user, mode, t.memberships); act != pathzpb.Action_ACTION_UNSPECIFIED {
			if len(walkPath.Elem) > longestPolicyLen {
				matchingPolicies = nil
				longestPolicyLen = len(walkPath.Elem)
			}
			matchingPolicies = append(matchingPolicies, node)
		}

		return true, nil
	})

	if len(matchingPolicies) == 0 {
		return pathzpb.Action_ACTION_DENY
	}

	// Pick the policies with the largest number of definite keys.
	var mostSpecificPolicies []*trieNode
	maxSetKeys := -1
	for _, p := range matchingPolicies {
		if p.totalSetKeys >= maxSetKeys {
			maxSetKeys = p.totalSetKeys
			mostSpecificPolicies = append(mostSpecificPolicies, p)
		}
	}

	var finalAction pathzpb.Action
	// Prefer user over groups and DENY over permit.
	for _, n := range mostSpecificPolicies {
		act, isUser := n.getAction(user, mode, t.memberships)
		switch {
		case isUser && act == pathzpb.Action_ACTION_DENY: // Prefer user and deny, so return immediately.
			return act
		case isUser: // Prefer a user over group.
			finalAction = act
		case finalAction != pathzpb.Action_ACTION_DENY: // Prefer deny over allow.
			finalAction = act
		}
	}

	return finalAction
}

// walk explores the trie in depth first order and invokes walkFn on every node.
// To continue exploring children of the node, the walkFn must return true.
// Note: the path object is modified between calls.
func (t *Trie) walk(walkFn func(*trieNode, *gpb.Path) (bool, error)) error {
	type traversalNode struct {
		node  *trieNode
		depth int
	}

	path := &gpb.Path{}
	stack := []*traversalNode{{node: t.root, depth: 0}}

	for len(stack) > 0 {
		last := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if last.depth == 0 {
			path.Elem = []*gpb.PathElem{}
		} else {
			path.Elem = append(path.Elem[:last.depth-1], last.node.elem)
		}

		cont, err := walkFn(last.node, path)
		if err != nil {
			return err
		}
		if !cont {
			continue
		}

		for _, c := range last.node.children {
			stack = append(stack, &traversalNode{node: c, depth: last.depth + 1})
		}
	}
	return nil
}

// elemToString returns a formatted string representation of a single path elem.
// wildcard keys are pruned from the resulting string and the number of non-wildcard keys are returned.
// TODO: upstream to ygot.
func elemToString(name string, kv map[string]string) (string, int, error) {
	if name == "" {
		return "", 0, errors.New("empty name for PathElem")
	}
	if len(kv) == 0 {
		return name, 0, nil
	}

	var keys []string
	for k, v := range kv {
		if k == "" {
			return "", 0, fmt.Errorf("empty key name (value: %s) in element %s", v, name)
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

	return name, len(keys), nil
}
