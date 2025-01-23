load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")
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
    srcs = ["converters.go"],
    embed = [":routing_go_proto"],
    importpath = "github.com/openconfig/lemming/proto/routing",
    visibility = ["//visibility:public"],
    deps = [
        "@com_github_google_gopacket//:gopacket",
        "@com_github_google_gopacket//layers",
    ],
)

go_test(
    name = "routing_test",
    srcs = ["converters_test.go"],
    embed = [":routing"],
    deps = [
        "@com_github_google_go_cmp//cmp",
        "@com_github_google_go_cmp//cmp/cmpopts",
        "@com_github_google_gopacket//:gopacket",
        "@com_github_google_gopacket//layers",
    ],
)
