load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "porttestutil",
    srcs = ["testutil.go"],
    importpath = "github.com/openconfig/lemming/dataplane/forwarding/fwdport/porttestutil",
    visibility = ["//visibility:public"],
    deps = [
        "//dataplane/forwarding/fwdport",
        "//dataplane/forwarding/infra/fwdcontext",
        "//dataplane/forwarding/infra/fwdobject",
        "//proto/forwarding",
    ],
)
