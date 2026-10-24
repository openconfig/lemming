load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

# Hedron's Compile Commands Extractor for Bazel
# https://github.com/hedronvision/bazel-compile-commands-extractor
http_archive(
    name = "hedron_compile_commands",
    sha256 = "1b08abffbfbe89f6dbee6a5b33753792e8004f6a36f37c0f72115bec86e68724",
    strip_prefix = "bazel-compile-commands-extractor-abb61a688167623088f8768cc9264798df6a9d10",
    url = "https://github.com/hedronvision/bazel-compile-commands-extractor/archive/abb61a688167623088f8768cc9264798df6a9d10.tar.gz",
)

load("@hedron_compile_commands//:workspace_setup.bzl", "hedron_compile_commands_setup")

hedron_compile_commands_setup()

load("@hedron_compile_commands//:workspace_setup_transitive.bzl", "hedron_compile_commands_setup_transitive")

hedron_compile_commands_setup_transitive()

load("@hedron_compile_commands//:workspace_setup_transitive_transitive.bzl", "hedron_compile_commands_setup_transitive_transitive")

hedron_compile_commands_setup_transitive_transitive()

load("@hedron_compile_commands//:workspace_setup_transitive_transitive_transitive.bzl", "hedron_compile_commands_setup_transitive_transitive_transitive")

hedron_compile_commands_setup_transitive_transitive_transitive()

# External tools and libraries
http_archive(
    name = "com_github_grpc_grpc",
    strip_prefix = "grpc-1.70.0",
    urls = [
        "https://github.com/grpc/grpc/archive/refs/tags/v1.70.0.tar.gz",
    ],
)

http_archive(
    name = "com_google_protobuf",
    sha256 = "4717ac5e9b69bfed18b1cbc39b8fc87921e0518b68c4a84643d547720a7adbd5",
    strip_prefix = "protobuf-36.1",
    url = "https://github.com/protocolbuffers/protobuf/archive/refs/tags/v36.1.zip",
)

http_archive(
    name = "rules_cc",
    sha256 = "6bcf1b46d9e629cbe4c96d7bfc56f61034c2193b7b16001201282ea50bef5d91",
    strip_prefix = "rules_cc-0.2.27",
    urls = ["https://github.com/bazelbuild/rules_cc/archive/refs/tags/0.2.27.tar.gz"],
)

http_archive(
    name = "rules_python",
    sha256 = "05088402fa356af097064dbccd7082a061171ef39fbd748fab946e9775f46a16",
    strip_prefix = "rules_python-2.4.5",
    url = "https://github.com/bazel-contrib/rules_python/releases/download/2.4.5/rules_python-2.4.5.tar.gz",
)

http_archive(
    name = "rules_pkg",
    sha256 = "b7215c636f22c1849f1c3142c72f4b954bb12bb8dcf3cbe229ae6e69cc6479db",
    urls = [
        "https://github.com/bazelbuild/rules_pkg/releases/download/1.3/rules_pkg-1.1.0.tar.gz",
    ],
)

http_archive(
    name = "com_google_absl",
    sha256 = "f9148fb00ec98a2396bdf875c99a78e6a70afa662b107862d92b285d857a8320",
    strip_prefix = "abseil-cpp-20250814.2",
    urls = [
        "https://github.com/abseil/abseil-cpp/releases/download/20250814.2/abseil-cpp-20250814.2.tar.gz",
    ],
)

http_archive(
    name = "build_bazel_rules_swift",
    sha256 = "d2a96eed71684877051e88dd2c2a772edbdb754de488cee066dd72cce240cfc6",
    url = "https://github.com/bazelbuild/rules_swift/releases/download/4.2.2/rules_swift.4.2.2.tar.gz",
)

http_archive(
    name = "rules_proto_grpc",
    sha256 = "a25b9992ded30441b0250b2a3e69638e94e664bffcbd77daac6dd355fc92dd17",
    strip_prefix = "rules_proto_grpc-5.8.0",
    urls = ["https://github.com/rules-proto-grpc/rules_proto_grpc/archive/5.8.0.tar.gz"],
)

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "d4a6e93dc63944172bc3b915a11bfa8d7391c9f1fdc11d7e43724f497617f0e5",
    urls = [
        "https://mirror.bazel.build/github.com/bazel-contrib/rules_go/releases/download/v0.65.0/rules_go-v0.65.0.zip",
        "https://github.com/bazel-contrib/rules_go/releases/download/v0.65.0/rules_go-v0.65.0.zip",
    ],
)

