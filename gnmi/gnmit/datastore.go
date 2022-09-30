package gnmit

import (
	"context"
	"sync"

	log "github.com/golang/glog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	gpb "github.com/openconfig/gnmi/proto/gnmi"
	"github.com/openconfig/lemming/gnmi/oc"
	"github.com/openconfig/ygot/ytypes"
)

const (
	// Unix socket for central datastore.
	DatastoreAddress = "/tmp/datastore.api"
)

type DatastoreServer struct {
	gpb.UnimplementedGNMIServer // For forward-compatibility

	smu    sync.Mutex
	schema *ytypes.Schema

	gnmiServer *GNMIServer
}

func NewDatastoreServer(gnmiServer *GNMIServer) (*DatastoreServer, error) {
	schema, err := oc.Schema()
	if err != nil {
		return nil, err
	}
	if err := SetupSchema(schema); err != nil {
		return nil, err
	}
	return &DatastoreServer{
		gnmiServer: gnmiServer,
		schema:     schema,
	}, nil
}

func (d *DatastoreServer) Set(ctx context.Context, req *gpb.SetRequest) (*gpb.SetResponse, error) {
	d.smu.Lock()
	defer d.smu.Unlock()

	log.V(1).Infof("operational state datastore service received SetRequest: %v", req)
	if d.schema == nil {
		return d.gnmiServer.UnimplementedGNMIServer.Set(ctx, req)
	}
	// TODO(wenbli): Reject values that modify config values. We only allow modifying state through this server.
	if err := set(d.schema, d.gnmiServer.c.cache, d.gnmiServer.c.name, req, false); err != nil {
		return &gpb.SetResponse{}, status.Errorf(codes.Aborted, "%v", err)
	}

	return &gpb.SetResponse{}, nil
}
