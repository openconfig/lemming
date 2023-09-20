cc_library(
    name = "sai",
    srcs = glob(["*.cc"]),
    hdrs = glob(
        ["*.h"],
        exclude = ["entry.h"],
    ),
    visibility = ["//visibility:public"],
    deps = [
        "//dataplane/standalone/proto:sai_cc_grpc",
        "//dataplane/standalone/proto:sai_cc_proto",
        "@com_github_google_glog//:glog",
        "@com_github_opencomputeproject_sai//:sai",
    ],
)
