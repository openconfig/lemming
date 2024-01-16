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
	"io"
	"net"
	"net/netip"
	"path/filepath"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/uuid"
	"github.com/open-traffic-generator/snappi/gosnappi"
	"github.com/open-traffic-generator/snappi/gosnappi/otg"
	"github.com/openconfig/gnoigo"
	"github.com/openconfig/magna/flows/common"
	"github.com/openconfig/magna/flows/ip"
	"github.com/openconfig/magna/flows/mpls"
	"github.com/openconfig/magna/intf"
	"github.com/openconfig/magna/lwotg"
	"github.com/openconfig/magna/lwotgtelem"
	"github.com/openconfig/magna/telemetry/arp"
	"k8s.io/klog"

	"github.com/openconfig/ondatra/binding"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/local"
	"google.golang.org/grpc/reflection"

	"github.com/openconfig/lemming"
	"github.com/openconfig/lemming/dataplane/dplaneopts"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdcontext"
	"github.com/openconfig/lemming/dataplane/forwarding/util/queue"

	gpb "github.com/openconfig/gnmi/proto/gnmi"
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

// gnsiClient implemts the binding.GNSIClients interface
type gnsiClient struct {
	binding.AbstractGNSIClients
	conn grpc.ClientConnInterface
}

// Acctz returns the gNSI acctz client.
func (c *gnsiClient) Acctz() acctzpb.AcctzClient { return acctzpb.NewAcctzClient(c.conn) }

// Authz returns the gNSI authz client.
func (c *gnsiClient) Authz() authzpb.AuthzClient { return authzpb.NewAuthzClient(c.conn) }

// Certz returns the gNSI certz client.
func (c *gnsiClient) Certz() certzpb.CertzClient { return certzpb.NewCertzClient(c.conn) }

// Credentialz returns the gNSI credentialz client.
func (c *gnsiClient) Credentialz() credzpb.CredentialzClient {
	return credzpb.NewCredentialzClient(c.conn)
}

// Pathz returns the gNSI pathz client.
func (c *gnsiClient) Pathz() pathzpb.PathzClient { return pathzpb.NewPathzClient(c.conn) }

// DialGNMI returns a gNMI client for the dut.
func (l *localLemming) DialGNMI(ctx context.Context, opts ...grpc.DialOption) (gpb.GNMIClient, error) {
	opts = append(opts, grpc.WithTransportCredentials(local.NewCredentials()))
	conn, err := grpc.DialContext(ctx, net.JoinHostPort(l.addr, fmt.Sprint(gnmiPort)), opts...)
	if err != nil {
		return nil, err
	}
	return gpb.NewGNMIClient(conn), nil
}

// DialGNOI returns a gNOI client for the dut.
func (l *localLemming) DialGNOI(ctx context.Context, opts ...grpc.DialOption) (gnoigo.Clients, error) {
	opts = append(opts, grpc.WithTransportCredentials(local.NewCredentials()))
	conn, err := grpc.DialContext(ctx, net.JoinHostPort(l.addr, fmt.Sprint(gnmiPort)), opts...)
	if err != nil {
		return nil, err
	}
	return gnoigo.NewClients(conn), nil
}

// DialGNSI returns a gNSI client for the dut.
func (l *localLemming) DialGNSI(ctx context.Context, opts ...grpc.DialOption) (binding.GNSIClients, error) {
	opts = append(opts, grpc.WithTransportCredentials(local.NewCredentials()))
	conn, err := grpc.DialContext(ctx, net.JoinHostPort(l.addr, fmt.Sprint(gnmiPort)), opts...)
	if err != nil {
		return nil, err
	}
	return &gnsiClient{conn: conn}, nil
}

// DialGRIBI returns a gRIBI client for the dut.
func (l *localLemming) DialGRIBI(ctx context.Context, opts ...grpc.DialOption) (grpb.GRIBIClient, error) {
	opts = append(opts, grpc.WithTransportCredentials(local.NewCredentials()))
	conn, err := grpc.DialContext(ctx, net.JoinHostPort(l.addr, fmt.Sprint(gribiPort)), opts...)
	if err != nil {
		return nil, err
	}
	return grpb.NewGRIBIClient(conn), nil
}

type localMagna struct {
	binding.AbstractATE
	addr string
}

