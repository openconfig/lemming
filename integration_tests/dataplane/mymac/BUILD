load("@io_bazel_rules_go//go:def.bzl", "go_test")

go_test(
    name = "mymac_test",
    srcs = ["mymac_test.go"],
    data = [":testbed.pb.txt"],
    deps = [
        "//dataplane/proto/sai",
        "//internal/attrs",
        "//internal/binding",
        "@com_github_open_traffic_generator_snappi_gosnappi//:gosnappi",
        "@com_github_openconfig_ondatra//:ondatra",
        "@com_github_openconfig_ondatra//binding",
        "@com_github_openconfig_ondatra//gnmi",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_protobuf//proto",
    ],
)
