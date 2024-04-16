package gnoi

import (
	"context"
	"net"
	"testing"
	"time"

	spb "github.com/openconfig/gnoi/system"
	"github.com/openconfig/lemming/gnmi"
	"github.com/openconfig/lemming/gnmi/fakedevice"
	"github.com/openconfig/lemming/gnmi/oc/ocpath"
	"github.com/openconfig/ygnmi/ygnmi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestReboot(t *testing.T) {
	grpcServer := grpc.NewServer()
	gnmiServer, err := gnmi.New(grpcServer, "local", nil)
	if err != nil {
		t.Fatal(err)
	}
	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Failed to start listener: %v", err)
	}
	go func() {
		grpcServer.Serve(lis)
	}()

	// Update the interface configuration on the gNMI server.
	client := gnmiServer.LocalClient()

	c, err := ygnmi.NewClient(client, ygnmi.WithTarget("local"))
	if err != nil {
		t.Fatalf("cannot create ygnmi client: %v", err)
	}

	s := newSystem(c)

	ctx := context.Background()
	fakedevice.NewBootTimeTask().Start(ctx, client, "local")

	t.Run("zero-delay", func(t *testing.T) {
		prevTime, err := ygnmi.Get(context.Background(), c, ocpath.Root().System().BootTime().State())
		if err != nil {
			t.Fatal(err)
		}
		if _, err := s.Reboot(ctx, &spb.RebootRequest{}); err != nil {
			t.Fatal(err)
		}
		afterTime, err := ygnmi.Get(context.Background(), c, ocpath.Root().System().BootTime().State())
		if err != nil {
			t.Fatal(err)
		}
		if !(prevTime < afterTime) {
			t.Errorf("boot time did not update after reboot")
		}
	})

	t.Run("one-second-delay", func(t *testing.T) {
		prevTime, err := ygnmi.Get(context.Background(), c, ocpath.Root().System().BootTime().State())
		if err != nil {
			t.Fatal(err)
		}
		const delay = 1e9
		if _, err := s.Reboot(ctx, &spb.RebootRequest{Delay: delay}); err != nil {
			t.Fatal(err)
		}

		var afterTime uint64
		tryN := 50
		for i := 0; i != tryN; i++ {
			var err error
			afterTime, err = ygnmi.Get(context.Background(), c, ocpath.Root().System().BootTime().State())
			if err != nil {
				t.Fatal(err)
			}
			if prevTime < afterTime {
				break
			}
			time.Sleep(100 * time.Millisecond)
		}
		if !(prevTime < afterTime) {
			t.Errorf("boot time did not update after reboot")
		}
	})

	t.Run("cancel-no-pending", func(t *testing.T) {
		if _, err := s.CancelReboot(ctx, &spb.CancelRebootRequest{}); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("cancel-pending", func(t *testing.T) {
		prevTime, err := ygnmi.Get(context.Background(), c, ocpath.Root().System().BootTime().State())
		if err != nil {
			t.Fatal(err)
		}
		const delay = 120e9
		if _, err := s.Reboot(ctx, &spb.RebootRequest{Delay: delay}); err != nil {
			t.Fatal(err)
		}
		if _, err := s.CancelReboot(ctx, &spb.CancelRebootRequest{}); err != nil {
			t.Fatal(err)
		}

		afterTime, err := ygnmi.Get(context.Background(), c, ocpath.Root().System().BootTime().State())
		if err != nil {
			t.Fatal(err)
		}
		if prevTime != afterTime {
			t.Errorf("boot did not get cancelled")
		}

		// Couple more no-op cancels.
		if _, err := s.CancelReboot(ctx, &spb.CancelRebootRequest{}); err != nil {
			t.Fatal(err)
		}
		if _, err := s.CancelReboot(ctx, &spb.CancelRebootRequest{}); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("reboot-while-pending", func(t *testing.T) {
		const delay = 120e9
		if _, err := s.Reboot(ctx, &spb.RebootRequest{Delay: delay}); err != nil {
			t.Fatal(err)
		}
		if _, err := s.Reboot(ctx, &spb.RebootRequest{}); status.Convert(err).Code() != codes.AlreadyExists {
			t.Fatalf("Expected AlreadyExists error, got %v", err)
		}
		if _, err := s.Reboot(ctx, &spb.RebootRequest{Delay: 1}); status.Convert(err).Code() != codes.AlreadyExists {
			t.Fatalf("Expected AlreadyExists error, got %v", err)
		}
		if _, err := s.CancelReboot(ctx, &spb.CancelRebootRequest{}); err != nil {
			t.Fatal(err)
		}
	})

	closeProximityTests := []struct {
		desc  string
		delay uint64
	}{{
		desc:  "cancel-while-possibly-pending-1ns",
		delay: 1e0,
	}, {
		desc:  "cancel-while-possibly-pending-10ns",
		delay: 1e1,
	}, {
		desc:  "cancel-while-possibly-pending-100ns",
		delay: 1e2,
	}, {
		desc:  "cancel-while-possibly-pending-1us",
		delay: 1e3,
	}, {
		desc:  "cancel-while-possibly-pending-10us",
		delay: 1e4,
	}, {
		desc:  "cancel-while-possibly-pending-100us",
		delay: 1e5,
	}, {
		desc:  "cancel-while-possibly-pending-1ms",
		delay: 1e6,
	}, {
		desc:  "cancel-while-possibly-pending-10ms",
		delay: 1e7,
	}, {
		desc:  "cancel-while-possibly-pending-100ms",
		delay: 1e8,
	}}

	for _, tt := range closeProximityTests {
		t.Run(tt.desc, func(t *testing.T) {
			if _, err := s.Reboot(ctx, &spb.RebootRequest{Delay: tt.delay}); err != nil {
				t.Fatal(err)
			}
			if _, err := s.CancelReboot(ctx, &spb.CancelRebootRequest{}); err != nil {
				t.Fatal(err)
			}
			if _, err := s.CancelReboot(ctx, &spb.CancelRebootRequest{}); err != nil {
				t.Fatal(err)
			}
			if _, err := s.CancelReboot(ctx, &spb.CancelRebootRequest{}); err != nil {
				t.Fatal(err)
			}
			if _, err := s.CancelReboot(ctx, &spb.CancelRebootRequest{}); err != nil {
				t.Fatal(err)
			}
			if _, err := s.CancelReboot(ctx, &spb.CancelRebootRequest{}); err != nil {
				t.Fatal(err)
			}
		})
	}
}
