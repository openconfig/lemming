load("@io_bazel_rules_go//go:def.bzl", "go_test")

go_test(
    name = "complete_chassis_reboot_test",
    size = "enormous",
    srcs = ["complete_chassis_reboot_test.go"],
    data = ["//integration_tests/onedut_tests:topology_testbed"],
    deps = [
        "//internal/binding",
        "@com_github_google_go_cmp//cmp",
        "@com_github_openconfig_gnoi//system",
        "@com_github_openconfig_ondatra//:ondatra",
        "@com_github_openconfig_ondatra//gnmi",
        "@com_github_openconfig_ondatra//gnmi/oc",
        "@com_github_openconfig_testt//:testt",
    ],
)
