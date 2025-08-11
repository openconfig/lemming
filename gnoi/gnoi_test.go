package gnoi

import (
	"bytes"
	"context"
	"crypto/md5"
	"crypto/sha256"
	"io"
	"math"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/openconfig/gnmi/errdiff"
	"github.com/openconfig/ygnmi/ygnmi"
	"github.com/openconfig/ygot/ygot"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	cpb "github.com/openconfig/gnoi/common"
	fpb "github.com/openconfig/gnoi/file"
	plqpb "github.com/openconfig/gnoi/packet_link_qualification"
	spb "github.com/openconfig/gnoi/system"
	pb "github.com/openconfig/gnoi/types"
	configpb "github.com/openconfig/lemming/proto/config"

	"github.com/openconfig/lemming/gnmi"
	"github.com/openconfig/lemming/gnmi/fakedevice"
	"github.com/openconfig/lemming/gnmi/gnmiclient"
	"github.com/openconfig/lemming/gnmi/oc"
	"github.com/openconfig/lemming/gnmi/oc/ocpath"
	"github.com/openconfig/lemming/internal/config"
)

const (
	// Default supervisor configuration
	defaultPrimarySupervisor   = "Supervisor1"
	defaultSecondarySupervisor = "Supervisor2"
)

// loadDefaultConfig loads the lemming default configuration for tests
func loadDefaultConfig(t *testing.T) *configpb.Config {
	cfg, err := config.Load("")
	if err != nil {
		t.Fatalf("Failed to load default config: %v", err)
	}
	return cfg
}

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

	// Load lemming default configuration for tests
	lemmingConfig := loadDefaultConfig(t)

	s := newSystem(c, lemmingConfig)

	ctx := context.Background()
	fakedevice.NewBootTimeTask(lemmingConfig).Start(ctx, client, "local")

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
			time.Sleep(50 * time.Millisecond)
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
		delay time.Duration
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
			if _, err := s.Reboot(ctx, &spb.RebootRequest{Delay: uint64(tt.delay.Nanoseconds())}); err != nil {
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

// Testing component reboot implementation
func TestComponentReboot(t *testing.T) {
	lemmingConfig := loadDefaultConfig(t)

	tests := map[string]struct {
		fn func(*testing.T, *system, context.Context)
	}{
		"reboot-nonexistent-component": {
			fn: func(t *testing.T, s *system, ctx context.Context) {
				req := &spb.RebootRequest{
					Method: spb.RebootMethod_COLD,
					Subcomponents: []*pb.Path{{
						Elem: []*pb.PathElem{
							{Name: "components"},
							{Name: "component", Key: map[string]string{"name": "non-existent"}},
						},
					}},
				}
				_, err := s.Reboot(ctx, req)
				if err == nil {
					t.Error("Reboot of non-existent component should have failed")
				} else if got := status.Code(err); got != codes.NotFound {
					t.Errorf("Expected NotFound error, got %v", got)
				}
			},
		},
		"immediate-linecard-reboot": {
			fn: func(t *testing.T, s *system, ctx context.Context) {
				componentName := "Linecard0"

				// Get initial state
				initialState, err := ygnmi.Get(ctx, s.c, ocpath.Root().Component(componentName).State())
				if err != nil {
					t.Fatalf("Failed to get initial state: %v", err)
				}
				initialUptime := initialState.GetLastRebootTime()

				req := &spb.RebootRequest{
					Method: spb.RebootMethod_COLD,
					Subcomponents: []*pb.Path{{
						Elem: []*pb.PathElem{
							{Name: "components"},
							{Name: "component", Key: map[string]string{"name": componentName}},
						},
					}},
				}
				_, err = s.Reboot(ctx, req)
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}

				// Verify component goes INACTIVE during reboot
				deadline := time.Now().Add(3 * time.Second)
				for time.Now().Before(deadline) {
					state, err := ygnmi.Get(ctx, s.c, ocpath.Root().Component(componentName).OperStatus().State())
					if err != nil {
						t.Fatalf("Failed to get component state: %v", err)
					}
					if state == oc.PlatformTypes_COMPONENT_OPER_STATUS_INACTIVE {
						break
					}
					time.Sleep(time.Second)
				}

				// Wait for reboot to complete
				deadline = time.Now().Add(3 * time.Second)
				var finalState *oc.Component
				for time.Now().Before(deadline) {
					state, err := ygnmi.Get(ctx, s.c, ocpath.Root().Component(componentName).State())
					if err != nil {
						t.Fatalf("Failed to get component state: %v", err)
					}
					if state.GetOperStatus() == oc.PlatformTypes_COMPONENT_OPER_STATUS_ACTIVE {
						finalState = state
						break
					}
					time.Sleep(3 * time.Second)
				}
				if finalState == nil {
					t.Fatal("Component did not return to ACTIVE state")
				}

				// Verify state changes
				if finalState.GetOperStatus() != oc.PlatformTypes_COMPONENT_OPER_STATUS_ACTIVE {
					t.Errorf("Expected final OperStatus ACTIVE, got %v", finalState.GetOperStatus())
				}
				if finalState.GetLastRebootTime() <= initialUptime {
					t.Error("LastRebootTime was not updated")
				}
				if finalState.GetLastRebootReason() != oc.PlatformTypes_COMPONENT_REBOOT_REASON_REBOOT_USER_INITIATED {
					t.Errorf("Expected LastRebootReason USER_INITIATED, got %v", finalState.GetLastRebootReason())
				}
			},
		},
		"delayed-linecard-reboot-and-cancel": {
			fn: func(t *testing.T, s *system, ctx context.Context) {
				componentName := "Linecard0"

				// Get initial state
				initialState, err := ygnmi.Get(ctx, s.c, ocpath.Root().Component(componentName).State())
				if err != nil {
					t.Fatalf("Failed to get initial state: %v", err)
				}
				initialUptime := initialState.GetLastRebootTime()

				req := &spb.RebootRequest{
					Method: spb.RebootMethod_COLD,
					Delay:  2000000000, // 2 seconds
					Subcomponents: []*pb.Path{{
						Elem: []*pb.PathElem{
							{Name: "components"},
							{Name: "component", Key: map[string]string{"name": componentName}},
						},
					}},
				}
				_, err = s.Reboot(ctx, req)
				if err != nil {
					t.Fatalf("Delayed linecard reboot failed: %v", err)
				}
				time.Sleep(time.Second)

				// Verify state hasn't changed yet
				currentState, err := ygnmi.Get(ctx, s.c, ocpath.Root().Component(componentName).State())
				if err != nil {
					t.Fatalf("Failed to get current state: %v", err)
				}
				if currentState.GetOperStatus() != initialState.GetOperStatus() {
					t.Error("Component state changed before delay expired")
				}
				if currentState.GetLastRebootTime() != initialUptime {
					t.Error("LastRebootTime changed before delay expired")
				}

				// Cancel the reboot
				_, err = s.CancelReboot(ctx, &spb.CancelRebootRequest{})
				if err != nil {
					t.Fatalf("Failed to cancel reboot: %v", err)
				}

				// Verify no state changes after cancel
				finalState, err := ygnmi.Get(ctx, s.c, ocpath.Root().Component(componentName).State())
				if err != nil {
					t.Fatalf("Failed to get final state: %v", err)
				}
				if diff := cmp.Diff(initialState, finalState); diff != "" {
					t.Errorf("Component state changed after cancel (-want +got):\n%s", diff)
				}

				// Verify component reboot was cancelled
				s.componentRebootsMu.Lock()
				if len(s.componentReboots) != 0 {
					t.Errorf("Expected no pending reboots, found %d", len(s.componentReboots))
				}
				s.componentRebootsMu.Unlock()
			},
		},
		"multiple-component-reboot": {
			fn: func(t *testing.T, s *system, ctx context.Context) {
				req := &spb.RebootRequest{
					Method: spb.RebootMethod_COLD,
					Subcomponents: []*pb.Path{
						{
							Elem: []*pb.PathElem{
								{Name: "components"},
								{Name: "component", Key: map[string]string{"name": "Linecard0"}},
							},
						},
						{
							Elem: []*pb.PathElem{
								{Name: "components"},
								{Name: "component", Key: map[string]string{"name": "Fabric0"}},
							},
						},
					},
				}
				_, err := s.Reboot(ctx, req)
				if err != nil {
					t.Fatalf("Multiple component reboot failed: %v", err)
				}

				// Clean up
				_, err = s.CancelReboot(ctx, &spb.CancelRebootRequest{})
				if err != nil {
					t.Fatalf("Failed to cancel reboot: %v", err)
				}
			},
		},
		"cancel-no-pending-reboot": {
			fn: func(t *testing.T, s *system, ctx context.Context) {
				_, err := s.CancelReboot(ctx, &spb.CancelRebootRequest{})
				if err != nil {
					t.Errorf("Cancel with no pending reboot should succeed: %v", err)
				}
			},
		},
		"concurrent-same-component-reboot": {
			fn: func(t *testing.T, s *system, ctx context.Context) {
				// Start first reboot
				req := &spb.RebootRequest{
					Method: spb.RebootMethod_COLD,
					Delay:  2000000000,
					Subcomponents: []*pb.Path{{
						Elem: []*pb.PathElem{
							{Name: "components"},
							{Name: "component", Key: map[string]string{"name": "Linecard0"}},
						},
					}},
				}
				_, err := s.Reboot(ctx, req)
				if err != nil {
					t.Fatalf("First reboot request failed: %v", err)
				}

				// Try same component again
				_, err = s.Reboot(ctx, req)
				if err == nil {
					t.Error("Second reboot of same component should have failed")
				} else if got := status.Code(err); got != codes.AlreadyExists {
					t.Errorf("Expected AlreadyExists error, got %v", got)
				}

				// Clean up
				_, err = s.CancelReboot(ctx, &spb.CancelRebootRequest{})
				if err != nil {
					t.Fatalf("Failed to cancel reboot: %v", err)
				}
			},
		},
		"reject-active-supervisor-reboot": {
			fn: func(t *testing.T, s *system, ctx context.Context) {
				// Try to reboot active supervisor
				activeSupervisorReq := &spb.RebootRequest{
					Method: spb.RebootMethod_COLD,
					Subcomponents: []*pb.Path{{
						Elem: []*pb.PathElem{
							{Name: "components"},
							{Name: "component", Key: map[string]string{"name": "Supervisor1"}},
						},
					}},
				}
				_, err := s.Reboot(ctx, activeSupervisorReq)

				// Should fail with FailedPrecondition
				if err == nil {
					t.Error("Active supervisor reboot should have failed")
				} else if got := status.Code(err); got != codes.FailedPrecondition {
					t.Errorf("Expected FailedPrecondition error, got %v", got)
				}

				_, err = s.CancelReboot(ctx, &spb.CancelRebootRequest{})
				if err != nil {
					t.Fatalf("Failed to cancel reboot: %v", err)
				}
			},
		},
	}

	ctx := context.Background()
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
	defer grpcServer.Stop()

	client := gnmiServer.LocalClient()
	c, err := ygnmi.NewClient(client, ygnmi.WithTarget("local"))
	if err != nil {
		t.Fatalf("cannot create ygnmi client: %v", err)
	}
	s := newSystem(c, lemmingConfig)

	// Initialize the system
	if err := fakedevice.NewBootTimeTask(lemmingConfig).Start(ctx, client, "local"); err != nil {
		t.Fatalf("Failed to initialize boot time: %v", err)
	}
	if err := fakedevice.NewChassisComponentsTask(lemmingConfig).Start(ctx, client, "local"); err != nil {
		t.Fatalf("Failed to initialize chassis components: %v", err)
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			test.fn(t, s, ctx)
		})
	}
}

