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
        "//dataplane/internal/engine",
        "//gnmi/oc",
        "//gnmi/reconciler",
        "//proto/dataplane",
        "//proto/forwarding",
        "@com_github_openconfig_gnmi//proto/gnmi",
        "@com_github_openconfig_ygnmi//ygnmi",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//credentials/local",
    ] + select({
        "@io_bazel_rules_go//go/platform:aix": [
            "//dataplane/handlers",
        ],
        "@io_bazel_rules_go//go/platform:android": [
            "//dataplane/handlers",
        ],
        "@io_bazel_rules_go//go/platform:darwin": [
            "//dataplane/handlers",
        ],
        "@io_bazel_rules_go//go/platform:dragonfly": [
            "//dataplane/handlers",
        ],
        "@io_bazel_rules_go//go/platform:freebsd": [
            "//dataplane/handlers",
        ],
        "@io_bazel_rules_go//go/platform:illumos": [
            "//dataplane/handlers",
        ],
        "@io_bazel_rules_go//go/platform:ios": [
            "//dataplane/handlers",
        ],
        "@io_bazel_rules_go//go/platform:js": [
            "//dataplane/handlers",
        ],
        "@io_bazel_rules_go//go/platform:linux": [
            "//dataplane/handlers",
        ],
        "@io_bazel_rules_go//go/platform:netbsd": [
            "//dataplane/handlers",
        ],
        "@io_bazel_rules_go//go/platform:openbsd": [
            "//dataplane/handlers",
        ],
        "@io_bazel_rules_go//go/platform:plan9": [
            "//dataplane/handlers",
        ],
        "@io_bazel_rules_go//go/platform:solaris": [
            "//dataplane/handlers",
        ],
        "@io_bazel_rules_go//go/platform:windows": [
            "//dataplane/handlers",
        ],
        "//conditions:default": [],
    }),
)
