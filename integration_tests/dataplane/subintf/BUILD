load("@io_bazel_rules_go//go:def.bzl", "go_test")

go_test(
    name = "subintf_test",
    srcs = ["subintf_test.go"],
    data = [
        "testbed.pb.txt",
    ],
    deps = [
        "//gnmi/oc",
        "//integration_tests/saiutil",
        "//internal/attrs",
        "//internal/binding",
        "@com_github_google_gopacket//:gopacket",
        "@com_github_google_gopacket//layers",
        "@com_github_openconfig_ondatra//:ondatra",
        "@org_golang_google_protobuf//proto",
    ],
)
