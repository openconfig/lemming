load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "sai",
    srcs = [
        "apidiscover.go",
        "sai.go",
    ],
    importpath = "github.com/openconfig/lemming/dataplane/luciusctl/sai",
    visibility = ["//visibility:public"],
    deps = [
        "@com_github_spf13_cobra//:cobra",
        "@com_github_spf13_viper//:viper",
        "@org_golang_google_grpc//:grpc",
        "@org_golang_google_grpc//credentials",
        "@org_golang_google_grpc//credentials/insecure",
        "@org_golang_google_grpc//reflection/grpc_reflection_v1",
        "@org_golang_google_protobuf//encoding/prototext",
        "@org_golang_google_protobuf//proto",
        "@org_golang_google_protobuf//reflect/protodesc",
        "@org_golang_google_protobuf//reflect/protoreflect",
        "@org_golang_google_protobuf//types/descriptorpb",
        "@org_golang_google_protobuf//types/dynamicpb",
    ],
)

go_test(
    name = "sai_test",
    srcs = ["apidiscover_test.go"],
    embed = [":sai"],
    deps = [
        "@com_github_google_go_cmp//cmp",
        "@org_golang_google_grpc//:grpc",
        "@org_golang_google_grpc//reflection/grpc_reflection_v1",
        "@org_golang_google_protobuf//proto",
        "@org_golang_google_protobuf//types/descriptorpb",
    ],
)
