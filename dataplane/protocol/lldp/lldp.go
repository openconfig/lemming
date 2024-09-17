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
	"net"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/openconfig/ygnmi/ygnmi"

	"github.com/openconfig/lemming/gnmi/gnmiclient"
	"github.com/openconfig/lemming/gnmi/oc"

	log "github.com/golang/glog"

	"github.com/openconfig/lemming/gnmi/oc/ocpath"

	"github.com/openconfig/lemming/dataplane/proto/packetio"
)

// Daemon is the implementation of the LLDP protocol.
type Daemon struct {
	enabled     bool                   // whether LLDP is enabled globally
	portEnabled map[string]bool        // contains the enabled ports
	portDaemons map[string]*portDaemon // tracks the active port daemons
}

// New returns the main procotol handler.
func New() *Daemon {
	return &Daemon{
		portEnabled: map[string]bool{},
		portDaemons: map[string]*portDaemon{},
	}
}

// Start starts the procotol handler.
func (d *Daemon) Start() {
	if d.portDaemons == nil {
		d.portDaemons = map[string]*portDaemon{}
	}
	d.enabled = true
}

// Stop stops the procotol handler by stopping all port daemons.
func (d *Daemon) Stop() {
	for _, p := range d.portDaemons {
		p.Stop()
	}
	d.enabled = false
	d.portDaemons = nil
}

// changePortState tries to change the port state and returns whether the state is acutally changed.
func (d *Daemon) changePortState(intf *oc.Lldp_Interface) bool {
	// Update the PD state.
	wantActive := d.enabled && intf.GetEnabled()
	pd, currActive := d.portDaemons[intf.GetName()]
	switch {
	case !currActive && wantActive:
		d.portDaemons[intf.GetName()] = newPortDaemon(intf.GetName())
	case currActive && !wantActive:
		pd.Stop()
		delete(d.portDaemons, intf.GetName())
	}

	// Update the enabled state.
	state, ok := d.portEnabled[intf.GetName()]
	if !ok {
		d.portEnabled[intf.GetName()] = false
	}
	if state != intf.GetEnabled() {
		d.portEnabled[intf.GetName()] = intf.GetEnabled()
		return true
	}
	return false
}

func (d *Daemon) reconcileLldpEnabled(sb *ygnmi.SetBatch, intent *oc.Root, c *ygnmi.Client) error {
	if wantEnabled := intent.Lldp.GetEnabled(); d.enabled != wantEnabled {
		d.enabled = wantEnabled
		gnmiclient.BatchUpdate(sb, ocpath.Root().Lldp().Enabled().State(), d.enabled)
	}
	return nil
}

func (d *Daemon) reconcileLldpInterfaceEnabled(sb *ygnmi.SetBatch, intent *oc.Root, c *ygnmi.Client) error {
	for _, intf := range intent.GetLldp().Interface {
		if changed := d.changePortState(intf); changed {
			gnmiclient.BatchUpdate(sb, ocpath.Root().Lldp().Interface(intf.GetName()).Enabled().State(), intf.GetEnabled())
		}
	}
	return nil
}

// Reconcile reconciles LLDP for all ports.
func (d *Daemon) Reconcile(ctx context.Context, intent *oc.Root, c *ygnmi.Client) error {
	sb := &ygnmi.SetBatch{}
	for _, f := range []func(*ygnmi.SetBatch, *oc.Root, *ygnmi.Client) error{
		d.reconcileLldpEnabled,
		d.reconcileLldpInterfaceEnabled,
	} {
		f(sb, intent, c)
	}
	if _, err := sb.Set(ctx, c); err != nil {
		return fmt.Errorf("failed to update LLDP enable state: %v", err)
	}
	return nil
}

// Matched returns true if this packet can be handled by this protocol.
func (d *Daemon) Matched(po *packetio.PacketOut) bool {
	return IsLldp(po)
}

// Process dispatches the packet to the corresponding port handler.
func (d *Daemon) Process(p *packetio.Packet) error {
	pd, ok := d.portDaemons[fmt.Sprintf("%d", p.HostPort)]
	if !ok {
		return fmt.Errorf("port %q not found", p.HostPort)
	}
	return pd.Process(p)
}

