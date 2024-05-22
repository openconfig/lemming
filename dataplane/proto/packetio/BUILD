load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")
load("@rules_proto//proto:defs.bzl", "proto_library")

proto_library(
    name = "packetio_proto",
    srcs = ["packetio.proto"],
    visibility = ["//visibility:public"],
    deps = ["@googleapis//google/rpc:status_proto"],
)

go_proto_library(
    name = "packetio_go_proto",
    compilers = ["@io_bazel_rules_go//proto:go_grpc"],
    importpath = "github.com/openconfig/lemming/dataplane/proto/packetio",
    proto = ":packetio_proto",
    visibility = ["//visibility:public"],
    deps = ["@org_golang_google_genproto_googleapis_rpc//status"],
)

go_library(
    name = "packetio",
    embed = [":packetio_go_proto"],
    importpath = "github.com/openconfig/lemming/dataplane/proto/packetio",
    visibility = ["//visibility:public"],
)
