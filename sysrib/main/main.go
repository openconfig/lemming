package main

import (
	"flag"
	"net"
	"os"

	log "github.com/golang/glog"
	"google.golang.org/grpc"

	pb "github.com/openconfig/lemming/proto/sysrib"
	"github.com/openconfig/lemming/sysrib"
)

func main() {
	flag.Parse()

	if err := os.RemoveAll(sysrib.SockAddr); err != nil {
		log.Fatal(err)
	}

	lis, err := net.Listen("unix", sysrib.SockAddr)
	if err != nil {
		log.Fatalf("listen error: %v", err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	s, err := sysrib.NewServer()
	if err != nil {
		log.Fatalf("error while creating sysrib server: %v", err)
	}
	s.AddInterface("eth0", 0, true, "192.0.0.0/8", "DEFAULT")
	pb.RegisterSysribServer(grpcServer, s)
	grpcServer.Serve(lis)
}
