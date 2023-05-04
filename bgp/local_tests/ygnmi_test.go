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

import (
	"context"
	"testing"

	"github.com/openconfig/ygnmi/ygnmi"
)

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
