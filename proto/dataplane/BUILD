load("@com_github_grpc_grpc//bazel:cc_grpc_library.bzl", "cc_grpc_library")
load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")
load("@rules_proto//proto:defs.bzl", "proto_library")

proto_library(
    name = "dataplane_proto",
    srcs = ["dataplane.proto"],
    visibility = ["//visibility:public"],
    deps = ["//proto/routing:routing_proto"],
)

go_proto_library(
    name = "dataplane_go_proto",
    importpath = "github.com/openconfig/lemming/proto/dataplane",
    proto = ":dataplane_proto",
    visibility = ["//visibility:public"],
    deps = ["//proto/routing"],
)

go_library(
    name = "dataplane",
    embed = [":dataplane_go_proto"],
    importpath = "github.com/openconfig/lemming/proto/dataplane",
    visibility = ["//visibility:public"],
)

cc_proto_library(
    name = "dataplane_cc_proto",
    visibility = ["//visibility:public"],
    deps = [":dataplane_proto"],
)

cc_grpc_library(
    name = "dataplane_cc_grpc",
    srcs = [":dataplane_proto"],
    grpc_only = True,
    visibility = ["//visibility:public"],
    deps = [":dataplane_cc_proto"],
)
