load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")
load("@rules_distroless//distroless:defs.bzl", "cacerts")
load("@rules_oci//oci:defs.bzl", "oci_image", "oci_tarball")
load("@rules_pkg//:pkg.bzl", "pkg_tar")

go_library(
    name = "lemming_lib",
    srcs = ["lemming.go"],
    importpath = "github.com/openconfig/lemming/cmd/lemming",
    visibility = ["//visibility:private"],
    deps = [
        "//:lemming",
        "//dataplane/dplaneopts",
        "//internal/telemetry",
        "@com_github_golang_glog//:glog",
        "@com_github_spf13_pflag//:pflag",
        "@com_github_spf13_viper//:viper",
        "@org_golang_google_grpc//credentials",
        "@org_golang_google_grpc//credentials/insecure",
    ],
)

go_binary(
    name = "lemming",
    embed = [":lemming_lib"],
    visibility = ["//visibility:public"],
)

pkg_tar(
    name = "bin-tar",
    srcs = [":lemming"],
    package_dir = "lemming",
)

cacerts(
    name = "cacerts",
    package = select({
        "@platforms//cpu:x86_64": "@bookworm//ca-certificates/amd64:data",
    }),
)

oci_image(
    name = "image",
    base = "@debian_bookworm",
    entrypoint = ["/lemming/lemming"],
    tars = [
        ":bin-tar",
        ":cacerts",
        "@bookworm//libpcap-dev/amd64",
    ],
)

# Note: This is image isn't suitable for release, as it doesn't contain license information.
oci_tarball(
    name = "image-tar",
    image = ":image",
    repo_tags = ["us-west1-docker.pkg.dev/openconfig-lemming/release/lemming:ga"],
)
