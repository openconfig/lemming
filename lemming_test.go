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

package lemming

import (
	"context"
	"testing"

	"github.com/openconfig/gnmi/errdiff"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"

	// gNMI
	gnmipb "github.com/openconfig/gnmi/proto/gnmi"

	// gNOI
	bpb "github.com/openconfig/gnoi/bgp"
	cmpb "github.com/openconfig/gnoi/cert"
	diagpb "github.com/openconfig/gnoi/diag"
	frpb "github.com/openconfig/gnoi/factory_reset"
	fpb "github.com/openconfig/gnoi/file"
	hpb "github.com/openconfig/gnoi/healthz"
	lpb "github.com/openconfig/gnoi/layer2"
	mpb "github.com/openconfig/gnoi/mpls"
	ospb "github.com/openconfig/gnoi/os"
	otpb "github.com/openconfig/gnoi/otdr"
	spb "github.com/openconfig/gnoi/system"
	wrpb "github.com/openconfig/gnoi/wavelength_router"
	// gNSI
	// authzpb "github.com/openconfig/gnsi/authz/authz_go_proto"
	// certpb "github.com/openconfig/gnsi/cert/cert_go_proto"
	// consolepb "github.com/openconfig/gnsi/console/console_go_proto"
	// pathzpb "github.com/openconfig/gnsi/pathz/pathz_go_proto"
	// sshpb "github.com/openconfig/gnsi/ssh/ssh_go_proto"
)

