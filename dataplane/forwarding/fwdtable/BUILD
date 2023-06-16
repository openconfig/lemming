load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "fwdtable",
    srcs = ["table.go"],
    importpath = "github.com/openconfig/lemming/dataplane/forwarding/fwdtable",
    visibility = ["//visibility:public"],
    deps = [
        "//dataplane/forwarding/fwdaction",
        "//dataplane/forwarding/infra/fwdcontext",
        "//dataplane/forwarding/infra/fwdobject",
        "//dataplane/forwarding/infra/fwdpacket",
        "//proto/forwarding",
    ],
)

go_test(
    name = "fwdtable_test",
    srcs = ["table_test.go"],
    embed = [":fwdtable"],
    deps = [
        "//dataplane/forwarding/fwdaction",
        "//dataplane/forwarding/infra/fwdcontext",
        "//dataplane/forwarding/infra/fwdobject",
        "//dataplane/forwarding/infra/fwdpacket",
        "//proto/forwarding",
    ],
)
