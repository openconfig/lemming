package gnmit

import (
	"context"
	"time"

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
	notif := &gpb.Notification{
		Prefix:    req.Prefix,
		Delete:    req.Delete,
		Timestamp: time.Now().UnixNano(),
	}
	for _, upd := range req.Replace {
		notif.Delete = append(notif.Delete, upd.Path)
		notif.Update = append(notif.Update, upd)
	}
	notif.Update = append(notif.Update, req.Update...)
	if err := d.gnmiServer.set(notif); err != nil {
		return &gpb.SetResponse{}, status.Errorf(codes.Aborted, "%v", err)
	}
	return &gpb.SetResponse{}, nil
}
