load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "protocol",
    srcs = [
        "attr.go",
        "doc.go",
        "handler.go",
        "packet.go",
        "parser.go",
        "util.go",
    ],
    importpath = "github.com/openconfig/lemming/dataplane/forwarding/protocol",
    visibility = ["//visibility:public"],
    deps = [
        "//dataplane/forwarding/infra/fwdattribute",
        "//dataplane/forwarding/infra/fwdpacket",
        "//dataplane/forwarding/util/frame",
        "//proto/forwarding",
        "@com_github_go_logr_logr//:logr",
        "@com_github_go_logr_logr//funcr",
        "@com_github_golang_glog//:glog",
    ],
)
