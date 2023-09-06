load("@io_bazel_rules_go//go:def.bzl", "go_test")

go_test(
    name = "prefix_set_test",
    size = "enormous",
    srcs = ["prefix_set_test.go"],
    data = ["//integration_tests/dut_policy_tests:topology_testbed"],
    tags = [
        "exclusive",
    ],
    deps = [
        "//bgp/tests/proto/policyval",
        "//internal/binding",
        "//policytest",
        "@com_github_openconfig_ondatra//:ondatra",
        "@com_github_openconfig_ondatra//gnmi",
        "@com_github_openconfig_ondatra//gnmi/oc",
    ],
)
