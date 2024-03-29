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

package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pktiopb "github.com/openconfig/lemming/dataplane/proto/packetio"
	"github.com/openconfig/lemming/dataplane/standalone/pkthandler/pktiohandler"

	log "github.com/golang/glog"
)

const (
	addr = "10.0.2.2:50000"
)

var portFile = flag.String("port_file", "/etc/sonic/pktio_ports.json", "File at which to include hostif info, for debugging only")

func main() {
	flag.Parse()

	ctx, cancelFn := context.WithTimeout(context.Background(), time.Minute)
	defer cancelFn()

	conn, err := grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Exit(err)
	}

	pktio := pktiopb.NewPacketIOClient(conn)

	h, err := pktiohandler.New(*portFile)
	if err != nil {
		log.Exit(err)
	}
	ctx, cancel := context.WithCancel(context.Background())

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		cancel()
	}()

	portCtl, err := pktio.HostPortControl(ctx)
	if err != nil {
		log.Exit(err)
	}
	packet, err := pktio.CPUPacketStream(ctx)
	if err != nil {
		log.Exit(err)
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		h.ManagePorts(portCtl)
		wg.Done()
	}()
	go func() {
		h.StreamPackets(packet)
		wg.Done()
	}()

	wg.Wait()
}
