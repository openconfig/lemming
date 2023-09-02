load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "ccgen",
    srcs = ["ccgen.go"],
    importpath = "github.com/openconfig/lemming/dataplane/standalone/apigen/ccgen",
    visibility = ["//visibility:public"],
    deps = [
        "//dataplane/standalone/apigen/docparser",
        "//dataplane/standalone/apigen/protogen",
        "//dataplane/standalone/apigen/saiast",
        "@com_github_stoewer_go_strcase//:go-strcase",
    ],
)
