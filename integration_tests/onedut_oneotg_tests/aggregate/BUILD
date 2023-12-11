load("@io_bazel_rules_go//go:def.bzl", "go_test")

go_test(
    name = "aggregate_test",
    size = "enormous",
    srcs = ["aggregate_test.go"],
    data = [
        "testbed.pb.txt",
        "topology.pb.txt",
    ],
    tags = [
        "exclusive",
    ],
    deps = [
        "//internal/attrs",
        "//internal/binding",
        "@com_github_open_traffic_generator_snappi_gosnappi//:gosnappi",
        "@com_github_openconfig_ondatra//:ondatra",
        "@com_github_openconfig_ondatra//gnmi",
        "@com_github_openconfig_ondatra//gnmi/oc",
        "@com_github_openconfig_ondatra//gnmi/otg",
        "@com_github_openconfig_ygnmi//ygnmi",
        "@com_github_openconfig_ygot//ygot",
    ],
)
