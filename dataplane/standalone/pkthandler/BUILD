load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")
load("@rules_pkg//pkg:deb.bzl", "pkg_deb")
load("@rules_pkg//pkg:tar.bzl", "pkg_tar")
load("@rules_pkg//pkg:mappings.bzl", "pkg_attributes", "pkg_files")

go_library(
    name = "pkthandler_lib",
    srcs = ["main.go"],
    importpath = "github.com/openconfig/lemming/dataplane/standalone/pkthandler",
    visibility = ["//visibility:private"],
    deps = [
        "//dataplane/standalone/pkthandler/pktiohandler",
        "//proto/forwarding",
        "@com_github_golang_glog//:glog",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//credentials/insecure",
    ],
)

go_binary(
    name = "pkthandler",
    embed = [":pkthandler_lib"],
    visibility = ["//visibility:public"],
)

pkg_files(
    name = "pkthandler-bin",
    srcs = [
        ":pkthandler",
    ],
    attributes = pkg_attributes(
        mode = "0755",
    ),
    prefix = "/usr/bin",
)

pkg_files(
    name = "pkthandler-service",
    srcs = [
        "lucius-pkthandler.service",
    ],
    prefix = "/etc/systemd/system/",
)

pkg_tar(
    name = "pkthandler-tar",
    srcs = [
        ":pkthandler-bin",
        ":pkthandler-service",
    ],
    tags = ["manual"],
)

pkg_deb(
    name = "pkthandler-deb",
    architecture = "amd64",
    data = ":pkthandler-tar",
    description = "Packet Handler for lucius",
    maintainer = "OpenConfig",
    package = "lucius-pkthandler",
    postinst = "postinst",
    postrm = "postrm",
    preinst = "preinst",
    prerm = "prerm",
    tags = ["manual"],
    version = "1.0-0",
)
