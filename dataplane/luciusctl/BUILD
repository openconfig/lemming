load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "luciusctl_lib",
    srcs = ["main.go"],
    importpath = "github.com/openconfig/lemming/dataplane/luciusctl",
    visibility = ["//visibility:private"],
    deps = [
        "//dataplane/luciusctl/info",
        "@com_github_spf13_cobra//:cobra",
        "@com_github_spf13_viper//:viper",
    ],
)

go_binary(
    name = "luciusctl",
    embed = [":luciusctl_lib"],
    visibility = ["//visibility:public"],
)
