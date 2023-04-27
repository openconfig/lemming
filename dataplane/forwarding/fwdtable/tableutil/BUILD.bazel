load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "tableutil",
    srcs = ["util.go"],
    importpath = "github.com/openconfig/lemming/dataplane/forwarding/fwdtable/tableutil",
    visibility = ["//visibility:public"],
    deps = [
        "//dataplane/forwarding/infra/fwdpacket",
        "//proto/forwarding",
    ],
)
