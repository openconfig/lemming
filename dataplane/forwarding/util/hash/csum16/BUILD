load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "csum16",
    srcs = ["csum16.go"],
    importpath = "github.com/openconfig/lemming/dataplane/forwarding/util/hash/csum16",
    visibility = ["//visibility:public"],
)

go_test(
    name = "csum16_test",
    size = "small",
    srcs = ["csum16_test.go"],
    embed = [":csum16"],
)
