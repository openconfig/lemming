load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "cpusink_lib",
    srcs = ["main.go"],
    cgo = True,
    importpath = "github.com/openconfig/lemming/dataplane/standalone/cpusink",
    visibility = ["//visibility:private"],
    deps = [
        "//dataplane/standalone/cpusink/sink",
        "//proto/forwarding",
        "@com_github_golang_glog//:glog",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//credentials/insecure",
    ],
)

go_binary(
    name = "cpusink",
    embed = [":cpusink_lib"],
    linkmode = "c-archive",
    visibility = ["//visibility:public"],
)
