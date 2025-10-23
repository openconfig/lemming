load("@io_bazel_rules_go//go:def.bzl", "go_test")

go_test(
    name = "gnsi_test_test",
    size = "enormous",
    srcs = ["gnsi_test.go"],
    data = ["//integration_tests/onedut_tests:topology_testbed"],
    deps = [
        "//gnmi/oc",
        "//gnmi/oc/ocpath",
        "//internal/binding",
        "@com_github_openconfig_gnsi//pathz",
        "@com_github_openconfig_ondatra//:ondatra",
        "@com_github_openconfig_ondatra//gnmi",
        "@com_github_openconfig_ygnmi//ygnmi",
        "@com_github_openconfig_ygot//ygot",
        "@openconfig_gnmi//errdiff",
        "@openconfig_gnmi//proto/gnmi",
        "@org_golang_google_grpc//metadata",
    ],
)
