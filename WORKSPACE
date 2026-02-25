load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

# Hedron's Compile Commands Extractor for Bazel
# https://github.com/hedronvision/bazel-compile-commands-extractor
http_archive(
    name = "hedron_compile_commands",
    sha256 = "658122cfb1f25be76ea212b00f5eb047d8e2adc8bcf923b918461f2b1e37cdf2",
    strip_prefix = "bazel-compile-commands-extractor-4f28899228fb3ad0126897876f147ca15026151e",
    url = "https://github.com/hedronvision/bazel-compile-commands-extractor/archive/4f28899228fb3ad0126897876f147ca15026151e.tar.gz",
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
    sha256 = "c5dddfcb1702737d97c8ae1becd33079d0b28a32688727d177be46915924f78e",
    strip_prefix = "protobuf-34.0",
    url = "https://github.com/protocolbuffers/protobuf/archive/refs/tags/v34.0.zip",
)

http_archive(
    name = "rules_cc",
    sha256 = "712d77868b3152dd618c4d64faaddefcc5965f90f5de6e6dd1d5ddcd0be82d42",
    strip_prefix = "rules_cc-0.1.1",
    urls = ["https://github.com/bazelbuild/rules_cc/archive/refs/tags/0.1.1.tar.gz"],
)

http_archive(
    name = "rules_python",
    sha256 = "2cc26bbd53854ceb76dd42a834b1002cd4ba7f8df35440cf03482e045affc244",
    strip_prefix = "rules_python-1.3.0",
    url = "https://github.com/bazel-contrib/rules_python/releases/download/1.3.0/rules_python-1.3.0.tar.gz",
)

http_archive(
    name = "rules_pkg",
    sha256 = "b7215c636f22c1849f1c3142c72f4b954bb12bb8dcf3cbe229ae6e69cc6479db",
    urls = [
        "https://github.com/bazelbuild/rules_pkg/releases/download/1.1.0/rules_pkg-1.1.0.tar.gz",
    ],
)

http_archive(
    name = "com_google_absl",
    sha256 = "b396401fd29e2e679cace77867481d388c807671dc2acc602a0259eeb79b7811",
    strip_prefix = "abseil-cpp-20250127.1",
    urls = [
        "https://github.com/abseil/abseil-cpp/releases/download/20250127.1/abseil-cpp-20250127.1.tar.gz",
    ],
)

http_archive(
    name = "build_bazel_rules_swift",
    sha256 = "4901feadef8e47ede930c95c40298dd38a83a81eb1ed5b74e62abfa546ff2d1b",
    url = "https://github.com/bazelbuild/rules_swift/releases/download/2.8.1/rules_swift.2.8.1.tar.gz",
)

http_archive(
    name = "rules_proto_grpc",
    sha256 = "fb7fc7a3c19a92b2f15ed7c4ffb2983e956625c1436f57a3430b897ba9864059",
    strip_prefix = "rules_proto_grpc-4.3.0",
    urls = ["https://github.com/rules-proto-grpc/rules_proto_grpc/archive/4.3.0.tar.gz"],
)

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "a729c8ed2447c90fe140077689079ca0acfb7580ec41637f312d650ce9d93d96",
    urls = [
        "https://mirror.bazel.build/github.com/bazel-contrib/rules_go/releases/download/v0.57.0/rules_go-v0.57.0.zip",
        "https://github.com/bazel-contrib/rules_go/releases/download/v0.57.0/rules_go-v0.57.0.zip",
    ],
)

http_archive(
    name = "bazel_gazelle",
    sha256 = "b760f7fe75173886007f7c2e616a21241208f3d90e8657dc65d36a771e916b6a",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.39.1/bazel-gazelle-v0.39.1.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.39.1/bazel-gazelle-v0.39.1.tar.gz",
    ],
)

http_archive(
    name = "rules_proto",
    sha256 = "303e86e722a520f6f326a50b41cfc16b98fe6d1955ce46642a5b7a67c11c0f5d",
    strip_prefix = "rules_proto-6.0.0",
    url = "https://github.com/bazelbuild/rules_proto/releases/download/6.0.0/rules_proto-6.0.0.tar.gz",
)

http_archive(
    name = "rules_oci",
    sha256 = "361c417e8c95cd7c3d8b5cf4b202e76bac8d41532131534ff8e6fa43aa161142",
    strip_prefix = "rules_oci-2.2.5",
    url = "https://github.com/bazel-contrib/rules_oci/releases/download/v2.2.5/rules_oci-v2.2.5.tar.gz",
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

http_archive(
    name = "rules_distroless",
    sha256 = "44c1e485723ad342212b48e410bae50306b5f8b39da65243e1db2f5b74faa8d6",
    strip_prefix = "rules_distroless-0.3.7",
    url = "https://github.com/GoogleContainerTools/rules_distroless/releases/download/v0.3.7/rules_distroless-v0.3.7.tar.gz",
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
    sha256 = "d7787da289a7fb497352211ad200ec9f698822a9e0757a4976fd9f713ff372b3",
    strip_prefix = "bazel_features-1.9.1",
    url = "https://github.com/bazel-contrib/bazel_features/releases/download/v1.9.1/bazel_features-v1.9.1.tar.gz",
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
cc_library(
    name = "sai",
    hdrs = glob(["inc/*.h","experimental/*.h"]),
    includes = ["inc", "experimental"],
    visibility = ["//visibility:public"],
)
""",
    patch_args = ["-p1"],
    patches = ["//patches:sai.patch"],
    sha256 = "4e3a1d010bda0c589db46e077725a2cd9624a5cc255c89d1caa79deb408d1fa7",
    strip_prefix = "SAI-1.14.0",
    urls = ["https://github.com/opencomputeproject/SAI/archive/refs/tags/v1.14.0.tar.gz"],
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
    digest = "sha256:a92ed51e0996d8e9de041ca05ce623d2c491444df6a535a566dabd5cb8336946",  # bookworm as of 06/20/24
    image = "debian",
    platforms = ["linux/amd64"],
)

oci_pull(
    name = "distroless_static_debug_nonroot",
    digest = "sha256:cb0459bf13af06cb3d3ee5dde5f1c5c34381cbce3a86bd08e1e7fd7a3ed28e59",  # debug-nonroot as of 06/20/24
    image = "gcr.io/distroless/static",
    platforms = ["linux/amd64"],
)
