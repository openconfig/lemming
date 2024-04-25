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

package sai

import (
	"crypto/tls"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use: "sai",
	}
	get := &cobra.Command{
		Use:     getAttr,
		Aliases: []string{"get-attr"},
		RunE: func(cmd *cobra.Command, args []string) error {
			conn, err := dial()
			if err != nil {
				return err
			}
			rpcs, err := discoverRPCs(conn)
			if err != nil {
				return err
			}
			m := rpcs[getAttr][args[0]]
			req := dynamicpb.NewMessage(m.input)

			if m.hasOID {
				oid, err := strconv.ParseUint(args[1], 10, 64)
				if err != nil {
					return err
				}
				req.Set(req.Descriptor().Fields().ByNumber(1), protoreflect.ValueOfUint64(oid))
			} else {
				if err := prototext.Unmarshal([]byte(args[1]), req.Mutable(req.Descriptor().Fields().ByNumber(1)).Message().Interface()); err != nil {
					return err
				}
			}

			for _, arg := range args[2:] {
				req.Mutable(req.Descriptor().Fields().ByNumber(2)).List().Append(protoreflect.ValueOf(m.args[arg].enumVal))
			}

			resp := dynamicpb.NewMessage(m.output)
			if err := conn.Invoke(cmd.Context(), m.fullName, req, resp); err != nil {
				return err
			}
			fmt.Println(resp)

			return nil
		},
		ValidArgsFunction: func(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			conn, err := dial()
			if err != nil {
				return nil, cobra.ShellCompDirectiveError
			}
			msg, err := discoverRPCs(conn)
			if err != nil {
				return nil, cobra.ShellCompDirectiveError
			}

			switch len(args) {
			case 0:
				return keyPrefixMatch(msg[getAttr], toComplete), cobra.ShellCompDirectiveNoFileComp
			case 1:
				return nil, cobra.ShellCompDirectiveNoFileComp
			default:
				return keyPrefixMatch(msg[getAttr][args[0]].args, toComplete), cobra.ShellCompDirectiveNoFileComp
			}
		},
	}

	cmd.AddCommand(get)
	return cmd
}

const (
	saiPrefix = "lemming.dataplane.sai"
	getAttr   = "get-attribute"
)

func keyPrefixMatch[T any](method map[string]T, prefix string) []string {
	matches := []string{}
	for m := range method {
		if strings.HasPrefix(strings.ToLower(m), prefix) {
			matches = append(matches, m)
		}
	}
	slices.Sort(matches)
	return matches
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

	return grpc.Dial(viper.GetString("address"), opts)
}
