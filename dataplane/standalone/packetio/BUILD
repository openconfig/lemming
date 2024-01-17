load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_binary(
    name = "packetio",
    embed = [":packetio_lib"],
    linkmode = "c-archive",
    visibility = ["//visibility:public"],
)

go_library(
    name = "packetio_lib",
    srcs = ["main.go"],
    cgo = True,
    importpath = "github.com/openconfig/lemming/dataplane/standalone/packetio",
    visibility = ["//visibility:private"],
    deps = [
        "//dataplane/cpusink",
        "//proto/forwarding",
        "@com_github_golang_glog//:glog",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//credentials/insecure",
    ],
)
