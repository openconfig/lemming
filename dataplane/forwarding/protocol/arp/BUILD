load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "arp",
    srcs = ["arp.go"],
    importpath = "github.com/openconfig/lemming/dataplane/forwarding/protocol/arp",
    visibility = ["//visibility:public"],
    deps = [
        "//dataplane/forwarding/infra/fwdpacket",
        "//dataplane/forwarding/protocol",
        "//dataplane/forwarding/util/frame",
        "//proto/forwarding",
    ],
)
