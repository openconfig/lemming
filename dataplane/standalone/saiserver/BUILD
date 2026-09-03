cc_library(
    name = "saiserver",
    srcs = glob(["*.cc"]),
    hdrs = glob(
        ["*.h"],
    ),
    visibility = ["//visibility:public"],
    deps = [
        "//dataplane/proto/sai:sai_cc_grpc",
        "//dataplane/proto/sai:sai_cc_proto",
        "@com_github_google_glog//:glog",
        "@com_github_opencomputeproject_sai//:sai",
    ],
)
