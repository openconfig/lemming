load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "packet",
    srcs = ["packet.go"],
    importpath = "github.com/openconfig/lemming/cmd/packet",
    visibility = ["//visibility:public"],
    deps = [
        "//internal/packetutil",
        "@com_github_google_gopacket//:gopacket",
        "@com_github_google_gopacket//layers",
        "@com_github_spf13_cobra//:cobra",
    ],
)
