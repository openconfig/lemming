load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "lldp",
    srcs = [
        "lldp-0.go",
        "lldp-1.go",
        "lldp-2.go",
        "lldp-3.go",
        "lldp-4.go",
        "lldp-5.go",
        "lldp-6.go",
        "lldp-7.go",
    ],
    importpath = "github.com/openconfig/lemming/gnmi/oc/lldp",
    visibility = ["//visibility:public"],
    deps = [
        "//gnmi/oc",
        "@com_github_openconfig_ygnmi//ygnmi",
        "@com_github_openconfig_ygot//ygot",
        "@com_github_openconfig_ygot//ytypes",
    ],
)
