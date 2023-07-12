load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "prefix",
    srcs = [
        "prefix.go",
        "prefixkey.go",
        "prefixutil.go",
    ],
    importpath = "github.com/openconfig/lemming/dataplane/forwarding/fwdtable/prefix",
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
    name = "prefix_test",
    srcs = [
        "prefix_test.go",
        "prefixkey_test.go",
        "prefixutil_test.go",
    ],
    embed = [":prefix"],
    deps = [
        "//dataplane/forwarding/fwdaction",
        "//dataplane/forwarding/fwdaction/actions",
        "//dataplane/forwarding/fwdtable",
        "//dataplane/forwarding/fwdtable/mock_fwdpacket",
        "//dataplane/forwarding/fwdtable/tabletestutil",
        "//dataplane/forwarding/infra/fwdcontext",
        "//dataplane/forwarding/infra/fwdobject",
        "//dataplane/forwarding/infra/fwdpacket",
        "//proto/forwarding",
        "@com_github_golang_glog//:glog",
        "@org_golang_google_protobuf//proto",
        "@org_uber_go_mock//gomock",
    ],
)
