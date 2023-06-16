load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "fwdport",
    srcs = ["port.go"],
    importpath = "github.com/openconfig/lemming/dataplane/forwarding/fwdport",
    visibility = ["//visibility:public"],
    deps = [
        "//dataplane/forwarding/fwdaction",
        "//dataplane/forwarding/infra/fwdattribute",
        "//dataplane/forwarding/infra/fwdcontext",
        "//dataplane/forwarding/infra/fwdobject",
        "//dataplane/forwarding/infra/fwdpacket",
        "//dataplane/forwarding/protocol",
        "//proto/forwarding",
        "@com_github_golang_glog//:glog",
    ],
)

go_test(
    name = "fwdport_test",
    srcs = ["port_test.go"],
    embed = [":fwdport"],
    deps = [
        "//dataplane/forwarding/fwdaction",
        "//dataplane/forwarding/infra/fwdcontext",
        "//dataplane/forwarding/infra/fwdobject",
        "//dataplane/forwarding/infra/fwdpacket",
        "//proto/forwarding",
    ],
)
