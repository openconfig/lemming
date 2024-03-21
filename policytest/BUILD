load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "policytest",
    srcs = ["policytest.go"],
    importpath = "github.com/openconfig/lemming/policytest",
    visibility = ["//visibility:public"],
    deps = [
        "//bgp",
        "//gnmi/fakedevice",
        "//gnmi/oc",
        "//gnmi/oc/netinstbgp",
        "//gnmi/oc/ocpath",
        "//internal/attrs",
        "//proto/policyval",
        "@com_github_google_go_cmp//cmp",
        "@com_github_openconfig_ondatra//:ondatra",
        "@com_github_openconfig_ondatra//gnmi",
        "@com_github_openconfig_ygnmi//ygnmi",
        "@com_github_openconfig_ygot//ygot",
    ],
)
