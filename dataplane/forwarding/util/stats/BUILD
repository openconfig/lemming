load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "stats",
    srcs = ["stats.go"],
    importpath = "github.com/openconfig/lemming/dataplane/forwarding/util/stats",
    visibility = ["//visibility:public"],
)

go_test(
    name = "stats_test",
    size = "small",
    srcs = ["stats_test.go"],
    embed = [":stats"],
)
