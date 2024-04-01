load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "gribi",
    srcs = ["gribi.go"],
    importpath = "github.com/openconfig/lemming/gribi",
    visibility = ["//visibility:public"],
    deps = [
        "//gnmi/gnmiclient",
        "//gnmi/oc",
        "//gnmi/oc/ocpath",
        "//proto/sysrib",
        "@com_github_golang_glog//:glog",
        "@com_github_openconfig_gnmi//proto/gnmi",
        "@com_github_openconfig_gribi//v1/proto/service",
        "@com_github_openconfig_gribigo//aft",
        "@com_github_openconfig_gribigo//afthelper",
        "@com_github_openconfig_gribigo//constants",
        "@com_github_openconfig_gribigo//rib",
        "@com_github_openconfig_gribigo//server",
        "@com_github_openconfig_ygnmi//ygnmi",
        "@com_github_openconfig_ygot//ygot",
        "@com_github_openconfig_ygot//ytypes",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//credentials/insecure",
    ],
)
