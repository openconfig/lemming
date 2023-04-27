load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "oc",
    srcs = [
        "enum.go",
        "enum_map.go",
        "schema.go",
        "structs-0.go",
        "union.go",
    ],
    importpath = "github.com/openconfig/lemming/gnmi/oc",
    visibility = ["//visibility:public"],
    deps = [
        "@com_github_openconfig_goyang//pkg/yang",
        "@com_github_openconfig_ygot//ygot",
        "@com_github_openconfig_ygot//ytypes",
    ],
)
