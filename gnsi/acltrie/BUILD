load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "acltrie",
    srcs = ["trie.go"],
    importpath = "github.com/openconfig/lemming/gnsi/acltrie",
    visibility = ["//visibility:public"],
    deps = [
        "@com_github_openconfig_gnmi//proto/gnmi",
        "@com_github_openconfig_gnsi//pathz",
        "@com_github_openconfig_ygot//util",
    ],
)

go_test(
    name = "acltrie_test",
    srcs = ["trie_test.go"],
    embed = [":acltrie"],
    deps = [
        "@com_github_google_go_cmp//cmp",
        "@com_github_openconfig_gnmi//errdiff",
        "@com_github_openconfig_gnmi//proto/gnmi",
        "@com_github_openconfig_gnsi//pathz",
        "@com_github_openconfig_ygot//ygot",
    ],
)
