load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "actions",
    srcs = [
        "continue.go",
        "debug.go",
        "decap.go",
        "doc.go",
        "drop.go",
        "encap.go",
        "evaluate.go",
        "flow_counter.go",
        "lookup.go",
        "mirror.go",
        "output.go",
        "ratelimit.go",
        "reparse.go",
        "select_action_list.go",
        "swapoutput.go",
        "transmit.go",
        "update.go",
    ],
    importpath = "github.com/openconfig/lemming/dataplane/forwarding/fwdaction/actions",
    visibility = ["//visibility:public"],
    deps = [
        "//dataplane/forwarding/attributes",
        "//dataplane/forwarding/fwdaction",
        "//dataplane/forwarding/fwdport",
        "//dataplane/forwarding/fwdtable",
        "//dataplane/forwarding/infra/fwdattribute",
        "//dataplane/forwarding/infra/fwdcontext",
        "//dataplane/forwarding/infra/fwdflowcounter",
        "//dataplane/forwarding/infra/fwdobject",
        "//dataplane/forwarding/infra/fwdpacket",
        "//dataplane/forwarding/util/hash/crc16",
        "//proto/forwarding",
        "@com_github_golang_glog//:glog",
    ],
)

go_test(
    name = "actions_test",
    size = "small",
    srcs = [
        "debug_test.go",
        "drop_test.go",
        "flowcounter_test.go",
        "lookup_test.go",
        "mirror_test.go",
        "ratelimit_test.go",
        "reparse_test.go",
        "select_action_list_test.go",
        "transmit_test.go",
        "update_test.go",
    ],
    embed = [":actions"],
    deps = [
        "//dataplane/forwarding/fwdaction",
        "//dataplane/forwarding/fwdaction/mock_fwdpacket",
        "//dataplane/forwarding/fwdaction/mock_fwdport",
        "//dataplane/forwarding/fwdaction/mock_fwdtable",
        "//dataplane/forwarding/fwdport",
        "//dataplane/forwarding/fwdtable",
        "//dataplane/forwarding/infra/fwdcontext",
        "//dataplane/forwarding/infra/fwdflowcounter",
        "//dataplane/forwarding/infra/fwdobject",
        "//dataplane/forwarding/infra/fwdpacket",
        "//dataplane/forwarding/protocol",
        "//dataplane/forwarding/protocol/arp",
        "//dataplane/forwarding/protocol/ethernet",
        "//dataplane/forwarding/protocol/metadata",
        "//dataplane/forwarding/protocol/opaque",
        "//proto/forwarding",
        "@com_github_go_logr_logr//testr",
        "@org_uber_go_mock//gomock",
    ],
)
