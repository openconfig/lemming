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

package sysrib

import (
	"context"

	"github.com/openconfig/lemming/dataplane/handlers"
	dpb "github.com/openconfig/lemming/proto/dataplane"
	"github.com/openconfig/ygnmi/ygnmi"
)

func init() {
	programRouteFns = append(programRouteFns, func(ctx context.Context, c *ygnmi.Client, rr *dpb.Route) error {
		_, err := ygnmi.Replace(ctx, c, handlers.RouteQuery(rr.GetPrefix().GetVrfId(), rr.GetPrefix().GetCidr()), rr, ygnmi.WithSetFallbackEncoding())
		return err
	})
}
