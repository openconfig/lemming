load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "packetutil",
    srcs = ["capture.go"],
    importpath = "github.com/openconfig/lemming/internal/packetutil",
    visibility = ["//:__subpackages__"],
    deps = [
        "@com_github_google_gopacket//:gopacket",
        "@com_github_google_gopacket//layers",
        "@com_github_google_gopacket//pcapgo",
    ],
)
