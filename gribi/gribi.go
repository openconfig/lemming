package gribi

import (
	gribipb "github.com/openconfig/gribi/v1/proto/service"
	"google.golang.org/grpc"
)

// Server is a fake gNSI implementation.
type Server struct {
	*gribipb.UnimplementedGRIBIServer
	s *grpc.Server
}

// New returns a new fake gNMI server.
func New(s *grpc.Server) *Server {
	srv := &Server{
		s: s,
	}
	gribipb.RegisterGRIBIServer(s, srv)

	return srv
}
