package client

import (
	"context"
	"fmt"

	gpb "github.com/openconfig/gnmi/proto/gnmi"
	"github.com/openconfig/lemming/gnmi/gnmit"
	"github.com/openconfig/ygnmi/ygnmi"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/local"
)

type localClient struct {
	gpb.GNMIClient
	setClient gpb.GNMIClient
}

func (m *localClient) Set(ctx context.Context, in *gpb.SetRequest, opts ...grpc.CallOption) (*gpb.SetResponse, error) {
	return m.setClient.Set(ctx, in, opts...)
}

// NewLocal returns a gNMI client connected to the local gNMI cache.
func NewLocal() (gpb.GNMIClient, error) {
	setConn, err := grpc.Dial(fmt.Sprintf("unix:///%s", gnmit.DatastoreAddress), grpc.WithTransportCredentials(local.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to dial unix socket: %v", err)
	}
	cacheConn, err := grpc.Dial(fmt.Sprintf("localhost:%d", viper.GetInt("port")), grpc.WithTransportCredentials(local.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to dial cache socket: %v", err)
	}
	return &localClient{
		GNMIClient: gpb.NewGNMIClient(cacheConn),
		setClient:  gpb.NewGNMIClient(setConn),
	}, nil
}

// NewYGNMIClient returns ygnmi client connected to the local gNMI cache.
func NewYGNMIClient() (*ygnmi.Client, error) {
	gClient, err := NewLocal()
	if err != nil {
		return nil, err
	}
	return ygnmi.NewClient(gClient, ygnmi.WithTarget(viper.GetString("target")))
}