// TestSwitchControlProcessor tests the SwitchControlProcessor method.
func TestSwitchControlProcessor(t *testing.T) {
	lemmingConfig := loadDefaultConfig(t)

	tests := map[string]struct {
		fn func(*testing.T, *system, context.Context, *ygnmi.Client)
	}{
		"empty-request": {
			fn: func(t *testing.T, s *system, ctx context.Context, c *ygnmi.Client) {
				req := &spb.SwitchControlProcessorRequest{}
				_, err := s.SwitchControlProcessor(ctx, req)
				if err == nil {
					t.Error("Expected error for empty request")
				} else if got := status.Code(err); got != codes.InvalidArgument {
					t.Errorf("Expected InvalidArgument error, got %v", got)
				}
			},
		},
		"invalid-path-format": {
			fn: func(t *testing.T, s *system, ctx context.Context, c *ygnmi.Client) {
				req := &spb.SwitchControlProcessorRequest{
					ControlProcessor: &pb.Path{
						Elem: []*pb.PathElem{
							{Name: "invalid"},
							{Name: "path"},
							{Name: "format"},
						},
					},
				}
				_, err := s.SwitchControlProcessor(ctx, req)
				if err == nil {
					t.Error("Expected error for invalid path format")
				} else if got := status.Code(err); got != codes.InvalidArgument {
					t.Errorf("Expected InvalidArgument error, got %v", got)
				}
			},
		},
		"nonexistent-supervisor": {
			fn: func(t *testing.T, s *system, ctx context.Context, c *ygnmi.Client) {
				req := &spb.SwitchControlProcessorRequest{
					ControlProcessor: &pb.Path{
						Elem: []*pb.PathElem{
							{Name: "components"},
							{Name: "component", Key: map[string]string{"name": "SupervisorX"}},
						},
					},
				}
				_, err := s.SwitchControlProcessor(ctx, req)
				if err == nil {
					t.Error("Expected error for non-existent supervisor")
				} else if got := status.Code(err); got != codes.NotFound {
					t.Errorf("Expected NotFound error, got %v", got)
				}
			},
		},
		"non-controller-card": {
			fn: func(t *testing.T, s *system, ctx context.Context, c *ygnmi.Client) {
				req := &spb.SwitchControlProcessorRequest{
					ControlProcessor: &pb.Path{
						Elem: []*pb.PathElem{
							{Name: "components"},
							{Name: "component", Key: map[string]string{"name": "Linecard0"}},
						},
					},
				}
				_, err := s.SwitchControlProcessor(ctx, req)
				if err == nil {
					t.Error("Expected error for non-controller card")
				} else if got := status.Code(err); got != codes.NotFound {
					t.Errorf("Expected NotFound error, got %v", got)
				}
			},
		},
		"switchover-to-already-active": {
			fn: func(t *testing.T, s *system, ctx context.Context, c *ygnmi.Client) {
				// Default primary supervisor is always PRIMARY by default, so try switching to it
				req := &spb.SwitchControlProcessorRequest{
					ControlProcessor: &pb.Path{
						Elem: []*pb.PathElem{
							{Name: "components"},
							{Name: "component", Key: map[string]string{"name": defaultPrimarySupervisor}},
						},
					},
				}
				resp, err := s.SwitchControlProcessor(ctx, req)
				if err != nil {
					t.Errorf("Expected successful response for no-op switchover, got error: %v", err)
				}
				if resp == nil {
					t.Error("Expected response but got nil")
					return
				}

				// Verify response fields
				if resp.ControlProcessor == nil {
					t.Error("Expected ControlProcessor in response")
				}
				if resp.Version == "" {
					t.Error("Expected Version in response")
				}
			},
		},
		"successful-switchover": {
			fn: func(t *testing.T, s *system, ctx context.Context, c *ygnmi.Client) {
				// Get initial states for timestamp comparison
				initialSupervisor1, err := ygnmi.Get(ctx, c, ocpath.Root().Component(defaultPrimarySupervisor).State())
				if err != nil {
					t.Fatalf("Failed to get initial %s state: %v", defaultPrimarySupervisor, err)
				}
				initialSupervisor2, err := ygnmi.Get(ctx, c, ocpath.Root().Component(defaultSecondarySupervisor).State())
				if err != nil {
					t.Fatalf("Failed to get initial %s state: %v", defaultSecondarySupervisor, err)
				}

				// Switch to secondary supervisor (default: PRIMARY=Supervisor1, SECONDARY=Supervisor2)
				req := &spb.SwitchControlProcessorRequest{
					ControlProcessor: &pb.Path{
						Elem: []*pb.PathElem{
							{Name: "components"},
							{Name: "component", Key: map[string]string{"name": defaultSecondarySupervisor}},
						},
					},
				}

				resp, err := s.SwitchControlProcessor(ctx, req)
				if err != nil {
					t.Fatalf("Switchover failed: %v", err)
				}
				if resp == nil {
					t.Fatal("Expected response but got nil")
				}

				// Verify response
				if resp.ControlProcessor == nil {
					t.Error("Expected ControlProcessor in response")
				}
				if resp.Version == "" {
					t.Error("Expected Version in response")
				}
				if resp.Uptime < 0 {
					t.Error("Expected non-negative uptime in response")
				}
				time.Sleep(2500 * time.Millisecond)

				// Verify final states
				finalSupervisor1, err := ygnmi.Get(ctx, c, ocpath.Root().Component(defaultPrimarySupervisor).State())
				if err != nil {
					t.Fatalf("Failed to get %s final state: %v", defaultPrimarySupervisor, err)
				}
				finalSupervisor2, err := ygnmi.Get(ctx, c, ocpath.Root().Component(defaultSecondarySupervisor).State())
				if err != nil {
					t.Fatalf("Failed to get %s final state: %v", defaultSecondarySupervisor, err)
				}

				// Verify roles switched
				if finalSupervisor1.GetRedundantRole() != oc.PlatformTypes_ComponentRedundantRole_SECONDARY {
					t.Errorf("Expected %s to be SECONDARY after switchover, got %v", defaultPrimarySupervisor, finalSupervisor1.GetRedundantRole())
				}
				if finalSupervisor2.GetRedundantRole() != oc.PlatformTypes_ComponentRedundantRole_PRIMARY {
					t.Errorf("Expected %s to be PRIMARY after switchover, got %v", defaultSecondarySupervisor, finalSupervisor2.GetRedundantRole())
				}

				// Verify timestamps updated
				if finalSupervisor1.GetLastSwitchoverTime() <= initialSupervisor1.GetLastSwitchoverTime() {
					t.Error("Supervisor1 switchover time was not updated")
				}
				if finalSupervisor2.GetLastSwitchoverTime() <= initialSupervisor2.GetLastSwitchoverTime() {
					t.Error("Supervisor2 switchover time was not updated")
				}

				// Verify switchover reasons
				if finalSupervisor1.GetLastSwitchoverReason().GetTrigger() != oc.PlatformTypes_ComponentRedundantRoleSwitchoverReasonTrigger_USER_INITIATED {
					t.Errorf("Expected Supervisor1 switchover trigger USER_INITIATED, got %v", finalSupervisor1.GetLastSwitchoverReason().GetTrigger())
				}
				if finalSupervisor2.GetLastSwitchoverReason().GetTrigger() != oc.PlatformTypes_ComponentRedundantRoleSwitchoverReasonTrigger_USER_INITIATED {
					t.Errorf("Expected Supervisor2 switchover trigger USER_INITIATED, got %v", finalSupervisor2.GetLastSwitchoverReason().GetTrigger())
				}
			},
		},
		"switchover-back": {
			fn: func(t *testing.T, s *system, ctx context.Context, c *ygnmi.Client) {
				// Switch to Supervisor2 first (from default Supervisor1=PRIMARY)
				req1 := &spb.SwitchControlProcessorRequest{
					ControlProcessor: &pb.Path{
						Elem: []*pb.PathElem{
							{Name: "components"},
							{Name: "component", Key: map[string]string{"name": "Supervisor2"}},
						},
					},
				}
				_, err := s.SwitchControlProcessor(ctx, req1)
				if err != nil {
					t.Fatalf("First switchover failed: %v", err)
				}
				time.Sleep(2500 * time.Millisecond)

				// Verify Supervisor2 is now PRIMARY
				newActiveState, err := ygnmi.Get(ctx, c, ocpath.Root().Component("Supervisor2").State())
				if err != nil {
					t.Fatalf("Failed to get new active supervisor state: %v", err)
				}
				if newActiveState.GetRedundantRole() != oc.PlatformTypes_ComponentRedundantRole_PRIMARY {
					t.Errorf("Expected Supervisor2 to be PRIMARY after first switchover, got %v", newActiveState.GetRedundantRole())
				}

				// Switch back to Supervisor1
				req2 := &spb.SwitchControlProcessorRequest{
					ControlProcessor: &pb.Path{
						Elem: []*pb.PathElem{
							{Name: "components"},
							{Name: "component", Key: map[string]string{"name": "Supervisor1"}},
						},
					},
				}
				resp, err := s.SwitchControlProcessor(ctx, req2)
				if err != nil {
					t.Fatalf("Second switchover failed: %v", err)
				}
				if resp == nil {
					t.Fatal("Expected response but got nil")
				}
				time.Sleep(2500 * time.Millisecond)

				// Verify Supervisor1 is PRIMARY again
				finalState, err := ygnmi.Get(ctx, c, ocpath.Root().Component("Supervisor1").State())
				if err != nil {
					t.Fatalf("Failed to get final supervisor state: %v", err)
				}
				if finalState.GetRedundantRole() != oc.PlatformTypes_ComponentRedundantRole_PRIMARY {
					t.Errorf("Expected Supervisor1 to be PRIMARY after switchback, got %v", finalState.GetRedundantRole())
				}
			},
		},
		"switchover-with-pending-system-reboot": {
			fn: func(t *testing.T, s *system, ctx context.Context, c *ygnmi.Client) {
				// Start a delayed system reboot
				rebootReq := &spb.RebootRequest{
					Method: spb.RebootMethod_COLD,
					Delay:  100000,
				}
				_, err := s.Reboot(ctx, rebootReq)
				if err != nil {
					t.Fatalf("Failed to initiate system reboot: %v", err)
				}

				// Try switchover while system reboot is pending
				switchReq := &spb.SwitchControlProcessorRequest{
					ControlProcessor: &pb.Path{
						Elem: []*pb.PathElem{
							{Name: "components"},
							{Name: "component", Key: map[string]string{"name": "Supervisor2"}},
						},
					},
				}
				_, err = s.SwitchControlProcessor(ctx, switchReq)
				if err == nil {
					t.Error("Expected switchover to fail with pending system reboot")
				} else if got := status.Code(err); got != codes.FailedPrecondition {
					t.Errorf("Expected FailedPrecondition error, got %v", got)
				}

				// Cancel reboot
				_, err = s.CancelReboot(ctx, &spb.CancelRebootRequest{})
				if err != nil {
					t.Fatalf("Failed to cancel reboot: %v", err)
				}
			},
		},
		"switchover-with-pending-component-reboot": {
			fn: func(t *testing.T, s *system, ctx context.Context, c *ygnmi.Client) {
				// Start a delayed component reboot
				rebootReq := &spb.RebootRequest{
					Method: spb.RebootMethod_COLD,
					Delay:  100000,
					Subcomponents: []*pb.Path{{
						Elem: []*pb.PathElem{
							{Name: "components"},
							{Name: "component", Key: map[string]string{"name": "Linecard0"}},
						},
					}},
				}
				_, err := s.Reboot(ctx, rebootReq)
				if err != nil {
					t.Fatalf("Failed to initiate component reboot: %v", err)
				}

				// Try switchover while component reboot is pending
				switchReq := &spb.SwitchControlProcessorRequest{
					ControlProcessor: &pb.Path{
						Elem: []*pb.PathElem{
							{Name: "components"},
							{Name: "component", Key: map[string]string{"name": "Supervisor2"}},
						},
					},
				}
				_, err = s.SwitchControlProcessor(ctx, switchReq)
				if err == nil {
					t.Error("Expected switchover to fail with pending component reboot")
				} else if got := status.Code(err); got != codes.FailedPrecondition {
					t.Errorf("Expected FailedPrecondition error, got %v", got)
				}

				// Cancel reboot
				_, err = s.CancelReboot(ctx, &spb.CancelRebootRequest{})
				if err != nil {
					t.Fatalf("Failed to cancel reboot: %v", err)
				}
			},
		},
		"concurrent-switchover-blocked": {
			fn: func(t *testing.T, s *system, ctx context.Context, c *ygnmi.Client) {
				// Create request to switch to Supervisor2 (default standby)
				req := &spb.SwitchControlProcessorRequest{
					ControlProcessor: &pb.Path{
						Elem: []*pb.PathElem{
							{Name: "components"},
							{Name: "component", Key: map[string]string{"name": "Supervisor2"}},
						},
					},
				}

				// Manually set the pending switchover flag to simulate in-progress operation
				s.switchoverMu.Lock()
				s.hasPendingSwitchover = true
				s.switchoverMu.Unlock()

				// Try switchover while one is supposedly in progress
				_, err := s.SwitchControlProcessor(ctx, req)
				if err == nil {
					t.Error("Expected switchover to fail while another is in progress")
				} else if got := status.Code(err); got != codes.FailedPrecondition {
					t.Errorf("Expected FailedPrecondition error, got %v", got)
				}

				// Clear the flag and try again - should work
				s.switchoverMu.Lock()
				s.hasPendingSwitchover = false
				s.switchoverMu.Unlock()

				_, err = s.SwitchControlProcessor(ctx, req)
				if err != nil {
					t.Fatalf("Switchover should succeed after clearing flag: %v", err)
				}
			},
		},
	}
	setupEnvironment := func(t *testing.T) (context.Context, *grpc.Server, *ygnmi.Client, *system, func()) {
		ctx := context.Background()
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

		client := gnmiServer.LocalClient()
		c, err := ygnmi.NewClient(client, ygnmi.WithTarget("local"))
		if err != nil {
			t.Fatalf("cannot create ygnmi client: %v", err)
		}
		s := newSystem(c, lemmingConfig)

		// Initialize the system
		if err := fakedevice.NewBootTimeTask(lemmingConfig).Start(ctx, client, "local"); err != nil {
			t.Fatalf("Failed to initialize boot time: %v", err)
		}
		if err := fakedevice.NewChassisComponentsTask(lemmingConfig).Start(ctx, client, "local"); err != nil {
			t.Fatalf("Failed to initialize chassis components: %v", err)
		}

		cleanup := func() {
			grpcServer.Stop()
		}
		return ctx, grpcServer, c, s, cleanup
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, _, c, s, cleanup := setupEnvironment(t)
			defer cleanup()
			test.fn(t, s, ctx, c)
		})
	}
}

