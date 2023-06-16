load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "kernel",
    srcs = ["kernel.go"],
    importpath = "github.com/openconfig/lemming/dataplane/internal/kernel",
    visibility = ["//dataplane:__subpackages__"],
    deps = [
        "@com_github_vishvananda_netlink//:netlink",
        "@org_golang_x_sys//unix",
    ],
)
