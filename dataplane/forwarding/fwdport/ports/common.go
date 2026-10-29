// Copyright 2022 Google LLC
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

//go:build linux

package ports

import (
	"net"
	"sync"

	log "github.com/golang/glog"
	"github.com/vishvananda/netlink"

	"github.com/openconfig/lemming/dataplane/forwarding/fwdport"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdcontext"
	"github.com/openconfig/lemming/dataplane/kernel"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

type stateIfaceMgr interface {
	LinkByName(name string) (netlink.Link, error)
	SetState(name string, up bool) error
}

func stateFromAttrs(attrs *netlink.LinkAttrs) (fwdpb.PortState, fwdpb.PortState) {
	var adminState fwdpb.PortState
	if attrs.Flags&net.FlagUp != 0 {
		adminState = fwdpb.PortState_PORT_STATE_ENABLED_UP
	} else {
		adminState = fwdpb.PortState_PORT_STATE_DISABLED_DOWN
	}
	var operState fwdpb.PortState
	switch attrs.OperState {
	case netlink.OperDown, netlink.OperLowerLayerDown:
		operState = fwdpb.PortState_PORT_STATE_DISABLED_DOWN
	case netlink.OperUp, netlink.OperUnknown:
		operState = fwdpb.PortState_PORT_STATE_ENABLED_UP
	}
	return adminState, operState
}

func getPortState(name string, ifMgr stateIfaceMgr) (fwdpb.PortState, fwdpb.PortState, error) {
	iface, err := ifMgr.LinkByName(name)
	if err != nil {
		return fwdpb.PortState_PORT_STATE_UNSPECIFIED, fwdpb.PortState_PORT_STATE_UNSPECIFIED, err
	}
	admin, oper := stateFromAttrs(iface.Attrs())
	return admin, oper, nil
}

// getAndSetState returns the state of the port and optionally sets the state to the new value.
// Note: the reply doesn't contain the updated oper-status (if applicable).
func getAndSetState(name string, ifMgr stateIfaceMgr, pi *fwdpb.PortInfo) (*fwdpb.PortStateReply, error) {
	adminState, operState, err := getPortState(name, ifMgr)
	if err != nil {
		return nil, err
	}
	reply := &fwdpb.PortStateReply{
		Status: &fwdpb.PortInfo{
			OperStatus:  operState,
			AdminStatus: adminState,
		},
	}
	if pi == nil {
		log.V(1).Infof("dataplane read port %q: admin state %v, oper state %v", name, adminState.String(), operState.String())
		return reply, nil
	}
	log.V(1).Infof("dataplane write port %q: admin state %v, oper state %v, setting to %v", name, adminState.String(), operState.String(), pi.AdminStatus.String())
	if adminState == pi.GetAdminStatus() {
		return reply, nil
	}
	switch pi.AdminStatus {
	case fwdpb.PortState_PORT_STATE_DISABLED_DOWN:
		err = ifMgr.SetState(name, false)
	case fwdpb.PortState_PORT_STATE_ENABLED_UP:
		err = ifMgr.SetState(name, true)
	default:
		return reply, nil
	}
	if err != nil {
		return nil, err
	}
	reply.Status.AdminStatus = pi.AdminStatus
	return reply, nil
}

type linkSubscriber interface {
	LinkSubscribe(ch chan<- netlink.LinkUpdate, done <-chan struct{}) error
}

type portTarget struct {
	port fwdport.Port
	ctx  *fwdcontext.Context
}

type linkStateWatcher struct {
	mu       sync.Mutex
	ifaceMgr linkSubscriber
	ports    map[string]portTarget
	doneCh   chan struct{}
}

var defaultLinkWatcher = &linkStateWatcher{
	ifaceMgr: &kernel.Interfaces{},
	ports:    make(map[string]portTarget),
}

func (w *linkStateWatcher) watch(devName string, port fwdport.Port, ctx *fwdcontext.Context) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if len(w.ports) == 0 {
		updCh := make(chan netlink.LinkUpdate)
		doneCh := make(chan struct{})
		if err := w.ifaceMgr.LinkSubscribe(updCh, doneCh); err != nil {
			return err
		}
		w.doneCh = doneCh
		go w.run(updCh, doneCh)
	}
	w.ports[devName] = portTarget{port: port, ctx: ctx}
	return nil
}

func (w *linkStateWatcher) unwatch(devName string) {
	w.mu.Lock()
	defer w.mu.Unlock()

	delete(w.ports, devName)
	if len(w.ports) == 0 && w.doneCh != nil {
		close(w.doneCh)
		w.doneCh = nil
	}
}

func (w *linkStateWatcher) run(updCh <-chan netlink.LinkUpdate, doneCh <-chan struct{}) {
	for {
		select {
		case <-doneCh:
			log.Warningf("shared link state watch stopping")
			return
		case upd, ok := <-updCh:
			if !ok {
				log.Infof("shared link update chan closed")
				return
			}
			if upd.Attrs() == nil {
				continue
			}
			devName := upd.Attrs().Name
			w.mu.Lock()
			target, ok := w.ports[devName]
			w.mu.Unlock()
			if !ok {
				continue
			}
			admin, oper := stateFromAttrs(upd.Attrs())
			log.V(1).Infof("dataplane receive link update: port %q, admin %v, oper %v", devName, admin.String(), oper.String())
			if err := target.ctx.Notify(&fwdpb.EventDesc{
				Event: fwdpb.Event_EVENT_PORT,
				Desc: &fwdpb.EventDesc_Port{
					Port: &fwdpb.PortEventDesc{
						Context: &fwdpb.ContextId{Id: target.ctx.ID},
						PortId:  &fwdpb.PortId{ObjectId: &fwdpb.ObjectId{Id: string(target.port.ID())}},
						PortInfo: &fwdpb.PortInfo{
							AdminStatus: admin,
							OperStatus:  oper,
						},
					},
				},
			}); err != nil {
				log.Errorf("failed to notify port state update: %v", err)
			}
		}
	}
}
