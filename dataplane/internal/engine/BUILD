load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "engine",
    srcs = [
        "engine.go",
        "helpers.go",
    ],
    importpath = "github.com/openconfig/lemming/dataplane/internal/engine",
    visibility = ["//dataplane:__subpackages__"],
    deps = [
        "//dataplane/forwarding",
        "//dataplane/forwarding/attributes",
        "//dataplane/forwarding/fwdconfig",
        "//dataplane/forwarding/infra/fwdcontext",
        "//proto/dataplane",
        "//proto/forwarding",
        "@com_github_golang_glog//:glog",
        "@org_golang_google_grpc//codes",
        "@org_golang_google_grpc//status",
    ],
)
