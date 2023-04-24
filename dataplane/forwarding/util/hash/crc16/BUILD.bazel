load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "crc16",
    srcs = ["crc16.go"],
    importpath = "github.com/openconfig/lemming/dataplane/forwarding/util/hash/crc16",
    visibility = ["//visibility:public"],
    deps = ["//dataplane/forwarding/util/hash/hash16"],
)

go_test(
    name = "crc16_test",
    srcs = ["crc16_test.go"],
    embed = [":crc16"],
    deps = ["//dataplane/forwarding/util/hash/hash16"],
)
