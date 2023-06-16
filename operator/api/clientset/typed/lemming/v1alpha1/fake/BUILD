load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "fake",
    srcs = [
        "doc.go",
        "fake_lemming.go",
        "fake_lemming_client.go",
    ],
    importpath = "github.com/openconfig/lemming/operator/api/clientset/typed/lemming/v1alpha1/fake",
    visibility = ["//visibility:public"],
    deps = [
        "//operator/api/clientset/typed/lemming/v1alpha1",
        "//operator/api/lemming/v1alpha1",
        "@io_k8s_apimachinery//pkg/apis/meta/v1:meta",
        "@io_k8s_apimachinery//pkg/labels",
        "@io_k8s_apimachinery//pkg/runtime/schema",
        "@io_k8s_apimachinery//pkg/types",
        "@io_k8s_apimachinery//pkg/watch",
        "@io_k8s_client_go//rest",
        "@io_k8s_client_go//testing",
    ],
)
