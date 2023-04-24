load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "packettestutil",
    srcs = ["packet_test_util.go"],
    importpath = "github.com/openconfig/lemming/dataplane/forwarding/protocol/packettestutil",
    visibility = ["//visibility:public"],
    deps = [
        "//dataplane/forwarding/infra/fwdpacket",
        "//dataplane/forwarding/util/frame",
        "//proto/forwarding",
    ],
)
