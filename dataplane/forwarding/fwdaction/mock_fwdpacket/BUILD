load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "mock_fwdpacket",
    srcs = [
        "gen.go",
        "mock.go",
    ],
    importpath = "github.com/openconfig/lemming/dataplane/forwarding/fwdaction/mock_fwdpacket",
    visibility = ["//visibility:public"],
    deps = [
        "//dataplane/forwarding/infra/fwdattribute",
        "//dataplane/forwarding/infra/fwdpacket",
        "//proto/forwarding",
        "@com_github_golang_mock//gomock",
    ],
)
