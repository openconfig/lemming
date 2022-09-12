package main

import (
	"flag"
	"fmt"
	"net"

	log "github.com/golang/glog"
	"github.com/openconfig/lemming"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	port   = pflag.Int("port", 6030, "localhost port to listen to.")
	target = pflag.String("target", "fakedut", "name of the fake target")
	// nolint:unused,varcheck
	enableDataplane = pflag.Bool("enable_dataplane", false, "Controls whether to enable dataplane")
)

func main() {
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("Failed to start listener: %v", err)
	}
	f, err := lemming.New(lis, *target)
	if err != nil {
		log.Fatalf("Failed to start lemming: %v", err)
	}
	defer f.Stop()

	log.Info("lemming initialization complete")
	select {}
}
