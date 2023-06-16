load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "icmp",
    srcs = [
        "icmp.go",
        "icmp4.go",
        "icmp6.go",
    ],
    importpath = "github.com/openconfig/lemming/dataplane/forwarding/protocol/icmp",
    visibility = ["//visibility:public"],
    deps = [
        "//dataplane/forwarding/infra/fwdpacket",
        "//dataplane/forwarding/protocol",
        "//dataplane/forwarding/util/frame",
        "//dataplane/forwarding/util/hash/csum16",
        "//proto/forwarding",
    ],
)