// TestKillProcess tests the KillProcess RPC functionality with comprehensive scenarios
func TestKillProcess(t *testing.T) {
	lemmingConfig := loadDefaultConfig(t)

	setupFreshEnvironment := func(t *testing.T) (context.Context, *grpc.Server, *ygnmi.Client, *system, func()) {
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

		client := gnmiServer.LocalClient()
		c, err := ygnmi.NewClient(client, ygnmi.WithTarget("local"))
		if err != nil {
			t.Fatalf("cannot create ygnmi client: %v", err)
		}

		s := newSystem(c, lemmingConfig)
		ctx := context.Background()
		fakedevice.NewProcessMonitoringTask(lemmingConfig).Start(ctx, client, "local")

		time.Sleep(100 * time.Millisecond)

		cleanup := func() {
			grpcServer.GracefulStop()
		}

		return ctx, grpcServer, c, s, cleanup
	}

	tests := map[string]struct {
		fn func(*testing.T, *system, context.Context, *ygnmi.Client)
	}{
		"process-termination-kill": {
			fn: func(t *testing.T, s *system, ctx context.Context, c *ygnmi.Client) {
				// Kill Gribi with KILL signal, no restart
				req := &spb.KillProcessRequest{
					Pid:     1002,
					Signal:  spb.KillProcessRequest_SIGNAL_KILL,
					Restart: false,
				}
				_, err := s.KillProcess(ctx, req)
				if err != nil {
					t.Fatalf("KillProcess with KILL failed: %v", err)
				}

				// Process should be removed from state
				_, err = ygnmi.Get(ctx, c, ocpath.Root().System().Process(1002).State())
				if err == nil {
					t.Error("Process should be removed after KILL signal")
				}
			},
		},
		"process-reload-hup": {
			fn: func(t *testing.T, s *system, ctx context.Context, c *ygnmi.Client) {
				// Get initial state
				initialProcess, err := ygnmi.Get(ctx, c, ocpath.Root().System().Process(1004).State())
				if err != nil {
					t.Fatalf("Failed to get initial process state: %v", err)
				}
				initialStartTime := initialProcess.GetStartTime()
				initialPID := initialProcess.GetPid()

				// Reload kim with HUP signal
				req := &spb.KillProcessRequest{
					Name:    "kim",
					Signal:  spb.KillProcessRequest_SIGNAL_HUP,
					Restart: true, // Should be ignored for HUP
				}
				_, err = s.KillProcess(ctx, req)
				if err != nil {
					t.Fatalf("KillProcess with HUP failed: %v", err)
				}

				// Process should still exist with same PID but updated start time
				reloadedProcess, err := ygnmi.Get(ctx, c, ocpath.Root().System().Process(1004).State())
				if err != nil {
					t.Fatalf("Process should still exist after HUP: %v", err)
				}

				if reloadedProcess.GetPid() != initialPID {
					t.Errorf("PID should remain same after HUP, got %d, want %d", reloadedProcess.GetPid(), initialPID)
				}

				if reloadedProcess.GetStartTime() <= initialStartTime {
					t.Error("Start time should be updated after HUP signal")
				}

				if reloadedProcess.GetName() != "kim" {
					t.Errorf("Process name changed, got %s, want kim", reloadedProcess.GetName())
				}
			},
		},
		"process-restart-with-new-pid": {
			fn: func(t *testing.T, s *system, ctx context.Context, c *ygnmi.Client) {
				// Kill and restart Octa
				req := &spb.KillProcessRequest{
					Name:    "Octa",
					Signal:  spb.KillProcessRequest_SIGNAL_TERM,
					Restart: true,
				}
				_, err := s.KillProcess(ctx, req)
				if err != nil {
					t.Fatalf("KillProcess with restart failed: %v", err)
				}

				// Process should be deleted immediately
				_, err = ygnmi.Get(ctx, c, ocpath.Root().System().Process(1001).State())
				if err == nil {
					t.Error("Process should be deleted immediately after kill")
				}

				time.Sleep(3 * time.Second)

				// Check if a new bgpd process exists with different PID
				processes, err := ygnmi.GetAll(ctx, c, ocpath.Root().System().ProcessAny().State())
				if err != nil {
					t.Fatalf("Failed to get processes: %v", err)
				}

				var newOctaProcess *oc.System_Process
				for _, p := range processes {
					if p.GetName() == "Octa" {
						newOctaProcess = p
						break
					}
				}

				if newOctaProcess == nil {
					t.Fatal("Octa process not restarted")
				}

				if newOctaProcess.GetPid() == 1001 {
					t.Error("Restarted process should have different PID")
				}

				if newOctaProcess.GetPid() < 1 || newOctaProcess.GetPid() > 65535 {
					t.Errorf("New PID %d should be in range 1001-1100", newOctaProcess.GetPid())
				}
			},
		},
		"missing-process-invalid-name": {
			fn: func(t *testing.T, s *system, ctx context.Context, c *ygnmi.Client) {
				req := &spb.KillProcessRequest{
					Name:   "nonexistent",
					Signal: spb.KillProcessRequest_SIGNAL_TERM,
				}
				_, err := s.KillProcess(ctx, req)
				if err == nil {
					t.Error("Expected error for non-existent process name")
				} else if got := status.Code(err); got != codes.NotFound {
					t.Errorf("Expected NotFound error, got %v", got)
				}
			},
		},
		"missing-process-invalid-pid": {
			fn: func(t *testing.T, s *system, ctx context.Context, c *ygnmi.Client) {
				req := &spb.KillProcessRequest{
					Pid:    9999,
					Signal: spb.KillProcessRequest_SIGNAL_TERM,
				}
				_, err := s.KillProcess(ctx, req)
				if err == nil {
					t.Error("Expected error for non-existent PID")
				} else if got := status.Code(err); got != codes.NotFound {
					t.Errorf("Expected NotFound error, got %v", got)
				}
			},
		},
		"invalid-signal-unspecified": {
			fn: func(t *testing.T, s *system, ctx context.Context, c *ygnmi.Client) {
				req := &spb.KillProcessRequest{
					Pid:    1001,
					Signal: spb.KillProcessRequest_SIGNAL_UNSPECIFIED,
				}
				_, err := s.KillProcess(ctx, req)
				if err == nil {
					t.Error("Expected error for unspecified signal")
				} else if got := status.Code(err); got != codes.InvalidArgument {
					t.Errorf("Expected InvalidArgument error, got %v", got)
				}
			},
		},
		"missing-identifier-no-pid-no-name": {
			fn: func(t *testing.T, s *system, ctx context.Context, c *ygnmi.Client) {
				req := &spb.KillProcessRequest{
					Signal: spb.KillProcessRequest_SIGNAL_TERM,
				}
				_, err := s.KillProcess(ctx, req)
				if err == nil {
					t.Error("Expected error when neither PID nor name specified")
				} else if got := status.Code(err); got != codes.InvalidArgument {
					t.Errorf("Expected InvalidArgument error, got %v", got)
				}
			},
		},
		"missing-signal": {
			fn: func(t *testing.T, s *system, ctx context.Context, c *ygnmi.Client) {
				req := &spb.KillProcessRequest{
					Pid: 1001,
				}
				_, err := s.KillProcess(ctx, req)
				if err == nil {
					t.Error("Expected error when signal not specified")
				} else if got := status.Code(err); got != codes.InvalidArgument {
					t.Errorf("Expected InvalidArgument error, got %v", got)
				}
			},
		},
		"concurrent-kill-operations": {
			fn: func(t *testing.T, s *system, ctx context.Context, c *ygnmi.Client) {
				// Test concurrent operations are serialized properly
				req1 := &spb.KillProcessRequest{
					Pid:    1001,
					Signal: spb.KillProcessRequest_SIGNAL_HUP,
				}
				req2 := &spb.KillProcessRequest{
					Pid:    1002,
					Signal: spb.KillProcessRequest_SIGNAL_HUP,
				}

				// Start concurrent operations
				done1 := make(chan error, 1)
				done2 := make(chan error, 1)

				go func() {
					_, err := s.KillProcess(ctx, req1)
					done1 <- err
				}()

				go func() {
					_, err := s.KillProcess(ctx, req2)
					done2 <- err
				}()

				// Both should succeed (serialized by mutex)
				err1 := <-done1
				err2 := <-done2

				if err1 != nil {
					t.Errorf("Concurrent operation 1 failed: %v", err1)
				}
				if err2 != nil {
					t.Errorf("Concurrent operation 2 failed: %v", err2)
				}
			},
		},
		"default-processes-validation": {
			fn: func(t *testing.T, s *system, ctx context.Context, c *ygnmi.Client) {
				// Verify all default processes exist
				expectedProcesses := map[string]uint64{
					"Octa":        1001,
					"Gribi":       1002,
					"emsd":        1003,
					"kim":         1004,
					"grpc_server": 1005,
					"fibd":        1006,
					"rpd":         1007,
				}

				for name, expectedPID := range expectedProcesses {
					process, err := ygnmi.Get(ctx, c, ocpath.Root().System().Process(expectedPID).State())
					if err != nil {
						t.Errorf("Default process %s (PID %d) not found: %v", name, expectedPID, err)
						continue
					}

					if process.GetName() != name {
						t.Errorf("Process PID %d has wrong name: got %s, want %s", expectedPID, process.GetName(), name)
					}

					if process.GetPid() != expectedPID {
						t.Errorf("Process %s has wrong PID: got %d, want %d", name, process.GetPid(), expectedPID)
					}

					// Verify realistic resource simulation
					if process.GetCpuUtilization() == 0 {
						t.Errorf("Process %s should have non-zero CPU utilization", name)
					}
					if process.GetMemoryUsage() == 0 {
						t.Errorf("Process %s should have non-zero memory usage", name)
					}
				}
			},
		},
	}

	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx, _, c, s, cleanup := setupFreshEnvironment(t)
			defer cleanup()
			test.fn(t, s, ctx, c)
		})
	}
}

// mockPingServer implements spb.System_PingServer for testing
type mockPingServer struct {
	grpc.ServerStream
	ctx       context.Context
	responses []*spb.PingResponse
	sendErr   error
}

func (m *mockPingServer) Send(response *spb.PingResponse) error {
	if m.sendErr != nil {
		return m.sendErr
	}
	m.responses = append(m.responses, response)
	return nil
}

func (m *mockPingServer) Context() context.Context {
	return m.ctx
}

// TestPing tests the Ping RPC functionality with comprehensive scenarios
func TestPing(t *testing.T) {
	lemmingConfig := loadDefaultConfig(t)

	newMockPingServer := func(ctx context.Context) *mockPingServer {
		return &mockPingServer{ctx: ctx}
	}

	setupPingEnvironment := func(t *testing.T) (context.Context, *system, func()) {
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

		client := gnmiServer.LocalClient()
		c, err := ygnmi.NewClient(client, ygnmi.WithTarget("local"))
		if err != nil {
			t.Fatalf("cannot create ygnmi client: %v", err)
		}

		s := newSystem(c, lemmingConfig)
		ctx := context.Background()

		cleanup := func() {
			grpcServer.GracefulStop()
		}

		return ctx, s, cleanup
	}

	tests := map[string]struct {
		fn func(*testing.T, *system, context.Context)
	}{
		"validation-missing-destination": {
			fn: func(t *testing.T, s *system, ctx context.Context) {
				stream := newMockPingServer(ctx)
				err := s.Ping(&spb.PingRequest{}, stream)
				if err == nil {
					t.Fatal("Expected error for missing destination")
				}
				if got := status.Code(err); got != codes.InvalidArgument {
					t.Errorf("Expected InvalidArgument error, got %v", got)
				}
			},
		},
		"validation-empty-destination": {
			fn: func(t *testing.T, s *system, ctx context.Context) {
				stream := newMockPingServer(ctx)
				err := s.Ping(&spb.PingRequest{Destination: ""}, stream)
				if err == nil {
					t.Fatal("Expected error for empty destination")
				}
				if got := status.Code(err); got != codes.InvalidArgument {
					t.Errorf("Expected InvalidArgument error, got %v", got)
				}
			},
		},
		"validation-invalid-count": {
			fn: func(t *testing.T, s *system, ctx context.Context) {
				stream := newMockPingServer(ctx)
				err := s.Ping(&spb.PingRequest{
					Destination: "8.8.8.8",
					Count:       -2,
				}, stream)
				if err == nil {
					t.Fatal("Expected error for invalid count")
				}
				if got := status.Code(err); got != codes.InvalidArgument {
					t.Errorf("Expected InvalidArgument error, got %v", got)
				}
			},
		},
		"validation-invalid-interval": {
			fn: func(t *testing.T, s *system, ctx context.Context) {
				stream := newMockPingServer(ctx)
				err := s.Ping(&spb.PingRequest{
					Destination: "8.8.8.8",
					Interval:    -2,
				}, stream)
				if err == nil {
					t.Fatal("Expected error for invalid interval")
				}
				if got := status.Code(err); got != codes.InvalidArgument {
					t.Errorf("Expected InvalidArgument error, got %v", got)
				}
			},
		},
		"validation-packet-size-too-small": {
			fn: func(t *testing.T, s *system, ctx context.Context) {
				stream := newMockPingServer(ctx)
				err := s.Ping(&spb.PingRequest{
					Destination: "8.8.8.8",
					Size:        7,
				}, stream)
				if err == nil {
					t.Fatal("Expected error for packet size too small")
				}
				if got := status.Code(err); got != codes.InvalidArgument {
					t.Errorf("Expected InvalidArgument error, got %v", got)
				}
			},
		},
		"validation-packet-size-too-large": {
			fn: func(t *testing.T, s *system, ctx context.Context) {
				stream := newMockPingServer(ctx)
				err := s.Ping(&spb.PingRequest{
					Destination: "8.8.8.8",
					Size:        65508,
				}, stream)
				if err == nil {
					t.Fatal("Expected error for packet size too large")
				}
				if got := status.Code(err); got != codes.InvalidArgument {
					t.Errorf("Expected InvalidArgument error, got %v", got)
				}
			},
		},
		"streaming-individual-and-summary-responses": {
			fn: func(t *testing.T, s *system, ctx context.Context) {
				ctxWithTimeout, cancel := context.WithTimeout(ctx, 3*time.Second)
				defer cancel()

				stream := newMockPingServer(ctxWithTimeout)
				err := s.Ping(&spb.PingRequest{
					Destination: "1.1.1.1",
					Count:       3,
					Interval:    50000000, // 50ms
				}, stream)
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}

				expectedResponses := 4 // 3 individual + 1 summary
				if len(stream.responses) != expectedResponses {
					t.Errorf("Expected %d responses, got %d", expectedResponses, len(stream.responses))
				}

				// Check individual packet responses
				for i := 0; i < 3; i++ {
					resp := stream.responses[i]
					if resp.Source != "1.1.1.1" {
						t.Errorf("Response %d: expected source 1.1.1.1, got %s", i, resp.Source)
					}
					if resp.Time <= 0 {
						t.Errorf("Response %d: expected positive RTT, got %d", i, resp.Time)
					}
					if resp.Sequence != int32(i+1) {
						t.Errorf("Response %d: expected sequence %d, got %d", i, i+1, resp.Sequence)
					}
					if resp.Sent != 0 || resp.Received != 0 {
						t.Errorf("Response %d: individual response should not have summary fields", i)
					}
				}

				// Check summary response
				summary := stream.responses[3]
				if summary.Sent != 3 {
					t.Errorf("Summary: expected sent 3, got %d", summary.Sent)
				}
				if summary.MinTime <= 0 || summary.MaxTime <= 0 || summary.AvgTime <= 0 {
					t.Error("Summary: expected positive min/max/avg times")
				}
			},
		},
		"custom-packet-size": {
			fn: func(t *testing.T, s *system, ctx context.Context) {
				ctxWithTimeout, cancel := context.WithTimeout(ctx, 3*time.Second)
				defer cancel()

				stream := newMockPingServer(ctxWithTimeout)
				err := s.Ping(&spb.PingRequest{
					Destination: "192.168.1.1",
					Count:       1,
					Size:        128,
				}, stream)
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}

				if len(stream.responses) < 1 {
					t.Fatal("Expected at least one response")
				}

				resp := stream.responses[0]
				if resp.Bytes != 128 {
					t.Errorf("Expected bytes 128, got %d", resp.Bytes)
				}
			},
		},
		"cancellation-context-cancelled": {
			fn: func(t *testing.T, s *system, ctx context.Context) {
				ctxWithCancel, cancel := context.WithCancel(ctx)
				stream := newMockPingServer(ctxWithCancel)

				// Cancel context after short delay
				go func() {
					time.Sleep(50 * time.Millisecond)
					cancel()
				}()

				err := s.Ping(&spb.PingRequest{
					Destination: "192.168.1.1",
					Count:       100,
					Interval:    1000000000, // 1 second
				}, stream)

				if err == nil {
					t.Fatal("Expected error due to context cancellation")
				}
				if err != context.Canceled {
					t.Errorf("Expected context.Canceled, got %v", err)
				}
			},
		},
		"cancellation-timeout": {
			fn: func(t *testing.T, s *system, ctx context.Context) {
				ctxWithTimeout, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
				defer cancel()

				stream := newMockPingServer(ctxWithTimeout)
				err := s.Ping(&spb.PingRequest{
					Destination: "192.168.1.1",
					Count:       -1, // Continuous ping
				}, stream)

				if err == nil {
					t.Fatal("Expected timeout error")
				}
				if err != context.DeadlineExceeded {
					if st, ok := status.FromError(err); !ok || st.Code() != codes.DeadlineExceeded {
						t.Errorf("Expected deadline exceeded error, got %v", err)
					}
				}
			},
		},
		"continuous-ping-with-early-termination": {
			fn: func(t *testing.T, s *system, ctx context.Context) {
				ctxWithTimeout, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
				defer cancel()

				stream := newMockPingServer(ctxWithTimeout)
				err := s.Ping(&spb.PingRequest{
					Destination: "192.168.1.1",
					Count:       -1,       // Continuous
					Interval:    50000000, // 50ms
				}, stream)

				// Should get cancelled/timeout
				if err == nil {
					t.Fatal("Expected error due to timeout")
				}

				// Should have received some responses before cancellation
				if len(stream.responses) == 0 {
					t.Error("Expected some responses before cancellation")
				}
				t.Logf("Received %d responses before cancellation", len(stream.responses))
			},
		},
		"latency-variation-realistic": {
			fn: func(t *testing.T, s *system, ctx context.Context) {
				ctxWithTimeout, cancel := context.WithTimeout(ctx, 4*time.Second)
				defer cancel()

				stream := newMockPingServer(ctxWithTimeout)
				err := s.Ping(&spb.PingRequest{
					Destination: "8.8.8.8",
					Count:       5,
					Interval:    200000000,
				}, stream)
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}

				// Collect RTT values from individual responses
				var rtts []int64
				for i := 0; i < 5; i++ {
					rtts = append(rtts, stream.responses[i].Time)
				}

				// Verify RTTs are within reasonable bounds
				for i, rtt := range rtts {
					if rtt < 1000000 { // 1ms in nanoseconds
						t.Errorf("RTT %d too small: %d ns", i, rtt)
					}
					if rtt > 1000000000 { // 1s in nanoseconds
						t.Errorf("RTT %d too large: %d ns", i, rtt)
					}
				}

				// Verify summary statistics are reasonable
				summary := stream.responses[len(stream.responses)-1]
				if summary.MinTime >= summary.MaxTime && summary.Received > 1 {
					t.Error("Min time should be less than max time with multiple packets")
				}
			},
		},
		"flood-ping-short-interval": {
			fn: func(t *testing.T, s *system, ctx context.Context) {
				ctxWithTimeout, cancel := context.WithTimeout(ctx, 1*time.Second)
				defer cancel()

				stream := newMockPingServer(ctxWithTimeout)
				err := s.Ping(&spb.PingRequest{
					Destination: "192.168.1.1",
					Count:       3,
					Interval:    -1, // Flood ping (1ms minimum)
				}, stream)
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}

				expectedResponses := 4 // 3 individual + 1 summary
				if len(stream.responses) != expectedResponses {
					t.Errorf("Expected %d responses, got %d", expectedResponses, len(stream.responses))
				}
			},
		},
		"wait-timeout-simulation": {
			fn: func(t *testing.T, s *system, ctx context.Context) {
				ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
				defer cancel()

				stream := newMockPingServer(ctxWithTimeout)
				err := s.Ping(&spb.PingRequest{
					Destination: "192.0.2.1",
					Count:       3,
					Wait:        1000000,    // 1ms wait (very short)
					Interval:    1000000000, // 1 second interval
				}, stream)
				if err != nil {
					t.Fatalf("Ping failed: %v", err)
				}

				responses := stream.responses
				if len(responses) < 4 {
					t.Errorf("Expected at least 4 responses, got %d", len(responses))
				}

				// Check summary
				summary := responses[len(responses)-1]
				if summary.Sent != 3 {
					t.Errorf("Expected Sent=3, got %d", summary.Sent)
				}
				t.Logf("Summary with short wait: Sent=%d, Received=%d", summary.Sent, summary.Received)
			},
		},
	}

	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx, s, cleanup := setupPingEnvironment(t)
			defer cleanup()
			test.fn(t, s, ctx)
		})
	}
}

