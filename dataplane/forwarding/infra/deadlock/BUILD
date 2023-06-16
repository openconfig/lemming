load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "deadlock",
    srcs = ["deadlock.go"],
    importpath = "github.com/openconfig/lemming/dataplane/forwarding/infra/deadlock",
    visibility = ["//visibility:public"],
    deps = ["@com_github_golang_glog//:glog"],
)
