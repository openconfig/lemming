load("@io_bazel_rules_go//go:def.bzl", "go_test")

go_test(
    name = "bgp_route_propagation_test",
    size = "enormous",
    srcs = ["bgp_route_propagation_test.go"],
    data = ["//integration_tests/onedut_oneotg_tests:topology_testbed"],
    tags = [
        "exclusive",
    ],
    deps = [
        "//gnmi/fakedevice",
        "//internal/attrs",
        "//internal/binding",
        "@com_github_open_traffic_generator_snappi_gosnappi//:gosnappi",
        "@com_github_openconfig_ondatra//:ondatra",
        "@com_github_openconfig_ondatra//gnmi",
        "@com_github_openconfig_ondatra//gnmi/oc",
        "@com_github_openconfig_ondatra//gnmi/otg",
        "@com_github_openconfig_ondatra//otg",
        "@com_github_openconfig_ygnmi//ygnmi",
        "@com_github_openconfig_ygot//ygot",
    ],
)
