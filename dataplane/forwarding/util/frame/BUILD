load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "frame",
    srcs = ["frame.go"],
    importpath = "github.com/openconfig/lemming/dataplane/forwarding/util/frame",
    visibility = ["//visibility:public"],
)

go_test(
    name = "frame_test",
    srcs = ["frame_test.go"],
    embed = [":frame"],
)
