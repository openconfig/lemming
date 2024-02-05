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

package main

import "C"

import (
	"context"

	log "github.com/golang/glog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/openconfig/lemming/dataplane/cpusink"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

const addr = "10.0.2.2:50000"

//export StartSink
func StartSink() {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	sink, err := cpusink.New(fwdpb.NewForwardingClient(conn))
	if err != nil {
		log.Fatal(err)
	}
	if err := sink.HandleIPUpdates(context.Background()); err != nil {
		log.Fatal(err)
	}
	go func() {
		if err := sink.ReceivePackets(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()
	log.Info("started sink")
}

func main() {
}
