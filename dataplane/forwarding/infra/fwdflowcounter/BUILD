load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "fwdflowcounter",
    srcs = ["flow_counter.go"],
    importpath = "github.com/openconfig/lemming/dataplane/forwarding/infra/fwdflowcounter",
    visibility = ["//visibility:public"],
    deps = [
        "//dataplane/forwarding/infra/fwdcontext",
        "//dataplane/forwarding/infra/fwdobject",
        "//proto/forwarding",
    ],
)

go_test(
    name = "fwdflowcounter_test",
    srcs = ["flow_counter_test.go"],
    embed = [":fwdflowcounter"],
    deps = [
        "//dataplane/forwarding/infra/fwdcontext",
        "//proto/forwarding",
    ],
)
