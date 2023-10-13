load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "action",
    srcs = ["action.go"],
    importpath = "github.com/openconfig/lemming/dataplane/forwarding/fwdtable/action",
    visibility = ["//visibility:public"],
    deps = [
        "//dataplane/forwarding/fwdaction",
        "//dataplane/forwarding/fwdtable",
        "//dataplane/forwarding/infra/fwdcontext",
        "//dataplane/forwarding/infra/fwdobject",
        "//dataplane/forwarding/infra/fwdpacket",
        "//proto/forwarding",
    ],
)
