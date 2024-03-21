load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "controllers",
    srcs = ["lemming_controller.go"],
    importpath = "github.com/openconfig/lemming/operator/controllers",
    visibility = ["//visibility:public"],
    deps = [
        "//operator/api/lemming/v1alpha1",
        "@com_github_google_go_cmp//cmp",
        "@io_k8s_api//core/v1:core",
        "@io_k8s_apimachinery//pkg/api/equality",
        "@io_k8s_apimachinery//pkg/api/errors",
        "@io_k8s_apimachinery//pkg/apis/meta/v1:meta",
        "@io_k8s_apimachinery//pkg/runtime",
        "@io_k8s_apimachinery//pkg/types",
        "@io_k8s_apimachinery//pkg/util/intstr",
        "@io_k8s_client_go//util/cert",
        "@io_k8s_sigs_controller_runtime//:controller-runtime",
        "@io_k8s_sigs_controller_runtime//pkg/client",
        "@io_k8s_sigs_controller_runtime//pkg/log",
        "@io_k8s_utils//pointer",
    ],
)

go_test(
    name = "controllers_test",
    size = "small",
    srcs = ["lemming_controller_test.go"],
    embed = [":controllers"],
    deps = [
        "//operator/api/clientset/scheme",
        "//operator/api/lemming/v1alpha1",
        "@com_github_google_go_cmp//cmp",
        "@io_k8s_api//core/v1:core",
        "@io_k8s_apimachinery//pkg/api/resource",
        "@io_k8s_apimachinery//pkg/apis/meta/v1:meta",
        "@io_k8s_apimachinery//pkg/types",
        "@io_k8s_apimachinery//pkg/util/intstr",
        "@io_k8s_client_go//kubernetes/scheme",
        "@io_k8s_sigs_controller_runtime//pkg/client",
        "@io_k8s_sigs_controller_runtime//pkg/client/fake",
        "@io_k8s_sigs_controller_runtime//pkg/log",
        "@io_k8s_sigs_controller_runtime//pkg/log/zap",
        "@io_k8s_sigs_controller_runtime//pkg/reconcile",
        "@io_k8s_utils//pointer",
    ],
)
