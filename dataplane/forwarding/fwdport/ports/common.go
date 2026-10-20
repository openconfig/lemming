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

	log "github.com/golang/glog"
	"github.com/vishvananda/netlink"

	"github.com/openconfig/lemming/dataplane/forwarding/fwdport"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdcontext"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

type stateIfaceMgr interface {
	LinkByName(name string) (netlink.Link, error)
	SetState(name string, up bool) error
}

func stateFromAttrs(attrs *netlink.LinkAttrs) (fwdpb.PortState, fwdpb.PortState) {
	adminState := fwdpb.PortState_PORT_STATE_DISABLED_DOWN
	if attrs.Flags&net.FlagUp != 0 {
		adminState = fwdpb.PortState_PORT_STATE_ENABLED_UP
	}

	var state fwdpb.PortState
	switch attrs.OperState {
	case netlink.OperDown, netlink.OperLowerLayerDown, netlink.OperUnknown:
		state = fwdpb.PortState_PORT_STATE_DISABLED_DOWN
	case netlink.OperUp:
		state = fwdpb.PortState_PORT_STATE_ENABLED_UP
	}
	return adminState, state
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
	if err == nil {
		reply.Status.AdminStatus = pi.AdminStatus
	}
	return reply, err
}

func startStateWatch(updCh chan netlink.LinkUpdate, doneCh chan struct{}, devName string, port fwdport.Port, ctx *fwdcontext.Context) {
	go func() {
		for {
			select {
			case <-doneCh:
				log.Warningf("state watch stopping for dev: %v", devName)
				return
			case upd, ok := <-updCh:
				if !ok {
					log.Infof("update chan close for dev: %v", devName)
					return
				}
				if upd.Link == nil || upd.Attrs() == nil || upd.Attrs().Name != devName {
					continue
				}
				admin, oper := stateFromAttrs(upd.Attrs())
				log.V(1).Infof("dataplane receive link update: port %q, admin %v, oper %v", devName, admin.String(), oper.String())
				ctx.Notify(&fwdpb.EventDesc{
					Event: fwdpb.Event_EVENT_PORT,
					Desc: &fwdpb.EventDesc_Port{
						Port: &fwdpb.PortEventDesc{
							Context: &fwdpb.ContextId{Id: ctx.ID},
							PortId:  &fwdpb.PortId{ObjectId: &fwdpb.ObjectId{Id: string(port.ID())}},
							PortInfo: &fwdpb.PortInfo{
								AdminStatus: admin,
								OperStatus:  oper,
							},
						},
					},
				})
			}
		}
	}()
}
