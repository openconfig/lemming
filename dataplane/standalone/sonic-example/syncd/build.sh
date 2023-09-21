set -xe

DIR=$(dirname "$0")

rm -f libsai.so
bazel build //dataplane/standalone:sai --config=docker-sandbox
cp $DIR/../../../../bazel-bin/dataplane/standalone/libsai.so .
DOCKER_BUILDKIT=1 docker build  . -f $DIR/Dockerfile.syncdlucius -t "docker-syncd:latest"
docker save docker-syncd:latest -o  $DIR/docker-syncd.tar.gz