// DialGNMI returns a gNMI client for the dut.
func (m *localMagna) DialGNMI(ctx context.Context, opts ...grpc.DialOption) (gpb.GNMIClient, error) {
	opts = append(opts, grpc.WithTransportCredentials(local.NewCredentials()))
	conn, err := grpc.DialContext(ctx, net.JoinHostPort(m.addr, fmt.Sprint(gnmiPort)), opts...)
	if err != nil {
		return nil, err
	}
	return gpb.NewGNMIClient(conn), nil
}

// DialOTG returns a OTH client for the dut.
func (m *localMagna) DialOTG(ctx context.Context, opts ...grpc.DialOption) (gosnappi.GosnappiApi, error) {
	opts = append(opts, grpc.WithTransportCredentials(local.NewCredentials()))
	conn, err := grpc.DialContext(ctx, net.JoinHostPort(m.addr, fmt.Sprint(otgPort)), opts...)
	if err != nil {
		return nil, err
	}
	api := gosnappi.NewApi()
	api.NewGrpcTransport().SetClientConnection(conn).SetRequestTimeout(30 * time.Second)

	return api, nil
}

const (
	gnmiPort      = 9339
	gribiPort     = 9340
	bgpPort       = 1179
	dataplanePort = 50000
	otgPort       = 50001
)

// TODO: Implement this.
type accessor struct {
	intf.NetworkAccessor
}

// ARPList lists the set of ARP neighbours on the system.
func (accessor) ARPList() ([]*intf.ARPEntry, error) {
	return nil, nil
}

// ARPSubscribe writes changes to the ARP table to the channel updates.
func (accessor) ARPSubscribe(chan intf.ARPUpdate, chan struct{}) error {
	return nil
}

// InterfaceAdddresses lists the IP addresses configured on a particular interface.
func (accessor) InterfaceAddresses(string) ([]*net.IPNet, error) {
	return nil, nil
}

// AddInterfaceIP adds address ip to the interface name.
func (accessor) AddInterfaceIP(string, *net.IPNet) error {
	return nil
}

// Reserve creates a new local binding.
func (lb *LocalBind) Reserve(ctx context.Context, tb *opb.Testbed, _, _ time.Duration, _ map[string]string) (*binding.Reservation, error) {
	resv := binding.Reservation{
		ID:   uuid.New().String(),
		DUTs: make(map[string]binding.DUT),
		ATEs: make(map[string]binding.ATE),
	}

	portMgr := &portMgr{
		dut: &dutPortMgr{
			ports: map[string]*chanPort{},
		},
		ate: &atePortMgr{
			ports: map[string]*chanPort{},
		},
	}

	intf.OverrideAccessor(&accessor{})
	common.OverridePortCreator(portMgr.ate)

	for _, ate := range tb.Ates {
		magna, err := lb.createATE(ctx, ate, portMgr.ate)
		if err != nil {
			return nil, err
		}
		resv.ATEs[ate.Id] = magna
	}

	for _, dut := range tb.Duts {
		lemming, err := lb.createDUT(ctx, dut, portMgr.dut)
		if err != nil {
			return nil, err
		}
		resv.DUTs[dut.Id] = lemming
	}

	for _, l := range tb.Links {
		if err := portMgr.linkPorts(l.A, l.B); err != nil {
			return nil, err
		}
	}

	return &resv, nil
}

func (lb *LocalBind) createDUT(ctx context.Context, dut *opb.Device, portMgr *dutPortMgr) (*localLemming, error) {
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
	return boundLemming, nil
}

