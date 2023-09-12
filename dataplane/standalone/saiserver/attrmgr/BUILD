load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "attrmgr",
    srcs = ["attrmgr.go"],
    importpath = "github.com/openconfig/lemming/dataplane/standalone/saiserver/attrmgr",
    visibility = ["//visibility:public"],
    deps = [
        "//dataplane/standalone/proto:sai",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//codes",
        "@org_golang_google_grpc//status",
        "@org_golang_google_protobuf//proto",
        "@org_golang_google_protobuf//reflect/protoreflect",
        "@org_golang_google_protobuf//reflect/protoregistry",
        "@org_golang_google_protobuf//types/descriptorpb",
    ],
)

go_test(
    name = "attrmgr_test",
    srcs = ["attrmgr_test.go"],
    embed = [":attrmgr"],
    deps = [
        "//dataplane/standalone/proto:sai",
        "@com_github_google_go_cmp//cmp",
        "@com_github_openconfig_gnmi//errdiff",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//codes",
        "@org_golang_google_grpc//status",
        "@org_golang_google_protobuf//proto",
        "@org_golang_google_protobuf//reflect/protoreflect",
        "@org_golang_google_protobuf//testing/protocmp",
    ],
)
