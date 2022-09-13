package gnmit

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	gpb "github.com/openconfig/gnmi/proto/gnmi"
)

const (
	// Unix socket for central datastore.
	DatastoreAddress = "/tmp/datastore.api"
)

type DatastoreServer struct {
	gpb.UnimplementedGNMIServer // For forward-compatibility

	gnmiServer *GNMIServer
}

func NewDatastoreServer(gnmiServer *GNMIServer) *DatastoreServer {
	return &DatastoreServer{gnmiServer: gnmiServer}
}

func (d *DatastoreServer) Set(_ context.Context, req *gpb.SetRequest) (*gpb.SetResponse, error) {
	if err := d.gnmiServer.set(req, false); err != nil {
		return &gpb.SetResponse{}, status.Errorf(codes.Aborted, "%v", err)
	}

	// SetRequest has been validated, so we update the cache.
	deletes := append([]*gpb.Path{}, req.Delete...)
	for _, update := range req.Replace {
		deletes = append(deletes, update.Path)
	}
	t := d.gnmiServer.c.cache.GetTarget(d.gnmiServer.c.name)
	t.GnmiUpdate(&gpb.Notification{
		Prefix: req.Prefix,
		Delete: deletes,
		Update: req.Update,
	})
	return &gpb.SetResponse{}, nil
}
