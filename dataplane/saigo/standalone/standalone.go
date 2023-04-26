package main

import "C"

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/openconfig/lemming/dataplane/forwarding"

	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

//export initialize
func initialize(port int) {
	fwdSrv := forwarding.New("engine")

	lis, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := grpc.NewServer(grpc.Creds(insecure.NewCredentials()))
	fwdpb.RegisterForwardingServer(srv, fwdSrv)
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to server: %v", err)
	}
}

func main() {
}
