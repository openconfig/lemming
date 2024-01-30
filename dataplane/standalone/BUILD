load("@rules_pkg//pkg:deb.bzl", "pkg_deb")
load("@rules_pkg//pkg:tar.bzl", "pkg_tar")

cc_library(
    name = "entrypoint",
    srcs = [
        "entrypoint.cc",
    ],
    visibility = ["//visibility:public"],
    deps = [
        "//dataplane/proto/sai:sai_cc_grpc",
        "//dataplane/standalone/packetio",
        "//dataplane/standalone/sai",
        "@com_github_google_glog//:glog",
        "@com_github_grpc_grpc//:grpc++",
    ],
    alwayslink = True,
)

cc_binary(
    name = "sai",
    linkshared = True,
    linkstatic = True,
    deps = [
        ":entrypoint",
    ],
)

cc_binary(
    name = "sai-dynamic",
    linkopts = [
        "-lprotobuf",
        "-lgrpc++",
        "-lglog",
    ],
    linkshared = True,
    linkstatic = True,
    deps = [
        ":entrypoint",
    ],
)

pkg_tar(
    name = "sai-tar",
    package_dir = "/usr/lib/x86_64-linux-gnu",
    files = {
       ":sai-dynamic": "libsaivs.so.0.0.0"
    },
    symlinks = {
        "/usr/lib/x86_64-linux-gnu/libsaivs.so.0": "libsaivs.so.0.0.0"
    },
)

pkg_deb(
    name = "sai-deb",
    architecture = "amd64",
    maintainer = "OpenConfig",
    description = "SAI implementation for lucius",
    data = ":sai-tar",
    package = "lucius-libsai",
    version = "1.0-3",
)

# By building a non-shared library, detect some linker issues.
cc_binary(
    name = "saibin",
    linkstatic = True,
    deps = [
        ":entrypoint",
    ],
)
