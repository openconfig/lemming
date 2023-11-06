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
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/openconfig/ygnmi/ygnmi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	awaitTimeLimit = 10 * time.Second
)

// Get fetches the value of a SingletonQuery with a ONCE subscription,
// failing the test fatally if no value is present at the path.
// Use Lookup to also get metadata or tolerate no value present.
func Get[T any](t testing.TB, dut *Device, q ygnmi.SingletonQuery[T]) T {
	t.Helper()
	c := dut.yc
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
func GetConfig[T any](t testing.TB, dut *Device, q ygnmi.ConfigQuery[T]) T {
	t.Helper()
	return Get[T](t, dut, q)
}

// GetAll fetches the value of a WildcardQuery with a ONCE subscription skipping any non-present paths.
// It fails the test fatally if no value is present at the path
// Use LookupAll to also get metadata or tolerate no values present.
func GetAll[T any](t testing.TB, dut *Device, q ygnmi.WildcardQuery[T]) []T {
	t.Helper()
	c := dut.yc
	v, err := ygnmi.GetAll(context.Background(), c, q)
	if err != nil {
		t.Fatalf("GetAll(t) on %s at %v: %v", c, q, err)
	}
	return v
}

// Update updates the configuration at the given query path with the val.
func Update[T any](t testing.TB, dut *Device, q ygnmi.ConfigQuery[T], val T) *ygnmi.Result {
	t.Helper()
	c := dut.yc
	res, err := ygnmi.Update(context.Background(), c, q, val)
	if err != nil {
		t.Fatalf("Update(t) on %v at %v: %v", c, q, err)
	}
	return res
}

// Replace replaces the configuration at the given query path with the val.
func Replace[T any](t testing.TB, dut *Device, q ygnmi.ConfigQuery[T], val T) *ygnmi.Result {
	t.Helper()
	c := dut.yc
	res, err := ygnmi.Replace(context.Background(), c, q, val)
	if err != nil {
		t.Fatalf("Replace(t) on %v at %v: %v", c, q, err)
	}
	return res
}

// ReplaceExpectFail replaces the configuration at the given query path with
// the val, expecting a failure.
func ReplaceExpectFail[T any](t testing.TB, dut *Device, q ygnmi.ConfigQuery[T], val T) {
	t.Helper()
	c := dut.yc
	_, err := ygnmi.Replace(context.Background(), c, q, val)
	if err == nil {
		t.Fatalf("Replace(t) on %v at %v: did not fail", c, q)
	}
}

// Delete deletes the configuration at the given query path.
func Delete[T any](t testing.TB, dut *Device, q ygnmi.ConfigQuery[T]) *ygnmi.Result {
	t.Helper()
	c := dut.yc
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
func Await[T any](t testing.TB, dut *Device, q ygnmi.SingletonQuery[T], val T) *ygnmi.Value[T] {
	t.Helper()
	c := dut.yc
	ctx, cancel := context.WithTimeout(context.Background(), awaitTimeLimit)
	defer cancel()
	v, err := ygnmi.Await(ctx, c, q, val)
	if err != nil {
		t.Fatalf("Await(t) on %v at %v: %v", c, q, err)
	}
	return v
}

// AwaitWithErr is the same as Await except an error is returned.
//
// Its purpose is to add a better error message.
func AwaitWithErr[T any](dut *Device, q ygnmi.SingletonQuery[T], val T) (*ygnmi.Value[T], error) {
	ctx, cancel := context.WithTimeout(context.Background(), awaitTimeLimit)
	defer cancel()
	c := dut.yc
	v, err := ygnmi.Await(ctx, c, q, val)
	if err != nil {
		return v, fmt.Errorf("Await(t) on %v at %v: %v", c, q, err)
	}
	return v, nil
}

type watchAwaiter[T any] interface {
	Await() (*ygnmi.Value[T], error)
}

// Watcher represents an ongoing watch of telemetry values.
type Watcher[T any] struct {
	watcher  watchAwaiter[T]
	cancelFn func()
	c        *ygnmi.Client
	query    ygnmi.AnyQuery[T]
}

func isContextErr(err error) bool {
	// https://pkg.go.dev/google.golang.org/grpc@v1.48.0/internal/status#Error
	var st interface {
		GRPCStatus() *status.Status
	}
	ok := errors.As(err, &st)
	if !ok {
		return errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled)
	}
	return st.GRPCStatus().Code() == codes.DeadlineExceeded || st.GRPCStatus().Code() == codes.Canceled
}

// Await waits for the watch to finish and returns the last received value
// and a boolean indicating whether the predicate evaluated to true.
// When Await returns the watcher is closed, and Await may not be called again.
func (w *Watcher[T]) Await(t testing.TB) (*ygnmi.Value[T], bool) {
	t.Helper()
	v, err := w.watcher.Await()
	if err != nil {
		if isContextErr(err) {
			return v, false
		}
		t.Fatalf("Await(t) on %s at %v: %v", w.c, w.query, err)
	}
	return v, true
}

// Cancel stops the watch immediately.
func (w *Watcher[T]) Cancel() {
	w.cancelFn()
}

// Watch starts an asynchronous STREAM subscription, evaluating each observed value with the
// specified predicate. The subscription completes when either the predicate is true
// or the timeout is reached. Calling Await on the returned Watcher waits for the subscription
// to complete. It returns the last observed value and a boolean that indicates whether
// that value satisfies the predicate.
func Watch[T any](t testing.TB, dut *Device, q ygnmi.SingletonQuery[T], timeout time.Duration, pred func(*ygnmi.Value[T]) bool) *Watcher[T] {
	t.Helper()
	c := dut.yc
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	w := ygnmi.Watch(ctx, c, q, func(v *ygnmi.Value[T]) error {
		if ok := pred(v); ok {
			return nil
		}
		return ygnmi.Continue
	})
	return &Watcher[T]{
		watcher:  w,
		cancelFn: cancel,
		c:        c,
		query:    q,
	}
}
