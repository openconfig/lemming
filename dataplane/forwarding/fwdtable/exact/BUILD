load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "exact",
    srcs = [
        "exact.go",
        "exactutil.go",
    ],
    importpath = "github.com/openconfig/lemming/dataplane/forwarding/fwdtable/exact",
    visibility = ["//visibility:public"],
    deps = [
        "//dataplane/forwarding/fwdaction",
        "//dataplane/forwarding/fwdtable",
        "//dataplane/forwarding/fwdtable/tableutil",
        "//dataplane/forwarding/infra/fwdcontext",
        "//dataplane/forwarding/infra/fwdobject",
        "//dataplane/forwarding/infra/fwdpacket",
        "//proto/forwarding",
    ],
)

go_test(
    name = "exact_test",
    srcs = [
        "exact_test.go",
        "exactutil_test.go",
    ],
    embed = [":exact"],
    deps = [
        "//dataplane/forwarding/fwdaction/actions",
        "//dataplane/forwarding/fwdtable",
        "//dataplane/forwarding/fwdtable/mock_fwdpacket",
        "//dataplane/forwarding/fwdtable/tabletestutil",
        "//dataplane/forwarding/infra/fwdcontext",
        "//dataplane/forwarding/infra/fwdobject",
        "//dataplane/forwarding/infra/fwdpacket",
        "//proto/forwarding",
        "@com_github_go_logr_logr//testr",
        "@org_golang_google_protobuf//proto",
        "@org_uber_go_mock//gomock",
    ],
)
