load("@com_github_grpc_grpc//bazel:cc_grpc_library.bzl", "cc_grpc_library")
load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")
load("@rules_proto//proto:defs.bzl", "proto_library")
load("//:doxygen.bzl", "doxygen_xml")
load("//:supported_apis.bzl", "SUPPORTED_APIS")

SUPPORTED_SAI_APIS = SUPPORTED_APIS

ALL_SAI_APIS = SUPPORTED_APIS

GEN_CLIENT_SOURCES = [f + ".cc" for f in SUPPORTED_SAI_APIS] + ["enum.cc"]

GEN_CLIENT_HEADERS = [f + ".h" for f in SUPPORTED_SAI_APIS] + ["enum.h"]

GEN_PROTOS = [f + ".proto" for f in ALL_SAI_APIS] + ["common.proto"]

GEN_OUTS = GEN_CLIENT_SOURCES + GEN_CLIENT_HEADERS + GEN_PROTOS

doxygen_xml(
    name = "sai_doxygen_xml",
    doxyfile_template = "@com_github_opencomputeproject_sai//:meta/Doxyfile",
    headers = [
        "@com_github_opencomputeproject_sai//:headers",
    ],
)

genrule(
    name = "generate_sai_code",
    srcs = [
        ":sai_doxygen_xml",
        "@com_github_opencomputeproject_sai//:headers",
    ],
    outs = GEN_OUTS,
    cmd = "xml_dir=$(location :sai_doxygen_xml); " +
          "sai_files=($(locations @com_github_opencomputeproject_sai//:headers)); " +
          "sai_dir=$$(dirname $$(dirname $${sai_files[0]})); " +
          "$(location //dataplane/apigen:apigen) " +
          "-sai_path=$$sai_dir " +
          "-client_out_dir=$(RULEDIR) " +
          "-server_out_dir=$$TMPDIR " +
          "-proto_out_dir=$(RULEDIR) " +
          "-proto_import_path=dataplane/proto/sai " +
          "-xml_path=$$xml_dir " +
          "-apis=" + ",".join(SUPPORTED_APIS),
    tools = ["//dataplane/apigen"],
    visibility = ["//visibility:public"],
)

filegroup(
    name = "sai_client_sources",
    srcs = [":" + f for f in GEN_CLIENT_SOURCES],
    visibility = ["//visibility:public"],
)

filegroup(
    name = "sai_client_headers",
    srcs = [":" + f for f in GEN_CLIENT_HEADERS],
    visibility = ["//visibility:public"],
)

proto_library(
    name = "sai_proto",
    # keep
    srcs = [":" + f for f in GEN_PROTOS],
    visibility = ["//visibility:public"],
    deps = [
        "@com_google_protobuf//:descriptor_proto",
        "@com_google_protobuf//:timestamp_proto",
    ],
)

# keep
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

# keep
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
