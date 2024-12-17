load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "sampling",
    srcs = [
        "sampling-0.go",
        "sampling-1.go",
        "sampling-2.go",
        "sampling-3.go",
        "sampling-4.go",
        "sampling-5.go",
        "sampling-6.go",
        "sampling-7.go",
    ],
    importpath = "github.com/openconfig/lemming/gnmi/oc/sampling",
    visibility = ["//visibility:public"],
    deps = [
        "//gnmi/oc",
        "@com_github_openconfig_ygnmi//ygnmi",
        "@com_github_openconfig_ygot//ygot",
        "@com_github_openconfig_ygot//ytypes",
    ],
)
