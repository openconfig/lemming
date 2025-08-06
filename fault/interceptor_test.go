package fault

import (
	"context"
	"fmt"
	"testing"

	statuspb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"

	spb "github.com/openconfig/gnoi/system"
	configpb "github.com/openconfig/lemming/proto/config"
	faultpb "github.com/openconfig/lemming/proto/fault"
)

const (
	rebootMethod = "/gnoi.system.System/Reboot"
	pingMethod   = "/gnoi.system.System/Ping"
)

// Test helper to configure faults and handle errors
func mustConfigureFaults(t *testing.T, interceptor *Interceptor, rpcMethod string, faults []*faultpb.FaultMessage) {
	t.Helper()
	if err := interceptor.configureFaults(rpcMethod, faults); err != nil {
		t.Fatalf("Failed to configure faults for %s: %v", rpcMethod, err)
	}
}

func TestInterceptorConfigureFaults(t *testing.T) {
	interceptor := NewInterceptor()
	faults := []*faultpb.FaultMessage{
		{
			MsgId: "test_fault_1",
			Status: &statuspb.Status{
				Code:    int32(codes.Internal),
				Message: "test error 1",
			},
		},
		{
			MsgId: "test_fault_2",
			Status: &statuspb.Status{
				Code:    int32(codes.Unavailable),
				Message: "test error 2",
			},
		},
	}

	mustConfigureFaults(t, interceptor, rebootMethod, faults)

	// Test fault selection order
	tests := []struct {
		name    string
		wantID  string
		wantNil bool
	}{
		{"first fault", "test_fault_1", false},
		{"second fault", "test_fault_2", false},
		{"exhausted - first call", "", true},
		{"exhausted - second call", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fault := interceptor.nextConfiguredFault(rebootMethod)
			if tt.wantNil {
				if fault != nil {
					t.Errorf("expected nil fault, got %v", fault)
				}
				return
			}

			if fault == nil {
				t.Fatal("expected fault, got nil")
			}
			if got := fault.GetMsgId(); got != tt.wantID {
				t.Errorf("expected fault ID %q, got %q", tt.wantID, got)
			}
		})
	}
}

func TestInterceptorConfigureFaultsEmpty(t *testing.T) {
	interceptor := NewInterceptor()
	mustConfigureFaults(t, interceptor, rebootMethod, []*faultpb.FaultMessage{})

	if fault := interceptor.nextConfiguredFault(rebootMethod); fault != nil {
		t.Errorf("expected no fault after clearing, got %v", fault)
	}
}

func TestFaultBehavior(t *testing.T) {
	interceptor := NewInterceptor()
	faults := []*faultpb.FaultMessage{
		{
			MsgId: "first_failure",
			Status: &statuspb.Status{
				Code:    int32(codes.PermissionDenied),
				Message: "First failure",
			},
		},
		{
			MsgId: "second_failure",
			Status: &statuspb.Status{
				Code:    int32(codes.ResourceExhausted),
				Message: "Second failure",
			},
		},
	}

	mustConfigureFaults(t, interceptor, rebootMethod, faults)

	var handlerCallCount int
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		handlerCallCount++
		return &spb.RebootResponse{}, nil
	}

	info := &grpc.UnaryServerInfo{FullMethod: rebootMethod}
	req := &spb.RebootRequest{Method: spb.RebootMethod_COLD}

	tests := []struct {
		name            string
		wantCode        codes.Code
		wantMessage     string
		wantHandlerCall bool
		wantError       bool
	}{
		{
			name:        "first fault",
			wantCode:    codes.PermissionDenied,
			wantMessage: "First failure",
			wantError:   true,
		},
		{
			name:        "second fault",
			wantCode:    codes.ResourceExhausted,
			wantMessage: "Second failure",
			wantError:   true,
		},
		{
			name:            "pass through after exhaustion",
			wantHandlerCall: true,
		},
		{
			name:            "continue pass through",
			wantHandlerCall: true,
		},
	}

	expectedHandlerCalls := 0
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := interceptor.Unary(context.Background(), req, info, handler)

			if tt.wantError {
				if err == nil {
					t.Fatal("expected error")
				}
				st := status.Convert(err)
				if st.Code() != tt.wantCode {
					t.Errorf("expected code %v, got %v", tt.wantCode, st.Code())
				}
				if st.Message() != tt.wantMessage {
					t.Errorf("expected message %q, got %q", tt.wantMessage, st.Message())
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if _, ok := resp.(*spb.RebootResponse); !ok {
					t.Errorf("expected RebootResponse, got %T", resp)
				}
			}

			if tt.wantHandlerCall {
				expectedHandlerCalls++
			}
			if handlerCallCount != expectedHandlerCalls {
				t.Errorf("expected %d handler calls, got %d", expectedHandlerCalls, handlerCallCount)
			}
		})
	}
}

