# Config-Based Fault Injection

## Overview

Config-based fault injection allows you to pre-configure fault responses for specific gNOI RPC methods directly in the Lemming configuration file. This provides deterministic, reproducible fault scenarios for integration testing without requiring external fault injection clients.

## Configuration

### Basic Structure

```protobuf
fault_config {
  gnoi_faults {
    rpc_method: "/gnoi.system.System/Reboot"
    faults {
      msg_id: "reboot_fault_1"
      status {
        code: 3  # INVALID_ARGUMENT
        message: "Simulated reboot failure"
      }
    }
  }
}
```

### Configuration Fields

- **`rpc_method`**: Full gRPC method name (e.g., `/gnoi.system.System/Reboot`)
- **`faults`**: Array of fault messages to inject
  - **`msg_id`**: Unique identifier for the fault
  - **`msg`**: Optional modified request/response message  
  - **`status`**: gRPC status to return (code and message)

## Fault Behavior

### Exhaustible Application

When multiple faults are configured for an RPC method, they are applied sequentially and then exhausted:

```protobuf
gnoi_faults {
  rpc_method: "/gnoi.system.System/Reboot"
  faults {
    msg_id: "fault_1"
    status { code: 3 message: "First fault" }
  }
  faults {
    msg_id: "fault_2" 
    status { code: 14 message: "Second fault" }
  }
}
```

- 1st call → fault_1 (INVALID_ARGUMENT)
- 2nd call → fault_2 (UNAVAILABLE)  
- 3rd call → **normal behavior** (continues normally)

### Fault Types

#### Status-Only Faults
Return an error without executing the RPC:

```protobuf
faults {
  msg_id: "permission_denied"
  status {
    code: 7  # PERMISSION_DENIED
    message: "Access denied in maintenance mode"
  }
}
```

#### Message Modification Faults
Modify the request before processing:

```protobuf
faults {  
  msg_id: "modify_request"
  msg {
    # Encoded modified request message
    type_url: "type.googleapis.com/gnoi.system.RebootRequest"
    value: "..."  # Modified request proto bytes
  }
  status {
    code: 0  # OK - process with modified request
  } 
}
```

### Running Lemming with Fault Config

```bash
# Start lemming with fault configuration
./lemming \
  --config_file=path_to_fault_config
```