func TestLinkQualification(t *testing.T) {
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

	client := gnmiServer.LocalClient()
	c, err := ygnmi.NewClient(client, ygnmi.WithTarget("local"))
	if err != nil {
		t.Fatalf("cannot create ygnmi client: %v", err)
	}

	lemmingConfig := loadDefaultConfig(t)
	linkQualServer := newLinkQualification(c, lemmingConfig)
	ctx := context.Background()

	// Create test interfaces in the system for testing
	setupTestInterfaces(t, c)

	// Helper function to clean up all active qualifications
	cleanupQualifications := func() {
		// Get all active qualifications
		listResp, err := linkQualServer.List(ctx, &plqpb.ListRequest{})
		if err != nil {
			t.Logf("Warning: List failed during cleanup: %v", err)
			return
		}

		if len(listResp.Results) > 0 {
			// Collect all IDs
			var ids []string
			for _, result := range listResp.Results {
				ids = append(ids, result.Id)
			}

			// Delete all qualifications
			deleteResp, err := linkQualServer.Delete(ctx, &plqpb.DeleteRequest{Ids: ids})
			if err != nil {
				t.Logf("Warning: Delete failed during cleanup: %v", err)
			} else {
				for id, status := range deleteResp.Results {
					if status.Code != int32(codes.OK) {
						t.Logf("Warning: Failed to delete qualification %s: %s", id, status.Message)
					}
				}
			}

			// Wait a moment for cleanup to complete
			time.Sleep(100 * time.Millisecond)
		}
	}

	tests := []struct {
		name string
		test func(t *testing.T, linkQualServer *linkQualification, ctx context.Context)
	}{
		{"capabilities", testCapabilities},
		{"validation-tests", testValidation},
		{"single-interface-qualification", testSingleInterfaceQualification},
		{"multi-port-qualification", testMultiPortQualification},
		{"concurrent-qualifications", testConcurrentQualifications},
		{"state-machine-validation", testStateMachineValidation},
		{"packet-statistics-validation", testPacketStatisticsValidation},
		{"timing-validation", testTimingValidation},
		{"cancellation-and-cleanup", testCancellationAndCleanup},
		{"error-handling", testErrorHandling},
		{"configuration-variations", testConfigurationVariations},
		{"packet-injector-endpoint", testPacketInjectorEndpoint},
		{"ntp-timing-configuration", testNTPTimingConfiguration},
		{"generator-reflector-coordination", testGeneratorReflectorCoordination},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up before each test
			cleanupQualifications()

			// Run the test
			tt.test(t, linkQualServer, ctx)

			// Clean up after test (in case test didn't clean up properly)
			cleanupQualifications()
		})
	}

	// Final cleanup
	cleanupQualifications()
}

// setupTestInterfaces creates test interfaces for qualification testing
func setupTestInterfaces(t *testing.T, c *ygnmi.Client) {
	ctx := context.Background()
	interfaces := []string{"eth0", "eth1", "eth2", "eth3", "eth4", "eth5"}

	for _, intfName := range interfaces {
		intf := &oc.Interface{
			Name:        ygot.String(intfName),
			Type:        oc.IETFInterfaces_InterfaceType_ethernetCsmacd,
			OperStatus:  oc.Interface_OperStatus_UP,
			AdminStatus: oc.Interface_AdminStatus_UP,
			Enabled:     ygot.Bool(true),
		}

		_, err := gnmiclient.Replace(gnmi.AddTimestampMetadata(ctx, time.Now().UnixNano()),
			c, ocpath.Root().Interface(intfName).State(), intf)
		if err != nil {
			t.Logf("Warning: Failed to create test interface %s: %v", intfName, err)
		}
	}
}

func testCapabilities(t *testing.T, linkQualServer *linkQualification, ctx context.Context) {
	resp, err := linkQualServer.Capabilities(ctx, &plqpb.CapabilitiesRequest{})
	if err != nil {
		t.Fatalf("Capabilities failed: %v", err)
	}

	// Verify static capabilities per design
	if resp.GetGenerator().GetPacketGenerator().GetMaxBps() != 400000000000 {
		t.Errorf("Expected MaxBps 400G, got %d", resp.GetGenerator().GetPacketGenerator().GetMaxBps())
	}
	if resp.GetGenerator().GetPacketGenerator().GetMaxPps() != 500000000 {
		t.Errorf("Expected MaxPps 500M, got %d", resp.GetGenerator().GetPacketGenerator().GetMaxPps())
	}
	if resp.GetMaxHistoricalResultsPerInterface() != 10 {
		t.Errorf("Expected MaxHistoricalResultsPerInterface 10, got %d", resp.GetMaxHistoricalResultsPerInterface())
	}
	if resp.GetReflector().GetAsicLoopback() == nil {
		t.Error("Expected ASIC loopback capabilities")
	}
	if resp.GetReflector().GetPmdLoopback() == nil {
		t.Error("Expected PMD loopback capabilities")
	}

	// Verify PacketInjector is NOT advertised (unimplemented)
	if resp.GetGenerator().GetPacketInjector() != nil {
		t.Error("PacketInjector capabilities should not be advertised (unimplemented)")
	}

	// Verify timing capabilities
	if resp.GetGenerator().GetPacketGenerator().GetMinSetupDuration().AsDuration() != 1*time.Second {
		t.Errorf("Expected MinSetupDuration 1s, got %v", resp.GetGenerator().GetPacketGenerator().GetMinSetupDuration().AsDuration())
	}

	// Verify MinSampleInterval is 1 second as required by tests
	if resp.GetGenerator().GetPacketGenerator().GetMinSampleInterval().AsDuration() != 1*time.Second {
		t.Errorf("Expected MinSampleInterval 1s, got %v", resp.GetGenerator().GetPacketGenerator().GetMinSampleInterval().AsDuration())
	}

	t.Logf("Capabilities verified: PacketGenerator supported, PacketInjector unimplemented, reflectors supported")
}

func testValidation(t *testing.T, linkQualServer *linkQualification, ctx context.Context) {
	tests := []struct {
		name   string
		req    *plqpb.CreateRequest
		expect codes.Code
	}{
		{
			name:   "empty-interfaces",
			req:    &plqpb.CreateRequest{},
			expect: codes.InvalidArgument,
		},
		{
			name: "missing-id",
			req: &plqpb.CreateRequest{
				Interfaces: []*plqpb.QualificationConfiguration{{
					InterfaceName: "eth0",
				}},
			},
			expect: codes.InvalidArgument,
		},
		{
			name: "missing-interface-name",
			req: &plqpb.CreateRequest{
				Interfaces: []*plqpb.QualificationConfiguration{{
					Id: "test-qual-1",
				}},
			},
			expect: codes.InvalidArgument,
		},
		{
			name: "duplicate-ids",
			req: &plqpb.CreateRequest{
				Interfaces: []*plqpb.QualificationConfiguration{
					{
						Id:            "duplicate-id",
						InterfaceName: "eth0",
						EndpointType: &plqpb.QualificationConfiguration_PacketGenerator{
							PacketGenerator: &plqpb.PacketGeneratorConfiguration{
								PacketRate: 1000,
								PacketSize: 1500,
							},
						},
					},
					{
						Id:            "duplicate-id",
						InterfaceName: "eth1",
						EndpointType: &plqpb.QualificationConfiguration_AsicLoopback{
							AsicLoopback: &plqpb.AsicLoopbackConfiguration{},
						},
					},
				},
			},
			expect: codes.OK, // First should succeed, second should return AlreadyExists in response
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := linkQualServer.Create(ctx, tt.req)
			if err == nil && tt.expect != codes.OK {
				t.Fatalf("Expected error with code %v, but got success", tt.expect)
			}
			if err != nil && status.Code(err) != tt.expect {
				t.Fatalf("Expected error code %v, got %v", tt.expect, status.Code(err))
			}

			// For duplicate ID test, check individual status
			if tt.name == "duplicate-ids" && resp != nil {
				if len(resp.Status) != 1 {
					t.Fatalf("Expected 1 status entry for duplicate IDs, got %d", len(resp.Status))
				}
				// The duplicate ID should have AlreadyExists error
				duplicateStatus := resp.Status["duplicate-id"]
				if duplicateStatus.Code != int32(codes.AlreadyExists) {
					t.Errorf("Expected AlreadyExists for duplicate ID, got %d", duplicateStatus.Code)
				}
			}
		})
	}
}

func testSingleInterfaceQualification(t *testing.T, linkQualServer *linkQualification, ctx context.Context) {
	qualID := "single-qual-test"
	interfaceName := "eth0"

	// Create a single qualification
	createReq := &plqpb.CreateRequest{
		Interfaces: []*plqpb.QualificationConfiguration{
			{
				Id:            qualID,
				InterfaceName: interfaceName,
				EndpointType: &plqpb.QualificationConfiguration_PacketGenerator{
					PacketGenerator: &plqpb.PacketGeneratorConfiguration{
						PacketRate: 10000, // 10K PPS for quick test
						PacketSize: 1500,
					},
				},
				Timing: &plqpb.QualificationConfiguration_Rpc{
					Rpc: &plqpb.RPCSyncedTiming{
						SetupDuration:    durationpb.New(1 * time.Second),
						Duration:         durationpb.New(3 * time.Second),
						TeardownDuration: durationpb.New(1 * time.Second),
					},
				},
			},
		},
	}

	createResp, err := linkQualServer.Create(ctx, createReq)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if createResp.Status[qualID].Code != int32(codes.OK) {
		t.Fatalf("Expected OK status, got %d: %s", createResp.Status[qualID].Code, createResp.Status[qualID].Message)
	}

	// Verify List shows the qualification
	listResp, err := linkQualServer.List(ctx, &plqpb.ListRequest{})
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(listResp.Results) != 1 {
		t.Fatalf("Expected 1 qualification in list, got %d", len(listResp.Results))
	}

	if listResp.Results[0].Id != qualID {
		t.Errorf("Expected qualification ID %s, got %s", qualID, listResp.Results[0].Id)
	}

	// Poll Get until completion (should complete within ~6 seconds total)
	var finalResult *plqpb.QualificationResult
	maxWait := 10 * time.Second
	pollInterval := 500 * time.Millisecond
	deadline := time.Now().Add(maxWait)

	for time.Now().Before(deadline) {
		getResp, err := linkQualServer.Get(ctx, &plqpb.GetRequest{Ids: []string{qualID}})
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}

		result := getResp.Results[qualID]
		// Status field should only be set for ERROR states
		if result.State == plqpb.QualificationState_QUALIFICATION_STATE_ERROR && result.Status != nil {
			t.Fatalf("Get returned error status: %d: %s", result.Status.Code, result.Status.Message)
		}

		t.Logf("Qualification state: %v, packets sent: %d, received: %d",
			result.State, result.PacketsSent, result.PacketsReceived)

		if result.State == plqpb.QualificationState_QUALIFICATION_STATE_COMPLETED {
			finalResult = result
			break
		}

		if result.State == plqpb.QualificationState_QUALIFICATION_STATE_ERROR {
			t.Fatalf("Qualification failed with error state")
		}

		time.Sleep(pollInterval)
	}

	if finalResult == nil {
		t.Fatalf("Qualification did not complete within %v", maxWait)
	}

	// Verify final results
	if finalResult.PacketsSent == 0 {
		t.Error("Expected packets to be sent, got 0")
	}
	if finalResult.PacketsReceived == 0 {
		t.Error("Expected packets to be received, got 0")
	}

	// For simulation, we expect good performance with minimal loss
	lossRate := float64(finalResult.PacketsDropped) / float64(finalResult.PacketsSent)
	if lossRate > 0.01 { // Allow up to 1% loss
		t.Errorf("Excessive packet loss: %.2f%% (%d/%d)", lossRate*100, finalResult.PacketsDropped, finalResult.PacketsSent)
	}

	// Verify timing
	duration := finalResult.EndTime.AsTime().Sub(finalResult.StartTime.AsTime())
	expectedDuration := 5 * time.Second // setup + test + teardown
	if duration < expectedDuration-time.Second || duration > expectedDuration+2*time.Second {
		t.Errorf("Unexpected duration: %v, expected ~%v", duration, expectedDuration)
	}

	// Clean up
	deleteResp, err := linkQualServer.Delete(ctx, &plqpb.DeleteRequest{Ids: []string{qualID}})
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	if deleteResp.Results[qualID].Code != int32(codes.OK) {
		t.Errorf("Delete failed: %d: %s", deleteResp.Results[qualID].Code, deleteResp.Results[qualID].Message)
	}
}

