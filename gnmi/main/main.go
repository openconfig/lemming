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

package main

import (
	"context"
	"flag"
	"fmt"

	log "github.com/golang/glog"
	"github.com/openconfig/lemming/gnmi/fakedevice"
)

var (
	port   = flag.Int("port", 1234, "localhost port to listen to.")
	target = flag.String("target", "fakedut", "name of the fake target")
)

func init() {
	flag.Parse()
}

func main() {
	_, _, err := fakedevice.NewTarget(context.Background(), fmt.Sprintf(":%d", *port), *target)
	if err != nil {
		log.Fatal(err)
	}

	select {}
}
