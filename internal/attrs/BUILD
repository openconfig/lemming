load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "attrs",
    srcs = ["attrs.go"],
    importpath = "github.com/openconfig/lemming/internal/attrs",
    visibility = ["//:__subpackages__"],
    deps = [
        "@com_github_open_traffic_generator_snappi_gosnappi//:gosnappi",
        "@com_github_openconfig_ondatra//:ondatra",
        "@com_github_openconfig_ondatra//gnmi/oc",
        "@com_github_openconfig_ygot//ygot",
    ],
)
