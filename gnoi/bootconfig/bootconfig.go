// Package bootconfig provides the bootconfig implementation
// for Lemming.
package bootconfig

import (
	bpb "github.com/openconfig/gnoi/bootconfig"
)

type Server struct {
	bpb.UnimplementedBootConfigServer
}

func New() *Server {
	return &Server{}
}
