cc_library(
    name = "sai",
    srcs = glob(["*.cc"]),
    hdrs = glob(
        ["*.h"],
        exclude = ["entry.h"],
    ),
    visibility = ["//visibility:public"],
    deps = [
        ":entry",
        "//dataplane/standalone:translator",
        "@com_github_google_glog//:glog",
        "@com_github_opencomputeproject_sai//:sai",
    ],
)

cc_library(
    name = "entry",
    hdrs = ["entry.h"],
    visibility = ["//visibility:public"],
    deps = [
        "@com_github_opencomputeproject_sai//:sai",
    ],
)
