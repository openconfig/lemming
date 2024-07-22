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

package findtrace

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "findtrace [search term]",
		Short: "Searches lucius logs for the specified term and returns the matching packet traces",
		RunE: func(cmd *cobra.Command, args []string) error {
			scanner := bufio.NewScanner(os.Stdin)
			var trace strings.Builder

			for scanner.Scan() {
				text := scanner.Text()
				if strings.Contains(text, "Packet Trace:") {
					str := trace.String()
					if len(os.Args) == 1 || strings.Contains(str, os.Args[1]) {
						out := strings.ReplaceAll(str, os.Args[1], color.RedString(os.Args[1]))
						fmt.Println(out)
					}
					trace = strings.Builder{}
					continue
				}
				if strings.Contains(text, "packet.go") {
					trace.WriteString(text)
					trace.WriteString("\n")
				}
			}
			str := trace.String()
			if len(os.Args) == 1 || strings.Contains(str, os.Args[1]) {
				out := strings.ReplaceAll(str, os.Args[1], color.RedString(os.Args[1]))
				fmt.Println(out)
			}
			return nil
		},
	}
	return cmd
}
