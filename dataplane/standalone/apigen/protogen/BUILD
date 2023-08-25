load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "protogen",
    srcs = ["protogen.go"],
    importpath = "github.com/openconfig/lemming/dataplane/standalone/apigen/protogen",
    visibility = ["//visibility:public"],
    deps = [
        "//dataplane/standalone/apigen/docparser",
        "//dataplane/standalone/apigen/saiast",
        "@com_github_golang_glog//:glog",
        "@com_github_stoewer_go_strcase//:go-strcase",
        "@org_golang_x_exp//maps",
        "@org_golang_x_exp//slices",
    ],
)

go_test(
    name = "protogen_test",
    srcs = ["protogen_test.go"],
    embed = [":protogen"],
    deps = [
        "//dataplane/standalone/apigen/docparser",
        "//dataplane/standalone/apigen/saiast",
        "@com_github_google_go_cmp//cmp",
        "@com_github_openconfig_gnmi//errdiff",
    ],
)
