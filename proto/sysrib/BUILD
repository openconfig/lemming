load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")
load("@rules_proto//proto:defs.bzl", "proto_library")

proto_library(
    name = "sysrib_proto",
    srcs = ["sysrib.proto"],
    visibility = ["//visibility:public"],
    deps = ["//proto/routing:routing_proto"],
)

go_proto_library(
    name = "sysrib_go_proto",
    compilers = ["@io_bazel_rules_go//proto:go_grpc"],
    importpath = "github.com/openconfig/lemming/proto/sysrib",
    proto = ":sysrib_proto",
    visibility = ["//visibility:public"],
    deps = ["//proto/routing"],
)

go_library(
    name = "sysrib",
    embed = [":sysrib_go_proto"],
    importpath = "github.com/openconfig/lemming/proto/sysrib",
    visibility = ["//visibility:public"],
)
