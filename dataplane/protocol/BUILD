load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "protocol",
    srcs = ["registry.go"],
    importpath = "github.com/openconfig/lemming/dataplane/protocol",
    visibility = ["//visibility:public"],
    deps = [
        "//dataplane/forwarding/util/queue",
        "//dataplane/kernel",
        "//dataplane/proto/packetio",
        "@com_github_golang_glog//:glog",
        "@org_golang_google_genproto_googleapis_rpc//status",
        "@org_golang_google_grpc//codes",
    ],
)

go_test(
    name = "registry_test",
    srcs = ["registry_test.go"],
    embed = [":protocol"],
    deps = [
        "//dataplane/kernel",
        "//dataplane/proto/packetio",
        "@com_github_google_go_cmp//cmp",
        "@com_github_openconfig_gnmi//errdiff",
        "@com_github_golang_glog//:glog",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//codes",
        "@org_golang_google_grpc//credentials/insecure",
        "@org_golang_google_protobuf//testing/protocmp",
    ],
)
