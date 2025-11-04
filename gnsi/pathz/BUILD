load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "pathz",
    srcs = ["pathz.go"],
    importpath = "github.com/openconfig/lemming/gnsi/pathz",
    visibility = ["//visibility:public"],
    deps = [
        "//gnsi/acltrie",
        "@com_github_openconfig_gnmi//proto/gnmi",
        "@com_github_openconfig_gnsi//pathz",
        "@org_golang_google_grpc//codes",
        "@org_golang_google_grpc//status",
    ],
)

go_test(
    name = "pathz_test",
    size = "small",
    srcs = ["pathz_test.go"],
    embed = [":pathz"],
    deps = [
        "@com_github_google_go_cmp//cmp",
        "@com_github_openconfig_gnmi//errdiff",
        "@com_github_openconfig_gnmi//proto/gnmi",
        "@com_github_openconfig_gnsi//pathz",
        "@org_golang_google_grpc//:grpc",
        "@org_golang_google_grpc//credentials/insecure",
        "@org_golang_google_protobuf//testing/protocmp",
    ],
)
