package gnoi

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	spb "github.com/openconfig/gnoi/system"
	pb "github.com/openconfig/gnoi/types"
	"github.com/openconfig/ygnmi/ygnmi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/openconfig/lemming/gnmi"
	"github.com/openconfig/lemming/gnmi/fakedevice"
	"github.com/openconfig/lemming/gnmi/oc"
	"github.com/openconfig/lemming/gnmi/oc/ocpath"
)

const (
	// Default supervisor configuration
	defaultPrimarySupervisor   = "Supervisor1"
	defaultSecondarySupervisor = "Supervisor2"
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
				deadline := time.Now().Add(30 * time.Second)
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
				deadline = time.Now().Add(3 * time.Minute)
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
					time.Sleep(10 * time.Second)
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
					Delay:  10000000000, // 10 seconds
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
					Delay:  5000000000,
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
	s := newSystem(c)

	// Initialize the system
	if err := fakedevice.NewBootTimeTask().Start(ctx, client, "local"); err != nil {
		t.Fatalf("Failed to initialize boot time: %v", err)
	}
	if err := fakedevice.NewChassisComponentsTask().Start(ctx, client, "local"); err != nil {
		t.Fatalf("Failed to initialize components: %v", err)
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			test.fn(t, s, ctx)
		})
	}
}

// TestSwitchControlProcessor tests the SwitchControlProcessor method.
func TestSwitchControlProcessor(t *testing.T) {
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
							{Name: "Supervisor1"},
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
					Delay:  5000000000, // 5 seconds
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
					Delay:  5000000000, // 5 seconds
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
		s := newSystem(c)

		// Initialize the system
		if err := fakedevice.NewBootTimeTask().Start(ctx, client, "local"); err != nil {
			t.Fatalf("Failed to initialize boot time: %v", err)
		}
		if err := fakedevice.NewChassisComponentsTask().Start(ctx, client, "local"); err != nil {
			t.Fatalf("Failed to initialize components: %v", err)
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

		s := newSystem(c)
		ctx := context.Background()
		fakedevice.NewProcessMonitoringTask().Start(ctx, client, "local")

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
				// Kill ospfd with KILL signal, no restart
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

				// Reload sysrib with HUP signal
				req := &spb.KillProcessRequest{
					Name:    "sysrib",
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

				if reloadedProcess.GetName() != "sysrib" {
					t.Errorf("Process name changed, got %s, want sysrib", reloadedProcess.GetName())
				}
			},
		},
		"process-restart-with-new-pid": {
			fn: func(t *testing.T, s *system, ctx context.Context, c *ygnmi.Client) {
				// Kill and restart bgpd
				req := &spb.KillProcessRequest{
					Name:    "bgpd",
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

				var newBgpdProcess *oc.System_Process
				for _, p := range processes {
					if p.GetName() == "bgpd" {
						newBgpdProcess = p
						break
					}
				}

				if newBgpdProcess == nil {
					t.Fatal("bgpd process not restarted")
				}

				if newBgpdProcess.GetPid() == 1001 {
					t.Error("Restarted process should have different PID")
				}

				if newBgpdProcess.GetPid() < 1001 || newBgpdProcess.GetPid() > 1100 {
					t.Errorf("New PID %d should be in range 1001-1100", newBgpdProcess.GetPid())
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
					"bgpd":        1001,
					"ospfd":       1002,
					"gnmi-server": 1003,
					"sysrib":      1004,
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
