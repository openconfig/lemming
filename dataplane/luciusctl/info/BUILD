load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "info",
    srcs = ["info.go"],
    importpath = "github.com/openconfig/lemming/dataplane/luciusctl/info",
    visibility = ["//visibility:public"],
    deps = [
        "//proto/forwarding",
        "@com_github_spf13_cobra//:cobra",
        "@com_github_spf13_viper//:viper",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//credentials",
        "@org_golang_google_grpc//credentials/insecure",
    ],
)
