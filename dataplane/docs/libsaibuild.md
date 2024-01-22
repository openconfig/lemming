# Building libsai.so

## Option 1: Bazel

Bazel is the preferred way to build Lemming. Bazel produces a mostly statically linked library. This can cause issues where some
libraries are linked twice (statically and dynamically), in which case alternative build system can be used.

1. `bazel build //dataplane/standalone/sai`
2. Shared library is located at `bazel-bin/dataplane/standalone/libsai.so`

See details [here](https://bazel.build/reference/be/c-cpp#cc_binary).

## Option 2: Make (and Docker)

The make rule creates a dynamically linked shared library using the system installed libraries for glog, grpc, and protobuf.

1. Option a: Build locally `make libsai.so`.
2. Option b: Build using docker container
   1. Build library in debian:bullseye container: `docker build -f Dockerfile.saibuilder . -t lemming-libsai`
   2. Create instance of built container `docker create --name lemming-temp lemming-libsai:latest`
   3. Copy deb file to local machine `docker cp lemming-temp:/build/pkg/lucius-libsai.deb .`
   4. Remove temporary container `docker rm -f lemming-temp`
