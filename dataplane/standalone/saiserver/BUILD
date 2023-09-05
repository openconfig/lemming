load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "saiserver",
    srcs = ["saiserver.go"],
    importpath = "github.com/openconfig/lemming/dataplane/standalone/saiserver",
    visibility = ["//visibility:public"],
    deps = [
        "//dataplane/standalone/proto:sai",
        "@org_golang_google_grpc//:go_default_library",
    ],
)
