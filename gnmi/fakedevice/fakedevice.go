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
	"github.com/openconfig/goyang/pkg/yang"
	"github.com/openconfig/lemming/gnmi/gnmiclient"
	"github.com/openconfig/lemming/gnmi/gnmit"
	"github.com/openconfig/lemming/gnmi/oc"
	"github.com/openconfig/lemming/gnmi/oc/ocpath"
	"github.com/openconfig/ygnmi/ygnmi"
	"github.com/openconfig/ygot/util"
	"github.com/openconfig/ygot/ygot"
	"github.com/openconfig/ygot/ygot/pathtranslate"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"

	gpb "github.com/openconfig/gnmi/proto/gnmi"
)

var PathTranslator *pathtranslate.PathTranslator

func init() {
	var schemas []*yang.Entry
	for _, s := range oc.SchemaTree {
		schemas = append(schemas, s)
	}
	var err error
	if PathTranslator, err = pathtranslate.NewPathTranslator(schemas); err != nil {
		panic(err)
	}
}

// bootTimeTask is a task that updates the boot-time leaf with the current
// time. It does not spawn any long-running threads.
func bootTimeTask(_ func() *oc.Root, _ gnmit.Queue, updateFn gnmit.UpdateFn, target string, remove func()) error {
	defer remove()
	pathBootTime, _, errs := ygnmi.ResolvePath(ocpath.Root().System().BootTime().State().PathStruct())
	if errs != nil {
		return fmt.Errorf("bootTimeTask failed to initialize due to error: %v", errs)
	}

	now, err := value.FromScalar(time.Now().UnixNano())
	if err != nil {
		return fmt.Errorf("bootTimeTask: %v", err)
	}
	log.V(2).Infof("bootTimeTask: %v, %v", pathBootTime, now)
	if err := updateFn(&gpb.Notification{
		Timestamp: time.Now().UnixNano(),
		Prefix: &gpb.Path{
			Origin: "openconfig",
			Target: target,
		},
		Update: []*gpb.Update{{
			Path: pathBootTime,
			Val:  now,
		}},
	}); err != nil {
		return err
	}

	return nil
}

