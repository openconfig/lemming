load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

# External tools and libraries

http_archive(
    name = "com_github_grpc_grpc",
    sha256 = "79e3ff93f7fa3c8433e2165f2550fa14889fce147c15d9828531cbfc7ad11e01",
    strip_prefix = "grpc-1.54.1",
    urls = [
        "https://github.com/grpc/grpc/archive/refs/tags/v1.54.1.tar.gz",
    ],
)

http_archive(
    name = "rules_proto_grpc",
    sha256 = "fb7fc7a3c19a92b2f15ed7c4ffb2983e956625c1436f57a3430b897ba9864059",
    strip_prefix = "rules_proto_grpc-4.3.0",
    urls = ["https://github.com/rules-proto-grpc/rules_proto_grpc/archive/4.3.0.tar.gz"],
)

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "91585017debb61982f7054c9688857a2ad1fd823fc3f9cb05048b0025c47d023",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.42.0/rules_go-v0.42.0.zip",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.42.0/rules_go-v0.42.0.zip",
    ],
)

http_archive(
    name = "bazel_gazelle",
    sha256 = "29218f8e0cebe583643cbf93cae6f971be8a2484cdcfa1e45057658df8d54002",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.32.0/bazel-gazelle-v0.32.0.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.32.0/bazel-gazelle-v0.32.0.tar.gz",
    ],
)

http_archive(
    name = "rules_oci",
    sha256 = "db57efd706f01eb3ce771468366baa1614b5b25f4cce99757e2b8d942155b8ec",
    strip_prefix = "rules_oci-1.0.0",
    url = "https://github.com/bazel-contrib/rules_oci/releases/download/v1.0.0/rules_oci-v1.0.0.tar.gz",
)

http_archive(
    name = "bazel_skylib",
    sha256 = "b8a1527901774180afc798aeb28c4634bdccf19c4d98e7bdd1ce79d1fe9aaad7",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-skylib/releases/download/1.4.1/bazel-skylib-1.4.1.tar.gz",
        "https://github.com/bazelbuild/bazel-skylib/releases/download/1.4.1/bazel-skylib-1.4.1.tar.gz",
    ],
)

http_archive(
    name = "googleapis",
    sha256 = "9d1a930e767c93c825398b8f8692eca3fe353b9aaadedfbcf1fca2282c85df88",
    strip_prefix = "googleapis-64926d52febbf298cb82a8f472ade4a3969ba922",
    urls = [
        "https://github.com/googleapis/googleapis/archive/64926d52febbf298cb82a8f472ade4a3969ba922.zip",
    ],
)

load("@bazel_skylib//:workspace.bzl", "bazel_skylib_workspace")

bazel_skylib_workspace()

# Go

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

go_rules_dependencies()

go_register_toolchains(version = "1.21.3")

# go_repositories

load("//:repositories.bzl", "go_repositories")

# gazelle:repository_macro repositories.bzl%go_repositories
go_repositories()

gazelle_dependencies()

# Protobuf and gRPC

load("@com_github_grpc_grpc//bazel:grpc_deps.bzl", "grpc_deps")

grpc_deps()

load("@rules_proto_grpc//:repositories.bzl", "rules_proto_grpc_repos", "rules_proto_grpc_toolchains")

rules_proto_grpc_toolchains()

rules_proto_grpc_repos()

load("@rules_proto//proto:repositories.bzl", "rules_proto_dependencies", "rules_proto_toolchains")

rules_proto_dependencies()

rules_proto_toolchains()

load("@com_google_protobuf//:protobuf_deps.bzl", "protobuf_deps")

protobuf_deps()

#load("@com_google_googleapis//:repository_rules.bzl", "switched_rules_by_language")
load("@googleapis//:repository_rules.bzl", "switched_rules_by_language")

switched_rules_by_language(
    name = "com_google_googleapis_imports",
    cc = True,
    grpc = True,
)

# OCI Container

load("@rules_oci//oci:dependencies.bzl", "rules_oci_dependencies")

rules_oci_dependencies()

load("@rules_oci//oci:repositories.bzl", "LATEST_CRANE_VERSION", "oci_register_toolchains")

oci_register_toolchains(
    name = "oci",
    crane_version = LATEST_CRANE_VERSION,
)

load("@rules_pkg//:deps.bzl", "rules_pkg_dependencies")

rules_pkg_dependencies()

# External non-Go or bazel friendly dependencies

http_archive(
    name = "com_github_p4lang_p4runtime",
    patch_args = ["-p1"],
    patches = ["//patches:p4.patch"],
    sha256 = "ba31fb9afce6e62ffe565b16bb909e144cd30d65d926cd90af25e99ee8de863a",
    strip_prefix = "p4runtime-1.4.0-rc.5/proto",
    urls = ["https://github.com/p4lang/p4runtime/archive/refs/tags/v1.4.0-rc.5.zip"],
)

http_archive(
    name = "com_github_opencomputeproject_sai",
    build_file_content = """
cc_library(
    name = "sai",
    hdrs = glob(["inc/*.h","experimental/*.h"]),
    includes = ["inc", "experimental"],
    visibility = ["//visibility:public"],
)
""",
    patch_args = ["-p1"],
    patches = ["//patches:sai.patch"],
    sha256 = "240d0211bbea2758faabfdbfa5e5488d837a47d42839bfe99b4bfbff52ab6c11",
    strip_prefix = "SAI-1.11.0",
    urls = ["https://github.com/opencomputeproject/SAI/archive/refs/tags/v1.11.0.tar.gz"],
)

http_archive(
    name = "com_github_gflags_gflags",
    sha256 = "34af2f15cf7367513b352bdcd2493ab14ce43692d2dcd9dfc499492966c64dcf",
    strip_prefix = "gflags-2.2.2",
    urls = ["https://github.com/gflags/gflags/archive/v2.2.2.tar.gz"],
)

http_archive(
    name = "com_github_google_glog",
    sha256 = "122fb6b712808ef43fbf80f75c52a21c9760683dae470154f02bddfc61135022",
    strip_prefix = "glog-0.6.0",
    urls = ["https://github.com/google/glog/archive/v0.6.0.zip"],
)

load("@rules_oci//oci:pull.bzl", "oci_pull")

oci_pull(
    name = "debian_bookworm",
    digest = "sha256:e7072ef5bbeaca98db3056a7d944d5dfb7a44d47770d10d54ee3f5a61144f049",  # bookworm as of 06/22/23
    image = "debian",
)

oci_pull(
    name = "distroless_static_debug_nonroot",
    digest = "sha256:312a533b1f5584141a7d212ddcc1d079259a84ef68a1a5b0f522017093e3afda",  # debug-nonroute as of 06/22/23
    image = "gcr.io/distroless/static",
    platforms = ["linux/amd64"],
)
