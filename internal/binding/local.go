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
	"sync"
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
	"github.com/openconfig/ondatra"
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

	saipb "github.com/openconfig/lemming/dataplane/proto/sai"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

// LocalBind is an Ondatra binding for in-process testbed. Only Lemming and Magna are supported.
type LocalBind struct {
	binding.Binding
	portMgr *PortMgr
	closers []func() error
}

type Option func(lb *LocalBind)

func WithOverridePortManager(mgr *PortMgr) Option {
	return func(lb *LocalBind) {
		lb.portMgr = mgr
	}
}

// Local is a local (in-process) binding for lemming and magna.
func Local(topoDir string, opts ...Option) func() (binding.Binding, error) {
	dir, _ := filepath.Abs(topoDir)
	testbedFile := filepath.Join(dir, "testbed.pb.txt")

	lb := &LocalBind{}
	for _, opt := range opts {
		opt(lb)
	}

	flag.Set("testbed", testbedFile)
	return func() (binding.Binding, error) {
		return lb, nil
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
	conn, err := grpc.NewClient(net.JoinHostPort(l.addr, fmt.Sprint(gnmiPort)), opts...)
	if err != nil {
		return nil, err
	}
	return gpb.NewGNMIClient(conn), nil
}

// DialGNOI returns a gNOI client for the dut.
func (l *localLemming) DialGNOI(ctx context.Context, opts ...grpc.DialOption) (gnoigo.Clients, error) {
	opts = append(opts, grpc.WithTransportCredentials(local.NewCredentials()))
	conn, err := grpc.NewClient(net.JoinHostPort(l.addr, fmt.Sprint(gnmiPort)), opts...)
	if err != nil {
		return nil, err
	}
	return gnoigo.NewClients(conn), nil
}

// DialGNSI returns a gNSI client for the dut.
func (l *localLemming) DialGNSI(ctx context.Context, opts ...grpc.DialOption) (binding.GNSIClients, error) {
	opts = append(opts, grpc.WithTransportCredentials(local.NewCredentials()))
	conn, err := grpc.NewClient(net.JoinHostPort(l.addr, fmt.Sprint(gnmiPort)), opts...)
	if err != nil {
		return nil, err
	}
	return &gnsiClient{conn: conn}, nil
}

// DialGRIBI returns a gRIBI client for the dut.
func (l *localLemming) DialGRIBI(ctx context.Context, opts ...grpc.DialOption) (grpb.GRIBIClient, error) {
	opts = append(opts, grpc.WithTransportCredentials(local.NewCredentials()))
	conn, err := grpc.NewClient(net.JoinHostPort(l.addr, fmt.Sprint(gribiPort)), opts...)
	if err != nil {
		return nil, err
	}
	return grpb.NewGRIBIClient(conn), nil
}

// DataplaneConn returns a gRPC conn for the dataplane
func (l *localLemming) DataplaneConn(ctx context.Context, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	opts = append(opts, grpc.WithTransportCredentials(local.NewCredentials()))
	return grpc.NewClient(net.JoinHostPort(l.addr, fmt.Sprint(dataplanePort)), opts...)
}

type localMagna struct {
	binding.AbstractATE
	addr string
}

// DialGNMI returns a gNMI client for the dut.
func (m *localMagna) DialGNMI(ctx context.Context, opts ...grpc.DialOption) (gpb.GNMIClient, error) {
	opts = append(opts, grpc.WithTransportCredentials(local.NewCredentials()))
	conn, err := grpc.NewClient(net.JoinHostPort(m.addr, fmt.Sprint(gnmiPort)), opts...)
	if err != nil {
		return nil, err
	}
	return gpb.NewGNMIClient(conn), nil
}

// DialOTG returns a OTH client for the dut.
func (m *localMagna) DialOTG(ctx context.Context, opts ...grpc.DialOption) (gosnappi.Api, error) {
	opts = append(opts, grpc.WithTransportCredentials(local.NewCredentials()))
	conn, err := grpc.NewClient(net.JoinHostPort(m.addr, fmt.Sprint(otgPort)), opts...)
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

	if lb.portMgr == nil {
		lb.portMgr = &PortMgr{}
	}
	if lb.portMgr.ports == nil {
		lb.portMgr.ports = map[string]*ChanPort{}
	}
	if lb.portMgr.dutLaneToPort == nil {
		lb.portMgr.dutLaneToPort = map[string]map[string]string{}
	}

	if err := lb.portMgr.createPorts(tb); err != nil {
		return nil, err
	}
	for _, l := range tb.Links {
		if err := lb.portMgr.linkPorts(l.A, l.B); err != nil {
			return nil, err
		}
	}

	intf.OverrideAccessor(&accessor{})
	common.OverrideHandleCreator(lb.portMgr)

	for _, ate := range tb.Ates {
		addr, closeFn, err := findAvailableLoopbackIP()
		if err != nil {
			return nil, err
		}
		lb.closers = append(lb.closers, closeFn)

		magna, err := lb.createATE(ctx, ate, addr, lb.portMgr)
		if err != nil {
			return nil, err
		}
		resv.ATEs[ate.Id] = magna
	}

	for _, dut := range tb.Duts {
		addr, closeFn, err := findAvailableLoopbackIP()
		if err != nil {
			return nil, err
		}
		lb.closers = append(lb.closers, closeFn)

		lemming, err := lb.createDUT(ctx, dut, addr, lb.portMgr)
		if err != nil {
			return nil, err
		}
		resv.DUTs[dut.Id] = lemming
	}

	return &resv, nil
}

func (lb *LocalBind) createDUT(ctx context.Context, dut *opb.Device, addr string, portMgr *PortMgr) (*localLemming, error) {
	dutID := uuid.New().String()

	l, err := lemming.New(dut.Id, fmt.Sprintf("unix:/tmp/zserv-test%s.api", dutID),
		lemming.WithBGPPort(bgpPort),
		lemming.WithGNMIAddr(net.JoinHostPort(addr, fmt.Sprint(gnmiPort))),
		lemming.WithGRIBIAddr(net.JoinHostPort(addr, fmt.Sprint(gribiPort))),
		lemming.WithTransportCreds(local.NewCredentials()),
		lemming.WithDataplane(true),
		lemming.WithSysribAddr(fmt.Sprintf("/tmp/sysrib-%s.api", dutID)),
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
		addr:  addr,
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
	fwdCtx.FakePortManager = portMgr.dutManager(dut.GetId())
	dplaneConn, err := boundLemming.l.Dataplane().Conn()
	if err != nil {
		return nil, err
	}
	pc := saipb.NewPortClient(dplaneConn)
	// For each port on the topology textproto.
	// Create a saipb port and accociate the ondatra dut ID and port ID with the saipb OID.
	for i, port := range dut.Ports {
		resp, err := pc.CreatePort(ctx, &saipb.CreatePortRequest{
			HwLaneList: []uint32{uint32(i)},
		})
		if err != nil {
			return nil, err
		}
		boundLemming.AbstractDUT.Dims.Ports[port.Id] = &binding.Port{
			Name: fmt.Sprint(resp.Oid),
		}
	}
	return boundLemming, nil
}

// TODO: this should probably be a library in magna.
func (lb *LocalBind) createATE(_ context.Context, ate *opb.Device, addr string, portMgr *PortMgr) (*localMagna, error) {
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
		_, err := portMgr.CreateHandle(name)
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
	for _, closer := range lb.closers {
		closer()
	}

	for name, port := range lb.portMgr.ports {
		fmt.Printf("port %s, tx enqueue %d, tx dequeue %d\n", name, port.TXQueue.EnqueueCount(), port.TXQueue.DequeueCount())
		fmt.Printf("port %s, rx enqueue %d, rx dequeue %d\n", name, port.RXQueue.EnqueueCount(), port.RXQueue.DequeueCount())
	}
	return nil
}

// findAvailableLoopbackIP finds an unused loopback IP by attempting to open a tcp socket on the gNMI port.
func findAvailableLoopbackIP() (string, func() error, error) {
	var addr string
	var closer func() error
	for i := byte(2); i < 254; i++ {
		addr = netip.AddrFrom4([4]byte{127, 0, 0, i}).String()
		list, err := net.Listen("tcp", net.JoinHostPort(addr, "54591")) // Listen on a random port.
		if err == nil {
			closer = list.Close
			break
		}
	}
	if addr == "" {
		return "", nil, fmt.Errorf("failed to find available ip")
	}
	return addr, closer, nil
}

func newPort(name string) (*ChanPort, error) {
	tx, err := queue.NewUnbounded(name + "_tx")
	if err != nil {
		return nil, err
	}
	rx, err := queue.NewUnbounded(name + "_rx")
	if err != nil {
		return nil, err
	}
	p := &ChanPort{
		TXQueue: tx,
		RXQueue: rx,
	}
	tx.Run()
	rx.Run()
	p.run()
	return p, nil
}

type PortMgr struct {
	ports         map[string]*ChanPort
	dutLaneToPort map[string]map[string]string
}

func (pm *PortMgr) GetPort(p *ondatra.Port) *ChanPort {
	return pm.ports[fmt.Sprintf("%s:%s", p.Device().ID(), p.ID())]
}

func (pm *PortMgr) createPorts(tb *opb.Testbed) error {
	for _, dut := range tb.Duts {
		pm.dutLaneToPort[dut.GetId()] = make(map[string]string)
		for i, port := range dut.Ports {
			name := fmt.Sprintf("%s:%s", dut.GetId(), port.GetId())
			p, err := newPort(name)
			if err != nil {
				return err
			}
			pm.dutLaneToPort[dut.GetId()][fmt.Sprint(i)] = name
			pm.ports[name] = p
		}
	}
	for _, ate := range tb.Ates {
		for _, port := range ate.Ports {
			name := fmt.Sprintf("%s:%s", ate.GetId(), port.GetId())
			p, err := newPort(name)
			if err != nil {
				return err
			}
			pm.ports[name] = p
		}
	}
	return nil
}

// CreateHandle implements  magna's API for creating handles.
func (pm *PortMgr) CreateHandle(name string) (common.Port, error) {
	port, ok := pm.ports[name]
	if !ok {
		return nil, fmt.Errorf("port %v not found", name)
	}
	return port.newHandle(), nil
}

func (pm *PortMgr) dutManager(dutID string) *dutManager {
	return &dutManager{
		mgr:   pm,
		dutID: dutID,
	}
}

func (pm *PortMgr) linkPorts(a, b string) error {
	aPort := pm.ports[a]
	bPort := pm.ports[b]

	if aPort == nil || bPort == nil {
		return fmt.Errorf("ports do not exist: a %v b %v", a, b)
	}
	aPort.mu.Lock()
	bPort.mu.Lock()
	aPort.RXQueue.Close()
	bPort.RXQueue.Close()
	aPort.RXQueue = bPort.TXQueue
	bPort.RXQueue = aPort.TXQueue
	aPort.mu.Unlock()
	bPort.mu.Unlock()

	return nil
}

// dutManager handles mapping from the SAI API hardware lane to Ondatra port ID.
type dutManager struct {
	mgr   *PortMgr
	dutID string
}

// CreatePort implements lemming's API for creating port.
func (dm *dutManager) CreatePort(name string) (fwdcontext.Port, error) {
	port, ok := dm.mgr.ports[dm.mgr.dutLaneToPort[dm.dutID][name]]
	if !ok {
		return nil, fmt.Errorf("port %v not found", name)
	}
	return port.newHandle(), nil
}

// ChanPort is a fake port implemented using channels.
type ChanPort struct {
	mu sync.RWMutex

	TXQueue *queue.Queue
	RXQueue *queue.Queue
	handles []*portHandle
}

func (p *ChanPort) newHandle() *portHandle {
	p.mu.Lock()
	defer p.mu.Unlock()

	h := &portHandle{
		rx: make(chan []byte, 1024),
		tx: p.TXQueue,
	}
	p.handles = append(p.handles, h)
	i := len(p.handles) - 1
	h.closeFn = func() {
		p.mu.Lock()
		defer p.mu.Unlock()
		close(h.rx)
		p.handles = append(p.handles[:i], p.handles[i+1:]...)
	}
	return h
}

func (p *ChanPort) run() {
	go func() {
		for {
			p.mu.RLock()
			if p.RXQueue == nil {
				p.mu.RUnlock()
				continue
			}
			q := p.RXQueue.Receive()
			p.mu.RUnlock()
			packet, ok := <-q
			if !ok {
				time.Sleep(time.Millisecond)
				continue
			}
			for _, h := range p.handles {
				h.rx <- packet.([]byte)
			}
		}
	}()
}

// portHandle provides an API for reading and writing to a chanPort.
// Each handle gets a copy of all packets received and writes to a common queue.
type portHandle struct {
	rx      chan []byte
	tx      *queue.Queue
	closeFn func()
}

func (ph *portHandle) ReadPacketData() (data []byte, ci gopacket.CaptureInfo, err error) {
	packet, ok := <-ph.rx
	if !ok {
		return nil, gopacket.CaptureInfo{}, io.EOF
	}
	return packet, gopacket.CaptureInfo{
		Timestamp:     time.Now(),
		CaptureLength: len(packet),
		Length:        len(packet),
	}, nil
}

func (ph *portHandle) WritePacketData(data []byte) error {
	return ph.tx.Write(data)
}

func (ph *portHandle) Close() {
	ph.closeFn()
}

func (ph *portHandle) LinkType() layers.LinkType {
	return layers.LinkTypeEthernet
}

func (ph *portHandle) SetBPFFilter(string) error {
	return nil
}