// portDaemon contains the required information for LLDP and processes the LLDP frames for a given hostif.
type portDaemon struct {
	Name      string
	info      *lldpInfo
	Interval  time.Duration
	doneCh    chan struct{}
	inCh      chan *packetio.Packet
	errRecvCh chan error
}

// newPortDaemon creates a port daemon to send/recv LLDP protocol.
func newPortDaemon(hostif string) *portDaemon {
	d := &portDaemon{
		Name:      hostif,
		doneCh:    make(chan struct{}),
		inCh:      make(chan *packetio.Packet),
		errRecvCh: make(chan error),
	}
	go func() {
		log.Infof("Start LLDP sender.")
		for {
			select {
			case <-d.doneCh:
				log.Infof("Stop sending LLDP frame")
				return
			default:
				_, err := d.info.Frame()
				if err != nil {
					log.Errorf("failed to create LLDP frame: %v", err)
					continue
				}
				// TODO: Send the packe to the hostif.
				log.Infof("Write LLDP frame to port: %+v", d.Name)
				time.Sleep(d.Interval)
			}
		}
	}()

	go func() {
		log.Infof("Start LLDP receiver.")
		for {
			select {
			case <-d.doneCh:
				log.Infof("Stop processing LLDP frame")
				return
			default:
				f := <-d.inCh
				log.Infof("Got LLDP Frame: %+v", f.String())
				pkt := gopacket.NewPacket(f.GetFrame(), layers.LayerTypeEthernet, gopacket.Default)
				for _, layer := range pkt.Layers() {
					if layer.LayerType() != layers.LayerTypeLinkLayerDiscoveryInfo {
						log.Infof("Skipped layer %v", layer.LayerType().String())
						continue
					}
					info, ok := layer.(*layers.LinkLayerDiscoveryInfo)
					if !ok {
						d.errRecvCh <- fmt.Errorf("packet is not LinkLayerDiscoveryInfo: %+v", layer)
						continue
					}
					d.info.RemoteSysName = info.SysName
					d.info.RemoteSysDesc = info.SysDescription
				}
			}
		}
	}()
	return d
}

// Process handles the packet from the hostif and update the remote information.
func (d *portDaemon) Process(p *packetio.Packet) error {
	if d.inCh == nil {
		return fmt.Errorf("failed to inject packet")
	}
	d.inCh <- p
	err := <-d.errRecvCh
	return err
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
	HardwareAddr   []byte
	Interval       uint16
}

// Frames returns the LLDP packet.
func (l *lldpInfo) Frame() ([]byte, error) {
	dstMac, err := net.ParseMAC("01:80:C2:00:00:0E")
	if err != nil {
		return nil, err
	}
	pktEth := &layers.Ethernet{
		SrcMAC:       net.HardwareAddr(l.HardwareAddr),
		DstMAC:       dstMac,
		EthernetType: layers.EthernetTypeLinkLayerDiscovery,
	}
	pktLldp := &layers.LinkLayerDiscovery{
		ChassisID: layers.LLDPChassisID{
			Subtype: layers.LLDPChassisIDSubTypeMACAddr,
			ID:      l.HardwareAddr,
		},
		PortID: layers.LLDPPortID{
			Subtype: layers.LLDPPortIDSubtypeIfaceName,
			ID:      []byte(l.PortName),
		},
		TTL: 2 * l.Interval,
		Values: []layers.LinkLayerDiscoveryValue{
			{
				Type:   layers.LLDPTLVPortDescription,
				Value:  []byte(l.PortDesc),
				Length: uint16(len(l.PortDesc)),
			}, {
				Type:   layers.LLDPTLVSysName,
				Value:  []byte(l.SysName),
				Length: uint16(len(l.SysName)),
			}, {
				Type:   layers.LLDPTLVSysDescription,
				Value:  []byte(l.SysDesc),
				Length: uint16(len(l.SysDesc)),
			},
		},
	}
	buf := gopacket.NewSerializeBuffer()
	if err := gopacket.SerializeLayers(buf, gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}, pktEth, pktLldp); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
