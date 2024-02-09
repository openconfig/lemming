# saipb, apigen, and attrmgr

## saipb

The [saipb](../proto/sai/) is a protobuf/gRPC API generated from the [SAI API](https://github.com/opencomputeproject/SAI). It is intended to be functionally equivalent to SAI API. Each SAI API struct (eg `sai_example_api_t`) is transformed in a gRPC Service (eg `service Example`).

### API transformation

Most SAI APIs follow a similar format: create, set, get, delete functions. They use a unique integer ID to refer to an object and an enum to define the configurable

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

The `ExampleAttr` enum contains all attributes, these are used to query specific attributes using Get.
The `ExampleAttribute` message contains all possible attribute values for the Example.

> Note: The proto field options are used to programatically associate the message field in the message to their enum value.

### Create

The `CreateExampleRequest` message contains only the attributes values that can be specified on creation, the `CreateExampleResponse` contains the ID of the created object.

> Note: Some attributes are CREATE_ONLY and some CREATE_SET.
> Note: Some objects are "switch-scoped" which mean they belong to an instance of a switch object.

### Delete

The object can be deleted by its ID.

### Set

The Set API is similiar to Create. While in SAI API only one attribute can be set, the protobuf API doesn't enforce that restriction.

### Get

In the `GetExampleAttributeRequest` contains list the attributes to query from the server. The GetExampleAttributeResponse contains all possible attributes, in which only the requested attributes will be populated by the server.

## apigen

apigen generate the protobuf and C++ source and headers based on the [SAI](https://github.com/opencomputeproject/SAI/tree/master/inc) headers. apigen parses both the C headers directly and the Doxygen xml to generate the protobufs.
The apigen packages read the Doxygen comments from the attribute enums to "flatten" the list of `sai_attributes` into messages that only include the relevant attributes with the correct types.

The apigen also generates a C++ implementation of the SAI API, which converts the C  SAI API to protobufs, performs the RPC, and converts the response. (ie a saipb Client library)

> Note: The reverse is not implemented. There is no library that convert from the gRPC and calls the C style API. (ie a saipb Server library)

## attrmgr

attrmgr manager is Go gRPC interceptor and key-value store. It optional addon to make writing the server code easier. It uses protobuf reflection to store all the attributes in a CreateFooRequest and SetFooRequest as a map from enum value to proto value.
In other words, attributes can be stored using the CreateRequest or SetRequest messages and retrieved using the GetRequest message.
Since attrmgr is also gRPC interceptor, that all attributes set using Create and Set RPC calls can be queried using a Get without implementing any code in the server.

The generated protobuf code contains proto options that associate each field in the message with its enum value, so it possible to programmatically determine the enum value for a given field.
