load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")
load("@rules_oci//oci:defs.bzl", "oci_image", "oci_tarball")
load("@rules_pkg//:pkg.bzl", "pkg_tar")

go_library(
    name = "lucius_lib",
    srcs = ["lucius.go"],
    importpath = "github.com/openconfig/lemming/dataplane/standalone/lucius",
    visibility = ["//visibility:private"],
    deps = [
        "//dataplane/internal/engine",
        "//dataplane/standalone/saiserver",
        "//dataplane/standalone/saiserver/attrmgr",
        "//proto/dataplane",
        "//proto/forwarding",
        "@com_github_golang_glog//:glog",
        "@com_github_grpc_ecosystem_go_grpc_middleware_v2//interceptors/logging",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//credentials/insecure",
    ],
)

go_binary(
    name = "lucius",
    embed = [":lucius_lib"],
    visibility = ["//visibility:public"],
)

pkg_tar(
    name = "bin-tar",
    srcs = [":lucius"],
    package_dir = "lucius",
)

oci_image(
    name = "image",
    base = "@debian_bookworm",
    entrypoint = ["/lucius/lucius"],
    tars = [":bin-tar"],
)

# Note: This is image isn't suitable for release, as it doesn't contain license information.
oci_tarball(
    name = "image-tar",
    image = ":image",
    repo_tags = ["us-west1-docker.pkg.dev/openconfig-lemming/release/lucius:ga"],
)
