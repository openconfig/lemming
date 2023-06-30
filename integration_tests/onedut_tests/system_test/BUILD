load("@io_bazel_rules_go//go:def.bzl", "go_test")

go_test(
    name = "system_test_test",
    size = "enormous",
    srcs = ["system_test.go"],
    data = ["//integration_tests/onedut_tests:topology_testbed"],
    tags = [
        "exclusive",
    ],
    deps = [
        "//gnmi/oc/ocpath",
        "//internal/binding",
        "@com_github_openconfig_ondatra//:ondatra",
        "@com_github_openconfig_ondatra//gnmi",
    ],
)
