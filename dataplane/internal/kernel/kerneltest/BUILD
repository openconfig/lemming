load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "kerneltest",
    srcs = ["kernel.go"],
    importpath = "github.com/openconfig/lemming/dataplane/internal/kernel/kerneltest",
    visibility = ["//dataplane:__subpackages__"],
    deps = ["@com_github_vishvananda_netlink//:netlink"],
)

go_test(
    name = "kerneltest_test",
    srcs = ["kernel_test.go"],
    embed = [":kerneltest"],
    deps = [
        "@com_github_google_go_cmp//cmp",
        "@com_github_openconfig_gnmi//errdiff",
        "@com_github_vishvananda_netlink//:netlink",
    ],
)
