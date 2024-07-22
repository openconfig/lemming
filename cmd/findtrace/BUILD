load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "findtrace",
    srcs = ["findtrace.go"],
    importpath = "github.com/openconfig/lemming/cmd/findtrace",
    visibility = ["//visibility:public"],
    deps = [
        "@com_github_fatih_color//:color",
        "@com_github_spf13_cobra//:cobra",
    ],
)
