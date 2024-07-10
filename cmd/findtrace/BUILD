load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "findtrace_lib",
    srcs = ["findtrace.go"],
    importpath = "github.com/openconfig/lemming/cmd/findtrace",
    visibility = ["//visibility:private"],
    deps = ["@com_github_fatih_color//:color"],
)

go_binary(
    name = "findtrace",
    embed = [":findtrace_lib"],
    visibility = ["//visibility:public"],
)
