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
        "//dataplane/dplaneopts",
        "//dataplane/kernel/tap",
        "//dataplane/proto/packetio",
        "//dataplane/proto/sai",
        "//dataplane/protocol",
        "//dataplane/saiserver",
        "//dataplane/saiserver/attrmgr",
        "//dataplane/standalone/pkthandler/pktiohandler",
        "//gnmi/oc",
        "//gnmi/reconciler",
        "@com_github_openconfig_gnmi//proto/gnmi",
        "@com_github_openconfig_ygnmi//ygnmi",
        "@org_golang_google_grpc//:grpc",
        "@org_golang_google_grpc//credentials/local",
        "@org_golang_google_grpc//reflection",
        "@org_golang_google_protobuf//proto",
    ] + select({
        "@io_bazel_rules_go//go/platform:aix": [
            "//dataplane/dplanerc",
        ],
        "@io_bazel_rules_go//go/platform:android": [
            "//dataplane/dplanerc",
        ],
        "@io_bazel_rules_go//go/platform:darwin": [
            "//dataplane/dplanerc",
        ],
        "@io_bazel_rules_go//go/platform:dragonfly": [
            "//dataplane/dplanerc",
        ],
        "@io_bazel_rules_go//go/platform:freebsd": [
            "//dataplane/dplanerc",
        ],
        "@io_bazel_rules_go//go/platform:illumos": [
            "//dataplane/dplanerc",
        ],
        "@io_bazel_rules_go//go/platform:ios": [
            "//dataplane/dplanerc",
        ],
        "@io_bazel_rules_go//go/platform:js": [
            "//dataplane/dplanerc",
        ],
        "@io_bazel_rules_go//go/platform:linux": [
            "//dataplane/dplanerc",
        ],
        "@io_bazel_rules_go//go/platform:netbsd": [
            "//dataplane/dplanerc",
        ],
        "@io_bazel_rules_go//go/platform:openbsd": [
            "//dataplane/dplanerc",
        ],
        "@io_bazel_rules_go//go/platform:osx": [
            "//dataplane/dplanerc",
        ],
        "@io_bazel_rules_go//go/platform:plan9": [
            "//dataplane/dplanerc",
        ],
        "@io_bazel_rules_go//go/platform:qnx": [
            "//dataplane/dplanerc",
        ],
        "@io_bazel_rules_go//go/platform:solaris": [
            "//dataplane/dplanerc",
        ],
        "@io_bazel_rules_go//go/platform:windows": [
            "//dataplane/dplanerc",
        ],
        "//conditions:default": [],
    }),
)
