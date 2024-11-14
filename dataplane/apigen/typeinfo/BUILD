load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "typeinfo",
    srcs = [
        "typeinfo.go",
        "types.go",
    ],
    importpath = "github.com/openconfig/lemming/dataplane/apigen/typeinfo",
    visibility = ["//visibility:public"],
    deps = [
        "//dataplane/apigen/docparser",
        "//dataplane/apigen/saiast",
        "@com_github_golang_glog//:glog",
        "@com_github_stoewer_go_strcase//:go-strcase",
    ],
)
