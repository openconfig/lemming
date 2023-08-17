load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "flow",
    srcs = [
        "flow.go",
        "flowmap.go",
        "flowutil.go",
    ],
    importpath = "github.com/openconfig/lemming/dataplane/forwarding/fwdtable/flow",
    visibility = ["//visibility:public"],
    deps = [
        "//dataplane/forwarding/fwdaction",
        "//dataplane/forwarding/fwdtable",
        "//dataplane/forwarding/infra/fwdcontext",
        "//dataplane/forwarding/infra/fwdobject",
        "//dataplane/forwarding/infra/fwdpacket",
        "//dataplane/forwarding/infra/fwdset",
        "//proto/forwarding",
        "@org_golang_google_protobuf//proto",
    ],
)

go_test(
    name = "flow_test",
    size = "small",
    srcs = ["flow_test.go"],
    embed = [":flow"],
    deps = [
        "//dataplane/forwarding/fwdaction",
        "//dataplane/forwarding/fwdtable",
        "//dataplane/forwarding/fwdtable/mock_fwdpacket",
        "//dataplane/forwarding/infra/fwdcontext",
        "//dataplane/forwarding/infra/fwdobject",
        "//dataplane/forwarding/infra/fwdpacket",
        "//dataplane/forwarding/infra/fwdset",
        "//proto/forwarding",
        "@com_github_go_logr_logr//testr",
        "@org_uber_go_mock//gomock",
    ],
)
