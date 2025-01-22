load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")
load("@rules_proto//proto:defs.bzl", "proto_library")

proto_library(
    name = "routing_proto",
    srcs = ["routing.proto"],
    visibility = ["//visibility:public"],
)

go_proto_library(
    name = "routing_go_proto",
    importpath = "github.com/openconfig/lemming/proto/routing",
    proto = ":routing_proto",
    visibility = ["//visibility:public"],
)

go_library(
    name = "routing",
    embed = [":routing_go_proto"],
    importpath = "github.com/openconfig/lemming/proto/routing",
    visibility = ["//visibility:public"],
)
