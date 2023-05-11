load("@io_bazel_rules_go//go:def.bzl", "go_test")

go_test(
    name = "bgp_establish_test",
    srcs = ["bgp_establish_test.go"],
    data = ["//integration_tests/twodut_tests:topology_testbed"],
    deps = [
        "//gnmi/fakedevice",
        "//gnmi/oc",
        "//gnmi/oc/ocpath",
        "//internal/binding",
        "@com_github_openconfig_ondatra//:ondatra",
        "@com_github_openconfig_ondatra//gnmi",
        "@com_github_openconfig_ygot//ygot",
    ],
)
