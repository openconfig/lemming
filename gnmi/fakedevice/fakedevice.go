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
	"time"

	log "github.com/golang/glog"
	"github.com/openconfig/goyang/pkg/yang"
	"github.com/openconfig/lemming/gnmi/gnmiclient"
	"github.com/openconfig/lemming/gnmi/oc"
	"github.com/openconfig/lemming/gnmi/oc/ocpath"
	"github.com/openconfig/ygnmi/ygnmi"
	"github.com/openconfig/ygot/ygot/pathtranslate"

	gpb "github.com/openconfig/gnmi/proto/gnmi"
)

const (
	DefaultNetworkInstance = "DEFAULT"
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

func StartBootTimeTask(ctx context.Context, gClient gpb.GNMIClient, target string) error {
	yclient, err := ygnmi.NewClient(gClient, ygnmi.WithTarget(target))
	if err != nil {
		return err
	}

	gnmiclient.Replace(ctx, yclient, ocpath.Root().System().BootTime().State(), uint64(time.Now().UnixNano()))

	return nil
}

func StartCurrentDateTimeTask(ctx context.Context, gClient gpb.GNMIClient, target string) error {
	yclient, err := ygnmi.NewClient(gClient, ygnmi.WithTarget(target))
	if err != nil {
		return err
	}

	tick := time.NewTicker(time.Second)

	periodic := func() error {
		_, err := gnmiclient.Replace(ctx, yclient, ocpath.Root().System().CurrentDatetime().State(), time.Now().Format(time.RFC3339))
		return err
	}

	if err := periodic(); err != nil {
		return err
	}

	go func() {
		for range tick.C {
			if err := periodic(); err != nil {
				log.Errorf("currentDateTimeTask error: %v", err)
				return
			}
		}
	}()

	return nil
}

// StartSystemBaseTask handles some of the logic for the base systems feature
// profile using ygnmi as the client.
func StartSystemBaseTask(ctx context.Context, gClient gpb.GNMIClient, target string) error {
	yclient, err := ygnmi.NewClient(gClient, ygnmi.WithTarget(target))
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

	go func() {
		if _, err := systemWatcher.Await(); err != nil {
			log.Warningf("Sysrib interface watcher has stopped: %v", err)
		}
	}()
	return nil
}
