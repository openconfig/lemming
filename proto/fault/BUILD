load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")
load("@rules_proto//proto:defs.bzl", "proto_library")

proto_library(
    name = "fault_proto",
    srcs = ["fault.proto"],
    visibility = ["//visibility:public"],
    deps = [
        "@com_google_protobuf//:any_proto",
        "@googleapis//google/rpc:status_proto",
    ],
)

go_proto_library(
    name = "fault_go_proto",
    compilers = ["@io_bazel_rules_go//proto:go_grpc"],
    importpath = "github.com/openconfig/lemming/proto/fault",
    proto = ":fault_proto",
    visibility = ["//visibility:public"],
    deps = ["@org_golang_google_genproto_googleapis_rpc//status"],
)

go_library(
    name = "fault",
    embed = [":fault_go_proto"],
    importpath = "github.com/openconfig/lemming/proto/fault",
    visibility = ["//visibility:public"],
)
