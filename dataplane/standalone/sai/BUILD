cc_library(
    name = "sai",
    srcs = [
        "common.cc",
    ] + ["//dataplane/proto/sai:sai_client_sources"],
    hdrs = [
        "common.h",
    ] + ["//dataplane/proto/sai:sai_client_headers"],
    includes = ["."],
    visibility = ["//visibility:public"],
    deps = [
        "//dataplane/proto/sai:sai_cc_grpc",
        "//dataplane/proto/sai:sai_cc_proto",
        "@com_github_google_glog//:glog",
        "@com_github_opencomputeproject_sai//:sai",
    ],
)
