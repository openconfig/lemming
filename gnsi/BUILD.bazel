load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "gnsi",
    srcs = ["gnsi.go"],
    importpath = "github.com/openconfig/lemming/gnsi",
    visibility = ["//visibility:public"],
    deps = [
        "//gnsi/pathz",
        "@com_github_openconfig_gnsi//authz",
        "@com_github_openconfig_gnsi//certz",
        "@com_github_openconfig_gnsi//credentialz",
        "@com_github_openconfig_gnsi//pathz",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//codes",
        "@org_golang_google_grpc//status",
    ],
)