// currentDateTimeTask updates the current-datetime leaf with the current time,
// and spawns a thread that wakes up every second to update the leaf.
func currentDateTimeTask(_ func() *oc.Root, _ gnmit.Queue, updateFn gnmit.UpdateFn, target string, remove func()) error {
	pathDatetime, _, err := ygnmi.ResolvePath(ocpath.Root().System().CurrentDatetime().State().PathStruct())
	if err != nil {
		return fmt.Errorf("currentDateTimeTask failed to initialize due to error: %v", err)
	}

	tick := time.NewTicker(time.Second)

	periodic := func() error {
		currentDatetime, err := value.FromScalar(time.Now().Format(time.RFC3339))
		if err != nil {
			return fmt.Errorf("currentDateTimeTask: %v", err)
		}
		log.V(2).Infof("currentDateTimeTask: %v, %v", pathDatetime, currentDatetime)
		if err := updateFn(&gpb.Notification{
			Timestamp: time.Now().UnixNano(),
			Prefix: &gpb.Path{
				Origin: "openconfig",
				Target: target,
			},
			Update: []*gpb.Update{{
				Path: pathDatetime,
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
		for range tick.C {
			if err := periodic(); err != nil {
				log.Errorf("currentDateTimeTask error: %v", err)
				return
			}
		}
	}()

	return nil
}

// matchingPath returns true iff the path matches the given matcher path in
// length and in values; wildcards are allowed in the matcher path.
func matchingPath(path, matcher *gpb.Path) bool {
	return len(path.Elem) == len(matcher.Elem) && util.PathMatchesQuery(path, matcher)
}

// toStatePath converts the given config path to a state path by replacing the
// last instance (if any) of "config" in the path to "state".
// OpenConfig specifies that any leaf other than list keys must reside in a
// config/state container, and that there shall only be one such container in
// the path.
func toStatePath(configPath *gpb.Path) *gpb.Path {
	path := proto.Clone(configPath).(*gpb.Path)
	for i := len(path.Elem) - 1; i >= 0; i-- {
		if path.Elem[i].Name == "config" {
			path.Elem[i].Name = "state"
			break
		}
	}
	return path
}

// systemBaseTask handles most of the logic for the base systems feature profile.
func systemBaseTask(_ func() *oc.Root, queue gnmit.Queue, updateFn gnmit.UpdateFn, target string, remove func()) error {
	hostnamePath, _, err := ygnmi.ResolvePath(ocpath.Root().System().Hostname().Config().PathStruct())
	if err != nil {
		log.Errorf("systemBaseTask failed to initialize due to error: %v", err)
	}
	domainNamePath, _, err := ygnmi.ResolvePath(ocpath.Root().System().DomainName().Config().PathStruct())
	if err != nil {
		log.Errorf("systemBaseTask failed to initialize due to error: %v", err)
	}
	motdBannerPath, _, err := ygnmi.ResolvePath(ocpath.Root().System().MotdBanner().Config().PathStruct())
	if err != nil {
		log.Errorf("systemBaseTask failed to initialize due to error: %v", err)
	}
	loginBannerPath, _, err := ygnmi.ResolvePath(ocpath.Root().System().LoginBanner().Config().PathStruct())
	if err != nil {
		log.Errorf("systemBaseTask failed to initialize due to error: %v", err)
	}

	go func() {
		defer remove()
		for {
			item, _, err := queue.Next(context.Background())
			if coalesce.IsClosedQueue(err) {
				return
			}
			n, ok := item.(*ctree.Leaf)
			if !ok || n == nil {
				log.Errorf("systemBaseTask invalid cache node: %#v", item)
				return
			}
			v := n.Value()
			noti, ok := v.(*gpb.Notification)
			if !ok || noti == nil {
				log.Errorf("systemBaseTask invalid cache node, expected non-nil *gpb.Notification type, got: %#v", v)
				return
			}
			for _, update := range noti.Update {
				switch {
				case matchingPath(update.Path, hostnamePath), matchingPath(update.Path, domainNamePath), matchingPath(update.Path, motdBannerPath), matchingPath(update.Path, loginBannerPath):
					log.Infof("systemBaseTask got recognized path: %s", update.Path)
					statePath := toStatePath(update.Path)
					if err := updateFn(&gpb.Notification{
						Timestamp: time.Now().UnixNano(),
						Prefix: &gpb.Path{
							Origin: "openconfig",
							Target: target,
						},
						Update: []*gpb.Update{{
							Path: statePath,
							Val:  update.Val,
						}},
					}); err != nil {
						log.Errorf("systemBaseTask: %v", err)
						return
					}
				default:
					log.Errorf("systemBaseTask: update path received isn't matched by any handlers: %s", prototext.Format(update.Path))
				}
			}
			for _, path := range noti.Delete {
				// Since gNMI still sends delete paths using the deprecated Element field, we need to translate it into path-elems first.
				// We also need to strip the first element for origin.
				if len(path.Element) == 0 { // nolint:staticcheck
					log.Errorf("Unexpected: Element field for delete path is empty: %s", prototext.Format(path))
					return
				}
				elems, err := PathTranslator.PathElem(path.Element[1:]) // nolint:staticcheck
				if err != nil {
					log.Errorf("systemBaseTask: failed to translate delete path: %s", prototext.Format(path))
					return
				}
				path.Elem = elems
				switch {
				case matchingPath(path, hostnamePath), matchingPath(path, domainNamePath), matchingPath(path, motdBannerPath), matchingPath(path, loginBannerPath):
					statePath := toStatePath(path)
					if err := updateFn(&gpb.Notification{
						Timestamp: time.Now().UnixNano(),
						Prefix: &gpb.Path{
							Origin: "openconfig",
							Target: target,
						},
						Delete: []*gpb.Path{
							statePath,
						},
					}); err != nil {
						log.Errorf("systemBaseTask: %v", err)
						return
					}
				default:
					log.Errorf("systemBaseTask: delete path received isn't matched by any handlers: %s", prototext.Format(path))
				}
			}
		}
	}()

	return nil
}

// StartSystemBaseTask2 handles most of the logic for the base systems feature profile but uses ygnmi as the client.
func StartSystemBaseTask2(ctx context.Context, port int, target string, enableTLS bool) error {
	yclient, err := gnmiclient.NewYGNMIClient(port, target, enableTLS)
	if err != nil {
		return err
	}

	b := &ocpath.Batch{}
	b.AddPaths(
		ocpath.Root().System().Hostname().Config().PathStruct(),
		ocpath.Root().System().DomainName().Config().PathStruct(),
		ocpath.Root().System().MotdBanner().Config().PathStruct(),
		ocpath.Root().System().LoginBanner().Config().PathStruct(),
	)

	var hostname, domainName, motdBanner, loginBanner string

	systemWatcher := ygnmi.Watch(
		context.Background(),
		yclient,
		b.Config(),
		func(root *ygnmi.Value[*oc.Root]) error {
			rootVal, ok := root.Val()
			if !ok {
				return ygnmi.Continue
			}
			// DEBUG print out root.
			js, err := ygot.Marshal7951(rootVal, ygot.JSONIndent("  "))
			if err != nil {
				log.Infof("systembase watch: %s", string(js))
			}
			system := rootVal.GetSystem()
			if system == nil {
				return ygnmi.Continue
			}
			if system.Hostname != nil && system.GetHostname() != hostname {
				if _, err := gnmiclient.Replace(ctx, yclient, ocpath.Root().System().Hostname().State(), system.GetHostname()); err != nil {
					log.Warningf("unable to update hostname: %v", err)
				} else {
					hostname = system.GetHostname()
					log.Infof("Successfully updated hostname to %q", hostname)
				}
			}
			if system.DomainName != nil && system.GetDomainName() != domainName {
				if _, err := gnmiclient.Replace(ctx, yclient, ocpath.Root().System().DomainName().State(), system.GetDomainName()); err != nil {
					log.Warningf("unable to update domainName: %v", err)
				} else {
					domainName = system.GetDomainName()
					log.Infof("Successfully updated domainName to %q", domainName)
				}
			}
			if system.MotdBanner != nil && system.GetMotdBanner() != motdBanner {
				if _, err := gnmiclient.Replace(ctx, yclient, ocpath.Root().System().MotdBanner().State(), system.GetMotdBanner()); err != nil {
					log.Warningf("unable to update motdBanner: %v", err)
				} else {
					motdBanner = system.GetMotdBanner()
					log.Infof("Successfully updated motdBanner to %q", motdBanner)
				}
			}
			if system.LoginBanner != nil && system.GetLoginBanner() != loginBanner {
				if _, err := gnmiclient.Replace(ctx, yclient, ocpath.Root().System().LoginBanner().State(), system.GetLoginBanner()); err != nil {
					log.Warningf("unable to update loginBanner: %v", err)
				} else {
					loginBanner = system.GetLoginBanner()
					log.Infof("Successfully updated loginBanner to %q", loginBanner)
				}
			}
			return ygnmi.Continue
		},
	)

	// TODO(wenbli): Ideally, this is implemented by watching more fine-grained paths.
	// TODO(wenbli): Support interface removal.
	go func() {
		if _, err := systemWatcher.Await(); err != nil {
			log.Warningf("Sysrib interface watcher has stopped: %v", err)
		}
	}()
	return nil
}

// syslogTask is a meaningless test task that monitors updates to the
// current-datetime leaf and writes updates to the syslog message leaf whenever
// the current-datetime leaf is updated.
func syslogTask(_ func() *oc.Root, queue gnmit.Queue, updateFn gnmit.UpdateFn, target string, remove func()) error {
	pathSystemMsg, _, err := ygnmi.ResolvePath(ocpath.Root().System().Messages().Message().Msg().State().PathStruct())
	if err != nil {
		log.Errorf("syslogTask failed to initialize due to error: %v", err)
	}

	go func() {
		defer remove()
		for {
			item, _, err := queue.Next(context.Background())
			if coalesce.IsClosedQueue(err) {
				return
			}
			n, ok := item.(*ctree.Leaf)
			if !ok || n == nil {
				log.Errorf("syslogTask invalid cache node: %#v", item)
				return
			}
			v := n.Value()
			noti, ok := v.(*gpb.Notification)
			if !ok || noti == nil {
				log.Errorf("syslogTask invalid cache node, expected non-nil *gpb.Notification type, got: %#v", v)
				return
			}
			for _, u := range noti.Update {
				scalarValue, err := value.ToScalar(u.Val)
				if err != nil {
					log.Errorf("syslogTask: %v", err)
					return
				}
				strv, ok := scalarValue.(string)
				if !ok {
					log.Errorf("syslogTask: cannot convert to string, got (%T, %v)", scalarValue, scalarValue)
					return
				}
				syslog, err := value.FromScalar("current date-time updated to " + strv)
				if err != nil {
					log.Errorf("syslogTask: %v", err)
					return
				}
				if err := updateFn(&gpb.Notification{
					Timestamp: time.Now().UnixNano(),
					Prefix: &gpb.Path{
						Origin: "openconfig",
						Target: target,
					},
					Update: []*gpb.Update{{
						Path: pathSystemMsg,
						Val:  syslog,
					}},
				}); err != nil {
					log.Errorf("syslogTask: %v", err)
					return
				}
			}
			for range noti.Delete {
			}
		}
	}()

	return nil
}

// Tasks returns the set of functions that should be called that may read
// and/or modify internal state.
//
// These tasks are invoked during the creation of each device's Subscribe
// server and may spawn long-running or future-running thread(s) that make
// modifications to a device's cache.
func Tasks(target string) []gnmit.Task {
	// TODO(wenbli): We should decentralize how we add tasks by adding a
	// register function that's called by various init() functions.
	return []gnmit.Task{{
		Run: currentDateTimeTask,
		// No paths means the task should periodically wake up itself if it needs to be run at a later time.
		Paths: []ygnmi.PathStruct{},
		Prefix: &gpb.Path{
			Origin: "openconfig",
			Target: target,
		},
	}, {
		Run:   bootTimeTask,
		Paths: []ygnmi.PathStruct{},
		Prefix: &gpb.Path{
			Origin: "openconfig",
			Target: target,
		},
	}, {
		Run: syslogTask,
		Paths: []ygnmi.PathStruct{
			ocpath.Root().System().CurrentDatetime().State().PathStruct(),
		},
		Prefix: &gpb.Path{
			Origin: "openconfig",
			Target: target,
		},
		//}, {
		//	Run: systemBaseTask,
		//	Paths: []ygnmi.PathStruct{
		//		ocpath.Root().System().Hostname().Config().PathStruct(),
		//		ocpath.Root().System().DomainName().Config().PathStruct(),
		//		ocpath.Root().System().MotdBanner().Config().PathStruct(),
		//		ocpath.Root().System().LoginBanner().Config().PathStruct(),
		//	},
		//	Prefix: &gpb.Path{
		//		Origin: "openconfig",
		//		Target: target,
		//	},
	}, {
		Run: goBgpTask,
		Paths: []ygnmi.PathStruct{
			ocpath.Root().NetworkInstance("default").Protocol(oc.PolicyTypes_INSTALL_PROTOCOL_TYPE_BGP, "BGP").Bgp(),
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
	configSchema, err := oc.Schema()
	if err != nil {
		return nil, "", fmt.Errorf("cannot create ygot schema object for gNMI target: %v", err)
	}
	c, addr, err := gnmit.NewSettable(ctx, addr, targetName, false, configSchema, Tasks(targetName))
	if err != nil {
		return nil, "", fmt.Errorf("cannot start server, got err: %v", err)
	}
	return c, addr, nil
}
