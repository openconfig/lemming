// Copyright 2023 Google LLC
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

package ports

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcapgo"
	"github.com/openconfig/gnmi/errdiff"

	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

const (
	fileHeader = "0a0d0d0a440000004d3c2b1a01000000ffffffffffffffff04000800676f7061636b657402000500616d643634000000030005006c696e757800000000000000440000000100000038000000010000000000000002000500696e7466300000000c0005006c696e757800000009000100090000000000000038000000"
)

func writeHex(t testing.TB, w io.Writer, hexStr string) []byte {
	h, err := hex.DecodeString(hexStr)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := w.Write(h); err != nil {
		t.Fatal(err)
	}
	return h
}

func TestFakeBuild(t *testing.T) {
	tests := []struct {
		desc       string
		openErr    error
		createErr  error
		inDesc     *fwdpb.PortDesc
		inContents string
		wantErr    string
	}{{
		desc:    "wrong desc type",
		inDesc:  &fwdpb.PortDesc{Port: &fwdpb.PortDesc_Cpu{}},
		wantErr: "invalid port type in proto",
	}, {
		desc:    "open error",
		inDesc:  &fwdpb.PortDesc{Port: &fwdpb.PortDesc_Fake{Fake: &fwdpb.FakePortDesc{}}},
		openErr: fmt.Errorf("open error"),
		wantErr: "open error",
	}, {
		desc:    "new reader err",
		inDesc:  &fwdpb.PortDesc{Port: &fwdpb.PortDesc_Fake{Fake: &fwdpb.FakePortDesc{}}},
		wantErr: "EOF",
	}, {
		desc:       "create err",
		inDesc:     &fwdpb.PortDesc{Port: &fwdpb.PortDesc_Fake{Fake: &fwdpb.FakePortDesc{}}},
		inContents: fileHeader,
		createErr:  fmt.Errorf("create error"),
	}, {
		desc:       "sucess",
		inDesc:     &fwdpb.PortDesc{Port: &fwdpb.PortDesc_Fake{Fake: &fwdpb.FakePortDesc{}}},
		inContents: fileHeader,
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			in := &bytes.Buffer{}
			writeHex(t, in, tt.inContents)
			out := &bytes.Buffer{}
			openFile = func(filename string) (io.Reader, error) {
				return in, tt.openErr
			}
			createFile = func(filename string) (io.Writer, error) {
				return out, tt.openErr
			}
			_, err := (fakeBuilder{}).Build(tt.inDesc, nil)
			if d := errdiff.Check(err, tt.wantErr); d != "" {
				t.Fatalf("Build() unexpected error diff: %s", d)
			}
		})
	}
}

type fakeWriter struct {
	bytes.Buffer
	writeErr error
}

func (fw *fakeWriter) Write(data []byte) (int, error) {
	if fw.writeErr == nil {
		return fw.Buffer.Write(data)
	}
	return 0, fw.writeErr
}

func TestFakeWrite(t *testing.T) {
	tests := []struct {
		desc     string
		writeErr error
		wantErr  string
		want     string
	}{{
		desc:     "write error",
		writeErr: fmt.Errorf("write err"),
		wantErr:  "write err",
	}, {
		desc: "success", // packet metadata                                           // packet content
		want: fileHeader + "060000005c000000000000000000000000ca9a3b3c0000003c000000" + "1111111111111111111111110000686900000000000000000000000000000000000000000000000000000000000000000000000000000000000000005c000000",
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			fw := &fakeWriter{
				writeErr: tt.writeErr,
			}
			w, err := pcapgo.NewNgWriter(fw, layers.LinkTypeEthernet)
			if err != nil {
				t.Fatal(err)
			}
			fp := &fakePort{
				packetDst: w,
			}
			timeNow = func() time.Time { return time.Unix(1, 0) }
			_, err = fp.Write(createEthPacket(t))
			err = w.Flush()
			if d := errdiff.Check(err, tt.wantErr); d != "" {
				t.Fatalf("Write() unexpected error diff: %s", d)
			}
			if d := cmp.Diff(fmt.Sprintf("%x", fw.Bytes()), tt.want); d != "" {
				t.Fatalf("Write() unexpected diff(-got,+want)\n:%s", d)
			}
		})
	}
}
