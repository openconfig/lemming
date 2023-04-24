load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "acl",
    srcs = ["acl-0.go"],
    importpath = "github.com/openconfig/lemming/gnmi/oc/acl",
    visibility = ["//visibility:public"],
    deps = [
        "//gnmi/oc",
        "@com_github_openconfig_ygnmi//ygnmi",
        "@com_github_openconfig_ygot//ygot",
        "@com_github_openconfig_ygot//ytypes",
    ],
)
