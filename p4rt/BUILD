load("@io_bazel_rules_go//go:def.bzl", "go_library")

# gazelle:resolve proto proto google/rpc/status.proto @googleapis//google/rpc:status_proto
# gazelle:resolve proto go google/rpc/status.proto  @org_golang_google_genproto_googleapis_rpc//status

go_library(
    name = "p4rt",
    srcs = ["p4rt.go"],
    importpath = "github.com/openconfig/lemming/p4rt",
    visibility = ["//visibility:public"],
    deps = [
        "@com_github_p4lang_p4runtime//:p4runtime_go_proto",
        "@org_golang_google_grpc//:go_default_library",
    ],
)
