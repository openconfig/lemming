load("@io_bazel_rules_go//go:def.bzl", "go_test")

go_test(
    name = "vlan_subintf_test",
    size = "enormous",
    srcs = ["vlan_subintf_test.go"],
    data = ["//integration_tests/onedut_oneotg_tests:topology_testbed"],
    deps = [
        "//gnmi/fakedevice",
        "//gnmi/oc",
        "//gnmi/oc/ocpath",
        "//internal/attrs",
        "//internal/binding",
        "@com_github_open_traffic_generator_snappi_gosnappi//:gosnappi",
        "@com_github_openconfig_gribigo//chk",
        "@com_github_openconfig_gribigo//client",
        "@com_github_openconfig_gribigo//constants",
        "@com_github_openconfig_gribigo//fluent",
        "@com_github_openconfig_ondatra//:ondatra",
        "@com_github_openconfig_ondatra//gnmi",
        "@com_github_openconfig_ondatra//gnmi/otg/otgpath",
        "@com_github_openconfig_ondatra//otg",
        "@com_github_openconfig_ygnmi//ygnmi",
        "@com_github_openconfig_ygot//ygot",
    ],
)
