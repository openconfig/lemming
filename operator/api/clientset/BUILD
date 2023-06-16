load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "clientset",
    srcs = [
        "clientset.go",
        "doc.go",
    ],
    importpath = "github.com/openconfig/lemming/operator/api/clientset",
    visibility = ["//visibility:public"],
    deps = [
        "//operator/api/clientset/typed/lemming/v1alpha1",
        "@io_k8s_client_go//discovery",
        "@io_k8s_client_go//rest",
        "@io_k8s_client_go//util/flowcontrol",
    ],
)
