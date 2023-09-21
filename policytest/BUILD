load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "policytest",
    srcs = ["policytest.go"],
    importpath = "github.com/openconfig/lemming/policytest",
    visibility = ["//visibility:public"],
    deps = [
        "//bgp/tests/proto/policyval",
        "//gnmi/fakedevice",
        "//internal/attrs",
        "@com_github_openconfig_ondatra//:ondatra",
        "@com_github_openconfig_ondatra//gnmi",
        "@com_github_openconfig_ondatra//gnmi/oc",
        "@com_github_openconfig_ondatra//gnmi/oc/netinstbgp",
        "@com_github_openconfig_ondatra//gnmi/oc/ocpath",
        "@com_github_openconfig_ygnmi//ygnmi",
        "@com_github_openconfig_ygot//ygot",
    ],
)
