load("@io_bazel_rules_go//go:def.bzl", "go_test")

go_test(
    name = "gnoi_file_test",
    size = "enormous",
    srcs = ["gnoi_file_test.go"],
    data = ["//integration_tests/onedut_tests:topology_testbed"],
    deps = [
        "//internal/binding",
        "@com_github_openconfig_gnoi//common",
        "@com_github_openconfig_gnoi//file",
        "@com_github_openconfig_gnoi//types",
        "@com_github_openconfig_ondatra//:ondatra",
    ],
)