func TestSingleFaultExhaustion(t *testing.T) {
	interceptor := NewInterceptor()
	fault := &faultpb.FaultMessage{
		MsgId: "single_fault",
		Status: &statuspb.Status{
			Code:    int32(codes.Unavailable),
			Message: "Single fault",
		},
	}

	mustConfigureFaults(t, interceptor, pingMethod, []*faultpb.FaultMessage{fault})

	var handlerCallCount int
	handler := func(srv interface{}, stream grpc.ServerStream) error {
		handlerCallCount++
		return nil
	}

	info := &grpc.StreamServerInfo{FullMethod: pingMethod}

	// First call: should get the fault
	err := interceptor.Stream(nil, nil, info, handler)
	if err == nil {
		t.Fatal("expected fault error")
	}

	st := status.Convert(err)
	if st.Code() != codes.Unavailable {
		t.Errorf("expected Unavailable, got %v", st.Code())
	}
	if handlerCallCount != 0 {
		t.Error("handler should not be called on fault")
	}

	// Second call: should pass through normally
	if err := interceptor.Stream(nil, nil, info, handler); err != nil {
		t.Errorf("expected success, got error: %v", err)
	}
	if handlerCallCount != 1 {
		t.Errorf("expected 1 handler call, got %d", handlerCallCount)
	}
}

func TestFaultResetOnReconfiguration(t *testing.T) {
	interceptor := NewInterceptor()

	// Configure initial fault
	initialFault := &faultpb.FaultMessage{
		MsgId: "initial_fault",
		Status: &statuspb.Status{
			Code:    int32(codes.Internal),
			Message: "Initial fault",
		},
	}
	mustConfigureFaults(t, interceptor, rebootMethod, []*faultpb.FaultMessage{initialFault})

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return &spb.RebootResponse{}, nil
	}
	info := &grpc.UnaryServerInfo{FullMethod: rebootMethod}
	req := &spb.RebootRequest{Method: spb.RebootMethod_COLD}

	// Consume initial fault
	if _, err := interceptor.Unary(context.Background(), req, info, handler); err == nil {
		t.Fatal("expected fault error")
	}

	// Verify fault is exhausted
	if _, err := interceptor.Unary(context.Background(), req, info, handler); err != nil {
		t.Errorf("expected success after exhaustion, got error: %v", err)
	}

	// Reconfigure with new fault
	newFault := &faultpb.FaultMessage{
		MsgId: "new_fault",
		Status: &statuspb.Status{
			Code:    int32(codes.FailedPrecondition),
			Message: "New fault",
		},
	}
	mustConfigureFaults(t, interceptor, rebootMethod, []*faultpb.FaultMessage{newFault})

	// Should get the new fault
	_, err := interceptor.Unary(context.Background(), req, info, handler)
	if err == nil {
		t.Fatal("expected new fault error after reconfiguration")
	}

	st := status.Convert(err)
	if st.Code() != codes.FailedPrecondition {
		t.Errorf("expected FailedPrecondition, got %v", st.Code())
	}
	if st.Message() != "New fault" {
		t.Errorf("expected 'New fault', got %q", st.Message())
	}

	// Should pass through after new fault exhausted
	if _, err := interceptor.Unary(context.Background(), req, info, handler); err != nil {
		t.Errorf("expected success after new fault exhaustion, got error: %v", err)
	}
}

func TestZeroFaultsConfiguration(t *testing.T) {
	interceptor := NewInterceptor()

	// Configure then clear faults
	fault := &faultpb.FaultMessage{
		MsgId: "will_be_cleared",
		Status: &statuspb.Status{
			Code:    int32(codes.Internal),
			Message: "Will be cleared",
		},
	}
	mustConfigureFaults(t, interceptor, rebootMethod, []*faultpb.FaultMessage{fault})
	mustConfigureFaults(t, interceptor, rebootMethod, []*faultpb.FaultMessage{})

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return &spb.RebootResponse{}, nil
	}
	info := &grpc.UnaryServerInfo{FullMethod: rebootMethod}
	req := &spb.RebootRequest{Method: spb.RebootMethod_COLD}

	if _, err := interceptor.Unary(context.Background(), req, info, handler); err != nil {
		t.Errorf("expected success with no faults configured, got error: %v", err)
	}
}

