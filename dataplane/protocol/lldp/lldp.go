// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lldp

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/openconfig/ygnmi/ygnmi"

	"github.com/openconfig/lemming/gnmi/gnmiclient"
	"github.com/openconfig/lemming/gnmi/oc"

	log "github.com/golang/glog"
	"github.com/mdlayher/lldp"

	"github.com/openconfig/lemming/gnmi/oc/ocpath"

	"github.com/openconfig/lemming/dataplane/proto/packetio"
)

// Daemon is the implementation of the LLDP protocol.
type Daemon struct {
	enabled bool
	mu      sync.Mutex             // guard the ports map
	ports   map[string]*portDaemon // key=hostif_id
}

// New returns the main procotol handler.
func New() *Daemon {
	return &Daemon{
		ports: map[string]*portDaemon{},
	}
}

// getOrCreatePortDaemon returns the state interface.
func (d *Daemon) getOrCreatePortDaemon(name string) *portDaemon {
	if _, ok := d.ports[name]; !ok {
		d.ports[name] = newPortDaemon(name)
	}
	return d.ports[name]
}

// Reconcile reconciles LLDP for all ports.
func (d *Daemon) Reconcile(ctx context.Context, intent *oc.Root, c *ygnmi.Client) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	sb := &ygnmi.SetBatch{}
	// Protocol enable.
	if d.enabled != intent.Lldp.GetEnabled() {
		d.enabled = intent.Lldp.GetEnabled()
		gnmiclient.BatchUpdate(sb, ocpath.Root().Lldp().Enabled().State(), d.enabled)
	}
	// Per-interface enable.
	for _, intf := range intent.GetLldp().Interface {
		pd := d.getOrCreatePortDaemon(intf.GetName())
		newState := intent.Lldp.GetEnabled() && intf.GetEnabled()
		if pd.enabled != newState {
			pd.SetEnabled(newState)
			gnmiclient.BatchUpdate(sb, ocpath.Root().Lldp().Interface(intf.GetName()).Enabled().State(), newState)
		}
	}
	if _, err := sb.Set(ctx, c); err != nil {
		return fmt.Errorf("failed to update LLDP enable state: %v", err)
	}
	return nil
}

// AddPort registers a port.
func (d *Daemon) AddPort(pn string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	_, ok := d.ports[pn]
	if ok {
		return fmt.Errorf("port %q exists already", pn)
	}
	d.ports[pn] = newPortDaemon(pn)
	return nil
}

// RemovePort deregisters a port.
func (d *Daemon) RemovePort(pn string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	_, ok := d.ports[pn]
	if !ok {
		return fmt.Errorf("port %q not found", pn)
	}
	delete(d.ports, pn)
	return nil
}

// Matched returns true if this packet can be handled by this protocol.
func (d *Daemon) Matched(po *packetio.PacketOut) bool {
	return IsLldp(po)
}

// Process dispatches the packet to the corresponding port handler.
func (d *Daemon) Process(p *packetio.Packet) error {
	pd, ok := d.ports[fmt.Sprintf("%d", p.HostPort)]
	if !ok {
		return fmt.Errorf("port %q not found", p.HostPort)
	}
	return pd.Process(p)
}

// Start starts the procotol handler.
func (d *Daemon) Start() error {
	return nil
}

// Stop stops the procotol handler.
func (d *Daemon) Stop() {
	for _, p := range d.ports {
		p.Stop()
	}
	d.ports = nil
}

// lldpInfo contains LLDP protocol information.
type lldpInfo struct {
	HostIfId       string
	SysName        string
	SysDesc        string
	PortName       string
	PortDesc       string
	RemoteSysName  string
	RemoteSysDesc  string
	RemotePortName string
	RemotePortDesc string
	HardwareAddr   string
	Interval       time.Duration
}

// Frames returns the LLDP packet.
func (l *lldpInfo) Frame() ([]byte, error) {
	lf := lldp.Frame{
		ChassisID: &lldp.ChassisID{
			Subtype: lldp.ChassisIDSubtypeMACAddress,
			ID:      []byte(l.HardwareAddr),
		},
		PortID: &lldp.PortID{
			Subtype: lldp.PortIDSubtypeInterfaceName,
			ID:      []byte(l.PortName),
		},
		TTL: 2 * l.Interval,
		Optional: []*lldp.TLV{
			{
				Type:   lldp.TLVTypePortDescription,
				Value:  []byte(l.PortDesc),
				Length: uint16(len(l.PortDesc)),
			},
			{
				Type:   lldp.TLVTypeSystemName,
				Value:  []byte(l.SysName),
				Length: uint16(len(l.SysName)),
			},
			{
				Type:   lldp.TLVTypeSystemDescription,
				Value:  []byte(l.SysDesc),
				Length: uint16(len(l.SysDesc)),
			},
		},
	}
	return lf.MarshalBinary()
}

