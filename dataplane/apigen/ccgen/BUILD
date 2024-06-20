load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "ccgen",
    srcs = ["ccgen.go"],
    importpath = "github.com/openconfig/lemming/dataplane/apigen/ccgen",
    visibility = ["//visibility:public"],
    deps = [
        "//dataplane/apigen/docparser",
        "//dataplane/apigen/saiast",
        "@com_github_stoewer_go_strcase//:go-strcase",
    ],
)
