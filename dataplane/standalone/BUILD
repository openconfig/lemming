cc_library(
    name = "entrypoint",
    srcs = [
        "entrypoint.cc",
    ],
    visibility = ["//visibility:public"],
    deps = [
        "//dataplane/standalone/packetio",
        "//dataplane/proto:sai_cc_grpc",
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

# By building a non-shared library, detect some linker issues.
cc_binary(
    name = "saibin",
    linkstatic = True,
    deps = [
        ":entrypoint",
    ],
)
