load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "ccgen",
    srcs = [
        "ccgen.go",
        "servergen.go",
    ],
    importpath = "github.com/openconfig/lemming/dataplane/apigen/ccgen",
    visibility = ["//visibility:public"],
    deps = [
        "//dataplane/apigen/docparser",
        "//dataplane/apigen/saiast",
        "//dataplane/apigen/typeinfo",
    ],
)

go_test(
    name = "ccgen_test",
    srcs = ["ccgen_test.go"],
    embed = [":ccgen"],
    deps = [
        "//dataplane/apigen/docparser",
        "//dataplane/apigen/saiast",
    ],
)
