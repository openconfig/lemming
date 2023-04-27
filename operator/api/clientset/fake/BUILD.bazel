load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "fake",
    srcs = [
        "clientset_generated.go",
        "doc.go",
        "register.go",
    ],
    importpath = "github.com/openconfig/lemming/operator/api/clientset/fake",
    visibility = ["//visibility:public"],
    deps = [
        "//operator/api/clientset",
        "//operator/api/clientset/typed/lemming/v1alpha1",
        "//operator/api/clientset/typed/lemming/v1alpha1/fake",
        "//operator/api/lemming/v1alpha1",
        "@io_k8s_apimachinery//pkg/apis/meta/v1:meta",
        "@io_k8s_apimachinery//pkg/runtime",
        "@io_k8s_apimachinery//pkg/runtime/schema",
        "@io_k8s_apimachinery//pkg/runtime/serializer",
        "@io_k8s_apimachinery//pkg/util/runtime",
        "@io_k8s_apimachinery//pkg/watch",
        "@io_k8s_client_go//discovery",
        "@io_k8s_client_go//discovery/fake",
        "@io_k8s_client_go//testing",
    ],
)
