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

package info

import (
	"context"
	"crypto/tls"
	"encoding/hex"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

// New returns a new info command.
func New() *cobra.Command {
	cmd := &cobra.Command{
		Use: "info",
	}
	get := &cobra.Command{
		Use:   "get name",
		Short: "Get the description of a Lucius object.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return infoElement(cmd.Context(), args[0], fwdpb.InfoType_INFO_TYPE_ALL, "")
		},
	}
	lookup := &cobra.Command{
		Use:   "lookup table packetdata",
		Short: "Processes the packet in the table and returns the resulting actions. Note: may inject the packet into network.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return infoElement(cmd.Context(), args[0], fwdpb.InfoType_INFO_TYPE_LOOKUP, args[1])
		},
	}
	portInput := &cobra.Command{
		Use:   "port-input port packetdata",
		Short: "Processes the packet in the port's input actions. Note: may inject the packet into network.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return infoElement(cmd.Context(), args[0], fwdpb.InfoType_INFO_TYPE_PORT_INPUT, args[1])
		},
	}
	portOutput := &cobra.Command{
		Use:   "port-output port packetdata",
		Short: "Processes the packet using the port's output actions. Note: may inject the packet into network.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return infoElement(cmd.Context(), args[0], fwdpb.InfoType_INFO_TYPE_PORT_OUTPUT, args[1])
		},
	}

	list := &cobra.Command{
		Use:  "list",
		RunE: listFn,
	}
	cmd.AddCommand(list, get, lookup, portInput, portOutput)

	return cmd
}

func listFn(cmd *cobra.Command, _ []string) error {
	conn, err := dial()
	if err != nil {
		return fmt.Errorf("failed to dial dataplane: %v", err)
	}

	info := fwdpb.NewInfoClient(conn)

	resp, err := info.InfoList(cmd.Context(), &fwdpb.InfoListRequest{})
	if err != nil {
		return err
	}
	for _, n := range resp.GetNames() {
		fmt.Println(n)
	}

	return nil
}

func infoElement(ctx context.Context, name string, infoType fwdpb.InfoType, packet string) error {
	conn, err := dial()
	if err != nil {
		return fmt.Errorf("failed to dial dataplane: %v", err)
	}

	bytes, err := hex.DecodeString(packet)
	if err != nil {
		return err
	}

	info := fwdpb.NewInfoClient(conn)

	resp, err := info.InfoElement(ctx, &fwdpb.InfoElementRequest{
		Name:        name,
		Type:        infoType,
		Frame:       bytes,
		StartHeader: fwdpb.PacketHeaderId_PACKET_HEADER_ID_ETHERNET,
	})
	if err != nil {
		return err
	}
	fmt.Println(resp.GetContent())
	return nil
}

func dial() (*grpc.ClientConn, error) {
	insec, tlsSkipVerify := viper.GetBool("insecure"), viper.GetBool("tls_skip_verify")
	if insec && tlsSkipVerify {
		return nil, fmt.Errorf("both insecure and tls skip verify are set")
	}
	opts := grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
		InsecureSkipVerify: tlsSkipVerify, // nolint:gosec
	}))
	if insec {
		opts = grpc.WithTransportCredentials(insecure.NewCredentials())
	}

	return grpc.NewClient(viper.GetString("address"), opts)
}
