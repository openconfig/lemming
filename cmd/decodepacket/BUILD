load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "decodepacket_lib",
    srcs = ["decodepacket.go"],
    importpath = "github.com/openconfig/lemming/cmd/decodepacket",
    visibility = ["//visibility:private"],
    deps = [
        "@com_github_google_gopacket//:gopacket",
        "@com_github_google_gopacket//layers",
    ],
)

go_binary(
    name = "decodepacket",
    embed = [":decodepacket_lib"],
    visibility = ["//visibility:public"],
)
