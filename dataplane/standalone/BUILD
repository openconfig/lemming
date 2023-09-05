cc_library(
    name = "translator",
    srcs = [
        "acl.cc",
        "bridge.cc",
        "buffer.cc",
        "common.cc",
        "dtel.cc",
        "hostif.cc",
        "neighbor.cc",
        "next_hop.cc",
        "next_hop_group.cc",
        "policer.cc",
        "port.cc",
        "route.cc",
        "router_interface.cc",
        "switch.cc",
        "translator.cc",
        "tunnel.cc",
        "vlan.cc",
    ],
    hdrs = [
        "acl.h",
        "bridge.h",
        "buffer.h",
        "common.h",
        "dtel.h",
        "hostif.h",
        "neighbor.h",
        "next_hop.h",
        "next_hop_group.h",
        "policer.h",
        "port.h",
        "route.h",
        "router_interface.h",
        "switch.h",
        "translator.h",
        "tunnel.h",
        "vlan.h",
    ],
    visibility = ["//visibility:public"],
    deps = [
        "//dataplane/standalone/lucius:lucius_clib",
        "//dataplane/standalone/sai:entry",
        "//proto/dataplane:dataplane_cc_grpc",
        "//proto/dataplane:dataplane_cc_proto",
        "//proto/forwarding:forwarding_cc_grpc",
        "//proto/forwarding:forwarding_cc_proto",
        "@com_github_google_glog//:glog",
        "@com_github_grpc_grpc//:grpc++",
        "@com_github_opencomputeproject_sai//:meta",
        "@com_github_opencomputeproject_sai//:sai",
    ],
)

cc_library(
    name = "entrypoint",
    srcs = [
        "entrypoint.cc",
    ],
    visibility = ["//visibility:public"],
    deps = [
        ":translator",
        "//dataplane/standalone/proto:sai_cc_grpc",
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