// TODO: this should probably be a library in magna.
func (lb *LocalBind) createATE(_ context.Context, ate *opb.Device, portMgr *atePortMgr) (*localMagna, error) {
	addr, err := findAvailableLoopbackIP()
	if err != nil {
		return nil, err
	}

	otgSrv := lwotg.New()
	telemSrv, err := lwotgtelem.New(context.Background(), ate.Id)
	if err != nil {
		return nil, fmt.Errorf("cannot set up telemetry server, %v", err)
	}

	fh, task, err := mpls.New()
	if err != nil {
		return nil, fmt.Errorf("cannot initialise MPLS flow handler, %v", err)
	}

	ipFH, ipTask, err := ip.New()
	if err != nil {
		return nil, fmt.Errorf("cannot initialise IP flow handler, %v", err)
	}

	otgSrv.AddFlowHandlers(fh)
	otgSrv.AddFlowHandlers(ipFH)
	telemSrv.AddTask(task)
	telemSrv.AddTask(ipTask)

	hintCh := make(chan lwotg.Hint, 100)
	otgSrv.SetHintChannel(hintCh)
	// otgSrv.SetProtocolHandler(gatewayPinger)
	telemSrv.SetHintChannel(context.Background(), hintCh)

	telemSrv.AddTask(arp.New(context.Background(), telemSrv.GetHints, time.Now().UnixNano))

	otgLis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", addr, otgPort))
	if err != nil {
		return nil, fmt.Errorf("cannot listen on port, err: %v", err)
	}

	gnmiLis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", addr, gnmiPort))
	if err != nil {
		return nil, fmt.Errorf("cannot listen on port err: %v", err)
	}

	otgS := grpc.NewServer(grpc.Creds(local.NewCredentials()))
	reflection.Register(otgS)
	otg.RegisterOpenapiServer(otgS, otgSrv)

	gnmiS := grpc.NewServer(grpc.Creds(local.NewCredentials()))
	reflection.Register(gnmiS)
	gpb.RegisterGNMIServer(gnmiS, telemSrv.GNMIServer)

	klog.Infof("OTG listening at %s", otgLis.Addr())
	klog.Infof("gNMI listening at %s", gnmiLis.Addr())
	go otgS.Serve(otgLis)
	go gnmiS.Serve(gnmiLis)

	m := &localMagna{
		AbstractATE: binding.AbstractATE{
			Dims: &binding.Dims{
				Name:          ate.Id,
				Vendor:        opb.Device_OPENCONFIG,
				HardwareModel: "MAGNA",
				Ports:         map[string]*binding.Port{},
			},
		},
		addr: addr,
	}
	for _, port := range ate.Ports {
		name := fmt.Sprintf("%s:%s", ate.Id, port.Id)
		_, err := portMgr.CreatePort(name)
		if err != nil {
			return nil, err
		}

		m.AbstractATE.Dims.Ports[port.Id] = &binding.Port{
			Name: name,
		}
	}

	return m, nil
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
	data, ok := <-p.rxQueue.Receive()
	if !ok {
		return nil, gopacket.CaptureInfo{}, io.EOF
	}
	packet := data.([]byte)
	return packet, gopacket.CaptureInfo{
		Timestamp:     time.Now(),
		CaptureLength: len(packet),
		Length:        len(packet),
	}, nil
}

func (p *chanPort) Close() {
	p.txQueue.Close()
}

func (p *chanPort) LinkType() layers.LinkType {
	return layers.LinkTypeEthernet
}

func (p *chanPort) SetBPFFilter(string) error {
	return nil
}

type portMgr struct {
	dut *dutPortMgr
	ate *atePortMgr
}

// dutPortMgr implements lemming's port creator interface.
// Note: this interface is different from magna, to avoid making the two project depending on eachother.
type dutPortMgr struct {
	ports map[string]*chanPort
	// TODO: lastCreatedPort is a hack because we don't currently model hardware ports correctly in lemming.
	lastCreatedPort *chanPort
}

func (lm *dutPortMgr) CreatePort(name string) (fwdcontext.Port, error) {
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

// atePortMgr implement magna's port creator interface.
// Note: this interface is different from lemming, to avoid making the two project depending on eachother.
type atePortMgr struct {
	ports map[string]*chanPort
}

func (lm *atePortMgr) CreatePort(name string) (common.Port, error) {
	if p, ok := lm.ports[name]; ok {
		return p, nil
	}
	q, err := queue.NewUnbounded(fmt.Sprintf("%s-tx", name))
	if err != nil {
		return nil, err
	}
	p := &chanPort{
		txQueue: q,
	}
	q.Run()
	fmt.Printf("create port %q", name)
	lm.ports[name] = p
	return p, nil
}

func (lm *portMgr) linkPorts(a, b string) error {
	aPort := lm.dut.ports[a]
	if aPort == nil {
		aPort = lm.ate.ports[a]
	}
	bPort := lm.dut.ports[b]
	if bPort == nil {
		bPort = lm.ate.ports[b]
	}

	if aPort == nil || bPort == nil {
		return fmt.Errorf("ports do not exist: a %v b %v", a, b)
	}
	aPort.rxQueue = bPort.txQueue
	bPort.rxQueue = aPort.txQueue

	return nil
}