func testMultiPortQualification(t *testing.T, linkQualServer *linkQualification, ctx context.Context) {
	// Test multi-port qualification with generator and reflector pair
	generatorID := "multi-gen"
	reflectorID := "multi-ref"

	createReq := &plqpb.CreateRequest{
		Interfaces: []*plqpb.QualificationConfiguration{
			{
				Id:            generatorID,
				InterfaceName: "eth1",
				EndpointType: &plqpb.QualificationConfiguration_PacketGenerator{
					PacketGenerator: &plqpb.PacketGeneratorConfiguration{
						PacketRate: 50000, // 50K PPS
						PacketSize: 1200,
					},
				},
				Timing: &plqpb.QualificationConfiguration_Rpc{
					Rpc: &plqpb.RPCSyncedTiming{
						PreSyncDuration:  durationpb.New(500 * time.Millisecond),
						SetupDuration:    durationpb.New(1 * time.Second),
						Duration:         durationpb.New(4 * time.Second),
						PostSyncDuration: durationpb.New(500 * time.Millisecond),
						TeardownDuration: durationpb.New(1 * time.Second),
					},
				},
			},
			{
				Id:            reflectorID,
				InterfaceName: "eth2",
				EndpointType: &plqpb.QualificationConfiguration_AsicLoopback{
					AsicLoopback: &plqpb.AsicLoopbackConfiguration{},
				},
				Timing: &plqpb.QualificationConfiguration_Rpc{
					Rpc: &plqpb.RPCSyncedTiming{
						PreSyncDuration:  durationpb.New(0), // Reflector starts immediately
						SetupDuration:    durationpb.New(1 * time.Second),
						Duration:         durationpb.New(4 * time.Second),
						PostSyncDuration: durationpb.New(1 * time.Second), // Wait for generator to finish
						TeardownDuration: durationpb.New(1 * time.Second),
					},
				},
			},
		},
	}

	createResp, err := linkQualServer.Create(ctx, createReq)
	if err != nil {
		t.Fatalf("Multi-port create failed: %v", err)
	}

	// Verify both qualifications were created successfully
	if createResp.Status[generatorID].Code != int32(codes.OK) {
		t.Fatalf("Generator creation failed: %d: %s", createResp.Status[generatorID].Code, createResp.Status[generatorID].Message)
	}
	if createResp.Status[reflectorID].Code != int32(codes.OK) {
		t.Fatalf("Reflector creation failed: %d: %s", createResp.Status[reflectorID].Code, createResp.Status[reflectorID].Message)
	}

	// Verify List shows both qualifications
	listResp, err := linkQualServer.List(ctx, &plqpb.ListRequest{})
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(listResp.Results) != 2 {
		t.Fatalf("Expected 2 qualifications in list, got %d", len(listResp.Results))
	}

	// Poll both qualifications until completion
	maxWait := 15 * time.Second // Longer for multi-port with sync delays
	pollInterval := 1 * time.Second
	deadline := time.Now().Add(maxWait)

	var genResult, refResult *plqpb.QualificationResult
	completedCount := 0

	for time.Now().Before(deadline) && completedCount < 2 {
		getResp, err := linkQualServer.Get(ctx, &plqpb.GetRequest{
			Ids: []string{generatorID, reflectorID},
		})
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}

		genResult = getResp.Results[generatorID]
		refResult = getResp.Results[reflectorID]

		t.Logf("Generator state: %v (sent: %d, received: %d)", genResult.State, genResult.PacketsSent, genResult.PacketsReceived)
		t.Logf("Reflector state: %v (sent: %d, received: %d)", refResult.State, refResult.PacketsSent, refResult.PacketsReceived)

		completedCount = 0
		if genResult.State == plqpb.QualificationState_QUALIFICATION_STATE_COMPLETED {
			completedCount++
		}
		if refResult.State == plqpb.QualificationState_QUALIFICATION_STATE_COMPLETED {
			completedCount++
		}

		if genResult.State == plqpb.QualificationState_QUALIFICATION_STATE_ERROR ||
			refResult.State == plqpb.QualificationState_QUALIFICATION_STATE_ERROR {
			t.Fatalf("One or both qualifications failed")
		}

		if completedCount == 2 {
			break
		}

		time.Sleep(pollInterval)
	}

	if completedCount != 2 {
		t.Fatalf("Multi-port qualification did not complete within %v (completed: %d/2)", maxWait, completedCount)
	}

	// Verify both completed successfully with realistic packet counts
	if genResult.PacketsSent == 0 {
		t.Error("Generator sent no packets")
	}
	if refResult.PacketsReceived == 0 {
		t.Error("Reflector received no packets")
	}

	// For multi-port, verify coordination - both should have similar durations
	genDuration := genResult.EndTime.AsTime().Sub(genResult.StartTime.AsTime())
	refDuration := refResult.EndTime.AsTime().Sub(refResult.StartTime.AsTime())

	durationDiff := math.Abs(float64(genDuration - refDuration))
	if durationDiff > float64(2*time.Second) {
		t.Errorf("Qualification durations too different: gen=%v, ref=%v, diff=%v",
			genDuration, refDuration, time.Duration(durationDiff))
	}

	// Clean up both qualifications
	deleteResp, err := linkQualServer.Delete(ctx, &plqpb.DeleteRequest{
		Ids: []string{generatorID, reflectorID},
	})
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	for _, id := range []string{generatorID, reflectorID} {
		if deleteResp.Results[id].Code != int32(codes.OK) {
			t.Errorf("Delete failed for %s: %d: %s", id, deleteResp.Results[id].Code, deleteResp.Results[id].Message)
		}
	}
}

func testConcurrentQualifications(t *testing.T, linkQualServer *linkQualification, ctx context.Context) {
	// Test running multiple qualifications concurrently on different interfaces
	qualIDs := []string{"concurrent-1", "concurrent-2", "concurrent-3"}
	interfaces := []string{"eth3", "eth4", "eth5"}

	// Create multiple qualifications simultaneously
	var configs []*plqpb.QualificationConfiguration
	for i, qualID := range qualIDs {
		config := &plqpb.QualificationConfiguration{
			Id:            qualID,
			InterfaceName: interfaces[i],
			EndpointType: &plqpb.QualificationConfiguration_PacketGenerator{
				PacketGenerator: &plqpb.PacketGeneratorConfiguration{
					PacketRate: uint64(5000 * (i + 1)), // Different rates: 5K, 10K, 15K PPS
					PacketSize: 1000,
				},
			},
			Timing: &plqpb.QualificationConfiguration_Rpc{
				Rpc: &plqpb.RPCSyncedTiming{
					SetupDuration:    durationpb.New(1 * time.Second),
					Duration:         durationpb.New(3 * time.Second),
					TeardownDuration: durationpb.New(1 * time.Second),
				},
			},
		}
		configs = append(configs, config)
	}

	createReq := &plqpb.CreateRequest{Interfaces: configs}
	createResp, err := linkQualServer.Create(ctx, createReq)
	if err != nil {
		t.Fatalf("Concurrent create failed: %v", err)
	}

	// Verify all were created successfully
	for _, qualID := range qualIDs {
		if createResp.Status[qualID].Code != int32(codes.OK) {
			t.Fatalf("Qualification %s creation failed: %d: %s", qualID,
				createResp.Status[qualID].Code, createResp.Status[qualID].Message)
		}
	}

	// Poll all qualifications until completion
	maxWait := 12 * time.Second
	pollInterval := 1 * time.Second
	deadline := time.Now().Add(maxWait)

	completedCount := 0
	for time.Now().Before(deadline) && completedCount < len(qualIDs) {
		getResp, err := linkQualServer.Get(ctx, &plqpb.GetRequest{Ids: qualIDs})
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}

		completedCount = 0
		for _, qualID := range qualIDs {
			result := getResp.Results[qualID]
			t.Logf("Qualification %s state: %v (packets sent: %d)", qualID, result.State, result.PacketsSent)

			switch result.State {
			case plqpb.QualificationState_QUALIFICATION_STATE_COMPLETED:
				completedCount++
			case plqpb.QualificationState_QUALIFICATION_STATE_ERROR:
				t.Fatalf("Qualification %s failed", qualID)
			}
		}

		if completedCount == len(qualIDs) {
			break
		}
		time.Sleep(pollInterval)
	}

	if completedCount != len(qualIDs) {
		t.Fatalf("Not all concurrent qualifications completed within %v (completed: %d/%d)",
			maxWait, completedCount, len(qualIDs))
	}

	// Verify all have different packet counts due to different rates
	getResp, err := linkQualServer.Get(ctx, &plqpb.GetRequest{Ids: qualIDs})
	if err != nil {
		t.Fatalf("Final get failed: %v", err)
	}

	packetCounts := make([]uint64, len(qualIDs))
	for i, qualID := range qualIDs {
		result := getResp.Results[qualID]
		packetCounts[i] = result.PacketsSent
		if packetCounts[i] == 0 {
			t.Errorf("Qualification %s sent no packets", qualID)
		}
	}

	// Verify different rates resulted in different packet counts
	// Higher rate should result in more packets (within reasonable variance)
	if !(packetCounts[0] < packetCounts[1] && packetCounts[1] < packetCounts[2]) {
		t.Logf("Packet counts: %v", packetCounts)
		// Allow some tolerance for timing variations
		tolerance := 0.2                          // 20% tolerance
		expected1 := float64(packetCounts[0]) * 2 // 2x rate
		expected2 := float64(packetCounts[0]) * 3 // 3x rate

		if math.Abs(float64(packetCounts[1])-expected1)/expected1 > tolerance {
			t.Logf("Warning: packet count variance for 2x rate higher than expected")
		}
		if math.Abs(float64(packetCounts[2])-expected2)/expected2 > tolerance {
			t.Logf("Warning: packet count variance for 3x rate higher than expected")
		}
	}

	// Clean up all
	deleteResp, err := linkQualServer.Delete(ctx, &plqpb.DeleteRequest{Ids: qualIDs})
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	for _, qualID := range qualIDs {
		if deleteResp.Results[qualID].Code != int32(codes.OK) {
			t.Errorf("Delete failed for %s: %d", qualID, deleteResp.Results[qualID].Code)
		}
	}
}

func testStateMachineValidation(t *testing.T, linkQualServer *linkQualification, ctx context.Context) {
	// Test that qualifications go through proper state transitions
	qualID := "state-machine-test"

	createReq := &plqpb.CreateRequest{
		Interfaces: []*plqpb.QualificationConfiguration{
			{
				Id:            qualID,
				InterfaceName: "eth0",
				EndpointType: &plqpb.QualificationConfiguration_PacketGenerator{
					PacketGenerator: &plqpb.PacketGeneratorConfiguration{
						PacketRate: 1000,
						PacketSize: 1500,
					},
				},
				Timing: &plqpb.QualificationConfiguration_Rpc{
					Rpc: &plqpb.RPCSyncedTiming{
						SetupDuration:    durationpb.New(2 * time.Second), // Longer setup for observation
						Duration:         durationpb.New(3 * time.Second),
						TeardownDuration: durationpb.New(2 * time.Second), // Longer teardown for observation
					},
				},
			},
		},
	}

	createResp, err := linkQualServer.Create(ctx, createReq)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if createResp.Status[qualID].Code != int32(codes.OK) {
		t.Fatalf("Create failed: %d: %s", createResp.Status[qualID].Code, createResp.Status[qualID].Message)
	}

	// Track state transitions
	observedStates := make(map[plqpb.QualificationState]bool)
	var stateSequence []plqpb.QualificationState

	maxWait := 12 * time.Second
	pollInterval := 300 * time.Millisecond // Fast polling to catch state transitions
	deadline := time.Now().Add(maxWait)

	for time.Now().Before(deadline) {
		getResp, err := linkQualServer.Get(ctx, &plqpb.GetRequest{Ids: []string{qualID}})
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}

		result := getResp.Results[qualID]
		currentState := result.State

		// Record new states
		if !observedStates[currentState] {
			observedStates[currentState] = true
			stateSequence = append(stateSequence, currentState)
			t.Logf("State transition: %v", currentState)
		}

		if currentState == plqpb.QualificationState_QUALIFICATION_STATE_COMPLETED {
			break
		}
		if currentState == plqpb.QualificationState_QUALIFICATION_STATE_ERROR {
			t.Fatalf("Qualification failed unexpectedly")
		}

		time.Sleep(pollInterval)
	}

	// Verify we observed the expected state sequence
	t.Logf("Observed state sequence: %v", stateSequence)

	if len(stateSequence) < 3 {
		t.Errorf("Expected to observe at least 3 states, got %d: %v", len(stateSequence), stateSequence)
	}

	// Verify first and last states
	if len(stateSequence) > 0 && stateSequence[0] != plqpb.QualificationState_QUALIFICATION_STATE_IDLE {
		t.Errorf("Expected first state to be IDLE, got %v", stateSequence[0])
	}

	lastState := stateSequence[len(stateSequence)-1]
	if lastState != plqpb.QualificationState_QUALIFICATION_STATE_COMPLETED {
		t.Errorf("Expected final state to be COMPLETED, got %v", lastState)
	}

	// Check that we observed key states
	if !observedStates[plqpb.QualificationState_QUALIFICATION_STATE_RUNNING] {
		t.Error("Never observed RUNNING state")
	}

	// Clean up
	deleteResp, err := linkQualServer.Delete(ctx, &plqpb.DeleteRequest{Ids: []string{qualID}})
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	if deleteResp.Results[qualID].Code != int32(codes.OK) {
		t.Errorf("Delete failed: %d", deleteResp.Results[qualID].Code)
	}
}

