load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "fakedevice",
    srcs = ["fakedevice.go"],
    importpath = "github.com/openconfig/lemming/gnmi/fakedevice",
    visibility = ["//visibility:public"],
    deps = [
        "//gnmi",
        "//gnmi/gnmiclient",
        "//gnmi/oc",
        "//gnmi/oc/ocpath",
        "//gnmi/reconciler",
        "@com_github_golang_glog//:glog",
        "@com_github_openconfig_ygnmi//ygnmi",
        "@com_github_openconfig_ygot//ygot",
    ],
)
