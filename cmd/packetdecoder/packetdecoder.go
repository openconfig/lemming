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
	"encoding/hex"
	"fmt"
	"log"
	"os"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("expected exactly 1 arg: the hex string of the packet to decode")
	}
	bytes, err := hex.DecodeString(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	packet := gopacket.NewPacket(bytes, layers.LayerTypeEthernet, gopacket.Default)
	fmt.Println(packet.Dump())
}
