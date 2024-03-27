load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "fwdcontext",
    srcs = ["context.go"],
    importpath = "github.com/openconfig/lemming/dataplane/forwarding/infra/fwdcontext",
    visibility = ["//visibility:public"],
    deps = [
        "//dataplane/forwarding/infra/deadlock",
        "//dataplane/forwarding/infra/fwdattribute",
        "//dataplane/forwarding/infra/fwdobject",
        "//dataplane/forwarding/util/queue",
        "//dataplane/proto/packetio",
        "//proto/forwarding",
        "@com_github_golang_glog//:glog",
        "@com_github_google_gopacket//:gopacket",
    ],
)
