load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "fault",
    srcs = [
        "client.go",
        "fault.go",
    ],
    importpath = "github.com/openconfig/lemming/fault",
    visibility = ["//visibility:public"],
    deps = [
        "//proto/config",
        "//proto/fault",
        "@com_github_golang_glog//:glog",
        "@com_github_google_uuid//:uuid",
        "@openconfig_gnmi//proto/gnmi",
        "@org_golang_google_grpc//:grpc",
        "@org_golang_google_grpc//codes",
        "@org_golang_google_grpc//status",
        "@org_golang_google_protobuf//proto",
        "@org_golang_google_protobuf//types/known/anypb",
    ],
)

go_test(
    name = "fault_test",
    srcs = [
        "fault_test.go",
        "interceptor_test.go",
    ],
    embed = [":fault"],
    deps = [
        "//fault/proto/test",
        "//proto/config",
        "//proto/fault",
        "@com_github_google_go_cmp//cmp",
        "@com_github_openconfig_gnoi//system",
        "@openconfig_gnmi//errdiff",
        "@org_golang_google_genproto_googleapis_rpc//status",
        "@org_golang_google_grpc//:grpc",
        "@org_golang_google_grpc//codes",
        "@org_golang_google_grpc//credentials/insecure",
        "@org_golang_google_grpc//status",
        "@org_golang_google_protobuf//proto",
        "@org_golang_google_protobuf//testing/protocmp",
        "@org_golang_google_protobuf//types/known/anypb",
    ],
)
