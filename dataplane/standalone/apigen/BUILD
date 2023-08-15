load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "apigen_lib",
    srcs = [
        "apigen.go",
        "ccgen.go",
    ],
    importpath = "github.com/openconfig/lemming/dataplane/standalone/apigen",
    visibility = ["//visibility:private"],
    deps = [
        "//dataplane/standalone/apigen/docparser",
        "//dataplane/standalone/apigen/protogen",
        "//dataplane/standalone/apigen/saiast",
        "@org_modernc_cc_v4//:cc",
    ],
)

go_binary(
    name = "apigen",
    embed = [":apigen_lib"],
    visibility = ["//visibility:public"],
)
