load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "gnoi",
    srcs = ["gnoi.go"],
    importpath = "github.com/openconfig/lemming/gnoi",
    visibility = ["//visibility:public"],
    deps = [
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
        "@org_golang_google_grpc//:go_default_library",
    ],
)
