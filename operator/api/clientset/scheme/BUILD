load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "scheme",
    srcs = [
        "doc.go",
        "register.go",
    ],
    importpath = "github.com/openconfig/lemming/operator/api/clientset/scheme",
    visibility = ["//visibility:public"],
    deps = [
        "//operator/api/lemming/v1alpha1",
        "@io_k8s_apimachinery//pkg/apis/meta/v1:meta",
        "@io_k8s_apimachinery//pkg/runtime",
        "@io_k8s_apimachinery//pkg/runtime/schema",
        "@io_k8s_apimachinery//pkg/runtime/serializer",
        "@io_k8s_apimachinery//pkg/util/runtime",
    ],
)
