load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "tabletestutil",
    srcs = ["testutil.go"],
    importpath = "github.com/openconfig/lemming/dataplane/forwarding/fwdtable/tabletestutil",
    visibility = ["//visibility:public"],
    deps = [
        "//dataplane/forwarding/fwdtable/tableutil",
        "//proto/forwarding",
        "@org_golang_google_protobuf//proto",
    ],
)
