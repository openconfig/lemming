load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "fwdconfig",
    srcs = [
        "action.go",
        "packet.go",
        "table.go",
    ],
    importpath = "github.com/openconfig/lemming/dataplane/forwarding/fwdconfig",
    visibility = ["//visibility:public"],
    deps = ["//proto/forwarding"],
)
