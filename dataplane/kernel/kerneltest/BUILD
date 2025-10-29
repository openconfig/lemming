load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "kerneltest",
    srcs = ["kernel.go"],
    importpath = "github.com/openconfig/lemming/dataplane/kernel/kerneltest",
    visibility = ["//dataplane:__subpackages__"],
    deps = select({
        "@io_bazel_rules_go//go/platform:android": [
            "@com_github_vishvananda_netlink//:netlink",
        ],
        "@io_bazel_rules_go//go/platform:linux": [
            "@com_github_vishvananda_netlink//:netlink",
        ],
        "//conditions:default": [],
    }),
)

go_test(
    name = "kerneltest_test",
    size = "small",
    srcs = ["kernel_test.go"],
    embed = [":kerneltest"],
    deps = select({
        "@io_bazel_rules_go//go/platform:android": [
            "@com_github_google_go_cmp//cmp",
            "@com_github_vishvananda_netlink//:netlink",
            "@openconfig_gnmi//errdiff",
        ],
        "@io_bazel_rules_go//go/platform:linux": [
            "@com_github_google_go_cmp//cmp",
            "@com_github_vishvananda_netlink//:netlink",
            "@openconfig_gnmi//errdiff",
        ],
        "//conditions:default": [],
    }),
)