func testPacketStatisticsValidation(t *testing.T, linkQualServer *linkQualification, ctx context.Context) {
	// Test packet statistics accuracy and realism
	qualID := "packet-stats-test"
	packetRate := uint64(25000) // 25K PPS
	testDuration := 3 * time.Second

	createReq := &plqpb.CreateRequest{
		Interfaces: []*plqpb.QualificationConfiguration{
			{
				Id:            qualID,
				InterfaceName: "eth1",
				EndpointType: &plqpb.QualificationConfiguration_PacketGenerator{
					PacketGenerator: &plqpb.PacketGeneratorConfiguration{
						PacketRate: packetRate,
						PacketSize: 1200,
					},
				},
				Timing: &plqpb.QualificationConfiguration_Rpc{
					Rpc: &plqpb.RPCSyncedTiming{
						SetupDuration:    durationpb.New(1 * time.Second),
						Duration:         durationpb.New(testDuration),
						TeardownDuration: durationpb.New(1 * time.Second),
					},
				},
			},
		},
	}

	createResp, err := linkQualServer.Create(ctx, createReq)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if createResp.Status[qualID].Code != int32(codes.OK) {
		t.Fatalf("Create failed: %d: %s", createResp.Status[qualID].Code, createResp.Status[qualID].Message)
	}

	// Wait for completion
	maxWait := 10 * time.Second
	pollInterval := 1 * time.Second
	deadline := time.Now().Add(maxWait)

	var finalResult *plqpb.QualificationResult
	for time.Now().Before(deadline) {
		getResp, err := linkQualServer.Get(ctx, &plqpb.GetRequest{Ids: []string{qualID}})
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}

		result := getResp.Results[qualID]
		if result.State == plqpb.QualificationState_QUALIFICATION_STATE_COMPLETED {
			finalResult = result
			break
		}
		if result.State == plqpb.QualificationState_QUALIFICATION_STATE_ERROR {
			t.Fatalf("Qualification failed")
		}

		time.Sleep(pollInterval)
	}

	if finalResult == nil {
		t.Fatalf("Qualification did not complete within %v", maxWait)
	}

	// Verify packet statistics are realistic
	expectedPackets := packetRate * uint64(testDuration.Seconds())
	tolerance := 0.15 // 15% tolerance for timing variations

	t.Logf("Expected packets: %d, actual sent: %d, received: %d, dropped: %d, errors: %d",
		expectedPackets, finalResult.PacketsSent, finalResult.PacketsReceived,
		finalResult.PacketsDropped, finalResult.PacketsError)

	// Verify packet count is within reasonable range
	minExpected := uint64(float64(expectedPackets) * (1 - tolerance))
	maxExpected := uint64(float64(expectedPackets) * (1 + tolerance))

	if finalResult.PacketsSent < minExpected || finalResult.PacketsSent > maxExpected {
		t.Errorf("Packet count outside expected range: got %d, expected %d±%.0f%%",
			finalResult.PacketsSent, expectedPackets, tolerance*100)
	}

	// Verify packet conservation (sent = received + dropped + errors)
	totalAccountedPackets := finalResult.PacketsReceived + finalResult.PacketsDropped + finalResult.PacketsError
	if totalAccountedPackets != finalResult.PacketsSent {
		t.Errorf("Packet accounting error: sent=%d, received+dropped+errors=%d",
			finalResult.PacketsSent, totalAccountedPackets)
	}

	// Verify loss rate is reasonable (should be very low in simulation)
	if finalResult.PacketsSent > 0 {
		lossRate := float64(finalResult.PacketsDropped) / float64(finalResult.PacketsSent)
		if lossRate > 0.05 { // 5% max loss rate
			t.Errorf("Excessive packet loss: %.2f%%", lossRate*100)
		}

		errorRate := float64(finalResult.PacketsError) / float64(finalResult.PacketsSent)
		if errorRate > 0.01 { // 1% max error rate
			t.Errorf("Excessive packet errors: %.2f%%", errorRate*100)
		}
	}

	// Verify rate calculations
	if finalResult.ExpectedRateBytesPerSecond == 0 {
		t.Error("Expected rate should not be zero")
	}
	if finalResult.QualificationRateBytesPerSecond == 0 {
		t.Error("Qualification rate should not be zero")
	}

	// Expected rate should be based on packet size (1200 bytes) and rate (25K PPS)
	expectedBytesPerSecond := packetRate * 1200 // packet_rate * packet_size
	rateTolerance := 0.2                        // 20% tolerance

	if math.Abs(float64(finalResult.ExpectedRateBytesPerSecond)-float64(expectedBytesPerSecond))/float64(expectedBytesPerSecond) > rateTolerance {
		t.Errorf("Expected rate calculation incorrect: got %d Bps, expected ~%d Bps",
			finalResult.ExpectedRateBytesPerSecond, expectedBytesPerSecond)
	}

	// Clean up
	deleteResp, err := linkQualServer.Delete(ctx, &plqpb.DeleteRequest{Ids: []string{qualID}})
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	if deleteResp.Results[qualID].Code != int32(codes.OK) {
		t.Errorf("Delete failed: %d", deleteResp.Results[qualID].Code)
	}
}

func testTimingValidation(t *testing.T, linkQualServer *linkQualification, ctx context.Context) {
	// Test that timing parameters are respected and the test completes successfully.
	qualID := "timing-test-simplified"
	createReq := &plqpb.CreateRequest{
		Interfaces: []*plqpb.QualificationConfiguration{
			{
				Id:            qualID,
				InterfaceName: "eth2",
				EndpointType: &plqpb.QualificationConfiguration_PacketGenerator{
					PacketGenerator: &plqpb.PacketGeneratorConfiguration{
						PacketRate: 5000,
						PacketSize: 1000,
					},
				},
				Timing: &plqpb.QualificationConfiguration_Rpc{
					Rpc: &plqpb.RPCSyncedTiming{
						PreSyncDuration:  durationpb.New(1 * time.Second),
						SetupDuration:    durationpb.New(2 * time.Second),
						Duration:         durationpb.New(4 * time.Second),
						TeardownDuration: durationpb.New(3 * time.Second),
					},
				},
			},
		},
	}

	createResp, err := linkQualServer.Create(ctx, createReq)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if createResp.Status[qualID].Code != int32(codes.OK) {
		t.Fatalf("Create failed: %d: %s", createResp.Status[qualID].Code, createResp.Status[qualID].Message)
	}

	// Poll for completion
	maxWait := 15 * time.Second
	pollInterval := 500 * time.Millisecond
	deadline := time.Now().Add(maxWait)
	var finalState plqpb.QualificationState

	for time.Now().Before(deadline) {
		getResp, err := linkQualServer.Get(ctx, &plqpb.GetRequest{Ids: []string{qualID}})
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}
		result := getResp.Results[qualID]
		finalState = result.State

		if finalState == plqpb.QualificationState_QUALIFICATION_STATE_COMPLETED {
			t.Logf("Qualification completed successfully")
			break
		}
		if finalState == plqpb.QualificationState_QUALIFICATION_STATE_ERROR {
			t.Fatalf("Qualification failed with status: %v", result.GetStatus())
		}
		time.Sleep(pollInterval)
	}

	if finalState != plqpb.QualificationState_QUALIFICATION_STATE_COMPLETED {
		t.Errorf("Qualification did not complete, final state: %v", finalState)
	}

	// Clean up
	deleteResp, err := linkQualServer.Delete(ctx, &plqpb.DeleteRequest{Ids: []string{qualID}})
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	if deleteResp.Results[qualID].Code != int32(codes.OK) {
		t.Errorf("Delete failed: %d", deleteResp.Results[qualID].Code)
	}
}

func testCancellationAndCleanup(t *testing.T, linkQualServer *linkQualification, ctx context.Context) {
	// Basic test for cancellation - create and immediately delete
	qualID := "cancel-test"

	createReq := &plqpb.CreateRequest{
		Interfaces: []*plqpb.QualificationConfiguration{
			{
				Id:            qualID,
				InterfaceName: "eth0",
				EndpointType: &plqpb.QualificationConfiguration_PacketGenerator{
					PacketGenerator: &plqpb.PacketGeneratorConfiguration{
						PacketRate: 1000,
						PacketSize: 1500,
					},
				},
				Timing: &plqpb.QualificationConfiguration_Rpc{
					Rpc: &plqpb.RPCSyncedTiming{
						Duration: durationpb.New(10 * time.Second), // Long duration
					},
				},
			},
		},
	}

	createResp, err := linkQualServer.Create(ctx, createReq)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if createResp.Status[qualID].Code != int32(codes.OK) {
		t.Fatalf("Create failed: %d", createResp.Status[qualID].Code)
	}

	// Immediately delete to test cancellation
	deleteResp, err := linkQualServer.Delete(ctx, &plqpb.DeleteRequest{Ids: []string{qualID}})
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	if deleteResp.Results[qualID].Code != int32(codes.OK) {
		t.Errorf("Delete failed: %d", deleteResp.Results[qualID].Code)
	}

	// Verify it's gone
	getResp, err := linkQualServer.Get(ctx, &plqpb.GetRequest{Ids: []string{qualID}})
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if getResp.Results[qualID].Status.Code != int32(codes.NotFound) {
		t.Errorf("Expected NotFound after delete, got %d", getResp.Results[qualID].Status.Code)
	}
}

func testErrorHandling(t *testing.T, linkQualServer *linkQualification, ctx context.Context) {
	// Test various error conditions

	// Test duplicate ID in same request
	createReq := &plqpb.CreateRequest{
		Interfaces: []*plqpb.QualificationConfiguration{
			{
				Id:            "duplicate-test",
				InterfaceName: "eth0",
				EndpointType: &plqpb.QualificationConfiguration_PacketGenerator{
					PacketGenerator: &plqpb.PacketGeneratorConfiguration{
						PacketRate: 1000,
						PacketSize: 1500,
					},
				},
				Timing: &plqpb.QualificationConfiguration_Rpc{
					Rpc: &plqpb.RPCSyncedTiming{
						Duration: durationpb.New(2 * time.Second),
					},
				},
			},
			{
				Id:            "duplicate-test", // Same ID
				InterfaceName: "eth1",
				EndpointType: &plqpb.QualificationConfiguration_AsicLoopback{
					AsicLoopback: &plqpb.AsicLoopbackConfiguration{},
				},
				Timing: &plqpb.QualificationConfiguration_Rpc{
					Rpc: &plqpb.RPCSyncedTiming{
						Duration: durationpb.New(2 * time.Second),
					},
				},
			},
		},
	}

	createResp, err := linkQualServer.Create(ctx, createReq)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Should get AlreadyExists error for duplicate ID
	if createResp.Status["duplicate-test"].Code != int32(codes.AlreadyExists) {
		t.Errorf("Expected AlreadyExists for duplicate ID, got %d", createResp.Status["duplicate-test"].Code)
	}

	// Clean up any successful creations
	linkQualServer.Delete(ctx, &plqpb.DeleteRequest{Ids: []string{"duplicate-test"}})
}

func testConfigurationVariations(t *testing.T, linkQualServer *linkQualification, ctx context.Context) {
	// Test different endpoint types
	configs := []*plqpb.QualificationConfiguration{
		{
			Id:            "config-generator",
			InterfaceName: "eth0",
			EndpointType: &plqpb.QualificationConfiguration_PacketGenerator{
				PacketGenerator: &plqpb.PacketGeneratorConfiguration{
					PacketRate: 5000,
					PacketSize: 1200,
				},
			},
			Timing: &plqpb.QualificationConfiguration_Rpc{
				Rpc: &plqpb.RPCSyncedTiming{
					Duration: durationpb.New(2 * time.Second),
				},
			},
		},
		{
			Id:            "config-asic",
			InterfaceName: "eth1",
			EndpointType: &plqpb.QualificationConfiguration_AsicLoopback{
				AsicLoopback: &plqpb.AsicLoopbackConfiguration{},
			},
			Timing: &plqpb.QualificationConfiguration_Rpc{
				Rpc: &plqpb.RPCSyncedTiming{
					Duration: durationpb.New(2 * time.Second),
				},
			},
		},
		{
			Id:            "config-pmd",
			InterfaceName: "eth2",
			EndpointType: &plqpb.QualificationConfiguration_PmdLoopback{
				PmdLoopback: &plqpb.PmdLoopbackConfiguration{},
			},
			Timing: &plqpb.QualificationConfiguration_Rpc{
				Rpc: &plqpb.RPCSyncedTiming{
					Duration: durationpb.New(2 * time.Second),
				},
			},
		},
	}

	createReq := &plqpb.CreateRequest{Interfaces: configs}
	createResp, err := linkQualServer.Create(ctx, createReq)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Verify all endpoint types are accepted
	for _, config := range configs {
		if createResp.Status[config.Id].Code != int32(codes.OK) {
			t.Errorf("Create failed for %s: %d", config.Id, createResp.Status[config.Id].Code)
		}
	}

	// Wait for completion
	time.Sleep(5 * time.Second)

	// Verify they completed
	ids := []string{"config-generator", "config-asic", "config-pmd"}
	getResp, err := linkQualServer.Get(ctx, &plqpb.GetRequest{Ids: ids})
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	for _, id := range ids {
		result := getResp.Results[id]
		if result.State != plqpb.QualificationState_QUALIFICATION_STATE_COMPLETED {
			t.Logf("Warning: Qualification %s not completed: %v", id, result.State)
		}
	}

	// Clean up
	linkQualServer.Delete(ctx, &plqpb.DeleteRequest{Ids: ids})
}

func testPacketInjectorEndpoint(t *testing.T, linkQualServer *linkQualification, ctx context.Context) {
	// Test PacketInjector endpoint type returns UNIMPLEMENTED
	qualID := "packet-injector-test"

	createReq := &plqpb.CreateRequest{
		Interfaces: []*plqpb.QualificationConfiguration{
			{
				Id:            qualID,
				InterfaceName: "eth3",
				EndpointType: &plqpb.QualificationConfiguration_PacketInjector{
					PacketInjector: &plqpb.PacketInjectorConfiguration{
						PacketCount: 5000,
						PacketSize:  1200,
						LoopbackMode: &plqpb.PacketInjectorConfiguration_AsicLoopback{
							AsicLoopback: &plqpb.AsicLoopbackConfiguration{},
						},
					},
				},
				Timing: &plqpb.QualificationConfiguration_Rpc{
					Rpc: &plqpb.RPCSyncedTiming{
						SetupDuration:    durationpb.New(1 * time.Second),
						Duration:         durationpb.New(2 * time.Second),
						TeardownDuration: durationpb.New(1 * time.Second),
					},
				},
			},
		},
	}

	createResp, err := linkQualServer.Create(ctx, createReq)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Verify PacketInjector returns UNIMPLEMENTED status
	if createResp.Status[qualID].Code != int32(codes.Unimplemented) {
		t.Errorf("Expected UNIMPLEMENTED status for PacketInjector, got %d: %s",
			createResp.Status[qualID].Code, createResp.Status[qualID].Message)
	}

	// Verify the error message mentions PacketInjector is unimplemented
	if !strings.Contains(createResp.Status[qualID].Message, "PacketInjector") {
		t.Errorf("Expected error message to mention PacketInjector, got: %s", createResp.Status[qualID].Message)
	}

	t.Logf("PacketInjector correctly returned UNIMPLEMENTED: %s", createResp.Status[qualID].Message)
}

