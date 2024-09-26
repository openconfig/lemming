load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "kernel",
    srcs = [
        "kernel.go",
        "tap.go",
    ],
    importpath = "github.com/openconfig/lemming/dataplane/kernel",
    visibility = ["//dataplane:__subpackages__"],
    deps = [
        "@com_github_vishvananda_netlink//:netlink",
    ] + select({
        "@io_bazel_rules_go//go/platform:android": [
            "@org_golang_x_sys//unix",
        ],
        "@io_bazel_rules_go//go/platform:linux": [
            "@org_golang_x_sys//unix",
        ],
        "//conditions:default": [],
    }),
)
