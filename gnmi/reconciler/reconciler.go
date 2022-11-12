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

// Package reconciler contains a common interface for gNMI reconciler.
package reconciler

import (
	"context"
	"fmt"

	"github.com/openconfig/gnmi/errlist"
	"github.com/openconfig/lemming/gnmi/oc"
	"github.com/openconfig/ygnmi/ygnmi"

	log "github.com/golang/glog"
	gpb "github.com/openconfig/gnmi/proto/gnmi"
)

// Reconciler is a common interface for gNMI reconciler.
// Reconcilers are responsible for watching config and state and performing any necessary reconciliation.
type Reconciler interface {
	// ID the id for the reconciler.
	ID() string
	// Start start the reconciliation loop.
	// The client and target are connected to the local gNMI cache.
	// An error returned during Start will cause lemming to exit, so Start should return non-retriable errors.
	Start(ctx context.Context, client gpb.GNMIClient, target string) error
	// Stop stops the reconciliation loop.
	Stop(context.Context) error
	// Validate is called after a SetRequest is checked for schema compliance, but before data is written to the cache.
	// Reconcilers can validate the intended config for semantic correctness (reject config that can never be reconciled).
	// The Validate func is only called if the SetRequest contains paths which match the ValidationPaths.
	Validate(intendedConfig *oc.Root) error
	// ValidationPaths returns the set of path prefixes that the reconciler can validate.
	ValidationPaths() []ygnmi.PathStruct
}

// Builder simplifies the creation of reconcilers and reduces some of the required boilerplate.
type Builder struct {
	br *BuiltReconciler
}

// NewBuilder creates a new reconciler builder.
func NewBuilder(id string) *Builder {
	return &Builder{
		br: &BuiltReconciler{
			id: id,
		},
	}
}

// Build returns the reconciler as configuration and resets the builder.
func (b *Builder) Build() *BuiltReconciler {
	reconciler := b.br
	b.br = &BuiltReconciler{}
	return reconciler
}

// WithStart appends a new start func to the reconciler.
func (b *Builder) WithStart(startFn func(context.Context, *ygnmi.Client) error) *Builder {
	if b.br == nil {
		b.br = &BuiltReconciler{}
	}
	b.br.startFns = append(b.br.startFns, func(ctx context.Context, client gpb.GNMIClient, target string) error {
		c, err := ygnmi.NewClient(client, ygnmi.WithTarget(target))
		if err != nil {
			return fmt.Errorf("failed to build client: %v", err)
		}
		return startFn(ctx, c)
	})
	return b
}

// WithStop appends a new stop func to the reconciler.
func (b *Builder) WithStop(stopFn func(context.Context) error) *Builder {
	if b.br == nil {
		b.br = &BuiltReconciler{}
	}
	b.br.stopFns = append(b.br.stopFns, stopFn)
	return b
}

// WithValidator appends a validator and validations paths to the reconciler.
// The Validate func is only called if the SetRequest contains paths which match the paths.
func (b *Builder) WithValidator(paths []ygnmi.PathStruct, validator func(*oc.Root) error) *Builder {
	if b.br == nil {
		b.br = &BuiltReconciler{}
	}
	b.br.validateFns = append(b.br.validateFns, validator)
	b.br.validationPaths = append(b.br.validationPaths, paths...)
	return b
}

// TypedBuilder is similar to builder except with a type parameter for use with ygnmi Queries.
type TypedBuilder[T any] struct {
	Builder
}

// NewTypedBuilder creates a new reconciler builder.
func NewTypedBuilder[T any](id string) *TypedBuilder[T] {
	return &TypedBuilder[T]{
		Builder: Builder{
			br: &BuiltReconciler{
				id: id,
			},
		},
	}
}

// WithWatch adds starting a watch to the reconciler's start funcs.
func (tb *TypedBuilder[T]) WithWatch(query ygnmi.SingletonQuery[T], predicate func(context.Context, *ygnmi.Client, *ygnmi.Value[T]) error) *TypedBuilder[T] {
	tb.WithStart(func(ctx context.Context, c *ygnmi.Client) error {
		watchCtx, cancel := context.WithCancel(ctx)
		tb.br.stopFns = append(tb.br.stopFns, func(context.Context) error { cancel(); return nil })
		w := ygnmi.Watch(watchCtx, c, query, func(v *ygnmi.Value[T]) error {
			return predicate(watchCtx, c, v)
		})

		// TODO: handle errors here.
		go func() {
			if _, err := w.Await(); err != nil {
				log.Errorf("Reconciler %q watch err: %v", tb.br.id, err)
			}
		}()
		return nil
	})
	return tb
}

// BuiltReconciler is an implementation of the reconciler interface returned by builders.
type BuiltReconciler struct {
	id              string
	startFns        []func(context.Context, gpb.GNMIClient, string) error
	stopFns         []func(context.Context) error
	validateFns     []func(*oc.Root) error
	validationPaths []ygnmi.PathStruct
}

func (bt *BuiltReconciler) ID() string {
	return bt.id
}

func (bt *BuiltReconciler) Start(ctx context.Context, client gpb.GNMIClient, target string) error {
	var l errlist.List
	for _, startFn := range bt.startFns {
		l.Add(startFn(ctx, client, target))
	}
	if err := l.Err(); err != nil {
		return fmt.Errorf("reconciler %q start errs: %v", bt.id, l.Err())
	}
	return nil
}

func (bt *BuiltReconciler) Stop(ctx context.Context) error {
	var l errlist.List
	for _, stopFn := range bt.stopFns {
		l.Add(stopFn(ctx))
	}
	if err := l.Err(); err != nil {
		return fmt.Errorf("reconciler %q stop errs: %v", bt.id, l.Err())
	}
	return nil
}

func (bt *BuiltReconciler) Validate(intendedConfig *oc.Root) error {
	var l errlist.List
	for _, validate := range bt.validateFns {
		l.Add(validate(intendedConfig))
	}
	if err := l.Err(); err != nil {
		return fmt.Errorf("reconciler %q validation errs: %v", bt.id, l.Err())
	}
	return nil
}

func (bt *BuiltReconciler) ValidationPaths() []ygnmi.PathStruct {
	return bt.validationPaths
}
