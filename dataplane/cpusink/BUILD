load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "cpusink",
    srcs = ["cpusink.go"],
    importpath = "github.com/openconfig/lemming/dataplane/cpusink",
    visibility = ["//visibility:public"],
    deps = [
        "//dataplane/dplaneopts",
        "//dataplane/forwarding/attributes",
        "//dataplane/forwarding/fwdconfig",
        "//dataplane/kernel",
        "//proto/forwarding",
        "@com_github_golang_glog//:glog",
        "@com_github_vishvananda_netlink//:netlink",
        "@org_golang_x_exp//maps",
    ],
)
