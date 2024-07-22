load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "lemctl_lib",
    srcs = ["main.go"],
    importpath = "github.com/openconfig/lemming/cmd/lemctl",
    visibility = ["//visibility:private"],
    deps = [
        "//cmd/findtrace",
        "//cmd/packet",
        "//cmd/release",
        "//dataplane/luciusctl",
        "@com_github_spf13_cobra//:cobra",
    ],
)

go_binary(
    name = "lemctl",
    embed = [":lemctl_lib"],
    visibility = ["//visibility:public"],
)
