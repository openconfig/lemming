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

package gnmi

import (
	"context"
	"errors"
	"fmt"
	"time"

	log "github.com/golang/glog"
	"github.com/openconfig/gnmi/cache"
	gpb "github.com/openconfig/gnmi/proto/gnmi"
	"github.com/openconfig/gnmi/subscribe"
	"google.golang.org/protobuf/proto"
)

var (
	// metadataUpdatePeriod is the period of time after which the metadata for the collector
	// is updated to the client.
	metadataUpdatePeriod = 30 * time.Second
	// sizeUpdatePeriod is the period of time after which the storage size information for
	// the collector is updated to the client.
	sizeUpdatePeriod = 30 * time.Second
)

// Collector is a basic gNMI target that supports only the Subscribe
// RPC, and acts as a cache for exactly one target.
type Collector struct {
	cache *cache.Cache
	// name is the cache target name.
	name string
	// inCh is a channel use to write new SubscribeResponses to the client.
	inCh chan *gpb.SubscribeResponse
	// stopFn is the function used to stop the server.
	stopFn func()
}

// NewCollector returns an initialized gNMI Collector implementation.
//
// To create a gNMI server that supports gnmi.Set as well, use New() instead.
func NewCollector(targetName string) *Collector {
	return &Collector{
		cache: cache.New([]string{targetName}),
		name:  targetName,
		inCh:  make(chan *gpb.SubscribeResponse),
	}
}

// Start starts the collector and returns a linked gNMI server that supports
// gnmi.Subscribe.
func (c *Collector) Start(ctx context.Context, sendMeta bool) (*subscribe.Server, error) {
	t := c.cache.GetTarget(c.name)

	subscribeSrv, err := subscribe.NewServer(c.cache)
	if err != nil {
		return nil, fmt.Errorf("could not instantiate gNMI server: %v", err)
	}
	c.cache.SetClient(subscribeSrv.Update)

	if sendMeta {
		go periodic(metadataUpdatePeriod, c.cache.UpdateMetadata)
		go periodic(sizeUpdatePeriod, c.cache.UpdateSize)
	}
	t.Connect()

	// start our single collector from the input channel.
	go func() {
		for {
			select {
			case msg := <-c.inCh:
				if err := c.handleUpdate(msg); err != nil {
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return subscribeSrv, nil
}

// TargetUpdate provides an input gNMI SubscribeResponse to update the
// cache and clients with.
func (c *Collector) TargetUpdate(m *gpb.SubscribeResponse) {
	c.inCh <- m
}

// Stop halts the running collector.
func (c *Collector) Stop() {
	c.stopFn()
}

// handleUpdate handles an input gNMI SubscribeResponse that is received by
// the target.
func (c *Collector) handleUpdate(resp *gpb.SubscribeResponse) error {
	switch v := resp.Response.(type) {
	case *gpb.SubscribeResponse_Update:
		return c.GnmiUpdate(v.Update)
	case *gpb.SubscribeResponse_SyncResponse:
		c.cache.GetTarget(c.name).Sync()
	case *gpb.SubscribeResponse_Error:
		return fmt.Errorf("error in response: %s", v)
	default:
		return fmt.Errorf("unknown response %T: %s", v, v)
	}
	return nil
}

// GnmiUpdate sends a pb.Notification into the target cache.
//
// It simply forwards it to the gNMI cache implementation, populating the
// target (without copying the message) if empty.
func (c *Collector) GnmiUpdate(n *gpb.Notification) error {
	t := c.cache.GetTarget(c.name)
	// If target is not specified, then set it to the initialized
	// value.
	if n.GetPrefix().GetTarget() == "" {
		if n.Prefix == nil {
			n.Prefix = &gpb.Path{}
		} else {
			n.Prefix = proto.Clone(n.Prefix).(*gpb.Path)
		}
		n.Prefix.Target = c.name
	}
	if n.Timestamp == 0 {
		// Set timestamp here in order to minimize latency and reduce chance for "update is stale" error.
		n.Timestamp = time.Now().UnixNano()
	}
	err := t.GnmiUpdate(n)
	if errors.Is(err, cache.ErrStale) {
		log.Errorf("Unexpected stale update while updating cache: %v, (current time: %v, notif time: %v)", err, time.Now().UnixNano(), n.Timestamp)
		return nil
	}
	return err
}

// periodic runs the function fn every period.
func periodic(period time.Duration, fn func()) {
	if period == 0 {
		return
	}
	t := time.NewTicker(period)
	defer t.Stop()
	for range t.C {
		fn()
	}
}
