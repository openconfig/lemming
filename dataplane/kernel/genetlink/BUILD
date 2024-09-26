load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "genetlink",
    srcs = [
        "genetlink.c",
        "genetlink.go",
        "genetlink.h",
    ],
    cgo = True,
    clinkopts = select({
        "@io_bazel_rules_go//go/platform:android": [
            "-lnl-3 -lnl-genl-3",
        ],
        "@io_bazel_rules_go//go/platform:linux": [
            "-lnl-3 -lnl-genl-3",
        ],
        "//conditions:default": [],
    }),
    copts = select({
        "@io_bazel_rules_go//go/platform:android": [
            "-I/usr/include/libnl3",
        ],
        "@io_bazel_rules_go//go/platform:linux": [
            "-I/usr/include/libnl3",
        ],
        "//conditions:default": [],
    }),
    importpath = "github.com/openconfig/lemming/dataplane/kernel/genetlink",
    visibility = ["//visibility:public"],
    deps = select({
        "@io_bazel_rules_go//go/platform:android": [
            "//dataplane/kernel",
            "@com_github_golang_glog//:glog",
        ],
        "@io_bazel_rules_go//go/platform:linux": [
            "//dataplane/kernel",
            "@com_github_golang_glog//:glog",
        ],
        "//conditions:default": [],
    }),
)
