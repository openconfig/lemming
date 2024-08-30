// Copyright 2024 Google LLC
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

// Package registry is the registry of protocol on the dataplane.
package protocol

import (
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/openconfig/lemming/dataplane/forwarding/util/queue"

	"github.com/openconfig/lemming/dataplane/proto/packetio"

	log "github.com/golang/glog"
	pktiopb "github.com/openconfig/lemming/dataplane/proto/packetio"
)

// Handler handles a protocol.
type Handler interface {
	Matched(*packetio.PacketOut) bool  // returns true if this protocol can handle the packet.
	Process(*packetio.PacketOut) error // processes the packet and update internal data structure if necessary.
}

// Registry is the repository of protocol handler.
type Registry struct {
	pktiopb.PacketIO_CPUPacketStreamClient
	bypassQ *queue.Queue                           // the queue as a buffer for packet peeking/filtering.
	psc     pktiopb.PacketIO_CPUPacketStreamClient // the original packet stream client.
	doneCh  chan struct{}
	mu      sync.Mutex
	reg     map[string]Handler // map the protocol name to its handler.
}

// NewRegistry takes a packet stream client and returns an empty registry.
func NewRegistry(psc pktiopb.PacketIO_CPUPacketStreamClient) (*Registry, error) {
	q, err := queue.NewUnbounded("recv")
	if err != nil {
		return nil, err
	}
	q.Run()
	pr := &Registry{
		bypassQ: q,
		psc:     psc,
		reg:     map[string]Handler{},
		doneCh:  make(chan struct{}),
	}
	return pr, nil
}

// Start starts a goroutine to intecept the packets from the psc stream client.
// The packet will be sent to the protocol handler if it is available in the
// registry, or sent to bypass queue so that it will be processed by the
// default packet IO manager.
func (r *Registry) Start() {
	go func() {
		for {
			select {
			case <-r.doneCh:
				log.Info("Protocol registry stopped.")
				r.bypassQ.Close()
				return
			default:
				pkt, err := r.psc.Recv()
				if err != nil {
					if errors.Is(err, io.EOF) {
						log.Warning("Received EOF from server, exiting")
						r.bypassQ.Close()
						return
					}
					log.Warningf("Received errors: %v", err)
					continue
				}
				processed := false
				for name, ph := range r.reg {
					if ph.Matched(pkt) {
						if err := ph.Process(pkt); err != nil {
							log.Warningf("Error occurred when processing %d packet: %v", name, err)
						}
						processed = true
						break
					}
				}
				if !processed {
					r.bypassQ.Write(pkt)
				}
				time.Sleep(time.Millisecond)
			}
		}
	}()
}

// Stop stops the registry.
func (r *Registry) Stop() {
	r.doneCh <- struct{}{}
}

// Recv returns the packet that is not processed by any protocol.
func (r *Registry) Recv() (*packetio.PacketOut, error) {
	pkt, ok := <-r.bypassQ.Receive()
	if !ok {
		// The queue is closed.
		return nil, nil
	}
	return &pktiopb.PacketOut{Packet: pkt.(*pktiopb.Packet)}, nil
}

// Send sends the packet via the streaming client it holds.
func (r *Registry) Send(pkt *packetio.PacketIn) error {
	return r.psc.Send(pkt)
}

// Register adds a new protocol handler to the registry.
func (r *Registry) Register(name string, h Handler) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.reg[name]; ok {
		return fmt.Errorf("Protocol %q is existing", name)
	}
	log.Infof("Register protocl %q", name)
	r.reg[name] = h
	return nil
}

// Deregister removes a protocol handler from the registry.
func (r *Registry) Deregister(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.reg[name]; !ok {
		return fmt.Errorf("Protocol %q not found", name)
	}
	log.Infof("Deregister protocl %q", name)
	delete(r.reg, name)
	return nil
}
