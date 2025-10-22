load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "protogen",
    srcs = ["protogen.go"],
    importpath = "github.com/openconfig/lemming/dataplane/apigen/protogen",
    visibility = ["//visibility:public"],
    deps = [
        "//dataplane/apigen/docparser",
        "//dataplane/apigen/saiast",
        "//dataplane/apigen/typeinfo",
        "//internal/lemmingutil",
    ],
)

go_test(
    name = "protogen_test",
    srcs = ["protogen_test.go"],
    embed = [":protogen"],
    deps = [
        "//dataplane/apigen/docparser",
        "//dataplane/apigen/saiast",
        "@com_github_google_go_cmp//cmp",
        "@openconfig_gnmi//errdiff",
    ],
)
