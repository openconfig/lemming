load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "mpls",
    srcs = ["mpls.go"],
    importpath = "github.com/openconfig/lemming/dataplane/forwarding/protocol/mpls",
    visibility = ["//visibility:public"],
    deps = [
        "//dataplane/forwarding/protocol",
        "//dataplane/forwarding/util/frame",
        "//proto/forwarding",
    ],
)
