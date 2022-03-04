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

package fakedevice

import (
	"context"
	"fmt"
	"time"

	log "github.com/golang/glog"
	"github.com/openconfig/gnmi/coalesce"
	"github.com/openconfig/gnmi/ctree"
	"github.com/openconfig/gnmi/value"
	"github.com/openconfig/lemming/gnmi/gnmit"
	"github.com/openconfig/ondatra/telemetry"
	"github.com/openconfig/ondatra/telemetry/device"
	"github.com/openconfig/ygot/ygot"

	gpb "github.com/openconfig/gnmi/proto/gnmi"
)

// bootTimeTask is a task that updates the boot-time leaf with the current
// time. It does not spawn any long-running threads.
func bootTimeTask(_ gnmit.Queue, update gnmit.UpdateFn, target string, remove func()) error {
	defer remove()
	p0, _, errs := ygot.ResolvePath(device.DeviceRoot("").System().BootTime())
	if errs != nil {
		return fmt.Errorf("bootTimeTask failed to initialize due to error: %v", errs)
	}

	now, err := value.FromScalar(time.Now().UnixNano())
	if err != nil {
		return fmt.Errorf("bootTimeTask: %v", err)
	}
	log.V(2).Infof("bootTimeTask: %v, %v", p0, now)
	if err := update(&gpb.Notification{
		Timestamp: time.Now().UnixNano(),
		Prefix: &gpb.Path{
			Origin: "openconfig",
			Target: target,
		},
		Update: []*gpb.Update{{
			Path: p0,
			Val:  now,
		}},
	}); err != nil {
		return err
	}

	return nil
}

// currentDateTimeTask updates the current-datetime leaf with the current time,
// and spawns a thread that wakes up every second to update the leaf.
func currentDateTimeTask(_ gnmit.Queue, update gnmit.UpdateFn, target string, remove func()) error {
	p0, _, err := ygot.ResolvePath(device.DeviceRoot("").System().CurrentDatetime())
	if err != nil {
		return fmt.Errorf("currentDateTimeTask failed to initialize due to error: %v", err)
	}

	tick := time.Tick(time.Second)
	if tick == nil {
		return fmt.Errorf("currentDateTimeTask: tick is nil")
	}

	periodic := func() error {
		currentDatetime, err := value.FromScalar(time.Now().Format(time.RFC3339))
		if err != nil {
			return fmt.Errorf("currentDateTimeTask: %v", err)
		}
		log.V(2).Infof("currentDateTimeTask: %v, %v", p0, currentDatetime)
		if err := update(&gpb.Notification{
			Timestamp: time.Now().UnixNano(),
			Prefix: &gpb.Path{
				Origin: "openconfig",
				Target: target,
			},
			Update: []*gpb.Update{{
				Path: p0,
				Val:  currentDatetime,
			}},
		}); err != nil {
			return err
		}
		return nil
	}

	if err := periodic(); err != nil {
		return err
	}

	go func() {
		defer remove()
		for range tick {
			if err := periodic(); err != nil {
				log.Errorf("currentDateTimeTask error: %v", err)
				return
			}
		}
	}()

	return nil
}

// syslogTask is a meaningless test task that monitors updates to the
// current-datetime leaf and writes updates to the syslog message leaf whenever
// the current-datetime leaf is updated.
func syslogTask(q gnmit.Queue, update gnmit.UpdateFn, target string, remove func()) error {
	p0, _, err := ygot.ResolvePath(device.DeviceRoot("").System().Messages().Message().Msg())
	if err != nil {
		log.Errorf("syslogTask failed to initialize due to error: %v", err)
	}

	go func() {
		defer remove()
		for {
			item, _, err := q.Next(context.Background())
			if coalesce.IsClosedQueue(err) {
				return
			}
			n, ok := item.(*ctree.Leaf)
			if !ok || n == nil {
				log.Errorf("syslogTask invalid cache node: %#v", item)
				return
			}
			v := n.Value()
			no, ok := v.(*gpb.Notification)
			if !ok || no == nil {
				log.Errorf("syslogTask invalid cache node, expected non-nil *gpb.Notification type, got: %#v", v)
				return
			}
			for _, u := range no.Update {
				sv, err := value.ToScalar(u.Val)
				if err != nil {
					log.Errorf("syslogTask: %v", err)
					return
				}
				strv, ok := sv.(string)
				if !ok {
					log.Errorf("syslogTask: cannot convert to string, got (%T, %v)", sv, sv)
					return
				}
				syslog, err := value.FromScalar("current date-time updated to " + strv)
				if err != nil {
					log.Errorf("syslogTask: %v", err)
					return
				}
				if err := update(&gpb.Notification{
					Timestamp: time.Now().UnixNano(),
					Prefix: &gpb.Path{
						Origin: "openconfig",
						Target: target,
					},
					Update: []*gpb.Update{{
						Path: p0,
						Val:  syslog,
					}},
				}); err != nil {
					log.Errorf("syslogTask: %v", err)
					return
				}
			}
			for _, _ = range no.Delete {
			}
		}
	}()

	return nil
}

// tasks returns the set of functions that should be called that may read
// and/or modify internal state.
//
// These tasks are invoked during the creation of each device's Subscribe
// server and may spawn long-running or future-running thread(s) that make
// modifications to a device's cache.
func tasks(target string) []gnmit.Task {
	p0, _, err := ygot.ResolvePath(device.DeviceRoot("").System().CurrentDatetime())
	if err != nil {
		panic(fmt.Sprintf("currentDateTimeTask failed to initialize due to error: %v", err))
	}

	// TODO(wenbli): We should decentralize how we add tasks by adding a
	// register function that's called by various init() functions.
	return []gnmit.Task{{
		Run: currentDateTimeTask,
		// No paths means the task should periodically wake up itself.
		Paths: []*gpb.Path{},
		Prefix: &gpb.Path{
			Origin: "openconfig",
			Target: target,
		},
	}, {
		Run: bootTimeTask,
		// No paths means the task should periodically wake up itself.
		Paths: []*gpb.Path{},
		Prefix: &gpb.Path{
			Origin: "openconfig",
			Target: target,
		},
	}, {
		Run: syslogTask,
		// No paths means the task should periodically wake up itself.
		Paths: []*gpb.Path{
			p0,
		},
		Prefix: &gpb.Path{
			Origin: "openconfig",
			Target: target,
		},
	}}
}

// NewTarget creates a new gNMI fake device object.
// This fake gNMI server simply mirrors whatever is set for its config leafs in
// its state leafs. It also has a mechanism for adding new "tasks", or go
// thread agents that can subscribe to particular values in ONDATRA's
// OpenConfig tree and write back values to it.
func NewTarget(ctx context.Context, addr, targetName string) (*gnmit.Collector, string, error) {
	schema, err := telemetry.Schema()
	if err != nil {
		return nil, "", err
	}
	c, addr, err := gnmit.NewSettable(ctx, addr, targetName, false, schema, tasks(targetName))
	if err != nil {
		return nil, "", fmt.Errorf("cannot start server, got err: %v", err)
	}
	return c, addr, nil
}
