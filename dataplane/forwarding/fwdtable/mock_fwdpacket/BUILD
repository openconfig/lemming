load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "mock_fwdpacket",
    srcs = [
        "gen.go",
        "mock.go",
    ],
    importpath = "github.com/openconfig/lemming/dataplane/forwarding/fwdtable/mock_fwdpacket",
    visibility = ["//visibility:public"],
    deps = [
        "//dataplane/forwarding/infra/fwdattribute",
        "//dataplane/forwarding/infra/fwdpacket",
        "//proto/forwarding",
        "@com_github_go_logr_logr//:logr",
        "@org_uber_go_mock//gomock",
    ],
)
