load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "lucius_lib",
    srcs = ["lucius.go"],
    cgo = True,
    importpath = "github.com/openconfig/lemming/dataplane/standalone/lucius",
    visibility = ["//visibility:private"],
    deps = [
        "//dataplane/internal/engine",
        "//proto/dataplane",
        "//proto/forwarding",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//credentials/insecure",
    ],
)

go_binary(
    name = "lucius_clib",
    embed = [":lucius_lib"],
    linkmode = "c-archive",
    visibility = ["//visibility:public"],
)