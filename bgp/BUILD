load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "bgp",
    srcs = [
        "bgpconfig.go",
        "debug.go",
        "gobgp.go",
        "ocgobgp.go",
    ],
    importpath = "github.com/openconfig/lemming/bgp",
    visibility = ["//visibility:public"],
    deps = [
        "//gnmi/fakedevice",
        "//gnmi/gnmiclient",
        "//gnmi/oc",
        "//gnmi/oc/ocpath",
        "//gnmi/reconciler",
        "@com_github_golang_glog//:glog",
        "@com_github_openconfig_ygnmi//ygnmi",
        "@com_github_openconfig_ygot//ygot",
        "@com_github_wenovus_gobgp_v3//api",
        "@com_github_wenovus_gobgp_v3//pkg/apiutil",
        "@com_github_wenovus_gobgp_v3//pkg/bgpconfig",
        "@com_github_wenovus_gobgp_v3//pkg/log",
        "@com_github_wenovus_gobgp_v3//pkg/packet/bgp",
        "@com_github_wenovus_gobgp_v3//pkg/server",
        "@com_github_wenovus_gobgp_v3//pkg/table",
        "@com_github_wenovus_gobgp_v3//pkg/zebra",
        "@org_golang_google_protobuf//types/known/anypb",
        "@org_golang_x_net//context",
    ],
)