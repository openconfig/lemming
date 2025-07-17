load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "config",
    srcs = [
        "helper.go",
        "loader.go",
    ],
    importpath = "github.com/openconfig/lemming/internal/config",
    visibility = ["//:__subpackages__"],
    deps = [
        "//configs",
        "//proto/config",
        "@com_github_golang_glog//:glog",
        "@org_golang_google_protobuf//encoding/prototext",
    ],
)

go_test(
    name = "config_test",
    srcs = [
        "helper_test.go",
        "loader_test.go",
    ],
    embed = [":config"],
    deps = [
        "//proto/config",
        "@com_github_google_go_cmp//cmp",
        "@org_golang_google_protobuf//encoding/prototext",
        "@org_golang_google_protobuf//testing/protocmp",
    ],
)
