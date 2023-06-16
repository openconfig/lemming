load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "bridge",
    srcs = [
        "bridge.go",
        "bridge_learn.go",
    ],
    importpath = "github.com/openconfig/lemming/dataplane/forwarding/fwdtable/bridge",
    visibility = ["//visibility:public"],
    deps = [
        "//dataplane/forwarding/fwdaction",
        "//dataplane/forwarding/fwdport",
        "//dataplane/forwarding/fwdtable",
        "//dataplane/forwarding/fwdtable/exact",
        "//dataplane/forwarding/infra/fwdcontext",
        "//dataplane/forwarding/infra/fwdobject",
        "//dataplane/forwarding/infra/fwdpacket",
        "//dataplane/forwarding/util/queue",
        "//proto/forwarding",
        "@com_github_golang_glog//:glog",
    ],
)

go_test(
    name = "bridge_test",
    srcs = ["bridge_test.go"],
    embed = [":bridge"],
    deps = [
        "//dataplane/forwarding/fwdaction",
        "//dataplane/forwarding/fwdaction/actions",
        "//dataplane/forwarding/fwdport",
        "//dataplane/forwarding/fwdport/ports",
        "//dataplane/forwarding/fwdport/porttestutil",
        "//dataplane/forwarding/fwdtable",
        "//dataplane/forwarding/fwdtable/mock_fwdpacket",
        "//dataplane/forwarding/infra/fwdattribute",
        "//dataplane/forwarding/infra/fwdcontext",
        "//dataplane/forwarding/infra/fwdobject",
        "//dataplane/forwarding/infra/fwdpacket",
        "//dataplane/forwarding/protocol",
        "//proto/forwarding",
        "@com_github_golang_mock//gomock",
    ],
)