func TestInterceptorUnaryWithConfiguredFaults(t *testing.T) {
	interceptor := NewInterceptor()
	fault := &faultpb.FaultMessage{
		MsgId: "reboot_denied",
		Status: &statuspb.Status{
			Code:    int32(codes.PermissionDenied),
			Message: "Reboot not allowed",
		},
	}
	mustConfigureFaults(t, interceptor, rebootMethod, []*faultpb.FaultMessage{fault})

	var handlerCalled bool
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		handlerCalled = true
		return &spb.RebootResponse{}, nil
	}

	info := &grpc.UnaryServerInfo{FullMethod: rebootMethod}
	req := &spb.RebootRequest{Method: spb.RebootMethod_COLD}

	resp, err := interceptor.Unary(context.Background(), req, info, handler)

	if handlerCalled {
		t.Error("handler should not have been called due to configured fault")
	}
	if err == nil {
		t.Fatal("expected error from configured fault")
	}

	st, ok := status.FromError(err)
	if !ok {
		t.Fatal("expected gRPC status error")
	}
	if st.Code() != codes.PermissionDenied {
		t.Errorf("expected PermissionDenied, got %v", st.Code())
	}
	if st.Message() != "Reboot not allowed" {
		t.Errorf("expected 'Reboot not allowed', got %q", st.Message())
	}

	// Response should be the original request when fault has no message override
	if resp != req {
		t.Error("expected response to be original request when fault has no message")
	}
}

func TestInterceptorUnaryWithConfiguredFaultMessage(t *testing.T) {
	interceptor := NewInterceptor()

	// Create modified request for the fault
	faultReq := &spb.RebootRequest{
		Method: spb.RebootMethod_POWERUP,
		Force:  true,
	}
	faultReqAny, err := anypb.New(faultReq)
	if err != nil {
		t.Fatalf("failed to create Any message: %v", err)
	}

	fault := &faultpb.FaultMessage{
		MsgId: "reboot_modified",
		Msg:   faultReqAny,
		Status: &statuspb.Status{
			Code:    int32(codes.OK),
			Message: "",
		},
	}
	mustConfigureFaults(t, interceptor, rebootMethod, []*faultpb.FaultMessage{fault})

	var receivedReq *spb.RebootRequest
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		receivedReq = req.(*spb.RebootRequest)
		return &spb.RebootResponse{}, nil
	}

	info := &grpc.UnaryServerInfo{FullMethod: rebootMethod}
	originalReq := &spb.RebootRequest{Method: spb.RebootMethod_COLD}

	_, err = interceptor.Unary(context.Background(), originalReq, info, handler)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if receivedReq == nil {
		t.Fatal("handler was not called")
	}
	if receivedReq.GetMethod() != spb.RebootMethod_POWERUP {
		t.Errorf("expected modified method POWERUP, got %v", receivedReq.GetMethod())
	}
	if !receivedReq.GetForce() {
		t.Error("expected modified force=true")
	}
}

func TestInterceptorUnaryWithoutConfiguredFaults(t *testing.T) {
	interceptor := NewInterceptor()

	var handlerCalled bool
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		handlerCalled = true
		return &spb.RebootResponse{}, nil
	}

	info := &grpc.UnaryServerInfo{FullMethod: rebootMethod}
	req := &spb.RebootRequest{Method: spb.RebootMethod_COLD}

	resp, err := interceptor.Unary(context.Background(), req, info, handler)

	if !handlerCalled {
		t.Error("handler should have been called when no faults configured")
	}
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if _, ok := resp.(*spb.RebootResponse); !ok {
		t.Errorf("expected RebootResponse, got %T", resp)
	}
}

