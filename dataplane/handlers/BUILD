load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "handlers",
    srcs = [
        "interface.go",
        "routes.go",
    ],
    importpath = "github.com/openconfig/lemming/dataplane/handlers",
    visibility = ["//visibility:public"],
    deps = [
        "//dataplane/internal/engine",
        "//gnmi",
        "//gnmi/reconciler",
        "//proto/dataplane",
        "@com_github_golang_glog//:glog",
        "@com_github_openconfig_ygnmi//schemaless",
        "@com_github_openconfig_ygnmi//ygnmi",
    ] + select({
        "@io_bazel_rules_go//go/platform:android": [
            "//dataplane/internal/kernel",
            "//gnmi/gnmiclient",
            "//gnmi/oc",
            "//gnmi/oc/ocpath",
            "//proto/forwarding",
            "@com_github_openconfig_ygot//ygot",
            "@com_github_vishvananda_netlink//:netlink",
            "@org_golang_x_sys//unix",
        ],
        "@io_bazel_rules_go//go/platform:linux": [
            "//dataplane/internal/kernel",
            "//gnmi/gnmiclient",
            "//gnmi/oc",
            "//gnmi/oc/ocpath",
            "//proto/forwarding",
            "@com_github_openconfig_ygot//ygot",
            "@com_github_vishvananda_netlink//:netlink",
            "@org_golang_x_sys//unix",
        ],
        "//conditions:default": [],
    }),
)
