// Package lemming provides reference device to be used with ondatra.
package lemming

import (
	"fmt"
	"net"

	fgnmi "github.com/openconfig/lemming/gnmi"
	fgnoi "github.com/openconfig/lemming/gnoi"
	fgnsi "github.com/openconfig/lemming/gnsi"
	fgribi "github.com/openconfig/lemming/gribi"
	fp4rt "github.com/openconfig/lemming/p4rt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Device is the reference device implementation.
type Device struct {
	s           *grpc.Server
	addr        string
	stop        func()
	gnmiServer  *fgnmi.Server
	gnoiServer  *fgnoi.Server
	gribiServer *fgribi.Server
	gnsiServer  *fgnsi.Server
	p4rtServer  *fp4rt.Server
}

// New returns a new initialized device.
func New(addr string, opts ...grpc.ServerOption) (*Device, error) {
	s := grpc.NewServer(opts...)
	d := &Device{
		addr:        addr,
		s:           s,
		gnmiServer:  fgnmi.New(s),
		gnoiServer:  fgnoi.New(s),
		gribiServer: fgribi.New(s),
		gnsiServer:  fgnsi.New(s),
		p4rtServer:  fp4rt.New(s),
	}
	reflection.Register(s)
	if err := d.startServer(); err != nil {
		return nil, fmt.Errorf("failed to start device: %v", err)
	}
	return d, nil
}

// Addr returns the currently configured ip:port for the listening services.
func (d *Device) Addr() string {
	return d.addr
}

// Stop stops the listening services.
func (d *Device) Stop() {
	if d.stop == nil {
		return
	}
	d.stop()
}

func (d *Device) startServer() error {
	lis, err := net.Listen("tcp", d.addr)
	if err != nil {
		return fmt.Errorf("error creating TCP listener: %v", err)
	}
	d.addr = lis.Addr().String()
	go d.s.Serve(lis)

	d.stop = func() {
		d.s.Stop()
		lis.Close()
	}
	return nil
}
