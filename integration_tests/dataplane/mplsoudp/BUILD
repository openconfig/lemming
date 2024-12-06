load("@io_bazel_rules_go//go:def.bzl", "go_test")

go_test(
    name = "mplsoudp_test",
    srcs = ["mplsoverudp_test.go"],
    data = [
        "testbed.pb.txt",
    ],
    deps = [
        "//dataplane/proto/sai",
        "//dataplane/saiserver",
        "//gnmi/oc",
        "//integration_tests/saiutil",
        "//internal/attrs",
        "//internal/binding",
        "//proto/forwarding",
        "@com_github_google_go_cmp//cmp",
        "@com_github_google_go_cmp//cmp/cmpopts",
        "@com_github_google_gopacket//:gopacket",
        "@com_github_google_gopacket//layers",
        "@com_github_openconfig_ondatra//:ondatra",
        "@com_github_openconfig_ondatra//binding",
        "@com_github_openconfig_ygot//ygot",
        "@org_golang_google_grpc//:grpc",
        "@org_golang_google_protobuf//proto",
    ],
)
