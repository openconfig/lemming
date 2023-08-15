load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "engine",
    srcs = [
        "engine.go",
        "engine_linux.go",
        "engine_nonlinux.go",
        "helpers.go",
    ],
    importpath = "github.com/openconfig/lemming/dataplane/internal/engine",
    visibility = ["//dataplane:__subpackages__"],
    deps = [
        "//dataplane/forwarding",
        "//dataplane/forwarding/attributes",
        "//dataplane/forwarding/fwdconfig",
        "//proto/dataplane",
        "//proto/forwarding",
        "@com_github_golang_glog//:glog",
        "@org_golang_google_grpc//codes",
        "@org_golang_google_grpc//status",
    ] + select({
        "@io_bazel_rules_go//go/platform:android": [
            "@com_github_vishvananda_netlink//:netlink",
        ],
        "@io_bazel_rules_go//go/platform:linux": [
            "@com_github_vishvananda_netlink//:netlink",
        ],
        "//conditions:default": [],
    }),
)
