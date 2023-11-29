load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "config",
    srcs = ["config.go"],
    importpath = "github.com/openconfig/lemming/dataplane/config",
    visibility = ["//visibility:public"],
    deps = [
        "//proto/forwarding",
        "@com_github_spf13_viper//:viper",
    ],
)
