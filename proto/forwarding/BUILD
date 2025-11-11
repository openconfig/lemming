load("@com_github_grpc_grpc//bazel:cc_grpc_library.bzl", "cc_grpc_library")
load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")
load("@rules_proto//proto:defs.bzl", "proto_library")

proto_library(
    name = "forwarding_proto",
    srcs = [
        "forwarding_action.proto",
        "forwarding_attribute.proto",
        "forwarding_common.proto",
        "forwarding_info.proto",
        "forwarding_notification.proto",
        "forwarding_operation.proto",
        "forwarding_packetsink.proto",
        "forwarding_port.proto",
        "forwarding_service.proto",
        "forwarding_table.proto",
    ],
    visibility = ["//visibility:public"],
    deps = [
        "@googleapis//google/rpc:status_proto",
    ],
)

go_proto_library(
    name = "forwarding_go_proto",
    compilers = [
        "@io_bazel_rules_go//proto:go_grpc_v2",
        "@io_bazel_rules_go//proto:go_proto",
    ],
    importpath = "github.com/openconfig/lemming/proto/forwarding",
    proto = ":forwarding_proto",
    visibility = ["//visibility:public"],
    deps = ["@org_golang_google_genproto_googleapis_rpc//status"],
)

go_library(
    name = "forwarding",
    embed = [":forwarding_go_proto"],
    importpath = "github.com/openconfig/lemming/proto/forwarding",
    visibility = ["//visibility:public"],
)

cc_proto_library(
    name = "forwarding_cc_proto",
    visibility = ["//visibility:public"],
    deps = [":forwarding_proto"],
)

cc_grpc_library(
    name = "forwarding_cc_grpc",
    srcs = [":forwarding_proto"],
    grpc_only = True,
    visibility = ["//visibility:public"],
    deps = [":forwarding_cc_proto"],
)
