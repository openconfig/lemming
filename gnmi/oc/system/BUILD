load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "system",
    srcs = [
        "system-0.go",
        "system-1.go",
        "system-2.go",
        "system-3.go",
        "system-4.go",
        "system-5.go",
        "system-6.go",
        "system-7.go",
    ],
    importpath = "github.com/openconfig/lemming/gnmi/oc/system",
    visibility = ["//visibility:public"],
    deps = [
        "//gnmi/oc",
        "@com_github_openconfig_ygnmi//ygnmi",
        "@com_github_openconfig_ygot//ygot",
        "@com_github_openconfig_ygot//ytypes",
    ],
)
