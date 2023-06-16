load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "reconciler",
    srcs = ["reconciler.go"],
    importpath = "github.com/openconfig/lemming/gnmi/reconciler",
    visibility = ["//visibility:public"],
    deps = [
        "//gnmi/oc",
        "@com_github_golang_glog//:glog",
        "@com_github_openconfig_gnmi//errlist",
        "@com_github_openconfig_gnmi//proto/gnmi",
        "@com_github_openconfig_ygnmi//ygnmi",
    ],
)

go_test(
    name = "reconciler_test",
    srcs = ["reconciler_test.go"],
    embed = [":reconciler"],
    deps = [
        "//gnmi/oc",
        "//gnmi/oc/ocpath",
        "@com_github_google_go_cmp//cmp",
        "@com_github_openconfig_gnmi//errdiff",
        "@com_github_openconfig_gnmi//proto/gnmi",
        "@com_github_openconfig_ygnmi//ygnmi",
        "@org_golang_google_protobuf//testing/protocmp",
    ],
)
