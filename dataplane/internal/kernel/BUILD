load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "kernel",
    srcs = ["kernel.go"],
    importpath = "github.com/openconfig/lemming/dataplane/internal/kernel",
    visibility = ["//dataplane:__subpackages__"],
    deps = select({
        "@io_bazel_rules_go//go/platform:android": [
            "@com_github_vishvananda_netlink//:netlink",
            "@org_golang_x_sys//unix",
        ],
        "@io_bazel_rules_go//go/platform:linux": [
            "@com_github_vishvananda_netlink//:netlink",
            "@org_golang_x_sys//unix",
        ],
        "//conditions:default": [],
    }),
)
