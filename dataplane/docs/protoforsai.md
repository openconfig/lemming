# saipb, apigen, and attrmgr

## saipb

The [saipb](../proto/sai/) is a protobuf/gRPC API generated from the [SAI API](https://github.com/opencomputeproject/SAI). It is intended to be functionally equivalent to SAI API. Each SAI API struct (eg `sai_example_api_t`) is transformed in a gRPC Service (eg `service Example`).

### API transformation

Most SAI APIs follow a similar format: create, set, get, delete functions. They use a unique integer ID to refer to an object and an enum to define the configurable attributes.

The sai_attribute_t struct is a tuple of an enum and value.
The header defines a  `sai_example_attr_t` which defines the attributes and type of the attribute. Not all attributes can be set on create (some can be read only). The value can one of a number of different types.

```c
sai_status_t create_example(sai_object_id_t *example_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list)
sai_status_t remove_example(sai_object_id_t example_id)
sai_status_t set_example_attribute(sai_object_id_t example_id, const sai_attribute_t *attr)
sai_status_t get_example_attribute(sai_object_id_t example_id, uint32_t attr_count, const sai_attribute_t *attr_list)

typedef enum _sai_example_attr_t {
    /**
     * @brief status
     *
     * @type uint64
     * @flags CREATE_AND_SET
     */
    SAI_EXAMPLE_ATTR_COUNT
    /**
     * @brief status
     *
     * @type boolean
     * @flags status
     */
    SAI_EXAMPLE_ATTR_STATUS
}
```

The corresponding protobuf API looks as like:

```proto
enum ExampleAttr {
  EXAMPLE_ATTR_UNSPECIFIED = 0;
  EXAMPLE_ATTR_COUNT = 1;
  EXAMPLE_ATTR_STATUS = 2;
}

message ExampleAttribute {
  optional uint64 count = 1 [(attr_enum_value) = 1];
  optional bool status = 2 [(attr_enum_value) = 2];
}

message CreateExampleRequest {
  option (sai_type) = OBJECT_TYPE_EXAMPLE;
  uint64 switch = 1;
  optional uint64 count = 2 [(attr_enum_value) = 1];
}

message CreateExampleResponse {
  uint64 oid = 1;
}

message RemoveExampleRequest {
  uint64 oid = 1;
}

message RemoveExampleResponse {}

message SetExampleAttributeRequest {
  uint64 oid = 1;
  optional uint64 count = 1 [(attr_enum_value) = 1];
}

message SetExampleAttributeResponse {}

message GetExampleAttributeRequest {
  uint64 oid = 1;
  repeated ExampleAttr attr_type = 2;
}

message GetExampleAttributeResponse {
  ExampleAttribute attr = 1;
}

service Example {
  rpc CreateExample(CreateExampleRequest) returns (CreateExampleResponse) {}
  rpc RemoveExample(RemoveExampleRequest) returns (RemoveExampleResponse) {}
  rpc SetExampleAttribute(SetExampleAttributeRequest)
      returns (SetExampleAttributeResponse) {}
  rpc GetExampleAttribute(GetExampleAttributeRequest)
      returns (GetExampleAttributeResponse) {}
}
```

The `ExampleAttr` is an enum that contains all attributes, these are primarily used to query specific attributes using Get.
The `ExampleAttribute` message contains all possible attribute values for the Example, this is returned in a GetResponse.

The difference between CreateExampleRequest, SetExampleAttributeRequest, and ExampleAttr are the fields they contain.
CreateExampleRequest contains fields marked CREATE_ONLY | CREATE_AND_SET, SetExampleAttributeRequest CREATE_AND_SET, and ExampleAttr  CREATE_ONLY | CREATE_AND_SET | READ_ONLY.

> Note: The proto field options `sai_type` and `attr_enum_value` metadata fields used to programatically associate the message field in the message to their enum value.
> They do not need to be not set by the client.

### Create RPC

The `CreateExampleRequest` message contains only the attributes values that can be specified on creation, the `CreateExampleResponse` contains the ID of the created object.

> Note: Some attributes are CREATE_ONLY and some are CREATE_SET.
> Note: Some objects are "switch-scoped" which mean they belong to an instance of a switch object.

#### Example

> Note: The following examples are in Go, but works the same from an language supported by protobuf/gRPC.
> See generated code for C++ [example](https://github.com/openconfig/lemming/blob/main/dataplane/standalone/sai/acl.cc#L834)

```go
resp, err := client.CreateRouterInterface(ctx, &saipb.CreateRouterInterfaceRequest{
  Switch:          1,
  Type:            saipb.RouterInterfaceType_ROUTER_INTERFACE_TYPE_PORT.Enum(),
  PortId:          proto.Bool(100),
})
fmt.Println(resp.Oid)
```

### Remove RPC

The object can be removed by its ID.

```go
_, err := client.RemoveRouterInterface(ctx, &saipb.RemoveRouterInterfaceRequest{
  Oid: 10,
})
```

### Set RPC

The Set API is similiar to Create. While in SAI API only one attribute can be set, the protobuf API doesn't enforce that restriction.

```go
_, err := client.SetRouterInterfaceAttribute(ctx, &saipb.SetRouterInterfaceAttributeRequest{
  Oid:           10,
  SrcMacAddress: []byte{10, 10, 10, 10, 10, 10},
})
```

### Get RPC

In the `GetExampleAttributeRequest` contains list the attributes to query from the server. The GetExampleAttributeResponse contains all possible attributes, in which only the requested attributes will be populated by the server.

```go
swAttrs, err := sw.GetSwitchAttribute(ctx, &saipb.GetSwitchAttributeRequest{
  Oid: swResp.Oid,
  AttrType: []saipb.SwitchAttr{
    saipb.SwitchAttr_SWITCH_ATTR_CPU_PORT,
  },
})
fmt.Println(swAttrs.GetAttr().GetCpuPort())
```

## apigen

apigen generates the protobuf and C++ source and headers based on the [SAI](https://github.com/opencomputeproject/SAI/tree/master/inc) headers. apigen parses both the C headers directly and the Doxygen xml to generate the protobufs.
The apigen packages read the Doxygen comments from the attribute enums to "flatten" the list of `sai_attributes` into messages that only include the relevant attributes with the correct types.

The apigen also generates a C++ implementation of the SAI API, which converts the C  SAI API to protobufs, performs the RPC, and converts the response. (ie a saipb Client library)
Currently, neither the gRPC API nor the C++ client library are fully complete.

> Note: The reverse is not implemented. There is no library that convert from the gRPC and calls the C style API. (ie a saipb Server library)

## attrmgr

attrmgr manager is Go gRPC interceptor and key-value store. It is an optional addon to make writing the server code easier. It uses protobuf reflection to map from fields in messages and to their respective enum values (using `attr_enum_value` annotation).
Attributes can be stored using the CreateRequest, SetRequest messages into a central map and retrieved later using the GetRequest message.
attrmgr is also gRPC interceptor, so all attributes set using Create and Set RPC calls can be queried using a Get without implementing any code in the server.

The generated protobuf code contains proto options that associate each field in the message with its enum value, so it possible to programmatically determine the enum value for a given field.
