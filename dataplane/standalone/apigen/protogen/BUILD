load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "protogen",
    srcs = ["protogen.go"],
    importpath = "github.com/openconfig/lemming/dataplane/standalone/apigen/protogen",
    visibility = ["//visibility:public"],
    deps = [
        "//dataplane/standalone/apigen/docparser",
        "//dataplane/standalone/apigen/saiast",
        "@com_github_golang_glog//:glog",
        "@com_github_stoewer_go_strcase//:go-strcase",
    ],
)
