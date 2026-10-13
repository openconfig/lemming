load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "luciusctl",
    srcs = ["luciusctl.go"],
    importpath = "github.com/openconfig/lemming/dataplane/luciusctl",
    visibility = ["//visibility:public"],
    deps = [
        "//dataplane/luciusctl/info",
        "//dataplane/luciusctl/sai",
        "@com_github_spf13_cobra//:cobra",
        "@com_github_spf13_viper//:viper",
    ],
)

go_test(
    name = "luciusctl_test",
    size = "small",
    srcs = ["luciusctl_test.go"],
    embed = [":luciusctl"],
    deps = [
        "//dataplane/proto/sai",
        "//proto/forwarding",
        "@com_github_spf13_cobra//:cobra",
        "@com_github_spf13_viper//:viper",
        "@org_golang_google_grpc//:grpc",
        "@org_golang_google_grpc//credentials",
        "@org_golang_google_grpc//peer",
        "@org_golang_google_grpc//reflection",
    ],
)
