load("@io_bazel_rules_go//go:def.bzl", "go_test")

go_test(
    name = "mpls_over_udp_scale_test",
    size = "enormous",
    srcs = ["mpls_over_udp_scale_test.go"],
    data = ["//integration_tests/onedut_oneotg_tests:topology_testbed"],
    deps = [
        "//gnmi/fakedevice",
        "//gnmi/oc",
        "//gnmi/oc/ocpath",
        "//integration_tests/onedut_oneotg_tests/mpls_over_udp_scale/util",
        "//internal/attrs",
        "//internal/binding",
        "@com_github_google_go_cmp//cmp",
        "@com_github_google_gopacket//:gopacket",
        "@com_github_google_gopacket//layers",
        "@com_github_google_gopacket//pcap",
        "@com_github_open_traffic_generator_snappi_gosnappi//:gosnappi",
        "@com_github_openconfig_gribi//v1/proto/service",
        "@com_github_openconfig_gribigo//constants",
        "@com_github_openconfig_gribigo//fluent",
        "@com_github_openconfig_ondatra//:ondatra",
        "@com_github_openconfig_ondatra//gnmi",
        "@com_github_openconfig_ondatra//otg",
        "@com_github_openconfig_ygnmi//ygnmi",
        "@com_github_openconfig_ygot//ygot",
    ],
)
