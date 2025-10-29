load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "gnmi",
    srcs = [
        "cache.go",
        "collector.go",
        "generate.go",
        "gnmi.go",
    ],
    importpath = "github.com/openconfig/lemming/gnmi",
    visibility = ["//visibility:public"],
    deps = [
        "//gnmi/oc",
        "//gnmi/reconciler",
        "@com_github_golang_glog//:glog",
        "@com_github_openconfig_ygnmi//app/ygnmi/cmd",
        "@com_github_openconfig_ygot//util",
        "@com_github_openconfig_ygot//ygot",
        "@com_github_openconfig_ygot//ytypes",
        "@openconfig_gnmi//cache",
        "@openconfig_gnmi//proto/gnmi",
        "@openconfig_gnmi//subscribe",
        "@org_golang_google_grpc//:grpc",
        "@org_golang_google_grpc//codes",
        "@org_golang_google_grpc//metadata",
        "@org_golang_google_grpc//peer",
        "@org_golang_google_grpc//status",
        "@org_golang_google_protobuf//encoding/prototext",
        "@org_golang_google_protobuf//proto",
    ],
)

go_test(
    name = "gnmi_test",
    size = "small",
    srcs = [
        "gnmi_bench_test.go",
        "gnmi_test.go",
    ],
    embed = [":gnmi"],
    deps = [
        "//gnmi/gnmiclient",
        "//gnmi/oc",
        "//gnmi/oc/ocpath",
        "@com_github_google_go_cmp//cmp",
        "@com_github_google_go_cmp//cmp/cmpopts",
        "@com_github_openconfig_ygnmi//schemaless",
        "@com_github_openconfig_ygnmi//ygnmi",
        "@com_github_openconfig_ygot//ygot",
        "@openconfig_gnmi//errdiff",
        "@openconfig_gnmi//proto/gnmi",
        "@openconfig_gnmi//value",
        "@org_golang_google_grpc//:grpc",
        "@org_golang_google_grpc//credentials/local",
        "@org_golang_google_grpc//metadata",
        "@org_golang_google_protobuf//encoding/prototext",
    ],
)
