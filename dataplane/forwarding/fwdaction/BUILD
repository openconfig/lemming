load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "fwdaction",
    srcs = ["action.go"],
    importpath = "github.com/openconfig/lemming/dataplane/forwarding/fwdaction",
    visibility = ["//visibility:public"],
    deps = [
        "//dataplane/forwarding/infra/fwdcontext",
        "//dataplane/forwarding/infra/fwdobject",
        "//dataplane/forwarding/infra/fwdpacket",
        "//proto/forwarding",
    ],
)

go_test(
    name = "fwdaction_test",
    srcs = ["action_test.go"],
    embed = [":fwdaction"],
    deps = [
        "//dataplane/forwarding/fwdaction/mock_fwdpacket",
        "//dataplane/forwarding/infra/fwdcontext",
        "//dataplane/forwarding/infra/fwdobject",
        "//dataplane/forwarding/infra/fwdpacket",
        "//proto/forwarding",
        "@com_github_go_logr_logr//testr",
        "@org_uber_go_mock//gomock",
    ],
)