// portDaemon contains the required information for LLDP and processes the LLDP frames for a given hostif.
type portDaemon struct {
	Name      string
	info      *lldpInfo
	enabled   bool
	enabledCh chan bool // indicates the state of this daemon.
	Interval  time.Duration
	doneCh    chan struct{}
	InCh      chan *packetio.Packet
	ErrSendCh chan error
	ErrRecvCh chan error
}

// newPortDaemon creates a port daemon to send/recv LLDP protocol.
func newPortDaemon(hostif string) *portDaemon {
	pd := &portDaemon{
		Name:      hostif,
		doneCh:    make(chan struct{}),
		InCh:      make(chan *packetio.Packet),
		ErrSendCh: make(chan error),
		ErrRecvCh: make(chan error),
	}

	go func() {
		log.Infof("Start LLDP sender.")
		enabled := false
		for {
			select {
			case <-pd.doneCh:
				log.Infof("Stop sending LLDP frame")
				return
			case e := <-pd.enabledCh:
				enabled = e
			default:
				if enabled {
					_, err := pd.info.Frame()
					if err != nil {
						pd.ErrSendCh <- fmt.Errorf("failed to create LLDP frame: %v", err)
						continue
					}
					// TODO: Send the packe to the hostif.
					log.Infof("Write LLDP frame to port: %+v", hostif)
					time.Sleep(pd.Interval * time.Second)
				} else {
					log.Infof("LLDP frame sending paused.")
					enabled = <-pd.enabledCh
				}
			}
		}
	}()

	go func() {
		log.Infof("Start LLDP receiver.")
		enabled := false
		for {
			select {
			case <-pd.doneCh:
				log.Infof("Stop processing LLDP frame")
				return
			case e := <-pd.enabledCh:
				enabled = e
			default:
				if enabled {
					f := <-pd.InCh
					log.Infof("Got LLDP Frame: %+v", f.String())
					pkt := gopacket.NewPacket(f.GetFrame(), layers.LayerTypeEthernet, gopacket.Default)
					for _, layer := range pkt.Layers() {
						if layer.LayerType() != layers.LayerTypeLinkLayerDiscoveryInfo {
							log.Infof("Skipped layer %v", layer.LayerType().String())
							continue
						}
						info, ok := layer.(*layers.LinkLayerDiscoveryInfo)
						if !ok {
							pd.ErrRecvCh <- fmt.Errorf("packet is not LinkLayerDiscoveryInfo: %+v", layer)
							continue
						}
						pd.info.RemoteSysName = info.SysName
						pd.info.RemoteSysDesc = info.SysDescription
					}
				} else {
					log.Infof("LLDP frame receiving paused.")
					enabled = <-pd.enabledCh
				}
			}
		}
	}()
	return pd
}

// Process handles the packet from the hostif and update the remote information.
func (d *portDaemon) Process(p *packetio.Packet) error {
	if d.InCh == nil {
		return fmt.Errorf("failed to inject packet")
	}
	d.InCh <- p
	err := <-d.ErrRecvCh
	return err
}

// SetEnabled starts/stops sending and receiving goroutines for this port.
func (d *portDaemon) SetEnabled(b bool) error {
	d.enabledCh <- b
	d.enabled = b
	return nil
}

// Stop stops the port daemon.
func (d *portDaemon) Stop() {
	d.doneCh <- struct{}{}
}

// IsLldp returns whether the packet is a LLDP frame.
func IsLldp(po *packetio.PacketOut) bool {
	ethPkt := gopacket.NewPacket(po.GetPacket().GetFrame(), layers.LayerTypeEthernet, gopacket.Default)
	ethLayer := ethPkt.Layer(layers.LayerTypeEthernet)
	if ethLayer == nil {
		return false
	}
	return ethLayer.(*layers.Ethernet).EthernetType == layers.EthernetType(layers.EthernetTypeLinkLayerDiscovery.LayerType())
}
