// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package binding

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/netip"
	"path/filepath"
	"time"

	"github.com/google/gopacket"
	"github.com/google/uuid"
	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/openconfig/gnoigo"

	"github.com/openconfig/ondatra/binding"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/local"

	"github.com/openconfig/lemming"
	"github.com/openconfig/lemming/dataplane/dplaneopts"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdcontext"
	"github.com/openconfig/lemming/dataplane/forwarding/util/queue"

	acctzpb "github.com/openconfig/gnsi/acctz"
	authzpb "github.com/openconfig/gnsi/authz"
	certzpb "github.com/openconfig/gnsi/certz"
	credzpb "github.com/openconfig/gnsi/credentialz"
	pathzpb "github.com/openconfig/gnsi/pathz"
	grpb "github.com/openconfig/gribi/v1/proto/service"
	opb "github.com/openconfig/ondatra/proto"

	saipb "github.com/openconfig/lemming/dataplane/proto"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

// LocalBind is an Ondatra binding for in-process testbed. Only Lemming and Magna are supported.
type LocalBind struct {
	binding.Binding
}

// Local is a local (in-process) binding for lemming and magna.
func Local(topoDir string) func() (binding.Binding, error) {
	dir, _ := filepath.Abs(topoDir)
	testbedFile := filepath.Join(dir, "testbed.pb.txt")

	flag.Set("testbed", testbedFile)
	return func() (binding.Binding, error) {
		return &LocalBind{}, nil
	}
}

type localLemming struct {
	binding.AbstractDUT
	l     *lemming.Device
	dutID string
	addr  string
}

type gnsiClient struct {
	binding.AbstractGNSIClients
	conn grpc.ClientConnInterface
}

func (c *gnsiClient) Acctz() acctzpb.AcctzClient { return acctzpb.NewAcctzClient(c.conn) }
func (c *gnsiClient) Authz() authzpb.AuthzClient { return authzpb.NewAuthzClient(c.conn) }
func (c *gnsiClient) Certz() certzpb.CertzClient { return certzpb.NewCertzClient(c.conn) }
func (c *gnsiClient) Credentialz() credzpb.CredentialzClient {
	return credzpb.NewCredentialzClient(c.conn)
}
func (c *gnsiClient) Pathz() pathzpb.PathzClient { return pathzpb.NewPathzClient(c.conn) }

func (l *localLemming) DialGNMI(ctx context.Context, opts ...grpc.DialOption) (gnmi.GNMIClient, error) {
	conn, err := grpc.DialContext(ctx, net.JoinHostPort(l.addr, fmt.Sprint(gnmiPort)), opts...)
	if err != nil {
		return nil, err
	}
	return gnmi.NewGNMIClient(conn), nil
}

func (l *localLemming) DialGNOI(ctx context.Context, opts ...grpc.DialOption) (gnoigo.Clients, error) {
	conn, err := grpc.DialContext(ctx, net.JoinHostPort(l.addr, fmt.Sprint(gnmiPort)), opts...)
	if err != nil {
		return nil, err
	}
	return gnoigo.NewClients(conn), nil
}

func (l *localLemming) DialGNSI(ctx context.Context, opts ...grpc.DialOption) (binding.GNSIClients, error) {
	conn, err := grpc.DialContext(ctx, net.JoinHostPort(l.addr, fmt.Sprint(gnmiPort)), opts...)
	if err != nil {
		return nil, err
	}
	return &gnsiClient{conn: conn}, nil
}

func (l *localLemming) DialGRIBI(ctx context.Context, opts ...grpc.DialOption) (grpb.GRIBIClient, error) {
	conn, err := grpc.DialContext(ctx, net.JoinHostPort(l.addr, fmt.Sprint(gribiPort)), opts...)
	if err != nil {
		return nil, err
	}
	return grpb.NewGRIBIClient(conn), nil
}

const (
	gnmiPort      = 9339
	gribiPort     = 9340
	bgpPort       = 1179
	dataplanePort = 50000
)

