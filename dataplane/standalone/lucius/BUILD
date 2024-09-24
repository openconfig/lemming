load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")
load("@rules_distroless//distroless:defs.bzl", "cacerts")
load("@rules_oci//oci:defs.bzl", "oci_image", "oci_tarball")
load("@rules_pkg//:pkg.bzl", "pkg_tar")

go_library(
    name = "lucius_lib",
    srcs = ["lucius.go"],
    importpath = "github.com/openconfig/lemming/dataplane/standalone/lucius",
    visibility = ["//visibility:private"],
    deps = [
        "//dataplane/dplaneopts",
        "//dataplane/saiserver",
        "//dataplane/saiserver/attrmgr",
        "//proto/forwarding",
        "@com_github_googlecloudplatform_opentelemetry_operations_go_exporter_trace//:trace",
        "@com_github_grpc_ecosystem_go_grpc_middleware_v2//interceptors/logging",
        "@io_opentelemetry_go_contrib_instrumentation_google_golang_org_grpc_otelgrpc//:otelgrpc",
        "@io_opentelemetry_go_otel//:otel",
        "@io_opentelemetry_go_otel//propagation",
        "@io_opentelemetry_go_otel_exporters_stdout_stdoutlog//:stdoutlog",
        "@io_opentelemetry_go_otel_log//global",
        "@io_opentelemetry_go_otel_sdk//resource",
        "@io_opentelemetry_go_otel_sdk//trace",
        "@io_opentelemetry_go_otel_sdk_log//:log",
        "@io_opentelemetry_go_otel_sdk_metric//:metric",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//credentials/insecure",
        "@org_golang_google_grpc//metadata",
        "@org_golang_google_grpc//reflection",
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

cacerts(
    name = "cacerts",
    package = select({
        "@platforms//cpu:x86_64": "@bookworm//ca-certificates/amd64:data",
    }),
)

oci_image(
    name = "image",
    base = "@debian_bookworm",
    entrypoint = ["/lucius/lucius"],
    tars = [
        ":bin-tar",
        ":cacerts",
    ],
)

# Note: This is image isn't suitable for release, as it doesn't contain license information.
oci_tarball(
    name = "image-tar",
    image = ":image",
    repo_tags = ["us-west1-docker.pkg.dev/openconfig-lemming/release/lucius:ga"],
)
