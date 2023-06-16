load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

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
