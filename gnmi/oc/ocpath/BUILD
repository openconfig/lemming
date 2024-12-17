load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "ocpath",
    srcs = ["ocpath.go"],
    importpath = "github.com/openconfig/lemming/gnmi/oc/ocpath",
    visibility = ["//visibility:public"],
    deps = [
        "//gnmi/oc",
        "//gnmi/oc/acl",
        "//gnmi/oc/bgpgue",
        "//gnmi/oc/definedsets",
        "//gnmi/oc/interfaces",
        "//gnmi/oc/keychain",
        "//gnmi/oc/lacp",
        "//gnmi/oc/lldp",
        "//gnmi/oc/networkinstance",
        "//gnmi/oc/platform",
        "//gnmi/oc/qos",
        "//gnmi/oc/relayagent",
        "//gnmi/oc/routingpolicy",
        "//gnmi/oc/sampling",
        "//gnmi/oc/system",
        "@com_github_openconfig_ygnmi//ygnmi",
        "@com_github_openconfig_ygot//ygot",
        "@com_github_openconfig_ygot//ytypes",
    ],
)
