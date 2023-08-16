load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "queue",
    srcs = ["queue.go"],
    importpath = "github.com/openconfig/lemming/dataplane/forwarding/util/queue",
    visibility = ["//visibility:public"],
    deps = ["@com_github_golang_glog//:glog"],
)

go_test(
    name = "queue_test",
    size = "small",
    srcs = ["queue_test.go"],
    embed = [":queue"],
)