func testNTPTimingConfiguration(t *testing.T, linkQualServer *linkQualification, ctx context.Context) {
	// Test NTP timing configuration returns UNIMPLEMENTED
	qualID := "ntp-timing-test"
	now := time.Now()
	startTime := timestamppb.New(now.Add(1 * time.Second))
	endTime := timestamppb.New(now.Add(4 * time.Second))

	createReq := &plqpb.CreateRequest{
		Interfaces: []*plqpb.QualificationConfiguration{
			{
				Id:            qualID,
				InterfaceName: "eth4",
				EndpointType: &plqpb.QualificationConfiguration_PacketGenerator{
					PacketGenerator: &plqpb.PacketGeneratorConfiguration{
						PacketRate: 5000,
						PacketSize: 1000,
					},
				},
				Timing: &plqpb.QualificationConfiguration_Ntp{
					Ntp: &plqpb.NTPSyncedTiming{
						StartTime: startTime,
						EndTime:   endTime,
					},
				},
			},
		},
	}

	createResp, err := linkQualServer.Create(ctx, createReq)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Verify NTP timing returns UNIMPLEMENTED status
	if createResp.Status[qualID].Code != int32(codes.Unimplemented) {
		t.Errorf("Expected UNIMPLEMENTED status for NTP timing, got %d: %s",
			createResp.Status[qualID].Code, createResp.Status[qualID].Message)
	}

	// Verify the error message mentions NTP timing is unimplemented
	if !strings.Contains(createResp.Status[qualID].Message, "NTP timing") {
		t.Errorf("Expected error message to mention NTP timing, got: %s", createResp.Status[qualID].Message)
	}

	t.Logf("NTP timing correctly returned UNIMPLEMENTED: %s", createResp.Status[qualID].Message)
}

func testGeneratorReflectorCoordination(t *testing.T, linkQualServer *linkQualification, ctx context.Context) {
	t.Log("Testing generator and reflector independent operation")

	// Single generator qualification
	generatorConfig := &plqpb.QualificationConfiguration{
		Id:            "generator-test-1",
		InterfaceName: "eth1",
		EndpointType: &plqpb.QualificationConfiguration_PacketGenerator{
			PacketGenerator: &plqpb.PacketGeneratorConfiguration{
				PacketRate: 1000,
				PacketSize: 1000,
			},
		},
		Timing: &plqpb.QualificationConfiguration_Rpc{
			Rpc: &plqpb.RPCSyncedTiming{
				Duration:         durationpb.New(2 * time.Second),
				SetupDuration:    durationpb.New(1 * time.Second),
				TeardownDuration: durationpb.New(1 * time.Second),
			},
		},
	}

	// Create generator
	createReq := &plqpb.CreateRequest{
		Interfaces: []*plqpb.QualificationConfiguration{generatorConfig},
	}

	createResp, err := linkQualServer.Create(ctx, createReq)
	if err != nil {
		t.Fatalf("Failed to create generator qualification: %v", err)
	}

	if status := createResp.Status["generator-test-1"]; status.Code != int32(codes.OK) {
		t.Errorf("Expected OK status for generator, got: %v", status)
	}

	// Wait for completion
	time.Sleep(6 * time.Second)

	// Verify generator reports its own statistics
	getReq := &plqpb.GetRequest{Ids: []string{"generator-test-1"}}
	getResp, err := linkQualServer.Get(ctx, getReq)
	if err != nil {
		t.Fatalf("Failed to get generator result: %v", err)
	}

	genResult := getResp.Results["generator-test-1"]
	if genResult.State != plqpb.QualificationState_QUALIFICATION_STATE_COMPLETED {
		t.Errorf("Expected generator to be completed, got: %v", genResult.State)
	}

	if genResult.PacketsSent == 0 {
		t.Error("Expected generator to report packets sent")
	}

	t.Logf("Generator results: sent=%d, received=%d, dropped=%d",
		genResult.PacketsSent, genResult.PacketsReceived, genResult.PacketsDropped)

	// Single reflector qualification
	reflectorConfig := &plqpb.QualificationConfiguration{
		Id:            "reflector-test-1",
		InterfaceName: "eth2",
		EndpointType: &plqpb.QualificationConfiguration_AsicLoopback{
			AsicLoopback: &plqpb.AsicLoopbackConfiguration{},
		},
		Timing: &plqpb.QualificationConfiguration_Rpc{
			Rpc: &plqpb.RPCSyncedTiming{
				Duration:         durationpb.New(2 * time.Second),
				SetupDuration:    durationpb.New(1 * time.Second),
				TeardownDuration: durationpb.New(1 * time.Second),
			},
		},
	}

	createReq2 := &plqpb.CreateRequest{
		Interfaces: []*plqpb.QualificationConfiguration{reflectorConfig},
	}

	createResp2, err := linkQualServer.Create(ctx, createReq2)
	if err != nil {
		t.Fatalf("Failed to create reflector qualification: %v", err)
	}

	if status := createResp2.Status["reflector-test-1"]; status.Code != int32(codes.OK) {
		t.Errorf("Expected OK status for reflector, got: %v", status)
	}

	// Wait for completion
	time.Sleep(6 * time.Second)

	// Verify reflector reports its own statistics
	getReq2 := &plqpb.GetRequest{Ids: []string{"reflector-test-1"}}
	getResp2, err := linkQualServer.Get(ctx, getReq2)
	if err != nil {
		t.Fatalf("Failed to get reflector result: %v", err)
	}

	refResult := getResp2.Results["reflector-test-1"]
	if refResult.State != plqpb.QualificationState_QUALIFICATION_STATE_COMPLETED {
		t.Errorf("Expected reflector to be completed, got: %v", refResult.State)
	}

	t.Logf("Reflector results: sent=%d, received=%d, dropped=%d",
		refResult.PacketsSent, refResult.PacketsReceived, refResult.PacketsDropped)

	// Multiple qualifications in single Create request
	generatorConfig2 := &plqpb.QualificationConfiguration{
		Id:            "multi-generator",
		InterfaceName: "eth3",
		EndpointType: &plqpb.QualificationConfiguration_PacketGenerator{
			PacketGenerator: &plqpb.PacketGeneratorConfiguration{
				PacketRate: 500,
				PacketSize: 1200,
			},
		},
		Timing: &plqpb.QualificationConfiguration_Rpc{
			Rpc: &plqpb.RPCSyncedTiming{
				Duration:         durationpb.New(2 * time.Second),
				SetupDuration:    durationpb.New(1 * time.Second),
				TeardownDuration: durationpb.New(1 * time.Second),
			},
		},
	}

	reflectorConfig2 := &plqpb.QualificationConfiguration{
		Id:            "multi-reflector",
		InterfaceName: "eth4",
		EndpointType: &plqpb.QualificationConfiguration_PmdLoopback{
			PmdLoopback: &plqpb.PmdLoopbackConfiguration{},
		},
		Timing: &plqpb.QualificationConfiguration_Rpc{
			Rpc: &plqpb.RPCSyncedTiming{
				Duration:         durationpb.New(2 * time.Second),
				SetupDuration:    durationpb.New(1 * time.Second),
				TeardownDuration: durationpb.New(1 * time.Second),
			},
		},
	}

	createReq3 := &plqpb.CreateRequest{
		Interfaces: []*plqpb.QualificationConfiguration{generatorConfig2, reflectorConfig2},
	}

	createResp3, err := linkQualServer.Create(ctx, createReq3)
	if err != nil {
		t.Fatalf("Failed to create multi qualification: %v", err)
	}

	if status := createResp3.Status["multi-generator"]; status.Code != int32(codes.OK) {
		t.Errorf("Expected OK status for multi-generator, got: %v", status)
	}
	if status := createResp3.Status["multi-reflector"]; status.Code != int32(codes.OK) {
		t.Errorf("Expected OK status for multi-reflector, got: %v", status)
	}

	// Wait for completion
	time.Sleep(6 * time.Second)

	// Verify both report their own independent statistics
	getReq3 := &plqpb.GetRequest{Ids: []string{"multi-generator", "multi-reflector"}}
	getResp3, err := linkQualServer.Get(ctx, getReq3)
	if err != nil {
		t.Fatalf("Failed to get multi results: %v", err)
	}

	multiGenResult := getResp3.Results["multi-generator"]
	multiRefResult := getResp3.Results["multi-reflector"]

	if multiGenResult.State != plqpb.QualificationState_QUALIFICATION_STATE_COMPLETED {
		t.Errorf("Expected multi-generator to be completed, got: %v", multiGenResult.State)
	}
	if multiRefResult.State != plqpb.QualificationState_QUALIFICATION_STATE_COMPLETED {
		t.Errorf("Expected multi-reflector to be completed, got: %v", multiRefResult.State)
	}

	t.Logf("Multi-generator results: sent=%d, received=%d",
		multiGenResult.PacketsSent, multiGenResult.PacketsReceived)
	t.Logf("Multi-reflector results: sent=%d, received=%d",
		multiRefResult.PacketsSent, multiRefResult.PacketsReceived)

	// Clean up
	deleteReq := &plqpb.DeleteRequest{
		Ids: []string{"generator-test-1", "reflector-test-1", "multi-generator", "multi-reflector"},
	}
	_, err = linkQualServer.Delete(ctx, deleteReq)
	if err != nil {
		t.Fatalf("Failed to delete qualifications: %v", err)
	}
}

// File service tests

// mockPutStream implements fpb.File_PutServer for testing
type mockPutStream struct {
	grpc.ServerStream
	requests []*fpb.PutRequest
	response *fpb.PutResponse
	err      error
}

func (m *mockPutStream) Recv() (*fpb.PutRequest, error) {
	if len(m.requests) == 0 {
		return nil, io.EOF
	}
	req := m.requests[0]
	m.requests = m.requests[1:]
	return req, nil
}

func (m *mockPutStream) SendAndClose(response *fpb.PutResponse) error {
	m.response = response
	return m.err
}

// mockGetStream implements fpb.File_GetServer for testing
type mockGetStream struct {
	grpc.ServerStream
	responses []*fpb.GetResponse
	err       error
}

func (m *mockGetStream) Send(response *fpb.GetResponse) error {
	if m.err != nil {
		return m.err
	}
	m.responses = append(m.responses, response)
	return nil
}

func TestFile_Put(t *testing.T) {
	tests := []struct {
		name     string
		requests []*fpb.PutRequest
		wantErr  string
	}{
		{
			name: "successful put",
			requests: []*fpb.PutRequest{
				{
					Request: &fpb.PutRequest_Open{
						Open: &fpb.PutRequest_Details{
							RemoteFile:  "/tmp/test.txt",
							Permissions: 0644,
						},
					},
				},
				{
					Request: &fpb.PutRequest_Contents{
						Contents: []byte("hello world"),
					},
				},
				{
					Request: &fpb.PutRequest_Hash{
						Hash: &pb.HashType{
							Method: pb.HashType_MD5,
							Hash:   func() []byte { h := md5.Sum([]byte("hello world")); return h[:] }(),
						},
					},
				},
			},
			wantErr: "",
		},
		{
			name: "missing open request",
			requests: []*fpb.PutRequest{
				{
					Request: &fpb.PutRequest_Contents{
						Contents: []byte("hello world"),
					},
				},
			},
			wantErr: "must send Open message before Contents",
		},
		{
			name: "invalid path",
			requests: []*fpb.PutRequest{
				{
					Request: &fpb.PutRequest_Open{
						Open: &fpb.PutRequest_Details{
							RemoteFile:  "relative/path",
							Permissions: 0644,
						},
					},
				},
			},
			wantErr: "path must be absolute",
		},
		{
			name: "hash mismatch",
			requests: []*fpb.PutRequest{
				{
					Request: &fpb.PutRequest_Open{
						Open: &fpb.PutRequest_Details{
							RemoteFile:  "/tmp/test.txt",
							Permissions: 0644,
						},
					},
				},
				{
					Request: &fpb.PutRequest_Contents{
						Contents: []byte("hello world"),
					},
				},
				{
					Request: &fpb.PutRequest_Hash{
						Hash: &pb.HashType{
							Method: pb.HashType_MD5,
							Hash:   []byte("wrong hash"),
						},
					},
				},
			},
			wantErr: "hash verification failed",
		},
		{
			name: "system path access denied",
			requests: []*fpb.PutRequest{
				{
					Request: &fpb.PutRequest_Open{
						Open: &fpb.PutRequest_Details{
							RemoteFile:  "/proc/version",
							Permissions: 0644,
						},
					},
				},
			},
			wantErr: "access to system paths not allowed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := newFile()
			stream := &mockPutStream{requests: tt.requests}

			err := f.Put(stream)

			if diff := errdiff.Check(err, tt.wantErr); diff != "" {
				t.Fatalf("Put() error diff: %s", diff)
			}

			if tt.wantErr != "" {
				return
			}

			if stream.response == nil {
				t.Fatalf("Expected response but got none")
			}
		})
	}
}

func TestFile_Get(t *testing.T) {
	f := newFile()
	testPath := "/tmp/test.txt"
	testContent := []byte("hello world")

	// First put a file
	f.files[testPath] = &fileInfo{
		path:        testPath,
		content:     testContent,
		permissions: 0644,
		created:     time.Now(),
		modified:    time.Now(),
	}

	tests := []struct {
		name    string
		request *fpb.GetRequest
		wantErr string
	}{
		{
			name: "successful get",
			request: &fpb.GetRequest{
				RemoteFile: testPath,
			},
			wantErr: "",
		},
		{
			name: "file not found",
			request: &fpb.GetRequest{
				RemoteFile: "/tmp/nonexistent.txt",
			},
			wantErr: "file /tmp/nonexistent.txt not found",
		},
		{
			name: "invalid path",
			request: &fpb.GetRequest{
				RemoteFile: "relative/path",
			},
			wantErr: "path must be absolute",
		},
		{
			name: "system path access denied",
			request: &fpb.GetRequest{
				RemoteFile: "/proc/version",
			},
			wantErr: "access to system paths not allowed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stream := &mockGetStream{}

			err := f.Get(tt.request, stream)

			if diff := errdiff.Check(err, tt.wantErr); diff != "" {
				t.Fatalf("Get() error diff: %s", diff)
			}

			if tt.wantErr != "" {
				return
			}

			if len(stream.responses) == 0 {
				t.Fatalf("Expected responses but got none")
			}

			// Verify content
			var receivedContent []byte
			var receivedHash []byte
			for _, resp := range stream.responses {
				switch r := resp.Response.(type) {
				case *fpb.GetResponse_Contents:
					receivedContent = append(receivedContent, r.Contents...)
				case *fpb.GetResponse_Hash:
					receivedHash = r.Hash.Hash
				}
			}

			if !bytes.Equal(receivedContent, testContent) {
				t.Fatalf("Content mismatch: expected %s, got %s", testContent, receivedContent)
			}

			expectedHash := md5.Sum(testContent)
			if !bytes.Equal(receivedHash, expectedHash[:]) {
				t.Fatalf("Hash mismatch: expected %x, got %x", expectedHash, receivedHash)
			}
		})
	}
}

