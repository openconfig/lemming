load("@io_bazel_rules_go//go:def.bzl", "go_test")

go_test(
    name = "fault_test",
    srcs = ["fault_test.go"],
    deps = [
        "//internal/binding",
        "//proto/fault",
        "@com_github_openconfig_ondatra//:ondatra",
        "@org_golang_google_grpc//:grpc",
    ],
)
