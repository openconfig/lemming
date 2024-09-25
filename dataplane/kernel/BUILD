load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "kernel",
    srcs = [
        "genetlink.go",
        "kernel.go",
        "tap.go",
    ],
    importpath = "github.com/openconfig/lemming/dataplane/kernel",
    visibility = ["//dataplane:__subpackages__"],
    deps = [
        "@com_github_vishvananda_netlink//:netlink",
    ] + select({
        "@io_bazel_rules_go//go/platform:android": [
            "@com_github_golang_glog//:glog",
            "@com_github_mdlayher_genetlink//:genetlink",
            "@com_github_mdlayher_netlink//:netlink",
            "@org_golang_x_sys//unix",
        ],
        "@io_bazel_rules_go//go/platform:linux": [
            "@com_github_golang_glog//:glog",
            "@com_github_mdlayher_genetlink//:genetlink",
            "@com_github_mdlayher_netlink//:netlink",
            "@org_golang_x_sys//unix",
        ],
        "//conditions:default": [],
    }),
)
