load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "sink",
    srcs = ["cpusink.go"],
    importpath = "github.com/openconfig/lemming/dataplane/standalone/cpusink/sink",
    visibility = ["//visibility:public"],
    deps = [
        "//dataplane/forwarding/fwdconfig",
        "//dataplane/internal/kernel",
        "//proto/forwarding",
        "@com_github_golang_glog//:glog",
        "@com_github_vishvananda_netlink//:netlink",
    ],
)
