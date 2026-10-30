load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "udp",
    srcs = ["udp.go"],
    importpath = "github.com/openconfig/lemming/dataplane/forwarding/protocol/udp",
    visibility = ["//visibility:public"],
    deps = [
        "//dataplane/forwarding/infra/fwdpacket",
        "//dataplane/forwarding/protocol",
        "//dataplane/forwarding/protocol/vxlan",
        "//dataplane/forwarding/util/frame",
        "//dataplane/forwarding/util/hash/csum16",
        "//proto/forwarding",
    ],
)

go_test(
    name = "udp_test",
    srcs = ["udp_test.go"],
    embed = [":udp"],
    deps = [
        "//dataplane/forwarding/infra/fwdpacket",
        "//dataplane/forwarding/protocol",
        "//dataplane/forwarding/protocol/ethernet",
        "//dataplane/forwarding/protocol/ip",
        "//dataplane/forwarding/protocol/metadata",
        "//dataplane/forwarding/protocol/opaque",
        "//dataplane/forwarding/protocol/tcp",
        "//dataplane/forwarding/protocol/vxlan",
        "//dataplane/forwarding/util/frame",
        "//proto/forwarding",
        "@com_github_google_go_cmp//cmp",
    ],
)
