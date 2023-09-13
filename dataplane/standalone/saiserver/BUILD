load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "saiserver",
    srcs = [
        "ports.go",
        "saiserver.go",
        "switch.go",
    ],
    importpath = "github.com/openconfig/lemming/dataplane/standalone/saiserver",
    visibility = ["//visibility:public"],
    deps = [
        "//dataplane/standalone/proto:sai",
        "//dataplane/standalone/saiserver/attrmgr",
        "//proto/dataplane",
        "//proto/forwarding",
        "@com_github_golang_glog//:glog",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_protobuf//proto",
    ],
)

go_test(
    name = "saiserver_test",
    srcs = ["switch_test.go"],
    embed = [":saiserver"],
    deps = [
        "//dataplane/standalone/proto:sai",
        "//dataplane/standalone/saiserver/attrmgr",
        "@com_github_google_go_cmp//cmp",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//credentials/insecure",
        "@org_golang_google_protobuf//proto",
        "@org_golang_google_protobuf//testing/protocmp",
    ],
)
