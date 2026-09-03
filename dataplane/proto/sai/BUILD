load("@com_github_grpc_grpc//bazel:cc_grpc_library.bzl", "cc_grpc_library")
load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")
load("@rules_proto//proto:defs.bzl", "proto_library")

proto_library(
    name = "sai_proto",
    srcs = [
        "acl.proto",
        "ars.proto",
        "ars_profile.proto",
        "bfd.proto",
        "bmtor.proto",
        "bridge.proto",
        "buffer.proto",
        "common.proto",
        "counter.proto",
        "dash_acl.proto",
        "dash_appliance.proto",
        "dash_direction_lookup.proto",
        "dash_eni.proto",
        "dash_flow.proto",
        "dash_ha.proto",
        "dash_inbound_routing.proto",
        "dash_meter.proto",
        "dash_outbound_ca_to_pa.proto",
        "dash_outbound_routing.proto",
        "dash_pa_validation.proto",
        "dash_tunnel.proto",
        "dash_vip.proto",
        "dash_vnet.proto",
        "debug_counter.proto",
        "dtel.proto",
        "fdb.proto",
        "generic_programmable.proto",
        "hash.proto",
        "hostif.proto",
        "icmp_echo.proto",
        "ipmc.proto",
        "ipmc_group.proto",
        "ipsec.proto",
        "isolation_group.proto",
        "l2mc.proto",
        "l2mc_group.proto",
        "lag.proto",
        "macsec.proto",
        "mcast_fdb.proto",
        "mirror.proto",
        "mpls.proto",
        "my_mac.proto",
        "nat.proto",
        "neighbor.proto",
        "next_hop.proto",
        "next_hop_group.proto",
        "poe.proto",
        "policer.proto",
        "port.proto",
        "qos_map.proto",
        "queue.proto",
        "route.proto",
        "router_interface.proto",
        "rpf_group.proto",
        "samplepacket.proto",
        "scheduler.proto",
        "scheduler_group.proto",
        "srv6.proto",
        "stp.proto",
        "switch.proto",
        "system_port.proto",
        "tam.proto",
        "tunnel.proto",
        "twamp.proto",
        "udf.proto",
        "virtual_router.proto",
        "vlan.proto",
        "wred.proto",
    ],
    visibility = ["//visibility:public"],
    deps = [
        "@com_google_protobuf//:descriptor_proto",
        "@com_google_protobuf//:timestamp_proto",
    ],
)

go_proto_library(
    name = "sai_go_proto",
    compilers = [
        "@io_bazel_rules_go//proto:go_grpc_v2",
        "@io_bazel_rules_go//proto:go_proto",
    ],
    importpath = "github.com/openconfig/lemming/dataplane/proto/sai",
    proto = ":sai_proto",
    visibility = ["//visibility:public"],
)

go_library(
    name = "sai",
    embed = [":sai_go_proto"],
    importpath = "github.com/openconfig/lemming/dataplane/proto/sai",
    visibility = ["//visibility:public"],
)

cc_proto_library(
    name = "sai_cc_proto",
    visibility = ["//visibility:public"],
    deps = [":sai_proto"],
)

cc_grpc_library(
    name = "sai_cc_grpc",
    srcs = [":sai_proto"],
    grpc_only = True,
    visibility = ["//visibility:public"],
    deps = [":sai_cc_proto"],
)
