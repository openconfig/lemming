package gnmit

import (
	"context"

	log "github.com/golang/glog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	gpb "github.com/openconfig/gnmi/proto/gnmi"
	dspb "github.com/openconfig/lemming/proto/datastore"
)

const (
	// Unix socket for central datastore.
	DatastoreAddress = "/tmp/datastore.api"
)

type DatastoreServer struct {
	dspb.UnimplementedDatastoreServer // For forward-compatibility

	update UpdateFn
}

func NewDatastoreServer(updateFn UpdateFn) *DatastoreServer {
	return &DatastoreServer{update: updateFn}
}

func (d *DatastoreServer) Update(_ context.Context, n *gpb.Notification) (*dspb.UpdateResponse, error) {
	log.Info(n)
	err := d.update(n)
	if err != nil {
		return &dspb.UpdateResponse{}, status.Errorf(codes.Aborted, "%v", err)
	}
	return &dspb.UpdateResponse{}, nil
}
