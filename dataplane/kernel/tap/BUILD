load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "tap",
    srcs = ["tap.go"],
    importpath = "github.com/openconfig/lemming/dataplane/kernel/tap",
    visibility = ["//visibility:public"],
    deps = [
        "//dataplane/kernel",
        "//dataplane/proto/packetio",
        "//dataplane/standalone/pkthandler/pktiohandler",
        "@com_github_vishvananda_netlink//:netlink",
    ],
)
