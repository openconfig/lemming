load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")
load("@rules_proto//proto:defs.bzl", "proto_library")

proto_library(
    name = "test_proto",
    srcs = ["test.proto"],
    visibility = ["//visibility:public"],
)

go_proto_library(
    name = "test_go_proto",
    compilers = ["@io_bazel_rules_go//proto:go_grpc"],
    importpath = "github.com/openconfig/lemming/fault/proto/test",
    proto = ":test_proto",
    visibility = ["//visibility:public"],
)

go_library(
    name = "test",
    embed = [":test_go_proto"],
    importpath = "github.com/openconfig/lemming/fault/proto/test",
    visibility = ["//visibility:public"],
)