func TestInterceptorUnaryWithModifiedRequest(t *testing.T) {
	interceptor := NewInterceptor()

	// Create a modified request message
	modifiedReq := &spb.RebootRequest{Method: spb.RebootMethod_WARM}
	anyMsg, err := anypb.New(modifiedReq)
	if err != nil {
		t.Fatalf("Failed to create any message: %v", err)
	}

	fault := &faultpb.FaultMessage{
		MsgId: "modify_request",
		Msg:   anyMsg,
		Status: &statuspb.Status{
			Code:    0, // OK - should allow modified request to be processed
			Message: "Modified request",
		},
	}
	mustConfigureFaults(t, interceptor, rebootMethod, []*faultpb.FaultMessage{fault})

	var handlerCalled bool
	var receivedReq *spb.RebootRequest
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		handlerCalled = true
		if r, ok := req.(*spb.RebootRequest); ok {
			receivedReq = r
		}
		return &spb.RebootResponse{}, nil
	}

	info := &grpc.UnaryServerInfo{FullMethod: rebootMethod}
	originalReq := &spb.RebootRequest{Method: spb.RebootMethod_COLD}

	resp, err := interceptor.Unary(context.Background(), originalReq, info, handler)

	if !handlerCalled {
		t.Error("handler should have been called with modified request")
	}
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if receivedReq == nil {
		t.Error("handler should have received the modified request")
	}
	if receivedReq.Method != spb.RebootMethod_WARM {
		t.Errorf("expected WARM reboot method, got %v", receivedReq.Method)
	}
	if _, ok := resp.(*spb.RebootResponse); !ok {
		t.Errorf("expected RebootResponse, got %T", resp)
	}
}

func TestInterceptorUnaryWithModifiedRequestAndError(t *testing.T) {
	interceptor := NewInterceptor()

	// Create a modified request message
	modifiedReq := &spb.RebootRequest{Method: spb.RebootMethod_WARM}
	anyMsg, err := anypb.New(modifiedReq)
	if err != nil {
		t.Fatalf("Failed to create any message: %v", err)
	}

	fault := &faultpb.FaultMessage{
		MsgId: "modify_request_with_error",
		Msg:   anyMsg,
		Status: &statuspb.Status{
			Code:    int32(codes.PermissionDenied), // Non-zero code - should return error
			Message: "Permission denied",
		},
	}
	mustConfigureFaults(t, interceptor, rebootMethod, []*faultpb.FaultMessage{fault})

	var handlerCalled bool
	handler := func(ctx context.Context, req any) (any, error) {
		handlerCalled = true
		return &spb.RebootResponse{}, nil
	}

	info := &grpc.UnaryServerInfo{FullMethod: rebootMethod}
	originalReq := &spb.RebootRequest{Method: spb.RebootMethod_COLD}

	resp, err := interceptor.Unary(context.Background(), originalReq, info, handler)

	if handlerCalled {
		t.Error("handler should not have been called due to fault error")
	}
	if err == nil {
		t.Fatal("expected error from fault")
	}

	st := status.Convert(err)
	if st.Code() != codes.PermissionDenied {
		t.Errorf("expected PermissionDenied error, got %v", st.Code())
	}
	if st.Message() != "Permission denied" {
		t.Errorf("expected 'Permission denied' message, got %q", st.Message())
	}
	if resp == nil {
		t.Error("expected modified request as response")
	}
}

func TestInterceptorStreamWithConfiguredFaults(t *testing.T) {
	interceptor := NewInterceptor()
	fault := &faultpb.FaultMessage{
		MsgId: "stream_fault",
		Status: &statuspb.Status{
			Code:    int32(codes.Unavailable),
			Message: "Stream unavailable",
		},
	}
	mustConfigureFaults(t, interceptor, pingMethod, []*faultpb.FaultMessage{fault})

	var handlerCalled bool
	handler := func(srv interface{}, stream grpc.ServerStream) error {
		handlerCalled = true
		return nil
	}

	info := &grpc.StreamServerInfo{FullMethod: pingMethod}

	err := interceptor.Stream(nil, nil, info, handler)

	if handlerCalled {
		t.Error("handler should not have been called due to configured fault")
	}
	if err == nil {
		t.Fatal("expected error from configured stream fault")
	}

	st, ok := status.FromError(err)
	if !ok {
		t.Fatal("expected gRPC status error")
	}
	if st.Code() != codes.Unavailable {
		t.Errorf("expected Unavailable, got %v", st.Code())
	}
}

