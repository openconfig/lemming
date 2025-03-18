load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "bootconfig",
    srcs = ["bootconfig.go"],
    importpath = "github.com/openconfig/lemming/gnoi/bootconfig",
    visibility = ["//visibility:public"],
    deps = ["@com_github_openconfig_gnoi//bootconfig"],
)

go_test(
    name = "bootconfig_test",
    srcs = ["bootconfig_test.go"],
    embed = [":bootconfig"],
)