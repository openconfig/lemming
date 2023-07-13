load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "fwdpacket",
    srcs = [
        "field.go",
        "packet.go",
    ],
    importpath = "github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket",
    visibility = ["//visibility:public"],
    deps = [
        "//dataplane/forwarding/infra/fwdattribute",
        "//dataplane/forwarding/infra/fwdobject",
        "//dataplane/forwarding/util/frame",
        "//proto/forwarding",
        "@com_github_go_logr_logr//:logr",
        "@com_github_golang_glog//:glog",
    ],
)

go_test(
    name = "fwdpacket_test",
    srcs = ["field_test.go"],
    embed = [":fwdpacket"],
    deps = ["//proto/forwarding"],
)