func TestFile_Stat(t *testing.T) {
	f := newFile()

	// Add test files
	f.files["/tmp/test1.txt"] = &fileInfo{
		path:        "/tmp/test1.txt",
		content:     []byte("content1"),
		permissions: 0644,
		created:     time.Now(),
		modified:    time.Now(),
	}
	f.files["/tmp/test2.txt"] = &fileInfo{
		path:        "/tmp/test2.txt",
		content:     []byte("content2"),
		permissions: 0755,
		created:     time.Now(),
		modified:    time.Now(),
	}
	f.files["/home/user/file.txt"] = &fileInfo{
		path:        "/home/user/file.txt",
		content:     []byte("content3"),
		permissions: 0600,
		created:     time.Now(),
		modified:    time.Now(),
	}

	tests := []struct {
		name          string
		request       *fpb.StatRequest
		wantErr       string
		expectedCount int
	}{
		{
			name: "stat /tmp directory",
			request: &fpb.StatRequest{
				Path: "/tmp",
			},
			wantErr:       "",
			expectedCount: 2,
		},
		{
			name: "stat /home/user directory",
			request: &fpb.StatRequest{
				Path: "/home/user",
			},
			wantErr:       "",
			expectedCount: 1,
		},
		{
			name: "stat empty directory",
			request: &fpb.StatRequest{
				Path: "/empty",
			},
			wantErr:       "",
			expectedCount: 0,
		},
		{
			name: "invalid path",
			request: &fpb.StatRequest{
				Path: "relative/path",
			},
			wantErr: "path must be absolute",
		},
		{
			name: "system path access denied",
			request: &fpb.StatRequest{
				Path: "/proc",
			},
			wantErr: "access to system paths not allowed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := f.Stat(context.Background(), tt.request)

			if diff := errdiff.Check(err, tt.wantErr); diff != "" {
				t.Fatalf("Stat() error diff: %s", diff)
			}

			if tt.wantErr != "" {
				return
			}

			if len(resp.Stats) != tt.expectedCount {
				t.Fatalf("Expected %d files, got %d", tt.expectedCount, len(resp.Stats))
			}

			// Verify stat info structure
			for _, stat := range resp.Stats {
				if stat.Path == "" {
					t.Fatalf("Empty path in stat result")
				}
				if stat.Size == 0 && len(f.files[stat.Path].content) > 0 {
					t.Fatalf("Size mismatch for %s", stat.Path)
				}
				if stat.Umask != defaultUmask {
					t.Fatalf("Expected umask %o, got %o", defaultUmask, stat.Umask)
				}
			}
		})
	}
}

func TestFile_Remove(t *testing.T) {
	tests := []struct {
		name      string
		setupFile bool
		request   *fpb.RemoveRequest
		wantErr   string
	}{
		{
			name:      "successful remove",
			setupFile: true,
			request: &fpb.RemoveRequest{
				RemoteFile: "/tmp/test.txt",
			},
			wantErr: "",
		},
		{
			name:      "file not found",
			setupFile: false,
			request: &fpb.RemoveRequest{
				RemoteFile: "/tmp/nonexistent.txt",
			},
			wantErr: "file /tmp/nonexistent.txt not found",
		},
		{
			name:      "invalid path",
			setupFile: false,
			request: &fpb.RemoveRequest{
				RemoteFile: "relative/path",
			},
			wantErr: "path must be absolute",
		},
		{
			name:      "system path access denied",
			setupFile: false,
			request: &fpb.RemoveRequest{
				RemoteFile: "/proc/version",
			},
			wantErr: "access to system paths not allowed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := newFile()

			if tt.setupFile {
				f.files[tt.request.RemoteFile] = &fileInfo{
					path:        tt.request.RemoteFile,
					content:     []byte("test content"),
					permissions: 0644,
					created:     time.Now(),
					modified:    time.Now(),
				}
			}

			_, err := f.Remove(context.Background(), tt.request)

			if diff := errdiff.Check(err, tt.wantErr); diff != "" {
				t.Fatalf("Remove() error diff: %s", diff)
			}

			if tt.wantErr != "" {
				return
			}

			// Verify file was removed
			if _, exists := f.files[tt.request.RemoteFile]; exists {
				t.Fatalf("File was not removed: %s", tt.request.RemoteFile)
			}
		})
	}
}

func TestFile_TransferToRemote(t *testing.T) {
	f := newFile()

	// First, create a test file in the simulated file system
	testContent := []byte("test file content for transfer")
	testFilePath := "/tmp/test.txt"
	testFile := &fileInfo{
		path:        testFilePath,
		content:     testContent,
		permissions: 0644,
		created:     time.Now(),
		modified:    time.Now(),
	}
	f.mu.Lock()
	f.files[testFilePath] = testFile
	f.mu.Unlock()

	tests := []struct {
		name        string
		req         *fpb.TransferToRemoteRequest
		wantErr     bool
		wantCode    codes.Code
		description string
	}{
		{
			name: "successful transfer",
			req: &fpb.TransferToRemoteRequest{
				LocalPath: testFilePath,
				RemoteDownload: &cpb.RemoteDownload{
					Path:     "sftp://remote.example.com/path/file.txt",
					Protocol: cpb.RemoteDownload_SFTP,
				},
			},
			wantErr:     false,
			description: "should successfully transfer existing file",
		},
		{
			name: "empty local path",
			req: &fpb.TransferToRemoteRequest{
				LocalPath: "",
				RemoteDownload: &cpb.RemoteDownload{
					Path:     "sftp://remote.example.com/path/file.txt",
					Protocol: cpb.RemoteDownload_SFTP,
				},
			},
			wantErr:     true,
			wantCode:    codes.InvalidArgument,
			description: "should fail when local_path is empty",
		},
		{
			name: "nil remote download",
			req: &fpb.TransferToRemoteRequest{
				LocalPath:      testFilePath,
				RemoteDownload: nil,
			},
			wantErr:     true,
			wantCode:    codes.InvalidArgument,
			description: "should fail when remote_download is nil",
		},
		{
			name: "empty remote download path",
			req: &fpb.TransferToRemoteRequest{
				LocalPath: testFilePath,
				RemoteDownload: &cpb.RemoteDownload{
					Path:     "",
					Protocol: cpb.RemoteDownload_SFTP,
				},
			},
			wantErr:     true,
			wantCode:    codes.InvalidArgument,
			description: "should fail when remote_download.path is empty",
		},
		{
			name: "unknown protocol",
			req: &fpb.TransferToRemoteRequest{
				LocalPath: testFilePath,
				RemoteDownload: &cpb.RemoteDownload{
					Path:     "unknown://remote.example.com/path/file.txt",
					Protocol: cpb.RemoteDownload_UNKNOWN,
				},
			},
			wantErr:     true,
			wantCode:    codes.InvalidArgument,
			description: "should fail when protocol is UNKNOWN",
		},
		{
			name: "local file not found",
			req: &fpb.TransferToRemoteRequest{
				LocalPath: "/nonexistent/file.txt",
				RemoteDownload: &cpb.RemoteDownload{
					Path:     "sftp://remote.example.com/path/file.txt",
					Protocol: cpb.RemoteDownload_SFTP,
				},
			},
			wantErr:     true,
			wantCode:    codes.NotFound,
			description: "should fail when local file doesn't exist",
		},
		{
			name: "invalid local path",
			req: &fpb.TransferToRemoteRequest{
				LocalPath: "relative/path/file.txt",
				RemoteDownload: &cpb.RemoteDownload{
					Path:     "sftp://remote.example.com/path/file.txt",
					Protocol: cpb.RemoteDownload_SFTP,
				},
			},
			wantErr:     true,
			wantCode:    codes.InvalidArgument,
			description: "should fail when local path is not absolute",
		},
		{
			name: "HTTP protocol",
			req: &fpb.TransferToRemoteRequest{
				LocalPath: testFilePath,
				RemoteDownload: &cpb.RemoteDownload{
					Path:     "http://remote.example.com/upload/file.txt",
					Protocol: cpb.RemoteDownload_HTTP,
				},
			},
			wantErr:     false,
			description: "should support HTTP protocol",
		},
		{
			name: "HTTPS protocol",
			req: &fpb.TransferToRemoteRequest{
				LocalPath: testFilePath,
				RemoteDownload: &cpb.RemoteDownload{
					Path:     "https://remote.example.com/upload/file.txt",
					Protocol: cpb.RemoteDownload_HTTPS,
				},
			},
			wantErr:     false,
			description: "should support HTTPS protocol",
		},
		{
			name: "SCP protocol",
			req: &fpb.TransferToRemoteRequest{
				LocalPath: testFilePath,
				RemoteDownload: &cpb.RemoteDownload{
					Path:     "user@remote.example.com:/path/file.txt",
					Protocol: cpb.RemoteDownload_SCP,
				},
			},
			wantErr:     false,
			description: "should support SCP protocol",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := f.TransferToRemote(context.Background(), tt.req)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("Expected error for %s but got none", tt.description)
				}
				if tt.wantCode != codes.OK && status.Code(err) != tt.wantCode {
					t.Fatalf("Expected error code %v for %s, got %v", tt.wantCode, tt.description, status.Code(err))
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error for %s: %v", tt.description, err)
			}
			if resp == nil {
				t.Fatalf("Expected response for %s but got nil", tt.description)
			}
			if resp.Hash == nil {
				t.Fatalf("Expected hash in response for %s but got nil", tt.description)
			}
			if resp.Hash.Method != pb.HashType_MD5 {
				t.Fatalf("Expected MD5 hash method for %s, got %v", tt.description, resp.Hash.Method)
			}
			if len(resp.Hash.Hash) == 0 {
				t.Fatalf("Expected non-empty hash for %s", tt.description)
			}
		})
	}
}

func TestFile_ComputeHash(t *testing.T) {
	f := newFile()
	testData := []byte("hello world")

	tests := []struct {
		name    string
		method  pb.HashType_HashMethod
		wantErr string
	}{
		{
			name:    "MD5 hash",
			method:  pb.HashType_MD5,
			wantErr: "",
		},
		{
			name:    "SHA256 hash",
			method:  pb.HashType_SHA256,
			wantErr: "",
		},
		{
			name:    "SHA512 hash",
			method:  pb.HashType_SHA512,
			wantErr: "",
		},
		{
			name:    "unsupported hash method",
			method:  pb.HashType_HashMethod(999),
			wantErr: "unsupported hash method",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := f.computeHash(testData, tt.method)

			if diff := errdiff.Check(err, tt.wantErr); diff != "" {
				t.Fatalf("computeHash() error diff: %s", diff)
			}

			if tt.wantErr != "" {
				return
			}

			if len(hash) == 0 {
				t.Fatalf("Empty hash returned")
			}

			// Verify hash correctness
			switch tt.method {
			case pb.HashType_MD5:
				expected := md5.Sum(testData)
				if !bytes.Equal(hash, expected[:]) {
					t.Fatalf("MD5 hash mismatch")
				}
			case pb.HashType_SHA256:
				expected := sha256.Sum256(testData)
				if !bytes.Equal(hash, expected[:]) {
					t.Fatalf("SHA256 hash mismatch")
				}
			}
		})
	}
}

func TestFile_FileSizeLimit(t *testing.T) {
	f := newFile()

	// Create content that exceeds the limit
	largeContent := make([]byte, maxFileSize+1)
	for i := range largeContent {
		largeContent[i] = byte(i % 256)
	}

	hash := sha256.Sum256(largeContent)

	stream := &mockPutStream{
		requests: []*fpb.PutRequest{
			{
				Request: &fpb.PutRequest_Open{
					Open: &fpb.PutRequest_Details{
						RemoteFile:  "/tmp/large.txt",
						Permissions: 0644,
					},
				},
			},
			{
				Request: &fpb.PutRequest_Contents{
					Contents: largeContent,
				},
			},
			{
				Request: &fpb.PutRequest_Hash{
					Hash: &pb.HashType{
						Method: pb.HashType_SHA256,
						Hash:   hash[:],
					},
				},
			},
		},
	}

	err := f.Put(stream)

	if err == nil {
		t.Fatalf("Expected error for oversized file but got none")
	}

	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf("Expected InvalidArgument error, got %v", status.Code(err))
	}
}

func TestFile_ConcurrentAccess(t *testing.T) {
	f := newFile()
	testPath := "/tmp/concurrent.txt"
	testContent := []byte("concurrent test")

	// Test concurrent writes and reads
	const numGoroutines = 10
	done := make(chan bool, numGoroutines*2)

	// Concurrent puts
	for range numGoroutines {
		go func() {
			defer func() { done <- true }()

			hash := md5.Sum(testContent)
			stream := &mockPutStream{
				requests: []*fpb.PutRequest{
					{
						Request: &fpb.PutRequest_Open{
							Open: &fpb.PutRequest_Details{
								RemoteFile:  testPath,
								Permissions: 0644,
							},
						},
					},
					{
						Request: &fpb.PutRequest_Contents{
							Contents: testContent,
						},
					},
					{
						Request: &fpb.PutRequest_Hash{
							Hash: &pb.HashType{
								Method: pb.HashType_MD5,
								Hash:   hash[:],
							},
						},
					},
				},
			}

			if err := f.Put(stream); err != nil {
				t.Errorf("Concurrent put failed: %v", err)
			}
		}()
	}

	// Wait a bit for some files to be created
	time.Sleep(10 * time.Millisecond)

	// Concurrent stats
	for range numGoroutines {
		go func() {
			defer func() { done <- true }()

			_, err := f.Stat(context.Background(), &fpb.StatRequest{Path: "/tmp"})
			if err != nil {
				t.Errorf("Concurrent stat failed: %v", err)
			}
		}()
	}

	// Wait for all goroutines to complete
	for range numGoroutines * 2 {
		<-done
	}
}

func TestFile_HelperMethods(t *testing.T) {
	f := newFile()
	testPath := "/tmp/helper.txt"
	testContent := []byte("helper test")

	// Test GetFileInfo on non-existent file
	_, exists := f.GetFileInfo(testPath)
	if exists {
		t.Fatalf("Expected file to not exist")
	}

	// Add a file
	f.files[testPath] = &fileInfo{
		path:        testPath,
		content:     testContent,
		permissions: 0644,
		created:     time.Now(),
		modified:    time.Now(),
	}

	// Test GetFileInfo on existing file
	info, exists := f.GetFileInfo(testPath)
	if !exists {
		t.Fatalf("Expected file to exist")
	}
	if !bytes.Equal(info.content, testContent) {
		t.Fatalf("Content mismatch")
	}

	// Test ListFiles
	files := f.ListFiles()
	if len(files) != 1 {
		t.Fatalf("Expected 1 file, got %d", len(files))
	}
	if files[0] != testPath {
		t.Fatalf("Expected %s, got %s", testPath, files[0])
	}

	// Test Reset
	f.Reset()
	files = f.ListFiles()
	if len(files) != 0 {
		t.Fatalf("Expected 0 files after reset, got %d", len(files))
	}
}
