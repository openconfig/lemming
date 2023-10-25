load("@io_bazel_rules_go//go:def.bzl", "go_test")

go_test(
    name = "traffic_test_test",
    size = "enormous",
    srcs = ["traffic_test.go"],
    data = ["//integration_tests/onedut_oneotg_tests:topology_testbed"],
    tags = [
        "exclusive",
    ],
    deps = [
        "//internal/attrs",
        "//internal/binding",
        "@com_github_open_traffic_generator_snappi_gosnappi//:gosnappi",
        "@com_github_openconfig_gribi//v1/proto/service",
        "@com_github_openconfig_gribigo//chk",
        "@com_github_openconfig_gribigo//client",
        "@com_github_openconfig_gribigo//constants",
        "@com_github_openconfig_gribigo//fluent",
        "@com_github_openconfig_ondatra//:ondatra",
        "@com_github_openconfig_ondatra//gnmi",
        "@com_github_openconfig_ondatra//gnmi/otg/otgpath",
        "@com_github_openconfig_ygnmi//ygnmi",
    ],
)
