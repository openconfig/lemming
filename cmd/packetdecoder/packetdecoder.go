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
	bytes, err := hex.DecodeString(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	packet := gopacket.NewPacket(bytes, layers.LayerTypeEthernet, gopacket.Default)
	fmt.Println(packet.Dump())
}
