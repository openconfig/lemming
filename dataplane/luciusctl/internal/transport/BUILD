load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "transport",
    srcs = ["transport.go"],
    importpath = "github.com/openconfig/lemming/dataplane/luciusctl/internal/transport",
    visibility = ["//dataplane/luciusctl:__subpackages__"],
    deps = [
        "@com_github_spf13_viper//:viper",
        "@org_golang_google_grpc//:grpc",
        "@org_golang_google_grpc//credentials",
        "@org_golang_google_grpc//credentials/insecure",
    ],
)

go_test(
    name = "transport_test",
    size = "small",
    srcs = ["transport_test.go"],
    embed = [":transport"],
    deps = ["@com_github_spf13_viper//:viper"],
)
