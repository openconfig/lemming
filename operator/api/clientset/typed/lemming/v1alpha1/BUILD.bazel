load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "v1alpha1",
    srcs = [
        "doc.go",
        "generated_expansion.go",
        "lemming.go",
        "lemming_client.go",
    ],
    importpath = "github.com/openconfig/lemming/operator/api/clientset/typed/lemming/v1alpha1",
    visibility = ["//visibility:public"],
    deps = [
        "//operator/api/clientset/scheme",
        "//operator/api/lemming/v1alpha1",
        "@io_k8s_apimachinery//pkg/apis/meta/v1:meta",
        "@io_k8s_apimachinery//pkg/types",
        "@io_k8s_apimachinery//pkg/watch",
        "@io_k8s_client_go//rest",
    ],
)
