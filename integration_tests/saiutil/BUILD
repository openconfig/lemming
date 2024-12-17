load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "saiutil",
    srcs = [
        "runner.go",
        "saiutil.go",
    ],
    importpath = "github.com/openconfig/lemming/integration_tests/saiutil",
    visibility = ["//visibility:public"],
    deps = [
        "//dataplane/proto/sai",
        "//gnmi/oc",
        "//internal/binding",
        "@com_github_google_go_cmp//cmp",
        "@com_github_google_go_cmp//cmp/cmpopts",
        "@com_github_google_gopacket//:gopacket",
        "@com_github_google_gopacket//layers",
        "@com_github_openconfig_ondatra//:ondatra",
        "@com_github_openconfig_ondatra//binding",
        "@org_golang_google_grpc//:grpc",
        "@org_golang_google_protobuf//proto",
    ],
)
