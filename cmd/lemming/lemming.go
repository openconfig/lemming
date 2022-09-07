package main

import (
	"flag"
	"fmt"
	"net"

	log "github.com/golang/glog"
	"github.com/openconfig/lemming"
)

var (
	port   = flag.Int("port", 6030, "localhost port to listen to.")
	target = flag.String("target", "fakedut", "name of the fake target")
)

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("Failed to start listener: %v", err)
	}
	f, err := lemming.New(lis, *target)
	if err != nil {
		log.Fatalf("Failed to start lemming: %v", err)
	}
	defer f.Stop()

	select {}
}
