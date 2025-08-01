load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "configs",
    srcs = ["configs.go"],
    embedsrcs = [
        "arista_default.textproto",
        "lemming_default.textproto",
        "nokia_default.textproto",
        "fault_config.textproto",
    ],
    importpath = "github.com/openconfig/lemming/configs",
    visibility = ["//visibility:public"],
)
