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
    sha256 = "85803e01f347141e16a2f770213a496f808fff9f0138c7c0e0c9dfa708b0da92",
    strip_prefix = "protobuf-29.3",
    url = "https://github.com/protocolbuffers/protobuf/archive/refs/tags/v29.3.zip",
)

http_archive(
    name = "rules_cc",
    sha256 = "a2fdfde2ab9b2176bd6a33afca14458039023edb1dd2e73e6823810809df4027",
    strip_prefix = "rules_cc-0.2.14",
    urls = ["https://github.com/bazelbuild/rules_cc/archive/refs/tags/0.2.14.tar.gz"],
)

http_archive(
    name = "rules_python",
    sha256 = "f609f341d6e9090b981b3f45324d05a819fd7a5a56434f849c761971ce2c47da",
    strip_prefix = "rules_python-1.7.0",
    url = "https://github.com/bazel-contrib/rules_python/releases/download/1.7.0/rules_python-1.7.0.tar.gz",
)

http_archive(
    name = "rules_pkg",
    sha256 = "b5c9184a23bb0bcff241981fd9d9e2a97638a1374c9953bb1808836ce711f990",
    urls = [
        "https://github.com/bazelbuild/rules_pkg/releases/download/1.2.0/rules_pkg-1.2.0.tar.gz",
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
    sha256 = "f7a67197cd8a79debfe70b8cef4dc19d03039af02cc561e31e0718e98cad83ac",
    url = "https://github.com/bazelbuild/rules_swift/releases/download/2.9.0/rules_swift.2.9.0.tar.gz",
)

http_archive(
    name = "rules_proto_grpc",
    sha256 = "c0d718f4d892c524025504e67a5bfe83360b3a982e654bc71fed7514eb8ac8ad",
    strip_prefix = "rules_proto_grpc-4.6.0",
    urls = ["https://github.com/rules-proto-grpc/rules_proto_grpc/archive/4.6.0.tar.gz"],
)

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "68af54cb97fbdee5e5e8fe8d210d15a518f9d62abfd71620c3eaff3b26a5ff86",
    urls = [
        "https://mirror.bazel.build/github.com/bazel-contrib/rules_go/releases/download/v0.59.0/rules_go-v0.59.0.zip",
        "https://github.com/bazel-contrib/rules_go/releases/download/v0.59.0/rules_go-v0.59.0.zip",
    ],
)

http_archive(
    name = "bazel_gazelle",
    sha256 = "675114d8b433d0a9f54d81171833be96ebc4113115664b791e6f204d58e93446",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.47.0/bazel-gazelle-v0.47.0.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.47.0/bazel-gazelle-v0.47.0.tar.gz",
    ],
)

http_archive(
    name = "rules_proto",
    sha256 = "6fb6767d1bef535310547e03247f7518b03487740c11b6c6adb7952033fe1295",
    strip_prefix = "rules_proto-6.0.2",
    url = "https://github.com/bazelbuild/rules_proto/releases/download/6.0.2/rules_proto-6.0.2.tar.gz",
)

http_archive(
    name = "rules_oci",
    sha256 = "b8db7ab889d501db33313620b2c8040dbb07e95c26a0fefe06004b35baf80e08",
    strip_prefix = "rules_oci-2.2.7",
    url = "https://github.com/bazel-contrib/rules_oci/releases/download/v2.2.7/rules_oci-v2.2.7.tar.gz",
)

http_archive(
    name = "bazel_skylib",
    sha256 = "6e78f0e57de26801f6f564fa7c4a48dc8b36873e416257a92bbb0937eeac8446",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-skylib/releases/download/1.8.2/bazel-skylib-1.8.2.tar.gz",
        "https://github.com/bazelbuild/bazel-skylib/releases/download/1.8.2/bazel-skylib-1.8.2.tar.gz",
    ],
)

http_archive(
    name = "googleapis",
    sha256 = "7d094c97e6fcbc8053f3d7b9c858a3fe23d2c8736bb48d140e6f7b39330ee428",
    strip_prefix = "googleapis-b81889fefc4a617605e3d779c10a2dc6092db671",
    urls = [
        "https://github.com/googleapis/googleapis/archive/b81889fefc4a617605e3d779c10a2dc6092db671.zip",
    ],
)

http_archive(
    name = "rules_distroless",
    sha256 = "959ea166d5161834292bd7a43eefd0385578b69ab85641cba0e279afa747f933",
    strip_prefix = "rules_distroless-0.6.1",
    url = "https://github.com/GoogleContainerTools/rules_distroless/releases/download/v0.6.1/rules_distroless-v0.6.1.tar.gz",
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
    sha256 = "07271d0f6b12633777b69020c4cb1eb67b1939c0cf84bb3944dc85cc250c0c01",
    strip_prefix = "bazel_features-1.38.0",
    url = "https://github.com/bazel-contrib/bazel_features/releases/download/v1.38.0/bazel_features-v1.38.0.tar.gz",
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
    sha256 = "05411b13b32abcc50f2f2b78e491e503b2b05e5a1503699abd4cc1b81f90d1ae",
    strip_prefix = "SAI-1.17.1",
    urls = ["https://github.com/opencomputeproject/SAI/archive/refs/tags/v1.17.1.tar.gz"],
)

http_archive(
    name = "com_github_gflags_gflags",
    sha256 = "f619a51371f41c0ad6837b2a98af9d4643b3371015d873887f7e8d3237320b2f",
    strip_prefix = "gflags-2.3.0",
    urls = ["https://github.com/gflags/gflags/archive/v2.3.0.tar.gz"],
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
    digest = "sha256:0d01188e8dd0ac63bf155900fad49279131a876a1ea7fac917c62e87ccb2732d",  # bookworm as of 06/20/24
    image = "debian",
    platforms = ["linux/amd64"],
)

oci_pull(
    name = "distroless_static_debug_nonroot",
    digest = "sha256:4b2a093ef4649bccd586625090a3c668b254cfe180dee54f4c94f3e9bd7e381e",  # debug-nonroot as of 06/20/24
    image = "gcr.io/distroless/static",
    platforms = ["linux/amd64"],
)
