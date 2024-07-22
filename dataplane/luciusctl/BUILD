load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "luciusctl",
    srcs = ["luciusctl.go"],
    importpath = "github.com/openconfig/lemming/dataplane/luciusctl",
    visibility = ["//visibility:public"],
    deps = [
        "//dataplane/luciusctl/info",
        "//dataplane/luciusctl/sai",
        "@com_github_spf13_cobra//:cobra",
        "@com_github_spf13_viper//:viper",
    ],
)
