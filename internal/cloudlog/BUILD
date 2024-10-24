load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "cloudlog",
    srcs = ["cloudlog.go"],
    importpath = "github.com/openconfig/lemming/internal/cloudlog",
    visibility = ["//:__subpackages__"],
    deps = [
        "@com_google_cloud_go_logging//:logging",
        "@com_google_cloud_go_logging//apiv2/loggingpb",
        "@io_opentelemetry_go_otel_trace//:trace",
    ],
)
