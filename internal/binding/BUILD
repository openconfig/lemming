load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "binding",
    srcs = [
        "binding.go",
        "local.go",
    ],
    importpath = "github.com/openconfig/lemming/internal/binding",
    visibility = ["//:__subpackages__"],
    deps = [
        "//:lemming",
        "//dataplane/dplaneopts",
        "//dataplane/forwarding/infra/fwdcontext",
        "//dataplane/forwarding/util/queue",
        "//dataplane/proto/sai",
        "//proto/forwarding",
        "@com_github_google_gopacket//:gopacket",
        "@com_github_google_uuid//:uuid",
        "@com_github_openconfig_kne//topo",
        "@com_github_openconfig_ondatra//:ondatra",
        "@com_github_openconfig_ondatra//binding",
        "@com_github_openconfig_ondatra//eventlis",
        "@com_github_openconfig_ondatra//knebind/init",
        "@com_github_openconfig_ondatra//proto",
        "@io_k8s_apimachinery//pkg/api/errors",
        "@io_k8s_apimachinery//pkg/apis/meta/v1:meta",
        "@io_k8s_client_go//kubernetes",
        "@io_k8s_client_go//tools/clientcmd",
        "@io_k8s_client_go//tools/leaderelection",
        "@io_k8s_client_go//tools/leaderelection/resourcelock",
        "@org_golang_google_grpc//credentials/local",
    ],
)
