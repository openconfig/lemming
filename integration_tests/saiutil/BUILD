load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "saiutil",
    srcs = ["saiutil.go"],
    importpath = "github.com/openconfig/lemming/integration_tests/saiutil",
    visibility = ["//visibility:public"],
    deps = [
        "//dataplane/proto/sai",
        "@com_github_openconfig_ondatra//:ondatra",
        "@com_github_openconfig_ondatra//binding",
        "@org_golang_google_grpc//:grpc",
        "@org_golang_google_protobuf//proto",
    ],
)
