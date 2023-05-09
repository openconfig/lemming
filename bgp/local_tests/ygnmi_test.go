// Copyright 2023 Google LLC
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

package local_test

// TODO: Consider introducing an Ondatra binding for non-dataplane tests of lemming.

import (
	"context"
	"testing"

	"github.com/openconfig/ygnmi/ygnmi"
)

// Get fetches the value of a SingletonQuery with a ONCE subscription,
// failing the test fatally if no value is present at the path.
// Use Lookup to also get metadata or tolerate no value present.
func Get[T any](t testing.TB, c *ygnmi.Client, q ygnmi.SingletonQuery[T]) T {
	t.Helper()
	v, err := ygnmi.Get(context.Background(), c, q)
	if err != nil {
		t.Fatalf("Get(t) on %s at %v: %v", c, q, err)
	}
	return v
}

// GetConfig fetches the value of a SingletonQuery with a ONCE subscription,
// failing the test fatally if no value is present at the path.
// Use Lookup to also get metadata or tolerate no value present.
// Note: This is a workaround for Go's type inference not working for this use case and may be removed in a subsequent release.
// Note: This is equivalent to calling Get with a ConfigQuery and providing a fully-qualified type parameter.
func GetConfig[T any](t testing.TB, c *ygnmi.Client, q ygnmi.ConfigQuery[T]) T {
	t.Helper()
	return Get[T](t, c, q)
}

// GetAll fetches the value of a WildcardQuery with a ONCE subscription skipping any non-present paths.
// It fails the test fatally if no value is present at the path
// Use LookupAll to also get metadata or tolerate no values present.
func GetAll[T any](t testing.TB, c *ygnmi.Client, q ygnmi.WildcardQuery[T]) []T {
	t.Helper()
	v, err := ygnmi.GetAll(context.Background(), c, q)
	if err != nil {
		t.Fatalf("GetAll(t) on %s at %v: %v", c, q, err)
	}
	return v
}

// Update updates the configuration at the given query path with the val.
func Update[T any](t testing.TB, c *ygnmi.Client, q ygnmi.ConfigQuery[T], val T) *ygnmi.Result {
	t.Helper()
	res, err := ygnmi.Update(context.Background(), c, q, val)
	if err != nil {
		t.Fatalf("Update(t) on %v at %v: %v", c, q, err)
	}
	return res
}

// Replace replaces the configuration at the given query path with the val.
func Replace[T any](t testing.TB, c *ygnmi.Client, q ygnmi.ConfigQuery[T], val T) *ygnmi.Result {
	t.Helper()
	res, err := ygnmi.Replace(context.Background(), c, q, val)
	if err != nil {
		t.Fatalf("Replace(t) on %v at %v: %v", c, q, err)
	}
	return res
}

// Delete deletes the configuration at the given query path.
func Delete[T any](t testing.TB, c *ygnmi.Client, q ygnmi.ConfigQuery[T]) *ygnmi.Result {
	t.Helper()
	res, err := ygnmi.Delete(context.Background(), c, q)
	if err != nil {
		t.Fatalf("Delete(t) on %v at %v: %v", c, q, err)
	}
	return res
}

// Await observes values at Query with a STREAM subscription,
// blocking until a value that is deep equal to the specified val is received
// or the timeout is reached. To wait for a generic predicate, or to make a
// non-blocking call, use the Watch method instead.
func Await[T any](t testing.TB, c *ygnmi.Client, q ygnmi.SingletonQuery[T], val T) *ygnmi.Value[T] {
	t.Helper()
	v, err := ygnmi.Await(context.Background(), c, q, val)
	if err != nil {
		t.Fatalf("Await(t) on %v at %v: %v", c, q, err)
	}
	return v
}