func TestFakeGNMI(t *testing.T) {
	f := startLemming(t)
	defer f.Stop()
	conn, err := grpc.Dial(f.GNMIAddr(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("failed to Dial fake: %v", err)
	}
	want := &gnmipb.GetResponse{
		Notification: []*gnmipb.Notification{{
			Update: []*gnmipb.Update{{
				Path: &gnmipb.Path{
					Elem: []*gnmipb.PathElem{
						{Name: "intefaces"},
						{Name: "inteface", Key: map[string]string{"name": "eth0"}},
						{Name: "mtu"},
					},
				},
				Val: &gnmipb.TypedValue{
					Value: &gnmipb.TypedValue_IntVal{
						IntVal: 1500,
					},
				},
			}},
		}},
	}
	f.gnmiServer.GetResponses = []interface{}{want}
	cGNMI := gnmipb.NewGNMIClient(conn)
	resp, err := cGNMI.Get(context.Background(), &gnmipb.GetRequest{})
	if err != nil {
		t.Fatalf("gnmi.Get failed: %v", err)
	}
	if !proto.Equal(resp, want) {
		t.Fatalf("gnmi.Get failed got %v, want %v", resp, want)
	}
}

func TestStop(t *testing.T) {
	t.Run("errors", func(t *testing.T) {
		f := startLemming(t)
		// Close the listener so the get must fail. Sleep to ensure listener is closed before Get.
		f.GNMIListener().Close()
		err := f.Stop()
		if s := errdiff.Check(err, "use of closed network connection"); s != "" {
			t.Fatalf("failed to get error on close: %s", s)
		}
	})

	t.Run("success", func(t *testing.T) {
		f := startLemming(t)
		if err := f.Stop(); err != nil {
			t.Fatalf("did not get nil error on stop, got: %v", err)
		}
	})
}

func TestFakeGNOI(t *testing.T) {
	f := startLemming(t)
	defer f.stop()
	conn, err := grpc.Dial(f.GNMIAddr(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("failed to Dial fake: %v", err)
	}
	cBGP := bpb.NewBGPClient(conn)
	_, err = cBGP.ClearBGPNeighbor(context.Background(), &bpb.ClearBGPNeighborRequest{})
	if err == nil {
		t.Errorf("gnoi.BGP.ClearBGPNeighbor failed to return error")
	}

	cCertMgmt := cmpb.NewCertificateManagementClient(conn)
	_, err = cCertMgmt.CanGenerateCSR(context.Background(), &cmpb.CanGenerateCSRRequest{})
	if err == nil {
		t.Errorf("gnoi.Cert.CanGenerateCSR failed to return error")
	}

	cDiag := diagpb.NewDiagClient(conn)
	_, err = cDiag.GetBERTResult(context.Background(), &diagpb.GetBERTResultRequest{})
	if err == nil {
		t.Errorf("gnoi.Diag.GetBERTResult failed to return error")
	}

	cReset := frpb.NewFactoryResetClient(conn)
	_, err = cReset.Start(context.Background(), &frpb.StartRequest{})
	if err == nil {
		t.Errorf("gnoi.FactoryReset.Start failed to return error")
	}

	cFile := fpb.NewFileClient(conn)
	scFile, err := cFile.Get(context.Background(), &fpb.GetRequest{})
	if err != nil {
		t.Errorf("gnoi.File.Get failed to get stream client: %v", err)
	}
	_, err = scFile.Recv()
	if err == nil {
		t.Errorf("gnoi.File.Get failed to return error")
	}

	cHealthz := hpb.NewHealthzClient(conn)
	_, err = cHealthz.Get(context.Background(), &hpb.GetRequest{})
	if err == nil {
		t.Errorf("gnoi.Healthz.Get failed to return error")
	}

	cLayer2 := lpb.NewLayer2Client(conn)
	_, err = cLayer2.ClearLLDPInterface(context.Background(), &lpb.ClearLLDPInterfaceRequest{})
	if err == nil {
		t.Errorf("gnoi.Layer2.ClearLLDPInterface failed to return error")
	}

	cMPLS := mpb.NewMPLSClient(conn)
	_, err = cMPLS.ClearLSP(context.Background(), &mpb.ClearLSPRequest{})
	if err == nil {
		t.Errorf("gnoi.MPLS.ClearLSP failed to return error")
	}

	cOS := ospb.NewOSClient(conn)
	_, err = cOS.Activate(context.Background(), &ospb.ActivateRequest{})
	if err == nil {
		t.Errorf("gnoi.OS.Activate failed to return error")
	}

	cOTDR := otpb.NewOTDRClient(conn)
	scOTDR, err := cOTDR.Initiate(context.Background(), &otpb.InitiateRequest{})
	if err != nil {
		t.Errorf("gnoi.OTDR.Initiate failed to get stream client: %v", err)
	}
	_, err = scOTDR.Recv()
	if err == nil {
		t.Errorf("gnoi.OTDR.Initiate failed to return error")
	}

	cSystem := spb.NewSystemClient(conn)
	_, err = cSystem.Reboot(context.Background(), &spb.RebootRequest{})
	if err == nil {
		t.Errorf("gnoi.System.Reboot failed to return error")
	}

	cWaveLengthRouter := wrpb.NewWavelengthRouterClient(conn)
	scWaveLengthRouter, err := cWaveLengthRouter.AdjustSpectrum(context.Background(), &wrpb.AdjustSpectrumRequest{})
	if err != nil {
		t.Errorf("gnoi.WaveLengthRouter.AdjustSpectrum failed to get stream client: %v", err)
	}
	_, err = scWaveLengthRouter.Recv()
	if err == nil {
		t.Errorf("gnoi.WaveLengthRouter.AdjustSpectrum failed to return error")
	}
}

/*
func TestGNSI(t *testing.T) {
	desc := "gnsi.Authz.Rotate"
	t.Run(desc, func(t *testing.T) {
		f := startLemming(t)
		defer f.stop()
		conn, err := grpc.Dial(f.Addr(), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			t.Fatalf("failed to Dial fake: %v", err)
		}
		cAuthz := authzpb.NewAuthzManagementClient(conn)
		scAuthz, err := cAuthz.Rotate(context.Background())
		if err != nil {
			t.Errorf("%s failed to get stream client: %v", desc, err)
		}
		_, err = scAuthz.Recv()
		if err == nil {
			t.Errorf("%s failed to return error", desc)
		}
	})

	desc = "gnsi.Cert.Install"
	t.Run(desc, func(t *testing.T) {
		f := startLemming(t)
		defer f.stop()
		conn, err := grpc.Dial(f.Addr(), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			t.Fatalf("failed to Dial fake: %v", err)
		}
		cCert := certpb.NewCertificateManagementClient(conn)
		scCert, err := cCert.Install(context.Background())
		if err != nil {
			t.Errorf("%s failed to get stream client: %v", desc, err)
		}
		_, err = scCert.Recv()
		if err == nil {
			t.Errorf("%s failed to return error", desc)
		}
	})

	desc = "gnsi.Console.MutateAccountPassword"
	t.Run(desc, func(t *testing.T) {
		f := startLemming(t)
		defer f.stop()
		conn, err := grpc.Dial(f.Addr(), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			t.Fatalf("failed to Dial fake: %v", err)
		}
		cConsole := consolepb.NewConsoleClient(conn)
		scConsole, err := cConsole.MutateAccountPassword(context.Background())
		if err != nil {
			t.Errorf("%s failed to get stream client: %v", desc, err)
		}
		_, err = scConsole.Recv()
		if err == nil {
			t.Errorf("%s failed to return error", desc)
		}
	})

	desc = "gnsi.Pathz.Install"
	t.Run(desc, func(t *testing.T) {
		f := startLemming(t)
		defer f.stop()
		conn, err := grpc.Dial(f.Addr(), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			t.Fatalf("failed to Dial fake: %v", err)
		}
		cPathz := pathzpb.NewPathzManagementClient(conn)
		scPathz, err := cPathz.Install(context.Background())
		if err != nil {
			t.Errorf("%s failed to get stream client: %v", desc, err)
		}
		_, err = scPathz.Recv()
		if err == nil {
			t.Errorf("%s failed to return error", desc)
		}
	})

	desc = "gnsi.SSH.MutateAccountCredentials"
	t.Run(desc, func(t *testing.T) {
		f := startLemming(t)
		defer f.stop()
		conn, err := grpc.Dial(f.Addr(), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			t.Fatalf("failed to Dial fake: %v", err)
		}
		cSSH := sshpb.NewSshClient(conn)
		scSSH, err := cSSH.MutateAccountCredentials(context.Background())
		if err != nil {
			t.Fatalf("%s failed to get stream client: %v", desc, err)
		}
		_, err = scSSH.Recv()
		if err == nil {
			t.Fatalf("%s failed to return error", desc)
		}

	})
}
*/

func startLemming(t *testing.T, opts ...Option) *Device {
	f, err := New("fakedevice", "unix:/tmp/zserv.api", opts...)
	if err != nil {
		t.Fatalf("Failed to start lemming: %v", err)
	}
	return f
}
