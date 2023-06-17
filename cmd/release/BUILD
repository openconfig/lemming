load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "release_lib",
    srcs = ["release.go"],
    importpath = "github.com/openconfig/lemming/cmd/release",
    visibility = ["//visibility:private"],
    deps = [
        "@com_github_spf13_cobra//:cobra",
        "@com_google_cloud_go_cloudbuild//apiv1/v2:apiv1",
        "@com_google_cloud_go_cloudbuild//apiv1/v2/cloudbuildpb",
        "@org_golang_google_api//option",
    ],
)

go_binary(
    name = "release",
    embed = [":release_lib"],
    visibility = ["//visibility:public"],
)