http_archive(
    name = "bazel_gazelle",
    sha256 = "9dae87c03ea06b2990b8836256909a18367e191f885d5bb1d030e8ba80491b42",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.54.0/bazel-gazelle-v0.54.0.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.54.0/bazel-gazelle-v0.54.0.tar.gz",
    ],
)

http_archive(
    name = "rules_proto",
    sha256 = "14a225870ab4e91869652cfd69ef2028277fc1dc4910d65d353b62d6e0ae21f4",
    strip_prefix = "rules_proto-7.1.0",
    url = "https://github.com/bazelbuild/rules_proto/releases/download/7.1.0/rules_proto-7.1.0.tar.gz",
)

http_archive(
    name = "rules_oci",
    sha256 = "243eac87ec41d68d1d9e6d21884a50ec28e031f702b0b9fcdc7f97c3f3d5389a",
    strip_prefix = "rules_oci-2.3.3",
    url = "https://github.com/bazel-contrib/rules_oci/releases/download/v2.3.3/rules_oci-v2.3.3.tar.gz",
)

http_archive(
    name = "bazel_skylib",
    sha256 = "37cdfbc6faefea94f7b37760a305c98c08981116c2bc9e821e3b423221fad8c8",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-skylib/releases/download/1.9.2/bazel-skylib-1.9.2.tar.gz",
        "https://github.com/bazelbuild/bazel-skylib/releases/download/1.9.2/bazel-skylib-1.9.2.tar.gz",
    ],
)

http_archive(
    name = "googleapis",
    sha256 = "87f3ec7863f65322956ba85ebe3b8efb3a5774321be8a2dab8a3a180855c4ca0",
    strip_prefix = "googleapis-dff1cd0c7ff8c26ff4c9f7a87cbe175fb5124d4e",
    urls = [
        "https://github.com/googleapis/googleapis/archive/dff1cd0c7ff8c26ff4c9f7a87cbe175fb5124d4e.zip",
    ],
)

http_archive(
    name = "rules_distroless",
    sha256 = "1d74bebf851385a92858337acd98f612bea446588d6b1ea7a46400bb13c9a50c",
    strip_prefix = "rules_distroless-0.11.0",
    url = "https://github.com/GoogleContainerTools/rules_distroless/releases/download/v0.11.0/rules_distroless-v0.11.0.tar.gz",
)

http_archive(
    name = "openconfig_gnmi",
    integrity = "sha256-gT+KUt+gbdG5osd1smxC02oFWV36b7CoXbrq1GtcQ6M=",
    strip_prefix = "gnmi-0.14.1",
    url = "https://github.com/openconfig/gnmi/archive/refs/tags/v0.14.1.tar.gz",
)

# The non-polyfill version of this is needed by rules_proto below.
http_archive(
    name = "bazel_features",
    sha256 = "99232c34de353633c384188afa372d86a9dedbd19c4b46f99ac88a1ff11bb299",
    strip_prefix = "bazel_features-1.53.0",
    url = "https://github.com/bazel-contrib/bazel_features/releases/download/v1.53.0/bazel_features-v1.53.0.tar.gz",
)

load("@bazel_features//:deps.bzl", "bazel_features_deps")

bazel_features_deps()

load("@bazel_skylib//:workspace.bzl", "bazel_skylib_workspace")

bazel_skylib_workspace()

# Go

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")
load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains(version = "1.25.0")

# Create the host platform repository transitively required by rules_go.
load("@bazel_tools//tools/build_defs/repo:utils.bzl", "maybe")
load("@platforms//host:extension.bzl", "host_platform_repo")

maybe(
    host_platform_repo,
    name = "host_platform",
)

# go_repositories

load("//:repositories.bzl", "go_dependencies", "go_repositories")

