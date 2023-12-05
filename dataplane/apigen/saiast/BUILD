load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "saiast",
    srcs = ["saiast.go"],
    importpath = "github.com/openconfig/lemming/dataplane/apigen/saiast",
    visibility = ["//visibility:public"],
    deps = [
        "@com_github_stoewer_go_strcase//:go-strcase",
        "@org_modernc_cc_v4//:cc",
    ],
)
