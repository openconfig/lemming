load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "dplanerc",
    srcs = [
        "interface.go",
        "routes.go",
    ],
    importpath = "github.com/openconfig/lemming/dataplane/dplanerc",
    visibility = ["//visibility:public"],
    deps = [
        "//dataplane/proto/sai",
        "//dataplane/saiserver",
        "//gnmi",
        "//proto/dataplane",
        "//proto/forwarding",
        "//proto/routing",
        "@com_github_golang_glog//:glog",
        "@com_github_google_gopacket//:gopacket",
        "@com_github_google_gopacket//layers",
        "@com_github_openconfig_ygnmi//schemaless",
        "@com_github_openconfig_ygnmi//ygnmi",
        "@org_golang_google_protobuf//proto",
    ] + select({
        "@io_bazel_rules_go//go/platform:android": [
            "//dataplane/kernel",
            "//dataplane/protocol/lldp",
            "//gnmi/gnmiclient",
            "//gnmi/oc",
            "//gnmi/oc/ocpath",
            "@com_github_openconfig_ygot//ygot",
            "@com_github_vishvananda_netlink//:netlink",
            "@org_golang_google_grpc//:grpc",
            "@org_golang_x_sys//unix",
        ],
        "@io_bazel_rules_go//go/platform:linux": [
            "//dataplane/kernel",
            "//dataplane/protocol/lldp",
            "//gnmi/gnmiclient",
            "//gnmi/oc",
            "//gnmi/oc/ocpath",
            "@com_github_openconfig_ygot//ygot",
            "@com_github_vishvananda_netlink//:netlink",
            "@org_golang_google_grpc//:grpc",
            "@org_golang_x_sys//unix",
        ],
        "//conditions:default": [],
    }),
)