# gazelle:repository_macro repositories.bzl%go_dependencies
go_dependencies()

# gazelle:repository_macro repositories.bzl%go_repositories
go_repositories()

gazelle_dependencies()

# Protobuf and gRPC

load("@com_github_grpc_grpc//bazel:grpc_deps.bzl", "grpc_deps")

grpc_deps()

load("@rules_proto_grpc//:repositories.bzl", "rules_proto_grpc_repos", "rules_proto_grpc_toolchains")

rules_proto_grpc_toolchains()

rules_proto_grpc_repos()

load("@rules_proto//proto:repositories.bzl", "rules_proto_dependencies")

rules_proto_dependencies()

load("@rules_proto//proto:toolchains.bzl", "rules_proto_toolchains")

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

# Distroless

load("@rules_distroless//distroless:dependencies.bzl", "distroless_dependencies")

distroless_dependencies()

load("@rules_distroless//distroless:toolchains.bzl", "distroless_register_toolchains")

distroless_register_toolchains()

load("@rules_distroless//apt:index.bzl", "deb_index")

# bazel run @bookworm//:lock
deb_index(
    name = "bookworm",
    lock = "//:bookworm.lock.json",
    manifest = "//:bookworm.yaml",
)

load("@bookworm//:packages.bzl", "bookworm_packages")

bookworm_packages()

# OCI Container

load("@rules_oci//oci:dependencies.bzl", "rules_oci_dependencies")

rules_oci_dependencies()

load("@rules_oci//oci:repositories.bzl", "oci_register_toolchains")

# Crane was removed in rules_oci v2.x so digests from v1.x won't match v2.x.
oci_register_toolchains(name = "oci")

load("@rules_python//python:repositories.bzl", "py_repositories")

py_repositories()

load("@rules_pkg//:deps.bzl", "rules_pkg_dependencies")

rules_pkg_dependencies()

load(
    "@build_bazel_rules_swift//swift:repositories.bzl",
    "swift_rules_dependencies",
)

swift_rules_dependencies()

# External non-Go or bazel friendly dependencies

http_archive(
    name = "com_github_opencomputeproject_sai",
    build_file_content = """
load("@rules_cc//cc:defs.bzl", "cc_library")

cc_library(
    name = "sai",
    hdrs = glob(["inc/*.h","experimental/*.h"]),
    includes = ["inc", "experimental"],
    visibility = ["//visibility:public"],
)
""",
    patch_args = ["-p1"],
    patches = ["//patches:sai.patch"],
    sha256 = "8c6551c36f95fd0e6bc290b94be1f18ff048d24d3aaf996946ba8b3dde44904b",
    strip_prefix = "SAI-1.19.2",
    urls = ["https://github.com/opencomputeproject/SAI/archive/refs/tags/v1.19.2.tar.gz"],
)

http_archive(
    name = "com_github_gflags_gflags",
    sha256 = "1b5e0648d7b94021895086e65479b4eaca7935eecfffd7dd9512eb576181c53d",
    strip_prefix = "gflags-2.3.1",
    urls = ["https://github.com/gflags/gflags/archive/v2.3.1.tar.gz"],
)

http_archive(
    name = "com_github_google_glog",
    sha256 = "c17d85c03ad9630006ef32c7be7c65656aba2e7e2fbfc82226b7e680c771fc88",
    strip_prefix = "glog-0.7.1",
    urls = ["https://github.com/google/glog/archive/v0.7.1.zip"],
)

load("@rules_oci//oci:pull.bzl", "oci_pull")

oci_pull(
    name = "debian_bookworm",
    digest = "sha256:4c25275446ca371766b8b093c696782eabce06263fccacb8d5d813b3dc44b308",  # bookworm as of 06/20/24
    image = "debian",
    platforms = ["linux/amd64"],
)

oci_pull(
    name = "distroless_static_debug_nonroot",
    digest = "sha256:15f8637509792e05b97d83451b0b6ccfa4abd548f997e406ac1e72c87fe1fa78",  # debug-nonroot as of 06/20/24
    image = "gcr.io/distroless/static",
    platforms = ["linux/amd64"],
)
