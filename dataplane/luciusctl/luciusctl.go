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

package luciusctl

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/openconfig/lemming/dataplane/luciusctl/info"
	"github.com/openconfig/lemming/dataplane/luciusctl/sai"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lucius",
		Short: "lucius queries info from a running instance of lucius.",
	}
	cmd.PersistentFlags().StringP("address", "a", "", "Address of dataplane server")
	cmd.PersistentFlags().Bool("insecure", false, "Use plaintext to dial server")
	cmd.PersistentFlags().Bool("tls_skip_verify", false, "Use TLS without cert verification")

	cobra.OnInitialize(func() { viper.BindPFlags(cmd.Flags()) })
	viper.BindPFlags(cmd.Flags())

	cmd.AddCommand(info.New(), sai.New())

	return cmd
}
