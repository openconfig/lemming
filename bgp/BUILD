load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "bgp",
    srcs = [
        "config.go",
        "gobgp.go",
        "ocgobgp.go",
        "util.go",
    ],
    importpath = "github.com/openconfig/lemming/bgp",
    visibility = ["//visibility:public"],
    deps = [
        "//gnmi/fakedevice",
        "//gnmi/gnmiclient",
        "//gnmi/oc",
        "//gnmi/oc/ocpath",
        "//gnmi/reconciler",
        "//internal/lemmingutil",
        "@com_github_golang_glog//:glog",
        "@com_github_openconfig_gnmi//errlist",
        "@com_github_openconfig_ygnmi//ygnmi",
        "@com_github_openconfig_ygot//ygot",
        "@com_github_osrg_gobgp_v3//api",
        "@com_github_osrg_gobgp_v3//pkg/config",
        "@com_github_osrg_gobgp_v3//pkg/config/oc",
        "@com_github_osrg_gobgp_v3//pkg/server",
        "@com_github_osrg_gobgp_v3//pkg/zebra",
    ],
)

go_test(
    name = "bgp_test",
    srcs = [
        "config_test.go",
        "gobgp_test.go",
    ],
    embed = [":bgp"],
    deps = [
        "//gnmi/fakedevice",
        "//gnmi/oc",
        "@com_github_google_go_cmp//cmp",
        "@com_github_openconfig_ygot//ygot",
        "@com_github_osrg_gobgp_v3//pkg/config/oc",
    ],
)
