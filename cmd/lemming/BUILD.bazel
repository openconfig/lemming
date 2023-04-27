load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "lemming_lib",
    srcs = ["lemming.go"],
    importpath = "github.com/openconfig/lemming/cmd/lemming",
    visibility = ["//visibility:private"],
    deps = [
        "//:lemming",
        "//sysrib",
        "@com_github_golang_glog//:glog",
        "@com_github_spf13_pflag//:pflag",
        "@com_github_spf13_viper//:viper",
        "@org_golang_google_grpc//credentials",
        "@org_golang_google_grpc//credentials/insecure",
    ],
)

go_binary(
    name = "lemming",
    embed = [":lemming_lib"],
    visibility = ["//visibility:public"],
)
