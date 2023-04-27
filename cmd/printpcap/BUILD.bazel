load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "printpcap_lib",
    srcs = ["printpcap.go"],
    importpath = "github.com/openconfig/lemming/cmd/printpcap",
    visibility = ["//visibility:private"],
    deps = ["//internal/packetutil"],
)

go_binary(
    name = "printpcap",
    embed = [":printpcap_lib"],
    visibility = ["//visibility:public"],
)
