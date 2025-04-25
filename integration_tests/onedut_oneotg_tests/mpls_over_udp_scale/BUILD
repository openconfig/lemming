load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "mpls_over_udp_scale",
    srcs = ["util.go"],
    importpath = "github.com/openconfig/lemming/integration_tests/onedut_oneotg_tests/mpls_over_udp_scale",
    visibility = ["//visibility:public"],
    deps = ["@com_github_openconfig_gribigo//fluent"],
)

go_test(
    name = "mpls_over_udp_scale_test",
    srcs = ["util_test.go"],
    embed = [":mpls_over_udp_scale"],
    deps = [
        "@com_github_google_go_cmp//cmp",
        "@com_github_openconfig_gribi//v1/proto/service",
        "@com_github_openconfig_gribigo//fluent",
        "@org_golang_google_protobuf//testing/protocmp",
    ],
)
