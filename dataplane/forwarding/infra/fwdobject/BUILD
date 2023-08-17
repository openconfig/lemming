load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "fwdobject",
    srcs = [
        "object.go",
        "object_id.go",
        "object_table.go",
    ],
    importpath = "github.com/openconfig/lemming/dataplane/forwarding/infra/fwdobject",
    visibility = ["//visibility:public"],
    deps = [
        "//dataplane/forwarding/infra/fwdattribute",
        "//dataplane/forwarding/util/stats",
        "//proto/forwarding",
        "@com_github_golang_glog//:glog",
    ],
)

go_test(
    name = "fwdobject_test",
    size = "small",
    srcs = [
        "object_id_test.go",
        "object_table_test.go",
    ],
    embed = [":fwdobject"],
    deps = ["//proto/forwarding"],
)
