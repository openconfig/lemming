load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "ip",
    srcs = [
        "gre.go",
        "ip.go",
        "ip4.go",
        "ip6.go",
    ],
    importpath = "github.com/openconfig/lemming/dataplane/forwarding/protocol/ip",
    visibility = ["//visibility:public"],
    deps = [
        "//dataplane/forwarding/infra/fwdpacket",
        "//dataplane/forwarding/protocol",
        "//dataplane/forwarding/protocol/ethernet",
        "//dataplane/forwarding/util/frame",
        "//dataplane/forwarding/util/hash/csum16",
        "//proto/forwarding",
    ],
)

go_test(
    name = "ip_test",
    srcs = ["ip_test.go"],
    embed = [":ip"],
    deps = [
        "//dataplane/forwarding/util/frame",
        "//proto/forwarding",
    ],
)
