load("@io_bazel_rules_go//go:def.bzl", "go_test")

go_test(
    name = "basictraffic_test",
    srcs = ["basic_traffic_test.go"],
    data = [":testbed.pb.txt"],
    deps = [
        "//internal/binding",
        "@com_github_openconfig_ondatra//:ondatra",
    ],
)