// Reserve creates a new local binding.
func (lb *LocalBind) Reserve(ctx context.Context, tb *opb.Testbed, _, _ time.Duration, _ map[string]string) (*binding.Reservation, error) {
	resv := binding.Reservation{
		ID:   uuid.New().String(),
		DUTs: make(map[string]binding.DUT),
		ATEs: make(map[string]binding.ATE),
	}

	portMgr := &portMgr{
		ports: map[string]*chanPort{},
	}

	// for _, ate := range tb.Ates {
	// }

	for _, dut := range tb.Duts {
		addr, err := findAvailableLoopbackIP()
		if err != nil {
			return nil, err
		}
		dutID := uuid.New().String()

		l, err := lemming.New(dut.Id, fmt.Sprintf("unix:/tmp/zserv-test%s.api", dutID),
			lemming.WithBGPPort(bgpPort),
			lemming.WithGNMIAddr(net.JoinHostPort(addr, fmt.Sprint(gnmiPort))),
			lemming.WithGRIBIAddr(net.JoinHostPort(addr, fmt.Sprint(gribiPort))),
			lemming.WithTransportCreds(local.NewCredentials()),
			lemming.WithDataplane(true),
			lemming.WithDataplaneOpts(
				dplaneopts.WithAddrPort(net.JoinHostPort(addr, fmt.Sprint(dataplanePort))),
				dplaneopts.WithReconcilation(false),
				dplaneopts.WithPortType(fwdpb.PortType_PORT_TYPE_FAKE),
			),
		)
		if err != nil {
			return nil, err
		}
		boundLemming := &localLemming{
			l:     l,
			dutID: dutID,
			AbstractDUT: binding.AbstractDUT{
				Dims: &binding.Dims{
					Name:          dut.Id,
					Vendor:        opb.Device_OPENCONFIG,
					HardwareModel: "LEMMING",
					Ports:         make(map[string]*binding.Port),
				},
			},
		}

		fwdCtx, err := boundLemming.l.Dataplane().SaiServer().FindContext(&fwdpb.ContextId{Id: "lucius"})
		if err != nil {
			return nil, err
		}

		// In the dut, use the portMgr to create ports.
		// Ondatra Port ID: Value from the topology textproto.
		// Ondatra Port Name: Dataplane port id.
		fwdCtx.FakePortManager = portMgr
		dplaneConn, err := boundLemming.l.Dataplane().Conn()
		if err != nil {
			return nil, err
		}
		pc := saipb.NewPortClient(dplaneConn)
		// For each port on the topology textproto.
		// Create a saipb port and accociate the ondatra dut ID and port ID with the saipb OID.
		for _, port := range dut.Ports {
			resp, err := pc.CreatePort(ctx, &saipb.CreatePortRequest{})
			if err != nil {
				return nil, err
			}
			portMgr.ports[fmt.Sprintf("%s:%s", dut.Id, port.Id)] = portMgr.lastCreatedPort
			boundLemming.AbstractDUT.Dims.Ports[port.Id] = &binding.Port{
				Name: fmt.Sprint(resp.Oid),
			}
		}
		resv.DUTs[dut.Id] = boundLemming
	}

	for _, l := range tb.Links {
		if err := portMgr.linkPorts(l.A, l.B); err != nil {
			return nil, err
		}
	}

	return &resv, nil
}

// Release releases the reserved testbed.
func (lb *LocalBind) Release(context.Context) error {
	return nil
}

// findAvailableLoopbackIP finds an unused loopback IP by attempting to open a tcp socket on the gNMI port.
func findAvailableLoopbackIP() (string, error) {
	var addr string
	for i := byte(2); i < 254; i++ {
		addr = netip.AddrFrom4([4]byte{127, 0, 0, i}).String()
		ln, err := net.Listen("tcp", net.JoinHostPort(addr, fmt.Sprint(gnmiPort)))
		if err == nil {
			ln.Close()
			break
		}
	}
	if addr == "" {
		return "", fmt.Errorf("failed to find available ip")
	}
	return addr, nil
}

type chanPort struct {
	txQueue *queue.Queue
	rxQueue *queue.Queue
}

func (p *chanPort) WritePacketData(data []byte) error {
	return p.txQueue.Write(data)
}

func (p *chanPort) ReadPacketData() ([]byte, gopacket.CaptureInfo, error) {
	if p.rxQueue == nil {
		return nil, gopacket.CaptureInfo{}, nil
	}
	data := (<-p.rxQueue.Receive()).([]byte)
	return data, gopacket.CaptureInfo{
		Timestamp:     time.Now(),
		CaptureLength: len(data),
		Length:        len(data),
	}, nil
}

type portMgr struct {
	ports           map[string]*chanPort
	lastCreatedPort *chanPort
}

func (lm *portMgr) CreatePort(name string) (fwdcontext.Port, error) {
	q, err := queue.NewUnbounded(fmt.Sprintf("%s-tx", name))
	if err != nil {
		return nil, err
	}
	p := &chanPort{
		txQueue: q,
	}
	lm.lastCreatedPort = p
	q.Run()
	return p, nil
}

func (lm *portMgr) linkPorts(a, b string) error {
	aPort := lm.ports[a]
	bPort := lm.ports[b]
	if aPort == nil || bPort == nil {
		return fmt.Errorf("ports do not exist")
	}
	aPort.rxQueue = bPort.txQueue
	bPort.rxQueue = aPort.txQueue

	return nil
}
