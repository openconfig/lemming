// Copyright 2026 Google LLC
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

package ccgen

import (
	"strings"
	"testing"

	"github.com/openconfig/lemming/dataplane/apigen/docparser"
	"github.com/openconfig/lemming/dataplane/apigen/saiast"
)

// TestGenClientEnumListToSAIRespectsBuffer guards against regressing the SAI
// two-pass list protocol in the generated convert_list_*_to_sai helpers.
//
// SONiC (e.g. DebugCounterOrch querying
// SAI_SWITCH_ATTR_SUPPORTED_INGRESS_DROP_REASON_LIST) first calls
// get_switch_attribute with a null list and a count of 0 purely to learn how
// many elements it must allocate. If the generated helper writes into the
// caller's buffer unconditionally it dereferences a null pointer, and the
// attribute is never read, so the corresponding debug counter is never created.
func TestGenClientEnumListToSAIRespectsBuffer(t *testing.T) {
	doc := &docparser.SAIInfo{
		Attrs: map[string]*docparser.Attr{},
		Enums: map[string][]*docparser.Enum{
			"sai_in_drop_reason_t": {
				{Name: "SAI_IN_DROP_REASON_LPM4_MISS", Value: 0},
				{Name: "SAI_IN_DROP_REASON_LPM6_MISS", Value: 1},
			},
		},
	}

	files, err := GenClient(doc, &saiast.SAIAPI{}, "dataplane/proto/sai", "dataplane/standalone/sai", nil)
	if err != nil {
		t.Fatalf("GenClient() got unexpected error: %v", err)
	}

	got, ok := files["enum.cc"]
	if !ok {
		t.Fatalf("GenClient() did not generate enum.cc, got files: %v", files)
	}

	for _, want := range []string{
		"void convert_list_sai_in_drop_reason_t_to_sai(",
		"const uint32_t required = static_cast<uint32_t>(proto_list.size());",
		"if (list != nullptr && *count >= required) {",
		"*count = required;",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("GenClient() enum.cc missing %q", want)
		}
	}

	// The unguarded loop must not come back: it writes to list[i] even when the
	// caller passed a null or undersized buffer.
	if strings.Contains(got, "for (int i = 0; i < proto_list.size(); i++) {") {
		t.Errorf("GenClient() enum.cc writes to the caller's buffer without checking its size")
	}
}

func TestGenClientStatsExtGeneratesRPC(t *testing.T) {
	doc := &docparser.SAIInfo{
		Attrs: map[string]*docparser.Attr{
			"SWITCH": {},
		},
		Enums: map[string][]*docparser.Enum{},
	}
	sai := &saiast.SAIAPI{
		Ifaces: []*saiast.SAIInterface{
			{
				Name: "sai_switch_api_t",
				Funcs: []*saiast.TypeDecl{
					{
						Name: "sai_get_switch_stats_ext_fn",
						Typ:  "sai_get_switch_stats_ext_fn",
					},
				},
			},
		},
		Funcs: map[string]*saiast.SAIFunc{
			"sai_get_switch_stats_ext_fn": {
				Name:       "get_switch_stats_ext",
				ReturnType: "sai_status_t",
				Params: []saiast.TypeDecl{
					{Name: "switch_id", Typ: "sai_object_id_t"},
					{Name: "number_of_counters", Typ: "uint32_t"},
					{Name: "counter_ids", Typ: "const sai_stat_id_t *"},
					{Name: "mode", Typ: "sai_stats_mode_t"},
					{Name: "counters", Typ: "uint64_t *"},
				},
			},
		},
	}

	files, err := GenClient(doc, sai, "dataplane/proto/sai", "dataplane/standalone/sai", nil)
	if err != nil {
		t.Fatalf("GenClient() got unexpected error: %v", err)
	}

	got, ok := files["switch.cc"]
	if !ok {
		t.Fatalf("GenClient() did not generate switch.cc, got files: %v", files)
	}

	for _, want := range []string{
		"sai_status_t l_get_switch_stats_ext(",
		"lemming::dataplane::sai::GetSwitchStatsRequest req;",
		"lemming::dataplane::sai::GetSwitchStatsResponse resp;",
		"switch_->GetSwitchStats(&context, req, &resp);",
		"counters[i] = resp.values(i);",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("GenClient() switch.cc missing %q", want)
		}
	}
}
