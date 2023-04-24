load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "gnmiclient",
    srcs = ["client.go"],
    importpath = "github.com/openconfig/lemming/gnmi/gnmiclient",
    visibility = ["//visibility:public"],
    deps = ["@com_github_openconfig_ygnmi//ygnmi"],
)
