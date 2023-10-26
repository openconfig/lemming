load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")
load("@bazel_gazelle//:def.bzl", "gazelle")

# gazelle:prefix github.com/openconfig/lemming
# gazelle:resolve proto proto google/rpc/status.proto @googleapis//google/rpc:status_proto
# gazelle:resolve proto go google/rpc/status.proto  @org_golang_google_genproto_googleapis_rpc//status
# gazelle:resolve go github.com/p4lang/p4runtime/go/p4/v1 @com_github_p4lang_p4runtime//:p4runtime_go_proto
# gazelle:build_file_name BUILD
gazelle(
    name = "gazelle",
)

go_library(
    name = "lemming",
    srcs = ["lemming.go"],
    importpath = "github.com/openconfig/lemming",
    visibility = ["//visibility:public"],
    deps = [
        "//bgp",
        "//dataplane",
        "//gnmi",
        "//gnmi/fakedevice",
        "//gnmi/oc",
        "//gnmi/reconciler",
        "//gnoi",
        "//gnsi",
        "//gribi",
        "//p4rt",
        "//sysrib",
        "@com_github_golang_glog//:glog",
        "@com_github_spf13_viper//:viper",
        "@io_k8s_klog_v2//:klog",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//credentials",
        "@org_golang_google_grpc//reflection",
    ],
)

go_test(
    name = "lemming_test",
    size = "medium",
    srcs = ["lemming_test.go"],
    embed = [":lemming"],
    deps = [
        "@com_github_openconfig_gnmi//errdiff",
        "@com_github_openconfig_gnmi//proto/gnmi",
        "@com_github_openconfig_gnoi//bgp",
        "@com_github_openconfig_gnoi//cert",
        "@com_github_openconfig_gnoi//diag",
        "@com_github_openconfig_gnoi//factory_reset",
        "@com_github_openconfig_gnoi//file",
        "@com_github_openconfig_gnoi//healthz",
        "@com_github_openconfig_gnoi//layer2",
        "@com_github_openconfig_gnoi//mpls",
        "@com_github_openconfig_gnoi//os",
        "@com_github_openconfig_gnoi//otdr",
        "@com_github_openconfig_gnoi//system",
        "@com_github_openconfig_gnoi//wavelength_router",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//credentials/insecure",
        "@org_golang_google_protobuf//proto",
    ],
)
