load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "keychain",
    srcs = ["keychain-0.go"],
    importpath = "github.com/openconfig/lemming/gnmi/oc/keychain",
    visibility = ["//visibility:public"],
    deps = [
        "//gnmi/oc",
        "@com_github_openconfig_ygnmi//ygnmi",
        "@com_github_openconfig_ygot//ygot",
        "@com_github_openconfig_ygot//ytypes",
    ],
)
