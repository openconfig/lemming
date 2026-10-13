load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "info",
    srcs = ["info.go"],
    importpath = "github.com/openconfig/lemming/dataplane/luciusctl/info",
    visibility = ["//visibility:public"],
    deps = [
        "//dataplane/luciusctl/internal/transport",
        "//proto/forwarding",
        "@com_github_spf13_cobra//:cobra",
    ],
)
