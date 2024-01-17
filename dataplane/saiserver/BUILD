load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "saiserver",
    srcs = [
        "acl.go",
        "hostif.go",
        "ports.go",
        "routing.go",
        "saiserver.go",
        "switch.go",
    ],
    importpath = "github.com/openconfig/lemming/dataplane/saiserver",
    visibility = ["//visibility:public"],
    deps = [
        "//dataplane/cpusink",
        "//dataplane/dplaneopts",
        "//dataplane/forwarding",
        "//dataplane/forwarding/attributes",
        "//dataplane/forwarding/fwdconfig",
        "//dataplane/forwarding/infra/fwdcontext",
        "//dataplane/proto/sai",
        "//dataplane/saiserver/attrmgr",
        "//proto/forwarding",
        "@com_github_golang_glog//:glog",
        "@com_github_openconfig_gnmi//errlist",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//codes",
        "@org_golang_google_grpc//status",
        "@org_golang_google_protobuf//proto",
    ],
)

go_test(
    name = "saiserver_test",
    srcs = [
        "acl_test.go",
        "hostif_test.go",
        "ports_test.go",
        "routing_test.go",
        "switch_test.go",
    ],
    embed = [":saiserver"],
    deps = [
        "//dataplane/dplaneopts",
        "//dataplane/forwarding/infra/fwdcontext",
        "//dataplane/forwarding/infra/fwdobject",
        "//dataplane/proto/sai",
        "//dataplane/saiserver/attrmgr",
        "//proto/forwarding",
        "@com_github_google_go_cmp//cmp",
        "@com_github_openconfig_gnmi//errdiff",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//credentials/insecure",
        "@org_golang_google_protobuf//proto",
        "@org_golang_google_protobuf//testing/protocmp",
    ],
)
