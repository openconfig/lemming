load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "dataplane",
    srcs = ["server.go"],
    importpath = "github.com/openconfig/lemming/dataplane",
    visibility = ["//visibility:public"],
    deps = [
        "//dataplane/handlers",
        "//dataplane/internal/engine",
        "//gnmi/oc",
        "//gnmi/reconciler",
        "//proto/dataplane",
        "//proto/forwarding",
        "@com_github_openconfig_gnmi//proto/gnmi",
        "@com_github_openconfig_ygnmi//ygnmi",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//credentials/local",
    ],
)
