load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "dplaneopts",
    srcs = ["dplaneopts.go"],
    importpath = "github.com/openconfig/lemming/dataplane/dplaneopts",
    visibility = ["//visibility:public"],
    deps = [
        "//proto/forwarding",
        "@in_gopkg_yaml_v3//:yaml_v3",
    ],
)
