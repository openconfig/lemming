load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "binding",
    srcs = ["binding.go"],
    importpath = "github.com/openconfig/lemming/internal/binding",
    visibility = ["//:__subpackages__"],
    deps = [
        "@com_github_google_uuid//:uuid",
        "@com_github_openconfig_kne//topo",
        "@com_github_openconfig_ondatra//binding",
        "@com_github_openconfig_ondatra//knebind/init",
        "@com_github_openconfig_ondatra//proto",
        "@io_k8s_apimachinery//pkg/api/errors",
        "@io_k8s_apimachinery//pkg/apis/meta/v1:meta",
        "@io_k8s_client_go//kubernetes",
        "@io_k8s_client_go//tools/clientcmd",
        "@io_k8s_client_go//tools/leaderelection",
        "@io_k8s_client_go//tools/leaderelection/resourcelock",
    ],
)
