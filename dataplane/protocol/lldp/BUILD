load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "lldp",
    srcs = ["lldp.go"],
    importpath = "github.com/openconfig/lemming/dataplane/protocol/lldp",
    visibility = ["//visibility:public"],
    deps = [
        "//dataplane/proto/packetio",
        "//gnmi/gnmiclient",
        "//gnmi/oc",
        "//gnmi/oc/ocpath",
        "@com_github_golang_glog//:glog",
        "@com_github_google_gopacket//:gopacket",
        "@com_github_google_gopacket//layers",
        "@com_github_openconfig_ygnmi//ygnmi",
    ],
)

go_test(
    name = "lldp_test",
    srcs = ["lldp_test.go"],
    embed = [":lldp"],
    deps = [
        "//dataplane/proto/packetio",
        "@com_github_google_go_cmp//cmp",
        "@com_github_openconfig_gnmi//errdiff",
    ],
)
