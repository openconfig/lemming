load("@io_bazel_rules_go//go:def.bzl", "go_test")

go_test(
    name = "local_tests_test",
    size = "medium",
    srcs = [
        "policy_test.go",
        "prefix_set_test.go",
        "route_propagation_test.go",
        "session_establish_test.go",
        "ygnmi_test.go",
    ],
    deps = [
        "//:lemming",
        "//bgp",
        "//bgp/tests/proto/policyval",
        "//gnmi",
        "//gnmi/fakedevice",
        "//gnmi/gnmiclient",
        "//gnmi/oc",
        "//gnmi/oc/ocpath",
        "@com_github_openconfig_gnmi//proto/gnmi",
        "@com_github_openconfig_ygnmi//ygnmi",
        "@com_github_openconfig_ygot//ygot",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//codes",
        "@org_golang_google_grpc//credentials/insecure",
        "@org_golang_google_grpc//credentials/local",
        "@org_golang_google_grpc//status",
    ],
)
