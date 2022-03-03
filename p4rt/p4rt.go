package p4rt

import (
	"google.golang.org/grpc"

	p4rtpb "github.com/p4lang/p4runtime/go/p4/v1"
)

// Server is a fake p4rt implementation.
type Server struct {
	*p4rtpb.UnimplementedP4RuntimeServer
	s *grpc.Server
}

// New returns a new fake p4rt server.
func New(s *grpc.Server) *Server {
	srv := &Server{
		s: s,
	}
	p4rtpb.RegisterP4RuntimeServer(s, srv)

	return srv
}
