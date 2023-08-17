load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "fwdset",
    srcs = ["set.go"],
    importpath = "github.com/openconfig/lemming/dataplane/forwarding/infra/fwdset",
    visibility = ["//visibility:public"],
    deps = [
        "//dataplane/forwarding/infra/fwdcontext",
        "//dataplane/forwarding/infra/fwdobject",
        "//proto/forwarding",
    ],
)

go_test(
    name = "fwdset_test",
    size = "small",
    srcs = ["set_test.go"],
    embed = [":fwdset"],
    deps = [
        "//dataplane/forwarding/infra/fwdcontext",
        "//proto/forwarding",
    ],
)
