load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "gnoi",
    srcs = ["gnoi.go"],
    importpath = "github.com/openconfig/lemming/gnoi",
    visibility = ["//visibility:public"],
    deps = [
        "//gnmi/fakedevice",
        "@com_github_golang_glog//:glog",
        "@com_github_openconfig_gnmi//proto/gnmi",
        "@com_github_openconfig_gnoi//bgp",
        "@com_github_openconfig_gnoi//cert",
        "@com_github_openconfig_gnoi//diag",
        "@com_github_openconfig_gnoi//factory_reset",
        "@com_github_openconfig_gnoi//file",
        "@com_github_openconfig_gnoi//healthz",
        "@com_github_openconfig_gnoi//layer2",
        "@com_github_openconfig_gnoi//mpls",
        "@com_github_openconfig_gnoi//os",
        "@com_github_openconfig_gnoi//otdr",
        "@com_github_openconfig_gnoi//system",
        "@com_github_openconfig_gnoi//wavelength_router",
        "@com_github_openconfig_ygnmi//ygnmi",
        "@org_golang_google_grpc//:grpc",
        "@org_golang_google_grpc//codes",
        "@org_golang_google_grpc//status",
    ],
)

go_test(
    name = "gnoi_test",
    srcs = ["gnoi_test.go"],
    embed = [":gnoi"],
    deps = [
        "//gnmi",
        "//gnmi/fakedevice",
        "//gnmi/oc/ocpath",
        "@com_github_openconfig_gnoi//system",
        "@com_github_openconfig_ygnmi//ygnmi",
        "@org_golang_google_grpc//:grpc",
        "@org_golang_google_grpc//codes",
        "@org_golang_google_grpc//status",
    ],
)
