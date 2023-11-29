load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "dataplane",
    srcs = [
        "reconcilers_linux.go",
        "reconcilers_nonlinux.go",
        "server.go",
    ],
    importpath = "github.com/openconfig/lemming/dataplane",
    visibility = ["//visibility:public"],
    deps = [
        "//dataplane/standalone/proto:sai",
        "//dataplane/standalone/saiserver",
        "//dataplane/standalone/saiserver/attrmgr",
        "//gnmi/oc",
        "//gnmi/reconciler",
        "@com_github_openconfig_gnmi//proto/gnmi",
        "@com_github_openconfig_ygnmi//ygnmi",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//credentials/local",
    ] + select({
        "@io_bazel_rules_go//go/platform:aix": [
            "//dataplane/reconcilers:handlers",
        ],
        "@io_bazel_rules_go//go/platform:android": [
            "//dataplane/reconcilers:handlers",
        ],
        "@io_bazel_rules_go//go/platform:darwin": [
            "//dataplane/reconcilers:handlers",
        ],
        "@io_bazel_rules_go//go/platform:dragonfly": [
            "//dataplane/reconcilers:handlers",
        ],
        "@io_bazel_rules_go//go/platform:freebsd": [
            "//dataplane/reconcilers:handlers",
        ],
        "@io_bazel_rules_go//go/platform:illumos": [
            "//dataplane/reconcilers:handlers",
        ],
        "@io_bazel_rules_go//go/platform:ios": [
            "//dataplane/reconcilers:handlers",
        ],
        "@io_bazel_rules_go//go/platform:js": [
            "//dataplane/reconcilers:handlers",
        ],
        "@io_bazel_rules_go//go/platform:linux": [
            "//dataplane/reconcilers:handlers",
        ],
        "@io_bazel_rules_go//go/platform:netbsd": [
            "//dataplane/reconcilers:handlers",
        ],
        "@io_bazel_rules_go//go/platform:openbsd": [
            "//dataplane/reconcilers:handlers",
        ],
        "@io_bazel_rules_go//go/platform:plan9": [
            "//dataplane/reconcilers:handlers",
        ],
        "@io_bazel_rules_go//go/platform:solaris": [
            "//dataplane/reconcilers:handlers",
        ],
        "@io_bazel_rules_go//go/platform:windows": [
            "//dataplane/reconcilers:handlers",
        ],
        "//conditions:default": [],
    }),
)
