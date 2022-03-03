package gnmi

import (
	"context"
	"time"

	gnmipb "github.com/openconfig/gnmi/proto/gnmi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// Server is a fake gNMI implementation.
type Server struct {
	s            *grpc.Server
	subscription int
	Responses    [][]*gnmipb.SubscribeResponse
	GetResponses []interface{}
	Errs         []error
}

// New returns a new fake gNMI server.
func New(s *grpc.Server) *Server {
	srv := &Server{
		s: s,
	}
	gnmipb.RegisterGNMIServer(s, srv)
	return srv
}

func (s *Server) Capabilities(ctx context.Context, req *gnmipb.CapabilityRequest) (*gnmipb.CapabilityResponse, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "Fake Unimplemented")
}

func (s *Server) Get(ctx context.Context, req *gnmipb.GetRequest) (*gnmipb.GetResponse, error) {
	if len(s.GetResponses) == 0 {
		return nil, grpc.Errorf(codes.Unimplemented, "Fake Unimplemented")
	}
	resp := s.GetResponses[0]
	s.GetResponses = s.GetResponses[1:]
	switch v := resp.(type) {
	case error:
		return nil, v
	case *gnmipb.GetResponse:
		return v, nil
	default:
		return nil, grpc.Errorf(codes.DataLoss, "Unknown message type: %T", resp)
	}
}

func (s *Server) Set(ctx context.Context, req *gnmipb.SetRequest) (*gnmipb.SetResponse, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "Fake Unimplemented")
}

func (s *Server) Subscribe(stream gnmipb.GNMI_SubscribeServer) error {
	_, err := stream.Recv()
	if err != nil {
		return err
	}
	srs := s.Responses[s.subscription]
	srErr := s.Errs[s.subscription]
	s.subscription++
	for _, sr := range srs {
		stream.Send(sr)
	}
	time.Sleep(5 * time.Second)
	return srErr
}
