load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "vxlan",
    srcs = ["vxlan.go"],
    importpath = "github.com/openconfig/lemming/dataplane/forwarding/protocol/vxlan",
    visibility = ["//visibility:public"],
    deps = [
        "//dataplane/forwarding/infra/fwdpacket",
        "//dataplane/forwarding/protocol",
        "//dataplane/forwarding/util/frame",
        "//proto/forwarding",
    ],
)

go_test(
    name = "vxlan_test",
    srcs = ["vxlan_test.go"],
    embed = [":vxlan"],
    deps = [
        "//dataplane/forwarding/infra/fwdpacket",
        "//dataplane/forwarding/protocol",
        "//dataplane/forwarding/util/frame",
        "//proto/forwarding",
        "@com_github_google_go_cmp//cmp",
    ],
)
