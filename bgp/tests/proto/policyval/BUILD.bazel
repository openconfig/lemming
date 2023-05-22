load("@rules_proto//proto:defs.bzl", "proto_library")
load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")

proto_library(
    name = "policyval_proto",
    srcs = ["policy_validation.proto"],
    visibility = ["//visibility:public"],
)

go_proto_library(
    name = "policyval_go_proto",
    importpath = "github.com/openconfig/lemming/bgp/tests/proto/policyval",
    proto = ":policyval_proto",
    visibility = ["//visibility:public"],
)

go_library(
    name = "policyval",
    embed = [":policyval_go_proto"],
    importpath = "github.com/openconfig/lemming/bgp/tests/proto/policyval",
    visibility = ["//visibility:public"],
)
