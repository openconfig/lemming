load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")
load("@rules_proto//proto:defs.bzl", "proto_library")

proto_library(
    name = "config_proto",
    srcs = ["lemming_config.proto"],
    visibility = ["//visibility:public"],
)

go_proto_library(
    name = "config_go_proto",
    importpath = "github.com/openconfig/lemming/proto/config",
    proto = ":config_proto",
    visibility = ["//visibility:public"],
)

go_library(
    name = "config",
    embed = [":config_go_proto"],
    importpath = "github.com/openconfig/lemming/proto/config",
    visibility = ["//visibility:public"],
)
