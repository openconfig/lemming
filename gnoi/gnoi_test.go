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
						Elem: []*pb.PathElem{{
							Name: "component",
							Key:  map[string]string{"name": "non-existent"},
						}},
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
						Elem: []*pb.PathElem{{
							Name: "component",
							Key:  map[string]string{"name": componentName},
						}},
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
					time.Sleep(10 *time.Second)
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
					Delay:  10000000000, //10 seconds
					Subcomponents: []*pb.Path{{
						Elem: []*pb.PathElem{{
							Name: "component",
							Key:  map[string]string{"name": componentName},
						}},
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
							Elem: []*pb.PathElem{{
								Name: "component",
								Key:  map[string]string{"name": "Linecard0"},
							}},
						},
						{
							Elem: []*pb.PathElem{{
								Name: "component",
								Key:  map[string]string{"name": "Fabric0"},
							}},
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
						Elem: []*pb.PathElem{{
							Name: "component",
							Key:  map[string]string{"name": "Linecard0"},
						}},
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
				// First start a linecard reboot
				linecardReq := &spb.RebootRequest{
					Method: spb.RebootMethod_COLD,
					Delay:  5000000000, // 5 second delay
					Subcomponents: []*pb.Path{{
						Elem: []*pb.PathElem{{
							Name: "component",
							Key:  map[string]string{"name": "Linecard0"},
						}},
					}},
				}
				_, err := s.Reboot(ctx, linecardReq)
				if err != nil {
					t.Fatalf("Failed to start linecard reboot: %v", err)
				}

				// Then try to reboot active supervisor
				activeSupervisorReq := &spb.RebootRequest{
					Method: spb.RebootMethod_COLD,
					Subcomponents: []*pb.Path{{
						Elem: []*pb.PathElem{{
							Name: "component",
							Key:  map[string]string{"name": "Supervisor0"},
						}},
					}},
				}
				_, err = s.Reboot(ctx, activeSupervisorReq)

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