func TestInterceptorStreamWithoutConfiguredFaults(t *testing.T) {
	interceptor := NewInterceptor()

	var handlerCalled bool
	handler := func(srv interface{}, stream grpc.ServerStream) error {
		handlerCalled = true
		return nil
	}

	info := &grpc.StreamServerInfo{FullMethod: pingMethod}

	err := interceptor.Stream(nil, nil, info, handler)

	if !handlerCalled {
		t.Error("handler should have been called when no faults configured")
	}
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// TestConfigureFaultsWithNil tests edge case with nil faults slice
func TestConfigureFaultsWithNil(t *testing.T) {
	interceptor := NewInterceptor()

	// Configure with nil should not panic
	mustConfigureFaults(t, interceptor, rebootMethod, nil)

	// Should pass through normally
	var handlerCalled bool
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		handlerCalled = true
		return &spb.RebootResponse{}, nil
	}

	info := &grpc.UnaryServerInfo{FullMethod: rebootMethod}
	req := &spb.RebootRequest{Method: spb.RebootMethod_COLD}

	_, err := interceptor.Unary(context.Background(), req, info, handler)

	if !handlerCalled {
		t.Error("handler should have been called with nil faults")
	}
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// TestConfigureFaultsInvalidArgument tests error handling for invalid arguments
func TestConfigureFaultsInvalidArgument(t *testing.T) {
	interceptor := NewInterceptor()

	// Test empty RPC method
	err := interceptor.configureFaults("", []*faultpb.FaultMessage{
		{MsgId: "test", Status: &statuspb.Status{Code: int32(codes.Internal)}},
	})

	if err == nil {
		t.Fatal("expected error for empty RPC method")
	}

	st, ok := status.FromError(err)
	if !ok {
		t.Fatal("expected gRPC status error")
	}

	if st.Code() != codes.InvalidArgument {
		t.Errorf("expected InvalidArgument, got %v", st.Code())
	}

	expectedMessage := "rpc_method cannot be empty"
	if st.Message() != expectedMessage {
		t.Errorf("expected error message %q, got %q", expectedMessage, st.Message())
	}
}

// TestInterceptorConcurrency tests concurrent access to interceptor
func TestInterceptorConcurrency(t *testing.T) {
	interceptor := NewInterceptor()
	fault := &faultpb.FaultMessage{
		MsgId: "concurrent_fault",
		Status: &statuspb.Status{
			Code:    int32(codes.Internal),
			Message: "Concurrent fault",
		},
	}

	// Configure faults concurrently
	const numGoroutines = 50
	done := make(chan bool, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer func() { done <- true }()

			// Configure and use faults
			method := fmt.Sprintf("/test.Service/Method%d", id%5)
			mustConfigureFaults(t, interceptor, method, []*faultpb.FaultMessage{fault})

			// Try to get next fault
			_ = interceptor.nextConfiguredFault(method)

			// Clear faults
			mustConfigureFaults(t, interceptor, method, []*faultpb.FaultMessage{})
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < numGoroutines; i++ {
		<-done
	}
}

// TestNewInterceptorFromConfig tests the new constructor that loads config
func TestNewInterceptorFromConfig(t *testing.T) {
	t.Parallel()

	// Test with nil config
	interceptor := NewInterceptorFromConfig(nil)
	if interceptor == nil {
		t.Fatal("expected non-nil interceptor with nil config")
	}

	// Should pass through normally with no config
	var handlerCalled bool
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		handlerCalled = true
		return &spb.RebootResponse{}, nil
	}

	info := &grpc.UnaryServerInfo{FullMethod: rebootMethod}
	req := &spb.RebootRequest{Method: spb.RebootMethod_COLD}

	_, err := interceptor.Unary(context.Background(), req, info, handler)
	if err != nil {
		t.Errorf("expected success with nil config, got error: %v", err)
	}
	if !handlerCalled {
		t.Error("handler should have been called with nil config")
	}

	// Test with actual config
	faultConfig := &configpb.FaultServiceConfiguration{
		GnoiFaults: []*configpb.GNOIFaults{
			{
				RpcMethod: rebootMethod,
				Faults: []*faultpb.FaultMessage{
					{
						MsgId: "config_fault",
						Status: &statuspb.Status{
							Code:    int32(codes.Internal),
							Message: "Config-based fault",
						},
					},
				},
			},
		},
	}

	interceptor = NewInterceptorFromConfig(faultConfig)
	handlerCalled = false

	// Should inject the configured fault
	_, err = interceptor.Unary(context.Background(), req, info, handler)
	if err == nil {
		t.Fatal("expected fault error from config")
	}

	st := status.Convert(err)
	if st.Code() != codes.Internal {
		t.Errorf("expected Internal error, got %v", st.Code())
	}
	if st.Message() != "Config-based fault" {
		t.Errorf("expected 'Config-based fault', got %q", st.Message())
	}
	if handlerCalled {
		t.Error("handler should not have been called due to fault")
	}
}
