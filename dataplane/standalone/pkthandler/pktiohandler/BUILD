load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "pktiohandler",
    srcs = ["pktiohandler.go"],
    importpath = "github.com/openconfig/lemming/dataplane/standalone/pkthandler/pktiohandler",
    visibility = ["//visibility:public"],
    deps = [
        "//dataplane/forwarding/util/queue",
        "//dataplane/internal/kernel",
        "//proto/forwarding",
        "@com_github_golang_glog//:glog",
        "@org_golang_google_genproto_googleapis_rpc//status",
        "@org_golang_google_grpc//codes",
    ],
)
