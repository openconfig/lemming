load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "pktiohandler",
    srcs = ["pktiohandler.go"],
    importpath = "github.com/openconfig/lemming/dataplane/standalone/pkthandler/pktiohandler",
    visibility = ["//visibility:public"],
    deps = [
        "//dataplane/forwarding/util/queue",
        "//dataplane/kernel",
        "//dataplane/proto/packetio",
        "@com_github_golang_glog//:glog",
        "@com_github_vishvananda_netlink//:netlink",
        "@org_golang_google_genproto_googleapis_rpc//status",
        "@org_golang_google_grpc//codes",
    ],
)

go_test(
    name = "pktiohandler_test",
    srcs = ["pktiohandler_test.go"],
    embed = [":pktiohandler"],
    deps = [
        "//dataplane/kernel",
        "//dataplane/proto/packetio",
        "@com_github_google_go_cmp//cmp",
        "@com_github_vishvananda_netlink//:netlink",
        "@org_golang_google_grpc//codes",
        "@org_golang_google_protobuf//testing/protocmp",
    ],
)
