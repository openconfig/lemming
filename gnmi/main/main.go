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
