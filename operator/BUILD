load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")
load("@rules_oci//oci:defs.bzl", "oci_image", "oci_tarball")
load("@rules_pkg//:pkg.bzl", "pkg_tar")

go_library(
    name = "operator_lib",
    srcs = ["main.go"],
    importpath = "github.com/openconfig/lemming/operator",
    visibility = ["//visibility:private"],
    deps = [
        "//operator/api/lemming/v1alpha1",
        "//operator/controllers",
        "@io_k8s_apimachinery//pkg/runtime",
        "@io_k8s_apimachinery//pkg/util/runtime",
        "@io_k8s_client_go//kubernetes/scheme",
        "@io_k8s_client_go//plugin/pkg/client/auth",
        "@io_k8s_sigs_controller_runtime//:controller-runtime",
        "@io_k8s_sigs_controller_runtime//pkg/healthz",
        "@io_k8s_sigs_controller_runtime//pkg/log/zap",
    ],
)

go_binary(
    name = "operator",
    embed = [":operator_lib"],
    visibility = ["//visibility:public"],
)

pkg_tar(
    name = "bin-tar",
    srcs = [":operator"],
)

oci_image(
    name = "image",
    base = "@distroless_static_debug_nonroot",
    entrypoint = ["/operator"],
    tars = [":bin-tar"],
)

# Note: This is image isn't suitable for release, as it doesn't contain license information.
oci_tarball(
    name = "image-tar",
    image = ":image",
    repo_tags = ["us-west1-docker.pkg.dev/openconfig-lemming/release/operator:ga"],
)
