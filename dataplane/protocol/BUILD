load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "protocol",
    srcs = ["registry.go"],
    importpath = "github.com/openconfig/lemming/dataplane/protocol",
    visibility = ["//visibility:public"],
    deps = [
        "//dataplane/forwarding/util/queue",
        "//dataplane/proto/packetio",
        "@com_github_golang_glog//:glog",
    ],
)

go_test(
    name = "protocol_test",
    srcs = ["registry_test.go"],
    embed = [":protocol"],
    deps = [
        "//dataplane/proto/packetio",
        "@com_github_google_go_cmp//cmp",
        "@com_github_openconfig_gnmi//errdiff",
        "@org_golang_google_grpc//:grpc",
        "@org_golang_google_protobuf//testing/protocmp",
    ],
)
