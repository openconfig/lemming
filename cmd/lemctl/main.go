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

package main

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"github.com/openconfig/lemming/cmd/findtrace"
	"github.com/openconfig/lemming/cmd/packet"
	"github.com/openconfig/lemming/cmd/release"
	"github.com/openconfig/lemming/dataplane/luciusctl"
)

func main() {
	cmd := &cobra.Command{
		Use:   "lemctl",
		Short: "lemctl contains a variety of tools used for debug lemming and lucius.",
	}

	internal := &cobra.Command{
		Use:   "internal",
		Short: "Internal contains tools for internal lemming team use.",
	}

	internal.AddCommand(release.NewRelease())
	cmd.AddCommand(internal, packet.New(), findtrace.New(), luciusctl.New())

	if err := cmd.ExecuteContext(context.Background()); err != nil {
		os.Exit(1)
	}
}
