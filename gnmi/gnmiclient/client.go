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

// Package gnmiclient contains a funcs to create gNMI for the local cache.
package gnmiclient

import (
	"context"

	"github.com/openconfig/ygnmi/ygnmi"
)

// Update updates the configuration at the given query path with the val.
func Update[T any](ctx context.Context, c *ygnmi.Client, q ygnmi.SingletonQuery[T], val T) (*ygnmi.Result, error) {
	return ygnmi.Update[T](ctx, c, &singletonAsConfig[T]{SingletonQuery: q}, val)
}

// Replace replaces the configuration at the given query path with the val.
func Replace[T any](ctx context.Context, c *ygnmi.Client, q ygnmi.SingletonQuery[T], val T) (*ygnmi.Result, error) {
	return ygnmi.Replace[T](ctx, c, &singletonAsConfig[T]{SingletonQuery: q}, val)
}

// Delete deletes the configuration at the given query path.
func Delete[T any](ctx context.Context, c *ygnmi.Client, q ygnmi.SingletonQuery[T]) (*ygnmi.Result, error) {
	return ygnmi.Delete[T](ctx, c, &singletonAsConfig[T]{SingletonQuery: q})
}

// BatchUpdate stores an update operation in the SetBatch.
func BatchUpdate[T any](sb *ygnmi.SetBatch, q ygnmi.SingletonQuery[T], val T) {
	ygnmi.BatchUpdate[T](sb, &singletonAsConfig[T]{SingletonQuery: q}, val)
}

// BatchReplace stores an replace operation in the SetBatch.
func BatchReplace[T any](sb *ygnmi.SetBatch, q ygnmi.SingletonQuery[T], val T) {
	ygnmi.BatchReplace[T](sb, &singletonAsConfig[T]{SingletonQuery: q}, val)
}

// BatchDelete stores an delete operation in the SetBatch.
func BatchDelete[T any](sb *ygnmi.SetBatch, q ygnmi.SingletonQuery[T]) {
	ygnmi.BatchDelete[T](sb, &singletonAsConfig[T]{SingletonQuery: q})
}

// singletonAsConfig turns a SingletonQuery into ConfigQuery.
type singletonAsConfig[T any] struct {
	ygnmi.SingletonQuery[T]
}

func (*singletonAsConfig[T]) IsConfig() {}
